package models

import (
	"github.com/jinzhu/gorm"
)

type Genre struct {
	gorm.Model
	Name string
}
