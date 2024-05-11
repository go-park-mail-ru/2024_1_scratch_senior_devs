// Code generated by MockGen. DO NOT EDIT.
// Source: note_grpc.pb.go

// Package mock_gen is a generated GoMock package.
package mock_gen

import (
	context "context"
	reflect "reflect"

	gen "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/note/delivery/grpc/gen"
	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockNoteClient is a mock of NoteClient interface.
type MockNoteClient struct {
	ctrl     *gomock.Controller
	recorder *MockNoteClientMockRecorder
}

// MockNoteClientMockRecorder is the mock recorder for MockNoteClient.
type MockNoteClientMockRecorder struct {
	mock *MockNoteClient
}

// NewMockNoteClient creates a new mock instance.
func NewMockNoteClient(ctrl *gomock.Controller) *MockNoteClient {
	mock := &MockNoteClient{ctrl: ctrl}
	mock.recorder = &MockNoteClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNoteClient) EXPECT() *MockNoteClientMockRecorder {
	return m.recorder
}

// AddCollaborator mocks base method.
func (m *MockNoteClient) AddCollaborator(ctx context.Context, in *gen.AddCollaboratorRequest, opts ...grpc.CallOption) (*gen.AddCollaboratorResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AddCollaborator", varargs...)
	ret0, _ := ret[0].(*gen.AddCollaboratorResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddCollaborator indicates an expected call of AddCollaborator.
func (mr *MockNoteClientMockRecorder) AddCollaborator(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCollaborator", reflect.TypeOf((*MockNoteClient)(nil).AddCollaborator), varargs...)
}

// AddNote mocks base method.
func (m *MockNoteClient) AddNote(ctx context.Context, in *gen.AddNoteRequest, opts ...grpc.CallOption) (*gen.AddNoteResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AddNote", varargs...)
	ret0, _ := ret[0].(*gen.AddNoteResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddNote indicates an expected call of AddNote.
func (mr *MockNoteClientMockRecorder) AddNote(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddNote", reflect.TypeOf((*MockNoteClient)(nil).AddNote), varargs...)
}

// AddTag mocks base method.
func (m *MockNoteClient) AddTag(ctx context.Context, in *gen.TagRequest, opts ...grpc.CallOption) (*gen.GetNoteResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AddTag", varargs...)
	ret0, _ := ret[0].(*gen.GetNoteResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddTag indicates an expected call of AddTag.
func (mr *MockNoteClientMockRecorder) AddTag(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddTag", reflect.TypeOf((*MockNoteClient)(nil).AddTag), varargs...)
}

// CheckPermissions mocks base method.
func (m *MockNoteClient) CheckPermissions(ctx context.Context, in *gen.CheckPermissionsRequest, opts ...grpc.CallOption) (*gen.CheckPermissionsResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CheckPermissions", varargs...)
	ret0, _ := ret[0].(*gen.CheckPermissionsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckPermissions indicates an expected call of CheckPermissions.
func (mr *MockNoteClientMockRecorder) CheckPermissions(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckPermissions", reflect.TypeOf((*MockNoteClient)(nil).CheckPermissions), varargs...)
}

// CreateSubNote mocks base method.
func (m *MockNoteClient) CreateSubNote(ctx context.Context, in *gen.CreateSubNoteRequest, opts ...grpc.CallOption) (*gen.CreateSubNoteResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateSubNote", varargs...)
	ret0, _ := ret[0].(*gen.CreateSubNoteResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSubNote indicates an expected call of CreateSubNote.
func (mr *MockNoteClientMockRecorder) CreateSubNote(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSubNote", reflect.TypeOf((*MockNoteClient)(nil).CreateSubNote), varargs...)
}

// DeleteNote mocks base method.
func (m *MockNoteClient) DeleteNote(ctx context.Context, in *gen.DeleteNoteRequest, opts ...grpc.CallOption) (*gen.DeleteNoteResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteNote", varargs...)
	ret0, _ := ret[0].(*gen.DeleteNoteResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteNote indicates an expected call of DeleteNote.
func (mr *MockNoteClientMockRecorder) DeleteNote(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteNote", reflect.TypeOf((*MockNoteClient)(nil).DeleteNote), varargs...)
}

// DeleteTag mocks base method.
func (m *MockNoteClient) DeleteTag(ctx context.Context, in *gen.TagRequest, opts ...grpc.CallOption) (*gen.GetNoteResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteTag", varargs...)
	ret0, _ := ret[0].(*gen.GetNoteResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteTag indicates an expected call of DeleteTag.
func (mr *MockNoteClientMockRecorder) DeleteTag(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTag", reflect.TypeOf((*MockNoteClient)(nil).DeleteTag), varargs...)
}

// ForgetTag mocks base method.
func (m *MockNoteClient) ForgetTag(ctx context.Context, in *gen.AllTagRequest, opts ...grpc.CallOption) (*gen.EmptyResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ForgetTag", varargs...)
	ret0, _ := ret[0].(*gen.EmptyResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ForgetTag indicates an expected call of ForgetTag.
func (mr *MockNoteClientMockRecorder) ForgetTag(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForgetTag", reflect.TypeOf((*MockNoteClient)(nil).ForgetTag), varargs...)
}

// GetAllNotes mocks base method.
func (m *MockNoteClient) GetAllNotes(ctx context.Context, in *gen.GetAllRequest, opts ...grpc.CallOption) (*gen.GetAllResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetAllNotes", varargs...)
	ret0, _ := ret[0].(*gen.GetAllResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllNotes indicates an expected call of GetAllNotes.
func (mr *MockNoteClientMockRecorder) GetAllNotes(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllNotes", reflect.TypeOf((*MockNoteClient)(nil).GetAllNotes), varargs...)
}

// GetNote mocks base method.
func (m *MockNoteClient) GetNote(ctx context.Context, in *gen.GetNoteRequest, opts ...grpc.CallOption) (*gen.GetNoteResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetNote", varargs...)
	ret0, _ := ret[0].(*gen.GetNoteResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNote indicates an expected call of GetNote.
func (mr *MockNoteClientMockRecorder) GetNote(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNote", reflect.TypeOf((*MockNoteClient)(nil).GetNote), varargs...)
}

// GetTags mocks base method.
func (m *MockNoteClient) GetTags(ctx context.Context, in *gen.GetTagsRequest, opts ...grpc.CallOption) (*gen.GetTagsResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetTags", varargs...)
	ret0, _ := ret[0].(*gen.GetTagsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTags indicates an expected call of GetTags.
func (mr *MockNoteClientMockRecorder) GetTags(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTags", reflect.TypeOf((*MockNoteClient)(nil).GetTags), varargs...)
}

// RememberTag mocks base method.
func (m *MockNoteClient) RememberTag(ctx context.Context, in *gen.AllTagRequest, opts ...grpc.CallOption) (*gen.EmptyResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "RememberTag", varargs...)
	ret0, _ := ret[0].(*gen.EmptyResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RememberTag indicates an expected call of RememberTag.
func (mr *MockNoteClientMockRecorder) RememberTag(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RememberTag", reflect.TypeOf((*MockNoteClient)(nil).RememberTag), varargs...)
}

// UpdateNote mocks base method.
func (m *MockNoteClient) UpdateNote(ctx context.Context, in *gen.UpdateNoteRequest, opts ...grpc.CallOption) (*gen.UpdateNoteResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateNote", varargs...)
	ret0, _ := ret[0].(*gen.UpdateNoteResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateNote indicates an expected call of UpdateNote.
func (mr *MockNoteClientMockRecorder) UpdateNote(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateNote", reflect.TypeOf((*MockNoteClient)(nil).UpdateNote), varargs...)
}

// MockNoteServer is a mock of NoteServer interface.
type MockNoteServer struct {
	ctrl     *gomock.Controller
	recorder *MockNoteServerMockRecorder
}

// MockNoteServerMockRecorder is the mock recorder for MockNoteServer.
type MockNoteServerMockRecorder struct {
	mock *MockNoteServer
}

// NewMockNoteServer creates a new mock instance.
func NewMockNoteServer(ctrl *gomock.Controller) *MockNoteServer {
	mock := &MockNoteServer{ctrl: ctrl}
	mock.recorder = &MockNoteServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNoteServer) EXPECT() *MockNoteServerMockRecorder {
	return m.recorder
}

// AddCollaborator mocks base method.
func (m *MockNoteServer) AddCollaborator(arg0 context.Context, arg1 *gen.AddCollaboratorRequest) (*gen.AddCollaboratorResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddCollaborator", arg0, arg1)
	ret0, _ := ret[0].(*gen.AddCollaboratorResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddCollaborator indicates an expected call of AddCollaborator.
func (mr *MockNoteServerMockRecorder) AddCollaborator(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCollaborator", reflect.TypeOf((*MockNoteServer)(nil).AddCollaborator), arg0, arg1)
}

// AddNote mocks base method.
func (m *MockNoteServer) AddNote(arg0 context.Context, arg1 *gen.AddNoteRequest) (*gen.AddNoteResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddNote", arg0, arg1)
	ret0, _ := ret[0].(*gen.AddNoteResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddNote indicates an expected call of AddNote.
func (mr *MockNoteServerMockRecorder) AddNote(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddNote", reflect.TypeOf((*MockNoteServer)(nil).AddNote), arg0, arg1)
}

// AddTag mocks base method.
func (m *MockNoteServer) AddTag(arg0 context.Context, arg1 *gen.TagRequest) (*gen.GetNoteResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddTag", arg0, arg1)
	ret0, _ := ret[0].(*gen.GetNoteResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddTag indicates an expected call of AddTag.
func (mr *MockNoteServerMockRecorder) AddTag(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddTag", reflect.TypeOf((*MockNoteServer)(nil).AddTag), arg0, arg1)
}

// CheckPermissions mocks base method.
func (m *MockNoteServer) CheckPermissions(arg0 context.Context, arg1 *gen.CheckPermissionsRequest) (*gen.CheckPermissionsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckPermissions", arg0, arg1)
	ret0, _ := ret[0].(*gen.CheckPermissionsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckPermissions indicates an expected call of CheckPermissions.
func (mr *MockNoteServerMockRecorder) CheckPermissions(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckPermissions", reflect.TypeOf((*MockNoteServer)(nil).CheckPermissions), arg0, arg1)
}

// CreateSubNote mocks base method.
func (m *MockNoteServer) CreateSubNote(arg0 context.Context, arg1 *gen.CreateSubNoteRequest) (*gen.CreateSubNoteResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSubNote", arg0, arg1)
	ret0, _ := ret[0].(*gen.CreateSubNoteResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSubNote indicates an expected call of CreateSubNote.
func (mr *MockNoteServerMockRecorder) CreateSubNote(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSubNote", reflect.TypeOf((*MockNoteServer)(nil).CreateSubNote), arg0, arg1)
}

// DeleteNote mocks base method.
func (m *MockNoteServer) DeleteNote(arg0 context.Context, arg1 *gen.DeleteNoteRequest) (*gen.DeleteNoteResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteNote", arg0, arg1)
	ret0, _ := ret[0].(*gen.DeleteNoteResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteNote indicates an expected call of DeleteNote.
func (mr *MockNoteServerMockRecorder) DeleteNote(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteNote", reflect.TypeOf((*MockNoteServer)(nil).DeleteNote), arg0, arg1)
}

// DeleteTag mocks base method.
func (m *MockNoteServer) DeleteTag(arg0 context.Context, arg1 *gen.TagRequest) (*gen.GetNoteResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTag", arg0, arg1)
	ret0, _ := ret[0].(*gen.GetNoteResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteTag indicates an expected call of DeleteTag.
func (mr *MockNoteServerMockRecorder) DeleteTag(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTag", reflect.TypeOf((*MockNoteServer)(nil).DeleteTag), arg0, arg1)
}

// ForgetTag mocks base method.
func (m *MockNoteServer) ForgetTag(arg0 context.Context, arg1 *gen.AllTagRequest) (*gen.EmptyResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ForgetTag", arg0, arg1)
	ret0, _ := ret[0].(*gen.EmptyResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ForgetTag indicates an expected call of ForgetTag.
func (mr *MockNoteServerMockRecorder) ForgetTag(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForgetTag", reflect.TypeOf((*MockNoteServer)(nil).ForgetTag), arg0, arg1)
}

// GetAllNotes mocks base method.
func (m *MockNoteServer) GetAllNotes(arg0 context.Context, arg1 *gen.GetAllRequest) (*gen.GetAllResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllNotes", arg0, arg1)
	ret0, _ := ret[0].(*gen.GetAllResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllNotes indicates an expected call of GetAllNotes.
func (mr *MockNoteServerMockRecorder) GetAllNotes(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllNotes", reflect.TypeOf((*MockNoteServer)(nil).GetAllNotes), arg0, arg1)
}

// GetNote mocks base method.
func (m *MockNoteServer) GetNote(arg0 context.Context, arg1 *gen.GetNoteRequest) (*gen.GetNoteResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNote", arg0, arg1)
	ret0, _ := ret[0].(*gen.GetNoteResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNote indicates an expected call of GetNote.
func (mr *MockNoteServerMockRecorder) GetNote(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNote", reflect.TypeOf((*MockNoteServer)(nil).GetNote), arg0, arg1)
}

// GetTags mocks base method.
func (m *MockNoteServer) GetTags(arg0 context.Context, arg1 *gen.GetTagsRequest) (*gen.GetTagsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTags", arg0, arg1)
	ret0, _ := ret[0].(*gen.GetTagsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTags indicates an expected call of GetTags.
func (mr *MockNoteServerMockRecorder) GetTags(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTags", reflect.TypeOf((*MockNoteServer)(nil).GetTags), arg0, arg1)
}

// RememberTag mocks base method.
func (m *MockNoteServer) RememberTag(arg0 context.Context, arg1 *gen.AllTagRequest) (*gen.EmptyResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RememberTag", arg0, arg1)
	ret0, _ := ret[0].(*gen.EmptyResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RememberTag indicates an expected call of RememberTag.
func (mr *MockNoteServerMockRecorder) RememberTag(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RememberTag", reflect.TypeOf((*MockNoteServer)(nil).RememberTag), arg0, arg1)
}

// UpdateNote mocks base method.
func (m *MockNoteServer) UpdateNote(arg0 context.Context, arg1 *gen.UpdateNoteRequest) (*gen.UpdateNoteResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateNote", arg0, arg1)
	ret0, _ := ret[0].(*gen.UpdateNoteResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateNote indicates an expected call of UpdateNote.
func (mr *MockNoteServerMockRecorder) UpdateNote(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateNote", reflect.TypeOf((*MockNoteServer)(nil).UpdateNote), arg0, arg1)
}

// mustEmbedUnimplementedNoteServer mocks base method.
func (m *MockNoteServer) mustEmbedUnimplementedNoteServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedNoteServer")
}

// mustEmbedUnimplementedNoteServer indicates an expected call of mustEmbedUnimplementedNoteServer.
func (mr *MockNoteServerMockRecorder) mustEmbedUnimplementedNoteServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedNoteServer", reflect.TypeOf((*MockNoteServer)(nil).mustEmbedUnimplementedNoteServer))
}

// MockUnsafeNoteServer is a mock of UnsafeNoteServer interface.
type MockUnsafeNoteServer struct {
	ctrl     *gomock.Controller
	recorder *MockUnsafeNoteServerMockRecorder
}

// MockUnsafeNoteServerMockRecorder is the mock recorder for MockUnsafeNoteServer.
type MockUnsafeNoteServerMockRecorder struct {
	mock *MockUnsafeNoteServer
}

// NewMockUnsafeNoteServer creates a new mock instance.
func NewMockUnsafeNoteServer(ctrl *gomock.Controller) *MockUnsafeNoteServer {
	mock := &MockUnsafeNoteServer{ctrl: ctrl}
	mock.recorder = &MockUnsafeNoteServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUnsafeNoteServer) EXPECT() *MockUnsafeNoteServerMockRecorder {
	return m.recorder
}

// mustEmbedUnimplementedNoteServer mocks base method.
func (m *MockUnsafeNoteServer) mustEmbedUnimplementedNoteServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedNoteServer")
}

// mustEmbedUnimplementedNoteServer indicates an expected call of mustEmbedUnimplementedNoteServer.
func (mr *MockUnsafeNoteServerMockRecorder) mustEmbedUnimplementedNoteServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedNoteServer", reflect.TypeOf((*MockUnsafeNoteServer)(nil).mustEmbedUnimplementedNoteServer))
}
