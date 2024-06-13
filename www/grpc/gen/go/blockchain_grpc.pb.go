// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             (unknown)
// source: blockchain.proto

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
	Blockchain_GetBlock_FullMethodName              = "/pactus.Blockchain/GetBlock"
	Blockchain_GetBlockHash_FullMethodName          = "/pactus.Blockchain/GetBlockHash"
	Blockchain_GetBlockHeight_FullMethodName        = "/pactus.Blockchain/GetBlockHeight"
	Blockchain_GetBlockchainInfo_FullMethodName     = "/pactus.Blockchain/GetBlockchainInfo"
	Blockchain_GetConsensusInfo_FullMethodName      = "/pactus.Blockchain/GetConsensusInfo"
	Blockchain_GetAccount_FullMethodName            = "/pactus.Blockchain/GetAccount"
	Blockchain_GetValidator_FullMethodName          = "/pactus.Blockchain/GetValidator"
	Blockchain_GetValidatorByNumber_FullMethodName  = "/pactus.Blockchain/GetValidatorByNumber"
	Blockchain_GetValidatorAddresses_FullMethodName = "/pactus.Blockchain/GetValidatorAddresses"
	Blockchain_GetPublicKey_FullMethodName          = "/pactus.Blockchain/GetPublicKey"
)

// BlockchainClient is the client API for Blockchain service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// Blockchain service defines RPC methods for interacting with the blockchain.
type BlockchainClient interface {
	// GetBlock retrieves information about a block based on the provided request
	// parameters.
	GetBlock(ctx context.Context, in *GetBlockRequest, opts ...grpc.CallOption) (*GetBlockResponse, error)
	// GetBlockHash retrieves the hash of a block at the specified height.
	GetBlockHash(ctx context.Context, in *GetBlockHashRequest, opts ...grpc.CallOption) (*GetBlockHashResponse, error)
	// GetBlockHeight retrieves the height of a block with the specified hash.
	GetBlockHeight(ctx context.Context, in *GetBlockHeightRequest, opts ...grpc.CallOption) (*GetBlockHeightResponse, error)
	// GetBlockchainInfo retrieves general information about the blockchain.
	GetBlockchainInfo(ctx context.Context, in *GetBlockchainInfoRequest, opts ...grpc.CallOption) (*GetBlockchainInfoResponse, error)
	// GetConsensusInfo retrieves information about the consensus instances.
	GetConsensusInfo(ctx context.Context, in *GetConsensusInfoRequest, opts ...grpc.CallOption) (*GetConsensusInfoResponse, error)
	// GetAccount retrieves information about an account based on the provided
	// address.
	GetAccount(ctx context.Context, in *GetAccountRequest, opts ...grpc.CallOption) (*GetAccountResponse, error)
	// GetValidator retrieves information about a validator based on the provided
	// address.
	GetValidator(ctx context.Context, in *GetValidatorRequest, opts ...grpc.CallOption) (*GetValidatorResponse, error)
	// GetValidatorByNumber retrieves information about a validator based on the
	// provided number.
	GetValidatorByNumber(ctx context.Context, in *GetValidatorByNumberRequest, opts ...grpc.CallOption) (*GetValidatorResponse, error)
	// GetValidatorAddresses retrieves a list of all validator addresses.
	GetValidatorAddresses(ctx context.Context, in *GetValidatorAddressesRequest, opts ...grpc.CallOption) (*GetValidatorAddressesResponse, error)
	// GetPublicKey retrieves the public key of an account based on the provided
	// address.
	GetPublicKey(ctx context.Context, in *GetPublicKeyRequest, opts ...grpc.CallOption) (*GetPublicKeyResponse, error)
}

type blockchainClient struct {
	cc grpc.ClientConnInterface
}

func NewBlockchainClient(cc grpc.ClientConnInterface) BlockchainClient {
	return &blockchainClient{cc}
}

