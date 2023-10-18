package db

import "github.com/elastic/go-elasticsearch/v7"

func ConnectElasticSearch() (*elasticsearch.Client, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
	}
	return elasticsearch.NewClient(cfg)
}
