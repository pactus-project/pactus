//
//  Generated code. Do not modify.
//  source: wallet.proto
//
// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_final_fields
// ignore_for_file: unnecessary_import, unnecessary_this, unused_import

import 'dart:convert' as $convert;
import 'dart:core' as $core;
import 'dart:typed_data' as $typed_data;

@$core.Deprecated('Use addressTypeDescriptor instead')
const AddressType$json = {
  '1': 'AddressType',
  '2': [
    {'1': 'ADDRESS_TYPE_TREASURY', '2': 0},
    {'1': 'ADDRESS_TYPE_VALIDATOR', '2': 1},
    {'1': 'ADDRESS_TYPE_BLS_ACCOUNT', '2': 2},
    {'1': 'ADDRESS_TYPE_ED25519_ACCOUNT', '2': 3},
  ],
};

/// Descriptor for `AddressType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List addressTypeDescriptor = $convert.base64Decode(
    'CgtBZGRyZXNzVHlwZRIZChVBRERSRVNTX1RZUEVfVFJFQVNVUlkQABIaChZBRERSRVNTX1RZUE'
    'VfVkFMSURBVE9SEAESHAoYQUREUkVTU19UWVBFX0JMU19BQ0NPVU5UEAISIAocQUREUkVTU19U'
    'WVBFX0VEMjU1MTlfQUNDT1VOVBAD');

@$core.Deprecated('Use addressInfoDescriptor instead')
const AddressInfo$json = {
  '1': 'AddressInfo',
  '2': [
    {'1': 'address', '3': 1, '4': 1, '5': 9, '10': 'address'},
    {'1': 'public_key', '3': 2, '4': 1, '5': 9, '10': 'publicKey'},
    {'1': 'label', '3': 3, '4': 1, '5': 9, '10': 'label'},
    {'1': 'path', '3': 4, '4': 1, '5': 9, '10': 'path'},
  ],
};

/// Descriptor for `AddressInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List addressInfoDescriptor = $convert.base64Decode(
    'CgtBZGRyZXNzSW5mbxIYCgdhZGRyZXNzGAEgASgJUgdhZGRyZXNzEh0KCnB1YmxpY19rZXkYAi'
    'ABKAlSCXB1YmxpY0tleRIUCgVsYWJlbBgDIAEoCVIFbGFiZWwSEgoEcGF0aBgEIAEoCVIEcGF0'
    'aA==');

@$core.Deprecated('Use historyInfoDescriptor instead')
const HistoryInfo$json = {
  '1': 'HistoryInfo',
  '2': [
    {'1': 'transaction_id', '3': 1, '4': 1, '5': 9, '10': 'transactionId'},
    {'1': 'time', '3': 2, '4': 1, '5': 13, '10': 'time'},
    {'1': 'payload_type', '3': 3, '4': 1, '5': 9, '10': 'payloadType'},
    {'1': 'description', '3': 4, '4': 1, '5': 9, '10': 'description'},
    {'1': 'amount', '3': 5, '4': 1, '5': 3, '10': 'amount'},
  ],
};

/// Descriptor for `HistoryInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List historyInfoDescriptor = $convert.base64Decode(
    'CgtIaXN0b3J5SW5mbxIlCg50cmFuc2FjdGlvbl9pZBgBIAEoCVINdHJhbnNhY3Rpb25JZBISCg'
    'R0aW1lGAIgASgNUgR0aW1lEiEKDHBheWxvYWRfdHlwZRgDIAEoCVILcGF5bG9hZFR5cGUSIAoL'
    'ZGVzY3JpcHRpb24YBCABKAlSC2Rlc2NyaXB0aW9uEhYKBmFtb3VudBgFIAEoA1IGYW1vdW50');

@$core.Deprecated('Use getAddressHistoryRequestDescriptor instead')
const GetAddressHistoryRequest$json = {
  '1': 'GetAddressHistoryRequest',
  '2': [
    {'1': 'wallet_name', '3': 1, '4': 1, '5': 9, '10': 'walletName'},
    {'1': 'address', '3': 2, '4': 1, '5': 9, '10': 'address'},
  ],
};

/// Descriptor for `GetAddressHistoryRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getAddressHistoryRequestDescriptor = $convert.base64Decode(
    'ChhHZXRBZGRyZXNzSGlzdG9yeVJlcXVlc3QSHwoLd2FsbGV0X25hbWUYASABKAlSCndhbGxldE'
    '5hbWUSGAoHYWRkcmVzcxgCIAEoCVIHYWRkcmVzcw==');

@$core.Deprecated('Use getAddressHistoryResponseDescriptor instead')
const GetAddressHistoryResponse$json = {
  '1': 'GetAddressHistoryResponse',
  '2': [
    {'1': 'history_info', '3': 1, '4': 3, '5': 11, '6': '.pactus.HistoryInfo', '10': 'historyInfo'},
  ],
};

