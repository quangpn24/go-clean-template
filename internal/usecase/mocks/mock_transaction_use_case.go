// Code generated by mockery v2.39.1. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// ITransactionUseCase is an autogenerated mock type for the ITransactionUseCase type
type ITransactionUseCase struct {
	mock.Mock
}

type ITransactionUseCase_Expecter struct {
	mock *mock.Mock
}

func (_m *ITransactionUseCase) EXPECT() *ITransactionUseCase_Expecter {
	return &ITransactionUseCase_Expecter{mock: &_m.Mock}
}

// Deposit provides a mock function with given fields: ctx, walletID, accountID, amount, currency, note
func (_m *ITransactionUseCase) Deposit(ctx context.Context, walletID string, accountID string, amount float64, currency string, note string) error {
	ret := _m.Called(ctx, walletID, accountID, amount, currency, note)

	if len(ret) == 0 {
		panic("no return value specified for Deposit")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, float64, string, string) error); ok {
		r0 = rf(ctx, walletID, accountID, amount, currency, note)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ITransactionUseCase_Deposit_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Deposit'
type ITransactionUseCase_Deposit_Call struct {
	*mock.Call
}

// Deposit is a helper method to define mock.On call
//   - ctx context.Context
//   - walletID string
//   - accountID string
//   - amount float64
//   - currency string
//   - note string
func (_e *ITransactionUseCase_Expecter) Deposit(ctx interface{}, walletID interface{}, accountID interface{}, amount interface{}, currency interface{}, note interface{}) *ITransactionUseCase_Deposit_Call {
	return &ITransactionUseCase_Deposit_Call{Call: _e.mock.On("Deposit", ctx, walletID, accountID, amount, currency, note)}
}

func (_c *ITransactionUseCase_Deposit_Call) Run(run func(ctx context.Context, walletID string, accountID string, amount float64, currency string, note string)) *ITransactionUseCase_Deposit_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(float64), args[4].(string), args[5].(string))
	})
	return _c
}

func (_c *ITransactionUseCase_Deposit_Call) Return(_a0 error) *ITransactionUseCase_Deposit_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ITransactionUseCase_Deposit_Call) RunAndReturn(run func(context.Context, string, string, float64, string, string) error) *ITransactionUseCase_Deposit_Call {
	_c.Call.Return(run)
	return _c
}

// PayTransaction provides a mock function with given fields: ctx, transID
func (_m *ITransactionUseCase) PayTransaction(ctx context.Context, transID string) error {
	ret := _m.Called(ctx, transID)

	if len(ret) == 0 {
		panic("no return value specified for PayTransaction")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, transID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ITransactionUseCase_PayTransaction_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PayTransaction'
type ITransactionUseCase_PayTransaction_Call struct {
	*mock.Call
}

// PayTransaction is a helper method to define mock.On call
//   - ctx context.Context
//   - transID string
func (_e *ITransactionUseCase_Expecter) PayTransaction(ctx interface{}, transID interface{}) *ITransactionUseCase_PayTransaction_Call {
	return &ITransactionUseCase_PayTransaction_Call{Call: _e.mock.On("PayTransaction", ctx, transID)}
}

func (_c *ITransactionUseCase_PayTransaction_Call) Run(run func(ctx context.Context, transID string)) *ITransactionUseCase_PayTransaction_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *ITransactionUseCase_PayTransaction_Call) Return(_a0 error) *ITransactionUseCase_PayTransaction_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ITransactionUseCase_PayTransaction_Call) RunAndReturn(run func(context.Context, string) error) *ITransactionUseCase_PayTransaction_Call {
	_c.Call.Return(run)
	return _c
}

// Withdraw provides a mock function with given fields: ctx, walletID, accountID, amount, currency, note
func (_m *ITransactionUseCase) Withdraw(ctx context.Context, walletID string, accountID string, amount float64, currency string, note string) error {
	ret := _m.Called(ctx, walletID, accountID, amount, currency, note)

	if len(ret) == 0 {
		panic("no return value specified for Withdraw")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, float64, string, string) error); ok {
		r0 = rf(ctx, walletID, accountID, amount, currency, note)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ITransactionUseCase_Withdraw_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Withdraw'
type ITransactionUseCase_Withdraw_Call struct {
	*mock.Call
}

// Withdraw is a helper method to define mock.On call
//   - ctx context.Context
//   - walletID string
//   - accountID string
//   - amount float64
//   - currency string
//   - note string
func (_e *ITransactionUseCase_Expecter) Withdraw(ctx interface{}, walletID interface{}, accountID interface{}, amount interface{}, currency interface{}, note interface{}) *ITransactionUseCase_Withdraw_Call {
	return &ITransactionUseCase_Withdraw_Call{Call: _e.mock.On("Withdraw", ctx, walletID, accountID, amount, currency, note)}
}

func (_c *ITransactionUseCase_Withdraw_Call) Run(run func(ctx context.Context, walletID string, accountID string, amount float64, currency string, note string)) *ITransactionUseCase_Withdraw_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(float64), args[4].(string), args[5].(string))
	})
	return _c
}

func (_c *ITransactionUseCase_Withdraw_Call) Return(_a0 error) *ITransactionUseCase_Withdraw_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ITransactionUseCase_Withdraw_Call) RunAndReturn(run func(context.Context, string, string, float64, string, string) error) *ITransactionUseCase_Withdraw_Call {
	_c.Call.Return(run)
	return _c
}

// NewITransactionUseCase creates a new instance of ITransactionUseCase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewITransactionUseCase(t interface {
	mock.TestingT
	Cleanup(func())
}) *ITransactionUseCase {
	mock := &ITransactionUseCase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
