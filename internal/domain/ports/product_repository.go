package ports

import "api-products-meli/internal/domain/models"

type ProductRepository interface {
	List() ([]models.Product, error)
	GetByID(id string) (*models.Product, error)
}
