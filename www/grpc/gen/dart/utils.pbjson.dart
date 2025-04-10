//
//  Generated code. Do not modify.
//  source: utils.proto
//
// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_final_fields
// ignore_for_file: unnecessary_import, unnecessary_this, unused_import

import 'dart:convert' as $convert;
import 'dart:core' as $core;
import 'dart:typed_data' as $typed_data;

@$core.Deprecated('Use signMessageWithPrivateKeyRequestDescriptor instead')
const SignMessageWithPrivateKeyRequest$json = {
  '1': 'SignMessageWithPrivateKeyRequest',
  '2': [
    {'1': 'private_key', '3': 1, '4': 1, '5': 9, '10': 'privateKey'},
    {'1': 'message', '3': 2, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `SignMessageWithPrivateKeyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List signMessageWithPrivateKeyRequestDescriptor = $convert.base64Decode(
    'CiBTaWduTWVzc2FnZVdpdGhQcml2YXRlS2V5UmVxdWVzdBIfCgtwcml2YXRlX2tleRgBIAEoCV'
    'IKcHJpdmF0ZUtleRIYCgdtZXNzYWdlGAIgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use signMessageWithPrivateKeyResponseDescriptor instead')
const SignMessageWithPrivateKeyResponse$json = {
  '1': 'SignMessageWithPrivateKeyResponse',
  '2': [
    {'1': 'signature', '3': 1, '4': 1, '5': 9, '10': 'signature'},
  ],
};

/// Descriptor for `SignMessageWithPrivateKeyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List signMessageWithPrivateKeyResponseDescriptor = $convert.base64Decode(
    'CiFTaWduTWVzc2FnZVdpdGhQcml2YXRlS2V5UmVzcG9uc2USHAoJc2lnbmF0dXJlGAEgASgJUg'
    'lzaWduYXR1cmU=');

@$core.Deprecated('Use verifyMessageRequestDescriptor instead')
const VerifyMessageRequest$json = {
  '1': 'VerifyMessageRequest',
  '2': [
    {'1': 'message', '3': 1, '4': 1, '5': 9, '10': 'message'},
    {'1': 'signature', '3': 2, '4': 1, '5': 9, '10': 'signature'},
    {'1': 'public_key', '3': 3, '4': 1, '5': 9, '10': 'publicKey'},
  ],
};

/// Descriptor for `VerifyMessageRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List verifyMessageRequestDescriptor = $convert.base64Decode(
    'ChRWZXJpZnlNZXNzYWdlUmVxdWVzdBIYCgdtZXNzYWdlGAEgASgJUgdtZXNzYWdlEhwKCXNpZ2'
    '5hdHVyZRgCIAEoCVIJc2lnbmF0dXJlEh0KCnB1YmxpY19rZXkYAyABKAlSCXB1YmxpY0tleQ==');

@$core.Deprecated('Use verifyMessageResponseDescriptor instead')
const VerifyMessageResponse$json = {
  '1': 'VerifyMessageResponse',
  '2': [
    {'1': 'is_valid', '3': 1, '4': 1, '5': 8, '10': 'isValid'},
  ],
};

/// Descriptor for `VerifyMessageResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List verifyMessageResponseDescriptor = $convert.base64Decode(
    'ChVWZXJpZnlNZXNzYWdlUmVzcG9uc2USGQoIaXNfdmFsaWQYASABKAhSB2lzVmFsaWQ=');

@$core.Deprecated('Use publicKeyAggregationRequestDescriptor instead')
const PublicKeyAggregationRequest$json = {
  '1': 'PublicKeyAggregationRequest',
  '2': [
    {'1': 'public_keys', '3': 1, '4': 3, '5': 9, '10': 'publicKeys'},
  ],
};

/// Descriptor for `PublicKeyAggregationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List publicKeyAggregationRequestDescriptor = $convert.base64Decode(
    'ChtQdWJsaWNLZXlBZ2dyZWdhdGlvblJlcXVlc3QSHwoLcHVibGljX2tleXMYASADKAlSCnB1Ym'
    'xpY0tleXM=');

@$core.Deprecated('Use publicKeyAggregationResponseDescriptor instead')
const PublicKeyAggregationResponse$json = {
  '1': 'PublicKeyAggregationResponse',
  '2': [
    {'1': 'public_key', '3': 1, '4': 1, '5': 9, '10': 'publicKey'},
    {'1': 'address', '3': 2, '4': 1, '5': 9, '10': 'address'},
  ],
};

/// Descriptor for `PublicKeyAggregationResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List publicKeyAggregationResponseDescriptor = $convert.base64Decode(
    'ChxQdWJsaWNLZXlBZ2dyZWdhdGlvblJlc3BvbnNlEh0KCnB1YmxpY19rZXkYASABKAlSCXB1Ym'
    'xpY0tleRIYCgdhZGRyZXNzGAIgASgJUgdhZGRyZXNz');

