package ParseScripts

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/vortgo/ma-parser/models"
	"github.com/vortgo/ma-parser/repositories"
	"github.com/vortgo/ma-parser/utils"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const reviewsListUrl = "https://www.metal-archives.com/review/ajax-list-browse/by/date/selection/%s?sEcho=3&iColumns=7&sColumns=&iDisplayStart=%d&iDisplayLength=200&mDataProp_0=0&mDataProp_1=1&mDataProp_2=2&mDataProp_3=3&mDataProp_4=4&mDataProp_5=5&mDataProp_6=6&iSortCol_0=6&sSortDir_0=desc&iSortingCols=1&bSortable_0=true&bSortable_1=false&bSortable_2=true&bSortable_3=false&bSortable_4=true&bSortable_5=true&bSortable_6=true&_=1588494183245"
const offsetStep = 200

func ParseReviews() {
	//ticker := time.NewTicker(time.Hour * time.Duration(12))

	runParseReview()
	//for range ticker.C {
	//	runParseReview()
	//}
}

func runParseReview() {
	parseDate, _ := time.Parse("2006-01", "2002-07")
	for {
		if parseDate.After(time.Now()) {
			break
		}
		offset := 0
		for {
			log.Println(fmt.Sprintf("------ start  - %s", parseDate.Format("2006-01")))
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
	urlParts := strings.Split(link, "/")
	if len(urlParts) < 8 {
		return
	}
	reviewPlatformId := urlParts[6] + urlParts[8]

	reviewRepository := repositories.MakeReviewRepository()
	review := reviewRepository.FindReviewByPlatformId(reviewPlatformId)

	if review.ID != 0 {
		return
	}

	requester := utils.NewClient()
	response := requester.MakeGetRequest(link)
	defer response.Body.Close()

	doc, _ := goquery.NewDocumentFromReader(response.Body)
	html, _ := doc.Html()

	r, _ := regexp.Compile(`<h1 class="album_name"><a .*\/([0-9]+)">`)
	result := r.FindStringSubmatch(html)
	if len(result) < 2 {
		return
	}
	albumPlatformId := result[1]

	var album *models.Album

	if platformId, err := strconv.Atoi(albumPlatformId); err == nil {
		albumRepo := repositories.MakeAlbumRepository()
		album = albumRepo.FindAlbumByPlatformId(platformId)
	}
	if album.ID == 0 {
		return
	}

	titleData := strings.Replace(doc.Find(".reviewBox .reviewTitle").Text(), "\n", "", 2)
	r, _ = regexp.Compile(`^(.+)- +([0-9]+)%$`)
	result = r.FindStringSubmatch(titleData)

	if len(result) < 3 {
		return
	}
	title := result[1]
	rating, _ := strconv.Atoi(result[2])

	text := doc.Find(".reviewContent").Text()
	author := doc.Find(".reviewBox .profileMenu").Text()

	reviewDate, _ := time.Parse("2006-01", parseDate.Format("2006-01"))

	review.Text = text
	review.AlbumID = album.ID
	review.PlatformID = reviewPlatformId
	review.Rating = rating
	review.Title = title
	review.Date = reviewDate
	review.Author = author

	reviewRepository.Save(&review)
}
