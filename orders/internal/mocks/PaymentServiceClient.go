// Code generated by mockery v2.34.2. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	grpcclients "github.com/wathuta/technical_test/orders/internal/grpc_clients"

	payment "github.com/wathuta/technical_test/protos_gen/payment"
)

// PaymentServiceClient is an autogenerated mock type for the PaymentServiceClient type
type PaymentServiceClient struct {
	mock.Mock
}

// CreatePaymentRequest provides a mock function with given fields: ctx, status
func (_m *PaymentServiceClient) CreatePaymentRequest(ctx context.Context, status *payment.CreatePaymentRequest) chan grpcclients.ServiceResult {
	ret := _m.Called(ctx, status)

	var r0 chan grpcclients.ServiceResult
	if rf, ok := ret.Get(0).(func(context.Context, *payment.CreatePaymentRequest) chan grpcclients.ServiceResult); ok {
		r0 = rf(ctx, status)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(chan grpcclients.ServiceResult)
		}
	}

	return r0
}

// NewPaymentServiceClient creates a new instance of PaymentServiceClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPaymentServiceClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *PaymentServiceClient {
	mock := &PaymentServiceClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
