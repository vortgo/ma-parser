package ParseScripts

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/vortgo/ma-parser/logger"
	"github.com/vortgo/ma-parser/models"
	"github.com/vortgo/ma-parser/repositories"
	"github.com/vortgo/ma-parser/utils/tor"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"
)

func ParseSongs(album *models.Album, albumPage *goquery.Document) {
	var wg sync.WaitGroup
	songRepo := repositories.MakeSongRepository()
	albumPage.Find(`table.table_lyrics tbody tr`).Each(func(i int, selection *goquery.Selection) {
		wg.Add(1)
		go func() {
			defer func() {
				if e := recover(); e != nil {
					//skip
					//log.SetContext(logger.Context{
					//	AlbumId: int(album.ID),
					//	Collection: "Parse song",
					//}).SetData(logger.Data{
					//	"stacktrace": string(debug.Stack()),
					//}).Error(e)
				}
				wg.Done()
			}()

			if selection.HasClass("displayNone") == true {
				return
			}

			var Song *models.Song
			position := goquery.NewDocumentFromNode(selection.Find(`td`).Get(0)).Text()
			songsName := goquery.NewDocumentFromNode(selection.Find(`td`).Get(1)).Text()
			duration := goquery.NewDocumentFromNode(selection.Find(`td`).Get(2)).Text()

			if songId, exist := selection.Find(`.anchor`).Attr(`name`); exist {
				platformId, _ := strconv.Atoi(songId)

				if platformId == 0 {
					return
				}

				Song = songRepo.LoadByPlatformId(platformId)
				Song.Lyrics = getlyrics(songId)
				Song.Position, _ = strconv.Atoi(strings.TrimSpace(strings.Replace(position, `.`, ``, -1)))
				Song.Name = strings.TrimSpace(songsName)
				Song.Time = duration
				Song.PlatformID = platformId
				Song.BandID = album.BandID
				Song.AlbumID = album.ID

				songRepo.Save(Song)
			}
		}()
	})
	wg.Wait()
}

func getlyrics(platformId string) string {
	var log = logger.New()
	defer func() {
		if e := recover(); e != nil {
			log.SetContext(logger.Context{
				Collection: "Parse lyrics",
			}).SetData(logger.Data{
				"stacktrace": string(debug.Stack()),
			}).Error(e)
		}
	}()

	url := `https://www.metal-archives.com/release/ajax-view-lyrics/id/` + platformId

	requester := tor.NewClient()
	response := requester.MakeGetRequest(url)

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.SetContext(logger.Context{
			Collection: "Parse lyrics",
		}).SetData(logger.Data{
			"url": url,
		}).Error(err)
	}

	lyrics := doc.Text()

	time.Sleep(time.Second)
	return lyrics
}
