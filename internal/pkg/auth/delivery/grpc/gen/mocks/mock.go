// Code generated by MockGen. DO NOT EDIT.
// Source: auth_grpc.pb.go

// Package mock_gen is a generated GoMock package.
package mock_gen

import (
	context "context"
	reflect "reflect"

	gen "github.com/go-park-mail-ru/2024_1_scratch_senior_devs/internal/pkg/auth/delivery/grpc/gen"
	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockAuthClient is a mock of AuthClient interface.
type MockAuthClient struct {
	ctrl     *gomock.Controller
	recorder *MockAuthClientMockRecorder
}

// MockAuthClientMockRecorder is the mock recorder for MockAuthClient.
type MockAuthClientMockRecorder struct {
	mock *MockAuthClient
}

// NewMockAuthClient creates a new mock instance.
func NewMockAuthClient(ctrl *gomock.Controller) *MockAuthClient {
	mock := &MockAuthClient{ctrl: ctrl}
	mock.recorder = &MockAuthClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthClient) EXPECT() *MockAuthClientMockRecorder {
	return m.recorder
}

// CheckUser mocks base method.
func (m *MockAuthClient) CheckUser(ctx context.Context, in *gen.CheckUserRequest, opts ...grpc.CallOption) (*gen.User, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CheckUser", varargs...)
	ret0, _ := ret[0].(*gen.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckUser indicates an expected call of CheckUser.
func (mr *MockAuthClientMockRecorder) CheckUser(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckUser", reflect.TypeOf((*MockAuthClient)(nil).CheckUser), varargs...)
}

// DeleteSecret mocks base method.
func (m *MockAuthClient) DeleteSecret(ctx context.Context, in *gen.SecretRequest, opts ...grpc.CallOption) (*gen.EmptyMessage, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteSecret", varargs...)
	ret0, _ := ret[0].(*gen.EmptyMessage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteSecret indicates an expected call of DeleteSecret.
func (mr *MockAuthClientMockRecorder) DeleteSecret(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSecret", reflect.TypeOf((*MockAuthClient)(nil).DeleteSecret), varargs...)
}

// GenerateAndUpdateSecret mocks base method.
func (m *MockAuthClient) GenerateAndUpdateSecret(ctx context.Context, in *gen.SecretRequest, opts ...grpc.CallOption) (*gen.GenerateAndUpdateSecretResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GenerateAndUpdateSecret", varargs...)
	ret0, _ := ret[0].(*gen.GenerateAndUpdateSecretResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateAndUpdateSecret indicates an expected call of GenerateAndUpdateSecret.
func (mr *MockAuthClientMockRecorder) GenerateAndUpdateSecret(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateAndUpdateSecret", reflect.TypeOf((*MockAuthClient)(nil).GenerateAndUpdateSecret), varargs...)
}

// GetUserByUsername mocks base method.
func (m *MockAuthClient) GetUserByUsername(ctx context.Context, in *gen.GetUserByUsernameRequest, opts ...grpc.CallOption) (*gen.User, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetUserByUsername", varargs...)
	ret0, _ := ret[0].(*gen.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByUsername indicates an expected call of GetUserByUsername.
func (mr *MockAuthClientMockRecorder) GetUserByUsername(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByUsername", reflect.TypeOf((*MockAuthClient)(nil).GetUserByUsername), varargs...)
}

// SignIn mocks base method.
func (m *MockAuthClient) SignIn(ctx context.Context, in *gen.UserFormData, opts ...grpc.CallOption) (*gen.SignInResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SignIn", varargs...)
	ret0, _ := ret[0].(*gen.SignInResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignIn indicates an expected call of SignIn.
func (mr *MockAuthClientMockRecorder) SignIn(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignIn", reflect.TypeOf((*MockAuthClient)(nil).SignIn), varargs...)
}

// SignUp mocks base method.
func (m *MockAuthClient) SignUp(ctx context.Context, in *gen.UserFormData, opts ...grpc.CallOption) (*gen.SignUpResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SignUp", varargs...)
	ret0, _ := ret[0].(*gen.SignUpResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignUp indicates an expected call of SignUp.
func (mr *MockAuthClientMockRecorder) SignUp(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignUp", reflect.TypeOf((*MockAuthClient)(nil).SignUp), varargs...)
}

// UpdateProfile mocks base method.
func (m *MockAuthClient) UpdateProfile(ctx context.Context, in *gen.UpdateProfileRequest, opts ...grpc.CallOption) (*gen.User, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateProfile", varargs...)
	ret0, _ := ret[0].(*gen.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateProfile indicates an expected call of UpdateProfile.
func (mr *MockAuthClientMockRecorder) UpdateProfile(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProfile", reflect.TypeOf((*MockAuthClient)(nil).UpdateProfile), varargs...)
}

// UpdateProfileAvatar mocks base method.
func (m *MockAuthClient) UpdateProfileAvatar(ctx context.Context, in *gen.UpdateProfileAvatarRequest, opts ...grpc.CallOption) (*gen.User, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateProfileAvatar", varargs...)
	ret0, _ := ret[0].(*gen.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateProfileAvatar indicates an expected call of UpdateProfileAvatar.
func (mr *MockAuthClientMockRecorder) UpdateProfileAvatar(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProfileAvatar", reflect.TypeOf((*MockAuthClient)(nil).UpdateProfileAvatar), varargs...)
}

// MockAuthServer is a mock of AuthServer interface.
type MockAuthServer struct {
	ctrl     *gomock.Controller
	recorder *MockAuthServerMockRecorder
}

// MockAuthServerMockRecorder is the mock recorder for MockAuthServer.
type MockAuthServerMockRecorder struct {
	mock *MockAuthServer
}

// NewMockAuthServer creates a new mock instance.
func NewMockAuthServer(ctrl *gomock.Controller) *MockAuthServer {
	mock := &MockAuthServer{ctrl: ctrl}
	mock.recorder = &MockAuthServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthServer) EXPECT() *MockAuthServerMockRecorder {
	return m.recorder
}

// CheckUser mocks base method.
func (m *MockAuthServer) CheckUser(arg0 context.Context, arg1 *gen.CheckUserRequest) (*gen.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckUser", arg0, arg1)
	ret0, _ := ret[0].(*gen.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckUser indicates an expected call of CheckUser.
func (mr *MockAuthServerMockRecorder) CheckUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckUser", reflect.TypeOf((*MockAuthServer)(nil).CheckUser), arg0, arg1)
}

// DeleteSecret mocks base method.
func (m *MockAuthServer) DeleteSecret(arg0 context.Context, arg1 *gen.SecretRequest) (*gen.EmptyMessage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSecret", arg0, arg1)
	ret0, _ := ret[0].(*gen.EmptyMessage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteSecret indicates an expected call of DeleteSecret.
func (mr *MockAuthServerMockRecorder) DeleteSecret(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSecret", reflect.TypeOf((*MockAuthServer)(nil).DeleteSecret), arg0, arg1)
}

// GenerateAndUpdateSecret mocks base method.
func (m *MockAuthServer) GenerateAndUpdateSecret(arg0 context.Context, arg1 *gen.SecretRequest) (*gen.GenerateAndUpdateSecretResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateAndUpdateSecret", arg0, arg1)
	ret0, _ := ret[0].(*gen.GenerateAndUpdateSecretResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateAndUpdateSecret indicates an expected call of GenerateAndUpdateSecret.
func (mr *MockAuthServerMockRecorder) GenerateAndUpdateSecret(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateAndUpdateSecret", reflect.TypeOf((*MockAuthServer)(nil).GenerateAndUpdateSecret), arg0, arg1)
}

// GetUserByUsername mocks base method.
func (m *MockAuthServer) GetUserByUsername(arg0 context.Context, arg1 *gen.GetUserByUsernameRequest) (*gen.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByUsername", arg0, arg1)
	ret0, _ := ret[0].(*gen.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByUsername indicates an expected call of GetUserByUsername.
func (mr *MockAuthServerMockRecorder) GetUserByUsername(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByUsername", reflect.TypeOf((*MockAuthServer)(nil).GetUserByUsername), arg0, arg1)
}

// SignIn mocks base method.
func (m *MockAuthServer) SignIn(arg0 context.Context, arg1 *gen.UserFormData) (*gen.SignInResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignIn", arg0, arg1)
	ret0, _ := ret[0].(*gen.SignInResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignIn indicates an expected call of SignIn.
func (mr *MockAuthServerMockRecorder) SignIn(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignIn", reflect.TypeOf((*MockAuthServer)(nil).SignIn), arg0, arg1)
}

// SignUp mocks base method.
func (m *MockAuthServer) SignUp(arg0 context.Context, arg1 *gen.UserFormData) (*gen.SignUpResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignUp", arg0, arg1)
	ret0, _ := ret[0].(*gen.SignUpResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignUp indicates an expected call of SignUp.
func (mr *MockAuthServerMockRecorder) SignUp(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignUp", reflect.TypeOf((*MockAuthServer)(nil).SignUp), arg0, arg1)
}

// UpdateProfile mocks base method.
func (m *MockAuthServer) UpdateProfile(arg0 context.Context, arg1 *gen.UpdateProfileRequest) (*gen.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateProfile", arg0, arg1)
	ret0, _ := ret[0].(*gen.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateProfile indicates an expected call of UpdateProfile.
func (mr *MockAuthServerMockRecorder) UpdateProfile(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProfile", reflect.TypeOf((*MockAuthServer)(nil).UpdateProfile), arg0, arg1)
}

// UpdateProfileAvatar mocks base method.
func (m *MockAuthServer) UpdateProfileAvatar(arg0 context.Context, arg1 *gen.UpdateProfileAvatarRequest) (*gen.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateProfileAvatar", arg0, arg1)
	ret0, _ := ret[0].(*gen.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateProfileAvatar indicates an expected call of UpdateProfileAvatar.
func (mr *MockAuthServerMockRecorder) UpdateProfileAvatar(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProfileAvatar", reflect.TypeOf((*MockAuthServer)(nil).UpdateProfileAvatar), arg0, arg1)
}

// mustEmbedUnimplementedAuthServer mocks base method.
func (m *MockAuthServer) mustEmbedUnimplementedAuthServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedAuthServer")
}

// mustEmbedUnimplementedAuthServer indicates an expected call of mustEmbedUnimplementedAuthServer.
func (mr *MockAuthServerMockRecorder) mustEmbedUnimplementedAuthServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedAuthServer", reflect.TypeOf((*MockAuthServer)(nil).mustEmbedUnimplementedAuthServer))
}

// MockUnsafeAuthServer is a mock of UnsafeAuthServer interface.
type MockUnsafeAuthServer struct {
	ctrl     *gomock.Controller
	recorder *MockUnsafeAuthServerMockRecorder
}

// MockUnsafeAuthServerMockRecorder is the mock recorder for MockUnsafeAuthServer.
type MockUnsafeAuthServerMockRecorder struct {
	mock *MockUnsafeAuthServer
}

// NewMockUnsafeAuthServer creates a new mock instance.
func NewMockUnsafeAuthServer(ctrl *gomock.Controller) *MockUnsafeAuthServer {
	mock := &MockUnsafeAuthServer{ctrl: ctrl}
	mock.recorder = &MockUnsafeAuthServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUnsafeAuthServer) EXPECT() *MockUnsafeAuthServerMockRecorder {
	return m.recorder
}

// mustEmbedUnimplementedAuthServer mocks base method.
func (m *MockUnsafeAuthServer) mustEmbedUnimplementedAuthServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedAuthServer")
}

// mustEmbedUnimplementedAuthServer indicates an expected call of mustEmbedUnimplementedAuthServer.
func (mr *MockUnsafeAuthServerMockRecorder) mustEmbedUnimplementedAuthServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedAuthServer", reflect.TypeOf((*MockUnsafeAuthServer)(nil).mustEmbedUnimplementedAuthServer))
}
