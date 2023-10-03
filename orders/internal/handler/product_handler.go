package handler

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/wathuta/technical_test/orders/internal/common"
	"github.com/wathuta/technical_test/orders/internal/common/fieldmask"
	"github.com/wathuta/technical_test/orders/internal/model"
	productspb "github.com/wathuta/technical_test/protos_gen/products"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) CreateProduct(ctx context.Context, req *productspb.CreateProductRequest) (*productspb.CreateProductResponse, error) {
	var err error
	if req == nil || req.Product == nil {
		slog.Error("invalid request", "error", errResourceRequired)
		return nil, errResourceRequired
	}
	slog.Debug("creating product", "product_name", req.Product.Name)

	resource := model.ProductFromProto(req.Product)

	resource.ProductID = uuid.New().String()
	resource.CreatedAt = time.Now()
	resource.DeletedAt = time.Time{}

	// validate Model
	if err = common.ValidateGeneric(resource); err != nil {
		slog.Error("failed to validate product resource", "error", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// persist in db
	resource, err = h.repo.CreateProduct(ctx, resource)
	if err != nil {
		slog.Error("failed to create product in db", "error", err)
		return nil, errInternal
	}
	slog.Debug("create product successful")
	return &productspb.CreateProductResponse{
		Product: resource.Proto(),
	}, nil
}

func (h *Handler) GetProductById(ctx context.Context, req *productspb.GetProductByIdRequest) (*productspb.GetProductByIdResponse, error) {
	if req == nil || len(req.ProductId) == 0 {
		slog.Error("invalid request", "error", errResourceRequired)
		return nil, errResourceRequired
	}
	slog.Debug("get product by id", "product_id", req.ProductId)

	productUUID, err := uuid.Parse(req.ProductId)
	if err != nil {
		slog.Error("invalid product uuid value", "error", err)
		return nil, errBadRequest
	}

	resource, err := h.repo.GetProductById(ctx, req.ProductId)
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Error("product with the given product_id not found", "product_id", productUUID, "error", err)
			return nil, errNotFound
		}
		slog.Error("failed to get product from db", "error", err)
		return nil, errInternal
	}

	slog.Debug("get product successful")
	return &productspb.GetProductByIdResponse{Product: resource.Proto()}, nil
}
func (h *Handler) UpdateProduct(ctx context.Context, req *productspb.UpdateProductRequest) (*productspb.UpdateProductResponse, error) {
	if req == nil || req.Product == nil {
		slog.Error("invalid request", "error", errResourceRequired)
		return nil, errResourceRequired
	}

	slog.Debug("update product", "product_id", req.Product.ProductId)

	// Allows update of specific fields
	mask, err := fieldmask.New(req.UpdateMask)
	if err != nil {
		slog.Error("invalid request inputs", "error", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	mask.RemoveOutputOnly()

	productUUID, err := uuid.Parse(req.Product.ProductId)
	if err != nil {
		slog.Error("invalid product uuid value", "productUUID", productUUID, "error", err, req)
		return nil, errBadRequest
	}

	product := model.ProductFromProto(req.Product)
	updateProductDetails := model.UpdateProductMapping(mask.Fields, *product)
	if len(mask.Fields) == 0 || len(updateProductDetails) == 0 {

		slog.Debug("no fields to update")
		product, err = h.repo.GetProductById(ctx, productUUID.String())
		if err != nil {
			if err == sql.ErrNoRows {
				slog.Error("product with the given id not found", "product_id", productUUID, "error", err)
				return nil, errNotFound
			}
			slog.Error("failed to get product from db", "error", err)
			return nil, errInternal
		}
	} else {
		//persist in db
		product, err = h.repo.UpdateProductFields(ctx, productUUID.String(), updateProductDetails)
		if err != nil {
			if err == sql.ErrNoRows {
				slog.Error("product with the given id not found", "product_id", productUUID, "error", err)
				return nil, errNotFound
			}
			slog.Error("failed to update product from db", "error", err)
			return nil, errInternal
		}
	}

	slog.Debug("update product successful")
	return &productspb.UpdateProductResponse{Product: product.Proto()}, nil
}

func (h *Handler) DeleteProduct(ctx context.Context, req *productspb.DeleteProductRequest) (*productspb.DeleteProductResponse, error) {
	if req == nil || len(req.ProductId) == 0 {
		slog.Error("invalid request", "error", errResourceRequired)
		return &productspb.DeleteProductResponse{Success: false}, errResourceRequired
	}
	slog.Debug("delete product", "product_id", req.ProductId)

	// check if supplied id is a valid uuid
	productUUID, err := uuid.Parse(req.ProductId)
	if err != nil {
		slog.Error("invalid product uuid value", "error", err)
		return &productspb.DeleteProductResponse{Success: false}, errBadRequest
	}

	resource, err := h.repo.DeleteProduct(ctx, productUUID.String())
	if err != nil {
		slog.Error("failed to delete product from db", "error", err)
		return &productspb.DeleteProductResponse{Success: false}, errInternal
	}
	if resource == nil {
		slog.Error("product with the given id not found", "product_id", productUUID, "error", err)
		return &productspb.DeleteProductResponse{Success: false}, errNotFound
	}
	slog.Debug("delete product successful")

	return &productspb.DeleteProductResponse{Success: true}, nil
}
