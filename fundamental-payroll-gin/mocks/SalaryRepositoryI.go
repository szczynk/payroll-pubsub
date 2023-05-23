// Code generated by mockery v2.27.1. DO NOT EDIT.

package mocks

import (
	model "fundamental-payroll-gin/model"

	mock "github.com/stretchr/testify/mock"
)

// SalaryRepositoryI is an autogenerated mock type for the SalaryRepositoryI type
type SalaryRepositoryI struct {
	mock.Mock
}

// List provides a mock function with given fields:
func (_m *SalaryRepositoryI) List() ([]model.SalaryMatrix, error) {
	ret := _m.Called()

	var r0 []model.SalaryMatrix
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]model.SalaryMatrix, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []model.SalaryMatrix); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.SalaryMatrix)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewSalaryRepositoryI interface {
	mock.TestingT
	Cleanup(func())
}

// NewSalaryRepositoryI creates a new instance of SalaryRepositoryI. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewSalaryRepositoryI(t mockConstructorTestingTNewSalaryRepositoryI) *SalaryRepositoryI {
	mock := &SalaryRepositoryI{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
