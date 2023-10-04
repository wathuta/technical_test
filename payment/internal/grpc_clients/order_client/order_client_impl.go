package orderclient

import (
	"context"
	"errors"

	grpcclients "github.com/wathuta/technical_test/payment/internal/grpc_clients"
	orderspb "github.com/wathuta/technical_test/protos_gen/orders"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

type orderClient struct {
	client orderspb.OrderServiceClient
}

func NewOrderClient(host, grpcAuthKey string) (grpcclients.OrderServiceClient, error) {
	conn, err := grpc.Dial(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := orderspb.NewOrderServiceClient(conn)
	return &orderClient{
		client: client,
	}, nil
}
func (oc *orderClient) UpdateOrderDetails(orderId string, status orderspb.OrderStatus) <-chan grpcclients.ServiceResult {
	output := make(chan grpcclients.ServiceResult)

	go func() {
		defer close(output)
		args := &orderspb.UpdateOrderRequest{
			Order: &orderspb.Order{
				OrderId:     orderId,
				OrderStatus: orderspb.OrderStatus_ORDER_STATUS_PROCESSING,
			},
			UpdateMask: &fieldmaskpb.FieldMask{
				Paths: []string{"order_status"},
			},
		}
		res, err := oc.client.UpdateOrder(context.Background(), args)
		if err != nil {
			output <- grpcclients.ServiceResult{Error: err}
		}

		if res.Order.OrderStatus != orderspb.OrderStatus_ORDER_STATUS_PROCESSING {
			output <- grpcclients.ServiceResult{Error: errors.New("failed to update errors")}
		}
		output <- grpcclients.ServiceResult{Result: res.Order, Error: nil}
	}()
	return output
}
