// Code generated by MockGen. DO NOT EDIT.
// Source: ../history/history.go

// Package mock is a generated GoMock package.
package mock

import (
	history "github.com/aleksanderaleksic/tgmigrate/history"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockHistory is a mock of History interface
type MockHistory struct {
	ctrl     *gomock.Controller
	recorder *MockHistoryMockRecorder
}

// MockHistoryMockRecorder is the mock recorder for MockHistory
type MockHistoryMockRecorder struct {
	mock *MockHistory
}

// NewMockHistory creates a new mock instance
func NewMockHistory(ctrl *gomock.Controller) *MockHistory {
	mock := &MockHistory{ctrl: ctrl}
	mock.recorder = &MockHistoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockHistory) EXPECT() *MockHistoryMockRecorder {
	return m.recorder
}

// IsMigrationApplied mocks base method
func (m *MockHistory) IsMigrationApplied(hash string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsMigrationApplied", hash)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsMigrationApplied indicates an expected call of IsMigrationApplied
func (mr *MockHistoryMockRecorder) IsMigrationApplied(hash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsMigrationApplied", reflect.TypeOf((*MockHistory)(nil).IsMigrationApplied), hash)
}

// InitializeHistory mocks base method
func (m *MockHistory) InitializeHistory() (*history.StorageHistory, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InitializeHistory")
	ret0, _ := ret[0].(*history.StorageHistory)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InitializeHistory indicates an expected call of InitializeHistory
func (mr *MockHistoryMockRecorder) InitializeHistory() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InitializeHistory", reflect.TypeOf((*MockHistory)(nil).InitializeHistory))
}

// StoreAppliedMigration mocks base method
func (m *MockHistory) StoreAppliedMigration(migration *history.AppliedStorageHistoryObject) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "StoreAppliedMigration", migration)
}

// StoreAppliedMigration indicates an expected call of StoreAppliedMigration
func (mr *MockHistoryMockRecorder) StoreAppliedMigration(migration interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreAppliedMigration", reflect.TypeOf((*MockHistory)(nil).StoreAppliedMigration), migration)
}

// StoreFailedMigration mocks base method
func (m *MockHistory) StoreFailedMigration(migration *history.FailedStorageHistoryObject) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "StoreFailedMigration", migration)
}

// StoreFailedMigration indicates an expected call of StoreFailedMigration
func (mr *MockHistoryMockRecorder) StoreFailedMigration(migration interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreFailedMigration", reflect.TypeOf((*MockHistory)(nil).StoreFailedMigration), migration)
}

// RemoveAppliedMigration mocks base method
func (m *MockHistory) RemoveAppliedMigration(migrationName string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RemoveAppliedMigration", migrationName)
}

// RemoveAppliedMigration indicates an expected call of RemoveAppliedMigration
func (mr *MockHistoryMockRecorder) RemoveAppliedMigration(migrationName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveAppliedMigration", reflect.TypeOf((*MockHistory)(nil).RemoveAppliedMigration), migrationName)
}

// WriteToStorage mocks base method
func (m *MockHistory) WriteToStorage() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteToStorage")
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteToStorage indicates an expected call of WriteToStorage
func (mr *MockHistoryMockRecorder) WriteToStorage() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteToStorage", reflect.TypeOf((*MockHistory)(nil).WriteToStorage))
}

// Cleanup mocks base method
func (m *MockHistory) Cleanup() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Cleanup")
}

// Cleanup indicates an expected call of Cleanup
func (mr *MockHistoryMockRecorder) Cleanup() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Cleanup", reflect.TypeOf((*MockHistory)(nil).Cleanup))
}
