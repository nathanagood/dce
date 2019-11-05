// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"
import time "time"
import usage "github.com/Optum/dce/pkg/usage"

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

// GetUsageByDateRange provides a mock function with given fields: startDate, endDate
func (_m *Service) GetUsageByDateRange(startDate time.Time, endDate time.Time) ([]*usage.Usage, error) {
	ret := _m.Called(startDate, endDate)

	var r0 []*usage.Usage
	if rf, ok := ret.Get(0).(func(time.Time, time.Time) []*usage.Usage); ok {
		r0 = rf(startDate, endDate)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*usage.Usage)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(time.Time, time.Time) error); ok {
		r1 = rf(startDate, endDate)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutUsage provides a mock function with given fields: input
func (_m *Service) PutUsage(input usage.Usage) error {
	ret := _m.Called(input)

	var r0 error
	if rf, ok := ret.Get(0).(func(usage.Usage) error); ok {
		r0 = rf(input)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
