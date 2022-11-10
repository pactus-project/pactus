// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        (unknown)
// source: transaction.proto

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

type PayloadType int32

const (
	PayloadType_UNKNOWN           PayloadType = 0
	PayloadType_SEND_PAYLOAD      PayloadType = 1
	PayloadType_BOND_PAYLOAD      PayloadType = 2
	PayloadType_SORTITION_PAYLOAD PayloadType = 3
	PayloadType_UNBOND_PAYLOAD    PayloadType = 4
	PayloadType_WITHDRAW_PAYLOAD  PayloadType = 5
)

// Enum value maps for PayloadType.
var (
	PayloadType_name = map[int32]string{
		0: "UNKNOWN",
		1: "SEND_PAYLOAD",
		2: "BOND_PAYLOAD",
		3: "SORTITION_PAYLOAD",
		4: "UNBOND_PAYLOAD",
		5: "WITHDRAW_PAYLOAD",
	}
	PayloadType_value = map[string]int32{
		"UNKNOWN":           0,
		"SEND_PAYLOAD":      1,
		"BOND_PAYLOAD":      2,
		"SORTITION_PAYLOAD": 3,
		"UNBOND_PAYLOAD":    4,
		"WITHDRAW_PAYLOAD":  5,
	}
)

func (x PayloadType) Enum() *PayloadType {
	p := new(PayloadType)
	*p = x
	return p
}