/// Descriptor for `GetAddressHistoryResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getAddressHistoryResponseDescriptor = $convert.base64Decode(
    'ChlHZXRBZGRyZXNzSGlzdG9yeVJlc3BvbnNlEjYKDGhpc3RvcnlfaW5mbxgBIAMoCzITLnBhY3'
    'R1cy5IaXN0b3J5SW5mb1ILaGlzdG9yeUluZm8=');

@$core.Deprecated('Use getNewAddressRequestDescriptor instead')
const GetNewAddressRequest$json = {
  '1': 'GetNewAddressRequest',
  '2': [
    {'1': 'wallet_name', '3': 1, '4': 1, '5': 9, '10': 'walletName'},
    {'1': 'address_type', '3': 2, '4': 1, '5': 14, '6': '.pactus.AddressType', '10': 'addressType'},
    {'1': 'label', '3': 3, '4': 1, '5': 9, '10': 'label'},
    {'1': 'password', '3': 4, '4': 1, '5': 9, '10': 'password'},
  ],
};

/// Descriptor for `GetNewAddressRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getNewAddressRequestDescriptor = $convert.base64Decode(
    'ChRHZXROZXdBZGRyZXNzUmVxdWVzdBIfCgt3YWxsZXRfbmFtZRgBIAEoCVIKd2FsbGV0TmFtZR'
    'I2CgxhZGRyZXNzX3R5cGUYAiABKA4yEy5wYWN0dXMuQWRkcmVzc1R5cGVSC2FkZHJlc3NUeXBl'
    'EhQKBWxhYmVsGAMgASgJUgVsYWJlbBIaCghwYXNzd29yZBgEIAEoCVIIcGFzc3dvcmQ=');

@$core.Deprecated('Use getNewAddressResponseDescriptor instead')
const GetNewAddressResponse$json = {
  '1': 'GetNewAddressResponse',
  '2': [
    {'1': 'wallet_name', '3': 1, '4': 1, '5': 9, '10': 'walletName'},
    {'1': 'address_info', '3': 2, '4': 1, '5': 11, '6': '.pactus.AddressInfo', '10': 'addressInfo'},
  ],
};

/// Descriptor for `GetNewAddressResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getNewAddressResponseDescriptor = $convert.base64Decode(
    'ChVHZXROZXdBZGRyZXNzUmVzcG9uc2USHwoLd2FsbGV0X25hbWUYASABKAlSCndhbGxldE5hbW'
    'USNgoMYWRkcmVzc19pbmZvGAIgASgLMhMucGFjdHVzLkFkZHJlc3NJbmZvUgthZGRyZXNzSW5m'
    'bw==');

@$core.Deprecated('Use restoreWalletRequestDescriptor instead')
const RestoreWalletRequest$json = {
  '1': 'RestoreWalletRequest',
  '2': [
    {'1': 'wallet_name', '3': 1, '4': 1, '5': 9, '10': 'walletName'},
    {'1': 'mnemonic', '3': 2, '4': 1, '5': 9, '10': 'mnemonic'},
    {'1': 'password', '3': 3, '4': 1, '5': 9, '10': 'password'},
  ],
};

/// Descriptor for `RestoreWalletRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List restoreWalletRequestDescriptor = $convert.base64Decode(
    'ChRSZXN0b3JlV2FsbGV0UmVxdWVzdBIfCgt3YWxsZXRfbmFtZRgBIAEoCVIKd2FsbGV0TmFtZR'
    'IaCghtbmVtb25pYxgCIAEoCVIIbW5lbW9uaWMSGgoIcGFzc3dvcmQYAyABKAlSCHBhc3N3b3Jk');

@$core.Deprecated('Use restoreWalletResponseDescriptor instead')
const RestoreWalletResponse$json = {
  '1': 'RestoreWalletResponse',
  '2': [
    {'1': 'wallet_name', '3': 1, '4': 1, '5': 9, '10': 'walletName'},
  ],
};

/// Descriptor for `RestoreWalletResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List restoreWalletResponseDescriptor = $convert.base64Decode(
    'ChVSZXN0b3JlV2FsbGV0UmVzcG9uc2USHwoLd2FsbGV0X25hbWUYASABKAlSCndhbGxldE5hbW'
    'U=');

@$core.Deprecated('Use createWalletRequestDescriptor instead')
const CreateWalletRequest$json = {
  '1': 'CreateWalletRequest',
  '2': [
    {'1': 'wallet_name', '3': 1, '4': 1, '5': 9, '10': 'walletName'},
    {'1': 'password', '3': 4, '4': 1, '5': 9, '10': 'password'},
  ],
};

/// Descriptor for `CreateWalletRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createWalletRequestDescriptor = $convert.base64Decode(
    'ChNDcmVhdGVXYWxsZXRSZXF1ZXN0Eh8KC3dhbGxldF9uYW1lGAEgASgJUgp3YWxsZXROYW1lEh'
    'oKCHBhc3N3b3JkGAQgASgJUghwYXNzd29yZA==');

