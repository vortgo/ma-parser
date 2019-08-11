package models

import (
	"github.com/jinzhu/gorm"
)

type LyricalTheme struct {
	gorm.Model
	Name string
}
