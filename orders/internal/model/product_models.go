package model

import (
	"time"

	productspb "github.com/wathuta/technical_test/protos_gen/products"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	Brand string  `validate:"required" db:"brand"`
	Model string  `validate:"required" db:"model"`
	Price float64 `validate:"required,numeric,gt=0" db:"price"`
}

// Product represents a product.
type Product struct {
	ProductID         string          `validate:"required" db:"product_id"`
	Name              string          `validate:"required" db:"name"`
	Sku               string          `validate:"required" db:"sku"`
	Category          ProductCategory `validate:"required" db:"category"`
	ProductAttributes `validate:"required"`
	StockQuantity     int32     `validate:"required,gt=0" db:"stock_quantity"`
	IsAvailable       bool      `validate:"required" db:"is_available"`
	CreatedAt         time.Time `validate:"-" db:"created_at"`
	UpdatedAt         time.Time `validate:"-" db:"updated_at"`
	DeletedAt         time.Time `validate:"-" db:"deleted_at"`
}

type UpdateProductRequest struct {
	ProductID     string          `validate:"omitempty"`
	Name          string          `validate:"omitempty"`
	Sku           string          `validate:"omitempty"`
	Category      ProductCategory `validate:"omitempty"`
	Brand         string          `validate:"omitempty"`
	Model         string          `validate:"omitempty"`
	Price         float64         `validate:"omitempty,gt=0"`
	StockQuantity int32           `validate:"omitempty,gt=0"`
	IsAvailable   bool            `validate:"omitempty"`
}

// ValidateProductCategory validates if a Product's Category is a valid enum value.
func ValidateProductCategory(category ProductCategory) bool {
	switch category {
	case Electronics, Clothing, Books, Food, Toys, Other:
		return true
	default:
		return false
	}
}

func ProductFromProto(e *productspb.Product) *Product {
	return &Product{
		Name:      e.Name,
		Sku:       e.Sku,
		ProductID: e.ProductId,
		Category:  ProductCategory(e.Category.String()),
		ProductAttributes: ProductAttributes{
			Brand: e.Attributes.Brand,
			Model: e.Attributes.Model,
			Price: e.Attributes.Price,
		},
		StockQuantity: e.StockQuantity,
		IsAvailable:   e.IsAvailable,
	}
}
func (c *Product) Proto() *productspb.Product {
	return &productspb.Product{
		ProductId: c.ProductID,
		Name:      c.Name,
		Sku:       c.Sku,
		Attributes: &productspb.ProductAttributes{
			Brand: c.Brand,
			Model: c.Model,
			Price: c.Price,
		},
		Category:      productspb.ProductCategory(productspb.ProductCategory_value[string(c.Category)]),
		StockQuantity: c.StockQuantity,
		IsAvailable:   c.IsAvailable,
		CreatedAt:     timestamppb.New(c.CreatedAt),
		UpdatedAt:     timestamppb.New(c.UpdatedAt),
		DeletedAt:     timestamppb.New(c.DeletedAt),
	}
}
func UpdateProductMapping(updateFields []string, product Product) map[string]interface{} {
	updatedProductValues := make(map[string]interface{})

	for _, updateField := range updateFields {
		if updateField == "name" {
			updatedProductValues[updateField] = product.Name
		}
		if updateField == "sku" {
			updatedProductValues[updateField] = product.Sku
		}
		if updateField == "category" {
			updatedProductValues[updateField] = product.Category
		}
		if updateField == "brand" {
			updatedProductValues[updateField] = product.Brand
		}

		if updateField == "model" {
			updatedProductValues[updateField] = product.Model
		}
		if updateField == "price" {
			updatedProductValues[updateField] = product.Price
		}
		if updateField == "stock_quantity" {
			updatedProductValues[updateField] = product.StockQuantity
		}
		if updateField == "is_available" {
			updatedProductValues[updateField] = product.IsAvailable
		}
	}
	return updatedProductValues
}
