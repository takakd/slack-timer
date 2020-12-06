// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/app/driver/repository/dynamodbwrapper.go

// Package repository is a generated GoMock package.
package repository

import (
	reflect "reflect"

	dynamodb "github.com/aws/aws-sdk-go/service/dynamodb"
	gomock "github.com/golang/mock/gomock"
)

// MockDynamoDbWrapper is a mock of DynamoDbWrapper interface
type MockDynamoDbWrapper struct {
	ctrl     *gomock.Controller
	recorder *MockDynamoDbWrapperMockRecorder
}

// MockDynamoDbWrapperMockRecorder is the mock recorder for MockDynamoDbWrapper
type MockDynamoDbWrapperMockRecorder struct {
	mock *MockDynamoDbWrapper
}

// NewMockDynamoDbWrapper creates a new mock instance
func NewMockDynamoDbWrapper(ctrl *gomock.Controller) *MockDynamoDbWrapper {
	mock := &MockDynamoDbWrapper{ctrl: ctrl}
	mock.recorder = &MockDynamoDbWrapperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDynamoDbWrapper) EXPECT() *MockDynamoDbWrapperMockRecorder {
	return m.recorder
}

// GetItem mocks base method
func (m *MockDynamoDbWrapper) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetItem", input)
	ret0, _ := ret[0].(*dynamodb.GetItemOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetItem indicates an expected call of GetItem
func (mr *MockDynamoDbWrapperMockRecorder) GetItem(input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetItem", reflect.TypeOf((*MockDynamoDbWrapper)(nil).GetItem), input)
}

// Query mocks base method
func (m *MockDynamoDbWrapper) Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Query", input)
	ret0, _ := ret[0].(*dynamodb.QueryOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Query indicates an expected call of Query
func (mr *MockDynamoDbWrapperMockRecorder) Query(input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Query", reflect.TypeOf((*MockDynamoDbWrapper)(nil).Query), input)
}

// PutItem mocks base method
func (m *MockDynamoDbWrapper) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PutItem", input)
	ret0, _ := ret[0].(*dynamodb.PutItemOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PutItem indicates an expected call of PutItem
func (mr *MockDynamoDbWrapperMockRecorder) PutItem(input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutItem", reflect.TypeOf((*MockDynamoDbWrapper)(nil).PutItem), input)
}

// UnmarshalMap mocks base method
func (m_2 *MockDynamoDbWrapper) UnmarshalMap(m map[string]*dynamodb.AttributeValue, out interface{}) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "UnmarshalMap", m, out)
	ret0, _ := ret[0].(error)
	return ret0
}

// UnmarshalMap indicates an expected call of UnmarshalMap
func (mr *MockDynamoDbWrapperMockRecorder) UnmarshalMap(m, out interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnmarshalMap", reflect.TypeOf((*MockDynamoDbWrapper)(nil).UnmarshalMap), m, out)
}

// UnmarshalListOfMaps mocks base method
func (m *MockDynamoDbWrapper) UnmarshalListOfMaps(l []map[string]*dynamodb.AttributeValue, out interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnmarshalListOfMaps", l, out)
	ret0, _ := ret[0].(error)
	return ret0
}

// UnmarshalListOfMaps indicates an expected call of UnmarshalListOfMaps
func (mr *MockDynamoDbWrapperMockRecorder) UnmarshalListOfMaps(l, out interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnmarshalListOfMaps", reflect.TypeOf((*MockDynamoDbWrapper)(nil).UnmarshalListOfMaps), l, out)
}

// MarshalMap mocks base method
func (m *MockDynamoDbWrapper) MarshalMap(in interface{}) (map[string]*dynamodb.AttributeValue, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MarshalMap", in)
	ret0, _ := ret[0].(map[string]*dynamodb.AttributeValue)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MarshalMap indicates an expected call of MarshalMap
func (mr *MockDynamoDbWrapperMockRecorder) MarshalMap(in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MarshalMap", reflect.TypeOf((*MockDynamoDbWrapper)(nil).MarshalMap), in)
}
