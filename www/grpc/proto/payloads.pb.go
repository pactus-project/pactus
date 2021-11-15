// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.18.1
// source: payloads.proto

package zarb

import (
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type SEND_PAYLOAD struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sender   string `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender,omitempty"`
	Receiver string `protobuf:"bytes,2,opt,name=receiver,proto3" json:"receiver,omitempty"`
	Amount   int64  `protobuf:"varint,3,opt,name=amount,proto3" json:"amount,omitempty"`
}

func (x *SEND_PAYLOAD) Reset() {
	*x = SEND_PAYLOAD{}
	if protoimpl.UnsafeEnabled {
		mi := &file_payloads_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SEND_PAYLOAD) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SEND_PAYLOAD) ProtoMessage() {}

func (x *SEND_PAYLOAD) ProtoReflect() protoreflect.Message {
	mi := &file_payloads_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SEND_PAYLOAD.ProtoReflect.Descriptor instead.
func (*SEND_PAYLOAD) Descriptor() ([]byte, []int) {
	return file_payloads_proto_rawDescGZIP(), []int{0}
}

func (x *SEND_PAYLOAD) GetSender() string {
	if x != nil {
		return x.Sender
	}
	return ""
}

func (x *SEND_PAYLOAD) GetReceiver() string {
	if x != nil {
		return x.Receiver
	}
	return ""
}

func (x *SEND_PAYLOAD) GetAmount() int64 {
	if x != nil {
		return x.Amount
	}
	return 0
}

type BOND_PAYLOAD struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Bonder    string `protobuf:"bytes,1,opt,name=bonder,proto3" json:"bonder,omitempty"`
	Validator string `protobuf:"bytes,2,opt,name=validator,proto3" json:"validator,omitempty"`
	Stake     int64  `protobuf:"varint,3,opt,name=stake,proto3" json:"stake,omitempty"`
}

func (x *BOND_PAYLOAD) Reset() {
	*x = BOND_PAYLOAD{}
	if protoimpl.UnsafeEnabled {
		mi := &file_payloads_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BOND_PAYLOAD) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BOND_PAYLOAD) ProtoMessage() {}

func (x *BOND_PAYLOAD) ProtoReflect() protoreflect.Message {
	mi := &file_payloads_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BOND_PAYLOAD.ProtoReflect.Descriptor instead.
func (*BOND_PAYLOAD) Descriptor() ([]byte, []int) {
	return file_payloads_proto_rawDescGZIP(), []int{1}
}

func (x *BOND_PAYLOAD) GetBonder() string {
	if x != nil {
		return x.Bonder
	}
	return ""
}

func (x *BOND_PAYLOAD) GetValidator() string {
	if x != nil {
		return x.Validator
	}
	return ""
}

func (x *BOND_PAYLOAD) GetStake() int64 {
	if x != nil {
		return x.Stake
	}
	return 0
}

type SORTITION_PAYLOAD struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	Proof   string `protobuf:"bytes,2,opt,name=proof,proto3" json:"proof,omitempty"`
}

func (x *SORTITION_PAYLOAD) Reset() {
	*x = SORTITION_PAYLOAD{}
	if protoimpl.UnsafeEnabled {
		mi := &file_payloads_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SORTITION_PAYLOAD) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SORTITION_PAYLOAD) ProtoMessage() {}

func (x *SORTITION_PAYLOAD) ProtoReflect() protoreflect.Message {
	mi := &file_payloads_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SORTITION_PAYLOAD.ProtoReflect.Descriptor instead.
func (*SORTITION_PAYLOAD) Descriptor() ([]byte, []int) {
	return file_payloads_proto_rawDescGZIP(), []int{2}
}

func (x *SORTITION_PAYLOAD) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *SORTITION_PAYLOAD) GetProof() string {
	if x != nil {
		return x.Proof
	}
	return ""
}

var File_payloads_proto protoreflect.FileDescriptor

var file_payloads_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x08, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x73, 0x22, 0x5a, 0x0a, 0x0c, 0x53, 0x45,
	0x4e, 0x44, 0x5f, 0x50, 0x41, 0x59, 0x4c, 0x4f, 0x41, 0x44, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x65,
	0x6e, 0x64, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x65, 0x6e, 0x64,
	0x65, 0x72, 0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x72, 0x12, 0x16,
	0x0a, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06,
	0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x5a, 0x0a, 0x0c, 0x42, 0x4f, 0x4e, 0x44, 0x5f, 0x50,
	0x41, 0x59, 0x4c, 0x4f, 0x41, 0x44, 0x12, 0x16, 0x0a, 0x06, 0x62, 0x6f, 0x6e, 0x64, 0x65, 0x72,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x62, 0x6f, 0x6e, 0x64, 0x65, 0x72, 0x12, 0x1c,
	0x0a, 0x09, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x12, 0x14, 0x0a, 0x05,
	0x73, 0x74, 0x61, 0x6b, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x73, 0x74, 0x61,
	0x6b, 0x65, 0x22, 0x43, 0x0a, 0x11, 0x53, 0x4f, 0x52, 0x54, 0x49, 0x54, 0x49, 0x4f, 0x4e, 0x5f,
	0x50, 0x41, 0x59, 0x4c, 0x4f, 0x41, 0x44, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65,
	0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73,
	0x73, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x72, 0x6f, 0x6f, 0x66, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x70, 0x72, 0x6f, 0x6f, 0x66, 0x42, 0x2c, 0x5a, 0x2a, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x7a, 0x61, 0x72, 0x62, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x2f,
	0x7a, 0x61, 0x72, 0x62, 0x2d, 0x67, 0x6f, 0x2f, 0x77, 0x77, 0x77, 0x2f, 0x67, 0x72, 0x70, 0x63,
	0x2f, 0x7a, 0x61, 0x72, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_payloads_proto_rawDescOnce sync.Once
	file_payloads_proto_rawDescData = file_payloads_proto_rawDesc
)

func file_payloads_proto_rawDescGZIP() []byte {
	file_payloads_proto_rawDescOnce.Do(func() {
		file_payloads_proto_rawDescData = protoimpl.X.CompressGZIP(file_payloads_proto_rawDescData)
	})
	return file_payloads_proto_rawDescData
}

var file_payloads_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_payloads_proto_goTypes = []interface{}{
	(*SEND_PAYLOAD)(nil),      // 0: payloads.SEND_PAYLOAD
	(*BOND_PAYLOAD)(nil),      // 1: payloads.BOND_PAYLOAD
	(*SORTITION_PAYLOAD)(nil), // 2: payloads.SORTITION_PAYLOAD
}
var file_payloads_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_payloads_proto_init() }
func file_payloads_proto_init() {
	if File_payloads_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_payloads_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SEND_PAYLOAD); i {
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
		file_payloads_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BOND_PAYLOAD); i {
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
		file_payloads_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SORTITION_PAYLOAD); i {
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
			RawDescriptor: file_payloads_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_payloads_proto_goTypes,
		DependencyIndexes: file_payloads_proto_depIdxs,
		MessageInfos:      file_payloads_proto_msgTypes,
	}.Build()
	File_payloads_proto = out.File
	file_payloads_proto_rawDesc = nil
	file_payloads_proto_goTypes = nil
	file_payloads_proto_depIdxs = nil
}
