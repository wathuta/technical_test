package model

import (
	"time"
)

// ProductCategory represents the category of a product.
type ProductCategory string

const (
	Electronics ProductCategory = "electronics"
	Clothing    ProductCategory = "clothing"
	Books       ProductCategory = "books"
	Food        ProductCategory = "food"
	Toys        ProductCategory = "toys"
	Other       ProductCategory = "other"
)

// ProductAttributes represents attributes of a product.
type ProductAttributes struct {
	Brand string  `validate:"required"`
	Model string  `validate:"required"`
	Price float64 `validate:"required"`
}

// Product represents a product.
type Product struct {
	ProductID     string             `validate:"required"`
	Name          string             `validate:"required"`
	Category      ProductCategory    `validate:"required"`
	Attributes    *ProductAttributes `validate:"required"`
	StockQuantity int32              `validate:"required"`
	IsAvailable   bool               `validate:"required"`
	CreatedAt     *time.Time         `validate:"-"`
	UpdatedAt     *time.Time         `validate:"-"`
}
