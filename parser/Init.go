package parser

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/vortgo/ma-parser/parser/ParseScripts"
	"sync"
)

func Run() {
	var wg sync.WaitGroup
	wg.Add(1) //daemon
	go ParseScripts.ParseUpcomingAlbums()
	go ParseScripts.ParseLastBandUpdate()
	//go ParseScripts.ParseBandList()

	go ParseScripts.ParseAllReviews()
	wg.Wait()
}
