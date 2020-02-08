package models

import (
	"encoding/json"
	"time"
)

type BandElastic struct {
	ID                uint      `json:"id"`
	Name              string    `json:"name"`
	FormedIn          int       `json:"formed_in"`
	AlbumsCount       int       `json:"albums_count"`
	DescriptionLength int       `json:"description_length"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	ReindexAt         time.Time `json:"reindex_at"`
}

type Band struct {
	ID            uint `gorm:"primary_key" json:"id"`
	Name          string
	Status        string
	CountryID     int `gorm:"default:null"`
	Country       *Country
	FormedIn      int
	YearsActive   string
	LabelID       int `gorm:"default:null"`
	Label         *Label
	Description   string
	ImageLogo     string
	ImageBand     string
	Genres        []*Genre        `gorm:"many2many:bands_genres"`
	LyricalThemes []*LyricalTheme `gorm:"many2many:bands_lyrical_themes"`
	PlatformID    string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time `sql:"index"`
}

func (band *Band) GetIndexJson(countAlbums int) string {
	document := BandElastic{
		ID:                band.ID,
		Name:              band.Name,
		AlbumsCount:       countAlbums,
		FormedIn:          band.FormedIn,
		DescriptionLength: len(band.Description),
		Status:            band.Status,
		CreatedAt:         band.CreatedAt,
		UpdatedAt:         band.UpdatedAt,
		ReindexAt:         time.Now(),
	}
	jsonDoc, _ := json.Marshal(document)

	return string(jsonDoc)
}

func (band *Band) GetId() int {
	return int(band.ID)
}

func (band *Band) GetIndexName() string {
	return `bands`
}

func (band *Band) GetTypeName() string {
	return `bands`
}
