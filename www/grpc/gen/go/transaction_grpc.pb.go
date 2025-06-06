// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
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
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Transaction_GetTransaction_FullMethodName                 = "/pactus.Transaction/GetTransaction"
	Transaction_CalculateFee_FullMethodName                   = "/pactus.Transaction/CalculateFee"
	Transaction_BroadcastTransaction_FullMethodName           = "/pactus.Transaction/BroadcastTransaction"
	Transaction_GetRawTransferTransaction_FullMethodName      = "/pactus.Transaction/GetRawTransferTransaction"
	Transaction_GetRawBondTransaction_FullMethodName          = "/pactus.Transaction/GetRawBondTransaction"
	Transaction_GetRawUnbondTransaction_FullMethodName        = "/pactus.Transaction/GetRawUnbondTransaction"
	Transaction_GetRawWithdrawTransaction_FullMethodName      = "/pactus.Transaction/GetRawWithdrawTransaction"
	Transaction_GetRawBatchTransferTransaction_FullMethodName = "/pactus.Transaction/GetRawBatchTransferTransaction"
	Transaction_DecodeRawTransaction_FullMethodName           = "/pactus.Transaction/DecodeRawTransaction"
)

// TransactionClient is the client API for Transaction service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// Transaction service defines various RPC methods for interacting with transactions.
type TransactionClient interface {
	// GetTransaction retrieves transaction details based on the provided request parameters.
	GetTransaction(ctx context.Context, in *GetTransactionRequest, opts ...grpc.CallOption) (*GetTransactionResponse, error)
	// CalculateFee calculates the transaction fee based on the specified amount and payload type.
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
	// GetRawBatchTransferTransaction retrieves raw details of batch transfer transaction.
	GetRawBatchTransferTransaction(ctx context.Context, in *GetRawBatchTransferTransactionRequest, opts ...grpc.CallOption) (*GetRawTransactionResponse, error)
	// DecodeRawTransaction accepts raw transaction and returns decoded transaction.
	DecodeRawTransaction(ctx context.Context, in *DecodeRawTransactionRequest, opts ...grpc.CallOption) (*DecodeRawTransactionResponse, error)
}

type transactionClient struct {
	cc grpc.ClientConnInterface
}

func NewTransactionClient(cc grpc.ClientConnInterface) TransactionClient {
	return &transactionClient{cc}
}

