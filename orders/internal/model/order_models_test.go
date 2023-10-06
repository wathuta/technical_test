package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	orderspb "github.com/wathuta/technical_test/protos_gen/orders" // Import your ordersPb package
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestOrderFromProto(t *testing.T) {
	// Create a mock proto order
	protoOrder := &orderspb.Order{
		CustomerId:                "123",
		OrderId:                   "456",
		ShippingMethod:            "Express",
		OrderStatus:               orderspb.OrderStatus_ORDER_STATUS_SHIPPED,
		TrackingNumber:            "789",
		PaymentMethod:             orderspb.PaymentMethod_PAYMENT_METHOD_MPESA,
		InvoiceNumber:             "INV123",
		ShippingCost:              15.0,
		SpecialInstructions:       "Fragile",
		ScheduledPickupDatetime:   timestamppb.Now(),
		ScheduledDeliveryDatetime: timestamppb.Now(),
	}

	// Convert proto to Order struct
	order := OrderFromProto(protoOrder)

	// Assertions
	assert.NotNil(t, order)
	assert.Equal(t, "123", order.CustomerID)
	assert.Equal(t, "456", order.OrderID)
	assert.Equal(t, "Express", order.ShippingMethod)
	assert.Equal(t, OrderStatus(orderspb.OrderStatus_ORDER_STATUS_SHIPPED.String()), order.OrderStatus)
	assert.Equal(t, "789", order.TrackingNumber)
	assert.Equal(t, PaymentMethod(orderspb.PaymentMethod_PAYMENT_METHOD_MPESA.String()), order.PaymentMethod)
	assert.Equal(t, "INV123", order.InvoiceNumber)
	assert.Equal(t, 15.0, order.ShippingCost)
	assert.Equal(t, "Fragile", order.SpecialInstructions)
}

func TestOrderProto(t *testing.T) {
	// Create a mock order
	order := &Order{
		CustomerID:                "789",
		OrderID:                   "101",
		ShippingMethod:            "Standard",
		OrderStatus:               OrderStatusDelivered,
		TrackingNumber:            "456",
		PaymentMethod:             PaymentMethodCreditCard,
		InvoiceNumber:             "INV456",
		ShippingCost:              10.0,
		SpecialInstructions:       "Handle with care",
		ScheduledPickupDatetime:   time.Now(),
		ScheduledDeliveryDatetime: time.Now(),
	}

	// Convert Order struct to proto
	protoOrder := order.Proto()

	// Assertions
	assert.NotNil(t, protoOrder)
	assert.Equal(t, "789", protoOrder.CustomerId)
	assert.Equal(t, "101", protoOrder.OrderId)
	assert.Equal(t, "Standard", protoOrder.ShippingMethod)
	assert.Equal(t, orderspb.OrderStatus_ORDER_STATUS_DELIVERED, protoOrder.OrderStatus)
	assert.Equal(t, "456", protoOrder.TrackingNumber)
	assert.Equal(t, 10.0, protoOrder.ShippingCost)
	assert.Equal(t, "INV456", protoOrder.InvoiceNumber)
	assert.Equal(t, "Handle with care", protoOrder.SpecialInstructions)
}

func TestUpdateOrderMaping(t *testing.T) {
	// Create a mock order
	order := &Order{
		OrderID:                   "123",
		PickupAddress:             Address{Street: "123 Main St", City: "Cityville", State: "State1", PostalCode: "12345", Country: "USA"},
		DeliveryAddress:           Address{Street: "456 Elm St", City: "Townsville", State: "State2", PostalCode: "67890", Country: "Canada"},
		ShippingMethod:            "Express",
		OrderStatus:               OrderStatusShipped,
		ScheduledPickupDatetime:   time.Now(),
		ScheduledDeliveryDatetime: time.Now(),
		PaymentMethod:             PaymentMethodMpesa,
		InvoiceNumber:             "INV789",
		SpecialInstructions:       "Fragile",
		ShippingCost:              20.0,
	}

	// Define the fields to update
	updateFields := []string{"pickup_address", "shipping_method", "special_instructions"}

	// Create a map of updated values
	updateValues := UpdateOrderMaping(updateFields, *order)

	// Assertions
	assert.NotNil(t, updateValues)
	assert.Len(t, updateValues, 3)
	assert.Equal(t, order.PickupAddress, updateValues["pickup_address"])
	assert.Equal(t, "Express", updateValues["shipping_method"])
	assert.Equal(t, "Fragile", updateValues["special_instructions"])
}
