package repository

import (
	"context"
	"strings"

	"github.com/wathuta/technical_test/orders/internal/model"
)

func (r *repository) CreateProduct(ctx context.Context, product *model.Product) (*model.Product, error) {
	query := `
		INSERT INTO products
		(product_id, name, sku, category , brand , model , price , stock_quantity , is_available , created_at , updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING product_id, name, sku, category, brand, model, price, stock_quantity, is_available, created_at, updated_at, deleted_at
	`
	err := r.connection.QueryRowContext(
		ctx, query,
		product.ProductID, product.Name, product.Sku, product.Category, product.Brand,
		product.Model, product.Price, product.StockQuantity, product.IsAvailable,
		product.CreatedAt, product.UpdatedAt, product.DeletedAt,
	).Scan(
		&product.ProductID, &product.Name, &product.Sku, &product.Category, &product.Brand,
		&product.Model, &product.Price, &product.StockQuantity, &product.IsAvailable,
		&product.CreatedAt, &product.UpdatedAt, &product.DeletedAt,
	)
	if err != nil {
		return nil, err
	}

	return product, nil
}
func (r *repository) GetProductById(ctx context.Context, productId string) (*model.Product, error) {
	product := model.Product{}

	query := `SELECT * FROM products WHERE product_id = $1`

	err := r.connection.GetContext(ctx, &product, query, productId)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *repository) UpdateProductFields(ctx context.Context, productID string, updateFields map[string]interface{}) (*model.Product, error) {
	// Start a SQL transaction
	tx, err := r.connection.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Prepare the UPDATE statement
	query := "UPDATE products SET "
	namedArgs := make(map[string]interface{})

	// Build the SET clause for each field in the updateFields map
	setClauses := []string{}
	for field, value := range updateFields {
		setClauses = append(setClauses, field+"=:"+field) // Use named placeholders
		namedArgs[field] = value
	}
	query += strings.Join(setClauses, ",") + " WHERE product_id = :product_id"
	namedArgs["product_id"] = productID

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

	// Return the updated product (you may need to fetch it from the database again)
	updatedProduct, err := r.GetProductById(ctx, productID)
	if err != nil {
		return nil, err
	}

	return updatedProduct, nil
}

func (r *repository) DeleteProduct(ctx context.Context, productId string) (*model.Product, error) {
	// Start a transaction
	tx, err := r.connection.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() // Rollback if there's an error

	query := `
        DELETE FROM products
        WHERE product_id = $1
        RETURNING *
    `

	var product model.Product

	// Use the transaction to execute the query and scan the result
	err = tx.QueryRowContext(ctx, query, productId).
		Scan(&product.ProductID, &product.Name, &product.Sku, &product.Category, &product.Brand, &product.Model, &product.Price, &product.StockQuantity, &product.IsAvailable, &product.CreatedAt, &product.UpdatedAt, &product.DeletedAt)

	if err != nil {
		// Rollback the transaction in case of an error
		tx.Rollback()
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &product, nil
}
