// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        (unknown)
// source: wallet.proto

// Define the package and Go package path for the generated code.

package pactus

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Request message for creating a new wallet.
type CreateWalletRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Name of the new wallet.
	WalletName string `protobuf:"bytes,1,opt,name=wallet_name,json=walletName,proto3" json:"wallet_name,omitempty"`
	// Mnemonic for wallet recovery.
	Mnemonic string `protobuf:"bytes,2,opt,name=mnemonic,proto3" json:"mnemonic,omitempty"`
	// Language for the mnemonic.
	Language string `protobuf:"bytes,3,opt,name=language,proto3" json:"language,omitempty"`
	// Password for securing the wallet.
	Password string `protobuf:"bytes,4,opt,name=password,proto3" json:"password,omitempty"`
}

func (x *CreateWalletRequest) Reset() {
	*x = CreateWalletRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wallet_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateWalletRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateWalletRequest) ProtoMessage() {}

func (x *CreateWalletRequest) ProtoReflect() protoreflect.Message {
	mi := &file_wallet_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateWalletRequest.ProtoReflect.Descriptor instead.
func (*CreateWalletRequest) Descriptor() ([]byte, []int) {
	return file_wallet_proto_rawDescGZIP(), []int{0}
}

func (x *CreateWalletRequest) GetWalletName() string {
	if x != nil {
		return x.WalletName
	}
	return ""
}

func (x *CreateWalletRequest) GetMnemonic() string {
	if x != nil {
		return x.Mnemonic
	}
	return ""
}

func (x *CreateWalletRequest) GetLanguage() string {
	if x != nil {
		return x.Language
	}
	return ""
}

func (x *CreateWalletRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

// Response message containing the name of the created wallet.
type CreateWalletResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Name of the created wallet.
	WalletName string `protobuf:"bytes,1,opt,name=wallet_name,json=walletName,proto3" json:"wallet_name,omitempty"`
}

func (x *CreateWalletResponse) Reset() {
	*x = CreateWalletResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wallet_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateWalletResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateWalletResponse) ProtoMessage() {}

func (x *CreateWalletResponse) ProtoReflect() protoreflect.Message {
	mi := &file_wallet_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateWalletResponse.ProtoReflect.Descriptor instead.
func (*CreateWalletResponse) Descriptor() ([]byte, []int) {
	return file_wallet_proto_rawDescGZIP(), []int{1}
}

func (x *CreateWalletResponse) GetWalletName() string {
	if x != nil {
		return x.WalletName
	}
	return ""
}

// Request message for loading an existing wallet.
type LoadWalletRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Name of the wallet to load.
	WalletName string `protobuf:"bytes,1,opt,name=wallet_name,json=walletName,proto3" json:"wallet_name,omitempty"`
}

func (x *LoadWalletRequest) Reset() {
	*x = LoadWalletRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wallet_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LoadWalletRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoadWalletRequest) ProtoMessage() {}

func (x *LoadWalletRequest) ProtoReflect() protoreflect.Message {
	mi := &file_wallet_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoadWalletRequest.ProtoReflect.Descriptor instead.
func (*LoadWalletRequest) Descriptor() ([]byte, []int) {
	return file_wallet_proto_rawDescGZIP(), []int{2}
}

func (x *LoadWalletRequest) GetWalletName() string {
	if x != nil {
		return x.WalletName
	}
	return ""
}

// Response message containing the name of the loaded wallet.
type LoadWalletResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Name of the loaded wallet.
	WalletName string `protobuf:"bytes,1,opt,name=wallet_name,json=walletName,proto3" json:"wallet_name,omitempty"`
}

func (x *LoadWalletResponse) Reset() {
	*x = LoadWalletResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wallet_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LoadWalletResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoadWalletResponse) ProtoMessage() {}

func (x *LoadWalletResponse) ProtoReflect() protoreflect.Message {
	mi := &file_wallet_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoadWalletResponse.ProtoReflect.Descriptor instead.
func (*LoadWalletResponse) Descriptor() ([]byte, []int) {
	return file_wallet_proto_rawDescGZIP(), []int{3}
}

func (x *LoadWalletResponse) GetWalletName() string {
	if x != nil {
		return x.WalletName
	}
	return ""
}

// Request message for unloading a currently loaded wallet.
type UnloadWalletRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Name of the wallet to unload.
	WalletName string `protobuf:"bytes,1,opt,name=wallet_name,json=walletName,proto3" json:"wallet_name,omitempty"`
}

func (x *UnloadWalletRequest) Reset() {
	*x = UnloadWalletRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wallet_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UnloadWalletRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UnloadWalletRequest) ProtoMessage() {}

func (x *UnloadWalletRequest) ProtoReflect() protoreflect.Message {
	mi := &file_wallet_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UnloadWalletRequest.ProtoReflect.Descriptor instead.
func (*UnloadWalletRequest) Descriptor() ([]byte, []int) {
	return file_wallet_proto_rawDescGZIP(), []int{4}
}

func (x *UnloadWalletRequest) GetWalletName() string {
	if x != nil {
		return x.WalletName
	}
	return ""
}

// Response message containing the name of the unloaded wallet.
type UnloadWalletResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Name of the unloaded wallet.
	WalletName string `protobuf:"bytes,1,opt,name=wallet_name,json=walletName,proto3" json:"wallet_name,omitempty"`
}

func (x *UnloadWalletResponse) Reset() {
	*x = UnloadWalletResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wallet_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UnloadWalletResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UnloadWalletResponse) ProtoMessage() {}

func (x *UnloadWalletResponse) ProtoReflect() protoreflect.Message {
	mi := &file_wallet_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UnloadWalletResponse.ProtoReflect.Descriptor instead.
func (*UnloadWalletResponse) Descriptor() ([]byte, []int) {
	return file_wallet_proto_rawDescGZIP(), []int{5}
}

func (x *UnloadWalletResponse) GetWalletName() string {
	if x != nil {
		return x.WalletName
	}
	return ""
}

// Request message for locking a currently loaded wallet.
type LockWalletRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Name of the wallet to lock.
	WalletName string `protobuf:"bytes,1,opt,name=wallet_name,json=walletName,proto3" json:"wallet_name,omitempty"`
}

func (x *LockWalletRequest) Reset() {
	*x = LockWalletRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wallet_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LockWalletRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LockWalletRequest) ProtoMessage() {}

func (x *LockWalletRequest) ProtoReflect() protoreflect.Message {
	mi := &file_wallet_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LockWalletRequest.ProtoReflect.Descriptor instead.
func (*LockWalletRequest) Descriptor() ([]byte, []int) {
	return file_wallet_proto_rawDescGZIP(), []int{6}
}

func (x *LockWalletRequest) GetWalletName() string {
	if x != nil {
		return x.WalletName
	}
	return ""
}

// Response message containing the name of the locked wallet.
type LockWalletResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Name of the locked wallet.
	WalletName string `protobuf:"bytes,1,opt,name=wallet_name,json=walletName,proto3" json:"wallet_name,omitempty"`
}

func (x *LockWalletResponse) Reset() {
	*x = LockWalletResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wallet_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LockWalletResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LockWalletResponse) ProtoMessage() {}

func (x *LockWalletResponse) ProtoReflect() protoreflect.Message {
	mi := &file_wallet_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LockWalletResponse.ProtoReflect.Descriptor instead.
func (*LockWalletResponse) Descriptor() ([]byte, []int) {
	return file_wallet_proto_rawDescGZIP(), []int{7}
}

func (x *LockWalletResponse) GetWalletName() string {
	if x != nil {
		return x.WalletName
	}
	return ""
}

// Request message for obtaining the validator address associated with a public
// key.
type GetValidatorAddressRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Public key for which the validator address is requested.
	PublicKey string `protobuf:"bytes,1,opt,name=public_key,json=publicKey,proto3" json:"public_key,omitempty"`
}

func (x *GetValidatorAddressRequest) Reset() {
	*x = GetValidatorAddressRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wallet_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetValidatorAddressRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetValidatorAddressRequest) ProtoMessage() {}

func (x *GetValidatorAddressRequest) ProtoReflect() protoreflect.Message {
	mi := &file_wallet_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetValidatorAddressRequest.ProtoReflect.Descriptor instead.
func (*GetValidatorAddressRequest) Descriptor() ([]byte, []int) {
	return file_wallet_proto_rawDescGZIP(), []int{8}
}

func (x *GetValidatorAddressRequest) GetPublicKey() string {
	if x != nil {
		return x.PublicKey
	}
	return ""
}

// Response message containing the validator address corresponding to a public
// key.
type GetValidatorAddressResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Validator address associated with the public key.
	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
}

func (x *GetValidatorAddressResponse) Reset() {
	*x = GetValidatorAddressResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wallet_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetValidatorAddressResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetValidatorAddressResponse) ProtoMessage() {}

func (x *GetValidatorAddressResponse) ProtoReflect() protoreflect.Message {
	mi := &file_wallet_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetValidatorAddressResponse.ProtoReflect.Descriptor instead.
func (*GetValidatorAddressResponse) Descriptor() ([]byte, []int) {
	return file_wallet_proto_rawDescGZIP(), []int{9}
}

func (x *GetValidatorAddressResponse) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

// Request message for unlocking a wallet.
type UnlockWalletRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Name of the wallet to unlock.
	WalletName string `protobuf:"bytes,1,opt,name=wallet_name,json=walletName,proto3" json:"wallet_name,omitempty"`
	// Password for unlocking the wallet.
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	// Timeout duration for the unlocked state.
	Timeout int32 `protobuf:"varint,3,opt,name=timeout,proto3" json:"timeout,omitempty"`
}

func (x *UnlockWalletRequest) Reset() {
	*x = UnlockWalletRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wallet_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UnlockWalletRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UnlockWalletRequest) ProtoMessage() {}

func (x *UnlockWalletRequest) ProtoReflect() protoreflect.Message {
	mi := &file_wallet_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UnlockWalletRequest.ProtoReflect.Descriptor instead.
func (*UnlockWalletRequest) Descriptor() ([]byte, []int) {
	return file_wallet_proto_rawDescGZIP(), []int{10}
}

func (x *UnlockWalletRequest) GetWalletName() string {
	if x != nil {
		return x.WalletName
	}
	return ""
}

func (x *UnlockWalletRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *UnlockWalletRequest) GetTimeout() int32 {
	if x != nil {
		return x.Timeout
	}
	return 0
}

// Response message containing the name of the unlocked wallet.
type UnlockWalletResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Name of the unlocked wallet.
	WalletName string `protobuf:"bytes,1,opt,name=wallet_name,json=walletName,proto3" json:"wallet_name,omitempty"`
}

func (x *UnlockWalletResponse) Reset() {
	*x = UnlockWalletResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wallet_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UnlockWalletResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UnlockWalletResponse) ProtoMessage() {}

func (x *UnlockWalletResponse) ProtoReflect() protoreflect.Message {
	mi := &file_wallet_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UnlockWalletResponse.ProtoReflect.Descriptor instead.
func (*UnlockWalletResponse) Descriptor() ([]byte, []int) {
	return file_wallet_proto_rawDescGZIP(), []int{11}
}

func (x *UnlockWalletResponse) GetWalletName() string {
	if x != nil {
		return x.WalletName
	}
	return ""
}

// Request message for signing a raw transaction.
type SignRawTransactionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Name of the wallet used for signing.
	WalletName string `protobuf:"bytes,1,opt,name=wallet_name,json=walletName,proto3" json:"wallet_name,omitempty"`
	// Raw transaction data to be signed.
	RawTransaction []byte `protobuf:"bytes,2,opt,name=raw_transaction,json=rawTransaction,proto3" json:"raw_transaction,omitempty"`
	// Password for unlocking the wallet for signing.
	Password string `protobuf:"bytes,3,opt,name=password,proto3" json:"password,omitempty"`
}

func (x *SignRawTransactionRequest) Reset() {
	*x = SignRawTransactionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wallet_proto_msgTypes[12]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SignRawTransactionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SignRawTransactionRequest) ProtoMessage() {}

func (x *SignRawTransactionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_wallet_proto_msgTypes[12]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SignRawTransactionRequest.ProtoReflect.Descriptor instead.
func (*SignRawTransactionRequest) Descriptor() ([]byte, []int) {
	return file_wallet_proto_rawDescGZIP(), []int{12}
}

func (x *SignRawTransactionRequest) GetWalletName() string {
	if x != nil {
		return x.WalletName
	}
	return ""
}

func (x *SignRawTransactionRequest) GetRawTransaction() []byte {
	if x != nil {
		return x.RawTransaction
	}
	return nil
}

func (x *SignRawTransactionRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

// Response message containing the transaction ID and signed raw transaction.
type SignRawTransactionResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// ID of the signed transaction.
	TransactionId []byte `protobuf:"bytes,1,opt,name=transaction_id,json=transactionId,proto3" json:"transaction_id,omitempty"`
	// Signed raw transaction data.
	SignedRawTransaction []byte `protobuf:"bytes,2,opt,name=signed_raw_transaction,json=signedRawTransaction,proto3" json:"signed_raw_transaction,omitempty"`
}

func (x *SignRawTransactionResponse) Reset() {
	*x = SignRawTransactionResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wallet_proto_msgTypes[13]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SignRawTransactionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SignRawTransactionResponse) ProtoMessage() {}

func (x *SignRawTransactionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_wallet_proto_msgTypes[13]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SignRawTransactionResponse.ProtoReflect.Descriptor instead.
func (*SignRawTransactionResponse) Descriptor() ([]byte, []int) {
	return file_wallet_proto_rawDescGZIP(), []int{13}
}

func (x *SignRawTransactionResponse) GetTransactionId() []byte {
	if x != nil {
		return x.TransactionId
	}
	return nil
}

func (x *SignRawTransactionResponse) GetSignedRawTransaction() []byte {
	if x != nil {
		return x.SignedRawTransaction
	}
	return nil
}

var File_wallet_proto protoreflect.FileDescriptor

var file_wallet_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06,
	0x70, 0x61, 0x63, 0x74, 0x75, 0x73, 0x1a, 0x11, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x8a, 0x01, 0x0a, 0x13, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x4e, 0x61,
	0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x6d, 0x6e, 0x65, 0x6d, 0x6f, 0x6e, 0x69, 0x63, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6d, 0x6e, 0x65, 0x6d, 0x6f, 0x6e, 0x69, 0x63, 0x12, 0x1a,
	0x0a, 0x08, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61,
	0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x61,
	0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x22, 0x37, 0x0a, 0x14, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1f,
	0x0a, 0x0b, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0a, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x22,
	0x34, 0x0a, 0x11, 0x4c, 0x6f, 0x61, 0x64, 0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x5f, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x77, 0x61, 0x6c, 0x6c, 0x65,
	0x74, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x35, 0x0a, 0x12, 0x4c, 0x6f, 0x61, 0x64, 0x57, 0x61, 0x6c,
	0x6c, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x77,
	0x61, 0x6c, 0x6c, 0x65, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0a, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x36, 0x0a, 0x13,
	0x55, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x5f, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74,
	0x4e, 0x61, 0x6d, 0x65, 0x22, 0x37, 0x0a, 0x14, 0x55, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x57, 0x61,
	0x6c, 0x6c, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1f, 0x0a, 0x0b,
	0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0a, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x34, 0x0a,
	0x11, 0x4c, 0x6f, 0x63, 0x6b, 0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x5f, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x4e,
	0x61, 0x6d, 0x65, 0x22, 0x35, 0x0a, 0x12, 0x4c, 0x6f, 0x63, 0x6b, 0x57, 0x61, 0x6c, 0x6c, 0x65,
	0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x77, 0x61, 0x6c,
	0x6c, 0x65, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a,
	0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x3b, 0x0a, 0x1a, 0x47, 0x65,
	0x74, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73,
	0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x75, 0x62, 0x6c,
	0x69, 0x63, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x70, 0x75,
	0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x22, 0x37, 0x0a, 0x1b, 0x47, 0x65, 0x74, 0x56, 0x61,
	0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73,
	0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x22, 0x6c, 0x0a, 0x13, 0x55, 0x6e, 0x6c, 0x6f, 0x63, 0x6b, 0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x77, 0x61, 0x6c, 0x6c, 0x65,
	0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x77, 0x61,
	0x6c, 0x6c, 0x65, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73,
	0x77, 0x6f, 0x72, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73,
	0x77, 0x6f, 0x72, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x22, 0x37,
	0x0a, 0x14, 0x55, 0x6e, 0x6c, 0x6f, 0x63, 0x6b, 0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74,
	0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x77, 0x61, 0x6c,
	0x6c, 0x65, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x81, 0x01, 0x0a, 0x19, 0x53, 0x69, 0x67, 0x6e,
	0x52, 0x61, 0x77, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x5f,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x77, 0x61, 0x6c, 0x6c,
	0x65, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x27, 0x0a, 0x0f, 0x72, 0x61, 0x77, 0x5f, 0x74, 0x72,
	0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x0e, 0x72, 0x61, 0x77, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x22, 0x79, 0x0a, 0x1a, 0x53,
	0x69, 0x67, 0x6e, 0x52, 0x61, 0x77, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x25, 0x0a, 0x0e, 0x74, 0x72, 0x61,
	0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x0d, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64,
	0x12, 0x34, 0x0a, 0x16, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x5f, 0x72, 0x61, 0x77, 0x5f, 0x74,
	0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x14, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x52, 0x61, 0x77, 0x54, 0x72, 0x61, 0x6e, 0x73,
	0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x32, 0xb0, 0x04, 0x0a, 0x06, 0x57, 0x61, 0x6c, 0x6c, 0x65,
	0x74, 0x12, 0x49, 0x0a, 0x0c, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x57, 0x61, 0x6c, 0x6c, 0x65,
	0x74, 0x12, 0x1b, 0x2e, 0x70, 0x61, 0x63, 0x74, 0x75, 0x73, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c,
	0x2e, 0x70, 0x61, 0x63, 0x74, 0x75, 0x73, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x57, 0x61,
	0x6c, 0x6c, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x43, 0x0a, 0x0a,
	0x4c, 0x6f, 0x61, 0x64, 0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x12, 0x19, 0x2e, 0x70, 0x61, 0x63,
	0x74, 0x75, 0x73, 0x2e, 0x4c, 0x6f, 0x61, 0x64, 0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x70, 0x61, 0x63, 0x74, 0x75, 0x73, 0x2e, 0x4c,
	0x6f, 0x61, 0x64, 0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x49, 0x0a, 0x0c, 0x55, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x57, 0x61, 0x6c, 0x6c, 0x65,
	0x74, 0x12, 0x1b, 0x2e, 0x70, 0x61, 0x63, 0x74, 0x75, 0x73, 0x2e, 0x55, 0x6e, 0x6c, 0x6f, 0x61,
	0x64, 0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c,
	0x2e, 0x70, 0x61, 0x63, 0x74, 0x75, 0x73, 0x2e, 0x55, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x57, 0x61,
	0x6c, 0x6c, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x43, 0x0a, 0x0a,
	0x4c, 0x6f, 0x63, 0x6b, 0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x12, 0x19, 0x2e, 0x70, 0x61, 0x63,
	0x74, 0x75, 0x73, 0x2e, 0x4c, 0x6f, 0x63, 0x6b, 0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x70, 0x61, 0x63, 0x74, 0x75, 0x73, 0x2e, 0x4c,
	0x6f, 0x63, 0x6b, 0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x49, 0x0a, 0x0c, 0x55, 0x6e, 0x6c, 0x6f, 0x63, 0x6b, 0x57, 0x61, 0x6c, 0x6c, 0x65,
	0x74, 0x12, 0x1b, 0x2e, 0x70, 0x61, 0x63, 0x74, 0x75, 0x73, 0x2e, 0x55, 0x6e, 0x6c, 0x6f, 0x63,
	0x6b, 0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c,
	0x2e, 0x70, 0x61, 0x63, 0x74, 0x75, 0x73, 0x2e, 0x55, 0x6e, 0x6c, 0x6f, 0x63, 0x6b, 0x57, 0x61,
	0x6c, 0x6c, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x5b, 0x0a, 0x12,
	0x53, 0x69, 0x67, 0x6e, 0x52, 0x61, 0x77, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x21, 0x2e, 0x70, 0x61, 0x63, 0x74, 0x75, 0x73, 0x2e, 0x53, 0x69, 0x67, 0x6e,
	0x52, 0x61, 0x77, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e, 0x70, 0x61, 0x63, 0x74, 0x75, 0x73, 0x2e, 0x53,
	0x69, 0x67, 0x6e, 0x52, 0x61, 0x77, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x5e, 0x0a, 0x13, 0x47, 0x65, 0x74,
	0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x12, 0x22, 0x2e, 0x70, 0x61, 0x63, 0x74, 0x75, 0x73, 0x2e, 0x47, 0x65, 0x74, 0x56, 0x61, 0x6c,
	0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x23, 0x2e, 0x70, 0x61, 0x63, 0x74, 0x75, 0x73, 0x2e, 0x47, 0x65,
	0x74, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73,
	0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x41, 0x0a, 0x0d, 0x70, 0x61, 0x63,
	0x74, 0x75, 0x73, 0x2e, 0x77, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x5a, 0x30, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x70, 0x61, 0x63, 0x74, 0x75, 0x73, 0x2d, 0x70, 0x72,
	0x6f, 0x6a, 0x65, 0x63, 0x74, 0x2f, 0x70, 0x61, 0x63, 0x74, 0x75, 0x73, 0x2f, 0x77, 0x77, 0x77,
	0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x70, 0x61, 0x63, 0x74, 0x75, 0x73, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_wallet_proto_rawDescOnce sync.Once
	file_wallet_proto_rawDescData = file_wallet_proto_rawDesc
)

func file_wallet_proto_rawDescGZIP() []byte {
	file_wallet_proto_rawDescOnce.Do(func() {
		file_wallet_proto_rawDescData = protoimpl.X.CompressGZIP(file_wallet_proto_rawDescData)
	})
	return file_wallet_proto_rawDescData
}

var file_wallet_proto_msgTypes = make([]protoimpl.MessageInfo, 14)
var file_wallet_proto_goTypes = []interface{}{
	(*CreateWalletRequest)(nil),         // 0: pactus.CreateWalletRequest
	(*CreateWalletResponse)(nil),        // 1: pactus.CreateWalletResponse
	(*LoadWalletRequest)(nil),           // 2: pactus.LoadWalletRequest
	(*LoadWalletResponse)(nil),          // 3: pactus.LoadWalletResponse
	(*UnloadWalletRequest)(nil),         // 4: pactus.UnloadWalletRequest
	(*UnloadWalletResponse)(nil),        // 5: pactus.UnloadWalletResponse
	(*LockWalletRequest)(nil),           // 6: pactus.LockWalletRequest
	(*LockWalletResponse)(nil),          // 7: pactus.LockWalletResponse
	(*GetValidatorAddressRequest)(nil),  // 8: pactus.GetValidatorAddressRequest
	(*GetValidatorAddressResponse)(nil), // 9: pactus.GetValidatorAddressResponse
	(*UnlockWalletRequest)(nil),         // 10: pactus.UnlockWalletRequest
	(*UnlockWalletResponse)(nil),        // 11: pactus.UnlockWalletResponse
	(*SignRawTransactionRequest)(nil),   // 12: pactus.SignRawTransactionRequest
	(*SignRawTransactionResponse)(nil),  // 13: pactus.SignRawTransactionResponse
}
var file_wallet_proto_depIdxs = []int32{
	0,  // 0: pactus.Wallet.CreateWallet:input_type -> pactus.CreateWalletRequest
	2,  // 1: pactus.Wallet.LoadWallet:input_type -> pactus.LoadWalletRequest
	4,  // 2: pactus.Wallet.UnloadWallet:input_type -> pactus.UnloadWalletRequest
	6,  // 3: pactus.Wallet.LockWallet:input_type -> pactus.LockWalletRequest
	10, // 4: pactus.Wallet.UnlockWallet:input_type -> pactus.UnlockWalletRequest
	12, // 5: pactus.Wallet.SignRawTransaction:input_type -> pactus.SignRawTransactionRequest
	8,  // 6: pactus.Wallet.GetValidatorAddress:input_type -> pactus.GetValidatorAddressRequest
	1,  // 7: pactus.Wallet.CreateWallet:output_type -> pactus.CreateWalletResponse
	3,  // 8: pactus.Wallet.LoadWallet:output_type -> pactus.LoadWalletResponse
	5,  // 9: pactus.Wallet.UnloadWallet:output_type -> pactus.UnloadWalletResponse
	7,  // 10: pactus.Wallet.LockWallet:output_type -> pactus.LockWalletResponse
	11, // 11: pactus.Wallet.UnlockWallet:output_type -> pactus.UnlockWalletResponse
	13, // 12: pactus.Wallet.SignRawTransaction:output_type -> pactus.SignRawTransactionResponse
	9,  // 13: pactus.Wallet.GetValidatorAddress:output_type -> pactus.GetValidatorAddressResponse
	7,  // [7:14] is the sub-list for method output_type
	0,  // [0:7] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_wallet_proto_init() }
func file_wallet_proto_init() {
	if File_wallet_proto != nil {
		return
	}
	file_transaction_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_wallet_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateWalletRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_wallet_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateWalletResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_wallet_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LoadWalletRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_wallet_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LoadWalletResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_wallet_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UnloadWalletRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_wallet_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UnloadWalletResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_wallet_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LockWalletRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_wallet_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LockWalletResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_wallet_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetValidatorAddressRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_wallet_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetValidatorAddressResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_wallet_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UnlockWalletRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_wallet_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UnlockWalletResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_wallet_proto_msgTypes[12].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SignRawTransactionRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_wallet_proto_msgTypes[13].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SignRawTransactionResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_wallet_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   14,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_wallet_proto_goTypes,
		DependencyIndexes: file_wallet_proto_depIdxs,
		MessageInfos:      file_wallet_proto_msgTypes,
	}.Build()
	File_wallet_proto = out.File
	file_wallet_proto_rawDesc = nil
	file_wallet_proto_goTypes = nil
	file_wallet_proto_depIdxs = nil
}
