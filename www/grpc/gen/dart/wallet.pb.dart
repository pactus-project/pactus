///
//  Generated code. Do not modify.
//  source: wallet.proto
//
// @dart = 2.12
// ignore_for_file: annotate_overrides,camel_case_types,constant_identifier_names,directives_ordering,library_prefixes,non_constant_identifier_names,prefer_final_fields,return_of_invalid_type,unnecessary_const,unnecessary_import,unnecessary_this,unused_import,unused_shown_name

import 'dart:async' as $async;
import 'dart:core' as $core;

import 'package:fixnum/fixnum.dart' as $fixnum;
import 'package:protobuf/protobuf.dart' as $pb;

import 'wallet.pbenum.dart';

export 'wallet.pbenum.dart';

class AddressInfo extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'AddressInfo', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'address')
    ..aOS(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'publicKey')
    ..aOS(3, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'label')
    ..aOS(4, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'path')
    ..hasRequiredFields = false
  ;

  AddressInfo._() : super();
  factory AddressInfo({
    $core.String? address,
    $core.String? publicKey,
    $core.String? label,
    $core.String? path,
  }) {
    final _result = create();
    if (address != null) {
      _result.address = address;
    }
    if (publicKey != null) {
      _result.publicKey = publicKey;
    }
    if (label != null) {
      _result.label = label;
    }
    if (path != null) {
      _result.path = path;
    }
    return _result;
  }
  factory AddressInfo.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory AddressInfo.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  AddressInfo clone() => AddressInfo()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  AddressInfo copyWith(void Function(AddressInfo) updates) => super.copyWith((message) => updates(message as AddressInfo)) as AddressInfo; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static AddressInfo create() => AddressInfo._();
  AddressInfo createEmptyInstance() => create();
  static $pb.PbList<AddressInfo> createRepeated() => $pb.PbList<AddressInfo>();
  @$core.pragma('dart2js:noInline')
  static AddressInfo getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<AddressInfo>(create);
  static AddressInfo? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get address => $_getSZ(0);
  @$pb.TagNumber(1)
  set address($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasAddress() => $_has(0);
  @$pb.TagNumber(1)
  void clearAddress() => clearField(1);

  @$pb.TagNumber(2)
  $core.String get publicKey => $_getSZ(1);
  @$pb.TagNumber(2)
  set publicKey($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasPublicKey() => $_has(1);
  @$pb.TagNumber(2)
  void clearPublicKey() => clearField(2);

  @$pb.TagNumber(3)
  $core.String get label => $_getSZ(2);
  @$pb.TagNumber(3)
  set label($core.String v) { $_setString(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasLabel() => $_has(2);
  @$pb.TagNumber(3)
  void clearLabel() => clearField(3);

  @$pb.TagNumber(4)
  $core.String get path => $_getSZ(3);
  @$pb.TagNumber(4)
  set path($core.String v) { $_setString(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasPath() => $_has(3);
  @$pb.TagNumber(4)
  void clearPath() => clearField(4);
}

class HistoryInfo extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'HistoryInfo', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'transactionId')
    ..a<$core.int>(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'time', $pb.PbFieldType.OU3)
    ..aOS(3, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'payloadType')
    ..aOS(4, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'description')
    ..aInt64(5, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'amount')
    ..hasRequiredFields = false
  ;

  HistoryInfo._() : super();
  factory HistoryInfo({
    $core.String? transactionId,
    $core.int? time,
    $core.String? payloadType,
    $core.String? description,
    $fixnum.Int64? amount,
  }) {
    final _result = create();
    if (transactionId != null) {
      _result.transactionId = transactionId;
    }
    if (time != null) {
      _result.time = time;
    }
    if (payloadType != null) {
      _result.payloadType = payloadType;
    }
    if (description != null) {
      _result.description = description;
    }
    if (amount != null) {
      _result.amount = amount;
    }
    return _result;
  }
  factory HistoryInfo.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory HistoryInfo.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  HistoryInfo clone() => HistoryInfo()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  HistoryInfo copyWith(void Function(HistoryInfo) updates) => super.copyWith((message) => updates(message as HistoryInfo)) as HistoryInfo; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static HistoryInfo create() => HistoryInfo._();
  HistoryInfo createEmptyInstance() => create();
  static $pb.PbList<HistoryInfo> createRepeated() => $pb.PbList<HistoryInfo>();
  @$core.pragma('dart2js:noInline')
  static HistoryInfo getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<HistoryInfo>(create);
  static HistoryInfo? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get transactionId => $_getSZ(0);
  @$pb.TagNumber(1)
  set transactionId($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasTransactionId() => $_has(0);
  @$pb.TagNumber(1)
  void clearTransactionId() => clearField(1);

  @$pb.TagNumber(2)
  $core.int get time => $_getIZ(1);
  @$pb.TagNumber(2)
  set time($core.int v) { $_setUnsignedInt32(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasTime() => $_has(1);
  @$pb.TagNumber(2)
  void clearTime() => clearField(2);

  @$pb.TagNumber(3)
  $core.String get payloadType => $_getSZ(2);
  @$pb.TagNumber(3)
  set payloadType($core.String v) { $_setString(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasPayloadType() => $_has(2);
  @$pb.TagNumber(3)
  void clearPayloadType() => clearField(3);

  @$pb.TagNumber(4)
  $core.String get description => $_getSZ(3);
  @$pb.TagNumber(4)
  set description($core.String v) { $_setString(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasDescription() => $_has(3);
  @$pb.TagNumber(4)
  void clearDescription() => clearField(4);

  @$pb.TagNumber(5)
  $fixnum.Int64 get amount => $_getI64(4);
  @$pb.TagNumber(5)
  set amount($fixnum.Int64 v) { $_setInt64(4, v); }
  @$pb.TagNumber(5)
  $core.bool hasAmount() => $_has(4);
  @$pb.TagNumber(5)
  void clearAmount() => clearField(5);
}

class GetAddressHistoryRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'GetAddressHistoryRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'walletName')
    ..aOS(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'address')
    ..hasRequiredFields = false
  ;

  GetAddressHistoryRequest._() : super();
  factory GetAddressHistoryRequest({
    $core.String? walletName,
    $core.String? address,
  }) {
    final _result = create();
    if (walletName != null) {
      _result.walletName = walletName;
    }
    if (address != null) {
      _result.address = address;
    }
    return _result;
  }
  factory GetAddressHistoryRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetAddressHistoryRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetAddressHistoryRequest clone() => GetAddressHistoryRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetAddressHistoryRequest copyWith(void Function(GetAddressHistoryRequest) updates) => super.copyWith((message) => updates(message as GetAddressHistoryRequest)) as GetAddressHistoryRequest; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static GetAddressHistoryRequest create() => GetAddressHistoryRequest._();
  GetAddressHistoryRequest createEmptyInstance() => create();
  static $pb.PbList<GetAddressHistoryRequest> createRepeated() => $pb.PbList<GetAddressHistoryRequest>();
  @$core.pragma('dart2js:noInline')
  static GetAddressHistoryRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetAddressHistoryRequest>(create);
  static GetAddressHistoryRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => clearField(1);

  @$pb.TagNumber(2)
  $core.String get address => $_getSZ(1);
  @$pb.TagNumber(2)
  set address($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasAddress() => $_has(1);
  @$pb.TagNumber(2)
  void clearAddress() => clearField(2);
}

class GetAddressHistoryResponse extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'GetAddressHistoryResponse', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..pc<HistoryInfo>(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'historyInfo', $pb.PbFieldType.PM, subBuilder: HistoryInfo.create)
    ..hasRequiredFields = false
  ;

  GetAddressHistoryResponse._() : super();
  factory GetAddressHistoryResponse({
    $core.Iterable<HistoryInfo>? historyInfo,
  }) {
    final _result = create();
    if (historyInfo != null) {
      _result.historyInfo.addAll(historyInfo);
    }
    return _result;
  }
  factory GetAddressHistoryResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetAddressHistoryResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetAddressHistoryResponse clone() => GetAddressHistoryResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetAddressHistoryResponse copyWith(void Function(GetAddressHistoryResponse) updates) => super.copyWith((message) => updates(message as GetAddressHistoryResponse)) as GetAddressHistoryResponse; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static GetAddressHistoryResponse create() => GetAddressHistoryResponse._();
  GetAddressHistoryResponse createEmptyInstance() => create();
  static $pb.PbList<GetAddressHistoryResponse> createRepeated() => $pb.PbList<GetAddressHistoryResponse>();
  @$core.pragma('dart2js:noInline')
  static GetAddressHistoryResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetAddressHistoryResponse>(create);
  static GetAddressHistoryResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.List<HistoryInfo> get historyInfo => $_getList(0);
}

class GetNewAddressRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'GetNewAddressRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'walletName')
    ..e<AddressType>(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'addressType', $pb.PbFieldType.OE, defaultOrMaker: AddressType.ADDRESS_TYPE_TREASURY, valueOf: AddressType.valueOf, enumValues: AddressType.values)
    ..aOS(3, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'label')
    ..aOS(4, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'password')
    ..hasRequiredFields = false
  ;

  GetNewAddressRequest._() : super();
  factory GetNewAddressRequest({
    $core.String? walletName,
    AddressType? addressType,
    $core.String? label,
    $core.String? password,
  }) {
    final _result = create();
    if (walletName != null) {
      _result.walletName = walletName;
    }
    if (addressType != null) {
      _result.addressType = addressType;
    }
    if (label != null) {
      _result.label = label;
    }
    if (password != null) {
      _result.password = password;
    }
    return _result;
  }
  factory GetNewAddressRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetNewAddressRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetNewAddressRequest clone() => GetNewAddressRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetNewAddressRequest copyWith(void Function(GetNewAddressRequest) updates) => super.copyWith((message) => updates(message as GetNewAddressRequest)) as GetNewAddressRequest; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static GetNewAddressRequest create() => GetNewAddressRequest._();
  GetNewAddressRequest createEmptyInstance() => create();
  static $pb.PbList<GetNewAddressRequest> createRepeated() => $pb.PbList<GetNewAddressRequest>();
  @$core.pragma('dart2js:noInline')
  static GetNewAddressRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetNewAddressRequest>(create);
  static GetNewAddressRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => clearField(1);

  @$pb.TagNumber(2)
  AddressType get addressType => $_getN(1);
  @$pb.TagNumber(2)
  set addressType(AddressType v) { setField(2, v); }
  @$pb.TagNumber(2)
  $core.bool hasAddressType() => $_has(1);
  @$pb.TagNumber(2)
  void clearAddressType() => clearField(2);

  @$pb.TagNumber(3)
  $core.String get label => $_getSZ(2);
  @$pb.TagNumber(3)
  set label($core.String v) { $_setString(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasLabel() => $_has(2);
  @$pb.TagNumber(3)
  void clearLabel() => clearField(3);

  @$pb.TagNumber(4)
  $core.String get password => $_getSZ(3);
  @$pb.TagNumber(4)
  set password($core.String v) { $_setString(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasPassword() => $_has(3);
  @$pb.TagNumber(4)
  void clearPassword() => clearField(4);
}

class GetNewAddressResponse extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'GetNewAddressResponse', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'walletName')
    ..aOM<AddressInfo>(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'addressInfo', subBuilder: AddressInfo.create)
    ..hasRequiredFields = false
  ;

  GetNewAddressResponse._() : super();
  factory GetNewAddressResponse({
    $core.String? walletName,
    AddressInfo? addressInfo,
  }) {
    final _result = create();
    if (walletName != null) {
      _result.walletName = walletName;
    }
    if (addressInfo != null) {
      _result.addressInfo = addressInfo;
    }
    return _result;
  }
  factory GetNewAddressResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetNewAddressResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetNewAddressResponse clone() => GetNewAddressResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetNewAddressResponse copyWith(void Function(GetNewAddressResponse) updates) => super.copyWith((message) => updates(message as GetNewAddressResponse)) as GetNewAddressResponse; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static GetNewAddressResponse create() => GetNewAddressResponse._();
  GetNewAddressResponse createEmptyInstance() => create();
  static $pb.PbList<GetNewAddressResponse> createRepeated() => $pb.PbList<GetNewAddressResponse>();
  @$core.pragma('dart2js:noInline')
  static GetNewAddressResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetNewAddressResponse>(create);
  static GetNewAddressResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => clearField(1);

  @$pb.TagNumber(2)
  AddressInfo get addressInfo => $_getN(1);
  @$pb.TagNumber(2)
  set addressInfo(AddressInfo v) { setField(2, v); }
  @$pb.TagNumber(2)
  $core.bool hasAddressInfo() => $_has(1);
  @$pb.TagNumber(2)
  void clearAddressInfo() => clearField(2);
  @$pb.TagNumber(2)
  AddressInfo ensureAddressInfo() => $_ensure(1);
}

class RestoreWalletRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'RestoreWalletRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'walletName')
    ..aOS(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'mnemonic')
    ..aOS(3, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'password')
    ..hasRequiredFields = false
  ;

  RestoreWalletRequest._() : super();
  factory RestoreWalletRequest({
    $core.String? walletName,
    $core.String? mnemonic,
    $core.String? password,
  }) {
    final _result = create();
    if (walletName != null) {
      _result.walletName = walletName;
    }
    if (mnemonic != null) {
      _result.mnemonic = mnemonic;
    }
    if (password != null) {
      _result.password = password;
    }
    return _result;
  }
  factory RestoreWalletRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory RestoreWalletRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  RestoreWalletRequest clone() => RestoreWalletRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  RestoreWalletRequest copyWith(void Function(RestoreWalletRequest) updates) => super.copyWith((message) => updates(message as RestoreWalletRequest)) as RestoreWalletRequest; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static RestoreWalletRequest create() => RestoreWalletRequest._();
  RestoreWalletRequest createEmptyInstance() => create();
  static $pb.PbList<RestoreWalletRequest> createRepeated() => $pb.PbList<RestoreWalletRequest>();
  @$core.pragma('dart2js:noInline')
  static RestoreWalletRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<RestoreWalletRequest>(create);
  static RestoreWalletRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => clearField(1);

  @$pb.TagNumber(2)
  $core.String get mnemonic => $_getSZ(1);
  @$pb.TagNumber(2)
  set mnemonic($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasMnemonic() => $_has(1);
  @$pb.TagNumber(2)
  void clearMnemonic() => clearField(2);

  @$pb.TagNumber(3)
  $core.String get password => $_getSZ(2);
  @$pb.TagNumber(3)
  set password($core.String v) { $_setString(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasPassword() => $_has(2);
  @$pb.TagNumber(3)
  void clearPassword() => clearField(3);
}

class RestoreWalletResponse extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'RestoreWalletResponse', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'walletName')
    ..hasRequiredFields = false
  ;

  RestoreWalletResponse._() : super();
  factory RestoreWalletResponse({
    $core.String? walletName,
  }) {
    final _result = create();
    if (walletName != null) {
      _result.walletName = walletName;
    }
    return _result;
  }
  factory RestoreWalletResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory RestoreWalletResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  RestoreWalletResponse clone() => RestoreWalletResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  RestoreWalletResponse copyWith(void Function(RestoreWalletResponse) updates) => super.copyWith((message) => updates(message as RestoreWalletResponse)) as RestoreWalletResponse; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static RestoreWalletResponse create() => RestoreWalletResponse._();
  RestoreWalletResponse createEmptyInstance() => create();
  static $pb.PbList<RestoreWalletResponse> createRepeated() => $pb.PbList<RestoreWalletResponse>();
  @$core.pragma('dart2js:noInline')
  static RestoreWalletResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<RestoreWalletResponse>(create);
  static RestoreWalletResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => clearField(1);
}

class CreateWalletRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'CreateWalletRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'walletName')
    ..aOS(4, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'password')
    ..hasRequiredFields = false
  ;

  CreateWalletRequest._() : super();
  factory CreateWalletRequest({
    $core.String? walletName,
    $core.String? password,
  }) {
    final _result = create();
    if (walletName != null) {
      _result.walletName = walletName;
    }
    if (password != null) {
      _result.password = password;
    }
    return _result;
  }
  factory CreateWalletRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory CreateWalletRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  CreateWalletRequest clone() => CreateWalletRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  CreateWalletRequest copyWith(void Function(CreateWalletRequest) updates) => super.copyWith((message) => updates(message as CreateWalletRequest)) as CreateWalletRequest; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static CreateWalletRequest create() => CreateWalletRequest._();
  CreateWalletRequest createEmptyInstance() => create();
  static $pb.PbList<CreateWalletRequest> createRepeated() => $pb.PbList<CreateWalletRequest>();
  @$core.pragma('dart2js:noInline')
  static CreateWalletRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<CreateWalletRequest>(create);
  static CreateWalletRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => clearField(1);

  @$pb.TagNumber(4)
  $core.String get password => $_getSZ(1);
  @$pb.TagNumber(4)
  set password($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(4)
  $core.bool hasPassword() => $_has(1);
  @$pb.TagNumber(4)
  void clearPassword() => clearField(4);
}

class CreateWalletResponse extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'CreateWalletResponse', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'mnemonic')
    ..hasRequiredFields = false
  ;

  CreateWalletResponse._() : super();
  factory CreateWalletResponse({
    $core.String? mnemonic,
  }) {
    final _result = create();
    if (mnemonic != null) {
      _result.mnemonic = mnemonic;
    }
    return _result;
  }
  factory CreateWalletResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory CreateWalletResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  CreateWalletResponse clone() => CreateWalletResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  CreateWalletResponse copyWith(void Function(CreateWalletResponse) updates) => super.copyWith((message) => updates(message as CreateWalletResponse)) as CreateWalletResponse; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static CreateWalletResponse create() => CreateWalletResponse._();
  CreateWalletResponse createEmptyInstance() => create();
  static $pb.PbList<CreateWalletResponse> createRepeated() => $pb.PbList<CreateWalletResponse>();
  @$core.pragma('dart2js:noInline')
  static CreateWalletResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<CreateWalletResponse>(create);
  static CreateWalletResponse? _defaultInstance;

  @$pb.TagNumber(2)
  $core.String get mnemonic => $_getSZ(0);
  @$pb.TagNumber(2)
  set mnemonic($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(2)
  $core.bool hasMnemonic() => $_has(0);
  @$pb.TagNumber(2)
  void clearMnemonic() => clearField(2);
}

class LoadWalletRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'LoadWalletRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'walletName')
    ..hasRequiredFields = false
  ;

  LoadWalletRequest._() : super();
  factory LoadWalletRequest({
    $core.String? walletName,
  }) {
    final _result = create();
    if (walletName != null) {
      _result.walletName = walletName;
    }
    return _result;
  }
  factory LoadWalletRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory LoadWalletRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  LoadWalletRequest clone() => LoadWalletRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  LoadWalletRequest copyWith(void Function(LoadWalletRequest) updates) => super.copyWith((message) => updates(message as LoadWalletRequest)) as LoadWalletRequest; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static LoadWalletRequest create() => LoadWalletRequest._();
  LoadWalletRequest createEmptyInstance() => create();
  static $pb.PbList<LoadWalletRequest> createRepeated() => $pb.PbList<LoadWalletRequest>();
  @$core.pragma('dart2js:noInline')
  static LoadWalletRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<LoadWalletRequest>(create);
  static LoadWalletRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => clearField(1);
}

class LoadWalletResponse extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'LoadWalletResponse', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'walletName')
    ..hasRequiredFields = false
  ;

  LoadWalletResponse._() : super();
  factory LoadWalletResponse({
    $core.String? walletName,
  }) {
    final _result = create();
    if (walletName != null) {
      _result.walletName = walletName;
    }
    return _result;
  }
  factory LoadWalletResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory LoadWalletResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  LoadWalletResponse clone() => LoadWalletResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  LoadWalletResponse copyWith(void Function(LoadWalletResponse) updates) => super.copyWith((message) => updates(message as LoadWalletResponse)) as LoadWalletResponse; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static LoadWalletResponse create() => LoadWalletResponse._();
  LoadWalletResponse createEmptyInstance() => create();
  static $pb.PbList<LoadWalletResponse> createRepeated() => $pb.PbList<LoadWalletResponse>();
  @$core.pragma('dart2js:noInline')
  static LoadWalletResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<LoadWalletResponse>(create);
  static LoadWalletResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => clearField(1);
}

class UnloadWalletRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'UnloadWalletRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'walletName')
    ..hasRequiredFields = false
  ;

  UnloadWalletRequest._() : super();
  factory UnloadWalletRequest({
    $core.String? walletName,
  }) {
    final _result = create();
    if (walletName != null) {
      _result.walletName = walletName;
    }
    return _result;
  }
  factory UnloadWalletRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory UnloadWalletRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  UnloadWalletRequest clone() => UnloadWalletRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  UnloadWalletRequest copyWith(void Function(UnloadWalletRequest) updates) => super.copyWith((message) => updates(message as UnloadWalletRequest)) as UnloadWalletRequest; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static UnloadWalletRequest create() => UnloadWalletRequest._();
  UnloadWalletRequest createEmptyInstance() => create();
  static $pb.PbList<UnloadWalletRequest> createRepeated() => $pb.PbList<UnloadWalletRequest>();
  @$core.pragma('dart2js:noInline')
  static UnloadWalletRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<UnloadWalletRequest>(create);
  static UnloadWalletRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => clearField(1);
}

class UnloadWalletResponse extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'UnloadWalletResponse', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'walletName')
    ..hasRequiredFields = false
  ;

  UnloadWalletResponse._() : super();
  factory UnloadWalletResponse({
    $core.String? walletName,
  }) {
    final _result = create();
    if (walletName != null) {
      _result.walletName = walletName;
    }
    return _result;
  }
  factory UnloadWalletResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory UnloadWalletResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  UnloadWalletResponse clone() => UnloadWalletResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  UnloadWalletResponse copyWith(void Function(UnloadWalletResponse) updates) => super.copyWith((message) => updates(message as UnloadWalletResponse)) as UnloadWalletResponse; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static UnloadWalletResponse create() => UnloadWalletResponse._();
  UnloadWalletResponse createEmptyInstance() => create();
  static $pb.PbList<UnloadWalletResponse> createRepeated() => $pb.PbList<UnloadWalletResponse>();
  @$core.pragma('dart2js:noInline')
  static UnloadWalletResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<UnloadWalletResponse>(create);
  static UnloadWalletResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => clearField(1);
}

class GetValidatorAddressRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'GetValidatorAddressRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'publicKey')
    ..hasRequiredFields = false
  ;

  GetValidatorAddressRequest._() : super();
  factory GetValidatorAddressRequest({
    $core.String? publicKey,
  }) {
    final _result = create();
    if (publicKey != null) {
      _result.publicKey = publicKey;
    }
    return _result;
  }
  factory GetValidatorAddressRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetValidatorAddressRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetValidatorAddressRequest clone() => GetValidatorAddressRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetValidatorAddressRequest copyWith(void Function(GetValidatorAddressRequest) updates) => super.copyWith((message) => updates(message as GetValidatorAddressRequest)) as GetValidatorAddressRequest; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static GetValidatorAddressRequest create() => GetValidatorAddressRequest._();
  GetValidatorAddressRequest createEmptyInstance() => create();
  static $pb.PbList<GetValidatorAddressRequest> createRepeated() => $pb.PbList<GetValidatorAddressRequest>();
  @$core.pragma('dart2js:noInline')
  static GetValidatorAddressRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetValidatorAddressRequest>(create);
  static GetValidatorAddressRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get publicKey => $_getSZ(0);
  @$pb.TagNumber(1)
  set publicKey($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasPublicKey() => $_has(0);
  @$pb.TagNumber(1)
  void clearPublicKey() => clearField(1);
}

class GetValidatorAddressResponse extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'GetValidatorAddressResponse', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'address')
    ..hasRequiredFields = false
  ;

  GetValidatorAddressResponse._() : super();
  factory GetValidatorAddressResponse({
    $core.String? address,
  }) {
    final _result = create();
    if (address != null) {
      _result.address = address;
    }
    return _result;
  }
  factory GetValidatorAddressResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetValidatorAddressResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetValidatorAddressResponse clone() => GetValidatorAddressResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetValidatorAddressResponse copyWith(void Function(GetValidatorAddressResponse) updates) => super.copyWith((message) => updates(message as GetValidatorAddressResponse)) as GetValidatorAddressResponse; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static GetValidatorAddressResponse create() => GetValidatorAddressResponse._();
  GetValidatorAddressResponse createEmptyInstance() => create();
  static $pb.PbList<GetValidatorAddressResponse> createRepeated() => $pb.PbList<GetValidatorAddressResponse>();
  @$core.pragma('dart2js:noInline')
  static GetValidatorAddressResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetValidatorAddressResponse>(create);
  static GetValidatorAddressResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get address => $_getSZ(0);
  @$pb.TagNumber(1)
  set address($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasAddress() => $_has(0);
  @$pb.TagNumber(1)
  void clearAddress() => clearField(1);
}

class SignRawTransactionRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'SignRawTransactionRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'walletName')
    ..aOS(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'rawTransaction')
    ..aOS(3, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'password')
    ..hasRequiredFields = false
  ;

  SignRawTransactionRequest._() : super();
  factory SignRawTransactionRequest({
    $core.String? walletName,
    $core.String? rawTransaction,
    $core.String? password,
  }) {
    final _result = create();
    if (walletName != null) {
      _result.walletName = walletName;
    }
    if (rawTransaction != null) {
      _result.rawTransaction = rawTransaction;
    }
    if (password != null) {
      _result.password = password;
    }
    return _result;
  }
  factory SignRawTransactionRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory SignRawTransactionRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  SignRawTransactionRequest clone() => SignRawTransactionRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  SignRawTransactionRequest copyWith(void Function(SignRawTransactionRequest) updates) => super.copyWith((message) => updates(message as SignRawTransactionRequest)) as SignRawTransactionRequest; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static SignRawTransactionRequest create() => SignRawTransactionRequest._();
  SignRawTransactionRequest createEmptyInstance() => create();
  static $pb.PbList<SignRawTransactionRequest> createRepeated() => $pb.PbList<SignRawTransactionRequest>();
  @$core.pragma('dart2js:noInline')
  static SignRawTransactionRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<SignRawTransactionRequest>(create);
  static SignRawTransactionRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => clearField(1);

  @$pb.TagNumber(2)
  $core.String get rawTransaction => $_getSZ(1);
  @$pb.TagNumber(2)
  set rawTransaction($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasRawTransaction() => $_has(1);
  @$pb.TagNumber(2)
  void clearRawTransaction() => clearField(2);

  @$pb.TagNumber(3)
  $core.String get password => $_getSZ(2);
  @$pb.TagNumber(3)
  set password($core.String v) { $_setString(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasPassword() => $_has(2);
  @$pb.TagNumber(3)
  void clearPassword() => clearField(3);
}

class SignRawTransactionResponse extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'SignRawTransactionResponse', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'transactionId')
    ..aOS(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'signedRawTransaction')
    ..hasRequiredFields = false
  ;

  SignRawTransactionResponse._() : super();
  factory SignRawTransactionResponse({
    $core.String? transactionId,
    $core.String? signedRawTransaction,
  }) {
    final _result = create();
    if (transactionId != null) {
      _result.transactionId = transactionId;
    }
    if (signedRawTransaction != null) {
      _result.signedRawTransaction = signedRawTransaction;
    }
    return _result;
  }
  factory SignRawTransactionResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory SignRawTransactionResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  SignRawTransactionResponse clone() => SignRawTransactionResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  SignRawTransactionResponse copyWith(void Function(SignRawTransactionResponse) updates) => super.copyWith((message) => updates(message as SignRawTransactionResponse)) as SignRawTransactionResponse; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static SignRawTransactionResponse create() => SignRawTransactionResponse._();
  SignRawTransactionResponse createEmptyInstance() => create();
  static $pb.PbList<SignRawTransactionResponse> createRepeated() => $pb.PbList<SignRawTransactionResponse>();
  @$core.pragma('dart2js:noInline')
  static SignRawTransactionResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<SignRawTransactionResponse>(create);
  static SignRawTransactionResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get transactionId => $_getSZ(0);
  @$pb.TagNumber(1)
  set transactionId($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasTransactionId() => $_has(0);
  @$pb.TagNumber(1)
  void clearTransactionId() => clearField(1);

  @$pb.TagNumber(2)
  $core.String get signedRawTransaction => $_getSZ(1);
  @$pb.TagNumber(2)
  set signedRawTransaction($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasSignedRawTransaction() => $_has(1);
  @$pb.TagNumber(2)
  void clearSignedRawTransaction() => clearField(2);
}

class GetTotalBalanceRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'GetTotalBalanceRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'walletName')
    ..hasRequiredFields = false
  ;

  GetTotalBalanceRequest._() : super();
  factory GetTotalBalanceRequest({
    $core.String? walletName,
  }) {
    final _result = create();
    if (walletName != null) {
      _result.walletName = walletName;
    }
    return _result;
  }
  factory GetTotalBalanceRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetTotalBalanceRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetTotalBalanceRequest clone() => GetTotalBalanceRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetTotalBalanceRequest copyWith(void Function(GetTotalBalanceRequest) updates) => super.copyWith((message) => updates(message as GetTotalBalanceRequest)) as GetTotalBalanceRequest; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static GetTotalBalanceRequest create() => GetTotalBalanceRequest._();
  GetTotalBalanceRequest createEmptyInstance() => create();
  static $pb.PbList<GetTotalBalanceRequest> createRepeated() => $pb.PbList<GetTotalBalanceRequest>();
  @$core.pragma('dart2js:noInline')
  static GetTotalBalanceRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetTotalBalanceRequest>(create);
  static GetTotalBalanceRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => clearField(1);
}

class GetTotalBalanceResponse extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'GetTotalBalanceResponse', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'walletName')
    ..aInt64(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'totalBalance')
    ..hasRequiredFields = false
  ;

  GetTotalBalanceResponse._() : super();
  factory GetTotalBalanceResponse({
    $core.String? walletName,
    $fixnum.Int64? totalBalance,
  }) {
    final _result = create();
    if (walletName != null) {
      _result.walletName = walletName;
    }
    if (totalBalance != null) {
      _result.totalBalance = totalBalance;
    }
    return _result;
  }
  factory GetTotalBalanceResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetTotalBalanceResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetTotalBalanceResponse clone() => GetTotalBalanceResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetTotalBalanceResponse copyWith(void Function(GetTotalBalanceResponse) updates) => super.copyWith((message) => updates(message as GetTotalBalanceResponse)) as GetTotalBalanceResponse; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static GetTotalBalanceResponse create() => GetTotalBalanceResponse._();
  GetTotalBalanceResponse createEmptyInstance() => create();
  static $pb.PbList<GetTotalBalanceResponse> createRepeated() => $pb.PbList<GetTotalBalanceResponse>();
  @$core.pragma('dart2js:noInline')
  static GetTotalBalanceResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetTotalBalanceResponse>(create);
  static GetTotalBalanceResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => clearField(1);

  @$pb.TagNumber(2)
  $fixnum.Int64 get totalBalance => $_getI64(1);
  @$pb.TagNumber(2)
  set totalBalance($fixnum.Int64 v) { $_setInt64(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasTotalBalance() => $_has(1);
  @$pb.TagNumber(2)
  void clearTotalBalance() => clearField(2);
}

class SignMessageRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'SignMessageRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'walletName')
    ..aOS(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'password')
    ..aOS(3, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'address')
    ..aOS(4, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'message')
    ..hasRequiredFields = false
  ;

  SignMessageRequest._() : super();
  factory SignMessageRequest({
    $core.String? walletName,
    $core.String? password,
    $core.String? address,
    $core.String? message,
  }) {
    final _result = create();
    if (walletName != null) {
      _result.walletName = walletName;
    }
    if (password != null) {
      _result.password = password;
    }
    if (address != null) {
      _result.address = address;
    }
    if (message != null) {
      _result.message = message;
    }
    return _result;
  }
  factory SignMessageRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory SignMessageRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  SignMessageRequest clone() => SignMessageRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  SignMessageRequest copyWith(void Function(SignMessageRequest) updates) => super.copyWith((message) => updates(message as SignMessageRequest)) as SignMessageRequest; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static SignMessageRequest create() => SignMessageRequest._();
  SignMessageRequest createEmptyInstance() => create();
  static $pb.PbList<SignMessageRequest> createRepeated() => $pb.PbList<SignMessageRequest>();
  @$core.pragma('dart2js:noInline')
  static SignMessageRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<SignMessageRequest>(create);
  static SignMessageRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => clearField(1);

  @$pb.TagNumber(2)
  $core.String get password => $_getSZ(1);
  @$pb.TagNumber(2)
  set password($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasPassword() => $_has(1);
  @$pb.TagNumber(2)
  void clearPassword() => clearField(2);

  @$pb.TagNumber(3)
  $core.String get address => $_getSZ(2);
  @$pb.TagNumber(3)
  set address($core.String v) { $_setString(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasAddress() => $_has(2);
  @$pb.TagNumber(3)
  void clearAddress() => clearField(3);

  @$pb.TagNumber(4)
  $core.String get message => $_getSZ(3);
  @$pb.TagNumber(4)
  set message($core.String v) { $_setString(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasMessage() => $_has(3);
  @$pb.TagNumber(4)
  void clearMessage() => clearField(4);
}

class SignMessageResponse extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'SignMessageResponse', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'signature')
    ..hasRequiredFields = false
  ;

  SignMessageResponse._() : super();
  factory SignMessageResponse({
    $core.String? signature,
  }) {
    final _result = create();
    if (signature != null) {
      _result.signature = signature;
    }
    return _result;
  }
  factory SignMessageResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory SignMessageResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  SignMessageResponse clone() => SignMessageResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  SignMessageResponse copyWith(void Function(SignMessageResponse) updates) => super.copyWith((message) => updates(message as SignMessageResponse)) as SignMessageResponse; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static SignMessageResponse create() => SignMessageResponse._();
  SignMessageResponse createEmptyInstance() => create();
  static $pb.PbList<SignMessageResponse> createRepeated() => $pb.PbList<SignMessageResponse>();
  @$core.pragma('dart2js:noInline')
  static SignMessageResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<SignMessageResponse>(create);
  static SignMessageResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get signature => $_getSZ(0);
  @$pb.TagNumber(1)
  set signature($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasSignature() => $_has(0);
  @$pb.TagNumber(1)
  void clearSignature() => clearField(1);
}

class GetTotalStakeRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'GetTotalStakeRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'walletName')
    ..hasRequiredFields = false
  ;

  GetTotalStakeRequest._() : super();
  factory GetTotalStakeRequest({
    $core.String? walletName,
  }) {
    final _result = create();
    if (walletName != null) {
      _result.walletName = walletName;
    }
    return _result;
  }
  factory GetTotalStakeRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetTotalStakeRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetTotalStakeRequest clone() => GetTotalStakeRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetTotalStakeRequest copyWith(void Function(GetTotalStakeRequest) updates) => super.copyWith((message) => updates(message as GetTotalStakeRequest)) as GetTotalStakeRequest; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static GetTotalStakeRequest create() => GetTotalStakeRequest._();
  GetTotalStakeRequest createEmptyInstance() => create();
  static $pb.PbList<GetTotalStakeRequest> createRepeated() => $pb.PbList<GetTotalStakeRequest>();
  @$core.pragma('dart2js:noInline')
  static GetTotalStakeRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetTotalStakeRequest>(create);
  static GetTotalStakeRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => clearField(1);
}

class GetTotalStakeResponse extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'GetTotalStakeResponse', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aInt64(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'totalStake')
    ..aOS(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'walletName')
    ..hasRequiredFields = false
  ;

  GetTotalStakeResponse._() : super();
  factory GetTotalStakeResponse({
    $fixnum.Int64? totalStake,
    $core.String? walletName,
  }) {
    final _result = create();
    if (totalStake != null) {
      _result.totalStake = totalStake;
    }
    if (walletName != null) {
      _result.walletName = walletName;
    }
    return _result;
  }
  factory GetTotalStakeResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetTotalStakeResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetTotalStakeResponse clone() => GetTotalStakeResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetTotalStakeResponse copyWith(void Function(GetTotalStakeResponse) updates) => super.copyWith((message) => updates(message as GetTotalStakeResponse)) as GetTotalStakeResponse; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static GetTotalStakeResponse create() => GetTotalStakeResponse._();
  GetTotalStakeResponse createEmptyInstance() => create();
  static $pb.PbList<GetTotalStakeResponse> createRepeated() => $pb.PbList<GetTotalStakeResponse>();
  @$core.pragma('dart2js:noInline')
  static GetTotalStakeResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetTotalStakeResponse>(create);
  static GetTotalStakeResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $fixnum.Int64 get totalStake => $_getI64(0);
  @$pb.TagNumber(1)
  set totalStake($fixnum.Int64 v) { $_setInt64(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasTotalStake() => $_has(0);
  @$pb.TagNumber(1)
  void clearTotalStake() => clearField(1);

  @$pb.TagNumber(2)
  $core.String get walletName => $_getSZ(1);
  @$pb.TagNumber(2)
  set walletName($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasWalletName() => $_has(1);
  @$pb.TagNumber(2)
  void clearWalletName() => clearField(2);
}

class GetAddressInfoRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'GetAddressInfoRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'walletName')
    ..aOS(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'address')
    ..hasRequiredFields = false
  ;

  GetAddressInfoRequest._() : super();
  factory GetAddressInfoRequest({
    $core.String? walletName,
    $core.String? address,
  }) {
    final _result = create();
    if (walletName != null) {
      _result.walletName = walletName;
    }
    if (address != null) {
      _result.address = address;
    }
    return _result;
  }
  factory GetAddressInfoRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetAddressInfoRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetAddressInfoRequest clone() => GetAddressInfoRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetAddressInfoRequest copyWith(void Function(GetAddressInfoRequest) updates) => super.copyWith((message) => updates(message as GetAddressInfoRequest)) as GetAddressInfoRequest; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static GetAddressInfoRequest create() => GetAddressInfoRequest._();
  GetAddressInfoRequest createEmptyInstance() => create();
  static $pb.PbList<GetAddressInfoRequest> createRepeated() => $pb.PbList<GetAddressInfoRequest>();
  @$core.pragma('dart2js:noInline')
  static GetAddressInfoRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetAddressInfoRequest>(create);
  static GetAddressInfoRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => clearField(1);

  @$pb.TagNumber(2)
  $core.String get address => $_getSZ(1);
  @$pb.TagNumber(2)
  set address($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasAddress() => $_has(1);
  @$pb.TagNumber(2)
  void clearAddress() => clearField(2);
}

class GetAddressInfoResponse extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'GetAddressInfoResponse', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'address')
    ..aOS(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'label')
    ..aOS(3, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'publicKey')
    ..aOS(4, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'path')
    ..aOS(5, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'walletName')
    ..hasRequiredFields = false
  ;

  GetAddressInfoResponse._() : super();
  factory GetAddressInfoResponse({
    $core.String? address,
    $core.String? label,
    $core.String? publicKey,
    $core.String? path,
    $core.String? walletName,
  }) {
    final _result = create();
    if (address != null) {
      _result.address = address;
    }
    if (label != null) {
      _result.label = label;
    }
    if (publicKey != null) {
      _result.publicKey = publicKey;
    }
    if (path != null) {
      _result.path = path;
    }
    if (walletName != null) {
      _result.walletName = walletName;
    }
    return _result;
  }
  factory GetAddressInfoResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetAddressInfoResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetAddressInfoResponse clone() => GetAddressInfoResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetAddressInfoResponse copyWith(void Function(GetAddressInfoResponse) updates) => super.copyWith((message) => updates(message as GetAddressInfoResponse)) as GetAddressInfoResponse; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static GetAddressInfoResponse create() => GetAddressInfoResponse._();
  GetAddressInfoResponse createEmptyInstance() => create();
  static $pb.PbList<GetAddressInfoResponse> createRepeated() => $pb.PbList<GetAddressInfoResponse>();
  @$core.pragma('dart2js:noInline')
  static GetAddressInfoResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetAddressInfoResponse>(create);
  static GetAddressInfoResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get address => $_getSZ(0);
  @$pb.TagNumber(1)
  set address($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasAddress() => $_has(0);
  @$pb.TagNumber(1)
  void clearAddress() => clearField(1);

  @$pb.TagNumber(2)
  $core.String get label => $_getSZ(1);
  @$pb.TagNumber(2)
  set label($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasLabel() => $_has(1);
  @$pb.TagNumber(2)
  void clearLabel() => clearField(2);

  @$pb.TagNumber(3)
  $core.String get publicKey => $_getSZ(2);
  @$pb.TagNumber(3)
  set publicKey($core.String v) { $_setString(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasPublicKey() => $_has(2);
  @$pb.TagNumber(3)
  void clearPublicKey() => clearField(3);

  @$pb.TagNumber(4)
  $core.String get path => $_getSZ(3);
  @$pb.TagNumber(4)
  set path($core.String v) { $_setString(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasPath() => $_has(3);
  @$pb.TagNumber(4)
  void clearPath() => clearField(4);

  @$pb.TagNumber(5)
  $core.String get walletName => $_getSZ(4);
  @$pb.TagNumber(5)
  set walletName($core.String v) { $_setString(4, v); }
  @$pb.TagNumber(5)
  $core.bool hasWalletName() => $_has(4);
  @$pb.TagNumber(5)
  void clearWalletName() => clearField(5);
}

class SetLabelRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'SetLabelRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'walletName')
    ..aOS(3, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'password')
    ..aOS(4, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'address')
    ..aOS(5, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'label')
    ..hasRequiredFields = false
  ;

  SetLabelRequest._() : super();
  factory SetLabelRequest({
    $core.String? walletName,
    $core.String? password,
    $core.String? address,
    $core.String? label,
  }) {
    final _result = create();
    if (walletName != null) {
      _result.walletName = walletName;
    }
    if (password != null) {
      _result.password = password;
    }
    if (address != null) {
      _result.address = address;
    }
    if (label != null) {
      _result.label = label;
    }
    return _result;
  }
  factory SetLabelRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory SetLabelRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  SetLabelRequest clone() => SetLabelRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  SetLabelRequest copyWith(void Function(SetLabelRequest) updates) => super.copyWith((message) => updates(message as SetLabelRequest)) as SetLabelRequest; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static SetLabelRequest create() => SetLabelRequest._();
  SetLabelRequest createEmptyInstance() => create();
  static $pb.PbList<SetLabelRequest> createRepeated() => $pb.PbList<SetLabelRequest>();
  @$core.pragma('dart2js:noInline')
  static SetLabelRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<SetLabelRequest>(create);
  static SetLabelRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => clearField(1);

  @$pb.TagNumber(3)
  $core.String get password => $_getSZ(1);
  @$pb.TagNumber(3)
  set password($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(3)
  $core.bool hasPassword() => $_has(1);
  @$pb.TagNumber(3)
  void clearPassword() => clearField(3);

  @$pb.TagNumber(4)
  $core.String get address => $_getSZ(2);
  @$pb.TagNumber(4)
  set address($core.String v) { $_setString(2, v); }
  @$pb.TagNumber(4)
  $core.bool hasAddress() => $_has(2);
  @$pb.TagNumber(4)
  void clearAddress() => clearField(4);

  @$pb.TagNumber(5)
  $core.String get label => $_getSZ(3);
  @$pb.TagNumber(5)
  set label($core.String v) { $_setString(3, v); }
  @$pb.TagNumber(5)
  $core.bool hasLabel() => $_has(3);
  @$pb.TagNumber(5)
  void clearLabel() => clearField(5);
}

class SetLabelResponse extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'SetLabelResponse', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..hasRequiredFields = false
  ;

  SetLabelResponse._() : super();
  factory SetLabelResponse() => create();
  factory SetLabelResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory SetLabelResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  SetLabelResponse clone() => SetLabelResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  SetLabelResponse copyWith(void Function(SetLabelResponse) updates) => super.copyWith((message) => updates(message as SetLabelResponse)) as SetLabelResponse; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static SetLabelResponse create() => SetLabelResponse._();
  SetLabelResponse createEmptyInstance() => create();
  static $pb.PbList<SetLabelResponse> createRepeated() => $pb.PbList<SetLabelResponse>();
  @$core.pragma('dart2js:noInline')
  static SetLabelResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<SetLabelResponse>(create);
  static SetLabelResponse? _defaultInstance;
}

class ListWalletRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'ListWalletRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..hasRequiredFields = false
  ;

  ListWalletRequest._() : super();
  factory ListWalletRequest() => create();
  factory ListWalletRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory ListWalletRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  ListWalletRequest clone() => ListWalletRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  ListWalletRequest copyWith(void Function(ListWalletRequest) updates) => super.copyWith((message) => updates(message as ListWalletRequest)) as ListWalletRequest; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static ListWalletRequest create() => ListWalletRequest._();
  ListWalletRequest createEmptyInstance() => create();
  static $pb.PbList<ListWalletRequest> createRepeated() => $pb.PbList<ListWalletRequest>();
  @$core.pragma('dart2js:noInline')
  static ListWalletRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ListWalletRequest>(create);
  static ListWalletRequest? _defaultInstance;
}

class ListWalletResponse extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'ListWalletResponse', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..pPS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'wallets')
    ..hasRequiredFields = false
  ;

  ListWalletResponse._() : super();
  factory ListWalletResponse({
    $core.Iterable<$core.String>? wallets,
  }) {
    final _result = create();
    if (wallets != null) {
      _result.wallets.addAll(wallets);
    }
    return _result;
  }
  factory ListWalletResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory ListWalletResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  ListWalletResponse clone() => ListWalletResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  ListWalletResponse copyWith(void Function(ListWalletResponse) updates) => super.copyWith((message) => updates(message as ListWalletResponse)) as ListWalletResponse; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static ListWalletResponse create() => ListWalletResponse._();
  ListWalletResponse createEmptyInstance() => create();
  static $pb.PbList<ListWalletResponse> createRepeated() => $pb.PbList<ListWalletResponse>();
  @$core.pragma('dart2js:noInline')
  static ListWalletResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ListWalletResponse>(create);
  static ListWalletResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.List<$core.String> get wallets => $_getList(0);
}

class GetWalletInfoRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'GetWalletInfoRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'walletName')
    ..hasRequiredFields = false
  ;

  GetWalletInfoRequest._() : super();
  factory GetWalletInfoRequest({
    $core.String? walletName,
  }) {
    final _result = create();
    if (walletName != null) {
      _result.walletName = walletName;
    }
    return _result;
  }
  factory GetWalletInfoRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetWalletInfoRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetWalletInfoRequest clone() => GetWalletInfoRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetWalletInfoRequest copyWith(void Function(GetWalletInfoRequest) updates) => super.copyWith((message) => updates(message as GetWalletInfoRequest)) as GetWalletInfoRequest; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static GetWalletInfoRequest create() => GetWalletInfoRequest._();
  GetWalletInfoRequest createEmptyInstance() => create();
  static $pb.PbList<GetWalletInfoRequest> createRepeated() => $pb.PbList<GetWalletInfoRequest>();
  @$core.pragma('dart2js:noInline')
  static GetWalletInfoRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetWalletInfoRequest>(create);
  static GetWalletInfoRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => clearField(1);
}

class GetWalletInfoResponse extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'GetWalletInfoResponse', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'walletName')
    ..aInt64(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'version')
    ..aOS(3, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'network')
    ..aOB(4, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'encrypted')
    ..aOS(5, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'uuid')
    ..aInt64(6, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'createdAt')
    ..hasRequiredFields = false
  ;

  GetWalletInfoResponse._() : super();
  factory GetWalletInfoResponse({
    $core.String? walletName,
    $fixnum.Int64? version,
    $core.String? network,
    $core.bool? encrypted,
    $core.String? uuid,
    $fixnum.Int64? createdAt,
  }) {
    final _result = create();
    if (walletName != null) {
      _result.walletName = walletName;
    }
    if (version != null) {
      _result.version = version;
    }
    if (network != null) {
      _result.network = network;
    }
    if (encrypted != null) {
      _result.encrypted = encrypted;
    }
    if (uuid != null) {
      _result.uuid = uuid;
    }
    if (createdAt != null) {
      _result.createdAt = createdAt;
    }
    return _result;
  }
  factory GetWalletInfoResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetWalletInfoResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetWalletInfoResponse clone() => GetWalletInfoResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetWalletInfoResponse copyWith(void Function(GetWalletInfoResponse) updates) => super.copyWith((message) => updates(message as GetWalletInfoResponse)) as GetWalletInfoResponse; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static GetWalletInfoResponse create() => GetWalletInfoResponse._();
  GetWalletInfoResponse createEmptyInstance() => create();
  static $pb.PbList<GetWalletInfoResponse> createRepeated() => $pb.PbList<GetWalletInfoResponse>();
  @$core.pragma('dart2js:noInline')
  static GetWalletInfoResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetWalletInfoResponse>(create);
  static GetWalletInfoResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => clearField(1);

  @$pb.TagNumber(2)
  $fixnum.Int64 get version => $_getI64(1);
  @$pb.TagNumber(2)
  set version($fixnum.Int64 v) { $_setInt64(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasVersion() => $_has(1);
  @$pb.TagNumber(2)
  void clearVersion() => clearField(2);

  @$pb.TagNumber(3)
  $core.String get network => $_getSZ(2);
  @$pb.TagNumber(3)
  set network($core.String v) { $_setString(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasNetwork() => $_has(2);
  @$pb.TagNumber(3)
  void clearNetwork() => clearField(3);

  @$pb.TagNumber(4)
  $core.bool get encrypted => $_getBF(3);
  @$pb.TagNumber(4)
  set encrypted($core.bool v) { $_setBool(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasEncrypted() => $_has(3);
  @$pb.TagNumber(4)
  void clearEncrypted() => clearField(4);

  @$pb.TagNumber(5)
  $core.String get uuid => $_getSZ(4);
  @$pb.TagNumber(5)
  set uuid($core.String v) { $_setString(4, v); }
  @$pb.TagNumber(5)
  $core.bool hasUuid() => $_has(4);
  @$pb.TagNumber(5)
  void clearUuid() => clearField(5);

  @$pb.TagNumber(6)
  $fixnum.Int64 get createdAt => $_getI64(5);
  @$pb.TagNumber(6)
  set createdAt($fixnum.Int64 v) { $_setInt64(5, v); }
  @$pb.TagNumber(6)
  $core.bool hasCreatedAt() => $_has(5);
  @$pb.TagNumber(6)
  void clearCreatedAt() => clearField(6);
}

class ListAddressRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'ListAddressRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'walletName')
    ..hasRequiredFields = false
  ;

  ListAddressRequest._() : super();
  factory ListAddressRequest({
    $core.String? walletName,
  }) {
    final _result = create();
    if (walletName != null) {
      _result.walletName = walletName;
    }
    return _result;
  }
  factory ListAddressRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory ListAddressRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  ListAddressRequest clone() => ListAddressRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  ListAddressRequest copyWith(void Function(ListAddressRequest) updates) => super.copyWith((message) => updates(message as ListAddressRequest)) as ListAddressRequest; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static ListAddressRequest create() => ListAddressRequest._();
  ListAddressRequest createEmptyInstance() => create();
  static $pb.PbList<ListAddressRequest> createRepeated() => $pb.PbList<ListAddressRequest>();
  @$core.pragma('dart2js:noInline')
  static ListAddressRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ListAddressRequest>(create);
  static ListAddressRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => clearField(1);
}

class ListAddressResponse extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'ListAddressResponse', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..pc<AddressInfo>(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'data', $pb.PbFieldType.PM, subBuilder: AddressInfo.create)
    ..hasRequiredFields = false
  ;

  ListAddressResponse._() : super();
  factory ListAddressResponse({
    $core.Iterable<AddressInfo>? data,
  }) {
    final _result = create();
    if (data != null) {
      _result.data.addAll(data);
    }
    return _result;
  }
  factory ListAddressResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory ListAddressResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  ListAddressResponse clone() => ListAddressResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  ListAddressResponse copyWith(void Function(ListAddressResponse) updates) => super.copyWith((message) => updates(message as ListAddressResponse)) as ListAddressResponse; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static ListAddressResponse create() => ListAddressResponse._();
  ListAddressResponse createEmptyInstance() => create();
  static $pb.PbList<ListAddressResponse> createRepeated() => $pb.PbList<ListAddressResponse>();
  @$core.pragma('dart2js:noInline')
  static ListAddressResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ListAddressResponse>(create);
  static ListAddressResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.List<AddressInfo> get data => $_getList(0);
}

class WalletApi {
  $pb.RpcClient _client;
  WalletApi(this._client);

  $async.Future<CreateWalletResponse> createWallet($pb.ClientContext? ctx, CreateWalletRequest request) {
    var emptyResponse = CreateWalletResponse();
    return _client.invoke<CreateWalletResponse>(ctx, 'Wallet', 'CreateWallet', request, emptyResponse);
  }
  $async.Future<RestoreWalletResponse> restoreWallet($pb.ClientContext? ctx, RestoreWalletRequest request) {
    var emptyResponse = RestoreWalletResponse();
    return _client.invoke<RestoreWalletResponse>(ctx, 'Wallet', 'RestoreWallet', request, emptyResponse);
  }
  $async.Future<LoadWalletResponse> loadWallet($pb.ClientContext? ctx, LoadWalletRequest request) {
    var emptyResponse = LoadWalletResponse();
    return _client.invoke<LoadWalletResponse>(ctx, 'Wallet', 'LoadWallet', request, emptyResponse);
  }
  $async.Future<UnloadWalletResponse> unloadWallet($pb.ClientContext? ctx, UnloadWalletRequest request) {
    var emptyResponse = UnloadWalletResponse();
    return _client.invoke<UnloadWalletResponse>(ctx, 'Wallet', 'UnloadWallet', request, emptyResponse);
  }
  $async.Future<GetTotalBalanceResponse> getTotalBalance($pb.ClientContext? ctx, GetTotalBalanceRequest request) {
    var emptyResponse = GetTotalBalanceResponse();
    return _client.invoke<GetTotalBalanceResponse>(ctx, 'Wallet', 'GetTotalBalance', request, emptyResponse);
  }
  $async.Future<SignRawTransactionResponse> signRawTransaction($pb.ClientContext? ctx, SignRawTransactionRequest request) {
    var emptyResponse = SignRawTransactionResponse();
    return _client.invoke<SignRawTransactionResponse>(ctx, 'Wallet', 'SignRawTransaction', request, emptyResponse);
  }
  $async.Future<GetValidatorAddressResponse> getValidatorAddress($pb.ClientContext? ctx, GetValidatorAddressRequest request) {
    var emptyResponse = GetValidatorAddressResponse();
    return _client.invoke<GetValidatorAddressResponse>(ctx, 'Wallet', 'GetValidatorAddress', request, emptyResponse);
  }
  $async.Future<GetNewAddressResponse> getNewAddress($pb.ClientContext? ctx, GetNewAddressRequest request) {
    var emptyResponse = GetNewAddressResponse();
    return _client.invoke<GetNewAddressResponse>(ctx, 'Wallet', 'GetNewAddress', request, emptyResponse);
  }
  $async.Future<GetAddressHistoryResponse> getAddressHistory($pb.ClientContext? ctx, GetAddressHistoryRequest request) {
    var emptyResponse = GetAddressHistoryResponse();
    return _client.invoke<GetAddressHistoryResponse>(ctx, 'Wallet', 'GetAddressHistory', request, emptyResponse);
  }
  $async.Future<SignMessageResponse> signMessage($pb.ClientContext? ctx, SignMessageRequest request) {
    var emptyResponse = SignMessageResponse();
    return _client.invoke<SignMessageResponse>(ctx, 'Wallet', 'SignMessage', request, emptyResponse);
  }
  $async.Future<GetTotalStakeResponse> getTotalStake($pb.ClientContext? ctx, GetTotalStakeRequest request) {
    var emptyResponse = GetTotalStakeResponse();
    return _client.invoke<GetTotalStakeResponse>(ctx, 'Wallet', 'GetTotalStake', request, emptyResponse);
  }
  $async.Future<GetAddressInfoResponse> getAddressInfo($pb.ClientContext? ctx, GetAddressInfoRequest request) {
    var emptyResponse = GetAddressInfoResponse();
    return _client.invoke<GetAddressInfoResponse>(ctx, 'Wallet', 'GetAddressInfo', request, emptyResponse);
  }
  $async.Future<SetLabelResponse> setAddressLabel($pb.ClientContext? ctx, SetLabelRequest request) {
    var emptyResponse = SetLabelResponse();
    return _client.invoke<SetLabelResponse>(ctx, 'Wallet', 'SetAddressLabel', request, emptyResponse);
  }
  $async.Future<ListWalletResponse> listWallet($pb.ClientContext? ctx, ListWalletRequest request) {
    var emptyResponse = ListWalletResponse();
    return _client.invoke<ListWalletResponse>(ctx, 'Wallet', 'ListWallet', request, emptyResponse);
  }
  $async.Future<GetWalletInfoResponse> getWalletInfo($pb.ClientContext? ctx, GetWalletInfoRequest request) {
    var emptyResponse = GetWalletInfoResponse();
    return _client.invoke<GetWalletInfoResponse>(ctx, 'Wallet', 'GetWalletInfo', request, emptyResponse);
  }
  $async.Future<ListAddressResponse> listAddress($pb.ClientContext? ctx, ListAddressRequest request) {
    var emptyResponse = ListAddressResponse();
    return _client.invoke<ListAddressResponse>(ctx, 'Wallet', 'ListAddress', request, emptyResponse);
  }
}

