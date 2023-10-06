package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	productspb "github.com/wathuta/technical_test/protos_gen/products"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestValidateProductCategory_ValidCategory(t *testing.T) {
	validCategories := []ProductCategory{
		Electronics,
		Clothing,
		Books,
		Food,
		Toys,
		Other,
	}

	for _, category := range validCategories {
		valid := ValidateProductCategory(category)
		assert.True(t, valid, "Expected category to be valid: %s", category)
	}
}

func TestValidateProductCategory_InvalidCategory(t *testing.T) {
	invalidCategory := ProductCategory("InvalidCategory")
	valid := ValidateProductCategory(invalidCategory)
	assert.False(t, valid, "Expected category to be invalid: %s", invalidCategory)
}

func TestProductFromProto(t *testing.T) {
	productProto := &productspb.Product{
		ProductId: "12345",
		Name:      "Sample Product",
		Sku:       "SKU123",
		Category:  productspb.ProductCategory_BOOKS, // Unknown category
		Attributes: &productspb.ProductAttributes{
			Brand: "Sample Brand",
			Model: "Model XYZ",
			Price: 99.99,
		},
		StockQuantity: 100,
		IsAvailable:   true,
		CreatedAt:     timestamppb.Now(),
		UpdatedAt:     timestamppb.Now(),
		DeletedAt:     timestamppb.Now(),
	}

	product := ProductFromProto(productProto)

	// Check if fields are correctly mapped
	require.NotNil(t, product)
	assert.Equal(t, productProto.ProductId, product.ProductID)
	assert.Equal(t, productProto.Name, product.Name)
	assert.Equal(t, productProto.Sku, product.Sku)
	assert.Equal(t, ProductCategory(Books), product.Category)
	assert.Equal(t, productProto.Attributes.Brand, product.Brand)
	assert.Equal(t, productProto.Attributes.Model, product.Model)
	assert.Equal(t, productProto.Attributes.Price, product.Price)
	assert.Equal(t, productProto.StockQuantity, product.StockQuantity)
	assert.Equal(t, productProto.IsAvailable, product.IsAvailable)
}

func TestProductProto(t *testing.T) {
	product := &Product{
		ProductID: "12345",
		Name:      "Sample Product",
		Sku:       "SKU123",
		Category:  Clothing,
		ProductAttributes: ProductAttributes{
			Brand: "Sample Brand",
			Model: "Model XYZ",
			Price: 99.99,
		},
		StockQuantity: 100,
		IsAvailable:   true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		DeletedAt:     time.Now(),
	}

	productProto := product.Proto()

	// Check if fields are correctly mapped
	require.NotNil(t, productProto)
	assert.Equal(t, product.ProductID, productProto.ProductId)
	assert.Equal(t, product.Name, productProto.Name)
	assert.Equal(t, product.Sku, productProto.Sku)
	assert.Equal(t, productspb.ProductCategory_value[string(product.Category)], productspb.ProductCategory_value[string(product.Category)])
	assert.Equal(t, product.Brand, productProto.Attributes.Brand)
	assert.Equal(t, product.Model, productProto.Attributes.Model)
	assert.Equal(t, product.Price, productProto.Attributes.Price)
	assert.Equal(t, product.StockQuantity, productProto.StockQuantity)
	assert.Equal(t, product.IsAvailable, productProto.IsAvailable)
}

func TestUpdateProductMapping(t *testing.T) {
	product := Product{
		Name:     "Updated Product Name",
		Sku:      "Updated SKU",
		Category: Clothing,
		ProductAttributes: ProductAttributes{
			Brand: "Updated Brand",
			Model: "Updated Model",
			Price: 49.99,
		},
		StockQuantity: 50,
		IsAvailable:   false,
	}

	updateFields := []string{
		"name",
		"sku",
		"category",
		"brand",
		"model",
		"price",
		"stock_quantity",
		"is_available",
	}

	updateValues := UpdateProductMapping(updateFields, product)

	expectedValues := map[string]interface{}{
		"name":           product.Name,
		"sku":            product.Sku,
		"category":       product.Category,
		"brand":          product.Brand,
		"model":          product.Model,
		"price":          product.Price,
		"stock_quantity": product.StockQuantity,
		"is_available":   product.IsAvailable,
	}

	assert.Equal(t, expectedValues, updateValues)
}
