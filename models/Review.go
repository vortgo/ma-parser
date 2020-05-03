package models

import (
	"time"
)

type Review struct {
	ID         uint `gorm:"primary_key"`
	AlbumID    uint
	Album      Album
	Rating     int
	Date       time.Time
	Text       string
	Title      string
	PlatformID string
	Author     string
}
