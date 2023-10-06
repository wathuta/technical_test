package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	grpcclients "github.com/wathuta/technical_test/payment/internal/grpc_clients"
	"github.com/wathuta/technical_test/payment/internal/model"
	orderspb "github.com/wathuta/technical_test/protos_gen/orders"
)

func (st *PaymentHandlerTestSuite) TestCallbackHandler_Success() {
	output := make(chan grpcclients.ServiceResult)
	go func() {
		output <- grpcclients.ServiceResult{Error: nil, Result: &orderspb.Order{
			OrderId:     st.testUUID.String(),
			OrderStatus: orderspb.OrderStatus_ORDER_STATUS_PROCESSING,
		}}
	}()
	// Mock repository to return a payment record
	st.repo.On("GetPaymentByMerchantRequestId", mock.Anything, mock.Anything).Return(
		&model.Payment{
			PaymentID: st.testUUID1.String(),
			OrderID:   st.testUUID.String(),
		}, nil,
	)
	st.orderclient.On("UpdateOrderDetails", st.testUUID.String(), orderspb.OrderStatus_ORDER_STATUS_PROCESSING).Return(
		output,
	)
	st.repo.On("UpdatePaymentStatus", mock.Anything, mock.Anything, mock.Anything).Return(
		&model.Payment{
			PaymentID: st.testUUID1.String(),
			OrderID:   st.testUUID.String(),
			Status:    model.PaymentStatus_COMPLETED,
		}, nil,
	)

	// Create a mock callback response
	callbackResponse := &model.CallbackResponse{
		Body: model.Body{
			StkCallback: model.StkCallback{
				MerchantRequestID: "123456",
				ResultCode:        0,
			},
		},
	}

	// Create a JSON request from the callback response
	requestJSON, _ := json.Marshal(callbackResponse)

	// Create a mock gin.Context with the JSON request
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest("POST", "/callback", bytes.NewReader(requestJSON))

	// Call the CallbackHandler function
	st.handler.CallbackHandler(ctx)

	// Check the response status code
	st.Require().Equal(http.StatusOK, ctx.Writer.Status())
}

func (st *PaymentHandlerTestSuite) TestCallbackHandler_UpdatePaymentStatusError() {
	output := make(chan grpcclients.ServiceResult)
	go func() {
		output <- grpcclients.ServiceResult{Error: nil, Result: &orderspb.Order{
			OrderId:     st.testUUID.String(),
			OrderStatus: orderspb.OrderStatus_ORDER_STATUS_PROCESSING,
		}}
	}()
	// Mock repository to return a payment record
	st.repo.On("GetPaymentByMerchantRequestId", mock.Anything, mock.Anything).Return(
		&model.Payment{
			PaymentID: st.testUUID1.String(),
			OrderID:   st.testUUID.String(),
		}, nil,
	)
	st.orderclient.On("UpdateOrderDetails", st.testUUID.String(), orderspb.OrderStatus_ORDER_STATUS_PROCESSING).Return(
		output,
	)
	st.repo.On("UpdatePaymentStatus", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("some error"))

	// Create a mock callback response
	callbackResponse := &model.CallbackResponse{
		Body: model.Body{
			StkCallback: model.StkCallback{
				MerchantRequestID: "123456",
				ResultCode:        0,
			},
		},
	}

	// Create a JSON request from the callback response
	requestJSON, _ := json.Marshal(callbackResponse)

	// Create a mock gin.Context with the JSON request
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest("POST", "/callback", bytes.NewReader(requestJSON))

	// Call the CallbackHandler function
	st.handler.CallbackHandler(ctx)

	// Check the response status code
	st.Require().Equal(http.StatusInternalServerError, ctx.Writer.Status())
}

