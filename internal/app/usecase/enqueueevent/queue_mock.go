// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/app/usecase/enqueueevent/queue.go

// Package enqueueevent is a generated GoMock package.
package enqueueevent

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockQueue is a mock of Queue interface
type MockQueue struct {
	ctrl     *gomock.Controller
	recorder *MockQueueMockRecorder
}

// MockQueueMockRecorder is the mock recorder for MockQueue
type MockQueueMockRecorder struct {
	mock *MockQueue
}

// NewMockQueue creates a new mock instance
func NewMockQueue(ctrl *gomock.Controller) *MockQueue {
	mock := &MockQueue{ctrl: ctrl}
	mock.recorder = &MockQueueMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockQueue) EXPECT() *MockQueueMockRecorder {
	return m.recorder
}

// Enqueue mocks base method
func (m *MockQueue) Enqueue(message *QueueMessage) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Enqueue", message)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Enqueue indicates an expected call of Enqueue
func (mr *MockQueueMockRecorder) Enqueue(message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Enqueue", reflect.TypeOf((*MockQueue)(nil).Enqueue), message)
}