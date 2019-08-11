package models

import "github.com/jinzhu/gorm"

type BandLink struct {
	gorm.Model
	Name string `gorm: "type:varchar(255);not_null"`
	Url  string `gorm: "type:varchar(255);not_null"`
}
