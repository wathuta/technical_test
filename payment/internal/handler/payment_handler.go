package handler

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/wathuta/technical_test/payment/internal/common"
	"github.com/wathuta/technical_test/payment/internal/config"
	"github.com/wathuta/technical_test/payment/internal/model"
	paymentpb "github.com/wathuta/technical_test/protos_gen/payment"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) CreatePayment(ctx context.Context, req *paymentpb.CreatePaymentRequest) (*paymentpb.CreatePaymentResponse, error) {
	if req == nil || len(req.OrderId) == 0 || model.CheckNotAValidEnum(req) {
		slog.Error("invalid request", "error", errResourceRequired)
		return nil, errResourceRequired
	}
	slog.Debug("create payment", "order_id", req.OrderId)

	payment := &model.Payment{
		PaymentID:     uuid.New().String(),
		OrderID:       req.OrderId,
		CustomerID:    req.CustomerId,
		Description:   fmt.Sprintf("Payment for order %s", req.OrderId),
		Currency:      string(model.KES),
		PaymentMethod: model.PaymentMethod(req.PaymentMethod.String()),
		Amount:        req.Amount,
		ShippingCost:  float64(req.ShippingFee),
		ProductCost:   float64(req.ProductCost),
		Status:        model.PaymentStatus_PENDING,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	validator := common.NewValidator()

	if err := validator.Struct(payment); err != nil {
		slog.Error("failed to validate payment", "error", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if math.Ceil(payment.Amount) != math.Ceil(float64(req.ShippingFee)+payment.ProductCost) {
		slog.Error("invalid payment values shipping cost + product price is not equal to amount")
		return nil, errBadRequest
	}
	// Format the current time as "yyyyMMddHHmmss"
	formattedTime := time.Now().Format("20060102150405")

	// generating dajara api password
	combinedValue := model.BusinessSortCode + os.Getenv(config.MpesaPassKeyEnv) + formattedTime
	password := base64.StdEncoding.EncodeToString([]byte(combinedValue))

	callbackURL := fmt.Sprintf("%s%s", os.Getenv(config.CallBackBaseURL), "/callback")
	fmt.Println(callbackURL)
	resp, err := h.mpesa.InitiateSTKPushRequest(&model.STKPushRequestBody{
		Timestamp:         formattedTime,
		Amount:            int(math.Ceil(req.Amount)),
		Password:          password,
		TransactionType:   string(model.CustomerPayBillOnline),
		BusinessShortCode: model.BusinessSortCode,
		PartyA:            strings.ReplaceAll(req.CustomerPhone, "+", ""),
		PhoneNumber:       strings.ReplaceAll(req.CustomerPhone, "+", ""),
		PartyB:            model.BusinessSortCode,
		//To do replace order id with tracking number
		TransactionDesc:  fmt.Sprintf("Payment for order %s", req.OrderId),
		CallBackURL:      callbackURL,
		AccountReference: "Technical test",
	})
	if err != nil {
		slog.Error("initiating Stk push failed", "error", err)
		return nil, errInternal
	}

	payment.MerchantRequestID = resp.MerchantRequestID

	payment, err = h.repo.CreatePayment(ctx, payment)
	if err != nil {
		slog.Error("failed to create payment in db", "error", err)
		return nil, errInternal
	}

	slog.Debug("create payment successful")
	return &paymentpb.CreatePaymentResponse{Payment: payment.Proto()}, nil
}

func (h *Handler) GetPaymentById(ctx context.Context, req *paymentpb.GetPaymentByIdRequest) (*paymentpb.GetPaymentByIdResponse, error) {
	if req == nil || len(req.Id) == 0 {
		slog.Error("invalid request", "error", errResourceRequired)
		return nil, errResourceRequired
	}
	slog.Debug("get payment by id", "payment_id", req.Id)

	paymentUUID, err := uuid.Parse(req.Id)
	if err != nil {
		slog.Error("invalid payment uuid value", "error", err)
		return nil, errBadRequest
	}

	resource, err := h.repo.GetPaymentById(ctx, req.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Error("payment with the given payment_id not found", "payment_id", paymentUUID, "error", err)
			return nil, errNotFound
		}
		slog.Error("failed to get payment from db", "error", err)
		return nil, errInternal
	}

	slog.Debug("get payment successful")
	return &paymentpb.GetPaymentByIdResponse{Payment: resource.Proto()}, nil
}
