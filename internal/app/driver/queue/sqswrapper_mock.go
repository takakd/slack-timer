// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/app/driver/queue/sqswrapper.go

// Package queue is a generated GoMock package.
package queue

import (
	sqs "github.com/aws/aws-sdk-go/service/sqs"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockSqsWrapper is a mock of SqsWrapper interface
type MockSqsWrapper struct {
	ctrl     *gomock.Controller
	recorder *MockSqsWrapperMockRecorder
}

// MockSqsWrapperMockRecorder is the mock recorder for MockSqsWrapper
type MockSqsWrapperMockRecorder struct {
	mock *MockSqsWrapper
}

// NewMockSqsWrapper creates a new mock instance
func NewMockSqsWrapper(ctrl *gomock.Controller) *MockSqsWrapper {
	mock := &MockSqsWrapper{ctrl: ctrl}
	mock.recorder = &MockSqsWrapperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSqsWrapper) EXPECT() *MockSqsWrapperMockRecorder {
	return m.recorder
}

// SendMessage mocks base method
func (m *MockSqsWrapper) SendMessage(input *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMessage", input)
	ret0, _ := ret[0].(*sqs.SendMessageOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendMessage indicates an expected call of SendMessage
func (mr *MockSqsWrapperMockRecorder) SendMessage(input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMessage", reflect.TypeOf((*MockSqsWrapper)(nil).SendMessage), input)
}