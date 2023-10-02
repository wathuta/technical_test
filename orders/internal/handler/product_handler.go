package handler

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/wathuta/technical_test/orders/internal/common"
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
	if req == nil || len(req.ProductId) == 0 {
		slog.Error("invalid request", "error", errResourceRequired)
		return nil, errResourceRequired
	}

	slog.Debug("update product", "product_id", req.ProductId)

	// Allows update of specific fields
	updateFields := make(map[string]interface{})

	// Check the field mask for each field and add it to the updateFields map
	// Check if req.UpdateMask is null or empty
	if req.UpdateMask == nil || len(req.UpdateMask.Paths) == 0 {
		// If no field mask is provided, assume all fields should be updated
		updateFields["name"] = req.Name
		updateFields["brand"] = req.Attributes.Brand
		updateFields["model"] = req.Attributes.Model
		updateFields["price"] = req.Attributes.Price
		updateFields["category"] = req.Category
		updateFields["is_available"] = req.IsAvailable
		updateFields["sku"] = req.Sku
		updateFields["stock_quantity"] = req.StockQuantity
		// Add other fields as needed
	} else {
		// If a field mask is provided, update only the specified fields
		mask := req.UpdateMask.Paths
		if common.IsInMask("name", mask) && req.Name != "" {
			updateFields["name"] = req.Name
		}
		if common.IsInMask("brand", mask) && req.Attributes.Brand != "" {
			updateFields["brand"] = req.Attributes.Brand
		}
		if common.IsInMask("model", mask) && req.Attributes.Model != "" {
			updateFields["model"] = req.Attributes.Model
		}
		if common.IsInMask("price", mask) && req.Attributes.Price > 0 {
			updateFields["price"] = req.Attributes.Price
		}
		if common.IsInMask("category", mask) {
			updateFields["category"] = req.Category
		}
		if common.IsInMask("is_available", mask) {
			updateFields["is_available"] = req.IsAvailable
		}
		if common.IsInMask("sku", mask) && req.Sku != "" {
			updateFields["sku"] = req.Sku
		}
		if common.IsInMask("stock_quantity", mask) && req.StockQuantity > 0 {
			updateFields["stock_quantity"] = req.StockQuantity
		}
		// Add other fields as needed
	}

	productUUID, err := uuid.Parse(req.ProductId)
	if err != nil {
		slog.Error("invalid product uuid value", "productUUID", productUUID, "error", err, req)
		return nil, errBadRequest
	}

	// Check if there are no fields to update the perform a get
	if len(updateFields) == 0 {
		slog.Debug("no fields to update")

		resource, err := h.repo.GetProductById(ctx, req.ProductId)
		if err != nil {
			if err == sql.ErrNoRows {
				slog.Error("product with the given id not found", "product_id", productUUID, "error", err)
				return nil, errNotFound
			}
			slog.Error("failed to get product from db", "error", err)
			return nil, errInternal
		}

		return &productspb.UpdateProductResponse{Product: resource.Proto()}, nil
	}

	// validate struct fields
	updateProductModel := model.UpdateProductToProto(req)
	if err := common.ValidateGeneric(updateProductModel); err != nil {
		slog.Error("failed to validate product resource", "error", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	//persist in db
	resource, err := h.repo.UpdateProductFields(ctx, productUUID.String(), updateFields)
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Error("product with the given id not found", "product_id", productUUID, "error", err)
			return nil, errNotFound
		}
		slog.Error("failed to update product from db", "error", err)
		return nil, errInternal
	}

	slog.Debug("update product successful")
	return &productspb.UpdateProductResponse{Product: resource.Proto()}, nil
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
