package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/vortgo/ma-parser/elasticsearch"
	"github.com/vortgo/ma-parser/models"
)

type BandRepository struct {
	*gorm.DB
}

func MakeBandRepository() *BandRepository {
	return &BandRepository{PostgresDB}
}

func (repo *BandRepository) FindBandByPlatformId(platformId string) *models.Band {
	band := models.Band{PlatformID: platformId}

	repo.Where(&band).First(&band)
	return &band
}

func (repo *BandRepository) SaveToElastic(band models.Band) {
	var count int
	repo.DB.Model(&models.Album{}).Where("band_id = ?", band.ID).Count(&count)
	elasticDocument := band.GetIndexJson(count)
	elasticsearch.IndexDataToElastic(&band, elasticDocument)
}
