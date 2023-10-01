package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/wathuta/technical_test/orders/internal/model"
)

func (r *repository) CreateProduct(ctx context.Context, product *model.Product) (*model.Product, error) {
	return &model.Product{}, nil
}
func (r *repository) GetProductById(ctx context.Context, productId uuid.UUID) (*model.Product, error) {
	return &model.Product{}, nil
}
func (r *repository) UpdateProduct(ctx context.Context, product *model.Product) (*model.Product, error) {
	return &model.Product{}, nil
}
func (r *repository) DeleteProduct(ctx context.Context, productId uuid.UUID) (bool, error) {
	return true, nil
}
