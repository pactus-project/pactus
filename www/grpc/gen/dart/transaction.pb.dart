// This is a generated file - do not edit.
//
// Generated from transaction.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:async' as $async;
import 'dart:core' as $core;

import 'package:fixnum/fixnum.dart' as $fixnum;
import 'package:protobuf/protobuf.dart' as $pb;

import 'transaction.pbenum.dart';

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

export 'transaction.pbenum.dart';

/// Request message for retrieving transaction details.
class GetTransactionRequest extends $pb.GeneratedMessage {
  factory GetTransactionRequest({
    $core.String? id,
    TransactionVerbosity? verbosity,
  }) {
    final result = create();
    if (id != null) result.id = id;
    if (verbosity != null) result.verbosity = verbosity;
    return result;
  }

  GetTransactionRequest._();

  factory GetTransactionRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetTransactionRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetTransactionRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aE<TransactionVerbosity>(2, _omitFieldNames ? '' : 'verbosity',
        enumValues: TransactionVerbosity.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetTransactionRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetTransactionRequest copyWith(
          void Function(GetTransactionRequest) updates) =>
      super.copyWith((message) => updates(message as GetTransactionRequest))
          as GetTransactionRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetTransactionRequest create() => GetTransactionRequest._();
  @$core.override
  GetTransactionRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetTransactionRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetTransactionRequest>(create);
  static GetTransactionRequest? _defaultInstance;

  /// The unique ID of the transaction to retrieve.
  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  /// The verbosity level for transaction details.
  @$pb.TagNumber(2)
  TransactionVerbosity get verbosity => $_getN(1);
  @$pb.TagNumber(2)
  set verbosity(TransactionVerbosity value) => $_setField(2, value);
  @$pb.TagNumber(2)
  $core.bool hasVerbosity() => $_has(1);
  @$pb.TagNumber(2)
  void clearVerbosity() => $_clearField(2);
}

/// Response message contains details of a transaction.
class GetTransactionResponse extends $pb.GeneratedMessage {
  factory GetTransactionResponse({
    $core.int? blockHeight,
    $core.int? blockTime,
    TransactionInfo? transaction,
  }) {
    final result = create();
    if (blockHeight != null) result.blockHeight = blockHeight;
    if (blockTime != null) result.blockTime = blockTime;
    if (transaction != null) result.transaction = transaction;
    return result;
  }

  GetTransactionResponse._();

  factory GetTransactionResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetTransactionResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetTransactionResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aI(1, _omitFieldNames ? '' : 'blockHeight',
        fieldType: $pb.PbFieldType.OU3)
    ..aI(2, _omitFieldNames ? '' : 'blockTime', fieldType: $pb.PbFieldType.OU3)
    ..aOM<TransactionInfo>(3, _omitFieldNames ? '' : 'transaction',
        subBuilder: TransactionInfo.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetTransactionResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetTransactionResponse copyWith(
          void Function(GetTransactionResponse) updates) =>
      super.copyWith((message) => updates(message as GetTransactionResponse))
          as GetTransactionResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetTransactionResponse create() => GetTransactionResponse._();
  @$core.override
  GetTransactionResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetTransactionResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetTransactionResponse>(create);
  static GetTransactionResponse? _defaultInstance;

  /// The height of the block containing the transaction.
  @$pb.TagNumber(1)
  $core.int get blockHeight => $_getIZ(0);
  @$pb.TagNumber(1)
  set blockHeight($core.int value) => $_setUnsignedInt32(0, value);
  @$pb.TagNumber(1)
  $core.bool hasBlockHeight() => $_has(0);
  @$pb.TagNumber(1)
  void clearBlockHeight() => $_clearField(1);

  /// The UNIX timestamp of the block containing the transaction.
  @$pb.TagNumber(2)
  $core.int get blockTime => $_getIZ(1);
  @$pb.TagNumber(2)
  set blockTime($core.int value) => $_setUnsignedInt32(1, value);
  @$pb.TagNumber(2)
  $core.bool hasBlockTime() => $_has(1);
  @$pb.TagNumber(2)
  void clearBlockTime() => $_clearField(2);

  /// Detailed information about the transaction.
  @$pb.TagNumber(3)
  TransactionInfo get transaction => $_getN(2);
  @$pb.TagNumber(3)
  set transaction(TransactionInfo value) => $_setField(3, value);
  @$pb.TagNumber(3)
  $core.bool hasTransaction() => $_has(2);
  @$pb.TagNumber(3)
  void clearTransaction() => $_clearField(3);
  @$pb.TagNumber(3)
  TransactionInfo ensureTransaction() => $_ensure(2);
}

/// Request message for calculating transaction fee.
class CalculateFeeRequest extends $pb.GeneratedMessage {
  factory CalculateFeeRequest({
    $fixnum.Int64? amount,
    PayloadType? payloadType,
    $core.bool? fixedAmount,
  }) {
    final result = create();
    if (amount != null) result.amount = amount;
    if (payloadType != null) result.payloadType = payloadType;
    if (fixedAmount != null) result.fixedAmount = fixedAmount;
    return result;
  }

  CalculateFeeRequest._();

  factory CalculateFeeRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CalculateFeeRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CalculateFeeRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aInt64(1, _omitFieldNames ? '' : 'amount')
    ..aE<PayloadType>(2, _omitFieldNames ? '' : 'payloadType',
        enumValues: PayloadType.values)
    ..aOB(3, _omitFieldNames ? '' : 'fixedAmount')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CalculateFeeRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CalculateFeeRequest copyWith(void Function(CalculateFeeRequest) updates) =>
      super.copyWith((message) => updates(message as CalculateFeeRequest))
          as CalculateFeeRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CalculateFeeRequest create() => CalculateFeeRequest._();
  @$core.override
  CalculateFeeRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CalculateFeeRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CalculateFeeRequest>(create);
  static CalculateFeeRequest? _defaultInstance;

  /// The amount involved in the transaction, specified in NanoPAC.
  @$pb.TagNumber(1)
  $fixnum.Int64 get amount => $_getI64(0);
  @$pb.TagNumber(1)
  set amount($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(1)
  $core.bool hasAmount() => $_has(0);
  @$pb.TagNumber(1)
  void clearAmount() => $_clearField(1);

  /// The type of transaction payload.
  @$pb.TagNumber(2)
  PayloadType get payloadType => $_getN(1);
  @$pb.TagNumber(2)
  set payloadType(PayloadType value) => $_setField(2, value);
  @$pb.TagNumber(2)
  $core.bool hasPayloadType() => $_has(1);
  @$pb.TagNumber(2)
  void clearPayloadType() => $_clearField(2);

  /// Indicates if the amount should be fixed and include the fee.
  @$pb.TagNumber(3)
  $core.bool get fixedAmount => $_getBF(2);
  @$pb.TagNumber(3)
  set fixedAmount($core.bool value) => $_setBool(2, value);
  @$pb.TagNumber(3)
  $core.bool hasFixedAmount() => $_has(2);
  @$pb.TagNumber(3)
  void clearFixedAmount() => $_clearField(3);
}

/// Response message contains the calculated transaction fee.
class CalculateFeeResponse extends $pb.GeneratedMessage {
  factory CalculateFeeResponse({
    $fixnum.Int64? amount,
    $fixnum.Int64? fee,
  }) {
    final result = create();
    if (amount != null) result.amount = amount;
    if (fee != null) result.fee = fee;
    return result;
  }

  CalculateFeeResponse._();

  factory CalculateFeeResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CalculateFeeResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CalculateFeeResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aInt64(1, _omitFieldNames ? '' : 'amount')
    ..aInt64(2, _omitFieldNames ? '' : 'fee')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CalculateFeeResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CalculateFeeResponse copyWith(void Function(CalculateFeeResponse) updates) =>
      super.copyWith((message) => updates(message as CalculateFeeResponse))
          as CalculateFeeResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CalculateFeeResponse create() => CalculateFeeResponse._();
  @$core.override
  CalculateFeeResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CalculateFeeResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CalculateFeeResponse>(create);
  static CalculateFeeResponse? _defaultInstance;

