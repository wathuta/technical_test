package handler

import (
	grpcclients "github.com/wathuta/technical_test/orders/internal/grpc_clients"
	"github.com/wathuta/technical_test/orders/internal/repository"
	"github.com/wathuta/technical_test/protos_gen/customers"
	"github.com/wathuta/technical_test/protos_gen/orders"
	"github.com/wathuta/technical_test/protos_gen/products"
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
	errBadRequest                 = status.Error(codes.InvalidArgument, "invalid request payload")
)

type Handler struct {
	customers.UnimplementedCustomerServiceServer
	orders.UnimplementedOrderServiceServer
	products.UnimplementedProductServiceServer

	repo repository.Repository

	paymentclients grpcclients.PaymentServiceClient
}

func New(
	repo repository.Repository,
	clients grpcclients.PaymentServiceClient,
) *Handler {
	return &Handler{
		repo:    repo,
		paymentclients: clients,
	}
}
