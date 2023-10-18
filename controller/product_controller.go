package controller

import (
	"golang-elastic-search/model"
	"golang-elastic-search/operation"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ProductController struct {
	operation *operation.Operation
}

func NewProductController(operation *operation.Operation) *ProductController {
	return &ProductController{operation: operation}
}

func (pc *ProductController) GetAllProductHandler(c *fiber.Ctx) error {
	listProduct, err := pc.operation.GetAllProducts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(listProduct)
}

func (pc *ProductController) GetProductByIdHandler(c *fiber.Ctx) error {
	// Extract the :id parameter from the URL path
	idStr := c.Params("id")

	// Convert the idStr to an integer
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid ID format")
	}

	productData, err := pc.operation.GetProductById(idInt)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(productData)
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
