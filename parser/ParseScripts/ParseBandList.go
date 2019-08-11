package ParseScripts

import (
	"encoding/json"
	"github.com/vortgo/ma-parser/logger"
	"github.com/vortgo/ma-parser/models"
	"github.com/vortgo/ma-parser/utils/tor"
	"io/ioutil"
	"os"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
)

type bandList struct {
	Data  [][]string `json:"aaData"`
	Total int        `json:"iTotalRecords"`
}

func ParseBandList() {
	offset := 0
	for {
		var wg sync.WaitGroup
		jobs := make(chan string, 100)
		bandLinks := getBandsLinks(offset)

		if len(*bandLinks) <= 0 {
			offset = 0
			continue
		}

		countGorutines, _ := strconv.Atoi(os.Getenv("COUNT_LIST_BAND_GORUTINES"))
		for w := 0; w < countGorutines; w++ {
			wg.Add(1)
			go parseBandWorker(jobs, &wg)
		}

		for _, bandLink := range *bandLinks {
			jobs <- bandLink.Url
		}

		close(jobs)

		wg.Wait()
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

	body, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		log.SetData(logger.Data{
			"url": url,
		}).Error(readErr)
	}

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
	return bandList
}

func parseBandWorker(jobs <-chan string, wg *sync.WaitGroup) {
	var log = logger.New()
	defer wg.Done()
	defer func() {
		if e := recover(); e != nil {
			log.SetData(logger.Data{
				"stacktrace": string(debug.Stack()),
			}).Error(e)
		}
	}()

	for url := range jobs {
		band := ParseBandByUrl(url)
		if band == nil {
			log.SetData(logger.Data{
				"url":        url,
				"stacktrace": string(debug.Stack()),
			}).Error("Received a nil after parse band")
			continue
		}
	}
}
