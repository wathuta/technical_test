package handler

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/wathuta/technical_test/orders/internal/config"
	"github.com/wathuta/technical_test/orders/internal/mocks"
	"github.com/wathuta/technical_test/orders/internal/model"
	productspb "github.com/wathuta/technical_test/protos_gen/products"
	"golang.org/x/exp/slog"
	"google.golang.org/genproto/protobuf/field_mask"
)

type ProductHandlerTestSuite struct {
	suite.Suite

	handler *Handler
	repo    *mocks.Repository

	testUUID  uuid.UUID
	testUUID1 uuid.UUID
}

func TestProductHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(ProductHandlerTestSuite))
}

func (st *ProductHandlerTestSuite) SetupSuite() {
	var programLevel = new(slog.LevelVar) // Info by default
	h := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: programLevel, AddSource: true})
	slog.SetDefault(slog.New(h))
	programLevel.Set(slog.LevelDebug)

	err := config.HasAllEnvVariables()
	st.Require().False(err)
}

func (st *ProductHandlerTestSuite) SetupTest() {
	repo := mocks.NewRepository(st.T())
	client := mocks.NewPaymentServiceClient(st.T())

	st.handler = New(repo, client)
	st.repo = repo
	st.testUUID = uuid.New()
	st.testUUID1 = uuid.New()
}

func (st *ProductHandlerTestSuite) TestCreateProduct_Success() {
	// Create a mock product request
	productRequest := &productspb.CreateProductRequest{
		Product: &productspb.Product{
			Name:     "Sample Product",
			Sku:      "SKU123",
			Category: productspb.ProductCategory_ELECTRONICS,
			Attributes: &productspb.ProductAttributes{
				Brand: "Sample Brand",
				Model: "Sample Model",
				Price: 100.00,
			},
			StockQuantity: 10,
			IsAvailable:   true,
		},
	}

	// Set up expectations for the mock repository
	st.repo.On("CreateProduct", mock.Anything, mock.Anything).Return(&model.Product{
		ProductID: st.testUUID.String(),
		Name:      productRequest.Product.Name,
		Sku:       productRequest.Product.Sku,
		Category:  model.ProductCategory(productRequest.Product.Category.String()),
		ProductAttributes: model.ProductAttributes{
			Brand: productRequest.Product.Attributes.Brand,
			Model: productRequest.Product.Attributes.Model,
			Price: productRequest.Product.Attributes.Price,
		},
		StockQuantity: productRequest.Product.StockQuantity,
		IsAvailable:   productRequest.Product.IsAvailable,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}, nil)

	// Call the CreateProduct function
	response, err := st.handler.CreateProduct(context.Background(), productRequest)

	// Assertions
	st.Require().NoError(err)
	st.Require().NotNil(response)
	st.Require().NotNil(response.Product)
	st.Require().NotEmpty(response.Product.ProductId)
	st.Require().Equal(productRequest.Product.Name, response.Product.Name)
	st.Require().Equal(productRequest.Product.Sku, response.Product.Sku)
	st.Require().Equal(productRequest.Product.Category, response.Product.Category)
	st.Require().Equal(productRequest.Product.Attributes.Brand, response.Product.Attributes.Brand)
	st.Require().Equal(productRequest.Product.Attributes.Model, response.Product.Attributes.Model)
	st.Require().Equal(productRequest.Product.Attributes.Price, response.Product.Attributes.Price)
	st.Require().Equal(productRequest.Product.StockQuantity, response.Product.StockQuantity)
	st.Require().Equal(productRequest.Product.IsAvailable, response.Product.IsAvailable)

	st.repo.AssertExpectations(st.T())
}

func (st *ProductHandlerTestSuite) TestCreateProduct_InvalidRequest() {
	// Create a mock product request with nil Product
	productRequest := &productspb.CreateProductRequest{
		Product: nil,
	}

	// Call the CreateProduct function
	response, err := st.handler.CreateProduct(context.Background(), productRequest)

	// Assertions
	st.Require().Error(err)
	st.Require().Nil(response)

	// Assert expectations for the mock repository (no calls expected)
	st.repo.AssertExpectations(st.T())
}

