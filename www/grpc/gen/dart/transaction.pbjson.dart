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
@$core.Deprecated('Use calculateFeeRequestDescriptor instead')
const CalculateFeeRequest$json = const {
  '1': 'CalculateFeeRequest',
  '2': const [
    const {'1': 'amount', '3': 1, '4': 1, '5': 3, '10': 'amount'},
    const {'1': 'payload_type', '3': 2, '4': 1, '5': 14, '6': '.pactus.PayloadType', '10': 'payloadType'},
  ],
};

/// Descriptor for `CalculateFeeRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List calculateFeeRequestDescriptor = $convert.base64Decode('ChNDYWxjdWxhdGVGZWVSZXF1ZXN0EhYKBmFtb3VudBgBIAEoA1IGYW1vdW50EjYKDHBheWxvYWRfdHlwZRgCIAEoDjITLnBhY3R1cy5QYXlsb2FkVHlwZVILcGF5bG9hZFR5cGU=');
@$core.Deprecated('Use calculateFeeResponseDescriptor instead')
const CalculateFeeResponse$json = const {
  '1': 'CalculateFeeResponse',
  '2': const [
    const {'1': 'fee', '3': 1, '4': 1, '5': 3, '10': 'fee'},
  ],
};

/// Descriptor for `CalculateFeeResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List calculateFeeResponseDescriptor = $convert.base64Decode('ChRDYWxjdWxhdGVGZWVSZXNwb25zZRIQCgNmZWUYASABKANSA2ZlZQ==');
@$core.Deprecated('Use broadcastTransactionRequestDescriptor instead')
const BroadcastTransactionRequest$json = const {
  '1': 'BroadcastTransactionRequest',
  '2': const [
    const {'1': 'signed_raw_transaction', '3': 1, '4': 1, '5': 12, '10': 'signedRawTransaction'},
  ],
};

/// Descriptor for `BroadcastTransactionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List broadcastTransactionRequestDescriptor = $convert.base64Decode('ChtCcm9hZGNhc3RUcmFuc2FjdGlvblJlcXVlc3QSNAoWc2lnbmVkX3Jhd190cmFuc2FjdGlvbhgBIAEoDFIUc2lnbmVkUmF3VHJhbnNhY3Rpb24=');
@$core.Deprecated('Use broadcastTransactionResponseDescriptor instead')
const BroadcastTransactionResponse$json = const {
  '1': 'BroadcastTransactionResponse',
  '2': const [
    const {'1': 'id', '3': 2, '4': 1, '5': 12, '10': 'id'},
  ],
};

/// Descriptor for `BroadcastTransactionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List broadcastTransactionResponseDescriptor = $convert.base64Decode('ChxCcm9hZGNhc3RUcmFuc2FjdGlvblJlc3BvbnNlEg4KAmlkGAIgASgMUgJpZA==');
@$core.Deprecated('Use getRawTransferTransactionRequestDescriptor instead')
const GetRawTransferTransactionRequest$json = const {
  '1': 'GetRawTransferTransactionRequest',
  '2': const [
    const {'1': 'lock_time', '3': 1, '4': 1, '5': 13, '10': 'lockTime'},
    const {'1': 'sender', '3': 2, '4': 1, '5': 9, '10': 'sender'},
    const {'1': 'receiver', '3': 3, '4': 1, '5': 9, '10': 'receiver'},
    const {'1': 'amount', '3': 4, '4': 1, '5': 3, '10': 'amount'},
    const {'1': 'fee', '3': 5, '4': 1, '5': 3, '10': 'fee'},
    const {'1': 'memo', '3': 6, '4': 1, '5': 9, '10': 'memo'},
  ],
};