@$core.Deprecated('Use createWalletResponseDescriptor instead')
const CreateWalletResponse$json = {
  '1': 'CreateWalletResponse',
  '2': [
    {'1': 'mnemonic', '3': 2, '4': 1, '5': 9, '10': 'mnemonic'},
  ],
};

/// Descriptor for `CreateWalletResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createWalletResponseDescriptor = $convert.base64Decode(
    'ChRDcmVhdGVXYWxsZXRSZXNwb25zZRIaCghtbmVtb25pYxgCIAEoCVIIbW5lbW9uaWM=');

@$core.Deprecated('Use loadWalletRequestDescriptor instead')
const LoadWalletRequest$json = {
  '1': 'LoadWalletRequest',
  '2': [
    {'1': 'wallet_name', '3': 1, '4': 1, '5': 9, '10': 'walletName'},
  ],
};

/// Descriptor for `LoadWalletRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List loadWalletRequestDescriptor = $convert.base64Decode(
    'ChFMb2FkV2FsbGV0UmVxdWVzdBIfCgt3YWxsZXRfbmFtZRgBIAEoCVIKd2FsbGV0TmFtZQ==');

@$core.Deprecated('Use loadWalletResponseDescriptor instead')
const LoadWalletResponse$json = {
  '1': 'LoadWalletResponse',
  '2': [
    {'1': 'wallet_name', '3': 1, '4': 1, '5': 9, '10': 'walletName'},
  ],
};

/// Descriptor for `LoadWalletResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List loadWalletResponseDescriptor = $convert.base64Decode(
    'ChJMb2FkV2FsbGV0UmVzcG9uc2USHwoLd2FsbGV0X25hbWUYASABKAlSCndhbGxldE5hbWU=');

@$core.Deprecated('Use unloadWalletRequestDescriptor instead')
const UnloadWalletRequest$json = {
  '1': 'UnloadWalletRequest',
  '2': [
    {'1': 'wallet_name', '3': 1, '4': 1, '5': 9, '10': 'walletName'},
  ],
};

/// Descriptor for `UnloadWalletRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List unloadWalletRequestDescriptor = $convert.base64Decode(
    'ChNVbmxvYWRXYWxsZXRSZXF1ZXN0Eh8KC3dhbGxldF9uYW1lGAEgASgJUgp3YWxsZXROYW1l');

@$core.Deprecated('Use unloadWalletResponseDescriptor instead')
const UnloadWalletResponse$json = {
  '1': 'UnloadWalletResponse',
  '2': [
    {'1': 'wallet_name', '3': 1, '4': 1, '5': 9, '10': 'walletName'},
  ],
};

/// Descriptor for `UnloadWalletResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List unloadWalletResponseDescriptor = $convert.base64Decode(
    'ChRVbmxvYWRXYWxsZXRSZXNwb25zZRIfCgt3YWxsZXRfbmFtZRgBIAEoCVIKd2FsbGV0TmFtZQ'
    '==');

@$core.Deprecated('Use getValidatorAddressRequestDescriptor instead')
const GetValidatorAddressRequest$json = {
  '1': 'GetValidatorAddressRequest',
  '2': [
    {'1': 'public_key', '3': 1, '4': 1, '5': 9, '10': 'publicKey'},
  ],
};

/// Descriptor for `GetValidatorAddressRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getValidatorAddressRequestDescriptor = $convert.base64Decode(
    'ChpHZXRWYWxpZGF0b3JBZGRyZXNzUmVxdWVzdBIdCgpwdWJsaWNfa2V5GAEgASgJUglwdWJsaW'
    'NLZXk=');

@$core.Deprecated('Use getValidatorAddressResponseDescriptor instead')
const GetValidatorAddressResponse$json = {
  '1': 'GetValidatorAddressResponse',
  '2': [
    {'1': 'address', '3': 1, '4': 1, '5': 9, '10': 'address'},
  ],
};

/// Descriptor for `GetValidatorAddressResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getValidatorAddressResponseDescriptor = $convert.base64Decode(
    'ChtHZXRWYWxpZGF0b3JBZGRyZXNzUmVzcG9uc2USGAoHYWRkcmVzcxgBIAEoCVIHYWRkcmVzcw'
    '==');

@$core.Deprecated('Use signRawTransactionRequestDescriptor instead')
const SignRawTransactionRequest$json = {
  '1': 'SignRawTransactionRequest',
  '2': [
    {'1': 'wallet_name', '3': 1, '4': 1, '5': 9, '10': 'walletName'},
    {'1': 'raw_transaction', '3': 2, '4': 1, '5': 9, '10': 'rawTransaction'},
    {'1': 'password', '3': 3, '4': 1, '5': 9, '10': 'password'},
  ],
};

