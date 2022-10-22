///
//  Generated code. Do not modify.
//  source: wallet.proto
//
// @dart = 2.12
// ignore_for_file: annotate_overrides,camel_case_types,constant_identifier_names,deprecated_member_use_from_same_package,directives_ordering,library_prefixes,non_constant_identifier_names,prefer_final_fields,return_of_invalid_type,unnecessary_const,unnecessary_import,unnecessary_this,unused_import,unused_shown_name

import 'dart:core' as $core;
import 'dart:convert' as $convert;
import 'dart:typed_data' as $typed_data;
@$core.Deprecated('Use generateMnemonicRequestDescriptor instead')
const GenerateMnemonicRequest$json = const {
  '1': 'GenerateMnemonicRequest',
  '2': const [
    const {'1': 'entropy', '3': 1, '4': 1, '5': 5, '10': 'entropy'},
    const {'1': 'language', '3': 2, '4': 1, '5': 9, '10': 'language'},
  ],
};

/// Descriptor for `GenerateMnemonicRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List generateMnemonicRequestDescriptor = $convert.base64Decode('ChdHZW5lcmF0ZU1uZW1vbmljUmVxdWVzdBIYCgdlbnRyb3B5GAEgASgFUgdlbnRyb3B5EhoKCGxhbmd1YWdlGAIgASgJUghsYW5ndWFnZQ==');
@$core.Deprecated('Use generateMnemonicResponseDescriptor instead')
const GenerateMnemonicResponse$json = const {
  '1': 'GenerateMnemonicResponse',
  '2': const [
    const {'1': 'mnemonic', '3': 1, '4': 1, '5': 9, '10': 'mnemonic'},
  ],
};

/// Descriptor for `GenerateMnemonicResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List generateMnemonicResponseDescriptor = $convert.base64Decode('ChhHZW5lcmF0ZU1uZW1vbmljUmVzcG9uc2USGgoIbW5lbW9uaWMYASABKAlSCG1uZW1vbmlj');
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
const $core.Map<$core.String, $core.dynamic> WalletServiceBase$json = const {
  '1': 'Wallet',
  '2': const [
    const {'1': 'GenerateMnemonic', '2': '.pactus.GenerateMnemonicRequest', '3': '.pactus.GenerateMnemonicResponse', '4': const {}},
    const {'1': 'CreateWallet', '2': '.pactus.CreateWalletRequest', '3': '.pactus.CreateWalletResponse', '4': const {}},
  ],
};

@$core.Deprecated('Use walletServiceDescriptor instead')
const $core.Map<$core.String, $core.Map<$core.String, $core.dynamic>> WalletServiceBase$messageJson = const {
  '.pactus.GenerateMnemonicRequest': GenerateMnemonicRequest$json,
  '.pactus.GenerateMnemonicResponse': GenerateMnemonicResponse$json,
  '.pactus.CreateWalletRequest': CreateWalletRequest$json,
  '.pactus.CreateWalletResponse': CreateWalletResponse$json,
};

/// Descriptor for `Wallet`. Decode as a `google.protobuf.ServiceDescriptorProto`.
final $typed_data.Uint8List walletServiceDescriptor = $convert.base64Decode('CgZXYWxsZXQSmAEKEEdlbmVyYXRlTW5lbW9uaWMSHy5wYWN0dXMuR2VuZXJhdGVNbmVtb25pY1JlcXVlc3QaIC5wYWN0dXMuR2VuZXJhdGVNbmVtb25pY1Jlc3BvbnNlIkGC0+STAjsSOS92MS93YWxsZXQvbW5lbW9uaWMvZW50cm9weS97ZW50cm9weX0vbGFuZ3VhZ2Uve2xhbmd1YWdlfRKsAQoMQ3JlYXRlV2FsbGV0EhsucGFjdHVzLkNyZWF0ZVdhbGxldFJlcXVlc3QaHC5wYWN0dXMuQ3JlYXRlV2FsbGV0UmVzcG9uc2UiYYLT5JMCWxJZL3YxL3dhbGxldC9jcmVhdGUvbmFtZS97bmFtZX0vbW5lbW9uaWMve21uZW1vbmljfS9sYW5ndWFnZS97bGFuZ3VhZ2V9L3Bhc3N3b3JkL3twYXNzd29yZH0=');
