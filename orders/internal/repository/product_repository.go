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
	// Check if there are fields to update
	if len(updateFields) == 0 {
		return nil, nil // Nothing to update
	}

	// Build the SQL query to update the product
	query := "UPDATE products SET "
	params := make(map[string]interface{})

	// Generate the SET clause for each field to update
	setClauses := []string{}
	i := 1
	for field, value := range updateFields {
		setClauses = append(setClauses, field+"=:"+field) // Remove the additional colons and strconv.Itoa(i)
		params[field] = value
		i++
	}

	query += strings.Join(setClauses, ", ") + " WHERE product_id=:product_id"
	params["product_id"] = productID

	// Execute the SQL query and return the product based on updated fields
	_, err := r.connection.NamedExecContext(ctx, query, params)
	if err != nil {
		return nil, err
	}

	// Construct the updated product based on the provided fields
	updatedProduct := &model.Product{
		ProductID: productID,
	}

	// Update the product fields based on the updateFields map
	if name, ok := updateFields["name"].(string); ok {
		updatedProduct.Name = name
	}
	if brand, ok := updateFields["brand"].(string); ok {
		updatedProduct.Brand = brand
	}
	if model, ok := updateFields["model"].(string); ok {
		updatedProduct.Model = model
	}
	if price, ok := updateFields["price"].(float64); ok {
		updatedProduct.Price = price
	}
	if category, ok := updateFields["category"].(model.ProductCategory); ok {
		updatedProduct.Category = category
	}
	if isAvailable, ok := updateFields["is_available"].(bool); ok {
		updatedProduct.IsAvailable = isAvailable
	}
	if sku, ok := updateFields["sku"].(string); ok {
		updatedProduct.Sku = sku
	}
	if stockQuantity, ok := updateFields["stock_quantity"].(int32); ok {
		updatedProduct.StockQuantity = stockQuantity
	}

	// Return the updated product
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
