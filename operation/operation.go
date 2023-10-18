package operation

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"golang-elastic-search/model"
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
		return fmt.Errorf("error: %s", res.Status())
	}

	return nil
}

func (p *Operation) GetProductById(id int) (*model.Product, error) {
	// Construct a search query to match a single product by ID
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"ID": id,
			},
		},
	}

	// Marshal the query into a JSON byte array
	queryBytes, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	// Perform a search request to retrieve the product by ID
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
		return nil, fmt.Errorf("Failed to retrieve product with ID %d", id)
	}

	hitsData, found := hits["hits"].([]interface{})
	if !found || len(hitsData) == 0 {
		return nil, fmt.Errorf("Product with ID %d not found", id)
	}

	firstHit, found := hitsData[0].(map[string]interface{})
	if !found {
		return nil, fmt.Errorf("Failed to retrieve product with ID %d", id)
	}

	source, found := firstHit["_source"].(map[string]interface{})
	if !found {
		return nil, fmt.Errorf("Failed to retrieve product with ID %d", id)
	}

	var product model.Product
	err = mapstructure.Decode(source, &product)
	if err != nil {
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

func (p *Operation) UpdateProductById(id int, updatedProduct model.Product) error {
	// delete
	if err := p.DeleteProductById(id); err != nil {
		return errors.New("error delete product by id: " + err.Error())
	}

	// and insert new data
	if err := p.CreateProduct(updatedProduct); err != nil {
		return errors.New("erorr create product: " + err.Error())
	}

	return nil
}

func (p *Operation) DeleteProductById(id int) error {
	// Create a query to match the product by its ID
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"ID": id,
			},
		},
	}

	// Marshal the query into a JSON byte array
	queryBytes, err := json.Marshal(query)
	if err != nil {
		return err
	}

	// Use the Elasticsearch Delete By Query API to remove the product by ID
	res, err := p.client.DeleteByQuery(
		[]string{"products"},        // Index names (an array for multi-index)
		bytes.NewReader(queryBytes), // Wrap queryBytes in an io.Reader
		p.client.DeleteByQuery.WithContext(context.Background()),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Check for errors in the response
	if res.IsError() {
		return fmt.Errorf("Elasticsearch error: %s", res.Status())
	}

	return nil
}
