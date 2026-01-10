// This is a generated file - do not edit.
//
// Generated from wallet.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports
// ignore_for_file: unused_import

import 'dart:convert' as $convert;
import 'dart:core' as $core;
import 'dart:typed_data' as $typed_data;

import 'transaction.pbjson.dart' as $0;

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

@$core.Deprecated('Use txDirectionDescriptor instead')
const TxDirection$json = {
  '1': 'TxDirection',
  '2': [
    {'1': 'TX_DIRECTION_ANY', '2': 0},
    {'1': 'TX_DIRECTION_INCOMING', '2': 1},
    {'1': 'TX_DIRECTION_OUTGOING', '2': 2},
  ],
};

/// Descriptor for `TxDirection`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List txDirectionDescriptor = $convert.base64Decode(
    'CgtUeERpcmVjdGlvbhIUChBUWF9ESVJFQ1RJT05fQU5ZEAASGQoVVFhfRElSRUNUSU9OX0lOQ0'
    '9NSU5HEAESGQoVVFhfRElSRUNUSU9OX09VVEdPSU5HEAI=');

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

@$core.Deprecated('Use getNewAddressRequestDescriptor instead')
const GetNewAddressRequest$json = {
  '1': 'GetNewAddressRequest',
  '2': [
    {'1': 'wallet_name', '3': 1, '4': 1, '5': 9, '10': 'walletName'},
    {
      '1': 'address_type',
      '3': 2,
      '4': 1,
      '5': 14,
      '6': '.pactus.AddressType',
      '10': 'addressType'
    },
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
    {
      '1': 'address_info',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.pactus.AddressInfo',
      '10': 'addressInfo'
    },
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
    {'1': 'password', '3': 2, '4': 1, '5': 9, '10': 'password'},
  ],
};

/// Descriptor for `CreateWalletRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createWalletRequestDescriptor = $convert.base64Decode(
    'ChNDcmVhdGVXYWxsZXRSZXF1ZXN0Eh8KC3dhbGxldF9uYW1lGAEgASgJUgp3YWxsZXROYW1lEh'
    'oKCHBhc3N3b3JkGAIgASgJUghwYXNzd29yZA==');

@$core.Deprecated('Use createWalletResponseDescriptor instead')
const CreateWalletResponse$json = {
  '1': 'CreateWalletResponse',
  '2': [
    {'1': 'wallet_name', '3': 1, '4': 1, '5': 9, '10': 'walletName'},
    {'1': 'mnemonic', '3': 2, '4': 1, '5': 9, '10': 'mnemonic'},
  ],
};

/// Descriptor for `CreateWalletResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createWalletResponseDescriptor = $convert.base64Decode(
    'ChRDcmVhdGVXYWxsZXRSZXNwb25zZRIfCgt3YWxsZXRfbmFtZRgBIAEoCVIKd2FsbGV0TmFtZR'
    'IaCghtbmVtb25pYxgCIAEoCVIIbW5lbW9uaWM=');

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
final $typed_data.Uint8List getValidatorAddressRequestDescriptor =
    $convert.base64Decode(
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
final $typed_data.Uint8List getValidatorAddressResponseDescriptor =
    $convert.base64Decode(
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
    {
      '1': 'signed_raw_transaction',
      '3': 2,
      '4': 1,
      '5': 9,
      '10': 'signedRawTransaction'
    },
  ],
};

