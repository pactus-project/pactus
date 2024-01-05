// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: transaction.proto

package pactus

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

// TransactionClient is the client API for Transaction service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TransactionClient interface {
	GetTransaction(ctx context.Context, in *GetTransactionRequest, opts ...grpc.CallOption) (*GetTransactionResponse, error)
	CalculateFee(ctx context.Context, in *CalculateFeeRequest, opts ...grpc.CallOption) (*CalculateFeeResponse, error)
	SendRawTransaction(ctx context.Context, in *SendRawTransactionRequest, opts ...grpc.CallOption) (*SendRawTransactionResponse, error)
	GetRawTransferTransaction(ctx context.Context, in *GetRawTransferTransactionRequest, opts ...grpc.CallOption) (*GetRawTransactionResponse, error)
	GetRawBondTransaction(ctx context.Context, in *GetRawBondTransactionRequest, opts ...grpc.CallOption) (*GetRawTransactionResponse, error)
	GetRawUnBondTransaction(ctx context.Context, in *GetRawUnBondTransactionRequest, opts ...grpc.CallOption) (*GetRawTransactionResponse, error)
	GetRawWithdrawTransaction(ctx context.Context, in *GetRawWithdrawTransactionRequest, opts ...grpc.CallOption) (*GetRawTransactionResponse, error)
}

type transactionClient struct {
	cc grpc.ClientConnInterface
}

func NewTransactionClient(cc grpc.ClientConnInterface) TransactionClient {
	return &transactionClient{cc}
}

