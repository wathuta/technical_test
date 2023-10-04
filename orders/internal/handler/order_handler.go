package handler

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/wathuta/technical_test/orders/internal/common"
	"github.com/wathuta/technical_test/orders/internal/common/fieldmask"
	"github.com/wathuta/technical_test/orders/internal/model"
	orderspb "github.com/wathuta/technical_test/protos_gen/orders"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Create a new order
func (h *Handler) CreateOrder(ctx context.Context, req *orderspb.CreateOrderRequest) (*orderspb.CreateOrderResponse, error) {
	if req == nil || req.Order == nil || len(req.ProductId) == 0 || req.ProductQuantity <= 0 {
		slog.Error("invalid request", "error", errResourceRequired)
		return nil, errResourceRequired
	}
	slog.Debug("create order and order details", req)

	resource := model.OrderFromProto(req.Order)
	resource.OrderID = uuid.New().String()
	resource.CreatedAt = time.Now()
	resource.DeletedAt = time.Time{}

	// validate request
	if err := common.ValidateGeneric(resource); err != nil {
		slog.Error("failed to validate product resource", "error", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	product, err := h.repo.GetProductById(ctx, req.ProductId)
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Error("product with the given product_id not found", "product_id", req.ProductId, "error", err)
			return nil, errNotFound
		}
		slog.Error("failed to get product from db", "error", err)
		return nil, errInternal
	}
	if product == nil {
		slog.Error("product with the given id not found", "product_id", req.ProductId, "error", err)
		return nil, errNotFound
	}

	customer, err := h.repo.GetCustomerById(ctx, req.Order.CustomerId)
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Error("customer with the given id not found", "customer_id", req.Order.CustomerId, "error", err)
			return nil, errNotFound
		}
		slog.Error("failed to get customer from db", "error", err)
		return nil, errInternal
	}

	if customer == nil {
		slog.Error("customer with the given id not found", "customer_id", req.Order.CustomerId, "error", err)
		return nil, errNotFound
	}
	orderdetails := &model.OrderDetails{
		OrderDetailsID: uuid.NewString(),
		OrderID:        resource.OrderID,
		ProductID:      req.ProductId,
		Quantity:       req.ProductQuantity,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Time{},
		DeletedAt:      time.Time{},
	}
	order, order_details, err := h.repo.CreateOrder(ctx, resource, orderdetails)
	if err != nil {
		slog.Error("failed to create order in db", "error", err)
		return nil, errInternal
	}

	slog.Debug("create order and order details successful")
	return &orderspb.CreateOrderResponse{
		Order:        order.Proto(),
		OrderDetails: order_details.Proto(),
	}, nil
}

// Get details of an order
func (h *Handler) GetOrderById(ctx context.Context, req *orderspb.GetOrderRequest) (*orderspb.GetOrderResponse, error) {
	if req == nil || len(req.OrderId) == 0 {
		slog.Error("invalid request", "error", errResourceRequired)
		return nil, errResourceRequired
	}
	slog.Debug("get order by id", "order_id", req.OrderId)

	orderUUID, err := uuid.Parse(req.OrderId)
	if err != nil {
		slog.Error("invalid order uuid value", "error", err)
		return nil, errBadRequest
	}

	order, err := h.repo.GetOrderById(ctx, orderUUID.String())
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Error("order with the given id not found", "order_id", orderUUID, "error", err)
			return nil, errNotFound
		}
		slog.Error("failed to get order from db", "error", err)
		return nil, errInternal
	}

	slog.Debug("get order successful")
	return &orderspb.GetOrderResponse{
		Order: order.Proto(),
	}, nil
}