func (st *PaymentHandlerTestSuite) TestCallbackHandler_UpdateOrderDetailsError() {
	output := make(chan grpcclients.ServiceResult)
	go func() {
		output <- grpcclients.ServiceResult{Error: errors.New("some error"), Result: nil}
	}()
	// Mock repository to return a payment record
	st.repo.On("GetPaymentByMerchantRequestId", mock.Anything, mock.Anything).Return(
		&model.Payment{
			PaymentID: st.testUUID1.String(),
			OrderID:   st.testUUID.String(),
		}, nil,
	)
	st.orderclient.On("UpdateOrderDetails", st.testUUID.String(), orderspb.OrderStatus_ORDER_STATUS_PROCESSING).Return(
		output,
	)

	// Create a mock callback response
	callbackResponse := &model.CallbackResponse{
		Body: model.Body{
			StkCallback: model.StkCallback{
				MerchantRequestID: "123456",
				ResultCode:        0,
			},
		},
	}

	// Create a JSON request from the callback response
	requestJSON, _ := json.Marshal(callbackResponse)

	// Create a mock gin.Context with the JSON request
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest("POST", "/callback", bytes.NewReader(requestJSON))

	// Call the CallbackHandler function
	st.handler.CallbackHandler(ctx)

	// Check the response status code
	st.Require().Equal(http.StatusInternalServerError, ctx.Writer.Status())
}

func (st *PaymentHandlerTestSuite) TestCallbackHandler_GetPaymentByMerchantRequestIdError() {

	// Mock repository to return a payment record
	st.repo.On("GetPaymentByMerchantRequestId", mock.Anything, mock.Anything).Return(
		&model.Payment{
			PaymentID: st.testUUID1.String(),
			OrderID:   st.testUUID.String(),
		}, errors.New("some error"),
	)

	// Create a mock callback response
	callbackResponse := &model.CallbackResponse{
		Body: model.Body{
			StkCallback: model.StkCallback{
				MerchantRequestID: "123456",
				ResultCode:        0,
			},
		},
	}

	// Create a JSON request from the callback response
	requestJSON, _ := json.Marshal(callbackResponse)

	// Create a mock gin.Context with the JSON request
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest("POST", "/callback", bytes.NewReader(requestJSON))

	// Call the CallbackHandler function
	st.handler.CallbackHandler(ctx)

	// Check the response status code
	st.Require().Equal(http.StatusInternalServerError, ctx.Writer.Status())
}

func (st *PaymentHandlerTestSuite) TestCallbackHandler_PaymentNotFound() {
	// Mock repository to return an error indicating payment not found
	st.repo.On("GetPaymentByMerchantRequestId", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)

	// Create a mock callback response
	callbackResponse := &model.CallbackResponse{
		Body: model.Body{
			StkCallback: model.StkCallback{
				MerchantRequestID: "123456",
				ResultCode:        0,
			},
		},
	}

	// Create a JSON request from the callback response
	requestJSON, _ := json.Marshal(callbackResponse)

	// Create a mock gin.Context with the JSON request
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest("POST", "/callback", bytes.NewReader(requestJSON))

	// Call the CallbackHandler function
	st.handler.CallbackHandler(ctx)

	// Check the response status code
	st.Require().Equal(http.StatusNotFound, ctx.Writer.Status())
}

func (st *PaymentHandlerTestSuite) TestCallbackHandler_InternalError() {
	// Mock repository to return an error indicating an internal error
	st.repo.On("GetPaymentByMerchantRequestId", mock.Anything, mock.Anything).Return(nil, errors.New("internal error"))

	// Create a mock callback response
	callbackResponse := &model.CallbackResponse{
		Body: model.Body{
			StkCallback: model.StkCallback{
				MerchantRequestID: "123456",
				ResultCode:        0,
			},
		},
	}

	// Create a JSON request from the callback response
	requestJSON, _ := json.Marshal(callbackResponse)

	// Create a mock gin.Context with the JSON request
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest("POST", "/callback", bytes.NewReader(requestJSON))

	// Call the CallbackHandler function
	st.handler.CallbackHandler(ctx)

	// Check the response status code
	st.Require().Equal(http.StatusInternalServerError, ctx.Writer.Status())
}

