package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/wathuta/technical_test/orders/internal/model"
	database "github.com/wathuta/technical_test/orders/internal/platform/postgres"
)

type Repository interface {
	CreateOrder(ctx context.Context, order model.Order)
	UpdateOrder(ctx context.Context, order model.Order)
	GetOrderById(ctx context.Context, orderId string)
	GetOrderByUserId(ctx context.Context, userId string)
	ListOrder(ctx context.Context, order model.Order)
	DeleteOrder(ctx context.Context, orderId string)
}

type PostgresRepository struct {
	connection *sqlx.DB
}

func NewPostgresRepository() (Repository, error) {
	c, err := database.OpenDBConnection()
	if err != nil {
		return nil, err
	}
	return &PostgresRepository{c}, nil
}
