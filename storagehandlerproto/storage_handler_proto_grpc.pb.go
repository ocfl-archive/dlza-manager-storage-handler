// protoc --go_out=. --go-grpc_out=. proto/copy.proto

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.24.4
// source: storage_handler_proto.proto

package storagehandlerproto

import (
	context "context"
	dlzamanagerproto "github.com/ocfl-archive/dlza-manager/dlzamanagerproto"
	proto "go.ub.unibas.ch/cloud/genericproto/v2/pkg/generic/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	ClerkStorageHandlerService_CreateStoragePartition_FullMethodName = "/storagehandlerproto.ClerkStorageHandlerService/CreateStoragePartition"
)

// ClerkStorageHandlerServiceClient is the client API for ClerkStorageHandlerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ClerkStorageHandlerServiceClient interface {
	CreateStoragePartition(ctx context.Context, in *dlzamanagerproto.StoragePartition, opts ...grpc.CallOption) (*dlzamanagerproto.Status, error)
}

type clerkStorageHandlerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewClerkStorageHandlerServiceClient(cc grpc.ClientConnInterface) ClerkStorageHandlerServiceClient {
	return &clerkStorageHandlerServiceClient{cc}
}

func (c *clerkStorageHandlerServiceClient) CreateStoragePartition(ctx context.Context, in *dlzamanagerproto.StoragePartition, opts ...grpc.CallOption) (*dlzamanagerproto.Status, error) {
	out := new(dlzamanagerproto.Status)
	err := c.cc.Invoke(ctx, ClerkStorageHandlerService_CreateStoragePartition_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ClerkStorageHandlerServiceServer is the server API for ClerkStorageHandlerService service.
// All implementations must embed UnimplementedClerkStorageHandlerServiceServer
// for forward compatibility
type ClerkStorageHandlerServiceServer interface {
	CreateStoragePartition(context.Context, *dlzamanagerproto.StoragePartition) (*dlzamanagerproto.Status, error)
	mustEmbedUnimplementedClerkStorageHandlerServiceServer()
}

// UnimplementedClerkStorageHandlerServiceServer must be embedded to have forward compatible implementations.
type UnimplementedClerkStorageHandlerServiceServer struct {
}

func (UnimplementedClerkStorageHandlerServiceServer) CreateStoragePartition(context.Context, *dlzamanagerproto.StoragePartition) (*dlzamanagerproto.Status, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateStoragePartition not implemented")
}
func (UnimplementedClerkStorageHandlerServiceServer) mustEmbedUnimplementedClerkStorageHandlerServiceServer() {
}

// UnsafeClerkStorageHandlerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ClerkStorageHandlerServiceServer will
// result in compilation errors.
type UnsafeClerkStorageHandlerServiceServer interface {
	mustEmbedUnimplementedClerkStorageHandlerServiceServer()
}

func RegisterClerkStorageHandlerServiceServer(s grpc.ServiceRegistrar, srv ClerkStorageHandlerServiceServer) {
	s.RegisterService(&ClerkStorageHandlerService_ServiceDesc, srv)
}

func _ClerkStorageHandlerService_CreateStoragePartition_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(dlzamanagerproto.StoragePartition)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClerkStorageHandlerServiceServer).CreateStoragePartition(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ClerkStorageHandlerService_CreateStoragePartition_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClerkStorageHandlerServiceServer).CreateStoragePartition(ctx, req.(*dlzamanagerproto.StoragePartition))
	}
	return interceptor(ctx, in, info, handler)
}

// ClerkStorageHandlerService_ServiceDesc is the grpc.ServiceDesc for ClerkStorageHandlerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ClerkStorageHandlerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "storagehandlerproto.ClerkStorageHandlerService",
	HandlerType: (*ClerkStorageHandlerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateStoragePartition",
			Handler:    _ClerkStorageHandlerService_CreateStoragePartition_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "storage_handler_proto.proto",
}

const (
	DispatcherStorageHandlerService_ConnectionCheck_FullMethodName = "/storagehandlerproto.DispatcherStorageHandlerService/ConnectionCheck"
	DispatcherStorageHandlerService_CopyArchiveTo_FullMethodName   = "/storagehandlerproto.DispatcherStorageHandlerService/CopyArchiveTo"
)

