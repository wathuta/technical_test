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
	grpcclients "github.com/wathuta/technical_test/orders/internal/grpc_clients"
	"github.com/wathuta/technical_test/orders/internal/mocks"
	"github.com/wathuta/technical_test/orders/internal/model"
	orderspb "github.com/wathuta/technical_test/protos_gen/orders"
	paymentpb "github.com/wathuta/technical_test/protos_gen/payment"
	"golang.org/x/exp/slog"
	"google.golang.org/genproto/protobuf/field_mask"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type OrderHandlerTestSuite struct {
	suite.Suite

	handler       *Handler
	repo          *mocks.Repository
	paymentclient *mocks.PaymentServiceClient

	testUUID  uuid.UUID
	testUUID1 uuid.UUID
	testUUID2 uuid.UUID
}

func TestOrderHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(OrderHandlerTestSuite))
}

func (st *OrderHandlerTestSuite) SetupSuite() {
	var programLevel = new(slog.LevelVar) // Info by default
	h := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: programLevel, AddSource: true})
	slog.SetDefault(slog.New(h))
	programLevel.Set(slog.LevelDebug)

	err := config.HasAllEnvVariables()
	st.Require().False(err)
}

func (st *OrderHandlerTestSuite) SetupTest() {
	repo := mocks.NewRepository(st.T())
	client := mocks.NewPaymentServiceClient(st.T())

	st.handler = New(repo, client)
	st.repo = repo
	st.paymentclient = client
	st.testUUID = uuid.New()
	st.testUUID1 = uuid.New()
	st.testUUID2 = uuid.New()

}