func (c *transactionClient) GetTransaction(ctx context.Context, in *GetTransactionRequest, opts ...grpc.CallOption) (*GetTransactionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetTransactionResponse)
	err := c.cc.Invoke(ctx, Transaction_GetTransaction_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *transactionClient) CalculateFee(ctx context.Context, in *CalculateFeeRequest, opts ...grpc.CallOption) (*CalculateFeeResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CalculateFeeResponse)
	err := c.cc.Invoke(ctx, Transaction_CalculateFee_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *transactionClient) BroadcastTransaction(ctx context.Context, in *BroadcastTransactionRequest, opts ...grpc.CallOption) (*BroadcastTransactionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BroadcastTransactionResponse)
	err := c.cc.Invoke(ctx, Transaction_BroadcastTransaction_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *transactionClient) GetRawTransferTransaction(ctx context.Context, in *GetRawTransferTransactionRequest, opts ...grpc.CallOption) (*GetRawTransactionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetRawTransactionResponse)
	err := c.cc.Invoke(ctx, Transaction_GetRawTransferTransaction_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *transactionClient) GetRawBondTransaction(ctx context.Context, in *GetRawBondTransactionRequest, opts ...grpc.CallOption) (*GetRawTransactionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetRawTransactionResponse)
	err := c.cc.Invoke(ctx, Transaction_GetRawBondTransaction_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *transactionClient) GetRawUnbondTransaction(ctx context.Context, in *GetRawUnbondTransactionRequest, opts ...grpc.CallOption) (*GetRawTransactionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetRawTransactionResponse)
	err := c.cc.Invoke(ctx, Transaction_GetRawUnbondTransaction_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *transactionClient) GetRawWithdrawTransaction(ctx context.Context, in *GetRawWithdrawTransactionRequest, opts ...grpc.CallOption) (*GetRawTransactionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetRawTransactionResponse)
	err := c.cc.Invoke(ctx, Transaction_GetRawWithdrawTransaction_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *transactionClient) GetRawBatchTransferTransaction(ctx context.Context, in *GetRawBatchTransferTransactionRequest, opts ...grpc.CallOption) (*GetRawTransactionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetRawTransactionResponse)
	err := c.cc.Invoke(ctx, Transaction_GetRawBatchTransferTransaction_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *transactionClient) DecodeRawTransaction(ctx context.Context, in *DecodeRawTransactionRequest, opts ...grpc.CallOption) (*DecodeRawTransactionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DecodeRawTransactionResponse)
	err := c.cc.Invoke(ctx, Transaction_DecodeRawTransaction_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TransactionServer is the server API for Transaction service.
// All implementations should embed UnimplementedTransactionServer
// for forward compatibility.
//
// Transaction service defines various RPC methods for interacting with transactions.
type TransactionServer interface {
	// GetTransaction retrieves transaction details based on the provided request parameters.
	GetTransaction(context.Context, *GetTransactionRequest) (*GetTransactionResponse, error)
	// CalculateFee calculates the transaction fee based on the specified amount and payload type.
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
	// GetRawBatchTransferTransaction retrieves raw details of batch transfer transaction.
	GetRawBatchTransferTransaction(context.Context, *GetRawBatchTransferTransactionRequest) (*GetRawTransactionResponse, error)
	// DecodeRawTransaction accepts raw transaction and returns decoded transaction.
	DecodeRawTransaction(context.Context, *DecodeRawTransactionRequest) (*DecodeRawTransactionResponse, error)
}

// UnimplementedTransactionServer should be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedTransactionServer struct{}

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
func (UnimplementedTransactionServer) GetRawBatchTransferTransaction(context.Context, *GetRawBatchTransferTransactionRequest) (*GetRawTransactionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRawBatchTransferTransaction not implemented")
}
func (UnimplementedTransactionServer) DecodeRawTransaction(context.Context, *DecodeRawTransactionRequest) (*DecodeRawTransactionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DecodeRawTransaction not implemented")
}
func (UnimplementedTransactionServer) testEmbeddedByValue() {}

// UnsafeTransactionServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TransactionServer will
// result in compilation errors.
type UnsafeTransactionServer interface {
	mustEmbedUnimplementedTransactionServer()
}

func RegisterTransactionServer(s grpc.ServiceRegistrar, srv TransactionServer) {
	// If the following call pancis, it indicates UnimplementedTransactionServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
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
		FullMethod: Transaction_GetTransaction_FullMethodName,
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
		FullMethod: Transaction_CalculateFee_FullMethodName,
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
		FullMethod: Transaction_BroadcastTransaction_FullMethodName,
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
		FullMethod: Transaction_GetRawTransferTransaction_FullMethodName,
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
		FullMethod: Transaction_GetRawBondTransaction_FullMethodName,
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
		FullMethod: Transaction_GetRawUnbondTransaction_FullMethodName,
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
		FullMethod: Transaction_GetRawWithdrawTransaction_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TransactionServer).GetRawWithdrawTransaction(ctx, req.(*GetRawWithdrawTransactionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Transaction_GetRawBatchTransferTransaction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRawBatchTransferTransactionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TransactionServer).GetRawBatchTransferTransaction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Transaction_GetRawBatchTransferTransaction_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TransactionServer).GetRawBatchTransferTransaction(ctx, req.(*GetRawBatchTransferTransactionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Transaction_DecodeRawTransaction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DecodeRawTransactionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TransactionServer).DecodeRawTransaction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Transaction_DecodeRawTransaction_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TransactionServer).DecodeRawTransaction(ctx, req.(*DecodeRawTransactionRequest))
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
		{
			MethodName: "GetRawBatchTransferTransaction",
			Handler:    _Transaction_GetRawBatchTransferTransaction_Handler,
		},
		{
			MethodName: "DecodeRawTransaction",
			Handler:    _Transaction_DecodeRawTransaction_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "transaction.proto",
}
