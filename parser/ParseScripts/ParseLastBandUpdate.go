package ParseScripts

import (
	"github.com/vortgo/ma-parser/logger"
	"github.com/vortgo/ma-parser/models"
	"github.com/vortgo/ma-parser/repositories"
	"os"
	"regexp"
	"runtime/debug"
	"strconv"
	"sync"
	"time"
)

const lastBandUpdateUrl = "https://www.metal-archives.com/archives/ajax-band-list/selection/2019-08/by/modified//json/1?sEcho=1iDisplayStart=0&iDisplayLength=200&mDataProp_3=3&mDataProp_4=4&mDataProp_5=5&iSortCol_0=4&sSortDir_0=desc&iSortingCols=1&bSortable_0=true&bSortable_1=true&bSortable_2=true&bSortable_3=true&bSortable_4=true&bSortable_5=true&_=1565032791855"

func ParseLastBandUpdate() {
	var log = logger.New()
	lastBandUpdatePeriod, _ := strconv.Atoi(os.Getenv("PARSE_LAST_BAND_UPDATE_PERIOD_MINUTES"))
	ticker := time.NewTicker(time.Minute * time.Duration(lastBandUpdatePeriod))
	defer func() {
		if e := recover(); e != nil {
			log.SetContext(logger.Context{
				Collection: "ParseLastBandUpdate",
			}).SetData(logger.Data{
				"stacktrace": string(debug.Stack()),
			}).Error(e)
			return
		}
	}()

	for range ticker.C {
		var wg sync.WaitGroup
		jsonString := getJsonFromUrl(lastBandUpdateUrl)
		bandList := parseJson(jsonString)
		latestBandUpdRepo := repositories.MakeLatestBandUpdateRepository()
		list := bandList.Data[:10]

		for _, v := range list {
			wg.Add(1)
			go func() {
				r, _ := regexp.Compile(`<a href="(.*?)">`)
				link := r.FindStringSubmatch(v[1])[1]
				band := ParseBandByUrl(link)

				if band != nil {
					latestBandUpdRepo.Save(&models.LatestBandUpdate{BandID: band.ID})
				}
			}()
		}
		wg.Wait()
	}

}
