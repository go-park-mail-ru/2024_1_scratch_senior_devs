// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces.go

// Package mock_note is a generated GoMock package.
package mock_note

import (
	context "context"
	reflect "reflect"

	models "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/models"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/satori/uuid"
)

// MockNoteUsecase is a mock of NoteUsecase interface.
type MockNoteUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockNoteUsecaseMockRecorder
}

// MockNoteUsecaseMockRecorder is the mock recorder for MockNoteUsecase.
type MockNoteUsecaseMockRecorder struct {
	mock *MockNoteUsecase
}

// NewMockNoteUsecase creates a new mock instance.
func NewMockNoteUsecase(ctrl *gomock.Controller) *MockNoteUsecase {
	mock := &MockNoteUsecase{ctrl: ctrl}
	mock.recorder = &MockNoteUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNoteUsecase) EXPECT() *MockNoteUsecaseMockRecorder {
	return m.recorder
}

// CreateNote mocks base method.
func (m *MockNoteUsecase) CreateNote(arg0 context.Context, arg1 uuid.UUID, arg2 []byte) (models.Note, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateNote", arg0, arg1, arg2)
	ret0, _ := ret[0].(models.Note)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateNote indicates an expected call of CreateNote.
func (mr *MockNoteUsecaseMockRecorder) CreateNote(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateNote", reflect.TypeOf((*MockNoteUsecase)(nil).CreateNote), arg0, arg1, arg2)
}

// DeleteNote mocks base method.
func (m *MockNoteUsecase) DeleteNote(arg0 context.Context, arg1, arg2 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteNote", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteNote indicates an expected call of DeleteNote.
func (mr *MockNoteUsecaseMockRecorder) DeleteNote(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteNote", reflect.TypeOf((*MockNoteUsecase)(nil).DeleteNote), arg0, arg1, arg2)
}

// GetAllNotes mocks base method.
func (m *MockNoteUsecase) GetAllNotes(arg0 context.Context, arg1 uuid.UUID, arg2, arg3 int64, arg4 string) ([]models.Note, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllNotes", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].([]models.Note)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllNotes indicates an expected call of GetAllNotes.
func (mr *MockNoteUsecaseMockRecorder) GetAllNotes(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllNotes", reflect.TypeOf((*MockNoteUsecase)(nil).GetAllNotes), arg0, arg1, arg2, arg3, arg4)
}

// GetNote mocks base method.
func (m *MockNoteUsecase) GetNote(arg0 context.Context, arg1, arg2 uuid.UUID) (models.Note, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNote", arg0, arg1, arg2)
	ret0, _ := ret[0].(models.Note)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNote indicates an expected call of GetNote.
func (mr *MockNoteUsecaseMockRecorder) GetNote(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNote", reflect.TypeOf((*MockNoteUsecase)(nil).GetNote), arg0, arg1, arg2)
}

// UpdateNote mocks base method.
func (m *MockNoteUsecase) UpdateNote(arg0 context.Context, arg1, arg2 uuid.UUID, arg3 []byte) (models.Note, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateNote", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(models.Note)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateNote indicates an expected call of UpdateNote.
func (mr *MockNoteUsecaseMockRecorder) UpdateNote(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateNote", reflect.TypeOf((*MockNoteUsecase)(nil).UpdateNote), arg0, arg1, arg2, arg3)
}

// MockNoteBaseRepo is a mock of NoteBaseRepo interface.
type MockNoteBaseRepo struct {
	ctrl     *gomock.Controller
	recorder *MockNoteBaseRepoMockRecorder
}

// MockNoteBaseRepoMockRecorder is the mock recorder for MockNoteBaseRepo.
type MockNoteBaseRepoMockRecorder struct {
	mock *MockNoteBaseRepo
}

// NewMockNoteBaseRepo creates a new mock instance.
func NewMockNoteBaseRepo(ctrl *gomock.Controller) *MockNoteBaseRepo {
	mock := &MockNoteBaseRepo{ctrl: ctrl}
	mock.recorder = &MockNoteBaseRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNoteBaseRepo) EXPECT() *MockNoteBaseRepoMockRecorder {
	return m.recorder
}

// CreateNote mocks base method.
func (m *MockNoteBaseRepo) CreateNote(arg0 context.Context, arg1 models.Note) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateNote", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateNote indicates an expected call of CreateNote.
func (mr *MockNoteBaseRepoMockRecorder) CreateNote(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateNote", reflect.TypeOf((*MockNoteBaseRepo)(nil).CreateNote), arg0, arg1)
}

// DeleteNote mocks base method.
func (m *MockNoteBaseRepo) DeleteNote(arg0 context.Context, arg1 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteNote", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteNote indicates an expected call of DeleteNote.
func (mr *MockNoteBaseRepoMockRecorder) DeleteNote(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteNote", reflect.TypeOf((*MockNoteBaseRepo)(nil).DeleteNote), arg0, arg1)
}

// ReadAllNotes mocks base method.
func (m *MockNoteBaseRepo) ReadAllNotes(arg0 context.Context, arg1 uuid.UUID, arg2, arg3 int64) ([]models.Note, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadAllNotes", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].([]models.Note)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadAllNotes indicates an expected call of ReadAllNotes.
func (mr *MockNoteBaseRepoMockRecorder) ReadAllNotes(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadAllNotes", reflect.TypeOf((*MockNoteBaseRepo)(nil).ReadAllNotes), arg0, arg1, arg2, arg3)
}

// ReadNote mocks base method.
func (m *MockNoteBaseRepo) ReadNote(arg0 context.Context, arg1 uuid.UUID) (models.Note, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadNote", arg0, arg1)
	ret0, _ := ret[0].(models.Note)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadNote indicates an expected call of ReadNote.
func (mr *MockNoteBaseRepoMockRecorder) ReadNote(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadNote", reflect.TypeOf((*MockNoteBaseRepo)(nil).ReadNote), arg0, arg1)
}

// UpdateNote mocks base method.
func (m *MockNoteBaseRepo) UpdateNote(arg0 context.Context, arg1 models.Note) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateNote", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateNote indicates an expected call of UpdateNote.
func (mr *MockNoteBaseRepoMockRecorder) UpdateNote(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateNote", reflect.TypeOf((*MockNoteBaseRepo)(nil).UpdateNote), arg0, arg1)
}

// MockNoteSearchRepo is a mock of NoteSearchRepo interface.
type MockNoteSearchRepo struct {
	ctrl     *gomock.Controller
	recorder *MockNoteSearchRepoMockRecorder
}

// MockNoteSearchRepoMockRecorder is the mock recorder for MockNoteSearchRepo.
type MockNoteSearchRepoMockRecorder struct {
	mock *MockNoteSearchRepo
}

// NewMockNoteSearchRepo creates a new mock instance.
func NewMockNoteSearchRepo(ctrl *gomock.Controller) *MockNoteSearchRepo {
	mock := &MockNoteSearchRepo{ctrl: ctrl}
	mock.recorder = &MockNoteSearchRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNoteSearchRepo) EXPECT() *MockNoteSearchRepoMockRecorder {
	return m.recorder
}

// CreateNote mocks base method.
func (m *MockNoteSearchRepo) CreateNote(arg0 context.Context, arg1 models.Note) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateNote", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateNote indicates an expected call of CreateNote.
func (mr *MockNoteSearchRepoMockRecorder) CreateNote(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateNote", reflect.TypeOf((*MockNoteSearchRepo)(nil).CreateNote), arg0, arg1)
}

// DeleteNote mocks base method.
func (m *MockNoteSearchRepo) DeleteNote(arg0 context.Context, arg1 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteNote", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteNote indicates an expected call of DeleteNote.
func (mr *MockNoteSearchRepoMockRecorder) DeleteNote(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteNote", reflect.TypeOf((*MockNoteSearchRepo)(nil).DeleteNote), arg0, arg1)
}

// SearchNotes mocks base method.
func (m *MockNoteSearchRepo) SearchNotes(arg0 context.Context, arg1 uuid.UUID, arg2, arg3 int64, arg4 string) ([]models.Note, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchNotes", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].([]models.Note)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchNotes indicates an expected call of SearchNotes.
func (mr *MockNoteSearchRepoMockRecorder) SearchNotes(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchNotes", reflect.TypeOf((*MockNoteSearchRepo)(nil).SearchNotes), arg0, arg1, arg2, arg3, arg4)
}

// UpdateNote mocks base method.
func (m *MockNoteSearchRepo) UpdateNote(arg0 context.Context, arg1 models.Note) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateNote", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateNote indicates an expected call of UpdateNote.
func (mr *MockNoteSearchRepoMockRecorder) UpdateNote(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateNote", reflect.TypeOf((*MockNoteSearchRepo)(nil).UpdateNote), arg0, arg1)
}
