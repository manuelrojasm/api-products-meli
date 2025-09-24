package http

import "api-products-meli/internal/app"

type ListProductQueryDTO struct {
	Q        string   `form:"q"`
	MinPrice *float64 `form:"minPrice"`
	MaxPrice *float64 `form:"maxPrice"`
	Sort     string   `form:"sort"`
	Order    string   `form:"order"`
}

func (q ListProductQueryDTO) ToFilters() app.ListFilters {
	return app.ListFilters{
		Q: q.Q, MinPrice: q.MinPrice, MaxPrice: q.MaxPrice,
		SortBy: q.Sort, Order: q.Order,
	}
}
