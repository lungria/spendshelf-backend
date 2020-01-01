// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/lungria/spendshelf-backend/src/db (interfaces: WebHookDB)
//todo regenerate
// Package mock_db is a generated GoMock package.
package mock_db

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	models "github.com/lungria/spendshelf-backend/src/models"
)

// MockWebHookDB is a mock of WebHookDB interface
type MockWebHookDB struct {
	ctrl     *gomock.Controller
	recorder *MockWebHookDBMockRecorder
}

// MockWebHookDBMockRecorder is the mock recorder for MockWebHookDB
type MockWebHookDBMockRecorder struct {
	mock *MockWebHookDB
}

// NewMockWebHookDB creates a new mock instance
func NewMockWebHookDB(ctrl *gomock.Controller) *MockWebHookDB {
	mock := &MockWebHookDB{ctrl: ctrl}
	mock.recorder = &MockWebHookDBMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockWebHookDB) EXPECT() *MockWebHookDBMockRecorder {
	return m.recorder
}

// GetAllTransactions mocks base method
func (m *MockWebHookDB) GetAllTransactions(arg0 string) ([]models.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllTransactions", arg0)
	ret0, _ := ret[0].([]models.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllTransactions indicates an expected call of GetAllTransactions
func (mr *MockWebHookDBMockRecorder) GetAllTransactions(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllTransactions", reflect.TypeOf((*MockWebHookDB)(nil).GetAllTransactions), arg0)
}

// GetTransactionByID mocks base method
func (m *MockWebHookDB) GetTransactionByID(arg0 string) (models.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransactionByID", arg0)
	ret0, _ := ret[0].(models.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTransactionByID indicates an expected call of GetTransactionByID
func (mr *MockWebHookDBMockRecorder) GetTransactionByID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransactionByID", reflect.TypeOf((*MockWebHookDB)(nil).GetTransactionByID), arg0)
}

// SaveOneTransaction mocks base method
func (m *MockWebHookDB) SaveOneTransaction(arg0 *models.Transaction) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveOneTransaction", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveOneTransaction indicates an expected call of SaveOneTransaction
func (mr *MockWebHookDBMockRecorder) SaveOneTransaction(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveOneTransaction", reflect.TypeOf((*MockWebHookDB)(nil).SaveOneTransaction), arg0)
}
