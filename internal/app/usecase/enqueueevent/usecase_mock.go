// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/app/usecase/enqueueevent/usecase.go

// Package enqueueevent is a generated GoMock package.
package enqueueevent

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
	time "time"
)

// MockUsecase is a mock of Usecase interface
type MockUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockUsecaseMockRecorder
}

// MockUsecaseMockRecorder is the mock recorder for MockUsecase
type MockUsecaseMockRecorder struct {
	mock *MockUsecase
}

// NewMockUsecase creates a new mock instance
func NewMockUsecase(ctrl *gomock.Controller) *MockUsecase {
	mock := &MockUsecase{ctrl: ctrl}
	mock.recorder = &MockUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUsecase) EXPECT() *MockUsecaseMockRecorder {
	return m.recorder
}

// EnqueueEvent mocks base method
func (m *MockUsecase) EnqueueEvent(ctx context.Context, eventTime time.Time) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "EnqueueEvent", ctx, eventTime)
}

// EnqueueEvent indicates an expected call of EnqueueEvent
func (mr *MockUsecaseMockRecorder) EnqueueEvent(ctx, eventTime interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EnqueueEvent", reflect.TypeOf((*MockUsecase)(nil).EnqueueEvent), ctx, eventTime)
}

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
func (m *MockOutputPort) Output(data *OutputData) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Output", data)
}

// Output indicates an expected call of Output
func (mr *MockOutputPortMockRecorder) Output(data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Output", reflect.TypeOf((*MockOutputPort)(nil).Output), data)
}
