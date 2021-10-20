// Code generated by mockery v1.0.0. DO NOT EDIT.

package service

import context "context"
import mock "github.com/stretchr/testify/mock"
import model "myapp3.0/internal/model"

import uuid "github.com/google/uuid"

// Cats is an autogenerated mock type for the Cats type
type MockCats struct {
	mock.Mock
}

// AddC provides a mock function with given fields: c, rec
func (_m *MockCats) AddC(c context.Context, rec *model.Record) (uuid.UUID, error) {
	ret := _m.Called(c, rec)

	var r0 uuid.UUID
	if rf, ok := ret.Get(0).(func(context.Context, *model.Record) uuid.UUID); ok {
		r0 = rf(c, rec)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(uuid.UUID)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *model.Record) error); ok {
		r1 = rf(c, rec)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteC provides a mock function with given fields: c, id
func (_m *MockCats) DeleteC(c context.Context, id uuid.UUID) error {
	ret := _m.Called(c, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) error); ok {
		r0 = rf(c, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAllC provides a mock function with given fields: c
func (_m *MockCats) GetAllC(c context.Context) ([]*model.Record, error) {
	ret := _m.Called(c)

	var r0 []*model.Record
	if rf, ok := ret.Get(0).(func(context.Context) []*model.Record); ok {
		r0 = rf(c)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Record)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(c)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetC provides a mock function with given fields: c, id
func (_m *MockCats) GetC(c context.Context, id uuid.UUID) (*model.Record, error) {
	ret := _m.Called(c, id)

	var r0 model.Record
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) model.Record); ok {
		r0 = rf(c, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(model.Record)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(c, id)
	} else {
		r1 = ret.Error(1)
	}

	return &r0, r1
}

// UpdateC provides a mock function with given fields: c, rec
func (_m *MockCats) UpdateC(c context.Context, rec *model.Record) error {
	ret := _m.Called(c, rec)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.Record) error); ok {
		r0 = rf(c, rec)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
