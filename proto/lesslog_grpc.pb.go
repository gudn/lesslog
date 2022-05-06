// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package proto

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

// LesslogClient is the client API for Lesslog service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LesslogClient interface {
	Fetch(ctx context.Context, in *FetchRequest, opts ...grpc.CallOption) (*FetchResponse, error)
	Watch(ctx context.Context, in *FetchRequest, opts ...grpc.CallOption) (Lesslog_WatchClient, error)
	Push(ctx context.Context, in *PushRequest, opts ...grpc.CallOption) (*PushResponse, error)
	Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*PushResponse, error)
}

type lesslogClient struct {
	cc grpc.ClientConnInterface
}

func NewLesslogClient(cc grpc.ClientConnInterface) LesslogClient {
	return &lesslogClient{cc}
}

func (c *lesslogClient) Fetch(ctx context.Context, in *FetchRequest, opts ...grpc.CallOption) (*FetchResponse, error) {
	out := new(FetchResponse)
	err := c.cc.Invoke(ctx, "/Lesslog/Fetch", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lesslogClient) Watch(ctx context.Context, in *FetchRequest, opts ...grpc.CallOption) (Lesslog_WatchClient, error) {
	stream, err := c.cc.NewStream(ctx, &Lesslog_ServiceDesc.Streams[0], "/Lesslog/Watch", opts...)
	if err != nil {
		return nil, err
	}
	x := &lesslogWatchClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Lesslog_WatchClient interface {
	Recv() (*FetchResponse, error)
	grpc.ClientStream
}

type lesslogWatchClient struct {
	grpc.ClientStream
}

func (x *lesslogWatchClient) Recv() (*FetchResponse, error) {
	m := new(FetchResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *lesslogClient) Push(ctx context.Context, in *PushRequest, opts ...grpc.CallOption) (*PushResponse, error) {
	out := new(PushResponse)
	err := c.cc.Invoke(ctx, "/Lesslog/Push", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lesslogClient) Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*PushResponse, error) {
	out := new(PushResponse)
	err := c.cc.Invoke(ctx, "/Lesslog/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LesslogServer is the server API for Lesslog service.
// All implementations must embed UnimplementedLesslogServer
// for forward compatibility
type LesslogServer interface {
	Fetch(context.Context, *FetchRequest) (*FetchResponse, error)
	Watch(*FetchRequest, Lesslog_WatchServer) error
	Push(context.Context, *PushRequest) (*PushResponse, error)
	Create(context.Context, *CreateRequest) (*PushResponse, error)
	mustEmbedUnimplementedLesslogServer()
}

// UnimplementedLesslogServer must be embedded to have forward compatible implementations.
type UnimplementedLesslogServer struct {
}

func (UnimplementedLesslogServer) Fetch(context.Context, *FetchRequest) (*FetchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Fetch not implemented")
}
func (UnimplementedLesslogServer) Watch(*FetchRequest, Lesslog_WatchServer) error {
	return status.Errorf(codes.Unimplemented, "method Watch not implemented")
}
func (UnimplementedLesslogServer) Push(context.Context, *PushRequest) (*PushResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Push not implemented")
}
func (UnimplementedLesslogServer) Create(context.Context, *CreateRequest) (*PushResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedLesslogServer) mustEmbedUnimplementedLesslogServer() {}

// UnsafeLesslogServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LesslogServer will
// result in compilation errors.
type UnsafeLesslogServer interface {
	mustEmbedUnimplementedLesslogServer()
}

func RegisterLesslogServer(s grpc.ServiceRegistrar, srv LesslogServer) {
	s.RegisterService(&Lesslog_ServiceDesc, srv)
}

func _Lesslog_Fetch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FetchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LesslogServer).Fetch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Lesslog/Fetch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LesslogServer).Fetch(ctx, req.(*FetchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Lesslog_Watch_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(FetchRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(LesslogServer).Watch(m, &lesslogWatchServer{stream})
}

type Lesslog_WatchServer interface {
	Send(*FetchResponse) error
	grpc.ServerStream
}

type lesslogWatchServer struct {
	grpc.ServerStream
}

func (x *lesslogWatchServer) Send(m *FetchResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _Lesslog_Push_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PushRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LesslogServer).Push(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Lesslog/Push",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LesslogServer).Push(ctx, req.(*PushRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Lesslog_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LesslogServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Lesslog/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LesslogServer).Create(ctx, req.(*CreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Lesslog_ServiceDesc is the grpc.ServiceDesc for Lesslog service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Lesslog_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Lesslog",
	HandlerType: (*LesslogServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Fetch",
			Handler:    _Lesslog_Fetch_Handler,
		},
		{
			MethodName: "Push",
			Handler:    _Lesslog_Push_Handler,
		},
		{
			MethodName: "Create",
			Handler:    _Lesslog_Create_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Watch",
			Handler:       _Lesslog_Watch_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "lesslog.proto",
}
