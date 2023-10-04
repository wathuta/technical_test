package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/wathuta/technical_test/payment/internal/model"
)

type Repository interface {
	CreatePayment(ctx context.Context, payment *model.Payment) (*model.Payment, error)
	GetPaymentById(ctx context.Context, payment_id string) (*model.Payment, error)
	GetPaymentByMerchantRequestId(ctx context.Context, merchnt_request_id string) (*model.Payment, error)
	UpdatePaymentStatus(ctx context.Context, paymentStatus model.PaymentStatus, paymentId string) (*model.Payment, error)
}

type repository struct {
	connection *sqlx.DB
}

func NewRepository(connection *sqlx.DB) Repository {
	return &repository{
		connection: connection,
	}
}
