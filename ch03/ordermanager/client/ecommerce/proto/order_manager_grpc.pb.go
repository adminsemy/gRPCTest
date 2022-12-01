// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: proto/order_manager.proto

package ecommerce

import (
	context "context"
	wrappers "github.com/golang/protobuf/ptypes/wrappers"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// OrderManagerClient is the client API for OrderManager service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OrderManagerClient interface {
	GetOrder(ctx context.Context, in *wrappers.StringValue, opts ...grpc.CallOption) (*Order, error)
}

type orderManagerClient struct {
	cc grpc.ClientConnInterface
}

func NewOrderManagerClient(cc grpc.ClientConnInterface) OrderManagerClient {
	return &orderManagerClient{cc}
}

func (c *orderManagerClient) GetOrder(ctx context.Context, in *wrappers.StringValue, opts ...grpc.CallOption) (*Order, error) {
	out := new(Order)
	err := c.cc.Invoke(ctx, "/ecommerce.OrderManager/getOrder", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OrderManagerServer is the server API for OrderManager service.
// All implementations must embed UnimplementedOrderManagerServer
// for forward compatibility
type OrderManagerServer interface {
	GetOrder(context.Context, *wrappers.StringValue) (*Order, error)
	mustEmbedUnimplementedOrderManagerServer()
}

// UnimplementedOrderManagerServer must be embedded to have forward compatible implementations.
type UnimplementedOrderManagerServer struct {
}

func (UnimplementedOrderManagerServer) GetOrder(context.Context, *wrappers.StringValue) (*Order, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOrder not implemented")
}
func (UnimplementedOrderManagerServer) mustEmbedUnimplementedOrderManagerServer() {}

// UnsafeOrderManagerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OrderManagerServer will
// result in compilation errors.
type UnsafeOrderManagerServer interface {
	mustEmbedUnimplementedOrderManagerServer()
}

func RegisterOrderManagerServer(s grpc.ServiceRegistrar, srv OrderManagerServer) {
	s.RegisterService(&OrderManager_ServiceDesc, srv)
}

func _OrderManager_GetOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(wrappers.StringValue)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderManagerServer).GetOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ecommerce.OrderManager/getOrder",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderManagerServer).GetOrder(ctx, req.(*wrappers.StringValue))
	}
	return interceptor(ctx, in, info, handler)
}

// OrderManager_ServiceDesc is the grpc.ServiceDesc for OrderManager service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var OrderManager_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ecommerce.OrderManager",
	HandlerType: (*OrderManagerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "getOrder",
			Handler:    _OrderManager_GetOrder_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/order_manager.proto",
}
