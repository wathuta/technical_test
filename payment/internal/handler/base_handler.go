package handler

import (
	grpcclients "github.com/wathuta/technical_test/payment/internal/grpc_clients"
	"github.com/wathuta/technical_test/payment/internal/platform/mpesa"
	"github.com/wathuta/technical_test/payment/internal/repository"
	"github.com/wathuta/technical_test/protos_gen/payment"
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
	payment.UnimplementedPaymentServiceServer

	repo  repository.Repository
	mpesa mpesa.MpesaService

	grpcclients.OrderServiceClient
}

func New(
	repo repository.Repository,
	mpesaIntegration mpesa.MpesaService,

) *Handler {
	return &Handler{
		repo:  repo,
		mpesa: mpesaIntegration,
	}

}
