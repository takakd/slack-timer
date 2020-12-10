// Code generated by MockGen. DO NOT EDIT.
// Source: inputport.go

// Package updatetimerevent is a generated GoMock package.
package updatetimerevent

import (
	reflect "reflect"
	appcontext "slacktimer/internal/app/util/appcontext"
	time "time"

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

// UpdateNotificationTime mocks base method
func (m *MockInputPort) UpdateNotificationTime(ac appcontext.AppContext, userID string, notificationTime time.Time, presenter OutputPort) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UpdateNotificationTime", ac, userID, notificationTime, presenter)
}

// UpdateNotificationTime indicates an expected call of UpdateNotificationTime
func (mr *MockInputPortMockRecorder) UpdateNotificationTime(ac, userID, notificationTime, presenter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateNotificationTime", reflect.TypeOf((*MockInputPort)(nil).UpdateNotificationTime), ac, userID, notificationTime, presenter)
}

// SaveIntervalMin mocks base method
func (m *MockInputPort) SaveIntervalMin(ac appcontext.AppContext, userID string, currentTime time.Time, minutes int, presenter OutputPort) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SaveIntervalMin", ac, userID, currentTime, minutes, presenter)
}

// SaveIntervalMin indicates an expected call of SaveIntervalMin
func (mr *MockInputPortMockRecorder) SaveIntervalMin(ac, userID, currentTime, minutes, presenter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveIntervalMin", reflect.TypeOf((*MockInputPort)(nil).SaveIntervalMin), ac, userID, currentTime, minutes, presenter)
}
