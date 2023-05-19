///
//  Generated code. Do not modify.
//  source: transaction.proto
//
// @dart = 2.12
// ignore_for_file: annotate_overrides,camel_case_types,constant_identifier_names,deprecated_member_use_from_same_package,directives_ordering,library_prefixes,non_constant_identifier_names,prefer_final_fields,return_of_invalid_type,unnecessary_const,unnecessary_import,unnecessary_this,unused_import,unused_shown_name

import 'dart:core' as $core;
import 'dart:convert' as $convert;
import 'dart:typed_data' as $typed_data;
@$core.Deprecated('Use payloadTypeDescriptor instead')
const PayloadType$json = const {
  '1': 'PayloadType',
  '2': const [
    const {'1': 'UNKNOWN', '2': 0},
    const {'1': 'TRANSFER_PAYLOAD', '2': 1},
    const {'1': 'BOND_PAYLOAD', '2': 2},
    const {'1': 'SORTITION_PAYLOAD', '2': 3},
    const {'1': 'UNBOND_PAYLOAD', '2': 4},
    const {'1': 'WITHDRAW_PAYLOAD', '2': 5},
  ],
};

/// Descriptor for `PayloadType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List payloadTypeDescriptor = $convert.base64Decode('CgtQYXlsb2FkVHlwZRILCgdVTktOT1dOEAASFAoQVFJBTlNGRVJfUEFZTE9BRBABEhAKDEJPTkRfUEFZTE9BRBACEhUKEVNPUlRJVElPTl9QQVlMT0FEEAMSEgoOVU5CT05EX1BBWUxPQUQQBBIUChBXSVRIRFJBV19QQVlMT0FEEAU=');
@$core.Deprecated('Use transactionVerbosityDescriptor instead')
const TransactionVerbosity$json = const {
  '1': 'TransactionVerbosity',
  '2': const [
    const {'1': 'TRANSACTION_DATA', '2': 0},
    const {'1': 'TRANSACTION_INFO', '2': 1},
  ],
};

/// Descriptor for `TransactionVerbosity`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List transactionVerbosityDescriptor = $convert.base64Decode('ChRUcmFuc2FjdGlvblZlcmJvc2l0eRIUChBUUkFOU0FDVElPTl9EQVRBEAASFAoQVFJBTlNBQ1RJT05fSU5GTxAB');
@$core.Deprecated('Use getTransactionRequestDescriptor instead')
const GetTransactionRequest$json = const {
  '1': 'GetTransactionRequest',
  '2': const [
    const {'1': 'id', '3': 1, '4': 1, '5': 12, '10': 'id'},
    const {'1': 'verbosity', '3': 2, '4': 1, '5': 14, '6': '.pactus.TransactionVerbosity', '10': 'verbosity'},
  ],
};

/// Descriptor for `GetTransactionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getTransactionRequestDescriptor = $convert.base64Decode('ChVHZXRUcmFuc2FjdGlvblJlcXVlc3QSDgoCaWQYASABKAxSAmlkEjoKCXZlcmJvc2l0eRgCIAEoDjIcLnBhY3R1cy5UcmFuc2FjdGlvblZlcmJvc2l0eVIJdmVyYm9zaXR5');
@$core.Deprecated('Use getTransactionResponseDescriptor instead')
const GetTransactionResponse$json = const {
  '1': 'GetTransactionResponse',
  '2': const [
    const {'1': 'block_height', '3': 12, '4': 1, '5': 13, '10': 'blockHeight'},
    const {'1': 'block_time', '3': 13, '4': 1, '5': 13, '10': 'blockTime'},
    const {'1': 'transaction', '3': 3, '4': 1, '5': 11, '6': '.pactus.TransactionInfo', '10': 'transaction'},
  ],
};

