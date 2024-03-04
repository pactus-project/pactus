///
//  Generated code. Do not modify.
//  source: transaction.proto
//
// @dart = 2.12
// ignore_for_file: annotate_overrides,camel_case_types,constant_identifier_names,directives_ordering,library_prefixes,non_constant_identifier_names,prefer_final_fields,return_of_invalid_type,unnecessary_const,unnecessary_import,unnecessary_this,unused_import,unused_shown_name

import 'dart:async' as $async;
import 'dart:core' as $core;

import 'package:fixnum/fixnum.dart' as $fixnum;
import 'package:protobuf/protobuf.dart' as $pb;

import 'transaction.pbenum.dart';

export 'transaction.pbenum.dart';

class GetTransactionRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'GetTransactionRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..a<$core.List<$core.int>>(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'id', $pb.PbFieldType.OY)
    ..e<TransactionVerbosity>(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'verbosity', $pb.PbFieldType.OE, defaultOrMaker: TransactionVerbosity.TRANSACTION_DATA, valueOf: TransactionVerbosity.valueOf, enumValues: TransactionVerbosity.values)
    ..hasRequiredFields = false
  ;

  GetTransactionRequest._() : super();
  factory GetTransactionRequest({
    $core.List<$core.int>? id,
    TransactionVerbosity? verbosity,
  }) {
    final _result = create();
    if (id != null) {
      _result.id = id;
    }
    if (verbosity != null) {
      _result.verbosity = verbosity;
    }
    return _result;
  }
  factory GetTransactionRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetTransactionRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetTransactionRequest clone() => GetTransactionRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetTransactionRequest copyWith(void Function(GetTransactionRequest) updates) => super.copyWith((message) => updates(message as GetTransactionRequest)) as GetTransactionRequest; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static GetTransactionRequest create() => GetTransactionRequest._();
  GetTransactionRequest createEmptyInstance() => create();
  static $pb.PbList<GetTransactionRequest> createRepeated() => $pb.PbList<GetTransactionRequest>();
  @$core.pragma('dart2js:noInline')
  static GetTransactionRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetTransactionRequest>(create);
  static GetTransactionRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.List<$core.int> get id => $_getN(0);
  @$pb.TagNumber(1)
  set id($core.List<$core.int> v) { $_setBytes(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => clearField(1);

  @$pb.TagNumber(2)
  TransactionVerbosity get verbosity => $_getN(1);
  @$pb.TagNumber(2)
  set verbosity(TransactionVerbosity v) { setField(2, v); }
  @$pb.TagNumber(2)
  $core.bool hasVerbosity() => $_has(1);
  @$pb.TagNumber(2)
  void clearVerbosity() => clearField(2);
}

class GetTransactionResponse extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'GetTransactionResponse', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOM<TransactionInfo>(3, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'transaction', subBuilder: TransactionInfo.create)
    ..a<$core.int>(12, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'blockHeight', $pb.PbFieldType.OU3)
    ..a<$core.int>(13, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'blockTime', $pb.PbFieldType.OU3)
    ..hasRequiredFields = false
  ;

  GetTransactionResponse._() : super();
  factory GetTransactionResponse({
    TransactionInfo? transaction,
    $core.int? blockHeight,
    $core.int? blockTime,
  }) {
    final _result = create();
    if (transaction != null) {
      _result.transaction = transaction;
    }
    if (blockHeight != null) {
      _result.blockHeight = blockHeight;
    }
    if (blockTime != null) {
      _result.blockTime = blockTime;
    }
    return _result;
  }
  factory GetTransactionResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetTransactionResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetTransactionResponse clone() => GetTransactionResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetTransactionResponse copyWith(void Function(GetTransactionResponse) updates) => super.copyWith((message) => updates(message as GetTransactionResponse)) as GetTransactionResponse; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static GetTransactionResponse create() => GetTransactionResponse._();
  GetTransactionResponse createEmptyInstance() => create();
  static $pb.PbList<GetTransactionResponse> createRepeated() => $pb.PbList<GetTransactionResponse>();
  @$core.pragma('dart2js:noInline')
  static GetTransactionResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetTransactionResponse>(create);
  static GetTransactionResponse? _defaultInstance;

  @$pb.TagNumber(3)
  TransactionInfo get transaction => $_getN(0);
  @$pb.TagNumber(3)
  set transaction(TransactionInfo v) { setField(3, v); }
  @$pb.TagNumber(3)
  $core.bool hasTransaction() => $_has(0);
  @$pb.TagNumber(3)
  void clearTransaction() => clearField(3);
  @$pb.TagNumber(3)
  TransactionInfo ensureTransaction() => $_ensure(0);

  @$pb.TagNumber(12)
  $core.int get blockHeight => $_getIZ(1);
  @$pb.TagNumber(12)
  set blockHeight($core.int v) { $_setUnsignedInt32(1, v); }
  @$pb.TagNumber(12)
  $core.bool hasBlockHeight() => $_has(1);
  @$pb.TagNumber(12)
  void clearBlockHeight() => clearField(12);

  @$pb.TagNumber(13)
  $core.int get blockTime => $_getIZ(2);
  @$pb.TagNumber(13)
  set blockTime($core.int v) { $_setUnsignedInt32(2, v); }
  @$pb.TagNumber(13)
  $core.bool hasBlockTime() => $_has(2);
  @$pb.TagNumber(13)
  void clearBlockTime() => clearField(13);
}

class CalculateFeeRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'CalculateFeeRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aInt64(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'amount')
    ..e<PayloadType>(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'payloadType', $pb.PbFieldType.OE, defaultOrMaker: PayloadType.UNKNOWN, valueOf: PayloadType.valueOf, enumValues: PayloadType.values)
    ..hasRequiredFields = false
  ;

  CalculateFeeRequest._() : super();
  factory CalculateFeeRequest({
    $fixnum.Int64? amount,
    PayloadType? payloadType,
  }) {
    final _result = create();
    if (amount != null) {
      _result.amount = amount;
    }
    if (payloadType != null) {
      _result.payloadType = payloadType;
    }
    return _result;
  }
  factory CalculateFeeRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory CalculateFeeRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  CalculateFeeRequest clone() => CalculateFeeRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  CalculateFeeRequest copyWith(void Function(CalculateFeeRequest) updates) => super.copyWith((message) => updates(message as CalculateFeeRequest)) as CalculateFeeRequest; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static CalculateFeeRequest create() => CalculateFeeRequest._();
  CalculateFeeRequest createEmptyInstance() => create();
  static $pb.PbList<CalculateFeeRequest> createRepeated() => $pb.PbList<CalculateFeeRequest>();
  @$core.pragma('dart2js:noInline')
  static CalculateFeeRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<CalculateFeeRequest>(create);
  static CalculateFeeRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $fixnum.Int64 get amount => $_getI64(0);
  @$pb.TagNumber(1)
  set amount($fixnum.Int64 v) { $_setInt64(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasAmount() => $_has(0);
  @$pb.TagNumber(1)
  void clearAmount() => clearField(1);

  @$pb.TagNumber(2)
  PayloadType get payloadType => $_getN(1);
  @$pb.TagNumber(2)
  set payloadType(PayloadType v) { setField(2, v); }
  @$pb.TagNumber(2)
  $core.bool hasPayloadType() => $_has(1);
  @$pb.TagNumber(2)
  void clearPayloadType() => clearField(2);
}

class CalculateFeeResponse extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'CalculateFeeResponse', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aInt64(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'fee')
    ..hasRequiredFields = false
  ;

  CalculateFeeResponse._() : super();
  factory CalculateFeeResponse({
    $fixnum.Int64? fee,
  }) {
    final _result = create();
    if (fee != null) {
      _result.fee = fee;
    }
    return _result;
  }
  factory CalculateFeeResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory CalculateFeeResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  CalculateFeeResponse clone() => CalculateFeeResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  CalculateFeeResponse copyWith(void Function(CalculateFeeResponse) updates) => super.copyWith((message) => updates(message as CalculateFeeResponse)) as CalculateFeeResponse; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static CalculateFeeResponse create() => CalculateFeeResponse._();
  CalculateFeeResponse createEmptyInstance() => create();
  static $pb.PbList<CalculateFeeResponse> createRepeated() => $pb.PbList<CalculateFeeResponse>();
  @$core.pragma('dart2js:noInline')
  static CalculateFeeResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<CalculateFeeResponse>(create);
  static CalculateFeeResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $fixnum.Int64 get fee => $_getI64(0);
  @$pb.TagNumber(1)
  set fee($fixnum.Int64 v) { $_setInt64(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasFee() => $_has(0);
  @$pb.TagNumber(1)
  void clearFee() => clearField(1);
}

class BroadcastTransactionRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'BroadcastTransactionRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..a<$core.List<$core.int>>(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'signedRawTransaction', $pb.PbFieldType.OY)
    ..hasRequiredFields = false
  ;

  BroadcastTransactionRequest._() : super();
  factory BroadcastTransactionRequest({
    $core.List<$core.int>? signedRawTransaction,
  }) {
    final _result = create();
    if (signedRawTransaction != null) {
      _result.signedRawTransaction = signedRawTransaction;
    }
    return _result;
  }
  factory BroadcastTransactionRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory BroadcastTransactionRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  BroadcastTransactionRequest clone() => BroadcastTransactionRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  BroadcastTransactionRequest copyWith(void Function(BroadcastTransactionRequest) updates) => super.copyWith((message) => updates(message as BroadcastTransactionRequest)) as BroadcastTransactionRequest; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static BroadcastTransactionRequest create() => BroadcastTransactionRequest._();
  BroadcastTransactionRequest createEmptyInstance() => create();
  static $pb.PbList<BroadcastTransactionRequest> createRepeated() => $pb.PbList<BroadcastTransactionRequest>();
  @$core.pragma('dart2js:noInline')
  static BroadcastTransactionRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<BroadcastTransactionRequest>(create);
  static BroadcastTransactionRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.List<$core.int> get signedRawTransaction => $_getN(0);
  @$pb.TagNumber(1)
  set signedRawTransaction($core.List<$core.int> v) { $_setBytes(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasSignedRawTransaction() => $_has(0);
  @$pb.TagNumber(1)
  void clearSignedRawTransaction() => clearField(1);
}

class BroadcastTransactionResponse extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'BroadcastTransactionResponse', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..a<$core.List<$core.int>>(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'id', $pb.PbFieldType.OY)
    ..hasRequiredFields = false
  ;

  BroadcastTransactionResponse._() : super();
  factory BroadcastTransactionResponse({
    $core.List<$core.int>? id,
  }) {
    final _result = create();
    if (id != null) {
      _result.id = id;
    }
    return _result;
  }
  factory BroadcastTransactionResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory BroadcastTransactionResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  BroadcastTransactionResponse clone() => BroadcastTransactionResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  BroadcastTransactionResponse copyWith(void Function(BroadcastTransactionResponse) updates) => super.copyWith((message) => updates(message as BroadcastTransactionResponse)) as BroadcastTransactionResponse; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static BroadcastTransactionResponse create() => BroadcastTransactionResponse._();
  BroadcastTransactionResponse createEmptyInstance() => create();
  static $pb.PbList<BroadcastTransactionResponse> createRepeated() => $pb.PbList<BroadcastTransactionResponse>();
  @$core.pragma('dart2js:noInline')
  static BroadcastTransactionResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<BroadcastTransactionResponse>(create);
  static BroadcastTransactionResponse? _defaultInstance;

  @$pb.TagNumber(2)
  $core.List<$core.int> get id => $_getN(0);
  @$pb.TagNumber(2)
  set id($core.List<$core.int> v) { $_setBytes(0, v); }
  @$pb.TagNumber(2)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(2)
  void clearId() => clearField(2);
}

class GetRawTransferTransactionRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'GetRawTransferTransactionRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..a<$core.int>(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'lockTime', $pb.PbFieldType.OU3)
    ..aOS(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'sender')
    ..aOS(3, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'receiver')
    ..aInt64(4, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'amount')
    ..aInt64(5, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'fee')
    ..aOS(6, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'memo')
    ..hasRequiredFields = false
  ;

  GetRawTransferTransactionRequest._() : super();
  factory GetRawTransferTransactionRequest({
    $core.int? lockTime,
    $core.String? sender,
    $core.String? receiver,
    $fixnum.Int64? amount,
    $fixnum.Int64? fee,
    $core.String? memo,
  }) {
    final _result = create();
    if (lockTime != null) {
      _result.lockTime = lockTime;
    }
    if (sender != null) {
      _result.sender = sender;
    }
    if (receiver != null) {
      _result.receiver = receiver;
    }
    if (amount != null) {
      _result.amount = amount;
    }
    if (fee != null) {
      _result.fee = fee;
    }
    if (memo != null) {
      _result.memo = memo;
    }
    return _result;
  }
  factory GetRawTransferTransactionRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetRawTransferTransactionRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetRawTransferTransactionRequest clone() => GetRawTransferTransactionRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetRawTransferTransactionRequest copyWith(void Function(GetRawTransferTransactionRequest) updates) => super.copyWith((message) => updates(message as GetRawTransferTransactionRequest)) as GetRawTransferTransactionRequest; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static GetRawTransferTransactionRequest create() => GetRawTransferTransactionRequest._();
  GetRawTransferTransactionRequest createEmptyInstance() => create();
  static $pb.PbList<GetRawTransferTransactionRequest> createRepeated() => $pb.PbList<GetRawTransferTransactionRequest>();
  @$core.pragma('dart2js:noInline')
  static GetRawTransferTransactionRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetRawTransferTransactionRequest>(create);
  static GetRawTransferTransactionRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.int get lockTime => $_getIZ(0);
  @$pb.TagNumber(1)
  set lockTime($core.int v) { $_setUnsignedInt32(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasLockTime() => $_has(0);
  @$pb.TagNumber(1)
  void clearLockTime() => clearField(1);

  @$pb.TagNumber(2)
  $core.String get sender => $_getSZ(1);
  @$pb.TagNumber(2)
  set sender($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasSender() => $_has(1);
  @$pb.TagNumber(2)
  void clearSender() => clearField(2);

  @$pb.TagNumber(3)
  $core.String get receiver => $_getSZ(2);
  @$pb.TagNumber(3)
  set receiver($core.String v) { $_setString(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasReceiver() => $_has(2);
  @$pb.TagNumber(3)
  void clearReceiver() => clearField(3);

  @$pb.TagNumber(4)
  $fixnum.Int64 get amount => $_getI64(3);
  @$pb.TagNumber(4)
  set amount($fixnum.Int64 v) { $_setInt64(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasAmount() => $_has(3);
  @$pb.TagNumber(4)
  void clearAmount() => clearField(4);

  @$pb.TagNumber(5)
  $fixnum.Int64 get fee => $_getI64(4);
  @$pb.TagNumber(5)
  set fee($fixnum.Int64 v) { $_setInt64(4, v); }
  @$pb.TagNumber(5)
  $core.bool hasFee() => $_has(4);
  @$pb.TagNumber(5)
  void clearFee() => clearField(5);

  @$pb.TagNumber(6)
  $core.String get memo => $_getSZ(5);
  @$pb.TagNumber(6)
  set memo($core.String v) { $_setString(5, v); }
  @$pb.TagNumber(6)
  $core.bool hasMemo() => $_has(5);
  @$pb.TagNumber(6)
  void clearMemo() => clearField(6);
}

class GetRawBondTransactionRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'GetRawBondTransactionRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..a<$core.int>(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'lockTime', $pb.PbFieldType.OU3)
    ..aOS(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'sender')
    ..aOS(3, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'receiver')
    ..aInt64(4, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'stake')
    ..aOS(5, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'publicKey')
    ..aInt64(6, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'fee')
    ..aOS(7, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'memo')
    ..hasRequiredFields = false
  ;

  GetRawBondTransactionRequest._() : super();
  factory GetRawBondTransactionRequest({
    $core.int? lockTime,
    $core.String? sender,
    $core.String? receiver,
    $fixnum.Int64? stake,
    $core.String? publicKey,
    $fixnum.Int64? fee,
    $core.String? memo,
  }) {
    final _result = create();
    if (lockTime != null) {
      _result.lockTime = lockTime;
    }
    if (sender != null) {
      _result.sender = sender;
    }
    if (receiver != null) {
      _result.receiver = receiver;
    }
    if (stake != null) {
      _result.stake = stake;
    }
    if (publicKey != null) {
      _result.publicKey = publicKey;
    }
    if (fee != null) {
      _result.fee = fee;
    }
    if (memo != null) {
      _result.memo = memo;
    }
    return _result;
  }
  factory GetRawBondTransactionRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetRawBondTransactionRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetRawBondTransactionRequest clone() => GetRawBondTransactionRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetRawBondTransactionRequest copyWith(void Function(GetRawBondTransactionRequest) updates) => super.copyWith((message) => updates(message as GetRawBondTransactionRequest)) as GetRawBondTransactionRequest; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static GetRawBondTransactionRequest create() => GetRawBondTransactionRequest._();
  GetRawBondTransactionRequest createEmptyInstance() => create();
  static $pb.PbList<GetRawBondTransactionRequest> createRepeated() => $pb.PbList<GetRawBondTransactionRequest>();
  @$core.pragma('dart2js:noInline')
  static GetRawBondTransactionRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetRawBondTransactionRequest>(create);
  static GetRawBondTransactionRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.int get lockTime => $_getIZ(0);
  @$pb.TagNumber(1)
  set lockTime($core.int v) { $_setUnsignedInt32(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasLockTime() => $_has(0);
  @$pb.TagNumber(1)
  void clearLockTime() => clearField(1);

  @$pb.TagNumber(2)
  $core.String get sender => $_getSZ(1);
  @$pb.TagNumber(2)
  set sender($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasSender() => $_has(1);
  @$pb.TagNumber(2)
  void clearSender() => clearField(2);

  @$pb.TagNumber(3)
  $core.String get receiver => $_getSZ(2);
  @$pb.TagNumber(3)
  set receiver($core.String v) { $_setString(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasReceiver() => $_has(2);
  @$pb.TagNumber(3)
  void clearReceiver() => clearField(3);

  @$pb.TagNumber(4)
  $fixnum.Int64 get stake => $_getI64(3);
  @$pb.TagNumber(4)
  set stake($fixnum.Int64 v) { $_setInt64(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasStake() => $_has(3);
  @$pb.TagNumber(4)
  void clearStake() => clearField(4);

  @$pb.TagNumber(5)
  $core.String get publicKey => $_getSZ(4);
  @$pb.TagNumber(5)
  set publicKey($core.String v) { $_setString(4, v); }
  @$pb.TagNumber(5)
  $core.bool hasPublicKey() => $_has(4);
  @$pb.TagNumber(5)
  void clearPublicKey() => clearField(5);

  @$pb.TagNumber(6)
  $fixnum.Int64 get fee => $_getI64(5);
  @$pb.TagNumber(6)
  set fee($fixnum.Int64 v) { $_setInt64(5, v); }
  @$pb.TagNumber(6)
  $core.bool hasFee() => $_has(5);
  @$pb.TagNumber(6)
  void clearFee() => clearField(6);

  @$pb.TagNumber(7)
  $core.String get memo => $_getSZ(6);
  @$pb.TagNumber(7)
  set memo($core.String v) { $_setString(6, v); }
  @$pb.TagNumber(7)
  $core.bool hasMemo() => $_has(6);
  @$pb.TagNumber(7)
  void clearMemo() => clearField(7);
}

class GetRawUnBondTransactionRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'GetRawUnBondTransactionRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..a<$core.int>(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'lockTime', $pb.PbFieldType.OU3)
    ..aOS(3, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'validatorAddress')
    ..aOS(4, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'memo')
    ..hasRequiredFields = false
  ;

  GetRawUnBondTransactionRequest._() : super();
  factory GetRawUnBondTransactionRequest({
    $core.int? lockTime,
    $core.String? validatorAddress,
    $core.String? memo,
  }) {
    final _result = create();
    if (lockTime != null) {
      _result.lockTime = lockTime;
    }
    if (validatorAddress != null) {
      _result.validatorAddress = validatorAddress;
    }
    if (memo != null) {
      _result.memo = memo;
    }
    return _result;
  }
  factory GetRawUnBondTransactionRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetRawUnBondTransactionRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetRawUnBondTransactionRequest clone() => GetRawUnBondTransactionRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetRawUnBondTransactionRequest copyWith(void Function(GetRawUnBondTransactionRequest) updates) => super.copyWith((message) => updates(message as GetRawUnBondTransactionRequest)) as GetRawUnBondTransactionRequest; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static GetRawUnBondTransactionRequest create() => GetRawUnBondTransactionRequest._();
  GetRawUnBondTransactionRequest createEmptyInstance() => create();
  static $pb.PbList<GetRawUnBondTransactionRequest> createRepeated() => $pb.PbList<GetRawUnBondTransactionRequest>();
  @$core.pragma('dart2js:noInline')
  static GetRawUnBondTransactionRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetRawUnBondTransactionRequest>(create);
  static GetRawUnBondTransactionRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.int get lockTime => $_getIZ(0);
  @$pb.TagNumber(1)
  set lockTime($core.int v) { $_setUnsignedInt32(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasLockTime() => $_has(0);
  @$pb.TagNumber(1)
  void clearLockTime() => clearField(1);

  @$pb.TagNumber(3)
  $core.String get validatorAddress => $_getSZ(1);
  @$pb.TagNumber(3)
  set validatorAddress($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(3)
  $core.bool hasValidatorAddress() => $_has(1);
  @$pb.TagNumber(3)
  void clearValidatorAddress() => clearField(3);

  @$pb.TagNumber(4)
  $core.String get memo => $_getSZ(2);
  @$pb.TagNumber(4)
  set memo($core.String v) { $_setString(2, v); }
  @$pb.TagNumber(4)
  $core.bool hasMemo() => $_has(2);
  @$pb.TagNumber(4)
  void clearMemo() => clearField(4);
}

class GetRawWithdrawTransactionRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'GetRawWithdrawTransactionRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..a<$core.int>(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'lockTime', $pb.PbFieldType.OU3)
    ..aOS(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'validatorAddress')
    ..aOS(3, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'accountAddress')
    ..aInt64(4, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'fee')
    ..aInt64(5, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'amount')
    ..aOS(6, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'memo')
    ..hasRequiredFields = false
  ;

  GetRawWithdrawTransactionRequest._() : super();
  factory GetRawWithdrawTransactionRequest({
    $core.int? lockTime,
    $core.String? validatorAddress,
    $core.String? accountAddress,
    $fixnum.Int64? fee,
    $fixnum.Int64? amount,
    $core.String? memo,
  }) {
    final _result = create();
    if (lockTime != null) {
      _result.lockTime = lockTime;
    }
    if (validatorAddress != null) {
      _result.validatorAddress = validatorAddress;
    }
    if (accountAddress != null) {
      _result.accountAddress = accountAddress;
    }
    if (fee != null) {
      _result.fee = fee;
    }
    if (amount != null) {
      _result.amount = amount;
    }
    if (memo != null) {
      _result.memo = memo;
    }
    return _result;
  }
  factory GetRawWithdrawTransactionRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetRawWithdrawTransactionRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetRawWithdrawTransactionRequest clone() => GetRawWithdrawTransactionRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetRawWithdrawTransactionRequest copyWith(void Function(GetRawWithdrawTransactionRequest) updates) => super.copyWith((message) => updates(message as GetRawWithdrawTransactionRequest)) as GetRawWithdrawTransactionRequest; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static GetRawWithdrawTransactionRequest create() => GetRawWithdrawTransactionRequest._();
  GetRawWithdrawTransactionRequest createEmptyInstance() => create();
  static $pb.PbList<GetRawWithdrawTransactionRequest> createRepeated() => $pb.PbList<GetRawWithdrawTransactionRequest>();
  @$core.pragma('dart2js:noInline')
  static GetRawWithdrawTransactionRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetRawWithdrawTransactionRequest>(create);
  static GetRawWithdrawTransactionRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.int get lockTime => $_getIZ(0);
  @$pb.TagNumber(1)
  set lockTime($core.int v) { $_setUnsignedInt32(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasLockTime() => $_has(0);
  @$pb.TagNumber(1)
  void clearLockTime() => clearField(1);

  @$pb.TagNumber(2)
  $core.String get validatorAddress => $_getSZ(1);
  @$pb.TagNumber(2)
  set validatorAddress($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasValidatorAddress() => $_has(1);
  @$pb.TagNumber(2)
  void clearValidatorAddress() => clearField(2);

  @$pb.TagNumber(3)
  $core.String get accountAddress => $_getSZ(2);
  @$pb.TagNumber(3)
  set accountAddress($core.String v) { $_setString(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasAccountAddress() => $_has(2);
  @$pb.TagNumber(3)
  void clearAccountAddress() => clearField(3);

  @$pb.TagNumber(4)
  $fixnum.Int64 get fee => $_getI64(3);
  @$pb.TagNumber(4)
  set fee($fixnum.Int64 v) { $_setInt64(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasFee() => $_has(3);
  @$pb.TagNumber(4)
  void clearFee() => clearField(4);

  @$pb.TagNumber(5)
  $fixnum.Int64 get amount => $_getI64(4);
  @$pb.TagNumber(5)
  set amount($fixnum.Int64 v) { $_setInt64(4, v); }
  @$pb.TagNumber(5)
  $core.bool hasAmount() => $_has(4);
  @$pb.TagNumber(5)
  void clearAmount() => clearField(5);

  @$pb.TagNumber(6)
  $core.String get memo => $_getSZ(5);
  @$pb.TagNumber(6)
  set memo($core.String v) { $_setString(5, v); }
  @$pb.TagNumber(6)
  $core.bool hasMemo() => $_has(5);
  @$pb.TagNumber(6)
  void clearMemo() => clearField(6);
}

class GetRawTransactionResponse extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'GetRawTransactionResponse', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..a<$core.List<$core.int>>(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'rawTransaction', $pb.PbFieldType.OY)
    ..hasRequiredFields = false
  ;

  GetRawTransactionResponse._() : super();
  factory GetRawTransactionResponse({
    $core.List<$core.int>? rawTransaction,
  }) {
    final _result = create();
    if (rawTransaction != null) {
      _result.rawTransaction = rawTransaction;
    }
    return _result;
  }
  factory GetRawTransactionResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetRawTransactionResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetRawTransactionResponse clone() => GetRawTransactionResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetRawTransactionResponse copyWith(void Function(GetRawTransactionResponse) updates) => super.copyWith((message) => updates(message as GetRawTransactionResponse)) as GetRawTransactionResponse; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static GetRawTransactionResponse create() => GetRawTransactionResponse._();
  GetRawTransactionResponse createEmptyInstance() => create();
  static $pb.PbList<GetRawTransactionResponse> createRepeated() => $pb.PbList<GetRawTransactionResponse>();
  @$core.pragma('dart2js:noInline')
  static GetRawTransactionResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetRawTransactionResponse>(create);
  static GetRawTransactionResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.List<$core.int> get rawTransaction => $_getN(0);
  @$pb.TagNumber(1)
  set rawTransaction($core.List<$core.int> v) { $_setBytes(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasRawTransaction() => $_has(0);
  @$pb.TagNumber(1)
  void clearRawTransaction() => clearField(1);
}

class PayloadTransfer extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'PayloadTransfer', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'sender')
    ..aOS(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'receiver')
    ..aInt64(3, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'amount')
    ..hasRequiredFields = false
  ;

  PayloadTransfer._() : super();
  factory PayloadTransfer({
    $core.String? sender,
    $core.String? receiver,
    $fixnum.Int64? amount,
  }) {
    final _result = create();
    if (sender != null) {
      _result.sender = sender;
    }
    if (receiver != null) {
      _result.receiver = receiver;
    }
    if (amount != null) {
      _result.amount = amount;
    }
    return _result;
  }
  factory PayloadTransfer.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory PayloadTransfer.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  PayloadTransfer clone() => PayloadTransfer()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  PayloadTransfer copyWith(void Function(PayloadTransfer) updates) => super.copyWith((message) => updates(message as PayloadTransfer)) as PayloadTransfer; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static PayloadTransfer create() => PayloadTransfer._();
  PayloadTransfer createEmptyInstance() => create();
  static $pb.PbList<PayloadTransfer> createRepeated() => $pb.PbList<PayloadTransfer>();
  @$core.pragma('dart2js:noInline')
  static PayloadTransfer getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<PayloadTransfer>(create);
  static PayloadTransfer? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get sender => $_getSZ(0);
  @$pb.TagNumber(1)
  set sender($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasSender() => $_has(0);
  @$pb.TagNumber(1)
  void clearSender() => clearField(1);

  @$pb.TagNumber(2)
  $core.String get receiver => $_getSZ(1);
  @$pb.TagNumber(2)
  set receiver($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasReceiver() => $_has(1);
  @$pb.TagNumber(2)
  void clearReceiver() => clearField(2);

  @$pb.TagNumber(3)
  $fixnum.Int64 get amount => $_getI64(2);
  @$pb.TagNumber(3)
  set amount($fixnum.Int64 v) { $_setInt64(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasAmount() => $_has(2);
  @$pb.TagNumber(3)
  void clearAmount() => clearField(3);
}

class PayloadBond extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'PayloadBond', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'sender')
    ..aOS(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'receiver')
    ..aInt64(3, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'stake')
    ..hasRequiredFields = false
  ;

  PayloadBond._() : super();
  factory PayloadBond({
    $core.String? sender,
    $core.String? receiver,
    $fixnum.Int64? stake,
  }) {
    final _result = create();
    if (sender != null) {
      _result.sender = sender;
    }
    if (receiver != null) {
      _result.receiver = receiver;
    }
    if (stake != null) {
      _result.stake = stake;
    }
    return _result;
  }
  factory PayloadBond.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory PayloadBond.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  PayloadBond clone() => PayloadBond()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  PayloadBond copyWith(void Function(PayloadBond) updates) => super.copyWith((message) => updates(message as PayloadBond)) as PayloadBond; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static PayloadBond create() => PayloadBond._();
  PayloadBond createEmptyInstance() => create();
  static $pb.PbList<PayloadBond> createRepeated() => $pb.PbList<PayloadBond>();
  @$core.pragma('dart2js:noInline')
  static PayloadBond getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<PayloadBond>(create);
  static PayloadBond? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get sender => $_getSZ(0);
  @$pb.TagNumber(1)
  set sender($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasSender() => $_has(0);
  @$pb.TagNumber(1)
  void clearSender() => clearField(1);

  @$pb.TagNumber(2)
  $core.String get receiver => $_getSZ(1);
  @$pb.TagNumber(2)
  set receiver($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasReceiver() => $_has(1);
  @$pb.TagNumber(2)
  void clearReceiver() => clearField(2);

  @$pb.TagNumber(3)
  $fixnum.Int64 get stake => $_getI64(2);
  @$pb.TagNumber(3)
  set stake($fixnum.Int64 v) { $_setInt64(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasStake() => $_has(2);
  @$pb.TagNumber(3)
  void clearStake() => clearField(3);
}

class PayloadSortition extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'PayloadSortition', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'address')
    ..a<$core.List<$core.int>>(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'proof', $pb.PbFieldType.OY)
    ..hasRequiredFields = false
  ;

  PayloadSortition._() : super();
  factory PayloadSortition({
    $core.String? address,
    $core.List<$core.int>? proof,
  }) {
    final _result = create();
    if (address != null) {
      _result.address = address;
    }
    if (proof != null) {
      _result.proof = proof;
    }
    return _result;
  }
  factory PayloadSortition.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory PayloadSortition.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  PayloadSortition clone() => PayloadSortition()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  PayloadSortition copyWith(void Function(PayloadSortition) updates) => super.copyWith((message) => updates(message as PayloadSortition)) as PayloadSortition; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static PayloadSortition create() => PayloadSortition._();
  PayloadSortition createEmptyInstance() => create();
  static $pb.PbList<PayloadSortition> createRepeated() => $pb.PbList<PayloadSortition>();
  @$core.pragma('dart2js:noInline')
  static PayloadSortition getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<PayloadSortition>(create);
  static PayloadSortition? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get address => $_getSZ(0);
  @$pb.TagNumber(1)
  set address($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasAddress() => $_has(0);
  @$pb.TagNumber(1)
  void clearAddress() => clearField(1);

  @$pb.TagNumber(2)
  $core.List<$core.int> get proof => $_getN(1);
  @$pb.TagNumber(2)
  set proof($core.List<$core.int> v) { $_setBytes(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasProof() => $_has(1);
  @$pb.TagNumber(2)
  void clearProof() => clearField(2);
}

class PayloadUnbond extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'PayloadUnbond', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'validator')
    ..hasRequiredFields = false
  ;

  PayloadUnbond._() : super();
  factory PayloadUnbond({
    $core.String? validator,
  }) {
    final _result = create();
    if (validator != null) {
      _result.validator = validator;
    }
    return _result;
  }
  factory PayloadUnbond.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory PayloadUnbond.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  PayloadUnbond clone() => PayloadUnbond()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  PayloadUnbond copyWith(void Function(PayloadUnbond) updates) => super.copyWith((message) => updates(message as PayloadUnbond)) as PayloadUnbond; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static PayloadUnbond create() => PayloadUnbond._();
  PayloadUnbond createEmptyInstance() => create();
  static $pb.PbList<PayloadUnbond> createRepeated() => $pb.PbList<PayloadUnbond>();
  @$core.pragma('dart2js:noInline')
  static PayloadUnbond getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<PayloadUnbond>(create);
  static PayloadUnbond? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get validator => $_getSZ(0);
  @$pb.TagNumber(1)
  set validator($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasValidator() => $_has(0);
  @$pb.TagNumber(1)
  void clearValidator() => clearField(1);
}

class PayloadWithdraw extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'PayloadWithdraw', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'from')
    ..aOS(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'to')
    ..aInt64(3, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'amount')
    ..hasRequiredFields = false
  ;

  PayloadWithdraw._() : super();
  factory PayloadWithdraw({
    $core.String? from,
    $core.String? to,
    $fixnum.Int64? amount,
  }) {
    final _result = create();
    if (from != null) {
      _result.from = from;
    }
    if (to != null) {
      _result.to = to;
    }
    if (amount != null) {
      _result.amount = amount;
    }
    return _result;
  }
  factory PayloadWithdraw.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory PayloadWithdraw.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  PayloadWithdraw clone() => PayloadWithdraw()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  PayloadWithdraw copyWith(void Function(PayloadWithdraw) updates) => super.copyWith((message) => updates(message as PayloadWithdraw)) as PayloadWithdraw; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static PayloadWithdraw create() => PayloadWithdraw._();
  PayloadWithdraw createEmptyInstance() => create();
  static $pb.PbList<PayloadWithdraw> createRepeated() => $pb.PbList<PayloadWithdraw>();
  @$core.pragma('dart2js:noInline')
  static PayloadWithdraw getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<PayloadWithdraw>(create);
  static PayloadWithdraw? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get from => $_getSZ(0);
  @$pb.TagNumber(1)
  set from($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasFrom() => $_has(0);
  @$pb.TagNumber(1)
  void clearFrom() => clearField(1);

  @$pb.TagNumber(2)
  $core.String get to => $_getSZ(1);
  @$pb.TagNumber(2)
  set to($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasTo() => $_has(1);
  @$pb.TagNumber(2)
  void clearTo() => clearField(2);

  @$pb.TagNumber(3)
  $fixnum.Int64 get amount => $_getI64(2);
  @$pb.TagNumber(3)
  set amount($fixnum.Int64 v) { $_setInt64(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasAmount() => $_has(2);
  @$pb.TagNumber(3)
  void clearAmount() => clearField(3);
}

enum TransactionInfo_Payload {
  transfer, 
  bond, 
  sortition, 
  unbond, 
  withdraw, 
  notSet
}

class TransactionInfo extends $pb.GeneratedMessage {
  static const $core.Map<$core.int, TransactionInfo_Payload> _TransactionInfo_PayloadByTag = {
    30 : TransactionInfo_Payload.transfer,
    31 : TransactionInfo_Payload.bond,
    32 : TransactionInfo_Payload.sortition,
    33 : TransactionInfo_Payload.unbond,
    34 : TransactionInfo_Payload.withdraw,
    0 : TransactionInfo_Payload.notSet
  };
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'TransactionInfo', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..oo(0, [30, 31, 32, 33, 34])
    ..a<$core.List<$core.int>>(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'id', $pb.PbFieldType.OY)
    ..a<$core.List<$core.int>>(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'data', $pb.PbFieldType.OY)
    ..a<$core.int>(3, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'version', $pb.PbFieldType.O3)
    ..a<$core.int>(4, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'lockTime', $pb.PbFieldType.OU3)
    ..aInt64(5, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'value')
    ..aInt64(6, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'fee')
    ..e<PayloadType>(7, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'payloadType', $pb.PbFieldType.OE, defaultOrMaker: PayloadType.UNKNOWN, valueOf: PayloadType.valueOf, enumValues: PayloadType.values)
    ..aOS(8, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'memo')
    ..aOS(9, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'publicKey')
    ..a<$core.List<$core.int>>(10, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'signature', $pb.PbFieldType.OY)
    ..aOM<PayloadTransfer>(30, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'transfer', subBuilder: PayloadTransfer.create)
    ..aOM<PayloadBond>(31, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'bond', subBuilder: PayloadBond.create)
    ..aOM<PayloadSortition>(32, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'sortition', subBuilder: PayloadSortition.create)
    ..aOM<PayloadUnbond>(33, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'unbond', subBuilder: PayloadUnbond.create)
    ..aOM<PayloadWithdraw>(34, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'withdraw', subBuilder: PayloadWithdraw.create)
    ..hasRequiredFields = false
  ;

  TransactionInfo._() : super();
  factory TransactionInfo({
    $core.List<$core.int>? id,
    $core.List<$core.int>? data,
    $core.int? version,
    $core.int? lockTime,
    $fixnum.Int64? value,
    $fixnum.Int64? fee,
    PayloadType? payloadType,
    $core.String? memo,
    $core.String? publicKey,
    $core.List<$core.int>? signature,
    PayloadTransfer? transfer,
    PayloadBond? bond,
    PayloadSortition? sortition,
    PayloadUnbond? unbond,
    PayloadWithdraw? withdraw,
  }) {
    final _result = create();
    if (id != null) {
      _result.id = id;
    }
    if (data != null) {
      _result.data = data;
    }
    if (version != null) {
      _result.version = version;
    }
    if (lockTime != null) {
      _result.lockTime = lockTime;
    }
    if (value != null) {
      _result.value = value;
    }
    if (fee != null) {
      _result.fee = fee;
    }
    if (payloadType != null) {
      _result.payloadType = payloadType;
    }
    if (memo != null) {
      _result.memo = memo;
    }
    if (publicKey != null) {
      _result.publicKey = publicKey;
    }
    if (signature != null) {
      _result.signature = signature;
    }
    if (transfer != null) {
      _result.transfer = transfer;
    }
    if (bond != null) {
      _result.bond = bond;
    }
    if (sortition != null) {
      _result.sortition = sortition;
    }
    if (unbond != null) {
      _result.unbond = unbond;
    }
    if (withdraw != null) {
      _result.withdraw = withdraw;
    }
    return _result;
  }
  factory TransactionInfo.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory TransactionInfo.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  TransactionInfo clone() => TransactionInfo()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  TransactionInfo copyWith(void Function(TransactionInfo) updates) => super.copyWith((message) => updates(message as TransactionInfo)) as TransactionInfo; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static TransactionInfo create() => TransactionInfo._();
  TransactionInfo createEmptyInstance() => create();
  static $pb.PbList<TransactionInfo> createRepeated() => $pb.PbList<TransactionInfo>();
  @$core.pragma('dart2js:noInline')
  static TransactionInfo getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<TransactionInfo>(create);
  static TransactionInfo? _defaultInstance;

  TransactionInfo_Payload whichPayload() => _TransactionInfo_PayloadByTag[$_whichOneof(0)]!;
  void clearPayload() => clearField($_whichOneof(0));

  @$pb.TagNumber(1)
  $core.List<$core.int> get id => $_getN(0);
  @$pb.TagNumber(1)
  set id($core.List<$core.int> v) { $_setBytes(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(1)
  void clearId() => clearField(1);

  @$pb.TagNumber(2)
  $core.List<$core.int> get data => $_getN(1);
  @$pb.TagNumber(2)
  set data($core.List<$core.int> v) { $_setBytes(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasData() => $_has(1);
  @$pb.TagNumber(2)
  void clearData() => clearField(2);

  @$pb.TagNumber(3)
  $core.int get version => $_getIZ(2);
  @$pb.TagNumber(3)
  set version($core.int v) { $_setSignedInt32(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasVersion() => $_has(2);
  @$pb.TagNumber(3)
  void clearVersion() => clearField(3);

  @$pb.TagNumber(4)
  $core.int get lockTime => $_getIZ(3);
  @$pb.TagNumber(4)
  set lockTime($core.int v) { $_setUnsignedInt32(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasLockTime() => $_has(3);
  @$pb.TagNumber(4)
  void clearLockTime() => clearField(4);

  @$pb.TagNumber(5)
  $fixnum.Int64 get value => $_getI64(4);
  @$pb.TagNumber(5)
  set value($fixnum.Int64 v) { $_setInt64(4, v); }
  @$pb.TagNumber(5)
  $core.bool hasValue() => $_has(4);
  @$pb.TagNumber(5)
  void clearValue() => clearField(5);

  @$pb.TagNumber(6)
  $fixnum.Int64 get fee => $_getI64(5);
  @$pb.TagNumber(6)
  set fee($fixnum.Int64 v) { $_setInt64(5, v); }
  @$pb.TagNumber(6)
  $core.bool hasFee() => $_has(5);
  @$pb.TagNumber(6)
  void clearFee() => clearField(6);

  @$pb.TagNumber(7)
  PayloadType get payloadType => $_getN(6);
  @$pb.TagNumber(7)
  set payloadType(PayloadType v) { setField(7, v); }
  @$pb.TagNumber(7)
  $core.bool hasPayloadType() => $_has(6);
  @$pb.TagNumber(7)
  void clearPayloadType() => clearField(7);

  @$pb.TagNumber(8)
  $core.String get memo => $_getSZ(7);
  @$pb.TagNumber(8)
  set memo($core.String v) { $_setString(7, v); }
  @$pb.TagNumber(8)
  $core.bool hasMemo() => $_has(7);
  @$pb.TagNumber(8)
  void clearMemo() => clearField(8);

  @$pb.TagNumber(9)
  $core.String get publicKey => $_getSZ(8);
  @$pb.TagNumber(9)
  set publicKey($core.String v) { $_setString(8, v); }
  @$pb.TagNumber(9)
  $core.bool hasPublicKey() => $_has(8);
  @$pb.TagNumber(9)
  void clearPublicKey() => clearField(9);

  @$pb.TagNumber(10)
  $core.List<$core.int> get signature => $_getN(9);
  @$pb.TagNumber(10)
  set signature($core.List<$core.int> v) { $_setBytes(9, v); }
  @$pb.TagNumber(10)
  $core.bool hasSignature() => $_has(9);
  @$pb.TagNumber(10)
  void clearSignature() => clearField(10);

  @$pb.TagNumber(30)
  PayloadTransfer get transfer => $_getN(10);
  @$pb.TagNumber(30)
  set transfer(PayloadTransfer v) { setField(30, v); }
  @$pb.TagNumber(30)
  $core.bool hasTransfer() => $_has(10);
  @$pb.TagNumber(30)
  void clearTransfer() => clearField(30);
  @$pb.TagNumber(30)
  PayloadTransfer ensureTransfer() => $_ensure(10);

  @$pb.TagNumber(31)
  PayloadBond get bond => $_getN(11);
  @$pb.TagNumber(31)
  set bond(PayloadBond v) { setField(31, v); }
  @$pb.TagNumber(31)
  $core.bool hasBond() => $_has(11);
  @$pb.TagNumber(31)
  void clearBond() => clearField(31);
  @$pb.TagNumber(31)
  PayloadBond ensureBond() => $_ensure(11);

  @$pb.TagNumber(32)
  PayloadSortition get sortition => $_getN(12);
  @$pb.TagNumber(32)
  set sortition(PayloadSortition v) { setField(32, v); }
  @$pb.TagNumber(32)
  $core.bool hasSortition() => $_has(12);
  @$pb.TagNumber(32)
  void clearSortition() => clearField(32);
  @$pb.TagNumber(32)
  PayloadSortition ensureSortition() => $_ensure(12);

  @$pb.TagNumber(33)
  PayloadUnbond get unbond => $_getN(13);
  @$pb.TagNumber(33)
  set unbond(PayloadUnbond v) { setField(33, v); }
  @$pb.TagNumber(33)
  $core.bool hasUnbond() => $_has(13);
  @$pb.TagNumber(33)
  void clearUnbond() => clearField(33);
  @$pb.TagNumber(33)
  PayloadUnbond ensureUnbond() => $_ensure(13);

  @$pb.TagNumber(34)
  PayloadWithdraw get withdraw => $_getN(14);
  @$pb.TagNumber(34)
  set withdraw(PayloadWithdraw v) { setField(34, v); }
  @$pb.TagNumber(34)
  $core.bool hasWithdraw() => $_has(14);
  @$pb.TagNumber(34)
  void clearWithdraw() => clearField(34);
  @$pb.TagNumber(34)
  PayloadWithdraw ensureWithdraw() => $_ensure(14);
}

class TransactionApi {
  $pb.RpcClient _client;
  TransactionApi(this._client);

  $async.Future<GetTransactionResponse> getTransaction($pb.ClientContext? ctx, GetTransactionRequest request) {
    var emptyResponse = GetTransactionResponse();
    return _client.invoke<GetTransactionResponse>(ctx, 'Transaction', 'GetTransaction', request, emptyResponse);
  }
  $async.Future<CalculateFeeResponse> calculateFee($pb.ClientContext? ctx, CalculateFeeRequest request) {
    var emptyResponse = CalculateFeeResponse();
    return _client.invoke<CalculateFeeResponse>(ctx, 'Transaction', 'CalculateFee', request, emptyResponse);
  }
  $async.Future<BroadcastTransactionResponse> broadcastTransaction($pb.ClientContext? ctx, BroadcastTransactionRequest request) {
    var emptyResponse = BroadcastTransactionResponse();
    return _client.invoke<BroadcastTransactionResponse>(ctx, 'Transaction', 'BroadcastTransaction', request, emptyResponse);
  }
  $async.Future<GetRawTransactionResponse> getRawTransferTransaction($pb.ClientContext? ctx, GetRawTransferTransactionRequest request) {
    var emptyResponse = GetRawTransactionResponse();
    return _client.invoke<GetRawTransactionResponse>(ctx, 'Transaction', 'GetRawTransferTransaction', request, emptyResponse);
  }
  $async.Future<GetRawTransactionResponse> getRawBondTransaction($pb.ClientContext? ctx, GetRawBondTransactionRequest request) {
    var emptyResponse = GetRawTransactionResponse();
    return _client.invoke<GetRawTransactionResponse>(ctx, 'Transaction', 'GetRawBondTransaction', request, emptyResponse);
  }
  $async.Future<GetRawTransactionResponse> getRawUnBondTransaction($pb.ClientContext? ctx, GetRawUnBondTransactionRequest request) {
    var emptyResponse = GetRawTransactionResponse();
    return _client.invoke<GetRawTransactionResponse>(ctx, 'Transaction', 'GetRawUnBondTransaction', request, emptyResponse);
  }
  $async.Future<GetRawTransactionResponse> getRawWithdrawTransaction($pb.ClientContext? ctx, GetRawWithdrawTransactionRequest request) {
    var emptyResponse = GetRawTransactionResponse();
    return _client.invoke<GetRawTransactionResponse>(ctx, 'Transaction', 'GetRawWithdrawTransaction', request, emptyResponse);
  }
}