  /// The calculated amount in NanoPAC.
  @$pb.TagNumber(1)
  $fixnum.Int64 get amount => $_getI64(0);
  @$pb.TagNumber(1)
  set amount($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(1)
  $core.bool hasAmount() => $_has(0);
  @$pb.TagNumber(1)
  void clearAmount() => $_clearField(1);

  /// The calculated transaction fee in NanoPAC.
  @$pb.TagNumber(2)
  $fixnum.Int64 get fee => $_getI64(1);
  @$pb.TagNumber(2)
  set fee($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(2)
  $core.bool hasFee() => $_has(1);
  @$pb.TagNumber(2)
  void clearFee() => $_clearField(2);
}

/// Request message for broadcasting a signed transaction to the network.
class BroadcastTransactionRequest extends $pb.GeneratedMessage {
  factory BroadcastTransactionRequest({
    $core.String? signedRawTransaction,
  }) {
    final result = create();
    if (signedRawTransaction != null)
      result.signedRawTransaction = signedRawTransaction;
    return result;
  }

  BroadcastTransactionRequest._();

  factory BroadcastTransactionRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory BroadcastTransactionRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'BroadcastTransactionRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'signedRawTransaction')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BroadcastTransactionRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BroadcastTransactionRequest copyWith(
          void Function(BroadcastTransactionRequest) updates) =>
      super.copyWith(
              (message) => updates(message as BroadcastTransactionRequest))
          as BroadcastTransactionRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static BroadcastTransactionRequest create() =>
      BroadcastTransactionRequest._();
  @$core.override
  BroadcastTransactionRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static BroadcastTransactionRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<BroadcastTransactionRequest>(create);
  static BroadcastTransactionRequest? _defaultInstance;

  /// The signed raw transaction data to be broadcasted.
  @$pb.TagNumber(1)
  $core.String get signedRawTransaction => $_getSZ(0);
  @$pb.TagNumber(1)
  set signedRawTransaction($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasSignedRawTransaction() => $_has(0);
  @$pb.TagNumber(1)
  void clearSignedRawTransaction() => $_clearField(1);
}

/// Response message contains the ID of the broadcasted transaction.
class BroadcastTransactionResponse extends $pb.GeneratedMessage {
  factory BroadcastTransactionResponse({
    $core.String? id,
  }) {
    final result = create();
    if (id != null) result.id = id;
    return result;
  }

  BroadcastTransactionResponse._();

  factory BroadcastTransactionResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory BroadcastTransactionResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'BroadcastTransactionResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BroadcastTransactionResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BroadcastTransactionResponse copyWith(
          void Function(BroadcastTransactionResponse) updates) =>
      super.copyWith(
              (message) => updates(message as BroadcastTransactionResponse))
          as BroadcastTransactionResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static BroadcastTransactionResponse create() =>
      BroadcastTransactionResponse._();
  @$core.override
  BroadcastTransactionResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static BroadcastTransactionResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<BroadcastTransactionResponse>(create);
  static BroadcastTransactionResponse? _defaultInstance;

  /// The unique ID of the broadcasted transaction.
  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);
}

/// Request message for retrieving raw details of a transfer transaction.
class GetRawTransferTransactionRequest extends $pb.GeneratedMessage {
  factory GetRawTransferTransactionRequest({
    $core.int? lockTime,
    $core.String? sender,
    $core.String? receiver,
    $fixnum.Int64? amount,
    $fixnum.Int64? fee,
    $core.String? memo,
  }) {
    final result = create();
    if (lockTime != null) result.lockTime = lockTime;
    if (sender != null) result.sender = sender;
    if (receiver != null) result.receiver = receiver;
    if (amount != null) result.amount = amount;
    if (fee != null) result.fee = fee;
    if (memo != null) result.memo = memo;
    return result;
  }

  GetRawTransferTransactionRequest._();

  factory GetRawTransferTransactionRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetRawTransferTransactionRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetRawTransferTransactionRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aI(1, _omitFieldNames ? '' : 'lockTime', fieldType: $pb.PbFieldType.OU3)
    ..aOS(2, _omitFieldNames ? '' : 'sender')
    ..aOS(3, _omitFieldNames ? '' : 'receiver')
    ..aInt64(4, _omitFieldNames ? '' : 'amount')
    ..aInt64(5, _omitFieldNames ? '' : 'fee')
    ..aOS(6, _omitFieldNames ? '' : 'memo')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRawTransferTransactionRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRawTransferTransactionRequest copyWith(
          void Function(GetRawTransferTransactionRequest) updates) =>
      super.copyWith(
              (message) => updates(message as GetRawTransferTransactionRequest))
          as GetRawTransferTransactionRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetRawTransferTransactionRequest create() =>
      GetRawTransferTransactionRequest._();
  @$core.override
  GetRawTransferTransactionRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetRawTransferTransactionRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetRawTransferTransactionRequest>(
          create);
  static GetRawTransferTransactionRequest? _defaultInstance;

  /// The lock time for the transaction. If not set, defaults to the last block height.
  @$pb.TagNumber(1)
  $core.int get lockTime => $_getIZ(0);
  @$pb.TagNumber(1)
  set lockTime($core.int value) => $_setUnsignedInt32(0, value);
  @$pb.TagNumber(1)
  $core.bool hasLockTime() => $_has(0);
  @$pb.TagNumber(1)
  void clearLockTime() => $_clearField(1);

