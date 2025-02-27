// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/SuperMatch/pkg/db/dao (interfaces: UserVerificationOTPRepository)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	model "github.com/SuperMatch/model"
	gomock "github.com/golang/mock/gomock"
)

// MockUserVerificationOTPRepository is a mock of UserVerificationOTPRepository interface.
type MockUserVerificationOTPRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUserVerificationOTPRepositoryMockRecorder
}

// MockUserVerificationOTPRepositoryMockRecorder is the mock recorder for MockUserVerificationOTPRepository.
type MockUserVerificationOTPRepositoryMockRecorder struct {
	mock *MockUserVerificationOTPRepository
}

// NewMockUserVerificationOTPRepository creates a new mock instance.
func NewMockUserVerificationOTPRepository(ctrl *gomock.Controller) *MockUserVerificationOTPRepository {
	mock := &MockUserVerificationOTPRepository{ctrl: ctrl}
	mock.recorder = &MockUserVerificationOTPRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserVerificationOTPRepository) EXPECT() *MockUserVerificationOTPRepositoryMockRecorder {
	return m.recorder
}

// FindByPhoneAndOTP mocks base method.
func (m *MockUserVerificationOTPRepository) FindByPhoneAndOTP(arg0, arg1 string) (model.UserVerificationOTP, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByPhoneAndOTP", arg0, arg1)
	ret0, _ := ret[0].(model.UserVerificationOTP)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByPhoneAndOTP indicates an expected call of FindByPhoneAndOTP.
func (mr *MockUserVerificationOTPRepositoryMockRecorder) FindByPhoneAndOTP(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByPhoneAndOTP", reflect.TypeOf((*MockUserVerificationOTPRepository)(nil).FindByPhoneAndOTP), arg0, arg1)
}

// Insert mocks base method.
func (m *MockUserVerificationOTPRepository) Insert(arg0 model.UserVerificationOTP) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Insert indicates an expected call of Insert.
func (mr *MockUserVerificationOTPRepositoryMockRecorder) Insert(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockUserVerificationOTPRepository)(nil).Insert), arg0)
}