/// Descriptor for `GetRawTransferTransactionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getRawTransferTransactionRequestDescriptor = $convert.base64Decode('CiBHZXRSYXdUcmFuc2ZlclRyYW5zYWN0aW9uUmVxdWVzdBIbCglsb2NrX3RpbWUYASABKA1SCGxvY2tUaW1lEhYKBnNlbmRlchgCIAEoCVIGc2VuZGVyEhoKCHJlY2VpdmVyGAMgASgJUghyZWNlaXZlchIWCgZhbW91bnQYBCABKANSBmFtb3VudBIQCgNmZWUYBSABKANSA2ZlZRISCgRtZW1vGAYgASgJUgRtZW1v');
@$core.Deprecated('Use getRawBondTransactionRequestDescriptor instead')
const GetRawBondTransactionRequest$json = const {
  '1': 'GetRawBondTransactionRequest',
  '2': const [
    const {'1': 'lock_time', '3': 1, '4': 1, '5': 13, '10': 'lockTime'},
    const {'1': 'sender', '3': 2, '4': 1, '5': 9, '10': 'sender'},
    const {'1': 'receiver', '3': 3, '4': 1, '5': 9, '10': 'receiver'},
    const {'1': 'stake', '3': 4, '4': 1, '5': 3, '10': 'stake'},
    const {'1': 'public_key', '3': 5, '4': 1, '5': 9, '10': 'publicKey'},
    const {'1': 'fee', '3': 6, '4': 1, '5': 3, '10': 'fee'},
    const {'1': 'memo', '3': 7, '4': 1, '5': 9, '10': 'memo'},
  ],
};

/// Descriptor for `GetRawBondTransactionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getRawBondTransactionRequestDescriptor = $convert.base64Decode('ChxHZXRSYXdCb25kVHJhbnNhY3Rpb25SZXF1ZXN0EhsKCWxvY2tfdGltZRgBIAEoDVIIbG9ja1RpbWUSFgoGc2VuZGVyGAIgASgJUgZzZW5kZXISGgoIcmVjZWl2ZXIYAyABKAlSCHJlY2VpdmVyEhQKBXN0YWtlGAQgASgDUgVzdGFrZRIdCgpwdWJsaWNfa2V5GAUgASgJUglwdWJsaWNLZXkSEAoDZmVlGAYgASgDUgNmZWUSEgoEbWVtbxgHIAEoCVIEbWVtbw==');
@$core.Deprecated('Use getRawUnBondTransactionRequestDescriptor instead')
const GetRawUnBondTransactionRequest$json = const {
  '1': 'GetRawUnBondTransactionRequest',
  '2': const [
    const {'1': 'lock_time', '3': 1, '4': 1, '5': 13, '10': 'lockTime'},
    const {'1': 'validator_address', '3': 3, '4': 1, '5': 9, '10': 'validatorAddress'},
    const {'1': 'memo', '3': 4, '4': 1, '5': 9, '10': 'memo'},
  ],
};

/// Descriptor for `GetRawUnBondTransactionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getRawUnBondTransactionRequestDescriptor = $convert.base64Decode('Ch5HZXRSYXdVbkJvbmRUcmFuc2FjdGlvblJlcXVlc3QSGwoJbG9ja190aW1lGAEgASgNUghsb2NrVGltZRIrChF2YWxpZGF0b3JfYWRkcmVzcxgDIAEoCVIQdmFsaWRhdG9yQWRkcmVzcxISCgRtZW1vGAQgASgJUgRtZW1v');
@$core.Deprecated('Use getRawWithdrawTransactionRequestDescriptor instead')
const GetRawWithdrawTransactionRequest$json = const {
  '1': 'GetRawWithdrawTransactionRequest',
  '2': const [
    const {'1': 'lock_time', '3': 1, '4': 1, '5': 13, '10': 'lockTime'},
    const {'1': 'validator_address', '3': 2, '4': 1, '5': 9, '10': 'validatorAddress'},
    const {'1': 'account_address', '3': 3, '4': 1, '5': 9, '10': 'accountAddress'},
    const {'1': 'fee', '3': 4, '4': 1, '5': 3, '10': 'fee'},
    const {'1': 'amount', '3': 5, '4': 1, '5': 3, '10': 'amount'},
    const {'1': 'memo', '3': 6, '4': 1, '5': 9, '10': 'memo'},
  ],
};

/// Descriptor for `GetRawWithdrawTransactionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getRawWithdrawTransactionRequestDescriptor = $convert.base64Decode('CiBHZXRSYXdXaXRoZHJhd1RyYW5zYWN0aW9uUmVxdWVzdBIbCglsb2NrX3RpbWUYASABKA1SCGxvY2tUaW1lEisKEXZhbGlkYXRvcl9hZGRyZXNzGAIgASgJUhB2YWxpZGF0b3JBZGRyZXNzEicKD2FjY291bnRfYWRkcmVzcxgDIAEoCVIOYWNjb3VudEFkZHJlc3MSEAoDZmVlGAQgASgDUgNmZWUSFgoGYW1vdW50GAUgASgDUgZhbW91bnQSEgoEbWVtbxgGIAEoCVIEbWVtbw==');
@$core.Deprecated('Use getRawTransactionResponseDescriptor instead')
const GetRawTransactionResponse$json = const {
  '1': 'GetRawTransactionResponse',
  '2': const [
    const {'1': 'raw_transaction', '3': 1, '4': 1, '5': 12, '10': 'rawTransaction'},
  ],
};

