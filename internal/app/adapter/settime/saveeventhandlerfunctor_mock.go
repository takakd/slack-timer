// Code generated by MockGen. DO NOT EDIT.
// Source: ./saveeventhandlerfunctor.go

// Package settime is a generated GoMock package.
package settime

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockSaveEventHandler is a mock of SaveEventHandler interface
type MockSaveEventHandler struct {
	ctrl     *gomock.Controller
	recorder *MockSaveEventHandlerMockRecorder
}

// MockSaveEventHandlerMockRecorder is the mock recorder for MockSaveEventHandler
type MockSaveEventHandlerMockRecorder struct {
	mock *MockSaveEventHandler
}

// NewMockSaveEventHandler creates a new mock instance
func NewMockSaveEventHandler(ctrl *gomock.Controller) *MockSaveEventHandler {
	mock := &MockSaveEventHandler{ctrl: ctrl}
	mock.recorder = &MockSaveEventHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSaveEventHandler) EXPECT() *MockSaveEventHandlerMockRecorder {
	return m.recorder
}

// Handle mocks base method
func (m *MockSaveEventHandler) Handle(ctx context.Context, data EventCallbackData) *Response {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Handle", ctx, data)
	ret0, _ := ret[0].(*Response)
	return ret0
}

// Handle indicates an expected call of Handle
func (mr *MockSaveEventHandlerMockRecorder) Handle(ctx, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Handle", reflect.TypeOf((*MockSaveEventHandler)(nil).Handle), ctx, data)
}
