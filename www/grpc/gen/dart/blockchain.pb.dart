///
//  Generated code. Do not modify.
//  source: blockchain.proto
//
// @dart = 2.12
// ignore_for_file: annotate_overrides,camel_case_types,constant_identifier_names,directives_ordering,library_prefixes,non_constant_identifier_names,prefer_final_fields,return_of_invalid_type,unnecessary_const,unnecessary_import,unnecessary_this,unused_import,unused_shown_name

import 'dart:async' as $async;
import 'dart:core' as $core;

import 'package:fixnum/fixnum.dart' as $fixnum;
import 'package:protobuf/protobuf.dart' as $pb;

import 'transaction.pb.dart' as $0;

import 'blockchain.pbenum.dart';

export 'blockchain.pbenum.dart';

class GetAccountRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'GetAccountRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'address')
    ..hasRequiredFields = false
  ;

  GetAccountRequest._() : super();
  factory GetAccountRequest({
    $core.String? address,
  }) {
    final _result = create();
    if (address != null) {
      _result.address = address;
    }
    return _result;
  }
  factory GetAccountRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetAccountRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetAccountRequest clone() => GetAccountRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetAccountRequest copyWith(void Function(GetAccountRequest) updates) => super.copyWith((message) => updates(message as GetAccountRequest)) as GetAccountRequest; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static GetAccountRequest create() => GetAccountRequest._();
  GetAccountRequest createEmptyInstance() => create();
  static $pb.PbList<GetAccountRequest> createRepeated() => $pb.PbList<GetAccountRequest>();
  @$core.pragma('dart2js:noInline')
  static GetAccountRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetAccountRequest>(create);
  static GetAccountRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get address => $_getSZ(0);
  @$pb.TagNumber(1)
  set address($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasAddress() => $_has(0);
  @$pb.TagNumber(1)
  void clearAddress() => clearField(1);
}

class GetAccountByNumberRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'GetAccountByNumberRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..a<$core.int>(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'number', $pb.PbFieldType.O3)
    ..hasRequiredFields = false
  ;

  GetAccountByNumberRequest._() : super();
  factory GetAccountByNumberRequest({
    $core.int? number,
  }) {
    final _result = create();
    if (number != null) {
      _result.number = number;
    }
    return _result;
  }
  factory GetAccountByNumberRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetAccountByNumberRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetAccountByNumberRequest clone() => GetAccountByNumberRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetAccountByNumberRequest copyWith(void Function(GetAccountByNumberRequest) updates) => super.copyWith((message) => updates(message as GetAccountByNumberRequest)) as GetAccountByNumberRequest; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static GetAccountByNumberRequest create() => GetAccountByNumberRequest._();
  GetAccountByNumberRequest createEmptyInstance() => create();
  static $pb.PbList<GetAccountByNumberRequest> createRepeated() => $pb.PbList<GetAccountByNumberRequest>();
  @$core.pragma('dart2js:noInline')
  static GetAccountByNumberRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetAccountByNumberRequest>(create);
  static GetAccountByNumberRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.int get number => $_getIZ(0);
  @$pb.TagNumber(1)
  set number($core.int v) { $_setSignedInt32(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasNumber() => $_has(0);
  @$pb.TagNumber(1)
  void clearNumber() => clearField(1);
}

class GetAccountResponse extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'GetAccountResponse', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOM<AccountInfo>(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'account', subBuilder: AccountInfo.create)
    ..hasRequiredFields = false
  ;

  GetAccountResponse._() : super();
  factory GetAccountResponse({
    AccountInfo? account,
  }) {
    final _result = create();
    if (account != null) {
      _result.account = account;
    }
    return _result;
  }
  factory GetAccountResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetAccountResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetAccountResponse clone() => GetAccountResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetAccountResponse copyWith(void Function(GetAccountResponse) updates) => super.copyWith((message) => updates(message as GetAccountResponse)) as GetAccountResponse; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static GetAccountResponse create() => GetAccountResponse._();
  GetAccountResponse createEmptyInstance() => create();
  static $pb.PbList<GetAccountResponse> createRepeated() => $pb.PbList<GetAccountResponse>();
  @$core.pragma('dart2js:noInline')
  static GetAccountResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetAccountResponse>(create);
  static GetAccountResponse? _defaultInstance;

  @$pb.TagNumber(1)
  AccountInfo get account => $_getN(0);
  @$pb.TagNumber(1)
  set account(AccountInfo v) { setField(1, v); }
  @$pb.TagNumber(1)
  $core.bool hasAccount() => $_has(0);
  @$pb.TagNumber(1)
  void clearAccount() => clearField(1);
  @$pb.TagNumber(1)
  AccountInfo ensureAccount() => $_ensure(0);
}

class GetValidatorAddressesRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'GetValidatorAddressesRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..hasRequiredFields = false
  ;

  GetValidatorAddressesRequest._() : super();
  factory GetValidatorAddressesRequest() => create();
  factory GetValidatorAddressesRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetValidatorAddressesRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetValidatorAddressesRequest clone() => GetValidatorAddressesRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetValidatorAddressesRequest copyWith(void Function(GetValidatorAddressesRequest) updates) => super.copyWith((message) => updates(message as GetValidatorAddressesRequest)) as GetValidatorAddressesRequest; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static GetValidatorAddressesRequest create() => GetValidatorAddressesRequest._();
  GetValidatorAddressesRequest createEmptyInstance() => create();
  static $pb.PbList<GetValidatorAddressesRequest> createRepeated() => $pb.PbList<GetValidatorAddressesRequest>();
  @$core.pragma('dart2js:noInline')
  static GetValidatorAddressesRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetValidatorAddressesRequest>(create);
  static GetValidatorAddressesRequest? _defaultInstance;
}

class GetValidatorAddressesResponse extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'GetValidatorAddressesResponse', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..pPS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'addresses')
    ..hasRequiredFields = false
  ;

  GetValidatorAddressesResponse._() : super();
  factory GetValidatorAddressesResponse({
    $core.Iterable<$core.String>? addresses,
  }) {
    final _result = create();
    if (addresses != null) {
      _result.addresses.addAll(addresses);
    }
    return _result;
  }
  factory GetValidatorAddressesResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetValidatorAddressesResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetValidatorAddressesResponse clone() => GetValidatorAddressesResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetValidatorAddressesResponse copyWith(void Function(GetValidatorAddressesResponse) updates) => super.copyWith((message) => updates(message as GetValidatorAddressesResponse)) as GetValidatorAddressesResponse; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static GetValidatorAddressesResponse create() => GetValidatorAddressesResponse._();
  GetValidatorAddressesResponse createEmptyInstance() => create();
  static $pb.PbList<GetValidatorAddressesResponse> createRepeated() => $pb.PbList<GetValidatorAddressesResponse>();
  @$core.pragma('dart2js:noInline')
  static GetValidatorAddressesResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetValidatorAddressesResponse>(create);
  static GetValidatorAddressesResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.List<$core.String> get addresses => $_getList(0);
}

class GetValidatorRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'GetValidatorRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'address')
    ..hasRequiredFields = false
  ;

  GetValidatorRequest._() : super();
  factory GetValidatorRequest({
    $core.String? address,
  }) {
    final _result = create();
    if (address != null) {
      _result.address = address;
    }
    return _result;
  }
  factory GetValidatorRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetValidatorRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetValidatorRequest clone() => GetValidatorRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetValidatorRequest copyWith(void Function(GetValidatorRequest) updates) => super.copyWith((message) => updates(message as GetValidatorRequest)) as GetValidatorRequest; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static GetValidatorRequest create() => GetValidatorRequest._();
  GetValidatorRequest createEmptyInstance() => create();
  static $pb.PbList<GetValidatorRequest> createRepeated() => $pb.PbList<GetValidatorRequest>();
  @$core.pragma('dart2js:noInline')
  static GetValidatorRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetValidatorRequest>(create);
  static GetValidatorRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get address => $_getSZ(0);
  @$pb.TagNumber(1)
  set address($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasAddress() => $_has(0);
  @$pb.TagNumber(1)
  void clearAddress() => clearField(1);
}

class GetValidatorByNumberRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'GetValidatorByNumberRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..a<$core.int>(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'number', $pb.PbFieldType.O3)
    ..hasRequiredFields = false
  ;

  GetValidatorByNumberRequest._() : super();
  factory GetValidatorByNumberRequest({
    $core.int? number,
  }) {
    final _result = create();
    if (number != null) {
      _result.number = number;
    }
    return _result;
  }
  factory GetValidatorByNumberRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetValidatorByNumberRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetValidatorByNumberRequest clone() => GetValidatorByNumberRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetValidatorByNumberRequest copyWith(void Function(GetValidatorByNumberRequest) updates) => super.copyWith((message) => updates(message as GetValidatorByNumberRequest)) as GetValidatorByNumberRequest; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static GetValidatorByNumberRequest create() => GetValidatorByNumberRequest._();
  GetValidatorByNumberRequest createEmptyInstance() => create();
  static $pb.PbList<GetValidatorByNumberRequest> createRepeated() => $pb.PbList<GetValidatorByNumberRequest>();
  @$core.pragma('dart2js:noInline')
  static GetValidatorByNumberRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetValidatorByNumberRequest>(create);
  static GetValidatorByNumberRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.int get number => $_getIZ(0);
  @$pb.TagNumber(1)
  set number($core.int v) { $_setSignedInt32(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasNumber() => $_has(0);
  @$pb.TagNumber(1)
  void clearNumber() => clearField(1);
}

class GetValidatorResponse extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'GetValidatorResponse', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOM<ValidatorInfo>(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'validator', subBuilder: ValidatorInfo.create)
    ..hasRequiredFields = false
  ;

  GetValidatorResponse._() : super();
  factory GetValidatorResponse({
    ValidatorInfo? validator,
  }) {
    final _result = create();
    if (validator != null) {
      _result.validator = validator;
    }
    return _result;
  }
  factory GetValidatorResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetValidatorResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetValidatorResponse clone() => GetValidatorResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetValidatorResponse copyWith(void Function(GetValidatorResponse) updates) => super.copyWith((message) => updates(message as GetValidatorResponse)) as GetValidatorResponse; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static GetValidatorResponse create() => GetValidatorResponse._();
  GetValidatorResponse createEmptyInstance() => create();
  static $pb.PbList<GetValidatorResponse> createRepeated() => $pb.PbList<GetValidatorResponse>();
  @$core.pragma('dart2js:noInline')
  static GetValidatorResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetValidatorResponse>(create);
  static GetValidatorResponse? _defaultInstance;

  @$pb.TagNumber(1)
  ValidatorInfo get validator => $_getN(0);
  @$pb.TagNumber(1)
  set validator(ValidatorInfo v) { setField(1, v); }
  @$pb.TagNumber(1)
  $core.bool hasValidator() => $_has(0);
  @$pb.TagNumber(1)
  void clearValidator() => clearField(1);
  @$pb.TagNumber(1)
  ValidatorInfo ensureValidator() => $_ensure(0);
}

class GetBlockRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'GetBlockRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..a<$core.int>(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'height', $pb.PbFieldType.OU3)
    ..e<BlockVerbosity>(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'verbosity', $pb.PbFieldType.OE, defaultOrMaker: BlockVerbosity.BLOCK_DATA, valueOf: BlockVerbosity.valueOf, enumValues: BlockVerbosity.values)
    ..hasRequiredFields = false
  ;

  GetBlockRequest._() : super();
  factory GetBlockRequest({
    $core.int? height,
    BlockVerbosity? verbosity,
  }) {
    final _result = create();
    if (height != null) {
      _result.height = height;
    }
    if (verbosity != null) {
      _result.verbosity = verbosity;
    }
    return _result;
  }
  factory GetBlockRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetBlockRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetBlockRequest clone() => GetBlockRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetBlockRequest copyWith(void Function(GetBlockRequest) updates) => super.copyWith((message) => updates(message as GetBlockRequest)) as GetBlockRequest; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static GetBlockRequest create() => GetBlockRequest._();
  GetBlockRequest createEmptyInstance() => create();
  static $pb.PbList<GetBlockRequest> createRepeated() => $pb.PbList<GetBlockRequest>();
  @$core.pragma('dart2js:noInline')
  static GetBlockRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetBlockRequest>(create);
  static GetBlockRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.int get height => $_getIZ(0);
  @$pb.TagNumber(1)
  set height($core.int v) { $_setUnsignedInt32(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasHeight() => $_has(0);
  @$pb.TagNumber(1)
  void clearHeight() => clearField(1);

  @$pb.TagNumber(2)
  BlockVerbosity get verbosity => $_getN(1);
  @$pb.TagNumber(2)
  set verbosity(BlockVerbosity v) { setField(2, v); }
  @$pb.TagNumber(2)
  $core.bool hasVerbosity() => $_has(1);
  @$pb.TagNumber(2)
  void clearVerbosity() => clearField(2);
}

class GetBlockResponse extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'GetBlockResponse', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..a<$core.int>(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'height', $pb.PbFieldType.OU3)
    ..a<$core.List<$core.int>>(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'hash', $pb.PbFieldType.OY)
    ..a<$core.List<$core.int>>(3, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'data', $pb.PbFieldType.OY)
    ..a<$core.int>(4, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'blockTime', $pb.PbFieldType.OU3)
    ..aOM<BlockHeaderInfo>(5, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'header', subBuilder: BlockHeaderInfo.create)
    ..aOM<CertificateInfo>(6, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'prevCert', subBuilder: CertificateInfo.create)
    ..pc<$0.TransactionInfo>(7, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'txs', $pb.PbFieldType.PM, subBuilder: $0.TransactionInfo.create)
    ..hasRequiredFields = false
  ;

  GetBlockResponse._() : super();
  factory GetBlockResponse({
    $core.int? height,
    $core.List<$core.int>? hash,
    $core.List<$core.int>? data,
    $core.int? blockTime,
    BlockHeaderInfo? header,
    CertificateInfo? prevCert,
    $core.Iterable<$0.TransactionInfo>? txs,
  }) {
    final _result = create();
    if (height != null) {
      _result.height = height;
    }
    if (hash != null) {
      _result.hash = hash;
    }
    if (data != null) {
      _result.data = data;
    }
    if (blockTime != null) {
      _result.blockTime = blockTime;
    }
    if (header != null) {
      _result.header = header;
    }
    if (prevCert != null) {
      _result.prevCert = prevCert;
    }
    if (txs != null) {
      _result.txs.addAll(txs);
    }
    return _result;
  }
  factory GetBlockResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetBlockResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetBlockResponse clone() => GetBlockResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetBlockResponse copyWith(void Function(GetBlockResponse) updates) => super.copyWith((message) => updates(message as GetBlockResponse)) as GetBlockResponse; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static GetBlockResponse create() => GetBlockResponse._();
  GetBlockResponse createEmptyInstance() => create();
  static $pb.PbList<GetBlockResponse> createRepeated() => $pb.PbList<GetBlockResponse>();
  @$core.pragma('dart2js:noInline')
  static GetBlockResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetBlockResponse>(create);
  static GetBlockResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.int get height => $_getIZ(0);
  @$pb.TagNumber(1)
  set height($core.int v) { $_setUnsignedInt32(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasHeight() => $_has(0);
  @$pb.TagNumber(1)
  void clearHeight() => clearField(1);

  @$pb.TagNumber(2)
  $core.List<$core.int> get hash => $_getN(1);
  @$pb.TagNumber(2)
  set hash($core.List<$core.int> v) { $_setBytes(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasHash() => $_has(1);
  @$pb.TagNumber(2)
  void clearHash() => clearField(2);

  @$pb.TagNumber(3)
  $core.List<$core.int> get data => $_getN(2);
  @$pb.TagNumber(3)
  set data($core.List<$core.int> v) { $_setBytes(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasData() => $_has(2);
  @$pb.TagNumber(3)
  void clearData() => clearField(3);

  @$pb.TagNumber(4)
  $core.int get blockTime => $_getIZ(3);
  @$pb.TagNumber(4)
  set blockTime($core.int v) { $_setUnsignedInt32(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasBlockTime() => $_has(3);
  @$pb.TagNumber(4)
  void clearBlockTime() => clearField(4);

  @$pb.TagNumber(5)
  BlockHeaderInfo get header => $_getN(4);
  @$pb.TagNumber(5)
  set header(BlockHeaderInfo v) { setField(5, v); }
  @$pb.TagNumber(5)
  $core.bool hasHeader() => $_has(4);
  @$pb.TagNumber(5)
  void clearHeader() => clearField(5);
  @$pb.TagNumber(5)
  BlockHeaderInfo ensureHeader() => $_ensure(4);

  @$pb.TagNumber(6)
  CertificateInfo get prevCert => $_getN(5);
  @$pb.TagNumber(6)
  set prevCert(CertificateInfo v) { setField(6, v); }
  @$pb.TagNumber(6)
  $core.bool hasPrevCert() => $_has(5);
  @$pb.TagNumber(6)
  void clearPrevCert() => clearField(6);
  @$pb.TagNumber(6)
  CertificateInfo ensurePrevCert() => $_ensure(5);

  @$pb.TagNumber(7)
  $core.List<$0.TransactionInfo> get txs => $_getList(6);
}

class GetBlockHashRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'GetBlockHashRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..a<$core.int>(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'height', $pb.PbFieldType.OU3)
    ..hasRequiredFields = false
  ;

  GetBlockHashRequest._() : super();
  factory GetBlockHashRequest({
    $core.int? height,
  }) {
    final _result = create();
    if (height != null) {
      _result.height = height;
    }
    return _result;
  }
  factory GetBlockHashRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetBlockHashRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetBlockHashRequest clone() => GetBlockHashRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetBlockHashRequest copyWith(void Function(GetBlockHashRequest) updates) => super.copyWith((message) => updates(message as GetBlockHashRequest)) as GetBlockHashRequest; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static GetBlockHashRequest create() => GetBlockHashRequest._();
  GetBlockHashRequest createEmptyInstance() => create();
  static $pb.PbList<GetBlockHashRequest> createRepeated() => $pb.PbList<GetBlockHashRequest>();
  @$core.pragma('dart2js:noInline')
  static GetBlockHashRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetBlockHashRequest>(create);
  static GetBlockHashRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.int get height => $_getIZ(0);
  @$pb.TagNumber(1)
  set height($core.int v) { $_setUnsignedInt32(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasHeight() => $_has(0);
  @$pb.TagNumber(1)
  void clearHeight() => clearField(1);
}

class GetBlockHashResponse extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'GetBlockHashResponse', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..a<$core.List<$core.int>>(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'hash', $pb.PbFieldType.OY)
    ..hasRequiredFields = false
  ;

  GetBlockHashResponse._() : super();
  factory GetBlockHashResponse({
    $core.List<$core.int>? hash,
  }) {
    final _result = create();
    if (hash != null) {
      _result.hash = hash;
    }
    return _result;
  }
  factory GetBlockHashResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetBlockHashResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetBlockHashResponse clone() => GetBlockHashResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetBlockHashResponse copyWith(void Function(GetBlockHashResponse) updates) => super.copyWith((message) => updates(message as GetBlockHashResponse)) as GetBlockHashResponse; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static GetBlockHashResponse create() => GetBlockHashResponse._();
  GetBlockHashResponse createEmptyInstance() => create();
  static $pb.PbList<GetBlockHashResponse> createRepeated() => $pb.PbList<GetBlockHashResponse>();
  @$core.pragma('dart2js:noInline')
  static GetBlockHashResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetBlockHashResponse>(create);
  static GetBlockHashResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.List<$core.int> get hash => $_getN(0);
  @$pb.TagNumber(1)
  set hash($core.List<$core.int> v) { $_setBytes(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasHash() => $_has(0);
  @$pb.TagNumber(1)
  void clearHash() => clearField(1);
}

class GetBlockHeightRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'GetBlockHeightRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..a<$core.List<$core.int>>(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'hash', $pb.PbFieldType.OY)
    ..hasRequiredFields = false
  ;

  GetBlockHeightRequest._() : super();
  factory GetBlockHeightRequest({
    $core.List<$core.int>? hash,
  }) {
    final _result = create();
    if (hash != null) {
      _result.hash = hash;
    }
    return _result;
  }
  factory GetBlockHeightRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetBlockHeightRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetBlockHeightRequest clone() => GetBlockHeightRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetBlockHeightRequest copyWith(void Function(GetBlockHeightRequest) updates) => super.copyWith((message) => updates(message as GetBlockHeightRequest)) as GetBlockHeightRequest; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static GetBlockHeightRequest create() => GetBlockHeightRequest._();
  GetBlockHeightRequest createEmptyInstance() => create();
  static $pb.PbList<GetBlockHeightRequest> createRepeated() => $pb.PbList<GetBlockHeightRequest>();
  @$core.pragma('dart2js:noInline')
  static GetBlockHeightRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetBlockHeightRequest>(create);
  static GetBlockHeightRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.List<$core.int> get hash => $_getN(0);
  @$pb.TagNumber(1)
  set hash($core.List<$core.int> v) { $_setBytes(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasHash() => $_has(0);
  @$pb.TagNumber(1)
  void clearHash() => clearField(1);
}

class GetBlockHeightResponse extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'GetBlockHeightResponse', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..a<$core.int>(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'height', $pb.PbFieldType.OU3)
    ..hasRequiredFields = false
  ;

  GetBlockHeightResponse._() : super();
  factory GetBlockHeightResponse({
    $core.int? height,
  }) {
    final _result = create();
    if (height != null) {
      _result.height = height;
    }
    return _result;
  }
  factory GetBlockHeightResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetBlockHeightResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetBlockHeightResponse clone() => GetBlockHeightResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetBlockHeightResponse copyWith(void Function(GetBlockHeightResponse) updates) => super.copyWith((message) => updates(message as GetBlockHeightResponse)) as GetBlockHeightResponse; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static GetBlockHeightResponse create() => GetBlockHeightResponse._();
  GetBlockHeightResponse createEmptyInstance() => create();
  static $pb.PbList<GetBlockHeightResponse> createRepeated() => $pb.PbList<GetBlockHeightResponse>();
  @$core.pragma('dart2js:noInline')
  static GetBlockHeightResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetBlockHeightResponse>(create);
  static GetBlockHeightResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.int get height => $_getIZ(0);
  @$pb.TagNumber(1)
  set height($core.int v) { $_setUnsignedInt32(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasHeight() => $_has(0);
  @$pb.TagNumber(1)
  void clearHeight() => clearField(1);
}

class GetBlockchainInfoRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'GetBlockchainInfoRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..hasRequiredFields = false
  ;

  GetBlockchainInfoRequest._() : super();
  factory GetBlockchainInfoRequest() => create();
  factory GetBlockchainInfoRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetBlockchainInfoRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetBlockchainInfoRequest clone() => GetBlockchainInfoRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetBlockchainInfoRequest copyWith(void Function(GetBlockchainInfoRequest) updates) => super.copyWith((message) => updates(message as GetBlockchainInfoRequest)) as GetBlockchainInfoRequest; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static GetBlockchainInfoRequest create() => GetBlockchainInfoRequest._();
  GetBlockchainInfoRequest createEmptyInstance() => create();
  static $pb.PbList<GetBlockchainInfoRequest> createRepeated() => $pb.PbList<GetBlockchainInfoRequest>();
  @$core.pragma('dart2js:noInline')
  static GetBlockchainInfoRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetBlockchainInfoRequest>(create);
  static GetBlockchainInfoRequest? _defaultInstance;
}

class GetBlockchainInfoResponse extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'GetBlockchainInfoResponse', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..a<$core.int>(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'lastBlockHeight', $pb.PbFieldType.OU3)
    ..a<$core.List<$core.int>>(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'lastBlockHash', $pb.PbFieldType.OY)
    ..a<$core.int>(3, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'totalAccounts', $pb.PbFieldType.O3)
    ..a<$core.int>(4, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'totalValidators', $pb.PbFieldType.O3)
    ..aInt64(5, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'totalPower')
    ..aInt64(6, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'committeePower')
    ..pc<ValidatorInfo>(7, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'committeeValidators', $pb.PbFieldType.PM, subBuilder: ValidatorInfo.create)
    ..hasRequiredFields = false
  ;

  GetBlockchainInfoResponse._() : super();
  factory GetBlockchainInfoResponse({
    $core.int? lastBlockHeight,
    $core.List<$core.int>? lastBlockHash,
    $core.int? totalAccounts,
    $core.int? totalValidators,
    $fixnum.Int64? totalPower,
    $fixnum.Int64? committeePower,
    $core.Iterable<ValidatorInfo>? committeeValidators,
  }) {
    final _result = create();
    if (lastBlockHeight != null) {
      _result.lastBlockHeight = lastBlockHeight;
    }
    if (lastBlockHash != null) {
      _result.lastBlockHash = lastBlockHash;
    }
    if (totalAccounts != null) {
      _result.totalAccounts = totalAccounts;
    }
    if (totalValidators != null) {
      _result.totalValidators = totalValidators;
    }
    if (totalPower != null) {
      _result.totalPower = totalPower;
    }
    if (committeePower != null) {
      _result.committeePower = committeePower;
    }
    if (committeeValidators != null) {
      _result.committeeValidators.addAll(committeeValidators);
    }
    return _result;
  }
  factory GetBlockchainInfoResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetBlockchainInfoResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetBlockchainInfoResponse clone() => GetBlockchainInfoResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetBlockchainInfoResponse copyWith(void Function(GetBlockchainInfoResponse) updates) => super.copyWith((message) => updates(message as GetBlockchainInfoResponse)) as GetBlockchainInfoResponse; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static GetBlockchainInfoResponse create() => GetBlockchainInfoResponse._();
  GetBlockchainInfoResponse createEmptyInstance() => create();
  static $pb.PbList<GetBlockchainInfoResponse> createRepeated() => $pb.PbList<GetBlockchainInfoResponse>();
  @$core.pragma('dart2js:noInline')
  static GetBlockchainInfoResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetBlockchainInfoResponse>(create);
  static GetBlockchainInfoResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.int get lastBlockHeight => $_getIZ(0);
  @$pb.TagNumber(1)
  set lastBlockHeight($core.int v) { $_setUnsignedInt32(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasLastBlockHeight() => $_has(0);
  @$pb.TagNumber(1)
  void clearLastBlockHeight() => clearField(1);

  @$pb.TagNumber(2)
  $core.List<$core.int> get lastBlockHash => $_getN(1);
  @$pb.TagNumber(2)
  set lastBlockHash($core.List<$core.int> v) { $_setBytes(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasLastBlockHash() => $_has(1);
  @$pb.TagNumber(2)
  void clearLastBlockHash() => clearField(2);

  @$pb.TagNumber(3)
  $core.int get totalAccounts => $_getIZ(2);
  @$pb.TagNumber(3)
  set totalAccounts($core.int v) { $_setSignedInt32(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasTotalAccounts() => $_has(2);
  @$pb.TagNumber(3)
  void clearTotalAccounts() => clearField(3);

  @$pb.TagNumber(4)
  $core.int get totalValidators => $_getIZ(3);
  @$pb.TagNumber(4)
  set totalValidators($core.int v) { $_setSignedInt32(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasTotalValidators() => $_has(3);
  @$pb.TagNumber(4)
  void clearTotalValidators() => clearField(4);

  @$pb.TagNumber(5)
  $fixnum.Int64 get totalPower => $_getI64(4);
  @$pb.TagNumber(5)
  set totalPower($fixnum.Int64 v) { $_setInt64(4, v); }
  @$pb.TagNumber(5)
  $core.bool hasTotalPower() => $_has(4);
  @$pb.TagNumber(5)
  void clearTotalPower() => clearField(5);

  @$pb.TagNumber(6)
  $fixnum.Int64 get committeePower => $_getI64(5);
  @$pb.TagNumber(6)
  set committeePower($fixnum.Int64 v) { $_setInt64(5, v); }
  @$pb.TagNumber(6)
  $core.bool hasCommitteePower() => $_has(5);
  @$pb.TagNumber(6)
  void clearCommitteePower() => clearField(6);

  @$pb.TagNumber(7)
  $core.List<ValidatorInfo> get committeeValidators => $_getList(6);
}

class GetConsensusInfoRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'GetConsensusInfoRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..hasRequiredFields = false
  ;

  GetConsensusInfoRequest._() : super();
  factory GetConsensusInfoRequest() => create();
  factory GetConsensusInfoRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetConsensusInfoRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetConsensusInfoRequest clone() => GetConsensusInfoRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetConsensusInfoRequest copyWith(void Function(GetConsensusInfoRequest) updates) => super.copyWith((message) => updates(message as GetConsensusInfoRequest)) as GetConsensusInfoRequest; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static GetConsensusInfoRequest create() => GetConsensusInfoRequest._();
  GetConsensusInfoRequest createEmptyInstance() => create();
  static $pb.PbList<GetConsensusInfoRequest> createRepeated() => $pb.PbList<GetConsensusInfoRequest>();
  @$core.pragma('dart2js:noInline')
  static GetConsensusInfoRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetConsensusInfoRequest>(create);
  static GetConsensusInfoRequest? _defaultInstance;
}

class GetConsensusInfoResponse extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'GetConsensusInfoResponse', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..pc<ConsensusInfo>(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'instances', $pb.PbFieldType.PM, subBuilder: ConsensusInfo.create)
    ..hasRequiredFields = false
  ;

  GetConsensusInfoResponse._() : super();
  factory GetConsensusInfoResponse({
    $core.Iterable<ConsensusInfo>? instances,
  }) {
    final _result = create();
    if (instances != null) {
      _result.instances.addAll(instances);
    }
    return _result;
  }
  factory GetConsensusInfoResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetConsensusInfoResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetConsensusInfoResponse clone() => GetConsensusInfoResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetConsensusInfoResponse copyWith(void Function(GetConsensusInfoResponse) updates) => super.copyWith((message) => updates(message as GetConsensusInfoResponse)) as GetConsensusInfoResponse; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static GetConsensusInfoResponse create() => GetConsensusInfoResponse._();
  GetConsensusInfoResponse createEmptyInstance() => create();
  static $pb.PbList<GetConsensusInfoResponse> createRepeated() => $pb.PbList<GetConsensusInfoResponse>();
  @$core.pragma('dart2js:noInline')
  static GetConsensusInfoResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetConsensusInfoResponse>(create);
  static GetConsensusInfoResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.List<ConsensusInfo> get instances => $_getList(0);
}

class ValidatorInfo extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'ValidatorInfo', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..a<$core.List<$core.int>>(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'hash', $pb.PbFieldType.OY)
    ..a<$core.List<$core.int>>(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'data', $pb.PbFieldType.OY)
    ..aOS(3, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'publicKey')
    ..a<$core.int>(4, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'number', $pb.PbFieldType.O3)
    ..a<$core.int>(5, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'sequence', $pb.PbFieldType.O3)
    ..aInt64(6, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'stake')
    ..a<$core.int>(7, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'lastBondingHeight', $pb.PbFieldType.OU3)
    ..a<$core.int>(8, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'lastJoinedHeight', $pb.PbFieldType.OU3)
    ..a<$core.int>(9, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'unbondingHeight', $pb.PbFieldType.OU3)
    ..aOS(10, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'address')
    ..hasRequiredFields = false
  ;

  ValidatorInfo._() : super();
  factory ValidatorInfo({
    $core.List<$core.int>? hash,
    $core.List<$core.int>? data,
    $core.String? publicKey,
    $core.int? number,
    $core.int? sequence,
    $fixnum.Int64? stake,
    $core.int? lastBondingHeight,
    $core.int? lastJoinedHeight,
    $core.int? unbondingHeight,
    $core.String? address,
  }) {
    final _result = create();
    if (hash != null) {
      _result.hash = hash;
    }
    if (data != null) {
      _result.data = data;
    }
    if (publicKey != null) {
      _result.publicKey = publicKey;
    }
    if (number != null) {
      _result.number = number;
    }
    if (sequence != null) {
      _result.sequence = sequence;
    }
    if (stake != null) {
      _result.stake = stake;
    }
    if (lastBondingHeight != null) {
      _result.lastBondingHeight = lastBondingHeight;
    }
    if (lastJoinedHeight != null) {
      _result.lastJoinedHeight = lastJoinedHeight;
    }
    if (unbondingHeight != null) {
      _result.unbondingHeight = unbondingHeight;
    }
    if (address != null) {
      _result.address = address;
    }
    return _result;
  }
  factory ValidatorInfo.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory ValidatorInfo.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  ValidatorInfo clone() => ValidatorInfo()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  ValidatorInfo copyWith(void Function(ValidatorInfo) updates) => super.copyWith((message) => updates(message as ValidatorInfo)) as ValidatorInfo; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static ValidatorInfo create() => ValidatorInfo._();
  ValidatorInfo createEmptyInstance() => create();
  static $pb.PbList<ValidatorInfo> createRepeated() => $pb.PbList<ValidatorInfo>();
  @$core.pragma('dart2js:noInline')
  static ValidatorInfo getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ValidatorInfo>(create);
  static ValidatorInfo? _defaultInstance;

  @$pb.TagNumber(1)
  $core.List<$core.int> get hash => $_getN(0);
  @$pb.TagNumber(1)
  set hash($core.List<$core.int> v) { $_setBytes(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasHash() => $_has(0);
  @$pb.TagNumber(1)
  void clearHash() => clearField(1);

  @$pb.TagNumber(2)
  $core.List<$core.int> get data => $_getN(1);
  @$pb.TagNumber(2)
  set data($core.List<$core.int> v) { $_setBytes(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasData() => $_has(1);
  @$pb.TagNumber(2)
  void clearData() => clearField(2);

  @$pb.TagNumber(3)
  $core.String get publicKey => $_getSZ(2);
  @$pb.TagNumber(3)
  set publicKey($core.String v) { $_setString(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasPublicKey() => $_has(2);
  @$pb.TagNumber(3)
  void clearPublicKey() => clearField(3);

  @$pb.TagNumber(4)
  $core.int get number => $_getIZ(3);
  @$pb.TagNumber(4)
  set number($core.int v) { $_setSignedInt32(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasNumber() => $_has(3);
  @$pb.TagNumber(4)
  void clearNumber() => clearField(4);

  @$pb.TagNumber(5)
  $core.int get sequence => $_getIZ(4);
  @$pb.TagNumber(5)
  set sequence($core.int v) { $_setSignedInt32(4, v); }
  @$pb.TagNumber(5)
  $core.bool hasSequence() => $_has(4);
  @$pb.TagNumber(5)
  void clearSequence() => clearField(5);

  @$pb.TagNumber(6)
  $fixnum.Int64 get stake => $_getI64(5);
  @$pb.TagNumber(6)
  set stake($fixnum.Int64 v) { $_setInt64(5, v); }
  @$pb.TagNumber(6)
  $core.bool hasStake() => $_has(5);
  @$pb.TagNumber(6)
  void clearStake() => clearField(6);

  @$pb.TagNumber(7)
  $core.int get lastBondingHeight => $_getIZ(6);
  @$pb.TagNumber(7)
  set lastBondingHeight($core.int v) { $_setUnsignedInt32(6, v); }
  @$pb.TagNumber(7)
  $core.bool hasLastBondingHeight() => $_has(6);
  @$pb.TagNumber(7)
  void clearLastBondingHeight() => clearField(7);

  @$pb.TagNumber(8)
  $core.int get lastJoinedHeight => $_getIZ(7);
  @$pb.TagNumber(8)
  set lastJoinedHeight($core.int v) { $_setUnsignedInt32(7, v); }
  @$pb.TagNumber(8)
  $core.bool hasLastJoinedHeight() => $_has(7);
  @$pb.TagNumber(8)
  void clearLastJoinedHeight() => clearField(8);

  @$pb.TagNumber(9)
  $core.int get unbondingHeight => $_getIZ(8);
  @$pb.TagNumber(9)
  set unbondingHeight($core.int v) { $_setUnsignedInt32(8, v); }
  @$pb.TagNumber(9)
  $core.bool hasUnbondingHeight() => $_has(8);
  @$pb.TagNumber(9)
  void clearUnbondingHeight() => clearField(9);

  @$pb.TagNumber(10)
  $core.String get address => $_getSZ(9);
  @$pb.TagNumber(10)
  set address($core.String v) { $_setString(9, v); }
  @$pb.TagNumber(10)
  $core.bool hasAddress() => $_has(9);
  @$pb.TagNumber(10)
  void clearAddress() => clearField(10);
}

class AccountInfo extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'AccountInfo', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..a<$core.List<$core.int>>(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'hash', $pb.PbFieldType.OY)
    ..a<$core.List<$core.int>>(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'data', $pb.PbFieldType.OY)
    ..a<$core.int>(3, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'number', $pb.PbFieldType.O3)
    ..a<$core.int>(4, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'sequence', $pb.PbFieldType.O3)
    ..aInt64(5, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'balance')
    ..hasRequiredFields = false
  ;

  AccountInfo._() : super();
  factory AccountInfo({
    $core.List<$core.int>? hash,
    $core.List<$core.int>? data,
    $core.int? number,
    $core.int? sequence,
    $fixnum.Int64? balance,
  }) {
    final _result = create();
    if (hash != null) {
      _result.hash = hash;
    }
    if (data != null) {
      _result.data = data;
    }
    if (number != null) {
      _result.number = number;
    }
    if (sequence != null) {
      _result.sequence = sequence;
    }
    if (balance != null) {
      _result.balance = balance;
    }
    return _result;
  }
  factory AccountInfo.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory AccountInfo.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  AccountInfo clone() => AccountInfo()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  AccountInfo copyWith(void Function(AccountInfo) updates) => super.copyWith((message) => updates(message as AccountInfo)) as AccountInfo; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static AccountInfo create() => AccountInfo._();
  AccountInfo createEmptyInstance() => create();
  static $pb.PbList<AccountInfo> createRepeated() => $pb.PbList<AccountInfo>();
  @$core.pragma('dart2js:noInline')
  static AccountInfo getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<AccountInfo>(create);
  static AccountInfo? _defaultInstance;

  @$pb.TagNumber(1)
  $core.List<$core.int> get hash => $_getN(0);
  @$pb.TagNumber(1)
  set hash($core.List<$core.int> v) { $_setBytes(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasHash() => $_has(0);
  @$pb.TagNumber(1)
  void clearHash() => clearField(1);

  @$pb.TagNumber(2)
  $core.List<$core.int> get data => $_getN(1);
  @$pb.TagNumber(2)
  set data($core.List<$core.int> v) { $_setBytes(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasData() => $_has(1);
  @$pb.TagNumber(2)
  void clearData() => clearField(2);

  @$pb.TagNumber(3)
  $core.int get number => $_getIZ(2);
  @$pb.TagNumber(3)
  set number($core.int v) { $_setSignedInt32(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasNumber() => $_has(2);
  @$pb.TagNumber(3)
  void clearNumber() => clearField(3);

  @$pb.TagNumber(4)
  $core.int get sequence => $_getIZ(3);
  @$pb.TagNumber(4)
  set sequence($core.int v) { $_setSignedInt32(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasSequence() => $_has(3);
  @$pb.TagNumber(4)
  void clearSequence() => clearField(4);

  @$pb.TagNumber(5)
  $fixnum.Int64 get balance => $_getI64(4);
  @$pb.TagNumber(5)
  set balance($fixnum.Int64 v) { $_setInt64(4, v); }
  @$pb.TagNumber(5)
  $core.bool hasBalance() => $_has(4);
  @$pb.TagNumber(5)
  void clearBalance() => clearField(5);
}

class BlockHeaderInfo extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'BlockHeaderInfo', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..a<$core.int>(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'version', $pb.PbFieldType.O3)
    ..a<$core.List<$core.int>>(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'prevBlockHash', $pb.PbFieldType.OY)
    ..a<$core.List<$core.int>>(3, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'stateRoot', $pb.PbFieldType.OY)
    ..a<$core.List<$core.int>>(4, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'sortitionSeed', $pb.PbFieldType.OY)
    ..aOS(5, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'proposerAddress')
    ..hasRequiredFields = false
  ;

  BlockHeaderInfo._() : super();
  factory BlockHeaderInfo({
    $core.int? version,
    $core.List<$core.int>? prevBlockHash,
    $core.List<$core.int>? stateRoot,
    $core.List<$core.int>? sortitionSeed,
    $core.String? proposerAddress,
  }) {
    final _result = create();
    if (version != null) {
      _result.version = version;
    }
    if (prevBlockHash != null) {
      _result.prevBlockHash = prevBlockHash;
    }
    if (stateRoot != null) {
      _result.stateRoot = stateRoot;
    }
    if (sortitionSeed != null) {
      _result.sortitionSeed = sortitionSeed;
    }
    if (proposerAddress != null) {
      _result.proposerAddress = proposerAddress;
    }
    return _result;
  }
  factory BlockHeaderInfo.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory BlockHeaderInfo.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  BlockHeaderInfo clone() => BlockHeaderInfo()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  BlockHeaderInfo copyWith(void Function(BlockHeaderInfo) updates) => super.copyWith((message) => updates(message as BlockHeaderInfo)) as BlockHeaderInfo; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static BlockHeaderInfo create() => BlockHeaderInfo._();
  BlockHeaderInfo createEmptyInstance() => create();
  static $pb.PbList<BlockHeaderInfo> createRepeated() => $pb.PbList<BlockHeaderInfo>();
  @$core.pragma('dart2js:noInline')
  static BlockHeaderInfo getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<BlockHeaderInfo>(create);
  static BlockHeaderInfo? _defaultInstance;

  @$pb.TagNumber(1)
  $core.int get version => $_getIZ(0);
  @$pb.TagNumber(1)
  set version($core.int v) { $_setSignedInt32(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasVersion() => $_has(0);
  @$pb.TagNumber(1)
  void clearVersion() => clearField(1);

  @$pb.TagNumber(2)
  $core.List<$core.int> get prevBlockHash => $_getN(1);
  @$pb.TagNumber(2)
  set prevBlockHash($core.List<$core.int> v) { $_setBytes(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasPrevBlockHash() => $_has(1);
  @$pb.TagNumber(2)
  void clearPrevBlockHash() => clearField(2);

  @$pb.TagNumber(3)
  $core.List<$core.int> get stateRoot => $_getN(2);
  @$pb.TagNumber(3)
  set stateRoot($core.List<$core.int> v) { $_setBytes(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasStateRoot() => $_has(2);
  @$pb.TagNumber(3)
  void clearStateRoot() => clearField(3);

  @$pb.TagNumber(4)
  $core.List<$core.int> get sortitionSeed => $_getN(3);
  @$pb.TagNumber(4)
  set sortitionSeed($core.List<$core.int> v) { $_setBytes(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasSortitionSeed() => $_has(3);
  @$pb.TagNumber(4)
  void clearSortitionSeed() => clearField(4);

  @$pb.TagNumber(5)
  $core.String get proposerAddress => $_getSZ(4);
  @$pb.TagNumber(5)
  set proposerAddress($core.String v) { $_setString(4, v); }
  @$pb.TagNumber(5)
  $core.bool hasProposerAddress() => $_has(4);
  @$pb.TagNumber(5)
  void clearProposerAddress() => clearField(5);
}

class CertificateInfo extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'CertificateInfo', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..a<$core.List<$core.int>>(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'hash', $pb.PbFieldType.OY)
    ..a<$core.int>(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'round', $pb.PbFieldType.O3)
    ..p<$core.int>(3, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'committers', $pb.PbFieldType.K3)
    ..p<$core.int>(4, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'absentees', $pb.PbFieldType.K3)
    ..a<$core.List<$core.int>>(5, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'signature', $pb.PbFieldType.OY)
    ..hasRequiredFields = false
  ;

  CertificateInfo._() : super();
  factory CertificateInfo({
    $core.List<$core.int>? hash,
    $core.int? round,
    $core.Iterable<$core.int>? committers,
    $core.Iterable<$core.int>? absentees,
    $core.List<$core.int>? signature,
  }) {
    final _result = create();
    if (hash != null) {
      _result.hash = hash;
    }
    if (round != null) {
      _result.round = round;
    }
    if (committers != null) {
      _result.committers.addAll(committers);
    }
    if (absentees != null) {
      _result.absentees.addAll(absentees);
    }
    if (signature != null) {
      _result.signature = signature;
    }
    return _result;
  }
  factory CertificateInfo.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory CertificateInfo.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  CertificateInfo clone() => CertificateInfo()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  CertificateInfo copyWith(void Function(CertificateInfo) updates) => super.copyWith((message) => updates(message as CertificateInfo)) as CertificateInfo; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static CertificateInfo create() => CertificateInfo._();
  CertificateInfo createEmptyInstance() => create();
  static $pb.PbList<CertificateInfo> createRepeated() => $pb.PbList<CertificateInfo>();
  @$core.pragma('dart2js:noInline')
  static CertificateInfo getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<CertificateInfo>(create);
  static CertificateInfo? _defaultInstance;

  @$pb.TagNumber(1)
  $core.List<$core.int> get hash => $_getN(0);
  @$pb.TagNumber(1)
  set hash($core.List<$core.int> v) { $_setBytes(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasHash() => $_has(0);
  @$pb.TagNumber(1)
  void clearHash() => clearField(1);

  @$pb.TagNumber(2)
  $core.int get round => $_getIZ(1);
  @$pb.TagNumber(2)
  set round($core.int v) { $_setSignedInt32(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasRound() => $_has(1);
  @$pb.TagNumber(2)
  void clearRound() => clearField(2);

  @$pb.TagNumber(3)
  $core.List<$core.int> get committers => $_getList(2);

  @$pb.TagNumber(4)
  $core.List<$core.int> get absentees => $_getList(3);

  @$pb.TagNumber(5)
  $core.List<$core.int> get signature => $_getN(4);
  @$pb.TagNumber(5)
  set signature($core.List<$core.int> v) { $_setBytes(4, v); }
  @$pb.TagNumber(5)
  $core.bool hasSignature() => $_has(4);
  @$pb.TagNumber(5)
  void clearSignature() => clearField(5);
}

class VoteInfo extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'VoteInfo', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..e<VoteType>(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'type', $pb.PbFieldType.OE, defaultOrMaker: VoteType.VOTE_UNKNOWN, valueOf: VoteType.valueOf, enumValues: VoteType.values)
    ..aOS(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'voter')
    ..a<$core.List<$core.int>>(3, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'blockHash', $pb.PbFieldType.OY)
    ..a<$core.int>(4, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'round', $pb.PbFieldType.O3)
    ..hasRequiredFields = false
  ;

  VoteInfo._() : super();
  factory VoteInfo({
    VoteType? type,
    $core.String? voter,
    $core.List<$core.int>? blockHash,
    $core.int? round,
  }) {
    final _result = create();
    if (type != null) {
      _result.type = type;
    }
    if (voter != null) {
      _result.voter = voter;
    }
    if (blockHash != null) {
      _result.blockHash = blockHash;
    }
    if (round != null) {
      _result.round = round;
    }
    return _result;
  }
  factory VoteInfo.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory VoteInfo.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  VoteInfo clone() => VoteInfo()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  VoteInfo copyWith(void Function(VoteInfo) updates) => super.copyWith((message) => updates(message as VoteInfo)) as VoteInfo; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static VoteInfo create() => VoteInfo._();
  VoteInfo createEmptyInstance() => create();
  static $pb.PbList<VoteInfo> createRepeated() => $pb.PbList<VoteInfo>();
  @$core.pragma('dart2js:noInline')
  static VoteInfo getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<VoteInfo>(create);
  static VoteInfo? _defaultInstance;

  @$pb.TagNumber(1)
  VoteType get type => $_getN(0);
  @$pb.TagNumber(1)
  set type(VoteType v) { setField(1, v); }
  @$pb.TagNumber(1)
  $core.bool hasType() => $_has(0);
  @$pb.TagNumber(1)
  void clearType() => clearField(1);

  @$pb.TagNumber(2)
  $core.String get voter => $_getSZ(1);
  @$pb.TagNumber(2)
  set voter($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasVoter() => $_has(1);
  @$pb.TagNumber(2)
  void clearVoter() => clearField(2);

  @$pb.TagNumber(3)
  $core.List<$core.int> get blockHash => $_getN(2);
  @$pb.TagNumber(3)
  set blockHash($core.List<$core.int> v) { $_setBytes(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasBlockHash() => $_has(2);
  @$pb.TagNumber(3)
  void clearBlockHash() => clearField(3);

  @$pb.TagNumber(4)
  $core.int get round => $_getIZ(3);
  @$pb.TagNumber(4)
  set round($core.int v) { $_setSignedInt32(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasRound() => $_has(3);
  @$pb.TagNumber(4)
  void clearRound() => clearField(4);
}

class ConsensusInfo extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'ConsensusInfo', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'address')
    ..aOB(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'Active', protoName: 'Active')
    ..a<$core.int>(3, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'height', $pb.PbFieldType.OU3)
    ..a<$core.int>(4, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'round', $pb.PbFieldType.O3)
    ..pc<VoteInfo>(5, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'votes', $pb.PbFieldType.PM, subBuilder: VoteInfo.create)
    ..hasRequiredFields = false
  ;

  ConsensusInfo._() : super();
  factory ConsensusInfo({
    $core.String? address,
    $core.bool? active,
    $core.int? height,
    $core.int? round,
    $core.Iterable<VoteInfo>? votes,
  }) {
    final _result = create();
    if (address != null) {
      _result.address = address;
    }
    if (active != null) {
      _result.active = active;
    }
    if (height != null) {
      _result.height = height;
    }
    if (round != null) {
      _result.round = round;
    }
    if (votes != null) {
      _result.votes.addAll(votes);
    }
    return _result;
  }
  factory ConsensusInfo.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory ConsensusInfo.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  ConsensusInfo clone() => ConsensusInfo()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  ConsensusInfo copyWith(void Function(ConsensusInfo) updates) => super.copyWith((message) => updates(message as ConsensusInfo)) as ConsensusInfo; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static ConsensusInfo create() => ConsensusInfo._();
  ConsensusInfo createEmptyInstance() => create();
  static $pb.PbList<ConsensusInfo> createRepeated() => $pb.PbList<ConsensusInfo>();
  @$core.pragma('dart2js:noInline')
  static ConsensusInfo getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ConsensusInfo>(create);
  static ConsensusInfo? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get address => $_getSZ(0);
  @$pb.TagNumber(1)
  set address($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasAddress() => $_has(0);
  @$pb.TagNumber(1)
  void clearAddress() => clearField(1);

  @$pb.TagNumber(2)
  $core.bool get active => $_getBF(1);
  @$pb.TagNumber(2)
  set active($core.bool v) { $_setBool(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasActive() => $_has(1);
  @$pb.TagNumber(2)
  void clearActive() => clearField(2);

  @$pb.TagNumber(3)
  $core.int get height => $_getIZ(2);
  @$pb.TagNumber(3)
  set height($core.int v) { $_setUnsignedInt32(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasHeight() => $_has(2);
  @$pb.TagNumber(3)
  void clearHeight() => clearField(3);

  @$pb.TagNumber(4)
  $core.int get round => $_getIZ(3);
  @$pb.TagNumber(4)
  set round($core.int v) { $_setSignedInt32(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasRound() => $_has(3);
  @$pb.TagNumber(4)
  void clearRound() => clearField(4);

  @$pb.TagNumber(5)
  $core.List<VoteInfo> get votes => $_getList(4);
}

class BlockchainApi {
  $pb.RpcClient _client;
  BlockchainApi(this._client);

  $async.Future<GetBlockResponse> getBlock($pb.ClientContext? ctx, GetBlockRequest request) {
    var emptyResponse = GetBlockResponse();
    return _client.invoke<GetBlockResponse>(ctx, 'Blockchain', 'GetBlock', request, emptyResponse);
  }
  $async.Future<GetBlockHashResponse> getBlockHash($pb.ClientContext? ctx, GetBlockHashRequest request) {
    var emptyResponse = GetBlockHashResponse();
    return _client.invoke<GetBlockHashResponse>(ctx, 'Blockchain', 'GetBlockHash', request, emptyResponse);
  }
  $async.Future<GetBlockHeightResponse> getBlockHeight($pb.ClientContext? ctx, GetBlockHeightRequest request) {
    var emptyResponse = GetBlockHeightResponse();
    return _client.invoke<GetBlockHeightResponse>(ctx, 'Blockchain', 'GetBlockHeight', request, emptyResponse);
  }
  $async.Future<GetBlockchainInfoResponse> getBlockchainInfo($pb.ClientContext? ctx, GetBlockchainInfoRequest request) {
    var emptyResponse = GetBlockchainInfoResponse();
    return _client.invoke<GetBlockchainInfoResponse>(ctx, 'Blockchain', 'GetBlockchainInfo', request, emptyResponse);
  }
  $async.Future<GetConsensusInfoResponse> getConsensusInfo($pb.ClientContext? ctx, GetConsensusInfoRequest request) {
    var emptyResponse = GetConsensusInfoResponse();
    return _client.invoke<GetConsensusInfoResponse>(ctx, 'Blockchain', 'GetConsensusInfo', request, emptyResponse);
  }
  $async.Future<GetAccountResponse> getAccount($pb.ClientContext? ctx, GetAccountRequest request) {
    var emptyResponse = GetAccountResponse();
    return _client.invoke<GetAccountResponse>(ctx, 'Blockchain', 'GetAccount', request, emptyResponse);
  }
  $async.Future<GetAccountResponse> getAccountByNumber($pb.ClientContext? ctx, GetAccountByNumberRequest request) {
    var emptyResponse = GetAccountResponse();
    return _client.invoke<GetAccountResponse>(ctx, 'Blockchain', 'GetAccountByNumber', request, emptyResponse);
  }
  $async.Future<GetValidatorResponse> getValidator($pb.ClientContext? ctx, GetValidatorRequest request) {
    var emptyResponse = GetValidatorResponse();
    return _client.invoke<GetValidatorResponse>(ctx, 'Blockchain', 'GetValidator', request, emptyResponse);
  }
  $async.Future<GetValidatorResponse> getValidatorByNumber($pb.ClientContext? ctx, GetValidatorByNumberRequest request) {
    var emptyResponse = GetValidatorResponse();
    return _client.invoke<GetValidatorResponse>(ctx, 'Blockchain', 'GetValidatorByNumber', request, emptyResponse);
  }
  $async.Future<GetValidatorAddressesResponse> getValidatorAddresses($pb.ClientContext? ctx, GetValidatorAddressesRequest request) {
    var emptyResponse = GetValidatorAddressesResponse();
    return _client.invoke<GetValidatorAddressesResponse>(ctx, 'Blockchain', 'GetValidatorAddresses', request, emptyResponse);
  }
}

