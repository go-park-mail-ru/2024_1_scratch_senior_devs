// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces.go

// Package mock_hub is a generated GoMock package.
package mock_hub

import (
	context "context"
	reflect "reflect"

	models "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	hub "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/hub"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/satori/uuid"
)

// MockHubInterface is a mock of HubInterface interface.
type MockHubInterface struct {
	ctrl     *gomock.Controller
	recorder *MockHubInterfaceMockRecorder
}

// MockHubInterfaceMockRecorder is the mock recorder for MockHubInterface.
type MockHubInterfaceMockRecorder struct {
	mock *MockHubInterface
}

// NewMockHubInterface creates a new mock instance.
func NewMockHubInterface(ctrl *gomock.Controller) *MockHubInterface {
	mock := &MockHubInterface{ctrl: ctrl}
	mock.recorder = &MockHubInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHubInterface) EXPECT() *MockHubInterfaceMockRecorder {
	return m.recorder
}

// AddClient mocks base method.
func (m *MockHubInterface) AddClient(arg0 context.Context, arg1 uuid.UUID, arg2 *hub.CustomClient) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddClient", arg0, arg1, arg2)
}

// AddClient indicates an expected call of AddClient.
func (mr *MockHubInterfaceMockRecorder) AddClient(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddClient", reflect.TypeOf((*MockHubInterface)(nil).AddClient), arg0, arg1, arg2)
}

// AddClientMain mocks base method.
func (m *MockHubInterface) AddClientMain(arg0 context.Context, arg1 uuid.UUID, arg2 *hub.CustomClient) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddClientMain", arg0, arg1, arg2)
}

// AddClientMain indicates an expected call of AddClientMain.
func (mr *MockHubInterfaceMockRecorder) AddClientMain(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddClientMain", reflect.TypeOf((*MockHubInterface)(nil).AddClientMain), arg0, arg1, arg2)
}

// Run mocks base method.
func (m *MockHubInterface) Run(arg0 context.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Run", arg0)
}

// Run indicates an expected call of Run.
func (mr *MockHubInterfaceMockRecorder) Run(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockHubInterface)(nil).Run), arg0)
}

// StartCache mocks base method.
func (m *MockHubInterface) StartCache(arg0 context.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "StartCache", arg0)
}

// StartCache indicates an expected call of StartCache.
func (mr *MockHubInterfaceMockRecorder) StartCache(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartCache", reflect.TypeOf((*MockHubInterface)(nil).StartCache), arg0)
}

// StartCacheMain mocks base method.
func (m *MockHubInterface) StartCacheMain(ctx context.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "StartCacheMain", ctx)
}

// StartCacheMain indicates an expected call of StartCacheMain.
func (mr *MockHubInterfaceMockRecorder) StartCacheMain(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartCacheMain", reflect.TypeOf((*MockHubInterface)(nil).StartCacheMain), ctx)
}

// WriteToCache mocks base method.
func (m *MockHubInterface) WriteToCache(arg0 context.Context, arg1 models.CacheMessage) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WriteToCache", arg0, arg1)
}

// WriteToCache indicates an expected call of WriteToCache.
func (mr *MockHubInterfaceMockRecorder) WriteToCache(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteToCache", reflect.TypeOf((*MockHubInterface)(nil).WriteToCache), arg0, arg1)
}

// WriteToCacheMain mocks base method.
func (m *MockHubInterface) WriteToCacheMain(arg0 context.Context, arg1 uuid.UUID, arg2 models.InviteMessage) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WriteToCacheMain", arg0, arg1, arg2)
}

// WriteToCacheMain indicates an expected call of WriteToCacheMain.
func (mr *MockHubInterfaceMockRecorder) WriteToCacheMain(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteToCacheMain", reflect.TypeOf((*MockHubInterface)(nil).WriteToCacheMain), arg0, arg1, arg2)
}