/// Descriptor for `SignRawTransactionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List signRawTransactionRequestDescriptor = $convert.base64Decode(
    'ChlTaWduUmF3VHJhbnNhY3Rpb25SZXF1ZXN0Eh8KC3dhbGxldF9uYW1lGAEgASgJUgp3YWxsZX'
    'ROYW1lEicKD3Jhd190cmFuc2FjdGlvbhgCIAEoCVIOcmF3VHJhbnNhY3Rpb24SGgoIcGFzc3dv'
    'cmQYAyABKAlSCHBhc3N3b3Jk');

@$core.Deprecated('Use signRawTransactionResponseDescriptor instead')
const SignRawTransactionResponse$json = {
  '1': 'SignRawTransactionResponse',
  '2': [
    {'1': 'transaction_id', '3': 1, '4': 1, '5': 9, '10': 'transactionId'},
    {'1': 'signed_raw_transaction', '3': 2, '4': 1, '5': 9, '10': 'signedRawTransaction'},
  ],
};

/// Descriptor for `SignRawTransactionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List signRawTransactionResponseDescriptor = $convert.base64Decode(
    'ChpTaWduUmF3VHJhbnNhY3Rpb25SZXNwb25zZRIlCg50cmFuc2FjdGlvbl9pZBgBIAEoCVINdH'
    'JhbnNhY3Rpb25JZBI0ChZzaWduZWRfcmF3X3RyYW5zYWN0aW9uGAIgASgJUhRzaWduZWRSYXdU'
    'cmFuc2FjdGlvbg==');

@$core.Deprecated('Use getTotalBalanceRequestDescriptor instead')
const GetTotalBalanceRequest$json = {
  '1': 'GetTotalBalanceRequest',
  '2': [
    {'1': 'wallet_name', '3': 1, '4': 1, '5': 9, '10': 'walletName'},
  ],
};

/// Descriptor for `GetTotalBalanceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getTotalBalanceRequestDescriptor = $convert.base64Decode(
    'ChZHZXRUb3RhbEJhbGFuY2VSZXF1ZXN0Eh8KC3dhbGxldF9uYW1lGAEgASgJUgp3YWxsZXROYW'
    '1l');

@$core.Deprecated('Use getTotalBalanceResponseDescriptor instead')
const GetTotalBalanceResponse$json = {
  '1': 'GetTotalBalanceResponse',
  '2': [
    {'1': 'wallet_name', '3': 1, '4': 1, '5': 9, '10': 'walletName'},
    {'1': 'total_balance', '3': 2, '4': 1, '5': 3, '10': 'totalBalance'},
  ],
};

/// Descriptor for `GetTotalBalanceResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getTotalBalanceResponseDescriptor = $convert.base64Decode(
    'ChdHZXRUb3RhbEJhbGFuY2VSZXNwb25zZRIfCgt3YWxsZXRfbmFtZRgBIAEoCVIKd2FsbGV0Tm'
    'FtZRIjCg10b3RhbF9iYWxhbmNlGAIgASgDUgx0b3RhbEJhbGFuY2U=');

