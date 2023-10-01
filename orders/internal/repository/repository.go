package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/wathuta/technical_test/orders/internal/model"
)

type Repository interface {
	CreateOrder(ctx context.Context, order *model.Order)
	UpdateOrder(ctx context.Context, order *model.Order)
	GetOrderById(ctx context.Context, orderId string)
	GetOrderByUserId(ctx context.Context, userId string)
	ListOrder(ctx context.Context)
	DeleteOrder(ctx context.Context, orderId string)

	CreateCustomer(ctx context.Context, customer *model.Customer) (*model.Customer, error)
	GetCustomerById(ctx context.Context, customerID uuid.UUID) (*model.Customer, error)
	UpdateCustomer(ctx context.Context, updateCustomer *model.Customer) (*model.Customer, error)
	DeleteCustomer(ctx context.Context, customerID uuid.UUID) (*model.Customer, error)
}

type repository struct {
	connection *sqlx.DB
}

func NewRepository(connection *sqlx.DB) Repository {
	return &repository{
		connection: connection,
	}
}
