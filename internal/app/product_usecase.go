package app

import (
	"api-products-meli/internal/domain/models"
	"api-products-meli/internal/domain/ports"
	"sort"
	"strings"
)

type ListFilters struct {
	Q        string
	MinPrice *float64
	MaxPrice *float64
	SortBy   string // "price" | "rating"
	Order    string // "asc" | "desc"
}

type ProductUseCase struct {
	repo ports.ProductRepository
}

func NewProductUseCase(r ports.ProductRepository) *ProductUseCase {
	return &ProductUseCase{repo: r}
}

func (u *ProductUseCase) List(f ListFilters) ([]models.Product, error) {
	items, err := u.repo.List()
	if err != nil {
		return nil, err
	}

	// Filtro por texto simple
	if q := strings.TrimSpace(strings.ToLower(f.Q)); q != "" {
		dst := items[:0]
		for _, p := range items {
			if strings.Contains(strings.ToLower(p.Name), q) ||
				strings.Contains(strings.ToLower(p.Description), q) {
				dst = append(dst, p)
			}
		}
		items = dst
	}

	// Filtros numéricos
	if f.MinPrice != nil {
		dst := items[:0]
		for _, p := range items {
			if p.Price >= *f.MinPrice {
				dst = append(dst, p)
			}
		}
		items = dst
	}
	if f.MaxPrice != nil {
		dst := items[:0]
		for _, p := range items {
			if p.Price <= *f.MaxPrice {
				dst = append(dst, p)
			}
		}
		items = dst
	}

	// Orden
	sortBy, order := f.SortBy, f.Order
	if sortBy == "" {
		sortBy = "price"
	}
	if order == "" {
		order = "asc"
	}

	sort.SliceStable(items, func(i, j int) bool {
		switch sortBy {
		case "rating":
			if order == "desc" {
				return items[i].Rating > items[j].Rating
			}
			return items[i].Rating < items[j].Rating
		default: // price
			if order == "desc" {
				return items[i].Price > items[j].Price
			}
			return items[i].Price < items[j].Price
		}
	})

	return items, nil
}

func (u *ProductUseCase) Get(id string) (*models.Product, error) {
	return u.repo.GetByID(id)
}
