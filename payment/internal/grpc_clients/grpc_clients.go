package grpcclients

import "github.com/wathuta/technical_test/protos_gen/orders"

// ServiceResult
type ServiceResult struct {
	Result interface{}
	Error  error
}

type OrderServiceClient interface {
	UpdateOrderDetails(orderId string, status orders.OrderStatus) chan ServiceResult
}
