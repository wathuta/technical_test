package repository

import (
	"context"

	"github.com/wathuta/technical_test/orders/internal/model"
)

func (p *PostgresRepository) CreateOrder(ctx context.Context, order model.Order)  {}
func (p *PostgresRepository) UpdateOrder(ctx context.Context, order model.Order)  {}
func (p *PostgresRepository) GetOrderById(ctx context.Context, orderId string)    {}
func (p *PostgresRepository) GetOrderByUserId(ctx context.Context, userId string) {}
func (p *PostgresRepository) ListOrder(ctx context.Context, order model.Order)    {}
func (p *PostgresRepository) DeleteOrder(ctx context.Context, orderId string)     {}