// Update an order
func (h *Handler) UpdateOrder(ctx context.Context, req *orderspb.UpdateOrderRequest) (*orderspb.UpdateOrderResponse, error) {
	if req == nil || req.Order == nil {
		slog.Error("invalid request", "error", errResourceRequired)
		return nil, errResourceRequired
	}
	slog.Debug("update order", "order_id", req.Order.OrderId)

	mask, err := fieldmask.New(req.UpdateMask)
	if err != nil {
		slog.Error("invalid request inputs", "error", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	mask.RemoveOutputOnly()

	orderUUID, err := uuid.Parse(req.Order.OrderId)
	if err != nil {
		slog.Error("invalid order uuid value", "orderUUID", orderUUID, "error", err, req)
		return nil, errBadRequest
	}
	order := model.OrderFromProto(req.Order)
	updatedOrderDetails := model.UpdateOrderMaping(mask.Fields, *order)
	updatedOrderDetails["updated_at"] = time.Now()
	// if fieldmask is empty perfom get
	if len(mask.Fields) == 0 || len(updatedOrderDetails) == 0 {
		slog.Debug("no fields to update")
		order, err = h.repo.GetOrderById(ctx, orderUUID.String())
		if err != nil {
			if err == sql.ErrNoRows {
				slog.Error("order with the given id not found", "customer_id", orderUUID, "error", err)
				return nil, errNotFound
			}
			slog.Error("failed to get order from db", "error", err)
			return nil, errInternal
		}
		return &orderspb.UpdateOrderResponse{Order: order.Proto()}, nil
	}

	order, err = h.repo.UpdateOrder(ctx, orderUUID.String(), updatedOrderDetails)
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Error("customer with the given id not found", "customer_id", orderUUID, "error", err)
			return nil, errNotFound
		}
		slog.Error("failed to update customer from db", "error", err)
		return nil, errInternal
	}

	return &orderspb.UpdateOrderResponse{Order: order.Proto()}, nil
}

// Delete an order
func (h *Handler) DeleteOrder(ctx context.Context, req *orderspb.DeleteOrderRequest) (*orderspb.DeleteOrderResponse, error) {
	if req == nil || len(req.OrderId) == 0 {
		slog.Error("invalid request", "error", errResourceRequired)
		return &orderspb.DeleteOrderResponse{Success: false}, errResourceRequired
	}
	slog.Debug("delete order", "order_id", req.OrderId)

	// verify supplied uuid
	orderUUID, err := uuid.Parse(req.OrderId)
	if err != nil {
		slog.Error("invalid order uuid value", "error", err)
		return &orderspb.DeleteOrderResponse{Success: false}, errBadRequest
	}

	resource, err := h.repo.DeleteOrder(ctx, orderUUID.String())
	if err != nil {
		slog.Error("failed to order from db", "error", err)
		return &orderspb.DeleteOrderResponse{Success: false}, errInternal
	}
	if resource == nil {
		slog.Error("order with the given id not found", "order_id", orderUUID, "error", err)
		return &orderspb.DeleteOrderResponse{Success: false}, errNotFound
	}

	slog.Debug("delete order successful")
	return &orderspb.DeleteOrderResponse{
		Success: true,
	}, nil
}

// Get orders by customer ID
func (h *Handler) ListOrdersByCustomerId(ctx context.Context, req *orderspb.ListOrdersByCustomerIdRequest) (*orderspb.ListOrdersByCustomerIdResponse, error) {
	if req == nil || len(req.CustomerId) == 0 {
		slog.Error("invalid request", "error", errResourceRequired)
		return nil, errResourceRequired
	}
	slog.Debug("get order by customer id", "customer_id", req.CustomerId)

	customerUUID, err := uuid.Parse(req.CustomerId)
	if err != nil {
		slog.Error("invalid order uuid value", "error", err)
		return nil, errBadRequest
	}
	pagesize := common.SetPageSize(int(req.PageSize), defaultPageSize, maxPageSize)
	token := common.SetPageToken(int(req.PageToken))

	orders, err := h.repo.GetOrdersByCustomerId(ctx, customerUUID.String(), pagesize, token)
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Error("order with the given customer id not found", "customer_id", customerUUID, "error", err)
			return nil, errNotFound
		}
		slog.Error("failed to get order from db", "error", err)
		return nil, errInternal
	}

	returnOrders := []*orderspb.Order{}
	for _, order := range orders {
		newOrder := &order
		returnOrders = append(returnOrders, newOrder.Proto())
	}

	slog.Debug("get orders by customer id successful")
	return &orderspb.ListOrdersByCustomerIdResponse{Orders: returnOrders}, nil
}

