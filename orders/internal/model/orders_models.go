package model

import (
	"encoding/json"
	"errors"
	"time"

	orderspb "github.com/wathuta/technical_test/protos_gen/orders"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Address represents the Address message.
type Address struct {
	Street     string `validate:"required" json:"street" db:"street"`
	City       string `validate:"required" json:"city" db:"city"`
	State      string `validate:"required" json:"state" db:"state"`
	PostalCode string `validate:"required" json:"postal_code" db:"postal_code"`
	Country    string `validate:"required" json:"country" db:"country"`
}

// OrderStatus represents the possible order statuses.
type OrderStatus string

const (
	// OrderStatusPending represents the "Pending" order status.
	OrderStatusPending OrderStatus = "ORDER_STATUS_PENDING"
	// OrderStatusShipped represents the "Shipped" order status.
	OrderStatusShipped OrderStatus = "ORDER_STATUS_SHIPPED"
	// OrderStatusProcessing represents the "Processing" order status.
	OrderStatusProcessing OrderStatus = "ORDER_STATUS_PROCESSING"
	// OrderStatusCanceled represents the "Canceled" order status.
	OrderStatusCanceled OrderStatus = "ORDER_STATUS_CANCELLED"
	// OrderStatusDelivered represents the "Delivered" order status.
	OrderStatusDelivered OrderStatus = "ORDER_STATUS_DELIVERED"
	// Add more order statuses as needed.
	OrderUnspecified OrderStatus = "ORDER_STATUS_UNSPECIFIED"
)

// PaymentMethod represents the possible payment methods.
type PaymentMethod string

const (
	// PaymentMethodCreditCard represents the "Credit Card" payment method.
	PaymentMethodCreditCard PaymentMethod = "Credit Card"
	// Add more payment methods as needed.
	PaymentMethodMpesa PaymentMethod = "Mpesa"
)

// Order represents the Order message.
type Order struct {
	OrderID                   string        `validate:"required" db:"order_id"`
	CustomerID                string        `validate:"required" db:"customer_id"`
	PickupAddress             Address       `validate:"required" db:"pickup_address"`
	DeliveryAddress           Address       `validate:"required" db:"delivery_address"`
	ShippingMethod            string        `validate:"required" db:"shipping_method"`
	OrderStatus               OrderStatus   `validate:"required" db:"order_status"`
	ScheduledPickupDatetime   time.Time     `validate:"required" db:"scheduled_pickup_datetime"`
	ScheduledDeliveryDatetime time.Time     `validate:"required" db:"scheduled_delivery_datetime"`
	TrackingNumber            string        `validate:"required" db:"tracking_number"`
	PaymentMethod             PaymentMethod `validate:"required" db:"payment_method"`
	InvoiceNumber             string        `validate:"required" db:"invoice_number"`
	SpecialInstructions       string        `db:"special_instructions"`
	ShippingCost              float64       `validate:"required" db:"shipping_cost"`
	CreatedAt                 time.Time     `db:"created_at"`
	UpdatedAt                 time.Time     `db:"updated_at"`
	DeletedAt                 time.Time     `db:"deleted_at"`
	// Add more fields as needed for orders.
}

// OrderDetails represents the OrderDetails message.
type OrderDetails struct {
	OrderDetailsID string    `validate:"required" db:"order_details_id"`
	OrderID        string    `validate:"required" db:"order_id"`
	ProductID      string    `validate:"required" db:"product_id"`
	Quantity       int32     `validate:"required" db:"quantity"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
	DeletedAt      time.Time `db:"deleted_at"`
	// Add more fields as needed for order details.
}

func OrderFromProto(req *orderspb.Order) *Order {

	order := &Order{
		CustomerID:                req.CustomerId,
		OrderID:                   req.OrderId,
		ShippingMethod:            req.ShippingMethod,
		OrderStatus:               OrderStatus(req.OrderStatus.String()),
		TrackingNumber:            req.TrackingNumber,
		PaymentMethod:             PaymentMethod(req.PaymentMethod.String()),
		InvoiceNumber:             req.InvoiceNumber,
		ShippingCost:              req.ShippingCost,
		SpecialInstructions:       req.SpecialInstructions,
		ScheduledPickupDatetime:   req.ScheduledPickupDatetime.AsTime(),
		ScheduledDeliveryDatetime: req.ScheduledDeliveryDatetime.AsTime(),
	}
	if req.PickupAddress != nil {
		order.PickupAddress = Address{
			Street:     req.PickupAddress.Street,
			City:       req.PickupAddress.City,
			State:      req.PickupAddress.State,
			PostalCode: req.PickupAddress.PostalCode,
			Country:    req.PickupAddress.Country,
		}
	}

	if req.DeliveryAddress != nil {
		order.DeliveryAddress = Address{
			Street:     req.DeliveryAddress.Street,
			City:       req.DeliveryAddress.City,
			State:      req.DeliveryAddress.State,
			Country:    req.DeliveryAddress.Country,
			PostalCode: req.PickupAddress.PostalCode,
		}
	}
	return order
}

func (o *Order) Proto() *orderspb.Order {
	return &orderspb.Order{
		OrderId:    o.OrderID,
		CustomerId: o.CustomerID,
		PickupAddress: &orderspb.Address{
			Street:     o.PickupAddress.Street,
			City:       o.PickupAddress.City,
			Country:    o.PickupAddress.Country,
			PostalCode: o.PickupAddress.PostalCode,
			State:      o.PickupAddress.State,
		},
		DeliveryAddress: &orderspb.Address{
			Street:     o.DeliveryAddress.Street,
			City:       o.DeliveryAddress.City,
			Country:    o.DeliveryAddress.Country,
			PostalCode: o.DeliveryAddress.PostalCode,
			State:      o.DeliveryAddress.State,
		},
		ShippingMethod:            o.ShippingMethod,
		OrderStatus:               orderspb.OrderStatus(orderspb.OrderStatus_value[string(o.OrderStatus)]),
		TrackingNumber:            o.TrackingNumber,
		ShippingCost:              o.ShippingCost,
		InvoiceNumber:             o.InvoiceNumber,
		PaymentMethod:             orderspb.PaymentMethod(orderspb.PaymentMethod_value[string(o.PaymentMethod)]),
		SpecialInstructions:       o.SpecialInstructions,
		ScheduledPickupDatetime:   timestamppb.New(o.ScheduledPickupDatetime),
		ScheduledDeliveryDatetime: timestamppb.New(o.ScheduledDeliveryDatetime),
		CreatedAt:                 timestamppb.New(o.CreatedAt),
		UpdatedAt:                 timestamppb.New(o.UpdatedAt),
		DeletedAt:                 timestamppb.New(o.DeletedAt),
	}
}

func (od *OrderDetails) Proto() *orderspb.OrderDetails {
	return &orderspb.OrderDetails{
		OrderDetailsId:  od.OrderDetailsID,
		OrderId:         od.OrderID,
		ProductId:       od.ProductID,
		ProductQuantity: od.Quantity,
		CreatedAt:       timestamppb.New(od.CreatedAt),
		UpdatedAt:       timestamppb.New(od.UpdatedAt),
		DeletedAt:       timestamppb.New(od.DeletedAt),
	}
}

func (a *Address) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}

func UpdateOrderMaping(updateFields []string, order Order) map[string]interface{} {
	updateValues := make(map[string]interface{})
	for _, updateField := range updateFields {
		if updateField == "pickup_address" {
			updateValues[updateField] = order.PickupAddress
		}
		if updateField == "delivery_address" {
			updateValues[updateField] = order.DeliveryAddress
		}
		if updateField == "shipping_method" {
			updateValues[updateField] = order.ShippingMethod
		}

		if updateField == "order_status" {
			updateValues[updateField] = order.OrderStatus
		}
		if updateField == "scheduled_pickup_datetime" {
			updateValues[updateField] = order.ScheduledPickupDatetime
		}
		if updateField == "scheduled_delivery_datetime" {
			updateValues[updateField] = order.ScheduledDeliveryDatetime
		}
		if updateField == "payment_method" {
			updateValues[updateField] = order.PaymentMethod
		}
		if updateField == "invoice_number" {
			updateValues[updateField] = order.InvoiceNumber
		}
		if updateField == "special_instructions" {
			updateValues[updateField] = order.SpecialInstructions
		}
		if updateField == "shipping_cost" {
			updateValues[updateField] = order.ShippingCost
		}
	}
	return updateValues
}
