package model

import (
	"time"
)

type Customer struct {
	CustomerID  string    `validate:"required"`
	Name        string    `validate:"required"`
	Email       string    `validate:"required,email"`
	PhoneNumber string    `validate:"required"`
	Address     string    `validate:"required"`
	CreatedAt   time.Time `validate:"-"`
	UpdatedAt   time.Time `validate:"-"`
	DeletedAt   time.Time `validate:"-"`
	// Add more fields as needed for customers.
}
