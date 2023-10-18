package operation

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"golang-elastic-search/model"
	"net/http"
	"strings"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/mitchellh/mapstructure"
)

type Operation struct {
	client *elasticsearch.Client
}

func NewOperation(client *elasticsearch.Client) *Operation {
	return &Operation{client: client}
}

func (p *Operation) CreateProduct(product model.Product) error {
	// Create a new document in Elasticsearch
	res, err := p.client.Index(
		"products",
		strings.NewReader(fmt.Sprintf(`{"ID": %d, "Name": "%s", "ImageURL": "%s"}`, product.ID, product.Name, product.ImageURL)),
		p.client.Index.WithContext(context.Background()),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Check for any errors in the response
	if res.IsError() {
		return fmt.Errorf("Error: %s", res.Status())
	}

	return nil
}

func (p *Operation) GetProductById(id int) (*model.Product, error) {
	// Retrieve a product by ID from Elasticsearch
	res, err := p.client.Get("products", fmt.Sprintf("%d", id), p.client.Get.WithContext(context.Background()))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Check if the product doesn't exist
	if res.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("Product with ID %d not found", id)
	}

	// Check for other errors in the response
	if res.IsError() {
		return nil, fmt.Errorf("Error: %s", res.Status())
	}

	var product model.Product
	if err := json.NewDecoder(res.Body).Decode(&product); err != nil {
		return nil, err
	}

	return &product, nil
}

func (p *Operation) GetAllProducts() ([]model.Product, error) {
	// Create a query to match all products
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
	}

	// Marshal the query into a JSON byte array
	queryBytes, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	// Perform a search request to retrieve all products
	res, err := p.client.Search(
		p.client.Search.WithContext(context.Background()),
		p.client.Search.WithIndex("products"),
		p.client.Search.WithBody(bytes.NewReader(queryBytes)), // Wrap queryBytes in an io.Reader
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Check for errors in the response
	if res.IsError() {
		return nil, fmt.Errorf("Elasticsearch error: %s", res.Status())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	// Extract the product data from the response
	hits, found := result["hits"].(map[string]interface{})
	if !found {
		return nil, fmt.Errorf("Failed to retrieve products")
	}

	hitsData, found := hits["hits"].([]interface{})
	if !found {
		return nil, fmt.Errorf("Failed to retrieve products")
	}

	var products []model.Product
	for _, hitData := range hitsData {
		hit, found := hitData.(map[string]interface{})
		if !found {
			continue
		}

		source, found := hit["_source"].(map[string]interface{})
		if !found {
			continue
		}

		var product model.Product
		err := mapstructure.Decode(source, &product)
		if err != nil {
			continue
		}

		products = append(products, product)
	}

	return products, nil
}

func (p *Operation) UpdateProduct(product model.Product) error {
	// Update an existing product in Elasticsearch
	res, err := p.client.Update(
		"products",
		fmt.Sprintf("%d", product.ID),
		strings.NewReader(fmt.Sprintf(`{"doc": {"Name": "%s", "ImageURL": "%s"}}`, product.Name, product.ImageURL)),
		p.client.Update.WithContext(context.Background()),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Check for any errors in the response
	if res.IsError() {
		return fmt.Errorf("Error: %s", res.Status())
	}

	return nil
}

func (p *Operation) DeleteProduct(client *elasticsearch.Client, id int) error {
	// Delete a product by ID from Elasticsearch
	res, err := p.client.Delete("products", fmt.Sprintf("%d", id), p.client.Delete.WithContext(context.Background()))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Check for any errors in the response
	if res.IsError() {
		return fmt.Errorf("Error: %s", res.Status())
	}

	return nil
}