@$core.Deprecated('Use signMessageRequestDescriptor instead')
const SignMessageRequest$json = {
  '1': 'SignMessageRequest',
  '2': [
    {'1': 'wallet_name', '3': 1, '4': 1, '5': 9, '10': 'walletName'},
    {'1': 'password', '3': 2, '4': 1, '5': 9, '10': 'password'},
    {'1': 'address', '3': 3, '4': 1, '5': 9, '10': 'address'},
    {'1': 'message', '3': 4, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `SignMessageRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List signMessageRequestDescriptor = $convert.base64Decode(
    'ChJTaWduTWVzc2FnZVJlcXVlc3QSHwoLd2FsbGV0X25hbWUYASABKAlSCndhbGxldE5hbWUSGg'
    'oIcGFzc3dvcmQYAiABKAlSCHBhc3N3b3JkEhgKB2FkZHJlc3MYAyABKAlSB2FkZHJlc3MSGAoH'
    'bWVzc2FnZRgEIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use signMessageResponseDescriptor instead')
const SignMessageResponse$json = {
  '1': 'SignMessageResponse',
  '2': [
    {'1': 'signature', '3': 1, '4': 1, '5': 9, '10': 'signature'},
  ],
};

/// Descriptor for `SignMessageResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List signMessageResponseDescriptor = $convert.base64Decode(
    'ChNTaWduTWVzc2FnZVJlc3BvbnNlEhwKCXNpZ25hdHVyZRgBIAEoCVIJc2lnbmF0dXJl');

@$core.Deprecated('Use getTotalStakeRequestDescriptor instead')
const GetTotalStakeRequest$json = {
  '1': 'GetTotalStakeRequest',
  '2': [
    {'1': 'wallet_name', '3': 1, '4': 1, '5': 9, '10': 'walletName'},
  ],
};

/// Descriptor for `GetTotalStakeRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getTotalStakeRequestDescriptor = $convert.base64Decode(
    'ChRHZXRUb3RhbFN0YWtlUmVxdWVzdBIfCgt3YWxsZXRfbmFtZRgBIAEoCVIKd2FsbGV0TmFtZQ'
    '==');

@$core.Deprecated('Use getTotalStakeResponseDescriptor instead')
const GetTotalStakeResponse$json = {
  '1': 'GetTotalStakeResponse',
  '2': [
    {'1': 'wallet_name', '3': 1, '4': 1, '5': 9, '10': 'walletName'},
    {'1': 'total_stake', '3': 2, '4': 1, '5': 3, '10': 'totalStake'},
  ],
};

/// Descriptor for `GetTotalStakeResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getTotalStakeResponseDescriptor = $convert.base64Decode(
    'ChVHZXRUb3RhbFN0YWtlUmVzcG9uc2USHwoLd2FsbGV0X25hbWUYASABKAlSCndhbGxldE5hbW'
    'USHwoLdG90YWxfc3Rha2UYAiABKANSCnRvdGFsU3Rha2U=');

@$core.Deprecated('Use getAddressInfoRequestDescriptor instead')
const GetAddressInfoRequest$json = {
  '1': 'GetAddressInfoRequest',
  '2': [
    {'1': 'wallet_name', '3': 1, '4': 1, '5': 9, '10': 'walletName'},
    {'1': 'address', '3': 2, '4': 1, '5': 9, '10': 'address'},
  ],
};

/// Descriptor for `GetAddressInfoRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getAddressInfoRequestDescriptor = $convert.base64Decode(
    'ChVHZXRBZGRyZXNzSW5mb1JlcXVlc3QSHwoLd2FsbGV0X25hbWUYASABKAlSCndhbGxldE5hbW'
    'USGAoHYWRkcmVzcxgCIAEoCVIHYWRkcmVzcw==');

@$core.Deprecated('Use getAddressInfoResponseDescriptor instead')
const GetAddressInfoResponse$json = {
  '1': 'GetAddressInfoResponse',
  '2': [
    {'1': 'wallet_name', '3': 1, '4': 1, '5': 9, '10': 'walletName'},
    {'1': 'address', '3': 2, '4': 1, '5': 9, '10': 'address'},
    {'1': 'label', '3': 3, '4': 1, '5': 9, '10': 'label'},
    {'1': 'public_key', '3': 4, '4': 1, '5': 9, '10': 'publicKey'},
    {'1': 'path', '3': 5, '4': 1, '5': 9, '10': 'path'},
  ],
};

/// Descriptor for `GetAddressInfoResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getAddressInfoResponseDescriptor = $convert.base64Decode(
    'ChZHZXRBZGRyZXNzSW5mb1Jlc3BvbnNlEh8KC3dhbGxldF9uYW1lGAEgASgJUgp3YWxsZXROYW'
    '1lEhgKB2FkZHJlc3MYAiABKAlSB2FkZHJlc3MSFAoFbGFiZWwYAyABKAlSBWxhYmVsEh0KCnB1'
    'YmxpY19rZXkYBCABKAlSCXB1YmxpY0tleRISCgRwYXRoGAUgASgJUgRwYXRo');

@$core.Deprecated('Use setAddressLabelRequestDescriptor instead')
const SetAddressLabelRequest$json = {
  '1': 'SetAddressLabelRequest',
  '2': [
    {'1': 'wallet_name', '3': 1, '4': 1, '5': 9, '10': 'walletName'},
    {'1': 'password', '3': 3, '4': 1, '5': 9, '10': 'password'},
    {'1': 'address', '3': 4, '4': 1, '5': 9, '10': 'address'},
    {'1': 'label', '3': 5, '4': 1, '5': 9, '10': 'label'},
  ],
};

/// Descriptor for `SetAddressLabelRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List setAddressLabelRequestDescriptor = $convert.base64Decode(
    'ChZTZXRBZGRyZXNzTGFiZWxSZXF1ZXN0Eh8KC3dhbGxldF9uYW1lGAEgASgJUgp3YWxsZXROYW'
    '1lEhoKCHBhc3N3b3JkGAMgASgJUghwYXNzd29yZBIYCgdhZGRyZXNzGAQgASgJUgdhZGRyZXNz'
    'EhQKBWxhYmVsGAUgASgJUgVsYWJlbA==');

@$core.Deprecated('Use setAddressLabelResponseDescriptor instead')
const SetAddressLabelResponse$json = {
  '1': 'SetAddressLabelResponse',
};

/// Descriptor for `SetAddressLabelResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List setAddressLabelResponseDescriptor = $convert.base64Decode(
    'ChdTZXRBZGRyZXNzTGFiZWxSZXNwb25zZQ==');

