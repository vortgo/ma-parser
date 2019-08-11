package ParseScripts

import (
	"github.com/vortgo/ma-parser/models"
	"github.com/vortgo/ma-parser/repositories"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

const upcomingAlbumsUrl = "https://www.metal-archives.com/release/ajax-upcoming/json/1?sEcho=1&iColumns=5&sColumns=&iDisplayStart=0&iDisplayLength=100&mDataProp_0=0&mDataProp_1=1&mDataProp_2=2&mDataProp_3=3&mDataProp_4=4&iSortCol_0=4&sSortDir_0=asc&iSortingCols=1&bSortable_0=true&bSortable_1=true&bSortable_2=true&bSortable_3=true&bSortable_4=true&_=1565289584123"

func ParseUpcomingAlbums() {
	var wg sync.WaitGroup
	ticker := time.NewTicker(time.Minute * 30)
	wg.Add(1)
	go func() {
		for range ticker.C {
			latestBandUpdRepo := repositories.MakeLatestBandUpdateRepository()
			albumRepo := repositories.MakeAlbumRepository()
			upcomingAlbumRepo := repositories.MakeUpcomingAlbumRepository()
			jsonString := getJsonFromUrl(upcomingAlbumsUrl)
			albumList := parseJson(jsonString)
			list := albumList.Data[:10]

			for _, v := range list {
				r, _ := regexp.Compile(`<a href="(.*?)">`)
				link := r.FindStringSubmatch(v[0])[1]

				band := ParseBandByUrl(link)

				if band != nil {
					latestBandUpdate := models.LatestBandUpdate{BandID: band.ID}
					latestBandUpdRepo.Save(&latestBandUpdate)

					albumLink := r.FindStringSubmatch(v[1])[1]
					urlParts := strings.Split(albumLink, "/")
					albumPlatformId, _ := strconv.Atoi(urlParts[len(urlParts)-1])
					if album := albumRepo.FindAlbumByPlatformId(albumPlatformId); album != nil {
						upcomingAlbumRepo.Save(&models.UpcomingAlbum{Album: *album})
					}
				}
			}
		}
	}()

	wg.Wait()
}
