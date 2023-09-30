package model

import "github.com/google/uuid"

// Address struct
type Address struct {
	Street     string `validate:"required"`
	City       string `validate:"required"`
	State      string `validate:"required"`
	PostalCode string `validate:"required"`
	Country    string `validate:"required"`
}

// Order struct
type Order struct {
	OrderID                   uuid.UUID `validate:"required"`
	CustomerID                uuid.UUID    `validate:"required"`
	ItemIDs                   []string  `validate:"required"`
	PickupAddress             Address   `validate:"required"`
	DeliveryAddress           Address   `validate:"required"`
	ShippingMethod            string    `validate:"required"`
	OrderStatus               string    `validate:"required"`
	ScheduledPickupDatetime   string    `validate:"required"`
	ScheduledDeliveryDatetime string    `validate:"required"`
	TrackingNumber            string    `validate:"required"`
	PaymentMethod             string    `validate:"required"`
	InvoiceNumber             string    `validate:"required"`
	SpecialInstructions       string
	ShippingCost              float64
	InsuranceInformation      string
}
