// Code generated by MockGen. DO NOT EDIT.
// Source: offeventhandlerfunctor.go

// Package settime is a generated GoMock package.
package settime

import (
	reflect "reflect"
	appcontext "slacktimer/internal/app/util/appcontext"

	gomock "github.com/golang/mock/gomock"
)

// MockOffEventHandler is a mock of OffEventHandler interface
type MockOffEventHandler struct {
	ctrl     *gomock.Controller
	recorder *MockOffEventHandlerMockRecorder
}

// MockOffEventHandlerMockRecorder is the mock recorder for MockOffEventHandler
type MockOffEventHandlerMockRecorder struct {
	mock *MockOffEventHandler
}

// NewMockOffEventHandler creates a new mock instance
func NewMockOffEventHandler(ctrl *gomock.Controller) *MockOffEventHandler {
	mock := &MockOffEventHandler{ctrl: ctrl}
	mock.recorder = &MockOffEventHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockOffEventHandler) EXPECT() *MockOffEventHandlerMockRecorder {
	return m.recorder
}

// Handle mocks base method
func (m *MockOffEventHandler) Handle(ac appcontext.AppContext, data EventCallbackData) *Response {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Handle", ac, data)
	ret0, _ := ret[0].(*Response)
	return ret0
}

// Handle indicates an expected call of Handle
func (mr *MockOffEventHandlerMockRecorder) Handle(ac, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Handle", reflect.TypeOf((*MockOffEventHandler)(nil).Handle), ac, data)
}