// DispatcherStorageHandlerServiceClient is the client API for DispatcherStorageHandlerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DispatcherStorageHandlerServiceClient interface {
	ConnectionCheck(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*dlzamanagerproto.Id, error)
	CopyArchiveTo(ctx context.Context, in *dlzamanagerproto.CopyFromTo, opts ...grpc.CallOption) (*dlzamanagerproto.NoParam, error)
}

type dispatcherStorageHandlerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDispatcherStorageHandlerServiceClient(cc grpc.ClientConnInterface) DispatcherStorageHandlerServiceClient {
	return &dispatcherStorageHandlerServiceClient{cc}
}

func (c *dispatcherStorageHandlerServiceClient) ConnectionCheck(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*dlzamanagerproto.Id, error) {
	out := new(dlzamanagerproto.Id)
	err := c.cc.Invoke(ctx, DispatcherStorageHandlerService_ConnectionCheck_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dispatcherStorageHandlerServiceClient) CopyArchiveTo(ctx context.Context, in *dlzamanagerproto.CopyFromTo, opts ...grpc.CallOption) (*dlzamanagerproto.NoParam, error) {
	out := new(dlzamanagerproto.NoParam)
	err := c.cc.Invoke(ctx, DispatcherStorageHandlerService_CopyArchiveTo_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DispatcherStorageHandlerServiceServer is the server API for DispatcherStorageHandlerService service.
// All implementations must embed UnimplementedDispatcherStorageHandlerServiceServer
// for forward compatibility
type DispatcherStorageHandlerServiceServer interface {
	ConnectionCheck(context.Context, *emptypb.Empty) (*dlzamanagerproto.Id, error)
	CopyArchiveTo(context.Context, *dlzamanagerproto.CopyFromTo) (*dlzamanagerproto.NoParam, error)
	mustEmbedUnimplementedDispatcherStorageHandlerServiceServer()
}

// UnimplementedDispatcherStorageHandlerServiceServer must be embedded to have forward compatible implementations.
type UnimplementedDispatcherStorageHandlerServiceServer struct {
}

func (UnimplementedDispatcherStorageHandlerServiceServer) ConnectionCheck(context.Context, *emptypb.Empty) (*dlzamanagerproto.Id, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConnectionCheck not implemented")
}
func (UnimplementedDispatcherStorageHandlerServiceServer) CopyArchiveTo(context.Context, *dlzamanagerproto.CopyFromTo) (*dlzamanagerproto.NoParam, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CopyArchiveTo not implemented")
}
func (UnimplementedDispatcherStorageHandlerServiceServer) mustEmbedUnimplementedDispatcherStorageHandlerServiceServer() {
}

// UnsafeDispatcherStorageHandlerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DispatcherStorageHandlerServiceServer will
// result in compilation errors.
type UnsafeDispatcherStorageHandlerServiceServer interface {
	mustEmbedUnimplementedDispatcherStorageHandlerServiceServer()
}

func RegisterDispatcherStorageHandlerServiceServer(s grpc.ServiceRegistrar, srv DispatcherStorageHandlerServiceServer) {
	s.RegisterService(&DispatcherStorageHandlerService_ServiceDesc, srv)
}

func _DispatcherStorageHandlerService_ConnectionCheck_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DispatcherStorageHandlerServiceServer).ConnectionCheck(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DispatcherStorageHandlerService_ConnectionCheck_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DispatcherStorageHandlerServiceServer).ConnectionCheck(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _DispatcherStorageHandlerService_CopyArchiveTo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(dlzamanagerproto.CopyFromTo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DispatcherStorageHandlerServiceServer).CopyArchiveTo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DispatcherStorageHandlerService_CopyArchiveTo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DispatcherStorageHandlerServiceServer).CopyArchiveTo(ctx, req.(*dlzamanagerproto.CopyFromTo))
	}
	return interceptor(ctx, in, info, handler)
}

// DispatcherStorageHandlerService_ServiceDesc is the grpc.ServiceDesc for DispatcherStorageHandlerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DispatcherStorageHandlerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "storagehandlerproto.DispatcherStorageHandlerService",
	HandlerType: (*DispatcherStorageHandlerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ConnectionCheck",
			Handler:    _DispatcherStorageHandlerService_ConnectionCheck_Handler,
		},
		{
			MethodName: "CopyArchiveTo",
			Handler:    _DispatcherStorageHandlerService_CopyArchiveTo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "storage_handler_proto.proto",
}