/// Descriptor for `GetTransactionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getTransactionResponseDescriptor = $convert.base64Decode('ChZHZXRUcmFuc2FjdGlvblJlc3BvbnNlEiEKDGJsb2NrX2hlaWdodBgMIAEoDVILYmxvY2tIZWlnaHQSHQoKYmxvY2tfdGltZRgNIAEoDVIJYmxvY2tUaW1lEjkKC3RyYW5zYWN0aW9uGAMgASgLMhcucGFjdHVzLlRyYW5zYWN0aW9uSW5mb1ILdHJhbnNhY3Rpb24=');
@$core.Deprecated('Use sendRawTransactionRequestDescriptor instead')
const SendRawTransactionRequest$json = const {
  '1': 'SendRawTransactionRequest',
  '2': const [
    const {'1': 'data', '3': 1, '4': 1, '5': 12, '10': 'data'},
  ],
};

/// Descriptor for `SendRawTransactionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sendRawTransactionRequestDescriptor = $convert.base64Decode('ChlTZW5kUmF3VHJhbnNhY3Rpb25SZXF1ZXN0EhIKBGRhdGEYASABKAxSBGRhdGE=');
@$core.Deprecated('Use sendRawTransactionResponseDescriptor instead')
const SendRawTransactionResponse$json = const {
  '1': 'SendRawTransactionResponse',
  '2': const [
    const {'1': 'id', '3': 2, '4': 1, '5': 12, '10': 'id'},
  ],
};

/// Descriptor for `SendRawTransactionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sendRawTransactionResponseDescriptor = $convert.base64Decode('ChpTZW5kUmF3VHJhbnNhY3Rpb25SZXNwb25zZRIOCgJpZBgCIAEoDFICaWQ=');
@$core.Deprecated('Use payloadTransferDescriptor instead')
const PayloadTransfer$json = const {
  '1': 'PayloadTransfer',
  '2': const [
    const {'1': 'sender', '3': 1, '4': 1, '5': 9, '10': 'sender'},
    const {'1': 'receiver', '3': 2, '4': 1, '5': 9, '10': 'receiver'},
    const {'1': 'amount', '3': 3, '4': 1, '5': 3, '10': 'amount'},
  ],
};

/// Descriptor for `PayloadTransfer`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List payloadTransferDescriptor = $convert.base64Decode('Cg9QYXlsb2FkVHJhbnNmZXISFgoGc2VuZGVyGAEgASgJUgZzZW5kZXISGgoIcmVjZWl2ZXIYAiABKAlSCHJlY2VpdmVyEhYKBmFtb3VudBgDIAEoA1IGYW1vdW50');
@$core.Deprecated('Use payloadBondDescriptor instead')
const PayloadBond$json = const {
  '1': 'PayloadBond',
  '2': const [
    const {'1': 'sender', '3': 1, '4': 1, '5': 9, '10': 'sender'},
    const {'1': 'receiver', '3': 2, '4': 1, '5': 9, '10': 'receiver'},
    const {'1': 'stake', '3': 3, '4': 1, '5': 3, '10': 'stake'},
  ],
};

/// Descriptor for `PayloadBond`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List payloadBondDescriptor = $convert.base64Decode('CgtQYXlsb2FkQm9uZBIWCgZzZW5kZXIYASABKAlSBnNlbmRlchIaCghyZWNlaXZlchgCIAEoCVIIcmVjZWl2ZXISFAoFc3Rha2UYAyABKANSBXN0YWtl');
@$core.Deprecated('Use payloadSortitionDescriptor instead')
const PayloadSortition$json = const {
  '1': 'PayloadSortition',
  '2': const [
    const {'1': 'address', '3': 1, '4': 1, '5': 9, '10': 'address'},
    const {'1': 'proof', '3': 2, '4': 1, '5': 12, '10': 'proof'},
  ],
};

