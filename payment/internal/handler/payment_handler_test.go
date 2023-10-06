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
	"github.com/wathuta/technical_test/payment/internal/config"
	"github.com/wathuta/technical_test/payment/internal/mocks"
	"github.com/wathuta/technical_test/payment/internal/model"
	paymentpb "github.com/wathuta/technical_test/protos_gen/payment"
	"golang.org/x/exp/slog"
)

type PaymentHandlerTestSuite struct {
	suite.Suite

	handler      *Handler
	repo         *mocks.Repository
	mpesaService *mocks.MpesaService
	orderclient  *mocks.OrderServiceClient

	testUUID  uuid.UUID
	testUUID1 uuid.UUID
}

func TestPaymentHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(PaymentHandlerTestSuite))
}

func (st *PaymentHandlerTestSuite) SetupSuite() {
	var programLevel = new(slog.LevelVar) // Info by default
	h := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: programLevel, AddSource: true})
	slog.SetDefault(slog.New(h))
	programLevel.Set(slog.LevelDebug)

	err := config.HasAllEnvVariables()
	st.Require().False(err)
}
func (st *PaymentHandlerTestSuite) SetupTest() {
	repo := mocks.NewRepository(st.T())
	mpesa := mocks.NewMpesaService(st.T())
	client := mocks.NewOrderServiceClient(st.T())

	st.handler = New(repo, mpesa, client)
	st.repo = repo
	st.mpesaService = mpesa
	st.orderclient = client
	st.testUUID = uuid.New()
	st.testUUID1 = uuid.New()
}

func (st *PaymentHandlerTestSuite) TestCreatePayment_Success() {
	payment := &paymentpb.CreatePaymentRequest{
		OrderId:       st.testUUID.String(),
		CustomerId:    st.testUUID1.String(),
		PaymentMethod: 2,
		Amount:        10,
		CustomerPhone: "+254724396746",
		ProductCost:   5,
		ShippingFee:   5,
	}
	st.mpesaService.On("InitiateSTKPushRequest", mock.Anything).Return(&model.STKPushRequestResponse{
		MerchantRequestID:   "29115-34620561-1",
		CheckoutRequestID:   "ws_CO_191220191020363925",
		ResponseCode:        "0",
		ResponseDescription: "Success. Request accepted for processing",
		CustomerMessage:     "Success. Request accepted for processing",
	}, nil)
	st.repo.On("CreatePayment", mock.Anything, mock.Anything).Return(&model.Payment{
		OrderID:           st.testUUID.String(),
		CustomerID:        st.testUUID1.String(),
		PaymentMethod:     model.PaymentMethod_MPESA,
		Amount:            10,
		ProductCost:       5,
		ShippingCost:      5,
		PaymentID:         st.testUUID1.String(),
		MerchantRequestID: "29115-34620561-1",
		Currency:          "KES",
		Status:            model.PaymentStatus_COMPLETED,
		Description:       "payment for order qerty-erty-cvbn-yh9ik",
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}, nil)

	resp, err := st.handler.CreatePayment(context.Background(), payment)
	st.Require().Nil(err)
	st.Require().NotEmpty(resp.Payment.Id)
}

func (st *PaymentHandlerTestSuite) TestCreatePayment_CreatePaymentError() {
	payment := &paymentpb.CreatePaymentRequest{
		OrderId:       st.testUUID.String(),
		CustomerId:    st.testUUID1.String(),
		PaymentMethod: 2,
		Amount:        10,
		CustomerPhone: "+254724396746",
		ProductCost:   5,
		ShippingFee:   5,
	}
	st.mpesaService.On("InitiateSTKPushRequest", mock.Anything).Return(&model.STKPushRequestResponse{
		MerchantRequestID:   "29115-34620561-1",
		CheckoutRequestID:   "ws_CO_191220191020363925",
		ResponseCode:        "0",
		ResponseDescription: "Success. Request accepted for processing",
		CustomerMessage:     "Success. Request accepted for processing",
	}, nil)
	st.repo.On("CreatePayment", mock.Anything, mock.Anything).Return(nil, errors.New("some error"))

	resp, err := st.handler.CreatePayment(context.Background(), payment)
	st.Require().NotNil(err)
	st.Require().Nil(resp)
}

