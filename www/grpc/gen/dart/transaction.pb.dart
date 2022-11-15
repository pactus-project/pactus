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

class SendRawTransactionRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'SendRawTransactionRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..a<$core.List<$core.int>>(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'data', $pb.PbFieldType.OY)
    ..hasRequiredFields = false
  ;

  SendRawTransactionRequest._() : super();
  factory SendRawTransactionRequest({
    $core.List<$core.int>? data,
  }) {
    final _result = create();
    if (data != null) {
      _result.data = data;
    }
    return _result;
  }
  factory SendRawTransactionRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory SendRawTransactionRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  SendRawTransactionRequest clone() => SendRawTransactionRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  SendRawTransactionRequest copyWith(void Function(SendRawTransactionRequest) updates) => super.copyWith((message) => updates(message as SendRawTransactionRequest)) as SendRawTransactionRequest; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static SendRawTransactionRequest create() => SendRawTransactionRequest._();
  SendRawTransactionRequest createEmptyInstance() => create();
  static $pb.PbList<SendRawTransactionRequest> createRepeated() => $pb.PbList<SendRawTransactionRequest>();
  @$core.pragma('dart2js:noInline')
  static SendRawTransactionRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<SendRawTransactionRequest>(create);
  static SendRawTransactionRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.List<$core.int> get data => $_getN(0);
  @$pb.TagNumber(1)
  set data($core.List<$core.int> v) { $_setBytes(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasData() => $_has(0);
  @$pb.TagNumber(1)
  void clearData() => clearField(1);
}

class SendRawTransactionResponse extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'SendRawTransactionResponse', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..a<$core.List<$core.int>>(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'id', $pb.PbFieldType.OY)
    ..hasRequiredFields = false
  ;

  SendRawTransactionResponse._() : super();
  factory SendRawTransactionResponse({
    $core.List<$core.int>? id,
  }) {
    final _result = create();
    if (id != null) {
      _result.id = id;
    }
    return _result;
  }
  factory SendRawTransactionResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory SendRawTransactionResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  SendRawTransactionResponse clone() => SendRawTransactionResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  SendRawTransactionResponse copyWith(void Function(SendRawTransactionResponse) updates) => super.copyWith((message) => updates(message as SendRawTransactionResponse)) as SendRawTransactionResponse; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static SendRawTransactionResponse create() => SendRawTransactionResponse._();
  SendRawTransactionResponse createEmptyInstance() => create();
  static $pb.PbList<SendRawTransactionResponse> createRepeated() => $pb.PbList<SendRawTransactionResponse>();
  @$core.pragma('dart2js:noInline')
  static SendRawTransactionResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<SendRawTransactionResponse>(create);
  static SendRawTransactionResponse? _defaultInstance;

  @$pb.TagNumber(2)
  $core.List<$core.int> get id => $_getN(0);
  @$pb.TagNumber(2)
  set id($core.List<$core.int> v) { $_setBytes(0, v); }
  @$pb.TagNumber(2)
  $core.bool hasId() => $_has(0);
  @$pb.TagNumber(2)
  void clearId() => clearField(2);
}