/// Descriptor for `SignRawTransactionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List signRawTransactionResponseDescriptor =
    $convert.base64Decode(
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
final $typed_data.Uint8List getTotalBalanceRequestDescriptor =
    $convert.base64Decode(
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
final $typed_data.Uint8List getTotalBalanceResponseDescriptor =
    $convert.base64Decode(
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
final $typed_data.Uint8List signMessageResponseDescriptor =
    $convert.base64Decode(
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
    {
      '1': 'address_info',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.pactus.AddressInfo',
      '10': 'addressInfo'
    },
  ],
};

/// Descriptor for `GetAddressInfoResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getAddressInfoResponseDescriptor = $convert.base64Decode(
    'ChZHZXRBZGRyZXNzSW5mb1Jlc3BvbnNlEh8KC3dhbGxldF9uYW1lGAEgASgJUgp3YWxsZXROYW'
    '1lEjYKDGFkZHJlc3NfaW5mbxgCIAEoCzITLnBhY3R1cy5BZGRyZXNzSW5mb1ILYWRkcmVzc0lu'
    'Zm8=');

@$core.Deprecated('Use setAddressLabelRequestDescriptor instead')
const SetAddressLabelRequest$json = {
  '1': 'SetAddressLabelRequest',
  '2': [
    {'1': 'wallet_name', '3': 1, '4': 1, '5': 9, '10': 'walletName'},
    {'1': 'password', '3': 2, '4': 1, '5': 9, '10': 'password'},
    {'1': 'address', '3': 3, '4': 1, '5': 9, '10': 'address'},
    {'1': 'label', '3': 4, '4': 1, '5': 9, '10': 'label'},
  ],
};

/// Descriptor for `SetAddressLabelRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List setAddressLabelRequestDescriptor = $convert.base64Decode(
    'ChZTZXRBZGRyZXNzTGFiZWxSZXF1ZXN0Eh8KC3dhbGxldF9uYW1lGAEgASgJUgp3YWxsZXROYW'
    '1lEhoKCHBhc3N3b3JkGAIgASgJUghwYXNzd29yZBIYCgdhZGRyZXNzGAMgASgJUgdhZGRyZXNz'
    'EhQKBWxhYmVsGAQgASgJUgVsYWJlbA==');

@$core.Deprecated('Use setAddressLabelResponseDescriptor instead')
const SetAddressLabelResponse$json = {
  '1': 'SetAddressLabelResponse',
  '2': [
    {'1': 'wallet_name', '3': 1, '4': 1, '5': 9, '10': 'walletName'},
    {'1': 'address', '3': 2, '4': 1, '5': 9, '10': 'address'},
    {'1': 'label', '3': 3, '4': 1, '5': 9, '10': 'label'},
  ],
};

/// Descriptor for `SetAddressLabelResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List setAddressLabelResponseDescriptor = $convert.base64Decode(
    'ChdTZXRBZGRyZXNzTGFiZWxSZXNwb25zZRIfCgt3YWxsZXRfbmFtZRgBIAEoCVIKd2FsbGV0Tm'
    'FtZRIYCgdhZGRyZXNzGAIgASgJUgdhZGRyZXNzEhQKBWxhYmVsGAMgASgJUgVsYWJlbA==');

@$core.Deprecated('Use listWalletsRequestDescriptor instead')
const ListWalletsRequest$json = {
  '1': 'ListWalletsRequest',
};

/// Descriptor for `ListWalletsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listWalletsRequestDescriptor =
    $convert.base64Decode('ChJMaXN0V2FsbGV0c1JlcXVlc3Q=');

@$core.Deprecated('Use listWalletsResponseDescriptor instead')
const ListWalletsResponse$json = {
  '1': 'ListWalletsResponse',
  '2': [
    {'1': 'wallets', '3': 1, '4': 3, '5': 9, '10': 'wallets'},
  ],
};

/// Descriptor for `ListWalletsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listWalletsResponseDescriptor =
    $convert.base64Decode(
        'ChNMaXN0V2FsbGV0c1Jlc3BvbnNlEhgKB3dhbGxldHMYASADKAlSB3dhbGxldHM=');

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
    {'1': 'driver', '3': 8, '4': 1, '5': 9, '10': 'driver'},
    {'1': 'path', '3': 9, '4': 1, '5': 9, '10': 'path'},
  ],
};

/// Descriptor for `GetWalletInfoResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getWalletInfoResponseDescriptor = $convert.base64Decode(
    'ChVHZXRXYWxsZXRJbmZvUmVzcG9uc2USHwoLd2FsbGV0X25hbWUYASABKAlSCndhbGxldE5hbW'
    'USGAoHdmVyc2lvbhgCIAEoBVIHdmVyc2lvbhIYCgduZXR3b3JrGAMgASgJUgduZXR3b3JrEhwK'
    'CWVuY3J5cHRlZBgEIAEoCFIJZW5jcnlwdGVkEhIKBHV1aWQYBSABKAlSBHV1aWQSHQoKY3JlYX'
    'RlZF9hdBgGIAEoA1IJY3JlYXRlZEF0Eh8KC2RlZmF1bHRfZmVlGAcgASgDUgpkZWZhdWx0RmVl'
    'EhYKBmRyaXZlchgIIAEoCVIGZHJpdmVyEhIKBHBhdGgYCSABKAlSBHBhdGg=');

