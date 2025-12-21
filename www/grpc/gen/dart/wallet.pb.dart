// This is a generated file - do not edit.
//
// Generated from wallet.proto.

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
    final result = create();
    if (address != null) result.address = address;
    if (publicKey != null) result.publicKey = publicKey;
    if (label != null) result.label = label;
    if (path != null) result.path = path;
    return result;
  }

  AddressInfo._();

  factory AddressInfo.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AddressInfo.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AddressInfo',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'address')
    ..aOS(2, _omitFieldNames ? '' : 'publicKey')
    ..aOS(3, _omitFieldNames ? '' : 'label')
    ..aOS(4, _omitFieldNames ? '' : 'path')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AddressInfo clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AddressInfo copyWith(void Function(AddressInfo) updates) =>
      super.copyWith((message) => updates(message as AddressInfo))
          as AddressInfo;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AddressInfo create() => AddressInfo._();
  @$core.override
  AddressInfo createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AddressInfo getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AddressInfo>(create);
  static AddressInfo? _defaultInstance;

  /// The address string.
  @$pb.TagNumber(1)
  $core.String get address => $_getSZ(0);
  @$pb.TagNumber(1)
  set address($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasAddress() => $_has(0);
  @$pb.TagNumber(1)
  void clearAddress() => $_clearField(1);

  /// The public key associated with the address.
  @$pb.TagNumber(2)
  $core.String get publicKey => $_getSZ(1);
  @$pb.TagNumber(2)
  set publicKey($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasPublicKey() => $_has(1);
  @$pb.TagNumber(2)
  void clearPublicKey() => $_clearField(2);

  /// A human-readable label associated with the address.
  @$pb.TagNumber(3)
  $core.String get label => $_getSZ(2);
  @$pb.TagNumber(3)
  set label($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasLabel() => $_has(2);
  @$pb.TagNumber(3)
  void clearLabel() => $_clearField(3);

  /// The Hierarchical Deterministic (HD) path of the address within the wallet.
  @$pb.TagNumber(4)
  $core.String get path => $_getSZ(3);
  @$pb.TagNumber(4)
  set path($core.String value) => $_setString(3, value);
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
    final result = create();
    if (transactionId != null) result.transactionId = transactionId;
    if (time != null) result.time = time;
    if (payloadType != null) result.payloadType = payloadType;
    if (description != null) result.description = description;
    if (amount != null) result.amount = amount;
    return result;
  }

  HistoryInfo._();

  factory HistoryInfo.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory HistoryInfo.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'HistoryInfo',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'transactionId')
    ..aI(2, _omitFieldNames ? '' : 'time', fieldType: $pb.PbFieldType.OU3)
    ..aOS(3, _omitFieldNames ? '' : 'payloadType')
    ..aOS(4, _omitFieldNames ? '' : 'description')
    ..aInt64(5, _omitFieldNames ? '' : 'amount')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HistoryInfo clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  HistoryInfo copyWith(void Function(HistoryInfo) updates) =>
      super.copyWith((message) => updates(message as HistoryInfo))
          as HistoryInfo;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static HistoryInfo create() => HistoryInfo._();
  @$core.override
  HistoryInfo createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static HistoryInfo getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<HistoryInfo>(create);
  static HistoryInfo? _defaultInstance;

  /// The transaction ID in hexadecimal format.
  @$pb.TagNumber(1)
  $core.String get transactionId => $_getSZ(0);
  @$pb.TagNumber(1)
  set transactionId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasTransactionId() => $_has(0);
  @$pb.TagNumber(1)
  void clearTransactionId() => $_clearField(1);

  /// Unix timestamp of when the transaction was confirmed.
  @$pb.TagNumber(2)
  $core.int get time => $_getIZ(1);
  @$pb.TagNumber(2)
  set time($core.int value) => $_setUnsignedInt32(1, value);
  @$pb.TagNumber(2)
  $core.bool hasTime() => $_has(1);
  @$pb.TagNumber(2)
  void clearTime() => $_clearField(2);

  /// The type of transaction payload.
  @$pb.TagNumber(3)
  $core.String get payloadType => $_getSZ(2);
  @$pb.TagNumber(3)
  set payloadType($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasPayloadType() => $_has(2);
  @$pb.TagNumber(3)
  void clearPayloadType() => $_clearField(3);

  /// Human-readable description of the transaction.
  @$pb.TagNumber(4)
  $core.String get description => $_getSZ(3);
  @$pb.TagNumber(4)
  set description($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasDescription() => $_has(3);
  @$pb.TagNumber(4)
  void clearDescription() => $_clearField(4);

  /// The transaction amount in NanoPAC.
  @$pb.TagNumber(5)
  $fixnum.Int64 get amount => $_getI64(4);
  @$pb.TagNumber(5)
  set amount($fixnum.Int64 value) => $_setInt64(4, value);
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
    final result = create();
    if (walletName != null) result.walletName = walletName;
    if (address != null) result.address = address;
    return result;
  }

  GetAddressHistoryRequest._();

  factory GetAddressHistoryRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetAddressHistoryRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetAddressHistoryRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'walletName')
    ..aOS(2, _omitFieldNames ? '' : 'address')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetAddressHistoryRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetAddressHistoryRequest copyWith(
          void Function(GetAddressHistoryRequest) updates) =>
      super.copyWith((message) => updates(message as GetAddressHistoryRequest))
          as GetAddressHistoryRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetAddressHistoryRequest create() => GetAddressHistoryRequest._();
  @$core.override
  GetAddressHistoryRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetAddressHistoryRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetAddressHistoryRequest>(create);
  static GetAddressHistoryRequest? _defaultInstance;

  /// The name of the wallet containing the address.
  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => $_clearField(1);

  /// The address to retrieve history for.
  @$pb.TagNumber(2)
  $core.String get address => $_getSZ(1);
  @$pb.TagNumber(2)
  set address($core.String value) => $_setString(1, value);
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
    final result = create();
    if (historyInfo != null) result.historyInfo.addAll(historyInfo);
    return result;
  }

  GetAddressHistoryResponse._();

  factory GetAddressHistoryResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetAddressHistoryResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetAddressHistoryResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..pPM<HistoryInfo>(1, _omitFieldNames ? '' : 'historyInfo',
        subBuilder: HistoryInfo.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetAddressHistoryResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetAddressHistoryResponse copyWith(
          void Function(GetAddressHistoryResponse) updates) =>
      super.copyWith((message) => updates(message as GetAddressHistoryResponse))
          as GetAddressHistoryResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetAddressHistoryResponse create() => GetAddressHistoryResponse._();
  @$core.override
  GetAddressHistoryResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetAddressHistoryResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetAddressHistoryResponse>(create);
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
    final result = create();
    if (walletName != null) result.walletName = walletName;
    if (addressType != null) result.addressType = addressType;
    if (label != null) result.label = label;
    if (password != null) result.password = password;
    return result;
  }

  GetNewAddressRequest._();

  factory GetNewAddressRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetNewAddressRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetNewAddressRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'walletName')
    ..aE<AddressType>(2, _omitFieldNames ? '' : 'addressType',
        enumValues: AddressType.values)
    ..aOS(3, _omitFieldNames ? '' : 'label')
    ..aOS(4, _omitFieldNames ? '' : 'password')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetNewAddressRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetNewAddressRequest copyWith(void Function(GetNewAddressRequest) updates) =>
      super.copyWith((message) => updates(message as GetNewAddressRequest))
          as GetNewAddressRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetNewAddressRequest create() => GetNewAddressRequest._();
  @$core.override
  GetNewAddressRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetNewAddressRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetNewAddressRequest>(create);
  static GetNewAddressRequest? _defaultInstance;

  /// The name of the wallet to generate a new address.
  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => $_clearField(1);

  /// The type of address to generate.
  @$pb.TagNumber(2)
  AddressType get addressType => $_getN(1);
  @$pb.TagNumber(2)
  set addressType(AddressType value) => $_setField(2, value);
  @$pb.TagNumber(2)
  $core.bool hasAddressType() => $_has(1);
  @$pb.TagNumber(2)
  void clearAddressType() => $_clearField(2);

  /// A label for the new address.
  @$pb.TagNumber(3)
  $core.String get label => $_getSZ(2);
  @$pb.TagNumber(3)
  set label($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasLabel() => $_has(2);
  @$pb.TagNumber(3)
  void clearLabel() => $_clearField(3);

  /// Password for the new address. It's required when address_type is Ed25519 type.
  @$pb.TagNumber(4)
  $core.String get password => $_getSZ(3);
  @$pb.TagNumber(4)
  set password($core.String value) => $_setString(3, value);
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
    final result = create();
    if (walletName != null) result.walletName = walletName;
    if (addressInfo != null) result.addressInfo = addressInfo;
    return result;
  }

  GetNewAddressResponse._();

  factory GetNewAddressResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetNewAddressResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetNewAddressResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'walletName')
    ..aOM<AddressInfo>(2, _omitFieldNames ? '' : 'addressInfo',
        subBuilder: AddressInfo.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetNewAddressResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetNewAddressResponse copyWith(
          void Function(GetNewAddressResponse) updates) =>
      super.copyWith((message) => updates(message as GetNewAddressResponse))
          as GetNewAddressResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetNewAddressResponse create() => GetNewAddressResponse._();
  @$core.override
  GetNewAddressResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetNewAddressResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetNewAddressResponse>(create);
  static GetNewAddressResponse? _defaultInstance;

  /// The name of the wallet where address was generated.
  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => $_clearField(1);

  /// Detailed information about the new address.
  @$pb.TagNumber(2)
  AddressInfo get addressInfo => $_getN(1);
  @$pb.TagNumber(2)
  set addressInfo(AddressInfo value) => $_setField(2, value);
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
    final result = create();
    if (walletName != null) result.walletName = walletName;
    if (mnemonic != null) result.mnemonic = mnemonic;
    if (password != null) result.password = password;
    return result;
  }

  RestoreWalletRequest._();

  factory RestoreWalletRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RestoreWalletRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RestoreWalletRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'walletName')
    ..aOS(2, _omitFieldNames ? '' : 'mnemonic')
    ..aOS(3, _omitFieldNames ? '' : 'password')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RestoreWalletRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RestoreWalletRequest copyWith(void Function(RestoreWalletRequest) updates) =>
      super.copyWith((message) => updates(message as RestoreWalletRequest))
          as RestoreWalletRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RestoreWalletRequest create() => RestoreWalletRequest._();
  @$core.override
  RestoreWalletRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RestoreWalletRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RestoreWalletRequest>(create);
  static RestoreWalletRequest? _defaultInstance;

  /// The name for the restored wallet.
  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => $_clearField(1);

  /// The mnemonic (seed phrase) for wallet recovery.
  @$pb.TagNumber(2)
  $core.String get mnemonic => $_getSZ(1);
  @$pb.TagNumber(2)
  set mnemonic($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasMnemonic() => $_has(1);
  @$pb.TagNumber(2)
  void clearMnemonic() => $_clearField(2);

  /// Password to secure the restored wallet.
  @$pb.TagNumber(3)
  $core.String get password => $_getSZ(2);
  @$pb.TagNumber(3)
  set password($core.String value) => $_setString(2, value);
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
    final result = create();
    if (walletName != null) result.walletName = walletName;
    return result;
  }

  RestoreWalletResponse._();

  factory RestoreWalletResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory RestoreWalletResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'RestoreWalletResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'walletName')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RestoreWalletResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  RestoreWalletResponse copyWith(
          void Function(RestoreWalletResponse) updates) =>
      super.copyWith((message) => updates(message as RestoreWalletResponse))
          as RestoreWalletResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static RestoreWalletResponse create() => RestoreWalletResponse._();
  @$core.override
  RestoreWalletResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static RestoreWalletResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<RestoreWalletResponse>(create);
  static RestoreWalletResponse? _defaultInstance;

  /// The name of the restored wallet.
  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String value) => $_setString(0, value);
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
    final result = create();
    if (walletName != null) result.walletName = walletName;
    if (password != null) result.password = password;
    return result;
  }

  CreateWalletRequest._();

  factory CreateWalletRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateWalletRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateWalletRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'walletName')
    ..aOS(2, _omitFieldNames ? '' : 'password')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateWalletRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateWalletRequest copyWith(void Function(CreateWalletRequest) updates) =>
      super.copyWith((message) => updates(message as CreateWalletRequest))
          as CreateWalletRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateWalletRequest create() => CreateWalletRequest._();
  @$core.override
  CreateWalletRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateWalletRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateWalletRequest>(create);
  static CreateWalletRequest? _defaultInstance;

  /// The name for the new wallet.
  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => $_clearField(1);

  /// Password to secure the new wallet.
  @$pb.TagNumber(2)
  $core.String get password => $_getSZ(1);
  @$pb.TagNumber(2)
  set password($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasPassword() => $_has(1);
  @$pb.TagNumber(2)
  void clearPassword() => $_clearField(2);
}

