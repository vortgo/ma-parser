package models

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"time"
)

type SongElastic struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ReindexAt time.Time `json:"reindex_at"`
}

type Song struct {
	gorm.Model
	BandID     uint
	Band       Band
	AlbumID    uint
	Album      Album
	Name       string
	PlatformID int
	Time       string
	Lyrics     string
	Position   int
}

func (song *Song) GetIndexJson() string {
	document := SongElastic{song.ID, song.Name, song.CreatedAt, song.UpdatedAt, time.Now()}
	jsonDoc, _ := json.Marshal(document)

	return string(jsonDoc)
}

func (song *Song) GetId() int {
	return int(song.ID)
}

func (song *Song) GetIndexName() string {
	return `songs`
}

func (song *Song) GetTypeName() string {
	return `songs`
}
