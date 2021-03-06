// Code generated by MockGen. DO NOT EDIT.
// Source: ./inputport.go

// Package enqueueevent is a generated GoMock package.
package enqueueevent

import (
	reflect "reflect"
	appcontext "slacktimer/internal/app/util/appcontext"

	gomock "github.com/golang/mock/gomock"
)

// MockInputPort is a mock of InputPort interface
type MockInputPort struct {
	ctrl     *gomock.Controller
	recorder *MockInputPortMockRecorder
}

// MockInputPortMockRecorder is the mock recorder for MockInputPort
type MockInputPortMockRecorder struct {
	mock *MockInputPort
}

// NewMockInputPort creates a new mock instance
func NewMockInputPort(ctrl *gomock.Controller) *MockInputPort {
	mock := &MockInputPort{ctrl: ctrl}
	mock.recorder = &MockInputPortMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockInputPort) EXPECT() *MockInputPortMockRecorder {
	return m.recorder
}

// EnqueueEvent mocks base method
func (m *MockInputPort) EnqueueEvent(ac appcontext.AppContext, data InputData) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "EnqueueEvent", ac, data)
}

// EnqueueEvent indicates an expected call of EnqueueEvent
func (mr *MockInputPortMockRecorder) EnqueueEvent(ac, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EnqueueEvent", reflect.TypeOf((*MockInputPort)(nil).EnqueueEvent), ac, data)
}
