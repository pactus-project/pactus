// This is a generated file - do not edit.
//
// Generated from blockchain.proto.

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

import 'blockchain.pbenum.dart';
import 'transaction.pb.dart' as $0;

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

export 'blockchain.pbenum.dart';

/// Request message for retrieving account information.
class GetAccountRequest extends $pb.GeneratedMessage {
  factory GetAccountRequest({
    $core.String? address,
  }) {
    final result = create();
    if (address != null) result.address = address;
    return result;
  }

  GetAccountRequest._();

  factory GetAccountRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetAccountRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetAccountRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'address')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetAccountRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetAccountRequest copyWith(void Function(GetAccountRequest) updates) =>
      super.copyWith((message) => updates(message as GetAccountRequest))
          as GetAccountRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetAccountRequest create() => GetAccountRequest._();
  @$core.override
  GetAccountRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetAccountRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetAccountRequest>(create);
  static GetAccountRequest? _defaultInstance;

  /// The address of the account to retrieve information for.
  @$pb.TagNumber(1)
  $core.String get address => $_getSZ(0);
  @$pb.TagNumber(1)
  set address($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasAddress() => $_has(0);
  @$pb.TagNumber(1)
  void clearAddress() => $_clearField(1);
}

/// Response message contains account information.
class GetAccountResponse extends $pb.GeneratedMessage {
  factory GetAccountResponse({
    AccountInfo? account,
  }) {
    final result = create();
    if (account != null) result.account = account;
    return result;
  }

  GetAccountResponse._();

