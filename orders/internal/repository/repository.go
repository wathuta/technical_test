package repository

import (
	"context"

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
	GetCustomerById(ctx context.Context, customerID string) (*model.Customer, error)
	UpdateCustomerFields(ctx context.Context, customerID string, updateFields map[string]interface{}) (*model.Customer, error)
	DeleteCustomer(ctx context.Context, customerID string) (*model.Customer, error)

	CreateProduct(ctx context.Context, product *model.Product) (*model.Product, error)
	GetProductById(ctx context.Context, productId string) (*model.Product, error)
	DeleteProduct(ctx context.Context, productId string) (*model.Product, error)
	UpdateProductFields(ctx context.Context, productId string, updateFields map[string]interface{}) (*model.Product, error)
}

type repository struct {
	connection *sqlx.DB
}

func NewRepository(connection *sqlx.DB) Repository {
	return &repository{
		connection: connection,
	}
}
