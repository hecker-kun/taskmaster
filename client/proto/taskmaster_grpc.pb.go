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

// TaskmasterClient is the client API for Taskmaster service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TaskmasterClient interface {
	CreateTask(ctx context.Context, in *AddTask, opts ...grpc.CallOption) (*Task, error)
	DeleteTask(ctx context.Context, in *DeleteParams, opts ...grpc.CallOption) (*Empty, error)
	//rpc getTask(getTaskReq) returns (getTaskRes);
	DeleteAllTasks(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	GetAllTasks(ctx context.Context, in *Empty, opts ...grpc.CallOption) (Taskmaster_GetAllTasksClient, error)
}

type taskmasterClient struct {
	cc grpc.ClientConnInterface
}

func NewTaskmasterClient(cc grpc.ClientConnInterface) TaskmasterClient {
	return &taskmasterClient{cc}
}

func (c *taskmasterClient) CreateTask(ctx context.Context, in *AddTask, opts ...grpc.CallOption) (*Task, error) {
	out := new(Task)
	err := c.cc.Invoke(ctx, "/proto.Taskmaster/createTask", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskmasterClient) DeleteTask(ctx context.Context, in *DeleteParams, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/proto.Taskmaster/deleteTask", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskmasterClient) DeleteAllTasks(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/proto.Taskmaster/deleteAllTasks", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskmasterClient) GetAllTasks(ctx context.Context, in *Empty, opts ...grpc.CallOption) (Taskmaster_GetAllTasksClient, error) {
	stream, err := c.cc.NewStream(ctx, &Taskmaster_ServiceDesc.Streams[0], "/proto.Taskmaster/getAllTasks", opts...)
	if err != nil {
		return nil, err
	}
	x := &taskmasterGetAllTasksClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Taskmaster_GetAllTasksClient interface {
	Recv() (*Task, error)
	grpc.ClientStream
}

type taskmasterGetAllTasksClient struct {
	grpc.ClientStream
}

func (x *taskmasterGetAllTasksClient) Recv() (*Task, error) {
	m := new(Task)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// TaskmasterServer is the server API for Taskmaster service.
// All implementations must embed UnimplementedTaskmasterServer
// for forward compatibility
type TaskmasterServer interface {
	CreateTask(context.Context, *AddTask) (*Task, error)
	DeleteTask(context.Context, *DeleteParams) (*Empty, error)
	//rpc getTask(getTaskReq) returns (getTaskRes);
	DeleteAllTasks(context.Context, *Empty) (*Empty, error)
	GetAllTasks(*Empty, Taskmaster_GetAllTasksServer) error
	mustEmbedUnimplementedTaskmasterServer()
}

// UnimplementedTaskmasterServer must be embedded to have forward compatible implementations.
type UnimplementedTaskmasterServer struct {
}

func (UnimplementedTaskmasterServer) CreateTask(context.Context, *AddTask) (*Task, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateTask not implemented")
}
func (UnimplementedTaskmasterServer) DeleteTask(context.Context, *DeleteParams) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteTask not implemented")
}
func (UnimplementedTaskmasterServer) DeleteAllTasks(context.Context, *Empty) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteAllTasks not implemented")
}
func (UnimplementedTaskmasterServer) GetAllTasks(*Empty, Taskmaster_GetAllTasksServer) error {
	return status.Errorf(codes.Unimplemented, "method GetAllTasks not implemented")
}
func (UnimplementedTaskmasterServer) mustEmbedUnimplementedTaskmasterServer() {}

// UnsafeTaskmasterServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TaskmasterServer will
// result in compilation errors.
type UnsafeTaskmasterServer interface {
	mustEmbedUnimplementedTaskmasterServer()
}

func RegisterTaskmasterServer(s grpc.ServiceRegistrar, srv TaskmasterServer) {
	s.RegisterService(&Taskmaster_ServiceDesc, srv)
}

func _Taskmaster_CreateTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddTask)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskmasterServer).CreateTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Taskmaster/createTask",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskmasterServer).CreateTask(ctx, req.(*AddTask))
	}
	return interceptor(ctx, in, info, handler)
}

func _Taskmaster_DeleteTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskmasterServer).DeleteTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Taskmaster/deleteTask",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskmasterServer).DeleteTask(ctx, req.(*DeleteParams))
	}
	return interceptor(ctx, in, info, handler)
}

func _Taskmaster_DeleteAllTasks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskmasterServer).DeleteAllTasks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Taskmaster/deleteAllTasks",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskmasterServer).DeleteAllTasks(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Taskmaster_GetAllTasks_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(TaskmasterServer).GetAllTasks(m, &taskmasterGetAllTasksServer{stream})
}

type Taskmaster_GetAllTasksServer interface {
	Send(*Task) error
	grpc.ServerStream
}

type taskmasterGetAllTasksServer struct {
	grpc.ServerStream
}

func (x *taskmasterGetAllTasksServer) Send(m *Task) error {
	return x.ServerStream.SendMsg(m)
}

// Taskmaster_ServiceDesc is the grpc.ServiceDesc for Taskmaster service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Taskmaster_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Taskmaster",
	HandlerType: (*TaskmasterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "createTask",
			Handler:    _Taskmaster_CreateTask_Handler,
		},
		{
			MethodName: "deleteTask",
			Handler:    _Taskmaster_DeleteTask_Handler,
		},
		{
			MethodName: "deleteAllTasks",
			Handler:    _Taskmaster_DeleteAllTasks_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "getAllTasks",
			Handler:       _Taskmaster_GetAllTasks_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "proto/taskmaster.proto",
}
