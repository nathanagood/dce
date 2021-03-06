// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import email "github.com/Optum/dce/pkg/email"
import mock "github.com/stretchr/testify/mock"

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

// SendEmail provides a mock function with given fields: input
func (_m *Service) SendEmail(input *email.SendEmailInput) error {
	ret := _m.Called(input)

	var r0 error
	if rf, ok := ret.Get(0).(func(*email.SendEmailInput) error); ok {
		r0 = rf(input)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SendRawEmailWithAttachment provides a mock function with given fields: input
func (_m *Service) SendRawEmailWithAttachment(input *email.SendEmailWithAttachmentInput) error {
	ret := _m.Called(input)

	var r0 error
	if rf, ok := ret.Get(0).(func(*email.SendEmailWithAttachmentInput) error); ok {
		r0 = rf(input)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
