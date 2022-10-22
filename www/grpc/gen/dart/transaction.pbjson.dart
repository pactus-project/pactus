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
    const {'1': 'SEND_PAYLOAD', '2': 1},
    const {'1': 'BOND_PAYLOAD', '2': 2},
    const {'1': 'SORTITION_PAYLOAD', '2': 3},
    const {'1': 'UNBOND_PAYLOAD', '2': 4},
    const {'1': 'WITHDRAW_PAYLOAD', '2': 5},
  ],
};

/// Descriptor for `PayloadType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List payloadTypeDescriptor = $convert.base64Decode('CgtQYXlsb2FkVHlwZRILCgdVTktOT1dOEAASEAoMU0VORF9QQVlMT0FEEAESEAoMQk9ORF9QQVlMT0FEEAISFQoRU09SVElUSU9OX1BBWUxPQUQQAxISCg5VTkJPTkRfUEFZTE9BRBAEEhQKEFdJVEhEUkFXX1BBWUxPQUQQBQ==');
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
@$core.Deprecated('Use transactionRequestDescriptor instead')
const TransactionRequest$json = const {
  '1': 'TransactionRequest',
  '2': const [
    const {'1': 'id', '3': 1, '4': 1, '5': 12, '10': 'id'},
    const {'1': 'verbosity', '3': 2, '4': 1, '5': 14, '6': '.pactus.TransactionVerbosity', '10': 'verbosity'},
  ],
};

/// Descriptor for `TransactionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List transactionRequestDescriptor = $convert.base64Decode('ChJUcmFuc2FjdGlvblJlcXVlc3QSDgoCaWQYASABKAxSAmlkEjoKCXZlcmJvc2l0eRgCIAEoDjIcLnBhY3R1cy5UcmFuc2FjdGlvblZlcmJvc2l0eVIJdmVyYm9zaXR5');
@$core.Deprecated('Use transactionResponseDescriptor instead')
const TransactionResponse$json = const {
  '1': 'TransactionResponse',
  '2': const [
    const {'1': 'block_height', '3': 12, '4': 1, '5': 13, '10': 'blockHeight'},
    const {'1': 'block_time', '3': 13, '4': 1, '5': 13, '10': 'blockTime'},
    const {'1': 'transaction', '3': 3, '4': 1, '5': 11, '6': '.pactus.TransactionInfo', '10': 'transaction'},
  ],
};

/// Descriptor for `TransactionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List transactionResponseDescriptor = $convert.base64Decode('ChNUcmFuc2FjdGlvblJlc3BvbnNlEiEKDGJsb2NrX2hlaWdodBgMIAEoDVILYmxvY2tIZWlnaHQSHQoKYmxvY2tfdGltZRgNIAEoDVIJYmxvY2tUaW1lEjkKC3RyYW5zYWN0aW9uGAMgASgLMhcucGFjdHVzLlRyYW5zYWN0aW9uSW5mb1ILdHJhbnNhY3Rpb24=');
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
@$core.Deprecated('Use payloadSendDescriptor instead')
const PayloadSend$json = const {
  '1': 'PayloadSend',
  '2': const [
    const {'1': 'sender', '3': 1, '4': 1, '5': 9, '10': 'sender'},
    const {'1': 'receiver', '3': 2, '4': 1, '5': 9, '10': 'receiver'},
    const {'1': 'amount', '3': 3, '4': 1, '5': 3, '10': 'amount'},
  ],
};

/// Descriptor for `PayloadSend`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List payloadSendDescriptor = $convert.base64Decode('CgtQYXlsb2FkU2VuZBIWCgZzZW5kZXIYASABKAlSBnNlbmRlchIaCghyZWNlaXZlchgCIAEoCVIIcmVjZWl2ZXISFgoGYW1vdW50GAMgASgDUgZhbW91bnQ=');
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
    const {'1': 'Type', '3': 8, '4': 1, '5': 14, '6': '.pactus.PayloadType', '10': 'Type'},
    const {'1': 'send', '3': 30, '4': 1, '5': 11, '6': '.pactus.PayloadSend', '9': 0, '10': 'send'},
    const {'1': 'bond', '3': 31, '4': 1, '5': 11, '6': '.pactus.PayloadBond', '9': 0, '10': 'bond'},
    const {'1': 'sortition', '3': 32, '4': 1, '5': 11, '6': '.pactus.PayloadSortition', '9': 0, '10': 'sortition'},
    const {'1': 'memo', '3': 9, '4': 1, '5': 9, '10': 'memo'},
    const {'1': 'public_key', '3': 10, '4': 1, '5': 9, '10': 'publicKey'},
    const {'1': 'signature', '3': 11, '4': 1, '5': 12, '10': 'signature'},
  ],
  '8': const [
    const {'1': 'Payload'},
  ],
};

