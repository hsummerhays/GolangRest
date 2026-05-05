package handlers

import (
	"encoding/json"
	"errors"
	"golangrest/internal/models"
	"golangrest/internal/services"
	"log/slog"
	"net/http"
)

type productService interface {
	GetAllProducts() ([]models.Product, error)
	CreateProduct(product models.Product) (models.Product, error)
}

type ProductHandler struct {
	service productService
}

func NewProductHandler(svc productService) *ProductHandler {
	return &ProductHandler{service: svc}
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
	products, err := h.service.GetAllProducts()
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

	createdProduct, err := h.service.CreateProduct(p)
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
