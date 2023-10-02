package model

import (
	"time"

	customersPb "github.com/wathuta/technical_test/protos_gen/customers"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Customer struct {
	CustomerID  string    `validate:"required" db:"customer_id"`
	Name        string    `validate:"required" db:"name"`
	Email       string    `validate:"required,email" db:"email"`
	PhoneNumber string    `validate:"required,e164" db:"phone_number"`
	Address     string    `validate:"required" db:"address"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	DeletedAt   time.Time `db:"deleted_at"`
	// Add more fields as needed for customers.
}

type UpdateCustomer struct {
	Name        string
	Email       string `validate:"email"`
	PhoneNumber string `validate:"e164"`
	Address     string
}

func CustomerFromProto(e *customersPb.Customer) *Customer {
	return &Customer{
		Name:        e.Name,
		Email:       e.Email,
		PhoneNumber: e.PhoneNumber,
		Address:     e.Address,
	}
}
func (c *Customer) Proto() *customersPb.Customer {

	return &customersPb.Customer{
		CustomerId:  c.CustomerID,
		Name:        c.Name,
		Email:       c.Email,
		PhoneNumber: c.PhoneNumber,
		Address:     c.Address,
		CreatedAt:   timestamppb.New(c.CreatedAt),
		UpdatedAt:   timestamppb.New(c.UpdatedAt),
	}
}

func UpdateCustomerFromProto(e *customersPb.UpdateCustomerRequest) *UpdateCustomer {
	updatedCustomer := &UpdateCustomer{}

	if e.Name != "" {
		updatedCustomer.Name = e.Name
	}
	if e.Email != "" {
		updatedCustomer.Email = e.Email
	}
	if e.PhoneNumber != "" {
		updatedCustomer.PhoneNumber = e.PhoneNumber
	}
	if e.Address != "" {
		updatedCustomer.Address = e.Address
	}

	return updatedCustomer
}
