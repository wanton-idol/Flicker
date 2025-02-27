// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/SuperMatch/pkg/db/dao (interfaces: UserMatchDao)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	model "github.com/SuperMatch/model"
	gomock "github.com/golang/mock/gomock"
)

// MockUserMatchDao is a mock of UserMatchDao interface.
type MockUserMatchDao struct {
	ctrl     *gomock.Controller
	recorder *MockUserMatchDaoMockRecorder
}

// MockUserMatchDaoMockRecorder is the mock recorder for MockUserMatchDao.
type MockUserMatchDaoMockRecorder struct {
	mock *MockUserMatchDao
}

// NewMockUserMatchDao creates a new mock instance.
func NewMockUserMatchDao(ctrl *gomock.Controller) *MockUserMatchDao {
	mock := &MockUserMatchDao{ctrl: ctrl}
	mock.recorder = &MockUserMatchDaoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserMatchDao) EXPECT() *MockUserMatchDaoMockRecorder {
	return m.recorder
}

// DeleteByUserID mocks base method.
func (m *MockUserMatchDao) DeleteByUserID(arg0 context.Context, arg1, arg2 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteByUserID", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteByUserID indicates an expected call of DeleteByUserID.
func (mr *MockUserMatchDaoMockRecorder) DeleteByUserID(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteByUserID", reflect.TypeOf((*MockUserMatchDao)(nil).DeleteByUserID), arg0, arg1, arg2)
}

// FindByUserId mocks base method.
func (m *MockUserMatchDao) FindByUserId(arg0 context.Context, arg1 int) ([]model.UserMatch, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByUserId", arg0, arg1)
	ret0, _ := ret[0].([]model.UserMatch)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByUserId indicates an expected call of FindByUserId.
func (mr *MockUserMatchDaoMockRecorder) FindByUserId(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByUserId", reflect.TypeOf((*MockUserMatchDao)(nil).FindByUserId), arg0, arg1)
}

// FindByUserIdMatchId mocks base method.
func (m *MockUserMatchDao) FindByUserIdMatchId(arg0 context.Context, arg1, arg2 int) (model.UserMatch, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByUserIdMatchId", arg0, arg1, arg2)
	ret0, _ := ret[0].(model.UserMatch)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByUserIdMatchId indicates an expected call of FindByUserIdMatchId.
func (mr *MockUserMatchDaoMockRecorder) FindByUserIdMatchId(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByUserIdMatchId", reflect.TypeOf((*MockUserMatchDao)(nil).FindByUserIdMatchId), arg0, arg1, arg2)
}

// Insert mocks base method.
func (m *MockUserMatchDao) Insert(arg0 context.Context, arg1 model.UserMatch) (model.UserMatch, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", arg0, arg1)
	ret0, _ := ret[0].(model.UserMatch)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Insert indicates an expected call of Insert.
func (mr *MockUserMatchDaoMockRecorder) Insert(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockUserMatchDao)(nil).Insert), arg0, arg1)
}

// InsertMany mocks base method.
func (m *MockUserMatchDao) InsertMany(arg0 context.Context, arg1 []model.UserMatch) ([]model.UserMatch, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertMany", arg0, arg1)
	ret0, _ := ret[0].([]model.UserMatch)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InsertMany indicates an expected call of InsertMany.
func (mr *MockUserMatchDaoMockRecorder) InsertMany(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertMany", reflect.TypeOf((*MockUserMatchDao)(nil).InsertMany), arg0, arg1)
}
