// Code generated by mockery v2.34.2. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	model "github.com/wathuta/technical_test/orders/internal/model"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// CreateCustomer provides a mock function with given fields: ctx, customer
func (_m *Repository) CreateCustomer(ctx context.Context, customer *model.Customer) (*model.Customer, error) {
	ret := _m.Called(ctx, customer)

	var r0 *model.Customer
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.Customer) (*model.Customer, error)); ok {
		return rf(ctx, customer)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.Customer) *model.Customer); ok {
		r0 = rf(ctx, customer)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Customer)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.Customer) error); ok {
		r1 = rf(ctx, customer)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateOrder provides a mock function with given fields: ctx, order, order_details
func (_m *Repository) CreateOrder(ctx context.Context, order *model.Order, order_details *model.OrderDetails) (*model.Order, *model.OrderDetails, error) {
	ret := _m.Called(ctx, order, order_details)

	var r0 *model.Order
	var r1 *model.OrderDetails
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.Order, *model.OrderDetails) (*model.Order, *model.OrderDetails, error)); ok {
		return rf(ctx, order, order_details)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.Order, *model.OrderDetails) *model.Order); ok {
		r0 = rf(ctx, order, order_details)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Order)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.Order, *model.OrderDetails) *model.OrderDetails); ok {
		r1 = rf(ctx, order, order_details)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.OrderDetails)
		}
	}

	if rf, ok := ret.Get(2).(func(context.Context, *model.Order, *model.OrderDetails) error); ok {
		r2 = rf(ctx, order, order_details)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// CreateProduct provides a mock function with given fields: ctx, product
func (_m *Repository) CreateProduct(ctx context.Context, product *model.Product) (*model.Product, error) {
	ret := _m.Called(ctx, product)

	var r0 *model.Product
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.Product) (*model.Product, error)); ok {
		return rf(ctx, product)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.Product) *model.Product); ok {
		r0 = rf(ctx, product)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Product)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.Product) error); ok {
		r1 = rf(ctx, product)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteCustomer provides a mock function with given fields: ctx, customerID
func (_m *Repository) DeleteCustomer(ctx context.Context, customerID string) (*model.Customer, error) {
	ret := _m.Called(ctx, customerID)

	var r0 *model.Customer
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*model.Customer, error)); ok {
		return rf(ctx, customerID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.Customer); ok {
		r0 = rf(ctx, customerID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Customer)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, customerID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteOrder provides a mock function with given fields: ctx, orderId
func (_m *Repository) DeleteOrder(ctx context.Context, orderId string) (*model.Order, error) {
	ret := _m.Called(ctx, orderId)

	var r0 *model.Order
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*model.Order, error)); ok {
		return rf(ctx, orderId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.Order); ok {
		r0 = rf(ctx, orderId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Order)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, orderId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteProduct provides a mock function with given fields: ctx, productId
func (_m *Repository) DeleteProduct(ctx context.Context, productId string) (*model.Product, error) {
	ret := _m.Called(ctx, productId)

	var r0 *model.Product
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*model.Product, error)); ok {
		return rf(ctx, productId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.Product); ok {
		r0 = rf(ctx, productId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Product)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, productId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCustomerById provides a mock function with given fields: ctx, customerID
func (_m *Repository) GetCustomerById(ctx context.Context, customerID string) (*model.Customer, error) {
	ret := _m.Called(ctx, customerID)

	var r0 *model.Customer
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*model.Customer, error)); ok {
		return rf(ctx, customerID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.Customer); ok {
		r0 = rf(ctx, customerID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Customer)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, customerID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOrderById provides a mock function with given fields: ctx, orderId
func (_m *Repository) GetOrderById(ctx context.Context, orderId string) (*model.Order, error) {
	ret := _m.Called(ctx, orderId)

	var r0 *model.Order
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*model.Order, error)); ok {
		return rf(ctx, orderId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.Order); ok {
		r0 = rf(ctx, orderId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Order)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, orderId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOrderDetailsById provides a mock function with given fields: ctx, orderDetailsId
func (_m *Repository) GetOrderDetailsById(ctx context.Context, orderDetailsId string) (*model.OrderDetails, error) {
	ret := _m.Called(ctx, orderDetailsId)

	var r0 *model.OrderDetails
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*model.OrderDetails, error)); ok {
		return rf(ctx, orderDetailsId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.OrderDetails); ok {
		r0 = rf(ctx, orderDetailsId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.OrderDetails)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, orderDetailsId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOrderDetailsByOrderId provides a mock function with given fields: ctx, orderId, limit, offset
func (_m *Repository) GetOrderDetailsByOrderId(ctx context.Context, orderId string, limit int, offset int) ([]model.OrderDetails, error) {
	ret := _m.Called(ctx, orderId, limit, offset)

	var r0 []model.OrderDetails
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, int, int) ([]model.OrderDetails, error)); ok {
		return rf(ctx, orderId, limit, offset)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, int, int) []model.OrderDetails); ok {
		r0 = rf(ctx, orderId, limit, offset)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.OrderDetails)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, int, int) error); ok {
		r1 = rf(ctx, orderId, limit, offset)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOrderDetailsByProductId provides a mock function with given fields: ctx, productId, limit, offset
func (_m *Repository) GetOrderDetailsByProductId(ctx context.Context, productId string, limit int, offset int) ([]model.OrderDetails, error) {
	ret := _m.Called(ctx, productId, limit, offset)

	var r0 []model.OrderDetails
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, int, int) ([]model.OrderDetails, error)); ok {
		return rf(ctx, productId, limit, offset)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, int, int) []model.OrderDetails); ok {
		r0 = rf(ctx, productId, limit, offset)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.OrderDetails)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, int, int) error); ok {
		r1 = rf(ctx, productId, limit, offset)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOrdersByCustomerId provides a mock function with given fields: ctx, customerId, limit, offset
func (_m *Repository) GetOrdersByCustomerId(ctx context.Context, customerId string, limit int, offset int) ([]model.Order, error) {
	ret := _m.Called(ctx, customerId, limit, offset)

	var r0 []model.Order
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, int, int) ([]model.Order, error)); ok {
		return rf(ctx, customerId, limit, offset)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, int, int) []model.Order); ok {
		r0 = rf(ctx, customerId, limit, offset)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Order)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, int, int) error); ok {
		r1 = rf(ctx, customerId, limit, offset)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetProductById provides a mock function with given fields: ctx, productId
func (_m *Repository) GetProductById(ctx context.Context, productId string) (*model.Product, error) {
	ret := _m.Called(ctx, productId)

	var r0 *model.Product
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*model.Product, error)); ok {
		return rf(ctx, productId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.Product); ok {
		r0 = rf(ctx, productId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Product)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, productId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateCustomerFields provides a mock function with given fields: ctx, customerID, updateFields
func (_m *Repository) UpdateCustomerFields(ctx context.Context, customerID string, updateFields map[string]interface{}) (*model.Customer, error) {
	ret := _m.Called(ctx, customerID, updateFields)

	var r0 *model.Customer
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, map[string]interface{}) (*model.Customer, error)); ok {
		return rf(ctx, customerID, updateFields)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, map[string]interface{}) *model.Customer); ok {
		r0 = rf(ctx, customerID, updateFields)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Customer)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, map[string]interface{}) error); ok {
		r1 = rf(ctx, customerID, updateFields)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateOrder provides a mock function with given fields: ctx, orderId, updateFields
func (_m *Repository) UpdateOrder(ctx context.Context, orderId string, updateFields map[string]interface{}) (*model.Order, error) {
	ret := _m.Called(ctx, orderId, updateFields)

	var r0 *model.Order
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, map[string]interface{}) (*model.Order, error)); ok {
		return rf(ctx, orderId, updateFields)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, map[string]interface{}) *model.Order); ok {
		r0 = rf(ctx, orderId, updateFields)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Order)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, map[string]interface{}) error); ok {
		r1 = rf(ctx, orderId, updateFields)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateProductFields provides a mock function with given fields: ctx, productId, updateFields
func (_m *Repository) UpdateProductFields(ctx context.Context, productId string, updateFields map[string]interface{}) (*model.Product, error) {
	ret := _m.Called(ctx, productId, updateFields)

	var r0 *model.Product
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, map[string]interface{}) (*model.Product, error)); ok {
		return rf(ctx, productId, updateFields)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, map[string]interface{}) *model.Product); ok {
		r0 = rf(ctx, productId, updateFields)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Product)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, map[string]interface{}) error); ok {
		r1 = rf(ctx, productId, updateFields)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewRepository creates a new instance of Repository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *Repository {
	mock := &Repository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}