// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.1
// source: count.proto

package countclient

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
	Count_IncFollowerCount_FullMethodName = "/count.Count/IncFollowerCount"
	Count_DecFollowerCount_FullMethodName = "/count.Count/DecFollowerCount"
	Count_GetFollowerCount_FullMethodName = "/count.Count/GetFollowerCount"
)

// CountClient is the client API for Count service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CountClient interface {
	IncFollowerCount(ctx context.Context, in *IncFollowerCountRequest, opts ...grpc.CallOption) (*IncFollowerCountResponse, error)
	DecFollowerCount(ctx context.Context, in *DecFollowerCountRequest, opts ...grpc.CallOption) (*DecFollowerCountResponse, error)
	GetFollowerCount(ctx context.Context, in *GetFollowerCountRequest, opts ...grpc.CallOption) (*GetFollowerCountResponse, error)
}

type countClient struct {
	cc grpc.ClientConnInterface
}

func NewCountClient(cc grpc.ClientConnInterface) CountClient {
	return &countClient{cc}
}

func (c *countClient) IncFollowerCount(ctx context.Context, in *IncFollowerCountRequest, opts ...grpc.CallOption) (*IncFollowerCountResponse, error) {
	out := new(IncFollowerCountResponse)
	err := c.cc.Invoke(ctx, Count_IncFollowerCount_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *countClient) DecFollowerCount(ctx context.Context, in *DecFollowerCountRequest, opts ...grpc.CallOption) (*DecFollowerCountResponse, error) {
	out := new(DecFollowerCountResponse)
	err := c.cc.Invoke(ctx, Count_DecFollowerCount_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *countClient) GetFollowerCount(ctx context.Context, in *GetFollowerCountRequest, opts ...grpc.CallOption) (*GetFollowerCountResponse, error) {
	out := new(GetFollowerCountResponse)
	err := c.cc.Invoke(ctx, Count_GetFollowerCount_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CountServer is the server API for Count service.
// All implementations must embed UnimplementedCountServer
// for forward compatibility
type CountServer interface {
	IncFollowerCount(context.Context, *IncFollowerCountRequest) (*IncFollowerCountResponse, error)
	DecFollowerCount(context.Context, *DecFollowerCountRequest) (*DecFollowerCountResponse, error)
	GetFollowerCount(context.Context, *GetFollowerCountRequest) (*GetFollowerCountResponse, error)
	mustEmbedUnimplementedCountServer()
}

// UnimplementedCountServer must be embedded to have forward compatible implementations.
type UnimplementedCountServer struct {
}

func (UnimplementedCountServer) IncFollowerCount(context.Context, *IncFollowerCountRequest) (*IncFollowerCountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IncFollowerCount not implemented")
}
func (UnimplementedCountServer) DecFollowerCount(context.Context, *DecFollowerCountRequest) (*DecFollowerCountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DecFollowerCount not implemented")
}
func (UnimplementedCountServer) GetFollowerCount(context.Context, *GetFollowerCountRequest) (*GetFollowerCountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFollowerCount not implemented")
}
func (UnimplementedCountServer) mustEmbedUnimplementedCountServer() {}

// UnsafeCountServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CountServer will
// result in compilation errors.
type UnsafeCountServer interface {
	mustEmbedUnimplementedCountServer()
}

func RegisterCountServer(s grpc.ServiceRegistrar, srv CountServer) {
	s.RegisterService(&Count_ServiceDesc, srv)
}

func _Count_IncFollowerCount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IncFollowerCountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CountServer).IncFollowerCount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Count_IncFollowerCount_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CountServer).IncFollowerCount(ctx, req.(*IncFollowerCountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Count_DecFollowerCount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DecFollowerCountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CountServer).DecFollowerCount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Count_DecFollowerCount_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CountServer).DecFollowerCount(ctx, req.(*DecFollowerCountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Count_GetFollowerCount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFollowerCountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CountServer).GetFollowerCount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Count_GetFollowerCount_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CountServer).GetFollowerCount(ctx, req.(*GetFollowerCountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Count_ServiceDesc is the grpc.ServiceDesc for Count service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Count_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "count.Count",
	HandlerType: (*CountServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "IncFollowerCount",
			Handler:    _Count_IncFollowerCount_Handler,
		},
		{
			MethodName: "DecFollowerCount",
			Handler:    _Count_DecFollowerCount_Handler,
		},
		{
			MethodName: "GetFollowerCount",
			Handler:    _Count_GetFollowerCount_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "count.proto",
}