// CreatePayment function test cases are already provided in a previous response.
func (st *PaymentHandlerTestSuite) TestCallbackHandler_PaymentCanceled() {

	// Mock repository to return a payment record
	st.repo.On("GetPaymentByMerchantRequestId", mock.Anything, mock.Anything).Return(
		&model.Payment{
			PaymentID: st.testUUID1.String(),
			OrderID:   st.testUUID.String(),
		}, nil,
	)

	st.repo.On("UpdatePaymentStatus", mock.Anything, mock.Anything, mock.Anything).Return(
		&model.Payment{
			PaymentID: st.testUUID1.String(),
			OrderID:   st.testUUID.String(),
			Status:    model.PaymentStatus_COMPLETED,
		}, nil,
	)

	// Create a mock callback response
	callbackResponse := &model.CallbackResponse{
		Body: model.Body{
			StkCallback: model.StkCallback{
				MerchantRequestID: "123456",
				ResultCode:        1032,
			},
		},
	}

	// Create a JSON request from the callback response
	requestJSON, _ := json.Marshal(callbackResponse)

	// Create a mock gin.Context with the JSON request
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest("POST", "/callback", bytes.NewReader(requestJSON))

	// Call the CallbackHandler function
	st.handler.CallbackHandler(ctx)

	// Check the response status code
	st.Require().Equal(http.StatusPaymentRequired, ctx.Writer.Status())
}

func (st *PaymentHandlerTestSuite) TestCallbackHandler_PaymentCanceledInternalError() {

	// Mock repository to return a payment record
	st.repo.On("GetPaymentByMerchantRequestId", mock.Anything, mock.Anything).Return(
		&model.Payment{
			PaymentID: st.testUUID1.String(),
			OrderID:   st.testUUID.String(),
		}, nil,
	)

	st.repo.On("UpdatePaymentStatus", mock.Anything, mock.Anything, mock.Anything).Return(
		nil, errors.New("some error"),
	)

	// Create a mock callback response
	callbackResponse := &model.CallbackResponse{
		Body: model.Body{
			StkCallback: model.StkCallback{
				MerchantRequestID: "123456",
				ResultCode:        1032,
			},
		},
	}

	// Create a JSON request from the callback response
	requestJSON, _ := json.Marshal(callbackResponse)

	// Create a mock gin.Context with the JSON request
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest("POST", "/callback", bytes.NewReader(requestJSON))

	// Call the CallbackHandler function
	st.handler.CallbackHandler(ctx)

	// Check the response status code
	st.Require().Equal(http.StatusInternalServerError, ctx.Writer.Status())
}

func (st *PaymentHandlerTestSuite) TestCallbackHandler_STKCallbackResultCodeError() {

	// Mock repository to return a payment record
	st.repo.On("GetPaymentByMerchantRequestId", mock.Anything, mock.Anything).Return(
		&model.Payment{
			PaymentID: st.testUUID1.String(),
			OrderID:   st.testUUID.String(),
		}, nil,
	)

	st.repo.On("UpdatePaymentStatus", mock.Anything, mock.Anything, mock.Anything).Return(
		nil, errors.New("some error"),
	)

	// Create a mock callback response
	callbackResponse := &model.CallbackResponse{
		Body: model.Body{
			StkCallback: model.StkCallback{
				MerchantRequestID: "123456",
				ResultCode:        10,
			},
		},
	}

	// Create a JSON request from the callback response
	requestJSON, _ := json.Marshal(callbackResponse)

	// Create a mock gin.Context with the JSON request
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = httptest.NewRequest("POST", "/callback", bytes.NewReader(requestJSON))

	// Call the CallbackHandler function
	st.handler.CallbackHandler(ctx)

	// Check the response status code
	st.Require().Equal(http.StatusInternalServerError, ctx.Writer.Status())
}
