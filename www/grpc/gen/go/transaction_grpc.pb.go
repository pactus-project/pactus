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
	// GetTransaction retrieves transaction details based on the provided request
	// parameters.
	GetTransaction(ctx context.Context, in *GetTransactionRequest, opts ...grpc.CallOption) (*GetTransactionResponse, error)
	// CalculateFee calculates the transaction fee based on the specified amount
	// and payload type.
	CalculateFee(ctx context.Context, in *CalculateFeeRequest, opts ...grpc.CallOption) (*CalculateFeeResponse, error)
	// BroadcastTransaction broadcasts a signed transaction to the network.
	BroadcastTransaction(ctx context.Context, in *BroadcastTransactionRequest, opts ...grpc.CallOption) (*BroadcastTransactionResponse, error)
	// GetRawTransferTransaction retrieves raw details of a transfer transaction.
	GetRawTransferTransaction(ctx context.Context, in *GetRawTransferTransactionRequest, opts ...grpc.CallOption) (*GetRawTransactionResponse, error)
	// GetRawBondTransaction retrieves raw details of a bond transaction.
	GetRawBondTransaction(ctx context.Context, in *GetRawBondTransactionRequest, opts ...grpc.CallOption) (*GetRawTransactionResponse, error)
	// GetRawUnbondTransaction retrieves raw details of an unbond transaction.
	GetRawUnbondTransaction(ctx context.Context, in *GetRawUnbondTransactionRequest, opts ...grpc.CallOption) (*GetRawTransactionResponse, error)
	// GetRawWithdrawTransaction retrieves raw details of a withdraw transaction.
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

func (c *transactionClient) BroadcastTransaction(ctx context.Context, in *BroadcastTransactionRequest, opts ...grpc.CallOption) (*BroadcastTransactionResponse, error) {
	out := new(BroadcastTransactionResponse)
	err := c.cc.Invoke(ctx, "/pactus.Transaction/BroadcastTransaction", in, out, opts...)
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

func (c *transactionClient) GetRawUnbondTransaction(ctx context.Context, in *GetRawUnbondTransactionRequest, opts ...grpc.CallOption) (*GetRawTransactionResponse, error) {
	out := new(GetRawTransactionResponse)
	err := c.cc.Invoke(ctx, "/pactus.Transaction/GetRawUnbondTransaction", in, out, opts...)
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
	// GetTransaction retrieves transaction details based on the provided request
	// parameters.
	GetTransaction(context.Context, *GetTransactionRequest) (*GetTransactionResponse, error)
	// CalculateFee calculates the transaction fee based on the specified amount
	// and payload type.
	CalculateFee(context.Context, *CalculateFeeRequest) (*CalculateFeeResponse, error)
	// BroadcastTransaction broadcasts a signed transaction to the network.
	BroadcastTransaction(context.Context, *BroadcastTransactionRequest) (*BroadcastTransactionResponse, error)
	// GetRawTransferTransaction retrieves raw details of a transfer transaction.
	GetRawTransferTransaction(context.Context, *GetRawTransferTransactionRequest) (*GetRawTransactionResponse, error)
	// GetRawBondTransaction retrieves raw details of a bond transaction.
	GetRawBondTransaction(context.Context, *GetRawBondTransactionRequest) (*GetRawTransactionResponse, error)
	// GetRawUnbondTransaction retrieves raw details of an unbond transaction.
	GetRawUnbondTransaction(context.Context, *GetRawUnbondTransactionRequest) (*GetRawTransactionResponse, error)
	// GetRawWithdrawTransaction retrieves raw details of a withdraw transaction.
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
func (UnimplementedTransactionServer) BroadcastTransaction(context.Context, *BroadcastTransactionRequest) (*BroadcastTransactionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BroadcastTransaction not implemented")
}
func (UnimplementedTransactionServer) GetRawTransferTransaction(context.Context, *GetRawTransferTransactionRequest) (*GetRawTransactionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRawTransferTransaction not implemented")
}
func (UnimplementedTransactionServer) GetRawBondTransaction(context.Context, *GetRawBondTransactionRequest) (*GetRawTransactionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRawBondTransaction not implemented")
}
func (UnimplementedTransactionServer) GetRawUnbondTransaction(context.Context, *GetRawUnbondTransactionRequest) (*GetRawTransactionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRawUnbondTransaction not implemented")
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

func _Transaction_BroadcastTransaction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BroadcastTransactionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TransactionServer).BroadcastTransaction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pactus.Transaction/BroadcastTransaction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TransactionServer).BroadcastTransaction(ctx, req.(*BroadcastTransactionRequest))
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

func _Transaction_GetRawUnbondTransaction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRawUnbondTransactionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TransactionServer).GetRawUnbondTransaction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pactus.Transaction/GetRawUnbondTransaction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TransactionServer).GetRawUnbondTransaction(ctx, req.(*GetRawUnbondTransactionRequest))
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
			MethodName: "BroadcastTransaction",
			Handler:    _Transaction_BroadcastTransaction_Handler,
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
			MethodName: "GetRawUnbondTransaction",
			Handler:    _Transaction_GetRawUnbondTransaction_Handler,
		},
		{
			MethodName: "GetRawWithdrawTransaction",
			Handler:    _Transaction_GetRawWithdrawTransaction_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "transaction.proto",
}