func (x PayloadType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (PayloadType) Descriptor() protoreflect.EnumDescriptor {
	return file_transaction_proto_enumTypes[0].Descriptor()
}

func (PayloadType) Type() protoreflect.EnumType {
	return &file_transaction_proto_enumTypes[0]
}

func (x PayloadType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use PayloadType.Descriptor instead.
func (PayloadType) EnumDescriptor() ([]byte, []int) {
	return file_transaction_proto_rawDescGZIP(), []int{0}
}

type TransactionVerbosity int32

const (
	TransactionVerbosity_TRANSACTION_DATA TransactionVerbosity = 0
	TransactionVerbosity_TRANSACTION_INFO TransactionVerbosity = 1
)

// Enum value maps for TransactionVerbosity.
var (
	TransactionVerbosity_name = map[int32]string{
		0: "TRANSACTION_DATA",
		1: "TRANSACTION_INFO",
	}
	TransactionVerbosity_value = map[string]int32{
		"TRANSACTION_DATA": 0,
		"TRANSACTION_INFO": 1,
	}
)

func (x TransactionVerbosity) Enum() *TransactionVerbosity {
	p := new(TransactionVerbosity)
	*p = x
	return p
}

func (x TransactionVerbosity) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (TransactionVerbosity) Descriptor() protoreflect.EnumDescriptor {
	return file_transaction_proto_enumTypes[1].Descriptor()
}

func (TransactionVerbosity) Type() protoreflect.EnumType {
	return &file_transaction_proto_enumTypes[1]
}

func (x TransactionVerbosity) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use TransactionVerbosity.Descriptor instead.
func (TransactionVerbosity) EnumDescriptor() ([]byte, []int) {
	return file_transaction_proto_rawDescGZIP(), []int{1}
}

type GetTransactionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        []byte               `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Verbosity TransactionVerbosity `protobuf:"varint,2,opt,name=verbosity,proto3,enum=pactus.TransactionVerbosity" json:"verbosity,omitempty"`
}

func (x *GetTransactionRequest) Reset() {
	*x = GetTransactionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_transaction_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetTransactionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetTransactionRequest) ProtoMessage() {}

func (x *GetTransactionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_transaction_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetTransactionRequest.ProtoReflect.Descriptor instead.
func (*GetTransactionRequest) Descriptor() ([]byte, []int) {
	return file_transaction_proto_rawDescGZIP(), []int{0}
}

func (x *GetTransactionRequest) GetId() []byte {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *GetTransactionRequest) GetVerbosity() TransactionVerbosity {
	if x != nil {
		return x.Verbosity
	}
	return TransactionVerbosity_TRANSACTION_DATA
}

type GetTransactionResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BlockHeight uint32           `protobuf:"varint,12,opt,name=block_height,json=blockHeight,proto3" json:"block_height,omitempty"`
	BlockTime   uint32           `protobuf:"varint,13,opt,name=block_time,json=blockTime,proto3" json:"block_time,omitempty"`
	Transaction *TransactionInfo `protobuf:"bytes,3,opt,name=transaction,proto3" json:"transaction,omitempty"`
}

func (x *GetTransactionResponse) Reset() {
	*x = GetTransactionResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_transaction_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetTransactionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetTransactionResponse) ProtoMessage() {}

func (x *GetTransactionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_transaction_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetTransactionResponse.ProtoReflect.Descriptor instead.
func (*GetTransactionResponse) Descriptor() ([]byte, []int) {
	return file_transaction_proto_rawDescGZIP(), []int{1}
}

func (x *GetTransactionResponse) GetBlockHeight() uint32 {
	if x != nil {
		return x.BlockHeight
	}
	return 0
}

func (x *GetTransactionResponse) GetBlockTime() uint32 {
	if x != nil {
		return x.BlockTime
	}
	return 0
}

func (x *GetTransactionResponse) GetTransaction() *TransactionInfo {
	if x != nil {
		return x.Transaction
	}
	return nil
}

type SendRawTransactionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data []byte `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *SendRawTransactionRequest) Reset() {
	*x = SendRawTransactionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_transaction_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendRawTransactionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendRawTransactionRequest) ProtoMessage() {}

func (x *SendRawTransactionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_transaction_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendRawTransactionRequest.ProtoReflect.Descriptor instead.
func (*SendRawTransactionRequest) Descriptor() ([]byte, []int) {
	return file_transaction_proto_rawDescGZIP(), []int{2}
}

func (x *SendRawTransactionRequest) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type SendRawTransactionResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id []byte `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *SendRawTransactionResponse) Reset() {
	*x = SendRawTransactionResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_transaction_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendRawTransactionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendRawTransactionResponse) ProtoMessage() {}

func (x *SendRawTransactionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_transaction_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendRawTransactionResponse.ProtoReflect.Descriptor instead.
func (*SendRawTransactionResponse) Descriptor() ([]byte, []int) {
	return file_transaction_proto_rawDescGZIP(), []int{3}
}

func (x *SendRawTransactionResponse) GetId() []byte {
	if x != nil {
		return x.Id
	}
	return nil
}

type PayloadSend struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sender   string `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender,omitempty"`
	Receiver string `protobuf:"bytes,2,opt,name=receiver,proto3" json:"receiver,omitempty"`
	Amount   int64  `protobuf:"varint,3,opt,name=amount,proto3" json:"amount,omitempty"`
}

func (x *PayloadSend) Reset() {
	*x = PayloadSend{}
	if protoimpl.UnsafeEnabled {
		mi := &file_transaction_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PayloadSend) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PayloadSend) ProtoMessage() {}

func (x *PayloadSend) ProtoReflect() protoreflect.Message {
	mi := &file_transaction_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PayloadSend.ProtoReflect.Descriptor instead.
func (*PayloadSend) Descriptor() ([]byte, []int) {
	return file_transaction_proto_rawDescGZIP(), []int{4}
}

func (x *PayloadSend) GetSender() string {
	if x != nil {
		return x.Sender
	}
	return ""
}

func (x *PayloadSend) GetReceiver() string {
	if x != nil {
		return x.Receiver
	}
	return ""
}

func (x *PayloadSend) GetAmount() int64 {
	if x != nil {
		return x.Amount
	}
	return 0
}

type PayloadBond struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sender   string `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender,omitempty"`
	Receiver string `protobuf:"bytes,2,opt,name=receiver,proto3" json:"receiver,omitempty"`
	Stake    int64  `protobuf:"varint,3,opt,name=stake,proto3" json:"stake,omitempty"`
}

func (x *PayloadBond) Reset() {
	*x = PayloadBond{}
	if protoimpl.UnsafeEnabled {
		mi := &file_transaction_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PayloadBond) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PayloadBond) ProtoMessage() {}

func (x *PayloadBond) ProtoReflect() protoreflect.Message {
	mi := &file_transaction_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PayloadBond.ProtoReflect.Descriptor instead.
func (*PayloadBond) Descriptor() ([]byte, []int) {
	return file_transaction_proto_rawDescGZIP(), []int{5}
}

func (x *PayloadBond) GetSender() string {
	if x != nil {
		return x.Sender
	}
	return ""
}

func (x *PayloadBond) GetReceiver() string {
	if x != nil {
		return x.Receiver
	}
	return ""
}

func (x *PayloadBond) GetStake() int64 {
	if x != nil {
		return x.Stake
	}
	return 0
}

type PayloadSortition struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	Proof   []byte `protobuf:"bytes,2,opt,name=proof,proto3" json:"proof,omitempty"`
}

func (x *PayloadSortition) Reset() {
	*x = PayloadSortition{}
	if protoimpl.UnsafeEnabled {
		mi := &file_transaction_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PayloadSortition) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PayloadSortition) ProtoMessage() {}

func (x *PayloadSortition) ProtoReflect() protoreflect.Message {
	mi := &file_transaction_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PayloadSortition.ProtoReflect.Descriptor instead.
func (*PayloadSortition) Descriptor() ([]byte, []int) {
	return file_transaction_proto_rawDescGZIP(), []int{6}
}

func (x *PayloadSortition) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *PayloadSortition) GetProof() []byte {
	if x != nil {
		return x.Proof
	}
	return nil
}

type TransactionInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       []byte      `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Data     []byte      `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	Version  int32       `protobuf:"varint,3,opt,name=version,proto3" json:"version,omitempty"`
	Stamp    []byte      `protobuf:"bytes,4,opt,name=stamp,proto3" json:"stamp,omitempty"`
	Sequence int32       `protobuf:"varint,5,opt,name=sequence,proto3" json:"sequence,omitempty"`
	Value    int64       `protobuf:"varint,6,opt,name=value,proto3" json:"value,omitempty"`
	Fee      int64       `protobuf:"varint,7,opt,name=fee,proto3" json:"fee,omitempty"`
	Type     PayloadType `protobuf:"varint,8,opt,name=Type,proto3,enum=pactus.PayloadType" json:"Type,omitempty"`
	// Types that are assignable to Payload:
	//
	//	*TransactionInfo_Send
	//	*TransactionInfo_Bond
	//	*TransactionInfo_Sortition
	Payload   isTransactionInfo_Payload `protobuf_oneof:"Payload"`
	Memo      string                    `protobuf:"bytes,9,opt,name=memo,proto3" json:"memo,omitempty"`
	PublicKey string                    `protobuf:"bytes,10,opt,name=public_key,json=publicKey,proto3" json:"public_key,omitempty"`
	Signature []byte                    `protobuf:"bytes,11,opt,name=signature,proto3" json:"signature,omitempty"`
}

