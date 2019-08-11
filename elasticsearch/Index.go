package elasticsearch

import (
	"context"
	"github.com/olivere/elastic"
	"github.com/vortgo/ma-parser/utils"
	"log"
	"os"
	"strconv"
)

func IndexDataToElastic(model IndexingModel) {
	ctx := context.Background()
	id := model.GetId()

	jsonDoc := model.GetIndexJson()
	client, err := elastic.NewClient(elastic.SetHttpClient(utils.CustomHttpClient), elastic.SetSniff(false), elastic.SetHealthcheck(false), elastic.SetURL(os.Getenv("ELASTIC_URL")))
	if err != nil {
		log.Printf("Elastic: %s\n", err)
		return
	}

	_, err = client.Index().
		Index(model.GetIndexName()).
		Type(model.GetTypeName()).
		Id(strconv.Itoa(id)).
		BodyString(string(jsonDoc)).
		Do(ctx)

	if err != nil {
		log.Printf("Elastic indexing: %s\n", err)
		return
	}
}