@$core.Deprecated('Use listWalletRequestDescriptor instead')
const ListWalletRequest$json = {
  '1': 'ListWalletRequest',
};

/// Descriptor for `ListWalletRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listWalletRequestDescriptor = $convert.base64Decode(
    'ChFMaXN0V2FsbGV0UmVxdWVzdA==');

@$core.Deprecated('Use listWalletResponseDescriptor instead')
const ListWalletResponse$json = {
  '1': 'ListWalletResponse',
  '2': [
    {'1': 'wallets', '3': 1, '4': 3, '5': 9, '10': 'wallets'},
  ],
};

/// Descriptor for `ListWalletResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listWalletResponseDescriptor = $convert.base64Decode(
    'ChJMaXN0V2FsbGV0UmVzcG9uc2USGAoHd2FsbGV0cxgBIAMoCVIHd2FsbGV0cw==');

@$core.Deprecated('Use getWalletInfoRequestDescriptor instead')
const GetWalletInfoRequest$json = {
  '1': 'GetWalletInfoRequest',
  '2': [
    {'1': 'wallet_name', '3': 1, '4': 1, '5': 9, '10': 'walletName'},
  ],
};

/// Descriptor for `GetWalletInfoRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getWalletInfoRequestDescriptor = $convert.base64Decode(
    'ChRHZXRXYWxsZXRJbmZvUmVxdWVzdBIfCgt3YWxsZXRfbmFtZRgBIAEoCVIKd2FsbGV0TmFtZQ'
    '==');

@$core.Deprecated('Use getWalletInfoResponseDescriptor instead')
const GetWalletInfoResponse$json = {
  '1': 'GetWalletInfoResponse',
  '2': [
    {'1': 'wallet_name', '3': 1, '4': 1, '5': 9, '10': 'walletName'},
    {'1': 'version', '3': 2, '4': 1, '5': 5, '10': 'version'},
    {'1': 'network', '3': 3, '4': 1, '5': 9, '10': 'network'},
    {'1': 'encrypted', '3': 4, '4': 1, '5': 8, '10': 'encrypted'},
    {'1': 'uuid', '3': 5, '4': 1, '5': 9, '10': 'uuid'},
    {'1': 'created_at', '3': 6, '4': 1, '5': 3, '10': 'createdAt'},
    {'1': 'default_fee', '3': 7, '4': 1, '5': 3, '10': 'defaultFee'},
  ],
};

/// Descriptor for `GetWalletInfoResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getWalletInfoResponseDescriptor = $convert.base64Decode(
    'ChVHZXRXYWxsZXRJbmZvUmVzcG9uc2USHwoLd2FsbGV0X25hbWUYASABKAlSCndhbGxldE5hbW'
    'USGAoHdmVyc2lvbhgCIAEoBVIHdmVyc2lvbhIYCgduZXR3b3JrGAMgASgJUgduZXR3b3JrEhwK'
    'CWVuY3J5cHRlZBgEIAEoCFIJZW5jcnlwdGVkEhIKBHV1aWQYBSABKAlSBHV1aWQSHQoKY3JlYX'
    'RlZF9hdBgGIAEoA1IJY3JlYXRlZEF0Eh8KC2RlZmF1bHRfZmVlGAcgASgDUgpkZWZhdWx0RmVl');

@$core.Deprecated('Use listAddressRequestDescriptor instead')
const ListAddressRequest$json = {
  '1': 'ListAddressRequest',
  '2': [
    {'1': 'wallet_name', '3': 1, '4': 1, '5': 9, '10': 'walletName'},
  ],
};

/// Descriptor for `ListAddressRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listAddressRequestDescriptor = $convert.base64Decode(
    'ChJMaXN0QWRkcmVzc1JlcXVlc3QSHwoLd2FsbGV0X25hbWUYASABKAlSCndhbGxldE5hbWU=');

@$core.Deprecated('Use listAddressResponseDescriptor instead')
const ListAddressResponse$json = {
  '1': 'ListAddressResponse',
  '2': [
    {'1': 'wallet_name', '3': 1, '4': 1, '5': 9, '10': 'walletName'},
    {'1': 'data', '3': 2, '4': 3, '5': 11, '6': '.pactus.AddressInfo', '10': 'data'},
  ],
};

/// Descriptor for `ListAddressResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listAddressResponseDescriptor = $convert.base64Decode(
    'ChNMaXN0QWRkcmVzc1Jlc3BvbnNlEh8KC3dhbGxldF9uYW1lGAEgASgJUgp3YWxsZXROYW1lEi'
    'cKBGRhdGEYAiADKAsyEy5wYWN0dXMuQWRkcmVzc0luZm9SBGRhdGE=');