/// Response message contains wallet recovery mnemonic (seed phrase).
class CreateWalletResponse extends $pb.GeneratedMessage {
  factory CreateWalletResponse({
    $core.String? walletName,
    $core.String? mnemonic,
  }) {
    final result = create();
    if (walletName != null) result.walletName = walletName;
    if (mnemonic != null) result.mnemonic = mnemonic;
    return result;
  }

  CreateWalletResponse._();

  factory CreateWalletResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CreateWalletResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CreateWalletResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'walletName')
    ..aOS(2, _omitFieldNames ? '' : 'mnemonic')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateWalletResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CreateWalletResponse copyWith(void Function(CreateWalletResponse) updates) =>
      super.copyWith((message) => updates(message as CreateWalletResponse))
          as CreateWalletResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CreateWalletResponse create() => CreateWalletResponse._();
  @$core.override
  CreateWalletResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CreateWalletResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CreateWalletResponse>(create);
  static CreateWalletResponse? _defaultInstance;

  /// The name for the new wallet.
  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => $_clearField(1);

  /// The mnemonic (seed phrase) for wallet recovery.
  @$pb.TagNumber(2)
  $core.String get mnemonic => $_getSZ(1);
  @$pb.TagNumber(2)
  set mnemonic($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasMnemonic() => $_has(1);
  @$pb.TagNumber(2)
  void clearMnemonic() => $_clearField(2);
}

