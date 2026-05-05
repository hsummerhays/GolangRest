package services

import (
	"context"
	"errors"
	"golangrest/internal/models"
	"golangrest/internal/repository"
	"testing"
)

func newTestService() *ProductService {
	return NewProductService(repository.NewInMemoryProductRepository())
}

func TestProductService_GetAllProducts(t *testing.T) {
	svc := newTestService()
	products, err := svc.GetAllProducts(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(products) != 3 {
		t.Errorf("expected 3 seeded products, got %d", len(products))
	}
}

func TestProductService_CreateProduct_Valid(t *testing.T) {
	svc := newTestService()
	created, err := svc.CreateProduct(context.Background(), models.Product{Name: "Tea", Price: 2.50})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if created.ID == 0 {
		t.Errorf("expected product ID to be set")
	}
}

func TestProductService_CreateProduct_EmptyName(t *testing.T) {
	svc := newTestService()
	_, err := svc.CreateProduct(context.Background(), models.Product{Name: "", Price: 2.50})
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestProductService_CreateProduct_NameTooShort(t *testing.T) {
	svc := newTestService()
	_, err := svc.CreateProduct(context.Background(), models.Product{Name: "A", Price: 2.50})
	if err == nil {
		t.Fatal("expected validation error for name shorter than 2 chars, got nil")
	}
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestProductService_CreateProduct_ZeroPrice(t *testing.T) {
	svc := newTestService()
	_, err := svc.CreateProduct(context.Background(), models.Product{Name: "Tea", Price: 0})
	if err == nil {
		t.Fatal("expected validation error for zero price, got nil")
	}
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestProductService_CreateProduct_NegativePrice(t *testing.T) {
	svc := newTestService()
	_, err := svc.CreateProduct(context.Background(), models.Product{Name: "Tea", Price: -1.0})
	if err == nil {
		t.Fatal("expected validation error for negative price, got nil")
	}
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}
