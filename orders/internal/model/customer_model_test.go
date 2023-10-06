package model

import (
	"testing"
	"time"

	customerspb "github.com/wathuta/technical_test/protos_gen/customers"

	"github.com/stretchr/testify/assert"
)

func TestCustomerFromProto(t *testing.T) {
	// Create a mock proto customer
	protoCustomer := &customerspb.Customer{
		Name:        "John Doe",
		Email:       "johndoe@example.com",
		PhoneNumber: "+1234567890",
		Address:     "123 Main St",
	}

	// Convert proto to Customer struct
	customer := CustomerFromProto(protoCustomer)

	// Assertions
	assert.NotNil(t, customer)
	assert.Equal(t, "John Doe", customer.Name)
	assert.Equal(t, "johndoe@example.com", customer.Email)
	assert.Equal(t, "+1234567890", customer.PhoneNumber)
	assert.Equal(t, "123 Main St", customer.Address)
}

func TestCustomerProto(t *testing.T) {
	// Create a mock customer
	customer := &Customer{
		CustomerID:  "123",
		Name:        "Alice",
		Email:       "alice@example.com",
		PhoneNumber: "+9876543210",
		Address:     "456 Elm St",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Convert Customer struct to proto
	protoCustomer := customer.Proto()

	// Assertions
	assert.NotNil(t, protoCustomer)
	assert.Equal(t, "123", protoCustomer.CustomerId)
	assert.Equal(t, "Alice", protoCustomer.Name)
	assert.Equal(t, "alice@example.com", protoCustomer.Email)
	assert.Equal(t, "+9876543210", protoCustomer.PhoneNumber)
	assert.Equal(t, "456 Elm St", protoCustomer.Address)
}

func TestUpdateCustomerMapping(t *testing.T) {
	// Create a mock customer
	customer := &Customer{
		Name:        "Bob",
		Email:       "bob@example.com",
		PhoneNumber: "+5555555555",
		Address:     "789 Oak St",
	}

	// Define the fields to update
	updateFields := []string{"name", "email"}

	// Create a map of updated values
	updateValues := UpdateCustomerMapping(updateFields, *customer)

	// Assertions
	assert.NotNil(t, updateValues)
	assert.Len(t, updateValues, 2)
	assert.Equal(t, "Bob", updateValues["name"])
	assert.Equal(t, "bob@example.com", updateValues["email"])
}