  factory GetAccountResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetAccountResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetAccountResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOM<AccountInfo>(1, _omitFieldNames ? '' : 'account',
        subBuilder: AccountInfo.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetAccountResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetAccountResponse copyWith(void Function(GetAccountResponse) updates) =>
      super.copyWith((message) => updates(message as GetAccountResponse))
          as GetAccountResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetAccountResponse create() => GetAccountResponse._();
  @$core.override
  GetAccountResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetAccountResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetAccountResponse>(create);
  static GetAccountResponse? _defaultInstance;

  /// Detailed information about the account.
  @$pb.TagNumber(1)
  AccountInfo get account => $_getN(0);
  @$pb.TagNumber(1)
  set account(AccountInfo value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasAccount() => $_has(0);
  @$pb.TagNumber(1)
  void clearAccount() => $_clearField(1);
  @$pb.TagNumber(1)
  AccountInfo ensureAccount() => $_ensure(0);
}

/// Request message for retrieving validator addresses.
class GetValidatorAddressesRequest extends $pb.GeneratedMessage {
  factory GetValidatorAddressesRequest() => create();

  GetValidatorAddressesRequest._();

  factory GetValidatorAddressesRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetValidatorAddressesRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetValidatorAddressesRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetValidatorAddressesRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetValidatorAddressesRequest copyWith(
          void Function(GetValidatorAddressesRequest) updates) =>
      super.copyWith(
              (message) => updates(message as GetValidatorAddressesRequest))
          as GetValidatorAddressesRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetValidatorAddressesRequest create() =>
      GetValidatorAddressesRequest._();
  @$core.override
  GetValidatorAddressesRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetValidatorAddressesRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetValidatorAddressesRequest>(create);
  static GetValidatorAddressesRequest? _defaultInstance;
}

/// Response message contains list of validator addresses.
class GetValidatorAddressesResponse extends $pb.GeneratedMessage {
  factory GetValidatorAddressesResponse({
    $core.Iterable<$core.String>? addresses,
  }) {
    final result = create();
    if (addresses != null) result.addresses.addAll(addresses);
    return result;
  }

  GetValidatorAddressesResponse._();

  factory GetValidatorAddressesResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetValidatorAddressesResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetValidatorAddressesResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..pPS(1, _omitFieldNames ? '' : 'addresses')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetValidatorAddressesResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetValidatorAddressesResponse copyWith(
          void Function(GetValidatorAddressesResponse) updates) =>
      super.copyWith(
              (message) => updates(message as GetValidatorAddressesResponse))
          as GetValidatorAddressesResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetValidatorAddressesResponse create() =>
      GetValidatorAddressesResponse._();
  @$core.override
  GetValidatorAddressesResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetValidatorAddressesResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetValidatorAddressesResponse>(create);
  static GetValidatorAddressesResponse? _defaultInstance;

  /// List of validator addresses.
  @$pb.TagNumber(1)
  $pb.PbList<$core.String> get addresses => $_getList(0);
}

/// Request message for retrieving validator information by address.
class GetValidatorRequest extends $pb.GeneratedMessage {
  factory GetValidatorRequest({
    $core.String? address,
  }) {
    final result = create();
    if (address != null) result.address = address;
    return result;
  }

  GetValidatorRequest._();

  factory GetValidatorRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetValidatorRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetValidatorRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'address')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetValidatorRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetValidatorRequest copyWith(void Function(GetValidatorRequest) updates) =>
      super.copyWith((message) => updates(message as GetValidatorRequest))
          as GetValidatorRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetValidatorRequest create() => GetValidatorRequest._();
  @$core.override
  GetValidatorRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetValidatorRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetValidatorRequest>(create);
  static GetValidatorRequest? _defaultInstance;

  /// The address of the validator to retrieve information for.
  @$pb.TagNumber(1)
  $core.String get address => $_getSZ(0);
  @$pb.TagNumber(1)
  set address($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasAddress() => $_has(0);
  @$pb.TagNumber(1)
  void clearAddress() => $_clearField(1);
}

/// Request message for retrieving validator information by number.
class GetValidatorByNumberRequest extends $pb.GeneratedMessage {
  factory GetValidatorByNumberRequest({
    $core.int? number,
  }) {
    final result = create();
    if (number != null) result.number = number;
    return result;
  }

  GetValidatorByNumberRequest._();

  factory GetValidatorByNumberRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetValidatorByNumberRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetValidatorByNumberRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aI(1, _omitFieldNames ? '' : 'number')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetValidatorByNumberRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetValidatorByNumberRequest copyWith(
          void Function(GetValidatorByNumberRequest) updates) =>
      super.copyWith(
              (message) => updates(message as GetValidatorByNumberRequest))
          as GetValidatorByNumberRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetValidatorByNumberRequest create() =>
      GetValidatorByNumberRequest._();
  @$core.override
  GetValidatorByNumberRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetValidatorByNumberRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetValidatorByNumberRequest>(create);
  static GetValidatorByNumberRequest? _defaultInstance;

  /// The unique number of the validator to retrieve information for.
  @$pb.TagNumber(1)
  $core.int get number => $_getIZ(0);
  @$pb.TagNumber(1)
  set number($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(1)
  $core.bool hasNumber() => $_has(0);
  @$pb.TagNumber(1)
  void clearNumber() => $_clearField(1);
}

/// Response message contains validator information.
class GetValidatorResponse extends $pb.GeneratedMessage {
  factory GetValidatorResponse({
    ValidatorInfo? validator,
  }) {
    final result = create();
    if (validator != null) result.validator = validator;
    return result;
  }

  GetValidatorResponse._();

  factory GetValidatorResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetValidatorResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetValidatorResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOM<ValidatorInfo>(1, _omitFieldNames ? '' : 'validator',
        subBuilder: ValidatorInfo.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetValidatorResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetValidatorResponse copyWith(void Function(GetValidatorResponse) updates) =>
      super.copyWith((message) => updates(message as GetValidatorResponse))
          as GetValidatorResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetValidatorResponse create() => GetValidatorResponse._();
  @$core.override
  GetValidatorResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetValidatorResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetValidatorResponse>(create);
  static GetValidatorResponse? _defaultInstance;

  /// Detailed information about the validator.
  @$pb.TagNumber(1)
  ValidatorInfo get validator => $_getN(0);
  @$pb.TagNumber(1)
  set validator(ValidatorInfo value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasValidator() => $_has(0);
  @$pb.TagNumber(1)
  void clearValidator() => $_clearField(1);
  @$pb.TagNumber(1)
  ValidatorInfo ensureValidator() => $_ensure(0);
}

/// Request message for retrieving public key by address.
class GetPublicKeyRequest extends $pb.GeneratedMessage {
  factory GetPublicKeyRequest({
    $core.String? address,
  }) {
    final result = create();
    if (address != null) result.address = address;
    return result;
  }

  GetPublicKeyRequest._();

  factory GetPublicKeyRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetPublicKeyRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetPublicKeyRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'address')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetPublicKeyRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetPublicKeyRequest copyWith(void Function(GetPublicKeyRequest) updates) =>
      super.copyWith((message) => updates(message as GetPublicKeyRequest))
          as GetPublicKeyRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetPublicKeyRequest create() => GetPublicKeyRequest._();
  @$core.override
  GetPublicKeyRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetPublicKeyRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetPublicKeyRequest>(create);
  static GetPublicKeyRequest? _defaultInstance;

  /// The address for which to retrieve the public key.
  @$pb.TagNumber(1)
  $core.String get address => $_getSZ(0);
  @$pb.TagNumber(1)
  set address($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasAddress() => $_has(0);
  @$pb.TagNumber(1)
  void clearAddress() => $_clearField(1);
}

/// Response message contains public key information.
class GetPublicKeyResponse extends $pb.GeneratedMessage {
  factory GetPublicKeyResponse({
    $core.String? publicKey,
  }) {
    final result = create();
    if (publicKey != null) result.publicKey = publicKey;
    return result;
  }

  GetPublicKeyResponse._();

  factory GetPublicKeyResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetPublicKeyResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetPublicKeyResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'publicKey')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetPublicKeyResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetPublicKeyResponse copyWith(void Function(GetPublicKeyResponse) updates) =>
      super.copyWith((message) => updates(message as GetPublicKeyResponse))
          as GetPublicKeyResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetPublicKeyResponse create() => GetPublicKeyResponse._();
  @$core.override
  GetPublicKeyResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetPublicKeyResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetPublicKeyResponse>(create);
  static GetPublicKeyResponse? _defaultInstance;

  /// The public key associated with the provided address.
  @$pb.TagNumber(1)
  $core.String get publicKey => $_getSZ(0);
  @$pb.TagNumber(1)
  set publicKey($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasPublicKey() => $_has(0);
  @$pb.TagNumber(1)
  void clearPublicKey() => $_clearField(1);
}

/// Request message for retrieving block information based on height and verbosity level.
class GetBlockRequest extends $pb.GeneratedMessage {
  factory GetBlockRequest({
    $core.int? height,
    BlockVerbosity? verbosity,
  }) {
    final result = create();
    if (height != null) result.height = height;
    if (verbosity != null) result.verbosity = verbosity;
    return result;
  }

  GetBlockRequest._();

  factory GetBlockRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetBlockRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetBlockRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aI(1, _omitFieldNames ? '' : 'height', fieldType: $pb.PbFieldType.OU3)
    ..aE<BlockVerbosity>(2, _omitFieldNames ? '' : 'verbosity',
        enumValues: BlockVerbosity.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetBlockRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetBlockRequest copyWith(void Function(GetBlockRequest) updates) =>
      super.copyWith((message) => updates(message as GetBlockRequest))
          as GetBlockRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetBlockRequest create() => GetBlockRequest._();
  @$core.override
  GetBlockRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetBlockRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetBlockRequest>(create);
  static GetBlockRequest? _defaultInstance;

  /// The height of the block to retrieve.
  @$pb.TagNumber(1)
  $core.int get height => $_getIZ(0);
  @$pb.TagNumber(1)
  set height($core.int value) => $_setUnsignedInt32(0, value);
  @$pb.TagNumber(1)
  $core.bool hasHeight() => $_has(0);
  @$pb.TagNumber(1)
  void clearHeight() => $_clearField(1);

  /// The verbosity level for block information.
  @$pb.TagNumber(2)
  BlockVerbosity get verbosity => $_getN(1);
  @$pb.TagNumber(2)
  set verbosity(BlockVerbosity value) => $_setField(2, value);
  @$pb.TagNumber(2)
  $core.bool hasVerbosity() => $_has(1);
  @$pb.TagNumber(2)
  void clearVerbosity() => $_clearField(2);
}

/// Response message contains block information.
class GetBlockResponse extends $pb.GeneratedMessage {
  factory GetBlockResponse({
    $core.int? height,
    $core.String? hash,
    $core.String? data,
    $core.int? blockTime,
    BlockHeaderInfo? header,
    CertificateInfo? prevCert,
    $core.Iterable<$0.TransactionInfo>? txs,
  }) {
    final result = create();
    if (height != null) result.height = height;
    if (hash != null) result.hash = hash;
    if (data != null) result.data = data;
    if (blockTime != null) result.blockTime = blockTime;
    if (header != null) result.header = header;
    if (prevCert != null) result.prevCert = prevCert;
    if (txs != null) result.txs.addAll(txs);
    return result;
  }

  GetBlockResponse._();

  factory GetBlockResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetBlockResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetBlockResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aI(1, _omitFieldNames ? '' : 'height', fieldType: $pb.PbFieldType.OU3)
    ..aOS(2, _omitFieldNames ? '' : 'hash')
    ..aOS(3, _omitFieldNames ? '' : 'data')
    ..aI(4, _omitFieldNames ? '' : 'blockTime', fieldType: $pb.PbFieldType.OU3)
    ..aOM<BlockHeaderInfo>(5, _omitFieldNames ? '' : 'header',
        subBuilder: BlockHeaderInfo.create)
    ..aOM<CertificateInfo>(6, _omitFieldNames ? '' : 'prevCert',
        subBuilder: CertificateInfo.create)
    ..pPM<$0.TransactionInfo>(7, _omitFieldNames ? '' : 'txs',
        subBuilder: $0.TransactionInfo.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetBlockResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetBlockResponse copyWith(void Function(GetBlockResponse) updates) =>
      super.copyWith((message) => updates(message as GetBlockResponse))
          as GetBlockResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetBlockResponse create() => GetBlockResponse._();
  @$core.override
  GetBlockResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetBlockResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetBlockResponse>(create);
  static GetBlockResponse? _defaultInstance;

  /// The height of the block.
  @$pb.TagNumber(1)
  $core.int get height => $_getIZ(0);
  @$pb.TagNumber(1)
  set height($core.int value) => $_setUnsignedInt32(0, value);
  @$pb.TagNumber(1)
  $core.bool hasHeight() => $_has(0);
  @$pb.TagNumber(1)
  void clearHeight() => $_clearField(1);

  /// The hash of the block.
  @$pb.TagNumber(2)
  $core.String get hash => $_getSZ(1);
  @$pb.TagNumber(2)
  set hash($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasHash() => $_has(1);
  @$pb.TagNumber(2)
  void clearHash() => $_clearField(2);

  /// Block data, available only if verbosity level is set to BLOCK_VERBOSITY_DATA.
  @$pb.TagNumber(3)
  $core.String get data => $_getSZ(2);
  @$pb.TagNumber(3)
  set data($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasData() => $_has(2);
  @$pb.TagNumber(3)
  void clearData() => $_clearField(3);

  /// The timestamp of the block.
  @$pb.TagNumber(4)
  $core.int get blockTime => $_getIZ(3);
  @$pb.TagNumber(4)
  set blockTime($core.int value) => $_setUnsignedInt32(3, value);
  @$pb.TagNumber(4)
  $core.bool hasBlockTime() => $_has(3);
  @$pb.TagNumber(4)
  void clearBlockTime() => $_clearField(4);

  /// Header information of the block.
  @$pb.TagNumber(5)
  BlockHeaderInfo get header => $_getN(4);
  @$pb.TagNumber(5)
  set header(BlockHeaderInfo value) => $_setField(5, value);
  @$pb.TagNumber(5)
  $core.bool hasHeader() => $_has(4);
  @$pb.TagNumber(5)
  void clearHeader() => $_clearField(5);
  @$pb.TagNumber(5)
  BlockHeaderInfo ensureHeader() => $_ensure(4);

  /// Certificate information of the previous block.
  @$pb.TagNumber(6)
  CertificateInfo get prevCert => $_getN(5);
  @$pb.TagNumber(6)
  set prevCert(CertificateInfo value) => $_setField(6, value);
  @$pb.TagNumber(6)
  $core.bool hasPrevCert() => $_has(5);
  @$pb.TagNumber(6)
  void clearPrevCert() => $_clearField(6);
  @$pb.TagNumber(6)
  CertificateInfo ensurePrevCert() => $_ensure(5);

  /// List of transactions in the block, available when verbosity level is set to
  /// BLOCK_VERBOSITY_TRANSACTIONS.
  @$pb.TagNumber(7)
  $pb.PbList<$0.TransactionInfo> get txs => $_getList(6);
}

/// Request message for retrieving block hash by height.
class GetBlockHashRequest extends $pb.GeneratedMessage {
  factory GetBlockHashRequest({
    $core.int? height,
  }) {
    final result = create();
    if (height != null) result.height = height;
    return result;
  }

  GetBlockHashRequest._();

  factory GetBlockHashRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetBlockHashRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetBlockHashRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aI(1, _omitFieldNames ? '' : 'height', fieldType: $pb.PbFieldType.OU3)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetBlockHashRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetBlockHashRequest copyWith(void Function(GetBlockHashRequest) updates) =>
      super.copyWith((message) => updates(message as GetBlockHashRequest))
          as GetBlockHashRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetBlockHashRequest create() => GetBlockHashRequest._();
  @$core.override
  GetBlockHashRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetBlockHashRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetBlockHashRequest>(create);
  static GetBlockHashRequest? _defaultInstance;

  /// The height of the block to retrieve the hash for.
  @$pb.TagNumber(1)
  $core.int get height => $_getIZ(0);
  @$pb.TagNumber(1)
  set height($core.int value) => $_setUnsignedInt32(0, value);
  @$pb.TagNumber(1)
  $core.bool hasHeight() => $_has(0);
  @$pb.TagNumber(1)
  void clearHeight() => $_clearField(1);
}

/// Response message contains block hash.
class GetBlockHashResponse extends $pb.GeneratedMessage {
  factory GetBlockHashResponse({
    $core.String? hash,
  }) {
    final result = create();
    if (hash != null) result.hash = hash;
    return result;
  }

  GetBlockHashResponse._();

  factory GetBlockHashResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetBlockHashResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetBlockHashResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'hash')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetBlockHashResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetBlockHashResponse copyWith(void Function(GetBlockHashResponse) updates) =>
      super.copyWith((message) => updates(message as GetBlockHashResponse))
          as GetBlockHashResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetBlockHashResponse create() => GetBlockHashResponse._();
  @$core.override
  GetBlockHashResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetBlockHashResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetBlockHashResponse>(create);
  static GetBlockHashResponse? _defaultInstance;

  /// The hash of the block.
  @$pb.TagNumber(1)
  $core.String get hash => $_getSZ(0);
  @$pb.TagNumber(1)
  set hash($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasHash() => $_has(0);
  @$pb.TagNumber(1)
  void clearHash() => $_clearField(1);
}

/// Request message for retrieving block height by hash.
class GetBlockHeightRequest extends $pb.GeneratedMessage {
  factory GetBlockHeightRequest({
    $core.String? hash,
  }) {
    final result = create();
    if (hash != null) result.hash = hash;
    return result;
  }

  GetBlockHeightRequest._();

  factory GetBlockHeightRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetBlockHeightRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetBlockHeightRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'hash')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetBlockHeightRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetBlockHeightRequest copyWith(
          void Function(GetBlockHeightRequest) updates) =>
      super.copyWith((message) => updates(message as GetBlockHeightRequest))
          as GetBlockHeightRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetBlockHeightRequest create() => GetBlockHeightRequest._();
  @$core.override
  GetBlockHeightRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetBlockHeightRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetBlockHeightRequest>(create);
  static GetBlockHeightRequest? _defaultInstance;

  /// The hash of the block to retrieve the height for.
  @$pb.TagNumber(1)
  $core.String get hash => $_getSZ(0);
  @$pb.TagNumber(1)
  set hash($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasHash() => $_has(0);
  @$pb.TagNumber(1)
  void clearHash() => $_clearField(1);
}

/// Response message contains block height.
class GetBlockHeightResponse extends $pb.GeneratedMessage {
  factory GetBlockHeightResponse({
    $core.int? height,
  }) {
    final result = create();
    if (height != null) result.height = height;
    return result;
  }

  GetBlockHeightResponse._();

  factory GetBlockHeightResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetBlockHeightResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetBlockHeightResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aI(1, _omitFieldNames ? '' : 'height', fieldType: $pb.PbFieldType.OU3)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetBlockHeightResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetBlockHeightResponse copyWith(
          void Function(GetBlockHeightResponse) updates) =>
      super.copyWith((message) => updates(message as GetBlockHeightResponse))
          as GetBlockHeightResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetBlockHeightResponse create() => GetBlockHeightResponse._();
  @$core.override
  GetBlockHeightResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetBlockHeightResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetBlockHeightResponse>(create);
  static GetBlockHeightResponse? _defaultInstance;

  /// The height of the block.
  @$pb.TagNumber(1)
  $core.int get height => $_getIZ(0);
  @$pb.TagNumber(1)
  set height($core.int value) => $_setUnsignedInt32(0, value);
  @$pb.TagNumber(1)
  $core.bool hasHeight() => $_has(0);
  @$pb.TagNumber(1)
  void clearHeight() => $_clearField(1);
}

/// Request message for retrieving blockchain information.
class GetBlockchainInfoRequest extends $pb.GeneratedMessage {
  factory GetBlockchainInfoRequest() => create();

  GetBlockchainInfoRequest._();

  factory GetBlockchainInfoRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetBlockchainInfoRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetBlockchainInfoRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetBlockchainInfoRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetBlockchainInfoRequest copyWith(
          void Function(GetBlockchainInfoRequest) updates) =>
      super.copyWith((message) => updates(message as GetBlockchainInfoRequest))
          as GetBlockchainInfoRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetBlockchainInfoRequest create() => GetBlockchainInfoRequest._();
  @$core.override
  GetBlockchainInfoRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetBlockchainInfoRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetBlockchainInfoRequest>(create);
  static GetBlockchainInfoRequest? _defaultInstance;
}

/// Response message contains general blockchain information.
class GetBlockchainInfoResponse extends $pb.GeneratedMessage {
  factory GetBlockchainInfoResponse({
    $core.int? lastBlockHeight,
    $core.String? lastBlockHash,
    $core.int? totalAccounts,
    $core.int? totalValidators,
    $fixnum.Int64? totalPower,
    $fixnum.Int64? committeePower,
    $core.bool? isPruned,
    $core.int? pruningHeight,
    $fixnum.Int64? lastBlockTime,
    $core.int? activeValidators,
    $core.bool? inCommittee,
  }) {
    final result = create();
    if (lastBlockHeight != null) result.lastBlockHeight = lastBlockHeight;
    if (lastBlockHash != null) result.lastBlockHash = lastBlockHash;
    if (totalAccounts != null) result.totalAccounts = totalAccounts;
    if (totalValidators != null) result.totalValidators = totalValidators;
    if (totalPower != null) result.totalPower = totalPower;
    if (committeePower != null) result.committeePower = committeePower;
    if (isPruned != null) result.isPruned = isPruned;
    if (pruningHeight != null) result.pruningHeight = pruningHeight;
    if (lastBlockTime != null) result.lastBlockTime = lastBlockTime;
    if (activeValidators != null) result.activeValidators = activeValidators;
    if (inCommittee != null) result.inCommittee = inCommittee;
    return result;
  }

  GetBlockchainInfoResponse._();

  factory GetBlockchainInfoResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetBlockchainInfoResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetBlockchainInfoResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aI(1, _omitFieldNames ? '' : 'lastBlockHeight',
        fieldType: $pb.PbFieldType.OU3)
    ..aOS(2, _omitFieldNames ? '' : 'lastBlockHash')
    ..aI(3, _omitFieldNames ? '' : 'totalAccounts')
    ..aI(4, _omitFieldNames ? '' : 'totalValidators')
    ..aInt64(5, _omitFieldNames ? '' : 'totalPower')
    ..aInt64(6, _omitFieldNames ? '' : 'committeePower')
    ..aOB(8, _omitFieldNames ? '' : 'isPruned')
    ..aI(9, _omitFieldNames ? '' : 'pruningHeight',
        fieldType: $pb.PbFieldType.OU3)
    ..aInt64(10, _omitFieldNames ? '' : 'lastBlockTime')
    ..aI(12, _omitFieldNames ? '' : 'activeValidators')
    ..aOB(13, _omitFieldNames ? '' : 'inCommittee')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetBlockchainInfoResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetBlockchainInfoResponse copyWith(
          void Function(GetBlockchainInfoResponse) updates) =>
      super.copyWith((message) => updates(message as GetBlockchainInfoResponse))
          as GetBlockchainInfoResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetBlockchainInfoResponse create() => GetBlockchainInfoResponse._();
  @$core.override
  GetBlockchainInfoResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetBlockchainInfoResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetBlockchainInfoResponse>(create);
  static GetBlockchainInfoResponse? _defaultInstance;

  /// The height of the last block in the blockchain.
  @$pb.TagNumber(1)
  $core.int get lastBlockHeight => $_getIZ(0);
  @$pb.TagNumber(1)
  set lastBlockHeight($core.int value) => $_setUnsignedInt32(0, value);
  @$pb.TagNumber(1)
  $core.bool hasLastBlockHeight() => $_has(0);
  @$pb.TagNumber(1)
  void clearLastBlockHeight() => $_clearField(1);

  /// The hash of the last block in the blockchain.
  @$pb.TagNumber(2)
  $core.String get lastBlockHash => $_getSZ(1);
  @$pb.TagNumber(2)
  set lastBlockHash($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasLastBlockHash() => $_has(1);
  @$pb.TagNumber(2)
  void clearLastBlockHash() => $_clearField(2);

  /// The total number of accounts in the blockchain.
  @$pb.TagNumber(3)
  $core.int get totalAccounts => $_getIZ(2);
  @$pb.TagNumber(3)
  set totalAccounts($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(3)
  $core.bool hasTotalAccounts() => $_has(2);
  @$pb.TagNumber(3)
  void clearTotalAccounts() => $_clearField(3);

  /// The total number of validators in the blockchain.
  @$pb.TagNumber(4)
  $core.int get totalValidators => $_getIZ(3);
  @$pb.TagNumber(4)
  set totalValidators($core.int value) => $_setSignedInt32(3, value);
  @$pb.TagNumber(4)
  $core.bool hasTotalValidators() => $_has(3);
  @$pb.TagNumber(4)
  void clearTotalValidators() => $_clearField(4);

  /// The total power of the blockchain.
  @$pb.TagNumber(5)
  $fixnum.Int64 get totalPower => $_getI64(4);
  @$pb.TagNumber(5)
  set totalPower($fixnum.Int64 value) => $_setInt64(4, value);
  @$pb.TagNumber(5)
  $core.bool hasTotalPower() => $_has(4);
  @$pb.TagNumber(5)
  void clearTotalPower() => $_clearField(5);

  /// The power of the committee.
  @$pb.TagNumber(6)
  $fixnum.Int64 get committeePower => $_getI64(5);
  @$pb.TagNumber(6)
  set committeePower($fixnum.Int64 value) => $_setInt64(5, value);
  @$pb.TagNumber(6)
  $core.bool hasCommitteePower() => $_has(5);
  @$pb.TagNumber(6)
  void clearCommitteePower() => $_clearField(6);

  /// If the blocks are subject to pruning.
  @$pb.TagNumber(8)
  $core.bool get isPruned => $_getBF(6);
  @$pb.TagNumber(8)
  set isPruned($core.bool value) => $_setBool(6, value);
  @$pb.TagNumber(8)
  $core.bool hasIsPruned() => $_has(6);
  @$pb.TagNumber(8)
  void clearIsPruned() => $_clearField(8);

  /// Lowest-height block stored (only present if pruning is enabled)
  @$pb.TagNumber(9)
  $core.int get pruningHeight => $_getIZ(7);
  @$pb.TagNumber(9)
  set pruningHeight($core.int value) => $_setUnsignedInt32(7, value);
  @$pb.TagNumber(9)
  $core.bool hasPruningHeight() => $_has(7);
  @$pb.TagNumber(9)
  void clearPruningHeight() => $_clearField(9);

  /// The timestamp of the last block in Unix format.
  @$pb.TagNumber(10)
  $fixnum.Int64 get lastBlockTime => $_getI64(8);
  @$pb.TagNumber(10)
  set lastBlockTime($fixnum.Int64 value) => $_setInt64(8, value);
  @$pb.TagNumber(10)
  $core.bool hasLastBlockTime() => $_has(8);
  @$pb.TagNumber(10)
  void clearLastBlockTime() => $_clearField(10);

  /// The number of active (not unbonded) validators in the blockchain.
  @$pb.TagNumber(12)
  $core.int get activeValidators => $_getIZ(9);
  @$pb.TagNumber(12)
  set activeValidators($core.int value) => $_setSignedInt32(9, value);
  @$pb.TagNumber(12)
  $core.bool hasActiveValidators() => $_has(9);
  @$pb.TagNumber(12)
  void clearActiveValidators() => $_clearField(12);

  /// Indicates whether this node participates in consensus: true if at least one
  /// of its running validators is a member of the current committee.
  @$pb.TagNumber(13)
  $core.bool get inCommittee => $_getBF(10);
  @$pb.TagNumber(13)
  set inCommittee($core.bool value) => $_setBool(10, value);
  @$pb.TagNumber(13)
  $core.bool hasInCommittee() => $_has(10);
  @$pb.TagNumber(13)
  void clearInCommittee() => $_clearField(13);
}

/// Request message for retrieving committee information.
class GetCommitteeInfoRequest extends $pb.GeneratedMessage {
  factory GetCommitteeInfoRequest() => create();

  GetCommitteeInfoRequest._();

  factory GetCommitteeInfoRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetCommitteeInfoRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetCommitteeInfoRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetCommitteeInfoRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetCommitteeInfoRequest copyWith(
          void Function(GetCommitteeInfoRequest) updates) =>
      super.copyWith((message) => updates(message as GetCommitteeInfoRequest))
          as GetCommitteeInfoRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetCommitteeInfoRequest create() => GetCommitteeInfoRequest._();
  @$core.override
  GetCommitteeInfoRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetCommitteeInfoRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetCommitteeInfoRequest>(create);
  static GetCommitteeInfoRequest? _defaultInstance;
}

/// Response message contains committee information.
class GetCommitteeInfoResponse extends $pb.GeneratedMessage {
  factory GetCommitteeInfoResponse({
    $fixnum.Int64? committeePower,
    $core.Iterable<ValidatorInfo>? validators,
    $core.Iterable<$core.MapEntry<$core.int, $core.double>>? protocolVersions,
  }) {
    final result = create();
    if (committeePower != null) result.committeePower = committeePower;
    if (validators != null) result.validators.addAll(validators);
    if (protocolVersions != null)
      result.protocolVersions.addEntries(protocolVersions);
    return result;
  }

  GetCommitteeInfoResponse._();

  factory GetCommitteeInfoResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetCommitteeInfoResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetCommitteeInfoResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aInt64(1, _omitFieldNames ? '' : 'committeePower')
    ..pPM<ValidatorInfo>(2, _omitFieldNames ? '' : 'validators',
        subBuilder: ValidatorInfo.create)
    ..m<$core.int, $core.double>(3, _omitFieldNames ? '' : 'protocolVersions',
        entryClassName: 'GetCommitteeInfoResponse.ProtocolVersionsEntry',
        keyFieldType: $pb.PbFieldType.O3,
        valueFieldType: $pb.PbFieldType.OD,
        packageName: const $pb.PackageName('pactus'))
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetCommitteeInfoResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetCommitteeInfoResponse copyWith(
          void Function(GetCommitteeInfoResponse) updates) =>
      super.copyWith((message) => updates(message as GetCommitteeInfoResponse))
          as GetCommitteeInfoResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetCommitteeInfoResponse create() => GetCommitteeInfoResponse._();
  @$core.override
  GetCommitteeInfoResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetCommitteeInfoResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetCommitteeInfoResponse>(create);
  static GetCommitteeInfoResponse? _defaultInstance;

  /// The power of the committee.
  @$pb.TagNumber(1)
  $fixnum.Int64 get committeePower => $_getI64(0);
  @$pb.TagNumber(1)
  set committeePower($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(1)
  $core.bool hasCommitteePower() => $_has(0);
  @$pb.TagNumber(1)
  void clearCommitteePower() => $_clearField(1);

  /// List of committee validators.
  @$pb.TagNumber(2)
  $pb.PbList<ValidatorInfo> get validators => $_getList(1);

  /// Map of protocol versions and their percentages in the committee.
  @$pb.TagNumber(3)
  $pb.PbMap<$core.int, $core.double> get protocolVersions => $_getMap(2);
}

/// Request message for retrieving consensus information.
class GetConsensusInfoRequest extends $pb.GeneratedMessage {
  factory GetConsensusInfoRequest() => create();

  GetConsensusInfoRequest._();

  factory GetConsensusInfoRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetConsensusInfoRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetConsensusInfoRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetConsensusInfoRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetConsensusInfoRequest copyWith(
          void Function(GetConsensusInfoRequest) updates) =>
      super.copyWith((message) => updates(message as GetConsensusInfoRequest))
          as GetConsensusInfoRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetConsensusInfoRequest create() => GetConsensusInfoRequest._();
  @$core.override
  GetConsensusInfoRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetConsensusInfoRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetConsensusInfoRequest>(create);
  static GetConsensusInfoRequest? _defaultInstance;
}

/// Response message contains consensus information.
class GetConsensusInfoResponse extends $pb.GeneratedMessage {
  factory GetConsensusInfoResponse({
    ProposalInfo? proposal,
    $core.Iterable<ConsensusInfo>? instances,
  }) {
    final result = create();
    if (proposal != null) result.proposal = proposal;
    if (instances != null) result.instances.addAll(instances);
    return result;
  }

  GetConsensusInfoResponse._();

  factory GetConsensusInfoResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetConsensusInfoResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetConsensusInfoResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOM<ProposalInfo>(1, _omitFieldNames ? '' : 'proposal',
        subBuilder: ProposalInfo.create)
    ..pPM<ConsensusInfo>(2, _omitFieldNames ? '' : 'instances',
        subBuilder: ConsensusInfo.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetConsensusInfoResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetConsensusInfoResponse copyWith(
          void Function(GetConsensusInfoResponse) updates) =>
      super.copyWith((message) => updates(message as GetConsensusInfoResponse))
          as GetConsensusInfoResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetConsensusInfoResponse create() => GetConsensusInfoResponse._();
  @$core.override
  GetConsensusInfoResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetConsensusInfoResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetConsensusInfoResponse>(create);
  static GetConsensusInfoResponse? _defaultInstance;

  /// The proposal of the consensus info.
  @$pb.TagNumber(1)
  ProposalInfo get proposal => $_getN(0);
  @$pb.TagNumber(1)
  set proposal(ProposalInfo value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasProposal() => $_has(0);
  @$pb.TagNumber(1)
  void clearProposal() => $_clearField(1);
  @$pb.TagNumber(1)
  ProposalInfo ensureProposal() => $_ensure(0);

  /// List of consensus instances.
  @$pb.TagNumber(2)
  $pb.PbList<ConsensusInfo> get instances => $_getList(1);
}

/// Request message for retrieving transactions in the transaction pool.
class GetTxPoolContentRequest extends $pb.GeneratedMessage {
  factory GetTxPoolContentRequest({
    $0.PayloadType? payloadType,
  }) {
    final result = create();
    if (payloadType != null) result.payloadType = payloadType;
    return result;
  }

  GetTxPoolContentRequest._();

  factory GetTxPoolContentRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetTxPoolContentRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetTxPoolContentRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aE<$0.PayloadType>(1, _omitFieldNames ? '' : 'payloadType',
        enumValues: $0.PayloadType.values)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetTxPoolContentRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetTxPoolContentRequest copyWith(
          void Function(GetTxPoolContentRequest) updates) =>
      super.copyWith((message) => updates(message as GetTxPoolContentRequest))
          as GetTxPoolContentRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetTxPoolContentRequest create() => GetTxPoolContentRequest._();
  @$core.override
  GetTxPoolContentRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetTxPoolContentRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetTxPoolContentRequest>(create);
  static GetTxPoolContentRequest? _defaultInstance;

  /// The type of transactions to retrieve from the transaction pool. 0 means all types.
  @$pb.TagNumber(1)
  $0.PayloadType get payloadType => $_getN(0);
  @$pb.TagNumber(1)
  set payloadType($0.PayloadType value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasPayloadType() => $_has(0);
  @$pb.TagNumber(1)
  void clearPayloadType() => $_clearField(1);
}

/// Response message contains transactions in the transaction pool.
class GetTxPoolContentResponse extends $pb.GeneratedMessage {
  factory GetTxPoolContentResponse({
    $core.Iterable<$0.TransactionInfo>? txs,
  }) {
    final result = create();
    if (txs != null) result.txs.addAll(txs);
    return result;
  }

  GetTxPoolContentResponse._();

  factory GetTxPoolContentResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetTxPoolContentResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetTxPoolContentResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..pPM<$0.TransactionInfo>(1, _omitFieldNames ? '' : 'txs',
        subBuilder: $0.TransactionInfo.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetTxPoolContentResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetTxPoolContentResponse copyWith(
          void Function(GetTxPoolContentResponse) updates) =>
      super.copyWith((message) => updates(message as GetTxPoolContentResponse))
          as GetTxPoolContentResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetTxPoolContentResponse create() => GetTxPoolContentResponse._();
  @$core.override
  GetTxPoolContentResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetTxPoolContentResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetTxPoolContentResponse>(create);
  static GetTxPoolContentResponse? _defaultInstance;

  /// List of transactions currently in the pool.
  @$pb.TagNumber(1)
  $pb.PbList<$0.TransactionInfo> get txs => $_getList(0);
}

/// Message contains information about a validator.
class ValidatorInfo extends $pb.GeneratedMessage {
  factory ValidatorInfo({
    $core.String? hash,
    $core.String? data,
    $core.String? publicKey,
    $core.int? number,
    $fixnum.Int64? stake,
    $core.int? lastBondingHeight,
    $core.int? lastSortitionHeight,
    $core.int? unbondingHeight,
    $core.String? address,
    $core.double? availabilityScore,
    $core.int? protocolVersion,
  }) {
    final result = create();
    if (hash != null) result.hash = hash;
    if (data != null) result.data = data;
    if (publicKey != null) result.publicKey = publicKey;
    if (number != null) result.number = number;
    if (stake != null) result.stake = stake;
    if (lastBondingHeight != null) result.lastBondingHeight = lastBondingHeight;
    if (lastSortitionHeight != null)
      result.lastSortitionHeight = lastSortitionHeight;
    if (unbondingHeight != null) result.unbondingHeight = unbondingHeight;
    if (address != null) result.address = address;
    if (availabilityScore != null) result.availabilityScore = availabilityScore;
    if (protocolVersion != null) result.protocolVersion = protocolVersion;
    return result;
  }

  ValidatorInfo._();

  factory ValidatorInfo.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ValidatorInfo.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ValidatorInfo',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'hash')
    ..aOS(2, _omitFieldNames ? '' : 'data')
    ..aOS(3, _omitFieldNames ? '' : 'publicKey')
    ..aI(4, _omitFieldNames ? '' : 'number')
    ..aInt64(5, _omitFieldNames ? '' : 'stake')
    ..aI(6, _omitFieldNames ? '' : 'lastBondingHeight',
        fieldType: $pb.PbFieldType.OU3)
    ..aI(7, _omitFieldNames ? '' : 'lastSortitionHeight',
        fieldType: $pb.PbFieldType.OU3)
    ..aI(8, _omitFieldNames ? '' : 'unbondingHeight',
        fieldType: $pb.PbFieldType.OU3)
    ..aOS(9, _omitFieldNames ? '' : 'address')
    ..aD(10, _omitFieldNames ? '' : 'availabilityScore')
    ..aI(11, _omitFieldNames ? '' : 'protocolVersion')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ValidatorInfo clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ValidatorInfo copyWith(void Function(ValidatorInfo) updates) =>
      super.copyWith((message) => updates(message as ValidatorInfo))
          as ValidatorInfo;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ValidatorInfo create() => ValidatorInfo._();
  @$core.override
  ValidatorInfo createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ValidatorInfo getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ValidatorInfo>(create);
  static ValidatorInfo? _defaultInstance;

  /// The hash of the validator.
  @$pb.TagNumber(1)
  $core.String get hash => $_getSZ(0);
  @$pb.TagNumber(1)
  set hash($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasHash() => $_has(0);
  @$pb.TagNumber(1)
  void clearHash() => $_clearField(1);

  /// The serialized data of the validator.
  @$pb.TagNumber(2)
  $core.String get data => $_getSZ(1);
  @$pb.TagNumber(2)
  set data($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasData() => $_has(1);
  @$pb.TagNumber(2)
  void clearData() => $_clearField(2);

  /// The public key of the validator.
  @$pb.TagNumber(3)
  $core.String get publicKey => $_getSZ(2);
  @$pb.TagNumber(3)
  set publicKey($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasPublicKey() => $_has(2);
  @$pb.TagNumber(3)
  void clearPublicKey() => $_clearField(3);

  /// The unique number assigned to the validator.
  @$pb.TagNumber(4)
  $core.int get number => $_getIZ(3);
  @$pb.TagNumber(4)
  set number($core.int value) => $_setSignedInt32(3, value);
  @$pb.TagNumber(4)
  $core.bool hasNumber() => $_has(3);
  @$pb.TagNumber(4)
  void clearNumber() => $_clearField(4);

  /// The stake of the validator in NanoPAC.
  @$pb.TagNumber(5)
  $fixnum.Int64 get stake => $_getI64(4);
  @$pb.TagNumber(5)
  set stake($fixnum.Int64 value) => $_setInt64(4, value);
  @$pb.TagNumber(5)
  $core.bool hasStake() => $_has(4);
  @$pb.TagNumber(5)
  void clearStake() => $_clearField(5);

  /// The height at which the validator last bonded.
  @$pb.TagNumber(6)
  $core.int get lastBondingHeight => $_getIZ(5);
  @$pb.TagNumber(6)
  set lastBondingHeight($core.int value) => $_setUnsignedInt32(5, value);
  @$pb.TagNumber(6)
  $core.bool hasLastBondingHeight() => $_has(5);
  @$pb.TagNumber(6)
  void clearLastBondingHeight() => $_clearField(6);

  /// The height at which the validator last participated in sortition.
  @$pb.TagNumber(7)
  $core.int get lastSortitionHeight => $_getIZ(6);
  @$pb.TagNumber(7)
  set lastSortitionHeight($core.int value) => $_setUnsignedInt32(6, value);
  @$pb.TagNumber(7)
  $core.bool hasLastSortitionHeight() => $_has(6);
  @$pb.TagNumber(7)
  void clearLastSortitionHeight() => $_clearField(7);

  /// The height at which the validator will unbond.
  @$pb.TagNumber(8)
  $core.int get unbondingHeight => $_getIZ(7);
  @$pb.TagNumber(8)
  set unbondingHeight($core.int value) => $_setUnsignedInt32(7, value);
  @$pb.TagNumber(8)
  $core.bool hasUnbondingHeight() => $_has(7);
  @$pb.TagNumber(8)
  void clearUnbondingHeight() => $_clearField(8);

  /// The address of the validator.
  @$pb.TagNumber(9)
  $core.String get address => $_getSZ(8);
  @$pb.TagNumber(9)
  set address($core.String value) => $_setString(8, value);
  @$pb.TagNumber(9)
  $core.bool hasAddress() => $_has(8);
  @$pb.TagNumber(9)
  void clearAddress() => $_clearField(9);

  /// The availability score of the validator.
  @$pb.TagNumber(10)
  $core.double get availabilityScore => $_getN(9);
  @$pb.TagNumber(10)
  set availabilityScore($core.double value) => $_setDouble(9, value);
  @$pb.TagNumber(10)
  $core.bool hasAvailabilityScore() => $_has(9);
  @$pb.TagNumber(10)
  void clearAvailabilityScore() => $_clearField(10);

  /// The protocol version of the validator.
  @$pb.TagNumber(11)
  $core.int get protocolVersion => $_getIZ(10);
  @$pb.TagNumber(11)
  set protocolVersion($core.int value) => $_setSignedInt32(10, value);
  @$pb.TagNumber(11)
  $core.bool hasProtocolVersion() => $_has(10);
  @$pb.TagNumber(11)
  void clearProtocolVersion() => $_clearField(11);
}

/// Message contains information about an account.
class AccountInfo extends $pb.GeneratedMessage {
  factory AccountInfo({
    $core.String? hash,
    $core.String? data,
    $core.int? number,
    $fixnum.Int64? balance,
    $core.String? address,
  }) {
    final result = create();
    if (hash != null) result.hash = hash;
    if (data != null) result.data = data;
    if (number != null) result.number = number;
    if (balance != null) result.balance = balance;
    if (address != null) result.address = address;
    return result;
  }

  AccountInfo._();

  factory AccountInfo.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory AccountInfo.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'AccountInfo',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'hash')
    ..aOS(2, _omitFieldNames ? '' : 'data')
    ..aI(3, _omitFieldNames ? '' : 'number')
    ..aInt64(4, _omitFieldNames ? '' : 'balance')
    ..aOS(5, _omitFieldNames ? '' : 'address')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AccountInfo clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  AccountInfo copyWith(void Function(AccountInfo) updates) =>
      super.copyWith((message) => updates(message as AccountInfo))
          as AccountInfo;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AccountInfo create() => AccountInfo._();
  @$core.override
  AccountInfo createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static AccountInfo getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<AccountInfo>(create);
  static AccountInfo? _defaultInstance;

  /// The hash of the account.
  @$pb.TagNumber(1)
  $core.String get hash => $_getSZ(0);
  @$pb.TagNumber(1)
  set hash($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasHash() => $_has(0);
  @$pb.TagNumber(1)
  void clearHash() => $_clearField(1);

  /// The serialized data of the account.
  @$pb.TagNumber(2)
  $core.String get data => $_getSZ(1);
  @$pb.TagNumber(2)
  set data($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasData() => $_has(1);
  @$pb.TagNumber(2)
  void clearData() => $_clearField(2);

  /// The unique number assigned to the account.
  @$pb.TagNumber(3)
  $core.int get number => $_getIZ(2);
  @$pb.TagNumber(3)
  set number($core.int value) => $_setSignedInt32(2, value);
  @$pb.TagNumber(3)
  $core.bool hasNumber() => $_has(2);
  @$pb.TagNumber(3)
  void clearNumber() => $_clearField(3);

  /// The balance of the account in NanoPAC.
  @$pb.TagNumber(4)
  $fixnum.Int64 get balance => $_getI64(3);
  @$pb.TagNumber(4)
  set balance($fixnum.Int64 value) => $_setInt64(3, value);
  @$pb.TagNumber(4)
  $core.bool hasBalance() => $_has(3);
  @$pb.TagNumber(4)
  void clearBalance() => $_clearField(4);

  /// The address of the account.
  @$pb.TagNumber(5)
  $core.String get address => $_getSZ(4);
  @$pb.TagNumber(5)
  set address($core.String value) => $_setString(4, value);
  @$pb.TagNumber(5)
  $core.bool hasAddress() => $_has(4);
  @$pb.TagNumber(5)
  void clearAddress() => $_clearField(5);
}

/// Message contains information about the header of a block.
class BlockHeaderInfo extends $pb.GeneratedMessage {
  factory BlockHeaderInfo({
    $core.int? version,
    $core.String? prevBlockHash,
    $core.String? stateRoot,
    $core.String? sortitionSeed,
    $core.String? proposerAddress,
  }) {
    final result = create();
    if (version != null) result.version = version;
    if (prevBlockHash != null) result.prevBlockHash = prevBlockHash;
    if (stateRoot != null) result.stateRoot = stateRoot;
    if (sortitionSeed != null) result.sortitionSeed = sortitionSeed;
    if (proposerAddress != null) result.proposerAddress = proposerAddress;
    return result;
  }

  BlockHeaderInfo._();

  factory BlockHeaderInfo.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory BlockHeaderInfo.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'BlockHeaderInfo',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aI(1, _omitFieldNames ? '' : 'version')
    ..aOS(2, _omitFieldNames ? '' : 'prevBlockHash')
    ..aOS(3, _omitFieldNames ? '' : 'stateRoot')
    ..aOS(4, _omitFieldNames ? '' : 'sortitionSeed')
    ..aOS(5, _omitFieldNames ? '' : 'proposerAddress')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BlockHeaderInfo clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  BlockHeaderInfo copyWith(void Function(BlockHeaderInfo) updates) =>
      super.copyWith((message) => updates(message as BlockHeaderInfo))
          as BlockHeaderInfo;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static BlockHeaderInfo create() => BlockHeaderInfo._();
  @$core.override
  BlockHeaderInfo createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static BlockHeaderInfo getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<BlockHeaderInfo>(create);
  static BlockHeaderInfo? _defaultInstance;

  /// The version of the block.
  @$pb.TagNumber(1)
  $core.int get version => $_getIZ(0);
  @$pb.TagNumber(1)
  set version($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(1)
  $core.bool hasVersion() => $_has(0);
  @$pb.TagNumber(1)
  void clearVersion() => $_clearField(1);

  /// The hash of the previous block.
  @$pb.TagNumber(2)
  $core.String get prevBlockHash => $_getSZ(1);
  @$pb.TagNumber(2)
  set prevBlockHash($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasPrevBlockHash() => $_has(1);
  @$pb.TagNumber(2)
  void clearPrevBlockHash() => $_clearField(2);

  /// The state root hash of the blockchain.
  @$pb.TagNumber(3)
  $core.String get stateRoot => $_getSZ(2);
  @$pb.TagNumber(3)
  set stateRoot($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasStateRoot() => $_has(2);
  @$pb.TagNumber(3)
  void clearStateRoot() => $_clearField(3);

  /// The sortition seed of the block.
  @$pb.TagNumber(4)
  $core.String get sortitionSeed => $_getSZ(3);
  @$pb.TagNumber(4)
  set sortitionSeed($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasSortitionSeed() => $_has(3);
  @$pb.TagNumber(4)
  void clearSortitionSeed() => $_clearField(4);

  /// The address of the proposer of the block.
  @$pb.TagNumber(5)
  $core.String get proposerAddress => $_getSZ(4);
  @$pb.TagNumber(5)
  set proposerAddress($core.String value) => $_setString(4, value);
  @$pb.TagNumber(5)
  $core.bool hasProposerAddress() => $_has(4);
  @$pb.TagNumber(5)
  void clearProposerAddress() => $_clearField(5);
}

/// Message contains information about a certificate.
class CertificateInfo extends $pb.GeneratedMessage {
  factory CertificateInfo({
    $core.String? hash,
    $core.int? round,
    $core.Iterable<$core.int>? committers,
    $core.Iterable<$core.int>? absentees,
    $core.String? signature,
  }) {
    final result = create();
    if (hash != null) result.hash = hash;
    if (round != null) result.round = round;
    if (committers != null) result.committers.addAll(committers);
    if (absentees != null) result.absentees.addAll(absentees);
    if (signature != null) result.signature = signature;
    return result;
  }

  CertificateInfo._();

  factory CertificateInfo.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CertificateInfo.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CertificateInfo',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'hash')
    ..aI(2, _omitFieldNames ? '' : 'round')
    ..p<$core.int>(3, _omitFieldNames ? '' : 'committers', $pb.PbFieldType.K3)
    ..p<$core.int>(4, _omitFieldNames ? '' : 'absentees', $pb.PbFieldType.K3)
    ..aOS(5, _omitFieldNames ? '' : 'signature')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CertificateInfo clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CertificateInfo copyWith(void Function(CertificateInfo) updates) =>
      super.copyWith((message) => updates(message as CertificateInfo))
          as CertificateInfo;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CertificateInfo create() => CertificateInfo._();
  @$core.override
  CertificateInfo createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CertificateInfo getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CertificateInfo>(create);
  static CertificateInfo? _defaultInstance;

  /// The hash of the certificate.
  @$pb.TagNumber(1)
  $core.String get hash => $_getSZ(0);
  @$pb.TagNumber(1)
  set hash($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasHash() => $_has(0);
  @$pb.TagNumber(1)
  void clearHash() => $_clearField(1);

  /// The round of the certificate.
  @$pb.TagNumber(2)
  $core.int get round => $_getIZ(1);
  @$pb.TagNumber(2)
  set round($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(2)
  $core.bool hasRound() => $_has(1);
  @$pb.TagNumber(2)
  void clearRound() => $_clearField(2);

  /// List of committers in the certificate.
  @$pb.TagNumber(3)
  $pb.PbList<$core.int> get committers => $_getList(2);

  /// List of absentees in the certificate.
  @$pb.TagNumber(4)
  $pb.PbList<$core.int> get absentees => $_getList(3);

  /// The signature of the certificate.
  @$pb.TagNumber(5)
  $core.String get signature => $_getSZ(4);
  @$pb.TagNumber(5)
  set signature($core.String value) => $_setString(4, value);
  @$pb.TagNumber(5)
  $core.bool hasSignature() => $_has(4);
  @$pb.TagNumber(5)
  void clearSignature() => $_clearField(5);
}

/// Message contains information about a vote.
class VoteInfo extends $pb.GeneratedMessage {
  factory VoteInfo({
    VoteType? type,
    $core.String? voter,
    $core.String? blockHash,
    $core.int? round,
    $core.int? cpRound,
    $core.int? cpValue,
  }) {
    final result = create();
    if (type != null) result.type = type;
    if (voter != null) result.voter = voter;
    if (blockHash != null) result.blockHash = blockHash;
    if (round != null) result.round = round;
    if (cpRound != null) result.cpRound = cpRound;
    if (cpValue != null) result.cpValue = cpValue;
    return result;
  }

  VoteInfo._();

  factory VoteInfo.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory VoteInfo.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'VoteInfo',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aE<VoteType>(1, _omitFieldNames ? '' : 'type',
        enumValues: VoteType.values)
    ..aOS(2, _omitFieldNames ? '' : 'voter')
    ..aOS(3, _omitFieldNames ? '' : 'blockHash')
    ..aI(4, _omitFieldNames ? '' : 'round')
    ..aI(5, _omitFieldNames ? '' : 'cpRound')
    ..aI(6, _omitFieldNames ? '' : 'cpValue')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  VoteInfo clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  VoteInfo copyWith(void Function(VoteInfo) updates) =>
      super.copyWith((message) => updates(message as VoteInfo)) as VoteInfo;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static VoteInfo create() => VoteInfo._();
  @$core.override
  VoteInfo createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static VoteInfo getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<VoteInfo>(create);
  static VoteInfo? _defaultInstance;

  /// The type of the vote.
  @$pb.TagNumber(1)
  VoteType get type => $_getN(0);
  @$pb.TagNumber(1)
  set type(VoteType value) => $_setField(1, value);
  @$pb.TagNumber(1)
  $core.bool hasType() => $_has(0);
  @$pb.TagNumber(1)
  void clearType() => $_clearField(1);

  /// The address of the voter.
  @$pb.TagNumber(2)
  $core.String get voter => $_getSZ(1);
  @$pb.TagNumber(2)
  set voter($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasVoter() => $_has(1);
  @$pb.TagNumber(2)
  void clearVoter() => $_clearField(2);

  /// The hash of the block being voted on.
  @$pb.TagNumber(3)
  $core.String get blockHash => $_getSZ(2);
  @$pb.TagNumber(3)
  set blockHash($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasBlockHash() => $_has(2);
  @$pb.TagNumber(3)
  void clearBlockHash() => $_clearField(3);

  /// The consensus round of the vote.
  @$pb.TagNumber(4)
  $core.int get round => $_getIZ(3);
  @$pb.TagNumber(4)
  set round($core.int value) => $_setSignedInt32(3, value);
  @$pb.TagNumber(4)
  $core.bool hasRound() => $_has(3);
  @$pb.TagNumber(4)
  void clearRound() => $_clearField(4);

  /// The change-proposer round of the vote.
  @$pb.TagNumber(5)
  $core.int get cpRound => $_getIZ(4);
  @$pb.TagNumber(5)
  set cpRound($core.int value) => $_setSignedInt32(4, value);
  @$pb.TagNumber(5)
  $core.bool hasCpRound() => $_has(4);
  @$pb.TagNumber(5)
  void clearCpRound() => $_clearField(5);

  /// The change-proposer value of the vote.
  @$pb.TagNumber(6)
  $core.int get cpValue => $_getIZ(5);
  @$pb.TagNumber(6)
  set cpValue($core.int value) => $_setSignedInt32(5, value);
  @$pb.TagNumber(6)
  $core.bool hasCpValue() => $_has(5);
  @$pb.TagNumber(6)
  void clearCpValue() => $_clearField(6);
}

/// Message contains information about a consensus instance.
class ConsensusInfo extends $pb.GeneratedMessage {
  factory ConsensusInfo({
    $core.String? address,
    $core.bool? active,
    $core.int? height,
    $core.int? round,
    $core.Iterable<VoteInfo>? votes,
  }) {
    final result = create();
    if (address != null) result.address = address;
    if (active != null) result.active = active;
    if (height != null) result.height = height;
    if (round != null) result.round = round;
    if (votes != null) result.votes.addAll(votes);
    return result;
  }

  ConsensusInfo._();

  factory ConsensusInfo.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ConsensusInfo.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ConsensusInfo',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'address')
    ..aOB(2, _omitFieldNames ? '' : 'active')
    ..aI(3, _omitFieldNames ? '' : 'height', fieldType: $pb.PbFieldType.OU3)
    ..aI(4, _omitFieldNames ? '' : 'round')
    ..pPM<VoteInfo>(5, _omitFieldNames ? '' : 'votes',
        subBuilder: VoteInfo.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ConsensusInfo clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ConsensusInfo copyWith(void Function(ConsensusInfo) updates) =>
      super.copyWith((message) => updates(message as ConsensusInfo))
          as ConsensusInfo;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ConsensusInfo create() => ConsensusInfo._();
  @$core.override
  ConsensusInfo createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ConsensusInfo getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ConsensusInfo>(create);
  static ConsensusInfo? _defaultInstance;

  /// The address of the consensus instance.
  @$pb.TagNumber(1)
  $core.String get address => $_getSZ(0);
  @$pb.TagNumber(1)
  set address($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasAddress() => $_has(0);
  @$pb.TagNumber(1)
  void clearAddress() => $_clearField(1);

  /// Indicates whether the consensus instance is active and part of the committee.
  @$pb.TagNumber(2)
  $core.bool get active => $_getBF(1);
  @$pb.TagNumber(2)
  set active($core.bool value) => $_setBool(1, value);
  @$pb.TagNumber(2)
  $core.bool hasActive() => $_has(1);
  @$pb.TagNumber(2)
  void clearActive() => $_clearField(2);

  /// The height of the consensus instance.
  @$pb.TagNumber(3)
  $core.int get height => $_getIZ(2);
  @$pb.TagNumber(3)
  set height($core.int value) => $_setUnsignedInt32(2, value);
  @$pb.TagNumber(3)
  $core.bool hasHeight() => $_has(2);
  @$pb.TagNumber(3)
  void clearHeight() => $_clearField(3);

  /// The round of the consensus instance.
  @$pb.TagNumber(4)
  $core.int get round => $_getIZ(3);
  @$pb.TagNumber(4)
  set round($core.int value) => $_setSignedInt32(3, value);
  @$pb.TagNumber(4)
  $core.bool hasRound() => $_has(3);
  @$pb.TagNumber(4)
  void clearRound() => $_clearField(4);

  /// List of votes in the consensus instance.
  @$pb.TagNumber(5)
  $pb.PbList<VoteInfo> get votes => $_getList(4);
}

/// Message contains information about a proposal.
class ProposalInfo extends $pb.GeneratedMessage {
  factory ProposalInfo({
    $core.int? height,
    $core.int? round,
    $core.String? blockData,
    $core.String? signature,
  }) {
    final result = create();
    if (height != null) result.height = height;
    if (round != null) result.round = round;
    if (blockData != null) result.blockData = blockData;
    if (signature != null) result.signature = signature;
    return result;
  }

  ProposalInfo._();

  factory ProposalInfo.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ProposalInfo.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ProposalInfo',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aI(1, _omitFieldNames ? '' : 'height', fieldType: $pb.PbFieldType.OU3)
    ..aI(2, _omitFieldNames ? '' : 'round')
    ..aOS(3, _omitFieldNames ? '' : 'blockData')
    ..aOS(4, _omitFieldNames ? '' : 'signature')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ProposalInfo clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ProposalInfo copyWith(void Function(ProposalInfo) updates) =>
      super.copyWith((message) => updates(message as ProposalInfo))
          as ProposalInfo;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ProposalInfo create() => ProposalInfo._();
  @$core.override
  ProposalInfo createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ProposalInfo getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ProposalInfo>(create);
  static ProposalInfo? _defaultInstance;

  /// The height of the proposal.
  @$pb.TagNumber(1)
  $core.int get height => $_getIZ(0);
  @$pb.TagNumber(1)
  set height($core.int value) => $_setUnsignedInt32(0, value);
  @$pb.TagNumber(1)
  $core.bool hasHeight() => $_has(0);
  @$pb.TagNumber(1)
  void clearHeight() => $_clearField(1);

  /// The round of the proposal.
  @$pb.TagNumber(2)
  $core.int get round => $_getIZ(1);
  @$pb.TagNumber(2)
  set round($core.int value) => $_setSignedInt32(1, value);
  @$pb.TagNumber(2)
  $core.bool hasRound() => $_has(1);
  @$pb.TagNumber(2)
  void clearRound() => $_clearField(2);

  /// The block data of the proposal.
  @$pb.TagNumber(3)
  $core.String get blockData => $_getSZ(2);
  @$pb.TagNumber(3)
  set blockData($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasBlockData() => $_has(2);
  @$pb.TagNumber(3)
  void clearBlockData() => $_clearField(3);

  /// The signature of the proposal, signed by the proposer.
  @$pb.TagNumber(4)
  $core.String get signature => $_getSZ(3);
  @$pb.TagNumber(4)
  set signature($core.String value) => $_setString(3, value);
  @$pb.TagNumber(4)
  $core.bool hasSignature() => $_has(3);
  @$pb.TagNumber(4)
  void clearSignature() => $_clearField(4);
}

/// Blockchain service defines RPC methods for interacting with the blockchain.
class BlockchainApi {
  final $pb.RpcClient _client;

  BlockchainApi(this._client);

  /// GetBlock retrieves information about a block based on the provided request parameters.
  $async.Future<GetBlockResponse> getBlock(
          $pb.ClientContext? ctx, GetBlockRequest request) =>
      _client.invoke<GetBlockResponse>(
          ctx, 'Blockchain', 'GetBlock', request, GetBlockResponse());

  /// GetBlockHash retrieves the hash of a block at the specified height.
  $async.Future<GetBlockHashResponse> getBlockHash(
          $pb.ClientContext? ctx, GetBlockHashRequest request) =>
      _client.invoke<GetBlockHashResponse>(
          ctx, 'Blockchain', 'GetBlockHash', request, GetBlockHashResponse());

  /// GetBlockHeight retrieves the height of a block with the specified hash.
  $async.Future<GetBlockHeightResponse> getBlockHeight(
          $pb.ClientContext? ctx, GetBlockHeightRequest request) =>
      _client.invoke<GetBlockHeightResponse>(ctx, 'Blockchain',
          'GetBlockHeight', request, GetBlockHeightResponse());

  /// GetBlockchainInfo retrieves general information about the blockchain.
  $async.Future<GetBlockchainInfoResponse> getBlockchainInfo(
          $pb.ClientContext? ctx, GetBlockchainInfoRequest request) =>
      _client.invoke<GetBlockchainInfoResponse>(ctx, 'Blockchain',
          'GetBlockchainInfo', request, GetBlockchainInfoResponse());

  /// GetCommitteeInfo retrieves information about the current committee.
  $async.Future<GetCommitteeInfoResponse> getCommitteeInfo(
          $pb.ClientContext? ctx, GetCommitteeInfoRequest request) =>
      _client.invoke<GetCommitteeInfoResponse>(ctx, 'Blockchain',
          'GetCommitteeInfo', request, GetCommitteeInfoResponse());

  /// GetConsensusInfo retrieves information about the consensus instances.
  $async.Future<GetConsensusInfoResponse> getConsensusInfo(
          $pb.ClientContext? ctx, GetConsensusInfoRequest request) =>
      _client.invoke<GetConsensusInfoResponse>(ctx, 'Blockchain',
          'GetConsensusInfo', request, GetConsensusInfoResponse());

  /// GetAccount retrieves information about an account based on the provided address.
  $async.Future<GetAccountResponse> getAccount(
          $pb.ClientContext? ctx, GetAccountRequest request) =>
      _client.invoke<GetAccountResponse>(
          ctx, 'Blockchain', 'GetAccount', request, GetAccountResponse());

  /// GetValidator retrieves information about a validator based on the provided address.
  $async.Future<GetValidatorResponse> getValidator(
          $pb.ClientContext? ctx, GetValidatorRequest request) =>
      _client.invoke<GetValidatorResponse>(
          ctx, 'Blockchain', 'GetValidator', request, GetValidatorResponse());

  /// GetValidatorByNumber retrieves information about a validator based on the provided number.
  $async.Future<GetValidatorResponse> getValidatorByNumber(
          $pb.ClientContext? ctx, GetValidatorByNumberRequest request) =>
      _client.invoke<GetValidatorResponse>(ctx, 'Blockchain',
          'GetValidatorByNumber', request, GetValidatorResponse());

  /// GetValidatorAddresses retrieves a list of all validator addresses.
  $async.Future<GetValidatorAddressesResponse> getValidatorAddresses(
          $pb.ClientContext? ctx, GetValidatorAddressesRequest request) =>
      _client.invoke<GetValidatorAddressesResponse>(ctx, 'Blockchain',
          'GetValidatorAddresses', request, GetValidatorAddressesResponse());

  /// GetPublicKey retrieves the public key of an account based on the provided address.
  $async.Future<GetPublicKeyResponse> getPublicKey(
          $pb.ClientContext? ctx, GetPublicKeyRequest request) =>
      _client.invoke<GetPublicKeyResponse>(
          ctx, 'Blockchain', 'GetPublicKey', request, GetPublicKeyResponse());

  /// GetTxPoolContent retrieves current transactions in the transaction pool.
  $async.Future<GetTxPoolContentResponse> getTxPoolContent(
          $pb.ClientContext? ctx, GetTxPoolContentRequest request) =>
      _client.invoke<GetTxPoolContentResponse>(ctx, 'Blockchain',
          'GetTxPoolContent', request, GetTxPoolContentResponse());
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