/// Descriptor for `GetRawTransactionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getRawTransactionResponseDescriptor = $convert.base64Decode('ChlHZXRSYXdUcmFuc2FjdGlvblJlc3BvbnNlEicKD3Jhd190cmFuc2FjdGlvbhgBIAEoDFIOcmF3VHJhbnNhY3Rpb24=');
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
    const {'1': 'lock_time', '3': 4, '4': 1, '5': 13, '10': 'lockTime'},
    const {'1': 'value', '3': 5, '4': 1, '5': 3, '10': 'value'},
    const {'1': 'fee', '3': 6, '4': 1, '5': 3, '10': 'fee'},
    const {'1': 'payload_type', '3': 7, '4': 1, '5': 14, '6': '.pactus.PayloadType', '10': 'payloadType'},
    const {'1': 'transfer', '3': 30, '4': 1, '5': 11, '6': '.pactus.PayloadTransfer', '9': 0, '10': 'transfer'},
    const {'1': 'bond', '3': 31, '4': 1, '5': 11, '6': '.pactus.PayloadBond', '9': 0, '10': 'bond'},
    const {'1': 'sortition', '3': 32, '4': 1, '5': 11, '6': '.pactus.PayloadSortition', '9': 0, '10': 'sortition'},
    const {'1': 'unbond', '3': 33, '4': 1, '5': 11, '6': '.pactus.PayloadUnbond', '9': 0, '10': 'unbond'},
    const {'1': 'withdraw', '3': 34, '4': 1, '5': 11, '6': '.pactus.PayloadWithdraw', '9': 0, '10': 'withdraw'},
    const {'1': 'memo', '3': 8, '4': 1, '5': 9, '10': 'memo'},
    const {'1': 'public_key', '3': 9, '4': 1, '5': 9, '10': 'publicKey'},
    const {'1': 'signature', '3': 10, '4': 1, '5': 12, '10': 'signature'},
  ],
  '8': const [
    const {'1': 'payload'},
  ],
};

