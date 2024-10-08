// Code generated by mockery v2.39.1. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// INotifier is an autogenerated mock type for the INotifier type
type INotifier struct {
	mock.Mock
}

type INotifier_Expecter struct {
	mock *mock.Mock
}

func (_m *INotifier) EXPECT() *INotifier_Expecter {
	return &INotifier_Expecter{mock: &_m.Mock}
}

// SendNotification provides a mock function with given fields: ctx, message
func (_m *INotifier) SendNotification(ctx context.Context, message string) {
	_m.Called(ctx, message)
}

// INotifier_SendNotification_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SendNotification'
type INotifier_SendNotification_Call struct {
	*mock.Call
}

// SendNotification is a helper method to define mock.On call
//   - ctx context.Context
//   - message string
func (_e *INotifier_Expecter) SendNotification(ctx interface{}, message interface{}) *INotifier_SendNotification_Call {
	return &INotifier_SendNotification_Call{Call: _e.mock.On("SendNotification", ctx, message)}
}

func (_c *INotifier_SendNotification_Call) Run(run func(ctx context.Context, message string)) *INotifier_SendNotification_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *INotifier_SendNotification_Call) Return() *INotifier_SendNotification_Call {
	_c.Call.Return()
	return _c
}

func (_c *INotifier_SendNotification_Call) RunAndReturn(run func(context.Context, string)) *INotifier_SendNotification_Call {
	_c.Call.Return(run)
	return _c
}

// NewINotifier creates a new instance of INotifier. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewINotifier(t interface {
	mock.TestingT
	Cleanup(func())
}) *INotifier {
	mock := &INotifier{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
