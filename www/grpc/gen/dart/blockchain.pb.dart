//
//  Generated code. Do not modify.
//  source: blockchain.proto
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

import 'blockchain.pbenum.dart';
import 'transaction.pb.dart' as $0;
import 'transaction.pbenum.dart' as $0;

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

export 'blockchain.pbenum.dart';

/// Request message for retrieving account information.
class GetAccountRequest extends $pb.GeneratedMessage {
  factory GetAccountRequest({
    $core.String? address,
  }) {
    final $result = create();
    if (address != null) {
      $result.address = address;
    }
    return $result;
  }
  GetAccountRequest._() : super();
  factory GetAccountRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetAccountRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetAccountRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'address')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetAccountRequest clone() => GetAccountRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetAccountRequest copyWith(void Function(GetAccountRequest) updates) => super.copyWith((message) => updates(message as GetAccountRequest)) as GetAccountRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetAccountRequest create() => GetAccountRequest._();
  GetAccountRequest createEmptyInstance() => create();
  static $pb.PbList<GetAccountRequest> createRepeated() => $pb.PbList<GetAccountRequest>();
  @$core.pragma('dart2js:noInline')
  static GetAccountRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetAccountRequest>(create);
  static GetAccountRequest? _defaultInstance;

  /// The address of the account to retrieve information for.
  @$pb.TagNumber(1)
  $core.String get address => $_getSZ(0);
  @$pb.TagNumber(1)
  set address($core.String v) { $_setString(0, v); }
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
    final $result = create();
    if (account != null) {
      $result.account = account;
    }
    return $result;
  }
  GetAccountResponse._() : super();
  factory GetAccountResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetAccountResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetAccountResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..aOM<AccountInfo>(1, _omitFieldNames ? '' : 'account', subBuilder: AccountInfo.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetAccountResponse clone() => GetAccountResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetAccountResponse copyWith(void Function(GetAccountResponse) updates) => super.copyWith((message) => updates(message as GetAccountResponse)) as GetAccountResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetAccountResponse create() => GetAccountResponse._();
  GetAccountResponse createEmptyInstance() => create();
  static $pb.PbList<GetAccountResponse> createRepeated() => $pb.PbList<GetAccountResponse>();
  @$core.pragma('dart2js:noInline')
  static GetAccountResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetAccountResponse>(create);
  static GetAccountResponse? _defaultInstance;

  /// Detailed information about the account.
  @$pb.TagNumber(1)
  AccountInfo get account => $_getN(0);
  @$pb.TagNumber(1)
  set account(AccountInfo v) { $_setField(1, v); }
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
  GetValidatorAddressesRequest._() : super();
  factory GetValidatorAddressesRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetValidatorAddressesRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetValidatorAddressesRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetValidatorAddressesRequest clone() => GetValidatorAddressesRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetValidatorAddressesRequest copyWith(void Function(GetValidatorAddressesRequest) updates) => super.copyWith((message) => updates(message as GetValidatorAddressesRequest)) as GetValidatorAddressesRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetValidatorAddressesRequest create() => GetValidatorAddressesRequest._();
  GetValidatorAddressesRequest createEmptyInstance() => create();
  static $pb.PbList<GetValidatorAddressesRequest> createRepeated() => $pb.PbList<GetValidatorAddressesRequest>();
  @$core.pragma('dart2js:noInline')
  static GetValidatorAddressesRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetValidatorAddressesRequest>(create);
  static GetValidatorAddressesRequest? _defaultInstance;
}