  /// The sender's account address.
  @$pb.TagNumber(2)
  $core.String get sender => $_getSZ(1);
  @$pb.TagNumber(2)
  set sender($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasSender() => $_has(1);
  @$pb.TagNumber(2)
  void clearSender() => $_clearField(2);

  /// The receiver's account address.
  @$pb.TagNumber(3)
  $core.String get receiver => $_getSZ(2);
  @$pb.TagNumber(3)
  set receiver($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasReceiver() => $_has(2);
  @$pb.TagNumber(3)
  void clearReceiver() => $_clearField(3);

  /// The amount to be transferred, specified in NanoPAC. Must be greater than 0.
  @$pb.TagNumber(4)
  $fixnum.Int64 get amount => $_getI64(3);
  @$pb.TagNumber(4)
  set amount($fixnum.Int64 value) => $_setInt64(3, value);
  @$pb.TagNumber(4)
  $core.bool hasAmount() => $_has(3);
  @$pb.TagNumber(4)
  void clearAmount() => $_clearField(4);

  /// The transaction fee in NanoPAC. If not set, it is set to the estimated fee.
  @$pb.TagNumber(5)
  $fixnum.Int64 get fee => $_getI64(4);
  @$pb.TagNumber(5)
  set fee($fixnum.Int64 value) => $_setInt64(4, value);
  @$pb.TagNumber(5)
  $core.bool hasFee() => $_has(4);
  @$pb.TagNumber(5)
  void clearFee() => $_clearField(5);

  /// A memo string for the transaction.
  @$pb.TagNumber(6)
  $core.String get memo => $_getSZ(5);
  @$pb.TagNumber(6)
  set memo($core.String value) => $_setString(5, value);
  @$pb.TagNumber(6)
  $core.bool hasMemo() => $_has(5);
  @$pb.TagNumber(6)
  void clearMemo() => $_clearField(6);
}

/// Request message for retrieving raw details of a bond transaction.
class GetRawBondTransactionRequest extends $pb.GeneratedMessage {
  factory GetRawBondTransactionRequest({
    $core.int? lockTime,
    $core.String? sender,
    $core.String? receiver,
    $fixnum.Int64? stake,
    $core.String? publicKey,
    $fixnum.Int64? fee,
    $core.String? memo,
  }) {
    final result = create();
    if (lockTime != null) result.lockTime = lockTime;
    if (sender != null) result.sender = sender;
    if (receiver != null) result.receiver = receiver;
    if (stake != null) result.stake = stake;
    if (publicKey != null) result.publicKey = publicKey;
    if (fee != null) result.fee = fee;
    if (memo != null) result.memo = memo;
    return result;
  }

  GetRawBondTransactionRequest._();

  factory GetRawBondTransactionRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetRawBondTransactionRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetRawBondTransactionRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aI(1, _omitFieldNames ? '' : 'lockTime', fieldType: $pb.PbFieldType.OU3)
    ..aOS(2, _omitFieldNames ? '' : 'sender')
    ..aOS(3, _omitFieldNames ? '' : 'receiver')
    ..aInt64(4, _omitFieldNames ? '' : 'stake')
    ..aOS(5, _omitFieldNames ? '' : 'publicKey')
    ..aInt64(6, _omitFieldNames ? '' : 'fee')
    ..aOS(7, _omitFieldNames ? '' : 'memo')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRawBondTransactionRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRawBondTransactionRequest copyWith(
          void Function(GetRawBondTransactionRequest) updates) =>
      super.copyWith(
              (message) => updates(message as GetRawBondTransactionRequest))
          as GetRawBondTransactionRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetRawBondTransactionRequest create() =>
      GetRawBondTransactionRequest._();
  @$core.override
  GetRawBondTransactionRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetRawBondTransactionRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetRawBondTransactionRequest>(create);
  static GetRawBondTransactionRequest? _defaultInstance;

  /// The lock time for the transaction. If not set, defaults to the last block height.
  @$pb.TagNumber(1)
  $core.int get lockTime => $_getIZ(0);
  @$pb.TagNumber(1)
  set lockTime($core.int value) => $_setUnsignedInt32(0, value);
  @$pb.TagNumber(1)
  $core.bool hasLockTime() => $_has(0);
  @$pb.TagNumber(1)
  void clearLockTime() => $_clearField(1);

  /// The sender's account address.
  @$pb.TagNumber(2)
  $core.String get sender => $_getSZ(1);
  @$pb.TagNumber(2)
  set sender($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasSender() => $_has(1);
  @$pb.TagNumber(2)
  void clearSender() => $_clearField(2);

  /// The receiver's validator address.
  @$pb.TagNumber(3)
  $core.String get receiver => $_getSZ(2);
  @$pb.TagNumber(3)
  set receiver($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasReceiver() => $_has(2);
  @$pb.TagNumber(3)
  void clearReceiver() => $_clearField(3);

  /// The stake amount in NanoPAC. Must be greater than 0.
  @$pb.TagNumber(4)
  $fixnum.Int64 get stake => $_getI64(3);
  @$pb.TagNumber(4)
  set stake($fixnum.Int64 value) => $_setInt64(3, value);
  @$pb.TagNumber(4)
  $core.bool hasStake() => $_has(3);
  @$pb.TagNumber(4)
  void clearStake() => $_clearField(4);

  /// The public key of the validator. Optional, but required when registering a new validator.
  @$pb.TagNumber(5)
  $core.String get publicKey => $_getSZ(4);
  @$pb.TagNumber(5)
  set publicKey($core.String value) => $_setString(4, value);
  @$pb.TagNumber(5)
  $core.bool hasPublicKey() => $_has(4);
  @$pb.TagNumber(5)
  void clearPublicKey() => $_clearField(5);

  /// The transaction fee in NanoPAC. If not set, it is set to the estimated fee.
  @$pb.TagNumber(6)
  $fixnum.Int64 get fee => $_getI64(5);
  @$pb.TagNumber(6)
  set fee($fixnum.Int64 value) => $_setInt64(5, value);
  @$pb.TagNumber(6)
  $core.bool hasFee() => $_has(5);
  @$pb.TagNumber(6)
  void clearFee() => $_clearField(6);

  /// A memo string for the transaction.
  @$pb.TagNumber(7)
  $core.String get memo => $_getSZ(6);
  @$pb.TagNumber(7)
  set memo($core.String value) => $_setString(6, value);
  @$pb.TagNumber(7)
  $core.bool hasMemo() => $_has(6);
  @$pb.TagNumber(7)
  void clearMemo() => $_clearField(7);
}

/// Request message for retrieving raw details of an unbond transaction.
class GetRawUnbondTransactionRequest extends $pb.GeneratedMessage {
  factory GetRawUnbondTransactionRequest({
    $core.int? lockTime,
    $core.String? validatorAddress,
    $core.String? memo,
  }) {
    final result = create();
    if (lockTime != null) result.lockTime = lockTime;
    if (validatorAddress != null) result.validatorAddress = validatorAddress;
    if (memo != null) result.memo = memo;
    return result;
  }

  GetRawUnbondTransactionRequest._();

  factory GetRawUnbondTransactionRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetRawUnbondTransactionRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetRawUnbondTransactionRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aI(1, _omitFieldNames ? '' : 'lockTime', fieldType: $pb.PbFieldType.OU3)
    ..aOS(3, _omitFieldNames ? '' : 'validatorAddress')
    ..aOS(4, _omitFieldNames ? '' : 'memo')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRawUnbondTransactionRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRawUnbondTransactionRequest copyWith(
          void Function(GetRawUnbondTransactionRequest) updates) =>
      super.copyWith(
              (message) => updates(message as GetRawUnbondTransactionRequest))
          as GetRawUnbondTransactionRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetRawUnbondTransactionRequest create() =>
      GetRawUnbondTransactionRequest._();
  @$core.override
  GetRawUnbondTransactionRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetRawUnbondTransactionRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetRawUnbondTransactionRequest>(create);
  static GetRawUnbondTransactionRequest? _defaultInstance;

  /// The lock time for the transaction. If not set, defaults to the last block height.
  @$pb.TagNumber(1)
  $core.int get lockTime => $_getIZ(0);
  @$pb.TagNumber(1)
  set lockTime($core.int value) => $_setUnsignedInt32(0, value);
  @$pb.TagNumber(1)
  $core.bool hasLockTime() => $_has(0);
  @$pb.TagNumber(1)
  void clearLockTime() => $_clearField(1);

  /// The address of the validator to unbond from.
  @$pb.TagNumber(3)
  $core.String get validatorAddress => $_getSZ(1);
  @$pb.TagNumber(3)
  set validatorAddress($core.String value) => $_setString(1, value);
  @$pb.TagNumber(3)
  $core.bool hasValidatorAddress() => $_has(1);
  @$pb.TagNumber(3)
  void clearValidatorAddress() => $_clearField(3);

  /// A memo string for the transaction.
  @$pb.TagNumber(4)
  $core.String get memo => $_getSZ(2);
  @$pb.TagNumber(4)
  set memo($core.String value) => $_setString(2, value);
  @$pb.TagNumber(4)
  $core.bool hasMemo() => $_has(2);
  @$pb.TagNumber(4)
  void clearMemo() => $_clearField(4);
}

/// Request message for retrieving raw details of a withdraw transaction.
class GetRawWithdrawTransactionRequest extends $pb.GeneratedMessage {
  factory GetRawWithdrawTransactionRequest({
    $core.int? lockTime,
    $core.String? validatorAddress,
    $core.String? accountAddress,
    $fixnum.Int64? amount,
    $fixnum.Int64? fee,
    $core.String? memo,
  }) {
    final result = create();
    if (lockTime != null) result.lockTime = lockTime;
    if (validatorAddress != null) result.validatorAddress = validatorAddress;
    if (accountAddress != null) result.accountAddress = accountAddress;
    if (amount != null) result.amount = amount;
    if (fee != null) result.fee = fee;
    if (memo != null) result.memo = memo;
    return result;
  }

  GetRawWithdrawTransactionRequest._();

  factory GetRawWithdrawTransactionRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetRawWithdrawTransactionRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetRawWithdrawTransactionRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aI(1, _omitFieldNames ? '' : 'lockTime', fieldType: $pb.PbFieldType.OU3)
    ..aOS(2, _omitFieldNames ? '' : 'validatorAddress')
    ..aOS(3, _omitFieldNames ? '' : 'accountAddress')
    ..aInt64(4, _omitFieldNames ? '' : 'amount')
    ..aInt64(5, _omitFieldNames ? '' : 'fee')
    ..aOS(6, _omitFieldNames ? '' : 'memo')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRawWithdrawTransactionRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRawWithdrawTransactionRequest copyWith(
          void Function(GetRawWithdrawTransactionRequest) updates) =>
      super.copyWith(
              (message) => updates(message as GetRawWithdrawTransactionRequest))
          as GetRawWithdrawTransactionRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetRawWithdrawTransactionRequest create() =>
      GetRawWithdrawTransactionRequest._();
  @$core.override
  GetRawWithdrawTransactionRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetRawWithdrawTransactionRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetRawWithdrawTransactionRequest>(
          create);
  static GetRawWithdrawTransactionRequest? _defaultInstance;

  /// The lock time for the transaction. If not set, defaults to the last block height.
  @$pb.TagNumber(1)
  $core.int get lockTime => $_getIZ(0);
  @$pb.TagNumber(1)
  set lockTime($core.int value) => $_setUnsignedInt32(0, value);
  @$pb.TagNumber(1)
  $core.bool hasLockTime() => $_has(0);
  @$pb.TagNumber(1)
  void clearLockTime() => $_clearField(1);

  /// The address of the validator to withdraw from.
  @$pb.TagNumber(2)
  $core.String get validatorAddress => $_getSZ(1);
  @$pb.TagNumber(2)
  set validatorAddress($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasValidatorAddress() => $_has(1);
  @$pb.TagNumber(2)
  void clearValidatorAddress() => $_clearField(2);

  /// The address of the account to withdraw to.
  @$pb.TagNumber(3)
  $core.String get accountAddress => $_getSZ(2);
  @$pb.TagNumber(3)
  set accountAddress($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasAccountAddress() => $_has(2);
  @$pb.TagNumber(3)
  void clearAccountAddress() => $_clearField(3);

  /// The withdrawal amount in NanoPAC. Must be greater than 0.
  @$pb.TagNumber(4)
  $fixnum.Int64 get amount => $_getI64(3);
  @$pb.TagNumber(4)
  set amount($fixnum.Int64 value) => $_setInt64(3, value);
  @$pb.TagNumber(4)
  $core.bool hasAmount() => $_has(3);
  @$pb.TagNumber(4)
  void clearAmount() => $_clearField(4);

  /// The transaction fee in NanoPAC. If not set, it is set to the estimated fee.
  @$pb.TagNumber(5)
  $fixnum.Int64 get fee => $_getI64(4);
  @$pb.TagNumber(5)
  set fee($fixnum.Int64 value) => $_setInt64(4, value);
  @$pb.TagNumber(5)
  $core.bool hasFee() => $_has(4);
  @$pb.TagNumber(5)
  void clearFee() => $_clearField(5);

  /// A memo string for the transaction.
  @$pb.TagNumber(6)
  $core.String get memo => $_getSZ(5);
  @$pb.TagNumber(6)
  set memo($core.String value) => $_setString(5, value);
  @$pb.TagNumber(6)
  $core.bool hasMemo() => $_has(5);
  @$pb.TagNumber(6)
  void clearMemo() => $_clearField(6);
}

/// Request message for retrieving raw details of a batch transfer transaction.
class GetRawBatchTransferTransactionRequest extends $pb.GeneratedMessage {
  factory GetRawBatchTransferTransactionRequest({
    $core.int? lockTime,
    $core.String? sender,
    $core.Iterable<Recipient>? recipients,
    $fixnum.Int64? fee,
    $core.String? memo,
  }) {
    final result = create();
    if (lockTime != null) result.lockTime = lockTime;
    if (sender != null) result.sender = sender;
    if (recipients != null) result.recipients.addAll(recipients);
    if (fee != null) result.fee = fee;
    if (memo != null) result.memo = memo;
    return result;
  }

  GetRawBatchTransferTransactionRequest._();

  factory GetRawBatchTransferTransactionRequest.fromBuffer(
          $core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetRawBatchTransferTransactionRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetRawBatchTransferTransactionRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aI(1, _omitFieldNames ? '' : 'lockTime', fieldType: $pb.PbFieldType.OU3)
    ..aOS(2, _omitFieldNames ? '' : 'sender')
    ..pPM<Recipient>(3, _omitFieldNames ? '' : 'recipients',
        subBuilder: Recipient.create)
    ..aInt64(4, _omitFieldNames ? '' : 'fee')
    ..aOS(5, _omitFieldNames ? '' : 'memo')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRawBatchTransferTransactionRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRawBatchTransferTransactionRequest copyWith(
          void Function(GetRawBatchTransferTransactionRequest) updates) =>
      super.copyWith((message) =>
              updates(message as GetRawBatchTransferTransactionRequest))
          as GetRawBatchTransferTransactionRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetRawBatchTransferTransactionRequest create() =>
      GetRawBatchTransferTransactionRequest._();
  @$core.override
  GetRawBatchTransferTransactionRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetRawBatchTransferTransactionRequest getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<
          GetRawBatchTransferTransactionRequest>(create);
  static GetRawBatchTransferTransactionRequest? _defaultInstance;

  /// The lock time for the transaction. If not set, defaults to the last block height.
  @$pb.TagNumber(1)
  $core.int get lockTime => $_getIZ(0);
  @$pb.TagNumber(1)
  set lockTime($core.int value) => $_setUnsignedInt32(0, value);
  @$pb.TagNumber(1)
  $core.bool hasLockTime() => $_has(0);
  @$pb.TagNumber(1)
  void clearLockTime() => $_clearField(1);

  /// The sender's account address.
  @$pb.TagNumber(2)
  $core.String get sender => $_getSZ(1);
  @$pb.TagNumber(2)
  set sender($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasSender() => $_has(1);
  @$pb.TagNumber(2)
  void clearSender() => $_clearField(2);

  /// The list of recipients with their amounts. Minimum 2 recipients required.
  @$pb.TagNumber(3)
  $pb.PbList<Recipient> get recipients => $_getList(2);

  /// The transaction fee in NanoPAC. If not set, it is set to the estimated fee.
  @$pb.TagNumber(4)
  $fixnum.Int64 get fee => $_getI64(3);
  @$pb.TagNumber(4)
  set fee($fixnum.Int64 value) => $_setInt64(3, value);
  @$pb.TagNumber(4)
  $core.bool hasFee() => $_has(3);
  @$pb.TagNumber(4)
  void clearFee() => $_clearField(4);

  /// A memo string for the transaction.
  @$pb.TagNumber(5)
  $core.String get memo => $_getSZ(4);
  @$pb.TagNumber(5)
  set memo($core.String value) => $_setString(4, value);
  @$pb.TagNumber(5)
  $core.bool hasMemo() => $_has(4);
  @$pb.TagNumber(5)
  void clearMemo() => $_clearField(5);
}

/// Response message contains raw transaction data.
class GetRawTransactionResponse extends $pb.GeneratedMessage {
  factory GetRawTransactionResponse({
    $core.String? rawTransaction,
    $core.String? id,
  }) {
    final result = create();
    if (rawTransaction != null) result.rawTransaction = rawTransaction;
    if (id != null) result.id = id;
    return result;
  }

  GetRawTransactionResponse._();

  factory GetRawTransactionResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetRawTransactionResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetRawTransactionResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'rawTransaction')
    ..aOS(2, _omitFieldNames ? '' : 'id')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRawTransactionResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetRawTransactionResponse copyWith(
          void Function(GetRawTransactionResponse) updates) =>
      super.copyWith((message) => updates(message as GetRawTransactionResponse))
          as GetRawTransactionResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetRawTransactionResponse create() => GetRawTransactionResponse._();
  @$core.override
  GetRawTransactionResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetRawTransactionResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetRawTransactionResponse>(create);
  static GetRawTransactionResponse? _defaultInstance;

  /// The raw transaction data in hexadecimal format.
  @$pb.TagNumber(1)
  $core.String get rawTransaction => $_getSZ(0);
  @$pb.TagNumber(1)
  set rawTransaction($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasRawTransaction() => $_has(0);
  @$pb.TagNumber(1)
  void clearRawTransaction() => $_clearField(1);

  /// The unique ID of the transaction.
  @$pb.TagNumber(2)
  $core.String get id => $_getSZ(1);
  @$pb.TagNumber(2)
  set id($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasId() => $_has(1);
  @$pb.TagNumber(2)
  void clearId() => $_clearField(2);
}

/// Payload for a transfer transaction.
class PayloadTransfer extends $pb.GeneratedMessage {
  factory PayloadTransfer({
    $core.String? sender,
    $core.String? receiver,
    $fixnum.Int64? amount,
  }) {
    final result = create();
    if (sender != null) result.sender = sender;
    if (receiver != null) result.receiver = receiver;
    if (amount != null) result.amount = amount;
    return result;
  }

  PayloadTransfer._();

  factory PayloadTransfer.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PayloadTransfer.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PayloadTransfer',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'sender')
    ..aOS(2, _omitFieldNames ? '' : 'receiver')
    ..aInt64(3, _omitFieldNames ? '' : 'amount')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PayloadTransfer clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PayloadTransfer copyWith(void Function(PayloadTransfer) updates) =>
      super.copyWith((message) => updates(message as PayloadTransfer))
          as PayloadTransfer;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PayloadTransfer create() => PayloadTransfer._();
  @$core.override
  PayloadTransfer createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PayloadTransfer getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PayloadTransfer>(create);
  static PayloadTransfer? _defaultInstance;

  /// The sender's address.
  @$pb.TagNumber(1)
  $core.String get sender => $_getSZ(0);
  @$pb.TagNumber(1)
  set sender($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasSender() => $_has(0);
  @$pb.TagNumber(1)
  void clearSender() => $_clearField(1);

  /// The receiver's address.
  @$pb.TagNumber(2)
  $core.String get receiver => $_getSZ(1);
  @$pb.TagNumber(2)
  set receiver($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasReceiver() => $_has(1);
  @$pb.TagNumber(2)
  void clearReceiver() => $_clearField(2);

  /// The amount to be transferred in NanoPAC.
  @$pb.TagNumber(3)
  $fixnum.Int64 get amount => $_getI64(2);
  @$pb.TagNumber(3)
  set amount($fixnum.Int64 value) => $_setInt64(2, value);
  @$pb.TagNumber(3)
  $core.bool hasAmount() => $_has(2);
  @$pb.TagNumber(3)
  void clearAmount() => $_clearField(3);
}

/// Payload for a bond transaction.
class PayloadBond extends $pb.GeneratedMessage {
  factory PayloadBond({
    $core.String? sender,
    $core.String? receiver,
    $fixnum.Int64? stake,
    $core.String? publicKey,
  }) {
    final result = create();
    if (sender != null) result.sender = sender;
    if (receiver != null) result.receiver = receiver;
    if (stake != null) result.stake = stake;
    if (publicKey != null) result.publicKey = publicKey;
    return result;
  }

  PayloadBond._();

  factory PayloadBond.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PayloadBond.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PayloadBond',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'sender')
    ..aOS(2, _omitFieldNames ? '' : 'receiver')
    ..aInt64(3, _omitFieldNames ? '' : 'stake')
    ..aOS(4, _omitFieldNames ? '' : 'publicKey')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PayloadBond clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PayloadBond copyWith(void Function(PayloadBond) updates) =>
      super.copyWith((message) => updates(message as PayloadBond))
          as PayloadBond;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PayloadBond create() => PayloadBond._();
  @$core.override
  PayloadBond createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PayloadBond getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PayloadBond>(create);
  static PayloadBond? _defaultInstance;

  /// The sender's address.
  @$pb.TagNumber(1)
  $core.String get sender => $_getSZ(0);
  @$pb.TagNumber(1)
  set sender($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasSender() => $_has(0);
  @$pb.TagNumber(1)
  void clearSender() => $_clearField(1);

  /// The receiver's address.
  @$pb.TagNumber(2)
  $core.String get receiver => $_getSZ(1);
  @$pb.TagNumber(2)
  set receiver($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasReceiver() => $_has(1);
  @$pb.TagNumber(2)
  void clearReceiver() => $_clearField(2);

  /// The stake amount in NanoPAC.
  @$pb.TagNumber(3)
  $fixnum.Int64 get stake => $_getI64(2);
  @$pb.TagNumber(3)
  set stake($fixnum.Int64 value) => $_setInt64(2, value);
  @$pb.TagNumber(3)
  $core.bool hasStake() => $_has(2);
  @$pb.TagNumber(3)
  void clearStake() => $_clearField(3);

  /// The public key of the validator.
  @$pb.TagNumber(4)
  $core.String get publicKey => $_getSZ(3);
  @$pb.TagNumber(4)
  set publicKey($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasPublicKey() => $_has(3);
  @$pb.TagNumber(4)
  void clearPublicKey() => $_clearField(4);
}

/// Payload for a sortition transaction.
class PayloadSortition extends $pb.GeneratedMessage {
  factory PayloadSortition({
    $core.String? address,
    $core.String? proof,
  }) {
    final result = create();
    if (address != null) result.address = address;
    if (proof != null) result.proof = proof;
    return result;
  }

  PayloadSortition._();

  factory PayloadSortition.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PayloadSortition.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PayloadSortition',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'address')
    ..aOS(2, _omitFieldNames ? '' : 'proof')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PayloadSortition clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PayloadSortition copyWith(void Function(PayloadSortition) updates) =>
      super.copyWith((message) => updates(message as PayloadSortition))
          as PayloadSortition;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PayloadSortition create() => PayloadSortition._();
  @$core.override
  PayloadSortition createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PayloadSortition getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PayloadSortition>(create);
  static PayloadSortition? _defaultInstance;

  /// The validator address associated with the sortition proof.
  @$pb.TagNumber(1)
  $core.String get address => $_getSZ(0);
  @$pb.TagNumber(1)
  set address($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasAddress() => $_has(0);
  @$pb.TagNumber(1)
  void clearAddress() => $_clearField(1);

  /// The proof for the sortition.
  @$pb.TagNumber(2)
  $core.String get proof => $_getSZ(1);
  @$pb.TagNumber(2)
  set proof($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasProof() => $_has(1);
  @$pb.TagNumber(2)
  void clearProof() => $_clearField(2);
}

/// Payload for an unbond transaction.
class PayloadUnbond extends $pb.GeneratedMessage {
  factory PayloadUnbond({
    $core.String? validator,
  }) {
    final result = create();
    if (validator != null) result.validator = validator;
    return result;
  }

  PayloadUnbond._();

  factory PayloadUnbond.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PayloadUnbond.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PayloadUnbond',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'validator')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PayloadUnbond clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PayloadUnbond copyWith(void Function(PayloadUnbond) updates) =>
      super.copyWith((message) => updates(message as PayloadUnbond))
          as PayloadUnbond;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PayloadUnbond create() => PayloadUnbond._();
  @$core.override
  PayloadUnbond createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PayloadUnbond getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PayloadUnbond>(create);
  static PayloadUnbond? _defaultInstance;

  /// The address of the validator to unbond from.
  @$pb.TagNumber(1)
  $core.String get validator => $_getSZ(0);
  @$pb.TagNumber(1)
  set validator($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasValidator() => $_has(0);
  @$pb.TagNumber(1)
  void clearValidator() => $_clearField(1);
}

/// Payload for a withdraw transaction.
class PayloadWithdraw extends $pb.GeneratedMessage {
  factory PayloadWithdraw({
    $core.String? validatorAddress,
    $core.String? accountAddress,
    $fixnum.Int64? amount,
  }) {
    final result = create();
    if (validatorAddress != null) result.validatorAddress = validatorAddress;
    if (accountAddress != null) result.accountAddress = accountAddress;
    if (amount != null) result.amount = amount;
    return result;
  }

  PayloadWithdraw._();

  factory PayloadWithdraw.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PayloadWithdraw.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PayloadWithdraw',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'validatorAddress')
    ..aOS(2, _omitFieldNames ? '' : 'accountAddress')
    ..aInt64(3, _omitFieldNames ? '' : 'amount')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PayloadWithdraw clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PayloadWithdraw copyWith(void Function(PayloadWithdraw) updates) =>
      super.copyWith((message) => updates(message as PayloadWithdraw))
          as PayloadWithdraw;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PayloadWithdraw create() => PayloadWithdraw._();
  @$core.override
  PayloadWithdraw createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PayloadWithdraw getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PayloadWithdraw>(create);
  static PayloadWithdraw? _defaultInstance;

  /// The address of the validator to withdraw from.
  @$pb.TagNumber(1)
  $core.String get validatorAddress => $_getSZ(0);
  @$pb.TagNumber(1)
  set validatorAddress($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasValidatorAddress() => $_has(0);
  @$pb.TagNumber(1)
  void clearValidatorAddress() => $_clearField(1);

  /// The address of the account to withdraw to.
  @$pb.TagNumber(2)
  $core.String get accountAddress => $_getSZ(1);
  @$pb.TagNumber(2)
  set accountAddress($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasAccountAddress() => $_has(1);
  @$pb.TagNumber(2)
  void clearAccountAddress() => $_clearField(2);

  /// The withdrawal amount in NanoPAC.
  @$pb.TagNumber(3)
  $fixnum.Int64 get amount => $_getI64(2);
  @$pb.TagNumber(3)
  set amount($fixnum.Int64 value) => $_setInt64(2, value);
  @$pb.TagNumber(3)
  $core.bool hasAmount() => $_has(2);
  @$pb.TagNumber(3)
  void clearAmount() => $_clearField(3);
}

/// Payload for a batch transfer transaction.
class PayloadBatchTransfer extends $pb.GeneratedMessage {
  factory PayloadBatchTransfer({
    $core.String? sender,
    $core.Iterable<Recipient>? recipients,
  }) {
    final result = create();
    if (sender != null) result.sender = sender;
    if (recipients != null) result.recipients.addAll(recipients);
    return result;
  }

  PayloadBatchTransfer._();

  factory PayloadBatchTransfer.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PayloadBatchTransfer.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PayloadBatchTransfer',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'sender')
    ..pPM<Recipient>(2, _omitFieldNames ? '' : 'recipients',
        subBuilder: Recipient.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PayloadBatchTransfer clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PayloadBatchTransfer copyWith(void Function(PayloadBatchTransfer) updates) =>
      super.copyWith((message) => updates(message as PayloadBatchTransfer))
          as PayloadBatchTransfer;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PayloadBatchTransfer create() => PayloadBatchTransfer._();
  @$core.override
  PayloadBatchTransfer createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PayloadBatchTransfer getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PayloadBatchTransfer>(create);
  static PayloadBatchTransfer? _defaultInstance;

  /// The sender's address.
  @$pb.TagNumber(1)
  $core.String get sender => $_getSZ(0);
  @$pb.TagNumber(1)
  set sender($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasSender() => $_has(0);
  @$pb.TagNumber(1)
  void clearSender() => $_clearField(1);

  /// The list of recipients with their amounts.
  @$pb.TagNumber(2)
  $pb.PbList<Recipient> get recipients => $_getList(1);
}

/// Recipient is receiver with amount.
class Recipient extends $pb.GeneratedMessage {
  factory Recipient({
    $core.String? receiver,
    $fixnum.Int64? amount,
  }) {
    final result = create();
    if (receiver != null) result.receiver = receiver;
    if (amount != null) result.amount = amount;
    return result;
  }

  Recipient._();

  factory Recipient.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory Recipient.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'Recipient',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'receiver')
    ..aInt64(2, _omitFieldNames ? '' : 'amount')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Recipient clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  Recipient copyWith(void Function(Recipient) updates) =>
      super.copyWith((message) => updates(message as Recipient)) as Recipient;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static Recipient create() => Recipient._();
  @$core.override
  Recipient createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static Recipient getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<Recipient>(create);
  static Recipient? _defaultInstance;

  /// The receiver's address.
  @$pb.TagNumber(1)
  $core.String get receiver => $_getSZ(0);
  @$pb.TagNumber(1)
  set receiver($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasReceiver() => $_has(0);
  @$pb.TagNumber(1)
  void clearReceiver() => $_clearField(1);

  /// The amount in NanoPAC.
  @$pb.TagNumber(2)
  $fixnum.Int64 get amount => $_getI64(1);
  @$pb.TagNumber(2)
  set amount($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(2)
  $core.bool hasAmount() => $_has(1);
  @$pb.TagNumber(2)
  void clearAmount() => $_clearField(2);
}

enum TransactionInfo_Payload {
  transfer,
  bond,
  sortition,
  unbond,
  withdraw,
  batchTransfer,
  notSet
}

/// Information about a transaction.
class TransactionInfo extends $pb.GeneratedMessage {
  factory TransactionInfo({
    $core.String? id,
    $core.String? data,
    $core.int? version,
    $core.int? lockTime,
    $fixnum.Int64? value,
    $fixnum.Int64? fee,
    PayloadType? payloadType,
    $core.String? memo,
    $core.String? publicKey,
    $core.String? signature,
    $core.int? blockHeight,
    $core.bool? confirmed,
    $core.int? confirmations,
    PayloadTransfer? transfer,
    PayloadBond? bond,
    PayloadSortition? sortition,
    PayloadUnbond? unbond,
    PayloadWithdraw? withdraw,
    PayloadBatchTransfer? batchTransfer,
  }) {
    final result = create();
    if (id != null) result.id = id;
    if (data != null) result.data = data;
    if (version != null) result.version = version;
    if (lockTime != null) result.lockTime = lockTime;
    if (value != null) result.value = value;
    if (fee != null) result.fee = fee;
    if (payloadType != null) result.payloadType = payloadType;
    if (memo != null) result.memo = memo;
    if (publicKey != null) result.publicKey = publicKey;
    if (signature != null) result.signature = signature;
    if (blockHeight != null) result.blockHeight = blockHeight;
    if (confirmed != null) result.confirmed = confirmed;
    if (confirmations != null) result.confirmations = confirmations;
    if (transfer != null) result.transfer = transfer;
    if (bond != null) result.bond = bond;
    if (sortition != null) result.sortition = sortition;
    if (unbond != null) result.unbond = unbond;
    if (withdraw != null) result.withdraw = withdraw;
    if (batchTransfer != null) result.batchTransfer = batchTransfer;
    return result;
  }

  TransactionInfo._();

  factory TransactionInfo.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory TransactionInfo.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static const $core.Map<$core.int, TransactionInfo_Payload>
      _TransactionInfo_PayloadByTag = {
    30: TransactionInfo_Payload.transfer,
    31: TransactionInfo_Payload.bond,
    32: TransactionInfo_Payload.sortition,
    33: TransactionInfo_Payload.unbond,
    34: TransactionInfo_Payload.withdraw,
    35: TransactionInfo_Payload.batchTransfer,
    0: TransactionInfo_Payload.notSet
  };
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'TransactionInfo',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..oo(0, [30, 31, 32, 33, 34, 35])
    ..aOS(1, _omitFieldNames ? '' : 'id')
    ..aOS(2, _omitFieldNames ? '' : 'data')
    ..aI(3, _omitFieldNames ? '' : 'version')
    ..aI(4, _omitFieldNames ? '' : 'lockTime', fieldType: $pb.PbFieldType.OU3)
    ..aInt64(5, _omitFieldNames ? '' : 'value')
    ..aInt64(6, _omitFieldNames ? '' : 'fee')
    ..aE<PayloadType>(7, _omitFieldNames ? '' : 'payloadType',
        enumValues: PayloadType.values)
    ..aOS(8, _omitFieldNames ? '' : 'memo')
    ..aOS(9, _omitFieldNames ? '' : 'publicKey')
    ..aOS(10, _omitFieldNames ? '' : 'signature')
    ..aI(11, _omitFieldNames ? '' : 'blockHeight',
        fieldType: $pb.PbFieldType.OU3)
    ..aOB(12, _omitFieldNames ? '' : 'confirmed')
    ..aI(13, _omitFieldNames ? '' : 'confirmations')
    ..aOM<PayloadTransfer>(30, _omitFieldNames ? '' : 'transfer',
        subBuilder: PayloadTransfer.create)
    ..aOM<PayloadBond>(31, _omitFieldNames ? '' : 'bond',
        subBuilder: PayloadBond.create)
    ..aOM<PayloadSortition>(32, _omitFieldNames ? '' : 'sortition',
        subBuilder: PayloadSortition.create)
    ..aOM<PayloadUnbond>(33, _omitFieldNames ? '' : 'unbond',
        subBuilder: PayloadUnbond.create)
    ..aOM<PayloadWithdraw>(34, _omitFieldNames ? '' : 'withdraw',
        subBuilder: PayloadWithdraw.create)
    ..aOM<PayloadBatchTransfer>(35, _omitFieldNames ? '' : 'batchTransfer',
        subBuilder: PayloadBatchTransfer.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TransactionInfo clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  TransactionInfo copyWith(void Function(TransactionInfo) updates) =>
      super.copyWith((message) => updates(message as TransactionInfo))
          as TransactionInfo;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static TransactionInfo create() => TransactionInfo._();
  @$core.override
  TransactionInfo createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static TransactionInfo getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<TransactionInfo>(create);
  static TransactionInfo? _defaultInstance;

  @$pb.TagNumber(30)
  @$pb.TagNumber(31)
  @$pb.TagNumber(32)
  @$pb.TagNumber(33)
  @$pb.TagNumber(34)
  @$pb.TagNumber(35)
  TransactionInfo_Payload whichPayload() =>
      _TransactionInfo_PayloadByTag[$_whichOneof(0)]!;
  @$pb.TagNumber(30)
  @$pb.TagNumber(31)
  @$pb.TagNumber(32)
  @$pb.TagNumber(33)
  @$pb.TagNumber(34)
  @$pb.TagNumber(35)
  void clearPayload() => $_clearField($_whichOneof(0));

  /// The unique ID of the transaction.
  @$pb.TagNumber(1)
  $core.String get id => $_getSZ(0);
  @$pb.TagNumber(1)
  set id($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => $_clearField(1);

  /// The raw transaction data in hexadecimal format.
  @$pb.TagNumber(2)
  $core.String get data => $_getSZ(1);
  @$pb.TagNumber(2)
  set data($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasData() => $_has(1);
  @$pb.TagNumber(2)
  void clearData() => $_clearField(2);

  /// The version of the transaction.
  @$pb.TagNumber(3)
  $core.int get version => $_getIZ(2);
  @$pb.TagNumber(3)
  set version($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(3)
  $core.bool hasVersion() => $_has(2);
  @$pb.TagNumber(3)
  void clearVersion() => $_clearField(3);

  /// The lock time for the transaction.
  @$pb.TagNumber(4)
  $core.int get lockTime => $_getIZ(3);
  @$pb.TagNumber(4)
  set lockTime($core.int value) => $_setUnsignedInt32(3, value);
  @$pb.TagNumber(4)
  $core.bool hasLockTime() => $_has(3);
  @$pb.TagNumber(4)
  void clearLockTime() => $_clearField(4);

  /// The value of the transaction in NanoPAC.
  @$pb.TagNumber(5)
  $fixnum.Int64 get value => $_getI64(4);
  @$pb.TagNumber(5)
  set value($fixnum.Int64 value) => $_setInt64(4, value);
  @$pb.TagNumber(5)
  $core.bool hasValue() => $_has(4);
  @$pb.TagNumber(5)
  void clearValue() => $_clearField(5);

  /// The fee for the transaction in NanoPAC.
  @$pb.TagNumber(6)
  $fixnum.Int64 get fee => $_getI64(5);
  @$pb.TagNumber(6)
  set fee($fixnum.Int64 value) => $_setInt64(5, value);
  @$pb.TagNumber(6)
  $core.bool hasFee() => $_has(5);
  @$pb.TagNumber(6)
  void clearFee() => $_clearField(6);

  /// The type of transaction payload.
  @$pb.TagNumber(7)
  PayloadType get payloadType => $_getN(6);
  @$pb.TagNumber(7)
  set payloadType(PayloadType value) => $_setField(7, value);
  @$pb.TagNumber(7)
  $core.bool hasPayloadType() => $_has(6);
  @$pb.TagNumber(7)
  void clearPayloadType() => $_clearField(7);

  /// A memo string for the transaction.
  @$pb.TagNumber(8)
  $core.String get memo => $_getSZ(7);
  @$pb.TagNumber(8)
  set memo($core.String value) => $_setString(7, value);
  @$pb.TagNumber(8)
  $core.bool hasMemo() => $_has(7);
  @$pb.TagNumber(8)
  void clearMemo() => $_clearField(8);

  /// The public key associated with the transaction.
  @$pb.TagNumber(9)
  $core.String get publicKey => $_getSZ(8);
  @$pb.TagNumber(9)
  set publicKey($core.String value) => $_setString(8, value);
  @$pb.TagNumber(9)
  $core.bool hasPublicKey() => $_has(8);
  @$pb.TagNumber(9)
  void clearPublicKey() => $_clearField(9);

  /// The signature for the transaction.
  @$pb.TagNumber(10)
  $core.String get signature => $_getSZ(9);
  @$pb.TagNumber(10)
  set signature($core.String value) => $_setString(9, value);
  @$pb.TagNumber(10)
  $core.bool hasSignature() => $_has(9);
  @$pb.TagNumber(10)
  void clearSignature() => $_clearField(10);

  /// The block height containing the transaction.
  /// A value of zero means the transaction is unconfirmed and may still in the transaction pool.
  @$pb.TagNumber(11)
  $core.int get blockHeight => $_getIZ(10);
  @$pb.TagNumber(11)
  set blockHeight($core.int value) => $_setUnsignedInt32(10, value);
  @$pb.TagNumber(11)
  $core.bool hasBlockHeight() => $_has(10);
  @$pb.TagNumber(11)
  void clearBlockHeight() => $_clearField(11);

  /// Indicates whether the transaction is confirmed.
  @$pb.TagNumber(12)
  $core.bool get confirmed => $_getBF(11);
  @$pb.TagNumber(12)
  set confirmed($core.bool value) => $_setBool(11, value);
  @$pb.TagNumber(12)
  $core.bool hasConfirmed() => $_has(11);
  @$pb.TagNumber(12)
  void clearConfirmed() => $_clearField(12);

  /// The number of blocks that have been added to the chain after this transaction was included in a block.
  /// A value of zero means the transaction is unconfirmed and may still in the transaction pool.
  @$pb.TagNumber(13)
  $core.int get confirmations => $_getIZ(12);
  @$pb.TagNumber(13)
  set confirmations($core.int value) => $_setSignedInt32(12, value);
  @$pb.TagNumber(13)
  $core.bool hasConfirmations() => $_has(12);
  @$pb.TagNumber(13)
  void clearConfirmations() => $_clearField(13);

  /// Transfer transaction payload.
  @$pb.TagNumber(30)
  PayloadTransfer get transfer => $_getN(13);
  @$pb.TagNumber(30)
  set transfer(PayloadTransfer value) => $_setField(30, value);
  @$pb.TagNumber(30)
  $core.bool hasTransfer() => $_has(13);
  @$pb.TagNumber(30)
  void clearTransfer() => $_clearField(30);
  @$pb.TagNumber(30)
  PayloadTransfer ensureTransfer() => $_ensure(13);

  /// Bond transaction payload.
  @$pb.TagNumber(31)
  PayloadBond get bond => $_getN(14);
  @$pb.TagNumber(31)
  set bond(PayloadBond value) => $_setField(31, value);
  @$pb.TagNumber(31)
  $core.bool hasBond() => $_has(14);
  @$pb.TagNumber(31)
  void clearBond() => $_clearField(31);
  @$pb.TagNumber(31)
  PayloadBond ensureBond() => $_ensure(14);

  /// Sortition transaction payload.
  @$pb.TagNumber(32)
  PayloadSortition get sortition => $_getN(15);
  @$pb.TagNumber(32)
  set sortition(PayloadSortition value) => $_setField(32, value);
  @$pb.TagNumber(32)
  $core.bool hasSortition() => $_has(15);
  @$pb.TagNumber(32)
  void clearSortition() => $_clearField(32);
  @$pb.TagNumber(32)
  PayloadSortition ensureSortition() => $_ensure(15);

  /// Unbond transaction payload.
  @$pb.TagNumber(33)
  PayloadUnbond get unbond => $_getN(16);
  @$pb.TagNumber(33)
  set unbond(PayloadUnbond value) => $_setField(33, value);
  @$pb.TagNumber(33)
  $core.bool hasUnbond() => $_has(16);
  @$pb.TagNumber(33)
  void clearUnbond() => $_clearField(33);
  @$pb.TagNumber(33)
  PayloadUnbond ensureUnbond() => $_ensure(16);

  /// Withdraw transaction payload.
  @$pb.TagNumber(34)
  PayloadWithdraw get withdraw => $_getN(17);
  @$pb.TagNumber(34)
  set withdraw(PayloadWithdraw value) => $_setField(34, value);
  @$pb.TagNumber(34)
  $core.bool hasWithdraw() => $_has(17);
  @$pb.TagNumber(34)
  void clearWithdraw() => $_clearField(34);
  @$pb.TagNumber(34)
  PayloadWithdraw ensureWithdraw() => $_ensure(17);

  /// Batch Transfer transaction payload.
  @$pb.TagNumber(35)
  PayloadBatchTransfer get batchTransfer => $_getN(18);
  @$pb.TagNumber(35)
  set batchTransfer(PayloadBatchTransfer value) => $_setField(35, value);
  @$pb.TagNumber(35)
  $core.bool hasBatchTransfer() => $_has(18);
  @$pb.TagNumber(35)
  void clearBatchTransfer() => $_clearField(35);
  @$pb.TagNumber(35)
  PayloadBatchTransfer ensureBatchTransfer() => $_ensure(18);
}

/// Request message for decoding a raw transaction.
class DecodeRawTransactionRequest extends $pb.GeneratedMessage {
  factory DecodeRawTransactionRequest({
    $core.String? rawTransaction,
  }) {
    final result = create();
    if (rawTransaction != null) result.rawTransaction = rawTransaction;
    return result;
  }

  DecodeRawTransactionRequest._();

  factory DecodeRawTransactionRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DecodeRawTransactionRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DecodeRawTransactionRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'rawTransaction')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DecodeRawTransactionRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DecodeRawTransactionRequest copyWith(
          void Function(DecodeRawTransactionRequest) updates) =>
      super.copyWith(
              (message) => updates(message as DecodeRawTransactionRequest))
          as DecodeRawTransactionRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DecodeRawTransactionRequest create() =>
      DecodeRawTransactionRequest._();
  @$core.override
  DecodeRawTransactionRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DecodeRawTransactionRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DecodeRawTransactionRequest>(create);
  static DecodeRawTransactionRequest? _defaultInstance;

  /// The raw transaction data in hexadecimal format.
  @$pb.TagNumber(1)
  $core.String get rawTransaction => $_getSZ(0);
  @$pb.TagNumber(1)
  set rawTransaction($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasRawTransaction() => $_has(0);
  @$pb.TagNumber(1)
  void clearRawTransaction() => $_clearField(1);
}

/// Response message contains the decoded transaction.
class DecodeRawTransactionResponse extends $pb.GeneratedMessage {
  factory DecodeRawTransactionResponse({
    TransactionInfo? transaction,
  }) {
    final result = create();
    if (transaction != null) result.transaction = transaction;
    return result;
  }

  DecodeRawTransactionResponse._();

  factory DecodeRawTransactionResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory DecodeRawTransactionResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'DecodeRawTransactionResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOM<TransactionInfo>(1, _omitFieldNames ? '' : 'transaction',
        subBuilder: TransactionInfo.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DecodeRawTransactionResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  DecodeRawTransactionResponse copyWith(
          void Function(DecodeRawTransactionResponse) updates) =>
      super.copyWith(
              (message) => updates(message as DecodeRawTransactionResponse))
          as DecodeRawTransactionResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static DecodeRawTransactionResponse create() =>
      DecodeRawTransactionResponse._();
  @$core.override
  DecodeRawTransactionResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static DecodeRawTransactionResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<DecodeRawTransactionResponse>(create);
  static DecodeRawTransactionResponse? _defaultInstance;

  /// The decoded transaction information.
  @$pb.TagNumber(1)
  TransactionInfo get transaction => $_getN(0);
  @$pb.TagNumber(1)
  set transaction(TransactionInfo value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasTransaction() => $_has(0);
  @$pb.TagNumber(1)
  void clearTransaction() => $_clearField(1);
  @$pb.TagNumber(1)
  TransactionInfo ensureTransaction() => $_ensure(0);
}

/// Transaction service defines various RPC methods for interacting with transactions.
class TransactionApi {
  final $pb.RpcClient _client;

  TransactionApi(this._client);

  /// GetTransaction retrieves transaction details based on the provided request parameters.
  $async.Future<GetTransactionResponse> getTransaction(
          $pb.ClientContext? ctx, GetTransactionRequest request) =>
      _client.invoke<GetTransactionResponse>(ctx, 'Transaction',
          'GetTransaction', request, GetTransactionResponse());

  /// CalculateFee calculates the transaction fee based on the specified amount and payload type.
  $async.Future<CalculateFeeResponse> calculateFee(
          $pb.ClientContext? ctx, CalculateFeeRequest request) =>
      _client.invoke<CalculateFeeResponse>(
          ctx, 'Transaction', 'CalculateFee', request, CalculateFeeResponse());

  /// BroadcastTransaction broadcasts a signed transaction to the network.
  $async.Future<BroadcastTransactionResponse> broadcastTransaction(
          $pb.ClientContext? ctx, BroadcastTransactionRequest request) =>
      _client.invoke<BroadcastTransactionResponse>(ctx, 'Transaction',
          'BroadcastTransaction', request, BroadcastTransactionResponse());

  /// GetRawTransferTransaction retrieves raw details of a transfer transaction.
  $async.Future<GetRawTransactionResponse> getRawTransferTransaction(
          $pb.ClientContext? ctx, GetRawTransferTransactionRequest request) =>
      _client.invoke<GetRawTransactionResponse>(ctx, 'Transaction',
          'GetRawTransferTransaction', request, GetRawTransactionResponse());

  /// GetRawBondTransaction retrieves raw details of a bond transaction.
  $async.Future<GetRawTransactionResponse> getRawBondTransaction(
          $pb.ClientContext? ctx, GetRawBondTransactionRequest request) =>
      _client.invoke<GetRawTransactionResponse>(ctx, 'Transaction',
          'GetRawBondTransaction', request, GetRawTransactionResponse());

  /// GetRawUnbondTransaction retrieves raw details of an unbond transaction.
  $async.Future<GetRawTransactionResponse> getRawUnbondTransaction(
          $pb.ClientContext? ctx, GetRawUnbondTransactionRequest request) =>
      _client.invoke<GetRawTransactionResponse>(ctx, 'Transaction',
          'GetRawUnbondTransaction', request, GetRawTransactionResponse());

  /// GetRawWithdrawTransaction retrieves raw details of a withdraw transaction.
  $async.Future<GetRawTransactionResponse> getRawWithdrawTransaction(
          $pb.ClientContext? ctx, GetRawWithdrawTransactionRequest request) =>
      _client.invoke<GetRawTransactionResponse>(ctx, 'Transaction',
          'GetRawWithdrawTransaction', request, GetRawTransactionResponse());

  /// GetRawBatchTransferTransaction retrieves raw details of batch transfer transaction.
  $async.Future<GetRawTransactionResponse> getRawBatchTransferTransaction(
          $pb.ClientContext? ctx,
          GetRawBatchTransferTransactionRequest request) =>
      _client.invoke<GetRawTransactionResponse>(
          ctx,
          'Transaction',
          'GetRawBatchTransferTransaction',
          request,
          GetRawTransactionResponse());

  /// DecodeRawTransaction accepts raw transaction and returns decoded transaction.
  $async.Future<DecodeRawTransactionResponse> decodeRawTransaction(
          $pb.ClientContext? ctx, DecodeRawTransactionRequest request) =>
      _client.invoke<DecodeRawTransactionResponse>(ctx, 'Transaction',
          'DecodeRawTransaction', request, DecodeRawTransactionResponse());
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
