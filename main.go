package main

import (
	"fmt"
	"golang-elastic-search/controller"
	"golang-elastic-search/db"
	"golang-elastic-search/model"
	"golang-elastic-search/operation"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	clientElasticSearch, err := db.ConnectElasticSearch()
	if err != nil {
		fmt.Println("ERROR CONNECT ELASTIC SEARCH: ", err.Error())
	}

	fmt.Println("Connecting to elastic search success")

	model.CreateIndex(clientElasticSearch, "products")

	fmt.Println("success create index, with name of products")

	//
	productOperation := operation.NewOperation(clientElasticSearch)
	productController := controller.NewProductController(productOperation)

	app := fiber.New()

	app.Post("/products", productController.CreateProductHandler)

	// Start the HTTP server
	err = app.Listen(":8085")
	if err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
