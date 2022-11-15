///
//  Generated code. Do not modify.
//  source: wallet.proto
//
// @dart = 2.12
// ignore_for_file: annotate_overrides,camel_case_types,constant_identifier_names,deprecated_member_use_from_same_package,directives_ordering,library_prefixes,non_constant_identifier_names,prefer_final_fields,return_of_invalid_type,unnecessary_const,unnecessary_import,unnecessary_this,unused_import,unused_shown_name

import 'dart:core' as $core;
import 'dart:convert' as $convert;
import 'dart:typed_data' as $typed_data;
@$core.Deprecated('Use createWalletRequestDescriptor instead')
const CreateWalletRequest$json = const {
  '1': 'CreateWalletRequest',
  '2': const [
    const {'1': 'name', '3': 1, '4': 1, '5': 9, '10': 'name'},
    const {'1': 'mnemonic', '3': 2, '4': 1, '5': 9, '10': 'mnemonic'},
    const {'1': 'language', '3': 3, '4': 1, '5': 9, '10': 'language'},
    const {'1': 'password', '3': 4, '4': 1, '5': 9, '10': 'password'},
  ],
};

/// Descriptor for `CreateWalletRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createWalletRequestDescriptor = $convert.base64Decode('ChNDcmVhdGVXYWxsZXRSZXF1ZXN0EhIKBG5hbWUYASABKAlSBG5hbWUSGgoIbW5lbW9uaWMYAiABKAlSCG1uZW1vbmljEhoKCGxhbmd1YWdlGAMgASgJUghsYW5ndWFnZRIaCghwYXNzd29yZBgEIAEoCVIIcGFzc3dvcmQ=');
@$core.Deprecated('Use createWalletResponseDescriptor instead')
const CreateWalletResponse$json = const {
  '1': 'CreateWalletResponse',
};

