package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/vortgo/ma-parser/elasticsearch"
)

type DbAndElasticRepository struct {
	*gorm.DB
}

func (model *DbAndElasticRepository) Save(value interface{}) {
	model.DB.Save(value)

	indexingModel, ok := value.(elasticsearch.IndexingModel)
	if ok {
		elasticsearch.IndexDataToElastic(indexingModel)
	}
}
