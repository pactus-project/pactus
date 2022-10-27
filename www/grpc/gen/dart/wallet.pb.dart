///
//  Generated code. Do not modify.
//  source: wallet.proto
//
// @dart = 2.12
// ignore_for_file: annotate_overrides,camel_case_types,constant_identifier_names,directives_ordering,library_prefixes,non_constant_identifier_names,prefer_final_fields,return_of_invalid_type,unnecessary_const,unnecessary_import,unnecessary_this,unused_import,unused_shown_name

import 'dart:async' as $async;
import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

class GenerateMnemonicRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'GenerateMnemonicRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..a<$core.int>(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'entropy', $pb.PbFieldType.O3)
    ..aOS(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'language')
    ..hasRequiredFields = false
  ;

  GenerateMnemonicRequest._() : super();
  factory GenerateMnemonicRequest({
    $core.int? entropy,
    $core.String? language,
  }) {
    final _result = create();
    if (entropy != null) {
      _result.entropy = entropy;
    }
    if (language != null) {
      _result.language = language;
    }
    return _result;
  }
  factory GenerateMnemonicRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GenerateMnemonicRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GenerateMnemonicRequest clone() => GenerateMnemonicRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GenerateMnemonicRequest copyWith(void Function(GenerateMnemonicRequest) updates) => super.copyWith((message) => updates(message as GenerateMnemonicRequest)) as GenerateMnemonicRequest; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static GenerateMnemonicRequest create() => GenerateMnemonicRequest._();
  GenerateMnemonicRequest createEmptyInstance() => create();
  static $pb.PbList<GenerateMnemonicRequest> createRepeated() => $pb.PbList<GenerateMnemonicRequest>();
  @$core.pragma('dart2js:noInline')
  static GenerateMnemonicRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GenerateMnemonicRequest>(create);
  static GenerateMnemonicRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.int get entropy => $_getIZ(0);
  @$pb.TagNumber(1)
  set entropy($core.int v) { $_setSignedInt32(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasEntropy() => $_has(0);
  @$pb.TagNumber(1)
  void clearEntropy() => clearField(1);

  @$pb.TagNumber(2)
  $core.String get language => $_getSZ(1);
  @$pb.TagNumber(2)
  set language($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasLanguage() => $_has(1);
  @$pb.TagNumber(2)
  void clearLanguage() => clearField(2);
}

class GenerateMnemonicResponse extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'GenerateMnemonicResponse', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'mnemonic')
    ..hasRequiredFields = false
  ;

  GenerateMnemonicResponse._() : super();
  factory GenerateMnemonicResponse({
    $core.String? mnemonic,
  }) {
    final _result = create();
    if (mnemonic != null) {
      _result.mnemonic = mnemonic;
    }
    return _result;
  }
  factory GenerateMnemonicResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GenerateMnemonicResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GenerateMnemonicResponse clone() => GenerateMnemonicResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GenerateMnemonicResponse copyWith(void Function(GenerateMnemonicResponse) updates) => super.copyWith((message) => updates(message as GenerateMnemonicResponse)) as GenerateMnemonicResponse; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static GenerateMnemonicResponse create() => GenerateMnemonicResponse._();
  GenerateMnemonicResponse createEmptyInstance() => create();
  static $pb.PbList<GenerateMnemonicResponse> createRepeated() => $pb.PbList<GenerateMnemonicResponse>();
  @$core.pragma('dart2js:noInline')
  static GenerateMnemonicResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GenerateMnemonicResponse>(create);
  static GenerateMnemonicResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get mnemonic => $_getSZ(0);
  @$pb.TagNumber(1)
  set mnemonic($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasMnemonic() => $_has(0);
  @$pb.TagNumber(1)
  void clearMnemonic() => clearField(1);
}

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

class WalletApi {
  $pb.RpcClient _client;
  WalletApi(this._client);

  $async.Future<GenerateMnemonicResponse> generateMnemonic($pb.ClientContext? ctx, GenerateMnemonicRequest request) {
    var emptyResponse = GenerateMnemonicResponse();
    return _client.invoke<GenerateMnemonicResponse>(ctx, 'Wallet', 'GenerateMnemonic', request, emptyResponse);
  }
  $async.Future<CreateWalletResponse> createWallet($pb.ClientContext? ctx, CreateWalletRequest request) {
    var emptyResponse = CreateWalletResponse();
    return _client.invoke<CreateWalletResponse>(ctx, 'Wallet', 'CreateWallet', request, emptyResponse);
  }
}

