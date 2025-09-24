package main

import (
	"log"
	"os"

	httpad "api-products-meli/internal/adapters/http"
	repoad "api-products-meli/internal/adapters/repo"
	"api-products-meli/internal/app"
)

func main() {
	jsonPath := envOr("PRODUCTS_JSON", "products.json")
	port := envOr("PORT", "8080")

	repo := repoad.NewJSONRepo(jsonPath)  // adapter de datos
	uc := app.NewProductUseCase(repo)     // casos de uso
	api := httpad.NewHandler(uc).Routes() // adapter HTTP

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
