// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.22.0
// source: checkout.proto

package product

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// CheckoutServiceClient is the client API for CheckoutService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CheckoutServiceClient interface {
	AddToCart(ctx context.Context, in *AddToCartRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	DeleteFromCart(ctx context.Context, in *DeleteFromCartRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	ListCart(ctx context.Context, in *ListCartRequest, opts ...grpc.CallOption) (*ListCartResponse, error)
	Purchase(ctx context.Context, in *PurchaseRequest, opts ...grpc.CallOption) (*PurchaseResponse, error)
}

type checkoutServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCheckoutServiceClient(cc grpc.ClientConnInterface) CheckoutServiceClient {
	return &checkoutServiceClient{cc}
}

func (c *checkoutServiceClient) AddToCart(ctx context.Context, in *AddToCartRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/route256.checkout.CheckoutService/AddToCart", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *checkoutServiceClient) DeleteFromCart(ctx context.Context, in *DeleteFromCartRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/route256.checkout.CheckoutService/DeleteFromCart", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *checkoutServiceClient) ListCart(ctx context.Context, in *ListCartRequest, opts ...grpc.CallOption) (*ListCartResponse, error) {
	out := new(ListCartResponse)
	err := c.cc.Invoke(ctx, "/route256.checkout.CheckoutService/ListCart", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *checkoutServiceClient) Purchase(ctx context.Context, in *PurchaseRequest, opts ...grpc.CallOption) (*PurchaseResponse, error) {
	out := new(PurchaseResponse)
	err := c.cc.Invoke(ctx, "/route256.checkout.CheckoutService/Purchase", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CheckoutServiceServer is the server API for CheckoutService service.
// All implementations must embed UnimplementedCheckoutServiceServer
// for forward compatibility
type CheckoutServiceServer interface {
	AddToCart(context.Context, *AddToCartRequest) (*emptypb.Empty, error)
	DeleteFromCart(context.Context, *DeleteFromCartRequest) (*emptypb.Empty, error)
	ListCart(context.Context, *ListCartRequest) (*ListCartResponse, error)
	Purchase(context.Context, *PurchaseRequest) (*PurchaseResponse, error)
	mustEmbedUnimplementedCheckoutServiceServer()
}

// UnimplementedCheckoutServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCheckoutServiceServer struct {
}

func (UnimplementedCheckoutServiceServer) AddToCart(context.Context, *AddToCartRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddToCart not implemented")
}
func (UnimplementedCheckoutServiceServer) DeleteFromCart(context.Context, *DeleteFromCartRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteFromCart not implemented")
}
func (UnimplementedCheckoutServiceServer) ListCart(context.Context, *ListCartRequest) (*ListCartResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListCart not implemented")
}
func (UnimplementedCheckoutServiceServer) Purchase(context.Context, *PurchaseRequest) (*PurchaseResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Purchase not implemented")
}
func (UnimplementedCheckoutServiceServer) mustEmbedUnimplementedCheckoutServiceServer() {}

// UnsafeCheckoutServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CheckoutServiceServer will
// result in compilation errors.
type UnsafeCheckoutServiceServer interface {
	mustEmbedUnimplementedCheckoutServiceServer()
}

func RegisterCheckoutServiceServer(s grpc.ServiceRegistrar, srv CheckoutServiceServer) {
	s.RegisterService(&CheckoutService_ServiceDesc, srv)
}

func _CheckoutService_AddToCart_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddToCartRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CheckoutServiceServer).AddToCart(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/route256.checkout.CheckoutService/AddToCart",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CheckoutServiceServer).AddToCart(ctx, req.(*AddToCartRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CheckoutService_DeleteFromCart_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteFromCartRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CheckoutServiceServer).DeleteFromCart(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/route256.checkout.CheckoutService/DeleteFromCart",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CheckoutServiceServer).DeleteFromCart(ctx, req.(*DeleteFromCartRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CheckoutService_ListCart_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListCartRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CheckoutServiceServer).ListCart(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/route256.checkout.CheckoutService/ListCart",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CheckoutServiceServer).ListCart(ctx, req.(*ListCartRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CheckoutService_Purchase_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PurchaseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CheckoutServiceServer).Purchase(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/route256.checkout.CheckoutService/Purchase",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CheckoutServiceServer).Purchase(ctx, req.(*PurchaseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CheckoutService_ServiceDesc is the grpc.ServiceDesc for CheckoutService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CheckoutService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "route256.checkout.CheckoutService",
	HandlerType: (*CheckoutServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddToCart",
			Handler:    _CheckoutService_AddToCart_Handler,
		},
		{
			MethodName: "DeleteFromCart",
			Handler:    _CheckoutService_DeleteFromCart_Handler,
		},
		{
			MethodName: "ListCart",
			Handler:    _CheckoutService_ListCart_Handler,
		},
		{
			MethodName: "Purchase",
			Handler:    _CheckoutService_Purchase_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "checkout.proto",
}
