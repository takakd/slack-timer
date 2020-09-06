// Code generated by MockGen. DO NOT EDIT.
// Source: internal/app/driver/mongodb.go

// Package driver is a generated GoMock package.
package driver

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	mongo "go.mongodb.org/mongo-driver/mongo"
	options "go.mongodb.org/mongo-driver/mongo/options"
	reflect "reflect"
)

// MockMongoDbConnector is a mock of MongoDbConnector interface
type MockMongoDbConnector struct {
	ctrl     *gomock.Controller
	recorder *MockMongoDbConnectorMockRecorder
}

// MockMongoDbConnectorMockRecorder is the mock recorder for MockMongoDbConnector
type MockMongoDbConnectorMockRecorder struct {
	mock *MockMongoDbConnector
}

// NewMockMongoDbConnector creates a new mock instance
func NewMockMongoDbConnector(ctrl *gomock.Controller) *MockMongoDbConnector {
	mock := &MockMongoDbConnector{ctrl: ctrl}
	mock.recorder = &MockMongoDbConnectorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMongoDbConnector) EXPECT() *MockMongoDbConnectorMockRecorder {
	return m.recorder
}

// GetDb mocks base method
func (m *MockMongoDbConnector) GetDb(ctx context.Context, mongoUrl string) (*mongo.Database, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDb", ctx, mongoUrl)
	ret0, _ := ret[0].(*mongo.Database)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDb indicates an expected call of GetDb
func (mr *MockMongoDbConnectorMockRecorder) GetDb(ctx, mongoUrl interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDb", reflect.TypeOf((*MockMongoDbConnector)(nil).GetDb), ctx, mongoUrl)
}

// GetCollection mocks base method
func (m *MockMongoDbConnector) GetCollection(db *mongo.Database, name string) *mongo.Collection {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCollection", db, name)
	ret0, _ := ret[0].(*mongo.Collection)
	return ret0
}

// GetCollection indicates an expected call of GetCollection
func (mr *MockMongoDbConnectorMockRecorder) GetCollection(db, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCollection", reflect.TypeOf((*MockMongoDbConnector)(nil).GetCollection), db, name)
}

// DisConnectClientFunc mocks base method
func (m *MockMongoDbConnector) DisConnectClientFunc(ctx context.Context, client *mongo.Client, f func(error)) func() {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DisConnectClientFunc", ctx, client, f)
	ret0, _ := ret[0].(func())
	return ret0
}

// DisConnectClientFunc indicates an expected call of DisConnectClientFunc
func (mr *MockMongoDbConnectorMockRecorder) DisConnectClientFunc(ctx, client, f interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DisConnectClientFunc", reflect.TypeOf((*MockMongoDbConnector)(nil).DisConnectClientFunc), ctx, client, f)
}

// MockMongoDatabase is a mock of MongoDatabase interface
type MockMongoDatabase struct {
	ctrl     *gomock.Controller
	recorder *MockMongoDatabaseMockRecorder
}

// MockMongoDatabaseMockRecorder is the mock recorder for MockMongoDatabase
type MockMongoDatabaseMockRecorder struct {
	mock *MockMongoDatabase
}

// NewMockMongoDatabase creates a new mock instance
func NewMockMongoDatabase(ctrl *gomock.Controller) *MockMongoDatabase {
	mock := &MockMongoDatabase{ctrl: ctrl}
	mock.recorder = &MockMongoDatabaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMongoDatabase) EXPECT() *MockMongoDatabaseMockRecorder {
	return m.recorder
}

// Client mocks base method
func (m *MockMongoDatabase) Client() *mongo.Client {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Client")
	ret0, _ := ret[0].(*mongo.Client)
	return ret0
}

// Client indicates an expected call of Client
func (mr *MockMongoDatabaseMockRecorder) Client() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Client", reflect.TypeOf((*MockMongoDatabase)(nil).Client))
}

// MockMongoCollection is a mock of MongoCollection interface
type MockMongoCollection struct {
	ctrl     *gomock.Controller
	recorder *MockMongoCollectionMockRecorder
}

// MockMongoCollectionMockRecorder is the mock recorder for MockMongoCollection
type MockMongoCollectionMockRecorder struct {
	mock *MockMongoCollection
}

// NewMockMongoCollection creates a new mock instance
func NewMockMongoCollection(ctrl *gomock.Controller) *MockMongoCollection {
	mock := &MockMongoCollection{ctrl: ctrl}
	mock.recorder = &MockMongoCollectionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMongoCollection) EXPECT() *MockMongoCollectionMockRecorder {
	return m.recorder
}

// FindOne mocks base method
func (m *MockMongoCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, filter}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "FindOne", varargs...)
	ret0, _ := ret[0].(*mongo.SingleResult)
	return ret0
}

// FindOne indicates an expected call of FindOne
func (mr *MockMongoCollectionMockRecorder) FindOne(ctx, filter interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, filter}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOne", reflect.TypeOf((*MockMongoCollection)(nil).FindOne), varargs...)
}

// MockMongoSingleResult is a mock of MongoSingleResult interface
type MockMongoSingleResult struct {
	ctrl     *gomock.Controller
	recorder *MockMongoSingleResultMockRecorder
}

// MockMongoSingleResultMockRecorder is the mock recorder for MockMongoSingleResult
type MockMongoSingleResultMockRecorder struct {
	mock *MockMongoSingleResult
}

// NewMockMongoSingleResult creates a new mock instance
func NewMockMongoSingleResult(ctrl *gomock.Controller) *MockMongoSingleResult {
	mock := &MockMongoSingleResult{ctrl: ctrl}
	mock.recorder = &MockMongoSingleResultMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMongoSingleResult) EXPECT() *MockMongoSingleResultMockRecorder {
	return m.recorder
}

// Err mocks base method
func (m *MockMongoSingleResult) Err() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Err")
	ret0, _ := ret[0].(error)
	return ret0
}

// Err indicates an expected call of Err
func (mr *MockMongoSingleResultMockRecorder) Err() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Err", reflect.TypeOf((*MockMongoSingleResult)(nil).Err))
}

// Decode mocks base method
func (m *MockMongoSingleResult) Decode(v interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Decode", v)
	ret0, _ := ret[0].(error)
	return ret0
}

// Decode indicates an expected call of Decode
func (mr *MockMongoSingleResultMockRecorder) Decode(v interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Decode", reflect.TypeOf((*MockMongoSingleResult)(nil).Decode), v)
}

// MockMongoClient is a mock of MongoClient interface
type MockMongoClient struct {
	ctrl     *gomock.Controller
	recorder *MockMongoClientMockRecorder
}

// MockMongoClientMockRecorder is the mock recorder for MockMongoClient
type MockMongoClientMockRecorder struct {
	mock *MockMongoClient
}

// NewMockMongoClient creates a new mock instance
func NewMockMongoClient(ctrl *gomock.Controller) *MockMongoClient {
	mock := &MockMongoClient{ctrl: ctrl}
	mock.recorder = &MockMongoClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMongoClient) EXPECT() *MockMongoClientMockRecorder {
	return m.recorder
}
