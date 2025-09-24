package repo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"api-products-meli/internal/domain/models"
)

type JSONRepo struct {
	path     string
	mu       sync.RWMutex
	modTime  time.Time
	products []models.Product
}

func NewJSONRepo(path string) *JSONRepo { return &JSONRepo{path: path} }

func (r *JSONRepo) loadIfChanged() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	info, err := os.Stat(r.path)
	if err != nil {
		return fmt.Errorf("stat file: %w", err)
	}

	// usa cache si no cambió el archivo
	if info.ModTime().Equal(r.modTime) && r.products != nil {
		return nil
	}

	f, err := os.Open(filepath.Clean(r.path))
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	var items []models.Product
	if err := json.NewDecoder(f).Decode(&items); err != nil {
		return fmt.Errorf("decode json: %w", err)
	}

	r.products = items
	r.modTime = info.ModTime()
	return nil
}

func (r *JSONRepo) List() ([]models.Product, error) {
	if err := r.loadIfChanged(); err != nil {
		return nil, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]models.Product, len(r.products))
	copy(out, r.products)
	return out, nil
}

func (r *JSONRepo) GetByID(id string) (*models.Product, error) {
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
