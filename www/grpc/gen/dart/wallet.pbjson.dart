///
//  Generated code. Do not modify.
//  source: wallet.proto
//
// @dart = 2.12
// ignore_for_file: annotate_overrides,camel_case_types,constant_identifier_names,deprecated_member_use_from_same_package,directives_ordering,library_prefixes,non_constant_identifier_names,prefer_final_fields,return_of_invalid_type,unnecessary_const,unnecessary_import,unnecessary_this,unused_import,unused_shown_name

import 'dart:core' as $core;
import 'dart:convert' as $convert;
import 'dart:typed_data' as $typed_data;
@$core.Deprecated('Use addressTypeDescriptor instead')
const AddressType$json = const {
  '1': 'AddressType',
  '2': const [
    const {'1': 'ADDRESS_TYPE_TREASURY', '2': 0},
    const {'1': 'ADDRESS_TYPE_VALIDATOR', '2': 1},
    const {'1': 'ADDRESS_TYPE_BLS_ACCOUNT', '2': 2},
    const {'1': 'ADDRESS_TYPE_ED25519_ACCOUNT', '2': 3},
  ],
};

/// Descriptor for `AddressType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List addressTypeDescriptor = $convert.base64Decode('CgtBZGRyZXNzVHlwZRIZChVBRERSRVNTX1RZUEVfVFJFQVNVUlkQABIaChZBRERSRVNTX1RZUEVfVkFMSURBVE9SEAESHAoYQUREUkVTU19UWVBFX0JMU19BQ0NPVU5UEAISIAocQUREUkVTU19UWVBFX0VEMjU1MTlfQUNDT1VOVBAD');
@$core.Deprecated('Use addressInfoDescriptor instead')
const AddressInfo$json = const {
  '1': 'AddressInfo',
  '2': const [
    const {'1': 'address', '3': 1, '4': 1, '5': 9, '10': 'address'},
    const {'1': 'public_key', '3': 2, '4': 1, '5': 9, '10': 'publicKey'},
    const {'1': 'label', '3': 3, '4': 1, '5': 9, '10': 'label'},
    const {'1': 'path', '3': 4, '4': 1, '5': 9, '10': 'path'},
  ],
};

/// Descriptor for `AddressInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List addressInfoDescriptor = $convert.base64Decode('CgtBZGRyZXNzSW5mbxIYCgdhZGRyZXNzGAEgASgJUgdhZGRyZXNzEh0KCnB1YmxpY19rZXkYAiABKAlSCXB1YmxpY0tleRIUCgVsYWJlbBgDIAEoCVIFbGFiZWwSEgoEcGF0aBgEIAEoCVIEcGF0aA==');
@$core.Deprecated('Use historyInfoDescriptor instead')
const HistoryInfo$json = const {
  '1': 'HistoryInfo',
  '2': const [
    const {'1': 'transaction_id', '3': 1, '4': 1, '5': 9, '10': 'transactionId'},
    const {'1': 'time', '3': 2, '4': 1, '5': 13, '10': 'time'},
    const {'1': 'payload_type', '3': 3, '4': 1, '5': 9, '10': 'payloadType'},
    const {'1': 'description', '3': 4, '4': 1, '5': 9, '10': 'description'},
    const {'1': 'amount', '3': 5, '4': 1, '5': 3, '10': 'amount'},
  ],
};

/// Descriptor for `HistoryInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List historyInfoDescriptor = $convert.base64Decode('CgtIaXN0b3J5SW5mbxIlCg50cmFuc2FjdGlvbl9pZBgBIAEoCVINdHJhbnNhY3Rpb25JZBISCgR0aW1lGAIgASgNUgR0aW1lEiEKDHBheWxvYWRfdHlwZRgDIAEoCVILcGF5bG9hZFR5cGUSIAoLZGVzY3JpcHRpb24YBCABKAlSC2Rlc2NyaXB0aW9uEhYKBmFtb3VudBgFIAEoA1IGYW1vdW50');
@$core.Deprecated('Use getAddressHistoryRequestDescriptor instead')
const GetAddressHistoryRequest$json = const {
  '1': 'GetAddressHistoryRequest',
  '2': const [
    const {'1': 'wallet_name', '3': 1, '4': 1, '5': 9, '10': 'walletName'},
    const {'1': 'address', '3': 2, '4': 1, '5': 9, '10': 'address'},
  ],
};