/// Response message contains list of validator addresses.
class GetValidatorAddressesResponse extends $pb.GeneratedMessage {
  factory GetValidatorAddressesResponse({
    $core.Iterable<$core.String>? addresses,
  }) {
    final $result = create();
    if (addresses != null) {
      $result.addresses.addAll(addresses);
    }
    return $result;
  }
  GetValidatorAddressesResponse._() : super();
  factory GetValidatorAddressesResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetValidatorAddressesResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetValidatorAddressesResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..pPS(1, _omitFieldNames ? '' : 'addresses')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetValidatorAddressesResponse clone() => GetValidatorAddressesResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetValidatorAddressesResponse copyWith(void Function(GetValidatorAddressesResponse) updates) => super.copyWith((message) => updates(message as GetValidatorAddressesResponse)) as GetValidatorAddressesResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetValidatorAddressesResponse create() => GetValidatorAddressesResponse._();
  GetValidatorAddressesResponse createEmptyInstance() => create();
  static $pb.PbList<GetValidatorAddressesResponse> createRepeated() => $pb.PbList<GetValidatorAddressesResponse>();
  @$core.pragma('dart2js:noInline')
  static GetValidatorAddressesResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetValidatorAddressesResponse>(create);
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
    final $result = create();
    if (address != null) {
      $result.address = address;
    }
    return $result;
  }
  GetValidatorRequest._() : super();
  factory GetValidatorRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetValidatorRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetValidatorRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'address')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetValidatorRequest clone() => GetValidatorRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetValidatorRequest copyWith(void Function(GetValidatorRequest) updates) => super.copyWith((message) => updates(message as GetValidatorRequest)) as GetValidatorRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetValidatorRequest create() => GetValidatorRequest._();
  GetValidatorRequest createEmptyInstance() => create();
  static $pb.PbList<GetValidatorRequest> createRepeated() => $pb.PbList<GetValidatorRequest>();
  @$core.pragma('dart2js:noInline')
  static GetValidatorRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetValidatorRequest>(create);
  static GetValidatorRequest? _defaultInstance;

  /// The address of the validator to retrieve information for.
  @$pb.TagNumber(1)
  $core.String get address => $_getSZ(0);
  @$pb.TagNumber(1)
  set address($core.String v) { $_setString(0, v); }
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
    final $result = create();
    if (number != null) {
      $result.number = number;
    }
    return $result;
  }
  GetValidatorByNumberRequest._() : super();
  factory GetValidatorByNumberRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetValidatorByNumberRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetValidatorByNumberRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..a<$core.int>(1, _omitFieldNames ? '' : 'number', $pb.PbFieldType.O3)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetValidatorByNumberRequest clone() => GetValidatorByNumberRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetValidatorByNumberRequest copyWith(void Function(GetValidatorByNumberRequest) updates) => super.copyWith((message) => updates(message as GetValidatorByNumberRequest)) as GetValidatorByNumberRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetValidatorByNumberRequest create() => GetValidatorByNumberRequest._();
  GetValidatorByNumberRequest createEmptyInstance() => create();
  static $pb.PbList<GetValidatorByNumberRequest> createRepeated() => $pb.PbList<GetValidatorByNumberRequest>();
  @$core.pragma('dart2js:noInline')
  static GetValidatorByNumberRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetValidatorByNumberRequest>(create);
  static GetValidatorByNumberRequest? _defaultInstance;

  /// The unique number of the validator to retrieve information for.
  @$pb.TagNumber(1)
  $core.int get number => $_getIZ(0);
  @$pb.TagNumber(1)
  set number($core.int v) { $_setSignedInt32(0, v); }
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
    final $result = create();
    if (validator != null) {
      $result.validator = validator;
    }
    return $result;
  }
  GetValidatorResponse._() : super();
  factory GetValidatorResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetValidatorResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetValidatorResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..aOM<ValidatorInfo>(1, _omitFieldNames ? '' : 'validator', subBuilder: ValidatorInfo.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetValidatorResponse clone() => GetValidatorResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetValidatorResponse copyWith(void Function(GetValidatorResponse) updates) => super.copyWith((message) => updates(message as GetValidatorResponse)) as GetValidatorResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetValidatorResponse create() => GetValidatorResponse._();
  GetValidatorResponse createEmptyInstance() => create();
  static $pb.PbList<GetValidatorResponse> createRepeated() => $pb.PbList<GetValidatorResponse>();
  @$core.pragma('dart2js:noInline')
  static GetValidatorResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetValidatorResponse>(create);
  static GetValidatorResponse? _defaultInstance;

  /// Detailed information about the validator.
  @$pb.TagNumber(1)
  ValidatorInfo get validator => $_getN(0);
  @$pb.TagNumber(1)
  set validator(ValidatorInfo v) { $_setField(1, v); }
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
    final $result = create();
    if (address != null) {
      $result.address = address;
    }
    return $result;
  }
  GetPublicKeyRequest._() : super();
  factory GetPublicKeyRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetPublicKeyRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetPublicKeyRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'address')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetPublicKeyRequest clone() => GetPublicKeyRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetPublicKeyRequest copyWith(void Function(GetPublicKeyRequest) updates) => super.copyWith((message) => updates(message as GetPublicKeyRequest)) as GetPublicKeyRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetPublicKeyRequest create() => GetPublicKeyRequest._();
  GetPublicKeyRequest createEmptyInstance() => create();
  static $pb.PbList<GetPublicKeyRequest> createRepeated() => $pb.PbList<GetPublicKeyRequest>();
  @$core.pragma('dart2js:noInline')
  static GetPublicKeyRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetPublicKeyRequest>(create);
  static GetPublicKeyRequest? _defaultInstance;

  /// The address for which to retrieve the public key.
  @$pb.TagNumber(1)
  $core.String get address => $_getSZ(0);
  @$pb.TagNumber(1)
  set address($core.String v) { $_setString(0, v); }
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
    final $result = create();
    if (publicKey != null) {
      $result.publicKey = publicKey;
    }
    return $result;
  }
  GetPublicKeyResponse._() : super();
  factory GetPublicKeyResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetPublicKeyResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetPublicKeyResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'publicKey')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetPublicKeyResponse clone() => GetPublicKeyResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetPublicKeyResponse copyWith(void Function(GetPublicKeyResponse) updates) => super.copyWith((message) => updates(message as GetPublicKeyResponse)) as GetPublicKeyResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetPublicKeyResponse create() => GetPublicKeyResponse._();
  GetPublicKeyResponse createEmptyInstance() => create();
  static $pb.PbList<GetPublicKeyResponse> createRepeated() => $pb.PbList<GetPublicKeyResponse>();
  @$core.pragma('dart2js:noInline')
  static GetPublicKeyResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetPublicKeyResponse>(create);
  static GetPublicKeyResponse? _defaultInstance;

  /// The public key associated with the provided address.
  @$pb.TagNumber(1)
  $core.String get publicKey => $_getSZ(0);
  @$pb.TagNumber(1)
  set publicKey($core.String v) { $_setString(0, v); }
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
    final $result = create();
    if (height != null) {
      $result.height = height;
    }
    if (verbosity != null) {
      $result.verbosity = verbosity;
    }
    return $result;
  }
  GetBlockRequest._() : super();
  factory GetBlockRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetBlockRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetBlockRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..a<$core.int>(1, _omitFieldNames ? '' : 'height', $pb.PbFieldType.OU3)
    ..e<BlockVerbosity>(2, _omitFieldNames ? '' : 'verbosity', $pb.PbFieldType.OE, defaultOrMaker: BlockVerbosity.BLOCK_VERBOSITY_DATA, valueOf: BlockVerbosity.valueOf, enumValues: BlockVerbosity.values)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetBlockRequest clone() => GetBlockRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetBlockRequest copyWith(void Function(GetBlockRequest) updates) => super.copyWith((message) => updates(message as GetBlockRequest)) as GetBlockRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetBlockRequest create() => GetBlockRequest._();
  GetBlockRequest createEmptyInstance() => create();
  static $pb.PbList<GetBlockRequest> createRepeated() => $pb.PbList<GetBlockRequest>();
  @$core.pragma('dart2js:noInline')
  static GetBlockRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetBlockRequest>(create);
  static GetBlockRequest? _defaultInstance;

  /// The height of the block to retrieve.
  @$pb.TagNumber(1)
  $core.int get height => $_getIZ(0);
  @$pb.TagNumber(1)
  set height($core.int v) { $_setUnsignedInt32(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasHeight() => $_has(0);
  @$pb.TagNumber(1)
  void clearHeight() => $_clearField(1);

  /// The verbosity level for block information.
  @$pb.TagNumber(2)
  BlockVerbosity get verbosity => $_getN(1);
  @$pb.TagNumber(2)
  set verbosity(BlockVerbosity v) { $_setField(2, v); }
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
    final $result = create();
    if (height != null) {
      $result.height = height;
    }
    if (hash != null) {
      $result.hash = hash;
    }
    if (data != null) {
      $result.data = data;
    }
    if (blockTime != null) {
      $result.blockTime = blockTime;
    }
    if (header != null) {
      $result.header = header;
    }
    if (prevCert != null) {
      $result.prevCert = prevCert;
    }
    if (txs != null) {
      $result.txs.addAll(txs);
    }
    return $result;
  }
  GetBlockResponse._() : super();
  factory GetBlockResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetBlockResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetBlockResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..a<$core.int>(1, _omitFieldNames ? '' : 'height', $pb.PbFieldType.OU3)
    ..aOS(2, _omitFieldNames ? '' : 'hash')
    ..aOS(3, _omitFieldNames ? '' : 'data')
    ..a<$core.int>(4, _omitFieldNames ? '' : 'blockTime', $pb.PbFieldType.OU3)
    ..aOM<BlockHeaderInfo>(5, _omitFieldNames ? '' : 'header', subBuilder: BlockHeaderInfo.create)
    ..aOM<CertificateInfo>(6, _omitFieldNames ? '' : 'prevCert', subBuilder: CertificateInfo.create)
    ..pc<$0.TransactionInfo>(7, _omitFieldNames ? '' : 'txs', $pb.PbFieldType.PM, subBuilder: $0.TransactionInfo.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetBlockResponse clone() => GetBlockResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetBlockResponse copyWith(void Function(GetBlockResponse) updates) => super.copyWith((message) => updates(message as GetBlockResponse)) as GetBlockResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetBlockResponse create() => GetBlockResponse._();
  GetBlockResponse createEmptyInstance() => create();
  static $pb.PbList<GetBlockResponse> createRepeated() => $pb.PbList<GetBlockResponse>();
  @$core.pragma('dart2js:noInline')
  static GetBlockResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetBlockResponse>(create);
  static GetBlockResponse? _defaultInstance;

  /// The height of the block.
  @$pb.TagNumber(1)
  $core.int get height => $_getIZ(0);
  @$pb.TagNumber(1)
  set height($core.int v) { $_setUnsignedInt32(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasHeight() => $_has(0);
  @$pb.TagNumber(1)
  void clearHeight() => $_clearField(1);

  /// The hash of the block.
  @$pb.TagNumber(2)
  $core.String get hash => $_getSZ(1);
  @$pb.TagNumber(2)
  set hash($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasHash() => $_has(1);
  @$pb.TagNumber(2)
  void clearHash() => $_clearField(2);

  /// Block data, available only if verbosity level is set to BLOCK_DATA.
  @$pb.TagNumber(3)
  $core.String get data => $_getSZ(2);
  @$pb.TagNumber(3)
  set data($core.String v) { $_setString(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasData() => $_has(2);
  @$pb.TagNumber(3)
  void clearData() => $_clearField(3);

  /// The timestamp of the block.
  @$pb.TagNumber(4)
  $core.int get blockTime => $_getIZ(3);
  @$pb.TagNumber(4)
  set blockTime($core.int v) { $_setUnsignedInt32(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasBlockTime() => $_has(3);
  @$pb.TagNumber(4)
  void clearBlockTime() => $_clearField(4);

  /// Header information of the block.
  @$pb.TagNumber(5)
  BlockHeaderInfo get header => $_getN(4);
  @$pb.TagNumber(5)
  set header(BlockHeaderInfo v) { $_setField(5, v); }
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
  set prevCert(CertificateInfo v) { $_setField(6, v); }
  @$pb.TagNumber(6)
  $core.bool hasPrevCert() => $_has(5);
  @$pb.TagNumber(6)
  void clearPrevCert() => $_clearField(6);
  @$pb.TagNumber(6)
  CertificateInfo ensurePrevCert() => $_ensure(5);

  /// List of transactions in the block, available when verbosity level is set to
  /// BLOCK_TRANSACTIONS.
  @$pb.TagNumber(7)
  $pb.PbList<$0.TransactionInfo> get txs => $_getList(6);
}

/// Request message for retrieving block hash by height.
class GetBlockHashRequest extends $pb.GeneratedMessage {
  factory GetBlockHashRequest({
    $core.int? height,
  }) {
    final $result = create();
    if (height != null) {
      $result.height = height;
    }
    return $result;
  }
  GetBlockHashRequest._() : super();
  factory GetBlockHashRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetBlockHashRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetBlockHashRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..a<$core.int>(1, _omitFieldNames ? '' : 'height', $pb.PbFieldType.OU3)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetBlockHashRequest clone() => GetBlockHashRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetBlockHashRequest copyWith(void Function(GetBlockHashRequest) updates) => super.copyWith((message) => updates(message as GetBlockHashRequest)) as GetBlockHashRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetBlockHashRequest create() => GetBlockHashRequest._();
  GetBlockHashRequest createEmptyInstance() => create();
  static $pb.PbList<GetBlockHashRequest> createRepeated() => $pb.PbList<GetBlockHashRequest>();
  @$core.pragma('dart2js:noInline')
  static GetBlockHashRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetBlockHashRequest>(create);
  static GetBlockHashRequest? _defaultInstance;

  /// The height of the block to retrieve the hash for.
  @$pb.TagNumber(1)
  $core.int get height => $_getIZ(0);
  @$pb.TagNumber(1)
  set height($core.int v) { $_setUnsignedInt32(0, v); }
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
    final $result = create();
    if (hash != null) {
      $result.hash = hash;
    }
    return $result;
  }
  GetBlockHashResponse._() : super();
  factory GetBlockHashResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetBlockHashResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetBlockHashResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'hash')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetBlockHashResponse clone() => GetBlockHashResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetBlockHashResponse copyWith(void Function(GetBlockHashResponse) updates) => super.copyWith((message) => updates(message as GetBlockHashResponse)) as GetBlockHashResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetBlockHashResponse create() => GetBlockHashResponse._();
  GetBlockHashResponse createEmptyInstance() => create();
  static $pb.PbList<GetBlockHashResponse> createRepeated() => $pb.PbList<GetBlockHashResponse>();
  @$core.pragma('dart2js:noInline')
  static GetBlockHashResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetBlockHashResponse>(create);
  static GetBlockHashResponse? _defaultInstance;

  /// The hash of the block.
  @$pb.TagNumber(1)
  $core.String get hash => $_getSZ(0);
  @$pb.TagNumber(1)
  set hash($core.String v) { $_setString(0, v); }
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
    final $result = create();
    if (hash != null) {
      $result.hash = hash;
    }
    return $result;
  }
  GetBlockHeightRequest._() : super();
  factory GetBlockHeightRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetBlockHeightRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetBlockHeightRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'hash')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetBlockHeightRequest clone() => GetBlockHeightRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetBlockHeightRequest copyWith(void Function(GetBlockHeightRequest) updates) => super.copyWith((message) => updates(message as GetBlockHeightRequest)) as GetBlockHeightRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetBlockHeightRequest create() => GetBlockHeightRequest._();
  GetBlockHeightRequest createEmptyInstance() => create();
  static $pb.PbList<GetBlockHeightRequest> createRepeated() => $pb.PbList<GetBlockHeightRequest>();
  @$core.pragma('dart2js:noInline')
  static GetBlockHeightRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetBlockHeightRequest>(create);
  static GetBlockHeightRequest? _defaultInstance;

  /// The hash of the block to retrieve the height for.
  @$pb.TagNumber(1)
  $core.String get hash => $_getSZ(0);
  @$pb.TagNumber(1)
  set hash($core.String v) { $_setString(0, v); }
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
    final $result = create();
    if (height != null) {
      $result.height = height;
    }
    return $result;
  }
  GetBlockHeightResponse._() : super();
  factory GetBlockHeightResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetBlockHeightResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetBlockHeightResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..a<$core.int>(1, _omitFieldNames ? '' : 'height', $pb.PbFieldType.OU3)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetBlockHeightResponse clone() => GetBlockHeightResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetBlockHeightResponse copyWith(void Function(GetBlockHeightResponse) updates) => super.copyWith((message) => updates(message as GetBlockHeightResponse)) as GetBlockHeightResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetBlockHeightResponse create() => GetBlockHeightResponse._();
  GetBlockHeightResponse createEmptyInstance() => create();
  static $pb.PbList<GetBlockHeightResponse> createRepeated() => $pb.PbList<GetBlockHeightResponse>();
  @$core.pragma('dart2js:noInline')
  static GetBlockHeightResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetBlockHeightResponse>(create);
  static GetBlockHeightResponse? _defaultInstance;

  /// The height of the block.
  @$pb.TagNumber(1)
  $core.int get height => $_getIZ(0);
  @$pb.TagNumber(1)
  set height($core.int v) { $_setUnsignedInt32(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasHeight() => $_has(0);
  @$pb.TagNumber(1)
  void clearHeight() => $_clearField(1);
}

/// Request message for retrieving blockchain information.
class GetBlockchainInfoRequest extends $pb.GeneratedMessage {
  factory GetBlockchainInfoRequest() => create();
  GetBlockchainInfoRequest._() : super();
  factory GetBlockchainInfoRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetBlockchainInfoRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetBlockchainInfoRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetBlockchainInfoRequest clone() => GetBlockchainInfoRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetBlockchainInfoRequest copyWith(void Function(GetBlockchainInfoRequest) updates) => super.copyWith((message) => updates(message as GetBlockchainInfoRequest)) as GetBlockchainInfoRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetBlockchainInfoRequest create() => GetBlockchainInfoRequest._();
  GetBlockchainInfoRequest createEmptyInstance() => create();
  static $pb.PbList<GetBlockchainInfoRequest> createRepeated() => $pb.PbList<GetBlockchainInfoRequest>();
  @$core.pragma('dart2js:noInline')
  static GetBlockchainInfoRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetBlockchainInfoRequest>(create);
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
    $core.Iterable<ValidatorInfo>? committeeValidators,
    $core.bool? isPruned,
    $core.int? pruningHeight,
    $fixnum.Int64? lastBlockTime,
  }) {
    final $result = create();
    if (lastBlockHeight != null) {
      $result.lastBlockHeight = lastBlockHeight;
    }
    if (lastBlockHash != null) {
      $result.lastBlockHash = lastBlockHash;
    }
    if (totalAccounts != null) {
      $result.totalAccounts = totalAccounts;
    }
    if (totalValidators != null) {
      $result.totalValidators = totalValidators;
    }
    if (totalPower != null) {
      $result.totalPower = totalPower;
    }
    if (committeePower != null) {
      $result.committeePower = committeePower;
    }
    if (committeeValidators != null) {
      $result.committeeValidators.addAll(committeeValidators);
    }
    if (isPruned != null) {
      $result.isPruned = isPruned;
    }
    if (pruningHeight != null) {
      $result.pruningHeight = pruningHeight;
    }
    if (lastBlockTime != null) {
      $result.lastBlockTime = lastBlockTime;
    }
    return $result;
  }
  GetBlockchainInfoResponse._() : super();
  factory GetBlockchainInfoResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetBlockchainInfoResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetBlockchainInfoResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..a<$core.int>(1, _omitFieldNames ? '' : 'lastBlockHeight', $pb.PbFieldType.OU3)
    ..aOS(2, _omitFieldNames ? '' : 'lastBlockHash')
    ..a<$core.int>(3, _omitFieldNames ? '' : 'totalAccounts', $pb.PbFieldType.O3)
    ..a<$core.int>(4, _omitFieldNames ? '' : 'totalValidators', $pb.PbFieldType.O3)
    ..aInt64(5, _omitFieldNames ? '' : 'totalPower')
    ..aInt64(6, _omitFieldNames ? '' : 'committeePower')
    ..pc<ValidatorInfo>(7, _omitFieldNames ? '' : 'committeeValidators', $pb.PbFieldType.PM, subBuilder: ValidatorInfo.create)
    ..aOB(8, _omitFieldNames ? '' : 'isPruned')
    ..a<$core.int>(9, _omitFieldNames ? '' : 'pruningHeight', $pb.PbFieldType.OU3)
    ..aInt64(10, _omitFieldNames ? '' : 'lastBlockTime')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetBlockchainInfoResponse clone() => GetBlockchainInfoResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetBlockchainInfoResponse copyWith(void Function(GetBlockchainInfoResponse) updates) => super.copyWith((message) => updates(message as GetBlockchainInfoResponse)) as GetBlockchainInfoResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetBlockchainInfoResponse create() => GetBlockchainInfoResponse._();
  GetBlockchainInfoResponse createEmptyInstance() => create();
  static $pb.PbList<GetBlockchainInfoResponse> createRepeated() => $pb.PbList<GetBlockchainInfoResponse>();
  @$core.pragma('dart2js:noInline')
  static GetBlockchainInfoResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetBlockchainInfoResponse>(create);
  static GetBlockchainInfoResponse? _defaultInstance;

  /// The height of the last block in the blockchain.
  @$pb.TagNumber(1)
  $core.int get lastBlockHeight => $_getIZ(0);
  @$pb.TagNumber(1)
  set lastBlockHeight($core.int v) { $_setUnsignedInt32(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasLastBlockHeight() => $_has(0);
  @$pb.TagNumber(1)
  void clearLastBlockHeight() => $_clearField(1);

  /// The hash of the last block in the blockchain.
  @$pb.TagNumber(2)
  $core.String get lastBlockHash => $_getSZ(1);
  @$pb.TagNumber(2)
  set lastBlockHash($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasLastBlockHash() => $_has(1);
  @$pb.TagNumber(2)
  void clearLastBlockHash() => $_clearField(2);

  /// The total number of accounts in the blockchain.
  @$pb.TagNumber(3)
  $core.int get totalAccounts => $_getIZ(2);
  @$pb.TagNumber(3)
  set totalAccounts($core.int v) { $_setSignedInt32(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasTotalAccounts() => $_has(2);
  @$pb.TagNumber(3)
  void clearTotalAccounts() => $_clearField(3);

  /// The total number of validators in the blockchain.
  @$pb.TagNumber(4)
  $core.int get totalValidators => $_getIZ(3);
  @$pb.TagNumber(4)
  set totalValidators($core.int v) { $_setSignedInt32(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasTotalValidators() => $_has(3);
  @$pb.TagNumber(4)
  void clearTotalValidators() => $_clearField(4);

  /// The total power of the blockchain.
  @$pb.TagNumber(5)
  $fixnum.Int64 get totalPower => $_getI64(4);
  @$pb.TagNumber(5)
  set totalPower($fixnum.Int64 v) { $_setInt64(4, v); }
  @$pb.TagNumber(5)
  $core.bool hasTotalPower() => $_has(4);
  @$pb.TagNumber(5)
  void clearTotalPower() => $_clearField(5);

  /// The power of the committee.
  @$pb.TagNumber(6)
  $fixnum.Int64 get committeePower => $_getI64(5);
  @$pb.TagNumber(6)
  set committeePower($fixnum.Int64 v) { $_setInt64(5, v); }
  @$pb.TagNumber(6)
  $core.bool hasCommitteePower() => $_has(5);
  @$pb.TagNumber(6)
  void clearCommitteePower() => $_clearField(6);

  /// List of committee validators.
  @$pb.TagNumber(7)
  $pb.PbList<ValidatorInfo> get committeeValidators => $_getList(6);

  /// If the blocks are subject to pruning.
  @$pb.TagNumber(8)
  $core.bool get isPruned => $_getBF(7);
  @$pb.TagNumber(8)
  set isPruned($core.bool v) { $_setBool(7, v); }
  @$pb.TagNumber(8)
  $core.bool hasIsPruned() => $_has(7);
  @$pb.TagNumber(8)
  void clearIsPruned() => $_clearField(8);

  /// Lowest-height block stored (only present if pruning is enabled)
  @$pb.TagNumber(9)
  $core.int get pruningHeight => $_getIZ(8);
  @$pb.TagNumber(9)
  set pruningHeight($core.int v) { $_setUnsignedInt32(8, v); }
  @$pb.TagNumber(9)
  $core.bool hasPruningHeight() => $_has(8);
  @$pb.TagNumber(9)
  void clearPruningHeight() => $_clearField(9);

  /// Timestamp of the last block in Unix format
  @$pb.TagNumber(10)
  $fixnum.Int64 get lastBlockTime => $_getI64(9);
  @$pb.TagNumber(10)
  set lastBlockTime($fixnum.Int64 v) { $_setInt64(9, v); }
  @$pb.TagNumber(10)
  $core.bool hasLastBlockTime() => $_has(9);
  @$pb.TagNumber(10)
  void clearLastBlockTime() => $_clearField(10);
}

/// Request message for retrieving consensus information.
class GetConsensusInfoRequest extends $pb.GeneratedMessage {
  factory GetConsensusInfoRequest() => create();
  GetConsensusInfoRequest._() : super();
  factory GetConsensusInfoRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetConsensusInfoRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetConsensusInfoRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetConsensusInfoRequest clone() => GetConsensusInfoRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetConsensusInfoRequest copyWith(void Function(GetConsensusInfoRequest) updates) => super.copyWith((message) => updates(message as GetConsensusInfoRequest)) as GetConsensusInfoRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetConsensusInfoRequest create() => GetConsensusInfoRequest._();
  GetConsensusInfoRequest createEmptyInstance() => create();
  static $pb.PbList<GetConsensusInfoRequest> createRepeated() => $pb.PbList<GetConsensusInfoRequest>();
  @$core.pragma('dart2js:noInline')
  static GetConsensusInfoRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetConsensusInfoRequest>(create);
  static GetConsensusInfoRequest? _defaultInstance;
}

/// Response message contains consensus information.
class GetConsensusInfoResponse extends $pb.GeneratedMessage {
  factory GetConsensusInfoResponse({
    ProposalInfo? proposal,
    $core.Iterable<ConsensusInfo>? instances,
  }) {
    final $result = create();
    if (proposal != null) {
      $result.proposal = proposal;
    }
    if (instances != null) {
      $result.instances.addAll(instances);
    }
    return $result;
  }
  GetConsensusInfoResponse._() : super();
  factory GetConsensusInfoResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetConsensusInfoResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetConsensusInfoResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..aOM<ProposalInfo>(1, _omitFieldNames ? '' : 'proposal', subBuilder: ProposalInfo.create)
    ..pc<ConsensusInfo>(2, _omitFieldNames ? '' : 'instances', $pb.PbFieldType.PM, subBuilder: ConsensusInfo.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetConsensusInfoResponse clone() => GetConsensusInfoResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetConsensusInfoResponse copyWith(void Function(GetConsensusInfoResponse) updates) => super.copyWith((message) => updates(message as GetConsensusInfoResponse)) as GetConsensusInfoResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetConsensusInfoResponse create() => GetConsensusInfoResponse._();
  GetConsensusInfoResponse createEmptyInstance() => create();
  static $pb.PbList<GetConsensusInfoResponse> createRepeated() => $pb.PbList<GetConsensusInfoResponse>();
  @$core.pragma('dart2js:noInline')
  static GetConsensusInfoResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetConsensusInfoResponse>(create);
  static GetConsensusInfoResponse? _defaultInstance;

  /// The proposal of the consensus info.
  @$pb.TagNumber(1)
  ProposalInfo get proposal => $_getN(0);
  @$pb.TagNumber(1)
  set proposal(ProposalInfo v) { $_setField(1, v); }
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
    final $result = create();
    if (payloadType != null) {
      $result.payloadType = payloadType;
    }
    return $result;
  }
  GetTxPoolContentRequest._() : super();
  factory GetTxPoolContentRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetTxPoolContentRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetTxPoolContentRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..e<$0.PayloadType>(1, _omitFieldNames ? '' : 'payloadType', $pb.PbFieldType.OE, defaultOrMaker: $0.PayloadType.PAYLOAD_TYPE_UNSPECIFIED, valueOf: $0.PayloadType.valueOf, enumValues: $0.PayloadType.values)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetTxPoolContentRequest clone() => GetTxPoolContentRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetTxPoolContentRequest copyWith(void Function(GetTxPoolContentRequest) updates) => super.copyWith((message) => updates(message as GetTxPoolContentRequest)) as GetTxPoolContentRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetTxPoolContentRequest create() => GetTxPoolContentRequest._();
  GetTxPoolContentRequest createEmptyInstance() => create();
  static $pb.PbList<GetTxPoolContentRequest> createRepeated() => $pb.PbList<GetTxPoolContentRequest>();
  @$core.pragma('dart2js:noInline')
  static GetTxPoolContentRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetTxPoolContentRequest>(create);
  static GetTxPoolContentRequest? _defaultInstance;

  /// The type of transactions to retrieve from the transaction pool. 0 means all types.
  @$pb.TagNumber(1)
  $0.PayloadType get payloadType => $_getN(0);
  @$pb.TagNumber(1)
  set payloadType($0.PayloadType v) { $_setField(1, v); }
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
    final $result = create();
    if (txs != null) {
      $result.txs.addAll(txs);
    }
    return $result;
  }
  GetTxPoolContentResponse._() : super();
  factory GetTxPoolContentResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetTxPoolContentResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetTxPoolContentResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..pc<$0.TransactionInfo>(1, _omitFieldNames ? '' : 'txs', $pb.PbFieldType.PM, subBuilder: $0.TransactionInfo.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetTxPoolContentResponse clone() => GetTxPoolContentResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetTxPoolContentResponse copyWith(void Function(GetTxPoolContentResponse) updates) => super.copyWith((message) => updates(message as GetTxPoolContentResponse)) as GetTxPoolContentResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetTxPoolContentResponse create() => GetTxPoolContentResponse._();
  GetTxPoolContentResponse createEmptyInstance() => create();
  static $pb.PbList<GetTxPoolContentResponse> createRepeated() => $pb.PbList<GetTxPoolContentResponse>();
  @$core.pragma('dart2js:noInline')
  static GetTxPoolContentResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetTxPoolContentResponse>(create);
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
    final $result = create();
    if (hash != null) {
      $result.hash = hash;
    }
    if (data != null) {
      $result.data = data;
    }
    if (publicKey != null) {
      $result.publicKey = publicKey;
    }
    if (number != null) {
      $result.number = number;
    }
    if (stake != null) {
      $result.stake = stake;
    }
    if (lastBondingHeight != null) {
      $result.lastBondingHeight = lastBondingHeight;
    }
    if (lastSortitionHeight != null) {
      $result.lastSortitionHeight = lastSortitionHeight;
    }
    if (unbondingHeight != null) {
      $result.unbondingHeight = unbondingHeight;
    }
    if (address != null) {
      $result.address = address;
    }
    if (availabilityScore != null) {
      $result.availabilityScore = availabilityScore;
    }
    if (protocolVersion != null) {
      $result.protocolVersion = protocolVersion;
    }
    return $result;
  }
  ValidatorInfo._() : super();
  factory ValidatorInfo.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory ValidatorInfo.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'ValidatorInfo', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'hash')
    ..aOS(2, _omitFieldNames ? '' : 'data')
    ..aOS(3, _omitFieldNames ? '' : 'publicKey')
    ..a<$core.int>(4, _omitFieldNames ? '' : 'number', $pb.PbFieldType.O3)
    ..aInt64(5, _omitFieldNames ? '' : 'stake')
    ..a<$core.int>(6, _omitFieldNames ? '' : 'lastBondingHeight', $pb.PbFieldType.OU3)
    ..a<$core.int>(7, _omitFieldNames ? '' : 'lastSortitionHeight', $pb.PbFieldType.OU3)
    ..a<$core.int>(8, _omitFieldNames ? '' : 'unbondingHeight', $pb.PbFieldType.OU3)
    ..aOS(9, _omitFieldNames ? '' : 'address')
    ..a<$core.double>(10, _omitFieldNames ? '' : 'availabilityScore', $pb.PbFieldType.OD)
    ..a<$core.int>(11, _omitFieldNames ? '' : 'protocolVersion', $pb.PbFieldType.O3)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  ValidatorInfo clone() => ValidatorInfo()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  ValidatorInfo copyWith(void Function(ValidatorInfo) updates) => super.copyWith((message) => updates(message as ValidatorInfo)) as ValidatorInfo;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ValidatorInfo create() => ValidatorInfo._();
  ValidatorInfo createEmptyInstance() => create();
  static $pb.PbList<ValidatorInfo> createRepeated() => $pb.PbList<ValidatorInfo>();
  @$core.pragma('dart2js:noInline')
  static ValidatorInfo getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ValidatorInfo>(create);
  static ValidatorInfo? _defaultInstance;

  /// The hash of the validator.
  @$pb.TagNumber(1)
  $core.String get hash => $_getSZ(0);
  @$pb.TagNumber(1)
  set hash($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasHash() => $_has(0);
  @$pb.TagNumber(1)
  void clearHash() => $_clearField(1);

  /// The serialized data of the validator.
  @$pb.TagNumber(2)
  $core.String get data => $_getSZ(1);
  @$pb.TagNumber(2)
  set data($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasData() => $_has(1);
  @$pb.TagNumber(2)
  void clearData() => $_clearField(2);

  /// The public key of the validator.
  @$pb.TagNumber(3)
  $core.String get publicKey => $_getSZ(2);
  @$pb.TagNumber(3)
  set publicKey($core.String v) { $_setString(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasPublicKey() => $_has(2);
  @$pb.TagNumber(3)
  void clearPublicKey() => $_clearField(3);

  /// The unique number assigned to the validator.
  @$pb.TagNumber(4)
  $core.int get number => $_getIZ(3);
  @$pb.TagNumber(4)
  set number($core.int v) { $_setSignedInt32(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasNumber() => $_has(3);
  @$pb.TagNumber(4)
  void clearNumber() => $_clearField(4);

  /// The stake of the validator in NanoPAC.
  @$pb.TagNumber(5)
  $fixnum.Int64 get stake => $_getI64(4);
  @$pb.TagNumber(5)
  set stake($fixnum.Int64 v) { $_setInt64(4, v); }
  @$pb.TagNumber(5)
  $core.bool hasStake() => $_has(4);
  @$pb.TagNumber(5)
  void clearStake() => $_clearField(5);

  /// The height at which the validator last bonded.
  @$pb.TagNumber(6)
  $core.int get lastBondingHeight => $_getIZ(5);
  @$pb.TagNumber(6)
  set lastBondingHeight($core.int v) { $_setUnsignedInt32(5, v); }
  @$pb.TagNumber(6)
  $core.bool hasLastBondingHeight() => $_has(5);
  @$pb.TagNumber(6)
  void clearLastBondingHeight() => $_clearField(6);

  /// The height at which the validator last participated in sortition.
  @$pb.TagNumber(7)
  $core.int get lastSortitionHeight => $_getIZ(6);
  @$pb.TagNumber(7)
  set lastSortitionHeight($core.int v) { $_setUnsignedInt32(6, v); }
  @$pb.TagNumber(7)
  $core.bool hasLastSortitionHeight() => $_has(6);
  @$pb.TagNumber(7)
  void clearLastSortitionHeight() => $_clearField(7);

  /// The height at which the validator will unbond.
  @$pb.TagNumber(8)
  $core.int get unbondingHeight => $_getIZ(7);
  @$pb.TagNumber(8)
  set unbondingHeight($core.int v) { $_setUnsignedInt32(7, v); }
  @$pb.TagNumber(8)
  $core.bool hasUnbondingHeight() => $_has(7);
  @$pb.TagNumber(8)
  void clearUnbondingHeight() => $_clearField(8);

  /// The address of the validator.
  @$pb.TagNumber(9)
  $core.String get address => $_getSZ(8);
  @$pb.TagNumber(9)
  set address($core.String v) { $_setString(8, v); }
  @$pb.TagNumber(9)
  $core.bool hasAddress() => $_has(8);
  @$pb.TagNumber(9)
  void clearAddress() => $_clearField(9);

  /// The availability score of the validator.
  @$pb.TagNumber(10)
  $core.double get availabilityScore => $_getN(9);
  @$pb.TagNumber(10)
  set availabilityScore($core.double v) { $_setDouble(9, v); }
  @$pb.TagNumber(10)
  $core.bool hasAvailabilityScore() => $_has(9);
  @$pb.TagNumber(10)
  void clearAvailabilityScore() => $_clearField(10);

  /// The protocol version of the validator.
  @$pb.TagNumber(11)
  $core.int get protocolVersion => $_getIZ(10);
  @$pb.TagNumber(11)
  set protocolVersion($core.int v) { $_setSignedInt32(10, v); }
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
    final $result = create();
    if (hash != null) {
      $result.hash = hash;
    }
    if (data != null) {
      $result.data = data;
    }
    if (number != null) {
      $result.number = number;
    }
    if (balance != null) {
      $result.balance = balance;
    }
    if (address != null) {
      $result.address = address;
    }
    return $result;
  }
  AccountInfo._() : super();
  factory AccountInfo.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory AccountInfo.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'AccountInfo', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'hash')
    ..aOS(2, _omitFieldNames ? '' : 'data')
    ..a<$core.int>(3, _omitFieldNames ? '' : 'number', $pb.PbFieldType.O3)
    ..aInt64(4, _omitFieldNames ? '' : 'balance')
    ..aOS(5, _omitFieldNames ? '' : 'address')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  AccountInfo clone() => AccountInfo()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  AccountInfo copyWith(void Function(AccountInfo) updates) => super.copyWith((message) => updates(message as AccountInfo)) as AccountInfo;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static AccountInfo create() => AccountInfo._();
  AccountInfo createEmptyInstance() => create();
  static $pb.PbList<AccountInfo> createRepeated() => $pb.PbList<AccountInfo>();
  @$core.pragma('dart2js:noInline')
  static AccountInfo getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<AccountInfo>(create);
  static AccountInfo? _defaultInstance;

  /// The hash of the account.
  @$pb.TagNumber(1)
  $core.String get hash => $_getSZ(0);
  @$pb.TagNumber(1)
  set hash($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasHash() => $_has(0);
  @$pb.TagNumber(1)
  void clearHash() => $_clearField(1);

  /// The serialized data of the account.
  @$pb.TagNumber(2)
  $core.String get data => $_getSZ(1);
  @$pb.TagNumber(2)
  set data($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasData() => $_has(1);
  @$pb.TagNumber(2)
  void clearData() => $_clearField(2);

  /// The unique number assigned to the account.
  @$pb.TagNumber(3)
  $core.int get number => $_getIZ(2);
  @$pb.TagNumber(3)
  set number($core.int v) { $_setSignedInt32(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasNumber() => $_has(2);
  @$pb.TagNumber(3)
  void clearNumber() => $_clearField(3);

  /// The balance of the account in NanoPAC.
  @$pb.TagNumber(4)
  $fixnum.Int64 get balance => $_getI64(3);
  @$pb.TagNumber(4)
  set balance($fixnum.Int64 v) { $_setInt64(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasBalance() => $_has(3);
  @$pb.TagNumber(4)
  void clearBalance() => $_clearField(4);

  /// The address of the account.
  @$pb.TagNumber(5)
  $core.String get address => $_getSZ(4);
  @$pb.TagNumber(5)
  set address($core.String v) { $_setString(4, v); }
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
    final $result = create();
    if (version != null) {
      $result.version = version;
    }
    if (prevBlockHash != null) {
      $result.prevBlockHash = prevBlockHash;
    }
    if (stateRoot != null) {
      $result.stateRoot = stateRoot;
    }
    if (sortitionSeed != null) {
      $result.sortitionSeed = sortitionSeed;
    }
    if (proposerAddress != null) {
      $result.proposerAddress = proposerAddress;
    }
    return $result;
  }
  BlockHeaderInfo._() : super();
  factory BlockHeaderInfo.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory BlockHeaderInfo.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'BlockHeaderInfo', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..a<$core.int>(1, _omitFieldNames ? '' : 'version', $pb.PbFieldType.O3)
    ..aOS(2, _omitFieldNames ? '' : 'prevBlockHash')
    ..aOS(3, _omitFieldNames ? '' : 'stateRoot')
    ..aOS(4, _omitFieldNames ? '' : 'sortitionSeed')
    ..aOS(5, _omitFieldNames ? '' : 'proposerAddress')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  BlockHeaderInfo clone() => BlockHeaderInfo()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  BlockHeaderInfo copyWith(void Function(BlockHeaderInfo) updates) => super.copyWith((message) => updates(message as BlockHeaderInfo)) as BlockHeaderInfo;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static BlockHeaderInfo create() => BlockHeaderInfo._();
  BlockHeaderInfo createEmptyInstance() => create();
  static $pb.PbList<BlockHeaderInfo> createRepeated() => $pb.PbList<BlockHeaderInfo>();
  @$core.pragma('dart2js:noInline')
  static BlockHeaderInfo getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<BlockHeaderInfo>(create);
  static BlockHeaderInfo? _defaultInstance;

  /// The version of the block.
  @$pb.TagNumber(1)
  $core.int get version => $_getIZ(0);
  @$pb.TagNumber(1)
  set version($core.int v) { $_setSignedInt32(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasVersion() => $_has(0);
  @$pb.TagNumber(1)
  void clearVersion() => $_clearField(1);

  /// The hash of the previous block.
  @$pb.TagNumber(2)
  $core.String get prevBlockHash => $_getSZ(1);
  @$pb.TagNumber(2)
  set prevBlockHash($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasPrevBlockHash() => $_has(1);
  @$pb.TagNumber(2)
  void clearPrevBlockHash() => $_clearField(2);

  /// The state root hash of the blockchain.
  @$pb.TagNumber(3)
  $core.String get stateRoot => $_getSZ(2);
  @$pb.TagNumber(3)
  set stateRoot($core.String v) { $_setString(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasStateRoot() => $_has(2);
  @$pb.TagNumber(3)
  void clearStateRoot() => $_clearField(3);

  /// The sortition seed of the block.
  @$pb.TagNumber(4)
  $core.String get sortitionSeed => $_getSZ(3);
  @$pb.TagNumber(4)
  set sortitionSeed($core.String v) { $_setString(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasSortitionSeed() => $_has(3);
  @$pb.TagNumber(4)
  void clearSortitionSeed() => $_clearField(4);

  /// The address of the proposer of the block.
  @$pb.TagNumber(5)
  $core.String get proposerAddress => $_getSZ(4);
  @$pb.TagNumber(5)
  set proposerAddress($core.String v) { $_setString(4, v); }
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
    final $result = create();
    if (hash != null) {
      $result.hash = hash;
    }
    if (round != null) {
      $result.round = round;
    }
    if (committers != null) {
      $result.committers.addAll(committers);
    }
    if (absentees != null) {
      $result.absentees.addAll(absentees);
    }
    if (signature != null) {
      $result.signature = signature;
    }
    return $result;
  }
  CertificateInfo._() : super();
  factory CertificateInfo.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory CertificateInfo.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'CertificateInfo', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'hash')
    ..a<$core.int>(2, _omitFieldNames ? '' : 'round', $pb.PbFieldType.O3)
    ..p<$core.int>(3, _omitFieldNames ? '' : 'committers', $pb.PbFieldType.K3)
    ..p<$core.int>(4, _omitFieldNames ? '' : 'absentees', $pb.PbFieldType.K3)
    ..aOS(5, _omitFieldNames ? '' : 'signature')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  CertificateInfo clone() => CertificateInfo()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  CertificateInfo copyWith(void Function(CertificateInfo) updates) => super.copyWith((message) => updates(message as CertificateInfo)) as CertificateInfo;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CertificateInfo create() => CertificateInfo._();
  CertificateInfo createEmptyInstance() => create();
  static $pb.PbList<CertificateInfo> createRepeated() => $pb.PbList<CertificateInfo>();
  @$core.pragma('dart2js:noInline')
  static CertificateInfo getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<CertificateInfo>(create);
  static CertificateInfo? _defaultInstance;

  /// The hash of the certificate.
  @$pb.TagNumber(1)
  $core.String get hash => $_getSZ(0);
  @$pb.TagNumber(1)
  set hash($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasHash() => $_has(0);
  @$pb.TagNumber(1)
  void clearHash() => $_clearField(1);

  /// The round of the certificate.
  @$pb.TagNumber(2)
  $core.int get round => $_getIZ(1);
  @$pb.TagNumber(2)
  set round($core.int v) { $_setSignedInt32(1, v); }
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
  set signature($core.String v) { $_setString(4, v); }
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
    final $result = create();
    if (type != null) {
      $result.type = type;
    }
    if (voter != null) {
      $result.voter = voter;
    }
    if (blockHash != null) {
      $result.blockHash = blockHash;
    }
    if (round != null) {
      $result.round = round;
    }
    if (cpRound != null) {
      $result.cpRound = cpRound;
    }
    if (cpValue != null) {
      $result.cpValue = cpValue;
    }
    return $result;
  }
  VoteInfo._() : super();
  factory VoteInfo.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory VoteInfo.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'VoteInfo', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..e<VoteType>(1, _omitFieldNames ? '' : 'type', $pb.PbFieldType.OE, defaultOrMaker: VoteType.VOTE_TYPE_UNSPECIFIED, valueOf: VoteType.valueOf, enumValues: VoteType.values)
    ..aOS(2, _omitFieldNames ? '' : 'voter')
    ..aOS(3, _omitFieldNames ? '' : 'blockHash')
    ..a<$core.int>(4, _omitFieldNames ? '' : 'round', $pb.PbFieldType.O3)
    ..a<$core.int>(5, _omitFieldNames ? '' : 'cpRound', $pb.PbFieldType.O3)
    ..a<$core.int>(6, _omitFieldNames ? '' : 'cpValue', $pb.PbFieldType.O3)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  VoteInfo clone() => VoteInfo()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  VoteInfo copyWith(void Function(VoteInfo) updates) => super.copyWith((message) => updates(message as VoteInfo)) as VoteInfo;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static VoteInfo create() => VoteInfo._();
  VoteInfo createEmptyInstance() => create();
  static $pb.PbList<VoteInfo> createRepeated() => $pb.PbList<VoteInfo>();
  @$core.pragma('dart2js:noInline')
  static VoteInfo getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<VoteInfo>(create);
  static VoteInfo? _defaultInstance;

  /// The type of the vote.
  @$pb.TagNumber(1)
  VoteType get type => $_getN(0);
  @$pb.TagNumber(1)
  set type(VoteType v) { $_setField(1, v); }
  @$pb.TagNumber(1)
  $core.bool hasType() => $_has(0);
  @$pb.TagNumber(1)
  void clearType() => $_clearField(1);

  /// The address of the voter.
  @$pb.TagNumber(2)
  $core.String get voter => $_getSZ(1);
  @$pb.TagNumber(2)
  set voter($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasVoter() => $_has(1);
  @$pb.TagNumber(2)
  void clearVoter() => $_clearField(2);

  /// The hash of the block being voted on.
  @$pb.TagNumber(3)
  $core.String get blockHash => $_getSZ(2);
  @$pb.TagNumber(3)
  set blockHash($core.String v) { $_setString(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasBlockHash() => $_has(2);
  @$pb.TagNumber(3)
  void clearBlockHash() => $_clearField(3);

  /// The consensus round of the vote.
  @$pb.TagNumber(4)
  $core.int get round => $_getIZ(3);
  @$pb.TagNumber(4)
  set round($core.int v) { $_setSignedInt32(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasRound() => $_has(3);
  @$pb.TagNumber(4)
  void clearRound() => $_clearField(4);

  /// The change-proposer round of the vote.
  @$pb.TagNumber(5)
  $core.int get cpRound => $_getIZ(4);
  @$pb.TagNumber(5)
  set cpRound($core.int v) { $_setSignedInt32(4, v); }
  @$pb.TagNumber(5)
  $core.bool hasCpRound() => $_has(4);
  @$pb.TagNumber(5)
  void clearCpRound() => $_clearField(5);

  /// The change-proposer value of the vote.
  @$pb.TagNumber(6)
  $core.int get cpValue => $_getIZ(5);
  @$pb.TagNumber(6)
  set cpValue($core.int v) { $_setSignedInt32(5, v); }
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
    final $result = create();
    if (address != null) {
      $result.address = address;
    }
    if (active != null) {
      $result.active = active;
    }
    if (height != null) {
      $result.height = height;
    }
    if (round != null) {
      $result.round = round;
    }
    if (votes != null) {
      $result.votes.addAll(votes);
    }
    return $result;
  }
  ConsensusInfo._() : super();
  factory ConsensusInfo.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory ConsensusInfo.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'ConsensusInfo', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'address')
    ..aOB(2, _omitFieldNames ? '' : 'active')
    ..a<$core.int>(3, _omitFieldNames ? '' : 'height', $pb.PbFieldType.OU3)
    ..a<$core.int>(4, _omitFieldNames ? '' : 'round', $pb.PbFieldType.O3)
    ..pc<VoteInfo>(5, _omitFieldNames ? '' : 'votes', $pb.PbFieldType.PM, subBuilder: VoteInfo.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  ConsensusInfo clone() => ConsensusInfo()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  ConsensusInfo copyWith(void Function(ConsensusInfo) updates) => super.copyWith((message) => updates(message as ConsensusInfo)) as ConsensusInfo;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ConsensusInfo create() => ConsensusInfo._();
  ConsensusInfo createEmptyInstance() => create();
  static $pb.PbList<ConsensusInfo> createRepeated() => $pb.PbList<ConsensusInfo>();
  @$core.pragma('dart2js:noInline')
  static ConsensusInfo getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ConsensusInfo>(create);
  static ConsensusInfo? _defaultInstance;

  /// The address of the consensus instance.
  @$pb.TagNumber(1)
  $core.String get address => $_getSZ(0);
  @$pb.TagNumber(1)
  set address($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasAddress() => $_has(0);
  @$pb.TagNumber(1)
  void clearAddress() => $_clearField(1);

  /// Indicates whether the consensus instance is active and part of the committee.
  @$pb.TagNumber(2)
  $core.bool get active => $_getBF(1);
  @$pb.TagNumber(2)
  set active($core.bool v) { $_setBool(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasActive() => $_has(1);
  @$pb.TagNumber(2)
  void clearActive() => $_clearField(2);

  /// The height of the consensus instance.
  @$pb.TagNumber(3)
  $core.int get height => $_getIZ(2);
  @$pb.TagNumber(3)
  set height($core.int v) { $_setUnsignedInt32(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasHeight() => $_has(2);
  @$pb.TagNumber(3)
  void clearHeight() => $_clearField(3);

  /// The round of the consensus instance.
  @$pb.TagNumber(4)
  $core.int get round => $_getIZ(3);
  @$pb.TagNumber(4)
  set round($core.int v) { $_setSignedInt32(3, v); }
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
    final $result = create();
    if (height != null) {
      $result.height = height;
    }
    if (round != null) {
      $result.round = round;
    }
    if (blockData != null) {
      $result.blockData = blockData;
    }
    if (signature != null) {
      $result.signature = signature;
    }
    return $result;
  }
  ProposalInfo._() : super();
  factory ProposalInfo.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory ProposalInfo.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'ProposalInfo', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..a<$core.int>(1, _omitFieldNames ? '' : 'height', $pb.PbFieldType.OU3)
    ..a<$core.int>(2, _omitFieldNames ? '' : 'round', $pb.PbFieldType.O3)
    ..aOS(3, _omitFieldNames ? '' : 'blockData')
    ..aOS(4, _omitFieldNames ? '' : 'signature')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  ProposalInfo clone() => ProposalInfo()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  ProposalInfo copyWith(void Function(ProposalInfo) updates) => super.copyWith((message) => updates(message as ProposalInfo)) as ProposalInfo;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ProposalInfo create() => ProposalInfo._();
  ProposalInfo createEmptyInstance() => create();
  static $pb.PbList<ProposalInfo> createRepeated() => $pb.PbList<ProposalInfo>();
  @$core.pragma('dart2js:noInline')
  static ProposalInfo getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ProposalInfo>(create);
  static ProposalInfo? _defaultInstance;

  /// The height of the proposal.
  @$pb.TagNumber(1)
  $core.int get height => $_getIZ(0);
  @$pb.TagNumber(1)
  set height($core.int v) { $_setUnsignedInt32(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasHeight() => $_has(0);
  @$pb.TagNumber(1)
  void clearHeight() => $_clearField(1);

  /// The round of the proposal.
  @$pb.TagNumber(2)
  $core.int get round => $_getIZ(1);
  @$pb.TagNumber(2)
  set round($core.int v) { $_setSignedInt32(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasRound() => $_has(1);
  @$pb.TagNumber(2)
  void clearRound() => $_clearField(2);

  /// The block data of the proposal.
  @$pb.TagNumber(3)
  $core.String get blockData => $_getSZ(2);
  @$pb.TagNumber(3)
  set blockData($core.String v) { $_setString(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasBlockData() => $_has(2);
  @$pb.TagNumber(3)
  void clearBlockData() => $_clearField(3);

  /// The signature of the proposal, signed by the proposer.
  @$pb.TagNumber(4)
  $core.String get signature => $_getSZ(3);
  @$pb.TagNumber(4)
  set signature($core.String v) { $_setString(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasSignature() => $_has(3);
  @$pb.TagNumber(4)
  void clearSignature() => $_clearField(4);
}

/// Blockchain service defines RPC methods for interacting with the blockchain.
class BlockchainApi {
  $pb.RpcClient _client;
  BlockchainApi(this._client);

  /// GetBlock retrieves information about a block based on the provided request parameters.
  $async.Future<GetBlockResponse> getBlock($pb.ClientContext? ctx, GetBlockRequest request) =>
    _client.invoke<GetBlockResponse>(ctx, 'Blockchain', 'GetBlock', request, GetBlockResponse())
  ;
  /// GetBlockHash retrieves the hash of a block at the specified height.
  $async.Future<GetBlockHashResponse> getBlockHash($pb.ClientContext? ctx, GetBlockHashRequest request) =>
    _client.invoke<GetBlockHashResponse>(ctx, 'Blockchain', 'GetBlockHash', request, GetBlockHashResponse())
  ;
  /// GetBlockHeight retrieves the height of a block with the specified hash.
  $async.Future<GetBlockHeightResponse> getBlockHeight($pb.ClientContext? ctx, GetBlockHeightRequest request) =>
    _client.invoke<GetBlockHeightResponse>(ctx, 'Blockchain', 'GetBlockHeight', request, GetBlockHeightResponse())
  ;
  /// GetBlockchainInfo retrieves general information about the blockchain.
  $async.Future<GetBlockchainInfoResponse> getBlockchainInfo($pb.ClientContext? ctx, GetBlockchainInfoRequest request) =>
    _client.invoke<GetBlockchainInfoResponse>(ctx, 'Blockchain', 'GetBlockchainInfo', request, GetBlockchainInfoResponse())
  ;
  /// GetConsensusInfo retrieves information about the consensus instances.
  $async.Future<GetConsensusInfoResponse> getConsensusInfo($pb.ClientContext? ctx, GetConsensusInfoRequest request) =>
    _client.invoke<GetConsensusInfoResponse>(ctx, 'Blockchain', 'GetConsensusInfo', request, GetConsensusInfoResponse())
  ;
  /// GetAccount retrieves information about an account based on the provided address.
  $async.Future<GetAccountResponse> getAccount($pb.ClientContext? ctx, GetAccountRequest request) =>
    _client.invoke<GetAccountResponse>(ctx, 'Blockchain', 'GetAccount', request, GetAccountResponse())
  ;
  /// GetValidator retrieves information about a validator based on the provided address.
  $async.Future<GetValidatorResponse> getValidator($pb.ClientContext? ctx, GetValidatorRequest request) =>
    _client.invoke<GetValidatorResponse>(ctx, 'Blockchain', 'GetValidator', request, GetValidatorResponse())
  ;
  /// GetValidatorByNumber retrieves information about a validator based on the provided number.
  $async.Future<GetValidatorResponse> getValidatorByNumber($pb.ClientContext? ctx, GetValidatorByNumberRequest request) =>
    _client.invoke<GetValidatorResponse>(ctx, 'Blockchain', 'GetValidatorByNumber', request, GetValidatorResponse())
  ;
  /// GetValidatorAddresses retrieves a list of all validator addresses.
  $async.Future<GetValidatorAddressesResponse> getValidatorAddresses($pb.ClientContext? ctx, GetValidatorAddressesRequest request) =>
    _client.invoke<GetValidatorAddressesResponse>(ctx, 'Blockchain', 'GetValidatorAddresses', request, GetValidatorAddressesResponse())
  ;
  /// GetPublicKey retrieves the public key of an account based on the provided address.
  $async.Future<GetPublicKeyResponse> getPublicKey($pb.ClientContext? ctx, GetPublicKeyRequest request) =>
    _client.invoke<GetPublicKeyResponse>(ctx, 'Blockchain', 'GetPublicKey', request, GetPublicKeyResponse())
  ;
  /// GetTxPoolContent retrieves current transactions in the transaction pool.
  $async.Future<GetTxPoolContentResponse> getTxPoolContent($pb.ClientContext? ctx, GetTxPoolContentRequest request) =>
    _client.invoke<GetTxPoolContentResponse>(ctx, 'Blockchain', 'GetTxPoolContent', request, GetTxPoolContentResponse())
  ;
}


const _omitFieldNames = $core.bool.fromEnvironment('protobuf.omit_field_names');
const _omitMessageNames = $core.bool.fromEnvironment('protobuf.omit_message_names');