/// Descriptor for `PayloadSortition`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List payloadSortitionDescriptor = $convert.base64Decode('ChBQYXlsb2FkU29ydGl0aW9uEhgKB2FkZHJlc3MYASABKAlSB2FkZHJlc3MSFAoFcHJvb2YYAiABKAxSBXByb29m');
@$core.Deprecated('Use payloadUnbondDescriptor instead')
const PayloadUnbond$json = const {
  '1': 'PayloadUnbond',
  '2': const [
    const {'1': 'validator', '3': 1, '4': 1, '5': 9, '10': 'validator'},
  ],
};

/// Descriptor for `PayloadUnbond`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List payloadUnbondDescriptor = $convert.base64Decode('Cg1QYXlsb2FkVW5ib25kEhwKCXZhbGlkYXRvchgBIAEoCVIJdmFsaWRhdG9y');
@$core.Deprecated('Use payloadWithdrawDescriptor instead')
const PayloadWithdraw$json = const {
  '1': 'PayloadWithdraw',
  '2': const [
    const {'1': 'from', '3': 1, '4': 1, '5': 9, '10': 'from'},
    const {'1': 'to', '3': 2, '4': 1, '5': 9, '10': 'to'},
    const {'1': 'amount', '3': 3, '4': 1, '5': 3, '10': 'amount'},
  ],
};

/// Descriptor for `PayloadWithdraw`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List payloadWithdrawDescriptor = $convert.base64Decode('Cg9QYXlsb2FkV2l0aGRyYXcSEgoEZnJvbRgBIAEoCVIEZnJvbRIOCgJ0bxgCIAEoCVICdG8SFgoGYW1vdW50GAMgASgDUgZhbW91bnQ=');
@$core.Deprecated('Use transactionInfoDescriptor instead')
const TransactionInfo$json = const {
  '1': 'TransactionInfo',
  '2': const [
    const {'1': 'id', '3': 1, '4': 1, '5': 12, '10': 'id'},
    const {'1': 'data', '3': 2, '4': 1, '5': 12, '10': 'data'},
    const {'1': 'version', '3': 3, '4': 1, '5': 5, '10': 'version'},
    const {'1': 'stamp', '3': 4, '4': 1, '5': 12, '10': 'stamp'},
    const {'1': 'sequence', '3': 5, '4': 1, '5': 5, '10': 'sequence'},
    const {'1': 'value', '3': 6, '4': 1, '5': 3, '10': 'value'},
    const {'1': 'fee', '3': 7, '4': 1, '5': 3, '10': 'fee'},
    const {'1': 'PayloadType', '3': 8, '4': 1, '5': 14, '6': '.pactus.PayloadType', '10': 'PayloadType'},
    const {'1': 'transfer', '3': 30, '4': 1, '5': 11, '6': '.pactus.PayloadTransfer', '9': 0, '10': 'transfer'},
    const {'1': 'bond', '3': 31, '4': 1, '5': 11, '6': '.pactus.PayloadBond', '9': 0, '10': 'bond'},
    const {'1': 'sortition', '3': 32, '4': 1, '5': 11, '6': '.pactus.PayloadSortition', '9': 0, '10': 'sortition'},
    const {'1': 'unbond', '3': 33, '4': 1, '5': 11, '6': '.pactus.PayloadUnbond', '9': 0, '10': 'unbond'},
    const {'1': 'withdraw', '3': 34, '4': 1, '5': 11, '6': '.pactus.PayloadWithdraw', '9': 0, '10': 'withdraw'},
    const {'1': 'memo', '3': 9, '4': 1, '5': 9, '10': 'memo'},
    const {'1': 'public_key', '3': 10, '4': 1, '5': 9, '10': 'publicKey'},
    const {'1': 'signature', '3': 11, '4': 1, '5': 12, '10': 'signature'},
  ],
  '8': const [
    const {'1': 'Payload'},
  ],
};

