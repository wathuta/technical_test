package repository

import (
	"context"

	"github.com/wathuta/technical_test/payment/internal/model"
)

func (r *repository) CreatePayment(ctx context.Context, payment *model.Payment) (*model.Payment, error) {
	// Define the SQL query with the RETURNING clause
	query := `
		INSERT INTO payments
		(id, order_id, customer_id, payment_method, merchant_request_id, amount, currency, status, description, shipping_cost, product_cost, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, order_id, customer_id, payment_method, merchant_request_id, amount, currency, status, description, shipping_cost, product_cost, created_at, updated_at
	`

	// Execute the SQL query and scan the result into the createdPayment struct
	err := r.connection.QueryRowContext(
		ctx, query,
		payment.PaymentID, payment.OrderID, payment.CustomerID, payment.PaymentMethod,
		payment.MerchantRequestID, payment.Amount, payment.Currency, payment.Status,
		payment.Description, payment.ShippingCost, payment.ProductCost,
		payment.CreatedAt, payment.UpdatedAt,
	).Scan(
		&payment.PaymentID, &payment.OrderID, &payment.CustomerID, &payment.PaymentMethod,
		&payment.MerchantRequestID, &payment.Amount, &payment.Currency, &payment.Status,
		&payment.Description, &payment.ShippingCost, &payment.ProductCost,
		&payment.CreatedAt, &payment.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	// Return the created payment
	return payment, nil
}

func (r *repository) GetPaymentById(ctx context.Context, paymentID string) (*model.Payment, error) {
	payment := model.Payment{}

	query := `SELECT * FROM payments WHERE id = $1`

	err := r.connection.Get(&payment, query, paymentID)
	if err != nil {
		return nil, err
	}

	// Return query result.
	return &payment, nil
}

func (r *repository) GetPaymentByMerchantRequestId(ctx context.Context, merchantRequestID string) (*model.Payment, error) {
	payment := model.Payment{}

	query := `SELECT * FROM payments WHERE merchant_request_id = $1`

	err := r.connection.Get(&payment, query, merchantRequestID)
	if err != nil {
		return nil, err
	}

	// Return query result.
	return &payment, nil
}

func (r *repository) UpdatePaymentStatus(ctx context.Context, paymentStatus model.PaymentStatus, paymentId string) (*model.Payment, error) {
	// Define the SQL query to update the payment status
	query := `
		UPDATE payments
		SET status = $1
		WHERE id = $2
		RETURNING id, order_id, payment_method, amount, currency, status, created_at, updated_at
	`

	// Execute the SQL query to update the payment status
	payment := model.Payment{}
	err := r.connection.GetContext(ctx, &payment, query, paymentStatus, paymentId)
	if err != nil {
		return nil, err
	}

	return &payment, nil
}
