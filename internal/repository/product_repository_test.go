package repository

import (
	"context"
	"golangrest/internal/models"
	"testing"
)

func TestInMemoryProductRepository_Create(t *testing.T) {
	repo := NewInMemoryProductRepository()
	
	newProduct := models.Product{Name: "Tea", Price: 2.50}
	created, err := repo.Create(context.Background(), newProduct)
	
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	
	if created.ID == 0 {
		t.Errorf("expected product ID to be set, got 0")
	}
	
	if created.Name != newProduct.Name || created.Price != newProduct.Price {
		t.Errorf("created product data does not match input")
	}
}

func TestInMemoryProductRepository_GetAll(t *testing.T) {
	repo := NewInMemoryProductRepository()
	
	// Initial data should have 3 items
	products, err := repo.GetAll(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	
	if len(products) != 3 {
		t.Errorf("expected 3 products initially, got %v", len(products))
	}
}
