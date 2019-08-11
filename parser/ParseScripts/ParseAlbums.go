package ParseScripts

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/araddon/dateparse"
	"github.com/vortgo/ma-parser/logger"
	"github.com/vortgo/ma-parser/models"
	"github.com/vortgo/ma-parser/repositories"
	"github.com/vortgo/ma-parser/utils/tor"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
)

func ParseAlbumsByBand(band *models.Band) ([]*models.Album, error) {
	var albums []*models.Album
	var url = `https://www.metal-archives.com/band/discography/id/` + band.PlatformID + `/tab/all`
	var wg sync.WaitGroup
	requester := tor.NewClient()
	response := requester.MakeGetRequest(url)

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
	}

	doc.Find(`tbody tr`).Each(func(i int, tr *goquery.Selection) {
		wg.Add(1)
		go func() {
			defer func() {
				wg.Done()
			}()
			node := tr.Find(`td`).Get(0)
			tdDoc := goquery.NewDocumentFromNode(node)
			if url, exists := tdDoc.Find(`a`).Attr(`href`); exists {
				album := parseAlbumWithSongs(band, url)
				albums = append(albums, album)
			}
		}()
	})

	wg.Wait()

	return albums, nil
}

func parseAlbumWithSongs(band *models.Band, albumUrl string) *models.Album {
	var log = logger.New()
	defer func() {
		if e := recover(); e != nil {
			log.SetContext(logger.Context{
				BandId:     int(band.ID),
				Collection: "album with songs",
			}).SetData(logger.Data{
				"album_url":  albumUrl,
				"stacktrace": string(debug.Stack()),
			}).Error(e)
		}
	}()

	requester := tor.NewClient()
	response := requester.MakeGetRequest(albumUrl)

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.SetContext(logger.Context{
			BandId:     int(band.ID),
			Collection: "album with songs",
		}).SetData(logger.Data{
			"album_url":  albumUrl,
			"stacktrace": string(debug.Stack()),
		}).Error(err)
		return nil
	}

	var Album *models.Album

	chunks := strings.Split(albumUrl, `/`)
	last := len(chunks) - 1
	albumId := chunks[last]

	if platformId, err := strconv.Atoi(albumId); err == nil {
		albumRepo := repositories.MakeAlbumRepository()
		Album = albumRepo.FindAlbumByPlatformId(platformId)
	}

	Album.BandID = band.ID
	Album.Name = doc.Find(`.album_name`).Text()

	position := make(map[string]int)
	findLeftSection := map[string]string{"type": "Type:", "release_date": "Release date:"}

	doc.Find(`#album_info .float_left dt`).Each(func(i int, selection *goquery.Selection) {
		for k, v := range findLeftSection {
			if v == selection.Text() {
				position[k] = i
			}
		}
	})

	Node := doc.Find(`#album_info .float_left dd`).Get(position["type"])
	Album.Type = goquery.NewDocumentFromNode(Node).Text()

	Node = doc.Find(`#album_info .float_left dd`).Get(position["release_date"])
	releaseDate := goquery.NewDocumentFromNode(Node).Text()
	if t, err := dateparse.ParseAny(releaseDate); err == nil {
		Album.ReleaseDate = t
		Album.Year = t.Year()
	}

	Node = doc.Find(`#album_info .float_right dd`).Get(1)
	labelRepo := repositories.MakeLabelRepository()
	Album.Label = labelRepo.FindOrCreateLabelByName(goquery.NewDocumentFromNode(Node).Text())
	if imgHref, exist := doc.Find(`#cover`).Attr(`href`); exist {
		Album.Image = imgHref
	}

	Album.TotalTime = doc.Find(`#album_tabs_tracklist strong`).Text()

	albumRepo := repositories.MakeAlbumRepository()
	albumRepo.Save(Album)

	ParseSongs(Album, doc)

	return Album
}
