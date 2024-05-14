// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.12.4
// source: note.proto

package gen

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Note_GetAllNotes_FullMethodName      = "/note.Note/GetAllNotes"
	Note_GetNote_FullMethodName          = "/note.Note/GetNote"
	Note_AddNote_FullMethodName          = "/note.Note/AddNote"
	Note_UpdateNote_FullMethodName       = "/note.Note/UpdateNote"
	Note_DeleteNote_FullMethodName       = "/note.Note/DeleteNote"
	Note_CreateSubNote_FullMethodName    = "/note.Note/CreateSubNote"
	Note_AddCollaborator_FullMethodName  = "/note.Note/AddCollaborator"
	Note_AddTag_FullMethodName           = "/note.Note/AddTag"
	Note_DeleteTag_FullMethodName        = "/note.Note/DeleteTag"
	Note_GetTags_FullMethodName          = "/note.Note/GetTags"
	Note_CheckPermissions_FullMethodName = "/note.Note/CheckPermissions"
	Note_RememberTag_FullMethodName      = "/note.Note/RememberTag"
	Note_ForgetTag_FullMethodName        = "/note.Note/ForgetTag"
	Note_UpdateTag_FullMethodName        = "/note.Note/UpdateTag"
	Note_SetIcon_FullMethodName          = "/note.Note/SetIcon"
	Note_SetHeader_FullMethodName        = "/note.Note/SetHeader"
)

// NoteClient is the client API for Note service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NoteClient interface {
	GetAllNotes(ctx context.Context, in *GetAllRequest, opts ...grpc.CallOption) (*GetAllResponse, error)
	GetNote(ctx context.Context, in *GetNoteRequest, opts ...grpc.CallOption) (*GetNoteResponse, error)
	AddNote(ctx context.Context, in *AddNoteRequest, opts ...grpc.CallOption) (*AddNoteResponse, error)
	UpdateNote(ctx context.Context, in *UpdateNoteRequest, opts ...grpc.CallOption) (*UpdateNoteResponse, error)
	DeleteNote(ctx context.Context, in *DeleteNoteRequest, opts ...grpc.CallOption) (*DeleteNoteResponse, error)
	CreateSubNote(ctx context.Context, in *CreateSubNoteRequest, opts ...grpc.CallOption) (*CreateSubNoteResponse, error)
	AddCollaborator(ctx context.Context, in *AddCollaboratorRequest, opts ...grpc.CallOption) (*AddCollaboratorResponse, error)
	AddTag(ctx context.Context, in *TagRequest, opts ...grpc.CallOption) (*GetNoteResponse, error)
	DeleteTag(ctx context.Context, in *TagRequest, opts ...grpc.CallOption) (*GetNoteResponse, error)
	GetTags(ctx context.Context, in *GetTagsRequest, opts ...grpc.CallOption) (*GetTagsResponse, error)
	CheckPermissions(ctx context.Context, in *CheckPermissionsRequest, opts ...grpc.CallOption) (*CheckPermissionsResponse, error)
	RememberTag(ctx context.Context, in *AllTagRequest, opts ...grpc.CallOption) (*EmptyResponse, error)
	ForgetTag(ctx context.Context, in *AllTagRequest, opts ...grpc.CallOption) (*EmptyResponse, error)
	UpdateTag(ctx context.Context, in *UpdateTagRequest, opts ...grpc.CallOption) (*EmptyResponse, error)
	SetIcon(ctx context.Context, in *SetIconRequest, opts ...grpc.CallOption) (*GetNoteResponse, error)
	SetHeader(ctx context.Context, in *SetHeaderRequest, opts ...grpc.CallOption) (*GetNoteResponse, error)
}

type noteClient struct {
	cc grpc.ClientConnInterface
}

func NewNoteClient(cc grpc.ClientConnInterface) NoteClient {
	return &noteClient{cc}
}