func (st *ProductHandlerTestSuite) TestCreateProduct_ValidationError() {
	// Create a mock product request with an invalid price
	productRequest := &productspb.CreateProductRequest{
		Product: &productspb.Product{
			Name:     "Sample Product",
			Sku:      "SKU123",
			Category: productspb.ProductCategory_BOOKS,
			Attributes: &productspb.ProductAttributes{
				Brand: "Sample Brand",
				Model: "Sample Model",
				Price: -100.00, // Invalid price (negative value)
			},
			StockQuantity: 10,
			IsAvailable:   true,
		},
	}

	// Call the CreateProduct function
	response, err := st.handler.CreateProduct(context.Background(), productRequest)

	// Assertions
	st.Require().Error(err)
	st.Require().Nil(response)

	// Assert expectations for the mock repository (no calls expected)
	st.repo.AssertExpectations(st.T())
}

func (st *ProductHandlerTestSuite) TestCreateProduct_CreateProductError() {
	// Create a mock product request
	productRequest := &productspb.CreateProductRequest{
		Product: &productspb.Product{
			Name:     "Sample Product",
			Sku:      "SKU123",
			Category: productspb.ProductCategory_BOOKS,
			Attributes: &productspb.ProductAttributes{
				Brand: "Sample Brand",
				Model: "Sample Model",
				Price: 100.00,
			},
			StockQuantity: 10,
			IsAvailable:   true,
		},
	}

	// Set up expectations for the mock repository to return an error
	st.repo.On("CreateProduct", mock.Anything, mock.Anything).Return(nil, errors.New("some error"))

	// Call the CreateProduct function
	response, err := st.handler.CreateProduct(context.Background(), productRequest)

	// Assertions
	st.Require().Error(err)
	st.Require().Nil(response)

	// Assert expectations for the mock repository
	st.repo.AssertExpectations(st.T())
}