/// Descriptor for `TransactionInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List transactionInfoDescriptor = $convert.base64Decode('Cg9UcmFuc2FjdGlvbkluZm8SDgoCaWQYASABKAxSAmlkEhIKBGRhdGEYAiABKAxSBGRhdGESGAoHdmVyc2lvbhgDIAEoBVIHdmVyc2lvbhIbCglsb2NrX3RpbWUYBCABKA1SCGxvY2tUaW1lEhQKBXZhbHVlGAUgASgDUgV2YWx1ZRIQCgNmZWUYBiABKANSA2ZlZRI2CgxwYXlsb2FkX3R5cGUYByABKA4yEy5wYWN0dXMuUGF5bG9hZFR5cGVSC3BheWxvYWRUeXBlEjUKCHRyYW5zZmVyGB4gASgLMhcucGFjdHVzLlBheWxvYWRUcmFuc2ZlckgAUgh0cmFuc2ZlchIpCgRib25kGB8gASgLMhMucGFjdHVzLlBheWxvYWRCb25kSABSBGJvbmQSOAoJc29ydGl0aW9uGCAgASgLMhgucGFjdHVzLlBheWxvYWRTb3J0aXRpb25IAFIJc29ydGl0aW9uEi8KBnVuYm9uZBghIAEoCzIVLnBhY3R1cy5QYXlsb2FkVW5ib25kSABSBnVuYm9uZBI1Cgh3aXRoZHJhdxgiIAEoCzIXLnBhY3R1cy5QYXlsb2FkV2l0aGRyYXdIAFIId2l0aGRyYXcSEgoEbWVtbxgIIAEoCVIEbWVtbxIdCgpwdWJsaWNfa2V5GAkgASgJUglwdWJsaWNLZXkSHAoJc2lnbmF0dXJlGAogASgMUglzaWduYXR1cmVCCQoHcGF5bG9hZA==');
const $core.Map<$core.String, $core.dynamic> TransactionServiceBase$json = const {
  '1': 'Transaction',
  '2': const [
    const {'1': 'GetTransaction', '2': '.pactus.GetTransactionRequest', '3': '.pactus.GetTransactionResponse'},
    const {'1': 'CalculateFee', '2': '.pactus.CalculateFeeRequest', '3': '.pactus.CalculateFeeResponse'},
    const {'1': 'BroadcastTransaction', '2': '.pactus.BroadcastTransactionRequest', '3': '.pactus.BroadcastTransactionResponse'},
    const {'1': 'GetRawTransferTransaction', '2': '.pactus.GetRawTransferTransactionRequest', '3': '.pactus.GetRawTransactionResponse'},
    const {'1': 'GetRawBondTransaction', '2': '.pactus.GetRawBondTransactionRequest', '3': '.pactus.GetRawTransactionResponse'},
    const {'1': 'GetRawUnBondTransaction', '2': '.pactus.GetRawUnBondTransactionRequest', '3': '.pactus.GetRawTransactionResponse'},
    const {'1': 'GetRawWithdrawTransaction', '2': '.pactus.GetRawWithdrawTransactionRequest', '3': '.pactus.GetRawTransactionResponse'},
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
  '.pactus.CalculateFeeRequest': CalculateFeeRequest$json,
  '.pactus.CalculateFeeResponse': CalculateFeeResponse$json,
  '.pactus.BroadcastTransactionRequest': BroadcastTransactionRequest$json,
  '.pactus.BroadcastTransactionResponse': BroadcastTransactionResponse$json,
  '.pactus.GetRawTransferTransactionRequest': GetRawTransferTransactionRequest$json,
  '.pactus.GetRawTransactionResponse': GetRawTransactionResponse$json,
  '.pactus.GetRawBondTransactionRequest': GetRawBondTransactionRequest$json,
  '.pactus.GetRawUnBondTransactionRequest': GetRawUnBondTransactionRequest$json,
  '.pactus.GetRawWithdrawTransactionRequest': GetRawWithdrawTransactionRequest$json,
};

/// Descriptor for `Transaction`. Decode as a `google.protobuf.ServiceDescriptorProto`.
final $typed_data.Uint8List transactionServiceDescriptor = $convert.base64Decode('CgtUcmFuc2FjdGlvbhJPCg5HZXRUcmFuc2FjdGlvbhIdLnBhY3R1cy5HZXRUcmFuc2FjdGlvblJlcXVlc3QaHi5wYWN0dXMuR2V0VHJhbnNhY3Rpb25SZXNwb25zZRJJCgxDYWxjdWxhdGVGZWUSGy5wYWN0dXMuQ2FsY3VsYXRlRmVlUmVxdWVzdBocLnBhY3R1cy5DYWxjdWxhdGVGZWVSZXNwb25zZRJhChRCcm9hZGNhc3RUcmFuc2FjdGlvbhIjLnBhY3R1cy5Ccm9hZGNhc3RUcmFuc2FjdGlvblJlcXVlc3QaJC5wYWN0dXMuQnJvYWRjYXN0VHJhbnNhY3Rpb25SZXNwb25zZRJoChlHZXRSYXdUcmFuc2ZlclRyYW5zYWN0aW9uEigucGFjdHVzLkdldFJhd1RyYW5zZmVyVHJhbnNhY3Rpb25SZXF1ZXN0GiEucGFjdHVzLkdldFJhd1RyYW5zYWN0aW9uUmVzcG9uc2USYAoVR2V0UmF3Qm9uZFRyYW5zYWN0aW9uEiQucGFjdHVzLkdldFJhd0JvbmRUcmFuc2FjdGlvblJlcXVlc3QaIS5wYWN0dXMuR2V0UmF3VHJhbnNhY3Rpb25SZXNwb25zZRJkChdHZXRSYXdVbkJvbmRUcmFuc2FjdGlvbhImLnBhY3R1cy5HZXRSYXdVbkJvbmRUcmFuc2FjdGlvblJlcXVlc3QaIS5wYWN0dXMuR2V0UmF3VHJhbnNhY3Rpb25SZXNwb25zZRJoChlHZXRSYXdXaXRoZHJhd1RyYW5zYWN0aW9uEigucGFjdHVzLkdldFJhd1dpdGhkcmF3VHJhbnNhY3Rpb25SZXF1ZXN0GiEucGFjdHVzLkdldFJhd1RyYW5zYWN0aW9uUmVzcG9uc2U=');
