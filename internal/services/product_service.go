package services

import (
	"errors"
	"fmt"
	"golangrest/internal/models"
	"golangrest/internal/repository"

	"github.com/go-playground/validator/v10"
)

var ErrValidation = errors.New("validation failed")

type ProductService struct {
	repo     repository.ProductRepository
	validate *validator.Validate
}

func NewProductService(repo repository.ProductRepository) *ProductService {
	return &ProductService{
		repo:     repo,
		validate: validator.New(),
	}
}

func (s *ProductService) GetAllProducts() ([]models.Product, error) {
	return s.repo.GetAll()
}

func (s *ProductService) CreateProduct(product models.Product) (models.Product, error) {
	if err := s.validate.Struct(product); err != nil {
		return models.Product{}, fmt.Errorf("%w: %s", ErrValidation, err.Error())
	}
	return s.repo.Create(product)
}
