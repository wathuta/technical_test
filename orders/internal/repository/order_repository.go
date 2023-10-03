package repository

import (
	"context"
	"strings"

	"github.com/wathuta/technical_test/orders/internal/common"
	"github.com/wathuta/technical_test/orders/internal/model"
)

func (r *repository) CreateOrder(ctx context.Context, order *model.Order, orderDetails *model.OrderDetails) (*model.Order, *model.OrderDetails, error) {
	// Start a transaction
	tx, err := r.connection.BeginTxx(ctx, nil)
	if err != nil {
		return nil, nil, err
	}

	// Defer rollback in case of error or return
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

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
        (order_id, customer_id, pickup_address, delivery_address, shipping_method, order_status,
        scheduled_pickup_datetime, scheduled_delivery_datetime, tracking_number, payment_method,
        invoice_number, special_instructions, shipping_cost, created_at, updated_at, deleted_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
        RETURNING order_id, customer_id, pickup_address, delivery_address, shipping_method, order_status,
        scheduled_pickup_datetime, scheduled_delivery_datetime, tracking_number, payment_method,
        invoice_number, special_instructions, shipping_cost, created_at, updated_at, deleted_at`

	// Execute the SQL query and scan the result into the createdCustomer struct
	err = tx.QueryRowContext(
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
        (order_details_id, order_id, product_id, quantity, created_at, updated_at, deleted_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING order_details_id, order_id, product_id, quantity, created_at, updated_at, deleted_at
    `

	// Execute the SQL query and scan the result into the createdCustomer struct
	err = tx.QueryRowContext(
		ctx, query,
		orderDetails.OrderDetailsID,
		orderDetails.OrderID,
		orderDetails.ProductID,
		orderDetails.Quantity,
		orderDetails.CreatedAt,
		orderDetails.UpdatedAt,
		orderDetails.DeletedAt,
	).Scan(
		&orderDetails.OrderDetailsID,
		&orderDetails.OrderID,
		&orderDetails.ProductID,
		&orderDetails.Quantity,
		&orderDetails.CreatedAt,
		&orderDetails.UpdatedAt,
		&orderDetails.DeletedAt,
	)
	if err != nil {
		return nil, nil, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, nil, err
	}

	// Return the created customer
	return order, orderDetails, nil
}

func (r *repository) UpdateOrder(ctx context.Context, orderId string, updateFields map[string]interface{}) (*model.Order, error) {
	// Start a SQL transaction
	tx, err := r.connection.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Prepare the UPDATE statement
	query := "UPDATE orders SET "
	namedArgs := make(map[string]interface{})

	// Build the SET clause for each field in the updateFields map
	setClauses := []string{}
	for field, value := range updateFields {
		setClauses = append(setClauses, field+"=:"+field) // Use named placeholders
		namedArgs[field] = value
	}
	query += strings.Join(setClauses, ",") + " WHERE order_id = :order_id"
	namedArgs["order_id"] = orderId

	// Execute the UPDATE statement
	_, err = tx.NamedExec(query, namedArgs)
	if err != nil {
		return nil, err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// Return the updated order (you may need to fetch it from the database again)
	updatedOrder, err := r.GetOrderById(ctx, orderId)
	if err != nil {
		return nil, err
	}

	return updatedOrder, nil
}

func (r *repository) GetOrderById(ctx context.Context, orderId string) (*model.Order, error) {
	order := model.Order{}
	// var addressFromDB interface{}
	query := `SELECT * FROM orders WHERE order_id = $1`

	err := r.connection.Get(&order, query, orderId)
	if err != nil {
		return nil, err
	}
	return &order, nil
}
func (r *repository) GetOrdersByCustomerId(ctx context.Context, customerId string, limit, offset int) ([]model.Order, error) {
	orders := []model.Order{}
	order := model.Order{}
	pickupAddr := &model.Address{}
	deliveryAddr := &model.Address{}

	var pickupAddressToDB interface{}
	var deliveryAddressToDB interface{}

	query := `SELECT * FROM orders WHERE customer_id = $1 LIMIT $2 OFFSET $3`

	rows, err := r.connection.Queryx(query, customerId, limit, offset)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(
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
			return nil, err
		}
		if err := pickupAddr.Scan(pickupAddressToDB); err != nil {
			return nil, err
		}
		if err := deliveryAddr.Scan(deliveryAddressToDB); err != nil {
			return nil, err
		}
		order.PickupAddress = *pickupAddr
		order.DeliveryAddress = *deliveryAddr
		orders = append(orders, order)
	}

	return orders, nil
}
func (r *repository) DeleteOrder(ctx context.Context, orderId string) (*model.Order, error) {
	tx, err := r.connection.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() // Rollback if there's an error

	query := `
        DELETE FROM orders
        WHERE order_id = $1
        RETURNING *
    `
	pickupAddr := &model.Address{}
	deliveryAddr := &model.Address{}

	var pickupAddressFromDB interface{}
	var deliveryAddressFromDB interface{}

	var order model.Order

	// Use the transaction to execute the query and scan the result
	err = tx.QueryRowContext(ctx, query, orderId).
		Scan(
			&order.OrderID,
			&order.CustomerID,
			&pickupAddressFromDB,
			&deliveryAddressFromDB,
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
		// Rollback the transaction in case of an error
		tx.Rollback()
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	if err := pickupAddr.Scan(pickupAddressFromDB); err != nil {
		return nil, err
	}
	if err := deliveryAddr.Scan(deliveryAddressFromDB); err != nil {
		return nil, err
	}
	order.PickupAddress = *pickupAddr
	order.DeliveryAddress = *deliveryAddr

	return &order, nil
}
func (r *repository) GetOrderDetailsById(ctx context.Context, orderDetailsId string) (*model.OrderDetails, error) {
	orderDetails := model.OrderDetails{}
	// var addressFromDB interface{}
	query := `SELECT * FROM order_details WHERE order_details_id = $1`

	err := r.connection.Get(&orderDetails, query, orderDetailsId)
	if err != nil {
		return nil, err
	}
	return &orderDetails, nil
}
func (r *repository) GetOrderDetailsByProductId(ctx context.Context, productId string, limit, offset int) ([]model.OrderDetails, error) {
	orderDetails := []model.OrderDetails{}
	orderDetail := model.OrderDetails{}

	query := `SELECT * FROM order_details WHERE product_id = $1 LIMIT $2 OFFSET $3`

	rows, err := r.connection.Queryx(query, productId, limit, offset)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(
			&orderDetail.OrderDetailsID,
			&orderDetail.OrderID,
			&orderDetail.ProductID,
			&orderDetail.Quantity,
			&orderDetail.CreatedAt,
			&orderDetail.UpdatedAt,
			&orderDetail.DeletedAt,
		)
		if err != nil {
			return nil, err
		}

		orderDetails = append(orderDetails, orderDetail)
	}

	return orderDetails, nil
}

func (r *repository) GetOrderDetailsByOrderId(ctx context.Context, orderId string, limit, offset int) ([]model.OrderDetails, error) {
	orderDetails := []model.OrderDetails{}
	orderDetail := model.OrderDetails{}

	query := `SELECT * FROM order_details WHERE order_id = $1 LIMIT $2 OFFSET $3`

	rows, err := r.connection.Queryx(query, orderId, limit, offset)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(
			&orderDetail.OrderDetailsID,
			&orderDetail.OrderID,
			&orderDetail.ProductID,
			&orderDetail.Quantity,
			&orderDetail.CreatedAt,
			&orderDetail.UpdatedAt,
			&orderDetail.DeletedAt,
		)
		if err != nil {
			return nil, err
		}

		orderDetails = append(orderDetails, orderDetail)
	}

	return orderDetails, nil
}
