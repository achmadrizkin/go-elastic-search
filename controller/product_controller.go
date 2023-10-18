package controller

import (
	"golang-elastic-search/model"
	"golang-elastic-search/operation"

	"github.com/gofiber/fiber/v2"
)

type ProductController struct {
	operation *operation.Operation
}

func NewProductController(operation *operation.Operation) *ProductController {
	return &ProductController{operation: operation}
}

// Define HTTP handler for creating a product
func (pc *ProductController) CreateProductHandler(c *fiber.Ctx) error {
	// Parse JSON request body into a Product struct
	var product model.Product
	if err := c.BodyParser(&product); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	// Create the product in Elasticsearch
	if err := pc.operation.CreateProduct(product); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// Respond with a success message or product details
	return c.Status(fiber.StatusCreated).JSON(product)
}