/// Descriptor for `TransactionInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List transactionInfoDescriptor = $convert.base64Decode('Cg9UcmFuc2FjdGlvbkluZm8SDgoCaWQYASABKAxSAmlkEhIKBGRhdGEYAiABKAxSBGRhdGESGAoHdmVyc2lvbhgDIAEoBVIHdmVyc2lvbhIUCgVzdGFtcBgEIAEoDFIFc3RhbXASGgoIc2VxdWVuY2UYBSABKAVSCHNlcXVlbmNlEhQKBXZhbHVlGAYgASgDUgV2YWx1ZRIQCgNmZWUYByABKANSA2ZlZRInCgRUeXBlGAggASgOMhMucGFjdHVzLlBheWxvYWRUeXBlUgRUeXBlEikKBHNlbmQYHiABKAsyEy5wYWN0dXMuUGF5bG9hZFNlbmRIAFIEc2VuZBIpCgRib25kGB8gASgLMhMucGFjdHVzLlBheWxvYWRCb25kSABSBGJvbmQSOAoJc29ydGl0aW9uGCAgASgLMhgucGFjdHVzLlBheWxvYWRTb3J0aXRpb25IAFIJc29ydGl0aW9uEhIKBG1lbW8YCSABKAlSBG1lbW8SHQoKcHVibGljX2tleRgKIAEoCVIJcHVibGljS2V5EhwKCXNpZ25hdHVyZRgLIAEoDFIJc2lnbmF0dXJlQgkKB1BheWxvYWQ=');
const $core.Map<$core.String, $core.dynamic> TransactionServiceBase$json = const {
  '1': 'Transaction',
  '2': const [
    const {'1': 'GetTransaction', '2': '.pactus.TransactionRequest', '3': '.pactus.TransactionResponse', '4': const {}},
    const {'1': 'SendRawTransaction', '2': '.pactus.SendRawTransactionRequest', '3': '.pactus.SendRawTransactionResponse', '4': const {}},
  ],
};

@$core.Deprecated('Use transactionServiceDescriptor instead')
const $core.Map<$core.String, $core.Map<$core.String, $core.dynamic>> TransactionServiceBase$messageJson = const {
  '.pactus.TransactionRequest': TransactionRequest$json,
  '.pactus.TransactionResponse': TransactionResponse$json,
  '.pactus.TransactionInfo': TransactionInfo$json,
  '.pactus.PayloadSend': PayloadSend$json,
  '.pactus.PayloadBond': PayloadBond$json,
  '.pactus.PayloadSortition': PayloadSortition$json,
  '.pactus.SendRawTransactionRequest': SendRawTransactionRequest$json,
  '.pactus.SendRawTransactionResponse': SendRawTransactionResponse$json,
};

/// Descriptor for `Transaction`. Decode as a `google.protobuf.ServiceDescriptorProto`.
final $typed_data.Uint8List transactionServiceDescriptor = $convert.base64Decode('CgtUcmFuc2FjdGlvbhKAAQoOR2V0VHJhbnNhY3Rpb24SGi5wYWN0dXMuVHJhbnNhY3Rpb25SZXF1ZXN0GhsucGFjdHVzLlRyYW5zYWN0aW9uUmVzcG9uc2UiNYLT5JMCLxItL3YxL3RyYW5zYWN0aW9uL2lkL3tpZH0vdmVyYm9zaXR5L3t2ZXJib3NpdHl9EpABChJTZW5kUmF3VHJhbnNhY3Rpb24SIS5wYWN0dXMuU2VuZFJhd1RyYW5zYWN0aW9uUmVxdWVzdBoiLnBhY3R1cy5TZW5kUmF3VHJhbnNhY3Rpb25SZXNwb25zZSIzgtPkkwItGisvdjEvdHJhbnNhY3Rpb24vc2VuZF9yYXdfdHJhbnNhY3Rpb24ve2RhdGF9');
