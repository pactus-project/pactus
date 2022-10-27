///
//  Generated code. Do not modify.
//  source: transaction.proto
//
// @dart = 2.12
// ignore_for_file: annotate_overrides,camel_case_types,constant_identifier_names,directives_ordering,library_prefixes,non_constant_identifier_names,prefer_final_fields,return_of_invalid_type,unnecessary_const,unnecessary_import,unnecessary_this,unused_import,unused_shown_name

// ignore_for_file: UNDEFINED_SHOWN_NAME
import 'dart:core' as $core;
import 'package:protobuf/protobuf.dart' as $pb;

class PayloadType extends $pb.ProtobufEnum {
  static const PayloadType UNKNOWN = PayloadType._(0, const $core.bool.fromEnvironment('protobuf.omit_enum_names') ? '' : 'UNKNOWN');
  static const PayloadType SEND_PAYLOAD = PayloadType._(1, const $core.bool.fromEnvironment('protobuf.omit_enum_names') ? '' : 'SEND_PAYLOAD');
  static const PayloadType BOND_PAYLOAD = PayloadType._(2, const $core.bool.fromEnvironment('protobuf.omit_enum_names') ? '' : 'BOND_PAYLOAD');
  static const PayloadType SORTITION_PAYLOAD = PayloadType._(3, const $core.bool.fromEnvironment('protobuf.omit_enum_names') ? '' : 'SORTITION_PAYLOAD');
  static const PayloadType UNBOND_PAYLOAD = PayloadType._(4, const $core.bool.fromEnvironment('protobuf.omit_enum_names') ? '' : 'UNBOND_PAYLOAD');
  static const PayloadType WITHDRAW_PAYLOAD = PayloadType._(5, const $core.bool.fromEnvironment('protobuf.omit_enum_names') ? '' : 'WITHDRAW_PAYLOAD');

  static const $core.List<PayloadType> values = <PayloadType> [
    UNKNOWN,
    SEND_PAYLOAD,
    BOND_PAYLOAD,
    SORTITION_PAYLOAD,
    UNBOND_PAYLOAD,
    WITHDRAW_PAYLOAD,
  ];

  static final $core.Map<$core.int, PayloadType> _byValue = $pb.ProtobufEnum.initByValue(values);
  static PayloadType? valueOf($core.int value) => _byValue[value];

  const PayloadType._($core.int v, $core.String n) : super(v, n);
}

class TransactionVerbosity extends $pb.ProtobufEnum {
  static const TransactionVerbosity TRANSACTION_DATA = TransactionVerbosity._(0, const $core.bool.fromEnvironment('protobuf.omit_enum_names') ? '' : 'TRANSACTION_DATA');
  static const TransactionVerbosity TRANSACTION_INFO = TransactionVerbosity._(1, const $core.bool.fromEnvironment('protobuf.omit_enum_names') ? '' : 'TRANSACTION_INFO');

  static const $core.List<TransactionVerbosity> values = <TransactionVerbosity> [
    TRANSACTION_DATA,
    TRANSACTION_INFO,
  ];

  static final $core.Map<$core.int, TransactionVerbosity> _byValue = $pb.ProtobufEnum.initByValue(values);
  static TransactionVerbosity? valueOf($core.int value) => _byValue[value];

  const TransactionVerbosity._($core.int v, $core.String n) : super(v, n);
}