/// Descriptor for `GetAddressHistoryRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getAddressHistoryRequestDescriptor = $convert.base64Decode('ChhHZXRBZGRyZXNzSGlzdG9yeVJlcXVlc3QSHwoLd2FsbGV0X25hbWUYASABKAlSCndhbGxldE5hbWUSGAoHYWRkcmVzcxgCIAEoCVIHYWRkcmVzcw==');
@$core.Deprecated('Use getAddressHistoryResponseDescriptor instead')
const GetAddressHistoryResponse$json = const {
  '1': 'GetAddressHistoryResponse',
  '2': const [
    const {'1': 'history_info', '3': 1, '4': 3, '5': 11, '6': '.pactus.HistoryInfo', '10': 'historyInfo'},
  ],
};

/// Descriptor for `GetAddressHistoryResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getAddressHistoryResponseDescriptor = $convert.base64Decode('ChlHZXRBZGRyZXNzSGlzdG9yeVJlc3BvbnNlEjYKDGhpc3RvcnlfaW5mbxgBIAMoCzITLnBhY3R1cy5IaXN0b3J5SW5mb1ILaGlzdG9yeUluZm8=');
@$core.Deprecated('Use getNewAddressRequestDescriptor instead')
const GetNewAddressRequest$json = const {
  '1': 'GetNewAddressRequest',
  '2': const [
    const {'1': 'wallet_name', '3': 1, '4': 1, '5': 9, '10': 'walletName'},
    const {'1': 'address_type', '3': 2, '4': 1, '5': 14, '6': '.pactus.AddressType', '10': 'addressType'},
    const {'1': 'label', '3': 3, '4': 1, '5': 9, '10': 'label'},
    const {'1': 'password', '3': 4, '4': 1, '5': 9, '10': 'password'},
  ],
};

/// Descriptor for `GetNewAddressRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getNewAddressRequestDescriptor = $convert.base64Decode('ChRHZXROZXdBZGRyZXNzUmVxdWVzdBIfCgt3YWxsZXRfbmFtZRgBIAEoCVIKd2FsbGV0TmFtZRI2CgxhZGRyZXNzX3R5cGUYAiABKA4yEy5wYWN0dXMuQWRkcmVzc1R5cGVSC2FkZHJlc3NUeXBlEhQKBWxhYmVsGAMgASgJUgVsYWJlbBIaCghwYXNzd29yZBgEIAEoCVIIcGFzc3dvcmQ=');
@$core.Deprecated('Use getNewAddressResponseDescriptor instead')
const GetNewAddressResponse$json = const {
  '1': 'GetNewAddressResponse',
  '2': const [
    const {'1': 'wallet_name', '3': 1, '4': 1, '5': 9, '10': 'walletName'},
    const {'1': 'address_info', '3': 2, '4': 1, '5': 11, '6': '.pactus.AddressInfo', '10': 'addressInfo'},
  ],
};

/// Descriptor for `GetNewAddressResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getNewAddressResponseDescriptor = $convert.base64Decode('ChVHZXROZXdBZGRyZXNzUmVzcG9uc2USHwoLd2FsbGV0X25hbWUYASABKAlSCndhbGxldE5hbWUSNgoMYWRkcmVzc19pbmZvGAIgASgLMhMucGFjdHVzLkFkZHJlc3NJbmZvUgthZGRyZXNzSW5mbw==');
@$core.Deprecated('Use restoreWalletRequestDescriptor instead')
const RestoreWalletRequest$json = const {
  '1': 'RestoreWalletRequest',
  '2': const [
    const {'1': 'wallet_name', '3': 1, '4': 1, '5': 9, '10': 'walletName'},
    const {'1': 'mnemonic', '3': 2, '4': 1, '5': 9, '10': 'mnemonic'},
    const {'1': 'password', '3': 3, '4': 1, '5': 9, '10': 'password'},
  ],
};

/// Descriptor for `RestoreWalletRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List restoreWalletRequestDescriptor = $convert.base64Decode('ChRSZXN0b3JlV2FsbGV0UmVxdWVzdBIfCgt3YWxsZXRfbmFtZRgBIAEoCVIKd2FsbGV0TmFtZRIaCghtbmVtb25pYxgCIAEoCVIIbW5lbW9uaWMSGgoIcGFzc3dvcmQYAyABKAlSCHBhc3N3b3Jk');
@$core.Deprecated('Use restoreWalletResponseDescriptor instead')
const RestoreWalletResponse$json = const {
  '1': 'RestoreWalletResponse',
  '2': const [
    const {'1': 'wallet_name', '3': 1, '4': 1, '5': 9, '10': 'walletName'},
  ],
};

