///
//  Generated code. Do not modify.
//  source: util.proto
//
// @dart = 2.12
// ignore_for_file: annotate_overrides,camel_case_types,constant_identifier_names,deprecated_member_use_from_same_package,directives_ordering,library_prefixes,non_constant_identifier_names,prefer_final_fields,return_of_invalid_type,unnecessary_const,unnecessary_import,unnecessary_this,unused_import,unused_shown_name

import 'dart:core' as $core;
import 'dart:convert' as $convert;
import 'dart:typed_data' as $typed_data;
@$core.Deprecated('Use signMessageWithPrivateKeyRequestDescriptor instead')
const SignMessageWithPrivateKeyRequest$json = const {
  '1': 'SignMessageWithPrivateKeyRequest',
  '2': const [
    const {'1': 'message', '3': 1, '4': 1, '5': 9, '10': 'message'},
    const {'1': 'private_key', '3': 2, '4': 1, '5': 9, '10': 'privateKey'},
  ],
};

/// Descriptor for `SignMessageWithPrivateKeyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List signMessageWithPrivateKeyRequestDescriptor = $convert.base64Decode('CiBTaWduTWVzc2FnZVdpdGhQcml2YXRlS2V5UmVxdWVzdBIYCgdtZXNzYWdlGAEgASgJUgdtZXNzYWdlEh8KC3ByaXZhdGVfa2V5GAIgASgJUgpwcml2YXRlS2V5');
@$core.Deprecated('Use signMessageWithPrivateKeyResponseDescriptor instead')
const SignMessageWithPrivateKeyResponse$json = const {
  '1': 'SignMessageWithPrivateKeyResponse',
  '2': const [
    const {'1': 'signature', '3': 1, '4': 1, '5': 9, '10': 'signature'},
  ],
};

/// Descriptor for `SignMessageWithPrivateKeyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List signMessageWithPrivateKeyResponseDescriptor = $convert.base64Decode('CiFTaWduTWVzc2FnZVdpdGhQcml2YXRlS2V5UmVzcG9uc2USHAoJc2lnbmF0dXJlGAEgASgJUglzaWduYXR1cmU=');
@$core.Deprecated('Use verifyMessageRequestDescriptor instead')
const VerifyMessageRequest$json = const {
  '1': 'VerifyMessageRequest',
  '2': const [
    const {'1': 'message', '3': 1, '4': 1, '5': 9, '10': 'message'},
    const {'1': 'signature', '3': 2, '4': 1, '5': 9, '10': 'signature'},
    const {'1': 'public_key', '3': 3, '4': 1, '5': 9, '10': 'publicKey'},
  ],
};

/// Descriptor for `VerifyMessageRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List verifyMessageRequestDescriptor = $convert.base64Decode('ChRWZXJpZnlNZXNzYWdlUmVxdWVzdBIYCgdtZXNzYWdlGAEgASgJUgdtZXNzYWdlEhwKCXNpZ25hdHVyZRgCIAEoCVIJc2lnbmF0dXJlEh0KCnB1YmxpY19rZXkYAyABKAlSCXB1YmxpY0tleQ==');
@$core.Deprecated('Use verifyMessageResponseDescriptor instead')
const VerifyMessageResponse$json = const {
  '1': 'VerifyMessageResponse',
  '2': const [
    const {'1': 'is_valid', '3': 1, '4': 1, '5': 8, '10': 'isValid'},
  ],
};

/// Descriptor for `VerifyMessageResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List verifyMessageResponseDescriptor = $convert.base64Decode('ChVWZXJpZnlNZXNzYWdlUmVzcG9uc2USGQoIaXNfdmFsaWQYASABKAhSB2lzVmFsaWQ=');
const $core.Map<$core.String, $core.dynamic> UtilServiceBase$json = const {
  '1': 'Util',
  '2': const [
    const {'1': 'SignMessageWithPrivateKey', '2': '.pactus.SignMessageWithPrivateKeyRequest', '3': '.pactus.SignMessageWithPrivateKeyResponse'},
    const {'1': 'VerifyMessage', '2': '.pactus.VerifyMessageRequest', '3': '.pactus.VerifyMessageResponse'},
  ],
};

@$core.Deprecated('Use utilServiceDescriptor instead')
const $core.Map<$core.String, $core.Map<$core.String, $core.dynamic>> UtilServiceBase$messageJson = const {
  '.pactus.SignMessageWithPrivateKeyRequest': SignMessageWithPrivateKeyRequest$json,
  '.pactus.SignMessageWithPrivateKeyResponse': SignMessageWithPrivateKeyResponse$json,
  '.pactus.VerifyMessageRequest': VerifyMessageRequest$json,
  '.pactus.VerifyMessageResponse': VerifyMessageResponse$json,
};

/// Descriptor for `Util`. Decode as a `google.protobuf.ServiceDescriptorProto`.
final $typed_data.Uint8List utilServiceDescriptor = $convert.base64Decode('CgRVdGlsEnAKGVNpZ25NZXNzYWdlV2l0aFByaXZhdGVLZXkSKC5wYWN0dXMuU2lnbk1lc3NhZ2VXaXRoUHJpdmF0ZUtleVJlcXVlc3QaKS5wYWN0dXMuU2lnbk1lc3NhZ2VXaXRoUHJpdmF0ZUtleVJlc3BvbnNlEkwKDVZlcmlmeU1lc3NhZ2USHC5wYWN0dXMuVmVyaWZ5TWVzc2FnZVJlcXVlc3QaHS5wYWN0dXMuVmVyaWZ5TWVzc2FnZVJlc3BvbnNl');
