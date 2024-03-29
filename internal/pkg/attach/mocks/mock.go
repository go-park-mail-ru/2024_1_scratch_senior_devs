// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces.go

// Package mock_attach is a generated GoMock package.
package mock_attach

import (
	context "context"
	io "io"
	reflect "reflect"

	models "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/satori/uuid"
)

// MockAttachUsecase is a mock of AttachUsecase interface.
type MockAttachUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockAttachUsecaseMockRecorder
}

// MockAttachUsecaseMockRecorder is the mock recorder for MockAttachUsecase.
type MockAttachUsecaseMockRecorder struct {
	mock *MockAttachUsecase
}

// NewMockAttachUsecase creates a new mock instance.
func NewMockAttachUsecase(ctrl *gomock.Controller) *MockAttachUsecase {
	mock := &MockAttachUsecase{ctrl: ctrl}
	mock.recorder = &MockAttachUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAttachUsecase) EXPECT() *MockAttachUsecaseMockRecorder {
	return m.recorder
}

// AddAttach mocks base method.
func (m *MockAttachUsecase) AddAttach(ctx context.Context, noteID uuid.UUID, attach io.ReadSeeker, extension string) (models.Attach, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddAttach", ctx, noteID, attach, extension)
	ret0, _ := ret[0].(models.Attach)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddAttach indicates an expected call of AddAttach.
func (mr *MockAttachUsecaseMockRecorder) AddAttach(ctx, noteID, attach, extension interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddAttach", reflect.TypeOf((*MockAttachUsecase)(nil).AddAttach), ctx, noteID, attach, extension)
}

// MockAttachRepo is a mock of AttachRepo interface.
type MockAttachRepo struct {
	ctrl     *gomock.Controller
	recorder *MockAttachRepoMockRecorder
}

// MockAttachRepoMockRecorder is the mock recorder for MockAttachRepo.
type MockAttachRepoMockRecorder struct {
	mock *MockAttachRepo
}

// NewMockAttachRepo creates a new mock instance.
func NewMockAttachRepo(ctrl *gomock.Controller) *MockAttachRepo {
	mock := &MockAttachRepo{ctrl: ctrl}
	mock.recorder = &MockAttachRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAttachRepo) EXPECT() *MockAttachRepoMockRecorder {
	return m.recorder
}

// AddAttach mocks base method.
func (m *MockAttachRepo) AddAttach(ctx context.Context, attach models.Attach) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddAttach", ctx, attach)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddAttach indicates an expected call of AddAttach.
func (mr *MockAttachRepoMockRecorder) AddAttach(ctx, attach interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddAttach", reflect.TypeOf((*MockAttachRepo)(nil).AddAttach), ctx, attach)
}
