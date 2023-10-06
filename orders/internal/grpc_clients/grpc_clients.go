package grpcclients

import (
	"context"

	paymentpb "github.com/wathuta/technical_test/protos_gen/payment"
)

// ServiceResult
type ServiceResult struct {
	Result interface{}
	Error  error
}

type PaymentServiceClient interface {
	CreatePaymentRequest(ctx context.Context, status *paymentpb.CreatePaymentRequest) chan ServiceResult
}
