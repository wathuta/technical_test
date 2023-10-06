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
	customersPb "github.com/wathuta/technical_test/protos_gen/customers"
	"golang.org/x/exp/slog"
	"google.golang.org/genproto/protobuf/field_mask"
)

type CustomerHandlerTestSuite struct {
	suite.Suite

	handler *Handler
	repo    *mocks.Repository

	testUUID  uuid.UUID
	testUUID1 uuid.UUID
}

func TestCustomerHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(CustomerHandlerTestSuite))
}

func (st *CustomerHandlerTestSuite) SetupSuite() {
	var programLevel = new(slog.LevelVar) // Info by default
	h := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: programLevel, AddSource: true})
	slog.SetDefault(slog.New(h))
	programLevel.Set(slog.LevelDebug)

	err := config.HasAllEnvVariables()
	st.Require().False(err)
}

func (st *CustomerHandlerTestSuite) SetupTest() {
	repo := mocks.NewRepository(st.T())
	client := mocks.NewPaymentServiceClient(st.T())

	st.handler = New(repo, client)
	st.repo = repo
	st.testUUID = uuid.New()
	st.testUUID1 = uuid.New()
}

func (st *CustomerHandlerTestSuite) TestCreateCustomer_Success() {
	// Create a mock customer request
	customerRequest := &customersPb.CreateCustomerRequest{
		Customer: &customersPb.Customer{
			Name:        "John Doe",
			Email:       "john@example.com",
			PhoneNumber: "+1234567890",
			Address:     "123 Main St",
		},
	}

	// Set up expectations for the mock repository
	st.repo.On("CreateCustomer", mock.Anything, mock.Anything).Return(&model.Customer{
		CustomerID:  st.testUUID.String(),
		Name:        customerRequest.Customer.Name,
		Email:       customerRequest.Customer.Email,
		PhoneNumber: customerRequest.Customer.PhoneNumber,
		Address:     customerRequest.Customer.Address,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil)

	// Call the CreateCustomer function
	response, err := st.handler.CreateCustomer(context.Background(), customerRequest)

	// Assertions
	st.Require().NoError(err)
	st.Require().NotNil(response)
	st.Require().NotNil(response.Customer)
	st.Require().NotEmpty(response.Customer.CustomerId)
	st.Require().Equal(customerRequest.Customer.Name, response.Customer.Name)
	st.Require().Equal(customerRequest.Customer.Email, response.Customer.Email)
	st.Require().Equal(customerRequest.Customer.PhoneNumber, response.Customer.PhoneNumber)
	st.Require().Equal(customerRequest.Customer.Address, response.Customer.Address)

	st.repo.AssertExpectations(st.T())
}

func (st *CustomerHandlerTestSuite) TestCreateCustomer_InvalidRequest() {
	// Create a mock customer request with nil Customer
	customerRequest := &customersPb.CreateCustomerRequest{
		Customer: nil,
	}

	// Call the CreateCustomer function
	response, err := st.handler.CreateCustomer(context.Background(), customerRequest)

	// Assertions
	st.Require().Error(err)
	st.Require().Nil(response)

	// Assert expectations for the mock repository (no calls expected)
	st.repo.AssertExpectations(st.T())
}

func (st *CustomerHandlerTestSuite) TestCreateCustomer_ValidationError() {
	// Create a mock customer request with an invalid email
	customerRequest := &customersPb.CreateCustomerRequest{
		Customer: &customersPb.Customer{
			Name:        "John Doe",
			Email:       "invalid-email",
			PhoneNumber: "+1234567890",
			Address:     "123 Main St",
		},
	}

	// Call the CreateCustomer function
	response, err := st.handler.CreateCustomer(context.Background(), customerRequest)

	// Assertions
	st.Require().Error(err)
	st.Require().Nil(response)

	// Assert expectations for the mock repository (no calls expected)
	st.repo.AssertExpectations(st.T())
}

func (st *CustomerHandlerTestSuite) TestCreateCustomer_CreateCustomerError() {
	// Create a mock customer request
	customerRequest := &customersPb.CreateCustomerRequest{
		Customer: &customersPb.Customer{
			Name:        "John Doe",
			Email:       "john@example.com",
			PhoneNumber: "+1234567890",
			Address:     "123 Main St",
		},
	}

	// Set up expectations for the mock repository to return an error
	st.repo.On("CreateCustomer", mock.Anything, mock.Anything).Return(nil, errors.New("some error"))

	// Call the CreateCustomer function
	response, err := st.handler.CreateCustomer(context.Background(), customerRequest)

	// Assertions
	st.Require().Error(err)
	st.Require().Nil(response)

	// Assert expectations for the mock repository
	st.repo.AssertExpectations(st.T())
}

func (st *CustomerHandlerTestSuite) TestGetCustomerById_InvalidUUIDError() {
	// Call the GetCustomerById function with an invalid UUID
	resp, err := st.handler.GetCustomerById(context.Background(), &customersPb.GetCustomerByIdRequest{
		CustomerId: "invalid-uuid",
	})

	// Assertions
	st.Require().Nil(resp)
	st.Require().NotNil(err)
	st.Require().Equal(errBadRequest, err)
}

func (st *CustomerHandlerTestSuite) TestGetCustomerById_DBError() {
	// Set up expectations for the mock repository to return an error
	st.repo.On("GetCustomerById", mock.Anything, mock.Anything).Return(nil, errors.New("some error"))

	// Call the GetCustomerById function
	resp, err := st.handler.GetCustomerById(context.Background(), &customersPb.GetCustomerByIdRequest{
		CustomerId: st.testUUID.String(),
	})

	// Assertions
	st.Require().Nil(resp)
	st.Require().NotNil(err)
	st.Require().Equal(errInternal, err)
}

func (st *CustomerHandlerTestSuite) TestGetCustomerById_NotFoundError() {
	// Set up expectations for the mock repository to return sql.ErrNoRows (not found)
	st.repo.On("GetCustomerById", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)

	// Call the GetCustomerById function
	resp, err := st.handler.GetCustomerById(context.Background(), &customersPb.GetCustomerByIdRequest{
		CustomerId: st.testUUID.String(),
	})

	// Assertions
	st.Require().Nil(resp)
	st.Require().NotNil(err)
	st.Require().Equal(errNotFound, err)
}

func (st *CustomerHandlerTestSuite) TestGetCustomerById_NilRequest() {
	// Call the GetCustomerById function with a nil request
	resp, err := st.handler.GetCustomerById(context.Background(), nil)

	// Assertions
	st.Require().Nil(resp)
	st.Require().NotNil(err)
}

func (st *CustomerHandlerTestSuite) TestGetCustomerById_Success() {
	// Set up expectations for the mock repository to return a customer
	st.repo.On("GetCustomerById", mock.Anything, mock.Anything).Return(&model.Customer{
		CustomerID:  st.testUUID.String(),
		Name:        "John Doe",
		Email:       "john@example.com",
		PhoneNumber: "+1234567890",
		Address:     "123 Main St",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil)

	// Call the GetCustomerById function
	resp, err := st.handler.GetCustomerById(context.Background(), &customersPb.GetCustomerByIdRequest{
		CustomerId: st.testUUID.String(),
	})

	// Assertions
	st.Require().NoError(err)
	st.Require().NotNil(resp)
	st.Require().NotNil(resp.Customer)
	st.Require().NotEmpty(resp.Customer.CustomerId)
	st.Require().Equal("John Doe", resp.Customer.Name)
	st.Require().Equal("john@example.com", resp.Customer.Email)
	st.Require().Equal("+1234567890", resp.Customer.PhoneNumber)
	st.Require().Equal("123 Main St", resp.Customer.Address)

	st.repo.AssertExpectations(st.T())
}

func (st *CustomerHandlerTestSuite) TestDeleteCustomer_Success() {
	// Set up expectations for the mock repository to delete a customer
	st.repo.On("DeleteCustomer", mock.Anything, mock.Anything).Return(&model.Customer{
		CustomerID: st.testUUID.String(),
	}, nil)

	// Call the DeleteCustomer function
	resp, err := st.handler.DeleteCustomer(context.Background(), &customersPb.DeleteCustomerRequest{
		CustomerId: st.testUUID.String(),
	})

	// Assertions
	st.Require().NoError(err)
	st.Require().NotNil(resp)
	st.Require().True(resp.Success)

	st.repo.AssertExpectations(st.T())
}

func (st *CustomerHandlerTestSuite) TestDeleteCustomer_CustomerNotFound() {
	// Set up expectations for the mock repository to delete a customer that is not found
	st.repo.On("DeleteCustomer", mock.Anything, mock.Anything).Return(nil, nil)

	// Call the DeleteCustomer function
	resp, err := st.handler.DeleteCustomer(context.Background(), &customersPb.DeleteCustomerRequest{
		CustomerId: st.testUUID.String(),
	})

	// Assertions
	st.Require().NotNil(err)
	st.Require().NotNil(resp)
	st.Require().False(resp.Success)

	st.repo.AssertExpectations(st.T())
}

func (st *CustomerHandlerTestSuite) TestDeleteCustomer_InvalidUUIDError() {
	// Call the DeleteCustomer function with an invalid UUID
	resp, err := st.handler.DeleteCustomer(context.Background(), &customersPb.DeleteCustomerRequest{
		CustomerId: "invalid-uuid",
	})

	// Assertions
	st.Require().False(resp.Success)
	st.Require().NotNil(err)
	st.Require().Equal(errBadRequest, err)
}

func (st *CustomerHandlerTestSuite) TestDeleteCustomer_DBError() {
	// Set up expectations for the mock repository to return an error when deleting a customer
	st.repo.On("DeleteCustomer", mock.Anything, mock.Anything).Return(nil, errors.New("some error"))

	// Call the DeleteCustomer function
	resp, err := st.handler.DeleteCustomer(context.Background(), &customersPb.DeleteCustomerRequest{
		CustomerId: st.testUUID.String(),
	})

	// Assertions
	st.Require().False(resp.Success)
	st.Require().NotNil(err)
	st.Require().Equal(errInternal, err)

	st.repo.AssertExpectations(st.T())
}

func (st *CustomerHandlerTestSuite) TestDeleteCustomer_NilRequest() {
	// Call the DeleteCustomer function with a nil request
	resp, err := st.handler.DeleteCustomer(context.Background(), nil)

	// Assertions
	st.Require().False(resp.Success)
	st.Require().NotNil(err)
}

func (st *CustomerHandlerTestSuite) TestUpdateCustomer_Success() {
	// Create a mock customer request for updating
	customerRequest := &customersPb.UpdateCustomerRequest{
		Customer: &customersPb.Customer{
			CustomerId:  st.testUUID.String(),
			Name:        "Updated Name",
			Email:       "updated@example.com",
			PhoneNumber: "+9876543210",
			Address:     "456 New St",
		},
		UpdateMask: &field_mask.FieldMask{
			Paths: []string{"name", "email", "phone_number", "address"},
		},
	}

	// Set up expectations for the mock repository to update the customer
	st.repo.On("UpdateCustomerFields", mock.Anything, mock.Anything, mock.Anything).Return(
		&model.Customer{
			CustomerID:  st.testUUID.String(),
			Name:        customerRequest.Customer.Name,
			Email:       customerRequest.Customer.Email,
			PhoneNumber: customerRequest.Customer.PhoneNumber,
			Address:     customerRequest.Customer.Address,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		nil,
	)

	// Call the UpdateCustomer function
	response, err := st.handler.UpdateCustomer(context.Background(), customerRequest)

	// Assertions
	st.Require().NoError(err)
	st.Require().NotNil(response)
	st.Require().NotNil(response.Customer)
	st.Require().NotEmpty(response.Customer.CustomerId)
	st.Require().Equal(customerRequest.Customer.Name, response.Customer.Name)
	st.Require().Equal(customerRequest.Customer.Email, response.Customer.Email)
	st.Require().Equal(customerRequest.Customer.PhoneNumber, response.Customer.PhoneNumber)
	st.Require().Equal(customerRequest.Customer.Address, response.Customer.Address)

	st.repo.AssertExpectations(st.T())
}

func (st *CustomerHandlerTestSuite) TestUpdateCustomer_EmptyUpdateMask() {
	// Create a mock customer request with an empty UpdateMask
	customerRequest := &customersPb.UpdateCustomerRequest{
		Customer: &customersPb.Customer{
			CustomerId:  st.testUUID.String(),
			Name:        "John Doe",
			Email:       "john@example.com",
			PhoneNumber: "+1234567890",
			Address:     "123 Main St",
		},
		UpdateMask: &field_mask.FieldMask{},
	}
	st.repo.On("GetCustomerById", mock.Anything, mock.Anything).Return(&model.Customer{
		CustomerID:  st.testUUID.String(),
		Name:        "John Doe",
		Email:       "john@example.com",
		PhoneNumber: "+1234567890",
		Address:     "123 Main St",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil)

	// Call the UpdateCustomer function
	response, err := st.handler.UpdateCustomer(context.Background(), customerRequest)

	// Assertions
	st.Require().NoError(err)
	st.Require().NotNil(response)
	st.Require().NotNil(response.Customer)
	st.Require().NotEmpty(response.Customer.CustomerId)
	st.Require().Equal(customerRequest.Customer.Name, response.Customer.Name)
	st.Require().Equal(customerRequest.Customer.Email, response.Customer.Email)
	st.Require().Equal(customerRequest.Customer.PhoneNumber, response.Customer.PhoneNumber)
	st.Require().Equal(customerRequest.Customer.Address, response.Customer.Address)

	st.repo.AssertExpectations(st.T())
}

func (st *CustomerHandlerTestSuite) TestUpdateCustomer_UpdateError() {
	// Create a mock customer request for updating
	customerRequest := &customersPb.UpdateCustomerRequest{
		Customer: &customersPb.Customer{
			CustomerId:  st.testUUID.String(),
			Name:        "Updated Name",
			Email:       "updated@example.com",
			PhoneNumber: "+9876543210",
			Address:     "456 New St",
		},
		UpdateMask: &field_mask.FieldMask{
			Paths: []string{"name", "email", "phone_number", "address"},
		},
	}

	// Set up expectations for the mock repository to return an update error
	st.repo.On("UpdateCustomerFields", mock.Anything, mock.Anything, mock.Anything).Return(
		nil, errors.New("update error"),
	)

	// Call the UpdateCustomer function
	response, err := st.handler.UpdateCustomer(context.Background(), customerRequest)

	// Assertions
	st.Require().Error(err)
	st.Require().Nil(response)

	st.repo.AssertExpectations(st.T())
}


