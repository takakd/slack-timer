// Code generated by MockGen. DO NOT EDIT.
// Source: ./repository.go

// Package notifyevent is a generated GoMock package.
package notifyevent

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
	enterpriserule "slacktimer/internal/app/enterpriserule"
)

// MockRepository is a mock of Repository interface
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// FindTimerEvent mocks base method
func (m *MockRepository) FindTimerEvent(ctx context.Context, userId string) (*enterpriserule.TimerEvent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindTimerEvent", ctx, userId)
	ret0, _ := ret[0].(*enterpriserule.TimerEvent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindTimerEvent indicates an expected call of FindTimerEvent
func (mr *MockRepositoryMockRecorder) FindTimerEvent(ctx, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindTimerEvent", reflect.TypeOf((*MockRepository)(nil).FindTimerEvent), ctx, userId)
}

// SaveTimerEvent mocks base method
func (m *MockRepository) SaveTimerEvent(ctx context.Context, event *enterpriserule.TimerEvent) (*enterpriserule.TimerEvent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveTimerEvent", ctx, event)
	ret0, _ := ret[0].(*enterpriserule.TimerEvent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SaveTimerEvent indicates an expected call of SaveTimerEvent
func (mr *MockRepositoryMockRecorder) SaveTimerEvent(ctx, event interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveTimerEvent", reflect.TypeOf((*MockRepository)(nil).SaveTimerEvent), ctx, event)
}