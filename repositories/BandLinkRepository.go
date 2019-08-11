package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/vortgo/ma-parser/models"
)

type BandLinkRepository struct {
	*gorm.DB
}

func MakeBandLinkRepository() *BandLinkRepository {
	return &BandLinkRepository{PostgresDB}
}

func CreateOrUpdate(link *models.BandLink) {

}
