// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import awsiface "github.com/Optum/dce/pkg/awsiface"
import client "github.com/aws/aws-sdk-go/aws/client"

import credentials "github.com/aws/aws-sdk-go/aws/credentials"
import mock "github.com/stretchr/testify/mock"
import sts "github.com/aws/aws-sdk-go/service/sts"

// TokenService is an autogenerated mock type for the TokenService type
type TokenService struct {
	mock.Mock
}

// AssumeRole provides a mock function with given fields: _a0
func (_m *TokenService) AssumeRole(_a0 *sts.AssumeRoleInput) (*sts.AssumeRoleOutput, error) {
	ret := _m.Called(_a0)

	var r0 *sts.AssumeRoleOutput
	if rf, ok := ret.Get(0).(func(*sts.AssumeRoleInput) *sts.AssumeRoleOutput); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sts.AssumeRoleOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*sts.AssumeRoleInput) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewCredentials provides a mock function with given fields: _a0, _a1
func (_m *TokenService) NewCredentials(_a0 client.ConfigProvider, _a1 string) *credentials.Credentials {
	ret := _m.Called(_a0, _a1)

	var r0 *credentials.Credentials
	if rf, ok := ret.Get(0).(func(client.ConfigProvider, string) *credentials.Credentials); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*credentials.Credentials)
		}
	}

	return r0
}

// NewSession provides a mock function with given fields: baseSession, roleArn
func (_m *TokenService) NewSession(baseSession awsiface.AwsSession, roleArn string) (awsiface.AwsSession, error) {
	ret := _m.Called(baseSession, roleArn)

	var r0 awsiface.AwsSession
	if rf, ok := ret.Get(0).(func(awsiface.AwsSession, string) awsiface.AwsSession); ok {
		r0 = rf(baseSession, roleArn)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(awsiface.AwsSession)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(awsiface.AwsSession, string) error); ok {
		r1 = rf(baseSession, roleArn)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
