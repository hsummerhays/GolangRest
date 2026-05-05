package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"golangrest/internal/models"
	"golangrest/internal/services"
	"net/http"
	"net/http/httptest"
	"testing"
	"golangrest/pkg/worker"
)

type mockProductService struct {
	products []models.Product
	err      error
}

func (m *mockProductService) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	return m.products, m.err
}

func (m *mockProductService) CreateProduct(ctx context.Context, product models.Product) (models.Product, error) {
	if m.err != nil {
		return models.Product{}, m.err
	}
	product.ID = 99
	return product, nil
}

func TestProductHandler_GetProducts(t *testing.T) {
	svc := &mockProductService{
		products: []models.Product{{ID: 1, Name: "Espresso", Price: 3.50}},
	}
	h := NewProductHandler(svc, worker.NewPool(1, 1))

	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	rr := httptest.NewRecorder()
	h.GetProducts(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
	var products []models.Product
	if err := json.NewDecoder(rr.Body).Decode(&products); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(products) != 1 {
		t.Errorf("expected 1 product, got %d", len(products))
	}
}

func TestProductHandler_GetProducts_ServiceError(t *testing.T) {
	h := NewProductHandler(&mockProductService{err: errors.New("db error")}, worker.NewPool(1, 1))

	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	rr := httptest.NewRecorder()
	h.GetProducts(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", rr.Code)
	}
}

func TestProductHandler_CreateProduct(t *testing.T) {
	h := NewProductHandler(&mockProductService{}, worker.NewPool(1, 1))

	body, _ := json.Marshal(models.Product{Name: "Tea", Price: 2.50})
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.CreateProduct(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", rr.Code)
	}
	var created models.Product
	if err := json.NewDecoder(rr.Body).Decode(&created); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if created.ID != 99 {
		t.Errorf("expected ID 99, got %d", created.ID)
	}
}

func TestProductHandler_CreateProduct_InvalidJSON(t *testing.T) {
	h := NewProductHandler(&mockProductService{}, worker.NewPool(1, 1))

	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader([]byte("not json")))
	rr := httptest.NewRecorder()
	h.CreateProduct(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestProductHandler_CreateProduct_ValidationError(t *testing.T) {
	svc := &mockProductService{
		err: fmt.Errorf("%w: name required", services.ErrValidation),
	}
	h := NewProductHandler(svc, worker.NewPool(1, 1))

	body, _ := json.Marshal(models.Product{Name: "T", Price: 2.50})
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.CreateProduct(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestProductHandler_CreateProduct_InternalError(t *testing.T) {
	svc := &mockProductService{err: errors.New("unexpected db failure")}
	h := NewProductHandler(svc, worker.NewPool(1, 1))

	body, _ := json.Marshal(models.Product{Name: "Tea", Price: 2.50})
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.CreateProduct(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", rr.Code)
	}
}
