// Code generated by mockery v1.0.0. DO NOT EDIT.

package transactions

import (
	models "github.com/lungria/spendshelf-backend/src/models"
	mock "github.com/stretchr/testify/mock"
	primitive "go.mongodb.org/mongo-driver/bson/primitive"
)

// MockRepository is an autogenerated mock type for the Repository type
type MockRepository struct {
	mock.Mock
}

// FindAll provides a mock function with given fields:
func (_m *MockRepository) FindAll() ([]models.Transaction, error) {
	ret := _m.Called()

	var r0 []models.Transaction
	if rf, ok := ret.Get(0).(func() []models.Transaction); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Transaction)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindAllByCategory provides a mock function with given fields: category
func (_m *MockRepository) FindAllByCategory(category string) ([]models.Transaction, error) {
	ret := _m.Called(category)

	var r0 []models.Transaction
	if rf, ok := ret.Get(0).(func(string) []models.Transaction); ok {
		r0 = rf(category)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Transaction)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(category)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindAllCategorized provides a mock function with given fields:
func (_m *MockRepository) FindAllCategorized() ([]models.Transaction, error) {
	ret := _m.Called()

	var r0 []models.Transaction
	if rf, ok := ret.Get(0).(func() []models.Transaction); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Transaction)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindAllUncategorized provides a mock function with given fields:
func (_m *MockRepository) FindAllUncategorized() ([]models.Transaction, error) {
	ret := _m.Called()

	var r0 []models.Transaction
	if rf, ok := ret.Get(0).(func() []models.Transaction); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Transaction)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateCategory provides a mock function with given fields: id, category
func (_m *MockRepository) UpdateCategory(id primitive.ObjectID, category string) error {
	ret := _m.Called(id, category)

	var r0 error
	if rf, ok := ret.Get(0).(func(primitive.ObjectID, string) error); ok {
		r0 = rf(id, category)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}