const (
	CheckerStorageHandlerService_Ping_FullMethodName                      = "/storagehandlerproto.CheckerStorageHandlerService/Ping"
	CheckerStorageHandlerService_GetObjectInstanceChecksum_FullMethodName = "/storagehandlerproto.CheckerStorageHandlerService/GetObjectInstanceChecksum"
)

// CheckerStorageHandlerServiceClient is the client API for CheckerStorageHandlerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CheckerStorageHandlerServiceClient interface {
	Ping(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*proto.DefaultResponse, error)
	GetObjectInstanceChecksum(ctx context.Context, in *dlzamanagerproto.ObjectInstance, opts ...grpc.CallOption) (*dlzamanagerproto.Id, error)
}

type checkerStorageHandlerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCheckerStorageHandlerServiceClient(cc grpc.ClientConnInterface) CheckerStorageHandlerServiceClient {
	return &checkerStorageHandlerServiceClient{cc}
}

func (c *checkerStorageHandlerServiceClient) Ping(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*proto.DefaultResponse, error) {
	out := new(proto.DefaultResponse)
	err := c.cc.Invoke(ctx, CheckerStorageHandlerService_Ping_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *checkerStorageHandlerServiceClient) GetObjectInstanceChecksum(ctx context.Context, in *dlzamanagerproto.ObjectInstance, opts ...grpc.CallOption) (*dlzamanagerproto.Id, error) {
	out := new(dlzamanagerproto.Id)
	err := c.cc.Invoke(ctx, CheckerStorageHandlerService_GetObjectInstanceChecksum_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CheckerStorageHandlerServiceServer is the server API for CheckerStorageHandlerService service.
// All implementations must embed UnimplementedCheckerStorageHandlerServiceServer
// for forward compatibility
type CheckerStorageHandlerServiceServer interface {
	Ping(context.Context, *emptypb.Empty) (*proto.DefaultResponse, error)
	GetObjectInstanceChecksum(context.Context, *dlzamanagerproto.ObjectInstance) (*dlzamanagerproto.Id, error)
	mustEmbedUnimplementedCheckerStorageHandlerServiceServer()
}

// UnimplementedCheckerStorageHandlerServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCheckerStorageHandlerServiceServer struct {
}

func (UnimplementedCheckerStorageHandlerServiceServer) Ping(context.Context, *emptypb.Empty) (*proto.DefaultResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedCheckerStorageHandlerServiceServer) GetObjectInstanceChecksum(context.Context, *dlzamanagerproto.ObjectInstance) (*dlzamanagerproto.Id, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetObjectInstanceChecksum not implemented")
}
func (UnimplementedCheckerStorageHandlerServiceServer) mustEmbedUnimplementedCheckerStorageHandlerServiceServer() {
}

// UnsafeCheckerStorageHandlerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CheckerStorageHandlerServiceServer will
// result in compilation errors.
type UnsafeCheckerStorageHandlerServiceServer interface {
	mustEmbedUnimplementedCheckerStorageHandlerServiceServer()
}

func RegisterCheckerStorageHandlerServiceServer(s grpc.ServiceRegistrar, srv CheckerStorageHandlerServiceServer) {
	s.RegisterService(&CheckerStorageHandlerService_ServiceDesc, srv)
}

func _CheckerStorageHandlerService_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CheckerStorageHandlerServiceServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CheckerStorageHandlerService_Ping_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CheckerStorageHandlerServiceServer).Ping(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _CheckerStorageHandlerService_GetObjectInstanceChecksum_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(dlzamanagerproto.ObjectInstance)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CheckerStorageHandlerServiceServer).GetObjectInstanceChecksum(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CheckerStorageHandlerService_GetObjectInstanceChecksum_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CheckerStorageHandlerServiceServer).GetObjectInstanceChecksum(ctx, req.(*dlzamanagerproto.ObjectInstance))
	}
	return interceptor(ctx, in, info, handler)
}

// CheckerStorageHandlerService_ServiceDesc is the grpc.ServiceDesc for CheckerStorageHandlerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CheckerStorageHandlerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "storagehandlerproto.CheckerStorageHandlerService",
	HandlerType: (*CheckerStorageHandlerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _CheckerStorageHandlerService_Ping_Handler,
		},
		{
			MethodName: "GetObjectInstanceChecksum",
			Handler:    _CheckerStorageHandlerService_GetObjectInstanceChecksum_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "storage_handler_proto.proto",
}
