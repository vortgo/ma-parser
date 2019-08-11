package repositories

import (
	"github.com/jinzhu/gorm"
)

type LatestBandUpdateRepository struct {
	*gorm.DB
}

func MakeLatestBandUpdateRepository() *LatestBandUpdateRepository {
	return &LatestBandUpdateRepository{PostgresDB}
}
