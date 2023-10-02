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
	if req == nil || len(req.CustomerId) == 0 {
		slog.Error("invalid request", "error", errResourceRequired)
		return nil, errResourceRequired
	}
	slog.Debug("update customer", "customer_id", req.CustomerId)

	updateFields := make(map[string]interface{})
	// Check the field mask for each field and add it to the updateFields map
	// Check if req.UpdateMask is null or empty
	if req.UpdateMask == nil || len(req.UpdateMask.Paths) == 0 {
		// If no field mask is provided, assume all fields should be updated
		updateFields["name"] = req.Name
		updateFields["email"] = req.Email
		updateFields["phone_number"] = req.PhoneNumber
		updateFields["address"] = req.Address
		// Add other fields as needed
	} else {
		// If a field mask is provided, update only the specified fields
		mask := req.UpdateMask.Paths
		if common.IsInMask("name", mask) && req.Name != "" {
			updateFields["name"] = req.Name
		}
		if common.IsInMask("email", mask) && req.Email != "" {
			updateFields["email"] = req.Email
		}
		if common.IsInMask("phone_number", mask) && req.PhoneNumber != "" {
			updateFields["phone_number"] = req.PhoneNumber
		}
		if common.IsInMask("address", mask) && req.Address != "" {
			updateFields["address"] = req.Address
		}
		// Add other fields as needed
	}

	customerUUID, err := uuid.Parse(req.CustomerId)
	if err != nil {
		slog.Error("invalid customer uuid value", "customerUUID", customerUUID, "error", err, req)
		return nil, errBadRequest
	}

	// Check if there are no fields to update
	if len(updateFields) == 0 {
		slog.Debug("no fields to update")

		resource, err := h.repo.GetCustomerById(ctx, req.CustomerId)
		if err != nil {
			if err == sql.ErrNoRows {
				slog.Error("customer with the given id not found", "customer_id", customerUUID, "error", err)
				return nil, errNotFound
			}
			slog.Error("failed to get customer from db", "error", err)
			return nil, errInternal
		}

		return &customersPb.UpdateCustomerResponse{Customer: resource.Proto()}, nil
	}

	updateCustomerModels := model.UpdateCustomerFromProto(req)
	if err := common.ValidateGeneric(updateCustomerModels); err != nil {
		slog.Error("failed to validate customer resource", "error", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	resource, err := h.repo.UpdateCustomerFields(ctx, customerUUID.String(), updateFields)
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Error("customer with the given id not found", "customer_id", customerUUID, "error", err)
			return nil, errNotFound
		}
		slog.Error("failed to update customer from db", "error", err)
		return nil, errInternal
	}

	slog.Debug("update customer successful")
	return &customersPb.UpdateCustomerResponse{Customer: resource.Proto()}, nil
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
