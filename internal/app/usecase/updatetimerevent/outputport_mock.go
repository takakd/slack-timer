// Code generated by MockGen. DO NOT EDIT.
// Source: outputport.go

// Package updatetimerevent is a generated GoMock package.
package updatetimerevent

import (
	reflect "reflect"
	appcontext "slacktimer/internal/app/util/appcontext"

	gomock "github.com/golang/mock/gomock"
)

// MockOutputPort is a mock of OutputPort interface
type MockOutputPort struct {
	ctrl     *gomock.Controller
	recorder *MockOutputPortMockRecorder
}

// MockOutputPortMockRecorder is the mock recorder for MockOutputPort
type MockOutputPortMockRecorder struct {
	mock *MockOutputPort
}

// NewMockOutputPort creates a new mock instance
func NewMockOutputPort(ctrl *gomock.Controller) *MockOutputPort {
	mock := &MockOutputPort{ctrl: ctrl}
	mock.recorder = &MockOutputPortMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockOutputPort) EXPECT() *MockOutputPortMockRecorder {
	return m.recorder
}

// Output mocks base method
func (m *MockOutputPort) Output(ac appcontext.AppContext, data OutputData) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Output", ac, data)
}

// Output indicates an expected call of Output
func (mr *MockOutputPortMockRecorder) Output(ac, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Output", reflect.TypeOf((*MockOutputPort)(nil).Output), ac, data)
}