func (x *TransactionInfo) Reset() {
	*x = TransactionInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_transaction_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TransactionInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TransactionInfo) ProtoMessage() {}

func (x *TransactionInfo) ProtoReflect() protoreflect.Message {
	mi := &file_transaction_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TransactionInfo.ProtoReflect.Descriptor instead.
func (*TransactionInfo) Descriptor() ([]byte, []int) {
	return file_transaction_proto_rawDescGZIP(), []int{7}
}

func (x *TransactionInfo) GetId() []byte {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *TransactionInfo) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *TransactionInfo) GetVersion() int32 {
	if x != nil {
		return x.Version
	}
	return 0
}

func (x *TransactionInfo) GetStamp() []byte {
	if x != nil {
		return x.Stamp
	}
	return nil
}

func (x *TransactionInfo) GetSequence() int32 {
	if x != nil {
		return x.Sequence
	}
	return 0
}

func (x *TransactionInfo) GetValue() int64 {
	if x != nil {
		return x.Value
	}
	return 0
}

func (x *TransactionInfo) GetFee() int64 {
	if x != nil {
		return x.Fee
	}
	return 0
}

func (x *TransactionInfo) GetType() PayloadType {
	if x != nil {
		return x.Type
	}
	return PayloadType_UNKNOWN
}

