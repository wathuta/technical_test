package handler

import (
	"context"

	productspb "github.com/wathuta/technical_test/protos_gen/products"
)

func (h *Handler) CreateProduct(context.Context, *productspb.CreateProductRequest) (*productspb.Product, error) {
	return &productspb.Product{}, nil
}
func (h *Handler) GetProductById(context.Context, *productspb.GetProductByIdRequest) (*productspb.Product, error) {
	return &productspb.Product{}, nil
}
func (h *Handler) UpdateProduct(context.Context, *productspb.UpdateProductRequest) (*productspb.Product, error) {
	return &productspb.Product{}, nil
}
func (h *Handler) DeleteProduct(context.Context, *productspb.DeleteProductRequest) (*productspb.DeleteProductResponse, error) {
	return &productspb.DeleteProductResponse{}, nil
}
