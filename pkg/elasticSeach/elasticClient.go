package elasticSeach

import (
	"github.com/SuperMatch/config"
	"github.com/opensearch-project/opensearch-go"
	"log"
)

var EsClient *opensearch.Client

func CreateElasticClient(conf config.Config) error {

	cfg := opensearch.Config{
		Addresses: conf.ElasticConfig.URL,
		Username:  conf.ElasticConfig.UserName,
		Password:  conf.ElasticConfig.Password,
	}

	es, err := opensearch.NewClient(cfg)

	if err != nil {
		log.Println("error in connecting elasticSearch")
		return err
	}

	EsClient = es
	return nil
}