/// Request message for loading an existing wallet.
class LoadWalletRequest extends $pb.GeneratedMessage {
  factory LoadWalletRequest({
    $core.String? walletName,
  }) {
    final result = create();
    if (walletName != null) result.walletName = walletName;
    return result;
  }

  LoadWalletRequest._();

  factory LoadWalletRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory LoadWalletRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'LoadWalletRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'walletName')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LoadWalletRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LoadWalletRequest copyWith(void Function(LoadWalletRequest) updates) =>
      super.copyWith((message) => updates(message as LoadWalletRequest))
          as LoadWalletRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static LoadWalletRequest create() => LoadWalletRequest._();
  @$core.override
  LoadWalletRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static LoadWalletRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<LoadWalletRequest>(create);
  static LoadWalletRequest? _defaultInstance;

  /// The name of the wallet to load.
  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String value) => $_setString(0, value);
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
    final result = create();
    if (walletName != null) result.walletName = walletName;
    return result;
  }

  LoadWalletResponse._();

  factory LoadWalletResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory LoadWalletResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'LoadWalletResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'walletName')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LoadWalletResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  LoadWalletResponse copyWith(void Function(LoadWalletResponse) updates) =>
      super.copyWith((message) => updates(message as LoadWalletResponse))
          as LoadWalletResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static LoadWalletResponse create() => LoadWalletResponse._();
  @$core.override
  LoadWalletResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static LoadWalletResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<LoadWalletResponse>(create);
  static LoadWalletResponse? _defaultInstance;

  /// The name of the loaded wallet.
  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String value) => $_setString(0, value);
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
    final result = create();
    if (walletName != null) result.walletName = walletName;
    return result;
  }

  UnloadWalletRequest._();

  factory UnloadWalletRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UnloadWalletRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UnloadWalletRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'walletName')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UnloadWalletRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UnloadWalletRequest copyWith(void Function(UnloadWalletRequest) updates) =>
      super.copyWith((message) => updates(message as UnloadWalletRequest))
          as UnloadWalletRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UnloadWalletRequest create() => UnloadWalletRequest._();
  @$core.override
  UnloadWalletRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UnloadWalletRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UnloadWalletRequest>(create);
  static UnloadWalletRequest? _defaultInstance;

  /// The name of the wallet to unload.
  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String value) => $_setString(0, value);
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
    final result = create();
    if (walletName != null) result.walletName = walletName;
    return result;
  }

  UnloadWalletResponse._();

  factory UnloadWalletResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory UnloadWalletResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'UnloadWalletResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'walletName')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UnloadWalletResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  UnloadWalletResponse copyWith(void Function(UnloadWalletResponse) updates) =>
      super.copyWith((message) => updates(message as UnloadWalletResponse))
          as UnloadWalletResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static UnloadWalletResponse create() => UnloadWalletResponse._();
  @$core.override
  UnloadWalletResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static UnloadWalletResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<UnloadWalletResponse>(create);
  static UnloadWalletResponse? _defaultInstance;

  /// The name of the unloaded wallet.
  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String value) => $_setString(0, value);
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
    final result = create();
    if (publicKey != null) result.publicKey = publicKey;
    return result;
  }

  GetValidatorAddressRequest._();

  factory GetValidatorAddressRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetValidatorAddressRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetValidatorAddressRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'publicKey')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetValidatorAddressRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetValidatorAddressRequest copyWith(
          void Function(GetValidatorAddressRequest) updates) =>
      super.copyWith(
              (message) => updates(message as GetValidatorAddressRequest))
          as GetValidatorAddressRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetValidatorAddressRequest create() => GetValidatorAddressRequest._();
  @$core.override
  GetValidatorAddressRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetValidatorAddressRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetValidatorAddressRequest>(create);
  static GetValidatorAddressRequest? _defaultInstance;

  /// The public key of the validator.
  @$pb.TagNumber(1)
  $core.String get publicKey => $_getSZ(0);
  @$pb.TagNumber(1)
  set publicKey($core.String value) => $_setString(0, value);
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
    final result = create();
    if (address != null) result.address = address;
    return result;
  }

  GetValidatorAddressResponse._();

  factory GetValidatorAddressResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetValidatorAddressResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetValidatorAddressResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'address')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetValidatorAddressResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetValidatorAddressResponse copyWith(
          void Function(GetValidatorAddressResponse) updates) =>
      super.copyWith(
              (message) => updates(message as GetValidatorAddressResponse))
          as GetValidatorAddressResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetValidatorAddressResponse create() =>
      GetValidatorAddressResponse._();
  @$core.override
  GetValidatorAddressResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetValidatorAddressResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetValidatorAddressResponse>(create);
  static GetValidatorAddressResponse? _defaultInstance;

  /// The validator address associated with the public key.
  @$pb.TagNumber(1)
  $core.String get address => $_getSZ(0);
  @$pb.TagNumber(1)
  set address($core.String value) => $_setString(0, value);
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
    final result = create();
    if (walletName != null) result.walletName = walletName;
    if (rawTransaction != null) result.rawTransaction = rawTransaction;
    if (password != null) result.password = password;
    return result;
  }

  SignRawTransactionRequest._();

  factory SignRawTransactionRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SignRawTransactionRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SignRawTransactionRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'walletName')
    ..aOS(2, _omitFieldNames ? '' : 'rawTransaction')
    ..aOS(3, _omitFieldNames ? '' : 'password')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SignRawTransactionRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SignRawTransactionRequest copyWith(
          void Function(SignRawTransactionRequest) updates) =>
      super.copyWith((message) => updates(message as SignRawTransactionRequest))
          as SignRawTransactionRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SignRawTransactionRequest create() => SignRawTransactionRequest._();
  @$core.override
  SignRawTransactionRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SignRawTransactionRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SignRawTransactionRequest>(create);
  static SignRawTransactionRequest? _defaultInstance;

  /// The name of the wallet used for signing.
  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => $_clearField(1);

  /// The raw transaction data to be signed.
  @$pb.TagNumber(2)
  $core.String get rawTransaction => $_getSZ(1);
  @$pb.TagNumber(2)
  set rawTransaction($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasRawTransaction() => $_has(1);
  @$pb.TagNumber(2)
  void clearRawTransaction() => $_clearField(2);

  /// Wallet password required for signing.
  @$pb.TagNumber(3)
  $core.String get password => $_getSZ(2);
  @$pb.TagNumber(3)
  set password($core.String value) => $_setString(2, value);
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
    final result = create();
    if (transactionId != null) result.transactionId = transactionId;
    if (signedRawTransaction != null)
      result.signedRawTransaction = signedRawTransaction;
    return result;
  }

  SignRawTransactionResponse._();

  factory SignRawTransactionResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SignRawTransactionResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SignRawTransactionResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'transactionId')
    ..aOS(2, _omitFieldNames ? '' : 'signedRawTransaction')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SignRawTransactionResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SignRawTransactionResponse copyWith(
          void Function(SignRawTransactionResponse) updates) =>
      super.copyWith(
              (message) => updates(message as SignRawTransactionResponse))
          as SignRawTransactionResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SignRawTransactionResponse create() => SignRawTransactionResponse._();
  @$core.override
  SignRawTransactionResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SignRawTransactionResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SignRawTransactionResponse>(create);
  static SignRawTransactionResponse? _defaultInstance;

  /// The ID of the signed transaction.
  @$pb.TagNumber(1)
  $core.String get transactionId => $_getSZ(0);
  @$pb.TagNumber(1)
  set transactionId($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasTransactionId() => $_has(0);
  @$pb.TagNumber(1)
  void clearTransactionId() => $_clearField(1);

  /// The signed raw transaction data.
  @$pb.TagNumber(2)
  $core.String get signedRawTransaction => $_getSZ(1);
  @$pb.TagNumber(2)
  set signedRawTransaction($core.String value) => $_setString(1, value);
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
    final result = create();
    if (walletName != null) result.walletName = walletName;
    return result;
  }

  GetTotalBalanceRequest._();

  factory GetTotalBalanceRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetTotalBalanceRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetTotalBalanceRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'walletName')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetTotalBalanceRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetTotalBalanceRequest copyWith(
          void Function(GetTotalBalanceRequest) updates) =>
      super.copyWith((message) => updates(message as GetTotalBalanceRequest))
          as GetTotalBalanceRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetTotalBalanceRequest create() => GetTotalBalanceRequest._();
  @$core.override
  GetTotalBalanceRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetTotalBalanceRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetTotalBalanceRequest>(create);
  static GetTotalBalanceRequest? _defaultInstance;

  /// The name of the wallet to get the total balance.
  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String value) => $_setString(0, value);
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
    final result = create();
    if (walletName != null) result.walletName = walletName;
    if (totalBalance != null) result.totalBalance = totalBalance;
    return result;
  }

  GetTotalBalanceResponse._();

  factory GetTotalBalanceResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetTotalBalanceResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetTotalBalanceResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'walletName')
    ..aInt64(2, _omitFieldNames ? '' : 'totalBalance')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetTotalBalanceResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetTotalBalanceResponse copyWith(
          void Function(GetTotalBalanceResponse) updates) =>
      super.copyWith((message) => updates(message as GetTotalBalanceResponse))
          as GetTotalBalanceResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetTotalBalanceResponse create() => GetTotalBalanceResponse._();
  @$core.override
  GetTotalBalanceResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetTotalBalanceResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetTotalBalanceResponse>(create);
  static GetTotalBalanceResponse? _defaultInstance;

  /// The name of the queried wallet.
  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => $_clearField(1);

  /// The total balance of the wallet in NanoPAC.
  @$pb.TagNumber(2)
  $fixnum.Int64 get totalBalance => $_getI64(1);
  @$pb.TagNumber(2)
  set totalBalance($fixnum.Int64 value) => $_setInt64(1, value);
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
    final result = create();
    if (walletName != null) result.walletName = walletName;
    if (password != null) result.password = password;
    if (address != null) result.address = address;
    if (message != null) result.message = message;
    return result;
  }

  SignMessageRequest._();

  factory SignMessageRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SignMessageRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SignMessageRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'walletName')
    ..aOS(2, _omitFieldNames ? '' : 'password')
    ..aOS(3, _omitFieldNames ? '' : 'address')
    ..aOS(4, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SignMessageRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SignMessageRequest copyWith(void Function(SignMessageRequest) updates) =>
      super.copyWith((message) => updates(message as SignMessageRequest))
          as SignMessageRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SignMessageRequest create() => SignMessageRequest._();
  @$core.override
  SignMessageRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SignMessageRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SignMessageRequest>(create);
  static SignMessageRequest? _defaultInstance;

  /// The name of the wallet to sign with.
  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => $_clearField(1);

  /// Wallet password required for signing.
  @$pb.TagNumber(2)
  $core.String get password => $_getSZ(1);
  @$pb.TagNumber(2)
  set password($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasPassword() => $_has(1);
  @$pb.TagNumber(2)
  void clearPassword() => $_clearField(2);

  /// The address whose private key should be used for signing the message.
  @$pb.TagNumber(3)
  $core.String get address => $_getSZ(2);
  @$pb.TagNumber(3)
  set address($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasAddress() => $_has(2);
  @$pb.TagNumber(3)
  void clearAddress() => $_clearField(3);

  /// The arbitrary message to be signed.
  @$pb.TagNumber(4)
  $core.String get message => $_getSZ(3);
  @$pb.TagNumber(4)
  set message($core.String value) => $_setString(3, value);
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
    final result = create();
    if (signature != null) result.signature = signature;
    return result;
  }

  SignMessageResponse._();

  factory SignMessageResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SignMessageResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SignMessageResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'signature')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SignMessageResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SignMessageResponse copyWith(void Function(SignMessageResponse) updates) =>
      super.copyWith((message) => updates(message as SignMessageResponse))
          as SignMessageResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SignMessageResponse create() => SignMessageResponse._();
  @$core.override
  SignMessageResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SignMessageResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SignMessageResponse>(create);
  static SignMessageResponse? _defaultInstance;

  /// The signature in hexadecimal format.
  @$pb.TagNumber(1)
  $core.String get signature => $_getSZ(0);
  @$pb.TagNumber(1)
  set signature($core.String value) => $_setString(0, value);
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
    final result = create();
    if (walletName != null) result.walletName = walletName;
    return result;
  }

  GetTotalStakeRequest._();

  factory GetTotalStakeRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetTotalStakeRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetTotalStakeRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'walletName')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetTotalStakeRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetTotalStakeRequest copyWith(void Function(GetTotalStakeRequest) updates) =>
      super.copyWith((message) => updates(message as GetTotalStakeRequest))
          as GetTotalStakeRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetTotalStakeRequest create() => GetTotalStakeRequest._();
  @$core.override
  GetTotalStakeRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetTotalStakeRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetTotalStakeRequest>(create);
  static GetTotalStakeRequest? _defaultInstance;

  /// The name of the wallet to get the total stake.
  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String value) => $_setString(0, value);
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
    final result = create();
    if (walletName != null) result.walletName = walletName;
    if (totalStake != null) result.totalStake = totalStake;
    return result;
  }

  GetTotalStakeResponse._();

  factory GetTotalStakeResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetTotalStakeResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetTotalStakeResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'walletName')
    ..aInt64(2, _omitFieldNames ? '' : 'totalStake')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetTotalStakeResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetTotalStakeResponse copyWith(
          void Function(GetTotalStakeResponse) updates) =>
      super.copyWith((message) => updates(message as GetTotalStakeResponse))
          as GetTotalStakeResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetTotalStakeResponse create() => GetTotalStakeResponse._();
  @$core.override
  GetTotalStakeResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetTotalStakeResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetTotalStakeResponse>(create);
  static GetTotalStakeResponse? _defaultInstance;

  /// The name of the queried wallet.
  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => $_clearField(1);

  /// The total stake amount in NanoPAC.
  @$pb.TagNumber(2)
  $fixnum.Int64 get totalStake => $_getI64(1);
  @$pb.TagNumber(2)
  set totalStake($fixnum.Int64 value) => $_setInt64(1, value);
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
    final result = create();
    if (walletName != null) result.walletName = walletName;
    if (address != null) result.address = address;
    return result;
  }

  GetAddressInfoRequest._();

  factory GetAddressInfoRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetAddressInfoRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetAddressInfoRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'walletName')
    ..aOS(2, _omitFieldNames ? '' : 'address')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetAddressInfoRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetAddressInfoRequest copyWith(
          void Function(GetAddressInfoRequest) updates) =>
      super.copyWith((message) => updates(message as GetAddressInfoRequest))
          as GetAddressInfoRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetAddressInfoRequest create() => GetAddressInfoRequest._();
  @$core.override
  GetAddressInfoRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetAddressInfoRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetAddressInfoRequest>(create);
  static GetAddressInfoRequest? _defaultInstance;

  /// The name of the wallet containing the address.
  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => $_clearField(1);

  /// The address to query.
  @$pb.TagNumber(2)
  $core.String get address => $_getSZ(1);
  @$pb.TagNumber(2)
  set address($core.String value) => $_setString(1, value);
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
    final result = create();
    if (walletName != null) result.walletName = walletName;
    if (address != null) result.address = address;
    if (label != null) result.label = label;
    if (publicKey != null) result.publicKey = publicKey;
    if (path != null) result.path = path;
    return result;
  }

  GetAddressInfoResponse._();

  factory GetAddressInfoResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetAddressInfoResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetAddressInfoResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'walletName')
    ..aOS(2, _omitFieldNames ? '' : 'address')
    ..aOS(3, _omitFieldNames ? '' : 'label')
    ..aOS(4, _omitFieldNames ? '' : 'publicKey')
    ..aOS(5, _omitFieldNames ? '' : 'path')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetAddressInfoResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetAddressInfoResponse copyWith(
          void Function(GetAddressInfoResponse) updates) =>
      super.copyWith((message) => updates(message as GetAddressInfoResponse))
          as GetAddressInfoResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetAddressInfoResponse create() => GetAddressInfoResponse._();
  @$core.override
  GetAddressInfoResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetAddressInfoResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetAddressInfoResponse>(create);
  static GetAddressInfoResponse? _defaultInstance;

  /// The name of the wallet containing the address.
  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => $_clearField(1);

  /// The queried address.
  @$pb.TagNumber(2)
  $core.String get address => $_getSZ(1);
  @$pb.TagNumber(2)
  set address($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasAddress() => $_has(1);
  @$pb.TagNumber(2)
  void clearAddress() => $_clearField(2);

  /// The address label.
  @$pb.TagNumber(3)
  $core.String get label => $_getSZ(2);
  @$pb.TagNumber(3)
  set label($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasLabel() => $_has(2);
  @$pb.TagNumber(3)
  void clearLabel() => $_clearField(3);

  /// The public key of the address.
  @$pb.TagNumber(4)
  $core.String get publicKey => $_getSZ(3);
  @$pb.TagNumber(4)
  set publicKey($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasPublicKey() => $_has(3);
  @$pb.TagNumber(4)
  void clearPublicKey() => $_clearField(4);

  /// The Hierarchical Deterministic (HD) path of the address.
  @$pb.TagNumber(5)
  $core.String get path => $_getSZ(4);
  @$pb.TagNumber(5)
  set path($core.String value) => $_setString(4, value);
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
    final result = create();
    if (walletName != null) result.walletName = walletName;
    if (password != null) result.password = password;
    if (address != null) result.address = address;
    if (label != null) result.label = label;
    return result;
  }

  SetAddressLabelRequest._();

  factory SetAddressLabelRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SetAddressLabelRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SetAddressLabelRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'walletName')
    ..aOS(2, _omitFieldNames ? '' : 'password')
    ..aOS(3, _omitFieldNames ? '' : 'address')
    ..aOS(4, _omitFieldNames ? '' : 'label')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SetAddressLabelRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SetAddressLabelRequest copyWith(
          void Function(SetAddressLabelRequest) updates) =>
      super.copyWith((message) => updates(message as SetAddressLabelRequest))
          as SetAddressLabelRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SetAddressLabelRequest create() => SetAddressLabelRequest._();
  @$core.override
  SetAddressLabelRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SetAddressLabelRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SetAddressLabelRequest>(create);
  static SetAddressLabelRequest? _defaultInstance;

  /// The name of the wallet containing the address.
  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => $_clearField(1);

  /// Wallet password required for modification.
  @$pb.TagNumber(2)
  $core.String get password => $_getSZ(1);
  @$pb.TagNumber(2)
  set password($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasPassword() => $_has(1);
  @$pb.TagNumber(2)
  void clearPassword() => $_clearField(2);

  /// The address to label.
  @$pb.TagNumber(3)
  $core.String get address => $_getSZ(2);
  @$pb.TagNumber(3)
  set address($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasAddress() => $_has(2);
  @$pb.TagNumber(3)
  void clearAddress() => $_clearField(3);

  /// The new label for the address.
  @$pb.TagNumber(4)
  $core.String get label => $_getSZ(3);
  @$pb.TagNumber(4)
  set label($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasLabel() => $_has(3);
  @$pb.TagNumber(4)
  void clearLabel() => $_clearField(4);
}

/// Response message for updated address label.
class SetAddressLabelResponse extends $pb.GeneratedMessage {
  factory SetAddressLabelResponse({
    $core.String? walletName,
    $core.String? address,
    $core.String? label,
  }) {
    final result = create();
    if (walletName != null) result.walletName = walletName;
    if (address != null) result.address = address;
    if (label != null) result.label = label;
    return result;
  }

  SetAddressLabelResponse._();

  factory SetAddressLabelResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory SetAddressLabelResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'SetAddressLabelResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'walletName')
    ..aOS(2, _omitFieldNames ? '' : 'address')
    ..aOS(3, _omitFieldNames ? '' : 'label')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SetAddressLabelResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  SetAddressLabelResponse copyWith(
          void Function(SetAddressLabelResponse) updates) =>
      super.copyWith((message) => updates(message as SetAddressLabelResponse))
          as SetAddressLabelResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SetAddressLabelResponse create() => SetAddressLabelResponse._();
  @$core.override
  SetAddressLabelResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static SetAddressLabelResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<SetAddressLabelResponse>(create);
  static SetAddressLabelResponse? _defaultInstance;

  /// The name of the wallet where the address label was updated.
  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => $_clearField(1);

  /// The address where the label was updated.
  @$pb.TagNumber(2)
  $core.String get address => $_getSZ(1);
  @$pb.TagNumber(2)
  set address($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasAddress() => $_has(1);
  @$pb.TagNumber(2)
  void clearAddress() => $_clearField(2);

  /// The new label for the address.
  @$pb.TagNumber(3)
  $core.String get label => $_getSZ(2);
  @$pb.TagNumber(3)
  set label($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasLabel() => $_has(2);
  @$pb.TagNumber(3)
  void clearLabel() => $_clearField(3);
}

/// Request message for listing all wallets.
class ListWalletsRequest extends $pb.GeneratedMessage {
  factory ListWalletsRequest() => create();

  ListWalletsRequest._();

  factory ListWalletsRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListWalletsRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListWalletsRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListWalletsRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListWalletsRequest copyWith(void Function(ListWalletsRequest) updates) =>
      super.copyWith((message) => updates(message as ListWalletsRequest))
          as ListWalletsRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListWalletsRequest create() => ListWalletsRequest._();
  @$core.override
  ListWalletsRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListWalletsRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListWalletsRequest>(create);
  static ListWalletsRequest? _defaultInstance;
}

/// Response message contains wallet names.
class ListWalletsResponse extends $pb.GeneratedMessage {
  factory ListWalletsResponse({
    $core.Iterable<$core.String>? wallets,
  }) {
    final result = create();
    if (wallets != null) result.wallets.addAll(wallets);
    return result;
  }

  ListWalletsResponse._();

  factory ListWalletsResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListWalletsResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListWalletsResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..pPS(1, _omitFieldNames ? '' : 'wallets')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListWalletsResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListWalletsResponse copyWith(void Function(ListWalletsResponse) updates) =>
      super.copyWith((message) => updates(message as ListWalletsResponse))
          as ListWalletsResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListWalletsResponse create() => ListWalletsResponse._();
  @$core.override
  ListWalletsResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListWalletsResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListWalletsResponse>(create);
  static ListWalletsResponse? _defaultInstance;

  /// Array of wallet names.
  @$pb.TagNumber(1)
  $pb.PbList<$core.String> get wallets => $_getList(0);
}

/// Request message for getting wallet information.
class GetWalletInfoRequest extends $pb.GeneratedMessage {
  factory GetWalletInfoRequest({
    $core.String? walletName,
  }) {
    final result = create();
    if (walletName != null) result.walletName = walletName;
    return result;
  }

  GetWalletInfoRequest._();

  factory GetWalletInfoRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetWalletInfoRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetWalletInfoRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'walletName')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetWalletInfoRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetWalletInfoRequest copyWith(void Function(GetWalletInfoRequest) updates) =>
      super.copyWith((message) => updates(message as GetWalletInfoRequest))
          as GetWalletInfoRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetWalletInfoRequest create() => GetWalletInfoRequest._();
  @$core.override
  GetWalletInfoRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetWalletInfoRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetWalletInfoRequest>(create);
  static GetWalletInfoRequest? _defaultInstance;

  /// The name of the wallet to query.
  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String value) => $_setString(0, value);
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
    final result = create();
    if (walletName != null) result.walletName = walletName;
    if (version != null) result.version = version;
    if (network != null) result.network = network;
    if (encrypted != null) result.encrypted = encrypted;
    if (uuid != null) result.uuid = uuid;
    if (createdAt != null) result.createdAt = createdAt;
    if (defaultFee != null) result.defaultFee = defaultFee;
    return result;
  }

  GetWalletInfoResponse._();

  factory GetWalletInfoResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetWalletInfoResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetWalletInfoResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'walletName')
    ..aI(2, _omitFieldNames ? '' : 'version')
    ..aOS(3, _omitFieldNames ? '' : 'network')
    ..aOB(4, _omitFieldNames ? '' : 'encrypted')
    ..aOS(5, _omitFieldNames ? '' : 'uuid')
    ..aInt64(6, _omitFieldNames ? '' : 'createdAt')
    ..aInt64(7, _omitFieldNames ? '' : 'defaultFee')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetWalletInfoResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetWalletInfoResponse copyWith(
          void Function(GetWalletInfoResponse) updates) =>
      super.copyWith((message) => updates(message as GetWalletInfoResponse))
          as GetWalletInfoResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetWalletInfoResponse create() => GetWalletInfoResponse._();
  @$core.override
  GetWalletInfoResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetWalletInfoResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetWalletInfoResponse>(create);
  static GetWalletInfoResponse? _defaultInstance;

  /// The name of the wallet to query.
  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => $_clearField(1);

  /// The wallet format version.
  @$pb.TagNumber(2)
  $core.int get version => $_getIZ(1);
  @$pb.TagNumber(2)
  set version($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(2)
  $core.bool hasVersion() => $_has(1);
  @$pb.TagNumber(2)
  void clearVersion() => $_clearField(2);

  /// The network the wallet is connected to (e.g., mainnet, testnet).
  @$pb.TagNumber(3)
  $core.String get network => $_getSZ(2);
  @$pb.TagNumber(3)
  set network($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasNetwork() => $_has(2);
  @$pb.TagNumber(3)
  void clearNetwork() => $_clearField(3);

  /// Indicates if the wallet is encrypted.
  @$pb.TagNumber(4)
  $core.bool get encrypted => $_getBF(3);
  @$pb.TagNumber(4)
  set encrypted($core.bool value) => $_setBool(3, value);
  @$pb.TagNumber(4)
  $core.bool hasEncrypted() => $_has(3);
  @$pb.TagNumber(4)
  void clearEncrypted() => $_clearField(4);

  /// A unique identifier of the wallet.
  @$pb.TagNumber(5)
  $core.String get uuid => $_getSZ(4);
  @$pb.TagNumber(5)
  set uuid($core.String value) => $_setString(4, value);
  @$pb.TagNumber(5)
  $core.bool hasUuid() => $_has(4);
  @$pb.TagNumber(5)
  void clearUuid() => $_clearField(5);

  /// Unix timestamp of wallet creation.
  @$pb.TagNumber(6)
  $fixnum.Int64 get createdAt => $_getI64(5);
  @$pb.TagNumber(6)
  set createdAt($fixnum.Int64 value) => $_setInt64(5, value);
  @$pb.TagNumber(6)
  $core.bool hasCreatedAt() => $_has(5);
  @$pb.TagNumber(6)
  void clearCreatedAt() => $_clearField(6);

  /// The default fee of the wallet.
  @$pb.TagNumber(7)
  $fixnum.Int64 get defaultFee => $_getI64(6);
  @$pb.TagNumber(7)
  set defaultFee($fixnum.Int64 value) => $_setInt64(6, value);
  @$pb.TagNumber(7)
  $core.bool hasDefaultFee() => $_has(6);
  @$pb.TagNumber(7)
  void clearDefaultFee() => $_clearField(7);
}

/// Request message for listing wallet addresses.
class ListAddressesRequest extends $pb.GeneratedMessage {
  factory ListAddressesRequest({
    $core.String? walletName,
  }) {
    final result = create();
    if (walletName != null) result.walletName = walletName;
    return result;
  }

  ListAddressesRequest._();

  factory ListAddressesRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListAddressesRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListAddressesRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'walletName')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListAddressesRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListAddressesRequest copyWith(void Function(ListAddressesRequest) updates) =>
      super.copyWith((message) => updates(message as ListAddressesRequest))
          as ListAddressesRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListAddressesRequest create() => ListAddressesRequest._();
  @$core.override
  ListAddressesRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListAddressesRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListAddressesRequest>(create);
  static ListAddressesRequest? _defaultInstance;

  /// The name of the queried wallet.
  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasWalletName() => $_has(0);
  @$pb.TagNumber(1)
  void clearWalletName() => $_clearField(1);
}

/// Response message contains wallet addresses.
class ListAddressesResponse extends $pb.GeneratedMessage {
  factory ListAddressesResponse({
    $core.String? walletName,
    $core.Iterable<AddressInfo>? data,
  }) {
    final result = create();
    if (walletName != null) result.walletName = walletName;
    if (data != null) result.data.addAll(data);
    return result;
  }

  ListAddressesResponse._();

  factory ListAddressesResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListAddressesResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListAddressesResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'walletName')
    ..pPM<AddressInfo>(2, _omitFieldNames ? '' : 'data',
        subBuilder: AddressInfo.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListAddressesResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListAddressesResponse copyWith(
          void Function(ListAddressesResponse) updates) =>
      super.copyWith((message) => updates(message as ListAddressesResponse))
          as ListAddressesResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListAddressesResponse create() => ListAddressesResponse._();
  @$core.override
  ListAddressesResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListAddressesResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListAddressesResponse>(create);
  static ListAddressesResponse? _defaultInstance;

  /// The name of the queried wallet.
  @$pb.TagNumber(1)
  $core.String get walletName => $_getSZ(0);
  @$pb.TagNumber(1)
  set walletName($core.String value) => $_setString(0, value);
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
  final $pb.RpcClient _client;

  WalletApi(this._client);

  /// CreateWallet creates a new wallet with the specified parameters.
  $async.Future<CreateWalletResponse> createWallet(
          $pb.ClientContext? ctx, CreateWalletRequest request) =>
      _client.invoke<CreateWalletResponse>(
          ctx, 'Wallet', 'CreateWallet', request, CreateWalletResponse());

  /// RestoreWallet restores an existing wallet with the given mnemonic.
  $async.Future<RestoreWalletResponse> restoreWallet(
          $pb.ClientContext? ctx, RestoreWalletRequest request) =>
      _client.invoke<RestoreWalletResponse>(
          ctx, 'Wallet', 'RestoreWallet', request, RestoreWalletResponse());

  /// LoadWallet loads an existing wallet with the given name.
  $async.Future<LoadWalletResponse> loadWallet(
          $pb.ClientContext? ctx, LoadWalletRequest request) =>
      _client.invoke<LoadWalletResponse>(
          ctx, 'Wallet', 'LoadWallet', request, LoadWalletResponse());

  /// UnloadWallet unloads a currently loaded wallet with the specified name.
  $async.Future<UnloadWalletResponse> unloadWallet(
          $pb.ClientContext? ctx, UnloadWalletRequest request) =>
      _client.invoke<UnloadWalletResponse>(
          ctx, 'Wallet', 'UnloadWallet', request, UnloadWalletResponse());

  /// GetTotalBalance returns the total available balance of the wallet.
  $async.Future<GetTotalBalanceResponse> getTotalBalance(
          $pb.ClientContext? ctx, GetTotalBalanceRequest request) =>
      _client.invoke<GetTotalBalanceResponse>(
          ctx, 'Wallet', 'GetTotalBalance', request, GetTotalBalanceResponse());

  /// SignRawTransaction signs a raw transaction for a specified wallet.
  $async.Future<SignRawTransactionResponse> signRawTransaction(
          $pb.ClientContext? ctx, SignRawTransactionRequest request) =>
      _client.invoke<SignRawTransactionResponse>(ctx, 'Wallet',
          'SignRawTransaction', request, SignRawTransactionResponse());

  /// GetValidatorAddress retrieves the validator address associated with a public key.
  /// Deprecated: Will move into utils.
  $async.Future<GetValidatorAddressResponse> getValidatorAddress(
          $pb.ClientContext? ctx, GetValidatorAddressRequest request) =>
      _client.invoke<GetValidatorAddressResponse>(ctx, 'Wallet',
          'GetValidatorAddress', request, GetValidatorAddressResponse());

  /// GetNewAddress generates a new address for the specified wallet.
  $async.Future<GetNewAddressResponse> getNewAddress(
          $pb.ClientContext? ctx, GetNewAddressRequest request) =>
      _client.invoke<GetNewAddressResponse>(
          ctx, 'Wallet', 'GetNewAddress', request, GetNewAddressResponse());

  /// GetAddressHistory retrieves the transaction history of an address.
  $async.Future<GetAddressHistoryResponse> getAddressHistory(
          $pb.ClientContext? ctx, GetAddressHistoryRequest request) =>
      _client.invoke<GetAddressHistoryResponse>(ctx, 'Wallet',
          'GetAddressHistory', request, GetAddressHistoryResponse());

  /// SignMessage signs an arbitrary message using a wallet's private key.
  $async.Future<SignMessageResponse> signMessage(
          $pb.ClientContext? ctx, SignMessageRequest request) =>
      _client.invoke<SignMessageResponse>(
          ctx, 'Wallet', 'SignMessage', request, SignMessageResponse());

  /// GetTotalStake returns the total stake amount in the wallet.
  $async.Future<GetTotalStakeResponse> getTotalStake(
          $pb.ClientContext? ctx, GetTotalStakeRequest request) =>
      _client.invoke<GetTotalStakeResponse>(
          ctx, 'Wallet', 'GetTotalStake', request, GetTotalStakeResponse());

  /// GetAddressInfo returns detailed information about a specific address.
  $async.Future<GetAddressInfoResponse> getAddressInfo(
          $pb.ClientContext? ctx, GetAddressInfoRequest request) =>
      _client.invoke<GetAddressInfoResponse>(
          ctx, 'Wallet', 'GetAddressInfo', request, GetAddressInfoResponse());

  /// SetAddressLabel sets or updates the label for a given address.
  $async.Future<SetAddressLabelResponse> setAddressLabel(
          $pb.ClientContext? ctx, SetAddressLabelRequest request) =>
      _client.invoke<SetAddressLabelResponse>(
          ctx, 'Wallet', 'SetAddressLabel', request, SetAddressLabelResponse());

  /// ListWallets returns a list of all available wallets.
  $async.Future<ListWalletsResponse> listWallets(
          $pb.ClientContext? ctx, ListWalletsRequest request) =>
      _client.invoke<ListWalletsResponse>(
          ctx, 'Wallet', 'ListWallets', request, ListWalletsResponse());

  /// GetWalletInfo returns detailed information about a specific wallet.
  $async.Future<GetWalletInfoResponse> getWalletInfo(
          $pb.ClientContext? ctx, GetWalletInfoRequest request) =>
      _client.invoke<GetWalletInfoResponse>(
          ctx, 'Wallet', 'GetWalletInfo', request, GetWalletInfoResponse());

  /// ListAddresses returns all addresses in the specified wallet.
  $async.Future<ListAddressesResponse> listAddresses(
          $pb.ClientContext? ctx, ListAddressesRequest request) =>
      _client.invoke<ListAddressesResponse>(
          ctx, 'Wallet', 'ListAddresses', request, ListAddressesResponse());
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