// Get orders by product ID
func (h *Handler) ListOrdersByProductId(ctx context.Context, req *orderspb.ListOrdersByProductIdRequest) (*orderspb.ListOrdersByProductIdResponse, error) {
	if req == nil || len(req.ProductId) == 0 {
		slog.Error("invalid request", "error", errResourceRequired)
		return nil, errResourceRequired
	}
	slog.Debug("get order by product_id", "order_id", req.ProductId)

	productUUID, err := uuid.Parse(req.ProductId)
	if err != nil {
		slog.Error("invalid product uuid value", "error", err)
		return nil, errBadRequest
	}
	pagesize := common.SetPageSize(int(req.PageSize), defaultPageSize, maxPageSize)
	token := common.SetPageToken(int(req.PageToken))

	orderDetails, err := h.repo.GetOrderDetailsByProductId(ctx, productUUID.String(), pagesize, token)
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Error("order details with the given id not found", "order_id", productUUID, "error", err)
			return nil, errNotFound
		}
		slog.Error("failed to get order details from db", "error", err)
		return nil, errInternal
	}

	var returnOrders []*orderspb.Order
	var returnOrderDetails []*orderspb.OrderDetails

	for _, orderDetail := range orderDetails {
		order, err := h.repo.GetOrderById(ctx, orderDetail.OrderID)
		if err != nil {
			if err == sql.ErrNoRows {
				slog.Error("order details with the given id not found", "order_id", productUUID, "error", err)
				return nil, errNotFound
			}
			slog.Error("failed to get order details from db", "error", err)
			return nil, errInternal
		}
		returnOrders = append(returnOrders, order.Proto())
		returnOrderDetails = append(returnOrderDetails, orderDetail.Proto())
	}

	slog.Debug("get order successful")
	return &orderspb.ListOrdersByProductIdResponse{
		Orders:       returnOrders,
		OrderDetails: returnOrderDetails,
	}, nil
}

// Get order details by ID
func (h *Handler) GetOrderDetailsById(ctx context.Context, req *orderspb.GetOrderDetailByIdRequest) (*orderspb.GetOrderDetailByIdResponse, error) {
	if req == nil || len(req.OrderDetailsId) == 0 {
		slog.Error("invalid request", "error", errResourceRequired)
		return nil, errResourceRequired
	}
	slog.Debug("get order details by id", "order_details_id", req.OrderDetailsId)

	orderUUID, err := uuid.Parse(req.OrderDetailsId)
	if err != nil {
		slog.Error("invalid order details uuid value", "error", err)
		return nil, errBadRequest
	}

	orderDetails, err := h.repo.GetOrderDetailsById(ctx, orderUUID.String())
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Error("order details with the given id not found", "order_id", orderUUID, "error", err)
			return nil, errNotFound
		}
		slog.Error("failed to get order details from db", "error", err)
		return nil, errInternal
	}

	slog.Debug("get order details successful")
	return &orderspb.GetOrderDetailByIdResponse{
		OrderDetails: orderDetails.Proto(),
	}, nil
}

// Update order details by ID
func (h *Handler) UpdateOrderDetails(ctx context.Context, req *orderspb.UpdateOrderDetailsRequest) (*orderspb.UpdateOrderDetailsResponse, error) {
	return &orderspb.UpdateOrderDetailsResponse{}, nil
}

func (h *Handler) ListOrderDetailsByOrderId(ctx context.Context, req *orderspb.ListOrderDetailsByOrderIdRequest) (*orderspb.ListOrderDetailsByOrderIdResponse, error) {
	if req == nil || len(req.OrderId) == 0 {
		slog.Error("invalid request", "error", errResourceRequired)
		return nil, errResourceRequired
	}
	slog.Debug("get order details by order id", "order_id", req.OrderId)

	orderUUID, err := uuid.Parse(req.OrderId)
	if err != nil {
		slog.Error("invalid order uuid value", "error", err)
		return nil, errBadRequest
	}
	pagesize := common.SetPageSize(int(req.PageSize), defaultPageSize, maxPageSize)
	token := common.SetPageToken(int(req.PageToken))

	orderDetails, err := h.repo.GetOrderDetailsByOrderId(ctx, orderUUID.String(), pagesize, token)
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Error("order details with the given order id not found", "order_id", orderUUID, "error", err)
			return nil, errNotFound
		}
		slog.Error("failed to get order details from db", "error", err)
		return nil, errInternal
	}

	var returnOrderDetails []*orderspb.OrderDetails
	for _, orderDetail := range orderDetails {
		newOrderDetail := &orderDetail
		returnOrderDetails = append(returnOrderDetails, newOrderDetail.Proto())
	}

	slog.Debug("get order details by order id successful")
	return &orderspb.ListOrderDetailsByOrderIdResponse{
		OrderDetails: returnOrderDetails,
	}, nil
}
