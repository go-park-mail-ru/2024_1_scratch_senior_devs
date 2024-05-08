// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces.go

// Package mock_note is a generated GoMock package.
package mock_note

import (
	context "context"
	reflect "reflect"
	time "time"

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

// AddCollaborator mocks base method.
func (m *MockNoteUsecase) AddCollaborator(arg0 context.Context, arg1, arg2, arg3 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddCollaborator", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddCollaborator indicates an expected call of AddCollaborator.
func (mr *MockNoteUsecaseMockRecorder) AddCollaborator(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCollaborator", reflect.TypeOf((*MockNoteUsecase)(nil).AddCollaborator), arg0, arg1, arg2, arg3)
}

// AddTag mocks base method.
func (m *MockNoteUsecase) AddTag(ctx context.Context, tagName string, noteId, userId uuid.UUID) (models.Note, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddTag", ctx, tagName, noteId, userId)
	ret0, _ := ret[0].(models.Note)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddTag indicates an expected call of AddTag.
func (mr *MockNoteUsecaseMockRecorder) AddTag(ctx, tagName, noteId, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddTag", reflect.TypeOf((*MockNoteUsecase)(nil).AddTag), ctx, tagName, noteId, userId)
}

// CheckPermissions mocks base method.
func (m *MockNoteUsecase) CheckPermissions(ctx context.Context, noteID, userID uuid.UUID) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckPermissions", ctx, noteID, userID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckPermissions indicates an expected call of CheckPermissions.
func (mr *MockNoteUsecaseMockRecorder) CheckPermissions(ctx, noteID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckPermissions", reflect.TypeOf((*MockNoteUsecase)(nil).CheckPermissions), ctx, noteID, userID)
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

// CreateSubNote mocks base method.
func (m *MockNoteUsecase) CreateSubNote(arg0 context.Context, arg1 uuid.UUID, arg2 []byte, arg3 uuid.UUID) (models.Note, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSubNote", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(models.Note)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSubNote indicates an expected call of CreateSubNote.
func (mr *MockNoteUsecaseMockRecorder) CreateSubNote(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSubNote", reflect.TypeOf((*MockNoteUsecase)(nil).CreateSubNote), arg0, arg1, arg2, arg3)
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

// DeleteTag mocks base method.
func (m *MockNoteUsecase) DeleteTag(ctx context.Context, tagName string, noteId, userId uuid.UUID) (models.Note, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTag", ctx, tagName, noteId, userId)
	ret0, _ := ret[0].(models.Note)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteTag indicates an expected call of DeleteTag.
func (mr *MockNoteUsecaseMockRecorder) DeleteTag(ctx, tagName, noteId, userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTag", reflect.TypeOf((*MockNoteUsecase)(nil).DeleteTag), ctx, tagName, noteId, userId)
}

// GetAllNotes mocks base method.
func (m *MockNoteUsecase) GetAllNotes(arg0 context.Context, arg1 uuid.UUID, arg2, arg3 int64, arg4 string, arg5 []string) ([]models.Note, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllNotes", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].([]models.Note)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllNotes indicates an expected call of GetAllNotes.
func (mr *MockNoteUsecaseMockRecorder) GetAllNotes(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllNotes", reflect.TypeOf((*MockNoteUsecase)(nil).GetAllNotes), arg0, arg1, arg2, arg3, arg4, arg5)
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

// GetTags mocks base method.
func (m *MockNoteUsecase) GetTags(ctx context.Context, userID uuid.UUID) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTags", ctx, userID)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTags indicates an expected call of GetTags.
func (mr *MockNoteUsecaseMockRecorder) GetTags(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTags", reflect.TypeOf((*MockNoteUsecase)(nil).GetTags), ctx, userID)
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

// AddCollaborator mocks base method.
func (m *MockNoteBaseRepo) AddCollaborator(arg0 context.Context, arg1, arg2 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddCollaborator", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddCollaborator indicates an expected call of AddCollaborator.
func (mr *MockNoteBaseRepoMockRecorder) AddCollaborator(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCollaborator", reflect.TypeOf((*MockNoteBaseRepo)(nil).AddCollaborator), arg0, arg1, arg2)
}

// AddSubNote mocks base method.
func (m *MockNoteBaseRepo) AddSubNote(arg0 context.Context, arg1, arg2 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddSubNote", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddSubNote indicates an expected call of AddSubNote.
func (mr *MockNoteBaseRepoMockRecorder) AddSubNote(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddSubNote", reflect.TypeOf((*MockNoteBaseRepo)(nil).AddSubNote), arg0, arg1, arg2)
}

// AddTag mocks base method.
func (m *MockNoteBaseRepo) AddTag(ctx context.Context, tagName string, noteId uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddTag", ctx, tagName, noteId)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddTag indicates an expected call of AddTag.
func (mr *MockNoteBaseRepoMockRecorder) AddTag(ctx, tagName, noteId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddTag", reflect.TypeOf((*MockNoteBaseRepo)(nil).AddTag), ctx, tagName, noteId)
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

// DeleteTag mocks base method.
func (m *MockNoteBaseRepo) DeleteTag(ctx context.Context, tagName string, noteId uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTag", ctx, tagName, noteId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTag indicates an expected call of DeleteTag.
func (mr *MockNoteBaseRepoMockRecorder) DeleteTag(ctx, tagName, noteId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTag", reflect.TypeOf((*MockNoteBaseRepo)(nil).DeleteTag), ctx, tagName, noteId)
}

// GetTags mocks base method.
func (m *MockNoteBaseRepo) GetTags(ctx context.Context, userID uuid.UUID) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTags", ctx, userID)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTags indicates an expected call of GetTags.
func (mr *MockNoteBaseRepoMockRecorder) GetTags(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTags", reflect.TypeOf((*MockNoteBaseRepo)(nil).GetTags), ctx, userID)
}

// GetUpdates mocks base method.
func (m *MockNoteBaseRepo) GetUpdates(arg0 context.Context, arg1 uuid.UUID, arg2 time.Time) ([]models.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUpdates", arg0, arg1, arg2)
	ret0, _ := ret[0].([]models.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUpdates indicates an expected call of GetUpdates.
func (mr *MockNoteBaseRepoMockRecorder) GetUpdates(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUpdates", reflect.TypeOf((*MockNoteBaseRepo)(nil).GetUpdates), arg0, arg1, arg2)
}

// ReadAllNotes mocks base method.
func (m *MockNoteBaseRepo) ReadAllNotes(arg0 context.Context, arg1 uuid.UUID, arg2, arg3 int64, arg4 []string) ([]models.Note, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadAllNotes", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].([]models.Note)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadAllNotes indicates an expected call of ReadAllNotes.
func (mr *MockNoteBaseRepoMockRecorder) ReadAllNotes(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadAllNotes", reflect.TypeOf((*MockNoteBaseRepo)(nil).ReadAllNotes), arg0, arg1, arg2, arg3, arg4)
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

// RemoveSubNote mocks base method.
func (m *MockNoteBaseRepo) RemoveSubNote(arg0 context.Context, arg1, arg2 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveSubNote", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveSubNote indicates an expected call of RemoveSubNote.
func (mr *MockNoteBaseRepoMockRecorder) RemoveSubNote(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveSubNote", reflect.TypeOf((*MockNoteBaseRepo)(nil).RemoveSubNote), arg0, arg1, arg2)
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

// AddCollaborator mocks base method.
func (m *MockNoteSearchRepo) AddCollaborator(arg0 context.Context, arg1, arg2 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddCollaborator", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddCollaborator indicates an expected call of AddCollaborator.
func (mr *MockNoteSearchRepoMockRecorder) AddCollaborator(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCollaborator", reflect.TypeOf((*MockNoteSearchRepo)(nil).AddCollaborator), arg0, arg1, arg2)
}

// AddSubNote mocks base method.
func (m *MockNoteSearchRepo) AddSubNote(arg0 context.Context, arg1, arg2 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddSubNote", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddSubNote indicates an expected call of AddSubNote.
func (mr *MockNoteSearchRepoMockRecorder) AddSubNote(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddSubNote", reflect.TypeOf((*MockNoteSearchRepo)(nil).AddSubNote), arg0, arg1, arg2)
}

// AddTag mocks base method.
func (m *MockNoteSearchRepo) AddTag(ctx context.Context, tagName string, noteID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddTag", ctx, tagName, noteID)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddTag indicates an expected call of AddTag.
func (mr *MockNoteSearchRepoMockRecorder) AddTag(ctx, tagName, noteID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddTag", reflect.TypeOf((*MockNoteSearchRepo)(nil).AddTag), ctx, tagName, noteID)
}

// CreateNote mocks base method.
func (m *MockNoteSearchRepo) CreateNote(arg0 context.Context, arg1 models.ElasticNote) error {
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

// DeleteTag mocks base method.
func (m *MockNoteSearchRepo) DeleteTag(ctx context.Context, tagName string, noteID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTag", ctx, tagName, noteID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTag indicates an expected call of DeleteTag.
func (mr *MockNoteSearchRepoMockRecorder) DeleteTag(ctx, tagName, noteID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTag", reflect.TypeOf((*MockNoteSearchRepo)(nil).DeleteTag), ctx, tagName, noteID)
}

// RemoveSubNote mocks base method.
func (m *MockNoteSearchRepo) RemoveSubNote(arg0 context.Context, arg1, arg2 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveSubNote", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveSubNote indicates an expected call of RemoveSubNote.
func (mr *MockNoteSearchRepoMockRecorder) RemoveSubNote(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveSubNote", reflect.TypeOf((*MockNoteSearchRepo)(nil).RemoveSubNote), arg0, arg1, arg2)
}

// SearchNotes mocks base method.
func (m *MockNoteSearchRepo) SearchNotes(arg0 context.Context, arg1 uuid.UUID, arg2, arg3 int64, arg4 string, arg5 []string) ([]models.Note, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchNotes", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].([]models.Note)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchNotes indicates an expected call of SearchNotes.
func (mr *MockNoteSearchRepoMockRecorder) SearchNotes(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchNotes", reflect.TypeOf((*MockNoteSearchRepo)(nil).SearchNotes), arg0, arg1, arg2, arg3, arg4, arg5)
}

// UpdateNote mocks base method.
func (m *MockNoteSearchRepo) UpdateNote(arg0 context.Context, arg1 models.ElasticNote) error {
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