@$core.Deprecated('Use signatureAggregationRequestDescriptor instead')
const SignatureAggregationRequest$json = {
  '1': 'SignatureAggregationRequest',
  '2': [
    {'1': 'signatures', '3': 1, '4': 3, '5': 9, '10': 'signatures'},
  ],
};

/// Descriptor for `SignatureAggregationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List signatureAggregationRequestDescriptor = $convert.base64Decode(
    'ChtTaWduYXR1cmVBZ2dyZWdhdGlvblJlcXVlc3QSHgoKc2lnbmF0dXJlcxgBIAMoCVIKc2lnbm'
    'F0dXJlcw==');

@$core.Deprecated('Use signatureAggregationResponseDescriptor instead')
const SignatureAggregationResponse$json = {
  '1': 'SignatureAggregationResponse',
  '2': [
    {'1': 'signature', '3': 1, '4': 1, '5': 9, '10': 'signature'},
  ],
};

/// Descriptor for `SignatureAggregationResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List signatureAggregationResponseDescriptor = $convert.base64Decode(
    'ChxTaWduYXR1cmVBZ2dyZWdhdGlvblJlc3BvbnNlEhwKCXNpZ25hdHVyZRgBIAEoCVIJc2lnbm'
    'F0dXJl');

const $core.Map<$core.String, $core.dynamic> UtilsServiceBase$json = {
  '1': 'Utils',
  '2': [
    {'1': 'SignMessageWithPrivateKey', '2': '.pactus.SignMessageWithPrivateKeyRequest', '3': '.pactus.SignMessageWithPrivateKeyResponse'},
    {'1': 'VerifyMessage', '2': '.pactus.VerifyMessageRequest', '3': '.pactus.VerifyMessageResponse'},
    {'1': 'PublicKeyAggregation', '2': '.pactus.PublicKeyAggregationRequest', '3': '.pactus.PublicKeyAggregationResponse'},
    {'1': 'SignatureAggregation', '2': '.pactus.SignatureAggregationRequest', '3': '.pactus.SignatureAggregationResponse'},
  ],
};

@$core.Deprecated('Use utilsServiceDescriptor instead')
const $core.Map<$core.String, $core.Map<$core.String, $core.dynamic>> UtilsServiceBase$messageJson = {
  '.pactus.SignMessageWithPrivateKeyRequest': SignMessageWithPrivateKeyRequest$json,
  '.pactus.SignMessageWithPrivateKeyResponse': SignMessageWithPrivateKeyResponse$json,
  '.pactus.VerifyMessageRequest': VerifyMessageRequest$json,
  '.pactus.VerifyMessageResponse': VerifyMessageResponse$json,
  '.pactus.PublicKeyAggregationRequest': PublicKeyAggregationRequest$json,
  '.pactus.PublicKeyAggregationResponse': PublicKeyAggregationResponse$json,
  '.pactus.SignatureAggregationRequest': SignatureAggregationRequest$json,
  '.pactus.SignatureAggregationResponse': SignatureAggregationResponse$json,
};

/// Descriptor for `Utils`. Decode as a `google.protobuf.ServiceDescriptorProto`.
final $typed_data.Uint8List utilsServiceDescriptor = $convert.base64Decode(
    'CgVVdGlscxJwChlTaWduTWVzc2FnZVdpdGhQcml2YXRlS2V5EigucGFjdHVzLlNpZ25NZXNzYW'
    'dlV2l0aFByaXZhdGVLZXlSZXF1ZXN0GikucGFjdHVzLlNpZ25NZXNzYWdlV2l0aFByaXZhdGVL'
    'ZXlSZXNwb25zZRJMCg1WZXJpZnlNZXNzYWdlEhwucGFjdHVzLlZlcmlmeU1lc3NhZ2VSZXF1ZX'
    'N0Gh0ucGFjdHVzLlZlcmlmeU1lc3NhZ2VSZXNwb25zZRJhChRQdWJsaWNLZXlBZ2dyZWdhdGlv'
    'bhIjLnBhY3R1cy5QdWJsaWNLZXlBZ2dyZWdhdGlvblJlcXVlc3QaJC5wYWN0dXMuUHVibGljS2'
    'V5QWdncmVnYXRpb25SZXNwb25zZRJhChRTaWduYXR1cmVBZ2dyZWdhdGlvbhIjLnBhY3R1cy5T'
    'aWduYXR1cmVBZ2dyZWdhdGlvblJlcXVlc3QaJC5wYWN0dXMuU2lnbmF0dXJlQWdncmVnYXRpb2'
    '5SZXNwb25zZQ==');