func (m *TransactionInfo) GetPayload() isTransactionInfo_Payload {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (x *TransactionInfo) GetSend() *PayloadSend {
	if x, ok := x.GetPayload().(*TransactionInfo_Send); ok {
		return x.Send
	}
	return nil
}

func (x *TransactionInfo) GetBond() *PayloadBond {
	if x, ok := x.GetPayload().(*TransactionInfo_Bond); ok {
		return x.Bond
	}
	return nil
}

func (x *TransactionInfo) GetSortition() *PayloadSortition {
	if x, ok := x.GetPayload().(*TransactionInfo_Sortition); ok {
		return x.Sortition
	}
	return nil
}

func (x *TransactionInfo) GetMemo() string {
	if x != nil {
		return x.Memo
	}
	return ""
}

func (x *TransactionInfo) GetPublicKey() string {
	if x != nil {
		return x.PublicKey
	}
	return ""
}

func (x *TransactionInfo) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

type isTransactionInfo_Payload interface {
	isTransactionInfo_Payload()
}

type TransactionInfo_Send struct {
	Send *PayloadSend `protobuf:"bytes,30,opt,name=send,proto3,oneof"`
}

type TransactionInfo_Bond struct {
	Bond *PayloadBond `protobuf:"bytes,31,opt,name=bond,proto3,oneof"`
}

type TransactionInfo_Sortition struct {
	Sortition *PayloadSortition `protobuf:"bytes,32,opt,name=sortition,proto3,oneof"`
}

func (*TransactionInfo_Send) isTransactionInfo_Payload() {}

func (*TransactionInfo_Bond) isTransactionInfo_Payload() {}

func (*TransactionInfo_Sortition) isTransactionInfo_Payload() {}

var File_transaction_proto protoreflect.FileDescriptor

var file_transaction_proto_rawDesc = []byte{
	0x0a, 0x11, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x06, 0x70, 0x61, 0x63, 0x74, 0x75, 0x73, 0x22, 0x63, 0x0a, 0x15, 0x47,
	0x65, 0x74, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x3a, 0x0a, 0x09, 0x76, 0x65, 0x72, 0x62, 0x6f, 0x73, 0x69, 0x74,
	0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1c, 0x2e, 0x70, 0x61, 0x63, 0x74, 0x75, 0x73,
	0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x56, 0x65, 0x72, 0x62,
	0x6f, 0x73, 0x69, 0x74, 0x79, 0x52, 0x09, 0x76, 0x65, 0x72, 0x62, 0x6f, 0x73, 0x69, 0x74, 0x79,
	0x22, 0x95, 0x01, 0x0a, 0x16, 0x47, 0x65, 0x74, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x62,
	0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x68, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x0c, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x0b, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x12, 0x1d,
	0x0a, 0x0a, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x0d, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x09, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x39, 0x0a,
	0x0b, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x17, 0x2e, 0x70, 0x61, 0x63, 0x74, 0x75, 0x73, 0x2e, 0x54, 0x72, 0x61, 0x6e,
	0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x0b, 0x74, 0x72, 0x61,
	0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x2f, 0x0a, 0x19, 0x53, 0x65, 0x6e, 0x64,
	0x52, 0x61, 0x77, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x2c, 0x0a, 0x1a, 0x53, 0x65, 0x6e,
	0x64, 0x52, 0x61, 0x77, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x02, 0x69, 0x64, 0x22, 0x59, 0x0a, 0x0b, 0x50, 0x61, 0x79, 0x6c, 0x6f,
	0x61, 0x64, 0x53, 0x65, 0x6e, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x12, 0x1a,
	0x0a, 0x08, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x6d,
	0x6f, 0x75, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x61, 0x6d, 0x6f, 0x75,
	0x6e, 0x74, 0x22, 0x57, 0x0a, 0x0b, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x42, 0x6f, 0x6e,
	0x64, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65, 0x63,
	0x65, 0x69, 0x76, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x72, 0x65, 0x63,
	0x65, 0x69, 0x76, 0x65, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x6b, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x73, 0x74, 0x61, 0x6b, 0x65, 0x22, 0x42, 0x0a, 0x10, 0x50,
	0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x53, 0x6f, 0x72, 0x74, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x72, 0x6f,
	0x6f, 0x66, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x70, 0x72, 0x6f, 0x6f, 0x66, 0x22,
	0xbe, 0x03, 0x0a, 0x0f, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49,
	0x6e, 0x66, 0x6f, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x05, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x65, 0x71, 0x75, 0x65,
	0x6e, 0x63, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x73, 0x65, 0x71, 0x75, 0x65,
	0x6e, 0x63, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x66, 0x65, 0x65,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x66, 0x65, 0x65, 0x12, 0x27, 0x0a, 0x04, 0x54,
	0x79, 0x70, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x13, 0x2e, 0x70, 0x61, 0x63, 0x74,
	0x75, 0x73, 0x2e, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04,
	0x54, 0x79, 0x70, 0x65, 0x12, 0x29, 0x0a, 0x04, 0x73, 0x65, 0x6e, 0x64, 0x18, 0x1e, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x13, 0x2e, 0x70, 0x61, 0x63, 0x74, 0x75, 0x73, 0x2e, 0x50, 0x61, 0x79, 0x6c,
	0x6f, 0x61, 0x64, 0x53, 0x65, 0x6e, 0x64, 0x48, 0x00, 0x52, 0x04, 0x73, 0x65, 0x6e, 0x64, 0x12,
	0x29, 0x0a, 0x04, 0x62, 0x6f, 0x6e, 0x64, 0x18, 0x1f, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e,
	0x70, 0x61, 0x63, 0x74, 0x75, 0x73, 0x2e, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x42, 0x6f,
	0x6e, 0x64, 0x48, 0x00, 0x52, 0x04, 0x62, 0x6f, 0x6e, 0x64, 0x12, 0x38, 0x0a, 0x09, 0x73, 0x6f,
	0x72, 0x74, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x20, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e,
	0x70, 0x61, 0x63, 0x74, 0x75, 0x73, 0x2e, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x53, 0x6f,
	0x72, 0x74, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x48, 0x00, 0x52, 0x09, 0x73, 0x6f, 0x72, 0x74, 0x69,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x6d, 0x65, 0x6d, 0x6f, 0x18, 0x09, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6d, 0x65, 0x6d, 0x6f, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x75, 0x62, 0x6c,
	0x69, 0x63, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x70, 0x75,
	0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61,
	0x74, 0x75, 0x72, 0x65, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x73, 0x69, 0x67, 0x6e,
	0x61, 0x74, 0x75, 0x72, 0x65, 0x42, 0x09, 0x0a, 0x07, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64,
	0x2a, 0x7f, 0x0a, 0x0b, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x10, 0x0a, 0x0c,
	0x53, 0x45, 0x4e, 0x44, 0x5f, 0x50, 0x41, 0x59, 0x4c, 0x4f, 0x41, 0x44, 0x10, 0x01, 0x12, 0x10,
	0x0a, 0x0c, 0x42, 0x4f, 0x4e, 0x44, 0x5f, 0x50, 0x41, 0x59, 0x4c, 0x4f, 0x41, 0x44, 0x10, 0x02,
	0x12, 0x15, 0x0a, 0x11, 0x53, 0x4f, 0x52, 0x54, 0x49, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x50, 0x41,
	0x59, 0x4c, 0x4f, 0x41, 0x44, 0x10, 0x03, 0x12, 0x12, 0x0a, 0x0e, 0x55, 0x4e, 0x42, 0x4f, 0x4e,
	0x44, 0x5f, 0x50, 0x41, 0x59, 0x4c, 0x4f, 0x41, 0x44, 0x10, 0x04, 0x12, 0x14, 0x0a, 0x10, 0x57,
	0x49, 0x54, 0x48, 0x44, 0x52, 0x41, 0x57, 0x5f, 0x50, 0x41, 0x59, 0x4c, 0x4f, 0x41, 0x44, 0x10,
	0x05, 0x2a, 0x42, 0x0a, 0x14, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x56, 0x65, 0x72, 0x62, 0x6f, 0x73, 0x69, 0x74, 0x79, 0x12, 0x14, 0x0a, 0x10, 0x54, 0x52, 0x41,
	0x4e, 0x53, 0x41, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x44, 0x41, 0x54, 0x41, 0x10, 0x00, 0x12,
	0x14, 0x0a, 0x10, 0x54, 0x52, 0x41, 0x4e, 0x53, 0x41, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x49,
	0x4e, 0x46, 0x4f, 0x10, 0x01, 0x32, 0xbb, 0x01, 0x0a, 0x0b, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x4f, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x54, 0x72, 0x61, 0x6e,
	0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1d, 0x2e, 0x70, 0x61, 0x63, 0x74, 0x75, 0x73,
	0x2e, 0x47, 0x65, 0x74, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x70, 0x61, 0x63, 0x74, 0x75, 0x73, 0x2e,
	0x47, 0x65, 0x74, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x5b, 0x0a, 0x12, 0x53, 0x65, 0x6e, 0x64, 0x52, 0x61,
	0x77, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x21, 0x2e, 0x70,
	0x61, 0x63, 0x74, 0x75, 0x73, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x52, 0x61, 0x77, 0x54, 0x72, 0x61,
	0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x22, 0x2e, 0x70, 0x61, 0x63, 0x74, 0x75, 0x73, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x52, 0x61, 0x77,
	0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x42, 0x46, 0x0a, 0x12, 0x70, 0x61, 0x63, 0x74, 0x75, 0x73, 0x2e, 0x74, 0x72,
	0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5a, 0x30, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x70, 0x61, 0x63, 0x74, 0x75, 0x73, 0x2d, 0x70, 0x72, 0x6f,
	0x6a, 0x65, 0x63, 0x74, 0x2f, 0x70, 0x61, 0x63, 0x74, 0x75, 0x73, 0x2f, 0x77, 0x77, 0x77, 0x2f,
	0x67, 0x72, 0x70, 0x63, 0x2f, 0x70, 0x61, 0x63, 0x74, 0x75, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_transaction_proto_rawDescOnce sync.Once
	file_transaction_proto_rawDescData = file_transaction_proto_rawDesc
)

func file_transaction_proto_rawDescGZIP() []byte {
	file_transaction_proto_rawDescOnce.Do(func() {
		file_transaction_proto_rawDescData = protoimpl.X.CompressGZIP(file_transaction_proto_rawDescData)
	})
	return file_transaction_proto_rawDescData
}

var file_transaction_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_transaction_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_transaction_proto_goTypes = []interface{}{
	(PayloadType)(0),                   // 0: pactus.PayloadType
	(TransactionVerbosity)(0),          // 1: pactus.TransactionVerbosity
	(*GetTransactionRequest)(nil),      // 2: pactus.GetTransactionRequest
	(*GetTransactionResponse)(nil),     // 3: pactus.GetTransactionResponse
	(*SendRawTransactionRequest)(nil),  // 4: pactus.SendRawTransactionRequest
	(*SendRawTransactionResponse)(nil), // 5: pactus.SendRawTransactionResponse
	(*PayloadSend)(nil),                // 6: pactus.PayloadSend
	(*PayloadBond)(nil),                // 7: pactus.PayloadBond
	(*PayloadSortition)(nil),           // 8: pactus.PayloadSortition
	(*TransactionInfo)(nil),            // 9: pactus.TransactionInfo
}
var file_transaction_proto_depIdxs = []int32{
	1, // 0: pactus.GetTransactionRequest.verbosity:type_name -> pactus.TransactionVerbosity
	9, // 1: pactus.GetTransactionResponse.transaction:type_name -> pactus.TransactionInfo
	0, // 2: pactus.TransactionInfo.Type:type_name -> pactus.PayloadType
	6, // 3: pactus.TransactionInfo.send:type_name -> pactus.PayloadSend
	7, // 4: pactus.TransactionInfo.bond:type_name -> pactus.PayloadBond
	8, // 5: pactus.TransactionInfo.sortition:type_name -> pactus.PayloadSortition
	2, // 6: pactus.Transaction.GetTransaction:input_type -> pactus.GetTransactionRequest
	4, // 7: pactus.Transaction.SendRawTransaction:input_type -> pactus.SendRawTransactionRequest
	3, // 8: pactus.Transaction.GetTransaction:output_type -> pactus.GetTransactionResponse
	5, // 9: pactus.Transaction.SendRawTransaction:output_type -> pactus.SendRawTransactionResponse
	8, // [8:10] is the sub-list for method output_type
	6, // [6:8] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_transaction_proto_init() }
func file_transaction_proto_init() {
	if File_transaction_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_transaction_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetTransactionRequest); i {
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
		file_transaction_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetTransactionResponse); i {
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
		file_transaction_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendRawTransactionRequest); i {
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
		file_transaction_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendRawTransactionResponse); i {
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
		file_transaction_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PayloadSend); i {
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
		file_transaction_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PayloadBond); i {
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
		file_transaction_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PayloadSortition); i {
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
		file_transaction_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TransactionInfo); i {
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
	file_transaction_proto_msgTypes[7].OneofWrappers = []interface{}{
		(*TransactionInfo_Send)(nil),
		(*TransactionInfo_Bond)(nil),
		(*TransactionInfo_Sortition)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_transaction_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_transaction_proto_goTypes,
		DependencyIndexes: file_transaction_proto_depIdxs,
		EnumInfos:         file_transaction_proto_enumTypes,
		MessageInfos:      file_transaction_proto_msgTypes,
	}.Build()
	File_transaction_proto = out.File
	file_transaction_proto_rawDesc = nil
	file_transaction_proto_goTypes = nil
	file_transaction_proto_depIdxs = nil
}
