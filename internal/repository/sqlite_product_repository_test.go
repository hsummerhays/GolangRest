package repository

import (
	"golangrest/internal/models"
	"testing"
)

func newTestSQLiteRepo(t *testing.T) *SQLiteProductRepository {
	t.Helper()
	repo, err := NewSQLiteProductRepository(":memory:")
	if err != nil {
		t.Fatalf("failed to create test repository: %v", err)
	}
	return repo
}

func TestSQLiteProductRepository_GetAll_Seeded(t *testing.T) {
	repo := newTestSQLiteRepo(t)
	products, err := repo.GetAll()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(products) != 3 {
		t.Errorf("expected 3 seeded products, got %d", len(products))
	}
}

func TestSQLiteProductRepository_Create(t *testing.T) {
	repo := newTestSQLiteRepo(t)
	product := models.Product{Name: "Mocha", Price: 5.00}

	created, err := repo.Create(product)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if created.ID == 0 {
		t.Errorf("expected product ID to be set after create")
	}
	if created.Name != product.Name || created.Price != product.Price {
		t.Errorf("created product data does not match input")
	}
}

func TestSQLiteProductRepository_GetAll_AfterCreate(t *testing.T) {
	repo := newTestSQLiteRepo(t)
	if _, err := repo.Create(models.Product{Name: "Mocha", Price: 5.00}); err != nil {
		t.Fatalf("unexpected error creating product: %v", err)
	}

	products, err := repo.GetAll()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(products) != 4 {
		t.Errorf("expected 4 products (3 seeded + 1 created), got %d", len(products))
	}
}
