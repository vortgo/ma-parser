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

func ParseAllReviews() {
	parseDate, err := time.Parse("2006-01", "2020-05")
	if err != nil {
		println(err)
	}

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

				parserReviewByLink(link)
			}
			offset += offsetStep
		}
		parseDate = parseDate.AddDate(0, 1, 0)
	}

}

func ParseLatestReviews() {
	currentDate := time.Now().Format("2006-01")

	link := fmt.Sprintf(reviewsListUrl, currentDate, 0)
	jsonString := getJsonFromUrl(link)
	reviewsList := parseJson(jsonString)

	for _, v := range reviewsList.Data {
		r, _ := regexp.Compile(`href="(.*?)"`)
		link := r.FindStringSubmatch(v[1])[1]

		parserReviewByLink(link)
	}
}

func parserReviewByLink(link string) {
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
	albumPlatformId := r.FindStringSubmatch(html)[1]

	var album *models.Album

	if platformId, err := strconv.Atoi(albumPlatformId); err == nil {
		albumRepo := repositories.MakeAlbumRepository()
		album = albumRepo.FindAlbumByPlatformId(platformId)
	}

	if album.ID == 0 {
		return
	}

	titleData := strings.Replace(doc.Find(".reviewBox .reviewTitle").Text(), "\n", "", 2)
	r, _ = regexp.Compile(`^([A-z0-9 ]+)- +([0-9]+)%$`)
	result = r.FindStringSubmatch(titleData)

	if len(result) < 3 {
		return
	}
	title := result[1]
	rating, _ := strconv.Atoi(result[2])

	text := doc.Find(".reviewContent").Text()
	author := doc.Find(".reviewBox .profileMenu").Text()
	date := strings.Replace(doc.Find(".reviewBox .profileMenu").Parent().Text(), author, "", 1)

	date = strings.Replace(date, "\n", "", -1)
	date = strings.TrimSpace(strings.Replace(date, ",", "", 1))
	layout := "Jan 1st, 2006"
	dateTime, err := time.Parse(layout, date)

	review.Text = text
	review.AlbumID = album.ID
	review.PlatformID = reviewPlatformId
	review.Rating = rating
	review.Title = title
	review.Date = dateTime
	review.Author = author

	reviewRepository.Save(&review)
}
