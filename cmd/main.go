package main

import (
	"log"
	"os"

	_ "api-products-meli/docs"

	// productos
	prodhttp "api-products-meli/internal/adapters/http"
	prodrepo "api-products-meli/internal/adapters/repo"
	"api-products-meli/internal/app"
)

func main() {
	jsonPath := envOr("PRODUCTS_JSON", "products.json")
	// csvPath := envOr("PRODUCTS_CSV", "products.csv")
	port := envOr("PORT", "8080")

	repo := prodrepo.NewJSONRepo(jsonPath) // adapter de datos JSON
	// repo := prodrepo.NewJSONRepo(csvPath) // adapter de datos CSV

	uc := app.NewProductUseCase(repo)              // casos de uso
	api := prodhttp.NewProductHandler(uc).Routes() // adapter HTTP

	log.Printf("Listening on :%s (data=%s)", port, jsonPath)
	if err := api.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

func envOr(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
