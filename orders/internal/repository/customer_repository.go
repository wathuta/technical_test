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
	tx, err := r.connection.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Prepare the UPDATE statement
	query := "UPDATE customers SET "
	namedArgs := make(map[string]interface{})

	// Build the SET clause for each field in the updateFields map
	setClauses := []string{}
	for field, value := range updateFields {
		setClauses = append(setClauses, field+"=:"+field) // Use named placeholders
		namedArgs[field] = value
	}
	query += strings.Join(setClauses, ",") + " WHERE customer_id = :customer_id"
	namedArgs["customer_id"] = customerID

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

	// Return the updated customer (you may need to fetch it from the database again)
	updatedCustomer, err := r.GetCustomerById(ctx, customerID)
	if err != nil {
		return nil, err
	}

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
