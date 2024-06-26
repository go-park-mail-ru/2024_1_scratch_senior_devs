// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces.go

// Package mock_metrics is a generated GoMock package.
package mock_metrics

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockDBMetrics is a mock of DBMetrics interface.
type MockDBMetrics struct {
	ctrl     *gomock.Controller
	recorder *MockDBMetricsMockRecorder
}

// MockDBMetricsMockRecorder is the mock recorder for MockDBMetrics.
type MockDBMetricsMockRecorder struct {
	mock *MockDBMetrics
}

// NewMockDBMetrics creates a new mock instance.
func NewMockDBMetrics(ctrl *gomock.Controller) *MockDBMetrics {
	mock := &MockDBMetrics{ctrl: ctrl}
	mock.recorder = &MockDBMetricsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDBMetrics) EXPECT() *MockDBMetricsMockRecorder {
	return m.recorder
}

// IncreaseErrors mocks base method.
func (m *MockDBMetrics) IncreaseErrors(queryName string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "IncreaseErrors", queryName)
}

// IncreaseErrors indicates an expected call of IncreaseErrors.
func (mr *MockDBMetricsMockRecorder) IncreaseErrors(queryName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IncreaseErrors", reflect.TypeOf((*MockDBMetrics)(nil).IncreaseErrors), queryName)
}

// ObserveResponseTime mocks base method.
func (m *MockDBMetrics) ObserveResponseTime(queryName string, observeTime float64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ObserveResponseTime", queryName, observeTime)
}

// ObserveResponseTime indicates an expected call of ObserveResponseTime.
func (mr *MockDBMetricsMockRecorder) ObserveResponseTime(queryName, observeTime interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ObserveResponseTime", reflect.TypeOf((*MockDBMetrics)(nil).ObserveResponseTime), queryName, observeTime)
}

// MockWSMetrics is a mock of WSMetrics interface.
type MockWSMetrics struct {
	ctrl     *gomock.Controller
	recorder *MockWSMetricsMockRecorder
}

// MockWSMetricsMockRecorder is the mock recorder for MockWSMetrics.
type MockWSMetricsMockRecorder struct {
	mock *MockWSMetrics
}

// NewMockWSMetrics creates a new mock instance.
func NewMockWSMetrics(ctrl *gomock.Controller) *MockWSMetrics {
	mock := &MockWSMetrics{ctrl: ctrl}
	mock.recorder = &MockWSMetricsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWSMetrics) EXPECT() *MockWSMetricsMockRecorder {
	return m.recorder
}

// DecreaseConnections mocks base method.
func (m *MockWSMetrics) DecreaseConnections() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "DecreaseConnections")
}

// DecreaseConnections indicates an expected call of DecreaseConnections.
func (mr *MockWSMetricsMockRecorder) DecreaseConnections() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DecreaseConnections", reflect.TypeOf((*MockWSMetrics)(nil).DecreaseConnections))
}

// IncreaseConnections mocks base method.
func (m *MockWSMetrics) IncreaseConnections() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "IncreaseConnections")
}

// IncreaseConnections indicates an expected call of IncreaseConnections.
func (mr *MockWSMetricsMockRecorder) IncreaseConnections() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IncreaseConnections", reflect.TypeOf((*MockWSMetrics)(nil).IncreaseConnections))
}
