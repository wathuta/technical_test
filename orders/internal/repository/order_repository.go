package repository

import (
	"context"

	"github.com/wathuta/technical_test/orders/internal/model"
)

func (r *repository) CreateOrder(ctx context.Context, order *model.Order) {}
func (r *repository) UpdateOrder(ctx context.Context, order *model.Order) {}
func (r *repository) GetOrderById(ctx context.Context, orderId string)    {}
func (r *repository) GetOrderByUserId(ctx context.Context, userId string) {}
func (r *repository) ListOrder(ctx context.Context)                       {}
func (r *repository) DeleteOrder(ctx context.Context, orderId string)     {}
