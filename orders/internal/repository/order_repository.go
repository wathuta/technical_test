package repository

import (
	"context"

	"github.com/wathuta/technical_test/orders/internal/common"
	"github.com/wathuta/technical_test/orders/internal/model"
)

func (r *repository) CreateOrder(ctx context.Context, order *model.Order, order_details *model.OrderDetails) (*model.Order, *model.OrderDetails, error) {
	deliveryAddressToDB, err := common.MarshalToBytes(order.DeliveryAddress)
	if err != nil {
		return nil, nil, err
	}

	pickupAddressToDB, err := common.MarshalToBytes(order.PickupAddress)
	if err != nil {
		return nil, nil, err
	}

	query := `
	INSERT INTO orders
	(order_id , customer_id , pickup_address , delivery_address , shipping_method ,
	order_status , scheduled_pickup_datetime , scheduled_delivery_datetime , tracking_number,
	payment_method , invoice_number , special_instructions , shipping_cost  , created_at , updated_at , deleted_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	RETURNING order_id , customer_id , pickup_address , delivery_address , shipping_method , order_status , scheduled_pickup_datetime , scheduled_delivery_datetime , tracking_number , payment_method , invoice_number , special_instructions , shipping_cost , created_at , updated_at , deleted_at` // Execute the SQL query and scan the result into the createdCustomer struct

	err = r.connection.QueryRowContext(
		ctx, query,
		order.OrderID,
		order.CustomerID,
		pickupAddressToDB,
		deliveryAddressToDB,
		order.ShippingMethod,
		order.OrderStatus,
		order.ScheduledPickupDatetime,
		order.ScheduledDeliveryDatetime,
		order.TrackingNumber,
		order.PaymentMethod,
		order.InvoiceNumber,
		order.SpecialInstructions,
		order.ShippingCost,
		order.CreatedAt,
		order.UpdatedAt,
		order.DeletedAt,
	).Scan(
		&order.OrderID,
		&order.CustomerID,
		&pickupAddressToDB,
		&deliveryAddressToDB,
		&order.ShippingMethod,
		&order.OrderStatus,
		&order.ScheduledPickupDatetime,
		&order.ScheduledDeliveryDatetime,
		&order.TrackingNumber,
		&order.PaymentMethod,
		&order.InvoiceNumber,
		&order.SpecialInstructions,
		&order.ShippingCost,
		&order.CreatedAt,
		&order.UpdatedAt,
		&order.DeletedAt,
	)
	if err != nil {
		return nil, nil, err
	}

	query = `
	INSERT INTO order_details
	( order_details_id , order_id , product_id , quantity , created_at , updated_at , deleted_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING  order_details_id , order_id , product_id , quantity , created_at , updated_at , deleted_at
`

	// Execute the SQL query and scan the result into the createdCustomer struct
	err = r.connection.QueryRowContext(
		ctx, query,
		order_details.OrderDetailsID,
		order_details.OrderID,
		order_details.ProductID,
		order_details.Quantity,
		order_details.CreatedAt,
		order_details.UpdatedAt,
		order_details.DeletedAt,
	).Scan(
		&order_details.OrderDetailsID,
		&order_details.OrderID,
		&order_details.ProductID,
		&order_details.Quantity,
		&order_details.CreatedAt,
		&order_details.UpdatedAt,
		&order_details.DeletedAt,
	)
	if err != nil {
		return nil, nil, err
	}
	// Return the created customer
	return order, order_details, nil
}
func (r *repository) UpdateOrder(ctx context.Context, orderId string, updateFields map[string]interface{}) (*model.Order, error) {
	return &model.Order{}, nil
}
func (r *repository) GetOrderById(ctx context.Context, orderId string) (*model.Order, error) {
	order := model.Order{}

	query := `SELECT * FROM orders WHERE order_id = $1`

	err := r.connection.GetContext(ctx, &order, query, orderId)
	if err != nil {
		return nil, err
	}
	return &order, nil
}
func (r *repository) GetOrdersByCustomerId(ctx context.Context, userId string) ([]model.Order, error) {
	return nil, nil
}
func (r *repository) DeleteOrder(ctx context.Context, orderId string) (*model.Order, error) {
	return &model.Order{}, nil
}
func (r *repository) GetOrderDetailsById(ctx context.Context, orderDetailsId string) (*model.OrderDetails, error) {
	return &model.OrderDetails{}, nil
}
func (r *repository) GetOrderDetailsByProductId(ctx context.Context, productId string) ([]model.OrderDetails, error) {
	return nil, nil
}

func (r *repository) GetOrderDetailsByOrderId(ctx context.Context, orderId string) ([]model.OrderDetails, error) {
	return nil, nil
}
