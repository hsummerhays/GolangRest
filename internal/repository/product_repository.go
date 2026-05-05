package repository

import (
	"golangrest/internal/models"
	"sync"
)

type ProductRepository interface {
	GetAll() ([]models.Product, error)
	Create(product models.Product) (models.Product, error)
}

type InMemoryProductRepository struct {
	mu       sync.RWMutex
	products []models.Product
	nextID   int
}

func NewInMemoryProductRepository() *InMemoryProductRepository {
	return &InMemoryProductRepository{
		products: []models.Product{
			{ID: 1, Name: "Espresso", Price: 3.50},
			{ID: 2, Name: "Latte", Price: 4.50},
			{ID: 3, Name: "Cappuccino", Price: 4.00},
		},
		nextID: 4,
	}
}

func (r *InMemoryProductRepository) GetAll() ([]models.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.products, nil
}

func (r *InMemoryProductRepository) Create(product models.Product) (models.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	product.ID = r.nextID
	r.nextID++
	r.products = append(r.products, product)
	return product, nil
}
