///
//  Generated code. Do not modify.
//  source: wallet.proto
//
// @dart = 2.12
// ignore_for_file: annotate_overrides,camel_case_types,constant_identifier_names,directives_ordering,library_prefixes,non_constant_identifier_names,prefer_final_fields,return_of_invalid_type,unnecessary_const,unnecessary_import,unnecessary_this,unused_import,unused_shown_name

import 'dart:async' as $async;
import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

class CreateWalletRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'CreateWalletRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'name')
    ..aOS(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'mnemonic')
    ..aOS(3, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'language')
    ..aOS(4, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'password')
    ..hasRequiredFields = false
  ;

  CreateWalletRequest._() : super();
  factory CreateWalletRequest({
    $core.String? name,
    $core.String? mnemonic,
    $core.String? language,
    $core.String? password,
  }) {
    final _result = create();
    if (name != null) {
      _result.name = name;
    }
    if (mnemonic != null) {
      _result.mnemonic = mnemonic;
    }
    if (language != null) {
      _result.language = language;
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
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(1)
  set name($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(1)
  void clearName() => clearField(1);

  @$pb.TagNumber(2)
  $core.String get mnemonic => $_getSZ(1);
  @$pb.TagNumber(2)
  set mnemonic($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasMnemonic() => $_has(1);
  @$pb.TagNumber(2)
  void clearMnemonic() => clearField(2);

  @$pb.TagNumber(3)
  $core.String get language => $_getSZ(2);
  @$pb.TagNumber(3)
  set language($core.String v) { $_setString(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasLanguage() => $_has(2);
  @$pb.TagNumber(3)
  void clearLanguage() => clearField(3);

  @$pb.TagNumber(4)
  $core.String get password => $_getSZ(3);
  @$pb.TagNumber(4)
  set password($core.String v) { $_setString(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasPassword() => $_has(3);
  @$pb.TagNumber(4)
  void clearPassword() => clearField(4);
}

class CreateWalletResponse extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'CreateWalletResponse', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..hasRequiredFields = false
  ;

  CreateWalletResponse._() : super();
  factory CreateWalletResponse() => create();
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
}

class LoadWalletRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'LoadWalletRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'name')
    ..hasRequiredFields = false
  ;

  LoadWalletRequest._() : super();
  factory LoadWalletRequest({
    $core.String? name,
  }) {
    final _result = create();
    if (name != null) {
      _result.name = name;
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
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(1)
  set name($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(1)
  void clearName() => clearField(1);
}

class LoadWalletResponse extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'LoadWalletResponse', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'name')
    ..hasRequiredFields = false
  ;

  LoadWalletResponse._() : super();
  factory LoadWalletResponse({
    $core.String? name,
  }) {
    final _result = create();
    if (name != null) {
      _result.name = name;
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
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(1)
  set name($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(1)
  void clearName() => clearField(1);
}

class UnloadWalletRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'UnloadWalletRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'name')
    ..hasRequiredFields = false
  ;

  UnloadWalletRequest._() : super();
  factory UnloadWalletRequest({
    $core.String? name,
  }) {
    final _result = create();
    if (name != null) {
      _result.name = name;
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
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(1)
  set name($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(1)
  void clearName() => clearField(1);
}

class UnloadWalletResponse extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'UnloadWalletResponse', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'name')
    ..hasRequiredFields = false
  ;

  UnloadWalletResponse._() : super();
  factory UnloadWalletResponse({
    $core.String? name,
  }) {
    final _result = create();
    if (name != null) {
      _result.name = name;
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
  $core.String get name => $_getSZ(0);
  @$pb.TagNumber(1)
  set name($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasName() => $_has(0);
  @$pb.TagNumber(1)
  void clearName() => clearField(1);
}

class LockWalletRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'LockWalletRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'password')
    ..a<$core.int>(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'timeout', $pb.PbFieldType.O3)
    ..hasRequiredFields = false
  ;

  LockWalletRequest._() : super();
  factory LockWalletRequest({
    $core.String? password,
    $core.int? timeout,
  }) {
    final _result = create();
    if (password != null) {
      _result.password = password;
    }
    if (timeout != null) {
      _result.timeout = timeout;
    }
    return _result;
  }
  factory LockWalletRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory LockWalletRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  LockWalletRequest clone() => LockWalletRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  LockWalletRequest copyWith(void Function(LockWalletRequest) updates) => super.copyWith((message) => updates(message as LockWalletRequest)) as LockWalletRequest; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static LockWalletRequest create() => LockWalletRequest._();
  LockWalletRequest createEmptyInstance() => create();
  static $pb.PbList<LockWalletRequest> createRepeated() => $pb.PbList<LockWalletRequest>();
  @$core.pragma('dart2js:noInline')
  static LockWalletRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<LockWalletRequest>(create);
  static LockWalletRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get password => $_getSZ(0);
  @$pb.TagNumber(1)
  set password($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasPassword() => $_has(0);
  @$pb.TagNumber(1)
  void clearPassword() => clearField(1);

  @$pb.TagNumber(2)
  $core.int get timeout => $_getIZ(1);
  @$pb.TagNumber(2)
  set timeout($core.int v) { $_setSignedInt32(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasTimeout() => $_has(1);
  @$pb.TagNumber(2)
  void clearTimeout() => clearField(2);
}

class LockWalletResponse extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'LockWalletResponse', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..hasRequiredFields = false
  ;

  LockWalletResponse._() : super();
  factory LockWalletResponse() => create();
  factory LockWalletResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory LockWalletResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  LockWalletResponse clone() => LockWalletResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  LockWalletResponse copyWith(void Function(LockWalletResponse) updates) => super.copyWith((message) => updates(message as LockWalletResponse)) as LockWalletResponse; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static LockWalletResponse create() => LockWalletResponse._();
  LockWalletResponse createEmptyInstance() => create();
  static $pb.PbList<LockWalletResponse> createRepeated() => $pb.PbList<LockWalletResponse>();
  @$core.pragma('dart2js:noInline')
  static LockWalletResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<LockWalletResponse>(create);
  static LockWalletResponse? _defaultInstance;
}

class UnlockWalletRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'UnlockWalletRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'password')
    ..a<$core.int>(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'timeout', $pb.PbFieldType.O3)
    ..hasRequiredFields = false
  ;

  UnlockWalletRequest._() : super();
  factory UnlockWalletRequest({
    $core.String? password,
    $core.int? timeout,
  }) {
    final _result = create();
    if (password != null) {
      _result.password = password;
    }
    if (timeout != null) {
      _result.timeout = timeout;
    }
    return _result;
  }
  factory UnlockWalletRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory UnlockWalletRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  UnlockWalletRequest clone() => UnlockWalletRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  UnlockWalletRequest copyWith(void Function(UnlockWalletRequest) updates) => super.copyWith((message) => updates(message as UnlockWalletRequest)) as UnlockWalletRequest; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static UnlockWalletRequest create() => UnlockWalletRequest._();
  UnlockWalletRequest createEmptyInstance() => create();
  static $pb.PbList<UnlockWalletRequest> createRepeated() => $pb.PbList<UnlockWalletRequest>();
  @$core.pragma('dart2js:noInline')
  static UnlockWalletRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<UnlockWalletRequest>(create);
  static UnlockWalletRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get password => $_getSZ(0);
  @$pb.TagNumber(1)
  set password($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasPassword() => $_has(0);
  @$pb.TagNumber(1)
  void clearPassword() => clearField(1);

  @$pb.TagNumber(2)
  $core.int get timeout => $_getIZ(1);
  @$pb.TagNumber(2)
  set timeout($core.int v) { $_setSignedInt32(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasTimeout() => $_has(1);
  @$pb.TagNumber(2)
  void clearTimeout() => clearField(2);
}

class UnlockWalletResponse extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'UnlockWalletResponse', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..hasRequiredFields = false
  ;

  UnlockWalletResponse._() : super();
  factory UnlockWalletResponse() => create();
  factory UnlockWalletResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory UnlockWalletResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  UnlockWalletResponse clone() => UnlockWalletResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  UnlockWalletResponse copyWith(void Function(UnlockWalletResponse) updates) => super.copyWith((message) => updates(message as UnlockWalletResponse)) as UnlockWalletResponse; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static UnlockWalletResponse create() => UnlockWalletResponse._();
  UnlockWalletResponse createEmptyInstance() => create();
  static $pb.PbList<UnlockWalletResponse> createRepeated() => $pb.PbList<UnlockWalletResponse>();
  @$core.pragma('dart2js:noInline')
  static UnlockWalletResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<UnlockWalletResponse>(create);
  static UnlockWalletResponse? _defaultInstance;
}

class WalletApi {
  $pb.RpcClient _client;
  WalletApi(this._client);

  $async.Future<CreateWalletResponse> createWallet($pb.ClientContext? ctx, CreateWalletRequest request) {
    var emptyResponse = CreateWalletResponse();
    return _client.invoke<CreateWalletResponse>(ctx, 'Wallet', 'CreateWallet', request, emptyResponse);
  }
  $async.Future<LoadWalletResponse> loadWallet($pb.ClientContext? ctx, LoadWalletRequest request) {
    var emptyResponse = LoadWalletResponse();
    return _client.invoke<LoadWalletResponse>(ctx, 'Wallet', 'LoadWallet', request, emptyResponse);
  }
  $async.Future<UnloadWalletResponse> unloadWallet($pb.ClientContext? ctx, UnloadWalletRequest request) {
    var emptyResponse = UnloadWalletResponse();
    return _client.invoke<UnloadWalletResponse>(ctx, 'Wallet', 'UnloadWallet', request, emptyResponse);
  }
  $async.Future<LockWalletResponse> lockWallet($pb.ClientContext? ctx, LockWalletRequest request) {
    var emptyResponse = LockWalletResponse();
    return _client.invoke<LockWalletResponse>(ctx, 'Wallet', 'LockWallet', request, emptyResponse);
  }
  $async.Future<UnlockWalletResponse> unlockWallet($pb.ClientContext? ctx, UnlockWalletRequest request) {
    var emptyResponse = UnlockWalletResponse();
    return _client.invoke<UnlockWalletResponse>(ctx, 'Wallet', 'UnlockWallet', request, emptyResponse);
  }
}