func (st *OrderHandlerTestSuite) TestCreateOrder_Success() {

	output := make(chan grpcclients.ServiceResult)
	go func() {
		output <- grpcclients.ServiceResult{
			Result: &paymentpb.CreatePaymentResponse{
				Payment: &paymentpb.Payment{
					PaymenId:   st.testUUID2.String(),
					OrderId:    st.testUUID.String(),
					CustomerId: st.testUUID.String(),
					Amount:     120.0,
					Currency:   "USD",
					Status:     paymentpb.PaymentStatus_COMPLETED,
					CreatedAt:  timestamppb.New(time.Now()),
				}}, Error: nil,
		}
	}()

	// Create a mock order request
	orderRequest := &orderspb.CreateOrderRequest{
		CustomerId:                st.testUUID.String(),
		ProductId:                 st.testUUID1.String(),
		ProductQuantity:           2,
		ShippingMethod:            "Express",
		PaymentMethod:             orderspb.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD,
		InvoiceNumber:             "INV12345",
		ShippingCost:              10.0,
		SpecialInstructions:       "Handle with care",
		ScheduledPickupDatetime:   timestamppb.New(time.Now().Add(24 * time.Hour)),
		ScheduledDeliveryDatetime: timestamppb.New(time.Now().Add(48 * time.Hour)),
		PickupAddress: &orderspb.Address{
			Street:     "123 Pickup St",
			City:       "Pickup City",
			State:      "Pickup State",
			PostalCode: "12345",
			Country:    "Pickup Country",
		},
		DeliveryAddress: &orderspb.Address{
			Street:     "123 Delivery St",
			City:       "Delivery City",
			State:      "Delivery State",
			PostalCode: "54321",
			Country:    "Delivery Country",
		},
	}

	// Set up expectations for the mock repository to get product and customer
	st.repo.On("GetProductById", mock.Anything, mock.Anything).Return(
		&model.Product{
			ProductID:     st.testUUID1.String(),
			Name:          "Sample Product",
			Sku:           "SKU123",
			Category:      "Electronics",
			StockQuantity: 100,
			ProductAttributes: model.ProductAttributes{
				Price: 100.0,
			},
			IsAvailable: true,
		},
		nil,
	)
	st.repo.On("GetCustomerById", mock.Anything, mock.Anything).Return(
		&model.Customer{
			CustomerID:  st.testUUID.String(),
			Name:        "John Doe",
			Email:       "johndoe@example.com",
			PhoneNumber: "+1234567890",
			Address:     "123 Main St",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		nil,
	)

	// Set up expectations for the mock repository to create an order
	st.repo.On("CreateOrder", mock.Anything, mock.Anything, mock.Anything).Return(
		&model.Order{
			OrderID:                   st.testUUID.String(),
			CustomerID:                orderRequest.CustomerId,
			ShippingMethod:            orderRequest.ShippingMethod,
			OrderStatus:               model.OrderStatus(orderspb.OrderStatus_ORDER_STATUS_PENDING.String()),
			TrackingNumber:            "1234567890",
			PaymentMethod:             model.PaymentMethod(orderspb.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD.String()),
			InvoiceNumber:             orderRequest.InvoiceNumber,
			ShippingCost:              orderRequest.ShippingCost,
			SpecialInstructions:       orderRequest.SpecialInstructions,
			ScheduledPickupDatetime:   orderRequest.ScheduledPickupDatetime.AsTime(),
			ScheduledDeliveryDatetime: orderRequest.ScheduledDeliveryDatetime.AsTime(),
			PickupAddress: model.Address{
				Street:     orderRequest.PickupAddress.Street,
				City:       orderRequest.PickupAddress.City,
				State:      orderRequest.PickupAddress.State,
				PostalCode: orderRequest.PickupAddress.PostalCode,
				Country:    orderRequest.PickupAddress.Country,
			},
			DeliveryAddress: model.Address{
				Street:     orderRequest.DeliveryAddress.Street,
				City:       orderRequest.DeliveryAddress.City,
				State:      orderRequest.DeliveryAddress.State,
				PostalCode: orderRequest.DeliveryAddress.PostalCode,
				Country:    orderRequest.DeliveryAddress.Country,
			},
			CreatedAt: time.Now(),
			DeletedAt: time.Time{},
		},
		&model.OrderDetails{
			OrderDetailsID: st.testUUID1.String(),
			OrderID:        st.testUUID.String(),
			ProductID:      orderRequest.ProductId,
			Quantity:       orderRequest.ProductQuantity,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Time{},
			DeletedAt:      time.Time{},
		},
		nil,
	)

	// Set up expectations for the mock payment service to create payment
	st.paymentclient.On("CreatePaymentRequest", mock.Anything, mock.Anything).Return(output)

	// Call the CreateOrder function
	response, err := st.handler.CreateOrder(context.Background(), orderRequest)

	// Assertions
	st.Require().NoError(err)
	st.Require().NotNil(response)
	st.Require().NotNil(response.Order)
	st.Require().NotNil(response.OrderDetails)
	st.Require().NotEmpty(response.Order.OrderId)
	st.Require().Equal(orderRequest.CustomerId, response.Order.CustomerId)
	st.Require().Equal(orderRequest.ShippingMethod, response.Order.ShippingMethod)
	st.Require().Equal(orderRequest.InvoiceNumber, response.Order.InvoiceNumber)
	st.Require().Equal(orderRequest.ShippingCost, response.Order.ShippingCost)
	st.Require().Equal(orderRequest.SpecialInstructions, response.Order.SpecialInstructions)
	st.Require().Equal(orderRequest.ScheduledPickupDatetime.AsTime(), response.Order.ScheduledPickupDatetime.AsTime())
	st.Require().Equal(orderRequest.ScheduledDeliveryDatetime.AsTime(), response.Order.ScheduledDeliveryDatetime.AsTime())

	st.repo.AssertExpectations(st.T())
	st.paymentclient.AssertExpectations(st.T())
}

func (st *OrderHandlerTestSuite) TestCreateOrder_InvalidRequest() {
	// Test case where the request is invalid (nil)
	response, err := st.handler.CreateOrder(context.Background(), nil)

	st.Require().Error(err)
	st.Require().Nil(response)
}

func (st *OrderHandlerTestSuite) TestCreateOrder_InvalidProductId() {
	// Test case where the request has an invalid ProductId
	orderRequest := &orderspb.CreateOrderRequest{
		CustomerId:      st.testUUID.String(),
		ProductId:       "", // Empty ProductId
		ProductQuantity: 2,
	}

	response, err := st.handler.CreateOrder(context.Background(), orderRequest)

	st.Require().Error(err)
	st.Require().Nil(response)
}

func (st *OrderHandlerTestSuite) TestCreateOrder_InvalidProductQuantity() {
	// Test case where the request has an invalid ProductQuantity
	orderRequest := &orderspb.CreateOrderRequest{
		CustomerId:      st.testUUID.String(),
		ProductId:       st.testUUID1.String(),
		ProductQuantity: -1, // Negative ProductQuantity
	}

	response, err := st.handler.CreateOrder(context.Background(), orderRequest)

	st.Require().Error(err)
	st.Require().Nil(response)
}

func (st *OrderHandlerTestSuite) TestCreateOrder_ProductNotFound() {
	// Test case where the requested product is not found
	orderRequest := &orderspb.CreateOrderRequest{
		CustomerId:                st.testUUID.String(),
		ProductId:                 st.testUUID1.String(),
		ProductQuantity:           2,
		ShippingMethod:            "Express",
		PaymentMethod:             orderspb.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD,
		InvoiceNumber:             "INV12345",
		ShippingCost:              10.0,
		SpecialInstructions:       "Handle with care",
		ScheduledPickupDatetime:   timestamppb.New(time.Now().Add(24 * time.Hour)),
		ScheduledDeliveryDatetime: timestamppb.New(time.Now().Add(48 * time.Hour)),
		PickupAddress: &orderspb.Address{
			Street:     "123 Pickup St",
			City:       "Pickup City",
			State:      "Pickup State",
			PostalCode: "12345",
			Country:    "Pickup Country",
		},
		DeliveryAddress: &orderspb.Address{
			Street:     "123 Delivery St",
			City:       "Delivery City",
			State:      "Delivery State",
			PostalCode: "54321",
			Country:    "Delivery Country",
		},
	}

	st.repo.On("GetProductById", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)

	response, err := st.handler.CreateOrder(context.Background(), orderRequest)

	st.Require().Error(err)
	st.Require().Nil(response)
}

func (st *OrderHandlerTestSuite) TestCreateOrder_CustomerNotFound() {
	// Test case where the customer is not found
	orderRequest := &orderspb.CreateOrderRequest{
		CustomerId:                st.testUUID.String(),
		ProductId:                 st.testUUID1.String(),
		ProductQuantity:           2,
		ShippingMethod:            "Express",
		PaymentMethod:             orderspb.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD,
		InvoiceNumber:             "INV12345",
		ShippingCost:              10.0,
		SpecialInstructions:       "Handle with care",
		ScheduledPickupDatetime:   timestamppb.New(time.Now().Add(24 * time.Hour)),
		ScheduledDeliveryDatetime: timestamppb.New(time.Now().Add(48 * time.Hour)),
		PickupAddress: &orderspb.Address{
			Street:     "123 Pickup St",
			City:       "Pickup City",
			State:      "Pickup State",
			PostalCode: "12345",
			Country:    "Pickup Country",
		},
		DeliveryAddress: &orderspb.Address{
			Street:     "123 Delivery St",
			City:       "Delivery City",
			State:      "Delivery State",
			PostalCode: "54321",
			Country:    "Delivery Country",
		},
	}

	st.repo.On("GetProductById", mock.Anything, mock.Anything).Return(
		&model.Product{
			ProductID: st.testUUID1.String(),
			Name:      "Sample Product",
		},
		nil,
	)
	st.repo.On("GetCustomerById", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)

	response, err := st.handler.CreateOrder(context.Background(), orderRequest)

	st.Require().Error(err)
	st.Require().Nil(response)
}

func (st *OrderHandlerTestSuite) TestCreateOrder_CreateOrderError() {
	// Test case where an error occurs while creating an order
	orderRequest := &orderspb.CreateOrderRequest{
		CustomerId:                st.testUUID.String(),
		ProductId:                 st.testUUID1.String(),
		ProductQuantity:           2,
		ShippingMethod:            "Express",
		PaymentMethod:             orderspb.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD,
		InvoiceNumber:             "INV12345",
		ShippingCost:              10.0,
		SpecialInstructions:       "Handle with care",
		ScheduledPickupDatetime:   timestamppb.New(time.Now().Add(24 * time.Hour)),
		ScheduledDeliveryDatetime: timestamppb.New(time.Now().Add(48 * time.Hour)),
		PickupAddress: &orderspb.Address{
			Street:     "123 Pickup St",
			City:       "Pickup City",
			State:      "Pickup State",
			PostalCode: "12345",
			Country:    "Pickup Country",
		},
		DeliveryAddress: &orderspb.Address{
			Street:     "123 Delivery St",
			City:       "Delivery City",
			State:      "Delivery State",
			PostalCode: "54321",
			Country:    "Delivery Country",
		},
	}

	st.repo.On("GetProductById", mock.Anything, mock.Anything).Return(
		&model.Product{
			ProductID: st.testUUID1.String(),
			Name:      "Sample Product",
		},
		nil,
	)
	st.repo.On("GetCustomerById", mock.Anything, mock.Anything).Return(
		&model.Customer{
			CustomerID: st.testUUID.String(),
			Name:       "John Doe",
		},
		nil,
	)
	st.repo.On("CreateOrder", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil, errors.New("order creation failed"))

	response, err := st.handler.CreateOrder(context.Background(), orderRequest)

	st.Require().Error(err)
	st.Require().Nil(response)
}

func (st *OrderHandlerTestSuite) TestCreateOrder_CreatePaymentError() {
	// Test case where an error occurs while creating a payment
	orderRequest := &orderspb.CreateOrderRequest{
		CustomerId:                st.testUUID.String(),
		ProductId:                 st.testUUID1.String(),
		ProductQuantity:           2,
		ShippingMethod:            "Express",
		PaymentMethod:             orderspb.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD,
		InvoiceNumber:             "INV12345",
		ShippingCost:              10.0,
		SpecialInstructions:       "Handle with care",
		ScheduledPickupDatetime:   timestamppb.New(time.Now().Add(24 * time.Hour)),
		ScheduledDeliveryDatetime: timestamppb.New(time.Now().Add(48 * time.Hour)),
		PickupAddress: &orderspb.Address{
			Street:     "123 Pickup St",
			City:       "Pickup City",
			State:      "Pickup State",
			PostalCode: "12345",
			Country:    "Pickup Country",
		},
		DeliveryAddress: &orderspb.Address{
			Street:     "123 Delivery St",
			City:       "Delivery City",
			State:      "Delivery State",
			PostalCode: "54321",
			Country:    "Delivery Country",
		},
	}

	output := make(chan grpcclients.ServiceResult)
	go func() {
		output <- grpcclients.ServiceResult{
			Result: nil, Error: errors.New("some error"),
		}
	}()

	st.repo.On("GetProductById", mock.Anything, mock.Anything).Return(
		&model.Product{
			ProductID: st.testUUID1.String(),
			Name:      "Sample Product",
			ProductAttributes: model.ProductAttributes{
				Price: 100.0,
			},
		},
		nil,
	)
	st.repo.On("GetCustomerById", mock.Anything, mock.Anything).Return(
		&model.Customer{
			CustomerID: st.testUUID.String(),
			Name:       "John Doe",
		},
		nil,
	)
	st.repo.On("CreateOrder", mock.Anything, mock.Anything, mock.Anything).Return(
		&model.Order{
			OrderID:                   st.testUUID.String(),
			CustomerID:                orderRequest.CustomerId,
			ShippingMethod:            orderRequest.ShippingMethod,
			OrderStatus:               model.OrderStatus(orderspb.OrderStatus_ORDER_STATUS_PENDING.String()),
			TrackingNumber:            "1234567890",
			PaymentMethod:             model.PaymentMethod(orderspb.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD.String()),
			InvoiceNumber:             orderRequest.InvoiceNumber,
			ShippingCost:              orderRequest.ShippingCost,
			SpecialInstructions:       orderRequest.SpecialInstructions,
			ScheduledPickupDatetime:   orderRequest.ScheduledPickupDatetime.AsTime(),
			ScheduledDeliveryDatetime: orderRequest.ScheduledDeliveryDatetime.AsTime(),
			PickupAddress: model.Address{
				Street:     orderRequest.PickupAddress.Street,
				City:       orderRequest.PickupAddress.City,
				State:      orderRequest.PickupAddress.State,
				PostalCode: orderRequest.PickupAddress.PostalCode,
				Country:    orderRequest.PickupAddress.Country,
			},
			DeliveryAddress: model.Address{
				Street:     orderRequest.DeliveryAddress.Street,
				City:       orderRequest.DeliveryAddress.City,
				State:      orderRequest.DeliveryAddress.State,
				PostalCode: orderRequest.DeliveryAddress.PostalCode,
				Country:    orderRequest.DeliveryAddress.Country,
			},
			CreatedAt: time.Now(),
			DeletedAt: time.Time{},
		},
		&model.OrderDetails{
			OrderDetailsID: st.testUUID1.String(),
			OrderID:        st.testUUID.String(),
			ProductID:      orderRequest.ProductId,
			Quantity:       orderRequest.ProductQuantity,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Time{},
			DeletedAt:      time.Time{},
		},
		nil,
	)
	st.paymentclient.On("CreatePaymentRequest", mock.Anything, mock.Anything).Return(output)

	response, err := st.handler.CreateOrder(context.Background(), orderRequest)

	st.Require().Error(err)
	st.Require().Nil(response)
}

func (st *OrderHandlerTestSuite) TestGetOrderById_Success() {
	// Create a mock order ID
	orderID := st.testUUID.String()

	// Set up expectations for the mock repository to return an order
	order := &model.Order{
		OrderID:                   orderID,
		CustomerID:                st.testUUID.String(),
		ShippingMethod:            "Express",
		OrderStatus:               model.OrderStatus(orderspb.OrderStatus_ORDER_STATUS_PENDING.String()),
		TrackingNumber:            "1234567890",
		PaymentMethod:             model.PaymentMethod(orderspb.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD.String()),
		InvoiceNumber:             "INV12345",
		ShippingCost:              10.0,
		SpecialInstructions:       "Handle with care",
		ScheduledPickupDatetime:   time.Now(),
		ScheduledDeliveryDatetime: time.Now().Add(24 * time.Hour),
		PickupAddress: model.Address{
			Street:     "123 Pickup St",
			City:       "Pickup City",
			State:      "Pickup State",
			PostalCode: "12345",
			Country:    "Pickup Country",
		},
		DeliveryAddress: model.Address{
			Street:     "123 Delivery St",
			City:       "Delivery City",
			State:      "Delivery State",
			PostalCode: "54321",
			Country:    "Delivery Country",
		},
		CreatedAt: time.Now(),
		DeletedAt: time.Time{},
	}

	st.repo.On("GetOrderById", mock.Anything, orderID).Return(order, nil)

	// Create a GetOrderRequest
	request := &orderspb.GetOrderRequest{
		OrderId: orderID,
	}

	// Call the GetOrderById function
	response, err := st.handler.GetOrderById(context.Background(), request)

	// Assertions
	st.Require().NoError(err)
	st.Require().NotNil(response)
	st.Require().NotNil(response.Order)
	st.Require().Equal(orderID, response.Order.OrderId)
	st.Require().Equal("Express", response.Order.ShippingMethod)
	st.Require().Equal("1234567890", response.Order.TrackingNumber)
	st.Require().Equal("INV12345", response.Order.InvoiceNumber)
	st.Require().Equal(10.0, response.Order.ShippingCost)
	st.Require().Equal("Handle with care", response.Order.SpecialInstructions)
	// Add more assertions as needed for other fields

	st.repo.AssertExpectations(st.T())
}

func (st *OrderHandlerTestSuite) TestGetOrderById_OrderNotFound() {
	// Create a mock order ID
	orderID := st.testUUID.String()

	// Set up expectations for the mock repository to return a "not found" error
	st.repo.On("GetOrderById", mock.Anything, orderID).Return(nil, sql.ErrNoRows)

	// Create a GetOrderRequest
	request := &orderspb.GetOrderRequest{
		OrderId: orderID,
	}

	// Call the GetOrderById function
	response, err := st.handler.GetOrderById(context.Background(), request)

	// Assertions
	st.Require().Error(err)
	st.Require().Nil(response)
	st.Require().Equal(errNotFound, err)

	st.repo.AssertExpectations(st.T())
}

func (st *OrderHandlerTestSuite) TestGetOrderById_InvalidOrderID() {
	// Create an invalid order ID (not a UUID)
	invalidOrderID := "invalid-id"

	// Create a GetOrderRequest with an invalid order ID
	request := &orderspb.GetOrderRequest{
		OrderId: invalidOrderID,
	}

	// Call the GetOrderById function
	response, err := st.handler.GetOrderById(context.Background(), request)

	// Assertions
	st.Require().Error(err)
	st.Require().Nil(response)
	st.Require().Equal(errBadRequest, err)
}

func (st *OrderHandlerTestSuite) TestGetOrderById_GetError() {
	// Create a mock order ID
	orderID := st.testUUID.String()

	// Set up expectations for the mock repository to return an error
	st.repo.On("GetOrderById", mock.Anything, orderID).Return(nil, errors.New("get error"))

	// Create a GetOrderRequest
	request := &orderspb.GetOrderRequest{
		OrderId: orderID,
	}

	// Call the GetOrderById function
	response, err := st.handler.GetOrderById(context.Background(), request)

	// Assertions
	st.Require().Error(err)
	st.Require().Nil(response)

	st.repo.AssertExpectations(st.T())
}

func (st *OrderHandlerTestSuite) TestDeleteOrder_Success() {
	// Create a mock order ID
	orderID := st.testUUID.String()

	// Set up expectations for the mock repository to return a successful deletion
	st.repo.On("DeleteOrder", mock.Anything, orderID).Return(&model.Order{}, nil)

	// Create a DeleteOrderRequest
	request := &orderspb.DeleteOrderRequest{
		OrderId: orderID,
	}

	// Call the DeleteOrder function
	response, err := st.handler.DeleteOrder(context.Background(), request)

	// Assertions
	st.Require().NoError(err)
	st.Require().True(response.Success)

	st.repo.AssertExpectations(st.T())
}

func (st *OrderHandlerTestSuite) TestDeleteOrder_InvalidOrderID() {
	// Create an invalid order ID (not a UUID)
	invalidOrderID := "invalid-id"

	// Create a DeleteOrderRequest with an invalid order ID
	request := &orderspb.DeleteOrderRequest{
		OrderId: invalidOrderID,
	}

	// Call the DeleteOrder function
	response, err := st.handler.DeleteOrder(context.Background(), request)

	// Assertions
	st.Require().Error(err)
	st.Require().False(response.Success)
	st.Require().Equal(errBadRequest, err)
}

func (st *OrderHandlerTestSuite) TestDeleteOrder_DeleteError() {
	// Create a mock order ID
	orderID := st.testUUID.String()

	// Set up expectations for the mock repository to return a deletion error
	st.repo.On("DeleteOrder", mock.Anything, orderID).Return(nil, errors.New("delete error"))

	// Create a DeleteOrderRequest
	request := &orderspb.DeleteOrderRequest{
		OrderId: orderID,
	}

	// Call the DeleteOrder function
	response, err := st.handler.DeleteOrder(context.Background(), request)

	// Assertions
	st.Require().Error(err)
	st.Require().False(response.Success)

	st.repo.AssertExpectations(st.T())
}

func (st *OrderHandlerTestSuite) TestUpdateOrder_Success() {
	// Create a mock order request for updating
	orderRequest := &orderspb.UpdateOrderRequest{
		Order: &orderspb.Order{
			OrderId:             st.testUUID.String(),
			ShippingMethod:      "Updated Shipping",
			OrderStatus:         orderspb.OrderStatus_ORDER_STATUS_SHIPPED,
			TrackingNumber:      "9876543210",
			PaymentMethod:       orderspb.PaymentMethod_PAYMENT_METHOD_MPESA,
			InvoiceNumber:       "UPDATED-INV56789",
			ShippingCost:        15.0,
			SpecialInstructions: "Fragile",
			// Add other fields as needed
		},
		UpdateMask: &field_mask.FieldMask{
			Paths: []string{"shipping_method", "order_status", "tracking_number", "payment_method", "invoice_number", "shipping_cost", "special_instructions"},
		},
	}

	// Set up expectations for the mock repository to update the order
	st.repo.On("UpdateOrder", mock.Anything, mock.Anything, mock.Anything).Return(
		&model.Order{
			OrderID:             st.testUUID.String(),
			ShippingMethod:      orderRequest.Order.ShippingMethod,
			OrderStatus:         model.OrderStatus(orderRequest.Order.OrderStatus.String()),
			TrackingNumber:      orderRequest.Order.TrackingNumber,
			PaymentMethod:       model.PaymentMethod(orderRequest.Order.PaymentMethod.String()),
			InvoiceNumber:       orderRequest.Order.InvoiceNumber,
			ShippingCost:        orderRequest.Order.ShippingCost,
			SpecialInstructions: orderRequest.Order.SpecialInstructions,
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
		},
		nil,
	)

	// Call the UpdateOrder function
	response, err := st.handler.UpdateOrder(context.Background(), orderRequest)

	// Assertions
	st.Require().NoError(err)
	st.Require().NotNil(response)
	st.Require().NotNil(response.Order)
	st.Require().NotEmpty(response.Order.OrderId)
	st.Require().Equal(orderRequest.Order.ShippingMethod, response.Order.ShippingMethod)
	st.Require().Equal(orderspb.OrderStatus_ORDER_STATUS_SHIPPED, response.Order.OrderStatus)
	st.Require().Equal(orderRequest.Order.TrackingNumber, response.Order.TrackingNumber)
	st.Require().Equal(orderRequest.Order.InvoiceNumber, response.Order.InvoiceNumber)
	st.Require().Equal(15.0, response.Order.ShippingCost)
	st.Require().Equal(orderRequest.Order.SpecialInstructions, response.Order.SpecialInstructions)

	st.repo.AssertExpectations(st.T())
}

func (st *OrderHandlerTestSuite) TestUpdateOrder_InvalidRequest() {
	// Create an invalid request with nil Order
	orderRequest := &orderspb.UpdateOrderRequest{
		Order:      nil,
		UpdateMask: nil, // Add a valid mask if needed
	}

	// Call the UpdateOrder function
	response, err := st.handler.UpdateOrder(context.Background(), orderRequest)

	// Assertions
	st.Require().Error(err)
	st.Require().Nil(response)
	st.Require().Equal(errResourceRequired, err)

	st.repo.AssertExpectations(st.T())
}

func (st *OrderHandlerTestSuite) TestUpdateOrder_UpdateError() {
	// Create a mock order request for updating
	orderRequest := &orderspb.UpdateOrderRequest{
		Order: &orderspb.Order{
			OrderId:             st.testUUID.String(),
			ShippingMethod:      "Updated Shipping",
			OrderStatus:         orderspb.OrderStatus_ORDER_STATUS_SHIPPED,
			TrackingNumber:      "9876543210",
			PaymentMethod:       orderspb.PaymentMethod_PAYMENT_METHOD_MPESA,
			InvoiceNumber:       "UPDATED-INV56789",
			ShippingCost:        15.0,
			SpecialInstructions: "Fragile",
			// Add other fields as needed
		},
		UpdateMask: &field_mask.FieldMask{
			Paths: []string{"shipping_method", "order_status", "tracking_number", "payment_method", "invoice_number", "shipping_cost", "special_instructions"},
		},
	}

	// Set up expectations for the mock repository to return an error during update
	st.repo.On("UpdateOrder", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("update error"))

	// Call the UpdateOrder function
	response, err := st.handler.UpdateOrder(context.Background(), orderRequest)

	// Assertions
	st.Require().Error(err)
	st.Require().Nil(response)

	st.repo.AssertExpectations(st.T())
}

func (st *OrderHandlerTestSuite) TestUpdateOrder_InvalidOrderID() {
	// Create an invalid order ID (not a UUID)
	invalidOrderID := "invalid-id"

	// Create an order request with an invalid order ID
	orderRequest := &orderspb.UpdateOrderRequest{
		Order: &orderspb.Order{
			OrderId: invalidOrderID,
			// Add other fields as needed
		},
		UpdateMask: &field_mask.FieldMask{
			Paths: []string{"shipping_method"}, // Add a valid field path
		},
	}

	// Call the UpdateOrder function
	response, err := st.handler.UpdateOrder(context.Background(), orderRequest)

	// Assertions
	st.Require().Error(err)
	st.Require().Nil(response)
	st.Require().Equal(errBadRequest, err)

	st.repo.AssertExpectations(st.T())
}

func (st *OrderHandlerTestSuite) TestUpdateOrder_OrderNotFound() {
	// Create a mock order request for updating
	orderRequest := &orderspb.UpdateOrderRequest{
		Order: &orderspb.Order{
			OrderId:             st.testUUID.String(),
			ShippingMethod:      "Updated Shipping",
			OrderStatus:         orderspb.OrderStatus_ORDER_STATUS_SHIPPED,
			TrackingNumber:      "9876543210",
			PaymentMethod:       orderspb.PaymentMethod_PAYMENT_METHOD_MPESA,
			InvoiceNumber:       "UPDATED-INV56789",
			ShippingCost:        15.0,
			SpecialInstructions: "Fragile",
			// Add other fields as needed
		},
		UpdateMask: &field_mask.FieldMask{
			Paths: []string{"shipping_method", "order_status", "tracking_number", "payment_method", "invoice_number", "shipping_cost", "special_instructions"},
		},
	}

	// Set up expectations for the mock repository to return an error indicating order not found
	st.repo.On("UpdateOrder", mock.Anything, mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)

	// Call the UpdateOrder function
	response, err := st.handler.UpdateOrder(context.Background(), orderRequest)

	// Assertions
	st.Require().Error(err)
	st.Require().Nil(response)
	st.Require().Equal(errNotFound, err)

	st.repo.AssertExpectations(st.T())
}

func (st *OrderHandlerTestSuite) TestUpdateOrder_EmptyUpdateFields() {
	// Create an order request with an empty update mask
	orderRequest := &orderspb.UpdateOrderRequest{
		Order: &orderspb.Order{
			OrderId: st.testUUID.String(),
			// Add other fields as needed
		},
		UpdateMask: &field_mask.FieldMask{
			Paths: nil, // Empty update fields
		},
	}
	st.repo.On("GetOrderById", mock.Anything, st.testUUID.String()).Return(&model.Order{
		OrderID: st.testUUID.String(),
	}, nil)

	// Call the UpdateOrder function
	response, err := st.handler.UpdateOrder(context.Background(), orderRequest)

	// Assertions
	st.Require().NoError(err)
	st.Require().NotNil(response)
	st.Require().NotNil(response.Order)
	st.Require().NotEmpty(response.Order.OrderId)
	// Verify that the order was not updated

	st.repo.AssertExpectations(st.T())
}

func (st *OrderHandlerTestSuite) TestListOrdersByCustomerId_Success() {
	// Create a mock customer ID
	customerID := st.testUUID.String()

	// Create a mock list of orders
	mockOrders := []model.Order{
		{
			OrderID:    st.testUUID.String(),
			CustomerID: customerID,
			// Add other order details as needed
		},
		{
			OrderID:    st.testUUID.String(),
			CustomerID: customerID,
			// Add other order details as needed
		},
	}

	st.repo.On("GetOrdersByCustomerId", mock.Anything, customerID, mock.Anything, mock.Anything).Return(mockOrders, nil)

	// Create a ListOrdersByCustomerIdRequest
	request := &orderspb.ListOrdersByCustomerIdRequest{
		CustomerId: customerID,
		PageSize:   10,
		PageToken:  0,
	}

	// Call the ListOrdersByCustomerId function
	response, err := st.handler.ListOrdersByCustomerId(context.Background(), request)

	// Assertions
	st.Require().NoError(err)
	st.Require().NotNil(response)
	st.Require().NotNil(response.Orders)
	st.Require().Len(response.Orders, len(mockOrders))

	st.repo.AssertExpectations(st.T())
}

func (st *OrderHandlerTestSuite) TestListOrdersByCustomerId_InvalidCustomerId() {
	// Create an invalid customer ID
	invalidCustomerID := "invalid-id"

	// Create a ListOrdersByCustomerIdRequest with an invalid customer ID
	request := &orderspb.ListOrdersByCustomerIdRequest{
		CustomerId: invalidCustomerID,
		PageSize:   10,
		PageToken:  0,
	}

	// Call the ListOrdersByCustomerId function
	response, err := st.handler.ListOrdersByCustomerId(context.Background(), request)

	// Assertions
	st.Require().Error(err)
	st.Require().Nil(response)
	st.Require().Equal(errBadRequest, err)
}

func (st *OrderHandlerTestSuite) TestListOrdersByCustomerId_NotFound() {

	request := &orderspb.ListOrdersByCustomerIdRequest{
		CustomerId: st.testUUID.String(),
		PageSize:   10,
		PageToken:  0,
	}

	st.repo.On("GetOrdersByCustomerId", mock.Anything, request.CustomerId, mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)

	// Call the ListOrdersByCustomerId function
	response, err := st.handler.ListOrdersByCustomerId(context.Background(), request)

	// Assertions
	st.Require().Error(err)
	st.Require().Nil(response)
	st.Require().Equal(errNotFound, err)
}

func (st *OrderHandlerTestSuite) TestListOrdersByCustomerId_InternalError() {

	request := &orderspb.ListOrdersByCustomerIdRequest{
		CustomerId: st.testUUID.String(),
		PageSize:   10,
		PageToken:  0,
	}

	st.repo.On("GetOrdersByCustomerId", mock.Anything, request.CustomerId, mock.Anything, mock.Anything).Return(nil, errors.New("some error"))

	// Call the ListOrdersByCustomerId function
	response, err := st.handler.ListOrdersByCustomerId(context.Background(), request)

	// Assertions
	st.Require().Error(err)
	st.Require().Nil(response)
	st.Require().Equal(errInternal, err)
}

// Happy Path Tests:
func (st *OrderHandlerTestSuite) TestListOrderDetailsByOrderId_Success() {
	// Create a mock order ID
	orderID := st.testUUID.String()

	// Create a mock list of order details
	mockOrderDetails := []model.OrderDetails{
		{
			OrderDetailsID: orderID + "-1",
			OrderID:        orderID,
			// Add other order details as needed
		},
		{
			OrderDetailsID: orderID + "-2",
			OrderID:        orderID,
			// Add other order details as needed
		},
	}

	st.repo.On("GetOrderDetailsByOrderId", mock.Anything, orderID, mock.Anything, mock.Anything).Return(mockOrderDetails, nil)

	// Create a ListOrderDetailsByOrderIdRequest
	request := &orderspb.ListOrderDetailsByOrderIdRequest{
		OrderId:  orderID,
		PageSize: 10,
	}

	// Call the ListOrderDetailsByOrderId function
	response, err := st.handler.ListOrderDetailsByOrderId(context.Background(), request)

	// Assertions
	st.Require().NoError(err)
	st.Require().NotNil(response)
	st.Require().NotNil(response.OrderDetails)
	st.Require().Len(response.OrderDetails, len(mockOrderDetails))

	st.repo.AssertExpectations(st.T())
}

// Sad Path Tests:
func (st *OrderHandlerTestSuite) TestListOrderDetailsByOrderId_InvalidOrderId() {
	// Create an invalid order ID
	invalidOrderID := "invalid-id"

	// Create a ListOrderDetailsByOrderIdRequest with an invalid order ID
	request := &orderspb.ListOrderDetailsByOrderIdRequest{
		OrderId:  invalidOrderID,
		PageSize: 10,
	}

	// Call the ListOrderDetailsByOrderId function
	response, err := st.handler.ListOrderDetailsByOrderId(context.Background(), request)

	// Assertions
	st.Require().Error(err)
	st.Require().Nil(response)
	st.Require().Equal(errBadRequest, err)
}

func (st *OrderHandlerTestSuite) TestListOrderDetailsByOrderId_OrderDetailsNotFound() {
	// Create a mock order ID
	orderID := st.testUUID.String()

	// Set up expectations for the mock repository to return a "order details not found" error
	st.repo.On("GetOrderDetailsByOrderId", mock.Anything, orderID, mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)

	// Create a ListOrderDetailsByOrderIdRequest
	request := &orderspb.ListOrderDetailsByOrderIdRequest{
		OrderId:  orderID,
		PageSize: 10,
	}

	// Call the ListOrderDetailsByOrderId function
	response, err := st.handler.ListOrderDetailsByOrderId(context.Background(), request)

	// Assertions
	st.Require().Error(err)
	st.Require().Nil(response)
	st.Require().Equal(errNotFound, err)

	st.repo.AssertExpectations(st.T())
}

func (st *OrderHandlerTestSuite) TestListOrderDetailsByOrderId_InternalError() {
	// Create a mock order ID
	orderID := st.testUUID.String()

	// Set up expectations for the mock repository to return an internal error
	st.repo.On("GetOrderDetailsByOrderId", mock.Anything, orderID, mock.Anything, mock.Anything).Return(nil, errors.New("some error"))

	// Create a ListOrderDetailsByOrderIdRequest
	request := &orderspb.ListOrderDetailsByOrderIdRequest{
		OrderId:  orderID,
		PageSize: 10,
	}

	// Call the ListOrderDetailsByOrderId function
	response, err := st.handler.ListOrderDetailsByOrderId(context.Background(), request)

	// Assertions
	st.Require().Error(err)
	st.Require().Nil(response)
	st.Require().Equal(errInternal, err)

	st.repo.AssertExpectations(st.T())
}

// Happy Path Tests:
func (st *OrderHandlerTestSuite) TestListOrdersByProductId_Success() {
	// Create a mock product ID
	productID := st.testUUID.String()

	// Create a mock list of order details
	mockOrderDetails := []model.OrderDetails{
		{
			OrderDetailsID: st.testUUID1.String(),
			OrderID:        st.testUUID.String(),
			ProductID:      productID,
			// Add other order detail fields as needed
		},
		{
			OrderDetailsID: st.testUUID1.String(),
			OrderID:        st.testUUID.String(),
			ProductID:      productID,
			// Add other order detail fields as needed
		},
	}

	// Set up expectations for the mock repository to return order details
	st.repo.On("GetOrderDetailsByProductId", mock.Anything, productID, mock.Anything, mock.Anything).Return(mockOrderDetails, nil)

	// Set up expectations for the mock repository to return orders for each order detail
	st.repo.On("GetOrderById", mock.Anything, mock.Anything).Return(&model.Order{
		OrderID: st.testUUID.String(),
		// Add other order fields as needed
	}, nil).Times(len(mockOrderDetails))

	// Create a ListOrdersByProductIdRequest
	request := &orderspb.ListOrdersByProductIdRequest{
		ProductId: productID,
		PageSize:  10,
	}

	// Call the ListOrdersByProductId function
	response, err := st.handler.ListOrdersByProductId(context.Background(), request)

	// Assertions
	st.Require().NoError(err)
	st.Require().NotNil(response)
	st.Require().NotNil(response.Orders)
	st.Require().NotNil(response.OrderDetails)
	st.Require().Len(response.Orders, len(mockOrderDetails))
	st.Require().Len(response.OrderDetails, len(mockOrderDetails))

	st.repo.AssertExpectations(st.T())
}

// Sad Path Tests:
func (st *OrderHandlerTestSuite) TestListOrdersByProductId_InvalidProductId() {
	// Create an invalid product ID
	invalidProductID := "invalid-id"

	// Create a ListOrdersByProductIdRequest with an invalid product ID
	request := &orderspb.ListOrdersByProductIdRequest{
		ProductId: invalidProductID,
		PageSize:  10,
	}

	// Call the ListOrdersByProductId function
	response, err := st.handler.ListOrdersByProductId(context.Background(), request)

	// Assertions
	st.Require().Error(err)
	st.Require().Nil(response)
	st.Require().Equal(errBadRequest, err)
}

func (st *OrderHandlerTestSuite) TestListOrdersByProductId_OrderDetailsNotFound() {
	// Create a mock product ID
	productID := st.testUUID.String()

	// Set up expectations for the mock repository to return "order details not found" error
	st.repo.On("GetOrderDetailsByProductId", mock.Anything, productID, mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)

	// Create a ListOrdersByProductIdRequest
	request := &orderspb.ListOrdersByProductIdRequest{
		ProductId: productID,
		PageSize:  10,
	}

	// Call the ListOrdersByProductId function
	response, err := st.handler.ListOrdersByProductId(context.Background(), request)

	// Assertions
	st.Require().Error(err)
	st.Require().Nil(response)
	st.Require().Equal(errNotFound, err)

	st.repo.AssertExpectations(st.T())
}

func (st *OrderHandlerTestSuite) TestListOrdersByProductId_OrderNotFound() {
	// Create a mock product ID
	productID := st.testUUID.String()

	// Set up expectations for the mock repository to return orders not found error for each order detail
	mockOrderDetails := []model.OrderDetails{
		{
			OrderDetailsID: st.testUUID1.String(),
			OrderID:        st.testUUID.String(),
			ProductID:      productID,
		},
		{
			OrderDetailsID: st.testUUID.String(),
			OrderID:        st.testUUID.String(),
			ProductID:      productID,
		},
	}

	st.repo.On("GetOrderDetailsByProductId", mock.Anything, productID, mock.Anything, mock.Anything).Return(mockOrderDetails, nil)

	st.repo.On("GetOrderById", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)

	// Create a ListOrdersByProductIdRequest
	request := &orderspb.ListOrdersByProductIdRequest{
		ProductId: productID,
		PageSize:  10,
	}

	// Call the ListOrdersByProductId function
	response, err := st.handler.ListOrdersByProductId(context.Background(), request)

	// Assertions
	st.Require().Error(err)
	st.Require().Nil(response)
	st.Require().Equal(errNotFound, err)

	st.repo.AssertExpectations(st.T())
}

func (st *OrderHandlerTestSuite) TestListOrdersByProductId_InternalError() {
	// Create a mock product ID
	productID := st.testUUID.String()

	// Set up expectations for the mock repository to return an internal error
	st.repo.On("GetOrderDetailsByProductId", mock.Anything, productID, mock.Anything, mock.Anything).Return(nil, errors.New("some error"))

	// Create a ListOrdersByProductIdRequest
	request := &orderspb.ListOrdersByProductIdRequest{
		ProductId: productID,
		PageSize:  10,
	}

	// Call the ListOrdersByProductId function
	response, err := st.handler.ListOrdersByProductId(context.Background(), request)

	// Assertions
	st.Require().Error(err)
	st.Require().Nil(response)
	st.Require().Equal(errInternal, err)

	st.repo.AssertExpectations(st.T())
}

// Happy Path Test:
func (st *OrderHandlerTestSuite) TestGetOrderDetailsById_Success() {
	// Create a mock order details ID
	orderDetailsID := st.testUUID.String()

	// Create a mock order details
	mockOrderDetails := &model.OrderDetails{
		OrderDetailsID: orderDetailsID,
		OrderID:        st.testUUID.String(),
		ProductID:      st.testUUID.String(),
		// Add other order details fields as needed
	}

	// Set up expectations for the mock repository to return order details
	st.repo.On("GetOrderDetailsById", mock.Anything, orderDetailsID).Return(mockOrderDetails, nil)

	// Create a GetOrderDetailByIdRequest
	request := &orderspb.GetOrderDetailByIdRequest{
		OrderDetailsId: orderDetailsID,
	}

	// Call the GetOrderDetailsById function
	response, err := st.handler.GetOrderDetailsById(context.Background(), request)

	// Assertions
	st.Require().NoError(err)
	st.Require().NotNil(response)
	st.Require().NotNil(response.OrderDetails)
	st.Require().Equal(orderDetailsID, response.OrderDetails.OrderDetailsId)

	st.repo.AssertExpectations(st.T())
}

// Sad Path Tests:
func (st *OrderHandlerTestSuite) TestGetOrderDetailsById_InvalidOrderDetailsId() {
	// Create an invalid order details ID
	invalidOrderDetailsID := "invalid-id"

	// Create a GetOrderDetailByIdRequest with an invalid order details ID
	request := &orderspb.GetOrderDetailByIdRequest{
		OrderDetailsId: invalidOrderDetailsID,
	}

	// Call the GetOrderDetailsById function
	response, err := st.handler.GetOrderDetailsById(context.Background(), request)

	// Assertions
	st.Require().Error(err)
	st.Require().Nil(response)
	st.Require().Equal(errBadRequest, err)
}

func (st *OrderHandlerTestSuite) TestGetOrderDetailsById_OrderDetailsNotFound() {
	// Create a mock order details ID
	orderDetailsID := st.testUUID.String()

	// Set up expectations for the mock repository to return "order details not found" error
	st.repo.On("GetOrderDetailsById", mock.Anything, orderDetailsID).Return(nil, sql.ErrNoRows)

	// Create a GetOrderDetailByIdRequest
	request := &orderspb.GetOrderDetailByIdRequest{
		OrderDetailsId: orderDetailsID,
	}

	// Call the GetOrderDetailsById function
	response, err := st.handler.GetOrderDetailsById(context.Background(), request)

	// Assertions
	st.Require().Error(err)
	st.Require().Nil(response)
	st.Require().Equal(errNotFound, err)

	st.repo.AssertExpectations(st.T())
}

func (st *OrderHandlerTestSuite) TestGetOrderDetailsById_InternalError() {
	// Create a mock order details ID
	orderDetailsID := st.testUUID.String()

	// Set up expectations for the mock repository to return an internal error
	st.repo.On("GetOrderDetailsById", mock.Anything, orderDetailsID).Return(nil, errors.New("some error"))

	// Create a GetOrderDetailByIdRequest
	request := &orderspb.GetOrderDetailByIdRequest{
		OrderDetailsId: orderDetailsID,
	}

	// Call the GetOrderDetailsById function
	response, err := st.handler.GetOrderDetailsById(context.Background(), request)

	// Assertions
	st.Require().Error(err)
	st.Require().Nil(response)
	st.Require().Equal(errInternal, err)

	st.repo.AssertExpectations(st.T())
}
