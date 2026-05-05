package services

import (
	"context"
	"errors"
	"fmt"
	"golangrest/internal/client"
	"golangrest/internal/models"
	"golangrest/internal/repository"

	"github.com/go-playground/validator/v10"
)

var ErrValidation = errors.New("validation failed")

type ProductService struct {
	repo          repository.ProductRepository
	pricingClient client.PricingClient
	validate      *validator.Validate
}

func NewProductService(repo repository.ProductRepository, pricingClient client.PricingClient) *ProductService {
	return &ProductService{
		repo:          repo,
		pricingClient: pricingClient,
		validate:      validator.New(),
	}
}

func (s *ProductService) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	return s.repo.GetAll(ctx)
}

func (s *ProductService) CreateProduct(ctx context.Context, product models.Product) (models.Product, error) {
	if err := s.validate.Struct(product); err != nil {
		return models.Product{}, fmt.Errorf("%w: %s", ErrValidation, err.Error())
	}

	// Simulation of an external cross-boundary call with retries and failure handling
	if err := s.pricingClient.ValidatePrice(ctx, product.Name, product.Price); err != nil {
		return models.Product{}, fmt.Errorf("external pricing validation failed: %w", err)
	}

	return s.repo.Create(ctx, product)
}