class PayloadSend extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'PayloadSend', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'sender')
    ..aOS(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'receiver')
    ..aInt64(3, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'amount')
    ..hasRequiredFields = false
  ;

  PayloadSend._() : super();
  factory PayloadSend({
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
  factory PayloadSend.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory PayloadSend.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  PayloadSend clone() => PayloadSend()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  PayloadSend copyWith(void Function(PayloadSend) updates) => super.copyWith((message) => updates(message as PayloadSend)) as PayloadSend; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static PayloadSend create() => PayloadSend._();
  PayloadSend createEmptyInstance() => create();
  static $pb.PbList<PayloadSend> createRepeated() => $pb.PbList<PayloadSend>();
  @$core.pragma('dart2js:noInline')
  static PayloadSend getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<PayloadSend>(create);
  static PayloadSend? _defaultInstance;

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
  send, 
  bond, 
  sortition, 
  unbond, 
  withdraw, 
  notSet
}

class TransactionInfo extends $pb.GeneratedMessage {
  static const $core.Map<$core.int, TransactionInfo_Payload> _TransactionInfo_PayloadByTag = {
    30 : TransactionInfo_Payload.send,
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
    ..a<$core.List<$core.int>>(4, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'stamp', $pb.PbFieldType.OY)
    ..a<$core.int>(5, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'sequence', $pb.PbFieldType.O3)
    ..aInt64(6, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'value')
    ..aInt64(7, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'fee')
    ..e<PayloadType>(8, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'Type', $pb.PbFieldType.OE, protoName: 'Type', defaultOrMaker: PayloadType.UNKNOWN, valueOf: PayloadType.valueOf, enumValues: PayloadType.values)
    ..aOS(9, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'memo')
    ..aOS(10, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'publicKey')
    ..a<$core.List<$core.int>>(11, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'signature', $pb.PbFieldType.OY)
    ..aOM<PayloadSend>(30, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'send', subBuilder: PayloadSend.create)
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
    $core.List<$core.int>? stamp,
    $core.int? sequence,
    $fixnum.Int64? value,
    $fixnum.Int64? fee,
    PayloadType? type,
    $core.String? memo,
    $core.String? publicKey,
    $core.List<$core.int>? signature,
    PayloadSend? send,
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
    if (stamp != null) {
      _result.stamp = stamp;
    }
    if (sequence != null) {
      _result.sequence = sequence;
    }
    if (value != null) {
      _result.value = value;
    }
    if (fee != null) {
      _result.fee = fee;
    }
    if (type != null) {
      _result.type = type;
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
    if (send != null) {
      _result.send = send;
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
  $core.List<$core.int> get stamp => $_getN(3);
  @$pb.TagNumber(4)
  set stamp($core.List<$core.int> v) { $_setBytes(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasStamp() => $_has(3);
  @$pb.TagNumber(4)
  void clearStamp() => clearField(4);

  @$pb.TagNumber(5)
  $core.int get sequence => $_getIZ(4);
  @$pb.TagNumber(5)
  set sequence($core.int v) { $_setSignedInt32(4, v); }
  @$pb.TagNumber(5)
  $core.bool hasSequence() => $_has(4);
  @$pb.TagNumber(5)
  void clearSequence() => clearField(5);

  @$pb.TagNumber(6)
  $fixnum.Int64 get value => $_getI64(5);
  @$pb.TagNumber(6)
  set value($fixnum.Int64 v) { $_setInt64(5, v); }
  @$pb.TagNumber(6)
  $core.bool hasValue() => $_has(5);
  @$pb.TagNumber(6)
  void clearValue() => clearField(6);

  @$pb.TagNumber(7)
  $fixnum.Int64 get fee => $_getI64(6);
  @$pb.TagNumber(7)
  set fee($fixnum.Int64 v) { $_setInt64(6, v); }
  @$pb.TagNumber(7)
  $core.bool hasFee() => $_has(6);
  @$pb.TagNumber(7)
  void clearFee() => clearField(7);

  @$pb.TagNumber(8)
  PayloadType get type => $_getN(7);
  @$pb.TagNumber(8)
  set type(PayloadType v) { setField(8, v); }
  @$pb.TagNumber(8)
  $core.bool hasType() => $_has(7);
  @$pb.TagNumber(8)
  void clearType() => clearField(8);

  @$pb.TagNumber(9)
  $core.String get memo => $_getSZ(8);
  @$pb.TagNumber(9)
  set memo($core.String v) { $_setString(8, v); }
  @$pb.TagNumber(9)
  $core.bool hasMemo() => $_has(8);
  @$pb.TagNumber(9)
  void clearMemo() => clearField(9);

  @$pb.TagNumber(10)
  $core.String get publicKey => $_getSZ(9);
  @$pb.TagNumber(10)
  set publicKey($core.String v) { $_setString(9, v); }
  @$pb.TagNumber(10)
  $core.bool hasPublicKey() => $_has(9);
  @$pb.TagNumber(10)
  void clearPublicKey() => clearField(10);

  @$pb.TagNumber(11)
  $core.List<$core.int> get signature => $_getN(10);
  @$pb.TagNumber(11)
  set signature($core.List<$core.int> v) { $_setBytes(10, v); }
  @$pb.TagNumber(11)
  $core.bool hasSignature() => $_has(10);
  @$pb.TagNumber(11)
  void clearSignature() => clearField(11);

  @$pb.TagNumber(30)
  PayloadSend get send => $_getN(11);
  @$pb.TagNumber(30)
  set send(PayloadSend v) { setField(30, v); }
  @$pb.TagNumber(30)
  $core.bool hasSend() => $_has(11);
  @$pb.TagNumber(30)
  void clearSend() => clearField(30);
  @$pb.TagNumber(30)
  PayloadSend ensureSend() => $_ensure(11);

  @$pb.TagNumber(31)
  PayloadBond get bond => $_getN(12);
  @$pb.TagNumber(31)
  set bond(PayloadBond v) { setField(31, v); }
  @$pb.TagNumber(31)
  $core.bool hasBond() => $_has(12);
  @$pb.TagNumber(31)
  void clearBond() => clearField(31);
  @$pb.TagNumber(31)
  PayloadBond ensureBond() => $_ensure(12);

  @$pb.TagNumber(32)
  PayloadSortition get sortition => $_getN(13);
  @$pb.TagNumber(32)
  set sortition(PayloadSortition v) { setField(32, v); }
  @$pb.TagNumber(32)
  $core.bool hasSortition() => $_has(13);
  @$pb.TagNumber(32)
  void clearSortition() => clearField(32);
  @$pb.TagNumber(32)
  PayloadSortition ensureSortition() => $_ensure(13);

  @$pb.TagNumber(33)
  PayloadUnbond get unbond => $_getN(14);
  @$pb.TagNumber(33)
  set unbond(PayloadUnbond v) { setField(33, v); }
  @$pb.TagNumber(33)
  $core.bool hasUnbond() => $_has(14);
  @$pb.TagNumber(33)
  void clearUnbond() => clearField(33);
  @$pb.TagNumber(33)
  PayloadUnbond ensureUnbond() => $_ensure(14);

  @$pb.TagNumber(34)
  PayloadWithdraw get withdraw => $_getN(15);
  @$pb.TagNumber(34)
  set withdraw(PayloadWithdraw v) { setField(34, v); }
  @$pb.TagNumber(34)
  $core.bool hasWithdraw() => $_has(15);
  @$pb.TagNumber(34)
  void clearWithdraw() => clearField(34);
  @$pb.TagNumber(34)
  PayloadWithdraw ensureWithdraw() => $_ensure(15);
}

class TransactionApi {
  $pb.RpcClient _client;
  TransactionApi(this._client);

  $async.Future<GetTransactionResponse> getTransaction($pb.ClientContext? ctx, GetTransactionRequest request) {
    var emptyResponse = GetTransactionResponse();
    return _client.invoke<GetTransactionResponse>(ctx, 'Transaction', 'GetTransaction', request, emptyResponse);
  }
  $async.Future<SendRawTransactionResponse> sendRawTransaction($pb.ClientContext? ctx, SendRawTransactionRequest request) {
    var emptyResponse = SendRawTransactionResponse();
    return _client.invoke<SendRawTransactionResponse>(ctx, 'Transaction', 'SendRawTransaction', request, emptyResponse);
  }
}