@$core.Deprecated('Use listAddressesRequestDescriptor instead')
const ListAddressesRequest$json = {
  '1': 'ListAddressesRequest',
  '2': [
    {'1': 'wallet_name', '3': 1, '4': 1, '5': 9, '10': 'walletName'},
    {
      '1': 'address_types',
      '3': 2,
      '4': 3,
      '5': 14,
      '6': '.pactus.AddressType',
      '10': 'addressTypes'
    },
  ],
};

/// Descriptor for `ListAddressesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listAddressesRequestDescriptor = $convert.base64Decode(
    'ChRMaXN0QWRkcmVzc2VzUmVxdWVzdBIfCgt3YWxsZXRfbmFtZRgBIAEoCVIKd2FsbGV0TmFtZR'
    'I4Cg1hZGRyZXNzX3R5cGVzGAIgAygOMhMucGFjdHVzLkFkZHJlc3NUeXBlUgxhZGRyZXNzVHlw'
    'ZXM=');

@$core.Deprecated('Use listAddressesResponseDescriptor instead')
const ListAddressesResponse$json = {
  '1': 'ListAddressesResponse',
  '2': [
    {'1': 'wallet_name', '3': 1, '4': 1, '5': 9, '10': 'walletName'},
    {
      '1': 'data',
      '3': 2,
      '4': 3,
      '5': 11,
      '6': '.pactus.AddressInfo',
      '10': 'data'
    },
  ],
};

/// Descriptor for `ListAddressesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listAddressesResponseDescriptor = $convert.base64Decode(
    'ChVMaXN0QWRkcmVzc2VzUmVzcG9uc2USHwoLd2FsbGV0X25hbWUYASABKAlSCndhbGxldE5hbW'
    'USJwoEZGF0YRgCIAMoCzITLnBhY3R1cy5BZGRyZXNzSW5mb1IEZGF0YQ==');

@$core.Deprecated('Use updatePasswordRequestDescriptor instead')
const UpdatePasswordRequest$json = {
  '1': 'UpdatePasswordRequest',
  '2': [
    {'1': 'wallet_name', '3': 1, '4': 1, '5': 9, '10': 'walletName'},
    {'1': 'old_password', '3': 2, '4': 1, '5': 9, '10': 'oldPassword'},
    {'1': 'new_password', '3': 3, '4': 1, '5': 9, '10': 'newPassword'},
  ],
};

/// Descriptor for `UpdatePasswordRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updatePasswordRequestDescriptor = $convert.base64Decode(
    'ChVVcGRhdGVQYXNzd29yZFJlcXVlc3QSHwoLd2FsbGV0X25hbWUYASABKAlSCndhbGxldE5hbW'
    'USIQoMb2xkX3Bhc3N3b3JkGAIgASgJUgtvbGRQYXNzd29yZBIhCgxuZXdfcGFzc3dvcmQYAyAB'
    'KAlSC25ld1Bhc3N3b3Jk');

@$core.Deprecated('Use updatePasswordResponseDescriptor instead')
const UpdatePasswordResponse$json = {
  '1': 'UpdatePasswordResponse',
  '2': [
    {'1': 'wallet_name', '3': 1, '4': 1, '5': 9, '10': 'walletName'},
  ],
};

/// Descriptor for `UpdatePasswordResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updatePasswordResponseDescriptor =
    $convert.base64Decode(
        'ChZVcGRhdGVQYXNzd29yZFJlc3BvbnNlEh8KC3dhbGxldF9uYW1lGAEgASgJUgp3YWxsZXROYW'
        '1l');

@$core.Deprecated('Use listTransactionsRequestDescriptor instead')
const ListTransactionsRequest$json = {
  '1': 'ListTransactionsRequest',
  '2': [
    {'1': 'wallet_name', '3': 1, '4': 1, '5': 9, '10': 'walletName'},
    {'1': 'address', '3': 2, '4': 1, '5': 9, '10': 'address'},
    {
      '1': 'direction',
      '3': 3,
      '4': 1,
      '5': 14,
      '6': '.pactus.TxDirection',
      '10': 'direction'
    },
    {'1': 'count', '3': 4, '4': 1, '5': 5, '10': 'count'},
    {'1': 'skip', '3': 5, '4': 1, '5': 5, '10': 'skip'},
  ],
};

