package handler

import (
	"context"

	customersPb "github.com/wathuta/technical_test/protos_gen/customers"
	"golang.org/x/exp/slog"
)

func (h *Handler) CreateCustomer(ctx context.Context, CreateCustomerRequest *customersPb.CreateCustomerRequest) (*customersPb.Customer, error) {
	slog.Info("hello")
	return &customersPb.Customer{}, nil
}
func (h *Handler) GetCustomerById(context.Context, *customersPb.GetCustomerByIdRequest) (*customersPb.Customer, error) {
	return &customersPb.Customer{}, nil
}
func (h *Handler) UpdateCustomer(context.Context, *customersPb.UpdateCustomerRequest) (*customersPb.Customer, error) {
	return &customersPb.Customer{}, nil
}
func (h *Handler) DeleteCustomer(context.Context, *customersPb.DeleteCustomerRequest) (*customersPb.DeleteCustomerResponse, error) {
	return &customersPb.DeleteCustomerResponse{}, nil
}
