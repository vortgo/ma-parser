package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/vortgo/ma-parser/models"
)

type ReviewRepository struct {
	*gorm.DB
}

func MakeReviewRepository() *ReviewRepository {
	return &ReviewRepository{PostgresDB}
}

func (repo *ReviewRepository) FindReviewByPlatformId(platformId string) *models.Review {
	Review := models.Review{PlatformID: platformId}

	repo.Where(&Review).First(&Review)
	return &Review
}
