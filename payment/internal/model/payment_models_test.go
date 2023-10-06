package model

import (
	"testing"
	"time"

	paymentpb "github.com/wathuta/technical_test/protos_gen/payment"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/stretchr/testify/assert"
)

func TestPaymentFromProto(t *testing.T) {
	testCases := []struct {
		name     string
		input    *paymentpb.Payment
		expected *Payment
	}{
		{
			name: "Valid Payment Proto",
			input: &paymentpb.Payment{
				OrderId:       "order123",
				PaymentMethod: paymentpb.PaymentMethod_CREDIT_CARD,
				Amount:        100.0,
				Currency:      "KES",
				Status:        paymentpb.PaymentStatus_COMPLETED,
			},
			expected: &Payment{
				OrderID:       "order123",
				PaymentMethod: PaymentMethod_CREDIT_CARD,
				Amount:        100.0,
				Currency:      "KES",
				Status:        PaymentStatus_COMPLETED,
			},
		},
		// Add more test cases for different scenarios
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := PaymentFromProto(tc.input)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestPayment_Proto(t *testing.T) {
	createdAt := time.Now()
	updatedAt := time.Now()

	tests := []struct {
		name    string
		payment Payment
		want    *paymentpb.Payment
	}{
		{
			name: "Valid Payment",
			payment: Payment{
				PaymentID:     "payment123",
				OrderID:       "order456",
				PaymentMethod: PaymentMethod_CREDIT_CARD,
				Amount:        100.0,
				Status:        PaymentStatus_COMPLETED,
				ProductCost:   80.0,
				CustomerID:    "customer789",
				ShippingCost:  20.0,
				Currency:      "USD",
				CreatedAt:     createdAt,
				UpdatedAt:     updatedAt,
			},
			want: &paymentpb.Payment{
				Id:            "payment123",
				OrderId:       "order456",
				PaymentMethod: paymentpb.PaymentMethod_CREDIT_CARD,
				Amount:        100.0,
				Status:        paymentpb.PaymentStatus_COMPLETED,
				ProductCost:   80,
				CustomerId:    "customer789",
				ShippingFee:   20,
				Currency:      "USD",
				CreatedAt:     timestamppb.New(createdAt),
				UpdatedAt:     timestamppb.New(updatedAt),
			},
		},
		// Add more test cases here as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.payment.Proto()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCheckEnums(t *testing.T) {
	testCases := []struct {
		name     string
		input    *paymentpb.CreatePaymentRequest
		expected bool
	}{
		{
			name: "Valid Payment Method",
			input: &paymentpb.CreatePaymentRequest{
				PaymentMethod: paymentpb.PaymentMethod_CREDIT_CARD,
			},
			expected: false,
		},
		{
			name: "Invalid Payment Method",
			input: &paymentpb.CreatePaymentRequest{
				PaymentMethod: paymentpb.PaymentMethod(1),
			},
			expected: true,
		},
		// Add more test cases for different scenarios
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := CheckNotAValidEnum(tc.input)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
