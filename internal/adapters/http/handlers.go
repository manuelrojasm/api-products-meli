package http

import (
	"net/http"
	"time"

	"api-products-meli/internal/app"

	"github.com/gin-gonic/gin"
)

type Handler struct{ uc *app.ProductUseCase }

func NewHandler(uc *app.ProductUseCase) *Handler { return &Handler{uc: uc} }

func (h *Handler) Routes() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		c.Writer.Header().Set("X-Response-Time", time.Since(start).String())
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/products", h.list)
	r.GET("/products/:id", h.get)
	return r
}

func (h *Handler) list(c *gin.Context) {
	var dto ListQueryDTO
	_ = c.Bind(&dto) // parse simple

	items, err := h.uc.List(dto.ToFilters())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": len(items), "products": items})
}

func (h *Handler) get(c *gin.Context) {
	id := c.Param("id")
	p, err := h.uc.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, p)
}
