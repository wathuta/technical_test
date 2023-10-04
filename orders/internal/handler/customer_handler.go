package handler

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/wathuta/technical_test/orders/internal/common"
	"github.com/wathuta/technical_test/orders/internal/common/fieldmask"
	"github.com/wathuta/technical_test/orders/internal/model"
	customersPb "github.com/wathuta/technical_test/protos_gen/customers"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) CreateCustomer(ctx context.Context, req *customersPb.CreateCustomerRequest) (*customersPb.CreateCustomerResponse, error) {
	if req == nil || req.Customer == nil {
		slog.Error("invalid request", "error", errResourceRequired)
		return nil, errResourceRequired
	}
	slog.Debug("creating customer", "customer_email", req.Customer.Email)

	resource := model.CustomerFromProto(req.Customer)
	// generate id for new resource
	resource.CustomerID = uuid.New().String()
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
	return &customersPb.CreateCustomerResponse{Customer: resource.Proto()}, nil
}
func (h *Handler) GetCustomerById(ctx context.Context, req *customersPb.GetCustomerByIdRequest) (*customersPb.GetCustomerByIdResponse, error) {
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
	resource, err := h.repo.GetCustomerById(ctx, customerUUID.String())
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Error("customer with the given id not found", "customer_id", customerUUID, "error", err)
			return nil, errNotFound
		}
		slog.Error("failed to get customer from db", "error", err)
		return nil, errInternal
	}
	slog.Debug("get customer successful")
	return &customersPb.GetCustomerByIdResponse{Customer: resource.Proto()}, nil
}
func (h *Handler) UpdateCustomer(ctx context.Context, req *customersPb.UpdateCustomerRequest) (*customersPb.UpdateCustomerResponse, error) {
	if req == nil || req.Customer == nil {
		slog.Error("invalid request", "error", errResourceRequired)
		return nil, errResourceRequired
	}
	slog.Debug("update customer", "customer_id", req.Customer.CustomerId)

	// check the mask
	mask, err := fieldmask.New(req.UpdateMask)
	if err != nil {
		slog.Error("invalid request inputs", "error", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	mask.RemoveOutputOnly()

	// validate customer UUID
	customerUUID, err := uuid.Parse(req.Customer.CustomerId)
	if err != nil {
		slog.Error("invalid customer uuid value", "customerUUID", customerUUID, "error", err, req)
		return nil, errBadRequest
	}

	customer := model.CustomerFromProto(req.Customer)
	updatedCustomerDetail := model.UpdateCustomerMapping(mask.Fields, *customer)
	updatedCustomerDetail["updated_at"] = time.Now()
	// if fieldmask is empty perfom get
	if len(mask.Fields) == 0 || len(updatedCustomerDetail) == 0 {
		slog.Debug("no fields to update")

		customer, err = h.repo.GetCustomerById(ctx, req.Customer.CustomerId)
		if err != nil {
			if err == sql.ErrNoRows {
				slog.Error("customer with the given id not found", "customer_id", customerUUID, "error", err)
				return nil, errNotFound
			}
			slog.Error("failed to get customer from db", "error", err)
			return nil, errInternal
		}
		slog.Debug("update customer successful")
		return &customersPb.UpdateCustomerResponse{Customer: customer.Proto()}, nil

	}

	// persist in db
	customer, err = h.repo.UpdateCustomerFields(ctx, customerUUID.String(), updatedCustomerDetail)
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Error("customer with the given id not found", "customer_id", customerUUID, "error", err)
			return nil, errNotFound
		}
		slog.Error("failed to update customer from db", "error", err)
		return nil, errInternal
	}

	slog.Debug("update customer successful")
	return &customersPb.UpdateCustomerResponse{Customer: customer.Proto()}, nil
}
func (h *Handler) DeleteCustomer(ctx context.Context, req *customersPb.DeleteCustomerRequest) (*customersPb.DeleteCustomerResponse, error) {
	if req == nil || len(req.CustomerId) == 0 {
		slog.Error("invalid request", "error", errResourceRequired)
		return &customersPb.DeleteCustomerResponse{Success: false}, errResourceRequired
	}
	slog.Debug("delete customer", "customer_id", req.CustomerId)

	// verify supplied uuid
	customerUUID, err := uuid.Parse(req.CustomerId)
	if err != nil {
		slog.Error("invalid customer uuid value", "error", err)
		return &customersPb.DeleteCustomerResponse{Success: false}, errBadRequest
	}

	resource, err := h.repo.DeleteCustomer(ctx, customerUUID.String())
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
