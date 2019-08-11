package repositories

import (
	"github.com/vortgo/ma-parser/models"
)

type BandRepository struct {
	*DbAndElasticRepository
}

func MakeBandRepository() *BandRepository {
	return &BandRepository{&DbAndElasticRepository{PostgresDB}}
}

func (repo *BandRepository) FindBandByPlatformId(platformId string) *models.Band {
	band := models.Band{PlatformID: platformId}

	repo.Where(&band).First(&band)
	return &band
}
