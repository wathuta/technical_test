package paymentclient

import (
	"context"

	grpcclients "github.com/wathuta/technical_test/orders/internal/grpc_clients"
	paymentpb "github.com/wathuta/technical_test/protos_gen/payment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type orderClient struct {
	client paymentpb.PaymentServiceClient
}

func NewPaymentClient(host string) (grpcclients.PaymentServiceClient, error) {
	conn, err := grpc.Dial(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := paymentpb.NewPaymentServiceClient(conn)
	return &orderClient{
		client: client,
	}, nil
}
func (oc *orderClient) CreatePaymentRequest(ctx context.Context, args *paymentpb.CreatePaymentRequest) chan grpcclients.ServiceResult {
	output := make(chan grpcclients.ServiceResult)

	go func() {
		defer close(output)
		res, err := oc.client.CreatePayment(context.Background(), args)
		if err != nil {
			output <- grpcclients.ServiceResult{Error: err}
		} else {
			output <- grpcclients.ServiceResult{Result: res.Payment, Error: nil}
		}
	}()
	return output
}
