//
//  Generated code. Do not modify.
//  source: wallet.proto
//
// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_final_fields
// ignore_for_file: unnecessary_import, unnecessary_this, unused_import

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

/// AddressType defines different types of blockchain addresses.
class AddressType extends $pb.ProtobufEnum {
  /// Treasury address type.
  /// Should not be used to generate new addresses.
  /// buf:lint:ignore ENUM_ZERO_VALUE_SUFFIX
  static const AddressType ADDRESS_TYPE_TREASURY = AddressType._(0, _omitEnumNames ? '' : 'ADDRESS_TYPE_TREASURY');
  /// Validator address type used for validator nodes.
  static const AddressType ADDRESS_TYPE_VALIDATOR = AddressType._(1, _omitEnumNames ? '' : 'ADDRESS_TYPE_VALIDATOR');
  /// Account address type with BLS signature scheme.
  static const AddressType ADDRESS_TYPE_BLS_ACCOUNT = AddressType._(2, _omitEnumNames ? '' : 'ADDRESS_TYPE_BLS_ACCOUNT');
  /// Account address type with Ed25519 signature scheme.
  /// Note: Generating a new Ed25519 address requires the wallet password.
  static const AddressType ADDRESS_TYPE_ED25519_ACCOUNT = AddressType._(3, _omitEnumNames ? '' : 'ADDRESS_TYPE_ED25519_ACCOUNT');

  static const $core.List<AddressType> values = <AddressType> [
    ADDRESS_TYPE_TREASURY,
    ADDRESS_TYPE_VALIDATOR,
    ADDRESS_TYPE_BLS_ACCOUNT,
    ADDRESS_TYPE_ED25519_ACCOUNT,
  ];

  static final $core.Map<$core.int, AddressType> _byValue = $pb.ProtobufEnum.initByValue(values);
  static AddressType? valueOf($core.int value) => _byValue[value];

  const AddressType._(super.v, super.n);
}


const _omitEnumNames = $core.bool.fromEnvironment('protobuf.omit_enum_names');