func (c *transactionClient) GetTransaction(ctx context.Context, in *GetTransactionRequest, opts ...grpc.CallOption) (*GetTransactionResponse, error) {
	out := new(GetTransactionResponse)
	err := c.cc.Invoke(ctx, "/pactus.Transaction/GetTransaction", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *transactionClient) CalculateFee(ctx context.Context, in *CalculateFeeRequest, opts ...grpc.CallOption) (*CalculateFeeResponse, error) {
	out := new(CalculateFeeResponse)
	err := c.cc.Invoke(ctx, "/pactus.Transaction/CalculateFee", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *transactionClient) SendRawTransaction(ctx context.Context, in *SendRawTransactionRequest, opts ...grpc.CallOption) (*SendRawTransactionResponse, error) {
	out := new(SendRawTransactionResponse)
	err := c.cc.Invoke(ctx, "/pactus.Transaction/SendRawTransaction", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *transactionClient) GetRawTransferTransaction(ctx context.Context, in *GetRawTransferTransactionRequest, opts ...grpc.CallOption) (*GetRawTransactionResponse, error) {
	out := new(GetRawTransactionResponse)
	err := c.cc.Invoke(ctx, "/pactus.Transaction/GetRawTransferTransaction", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *transactionClient) GetRawBondTransaction(ctx context.Context, in *GetRawBondTransactionRequest, opts ...grpc.CallOption) (*GetRawTransactionResponse, error) {
	out := new(GetRawTransactionResponse)
	err := c.cc.Invoke(ctx, "/pactus.Transaction/GetRawBondTransaction", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *transactionClient) GetRawUnBondTransaction(ctx context.Context, in *GetRawUnBondTransactionRequest, opts ...grpc.CallOption) (*GetRawTransactionResponse, error) {
	out := new(GetRawTransactionResponse)
	err := c.cc.Invoke(ctx, "/pactus.Transaction/GetRawUnBondTransaction", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *transactionClient) GetRawWithdrawTransaction(ctx context.Context, in *GetRawWithdrawTransactionRequest, opts ...grpc.CallOption) (*GetRawTransactionResponse, error) {
	out := new(GetRawTransactionResponse)
	err := c.cc.Invoke(ctx, "/pactus.Transaction/GetRawWithdrawTransaction", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TransactionServer is the server API for Transaction service.
// All implementations should embed UnimplementedTransactionServer
// for forward compatibility
type TransactionServer interface {
	GetTransaction(context.Context, *GetTransactionRequest) (*GetTransactionResponse, error)
	CalculateFee(context.Context, *CalculateFeeRequest) (*CalculateFeeResponse, error)
	SendRawTransaction(context.Context, *SendRawTransactionRequest) (*SendRawTransactionResponse, error)
	GetRawTransferTransaction(context.Context, *GetRawTransferTransactionRequest) (*GetRawTransactionResponse, error)
	GetRawBondTransaction(context.Context, *GetRawBondTransactionRequest) (*GetRawTransactionResponse, error)
	GetRawUnBondTransaction(context.Context, *GetRawUnBondTransactionRequest) (*GetRawTransactionResponse, error)
	GetRawWithdrawTransaction(context.Context, *GetRawWithdrawTransactionRequest) (*GetRawTransactionResponse, error)
}

// UnimplementedTransactionServer should be embedded to have forward compatible implementations.
type UnimplementedTransactionServer struct {
}

func (UnimplementedTransactionServer) GetTransaction(context.Context, *GetTransactionRequest) (*GetTransactionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTransaction not implemented")
}
func (UnimplementedTransactionServer) CalculateFee(context.Context, *CalculateFeeRequest) (*CalculateFeeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CalculateFee not implemented")
}
func (UnimplementedTransactionServer) SendRawTransaction(context.Context, *SendRawTransactionRequest) (*SendRawTransactionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendRawTransaction not implemented")
}
func (UnimplementedTransactionServer) GetRawTransferTransaction(context.Context, *GetRawTransferTransactionRequest) (*GetRawTransactionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRawTransferTransaction not implemented")
}
func (UnimplementedTransactionServer) GetRawBondTransaction(context.Context, *GetRawBondTransactionRequest) (*GetRawTransactionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRawBondTransaction not implemented")
}
func (UnimplementedTransactionServer) GetRawUnBondTransaction(context.Context, *GetRawUnBondTransactionRequest) (*GetRawTransactionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRawUnBondTransaction not implemented")
}
func (UnimplementedTransactionServer) GetRawWithdrawTransaction(context.Context, *GetRawWithdrawTransactionRequest) (*GetRawTransactionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRawWithdrawTransaction not implemented")
}

// UnsafeTransactionServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TransactionServer will
// result in compilation errors.
type UnsafeTransactionServer interface {
	mustEmbedUnimplementedTransactionServer()
}

func RegisterTransactionServer(s grpc.ServiceRegistrar, srv TransactionServer) {
	s.RegisterService(&Transaction_ServiceDesc, srv)
}

func _Transaction_GetTransaction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTransactionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TransactionServer).GetTransaction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pactus.Transaction/GetTransaction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TransactionServer).GetTransaction(ctx, req.(*GetTransactionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Transaction_CalculateFee_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CalculateFeeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TransactionServer).CalculateFee(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pactus.Transaction/CalculateFee",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TransactionServer).CalculateFee(ctx, req.(*CalculateFeeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Transaction_SendRawTransaction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendRawTransactionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TransactionServer).SendRawTransaction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pactus.Transaction/SendRawTransaction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TransactionServer).SendRawTransaction(ctx, req.(*SendRawTransactionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Transaction_GetRawTransferTransaction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRawTransferTransactionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TransactionServer).GetRawTransferTransaction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pactus.Transaction/GetRawTransferTransaction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TransactionServer).GetRawTransferTransaction(ctx, req.(*GetRawTransferTransactionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Transaction_GetRawBondTransaction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRawBondTransactionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TransactionServer).GetRawBondTransaction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pactus.Transaction/GetRawBondTransaction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TransactionServer).GetRawBondTransaction(ctx, req.(*GetRawBondTransactionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Transaction_GetRawUnBondTransaction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRawUnBondTransactionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TransactionServer).GetRawUnBondTransaction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pactus.Transaction/GetRawUnBondTransaction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TransactionServer).GetRawUnBondTransaction(ctx, req.(*GetRawUnBondTransactionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Transaction_GetRawWithdrawTransaction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRawWithdrawTransactionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TransactionServer).GetRawWithdrawTransaction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pactus.Transaction/GetRawWithdrawTransaction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TransactionServer).GetRawWithdrawTransaction(ctx, req.(*GetRawWithdrawTransactionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Transaction_ServiceDesc is the grpc.ServiceDesc for Transaction service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Transaction_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pactus.Transaction",
	HandlerType: (*TransactionServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetTransaction",
			Handler:    _Transaction_GetTransaction_Handler,
		},
		{
			MethodName: "CalculateFee",
			Handler:    _Transaction_CalculateFee_Handler,
		},
		{
			MethodName: "SendRawTransaction",
			Handler:    _Transaction_SendRawTransaction_Handler,
		},
		{
			MethodName: "GetRawTransferTransaction",
			Handler:    _Transaction_GetRawTransferTransaction_Handler,
		},
		{
			MethodName: "GetRawBondTransaction",
			Handler:    _Transaction_GetRawBondTransaction_Handler,
		},
		{
			MethodName: "GetRawUnBondTransaction",
			Handler:    _Transaction_GetRawUnBondTransaction_Handler,
		},
		{
			MethodName: "GetRawWithdrawTransaction",
			Handler:    _Transaction_GetRawWithdrawTransaction_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "transaction.proto",
}
