package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/vortgo/ma-parser/models"
)

type LatestBandUpdateRepository struct {
	*gorm.DB
}

func MakeLatestBandUpdateRepository() *LatestBandUpdateRepository {
	return &LatestBandUpdateRepository{PostgresDB}
}

func (repo *LatestBandUpdateRepository) FindByBandId(bandId uint) *models.LatestBandUpdate {
	latestBand := models.LatestBandUpdate{BandID: bandId}

	repo.Where(&latestBand).First(&latestBand)
	return &latestBand
}
