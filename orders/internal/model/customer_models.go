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

func UpdateCustomerMapping(updateFields []string, customer Customer) map[string]interface{} {
	updateValues := make(map[string]interface{})
	for _, updateField := range updateFields {
		if updateField == "name" {
			updateValues[updateField] = customer.Name
		}
		if updateField == "email" {
			updateValues[updateField] = customer.Email
		}
		if updateField == "phone_number" {
			updateValues[updateField] = customer.PhoneNumber
		}
		if updateField == "address" {
			updateValues[updateField] = customer.Address
		}
	}
	return updateValues
}