/// Descriptor for `RestoreWalletResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List restoreWalletResponseDescriptor = $convert.base64Decode('ChVSZXN0b3JlV2FsbGV0UmVzcG9uc2USHwoLd2FsbGV0X25hbWUYASABKAlSCndhbGxldE5hbWU=');
@$core.Deprecated('Use createWalletRequestDescriptor instead')
const CreateWalletRequest$json = const {
  '1': 'CreateWalletRequest',
  '2': const [
    const {'1': 'wallet_name', '3': 1, '4': 1, '5': 9, '10': 'walletName'},
    const {'1': 'password', '3': 4, '4': 1, '5': 9, '10': 'password'},
  ],
};

/// Descriptor for `CreateWalletRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createWalletRequestDescriptor = $convert.base64Decode('ChNDcmVhdGVXYWxsZXRSZXF1ZXN0Eh8KC3dhbGxldF9uYW1lGAEgASgJUgp3YWxsZXROYW1lEhoKCHBhc3N3b3JkGAQgASgJUghwYXNzd29yZA==');
@$core.Deprecated('Use createWalletResponseDescriptor instead')
const CreateWalletResponse$json = const {
  '1': 'CreateWalletResponse',
  '2': const [
    const {'1': 'mnemonic', '3': 2, '4': 1, '5': 9, '10': 'mnemonic'},
  ],
};

/// Descriptor for `CreateWalletResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createWalletResponseDescriptor = $convert.base64Decode('ChRDcmVhdGVXYWxsZXRSZXNwb25zZRIaCghtbmVtb25pYxgCIAEoCVIIbW5lbW9uaWM=');
@$core.Deprecated('Use loadWalletRequestDescriptor instead')
const LoadWalletRequest$json = const {
  '1': 'LoadWalletRequest',
  '2': const [
    const {'1': 'wallet_name', '3': 1, '4': 1, '5': 9, '10': 'walletName'},
  ],
};

/// Descriptor for `LoadWalletRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List loadWalletRequestDescriptor = $convert.base64Decode('ChFMb2FkV2FsbGV0UmVxdWVzdBIfCgt3YWxsZXRfbmFtZRgBIAEoCVIKd2FsbGV0TmFtZQ==');
@$core.Deprecated('Use loadWalletResponseDescriptor instead')
const LoadWalletResponse$json = const {
  '1': 'LoadWalletResponse',
  '2': const [
    const {'1': 'wallet_name', '3': 1, '4': 1, '5': 9, '10': 'walletName'},
  ],
};

/// Descriptor for `LoadWalletResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List loadWalletResponseDescriptor = $convert.base64Decode('ChJMb2FkV2FsbGV0UmVzcG9uc2USHwoLd2FsbGV0X25hbWUYASABKAlSCndhbGxldE5hbWU=');
@$core.Deprecated('Use unloadWalletRequestDescriptor instead')
const UnloadWalletRequest$json = const {
  '1': 'UnloadWalletRequest',
  '2': const [
    const {'1': 'wallet_name', '3': 1, '4': 1, '5': 9, '10': 'walletName'},
  ],
};

/// Descriptor for `UnloadWalletRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List unloadWalletRequestDescriptor = $convert.base64Decode('ChNVbmxvYWRXYWxsZXRSZXF1ZXN0Eh8KC3dhbGxldF9uYW1lGAEgASgJUgp3YWxsZXROYW1l');
@$core.Deprecated('Use unloadWalletResponseDescriptor instead')
const UnloadWalletResponse$json = const {
  '1': 'UnloadWalletResponse',
  '2': const [
    const {'1': 'wallet_name', '3': 1, '4': 1, '5': 9, '10': 'walletName'},
  ],
};

/// Descriptor for `UnloadWalletResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List unloadWalletResponseDescriptor = $convert.base64Decode('ChRVbmxvYWRXYWxsZXRSZXNwb25zZRIfCgt3YWxsZXRfbmFtZRgBIAEoCVIKd2FsbGV0TmFtZQ==');
@$core.Deprecated('Use getValidatorAddressRequestDescriptor instead')
const GetValidatorAddressRequest$json = const {
  '1': 'GetValidatorAddressRequest',
  '2': const [
    const {'1': 'public_key', '3': 1, '4': 1, '5': 9, '10': 'publicKey'},
  ],
};

