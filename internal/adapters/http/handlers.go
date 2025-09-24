package http

import (
	"net/http"
	"time"

	"api-products-meli/internal/app"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}

// @Summary      Lista productos
// @Description  Devuelve productos con filtros y orden
// @Tags         products
// @Produce      json
// @Param        q         query   string  false  "Texto a buscar"      example: "celular"
// @Param        minPrice  query   number  false  "Precio mínimo"       example: 1000
// @Param        maxPrice  query   number  false  "Precio máximo"       example: 5000
// @Param        sort      query   string  false  "Campo para ordenar"  Enums(price, rating)  example: price
// @Param        order     query   string  false  "Orden"               Enums(asc, desc)      example: asc
// @Success      200  {object}  map[string]interface{}  "Lista de productos"
// @Failure      500  {object}  map[string]string       "Error interno"
// @Router       /products [get]
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

// @Summary      Detalle de producto
// @Description  Devuelve el detalle de un producto por ID
// @Tags         products
// @Produce      json
// @Param        id   path   string  true  "ID del producto"  example: "MLA12345"
// @Success      200  {object}  map[string]interface{}  "Producto encontrado"
// @Failure      404  {object}  map[string]string       "Producto no encontrado"
// @Router       /products/{id} [get]
func (h *Handler) get(c *gin.Context) {
	id := c.Param("id")
	p, err := h.uc.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, p)
}
