// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             (unknown)
// source: util.proto

package pactus

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	Util_SignMessageWithPrivateKey_FullMethodName = "/pactus.Util/SignMessageWithPrivateKey"
	Util_VerifyMessage_FullMethodName             = "/pactus.Util/VerifyMessage"
)

// UtilClient is the client API for Util service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// Util service defines various RPC methods for interacting with
// Utils.
type UtilClient interface {
	// SignMessageWithPrivateKey
	SignMessageWithPrivateKey(ctx context.Context, in *SignMessageWithPrivateKeyRequest, opts ...grpc.CallOption) (*SignMessageWithPrivateKeyResponse, error)
	// VerifyMessage
	VerifyMessage(ctx context.Context, in *VerifyMessageRequest, opts ...grpc.CallOption) (*VerifyMessageResponse, error)
}

type utilClient struct {
	cc grpc.ClientConnInterface
}

func NewUtilClient(cc grpc.ClientConnInterface) UtilClient {
	return &utilClient{cc}
}

func (c *utilClient) SignMessageWithPrivateKey(ctx context.Context, in *SignMessageWithPrivateKeyRequest, opts ...grpc.CallOption) (*SignMessageWithPrivateKeyResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SignMessageWithPrivateKeyResponse)
	err := c.cc.Invoke(ctx, Util_SignMessageWithPrivateKey_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *utilClient) VerifyMessage(ctx context.Context, in *VerifyMessageRequest, opts ...grpc.CallOption) (*VerifyMessageResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(VerifyMessageResponse)
	err := c.cc.Invoke(ctx, Util_VerifyMessage_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UtilServer is the server API for Util service.
// All implementations should embed UnimplementedUtilServer
// for forward compatibility
//
// Util service defines various RPC methods for interacting with
// Utils.
type UtilServer interface {
	// SignMessageWithPrivateKey
	SignMessageWithPrivateKey(context.Context, *SignMessageWithPrivateKeyRequest) (*SignMessageWithPrivateKeyResponse, error)
	// VerifyMessage
	VerifyMessage(context.Context, *VerifyMessageRequest) (*VerifyMessageResponse, error)
}

// UnimplementedUtilServer should be embedded to have forward compatible implementations.
type UnimplementedUtilServer struct {
}

func (UnimplementedUtilServer) SignMessageWithPrivateKey(context.Context, *SignMessageWithPrivateKeyRequest) (*SignMessageWithPrivateKeyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SignMessageWithPrivateKey not implemented")
}
func (UnimplementedUtilServer) VerifyMessage(context.Context, *VerifyMessageRequest) (*VerifyMessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VerifyMessage not implemented")
}

// UnsafeUtilServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UtilServer will
// result in compilation errors.
type UnsafeUtilServer interface {
	mustEmbedUnimplementedUtilServer()
}

func RegisterUtilServer(s grpc.ServiceRegistrar, srv UtilServer) {
	s.RegisterService(&Util_ServiceDesc, srv)
}

func _Util_SignMessageWithPrivateKey_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SignMessageWithPrivateKeyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UtilServer).SignMessageWithPrivateKey(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Util_SignMessageWithPrivateKey_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UtilServer).SignMessageWithPrivateKey(ctx, req.(*SignMessageWithPrivateKeyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Util_VerifyMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VerifyMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UtilServer).VerifyMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Util_VerifyMessage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UtilServer).VerifyMessage(ctx, req.(*VerifyMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Util_ServiceDesc is the grpc.ServiceDesc for Util service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Util_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pactus.Util",
	HandlerType: (*UtilServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SignMessageWithPrivateKey",
			Handler:    _Util_SignMessageWithPrivateKey_Handler,
		},
		{
			MethodName: "VerifyMessage",
			Handler:    _Util_VerifyMessage_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "util.proto",
}