func (c *noteClient) GetAllNotes(ctx context.Context, in *GetAllRequest, opts ...grpc.CallOption) (*GetAllResponse, error) {
	out := new(GetAllResponse)
	err := c.cc.Invoke(ctx, Note_GetAllNotes_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *noteClient) GetNote(ctx context.Context, in *GetNoteRequest, opts ...grpc.CallOption) (*GetNoteResponse, error) {
	out := new(GetNoteResponse)
	err := c.cc.Invoke(ctx, Note_GetNote_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *noteClient) AddNote(ctx context.Context, in *AddNoteRequest, opts ...grpc.CallOption) (*AddNoteResponse, error) {
	out := new(AddNoteResponse)
	err := c.cc.Invoke(ctx, Note_AddNote_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *noteClient) UpdateNote(ctx context.Context, in *UpdateNoteRequest, opts ...grpc.CallOption) (*UpdateNoteResponse, error) {
	out := new(UpdateNoteResponse)
	err := c.cc.Invoke(ctx, Note_UpdateNote_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *noteClient) DeleteNote(ctx context.Context, in *DeleteNoteRequest, opts ...grpc.CallOption) (*DeleteNoteResponse, error) {
	out := new(DeleteNoteResponse)
	err := c.cc.Invoke(ctx, Note_DeleteNote_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *noteClient) CreateSubNote(ctx context.Context, in *CreateSubNoteRequest, opts ...grpc.CallOption) (*CreateSubNoteResponse, error) {
	out := new(CreateSubNoteResponse)
	err := c.cc.Invoke(ctx, Note_CreateSubNote_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *noteClient) AddCollaborator(ctx context.Context, in *AddCollaboratorRequest, opts ...grpc.CallOption) (*AddCollaboratorResponse, error) {
	out := new(AddCollaboratorResponse)
	err := c.cc.Invoke(ctx, Note_AddCollaborator_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *noteClient) AddTag(ctx context.Context, in *TagRequest, opts ...grpc.CallOption) (*GetNoteResponse, error) {
	out := new(GetNoteResponse)
	err := c.cc.Invoke(ctx, Note_AddTag_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *noteClient) DeleteTag(ctx context.Context, in *TagRequest, opts ...grpc.CallOption) (*GetNoteResponse, error) {
	out := new(GetNoteResponse)
	err := c.cc.Invoke(ctx, Note_DeleteTag_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *noteClient) GetTags(ctx context.Context, in *GetTagsRequest, opts ...grpc.CallOption) (*GetTagsResponse, error) {
	out := new(GetTagsResponse)
	err := c.cc.Invoke(ctx, Note_GetTags_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *noteClient) CheckPermissions(ctx context.Context, in *CheckPermissionsRequest, opts ...grpc.CallOption) (*CheckPermissionsResponse, error) {
	out := new(CheckPermissionsResponse)
	err := c.cc.Invoke(ctx, Note_CheckPermissions_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *noteClient) RememberTag(ctx context.Context, in *AllTagRequest, opts ...grpc.CallOption) (*EmptyResponse, error) {
	out := new(EmptyResponse)
	err := c.cc.Invoke(ctx, Note_RememberTag_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *noteClient) ForgetTag(ctx context.Context, in *AllTagRequest, opts ...grpc.CallOption) (*EmptyResponse, error) {
	out := new(EmptyResponse)
	err := c.cc.Invoke(ctx, Note_ForgetTag_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *noteClient) UpdateTag(ctx context.Context, in *UpdateTagRequest, opts ...grpc.CallOption) (*EmptyResponse, error) {
	out := new(EmptyResponse)
	err := c.cc.Invoke(ctx, Note_UpdateTag_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *noteClient) SetIcon(ctx context.Context, in *SetIconRequest, opts ...grpc.CallOption) (*GetNoteResponse, error) {
	out := new(GetNoteResponse)
	err := c.cc.Invoke(ctx, Note_SetIcon_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *noteClient) SetHeader(ctx context.Context, in *SetHeaderRequest, opts ...grpc.CallOption) (*GetNoteResponse, error) {
	out := new(GetNoteResponse)
	err := c.cc.Invoke(ctx, Note_SetHeader_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NoteServer is the server API for Note service.
// All implementations must embed UnimplementedNoteServer
// for forward compatibility
type NoteServer interface {
	GetAllNotes(context.Context, *GetAllRequest) (*GetAllResponse, error)
	GetNote(context.Context, *GetNoteRequest) (*GetNoteResponse, error)
	AddNote(context.Context, *AddNoteRequest) (*AddNoteResponse, error)
	UpdateNote(context.Context, *UpdateNoteRequest) (*UpdateNoteResponse, error)
	DeleteNote(context.Context, *DeleteNoteRequest) (*DeleteNoteResponse, error)
	CreateSubNote(context.Context, *CreateSubNoteRequest) (*CreateSubNoteResponse, error)
	AddCollaborator(context.Context, *AddCollaboratorRequest) (*AddCollaboratorResponse, error)
	AddTag(context.Context, *TagRequest) (*GetNoteResponse, error)
	DeleteTag(context.Context, *TagRequest) (*GetNoteResponse, error)
	GetTags(context.Context, *GetTagsRequest) (*GetTagsResponse, error)
	CheckPermissions(context.Context, *CheckPermissionsRequest) (*CheckPermissionsResponse, error)
	RememberTag(context.Context, *AllTagRequest) (*EmptyResponse, error)
	ForgetTag(context.Context, *AllTagRequest) (*EmptyResponse, error)
	UpdateTag(context.Context, *UpdateTagRequest) (*EmptyResponse, error)
	SetIcon(context.Context, *SetIconRequest) (*GetNoteResponse, error)
	SetHeader(context.Context, *SetHeaderRequest) (*GetNoteResponse, error)
	mustEmbedUnimplementedNoteServer()
}

// UnimplementedNoteServer must be embedded to have forward compatible implementations.
type UnimplementedNoteServer struct {
}

func (UnimplementedNoteServer) GetAllNotes(context.Context, *GetAllRequest) (*GetAllResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllNotes not implemented")
}
func (UnimplementedNoteServer) GetNote(context.Context, *GetNoteRequest) (*GetNoteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNote not implemented")
}
func (UnimplementedNoteServer) AddNote(context.Context, *AddNoteRequest) (*AddNoteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddNote not implemented")
}
func (UnimplementedNoteServer) UpdateNote(context.Context, *UpdateNoteRequest) (*UpdateNoteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateNote not implemented")
}
func (UnimplementedNoteServer) DeleteNote(context.Context, *DeleteNoteRequest) (*DeleteNoteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteNote not implemented")
}
func (UnimplementedNoteServer) CreateSubNote(context.Context, *CreateSubNoteRequest) (*CreateSubNoteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateSubNote not implemented")
}
func (UnimplementedNoteServer) AddCollaborator(context.Context, *AddCollaboratorRequest) (*AddCollaboratorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddCollaborator not implemented")
}
func (UnimplementedNoteServer) AddTag(context.Context, *TagRequest) (*GetNoteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddTag not implemented")
}
func (UnimplementedNoteServer) DeleteTag(context.Context, *TagRequest) (*GetNoteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteTag not implemented")
}
func (UnimplementedNoteServer) GetTags(context.Context, *GetTagsRequest) (*GetTagsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTags not implemented")
}
func (UnimplementedNoteServer) CheckPermissions(context.Context, *CheckPermissionsRequest) (*CheckPermissionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckPermissions not implemented")
}
func (UnimplementedNoteServer) RememberTag(context.Context, *AllTagRequest) (*EmptyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RememberTag not implemented")
}
func (UnimplementedNoteServer) ForgetTag(context.Context, *AllTagRequest) (*EmptyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ForgetTag not implemented")
}
func (UnimplementedNoteServer) UpdateTag(context.Context, *UpdateTagRequest) (*EmptyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateTag not implemented")
}
func (UnimplementedNoteServer) SetIcon(context.Context, *SetIconRequest) (*GetNoteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetIcon not implemented")
}
func (UnimplementedNoteServer) SetHeader(context.Context, *SetHeaderRequest) (*GetNoteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetHeader not implemented")
}
func (UnimplementedNoteServer) mustEmbedUnimplementedNoteServer() {}

// UnsafeNoteServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NoteServer will
// result in compilation errors.
type UnsafeNoteServer interface {
	mustEmbedUnimplementedNoteServer()
}

func RegisterNoteServer(s grpc.ServiceRegistrar, srv NoteServer) {
	s.RegisterService(&Note_ServiceDesc, srv)
}

func _Note_GetAllNotes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NoteServer).GetAllNotes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Note_GetAllNotes_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NoteServer).GetAllNotes(ctx, req.(*GetAllRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Note_GetNote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetNoteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NoteServer).GetNote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Note_GetNote_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NoteServer).GetNote(ctx, req.(*GetNoteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Note_AddNote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddNoteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NoteServer).AddNote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Note_AddNote_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NoteServer).AddNote(ctx, req.(*AddNoteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Note_UpdateNote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateNoteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NoteServer).UpdateNote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Note_UpdateNote_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NoteServer).UpdateNote(ctx, req.(*UpdateNoteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Note_DeleteNote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteNoteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NoteServer).DeleteNote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Note_DeleteNote_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NoteServer).DeleteNote(ctx, req.(*DeleteNoteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Note_CreateSubNote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateSubNoteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NoteServer).CreateSubNote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Note_CreateSubNote_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NoteServer).CreateSubNote(ctx, req.(*CreateSubNoteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Note_AddCollaborator_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddCollaboratorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NoteServer).AddCollaborator(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Note_AddCollaborator_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NoteServer).AddCollaborator(ctx, req.(*AddCollaboratorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Note_AddTag_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TagRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NoteServer).AddTag(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Note_AddTag_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NoteServer).AddTag(ctx, req.(*TagRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Note_DeleteTag_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TagRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NoteServer).DeleteTag(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Note_DeleteTag_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NoteServer).DeleteTag(ctx, req.(*TagRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Note_GetTags_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTagsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NoteServer).GetTags(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Note_GetTags_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NoteServer).GetTags(ctx, req.(*GetTagsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Note_CheckPermissions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckPermissionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NoteServer).CheckPermissions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Note_CheckPermissions_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NoteServer).CheckPermissions(ctx, req.(*CheckPermissionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Note_RememberTag_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AllTagRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NoteServer).RememberTag(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Note_RememberTag_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NoteServer).RememberTag(ctx, req.(*AllTagRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Note_ForgetTag_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AllTagRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NoteServer).ForgetTag(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Note_ForgetTag_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NoteServer).ForgetTag(ctx, req.(*AllTagRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Note_UpdateTag_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateTagRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NoteServer).UpdateTag(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Note_UpdateTag_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NoteServer).UpdateTag(ctx, req.(*UpdateTagRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Note_SetIcon_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetIconRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NoteServer).SetIcon(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Note_SetIcon_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NoteServer).SetIcon(ctx, req.(*SetIconRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Note_SetHeader_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetHeaderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NoteServer).SetHeader(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Note_SetHeader_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NoteServer).SetHeader(ctx, req.(*SetHeaderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Note_ServiceDesc is the grpc.ServiceDesc for Note service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Note_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "note.Note",
	HandlerType: (*NoteServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetAllNotes",
			Handler:    _Note_GetAllNotes_Handler,
		},
		{
			MethodName: "GetNote",
			Handler:    _Note_GetNote_Handler,
		},
		{
			MethodName: "AddNote",
			Handler:    _Note_AddNote_Handler,
		},
		{
			MethodName: "UpdateNote",
			Handler:    _Note_UpdateNote_Handler,
		},
		{
			MethodName: "DeleteNote",
			Handler:    _Note_DeleteNote_Handler,
		},
		{
			MethodName: "CreateSubNote",
			Handler:    _Note_CreateSubNote_Handler,
		},
		{
			MethodName: "AddCollaborator",
			Handler:    _Note_AddCollaborator_Handler,
		},
		{
			MethodName: "AddTag",
			Handler:    _Note_AddTag_Handler,
		},
		{
			MethodName: "DeleteTag",
			Handler:    _Note_DeleteTag_Handler,
		},
		{
			MethodName: "GetTags",
			Handler:    _Note_GetTags_Handler,
		},
		{
			MethodName: "CheckPermissions",
			Handler:    _Note_CheckPermissions_Handler,
		},
		{
			MethodName: "RememberTag",
			Handler:    _Note_RememberTag_Handler,
		},
		{
			MethodName: "ForgetTag",
			Handler:    _Note_ForgetTag_Handler,
		},
		{
			MethodName: "UpdateTag",
			Handler:    _Note_UpdateTag_Handler,
		},
		{
			MethodName: "SetIcon",
			Handler:    _Note_SetIcon_Handler,
		},
		{
			MethodName: "SetHeader",
			Handler:    _Note_SetHeader_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "note.proto",
}
