package handler

import (
	"github.com/wathuta/technical_test/orders/internal/repository"
	"github.com/wathuta/technical_test/protos_gen/customers"
	"github.com/wathuta/technical_test/protos_gen/orders"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	defaultPageSize = 10
	maxPageSize     = 1000
	maxUpdateRetry  = 5
)

var (
	errInternal                   = status.Error(codes.Internal, "internal error")
	errNotFound                   = status.Error(codes.InvalidArgument, "resource not found")
	errResourceRequired           = status.Error(codes.InvalidArgument, "resource required")
	errResourceUpdateMaskRequired = status.Error(codes.InvalidArgument, "resource update mask required")
)

type Handler struct {
	customers.UnimplementedCustomerServiceServer
	orders.UnimplementedOrderServiceServer

	repo repository.Repository
}

func New(
	repo repository.Repository,
) *Handler {
	return &Handler{
		repo: repo,
	}
}
