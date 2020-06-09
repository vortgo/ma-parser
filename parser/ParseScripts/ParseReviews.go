package ParseScripts

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/vortgo/ma-parser/logger"
	"github.com/vortgo/ma-parser/models"
	"github.com/vortgo/ma-parser/repositories"
	"github.com/vortgo/ma-parser/utils"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

const reviewsListUrl = "https://www.metal-archives.com/review/ajax-list-browse/by/date/selection/%s?sEcho=3&iColumns=7&sColumns=&iDisplayStart=%d&iDisplayLength=200&mDataProp_0=0&mDataProp_1=1&mDataProp_2=2&mDataProp_3=3&mDataProp_4=4&mDataProp_5=5&mDataProp_6=6&iSortCol_0=6&sSortDir_0=desc&iSortingCols=1&bSortable_0=true&bSortable_1=false&bSortable_2=true&bSortable_3=false&bSortable_4=true&bSortable_5=true&bSortable_6=true&_=1588494183245"
const offsetStep = 200

func ParseReviews() {
	ticker := time.NewTicker(time.Hour * time.Duration(12))

	runParseReview()
	for range ticker.C {
		runParseReview()
	}
}

func runParseReview() {
	parseDate := time.Now()
	for {
		if parseDate.After(time.Now()) {
			break
		}
		offset := 0
		for {
			link := fmt.Sprintf(reviewsListUrl, parseDate.Format("2006-01"), offset)
			jsonString := getJsonFromUrl(link)
			reviewsList := parseJson(jsonString)

			if len(reviewsList.Data) == 0 {
				break
			}

			for _, v := range reviewsList.Data {
				r, _ := regexp.Compile(`href="(.*?)"`)
				link := r.FindStringSubmatch(v[1])[1]

				parserReviewByLink(link, parseDate)
			}
			offset += offsetStep
		}
		parseDate = parseDate.AddDate(0, 1, 0)
	}
}

func parserReviewByLink(link string, parseDate time.Time) {
	r, _ := regexp.Compile(`\/([0-9]+)$`)
	result := r.FindStringSubmatch(link)
	if len(result) < 2 {
		return
	}
	reviewPlatformId := result[1]

	reviewRepository := repositories.MakeReviewRepository()
	review := reviewRepository.FindReviewByPlatformId(reviewPlatformId)

	if review.ID != 0 {
		return
	}

	var log = logger.New()
	requester := utils.NewClient()
	response := requester.MakeGetRequest(link)
	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)
	html, _ := doc.Html()
	if err != nil {
		log.SetContext(logger.Context{
			Collection: "Parser review by url",
		}).SetData(logger.Data{
			"url":        link,
			"stacktrace": string(debug.Stack()),
		}).Error(err)
	}

	r, _ = regexp.Compile(`<h1 class="album_name"><a .*\/([0-9]+)">`)
	result = r.FindStringSubmatch(html)
	if len(result) < 2 {
		println(link)
		println("invalid album id")
		return
	}
	albumPlatformId := result[1]

	var album *models.Album

	if platformId, err := strconv.Atoi(albumPlatformId); err == nil {
		albumRepo := repositories.MakeAlbumRepository()
		album = albumRepo.FindAlbumByPlatformId(platformId)
	}
	if album.ID == 0 {
		println(link)
		println("no found album " + albumPlatformId)
		return
	}

	titleData := strings.Replace(doc.Find(".reviewBox .reviewTitle").Text(), "\n", "", 2)
	r, _ = regexp.Compile(`^(.+)- +([0-9]+)%$`)
	result = r.FindStringSubmatch(titleData)

	if len(result) < 3 {
		println(link)
		println("invalid title")
		return
	}
	title := result[1]
	rating, _ := strconv.Atoi(result[2])

	text := doc.Find(".reviewContent").Text()
	author := doc.Find(".reviewBox .profileMenu").Text()

	reviewDate, err := time.Parse("2006-01", parseDate.Format("2006-01"))

	review.Text = text
	review.AlbumID = album.ID
	review.PlatformID = reviewPlatformId
	review.Rating = rating
	review.Title = title
	review.Date = reviewDate
	review.Author = author

	reviewRepository.Save(&review)

	time.Sleep(time.Second * time.Duration(7))
}
