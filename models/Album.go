package models

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"time"
)

type AlbumElastic struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	Year         int       `json:"year"`
	BandFormedIn int       `json:"band_formed_in"`
	Type         string    `json:"type"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	ReindexAt    time.Time `json:"reindex_at"`
}

type Album struct {
	gorm.Model
	BandID      uint
	Band        Band
	Type        string
	Name        string
	Year        int
	ReleaseDate time.Time
	PlatformID  int
	LabelID     int `gorm:"default:null"`
	Label       *Label
	Image       string
	TotalTime   string
}

func (album *Album) GetIndexJson(band Band) string {
	document := AlbumElastic{
		ID:           album.ID,
		Name:         album.Name,
		Year:         album.Year,
		Type:         album.Type,
		BandFormedIn: band.FormedIn,
		CreatedAt:    album.CreatedAt,
		UpdatedAt:    album.UpdatedAt,
		ReindexAt:    time.Now()}
	jsonDoc, _ := json.Marshal(document)

	return string(jsonDoc)
}

func (album *Album) GetId() int {
	return int(album.ID)
}

func (album *Album) GetIndexName() string {
	return `albums`
}

func (album *Album) GetTypeName() string {
	return `albums`
}