const $core.Map<$core.String, $core.dynamic> WalletServiceBase$json = {
  '1': 'Wallet',
  '2': [
    {'1': 'CreateWallet', '2': '.pactus.CreateWalletRequest', '3': '.pactus.CreateWalletResponse'},
    {'1': 'RestoreWallet', '2': '.pactus.RestoreWalletRequest', '3': '.pactus.RestoreWalletResponse'},
    {'1': 'LoadWallet', '2': '.pactus.LoadWalletRequest', '3': '.pactus.LoadWalletResponse'},
    {'1': 'UnloadWallet', '2': '.pactus.UnloadWalletRequest', '3': '.pactus.UnloadWalletResponse'},
    {'1': 'GetTotalBalance', '2': '.pactus.GetTotalBalanceRequest', '3': '.pactus.GetTotalBalanceResponse'},
    {'1': 'SignRawTransaction', '2': '.pactus.SignRawTransactionRequest', '3': '.pactus.SignRawTransactionResponse'},
    {'1': 'GetValidatorAddress', '2': '.pactus.GetValidatorAddressRequest', '3': '.pactus.GetValidatorAddressResponse'},
    {'1': 'GetNewAddress', '2': '.pactus.GetNewAddressRequest', '3': '.pactus.GetNewAddressResponse'},
    {'1': 'GetAddressHistory', '2': '.pactus.GetAddressHistoryRequest', '3': '.pactus.GetAddressHistoryResponse'},
    {'1': 'SignMessage', '2': '.pactus.SignMessageRequest', '3': '.pactus.SignMessageResponse'},
    {'1': 'GetTotalStake', '2': '.pactus.GetTotalStakeRequest', '3': '.pactus.GetTotalStakeResponse'},
    {'1': 'GetAddressInfo', '2': '.pactus.GetAddressInfoRequest', '3': '.pactus.GetAddressInfoResponse'},
    {'1': 'SetAddressLabel', '2': '.pactus.SetAddressLabelRequest', '3': '.pactus.SetAddressLabelResponse'},
    {'1': 'ListWallet', '2': '.pactus.ListWalletRequest', '3': '.pactus.ListWalletResponse'},
    {'1': 'GetWalletInfo', '2': '.pactus.GetWalletInfoRequest', '3': '.pactus.GetWalletInfoResponse'},
    {'1': 'ListAddress', '2': '.pactus.ListAddressRequest', '3': '.pactus.ListAddressResponse'},
  ],
};

@$core.Deprecated('Use walletServiceDescriptor instead')
const $core.Map<$core.String, $core.Map<$core.String, $core.dynamic>> WalletServiceBase$messageJson = {
  '.pactus.CreateWalletRequest': CreateWalletRequest$json,
  '.pactus.CreateWalletResponse': CreateWalletResponse$json,
  '.pactus.RestoreWalletRequest': RestoreWalletRequest$json,
  '.pactus.RestoreWalletResponse': RestoreWalletResponse$json,
  '.pactus.LoadWalletRequest': LoadWalletRequest$json,
  '.pactus.LoadWalletResponse': LoadWalletResponse$json,
  '.pactus.UnloadWalletRequest': UnloadWalletRequest$json,
  '.pactus.UnloadWalletResponse': UnloadWalletResponse$json,
  '.pactus.GetTotalBalanceRequest': GetTotalBalanceRequest$json,
  '.pactus.GetTotalBalanceResponse': GetTotalBalanceResponse$json,
  '.pactus.SignRawTransactionRequest': SignRawTransactionRequest$json,
  '.pactus.SignRawTransactionResponse': SignRawTransactionResponse$json,
  '.pactus.GetValidatorAddressRequest': GetValidatorAddressRequest$json,
  '.pactus.GetValidatorAddressResponse': GetValidatorAddressResponse$json,
  '.pactus.GetNewAddressRequest': GetNewAddressRequest$json,
  '.pactus.GetNewAddressResponse': GetNewAddressResponse$json,
  '.pactus.AddressInfo': AddressInfo$json,
  '.pactus.GetAddressHistoryRequest': GetAddressHistoryRequest$json,
  '.pactus.GetAddressHistoryResponse': GetAddressHistoryResponse$json,
  '.pactus.HistoryInfo': HistoryInfo$json,
  '.pactus.SignMessageRequest': SignMessageRequest$json,
  '.pactus.SignMessageResponse': SignMessageResponse$json,
  '.pactus.GetTotalStakeRequest': GetTotalStakeRequest$json,
  '.pactus.GetTotalStakeResponse': GetTotalStakeResponse$json,
  '.pactus.GetAddressInfoRequest': GetAddressInfoRequest$json,
  '.pactus.GetAddressInfoResponse': GetAddressInfoResponse$json,
  '.pactus.SetAddressLabelRequest': SetAddressLabelRequest$json,
  '.pactus.SetAddressLabelResponse': SetAddressLabelResponse$json,
  '.pactus.ListWalletRequest': ListWalletRequest$json,
  '.pactus.ListWalletResponse': ListWalletResponse$json,
  '.pactus.GetWalletInfoRequest': GetWalletInfoRequest$json,
  '.pactus.GetWalletInfoResponse': GetWalletInfoResponse$json,
  '.pactus.ListAddressRequest': ListAddressRequest$json,
  '.pactus.ListAddressResponse': ListAddressResponse$json,
};

