package model

import (
	"fmt"
	"time"

	paymentpb "github.com/wathuta/technical_test/protos_gen/payment"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PaymentMethod string

const (
	PaymentMethod_CREDIT_CARD PaymentMethod = "CREDIT_CARD"
	PaymentMethod_PAYPAL      PaymentMethod = "PAYPAL"
	PaymentMethod_MPESA       PaymentMethod = "MPESA"
)

// PaymentStatus represents possible payment statuses.
type PaymentStatus string

const (
	PaymentStatus_PENDING   PaymentStatus = "PENDING"
	PaymentStatus_COMPLETED PaymentStatus = "COMPLETED"
	PaymentStatus_FAILED    PaymentStatus = "FAILED"
	PaymentStatus_CANCELED  PaymentStatus = "CANCELED"
)

type Payment struct {
	PaymentID         string        `validate:"required,uuid" db:"id"`
	OrderID           string        `validate:"required,uuid" db:"order_id"`
	CustomerID        string        `validate:"required,uuid" db:"customer_id"`
	PaymentMethod     PaymentMethod `validate:"required" db:"payment_method"`
	MerchantRequestID string        `validate:"omitempty" db:"merchant_request_id"`
	Amount            float64       `validate:"required" db:"amount"`
	Currency          string        `validate:"required" db:"currency"`
	Status            PaymentStatus `validate:"required" db:"status"`
	Description       string        `validate:"required" db:"description"`
	ShippingCost      float64       `validate:"required" db:"shipping_cost"`
	ProductCost       float64       `validate:"required" db:"product_cost"`
	CreatedAt         time.Time     `db:"created_at"`
	UpdatedAt         time.Time     `db:"updated_at"`
}

func PaymentFromProto(e *paymentpb.Payment) *Payment {
	return &Payment{
		OrderID:       e.OrderId,
		PaymentMethod: PaymentMethod(e.PaymentMethod.String()),
		Amount:        e.Amount,
		Currency:      e.Currency,
		Status:        PaymentStatus(e.Status.String()),
	}
}

func (p *Payment) Proto() *paymentpb.Payment {
	fmt.Println(p.Status)
	return &paymentpb.Payment{
		Id:            p.PaymentID,
		OrderId:       p.OrderID,
		PaymentMethod: paymentpb.PaymentMethod(paymentpb.PaymentMethod_value[string(p.PaymentMethod)]),
		Amount:        p.Amount,
		Status:        paymentpb.PaymentStatus(paymentpb.PaymentStatus_value[string(p.Status)]),
		ProductCost:   int64(p.ProductCost),
		CustomerId:    p.CustomerID,
		ShippingFee:   int64(p.ShippingCost),
		Currency:      p.Currency,
		CreatedAt:     timestamppb.New(p.CreatedAt),
		UpdatedAt:     timestamppb.New(p.UpdatedAt),
	}

}
func CheckNotAValidEnum(req *paymentpb.CreatePaymentRequest) bool {
	// Check if the PaymentMethod is neither CREDIT_CARD nor MPESA
	return req.PaymentMethod != paymentpb.PaymentMethod_CREDIT_CARD &&
		req.PaymentMethod != paymentpb.PaymentMethod_MPESA
}
