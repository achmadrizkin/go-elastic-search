package model

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v7"
)

type Product struct {
	Name     string `json:"name"`
	ID       int    `json:"id"`
	ImageURL string `json:"image_url"`
}

func CreateIndex(client *elasticsearch.Client, indexName string) error {
	// Define the index settings and mappings as a Go map
	indexSettings := map[string]interface{}{
		"settings": map[string]interface{}{
			"number_of_shards":   1,
			"number_of_replicas": 0,
		},
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"ID":       map[string]interface{}{"type": "integer"},
				"Name":     map[string]interface{}{"type": "text"},
				"ImageURL": map[string]interface{}{"type": "text"},
			},
		},
	}

	// Marshal the index settings to JSON
	indexSettingsJSON, err := json.Marshal(indexSettings)
	if err != nil {
		return err
	}

	// Create the index with the specified settings and mappings
	res, err := client.Indices.Create(
		indexName,
		client.Indices.Create.WithContext(context.Background()),
		client.Indices.Create.WithBody(strings.NewReader(string(indexSettingsJSON))))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Check for any errors in the response
	if res.IsError() {
		return fmt.Errorf("Error creating index: %s", res.Status())
	}

	return nil
}
