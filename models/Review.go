package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Review struct {
	gorm.Model
	AlbumID    uint
	Album      Album
	Rating     string
	date       time.Time
	Text       string
	PlatformID int
}
