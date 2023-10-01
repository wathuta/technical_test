package handler

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/wathuta/technical_test/orders/internal/common"
	"github.com/wathuta/technical_test/orders/internal/model"
	customersPb "github.com/wathuta/technical_test/protos_gen/customers"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) CreateCustomer(ctx context.Context, req *customersPb.CreateCustomerRequest) (*customersPb.Customer, error) {
	if req == nil || req.Customer == nil {
		slog.Error("invalid request", "error", errResourceRequired)
		return nil, errResourceRequired
	}
	slog.Debug("creating customer", "customer_email", req.Customer.Email)

	resource := model.CustomerFromProto(req.Customer)
	// generate id for new resource
	resource.CustomerID = uuid.New()
	resource.CreatedAt = time.Now()

	if err := common.ValidateGeneric(resource); err != nil {
		slog.Error("failed to validate customer resource", "error", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	//persist to db
	resource, err := h.repo.CreateCustomer(ctx, resource)
	if err != nil {
		slog.Error("failed to create customer in db", "error", err)
		return nil, errInternal
	}
	slog.Debug("create customer successful")
	return resource.Proto(), nil
}
func (h *Handler) GetCustomerById(ctx context.Context, req *customersPb.GetCustomerByIdRequest) (*customersPb.Customer, error) {
	if req == nil || len(req.CustomerId) == 0 {
		slog.Error("invalid request", "error", errResourceRequired)
		return nil, errResourceRequired
	}
	slog.Debug("get customer by id", "customer_id", req.CustomerId)
	customerUUID, err := uuid.Parse(req.CustomerId)
	if err != nil {
		slog.Error("invalid customer uuid value", "error", err)
		return nil, errBadRequest
	}

	//retrieve from db
	resource, err := h.repo.GetCustomerById(ctx, customerUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Error("customer with the given id not found", "customer_id", customerUUID, "error", err)
			return nil, errNotFound
		}
		slog.Error("failed to get customer from db", "error", err)
		return nil, errInternal
	}
	slog.Debug("get customer successful")
	return resource.Proto(), nil
}
func (h *Handler) UpdateCustomer(ctx context.Context, req *customersPb.UpdateCustomerRequest) (*customersPb.Customer, error) {
	if req == nil || req.Customer == nil {
		slog.Error("invalid request", "error", errResourceRequired)
		return nil, errResourceRequired
	}
	slog.Debug("update customer", "customer_email", req.Customer.Email)

	resource := model.CustomerFromProto(req.Customer)
	
	slog.Debug("update customer successful")
	return &customersPb.Customer{}, nil
}
func (h *Handler) DeleteCustomer(ctx context.Context, req *customersPb.DeleteCustomerRequest) (*customersPb.DeleteCustomerResponse, error) {
	if req == nil || len(req.CustomerId) == 0 {
		slog.Error("invalid request", "error", errResourceRequired)
		return &customersPb.DeleteCustomerResponse{Success: false}, errResourceRequired
	}
	slog.Debug("delete customer", "customer_id", req.CustomerId)

	customerUUID, err := uuid.Parse(req.CustomerId)
	if err != nil {
		slog.Error("invalid customer uuid value", "error", err)
		return &customersPb.DeleteCustomerResponse{Success: false}, errBadRequest
	}

	resource, err := h.repo.DeleteCustomer(ctx, customerUUID)
	if err != nil {
		slog.Error("failed to delete customer from db", "error", err)
		return &customersPb.DeleteCustomerResponse{Success: false}, errInternal
	}
	if resource == nil {
		slog.Error("customer with the given id not found", "customer_id", customerUUID, "error", err)
		return &customersPb.DeleteCustomerResponse{Success: false}, errNotFound
	}

	slog.Debug("delete customer successful")
	return &customersPb.DeleteCustomerResponse{Success: true}, nil
}