func (st *PaymentHandlerTestSuite) TestCreatePayment_ErrorInitiateSTKPushFailed() {
	payment := &paymentpb.CreatePaymentRequest{
		OrderId:       st.testUUID.String(),
		CustomerId:    st.testUUID1.String(),
		PaymentMethod: 2,
		Amount:        10,
		CustomerPhone: "+254724396746",
		ProductCost:   5,
		ShippingFee:   5,
	}
	st.mpesaService.On("InitiateSTKPushRequest", mock.Anything).Return(nil, errors.New("some error"))

	resp, err := st.handler.CreatePayment(context.Background(), payment)
	st.Require().NotNil(err)
	st.Require().Nil(resp)
}

func (st *PaymentHandlerTestSuite) TestCreatePayment_StructValidationError() {
	payment := &paymentpb.CreatePaymentRequest{
		PaymentMethod: 2,
		Amount:        10,
		CustomerPhone: "+254724396746",
		ProductCost:   5,
		ShippingFee:   5,
	}

	resp, err := st.handler.CreatePayment(context.Background(), payment)
	st.Require().NotNil(err)
	st.Require().Nil(resp)
}

func (st *PaymentHandlerTestSuite) TestCreatePayment_NilRequestError() {
	resp, err := st.handler.CreatePayment(context.Background(), nil)
	st.Require().NotNil(err)
	st.Require().Nil(resp)
}
func (st *PaymentHandlerTestSuite) TestGetPaymentById_Success() {
	payment := &model.Payment{
		OrderID:           st.testUUID.String(),
		CustomerID:        st.testUUID1.String(),
		PaymentMethod:     model.PaymentMethod_MPESA,
		Amount:            10,
		ProductCost:       5,
		ShippingCost:      5,
		PaymentID:         st.testUUID1.String(),
		MerchantRequestID: "29115-34620561-1",
		Currency:          "KES",
		Status:            model.PaymentStatus_COMPLETED,
		Description:       "payment for order qerty-erty-cvbn-yh9ik",
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	st.repo.On("GetPaymentById", mock.Anything, mock.Anything).Return(payment, nil)

	resp, err := st.handler.GetPaymentById(context.Background(), &paymentpb.GetPaymentByIdRequest{
		Id: st.testUUID.String(),
	})
	st.Require().Nil(err)
	st.Require().NotEmpty(resp.Payment.Id)
}
func (st *PaymentHandlerTestSuite) TestGetPaymentById_InvalidUUIDError() {

	resp, err := st.handler.GetPaymentById(context.Background(), &paymentpb.GetPaymentByIdRequest{
		Id: "some uuid",
	})

	st.Require().Nil(resp)
	st.Require().NotNil(err)
}

func (st *PaymentHandlerTestSuite) TestGetPaymentById_DBError() {

	st.repo.On("GetPaymentById", mock.Anything, mock.Anything).Return(nil, errors.New("some error"))

	resp, err := st.handler.GetPaymentById(context.Background(), &paymentpb.GetPaymentByIdRequest{
		Id: st.testUUID.String(),
	})
	st.Require().Nil(resp)
	st.Require().NotNil(err)
	st.Require().Equal(errInternal, err)
}

func (st *PaymentHandlerTestSuite) TestGetPaymentById_NotFoundError() {

	st.repo.On("GetPaymentById", mock.Anything, mock.Anything).Return(nil, sql.ErrNoRows)

	resp, err := st.handler.GetPaymentById(context.Background(), &paymentpb.GetPaymentByIdRequest{
		Id: st.testUUID.String(),
	})
	st.Require().Nil(resp)
	st.Require().NotNil(err)
	st.Require().Equal(errNotFound, err)

}

func (st *PaymentHandlerTestSuite) TestGetPaymentById_NilRequest() {

	resp, err := st.handler.GetPaymentById(context.Background(), nil)

	st.Require().Nil(resp)
	st.Require().NotNil(err)
}