/// Descriptor for `GetValidatorAddressRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getValidatorAddressRequestDescriptor = $convert.base64Decode('ChpHZXRWYWxpZGF0b3JBZGRyZXNzUmVxdWVzdBIdCgpwdWJsaWNfa2V5GAEgASgJUglwdWJsaWNLZXk=');
@$core.Deprecated('Use getValidatorAddressResponseDescriptor instead')
const GetValidatorAddressResponse$json = const {
  '1': 'GetValidatorAddressResponse',
  '2': const [
    const {'1': 'address', '3': 1, '4': 1, '5': 9, '10': 'address'},
  ],
};

/// Descriptor for `GetValidatorAddressResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getValidatorAddressResponseDescriptor = $convert.base64Decode('ChtHZXRWYWxpZGF0b3JBZGRyZXNzUmVzcG9uc2USGAoHYWRkcmVzcxgBIAEoCVIHYWRkcmVzcw==');
@$core.Deprecated('Use signRawTransactionRequestDescriptor instead')
const SignRawTransactionRequest$json = const {
  '1': 'SignRawTransactionRequest',
  '2': const [
    const {'1': 'wallet_name', '3': 1, '4': 1, '5': 9, '10': 'walletName'},
    const {'1': 'raw_transaction', '3': 2, '4': 1, '5': 9, '10': 'rawTransaction'},
    const {'1': 'password', '3': 3, '4': 1, '5': 9, '10': 'password'},
  ],
};

/// Descriptor for `SignRawTransactionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List signRawTransactionRequestDescriptor = $convert.base64Decode('ChlTaWduUmF3VHJhbnNhY3Rpb25SZXF1ZXN0Eh8KC3dhbGxldF9uYW1lGAEgASgJUgp3YWxsZXROYW1lEicKD3Jhd190cmFuc2FjdGlvbhgCIAEoCVIOcmF3VHJhbnNhY3Rpb24SGgoIcGFzc3dvcmQYAyABKAlSCHBhc3N3b3Jk');
@$core.Deprecated('Use signRawTransactionResponseDescriptor instead')
const SignRawTransactionResponse$json = const {
  '1': 'SignRawTransactionResponse',
  '2': const [
    const {'1': 'transaction_id', '3': 1, '4': 1, '5': 9, '10': 'transactionId'},
    const {'1': 'signed_raw_transaction', '3': 2, '4': 1, '5': 9, '10': 'signedRawTransaction'},
  ],
};

/// Descriptor for `SignRawTransactionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List signRawTransactionResponseDescriptor = $convert.base64Decode('ChpTaWduUmF3VHJhbnNhY3Rpb25SZXNwb25zZRIlCg50cmFuc2FjdGlvbl9pZBgBIAEoCVINdHJhbnNhY3Rpb25JZBI0ChZzaWduZWRfcmF3X3RyYW5zYWN0aW9uGAIgASgJUhRzaWduZWRSYXdUcmFuc2FjdGlvbg==');
@$core.Deprecated('Use getTotalBalanceRequestDescriptor instead')
const GetTotalBalanceRequest$json = const {
  '1': 'GetTotalBalanceRequest',
  '2': const [
    const {'1': 'wallet_name', '3': 1, '4': 1, '5': 9, '10': 'walletName'},
  ],
};

/// Descriptor for `GetTotalBalanceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getTotalBalanceRequestDescriptor = $convert.base64Decode('ChZHZXRUb3RhbEJhbGFuY2VSZXF1ZXN0Eh8KC3dhbGxldF9uYW1lGAEgASgJUgp3YWxsZXROYW1l');
@$core.Deprecated('Use getTotalBalanceResponseDescriptor instead')
const GetTotalBalanceResponse$json = const {
  '1': 'GetTotalBalanceResponse',
  '2': const [
    const {'1': 'wallet_name', '3': 1, '4': 1, '5': 9, '10': 'walletName'},
    const {'1': 'total_balance', '3': 2, '4': 1, '5': 3, '10': 'totalBalance'},
  ],
};

