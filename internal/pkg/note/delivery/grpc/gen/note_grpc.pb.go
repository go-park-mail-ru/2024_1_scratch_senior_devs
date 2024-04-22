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
	Note_GetAllNotes_FullMethodName = "/note.Note/GetAllNotes"
	Note_GetNote_FullMethodName     = "/note.Note/GetNote"
	Note_AddNote_FullMethodName     = "/note.Note/AddNote"
	Note_UpdateNote_FullMethodName  = "/note.Note/UpdateNote"
	Note_DeleteNote_FullMethodName  = "/note.Note/DeleteNote"
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

// NoteServer is the server API for Note service.
// All implementations must embed UnimplementedNoteServer
// for forward compatibility
type NoteServer interface {
	GetAllNotes(context.Context, *GetAllRequest) (*GetAllResponse, error)
	GetNote(context.Context, *GetNoteRequest) (*GetNoteResponse, error)
	AddNote(context.Context, *AddNoteRequest) (*AddNoteResponse, error)
	UpdateNote(context.Context, *UpdateNoteRequest) (*UpdateNoteResponse, error)
	DeleteNote(context.Context, *DeleteNoteRequest) (*DeleteNoteResponse, error)
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
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "note.proto",
}