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

class AccountRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'AccountRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'address')
    ..hasRequiredFields = false
  ;

  AccountRequest._() : super();
  factory AccountRequest({
    $core.String? address,
  }) {
    final _result = create();
    if (address != null) {
      _result.address = address;
    }
    return _result;
  }
  factory AccountRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory AccountRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  AccountRequest clone() => AccountRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  AccountRequest copyWith(void Function(AccountRequest) updates) => super.copyWith((message) => updates(message as AccountRequest)) as AccountRequest; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static AccountRequest create() => AccountRequest._();
  AccountRequest createEmptyInstance() => create();
  static $pb.PbList<AccountRequest> createRepeated() => $pb.PbList<AccountRequest>();
  @$core.pragma('dart2js:noInline')
  static AccountRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<AccountRequest>(create);
  static AccountRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get address => $_getSZ(0);
  @$pb.TagNumber(1)
  set address($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasAddress() => $_has(0);
  @$pb.TagNumber(1)
  void clearAddress() => clearField(1);
}

class AccountResponse extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'AccountResponse', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOM<AccountInfo>(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'account', subBuilder: AccountInfo.create)
    ..hasRequiredFields = false
  ;

  AccountResponse._() : super();
  factory AccountResponse({
    AccountInfo? account,
  }) {
    final _result = create();
    if (account != null) {
      _result.account = account;
    }
    return _result;
  }
  factory AccountResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory AccountResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  AccountResponse clone() => AccountResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  AccountResponse copyWith(void Function(AccountResponse) updates) => super.copyWith((message) => updates(message as AccountResponse)) as AccountResponse; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static AccountResponse create() => AccountResponse._();
  AccountResponse createEmptyInstance() => create();
  static $pb.PbList<AccountResponse> createRepeated() => $pb.PbList<AccountResponse>();
  @$core.pragma('dart2js:noInline')
  static AccountResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<AccountResponse>(create);
  static AccountResponse? _defaultInstance;

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

class ValidatorsRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'ValidatorsRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..hasRequiredFields = false
  ;

  ValidatorsRequest._() : super();
  factory ValidatorsRequest() => create();
  factory ValidatorsRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory ValidatorsRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  ValidatorsRequest clone() => ValidatorsRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  ValidatorsRequest copyWith(void Function(ValidatorsRequest) updates) => super.copyWith((message) => updates(message as ValidatorsRequest)) as ValidatorsRequest; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static ValidatorsRequest create() => ValidatorsRequest._();
  ValidatorsRequest createEmptyInstance() => create();
  static $pb.PbList<ValidatorsRequest> createRepeated() => $pb.PbList<ValidatorsRequest>();
  @$core.pragma('dart2js:noInline')
  static ValidatorsRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ValidatorsRequest>(create);
  static ValidatorsRequest? _defaultInstance;
}

class ValidatorRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'ValidatorRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'address')
    ..hasRequiredFields = false
  ;

  ValidatorRequest._() : super();
  factory ValidatorRequest({
    $core.String? address,
  }) {
    final _result = create();
    if (address != null) {
      _result.address = address;
    }
    return _result;
  }
  factory ValidatorRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory ValidatorRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  ValidatorRequest clone() => ValidatorRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  ValidatorRequest copyWith(void Function(ValidatorRequest) updates) => super.copyWith((message) => updates(message as ValidatorRequest)) as ValidatorRequest; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static ValidatorRequest create() => ValidatorRequest._();
  ValidatorRequest createEmptyInstance() => create();
  static $pb.PbList<ValidatorRequest> createRepeated() => $pb.PbList<ValidatorRequest>();
  @$core.pragma('dart2js:noInline')
  static ValidatorRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ValidatorRequest>(create);
  static ValidatorRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get address => $_getSZ(0);
  @$pb.TagNumber(1)
  set address($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasAddress() => $_has(0);
  @$pb.TagNumber(1)
  void clearAddress() => clearField(1);
}

class ValidatorByNumberRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'ValidatorByNumberRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..a<$core.int>(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'number', $pb.PbFieldType.O3)
    ..hasRequiredFields = false
  ;

  ValidatorByNumberRequest._() : super();
  factory ValidatorByNumberRequest({
    $core.int? number,
  }) {
    final _result = create();
    if (number != null) {
      _result.number = number;
    }
    return _result;
  }
  factory ValidatorByNumberRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory ValidatorByNumberRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  ValidatorByNumberRequest clone() => ValidatorByNumberRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  ValidatorByNumberRequest copyWith(void Function(ValidatorByNumberRequest) updates) => super.copyWith((message) => updates(message as ValidatorByNumberRequest)) as ValidatorByNumberRequest; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static ValidatorByNumberRequest create() => ValidatorByNumberRequest._();
  ValidatorByNumberRequest createEmptyInstance() => create();
  static $pb.PbList<ValidatorByNumberRequest> createRepeated() => $pb.PbList<ValidatorByNumberRequest>();
  @$core.pragma('dart2js:noInline')
  static ValidatorByNumberRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ValidatorByNumberRequest>(create);
  static ValidatorByNumberRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.int get number => $_getIZ(0);
  @$pb.TagNumber(1)
  set number($core.int v) { $_setSignedInt32(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasNumber() => $_has(0);
  @$pb.TagNumber(1)
  void clearNumber() => clearField(1);
}

class ValidatorsResponse extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'ValidatorsResponse', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..pc<ValidatorInfo>(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'validators', $pb.PbFieldType.PM, subBuilder: ValidatorInfo.create)
    ..hasRequiredFields = false
  ;

  ValidatorsResponse._() : super();
  factory ValidatorsResponse({
    $core.Iterable<ValidatorInfo>? validators,
  }) {
    final _result = create();
    if (validators != null) {
      _result.validators.addAll(validators);
    }
    return _result;
  }
  factory ValidatorsResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory ValidatorsResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  ValidatorsResponse clone() => ValidatorsResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  ValidatorsResponse copyWith(void Function(ValidatorsResponse) updates) => super.copyWith((message) => updates(message as ValidatorsResponse)) as ValidatorsResponse; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static ValidatorsResponse create() => ValidatorsResponse._();
  ValidatorsResponse createEmptyInstance() => create();
  static $pb.PbList<ValidatorsResponse> createRepeated() => $pb.PbList<ValidatorsResponse>();
  @$core.pragma('dart2js:noInline')
  static ValidatorsResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ValidatorsResponse>(create);
  static ValidatorsResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.List<ValidatorInfo> get validators => $_getList(0);
}

class ValidatorResponse extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'ValidatorResponse', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOM<ValidatorInfo>(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'validator', subBuilder: ValidatorInfo.create)
    ..hasRequiredFields = false
  ;

  ValidatorResponse._() : super();
  factory ValidatorResponse({
    ValidatorInfo? validator,
  }) {
    final _result = create();
    if (validator != null) {
      _result.validator = validator;
    }
    return _result;
  }
  factory ValidatorResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory ValidatorResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  ValidatorResponse clone() => ValidatorResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  ValidatorResponse copyWith(void Function(ValidatorResponse) updates) => super.copyWith((message) => updates(message as ValidatorResponse)) as ValidatorResponse; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static ValidatorResponse create() => ValidatorResponse._();
  ValidatorResponse createEmptyInstance() => create();
  static $pb.PbList<ValidatorResponse> createRepeated() => $pb.PbList<ValidatorResponse>();
  @$core.pragma('dart2js:noInline')
  static ValidatorResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ValidatorResponse>(create);
  static ValidatorResponse? _defaultInstance;

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

class BlockRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'BlockRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..a<$core.int>(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'height', $pb.PbFieldType.OU3)
    ..e<BlockVerbosity>(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'verbosity', $pb.PbFieldType.OE, defaultOrMaker: BlockVerbosity.BLOCK_DATA, valueOf: BlockVerbosity.valueOf, enumValues: BlockVerbosity.values)
    ..hasRequiredFields = false
  ;

  BlockRequest._() : super();
  factory BlockRequest({
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
  factory BlockRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory BlockRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  BlockRequest clone() => BlockRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  BlockRequest copyWith(void Function(BlockRequest) updates) => super.copyWith((message) => updates(message as BlockRequest)) as BlockRequest; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static BlockRequest create() => BlockRequest._();
  BlockRequest createEmptyInstance() => create();
  static $pb.PbList<BlockRequest> createRepeated() => $pb.PbList<BlockRequest>();
  @$core.pragma('dart2js:noInline')
  static BlockRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<BlockRequest>(create);
  static BlockRequest? _defaultInstance;

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

class BlockResponse extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'BlockResponse', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..a<$core.int>(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'height', $pb.PbFieldType.OU3)
    ..a<$core.List<$core.int>>(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'hash', $pb.PbFieldType.OY)
    ..a<$core.List<$core.int>>(3, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'data', $pb.PbFieldType.OY)
    ..a<$core.int>(4, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'blockTime', $pb.PbFieldType.OU3)
    ..aOM<BlockHeaderInfo>(5, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'header', subBuilder: BlockHeaderInfo.create)
    ..aOM<CertificateInfo>(6, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'prevCert', subBuilder: CertificateInfo.create)
    ..pc<$0.TransactionInfo>(7, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'txs', $pb.PbFieldType.PM, subBuilder: $0.TransactionInfo.create)
    ..hasRequiredFields = false
  ;

  BlockResponse._() : super();
  factory BlockResponse({
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
  factory BlockResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory BlockResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  BlockResponse clone() => BlockResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  BlockResponse copyWith(void Function(BlockResponse) updates) => super.copyWith((message) => updates(message as BlockResponse)) as BlockResponse; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static BlockResponse create() => BlockResponse._();
  BlockResponse createEmptyInstance() => create();
  static $pb.PbList<BlockResponse> createRepeated() => $pb.PbList<BlockResponse>();
  @$core.pragma('dart2js:noInline')
  static BlockResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<BlockResponse>(create);
  static BlockResponse? _defaultInstance;

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

class BlockHashRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'BlockHashRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..a<$core.int>(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'height', $pb.PbFieldType.OU3)
    ..hasRequiredFields = false
  ;

  BlockHashRequest._() : super();
  factory BlockHashRequest({
    $core.int? height,
  }) {
    final _result = create();
    if (height != null) {
      _result.height = height;
    }
    return _result;
  }
  factory BlockHashRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory BlockHashRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  BlockHashRequest clone() => BlockHashRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  BlockHashRequest copyWith(void Function(BlockHashRequest) updates) => super.copyWith((message) => updates(message as BlockHashRequest)) as BlockHashRequest; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static BlockHashRequest create() => BlockHashRequest._();
  BlockHashRequest createEmptyInstance() => create();
  static $pb.PbList<BlockHashRequest> createRepeated() => $pb.PbList<BlockHashRequest>();
  @$core.pragma('dart2js:noInline')
  static BlockHashRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<BlockHashRequest>(create);
  static BlockHashRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.int get height => $_getIZ(0);
  @$pb.TagNumber(1)
  set height($core.int v) { $_setUnsignedInt32(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasHeight() => $_has(0);
  @$pb.TagNumber(1)
  void clearHeight() => clearField(1);
}

class BlockHashResponse extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'BlockHashResponse', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..a<$core.List<$core.int>>(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'hash', $pb.PbFieldType.OY)
    ..hasRequiredFields = false
  ;

  BlockHashResponse._() : super();
  factory BlockHashResponse({
    $core.List<$core.int>? hash,
  }) {
    final _result = create();
    if (hash != null) {
      _result.hash = hash;
    }
    return _result;
  }
  factory BlockHashResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory BlockHashResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  BlockHashResponse clone() => BlockHashResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  BlockHashResponse copyWith(void Function(BlockHashResponse) updates) => super.copyWith((message) => updates(message as BlockHashResponse)) as BlockHashResponse; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static BlockHashResponse create() => BlockHashResponse._();
  BlockHashResponse createEmptyInstance() => create();
  static $pb.PbList<BlockHashResponse> createRepeated() => $pb.PbList<BlockHashResponse>();
  @$core.pragma('dart2js:noInline')
  static BlockHashResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<BlockHashResponse>(create);
  static BlockHashResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.List<$core.int> get hash => $_getN(0);
  @$pb.TagNumber(1)
  set hash($core.List<$core.int> v) { $_setBytes(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasHash() => $_has(0);
  @$pb.TagNumber(1)
  void clearHash() => clearField(1);
}

class BlockHeightRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'BlockHeightRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..a<$core.List<$core.int>>(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'hash', $pb.PbFieldType.OY)
    ..hasRequiredFields = false
  ;

  BlockHeightRequest._() : super();
  factory BlockHeightRequest({
    $core.List<$core.int>? hash,
  }) {
    final _result = create();
    if (hash != null) {
      _result.hash = hash;
    }
    return _result;
  }
  factory BlockHeightRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory BlockHeightRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  BlockHeightRequest clone() => BlockHeightRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  BlockHeightRequest copyWith(void Function(BlockHeightRequest) updates) => super.copyWith((message) => updates(message as BlockHeightRequest)) as BlockHeightRequest; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static BlockHeightRequest create() => BlockHeightRequest._();
  BlockHeightRequest createEmptyInstance() => create();
  static $pb.PbList<BlockHeightRequest> createRepeated() => $pb.PbList<BlockHeightRequest>();
  @$core.pragma('dart2js:noInline')
  static BlockHeightRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<BlockHeightRequest>(create);
  static BlockHeightRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.List<$core.int> get hash => $_getN(0);
  @$pb.TagNumber(1)
  set hash($core.List<$core.int> v) { $_setBytes(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasHash() => $_has(0);
  @$pb.TagNumber(1)
  void clearHash() => clearField(1);
}

class BlockHeightResponse extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'BlockHeightResponse', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..a<$core.int>(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'height', $pb.PbFieldType.OU3)
    ..hasRequiredFields = false
  ;

  BlockHeightResponse._() : super();
  factory BlockHeightResponse({
    $core.int? height,
  }) {
    final _result = create();
    if (height != null) {
      _result.height = height;
    }
    return _result;
  }
  factory BlockHeightResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory BlockHeightResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  BlockHeightResponse clone() => BlockHeightResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  BlockHeightResponse copyWith(void Function(BlockHeightResponse) updates) => super.copyWith((message) => updates(message as BlockHeightResponse)) as BlockHeightResponse; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static BlockHeightResponse create() => BlockHeightResponse._();
  BlockHeightResponse createEmptyInstance() => create();
  static $pb.PbList<BlockHeightResponse> createRepeated() => $pb.PbList<BlockHeightResponse>();
  @$core.pragma('dart2js:noInline')
  static BlockHeightResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<BlockHeightResponse>(create);
  static BlockHeightResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.int get height => $_getIZ(0);
  @$pb.TagNumber(1)
  set height($core.int v) { $_setUnsignedInt32(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasHeight() => $_has(0);
  @$pb.TagNumber(1)
  void clearHeight() => clearField(1);
}

class BlockchainInfoRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'BlockchainInfoRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..hasRequiredFields = false
  ;

  BlockchainInfoRequest._() : super();
  factory BlockchainInfoRequest() => create();
  factory BlockchainInfoRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory BlockchainInfoRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  BlockchainInfoRequest clone() => BlockchainInfoRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  BlockchainInfoRequest copyWith(void Function(BlockchainInfoRequest) updates) => super.copyWith((message) => updates(message as BlockchainInfoRequest)) as BlockchainInfoRequest; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static BlockchainInfoRequest create() => BlockchainInfoRequest._();
  BlockchainInfoRequest createEmptyInstance() => create();
  static $pb.PbList<BlockchainInfoRequest> createRepeated() => $pb.PbList<BlockchainInfoRequest>();
  @$core.pragma('dart2js:noInline')
  static BlockchainInfoRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<BlockchainInfoRequest>(create);
  static BlockchainInfoRequest? _defaultInstance;
}

class BlockchainInfoResponse extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'BlockchainInfoResponse', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..a<$core.int>(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'lastBlockHeight', $pb.PbFieldType.OU3)
    ..a<$core.List<$core.int>>(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'lastBlockHash', $pb.PbFieldType.OY)
    ..hasRequiredFields = false
  ;

  BlockchainInfoResponse._() : super();
  factory BlockchainInfoResponse({
    $core.int? lastBlockHeight,
    $core.List<$core.int>? lastBlockHash,
  }) {
    final _result = create();
    if (lastBlockHeight != null) {
      _result.lastBlockHeight = lastBlockHeight;
    }
    if (lastBlockHash != null) {
      _result.lastBlockHash = lastBlockHash;
    }
    return _result;
  }
  factory BlockchainInfoResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory BlockchainInfoResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  BlockchainInfoResponse clone() => BlockchainInfoResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  BlockchainInfoResponse copyWith(void Function(BlockchainInfoResponse) updates) => super.copyWith((message) => updates(message as BlockchainInfoResponse)) as BlockchainInfoResponse; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static BlockchainInfoResponse create() => BlockchainInfoResponse._();
  BlockchainInfoResponse createEmptyInstance() => create();
  static $pb.PbList<BlockchainInfoResponse> createRepeated() => $pb.PbList<BlockchainInfoResponse>();
  @$core.pragma('dart2js:noInline')
  static BlockchainInfoResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<BlockchainInfoResponse>(create);
  static BlockchainInfoResponse? _defaultInstance;

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
}

class ValidatorInfo extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'ValidatorInfo', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'publicKey')
    ..a<$core.int>(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'number', $pb.PbFieldType.O3)
    ..a<$core.int>(3, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'sequence', $pb.PbFieldType.O3)
    ..aInt64(4, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'stake')
    ..a<$core.int>(5, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'lastBondingHeight', $pb.PbFieldType.OU3)
    ..a<$core.int>(6, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'lastJoinedHeight', $pb.PbFieldType.OU3)
    ..a<$core.int>(7, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'unbondingHeight', $pb.PbFieldType.OU3)
    ..aOS(8, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'address')
    ..hasRequiredFields = false
  ;

  ValidatorInfo._() : super();
  factory ValidatorInfo({
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
  $core.String get publicKey => $_getSZ(0);
  @$pb.TagNumber(1)
  set publicKey($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasPublicKey() => $_has(0);
  @$pb.TagNumber(1)
  void clearPublicKey() => clearField(1);

  @$pb.TagNumber(2)
  $core.int get number => $_getIZ(1);
  @$pb.TagNumber(2)
  set number($core.int v) { $_setSignedInt32(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasNumber() => $_has(1);
  @$pb.TagNumber(2)
  void clearNumber() => clearField(2);

  @$pb.TagNumber(3)
  $core.int get sequence => $_getIZ(2);
  @$pb.TagNumber(3)
  set sequence($core.int v) { $_setSignedInt32(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasSequence() => $_has(2);
  @$pb.TagNumber(3)
  void clearSequence() => clearField(3);

  @$pb.TagNumber(4)
  $fixnum.Int64 get stake => $_getI64(3);
  @$pb.TagNumber(4)
  set stake($fixnum.Int64 v) { $_setInt64(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasStake() => $_has(3);
  @$pb.TagNumber(4)
  void clearStake() => clearField(4);

  @$pb.TagNumber(5)
  $core.int get lastBondingHeight => $_getIZ(4);
  @$pb.TagNumber(5)
  set lastBondingHeight($core.int v) { $_setUnsignedInt32(4, v); }
  @$pb.TagNumber(5)
  $core.bool hasLastBondingHeight() => $_has(4);
  @$pb.TagNumber(5)
  void clearLastBondingHeight() => clearField(5);

  @$pb.TagNumber(6)
  $core.int get lastJoinedHeight => $_getIZ(5);
  @$pb.TagNumber(6)
  set lastJoinedHeight($core.int v) { $_setUnsignedInt32(5, v); }
  @$pb.TagNumber(6)
  $core.bool hasLastJoinedHeight() => $_has(5);
  @$pb.TagNumber(6)
  void clearLastJoinedHeight() => clearField(6);

  @$pb.TagNumber(7)
  $core.int get unbondingHeight => $_getIZ(6);
  @$pb.TagNumber(7)
  set unbondingHeight($core.int v) { $_setUnsignedInt32(6, v); }
  @$pb.TagNumber(7)
  $core.bool hasUnbondingHeight() => $_has(6);
  @$pb.TagNumber(7)
  void clearUnbondingHeight() => clearField(7);

  @$pb.TagNumber(8)
  $core.String get address => $_getSZ(7);
  @$pb.TagNumber(8)
  set address($core.String v) { $_setString(7, v); }
  @$pb.TagNumber(8)
  $core.bool hasAddress() => $_has(7);
  @$pb.TagNumber(8)
  void clearAddress() => clearField(8);
}

class AccountInfo extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'AccountInfo', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'address')
    ..a<$core.int>(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'number', $pb.PbFieldType.O3)
    ..a<$core.int>(3, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'sequence', $pb.PbFieldType.O3)
    ..aInt64(4, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'Balance', protoName: 'Balance')
    ..hasRequiredFields = false
  ;

  AccountInfo._() : super();
  factory AccountInfo({
    $core.String? address,
    $core.int? number,
    $core.int? sequence,
    $fixnum.Int64? balance,
  }) {
    final _result = create();
    if (address != null) {
      _result.address = address;
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
  $core.String get address => $_getSZ(0);
  @$pb.TagNumber(1)
  set address($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasAddress() => $_has(0);
  @$pb.TagNumber(1)
  void clearAddress() => clearField(1);

  @$pb.TagNumber(2)
  $core.int get number => $_getIZ(1);
  @$pb.TagNumber(2)
  set number($core.int v) { $_setSignedInt32(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasNumber() => $_has(1);
  @$pb.TagNumber(2)
  void clearNumber() => clearField(2);

  @$pb.TagNumber(3)
  $core.int get sequence => $_getIZ(2);
  @$pb.TagNumber(3)
  set sequence($core.int v) { $_setSignedInt32(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasSequence() => $_has(2);
  @$pb.TagNumber(3)
  void clearSequence() => clearField(3);

  @$pb.TagNumber(4)
  $fixnum.Int64 get balance => $_getI64(3);
  @$pb.TagNumber(4)
  set balance($fixnum.Int64 v) { $_setInt64(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasBalance() => $_has(3);
  @$pb.TagNumber(4)
  void clearBalance() => clearField(4);
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
    ..a<$core.int>(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'round', $pb.PbFieldType.O3)
    ..p<$core.int>(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'committers', $pb.PbFieldType.K3)
    ..p<$core.int>(3, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'absentees', $pb.PbFieldType.K3)
    ..a<$core.List<$core.int>>(4, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'signature', $pb.PbFieldType.OY)
    ..hasRequiredFields = false
  ;

  CertificateInfo._() : super();
  factory CertificateInfo({
    $core.int? round,
    $core.Iterable<$core.int>? committers,
    $core.Iterable<$core.int>? absentees,
    $core.List<$core.int>? signature,
  }) {
    final _result = create();
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
  $core.int get round => $_getIZ(0);
  @$pb.TagNumber(1)
  set round($core.int v) { $_setSignedInt32(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasRound() => $_has(0);
  @$pb.TagNumber(1)
  void clearRound() => clearField(1);

  @$pb.TagNumber(2)
  $core.List<$core.int> get committers => $_getList(1);

  @$pb.TagNumber(3)
  $core.List<$core.int> get absentees => $_getList(2);

  @$pb.TagNumber(4)
  $core.List<$core.int> get signature => $_getN(3);
  @$pb.TagNumber(4)
  set signature($core.List<$core.int> v) { $_setBytes(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasSignature() => $_has(3);
  @$pb.TagNumber(4)
  void clearSignature() => clearField(4);
}

class BlockchainApi {
  $pb.RpcClient _client;
  BlockchainApi(this._client);

  $async.Future<BlockResponse> getBlock($pb.ClientContext? ctx, BlockRequest request) {
    var emptyResponse = BlockResponse();
    return _client.invoke<BlockResponse>(ctx, 'Blockchain', 'GetBlock', request, emptyResponse);
  }
  $async.Future<BlockHashResponse> getBlockHash($pb.ClientContext? ctx, BlockHashRequest request) {
    var emptyResponse = BlockHashResponse();
    return _client.invoke<BlockHashResponse>(ctx, 'Blockchain', 'GetBlockHash', request, emptyResponse);
  }
  $async.Future<BlockHeightResponse> getBlockHeight($pb.ClientContext? ctx, BlockHeightRequest request) {
    var emptyResponse = BlockHeightResponse();
    return _client.invoke<BlockHeightResponse>(ctx, 'Blockchain', 'GetBlockHeight', request, emptyResponse);
  }
  $async.Future<AccountResponse> getAccount($pb.ClientContext? ctx, AccountRequest request) {
    var emptyResponse = AccountResponse();
    return _client.invoke<AccountResponse>(ctx, 'Blockchain', 'GetAccount', request, emptyResponse);
  }
  $async.Future<ValidatorsResponse> getValidators($pb.ClientContext? ctx, ValidatorsRequest request) {
    var emptyResponse = ValidatorsResponse();
    return _client.invoke<ValidatorsResponse>(ctx, 'Blockchain', 'GetValidators', request, emptyResponse);
  }
  $async.Future<ValidatorResponse> getValidator($pb.ClientContext? ctx, ValidatorRequest request) {
    var emptyResponse = ValidatorResponse();
    return _client.invoke<ValidatorResponse>(ctx, 'Blockchain', 'GetValidator', request, emptyResponse);
  }
  $async.Future<ValidatorResponse> getValidatorByNumber($pb.ClientContext? ctx, ValidatorByNumberRequest request) {
    var emptyResponse = ValidatorResponse();
    return _client.invoke<ValidatorResponse>(ctx, 'Blockchain', 'GetValidatorByNumber', request, emptyResponse);
  }
  $async.Future<BlockchainInfoResponse> getBlockchainInfo($pb.ClientContext? ctx, BlockchainInfoRequest request) {
    var emptyResponse = BlockchainInfoResponse();
    return _client.invoke<BlockchainInfoResponse>(ctx, 'Blockchain', 'GetBlockchainInfo', request, emptyResponse);
  }
}