/// Descriptor for `ListTransactionsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTransactionsRequestDescriptor = $convert.base64Decode(
    'ChdMaXN0VHJhbnNhY3Rpb25zUmVxdWVzdBIfCgt3YWxsZXRfbmFtZRgBIAEoCVIKd2FsbGV0Tm'
    'FtZRIYCgdhZGRyZXNzGAIgASgJUgdhZGRyZXNzEjEKCWRpcmVjdGlvbhgDIAEoDjITLnBhY3R1'
    'cy5UeERpcmVjdGlvblIJZGlyZWN0aW9uEhQKBWNvdW50GAQgASgFUgVjb3VudBISCgRza2lwGA'
    'UgASgFUgRza2lw');

@$core.Deprecated('Use listTransactionsResponseDescriptor instead')
const ListTransactionsResponse$json = {
  '1': 'ListTransactionsResponse',
  '2': [
    {'1': 'wallet_name', '3': 1, '4': 1, '5': 9, '10': 'walletName'},
    {
      '1': 'txs',
      '3': 2,
      '4': 3,
      '5': 11,
      '6': '.pactus.TransactionInfo',
      '10': 'txs'
    },
  ],
};

/// Descriptor for `ListTransactionsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTransactionsResponseDescriptor =
    $convert.base64Decode(
        'ChhMaXN0VHJhbnNhY3Rpb25zUmVzcG9uc2USHwoLd2FsbGV0X25hbWUYASABKAlSCndhbGxldE'
        '5hbWUSKQoDdHhzGAIgAygLMhcucGFjdHVzLlRyYW5zYWN0aW9uSW5mb1IDdHhz');

const $core.Map<$core.String, $core.dynamic> WalletServiceBase$json = {
  '1': 'Wallet',
  '2': [
    {
      '1': 'CreateWallet',
      '2': '.pactus.CreateWalletRequest',
      '3': '.pactus.CreateWalletResponse'
    },
    {
      '1': 'RestoreWallet',
      '2': '.pactus.RestoreWalletRequest',
      '3': '.pactus.RestoreWalletResponse'
    },
    {
      '1': 'LoadWallet',
      '2': '.pactus.LoadWalletRequest',
      '3': '.pactus.LoadWalletResponse'
    },
    {
      '1': 'UnloadWallet',
      '2': '.pactus.UnloadWalletRequest',
      '3': '.pactus.UnloadWalletResponse'
    },
    {
      '1': 'ListWallets',
      '2': '.pactus.ListWalletsRequest',
      '3': '.pactus.ListWalletsResponse'
    },
    {
      '1': 'GetWalletInfo',
      '2': '.pactus.GetWalletInfoRequest',
      '3': '.pactus.GetWalletInfoResponse'
    },
    {
      '1': 'UpdatePassword',
      '2': '.pactus.UpdatePasswordRequest',
      '3': '.pactus.UpdatePasswordResponse'
    },
    {
      '1': 'GetTotalBalance',
      '2': '.pactus.GetTotalBalanceRequest',
      '3': '.pactus.GetTotalBalanceResponse'
    },
    {
      '1': 'GetTotalStake',
      '2': '.pactus.GetTotalStakeRequest',
      '3': '.pactus.GetTotalStakeResponse'
    },
    {
      '1': 'GetValidatorAddress',
      '2': '.pactus.GetValidatorAddressRequest',
      '3': '.pactus.GetValidatorAddressResponse'
    },
    {
      '1': 'GetAddressInfo',
      '2': '.pactus.GetAddressInfoRequest',
      '3': '.pactus.GetAddressInfoResponse'
    },
    {
      '1': 'SetAddressLabel',
      '2': '.pactus.SetAddressLabelRequest',
      '3': '.pactus.SetAddressLabelResponse'
    },
    {
      '1': 'GetNewAddress',
      '2': '.pactus.GetNewAddressRequest',
      '3': '.pactus.GetNewAddressResponse'
    },
    {
      '1': 'ListAddresses',
      '2': '.pactus.ListAddressesRequest',
      '3': '.pactus.ListAddressesResponse'
    },
    {
      '1': 'SignMessage',
      '2': '.pactus.SignMessageRequest',
      '3': '.pactus.SignMessageResponse'
    },
    {
      '1': 'SignRawTransaction',
      '2': '.pactus.SignRawTransactionRequest',
      '3': '.pactus.SignRawTransactionResponse'
    },
    {
      '1': 'ListTransactions',
      '2': '.pactus.ListTransactionsRequest',
      '3': '.pactus.ListTransactionsResponse'
    },
  ],
};