func (st *ProductHandlerTestSuite) TestGetProductById_Success() {
	// Create a mock product ID
	productID := st.testUUID.String()

	// Set up expectations for the mock repository to return a product
	product := &model.Product{
		ProductID:     productID,
		Name:          "Sample Product",
		Sku:           "SKU123",
		Category:      model.Electronics,
		StockQuantity: 100,
		IsAvailable:   true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	expectedProduct := product.Proto()
	st.repo.On("GetProductById", mock.Anything, productID).Return(product, nil)

	// Create a GetProductById request
	request := &productspb.GetProductByIdRequest{
		ProductId: productID,
	}

	// Call the GetProductById function
	response, err := st.handler.GetProductById(context.Background(), request)

	// Assertions
	st.Require().NoError(err)
	st.Require().NotNil(response)
	st.Require().NotNil(response.Product)
	st.Require().Equal(productID, response.Product.ProductId)
	st.Require().Equal(expectedProduct.Name, response.Product.Name)
	st.Require().Equal(expectedProduct.Sku, response.Product.Sku)
	st.Require().Equal(expectedProduct.Category, response.Product.Category)
	st.Require().Equal(expectedProduct.StockQuantity, response.Product.StockQuantity)
	st.Require().Equal(expectedProduct.IsAvailable, response.Product.IsAvailable)

	st.repo.AssertExpectations(st.T())
}

func (st *ProductHandlerTestSuite) TestGetProductById_ProductNotFound() {
	// Create a mock product ID
	productID := st.testUUID.String()

	// Set up expectations for the mock repository to return a "not found" error
	st.repo.On("GetProductById", mock.Anything, productID).Return(nil, sql.ErrNoRows)

	// Create a GetProductById request
	request := &productspb.GetProductByIdRequest{
		ProductId: productID,
	}

	// Call the GetProductById function
	response, err := st.handler.GetProductById(context.Background(), request)

	// Assertions
	st.Require().Error(err)
	st.Require().Nil(response)
	st.Require().Equal(errNotFound, err)

	st.repo.AssertExpectations(st.T())
}

func (st *ProductHandlerTestSuite) TestGetProductById_InvalidProductID() {
	// Create an invalid product ID (not a UUID)
	invalidProductID := "invalid-id"

	// Create a GetProductById request with an invalid product ID
	request := &productspb.GetProductByIdRequest{
		ProductId: invalidProductID,
	}

	// Call the GetProductById function
	response, err := st.handler.GetProductById(context.Background(), request)

	// Assertions
	st.Require().Error(err)
	st.Require().Nil(response)
	st.Require().Equal(errBadRequest, err)
}

func (st *ProductHandlerTestSuite) TestGetProductById_GetError() {
	// Create a mock product ID
	productID := st.testUUID.String()

	// Set up expectations for the mock repository to return an error
	st.repo.On("GetProductById", mock.Anything, productID).Return(nil, errors.New("get error"))

	// Create a GetProductById request
	request := &productspb.GetProductByIdRequest{
		ProductId: productID,
	}

	// Call the GetProductById function
	response, err := st.handler.GetProductById(context.Background(), request)

	// Assertions
	st.Require().Error(err)
	st.Require().Nil(response)

	st.repo.AssertExpectations(st.T())
}

func (st *ProductHandlerTestSuite) TestDeleteProduct_Success() {
	// Set up expectations for the mock repository to delete a product
	st.repo.On("DeleteProduct", mock.Anything, mock.Anything).Return(&model.Product{
		ProductID: st.testUUID.String(),
	}, nil)

	// Call the DeleteProduct function
	resp, err := st.handler.DeleteProduct(context.Background(), &productspb.DeleteProductRequest{
		ProductId: st.testUUID.String(),
	})

	// Assertions
	st.Require().NoError(err)
	st.Require().NotNil(resp)
	st.Require().True(resp.Success)

	st.repo.AssertExpectations(st.T())
}

func (st *ProductHandlerTestSuite) TestDeleteProduct_ProductNotFound() {
	// Set up expectations for the mock repository to delete a product that is not found
	st.repo.On("DeleteProduct", mock.Anything, mock.Anything).Return(nil, nil)

	// Call the DeleteProduct function
	resp, err := st.handler.DeleteProduct(context.Background(), &productspb.DeleteProductRequest{
		ProductId: st.testUUID.String(),
	})

	// Assertions
	st.Require().NotNil(err)
	st.Require().NotNil(resp)
	st.Require().False(resp.Success)

	st.repo.AssertExpectations(st.T())
}

func (st *ProductHandlerTestSuite) TestDeleteProduct_InvalidUUIDError() {
	// Call the DeleteProduct function with an invalid UUID
	resp, err := st.handler.DeleteProduct(context.Background(), &productspb.DeleteProductRequest{
		ProductId: "invalid-uuid",
	})

	// Assertions
	st.Require().False(resp.Success)
	st.Require().NotNil(err)
	st.Require().Equal(errBadRequest, err)
}

func (st *ProductHandlerTestSuite) TestDeleteProduct_DBError() {
	// Set up expectations for the mock repository to return an error when deleting a product
	st.repo.On("DeleteProduct", mock.Anything, mock.Anything).Return(nil, errors.New("some error"))

	// Call the DeleteProduct function
	resp, err := st.handler.DeleteProduct(context.Background(), &productspb.DeleteProductRequest{
		ProductId: st.testUUID.String(),
	})

	// Assertions
	st.Require().False(resp.Success)
	st.Require().NotNil(err)
	st.Require().Equal(errInternal, err)

	st.repo.AssertExpectations(st.T())
}

func (st *ProductHandlerTestSuite) TestDeleteProduct_NilRequest() {
	// Call the DeleteProduct function with a nil request
	resp, err := st.handler.DeleteProduct(context.Background(), nil)

	// Assertions
	st.Require().False(resp.Success)
	st.Require().NotNil(err)
}

func (st *ProductHandlerTestSuite) TestUpdateProduct_Success() {
	// Create a mock product request for updating
	productRequest := &productspb.UpdateProductRequest{
		Product: &productspb.Product{
			ProductId:     st.testUUID.String(),
			Name:          "Updated Product Name",
			Sku:           "UPDATED-SKU",
			Category:      productspb.ProductCategory_CLOTHING,
			StockQuantity: 50,
			IsAvailable:   true,
		},
		UpdateMask: &field_mask.FieldMask{
			Paths: []string{"name", "sku", "category", "stock_quantity", "is_available"},
		},
	}

	// Set up expectations for the mock repository to update the product
	st.repo.On("UpdateProductFields", mock.Anything, mock.Anything, mock.Anything).Return(
		&model.Product{
			ProductID:     st.testUUID.String(),
			Name:          productRequest.Product.Name,
			Sku:           productRequest.Product.Sku,
			Category:      model.ProductCategory(productRequest.Product.Category.String()),
			StockQuantity: productRequest.Product.StockQuantity,
			IsAvailable:   productRequest.Product.IsAvailable,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		nil,
	)

	// Call the UpdateProduct function
	response, err := st.handler.UpdateProduct(context.Background(), productRequest)

	// Assertions
	st.Require().NoError(err)
	st.Require().NotNil(response)
	st.Require().NotNil(response.Product)
	st.Require().NotEmpty(response.Product.ProductId)
	st.Require().Equal(productRequest.Product.Name, response.Product.Name)
	st.Require().Equal(productRequest.Product.Sku, response.Product.Sku)
	st.Require().Equal(productRequest.Product.Category, response.Product.Category)
	st.Require().Equal(productRequest.Product.StockQuantity, response.Product.StockQuantity)
	st.Require().Equal(productRequest.Product.IsAvailable, response.Product.IsAvailable)

	st.repo.AssertExpectations(st.T())
}

func (st *ProductHandlerTestSuite) TestUpdateProduct_ProductNotFound() {
	// Create a mock product request for updating
	productRequest := &productspb.UpdateProductRequest{
		Product: &productspb.Product{
			ProductId:     st.testUUID.String(),
			Name:          "Updated Product Name",
			Sku:           "UPDATED-SKU",
			Category:      productspb.ProductCategory_CLOTHING,
			StockQuantity: 50,
			IsAvailable:   true,
		},
		UpdateMask: &field_mask.FieldMask{
			Paths: []string{"name", "sku", "category", "stock_quantity", "is_available"},
		},
	}

	// Set up expectations for the mock repository to return a product not found error
	st.repo.On("UpdateProductFields", mock.Anything, mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)

	// Call the UpdateProduct function
	response, err := st.handler.UpdateProduct(context.Background(), productRequest)

	// Assertions
	st.Require().NotNil(err)
	st.Require().Nil(response)
	st.Require().Equal(errNotFound, err)

	st.repo.AssertExpectations(st.T())
}

func (st *ProductHandlerTestSuite) TestUpdateProduct_InvalidUUIDError() {
	// Call the UpdateProduct function with an invalid UUID
	productRequest := &productspb.UpdateProductRequest{
		Product: &productspb.Product{
			ProductId:     "invalid-uuid",
			Name:          "Updated Product Name",
			Sku:           "UPDATED-SKU",
			Category:      productspb.ProductCategory_CLOTHING,
			StockQuantity: 50,
			IsAvailable:   true,
		},
		UpdateMask: &field_mask.FieldMask{
			Paths: []string{"name", "sku", "category", "stock_quantity", "is_available"},
		},
	}

	response, err := st.handler.UpdateProduct(context.Background(), productRequest)

	// Assertions
	st.Require().NotNil(err)
	st.Require().Nil(response)
	st.Require().Equal(errBadRequest, err)
}

func (st *ProductHandlerTestSuite) TestUpdateProduct_DBError() {
	// Create a mock product request for updating
	productRequest := &productspb.UpdateProductRequest{
		Product: &productspb.Product{
			ProductId:     st.testUUID.String(),
			Name:          "Updated Product Name",
			Sku:           "UPDATED-SKU",
			Category:      productspb.ProductCategory_CLOTHING,
			StockQuantity: 50,
			IsAvailable:   true,
		},
		UpdateMask: &field_mask.FieldMask{
			Paths: []string{"name", "sku", "category", "stock_quantity", "is_available"},
		},
	}

	// Set up expectations for the mock repository to return an error
	st.repo.On("UpdateProductFields", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("some error"))

	// Call the UpdateProduct function
	response, err := st.handler.UpdateProduct(context.Background(), productRequest)

	// Assertions
	st.Require().NotNil(err)
	st.Require().Nil(response)
	st.Require().Equal(errInternal, err)

	st.repo.AssertExpectations(st.T())
}

func (st *ProductHandlerTestSuite) TestUpdateProduct_NilRequest() {
	// Call the UpdateProduct function with a nil request
	response, err := st.handler.UpdateProduct(context.Background(), nil)

	// Assertions
	st.Require().NotNil(err)
	st.Require().Nil(response)
}
