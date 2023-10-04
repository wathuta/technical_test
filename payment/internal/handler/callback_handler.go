package handler

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wathuta/technical_test/payment/internal/model"
	"golang.org/x/exp/slog"
)

func (h *Handler) CallbackHandler(ctx *gin.Context) {
	callbackResponse := model.CallbackResponse{}
	err := ctx.BindJSON(&callbackResponse)
	if err != nil {
		slog.Error("failed to unmarshal callback request")
		ctx.JSON(http.StatusBadRequest, map[string]string{"status": "failed"})
		return
	}

	payment, err := h.repo.GetPaymentByMerchantRequestId(ctx, callbackResponse.Body.StkCallback.MerchantRequestID)
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Error("failed to Get payment record in db", "error", err)
			ctx.JSON(http.StatusNotFound, map[string]string{"message": "payment record not found"})
			return
		}
		slog.Error("failed to Get payment record in db", "error", err)
		ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal error"})
		return
	}
	fmt.Println(payment, callbackResponse)
	switch callbackResponse.Body.StkCallback.ResultCode {
	case 0:
		// result := <-h.OrderServiceClient.UpdateOrderDetails(payment.OrderID, orderspb.OrderStatus_ORDER_STATUS_PROCESSING)
		// if result.Error != nil {
		// 	slog.Error("failed to update order record from in order service", "error", err)
		// 	ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal error"})
		// }
		payment, err = h.repo.UpdatePaymentStatus(ctx, model.PaymentStatus_PENDING, payment.PaymentID)
		if err != nil {
			slog.Error("failed to update payment status in db", "error", err, "payment_id", payment.PaymentID)
			ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal error"})
			return
		}
		// to do update payment status
		slog.Debug("Update order status successful")
	case 1032:
		payment, err := h.repo.UpdatePaymentStatus(ctx, model.PaymentStatus_PENDING, payment.PaymentID)
		if err != nil {
			slog.Error("failed to update payment status in db", "error", err)
			ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal error"})
			return
		}
		slog.Debug("Transaction canceled by user", "payment_id", payment.PaymentID)
		ctx.JSON(http.StatusPaymentRequired, map[string]string{"status": "payment canceled"})
		return
	default:
		payment, err := h.repo.UpdatePaymentStatus(ctx, model.PaymentStatus_FAILED, payment.PaymentID)
		if err != nil {
			slog.Error("failed to update payment status in db", "error", err)
			ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal error"})
			return
		}
		slog.Info("Transaction failed", "payment_id", payment.PaymentID)
		ctx.JSON(http.StatusPaymentRequired, map[string]string{"status": "payment required"})
		return

	}
	ctx.JSON(http.StatusOK, map[string]string{"status": "payment processing"})
}
