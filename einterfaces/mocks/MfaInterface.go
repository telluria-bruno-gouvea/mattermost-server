// Code generated by mockery v1.0.0. DO NOT EDIT.

// Regenerate this file using `make einterfaces-mocks`.

package mocks

import (
	model "github.com/telluria-bruno-gouvea/mattermost-server/v5/model"
	mock "github.com/stretchr/testify/mock"
)

// MfaInterface is an autogenerated mock type for the MfaInterface type
type MfaInterface struct {
	mock.Mock
}

// Activate provides a mock function with given fields: user, token
func (_m *MfaInterface) Activate(user *model.User, token string) *model.AppError {
	ret := _m.Called(user, token)

	var r0 *model.AppError
	if rf, ok := ret.Get(0).(func(*model.User, string) *model.AppError); ok {
		r0 = rf(user, token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AppError)
		}
	}

	return r0
}

// Deactivate provides a mock function with given fields: userId
func (_m *MfaInterface) Deactivate(userId string) *model.AppError {
	ret := _m.Called(userId)

	var r0 *model.AppError
	if rf, ok := ret.Get(0).(func(string) *model.AppError); ok {
		r0 = rf(userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AppError)
		}
	}

	return r0
}

// GenerateSecret provides a mock function with given fields: user
func (_m *MfaInterface) GenerateSecret(user *model.User) (string, []byte, *model.AppError) {
	ret := _m.Called(user)

	var r0 string
	if rf, ok := ret.Get(0).(func(*model.User) string); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 []byte
	if rf, ok := ret.Get(1).(func(*model.User) []byte); ok {
		r1 = rf(user)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]byte)
		}
	}

	var r2 *model.AppError
	if rf, ok := ret.Get(2).(func(*model.User) *model.AppError); ok {
		r2 = rf(user)
	} else {
		if ret.Get(2) != nil {
			r2 = ret.Get(2).(*model.AppError)
		}
	}

	return r0, r1, r2
}

// ValidateToken provides a mock function with given fields: secret, token
func (_m *MfaInterface) ValidateToken(secret string, token string) (bool, *model.AppError) {
	ret := _m.Called(secret, token)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string, string) bool); ok {
		r0 = rf(secret, token)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string, string) *model.AppError); ok {
		r1 = rf(secret, token)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}
