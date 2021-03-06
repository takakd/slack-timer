// Code generated by MockGen. DO NOT EDIT.
// Source: inputport.go

// Package timeronevent is a generated GoMock package.
package timeronevent

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

// SetEventOn mocks base method
func (m *MockInputPort) SetEventOn(ac appcontext.AppContext, input InputData, presenter OutputPort) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetEventOn", ac, input, presenter)
}

// SetEventOn indicates an expected call of SetEventOn
func (mr *MockInputPortMockRecorder) SetEventOn(ac, input, presenter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetEventOn", reflect.TypeOf((*MockInputPort)(nil).SetEventOn), ac, input, presenter)
}