/// Descriptor for `CreateWalletResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createWalletResponseDescriptor = $convert.base64Decode('ChRDcmVhdGVXYWxsZXRSZXNwb25zZQ==');
@$core.Deprecated('Use loadWalletRequestDescriptor instead')
const LoadWalletRequest$json = const {
  '1': 'LoadWalletRequest',
  '2': const [
    const {'1': 'name', '3': 1, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `LoadWalletRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List loadWalletRequestDescriptor = $convert.base64Decode('ChFMb2FkV2FsbGV0UmVxdWVzdBISCgRuYW1lGAEgASgJUgRuYW1l');
@$core.Deprecated('Use loadWalletResponseDescriptor instead')
const LoadWalletResponse$json = const {
  '1': 'LoadWalletResponse',
  '2': const [
    const {'1': 'name', '3': 1, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `LoadWalletResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List loadWalletResponseDescriptor = $convert.base64Decode('ChJMb2FkV2FsbGV0UmVzcG9uc2USEgoEbmFtZRgBIAEoCVIEbmFtZQ==');
@$core.Deprecated('Use unloadWalletRequestDescriptor instead')
const UnloadWalletRequest$json = const {
  '1': 'UnloadWalletRequest',
  '2': const [
    const {'1': 'name', '3': 1, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `UnloadWalletRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List unloadWalletRequestDescriptor = $convert.base64Decode('ChNVbmxvYWRXYWxsZXRSZXF1ZXN0EhIKBG5hbWUYASABKAlSBG5hbWU=');
@$core.Deprecated('Use unloadWalletResponseDescriptor instead')
const UnloadWalletResponse$json = const {
  '1': 'UnloadWalletResponse',
  '2': const [
    const {'1': 'name', '3': 1, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `UnloadWalletResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List unloadWalletResponseDescriptor = $convert.base64Decode('ChRVbmxvYWRXYWxsZXRSZXNwb25zZRISCgRuYW1lGAEgASgJUgRuYW1l');
@$core.Deprecated('Use lockWalletRequestDescriptor instead')
const LockWalletRequest$json = const {
  '1': 'LockWalletRequest',
  '2': const [
    const {'1': 'password', '3': 1, '4': 1, '5': 9, '10': 'password'},
    const {'1': 'timeout', '3': 2, '4': 1, '5': 5, '10': 'timeout'},
  ],
};

/// Descriptor for `LockWalletRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List lockWalletRequestDescriptor = $convert.base64Decode('ChFMb2NrV2FsbGV0UmVxdWVzdBIaCghwYXNzd29yZBgBIAEoCVIIcGFzc3dvcmQSGAoHdGltZW91dBgCIAEoBVIHdGltZW91dA==');
@$core.Deprecated('Use lockWalletResponseDescriptor instead')
const LockWalletResponse$json = const {
  '1': 'LockWalletResponse',
};

/// Descriptor for `LockWalletResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List lockWalletResponseDescriptor = $convert.base64Decode('ChJMb2NrV2FsbGV0UmVzcG9uc2U=');
@$core.Deprecated('Use unlockWalletRequestDescriptor instead')
const UnlockWalletRequest$json = const {
  '1': 'UnlockWalletRequest',
  '2': const [
    const {'1': 'password', '3': 1, '4': 1, '5': 9, '10': 'password'},
    const {'1': 'timeout', '3': 2, '4': 1, '5': 5, '10': 'timeout'},
  ],
};

/// Descriptor for `UnlockWalletRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List unlockWalletRequestDescriptor = $convert.base64Decode('ChNVbmxvY2tXYWxsZXRSZXF1ZXN0EhoKCHBhc3N3b3JkGAEgASgJUghwYXNzd29yZBIYCgd0aW1lb3V0GAIgASgFUgd0aW1lb3V0');
@$core.Deprecated('Use unlockWalletResponseDescriptor instead')
const UnlockWalletResponse$json = const {
  '1': 'UnlockWalletResponse',
};

/// Descriptor for `UnlockWalletResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List unlockWalletResponseDescriptor = $convert.base64Decode('ChRVbmxvY2tXYWxsZXRSZXNwb25zZQ==');
const $core.Map<$core.String, $core.dynamic> WalletServiceBase$json = const {
  '1': 'Wallet',
  '2': const [
    const {'1': 'CreateWallet', '2': '.pactus.CreateWalletRequest', '3': '.pactus.CreateWalletResponse'},
    const {'1': 'LoadWallet', '2': '.pactus.LoadWalletRequest', '3': '.pactus.LoadWalletResponse'},
    const {'1': 'UnloadWallet', '2': '.pactus.UnloadWalletRequest', '3': '.pactus.UnloadWalletResponse'},
    const {'1': 'LockWallet', '2': '.pactus.LockWalletRequest', '3': '.pactus.LockWalletResponse'},
    const {'1': 'UnlockWallet', '2': '.pactus.UnlockWalletRequest', '3': '.pactus.UnlockWalletResponse'},
  ],
};

@$core.Deprecated('Use walletServiceDescriptor instead')
const $core.Map<$core.String, $core.Map<$core.String, $core.dynamic>> WalletServiceBase$messageJson = const {
  '.pactus.CreateWalletRequest': CreateWalletRequest$json,
  '.pactus.CreateWalletResponse': CreateWalletResponse$json,
  '.pactus.LoadWalletRequest': LoadWalletRequest$json,
  '.pactus.LoadWalletResponse': LoadWalletResponse$json,
  '.pactus.UnloadWalletRequest': UnloadWalletRequest$json,
  '.pactus.UnloadWalletResponse': UnloadWalletResponse$json,
  '.pactus.LockWalletRequest': LockWalletRequest$json,
  '.pactus.LockWalletResponse': LockWalletResponse$json,
  '.pactus.UnlockWalletRequest': UnlockWalletRequest$json,
  '.pactus.UnlockWalletResponse': UnlockWalletResponse$json,
};

/// Descriptor for `Wallet`. Decode as a `google.protobuf.ServiceDescriptorProto`.
final $typed_data.Uint8List walletServiceDescriptor = $convert.base64Decode('CgZXYWxsZXQSSQoMQ3JlYXRlV2FsbGV0EhsucGFjdHVzLkNyZWF0ZVdhbGxldFJlcXVlc3QaHC5wYWN0dXMuQ3JlYXRlV2FsbGV0UmVzcG9uc2USQwoKTG9hZFdhbGxldBIZLnBhY3R1cy5Mb2FkV2FsbGV0UmVxdWVzdBoaLnBhY3R1cy5Mb2FkV2FsbGV0UmVzcG9uc2USSQoMVW5sb2FkV2FsbGV0EhsucGFjdHVzLlVubG9hZFdhbGxldFJlcXVlc3QaHC5wYWN0dXMuVW5sb2FkV2FsbGV0UmVzcG9uc2USQwoKTG9ja1dhbGxldBIZLnBhY3R1cy5Mb2NrV2FsbGV0UmVxdWVzdBoaLnBhY3R1cy5Mb2NrV2FsbGV0UmVzcG9uc2USSQoMVW5sb2NrV2FsbGV0EhsucGFjdHVzLlVubG9ja1dhbGxldFJlcXVlc3QaHC5wYWN0dXMuVW5sb2NrV2FsbGV0UmVzcG9uc2U=');
