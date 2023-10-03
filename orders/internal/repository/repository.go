package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/wathuta/technical_test/orders/internal/model"
)

type Repository interface {
	CreateOrder(ctx context.Context, order *model.Order, order_details *model.OrderDetails) (*model.Order, *model.OrderDetails, error)
	UpdateOrder(ctx context.Context, orderId string, updateFields map[string]interface{}) (*model.Order, error)
	GetOrderById(ctx context.Context, orderId string) (*model.Order, error)
	GetOrdersByCustomerId(ctx context.Context, customerId string, limit, offset int) ([]model.Order, error)
	DeleteOrder(ctx context.Context, orderId string) (*model.Order, error)
	GetOrderDetailsById(ctx context.Context, orderDetailsId string) (*model.OrderDetails, error)
	GetOrderDetailsByProductId(ctx context.Context, productId string, limit, offset int) ([]model.OrderDetails, error)
	GetOrderDetailsByOrderId(ctx context.Context, orderId string, limit, offset int) ([]model.OrderDetails, error)

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
