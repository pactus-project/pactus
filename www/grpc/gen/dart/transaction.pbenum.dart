//
//  Generated code. Do not modify.
//  source: transaction.proto
//
// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_final_fields
// ignore_for_file: unnecessary_import, unnecessary_this, unused_import

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

/// Enumeration for different types of transaction payloads.
class PayloadType extends $pb.ProtobufEnum {
  /// Unspecified payload type.
  static const PayloadType PAYLOAD_TYPE_UNSPECIFIED = PayloadType._(0, _omitEnumNames ? '' : 'PAYLOAD_TYPE_UNSPECIFIED');
  /// Transfer payload type.
  static const PayloadType PAYLOAD_TYPE_TRANSFER = PayloadType._(1, _omitEnumNames ? '' : 'PAYLOAD_TYPE_TRANSFER');
  /// Bond payload type.
  static const PayloadType PAYLOAD_TYPE_BOND = PayloadType._(2, _omitEnumNames ? '' : 'PAYLOAD_TYPE_BOND');
  /// Sortition payload type.
  static const PayloadType PAYLOAD_TYPE_SORTITION = PayloadType._(3, _omitEnumNames ? '' : 'PAYLOAD_TYPE_SORTITION');
  /// Unbond payload type.
  static const PayloadType PAYLOAD_TYPE_UNBOND = PayloadType._(4, _omitEnumNames ? '' : 'PAYLOAD_TYPE_UNBOND');
  /// Withdraw payload type.
  static const PayloadType PAYLOAD_TYPE_WITHDRAW = PayloadType._(5, _omitEnumNames ? '' : 'PAYLOAD_TYPE_WITHDRAW');
  /// Batch transfer payload type.
  static const PayloadType PAYLOAD_TYPE_BATCH_TRANSFER = PayloadType._(6, _omitEnumNames ? '' : 'PAYLOAD_TYPE_BATCH_TRANSFER');

  static const $core.List<PayloadType> values = <PayloadType> [
    PAYLOAD_TYPE_UNSPECIFIED,
    PAYLOAD_TYPE_TRANSFER,
    PAYLOAD_TYPE_BOND,
    PAYLOAD_TYPE_SORTITION,
    PAYLOAD_TYPE_UNBOND,
    PAYLOAD_TYPE_WITHDRAW,
    PAYLOAD_TYPE_BATCH_TRANSFER,
  ];

  static final $core.Map<$core.int, PayloadType> _byValue = $pb.ProtobufEnum.initByValue(values);
  static PayloadType? valueOf($core.int value) => _byValue[value];

  const PayloadType._(super.v, super.n);
}

/// Enumeration for verbosity levels when requesting transaction details.
class TransactionVerbosity extends $pb.ProtobufEnum {
  /// Request transaction data only.
  static const TransactionVerbosity TRANSACTION_VERBOSITY_DATA = TransactionVerbosity._(0, _omitEnumNames ? '' : 'TRANSACTION_VERBOSITY_DATA');
  /// Request detailed transaction information.
  static const TransactionVerbosity TRANSACTION_VERBOSITY_INFO = TransactionVerbosity._(1, _omitEnumNames ? '' : 'TRANSACTION_VERBOSITY_INFO');

  static const $core.List<TransactionVerbosity> values = <TransactionVerbosity> [
    TRANSACTION_VERBOSITY_DATA,
    TRANSACTION_VERBOSITY_INFO,
  ];

  static final $core.Map<$core.int, TransactionVerbosity> _byValue = $pb.ProtobufEnum.initByValue(values);
  static TransactionVerbosity? valueOf($core.int value) => _byValue[value];

  const TransactionVerbosity._(super.v, super.n);
}


const _omitEnumNames = $core.bool.fromEnvironment('protobuf.omit_enum_names');
