package ParseScripts

import (
	"fmt"
	"github.com/vortgo/ma-parser/logger"
	"github.com/vortgo/ma-parser/repositories"
	"os"
	"regexp"
	"runtime/debug"
	"strconv"
	"time"
)

const lastBandUpdateUrl = "https://www.metal-archives.com/archives/ajax-band-list/selection/%s/by/modified//json/1?sEcho=1iDisplayStart=0&iDisplayLength=200&mDataProp_3=3&mDataProp_4=4&mDataProp_5=5&iSortCol_0=4&sSortDir_0=desc&iSortingCols=1&bSortable_0=true&bSortable_1=true&bSortable_2=true&bSortable_3=true&bSortable_4=true&bSortable_5=true&_=1565032791855"

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

	dt := time.Now()
	url := fmt.Sprintf(lastBandUpdateUrl, dt.Format("2006-01"))
	for range ticker.C {

		jsonString := getJsonFromUrl(url)
		bandList := parseJson(jsonString)
		latestBandUpdRepo := repositories.MakeLatestBandUpdateRepository()
		list := bandList.Data[:10]

		for _, v := range list {
			r, _ := regexp.Compile(`<a href="(.*?)">`)
			link := r.FindStringSubmatch(v[1])[1]
			band := ParseBandByUrl(link)

			if band != nil {
				latestBand := latestBandUpdRepo.FindByBandId(band.ID)
				latestBandUpdRepo.Save(latestBand)
			}
		}
	}
}
