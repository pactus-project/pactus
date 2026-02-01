// This is a generated file - do not edit.
//
// Generated from wallet.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

/// AddressType defines different types of blockchain addresses.
class AddressType extends $pb.ProtobufEnum {
  /// Treasury address type.
  /// Should not be used to generate new addresses.
  static const AddressType ADDRESS_TYPE_TREASURY =
      AddressType._(0, _omitEnumNames ? '' : 'ADDRESS_TYPE_TREASURY');

  /// Validator address type used for validator nodes.
  static const AddressType ADDRESS_TYPE_VALIDATOR =
      AddressType._(1, _omitEnumNames ? '' : 'ADDRESS_TYPE_VALIDATOR');

  /// Account address type with BLS signature scheme.
  static const AddressType ADDRESS_TYPE_BLS_ACCOUNT =
      AddressType._(2, _omitEnumNames ? '' : 'ADDRESS_TYPE_BLS_ACCOUNT');

  /// Account address type with Ed25519 signature scheme.
  /// Note: Generating a new Ed25519 address requires the wallet password.
  static const AddressType ADDRESS_TYPE_ED25519_ACCOUNT =
      AddressType._(3, _omitEnumNames ? '' : 'ADDRESS_TYPE_ED25519_ACCOUNT');

  static const $core.List<AddressType> values = <AddressType>[
    ADDRESS_TYPE_TREASURY,
    ADDRESS_TYPE_VALIDATOR,
    ADDRESS_TYPE_BLS_ACCOUNT,
    ADDRESS_TYPE_ED25519_ACCOUNT,
  ];

  static final $core.List<AddressType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static AddressType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const AddressType._(super.value, super.name);
}

/// TxDirection indicates the direction of a transaction relative to the wallet.
class TxDirection extends $pb.ProtobufEnum {
  /// include both incoming and outgoing transactions.
  static const TxDirection TX_DIRECTION_ANY =
      TxDirection._(0, _omitEnumNames ? '' : 'TX_DIRECTION_ANY');

  /// Include only incoming transactions where the wallet receives funds.
  static const TxDirection TX_DIRECTION_INCOMING =
      TxDirection._(1, _omitEnumNames ? '' : 'TX_DIRECTION_INCOMING');

  /// Include only outgoing transactions where the wallet sends funds.
  static const TxDirection TX_DIRECTION_OUTGOING =
      TxDirection._(2, _omitEnumNames ? '' : 'TX_DIRECTION_OUTGOING');

  static const $core.List<TxDirection> values = <TxDirection>[
    TX_DIRECTION_ANY,
    TX_DIRECTION_INCOMING,
    TX_DIRECTION_OUTGOING,
  ];

  static final $core.List<TxDirection?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static TxDirection? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const TxDirection._(super.value, super.name);
}

/// TransactionStatus defines the status of a transaction.
class TransactionStatus extends $pb.ProtobufEnum {
  /// Pending status for transactions in the mempool.
  static const TransactionStatus TRANSACTION_STATUS_PENDING =
      TransactionStatus._(
          0, _omitEnumNames ? '' : 'TRANSACTION_STATUS_PENDING');

  /// Confirmed status for transactions included in a block.
  static const TransactionStatus TRANSACTION_STATUS_CONFIRMED =
      TransactionStatus._(
          1, _omitEnumNames ? '' : 'TRANSACTION_STATUS_CONFIRMED');

  /// Failed status for transactions that were not successful.
  static const TransactionStatus TRANSACTION_STATUS_FAILED =
      TransactionStatus._(
          -1, _omitEnumNames ? '' : 'TRANSACTION_STATUS_FAILED');

  static const $core.List<TransactionStatus> values = <TransactionStatus>[
    TRANSACTION_STATUS_PENDING,
    TRANSACTION_STATUS_CONFIRMED,
    TRANSACTION_STATUS_FAILED,
  ];

  static final $core.Map<$core.int, TransactionStatus> _byValue =
      $pb.ProtobufEnum.initByValue(values);
  static TransactionStatus? valueOf($core.int value) => _byValue[value];

  const TransactionStatus._(super.value, super.name);
}

const $core.bool _omitEnumNames =
    $core.bool.fromEnvironment('protobuf.omit_enum_names');
