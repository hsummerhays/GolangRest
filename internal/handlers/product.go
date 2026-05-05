package handlers

import (
	"encoding/json"
	"errors"
	"golangrest/internal/models"
	"golangrest/internal/services"
	"golangrest/pkg/worker"
	"context"
	"log/slog"
	"net/http"
)

type productService interface {
	GetAllProducts(ctx context.Context) ([]models.Product, error)
	CreateProduct(ctx context.Context, product models.Product) (models.Product, error)
}

type ProductHandler struct {
	service productService
	pool    *worker.Pool
}

func NewProductHandler(svc productService, pool *worker.Pool) *ProductHandler {
	return &ProductHandler{
		service: svc,
		pool:    pool,
	}
}

// GetProducts godoc
// @Summary      List products
// @Description  get products
// @Tags         products
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.Product
// @Failure      500  {string}  string "Internal Server Error"
// @Router       /products [get]
func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	slog.Info("Handling GET /products request")
	products, err := h.service.GetAllProducts(r.Context())
	if err != nil {
		slog.Error("Error getting products", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(products); err != nil {
		slog.Error("Error encoding products response", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// CreateProduct godoc
// @Summary      Add a product
// @Description  add by json product
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        product  body      models.Product  true  "Add product"
// @Success      201      {object}  models.Product
// @Failure      400      {string}  string "Bad Request"
// @Failure      500      {string}  string "Internal Server Error"
// @Router       /products [post]
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	slog.Info("Handling POST /products request")
	var p models.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		slog.Warn("Error decoding product request body", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	createdProduct, err := h.service.CreateProduct(r.Context(), p)
	if err != nil {
		if errors.Is(err, services.ErrValidation) {
			slog.Warn("Validation error on product creation", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		slog.Error("Error creating product", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createdProduct); err != nil {
		slog.Error("Error encoding created product response", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// CreateProductBatch godoc
// @Summary      Add multiple products in background
// @Description  Accepts an array of products and queues them for background processing
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        products  body      []models.Product  true  "Array of products"
// @Success      202       {string}  string "Accepted for processing"
// @Failure      400       {string}  string "Bad Request"
// @Router       /products/batch [post]
func (h *ProductHandler) CreateProductBatch(w http.ResponseWriter, r *http.Request) {
	slog.Info("Handling POST /products/batch request")
	var products []models.Product
	if err := json.NewDecoder(r.Body).Decode(&products); err != nil {
		slog.Warn("Error decoding batch request body", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Submit each product to the worker pool
	for _, p := range products {
		productToCreate := p // capture loop variable
		h.pool.Submit(func(ctx context.Context) error {
			// Pass the worker's context to the service layer for cancellation and timeouts
			_, err := h.service.CreateProduct(ctx, productToCreate)
			return err
		})
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"status":"Accepted for processing"}`))
}
