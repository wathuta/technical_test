package repository

import (
	"context"
	"strings"

	"github.com/wathuta/technical_test/orders/internal/model"
)

func (r *repository) CreateCustomer(ctx context.Context, customer *model.Customer) (*model.Customer, error) {
	// Define the SQL query with the RETURNING clause
	query := `
		INSERT INTO customers
		(customer_id, name, email, phone_number, address, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING customer_id, name, email, phone_number, address, created_at, updated_at, deleted_at
	`

	// Execute the SQL query and scan the result into the createdCustomer struct
	err := r.connection.QueryRowContext(
		ctx, query,
		customer.CustomerID, customer.Name, customer.Email, customer.PhoneNumber,
		customer.Address, customer.CreatedAt, customer.UpdatedAt, customer.DeletedAt,
	).Scan(
		&customer.CustomerID, &customer.Name, &customer.Email, &customer.PhoneNumber,
		&customer.Address, &customer.CreatedAt, &customer.UpdatedAt, &customer.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	// Return the created customer
	return customer, nil
}

func (r *repository) GetCustomerById(ctx context.Context, customerID string) (*model.Customer, error) {
	customer := model.Customer{}

	query := `SELECT * FROM customers WHERE customer_id = $1`

	err := r.connection.Get(&customer, query, customerID)
	if err != nil {
		return nil, err
	}

	// Return query result.
	return &customer, nil

}
func (r *repository) UpdateCustomerFields(ctx context.Context, customerID string, updateFields map[string]interface{}) (*model.Customer, error) {
	// Check if there are fields to update
	if len(updateFields) == 0 {
		return nil, nil // Nothing to update
	}

	// Build the SQL query to update the customer
	query := "UPDATE customers SET "
	params := make(map[string]interface{})

	// Generate the SET clause for each field to update
	setClauses := []string{}
	i := 1
	for field, value := range updateFields {
		setClauses = append(setClauses, field+"=:"+field) // Remove the additional colons and strconv.Itoa(i)
		params[field] = value
		i++
	}

	query += strings.Join(setClauses, ", ") + " WHERE customer_id=:customer_id"
	params["customer_id"] = customerID
	// Execute the SQL query and return the customer based on updated fields
	_, err := r.connection.NamedExecContext(ctx, query, params)
	if err != nil {
		return nil, err
	}

	// Construct the updated customer based on the provided fields
	updatedCustomer := &model.Customer{
		CustomerID: customerID,
	}

	// Update the customer fields based on the updateFields map
	if name, ok := updateFields["name"].(string); ok {
		updatedCustomer.Name = name
	}
	if email, ok := updateFields["email"].(string); ok {
		updatedCustomer.Email = email
	}
	if phoneNumber, ok := updateFields["phone_number"].(string); ok {
		updatedCustomer.PhoneNumber = phoneNumber
	}
	if address, ok := updateFields["address"].(string); ok {
		updatedCustomer.Address = address
	}

	// Return the updated customer
	return updatedCustomer, nil
}

func (r *repository) DeleteCustomer(ctx context.Context, customerID string) (*model.Customer, error) {
    // Start a transaction because this delete operation is not atomic. This ensures that all parts of the queries are a success before making changes
    // if one part fails then all the changes are not made
	tx, err := r.connection.BeginTx(ctx, nil)
    if err != nil {
        return nil, err
    }
    defer tx.Rollback() // Rollback if there's an error

    query := `
        DELETE FROM customers
        WHERE customer_id = $1
        RETURNING *
    `

    var customer model.Customer

    // Use the transaction to execute the query and scan the result
    err = tx.QueryRowContext(ctx, query, customerID).
        Scan(&customer.CustomerID, &customer.Name, &customer.Email, &customer.PhoneNumber, &customer.Address, &customer.CreatedAt, &customer.UpdatedAt, &customer.DeletedAt)

    if err != nil {
        // Rollback the transaction in case of an error
        tx.Rollback()
        return nil, err
    }

    // Commit the transaction
    if err := tx.Commit(); err != nil {
        return nil, err
    }

    return &customer, nil
}
