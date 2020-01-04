// Code generated by mockery v1.0.0. DO NOT EDIT.

package webhooks

import mock "github.com/stretchr/testify/mock"

// MockRepository is an autogenerated mock type for the Repository type
type MockRepository struct {
	mock.Mock
}

// GetAllHooks provides a mock function with given fields: accountID
func (_m *MockRepository) GetAllHooks(accountID string) ([]WebHook, error) {
	ret := _m.Called(accountID)

	var r0 []WebHook
	if rf, ok := ret.Get(0).(func(string) []WebHook); ok {
		r0 = rf(accountID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]WebHook)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(accountID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetHookByID provides a mock function with given fields: transactionID
func (_m *MockRepository) GetHookByID(transactionID string) (WebHook, error) {
	ret := _m.Called(transactionID)

	var r0 WebHook
	if rf, ok := ret.Get(0).(func(string) WebHook); ok {
		r0 = rf(transactionID)
	} else {
		r0 = ret.Get(0).(WebHook)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(transactionID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SaveOneHook provides a mock function with given fields: transaction
func (_m *MockRepository) SaveOneHook(transaction *WebHook) error {
	ret := _m.Called(transaction)

	var r0 error
	if rf, ok := ret.Get(0).(func(*WebHook) error); ok {
		r0 = rf(transaction)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}