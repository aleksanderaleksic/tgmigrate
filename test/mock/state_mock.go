// Code generated by MockGen. DO NOT EDIT.
// Source: ../state/state.go

// Package mock is a generated GoMock package.
package mock

import (
	state "github.com/aleksanderaleksic/tgmigrate/state"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockState is a mock of State interface
type MockState struct {
	ctrl     *gomock.Controller
	recorder *MockStateMockRecorder
}

// MockStateMockRecorder is the mock recorder for MockState
type MockStateMockRecorder struct {
	mock *MockState
}

// NewMockState creates a new mock instance
func NewMockState(ctrl *gomock.Controller) *MockState {
	mock := &MockState{ctrl: ctrl}
	mock.recorder = &MockStateMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockState) EXPECT() *MockStateMockRecorder {
	return m.recorder
}

// InitializeState mocks base method
func (m *MockState) InitializeState() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InitializeState")
	ret0, _ := ret[0].(error)
	return ret0
}

// InitializeState indicates an expected call of InitializeState
func (mr *MockStateMockRecorder) InitializeState() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InitializeState", reflect.TypeOf((*MockState)(nil).InitializeState))
}

// Complete mocks base method
func (m *MockState) Complete() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Complete")
	ret0, _ := ret[0].(error)
	return ret0
}

// Complete indicates an expected call of Complete
func (mr *MockStateMockRecorder) Complete() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Complete", reflect.TypeOf((*MockState)(nil).Complete))
}

// Move mocks base method
func (m *MockState) Move(from, to state.ResourceContext) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Move", from, to)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Move indicates an expected call of Move
func (mr *MockStateMockRecorder) Move(from, to interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Move", reflect.TypeOf((*MockState)(nil).Move), from, to)
}

// Remove mocks base method
func (m *MockState) Remove(resource state.ResourceContext) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Remove", resource)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Remove indicates an expected call of Remove
func (mr *MockStateMockRecorder) Remove(resource interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*MockState)(nil).Remove), resource)
}

// Cleanup mocks base method
func (m *MockState) Cleanup() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Cleanup")
}

// Cleanup indicates an expected call of Cleanup
func (mr *MockStateMockRecorder) Cleanup() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Cleanup", reflect.TypeOf((*MockState)(nil).Cleanup))
}