/// Descriptor for `TransactionInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List transactionInfoDescriptor = $convert.base64Decode('Cg9UcmFuc2FjdGlvbkluZm8SDgoCaWQYASABKAxSAmlkEhIKBGRhdGEYAiABKAxSBGRhdGESGAoHdmVyc2lvbhgDIAEoBVIHdmVyc2lvbhIUCgVzdGFtcBgEIAEoDFIFc3RhbXASGgoIc2VxdWVuY2UYBSABKAVSCHNlcXVlbmNlEhQKBXZhbHVlGAYgASgDUgV2YWx1ZRIQCgNmZWUYByABKANSA2ZlZRI1CgtQYXlsb2FkVHlwZRgIIAEoDjITLnBhY3R1cy5QYXlsb2FkVHlwZVILUGF5bG9hZFR5cGUSNQoIdHJhbnNmZXIYHiABKAsyFy5wYWN0dXMuUGF5bG9hZFRyYW5zZmVySABSCHRyYW5zZmVyEikKBGJvbmQYHyABKAsyEy5wYWN0dXMuUGF5bG9hZEJvbmRIAFIEYm9uZBI4Cglzb3J0aXRpb24YICABKAsyGC5wYWN0dXMuUGF5bG9hZFNvcnRpdGlvbkgAUglzb3J0aXRpb24SLwoGdW5ib25kGCEgASgLMhUucGFjdHVzLlBheWxvYWRVbmJvbmRIAFIGdW5ib25kEjUKCHdpdGhkcmF3GCIgASgLMhcucGFjdHVzLlBheWxvYWRXaXRoZHJhd0gAUgh3aXRoZHJhdxISCgRtZW1vGAkgASgJUgRtZW1vEh0KCnB1YmxpY19rZXkYCiABKAlSCXB1YmxpY0tleRIcCglzaWduYXR1cmUYCyABKAxSCXNpZ25hdHVyZUIJCgdQYXlsb2Fk');
const $core.Map<$core.String, $core.dynamic> TransactionServiceBase$json = const {
  '1': 'Transaction',
  '2': const [
    const {'1': 'GetTransaction', '2': '.pactus.GetTransactionRequest', '3': '.pactus.GetTransactionResponse'},
    const {'1': 'SendRawTransaction', '2': '.pactus.SendRawTransactionRequest', '3': '.pactus.SendRawTransactionResponse'},
  ],
};

@$core.Deprecated('Use transactionServiceDescriptor instead')
const $core.Map<$core.String, $core.Map<$core.String, $core.dynamic>> TransactionServiceBase$messageJson = const {
  '.pactus.GetTransactionRequest': GetTransactionRequest$json,
  '.pactus.GetTransactionResponse': GetTransactionResponse$json,
  '.pactus.TransactionInfo': TransactionInfo$json,
  '.pactus.PayloadTransfer': PayloadTransfer$json,
  '.pactus.PayloadBond': PayloadBond$json,
  '.pactus.PayloadSortition': PayloadSortition$json,
  '.pactus.PayloadUnbond': PayloadUnbond$json,
  '.pactus.PayloadWithdraw': PayloadWithdraw$json,
  '.pactus.SendRawTransactionRequest': SendRawTransactionRequest$json,
  '.pactus.SendRawTransactionResponse': SendRawTransactionResponse$json,
};

/// Descriptor for `Transaction`. Decode as a `google.protobuf.ServiceDescriptorProto`.
final $typed_data.Uint8List transactionServiceDescriptor = $convert.base64Decode('CgtUcmFuc2FjdGlvbhJPCg5HZXRUcmFuc2FjdGlvbhIdLnBhY3R1cy5HZXRUcmFuc2FjdGlvblJlcXVlc3QaHi5wYWN0dXMuR2V0VHJhbnNhY3Rpb25SZXNwb25zZRJbChJTZW5kUmF3VHJhbnNhY3Rpb24SIS5wYWN0dXMuU2VuZFJhd1RyYW5zYWN0aW9uUmVxdWVzdBoiLnBhY3R1cy5TZW5kUmF3VHJhbnNhY3Rpb25SZXNwb25zZQ==');
