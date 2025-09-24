package repo

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"api-products-meli/internal/domain/models"
)

// ProductCSVRepo lee productos desde un CSV y hace cache en memoria,
// invalidándolo si cambia el mtime del archivo.
//
// Requisitos mínimos de columnas (case-insensitive):
//
//	id, name, image_url, description, price, rating
//
// Además, cualquier columna que empiece por "spec_" se mapeará a Specs[k]=v
// (por ejemplo: spec_Bateria -> Specs["Bateria"])
type ProductCSVRepo struct {
	path     string
	mu       sync.RWMutex
	modTime  time.Time
	products []models.Product
}

func NewProductCSVRepo(path string) *ProductCSVRepo {
	return &ProductCSVRepo{path: path}
}

// loadIfChanged recarga el CSV solo si cambió el mtime.
func (r *ProductCSVRepo) loadIfChanged() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	info, err := os.Stat(r.path)
	if err != nil {
		return fmt.Errorf("stat file: %w", err)
	}
	// Usa cache si no cambió y ya tenemos datos
	if info.ModTime().Equal(r.modTime) && r.products != nil {
		return nil
	}

	f, err := os.Open(filepath.Clean(r.path))
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.FieldsPerRecord = -1 // permite filas con distinto número de columnas
	// Si usas ; como separador:
	// reader.Comma = ';'

	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("read csv: %w", err)
	}
	if len(records) == 0 {
		r.products = nil
		r.modTime = info.ModTime()
		return nil
	}

	// ---- Parse de header
	header := records[0]
	colIndex := func(name string) int {
		for i, h := range header {
			if strings.EqualFold(strings.TrimSpace(h), name) {
				return i
			}
		}
		return -1
	}

	idxID := colIndex("id")
	idxName := colIndex("name")
	idxImg := colIndex("image_url")
	idxDesc := colIndex("description")
	idxPrice := colIndex("price")
	idxRating := colIndex("rating")

	// recolectar columnas spec_*
	specCols := make([]int, 0)
	specNames := make([]string, 0)
	for i, h := range header {
		hTrim := strings.TrimSpace(h)
		if strings.HasPrefix(strings.ToLower(hTrim), "spec_") {
			specCols = append(specCols, i)
			// normalizamos la clave del spec sin el prefijo "spec_"
			specNames = append(specNames, strings.TrimPrefix(hTrim, "spec_"))
		}
	}

	// helper seguro
	get := func(row []string, idx int) string {
		if idx >= 0 && idx < len(row) {
			return strings.TrimSpace(row[idx])
		}
		return ""
	}

	// ---- Parse filas
	items := make([]models.Product, 0, len(records)-1)
	for i := 1; i < len(records); i++ {
		row := records[i]
		if len(row) == 0 {
			continue
		}
		id := get(row, idxID)
		if id == "" {
			continue
		}

		priceStr := get(row, idxPrice)
		ratingStr := get(row, idxRating)

		price, _ := strconv.ParseFloat(strings.ReplaceAll(priceStr, ",", "."), 64)
		rating, _ := strconv.ParseFloat(strings.ReplaceAll(ratingStr, ",", "."), 64)

		// specs dinámicos
		specs := map[string]string{}
		for k, c := range specCols {
			val := get(row, c)
			if val != "" {
				specs[specNames[k]] = val
			}
		}

		items = append(items, models.Product{
			ID:          id,
			Name:        get(row, idxName),
			ImageURL:    get(row, idxImg),
			Description: get(row, idxDesc),
			Price:       price,
			Rating:      rating,
			Specs:       specs, // si tu modelo no tiene Specs, quita esta línea
		})
	}

	r.products = items
	r.modTime = info.ModTime()
	return nil
}

func (r *ProductCSVRepo) List() ([]models.Product, error) {
	if err := r.loadIfChanged(); err != nil {
		return nil, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]models.Product, len(r.products))
	copy(out, r.products)
	return out, nil
}

func (r *ProductCSVRepo) GetByID(id string) (*models.Product, error) {
	if err := r.loadIfChanged(); err != nil {
		return nil, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, p := range r.products {
		if p.ID == id {
			cp := p
			return &cp, nil
		}
	}
	return nil, errors.New("product not found")
}