@$core.Deprecated('Use walletServiceDescriptor instead')
const $core.Map<$core.String, $core.Map<$core.String, $core.dynamic>>
    WalletServiceBase$messageJson = {
  '.pactus.CreateWalletRequest': CreateWalletRequest$json,
  '.pactus.CreateWalletResponse': CreateWalletResponse$json,
  '.pactus.RestoreWalletRequest': RestoreWalletRequest$json,
  '.pactus.RestoreWalletResponse': RestoreWalletResponse$json,
  '.pactus.LoadWalletRequest': LoadWalletRequest$json,
  '.pactus.LoadWalletResponse': LoadWalletResponse$json,
  '.pactus.UnloadWalletRequest': UnloadWalletRequest$json,
  '.pactus.UnloadWalletResponse': UnloadWalletResponse$json,
  '.pactus.ListWalletsRequest': ListWalletsRequest$json,
  '.pactus.ListWalletsResponse': ListWalletsResponse$json,
  '.pactus.GetWalletInfoRequest': GetWalletInfoRequest$json,
  '.pactus.GetWalletInfoResponse': GetWalletInfoResponse$json,
  '.pactus.UpdatePasswordRequest': UpdatePasswordRequest$json,
  '.pactus.UpdatePasswordResponse': UpdatePasswordResponse$json,
  '.pactus.GetTotalBalanceRequest': GetTotalBalanceRequest$json,
  '.pactus.GetTotalBalanceResponse': GetTotalBalanceResponse$json,
  '.pactus.GetTotalStakeRequest': GetTotalStakeRequest$json,
  '.pactus.GetTotalStakeResponse': GetTotalStakeResponse$json,
  '.pactus.GetValidatorAddressRequest': GetValidatorAddressRequest$json,
  '.pactus.GetValidatorAddressResponse': GetValidatorAddressResponse$json,
  '.pactus.GetAddressInfoRequest': GetAddressInfoRequest$json,
  '.pactus.GetAddressInfoResponse': GetAddressInfoResponse$json,
  '.pactus.AddressInfo': AddressInfo$json,
  '.pactus.SetAddressLabelRequest': SetAddressLabelRequest$json,
  '.pactus.SetAddressLabelResponse': SetAddressLabelResponse$json,
  '.pactus.GetNewAddressRequest': GetNewAddressRequest$json,
  '.pactus.GetNewAddressResponse': GetNewAddressResponse$json,
  '.pactus.ListAddressesRequest': ListAddressesRequest$json,
  '.pactus.ListAddressesResponse': ListAddressesResponse$json,
  '.pactus.SignMessageRequest': SignMessageRequest$json,
  '.pactus.SignMessageResponse': SignMessageResponse$json,
  '.pactus.SignRawTransactionRequest': SignRawTransactionRequest$json,
  '.pactus.SignRawTransactionResponse': SignRawTransactionResponse$json,
  '.pactus.ListTransactionsRequest': ListTransactionsRequest$json,
  '.pactus.ListTransactionsResponse': ListTransactionsResponse$json,
  '.pactus.TransactionInfo': $0.TransactionInfo$json,
  '.pactus.PayloadTransfer': $0.PayloadTransfer$json,
  '.pactus.PayloadBond': $0.PayloadBond$json,
  '.pactus.PayloadSortition': $0.PayloadSortition$json,
  '.pactus.PayloadUnbond': $0.PayloadUnbond$json,
  '.pactus.PayloadWithdraw': $0.PayloadWithdraw$json,
  '.pactus.PayloadBatchTransfer': $0.PayloadBatchTransfer$json,
  '.pactus.Recipient': $0.Recipient$json,
};

