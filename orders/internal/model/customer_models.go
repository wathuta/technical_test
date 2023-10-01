package model

import (
	"time"

	"github.com/google/uuid"
	customersPb "github.com/wathuta/technical_test/protos_gen/customers"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Customer struct {
	CustomerID  uuid.UUID `validate:"required"`
	Name        string    `validate:"required"`
	Email       string    `validate:"required,email"`
	PhoneNumber string    `validate:"required"`
	Address     string    `validate:"required"`
	CreatedAt   time.Time `validate:"-"`
	UpdatedAt   time.Time `validate:"-"`
	DeletedAt   time.Time `validate:"-"`
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
		CustomerId:  c.CustomerID.String(),
		Name:        c.Name,
		Email:       c.Email,
		PhoneNumber: c.PhoneNumber,
		Address:     c.Address,
		CreatedAt:   timestamppb.New(c.CreatedAt),
		UpdatedAt:   timestamppb.New(c.UpdatedAt),
	}
}
