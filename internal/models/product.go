package models

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name" validate:"required,min=2"`
	Price float64 `json:"price" validate:"required,gt=0"`
}