/// Descriptor for `Wallet`. Decode as a `google.protobuf.ServiceDescriptorProto`.
final $typed_data.Uint8List walletServiceDescriptor = $convert.base64Decode(
    'CgZXYWxsZXQSSQoMQ3JlYXRlV2FsbGV0EhsucGFjdHVzLkNyZWF0ZVdhbGxldFJlcXVlc3QaHC'
    '5wYWN0dXMuQ3JlYXRlV2FsbGV0UmVzcG9uc2USTAoNUmVzdG9yZVdhbGxldBIcLnBhY3R1cy5S'
    'ZXN0b3JlV2FsbGV0UmVxdWVzdBodLnBhY3R1cy5SZXN0b3JlV2FsbGV0UmVzcG9uc2USQwoKTG'
    '9hZFdhbGxldBIZLnBhY3R1cy5Mb2FkV2FsbGV0UmVxdWVzdBoaLnBhY3R1cy5Mb2FkV2FsbGV0'
    'UmVzcG9uc2USSQoMVW5sb2FkV2FsbGV0EhsucGFjdHVzLlVubG9hZFdhbGxldFJlcXVlc3QaHC'
    '5wYWN0dXMuVW5sb2FkV2FsbGV0UmVzcG9uc2USRgoLTGlzdFdhbGxldHMSGi5wYWN0dXMuTGlz'
    'dFdhbGxldHNSZXF1ZXN0GhsucGFjdHVzLkxpc3RXYWxsZXRzUmVzcG9uc2USTAoNR2V0V2FsbG'
    'V0SW5mbxIcLnBhY3R1cy5HZXRXYWxsZXRJbmZvUmVxdWVzdBodLnBhY3R1cy5HZXRXYWxsZXRJ'
    'bmZvUmVzcG9uc2USTwoOVXBkYXRlUGFzc3dvcmQSHS5wYWN0dXMuVXBkYXRlUGFzc3dvcmRSZX'
    'F1ZXN0Gh4ucGFjdHVzLlVwZGF0ZVBhc3N3b3JkUmVzcG9uc2USUgoPR2V0VG90YWxCYWxhbmNl'
    'Eh4ucGFjdHVzLkdldFRvdGFsQmFsYW5jZVJlcXVlc3QaHy5wYWN0dXMuR2V0VG90YWxCYWxhbm'
    'NlUmVzcG9uc2USTAoNR2V0VG90YWxTdGFrZRIcLnBhY3R1cy5HZXRUb3RhbFN0YWtlUmVxdWVz'
    'dBodLnBhY3R1cy5HZXRUb3RhbFN0YWtlUmVzcG9uc2USXgoTR2V0VmFsaWRhdG9yQWRkcmVzcx'
    'IiLnBhY3R1cy5HZXRWYWxpZGF0b3JBZGRyZXNzUmVxdWVzdBojLnBhY3R1cy5HZXRWYWxpZGF0'
    'b3JBZGRyZXNzUmVzcG9uc2USTwoOR2V0QWRkcmVzc0luZm8SHS5wYWN0dXMuR2V0QWRkcmVzc0'
    'luZm9SZXF1ZXN0Gh4ucGFjdHVzLkdldEFkZHJlc3NJbmZvUmVzcG9uc2USUgoPU2V0QWRkcmVz'
    'c0xhYmVsEh4ucGFjdHVzLlNldEFkZHJlc3NMYWJlbFJlcXVlc3QaHy5wYWN0dXMuU2V0QWRkcm'
    'Vzc0xhYmVsUmVzcG9uc2USTAoNR2V0TmV3QWRkcmVzcxIcLnBhY3R1cy5HZXROZXdBZGRyZXNz'
    'UmVxdWVzdBodLnBhY3R1cy5HZXROZXdBZGRyZXNzUmVzcG9uc2USTAoNTGlzdEFkZHJlc3Nlcx'
    'IcLnBhY3R1cy5MaXN0QWRkcmVzc2VzUmVxdWVzdBodLnBhY3R1cy5MaXN0QWRkcmVzc2VzUmVz'
    'cG9uc2USRgoLU2lnbk1lc3NhZ2USGi5wYWN0dXMuU2lnbk1lc3NhZ2VSZXF1ZXN0GhsucGFjdH'
    'VzLlNpZ25NZXNzYWdlUmVzcG9uc2USWwoSU2lnblJhd1RyYW5zYWN0aW9uEiEucGFjdHVzLlNp'
    'Z25SYXdUcmFuc2FjdGlvblJlcXVlc3QaIi5wYWN0dXMuU2lnblJhd1RyYW5zYWN0aW9uUmVzcG'
    '9uc2USVQoQTGlzdFRyYW5zYWN0aW9ucxIfLnBhY3R1cy5MaXN0VHJhbnNhY3Rpb25zUmVxdWVz'
    'dBogLnBhY3R1cy5MaXN0VHJhbnNhY3Rpb25zUmVzcG9uc2U=');
