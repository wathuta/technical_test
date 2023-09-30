package repository

import (
	"context"

	"github.com/wathuta/technical_test/orders/internal/model"
)

func (r *repository) CreateCustomer(ctx context.Context, customer *model.Customer) (*model.Customer, error) {
	_, err := r.connection.Exec(
		`INSERT INTO
		customers
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		customer.CustomerID, customer.Name, customer.Email, customer.PhoneNumber, customer.Address, customer.CreatedAt, customer.UpdatedAt, customer.DeletedAt,
	)
	if err != nil {
		return nil, err
	}

	return customer, nil
}
func (r *repository) GetCustomerById(ctx context.Context, customerID string) (*model.Customer, error) {
	customer := model.Customer{}

	query := `SELECT * FROM customers WHERE id = $1`

	err := r.connection.Get(&customer, query, customerID)
	if err != nil {
		return nil, err
	}

	// Return query result.
	return &customer, nil

}
func (r *repository) UpdateCustomer(ctx context.Context, updateCustomer *model.Customer) (*model.Customer, error) {
	// Define query string
	query := `UPDATE customers SET updated_at = $2, email = $3,  name = $4, phone = $5, address = $6 WHERE id = $1`
	// Send query to database.
	_, err := r.connection.Exec(query, updateCustomer.CustomerID, updateCustomer.UpdatedAt, updateCustomer.Email, updateCustomer.Name, updateCustomer.PhoneNumber, updateCustomer.Address)
	if err != nil {
		// Return only error.
		return nil, err
	}

	// This query returns nothing.
	return updateCustomer, nil

}
func (r *repository) DeleteCustomer(ctx context.Context, customerID string) (bool, error) {
	query := `DELETE FROM customers WHERE id = $1`
	res, err := r.connection.Exec(query, customerID)
	if err != nil {
		return false, err
	}
	if n, err := res.RowsAffected(); err == nil && n == 0 {
		return false, err
	}

	return true, nil

}
