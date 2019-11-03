package ParseScripts

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/vortgo/ma-parser/logger"
	"github.com/vortgo/ma-parser/models"
	"github.com/vortgo/ma-parser/repositories"
	"github.com/vortgo/ma-parser/utils"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

func ParseBandByUrl(url string) *models.Band {
	var log = logger.New()
	defer func() {
		if e := recover(); e != nil {
			log.SetData(logger.Data{
				"url":        url,
				"stacktrace": string(debug.Stack()),
			}).Error(e)
		}
	}()

	requester := utils.NewClient()
	response := requester.MakeGetRequest(url)
	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.SetContext(logger.Context{
			Collection: "Parse band by url",
		}).SetData(logger.Data{
			"url":        url,
			"stacktrace": string(debug.Stack()),
		}).Error(err)
	}

	band := parseBandInfo(doc)
	//log.Print(doc.Html())
	time.Sleep(time.Second)
	return band
}

func parseBandInfo(doc *goquery.Document) *models.Band {
	var log = logger.New()
	var band models.Band
	defer func() {
		if e := recover(); e != nil {
			log.SetData(logger.Data{
				"stacktrace": string(debug.Stack()),
			}).Error(e)
		}
	}()

	bandRepo := repositories.MakeBandRepository()
	url, exists := doc.Find(`#band_info .band_name a`).Attr(`href`)
	if exists {
		chunks := strings.Split(url, `/`)
		last := len(chunks) - 1
		band = *bandRepo.FindBandByPlatformId(chunks[last])
	}

	position := make(map[string]int)
	findLeftSection := map[string]string{"country": "Country of origin:", "status": "Status:", "formed": "Formed in:"}
	findRightSection := map[string]string{"genres": "Genre:", "lyrics": "Lyrical themes:", "label": "Current label:"}

	doc.Find(`#band_stats .float_left dt`).Each(func(i int, selection *goquery.Selection) {
		for k, v := range findLeftSection {
			if v == selection.Text() {
				position[k] = i
			}
		}
	})

	doc.Find(`#band_stats .float_right dt`).Each(func(i int, selection *goquery.Selection) {
		for k, v := range findRightSection {
			if v == selection.Text() {
				position[k] = i
			}
		}
	})

	bandName := doc.Find(`#band_info .band_name a`).Text()
	band.Name = strings.Replace(bandName, "\n", "", -1)

	Node := doc.Find(`#band_stats .float_left dd`).Get(position["formed"])
	band.FormedIn, _ = strconv.Atoi(goquery.NewDocumentFromNode(Node).Text())

	Node = doc.Find(`#band_stats .float_left dd`).Get(position["status"])
	band.Status = goquery.NewDocumentFromNode(Node).Text()

	Node = doc.Find(`#band_stats .float_left dd`).Get(position["country"])
	countryRepo := repositories.MakeCountryRepository()
	band.Country = countryRepo.FindOrCreatCountryByName(goquery.NewDocumentFromNode(Node).Text())

	Node = doc.Find(`#band_stats .float_right dd`).Get(position["label"])
	labelRepo := repositories.MakeLabelRepository()
	band.Label = labelRepo.FindOrCreateLabelByName(goquery.NewDocumentFromNode(Node).Text())

	Node = doc.Find(`#band_stats .float_right dd`).Get(position[`genres`])
	genres := strings.Split(goquery.NewDocumentFromNode(Node).Text(), `/`)

	Node = doc.Find(`#band_stats .float_right dd`).Get(position[`lyrics`])
	lyrics := strings.Split(goquery.NewDocumentFromNode(Node).Text(), `,`)

	if bandLogo, exists := doc.Find(`#logo`).Attr(`href`); exists {
		band.ImageLogo = bandLogo
	}

	if bandLogo, exists := doc.Find(`#photo`).Attr(`href`); exists {
		band.ImageBand = bandLogo
	}

	var genresObjs []*models.Genre

	for _, v := range genres {
		genreRepo := repositories.MakeGenreRepository()
		genreObj := genreRepo.FindOrCreatGenreByName(strings.Trim(v, ` `))
		genresObjs = append(genresObjs, genreObj)
	}

	var lyricsObjs []*models.LyricalTheme

	for _, name := range lyrics {
		lyricRepo := repositories.MakeLyricalThemeRepository()
		lyricObj := lyricRepo.FindOrCreatLyricalThemeByName(strings.Trim(name, ` `))
		lyricsObjs = append(lyricsObjs, lyricObj)
	}

	band.Genres = genresObjs
	band.LyricalThemes = lyricsObjs
	band.Description = parseDescription(band.PlatformID)

	bandRepo.Save(&band)

	ParseAlbumsByBand(&band)

	return &band
}

func parseDescription(platformId string) string {
	var url = `https://www.metal-archives.com/band/read-more/id/` + platformId
	var log = logger.New()
	defer func() {
		if e := recover(); e != nil {
			log.SetData(logger.Data{
				"url":        url,
				"stacktrace": string(debug.Stack()),
			}).Error(e)
		}
	}()

	requester := utils.NewClient()
	response := requester.MakeGetRequest(url)
	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.SetData(logger.Data{
			"response":   response,
			"url":        url,
			"stacktrace": string(debug.Stack()),
		}).Error(err)
	}

	return doc.Text()
}
