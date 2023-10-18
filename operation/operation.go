package operation

import (
	"context"
	"encoding/json"
	"fmt"
	"golang-elastic-search/model"
	"strings"

	"github.com/elastic/go-elasticsearch/v7"
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

func (p *Operation) GetProduct(id int) (*model.Product, error) {
	// Retrieve a product by ID from Elasticsearch
	res, err := p.client.Get("products", fmt.Sprintf("%d", id), p.client.Get.WithContext(context.Background()))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Check for any errors in the response
	if res.IsError() {
		return nil, fmt.Errorf("Error: %s", res.Status())
	}

	var product model.Product
	if err := json.NewDecoder(res.Body).Decode(&product); err != nil {
		return nil, err
	}

	return &product, nil
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
