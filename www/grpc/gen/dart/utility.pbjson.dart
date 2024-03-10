///
//  Generated code. Do not modify.
//  source: utility.proto
//
// @dart = 2.12
// ignore_for_file: annotate_overrides,camel_case_types,constant_identifier_names,deprecated_member_use_from_same_package,directives_ordering,library_prefixes,non_constant_identifier_names,prefer_final_fields,return_of_invalid_type,unnecessary_const,unnecessary_import,unnecessary_this,unused_import,unused_shown_name

import 'dart:core' as $core;
import 'dart:convert' as $convert;
import 'dart:typed_data' as $typed_data;
@$core.Deprecated('Use calculateFeeRequestDescriptor instead')
const CalculateFeeRequest$json = const {
  '1': 'CalculateFeeRequest',
  '2': const [
    const {'1': 'amount', '3': 1, '4': 1, '5': 3, '10': 'amount'},
    const {'1': 'payload_type', '3': 2, '4': 1, '5': 14, '6': '.pactus.PayloadType', '10': 'payloadType'},
    const {'1': 'fixed_amount', '3': 3, '4': 1, '5': 8, '10': 'fixedAmount'},
  ],
};

/// Descriptor for `CalculateFeeRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List calculateFeeRequestDescriptor = $convert.base64Decode('ChNDYWxjdWxhdGVGZWVSZXF1ZXN0EhYKBmFtb3VudBgBIAEoA1IGYW1vdW50EjYKDHBheWxvYWRfdHlwZRgCIAEoDjITLnBhY3R1cy5QYXlsb2FkVHlwZVILcGF5bG9hZFR5cGUSIQoMZml4ZWRfYW1vdW50GAMgASgIUgtmaXhlZEFtb3VudA==');
@$core.Deprecated('Use calculateFeeResponseDescriptor instead')
const CalculateFeeResponse$json = const {
  '1': 'CalculateFeeResponse',
  '2': const [
    const {'1': 'amount', '3': 1, '4': 1, '5': 3, '10': 'amount'},
    const {'1': 'fee', '3': 2, '4': 1, '5': 3, '10': 'fee'},
  ],
};

/// Descriptor for `CalculateFeeResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List calculateFeeResponseDescriptor = $convert.base64Decode('ChRDYWxjdWxhdGVGZWVSZXNwb25zZRIWCgZhbW91bnQYASABKANSBmFtb3VudBIQCgNmZWUYAiABKANSA2ZlZQ==');
const $core.Map<$core.String, $core.dynamic> UtilityServiceBase$json = const {
  '1': 'Utility',
  '2': const [
    const {'1': 'CalculateFee', '2': '.pactus.CalculateFeeRequest', '3': '.pactus.CalculateFeeResponse'},
  ],
};

@$core.Deprecated('Use utilityServiceDescriptor instead')
const $core.Map<$core.String, $core.Map<$core.String, $core.dynamic>> UtilityServiceBase$messageJson = const {
  '.pactus.CalculateFeeRequest': CalculateFeeRequest$json,
  '.pactus.CalculateFeeResponse': CalculateFeeResponse$json,
};

/// Descriptor for `Utility`. Decode as a `google.protobuf.ServiceDescriptorProto`.
final $typed_data.Uint8List utilityServiceDescriptor = $convert.base64Decode('CgdVdGlsaXR5EkkKDENhbGN1bGF0ZUZlZRIbLnBhY3R1cy5DYWxjdWxhdGVGZWVSZXF1ZXN0GhwucGFjdHVzLkNhbGN1bGF0ZUZlZVJlc3BvbnNl');