/// Descriptor for `GetTotalBalanceResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getTotalBalanceResponseDescriptor = $convert.base64Decode('ChdHZXRUb3RhbEJhbGFuY2VSZXNwb25zZRIfCgt3YWxsZXRfbmFtZRgBIAEoCVIKd2FsbGV0TmFtZRIjCg10b3RhbF9iYWxhbmNlGAIgASgDUgx0b3RhbEJhbGFuY2U=');
@$core.Deprecated('Use signMessageRequestDescriptor instead')
const SignMessageRequest$json = const {
  '1': 'SignMessageRequest',
  '2': const [
    const {'1': 'wallet_name', '3': 1, '4': 1, '5': 9, '10': 'walletName'},
    const {'1': 'password', '3': 2, '4': 1, '5': 9, '10': 'password'},
    const {'1': 'address', '3': 3, '4': 1, '5': 9, '10': 'address'},
    const {'1': 'message', '3': 4, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `SignMessageRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List signMessageRequestDescriptor = $convert.base64Decode('ChJTaWduTWVzc2FnZVJlcXVlc3QSHwoLd2FsbGV0X25hbWUYASABKAlSCndhbGxldE5hbWUSGgoIcGFzc3dvcmQYAiABKAlSCHBhc3N3b3JkEhgKB2FkZHJlc3MYAyABKAlSB2FkZHJlc3MSGAoHbWVzc2FnZRgEIAEoCVIHbWVzc2FnZQ==');
@$core.Deprecated('Use signMessageResponseDescriptor instead')
const SignMessageResponse$json = const {
  '1': 'SignMessageResponse',
  '2': const [
    const {'1': 'signature', '3': 1, '4': 1, '5': 9, '10': 'signature'},
  ],
};

/// Descriptor for `SignMessageResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List signMessageResponseDescriptor = $convert.base64Decode('ChNTaWduTWVzc2FnZVJlc3BvbnNlEhwKCXNpZ25hdHVyZRgBIAEoCVIJc2lnbmF0dXJl');
@$core.Deprecated('Use getTotalStakeRequestDescriptor instead')
const GetTotalStakeRequest$json = const {
  '1': 'GetTotalStakeRequest',
  '2': const [
    const {'1': 'wallet_name', '3': 1, '4': 1, '5': 9, '10': 'walletName'},
  ],
};

/// Descriptor for `GetTotalStakeRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getTotalStakeRequestDescriptor = $convert.base64Decode('ChRHZXRUb3RhbFN0YWtlUmVxdWVzdBIfCgt3YWxsZXRfbmFtZRgBIAEoCVIKd2FsbGV0TmFtZQ==');
@$core.Deprecated('Use getTotalStakeResponseDescriptor instead')
const GetTotalStakeResponse$json = const {
  '1': 'GetTotalStakeResponse',
  '2': const [
    const {'1': 'total_stake', '3': 1, '4': 1, '5': 3, '10': 'totalStake'},
    const {'1': 'wallet_name', '3': 2, '4': 1, '5': 9, '10': 'walletName'},
  ],
};

/// Descriptor for `GetTotalStakeResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getTotalStakeResponseDescriptor = $convert.base64Decode('ChVHZXRUb3RhbFN0YWtlUmVzcG9uc2USHwoLdG90YWxfc3Rha2UYASABKANSCnRvdGFsU3Rha2USHwoLd2FsbGV0X25hbWUYAiABKAlSCndhbGxldE5hbWU=');
@$core.Deprecated('Use getAddressInfoRequestDescriptor instead')
const GetAddressInfoRequest$json = const {
  '1': 'GetAddressInfoRequest',
  '2': const [
    const {'1': 'wallet_name', '3': 1, '4': 1, '5': 9, '10': 'walletName'},
    const {'1': 'address', '3': 2, '4': 1, '5': 9, '10': 'address'},
  ],
};

/// Descriptor for `GetAddressInfoRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getAddressInfoRequestDescriptor = $convert.base64Decode('ChVHZXRBZGRyZXNzSW5mb1JlcXVlc3QSHwoLd2FsbGV0X25hbWUYASABKAlSCndhbGxldE5hbWUSGAoHYWRkcmVzcxgCIAEoCVIHYWRkcmVzcw==');
@$core.Deprecated('Use getAddressInfoResponseDescriptor instead')
const GetAddressInfoResponse$json = const {
  '1': 'GetAddressInfoResponse',
  '2': const [
    const {'1': 'address', '3': 1, '4': 1, '5': 9, '10': 'address'},
    const {'1': 'label', '3': 2, '4': 1, '5': 9, '10': 'label'},
    const {'1': 'public_key', '3': 3, '4': 1, '5': 9, '10': 'publicKey'},
    const {'1': 'path', '3': 4, '4': 1, '5': 9, '10': 'path'},
    const {'1': 'wallet_name', '3': 5, '4': 1, '5': 9, '10': 'walletName'},
  ],
};

/// Descriptor for `GetAddressInfoResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getAddressInfoResponseDescriptor = $convert.base64Decode('ChZHZXRBZGRyZXNzSW5mb1Jlc3BvbnNlEhgKB2FkZHJlc3MYASABKAlSB2FkZHJlc3MSFAoFbGFiZWwYAiABKAlSBWxhYmVsEh0KCnB1YmxpY19rZXkYAyABKAlSCXB1YmxpY0tleRISCgRwYXRoGAQgASgJUgRwYXRoEh8KC3dhbGxldF9uYW1lGAUgASgJUgp3YWxsZXROYW1l');
@$core.Deprecated('Use setLabelRequestDescriptor instead')
const SetLabelRequest$json = const {
  '1': 'SetLabelRequest',
  '2': const [
    const {'1': 'wallet_name', '3': 1, '4': 1, '5': 9, '10': 'walletName'},
    const {'1': 'password', '3': 3, '4': 1, '5': 9, '10': 'password'},
    const {'1': 'address', '3': 4, '4': 1, '5': 9, '10': 'address'},
    const {'1': 'label', '3': 5, '4': 1, '5': 9, '10': 'label'},
  ],
};

/// Descriptor for `SetLabelRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List setLabelRequestDescriptor = $convert.base64Decode('Cg9TZXRMYWJlbFJlcXVlc3QSHwoLd2FsbGV0X25hbWUYASABKAlSCndhbGxldE5hbWUSGgoIcGFzc3dvcmQYAyABKAlSCHBhc3N3b3JkEhgKB2FkZHJlc3MYBCABKAlSB2FkZHJlc3MSFAoFbGFiZWwYBSABKAlSBWxhYmVs');
@$core.Deprecated('Use setLabelResponseDescriptor instead')
const SetLabelResponse$json = const {
  '1': 'SetLabelResponse',
};

/// Descriptor for `SetLabelResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List setLabelResponseDescriptor = $convert.base64Decode('ChBTZXRMYWJlbFJlc3BvbnNl');
@$core.Deprecated('Use listWalletRequestDescriptor instead')
const ListWalletRequest$json = const {
  '1': 'ListWalletRequest',
};

/// Descriptor for `ListWalletRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listWalletRequestDescriptor = $convert.base64Decode('ChFMaXN0V2FsbGV0UmVxdWVzdA==');
@$core.Deprecated('Use listWalletResponseDescriptor instead')
const ListWalletResponse$json = const {
  '1': 'ListWalletResponse',
  '2': const [
    const {'1': 'wallets', '3': 1, '4': 3, '5': 9, '10': 'wallets'},
  ],
};

/// Descriptor for `ListWalletResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listWalletResponseDescriptor = $convert.base64Decode('ChJMaXN0V2FsbGV0UmVzcG9uc2USGAoHd2FsbGV0cxgBIAMoCVIHd2FsbGV0cw==');
@$core.Deprecated('Use getWalletInfoRequestDescriptor instead')
const GetWalletInfoRequest$json = const {
  '1': 'GetWalletInfoRequest',
  '2': const [
    const {'1': 'wallet_name', '3': 1, '4': 1, '5': 9, '10': 'walletName'},
  ],
};

/// Descriptor for `GetWalletInfoRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getWalletInfoRequestDescriptor = $convert.base64Decode('ChRHZXRXYWxsZXRJbmZvUmVxdWVzdBIfCgt3YWxsZXRfbmFtZRgBIAEoCVIKd2FsbGV0TmFtZQ==');
@$core.Deprecated('Use getWalletInfoResponseDescriptor instead')
const GetWalletInfoResponse$json = const {
  '1': 'GetWalletInfoResponse',
  '2': const [
    const {'1': 'wallet_name', '3': 1, '4': 1, '5': 9, '10': 'walletName'},
    const {'1': 'version', '3': 2, '4': 1, '5': 3, '10': 'version'},
    const {'1': 'network', '3': 3, '4': 1, '5': 9, '10': 'network'},
    const {'1': 'encrypted', '3': 4, '4': 1, '5': 8, '10': 'encrypted'},
    const {'1': 'uuid', '3': 5, '4': 1, '5': 9, '10': 'uuid'},
    const {'1': 'created_at', '3': 6, '4': 1, '5': 3, '10': 'createdAt'},
  ],
};

/// Descriptor for `GetWalletInfoResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getWalletInfoResponseDescriptor = $convert.base64Decode('ChVHZXRXYWxsZXRJbmZvUmVzcG9uc2USHwoLd2FsbGV0X25hbWUYASABKAlSCndhbGxldE5hbWUSGAoHdmVyc2lvbhgCIAEoA1IHdmVyc2lvbhIYCgduZXR3b3JrGAMgASgJUgduZXR3b3JrEhwKCWVuY3J5cHRlZBgEIAEoCFIJZW5jcnlwdGVkEhIKBHV1aWQYBSABKAlSBHV1aWQSHQoKY3JlYXRlZF9hdBgGIAEoA1IJY3JlYXRlZEF0');
@$core.Deprecated('Use listAddressRequestDescriptor instead')
const ListAddressRequest$json = const {
  '1': 'ListAddressRequest',
  '2': const [
    const {'1': 'wallet_name', '3': 1, '4': 1, '5': 9, '10': 'walletName'},
  ],
};

/// Descriptor for `ListAddressRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listAddressRequestDescriptor = $convert.base64Decode('ChJMaXN0QWRkcmVzc1JlcXVlc3QSHwoLd2FsbGV0X25hbWUYASABKAlSCndhbGxldE5hbWU=');
@$core.Deprecated('Use listAddressResponseDescriptor instead')
const ListAddressResponse$json = const {
  '1': 'ListAddressResponse',
  '2': const [
    const {'1': 'data', '3': 1, '4': 3, '5': 11, '6': '.pactus.AddressInfo', '10': 'data'},
  ],
};

/// Descriptor for `ListAddressResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listAddressResponseDescriptor = $convert.base64Decode('ChNMaXN0QWRkcmVzc1Jlc3BvbnNlEicKBGRhdGEYASADKAsyEy5wYWN0dXMuQWRkcmVzc0luZm9SBGRhdGE=');
const $core.Map<$core.String, $core.dynamic> WalletServiceBase$json = const {
  '1': 'Wallet',
  '2': const [
    const {'1': 'CreateWallet', '2': '.pactus.CreateWalletRequest', '3': '.pactus.CreateWalletResponse'},
    const {'1': 'RestoreWallet', '2': '.pactus.RestoreWalletRequest', '3': '.pactus.RestoreWalletResponse'},
    const {'1': 'LoadWallet', '2': '.pactus.LoadWalletRequest', '3': '.pactus.LoadWalletResponse'},
    const {'1': 'UnloadWallet', '2': '.pactus.UnloadWalletRequest', '3': '.pactus.UnloadWalletResponse'},
    const {'1': 'GetTotalBalance', '2': '.pactus.GetTotalBalanceRequest', '3': '.pactus.GetTotalBalanceResponse'},
    const {'1': 'SignRawTransaction', '2': '.pactus.SignRawTransactionRequest', '3': '.pactus.SignRawTransactionResponse'},
    const {'1': 'GetValidatorAddress', '2': '.pactus.GetValidatorAddressRequest', '3': '.pactus.GetValidatorAddressResponse'},
    const {'1': 'GetNewAddress', '2': '.pactus.GetNewAddressRequest', '3': '.pactus.GetNewAddressResponse'},
    const {'1': 'GetAddressHistory', '2': '.pactus.GetAddressHistoryRequest', '3': '.pactus.GetAddressHistoryResponse'},
    const {'1': 'SignMessage', '2': '.pactus.SignMessageRequest', '3': '.pactus.SignMessageResponse'},
    const {'1': 'GetTotalStake', '2': '.pactus.GetTotalStakeRequest', '3': '.pactus.GetTotalStakeResponse'},
    const {'1': 'GetAddressInfo', '2': '.pactus.GetAddressInfoRequest', '3': '.pactus.GetAddressInfoResponse'},
    const {'1': 'SetAddressLabel', '2': '.pactus.SetLabelRequest', '3': '.pactus.SetLabelResponse'},
    const {'1': 'ListWallet', '2': '.pactus.ListWalletRequest', '3': '.pactus.ListWalletResponse'},
    const {'1': 'GetWalletInfo', '2': '.pactus.GetWalletInfoRequest', '3': '.pactus.GetWalletInfoResponse'},
    const {'1': 'ListAddress', '2': '.pactus.ListAddressRequest', '3': '.pactus.ListAddressResponse'},
  ],
};

@$core.Deprecated('Use walletServiceDescriptor instead')
const $core.Map<$core.String, $core.Map<$core.String, $core.dynamic>> WalletServiceBase$messageJson = const {
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
  '.pactus.SetLabelRequest': SetLabelRequest$json,
  '.pactus.SetLabelResponse': SetLabelResponse$json,
  '.pactus.ListWalletRequest': ListWalletRequest$json,
  '.pactus.ListWalletResponse': ListWalletResponse$json,
  '.pactus.GetWalletInfoRequest': GetWalletInfoRequest$json,
  '.pactus.GetWalletInfoResponse': GetWalletInfoResponse$json,
  '.pactus.ListAddressRequest': ListAddressRequest$json,
  '.pactus.ListAddressResponse': ListAddressResponse$json,
};

/// Descriptor for `Wallet`. Decode as a `google.protobuf.ServiceDescriptorProto`.
final $typed_data.Uint8List walletServiceDescriptor = $convert.base64Decode('CgZXYWxsZXQSSQoMQ3JlYXRlV2FsbGV0EhsucGFjdHVzLkNyZWF0ZVdhbGxldFJlcXVlc3QaHC5wYWN0dXMuQ3JlYXRlV2FsbGV0UmVzcG9uc2USTAoNUmVzdG9yZVdhbGxldBIcLnBhY3R1cy5SZXN0b3JlV2FsbGV0UmVxdWVzdBodLnBhY3R1cy5SZXN0b3JlV2FsbGV0UmVzcG9uc2USQwoKTG9hZFdhbGxldBIZLnBhY3R1cy5Mb2FkV2FsbGV0UmVxdWVzdBoaLnBhY3R1cy5Mb2FkV2FsbGV0UmVzcG9uc2USSQoMVW5sb2FkV2FsbGV0EhsucGFjdHVzLlVubG9hZFdhbGxldFJlcXVlc3QaHC5wYWN0dXMuVW5sb2FkV2FsbGV0UmVzcG9uc2USUgoPR2V0VG90YWxCYWxhbmNlEh4ucGFjdHVzLkdldFRvdGFsQmFsYW5jZVJlcXVlc3QaHy5wYWN0dXMuR2V0VG90YWxCYWxhbmNlUmVzcG9uc2USWwoSU2lnblJhd1RyYW5zYWN0aW9uEiEucGFjdHVzLlNpZ25SYXdUcmFuc2FjdGlvblJlcXVlc3QaIi5wYWN0dXMuU2lnblJhd1RyYW5zYWN0aW9uUmVzcG9uc2USXgoTR2V0VmFsaWRhdG9yQWRkcmVzcxIiLnBhY3R1cy5HZXRWYWxpZGF0b3JBZGRyZXNzUmVxdWVzdBojLnBhY3R1cy5HZXRWYWxpZGF0b3JBZGRyZXNzUmVzcG9uc2USTAoNR2V0TmV3QWRkcmVzcxIcLnBhY3R1cy5HZXROZXdBZGRyZXNzUmVxdWVzdBodLnBhY3R1cy5HZXROZXdBZGRyZXNzUmVzcG9uc2USWAoRR2V0QWRkcmVzc0hpc3RvcnkSIC5wYWN0dXMuR2V0QWRkcmVzc0hpc3RvcnlSZXF1ZXN0GiEucGFjdHVzLkdldEFkZHJlc3NIaXN0b3J5UmVzcG9uc2USRgoLU2lnbk1lc3NhZ2USGi5wYWN0dXMuU2lnbk1lc3NhZ2VSZXF1ZXN0GhsucGFjdHVzLlNpZ25NZXNzYWdlUmVzcG9uc2USTAoNR2V0VG90YWxTdGFrZRIcLnBhY3R1cy5HZXRUb3RhbFN0YWtlUmVxdWVzdBodLnBhY3R1cy5HZXRUb3RhbFN0YWtlUmVzcG9uc2USTwoOR2V0QWRkcmVzc0luZm8SHS5wYWN0dXMuR2V0QWRkcmVzc0luZm9SZXF1ZXN0Gh4ucGFjdHVzLkdldEFkZHJlc3NJbmZvUmVzcG9uc2USRAoPU2V0QWRkcmVzc0xhYmVsEhcucGFjdHVzLlNldExhYmVsUmVxdWVzdBoYLnBhY3R1cy5TZXRMYWJlbFJlc3BvbnNlEkMKCkxpc3RXYWxsZXQSGS5wYWN0dXMuTGlzdFdhbGxldFJlcXVlc3QaGi5wYWN0dXMuTGlzdFdhbGxldFJlc3BvbnNlEkwKDUdldFdhbGxldEluZm8SHC5wYWN0dXMuR2V0V2FsbGV0SW5mb1JlcXVlc3QaHS5wYWN0dXMuR2V0V2FsbGV0SW5mb1Jlc3BvbnNlEkYKC0xpc3RBZGRyZXNzEhoucGFjdHVzLkxpc3RBZGRyZXNzUmVxdWVzdBobLnBhY3R1cy5MaXN0QWRkcmVzc1Jlc3BvbnNl');