/// Descriptor for `Wallet`. Decode as a `google.protobuf.ServiceDescriptorProto`.
final $typed_data.Uint8List walletServiceDescriptor = $convert.base64Decode(
    'CgZXYWxsZXQSSQoMQ3JlYXRlV2FsbGV0EhsucGFjdHVzLkNyZWF0ZVdhbGxldFJlcXVlc3QaHC'
    '5wYWN0dXMuQ3JlYXRlV2FsbGV0UmVzcG9uc2USTAoNUmVzdG9yZVdhbGxldBIcLnBhY3R1cy5S'
    'ZXN0b3JlV2FsbGV0UmVxdWVzdBodLnBhY3R1cy5SZXN0b3JlV2FsbGV0UmVzcG9uc2USQwoKTG'
    '9hZFdhbGxldBIZLnBhY3R1cy5Mb2FkV2FsbGV0UmVxdWVzdBoaLnBhY3R1cy5Mb2FkV2FsbGV0'
    'UmVzcG9uc2USSQoMVW5sb2FkV2FsbGV0EhsucGFjdHVzLlVubG9hZFdhbGxldFJlcXVlc3QaHC'
    '5wYWN0dXMuVW5sb2FkV2FsbGV0UmVzcG9uc2USUgoPR2V0VG90YWxCYWxhbmNlEh4ucGFjdHVz'
    'LkdldFRvdGFsQmFsYW5jZVJlcXVlc3QaHy5wYWN0dXMuR2V0VG90YWxCYWxhbmNlUmVzcG9uc2'
    'USWwoSU2lnblJhd1RyYW5zYWN0aW9uEiEucGFjdHVzLlNpZ25SYXdUcmFuc2FjdGlvblJlcXVl'
    'c3QaIi5wYWN0dXMuU2lnblJhd1RyYW5zYWN0aW9uUmVzcG9uc2USXgoTR2V0VmFsaWRhdG9yQW'
    'RkcmVzcxIiLnBhY3R1cy5HZXRWYWxpZGF0b3JBZGRyZXNzUmVxdWVzdBojLnBhY3R1cy5HZXRW'
    'YWxpZGF0b3JBZGRyZXNzUmVzcG9uc2USTAoNR2V0TmV3QWRkcmVzcxIcLnBhY3R1cy5HZXROZX'
    'dBZGRyZXNzUmVxdWVzdBodLnBhY3R1cy5HZXROZXdBZGRyZXNzUmVzcG9uc2USWAoRR2V0QWRk'
    'cmVzc0hpc3RvcnkSIC5wYWN0dXMuR2V0QWRkcmVzc0hpc3RvcnlSZXF1ZXN0GiEucGFjdHVzLk'
    'dldEFkZHJlc3NIaXN0b3J5UmVzcG9uc2USRgoLU2lnbk1lc3NhZ2USGi5wYWN0dXMuU2lnbk1l'
    'c3NhZ2VSZXF1ZXN0GhsucGFjdHVzLlNpZ25NZXNzYWdlUmVzcG9uc2USTAoNR2V0VG90YWxTdG'
    'FrZRIcLnBhY3R1cy5HZXRUb3RhbFN0YWtlUmVxdWVzdBodLnBhY3R1cy5HZXRUb3RhbFN0YWtl'
    'UmVzcG9uc2USTwoOR2V0QWRkcmVzc0luZm8SHS5wYWN0dXMuR2V0QWRkcmVzc0luZm9SZXF1ZX'
    'N0Gh4ucGFjdHVzLkdldEFkZHJlc3NJbmZvUmVzcG9uc2USUgoPU2V0QWRkcmVzc0xhYmVsEh4u'
    'cGFjdHVzLlNldEFkZHJlc3NMYWJlbFJlcXVlc3QaHy5wYWN0dXMuU2V0QWRkcmVzc0xhYmVsUm'
    'VzcG9uc2USQwoKTGlzdFdhbGxldBIZLnBhY3R1cy5MaXN0V2FsbGV0UmVxdWVzdBoaLnBhY3R1'
    'cy5MaXN0V2FsbGV0UmVzcG9uc2USTAoNR2V0V2FsbGV0SW5mbxIcLnBhY3R1cy5HZXRXYWxsZX'
    'RJbmZvUmVxdWVzdBodLnBhY3R1cy5HZXRXYWxsZXRJbmZvUmVzcG9uc2USRgoLTGlzdEFkZHJl'
    'c3MSGi5wYWN0dXMuTGlzdEFkZHJlc3NSZXF1ZXN0GhsucGFjdHVzLkxpc3RBZGRyZXNzUmVzcG'
    '9uc2U=');

