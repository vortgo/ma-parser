package ParseScripts

import (
	"encoding/json"
	"github.com/vortgo/ma-parser/logger"
	"github.com/vortgo/ma-parser/models"
	"github.com/vortgo/ma-parser/utils/tor"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

type bandList struct {
	Data  [][]string `json:"aaData"`
	Total int        `json:"iTotalRecords"`
}

func ParseBandList() {
	var log = logger.New()
	offset := 0
	jobs := make(chan string, 100)

	countGorutines, _ := strconv.Atoi(os.Getenv("COUNT_LIST_BAND_GORUTINES"))
	for w := 0; w < countGorutines; w++ {
		go parseBandWorker(jobs)
	}
	for {
		log.Printf("ParseBandList.go ParseBandList start current offset %d", offset)
		bandLinks := getBandsLinks(offset)

		if len(*bandLinks) <= 0 {
			offset = 0
			continue
		}

		for _, bandLink := range *bandLinks {
			log.Info("Send link to parse chanel")
			jobs <- bandLink.Url
		}

		log.Println("ParseBandList.go ParseBandList +200")
		offset += 200
	}
}

func getBandsLinks(offset int) *[]models.BandLink {
	url := "https://www.metal-archives.com/search/ajax-advanced/searching/bands/"

	var bandsLinks []models.BandLink

	link := url + "?iDisplayStart=" + strconv.Itoa(offset)
	jsonString := getJsonFromUrl(link)
	bandList := parseJson(jsonString)

	extractLinksFromBandList(bandList, &bandsLinks)

	log.Println("ParseBandList.go getBandsLinks +200")
	return &bandsLinks
}

func extractLinksFromBandList(bandList bandList, bandUrls *[]models.BandLink) {
	for _, v := range bandList.Data {
		r, _ := regexp.Compile(`<a href="(.*?)">`)
		link := r.FindStringSubmatch(v[0])[1]

		r, _ = regexp.Compile(`">(.*?)</a>`)
		name := r.FindStringSubmatch(v[0])[1]

		entity := models.BandLink{}
		entity.Name = name
		entity.Url = link
		*bandUrls = append(*bandUrls, entity)
	}
}

func getJsonFromUrl(url string) string {
	var log = logger.New()
	requester := tor.NewClient()
	response := requester.MakeGetRequest(url)
	defer response.Body.Close()

	body, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		log.SetData(logger.Data{
			"url": url,
		}).Error(readErr)
	}

	log.Println("ParseBandList.go getJsonFromUrl")
	return strings.Replace(string(body), "\"sEcho\": ,\n", "", -1)
}

func parseJson(jsonData string) bandList {
	var log = logger.New()
	bandList := bandList{}

	body := []byte(jsonData)

	jsonErr := json.Unmarshal(body, &bandList)
	if jsonErr != nil {
		log.SetData(logger.Data{
			"json_string": jsonData,
		}).Error(jsonErr)
	}

	log.Println("ParseBandList.go parseJson")
	return bandList
}

func parseBandWorker(jobs <-chan string) {
	var log = logger.New()
	defer func() {
		if e := recover(); e != nil {
			log.SetData(logger.Data{
				"stacktrace": string(debug.Stack()),
			}).Error(e)
		}

		log.Println("ParseBandList.go parseBandWorker defer")
	}()

	for url := range jobs {
		log.Println("ParseBandList.go parseBandWorker start ParseBandByUrl")
		ParseBandByUrl(url)
		time.Sleep(time.Second)
		log.Println("ParseBandList.go parseBandWorker end ParseBandByUrl")
	}
}