func (c *blockchainClient) GetBlock(ctx context.Context, in *GetBlockRequest, opts ...grpc.CallOption) (*GetBlockResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetBlockResponse)
	err := c.cc.Invoke(ctx, Blockchain_GetBlock_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *blockchainClient) GetBlockHash(ctx context.Context, in *GetBlockHashRequest, opts ...grpc.CallOption) (*GetBlockHashResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetBlockHashResponse)
	err := c.cc.Invoke(ctx, Blockchain_GetBlockHash_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *blockchainClient) GetBlockHeight(ctx context.Context, in *GetBlockHeightRequest, opts ...grpc.CallOption) (*GetBlockHeightResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetBlockHeightResponse)
	err := c.cc.Invoke(ctx, Blockchain_GetBlockHeight_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *blockchainClient) GetBlockchainInfo(ctx context.Context, in *GetBlockchainInfoRequest, opts ...grpc.CallOption) (*GetBlockchainInfoResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetBlockchainInfoResponse)
	err := c.cc.Invoke(ctx, Blockchain_GetBlockchainInfo_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *blockchainClient) GetConsensusInfo(ctx context.Context, in *GetConsensusInfoRequest, opts ...grpc.CallOption) (*GetConsensusInfoResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetConsensusInfoResponse)
	err := c.cc.Invoke(ctx, Blockchain_GetConsensusInfo_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *blockchainClient) GetAccount(ctx context.Context, in *GetAccountRequest, opts ...grpc.CallOption) (*GetAccountResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetAccountResponse)
	err := c.cc.Invoke(ctx, Blockchain_GetAccount_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *blockchainClient) GetValidator(ctx context.Context, in *GetValidatorRequest, opts ...grpc.CallOption) (*GetValidatorResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetValidatorResponse)
	err := c.cc.Invoke(ctx, Blockchain_GetValidator_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *blockchainClient) GetValidatorByNumber(ctx context.Context, in *GetValidatorByNumberRequest, opts ...grpc.CallOption) (*GetValidatorResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetValidatorResponse)
	err := c.cc.Invoke(ctx, Blockchain_GetValidatorByNumber_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *blockchainClient) GetValidatorAddresses(ctx context.Context, in *GetValidatorAddressesRequest, opts ...grpc.CallOption) (*GetValidatorAddressesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetValidatorAddressesResponse)
	err := c.cc.Invoke(ctx, Blockchain_GetValidatorAddresses_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *blockchainClient) GetPublicKey(ctx context.Context, in *GetPublicKeyRequest, opts ...grpc.CallOption) (*GetPublicKeyResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetPublicKeyResponse)
	err := c.cc.Invoke(ctx, Blockchain_GetPublicKey_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BlockchainServer is the server API for Blockchain service.
// All implementations should embed UnimplementedBlockchainServer
// for forward compatibility
//
// Blockchain service defines RPC methods for interacting with the blockchain.
type BlockchainServer interface {
	// GetBlock retrieves information about a block based on the provided request
	// parameters.
	GetBlock(context.Context, *GetBlockRequest) (*GetBlockResponse, error)
	// GetBlockHash retrieves the hash of a block at the specified height.
	GetBlockHash(context.Context, *GetBlockHashRequest) (*GetBlockHashResponse, error)
	// GetBlockHeight retrieves the height of a block with the specified hash.
	GetBlockHeight(context.Context, *GetBlockHeightRequest) (*GetBlockHeightResponse, error)
	// GetBlockchainInfo retrieves general information about the blockchain.
	GetBlockchainInfo(context.Context, *GetBlockchainInfoRequest) (*GetBlockchainInfoResponse, error)
	// GetConsensusInfo retrieves information about the consensus instances.
	GetConsensusInfo(context.Context, *GetConsensusInfoRequest) (*GetConsensusInfoResponse, error)
	// GetAccount retrieves information about an account based on the provided
	// address.
	GetAccount(context.Context, *GetAccountRequest) (*GetAccountResponse, error)
	// GetValidator retrieves information about a validator based on the provided
	// address.
	GetValidator(context.Context, *GetValidatorRequest) (*GetValidatorResponse, error)
	// GetValidatorByNumber retrieves information about a validator based on the
	// provided number.
	GetValidatorByNumber(context.Context, *GetValidatorByNumberRequest) (*GetValidatorResponse, error)
	// GetValidatorAddresses retrieves a list of all validator addresses.
	GetValidatorAddresses(context.Context, *GetValidatorAddressesRequest) (*GetValidatorAddressesResponse, error)
	// GetPublicKey retrieves the public key of an account based on the provided
	// address.
	GetPublicKey(context.Context, *GetPublicKeyRequest) (*GetPublicKeyResponse, error)
}

// UnimplementedBlockchainServer should be embedded to have forward compatible implementations.
type UnimplementedBlockchainServer struct {
}

func (UnimplementedBlockchainServer) GetBlock(context.Context, *GetBlockRequest) (*GetBlockResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBlock not implemented")
}
func (UnimplementedBlockchainServer) GetBlockHash(context.Context, *GetBlockHashRequest) (*GetBlockHashResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBlockHash not implemented")
}
func (UnimplementedBlockchainServer) GetBlockHeight(context.Context, *GetBlockHeightRequest) (*GetBlockHeightResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBlockHeight not implemented")
}
func (UnimplementedBlockchainServer) GetBlockchainInfo(context.Context, *GetBlockchainInfoRequest) (*GetBlockchainInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBlockchainInfo not implemented")
}
func (UnimplementedBlockchainServer) GetConsensusInfo(context.Context, *GetConsensusInfoRequest) (*GetConsensusInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetConsensusInfo not implemented")
}
func (UnimplementedBlockchainServer) GetAccount(context.Context, *GetAccountRequest) (*GetAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAccount not implemented")
}
func (UnimplementedBlockchainServer) GetValidator(context.Context, *GetValidatorRequest) (*GetValidatorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetValidator not implemented")
}
func (UnimplementedBlockchainServer) GetValidatorByNumber(context.Context, *GetValidatorByNumberRequest) (*GetValidatorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetValidatorByNumber not implemented")
}
func (UnimplementedBlockchainServer) GetValidatorAddresses(context.Context, *GetValidatorAddressesRequest) (*GetValidatorAddressesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetValidatorAddresses not implemented")
}
func (UnimplementedBlockchainServer) GetPublicKey(context.Context, *GetPublicKeyRequest) (*GetPublicKeyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPublicKey not implemented")
}

// UnsafeBlockchainServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BlockchainServer will
// result in compilation errors.
type UnsafeBlockchainServer interface {
	mustEmbedUnimplementedBlockchainServer()
}

func RegisterBlockchainServer(s grpc.ServiceRegistrar, srv BlockchainServer) {
	s.RegisterService(&Blockchain_ServiceDesc, srv)
}

func _Blockchain_GetBlock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBlockRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlockchainServer).GetBlock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Blockchain_GetBlock_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlockchainServer).GetBlock(ctx, req.(*GetBlockRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Blockchain_GetBlockHash_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBlockHashRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlockchainServer).GetBlockHash(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Blockchain_GetBlockHash_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlockchainServer).GetBlockHash(ctx, req.(*GetBlockHashRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Blockchain_GetBlockHeight_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBlockHeightRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlockchainServer).GetBlockHeight(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Blockchain_GetBlockHeight_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlockchainServer).GetBlockHeight(ctx, req.(*GetBlockHeightRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Blockchain_GetBlockchainInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBlockchainInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlockchainServer).GetBlockchainInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Blockchain_GetBlockchainInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlockchainServer).GetBlockchainInfo(ctx, req.(*GetBlockchainInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Blockchain_GetConsensusInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetConsensusInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlockchainServer).GetConsensusInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Blockchain_GetConsensusInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlockchainServer).GetConsensusInfo(ctx, req.(*GetConsensusInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Blockchain_GetAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlockchainServer).GetAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Blockchain_GetAccount_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlockchainServer).GetAccount(ctx, req.(*GetAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Blockchain_GetValidator_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetValidatorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlockchainServer).GetValidator(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Blockchain_GetValidator_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlockchainServer).GetValidator(ctx, req.(*GetValidatorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Blockchain_GetValidatorByNumber_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetValidatorByNumberRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlockchainServer).GetValidatorByNumber(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Blockchain_GetValidatorByNumber_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlockchainServer).GetValidatorByNumber(ctx, req.(*GetValidatorByNumberRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Blockchain_GetValidatorAddresses_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetValidatorAddressesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlockchainServer).GetValidatorAddresses(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Blockchain_GetValidatorAddresses_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlockchainServer).GetValidatorAddresses(ctx, req.(*GetValidatorAddressesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Blockchain_GetPublicKey_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPublicKeyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlockchainServer).GetPublicKey(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Blockchain_GetPublicKey_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlockchainServer).GetPublicKey(ctx, req.(*GetPublicKeyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Blockchain_ServiceDesc is the grpc.ServiceDesc for Blockchain service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Blockchain_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pactus.Blockchain",
	HandlerType: (*BlockchainServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetBlock",
			Handler:    _Blockchain_GetBlock_Handler,
		},
		{
			MethodName: "GetBlockHash",
			Handler:    _Blockchain_GetBlockHash_Handler,
		},
		{
			MethodName: "GetBlockHeight",
			Handler:    _Blockchain_GetBlockHeight_Handler,
		},
		{
			MethodName: "GetBlockchainInfo",
			Handler:    _Blockchain_GetBlockchainInfo_Handler,
		},
		{
			MethodName: "GetConsensusInfo",
			Handler:    _Blockchain_GetConsensusInfo_Handler,
		},
		{
			MethodName: "GetAccount",
			Handler:    _Blockchain_GetAccount_Handler,
		},
		{
			MethodName: "GetValidator",
			Handler:    _Blockchain_GetValidator_Handler,
		},
		{
			MethodName: "GetValidatorByNumber",
			Handler:    _Blockchain_GetValidatorByNumber_Handler,
		},
		{
			MethodName: "GetValidatorAddresses",
			Handler:    _Blockchain_GetValidatorAddresses_Handler,
		},
		{
			MethodName: "GetPublicKey",
			Handler:    _Blockchain_GetPublicKey_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "blockchain.proto",
}
