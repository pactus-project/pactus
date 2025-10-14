//
//  Generated code. Do not modify.
//  source: wallet.proto
//
// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_final_fields
// ignore_for_file: unnecessary_import, unnecessary_this, unused_import

import 'dart:async' as $async;
import 'dart:core' as $core;

import 'package:fixnum/fixnum.dart' as $fixnum;
import 'package:protobuf/protobuf.dart' as $pb;

import 'wallet.pbenum.dart';

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

export 'wallet.pbenum.dart';

/// AddressInfo contains detailed information about a wallet address.
class AddressInfo extends $pb.GeneratedMessage {
  factory AddressInfo({
    $core.String? address,
    $core.String? publicKey,
    $core.String? label,
    $core.String? path,
  }) {
    final $result = create();
    if (address != null) {
      $result.address = address;
    }
    if (publicKey != null) {
      $result.publicKey = publicKey;
    }
    if (label != null) {
      $result.label = label;
    }
    if (path != null) {
      $result.path = path;
    }
    return $result;
  }
  AddressInfo._() : super();
  factory AddressInfo.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory AddressInfo.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'AddressInfo', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'address')
    ..aOS(2, _omitFieldNames ? '' : 'publicKey')
    ..aOS(3, _omitFieldNames ? '' : 'label')
    ..aOS(4, _omitFieldNames ? '' : 'path')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  AddressInfo clone() => AddressInfo()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  AddressInfo copyWith(void Function(AddressInfo) updates) => super.copyWith((message) => updates(message as AddressInfo)) as AddressInfo;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AddressInfo create() => AddressInfo._();
  AddressInfo createEmptyInstance() => create();
  static $pb.PbList<AddressInfo> createRepeated() => $pb.PbList<AddressInfo>();
  @$core.pragma('dart2js:noInline')
  static AddressInfo getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<AddressInfo>(create);
  static AddressInfo? _defaultInstance;

  /// The address string.
  @$pb.TagNumber(1)
  $core.String get address => $_getSZ(0);
  @$pb.TagNumber(1)
  set address($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasAddress() => $_has(0);
  @$pb.TagNumber(1)
  void clearAddress() => $_clearField(1);

  /// The public key associated with the address.
  @$pb.TagNumber(2)
  $core.String get publicKey => $_getSZ(1);
  @$pb.TagNumber(2)
  set publicKey($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasPublicKey() => $_has(1);
  @$pb.TagNumber(2)
  void clearPublicKey() => $_clearField(2);

  /// A human-readable label associated with the address.
  @$pb.TagNumber(3)
  $core.String get label => $_getSZ(2);
  @$pb.TagNumber(3)
  set label($core.String v) { $_setString(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasLabel() => $_has(2);
  @$pb.TagNumber(3)
  void clearLabel() => $_clearField(3);

  /// The Hierarchical Deterministic (HD) path of the address within the wallet.
  @$pb.TagNumber(4)
  $core.String get path => $_getSZ(3);
  @$pb.TagNumber(4)
  set path($core.String v) { $_setString(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasPath() => $_has(3);
  @$pb.TagNumber(4)
  void clearPath() => $_clearField(4);
}

/// HistoryInfo contains transaction history details for an address.
class HistoryInfo extends $pb.GeneratedMessage {
  factory HistoryInfo({
    $core.String? transactionId,
    $core.int? time,
    $core.String? payloadType,
    $core.String? description,
    $fixnum.Int64? amount,
  }) {
    final $result = create();
    if (transactionId != null) {
      $result.transactionId = transactionId;
    }
    if (time != null) {
      $result.time = time;
    }
    if (payloadType != null) {
      $result.payloadType = payloadType;
    }
    if (description != null) {
      $result.description = description;
    }
    if (amount != null) {
      $result.amount = amount;
    }
    return $result;
  }
  HistoryInfo._() : super();
  factory HistoryInfo.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory HistoryInfo.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'HistoryInfo', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'transactionId')
    ..a<$core.int>(2, _omitFieldNames ? '' : 'time', $pb.PbFieldType.OU3)
    ..aOS(3, _omitFieldNames ? '' : 'payloadType')
    ..aOS(4, _omitFieldNames ? '' : 'description')
    ..aInt64(5, _omitFieldNames ? '' : 'amount')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  HistoryInfo clone() => HistoryInfo()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  HistoryInfo copyWith(void Function(HistoryInfo) updates) => super.copyWith((message) => updates(message as HistoryInfo)) as HistoryInfo;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static HistoryInfo create() => HistoryInfo._();
  HistoryInfo createEmptyInstance() => create();
  static $pb.PbList<HistoryInfo> createRepeated() => $pb.PbList<HistoryInfo>();
  @$core.pragma('dart2js:noInline')
  static HistoryInfo getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<HistoryInfo>(create);
  static HistoryInfo? _defaultInstance;

  /// The transaction ID in hexadecimal format.
  @$pb.TagNumber(1)
  $core.String get transactionId => $_getSZ(0);
  @$pb.TagNumber(1)
  set transactionId($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasTransactionId() => $_has(0);
  @$pb.TagNumber(1)
  void clearTransactionId() => $_clearField(1);

  /// Unix timestamp of when the transaction was confirmed.
  @$pb.TagNumber(2)
  $core.int get time => $_getIZ(1);
  @$pb.TagNumber(2)
  set time($core.int v) { $_setUnsignedInt32(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasTime() => $_has(1);
  @$pb.TagNumber(2)
  void clearTime() => $_clearField(2);

  /// The type of transaction payload.
  @$pb.TagNumber(3)
  $core.String get payloadType => $_getSZ(2);
  @$pb.TagNumber(3)
  set payloadType($core.String v) { $_setString(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasPayloadType() => $_has(2);
  @$pb.TagNumber(3)
  void clearPayloadType() => $_clearField(3);

  /// Human-readable description of the transaction.
  @$pb.TagNumber(4)
  $core.String get description => $_getSZ(3);
  @$pb.TagNumber(4)
  set description($core.String v) { $_setString(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasDescription() => $_has(3);
  @$pb.TagNumber(4)
  void clearDescription() => $_clearField(4);

  /// The transaction amount in NanoPAC.
  @$pb.TagNumber(5)
  $fixnum.Int64 get amount => $_getI64(4);
  @$pb.TagNumber(5)
  set amount($fixnum.Int64 v) { $_setInt64(4, v); }
  @$pb.TagNumber(5)
  $core.bool hasAmount() => $_has(4);
  @$pb.TagNumber(5)
  void clearAmount() => $_clearField(5);
}

/// Request message for retrieving address transaction history.
class GetAddressHistoryRequest extends $pb.GeneratedMessage {
  factory GetAddressHistoryRequest({
    $core.String? walletName,
    $core.String? address,
  }) {
    final $result = create();
    if (walletName != null) {
      $result.walletName = walletName;
    }
    if (address != null) {
      $result.address = address;
    }
    return $result;
  }
  GetAddressHistoryRequest._() : super();
  factory GetAddressHistoryRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetAddressHistoryRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetAddressHistoryRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'walletName')
    ..aOS(2, _omitFieldNames ? '' : 'address')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetAddressHistoryRequest clone() => GetAddressHistoryRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetAddressHistoryRequest copyWith(void Function(GetAddressHistoryRequest) updates) => super.copyWith((message) => updates(message as GetAddressHistoryRequest)) as GetAddressHistoryRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetAddressHistoryRequest create() => GetAddressHistoryRequest._();
  GetAddressHistoryRequest createEmptyInstance() => create();
  static $pb.PbList<GetAddressHistoryRequest> createRepeated() => $pb.PbList<GetAddressHistoryRequest>();
  @$core.pragma('dart2js:noInline')
  static GetAddressHistoryRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetAddressHistoryRequest>(create);
  static GetAddressHistoryRequest? _defaultInstance;

  /// The name of the wallet containing the address.
  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => $_clearField(1);

  /// The address to retrieve history for.
  @$pb.TagNumber(2)
  $core.String get address => $_getSZ(1);
  @$pb.TagNumber(2)
  set address($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasAddress() => $_has(1);
  @$pb.TagNumber(2)
  void clearAddress() => $_clearField(2);
}

/// Response message contains address transaction history.
class GetAddressHistoryResponse extends $pb.GeneratedMessage {
  factory GetAddressHistoryResponse({
    $core.Iterable<HistoryInfo>? historyInfo,
  }) {
    final $result = create();
    if (historyInfo != null) {
      $result.historyInfo.addAll(historyInfo);
    }
    return $result;
  }
  GetAddressHistoryResponse._() : super();
  factory GetAddressHistoryResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetAddressHistoryResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetAddressHistoryResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..pc<HistoryInfo>(1, _omitFieldNames ? '' : 'historyInfo', $pb.PbFieldType.PM, subBuilder: HistoryInfo.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetAddressHistoryResponse clone() => GetAddressHistoryResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetAddressHistoryResponse copyWith(void Function(GetAddressHistoryResponse) updates) => super.copyWith((message) => updates(message as GetAddressHistoryResponse)) as GetAddressHistoryResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetAddressHistoryResponse create() => GetAddressHistoryResponse._();
  GetAddressHistoryResponse createEmptyInstance() => create();
  static $pb.PbList<GetAddressHistoryResponse> createRepeated() => $pb.PbList<GetAddressHistoryResponse>();
  @$core.pragma('dart2js:noInline')
  static GetAddressHistoryResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetAddressHistoryResponse>(create);
  static GetAddressHistoryResponse? _defaultInstance;

  /// List of all historical transactions associated with the address.
  @$pb.TagNumber(1)
  $pb.PbList<HistoryInfo> get historyInfo => $_getList(0);
}

/// Request message for generating a new wallet address.
class GetNewAddressRequest extends $pb.GeneratedMessage {
  factory GetNewAddressRequest({
    $core.String? walletName,
    AddressType? addressType,
    $core.String? label,
    $core.String? password,
  }) {
    final $result = create();
    if (walletName != null) {
      $result.walletName = walletName;
    }
    if (addressType != null) {
      $result.addressType = addressType;
    }
    if (label != null) {
      $result.label = label;
    }
    if (password != null) {
      $result.password = password;
    }
    return $result;
  }
  GetNewAddressRequest._() : super();
  factory GetNewAddressRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetNewAddressRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetNewAddressRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'walletName')
    ..e<AddressType>(2, _omitFieldNames ? '' : 'addressType', $pb.PbFieldType.OE, defaultOrMaker: AddressType.ADDRESS_TYPE_TREASURY, valueOf: AddressType.valueOf, enumValues: AddressType.values)
    ..aOS(3, _omitFieldNames ? '' : 'label')
    ..aOS(4, _omitFieldNames ? '' : 'password')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetNewAddressRequest clone() => GetNewAddressRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetNewAddressRequest copyWith(void Function(GetNewAddressRequest) updates) => super.copyWith((message) => updates(message as GetNewAddressRequest)) as GetNewAddressRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetNewAddressRequest create() => GetNewAddressRequest._();
  GetNewAddressRequest createEmptyInstance() => create();
  static $pb.PbList<GetNewAddressRequest> createRepeated() => $pb.PbList<GetNewAddressRequest>();
  @$core.pragma('dart2js:noInline')
  static GetNewAddressRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetNewAddressRequest>(create);
  static GetNewAddressRequest? _defaultInstance;

  /// The name of the wallet to generate a new address.
  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => $_clearField(1);

  /// The type of address to generate.
  @$pb.TagNumber(2)
  AddressType get addressType => $_getN(1);
  @$pb.TagNumber(2)
  set addressType(AddressType v) { $_setField(2, v); }
  @$pb.TagNumber(2)
  $core.bool hasAddressType() => $_has(1);
  @$pb.TagNumber(2)
  void clearAddressType() => $_clearField(2);

  /// A label for the new address.
  @$pb.TagNumber(3)
  $core.String get label => $_getSZ(2);
  @$pb.TagNumber(3)
  set label($core.String v) { $_setString(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasLabel() => $_has(2);
  @$pb.TagNumber(3)
  void clearLabel() => $_clearField(3);

  /// Password for the new address. It's required when address_type is Ed25519 type.
  @$pb.TagNumber(4)
  $core.String get password => $_getSZ(3);
  @$pb.TagNumber(4)
  set password($core.String v) { $_setString(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasPassword() => $_has(3);
  @$pb.TagNumber(4)
  void clearPassword() => $_clearField(4);
}

/// Response message contains newly generated address information.
class GetNewAddressResponse extends $pb.GeneratedMessage {
  factory GetNewAddressResponse({
    $core.String? walletName,
    AddressInfo? addressInfo,
  }) {
    final $result = create();
    if (walletName != null) {
      $result.walletName = walletName;
    }
    if (addressInfo != null) {
      $result.addressInfo = addressInfo;
    }
    return $result;
  }
  GetNewAddressResponse._() : super();
  factory GetNewAddressResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetNewAddressResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetNewAddressResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'walletName')
    ..aOM<AddressInfo>(2, _omitFieldNames ? '' : 'addressInfo', subBuilder: AddressInfo.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetNewAddressResponse clone() => GetNewAddressResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetNewAddressResponse copyWith(void Function(GetNewAddressResponse) updates) => super.copyWith((message) => updates(message as GetNewAddressResponse)) as GetNewAddressResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetNewAddressResponse create() => GetNewAddressResponse._();
  GetNewAddressResponse createEmptyInstance() => create();
  static $pb.PbList<GetNewAddressResponse> createRepeated() => $pb.PbList<GetNewAddressResponse>();
  @$core.pragma('dart2js:noInline')
  static GetNewAddressResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetNewAddressResponse>(create);
  static GetNewAddressResponse? _defaultInstance;

  /// The name of the wallet where address was generated.
  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => $_clearField(1);

  /// Detailed information about the new address.
  @$pb.TagNumber(2)
  AddressInfo get addressInfo => $_getN(1);
  @$pb.TagNumber(2)
  set addressInfo(AddressInfo v) { $_setField(2, v); }
  @$pb.TagNumber(2)
  $core.bool hasAddressInfo() => $_has(1);
  @$pb.TagNumber(2)
  void clearAddressInfo() => $_clearField(2);
  @$pb.TagNumber(2)
  AddressInfo ensureAddressInfo() => $_ensure(1);
}

/// Request message for restoring a wallet from mnemonic (seed phrase).
class RestoreWalletRequest extends $pb.GeneratedMessage {
  factory RestoreWalletRequest({
    $core.String? walletName,
    $core.String? mnemonic,
    $core.String? password,
  }) {
    final $result = create();
    if (walletName != null) {
      $result.walletName = walletName;
    }
    if (mnemonic != null) {
      $result.mnemonic = mnemonic;
    }
    if (password != null) {
      $result.password = password;
    }
    return $result;
  }
  RestoreWalletRequest._() : super();
  factory RestoreWalletRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory RestoreWalletRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'RestoreWalletRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'walletName')
    ..aOS(2, _omitFieldNames ? '' : 'mnemonic')
    ..aOS(3, _omitFieldNames ? '' : 'password')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  RestoreWalletRequest clone() => RestoreWalletRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  RestoreWalletRequest copyWith(void Function(RestoreWalletRequest) updates) => super.copyWith((message) => updates(message as RestoreWalletRequest)) as RestoreWalletRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RestoreWalletRequest create() => RestoreWalletRequest._();
  RestoreWalletRequest createEmptyInstance() => create();
  static $pb.PbList<RestoreWalletRequest> createRepeated() => $pb.PbList<RestoreWalletRequest>();
  @$core.pragma('dart2js:noInline')
  static RestoreWalletRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<RestoreWalletRequest>(create);
  static RestoreWalletRequest? _defaultInstance;

  /// The name for the restored wallet.
  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => $_clearField(1);

  /// The mnemonic (seed phrase) for wallet recovery.
  @$pb.TagNumber(2)
  $core.String get mnemonic => $_getSZ(1);
  @$pb.TagNumber(2)
  set mnemonic($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasMnemonic() => $_has(1);
  @$pb.TagNumber(2)
  void clearMnemonic() => $_clearField(2);

  /// Password to secure the restored wallet.
  @$pb.TagNumber(3)
  $core.String get password => $_getSZ(2);
  @$pb.TagNumber(3)
  set password($core.String v) { $_setString(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasPassword() => $_has(2);
  @$pb.TagNumber(3)
  void clearPassword() => $_clearField(3);
}

/// Response message confirming wallet restoration.
class RestoreWalletResponse extends $pb.GeneratedMessage {
  factory RestoreWalletResponse({
    $core.String? walletName,
  }) {
    final $result = create();
    if (walletName != null) {
      $result.walletName = walletName;
    }
    return $result;
  }
  RestoreWalletResponse._() : super();
  factory RestoreWalletResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory RestoreWalletResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'RestoreWalletResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'walletName')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  RestoreWalletResponse clone() => RestoreWalletResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  RestoreWalletResponse copyWith(void Function(RestoreWalletResponse) updates) => super.copyWith((message) => updates(message as RestoreWalletResponse)) as RestoreWalletResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RestoreWalletResponse create() => RestoreWalletResponse._();
  RestoreWalletResponse createEmptyInstance() => create();
  static $pb.PbList<RestoreWalletResponse> createRepeated() => $pb.PbList<RestoreWalletResponse>();
  @$core.pragma('dart2js:noInline')
  static RestoreWalletResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<RestoreWalletResponse>(create);
  static RestoreWalletResponse? _defaultInstance;

  /// The name of the restored wallet.
  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => $_clearField(1);
}

/// Request message for creating a new wallet.
class CreateWalletRequest extends $pb.GeneratedMessage {
  factory CreateWalletRequest({
    $core.String? walletName,
    $core.String? password,
  }) {
    final $result = create();
    if (walletName != null) {
      $result.walletName = walletName;
    }
    if (password != null) {
      $result.password = password;
    }
    return $result;
  }
  CreateWalletRequest._() : super();
  factory CreateWalletRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory CreateWalletRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'CreateWalletRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'walletName')
    ..aOS(4, _omitFieldNames ? '' : 'password')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  CreateWalletRequest clone() => CreateWalletRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  CreateWalletRequest copyWith(void Function(CreateWalletRequest) updates) => super.copyWith((message) => updates(message as CreateWalletRequest)) as CreateWalletRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateWalletRequest create() => CreateWalletRequest._();
  CreateWalletRequest createEmptyInstance() => create();
  static $pb.PbList<CreateWalletRequest> createRepeated() => $pb.PbList<CreateWalletRequest>();
  @$core.pragma('dart2js:noInline')
  static CreateWalletRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<CreateWalletRequest>(create);
  static CreateWalletRequest? _defaultInstance;

  /// The name for the new wallet.
  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => $_clearField(1);

  /// Password to secure the new wallet.
  @$pb.TagNumber(4)
  $core.String get password => $_getSZ(1);
  @$pb.TagNumber(4)
  set password($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(4)
  $core.bool hasPassword() => $_has(1);
  @$pb.TagNumber(4)
  void clearPassword() => $_clearField(4);
}

/// Response message contains wallet recovery mnemonic (seed phrase).
class CreateWalletResponse extends $pb.GeneratedMessage {
  factory CreateWalletResponse({
    $core.String? mnemonic,
  }) {
    final $result = create();
    if (mnemonic != null) {
      $result.mnemonic = mnemonic;
    }
    return $result;
  }
  CreateWalletResponse._() : super();
  factory CreateWalletResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory CreateWalletResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'CreateWalletResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(2, _omitFieldNames ? '' : 'mnemonic')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  CreateWalletResponse clone() => CreateWalletResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  CreateWalletResponse copyWith(void Function(CreateWalletResponse) updates) => super.copyWith((message) => updates(message as CreateWalletResponse)) as CreateWalletResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateWalletResponse create() => CreateWalletResponse._();
  CreateWalletResponse createEmptyInstance() => create();
  static $pb.PbList<CreateWalletResponse> createRepeated() => $pb.PbList<CreateWalletResponse>();
  @$core.pragma('dart2js:noInline')
  static CreateWalletResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<CreateWalletResponse>(create);
  static CreateWalletResponse? _defaultInstance;

  /// The mnemonic (seed phrase) for wallet recovery.
  @$pb.TagNumber(2)
  $core.String get mnemonic => $_getSZ(0);
  @$pb.TagNumber(2)
  set mnemonic($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(2)
  $core.bool hasMnemonic() => $_has(0);
  @$pb.TagNumber(2)
  void clearMnemonic() => $_clearField(2);
}

/// Request message for loading an existing wallet.
class LoadWalletRequest extends $pb.GeneratedMessage {
  factory LoadWalletRequest({
    $core.String? walletName,
  }) {
    final $result = create();
    if (walletName != null) {
      $result.walletName = walletName;
    }
    return $result;
  }
  LoadWalletRequest._() : super();
  factory LoadWalletRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory LoadWalletRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'LoadWalletRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'walletName')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  LoadWalletRequest clone() => LoadWalletRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  LoadWalletRequest copyWith(void Function(LoadWalletRequest) updates) => super.copyWith((message) => updates(message as LoadWalletRequest)) as LoadWalletRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static LoadWalletRequest create() => LoadWalletRequest._();
  LoadWalletRequest createEmptyInstance() => create();
  static $pb.PbList<LoadWalletRequest> createRepeated() => $pb.PbList<LoadWalletRequest>();
  @$core.pragma('dart2js:noInline')
  static LoadWalletRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<LoadWalletRequest>(create);
  static LoadWalletRequest? _defaultInstance;

  /// The name of the wallet to load.
  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => $_clearField(1);
}

/// Response message confirming wallet loaded.
class LoadWalletResponse extends $pb.GeneratedMessage {
  factory LoadWalletResponse({
    $core.String? walletName,
  }) {
    final $result = create();
    if (walletName != null) {
      $result.walletName = walletName;
    }
    return $result;
  }
  LoadWalletResponse._() : super();
  factory LoadWalletResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory LoadWalletResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'LoadWalletResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'walletName')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  LoadWalletResponse clone() => LoadWalletResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  LoadWalletResponse copyWith(void Function(LoadWalletResponse) updates) => super.copyWith((message) => updates(message as LoadWalletResponse)) as LoadWalletResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static LoadWalletResponse create() => LoadWalletResponse._();
  LoadWalletResponse createEmptyInstance() => create();
  static $pb.PbList<LoadWalletResponse> createRepeated() => $pb.PbList<LoadWalletResponse>();
  @$core.pragma('dart2js:noInline')
  static LoadWalletResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<LoadWalletResponse>(create);
  static LoadWalletResponse? _defaultInstance;

  /// The name of the loaded wallet.
  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => $_clearField(1);
}

/// Request message for unloading a wallet.
class UnloadWalletRequest extends $pb.GeneratedMessage {
  factory UnloadWalletRequest({
    $core.String? walletName,
  }) {
    final $result = create();
    if (walletName != null) {
      $result.walletName = walletName;
    }
    return $result;
  }
  UnloadWalletRequest._() : super();
  factory UnloadWalletRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory UnloadWalletRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'UnloadWalletRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'walletName')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  UnloadWalletRequest clone() => UnloadWalletRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  UnloadWalletRequest copyWith(void Function(UnloadWalletRequest) updates) => super.copyWith((message) => updates(message as UnloadWalletRequest)) as UnloadWalletRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UnloadWalletRequest create() => UnloadWalletRequest._();
  UnloadWalletRequest createEmptyInstance() => create();
  static $pb.PbList<UnloadWalletRequest> createRepeated() => $pb.PbList<UnloadWalletRequest>();
  @$core.pragma('dart2js:noInline')
  static UnloadWalletRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<UnloadWalletRequest>(create);
  static UnloadWalletRequest? _defaultInstance;

  /// The name of the wallet to unload.
  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => $_clearField(1);
}

/// Response message confirming wallet unloading.
class UnloadWalletResponse extends $pb.GeneratedMessage {
  factory UnloadWalletResponse({
    $core.String? walletName,
  }) {
    final $result = create();
    if (walletName != null) {
      $result.walletName = walletName;
    }
    return $result;
  }
  UnloadWalletResponse._() : super();
  factory UnloadWalletResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory UnloadWalletResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'UnloadWalletResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'walletName')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  UnloadWalletResponse clone() => UnloadWalletResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  UnloadWalletResponse copyWith(void Function(UnloadWalletResponse) updates) => super.copyWith((message) => updates(message as UnloadWalletResponse)) as UnloadWalletResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UnloadWalletResponse create() => UnloadWalletResponse._();
  UnloadWalletResponse createEmptyInstance() => create();
  static $pb.PbList<UnloadWalletResponse> createRepeated() => $pb.PbList<UnloadWalletResponse>();
  @$core.pragma('dart2js:noInline')
  static UnloadWalletResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<UnloadWalletResponse>(create);
  static UnloadWalletResponse? _defaultInstance;

  /// The name of the unloaded wallet.
  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => $_clearField(1);
}

/// Request message for obtaining the validator address associated with a public key.
class GetValidatorAddressRequest extends $pb.GeneratedMessage {
  factory GetValidatorAddressRequest({
    $core.String? publicKey,
  }) {
    final $result = create();
    if (publicKey != null) {
      $result.publicKey = publicKey;
    }
    return $result;
  }
  GetValidatorAddressRequest._() : super();
  factory GetValidatorAddressRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetValidatorAddressRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetValidatorAddressRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'publicKey')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetValidatorAddressRequest clone() => GetValidatorAddressRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetValidatorAddressRequest copyWith(void Function(GetValidatorAddressRequest) updates) => super.copyWith((message) => updates(message as GetValidatorAddressRequest)) as GetValidatorAddressRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetValidatorAddressRequest create() => GetValidatorAddressRequest._();
  GetValidatorAddressRequest createEmptyInstance() => create();
  static $pb.PbList<GetValidatorAddressRequest> createRepeated() => $pb.PbList<GetValidatorAddressRequest>();
  @$core.pragma('dart2js:noInline')
  static GetValidatorAddressRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetValidatorAddressRequest>(create);
  static GetValidatorAddressRequest? _defaultInstance;

  /// The public key of the validator.
  @$pb.TagNumber(1)
  $core.String get publicKey => $_getSZ(0);
  @$pb.TagNumber(1)
  set publicKey($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasPublicKey() => $_has(0);
  @$pb.TagNumber(1)
  void clearPublicKey() => $_clearField(1);
}

/// Response message containing the validator address corresponding to a public key.
class GetValidatorAddressResponse extends $pb.GeneratedMessage {
  factory GetValidatorAddressResponse({
    $core.String? address,
  }) {
    final $result = create();
    if (address != null) {
      $result.address = address;
    }
    return $result;
  }
  GetValidatorAddressResponse._() : super();
  factory GetValidatorAddressResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetValidatorAddressResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetValidatorAddressResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'address')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetValidatorAddressResponse clone() => GetValidatorAddressResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetValidatorAddressResponse copyWith(void Function(GetValidatorAddressResponse) updates) => super.copyWith((message) => updates(message as GetValidatorAddressResponse)) as GetValidatorAddressResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetValidatorAddressResponse create() => GetValidatorAddressResponse._();
  GetValidatorAddressResponse createEmptyInstance() => create();
  static $pb.PbList<GetValidatorAddressResponse> createRepeated() => $pb.PbList<GetValidatorAddressResponse>();
  @$core.pragma('dart2js:noInline')
  static GetValidatorAddressResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetValidatorAddressResponse>(create);
  static GetValidatorAddressResponse? _defaultInstance;

  /// The validator address associated with the public key.
  @$pb.TagNumber(1)
  $core.String get address => $_getSZ(0);
  @$pb.TagNumber(1)
  set address($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasAddress() => $_has(0);
  @$pb.TagNumber(1)
  void clearAddress() => $_clearField(1);
}

/// Request message for signing a raw transaction.
class SignRawTransactionRequest extends $pb.GeneratedMessage {
  factory SignRawTransactionRequest({
    $core.String? walletName,
    $core.String? rawTransaction,
    $core.String? password,
  }) {
    final $result = create();
    if (walletName != null) {
      $result.walletName = walletName;
    }
    if (rawTransaction != null) {
      $result.rawTransaction = rawTransaction;
    }
    if (password != null) {
      $result.password = password;
    }
    return $result;
  }
  SignRawTransactionRequest._() : super();
  factory SignRawTransactionRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory SignRawTransactionRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'SignRawTransactionRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'walletName')
    ..aOS(2, _omitFieldNames ? '' : 'rawTransaction')
    ..aOS(3, _omitFieldNames ? '' : 'password')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  SignRawTransactionRequest clone() => SignRawTransactionRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  SignRawTransactionRequest copyWith(void Function(SignRawTransactionRequest) updates) => super.copyWith((message) => updates(message as SignRawTransactionRequest)) as SignRawTransactionRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SignRawTransactionRequest create() => SignRawTransactionRequest._();
  SignRawTransactionRequest createEmptyInstance() => create();
  static $pb.PbList<SignRawTransactionRequest> createRepeated() => $pb.PbList<SignRawTransactionRequest>();
  @$core.pragma('dart2js:noInline')
  static SignRawTransactionRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<SignRawTransactionRequest>(create);
  static SignRawTransactionRequest? _defaultInstance;

  /// The name of the wallet used for signing.
  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => $_clearField(1);

  /// The raw transaction data to be signed.
  @$pb.TagNumber(2)
  $core.String get rawTransaction => $_getSZ(1);
  @$pb.TagNumber(2)
  set rawTransaction($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasRawTransaction() => $_has(1);
  @$pb.TagNumber(2)
  void clearRawTransaction() => $_clearField(2);

  /// Wallet password required for signing.
  @$pb.TagNumber(3)
  $core.String get password => $_getSZ(2);
  @$pb.TagNumber(3)
  set password($core.String v) { $_setString(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasPassword() => $_has(2);
  @$pb.TagNumber(3)
  void clearPassword() => $_clearField(3);
}

/// Response message contains the transaction ID and signed raw transaction.
class SignRawTransactionResponse extends $pb.GeneratedMessage {
  factory SignRawTransactionResponse({
    $core.String? transactionId,
    $core.String? signedRawTransaction,
  }) {
    final $result = create();
    if (transactionId != null) {
      $result.transactionId = transactionId;
    }
    if (signedRawTransaction != null) {
      $result.signedRawTransaction = signedRawTransaction;
    }
    return $result;
  }
  SignRawTransactionResponse._() : super();
  factory SignRawTransactionResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory SignRawTransactionResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'SignRawTransactionResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'transactionId')
    ..aOS(2, _omitFieldNames ? '' : 'signedRawTransaction')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  SignRawTransactionResponse clone() => SignRawTransactionResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  SignRawTransactionResponse copyWith(void Function(SignRawTransactionResponse) updates) => super.copyWith((message) => updates(message as SignRawTransactionResponse)) as SignRawTransactionResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SignRawTransactionResponse create() => SignRawTransactionResponse._();
  SignRawTransactionResponse createEmptyInstance() => create();
  static $pb.PbList<SignRawTransactionResponse> createRepeated() => $pb.PbList<SignRawTransactionResponse>();
  @$core.pragma('dart2js:noInline')
  static SignRawTransactionResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<SignRawTransactionResponse>(create);
  static SignRawTransactionResponse? _defaultInstance;

  /// The ID of the signed transaction.
  @$pb.TagNumber(1)
  $core.String get transactionId => $_getSZ(0);
  @$pb.TagNumber(1)
  set transactionId($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasTransactionId() => $_has(0);
  @$pb.TagNumber(1)
  void clearTransactionId() => $_clearField(1);

  /// The signed raw transaction data.
  @$pb.TagNumber(2)
  $core.String get signedRawTransaction => $_getSZ(1);
  @$pb.TagNumber(2)
  set signedRawTransaction($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasSignedRawTransaction() => $_has(1);
  @$pb.TagNumber(2)
  void clearSignedRawTransaction() => $_clearField(2);
}

/// Request message for obtaining the total available balance of a wallet.
class GetTotalBalanceRequest extends $pb.GeneratedMessage {
  factory GetTotalBalanceRequest({
    $core.String? walletName,
  }) {
    final $result = create();
    if (walletName != null) {
      $result.walletName = walletName;
    }
    return $result;
  }
  GetTotalBalanceRequest._() : super();
  factory GetTotalBalanceRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetTotalBalanceRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetTotalBalanceRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'walletName')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetTotalBalanceRequest clone() => GetTotalBalanceRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetTotalBalanceRequest copyWith(void Function(GetTotalBalanceRequest) updates) => super.copyWith((message) => updates(message as GetTotalBalanceRequest)) as GetTotalBalanceRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetTotalBalanceRequest create() => GetTotalBalanceRequest._();
  GetTotalBalanceRequest createEmptyInstance() => create();
  static $pb.PbList<GetTotalBalanceRequest> createRepeated() => $pb.PbList<GetTotalBalanceRequest>();
  @$core.pragma('dart2js:noInline')
  static GetTotalBalanceRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetTotalBalanceRequest>(create);
  static GetTotalBalanceRequest? _defaultInstance;

  /// The name of the wallet to get the total balance.
  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => $_clearField(1);
}

/// Response message contains the total available balance of the wallet.
class GetTotalBalanceResponse extends $pb.GeneratedMessage {
  factory GetTotalBalanceResponse({
    $core.String? walletName,
    $fixnum.Int64? totalBalance,
  }) {
    final $result = create();
    if (walletName != null) {
      $result.walletName = walletName;
    }
    if (totalBalance != null) {
      $result.totalBalance = totalBalance;
    }
    return $result;
  }
  GetTotalBalanceResponse._() : super();
  factory GetTotalBalanceResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetTotalBalanceResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetTotalBalanceResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'walletName')
    ..aInt64(2, _omitFieldNames ? '' : 'totalBalance')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetTotalBalanceResponse clone() => GetTotalBalanceResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetTotalBalanceResponse copyWith(void Function(GetTotalBalanceResponse) updates) => super.copyWith((message) => updates(message as GetTotalBalanceResponse)) as GetTotalBalanceResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetTotalBalanceResponse create() => GetTotalBalanceResponse._();
  GetTotalBalanceResponse createEmptyInstance() => create();
  static $pb.PbList<GetTotalBalanceResponse> createRepeated() => $pb.PbList<GetTotalBalanceResponse>();
  @$core.pragma('dart2js:noInline')
  static GetTotalBalanceResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetTotalBalanceResponse>(create);
  static GetTotalBalanceResponse? _defaultInstance;

  /// The name of the queried wallet.
  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => $_clearField(1);

  /// The total balance of the wallet in NanoPAC.
  @$pb.TagNumber(2)
  $fixnum.Int64 get totalBalance => $_getI64(1);
  @$pb.TagNumber(2)
  set totalBalance($fixnum.Int64 v) { $_setInt64(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasTotalBalance() => $_has(1);
  @$pb.TagNumber(2)
  void clearTotalBalance() => $_clearField(2);
}

/// Request message to sign an arbitrary message.
class SignMessageRequest extends $pb.GeneratedMessage {
  factory SignMessageRequest({
    $core.String? walletName,
    $core.String? password,
    $core.String? address,
    $core.String? message,
  }) {
    final $result = create();
    if (walletName != null) {
      $result.walletName = walletName;
    }
    if (password != null) {
      $result.password = password;
    }
    if (address != null) {
      $result.address = address;
    }
    if (message != null) {
      $result.message = message;
    }
    return $result;
  }
  SignMessageRequest._() : super();
  factory SignMessageRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory SignMessageRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'SignMessageRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'walletName')
    ..aOS(2, _omitFieldNames ? '' : 'password')
    ..aOS(3, _omitFieldNames ? '' : 'address')
    ..aOS(4, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  SignMessageRequest clone() => SignMessageRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  SignMessageRequest copyWith(void Function(SignMessageRequest) updates) => super.copyWith((message) => updates(message as SignMessageRequest)) as SignMessageRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SignMessageRequest create() => SignMessageRequest._();
  SignMessageRequest createEmptyInstance() => create();
  static $pb.PbList<SignMessageRequest> createRepeated() => $pb.PbList<SignMessageRequest>();
  @$core.pragma('dart2js:noInline')
  static SignMessageRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<SignMessageRequest>(create);
  static SignMessageRequest? _defaultInstance;

  /// The name of the wallet to sign with.
  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => $_clearField(1);

  /// Wallet password required for signing.
  @$pb.TagNumber(2)
  $core.String get password => $_getSZ(1);
  @$pb.TagNumber(2)
  set password($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasPassword() => $_has(1);
  @$pb.TagNumber(2)
  void clearPassword() => $_clearField(2);

  /// The address whose private key should be used for signing the message.
  @$pb.TagNumber(3)
  $core.String get address => $_getSZ(2);
  @$pb.TagNumber(3)
  set address($core.String v) { $_setString(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasAddress() => $_has(2);
  @$pb.TagNumber(3)
  void clearAddress() => $_clearField(3);

  /// The arbitrary message to be signed.
  @$pb.TagNumber(4)
  $core.String get message => $_getSZ(3);
  @$pb.TagNumber(4)
  set message($core.String v) { $_setString(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasMessage() => $_has(3);
  @$pb.TagNumber(4)
  void clearMessage() => $_clearField(4);
}

/// Response message contains message signature.
class SignMessageResponse extends $pb.GeneratedMessage {
  factory SignMessageResponse({
    $core.String? signature,
  }) {
    final $result = create();
    if (signature != null) {
      $result.signature = signature;
    }
    return $result;
  }
  SignMessageResponse._() : super();
  factory SignMessageResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory SignMessageResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'SignMessageResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'signature')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  SignMessageResponse clone() => SignMessageResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  SignMessageResponse copyWith(void Function(SignMessageResponse) updates) => super.copyWith((message) => updates(message as SignMessageResponse)) as SignMessageResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SignMessageResponse create() => SignMessageResponse._();
  SignMessageResponse createEmptyInstance() => create();
  static $pb.PbList<SignMessageResponse> createRepeated() => $pb.PbList<SignMessageResponse>();
  @$core.pragma('dart2js:noInline')
  static SignMessageResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<SignMessageResponse>(create);
  static SignMessageResponse? _defaultInstance;

  /// The signature in hexadecimal format.
  @$pb.TagNumber(1)
  $core.String get signature => $_getSZ(0);
  @$pb.TagNumber(1)
  set signature($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasSignature() => $_has(0);
  @$pb.TagNumber(1)
  void clearSignature() => $_clearField(1);
}

/// Request message for obtaining the total stake of a wallet.
class GetTotalStakeRequest extends $pb.GeneratedMessage {
  factory GetTotalStakeRequest({
    $core.String? walletName,
  }) {
    final $result = create();
    if (walletName != null) {
      $result.walletName = walletName;
    }
    return $result;
  }
  GetTotalStakeRequest._() : super();
  factory GetTotalStakeRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetTotalStakeRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetTotalStakeRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'walletName')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetTotalStakeRequest clone() => GetTotalStakeRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetTotalStakeRequest copyWith(void Function(GetTotalStakeRequest) updates) => super.copyWith((message) => updates(message as GetTotalStakeRequest)) as GetTotalStakeRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetTotalStakeRequest create() => GetTotalStakeRequest._();
  GetTotalStakeRequest createEmptyInstance() => create();
  static $pb.PbList<GetTotalStakeRequest> createRepeated() => $pb.PbList<GetTotalStakeRequest>();
  @$core.pragma('dart2js:noInline')
  static GetTotalStakeRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetTotalStakeRequest>(create);
  static GetTotalStakeRequest? _defaultInstance;

  /// The name of the wallet to get the total stake.
  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => $_clearField(1);
}

/// Response message contains the total stake of the wallet.
class GetTotalStakeResponse extends $pb.GeneratedMessage {
  factory GetTotalStakeResponse({
    $core.String? walletName,
    $fixnum.Int64? totalStake,
  }) {
    final $result = create();
    if (walletName != null) {
      $result.walletName = walletName;
    }
    if (totalStake != null) {
      $result.totalStake = totalStake;
    }
    return $result;
  }
  GetTotalStakeResponse._() : super();
  factory GetTotalStakeResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetTotalStakeResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetTotalStakeResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'walletName')
    ..aInt64(2, _omitFieldNames ? '' : 'totalStake')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetTotalStakeResponse clone() => GetTotalStakeResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetTotalStakeResponse copyWith(void Function(GetTotalStakeResponse) updates) => super.copyWith((message) => updates(message as GetTotalStakeResponse)) as GetTotalStakeResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetTotalStakeResponse create() => GetTotalStakeResponse._();
  GetTotalStakeResponse createEmptyInstance() => create();
  static $pb.PbList<GetTotalStakeResponse> createRepeated() => $pb.PbList<GetTotalStakeResponse>();
  @$core.pragma('dart2js:noInline')
  static GetTotalStakeResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetTotalStakeResponse>(create);
  static GetTotalStakeResponse? _defaultInstance;

  /// The name of the queried wallet.
  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => $_clearField(1);

  /// The total stake amount in NanoPAC.
  @$pb.TagNumber(2)
  $fixnum.Int64 get totalStake => $_getI64(1);
  @$pb.TagNumber(2)
  set totalStake($fixnum.Int64 v) { $_setInt64(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasTotalStake() => $_has(1);
  @$pb.TagNumber(2)
  void clearTotalStake() => $_clearField(2);
}

/// Request message for getting address information.
class GetAddressInfoRequest extends $pb.GeneratedMessage {
  factory GetAddressInfoRequest({
    $core.String? walletName,
    $core.String? address,
  }) {
    final $result = create();
    if (walletName != null) {
      $result.walletName = walletName;
    }
    if (address != null) {
      $result.address = address;
    }
    return $result;
  }
  GetAddressInfoRequest._() : super();
  factory GetAddressInfoRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetAddressInfoRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetAddressInfoRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'walletName')
    ..aOS(2, _omitFieldNames ? '' : 'address')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetAddressInfoRequest clone() => GetAddressInfoRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetAddressInfoRequest copyWith(void Function(GetAddressInfoRequest) updates) => super.copyWith((message) => updates(message as GetAddressInfoRequest)) as GetAddressInfoRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetAddressInfoRequest create() => GetAddressInfoRequest._();
  GetAddressInfoRequest createEmptyInstance() => create();
  static $pb.PbList<GetAddressInfoRequest> createRepeated() => $pb.PbList<GetAddressInfoRequest>();
  @$core.pragma('dart2js:noInline')
  static GetAddressInfoRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetAddressInfoRequest>(create);
  static GetAddressInfoRequest? _defaultInstance;

  /// The name of the wallet containing the address.
  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => $_clearField(1);

  /// The address to query.
  @$pb.TagNumber(2)
  $core.String get address => $_getSZ(1);
  @$pb.TagNumber(2)
  set address($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasAddress() => $_has(1);
  @$pb.TagNumber(2)
  void clearAddress() => $_clearField(2);
}

/// Response message contains address details.
class GetAddressInfoResponse extends $pb.GeneratedMessage {
  factory GetAddressInfoResponse({
    $core.String? walletName,
    $core.String? address,
    $core.String? label,
    $core.String? publicKey,
    $core.String? path,
  }) {
    final $result = create();
    if (walletName != null) {
      $result.walletName = walletName;
    }
    if (address != null) {
      $result.address = address;
    }
    if (label != null) {
      $result.label = label;
    }
    if (publicKey != null) {
      $result.publicKey = publicKey;
    }
    if (path != null) {
      $result.path = path;
    }
    return $result;
  }
  GetAddressInfoResponse._() : super();
  factory GetAddressInfoResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetAddressInfoResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetAddressInfoResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'walletName')
    ..aOS(2, _omitFieldNames ? '' : 'address')
    ..aOS(3, _omitFieldNames ? '' : 'label')
    ..aOS(4, _omitFieldNames ? '' : 'publicKey')
    ..aOS(5, _omitFieldNames ? '' : 'path')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetAddressInfoResponse clone() => GetAddressInfoResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetAddressInfoResponse copyWith(void Function(GetAddressInfoResponse) updates) => super.copyWith((message) => updates(message as GetAddressInfoResponse)) as GetAddressInfoResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetAddressInfoResponse create() => GetAddressInfoResponse._();
  GetAddressInfoResponse createEmptyInstance() => create();
  static $pb.PbList<GetAddressInfoResponse> createRepeated() => $pb.PbList<GetAddressInfoResponse>();
  @$core.pragma('dart2js:noInline')
  static GetAddressInfoResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetAddressInfoResponse>(create);
  static GetAddressInfoResponse? _defaultInstance;

  /// The name of the wallet containing the address.
  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => $_clearField(1);

  /// The queried address.
  @$pb.TagNumber(2)
  $core.String get address => $_getSZ(1);
  @$pb.TagNumber(2)
  set address($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasAddress() => $_has(1);
  @$pb.TagNumber(2)
  void clearAddress() => $_clearField(2);

  /// The address label.
  @$pb.TagNumber(3)
  $core.String get label => $_getSZ(2);
  @$pb.TagNumber(3)
  set label($core.String v) { $_setString(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasLabel() => $_has(2);
  @$pb.TagNumber(3)
  void clearLabel() => $_clearField(3);

  /// The public key of the address.
  @$pb.TagNumber(4)
  $core.String get publicKey => $_getSZ(3);
  @$pb.TagNumber(4)
  set publicKey($core.String v) { $_setString(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasPublicKey() => $_has(3);
  @$pb.TagNumber(4)
  void clearPublicKey() => $_clearField(4);

  /// The Hierarchical Deterministic (HD) path of the address.
  @$pb.TagNumber(5)
  $core.String get path => $_getSZ(4);
  @$pb.TagNumber(5)
  set path($core.String v) { $_setString(4, v); }
  @$pb.TagNumber(5)
  $core.bool hasPath() => $_has(4);
  @$pb.TagNumber(5)
  void clearPath() => $_clearField(5);
}

/// Request message for setting address label.
class SetAddressLabelRequest extends $pb.GeneratedMessage {
  factory SetAddressLabelRequest({
    $core.String? walletName,
    $core.String? password,
    $core.String? address,
    $core.String? label,
  }) {
    final $result = create();
    if (walletName != null) {
      $result.walletName = walletName;
    }
    if (password != null) {
      $result.password = password;
    }
    if (address != null) {
      $result.address = address;
    }
    if (label != null) {
      $result.label = label;
    }
    return $result;
  }
  SetAddressLabelRequest._() : super();
  factory SetAddressLabelRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory SetAddressLabelRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'SetAddressLabelRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'walletName')
    ..aOS(3, _omitFieldNames ? '' : 'password')
    ..aOS(4, _omitFieldNames ? '' : 'address')
    ..aOS(5, _omitFieldNames ? '' : 'label')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  SetAddressLabelRequest clone() => SetAddressLabelRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  SetAddressLabelRequest copyWith(void Function(SetAddressLabelRequest) updates) => super.copyWith((message) => updates(message as SetAddressLabelRequest)) as SetAddressLabelRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SetAddressLabelRequest create() => SetAddressLabelRequest._();
  SetAddressLabelRequest createEmptyInstance() => create();
  static $pb.PbList<SetAddressLabelRequest> createRepeated() => $pb.PbList<SetAddressLabelRequest>();
  @$core.pragma('dart2js:noInline')
  static SetAddressLabelRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<SetAddressLabelRequest>(create);
  static SetAddressLabelRequest? _defaultInstance;

  /// The name of the wallet containing the address.
  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => $_clearField(1);

  /// Wallet password required for modification.
  @$pb.TagNumber(3)
  $core.String get password => $_getSZ(1);
  @$pb.TagNumber(3)
  set password($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(3)
  $core.bool hasPassword() => $_has(1);
  @$pb.TagNumber(3)
  void clearPassword() => $_clearField(3);

  /// The address to label.
  @$pb.TagNumber(4)
  $core.String get address => $_getSZ(2);
  @$pb.TagNumber(4)
  set address($core.String v) { $_setString(2, v); }
  @$pb.TagNumber(4)
  $core.bool hasAddress() => $_has(2);
  @$pb.TagNumber(4)
  void clearAddress() => $_clearField(4);

  /// The new label for the address.
  @$pb.TagNumber(5)
  $core.String get label => $_getSZ(3);
  @$pb.TagNumber(5)
  set label($core.String v) { $_setString(3, v); }
  @$pb.TagNumber(5)
  $core.bool hasLabel() => $_has(3);
  @$pb.TagNumber(5)
  void clearLabel() => $_clearField(5);
}

/// Response message for address label update.
class SetAddressLabelResponse extends $pb.GeneratedMessage {
  factory SetAddressLabelResponse() => create();
  SetAddressLabelResponse._() : super();
  factory SetAddressLabelResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory SetAddressLabelResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'SetAddressLabelResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  SetAddressLabelResponse clone() => SetAddressLabelResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  SetAddressLabelResponse copyWith(void Function(SetAddressLabelResponse) updates) => super.copyWith((message) => updates(message as SetAddressLabelResponse)) as SetAddressLabelResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SetAddressLabelResponse create() => SetAddressLabelResponse._();
  SetAddressLabelResponse createEmptyInstance() => create();
  static $pb.PbList<SetAddressLabelResponse> createRepeated() => $pb.PbList<SetAddressLabelResponse>();
  @$core.pragma('dart2js:noInline')
  static SetAddressLabelResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<SetAddressLabelResponse>(create);
  static SetAddressLabelResponse? _defaultInstance;
}

/// Request message for listing all wallets.
class ListWalletRequest extends $pb.GeneratedMessage {
  factory ListWalletRequest() => create();
  ListWalletRequest._() : super();
  factory ListWalletRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory ListWalletRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'ListWalletRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  ListWalletRequest clone() => ListWalletRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  ListWalletRequest copyWith(void Function(ListWalletRequest) updates) => super.copyWith((message) => updates(message as ListWalletRequest)) as ListWalletRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListWalletRequest create() => ListWalletRequest._();
  ListWalletRequest createEmptyInstance() => create();
  static $pb.PbList<ListWalletRequest> createRepeated() => $pb.PbList<ListWalletRequest>();
  @$core.pragma('dart2js:noInline')
  static ListWalletRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ListWalletRequest>(create);
  static ListWalletRequest? _defaultInstance;
}

/// Response message contains wallet names.
class ListWalletResponse extends $pb.GeneratedMessage {
  factory ListWalletResponse({
    $core.Iterable<$core.String>? wallets,
  }) {
    final $result = create();
    if (wallets != null) {
      $result.wallets.addAll(wallets);
    }
    return $result;
  }
  ListWalletResponse._() : super();
  factory ListWalletResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory ListWalletResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'ListWalletResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..pPS(1, _omitFieldNames ? '' : 'wallets')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  ListWalletResponse clone() => ListWalletResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  ListWalletResponse copyWith(void Function(ListWalletResponse) updates) => super.copyWith((message) => updates(message as ListWalletResponse)) as ListWalletResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListWalletResponse create() => ListWalletResponse._();
  ListWalletResponse createEmptyInstance() => create();
  static $pb.PbList<ListWalletResponse> createRepeated() => $pb.PbList<ListWalletResponse>();
  @$core.pragma('dart2js:noInline')
  static ListWalletResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ListWalletResponse>(create);
  static ListWalletResponse? _defaultInstance;

  /// Array of wallet names.
  @$pb.TagNumber(1)
  $pb.PbList<$core.String> get wallets => $_getList(0);
}

/// Request message for getting wallet information.
class GetWalletInfoRequest extends $pb.GeneratedMessage {
  factory GetWalletInfoRequest({
    $core.String? walletName,
  }) {
    final $result = create();
    if (walletName != null) {
      $result.walletName = walletName;
    }
    return $result;
  }
  GetWalletInfoRequest._() : super();
  factory GetWalletInfoRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetWalletInfoRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetWalletInfoRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'walletName')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetWalletInfoRequest clone() => GetWalletInfoRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetWalletInfoRequest copyWith(void Function(GetWalletInfoRequest) updates) => super.copyWith((message) => updates(message as GetWalletInfoRequest)) as GetWalletInfoRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetWalletInfoRequest create() => GetWalletInfoRequest._();
  GetWalletInfoRequest createEmptyInstance() => create();
  static $pb.PbList<GetWalletInfoRequest> createRepeated() => $pb.PbList<GetWalletInfoRequest>();
  @$core.pragma('dart2js:noInline')
  static GetWalletInfoRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetWalletInfoRequest>(create);
  static GetWalletInfoRequest? _defaultInstance;

  /// The name of the wallet to query.
  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => $_clearField(1);
}

/// Response message contains wallet details.
class GetWalletInfoResponse extends $pb.GeneratedMessage {
  factory GetWalletInfoResponse({
    $core.String? walletName,
    $core.int? version,
    $core.String? network,
    $core.bool? encrypted,
    $core.String? uuid,
    $fixnum.Int64? createdAt,
    $fixnum.Int64? defaultFee,
  }) {
    final $result = create();
    if (walletName != null) {
      $result.walletName = walletName;
    }
    if (version != null) {
      $result.version = version;
    }
    if (network != null) {
      $result.network = network;
    }
    if (encrypted != null) {
      $result.encrypted = encrypted;
    }
    if (uuid != null) {
      $result.uuid = uuid;
    }
    if (createdAt != null) {
      $result.createdAt = createdAt;
    }
    if (defaultFee != null) {
      $result.defaultFee = defaultFee;
    }
    return $result;
  }
  GetWalletInfoResponse._() : super();
  factory GetWalletInfoResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetWalletInfoResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetWalletInfoResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'walletName')
    ..a<$core.int>(2, _omitFieldNames ? '' : 'version', $pb.PbFieldType.O3)
    ..aOS(3, _omitFieldNames ? '' : 'network')
    ..aOB(4, _omitFieldNames ? '' : 'encrypted')
    ..aOS(5, _omitFieldNames ? '' : 'uuid')
    ..aInt64(6, _omitFieldNames ? '' : 'createdAt')
    ..aInt64(7, _omitFieldNames ? '' : 'defaultFee')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetWalletInfoResponse clone() => GetWalletInfoResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetWalletInfoResponse copyWith(void Function(GetWalletInfoResponse) updates) => super.copyWith((message) => updates(message as GetWalletInfoResponse)) as GetWalletInfoResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetWalletInfoResponse create() => GetWalletInfoResponse._();
  GetWalletInfoResponse createEmptyInstance() => create();
  static $pb.PbList<GetWalletInfoResponse> createRepeated() => $pb.PbList<GetWalletInfoResponse>();
  @$core.pragma('dart2js:noInline')
  static GetWalletInfoResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetWalletInfoResponse>(create);
  static GetWalletInfoResponse? _defaultInstance;

  /// The name of the wallet to query.
  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => $_clearField(1);

  /// The wallet format version.
  @$pb.TagNumber(2)
  $core.int get version => $_getIZ(1);
  @$pb.TagNumber(2)
  set version($core.int v) { $_setSignedInt32(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasVersion() => $_has(1);
  @$pb.TagNumber(2)
  void clearVersion() => $_clearField(2);

  /// The network the wallet is connected to (e.g., mainnet, testnet).
  @$pb.TagNumber(3)
  $core.String get network => $_getSZ(2);
  @$pb.TagNumber(3)
  set network($core.String v) { $_setString(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasNetwork() => $_has(2);
  @$pb.TagNumber(3)
  void clearNetwork() => $_clearField(3);

  /// Indicates if the wallet is encrypted.
  @$pb.TagNumber(4)
  $core.bool get encrypted => $_getBF(3);
  @$pb.TagNumber(4)
  set encrypted($core.bool v) { $_setBool(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasEncrypted() => $_has(3);
  @$pb.TagNumber(4)
  void clearEncrypted() => $_clearField(4);

  /// A unique identifier of the wallet.
  @$pb.TagNumber(5)
  $core.String get uuid => $_getSZ(4);
  @$pb.TagNumber(5)
  set uuid($core.String v) { $_setString(4, v); }
  @$pb.TagNumber(5)
  $core.bool hasUuid() => $_has(4);
  @$pb.TagNumber(5)
  void clearUuid() => $_clearField(5);

  /// Unix timestamp of wallet creation.
  @$pb.TagNumber(6)
  $fixnum.Int64 get createdAt => $_getI64(5);
  @$pb.TagNumber(6)
  set createdAt($fixnum.Int64 v) { $_setInt64(5, v); }
  @$pb.TagNumber(6)
  $core.bool hasCreatedAt() => $_has(5);
  @$pb.TagNumber(6)
  void clearCreatedAt() => $_clearField(6);

  /// The default fee of the wallet.
  @$pb.TagNumber(7)
  $fixnum.Int64 get defaultFee => $_getI64(6);
  @$pb.TagNumber(7)
  set defaultFee($fixnum.Int64 v) { $_setInt64(6, v); }
  @$pb.TagNumber(7)
  $core.bool hasDefaultFee() => $_has(6);
  @$pb.TagNumber(7)
  void clearDefaultFee() => $_clearField(7);
}

/// Request message for listing wallet addresses.
class ListAddressRequest extends $pb.GeneratedMessage {
  factory ListAddressRequest({
    $core.String? walletName,
  }) {
    final $result = create();
    if (walletName != null) {
      $result.walletName = walletName;
    }
    return $result;
  }
  ListAddressRequest._() : super();
  factory ListAddressRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory ListAddressRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'ListAddressRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'walletName')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  ListAddressRequest clone() => ListAddressRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  ListAddressRequest copyWith(void Function(ListAddressRequest) updates) => super.copyWith((message) => updates(message as ListAddressRequest)) as ListAddressRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListAddressRequest create() => ListAddressRequest._();
  ListAddressRequest createEmptyInstance() => create();
  static $pb.PbList<ListAddressRequest> createRepeated() => $pb.PbList<ListAddressRequest>();
  @$core.pragma('dart2js:noInline')
  static ListAddressRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ListAddressRequest>(create);
  static ListAddressRequest? _defaultInstance;

  /// The name of the queried wallet.
  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => $_clearField(1);
}

/// Response message contains wallet addresses.
class ListAddressResponse extends $pb.GeneratedMessage {
  factory ListAddressResponse({
    $core.String? walletName,
    $core.Iterable<AddressInfo>? data,
  }) {
    final $result = create();
    if (walletName != null) {
      $result.walletName = walletName;
    }
    if (data != null) {
      $result.data.addAll(data);
    }
    return $result;
  }
  ListAddressResponse._() : super();
  factory ListAddressResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory ListAddressResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'ListAddressResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'walletName')
    ..pc<AddressInfo>(2, _omitFieldNames ? '' : 'data', $pb.PbFieldType.PM, subBuilder: AddressInfo.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  ListAddressResponse clone() => ListAddressResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  ListAddressResponse copyWith(void Function(ListAddressResponse) updates) => super.copyWith((message) => updates(message as ListAddressResponse)) as ListAddressResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListAddressResponse create() => ListAddressResponse._();
  ListAddressResponse createEmptyInstance() => create();
  static $pb.PbList<ListAddressResponse> createRepeated() => $pb.PbList<ListAddressResponse>();
  @$core.pragma('dart2js:noInline')
  static ListAddressResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ListAddressResponse>(create);
  static ListAddressResponse? _defaultInstance;

  /// The name of the queried wallet.
  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => $_clearField(1);

  /// List of all addresses in the wallet with their details.
  @$pb.TagNumber(2)
  $pb.PbList<AddressInfo> get data => $_getList(1);
}

/// Wallet service provides RPC methods for wallet management operations.
class WalletApi {
  $pb.RpcClient _client;
  WalletApi(this._client);

  /// CreateWallet creates a new wallet with the specified parameters.
  $async.Future<CreateWalletResponse> createWallet($pb.ClientContext? ctx, CreateWalletRequest request) =>
    _client.invoke<CreateWalletResponse>(ctx, 'Wallet', 'CreateWallet', request, CreateWalletResponse())
  ;
  /// RestoreWallet restores an existing wallet with the given mnemonic.
  $async.Future<RestoreWalletResponse> restoreWallet($pb.ClientContext? ctx, RestoreWalletRequest request) =>
    _client.invoke<RestoreWalletResponse>(ctx, 'Wallet', 'RestoreWallet', request, RestoreWalletResponse())
  ;
  /// LoadWallet loads an existing wallet with the given name.
  $async.Future<LoadWalletResponse> loadWallet($pb.ClientContext? ctx, LoadWalletRequest request) =>
    _client.invoke<LoadWalletResponse>(ctx, 'Wallet', 'LoadWallet', request, LoadWalletResponse())
  ;
  /// UnloadWallet unloads a currently loaded wallet with the specified name.
  $async.Future<UnloadWalletResponse> unloadWallet($pb.ClientContext? ctx, UnloadWalletRequest request) =>
    _client.invoke<UnloadWalletResponse>(ctx, 'Wallet', 'UnloadWallet', request, UnloadWalletResponse())
  ;
  /// GetTotalBalance returns the total available balance of the wallet.
  $async.Future<GetTotalBalanceResponse> getTotalBalance($pb.ClientContext? ctx, GetTotalBalanceRequest request) =>
    _client.invoke<GetTotalBalanceResponse>(ctx, 'Wallet', 'GetTotalBalance', request, GetTotalBalanceResponse())
  ;
  /// SignRawTransaction signs a raw transaction for a specified wallet.
  $async.Future<SignRawTransactionResponse> signRawTransaction($pb.ClientContext? ctx, SignRawTransactionRequest request) =>
    _client.invoke<SignRawTransactionResponse>(ctx, 'Wallet', 'SignRawTransaction', request, SignRawTransactionResponse())
  ;
  /// GetValidatorAddress retrieves the validator address associated with a public key.
  $async.Future<GetValidatorAddressResponse> getValidatorAddress($pb.ClientContext? ctx, GetValidatorAddressRequest request) =>
    _client.invoke<GetValidatorAddressResponse>(ctx, 'Wallet', 'GetValidatorAddress', request, GetValidatorAddressResponse())
  ;
  /// GetNewAddress generates a new address for the specified wallet.
  $async.Future<GetNewAddressResponse> getNewAddress($pb.ClientContext? ctx, GetNewAddressRequest request) =>
    _client.invoke<GetNewAddressResponse>(ctx, 'Wallet', 'GetNewAddress', request, GetNewAddressResponse())
  ;
  /// GetAddressHistory retrieves the transaction history of an address.
  $async.Future<GetAddressHistoryResponse> getAddressHistory($pb.ClientContext? ctx, GetAddressHistoryRequest request) =>
    _client.invoke<GetAddressHistoryResponse>(ctx, 'Wallet', 'GetAddressHistory', request, GetAddressHistoryResponse())
  ;
  /// SignMessage signs an arbitrary message using a wallet's private key.
  $async.Future<SignMessageResponse> signMessage($pb.ClientContext? ctx, SignMessageRequest request) =>
    _client.invoke<SignMessageResponse>(ctx, 'Wallet', 'SignMessage', request, SignMessageResponse())
  ;
  /// GetTotalStake returns the total stake amount in the wallet.
  $async.Future<GetTotalStakeResponse> getTotalStake($pb.ClientContext? ctx, GetTotalStakeRequest request) =>
    _client.invoke<GetTotalStakeResponse>(ctx, 'Wallet', 'GetTotalStake', request, GetTotalStakeResponse())
  ;
  /// GetAddressInfo returns detailed information about a specific address.
  $async.Future<GetAddressInfoResponse> getAddressInfo($pb.ClientContext? ctx, GetAddressInfoRequest request) =>
    _client.invoke<GetAddressInfoResponse>(ctx, 'Wallet', 'GetAddressInfo', request, GetAddressInfoResponse())
  ;
  /// SetAddressLabel sets or updates the label for a given address.
  $async.Future<SetAddressLabelResponse> setAddressLabel($pb.ClientContext? ctx, SetAddressLabelRequest request) =>
    _client.invoke<SetAddressLabelResponse>(ctx, 'Wallet', 'SetAddressLabel', request, SetAddressLabelResponse())
  ;
  /// ListWallet returns a list of all available wallets.
  $async.Future<ListWalletResponse> listWallet($pb.ClientContext? ctx, ListWalletRequest request) =>
    _client.invoke<ListWalletResponse>(ctx, 'Wallet', 'ListWallet', request, ListWalletResponse())
  ;
  /// GetWalletInfo returns detailed information about a specific wallet.
  $async.Future<GetWalletInfoResponse> getWalletInfo($pb.ClientContext? ctx, GetWalletInfoRequest request) =>
    _client.invoke<GetWalletInfoResponse>(ctx, 'Wallet', 'GetWalletInfo', request, GetWalletInfoResponse())
  ;
  /// ListAddress returns all addresses in the specified wallet.
  $async.Future<ListAddressResponse> listAddress($pb.ClientContext? ctx, ListAddressRequest request) =>
    _client.invoke<ListAddressResponse>(ctx, 'Wallet', 'ListAddress', request, ListAddressResponse())
  ;
}


const _omitFieldNames = $core.bool.fromEnvironment('protobuf.omit_field_names');
const _omitMessageNames = $core.bool.fromEnvironment('protobuf.omit_message_names');
