///
//  Generated code. Do not modify.
//  source: utils.proto
//
// @dart = 2.12
// ignore_for_file: annotate_overrides,camel_case_types,constant_identifier_names,directives_ordering,library_prefixes,non_constant_identifier_names,prefer_final_fields,return_of_invalid_type,unnecessary_const,unnecessary_import,unnecessary_this,unused_import,unused_shown_name

import 'dart:async' as $async;
import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

class SignMessageWithPrivateKeyRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'SignMessageWithPrivateKeyRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'privateKey')
    ..aOS(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'message')
    ..hasRequiredFields = false
  ;

  SignMessageWithPrivateKeyRequest._() : super();
  factory SignMessageWithPrivateKeyRequest({
    $core.String? privateKey,
    $core.String? message,
  }) {
    final _result = create();
    if (privateKey != null) {
      _result.privateKey = privateKey;
    }
    if (message != null) {
      _result.message = message;
    }
    return _result;
  }
  factory SignMessageWithPrivateKeyRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory SignMessageWithPrivateKeyRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  SignMessageWithPrivateKeyRequest clone() => SignMessageWithPrivateKeyRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  SignMessageWithPrivateKeyRequest copyWith(void Function(SignMessageWithPrivateKeyRequest) updates) => super.copyWith((message) => updates(message as SignMessageWithPrivateKeyRequest)) as SignMessageWithPrivateKeyRequest; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static SignMessageWithPrivateKeyRequest create() => SignMessageWithPrivateKeyRequest._();
  SignMessageWithPrivateKeyRequest createEmptyInstance() => create();
  static $pb.PbList<SignMessageWithPrivateKeyRequest> createRepeated() => $pb.PbList<SignMessageWithPrivateKeyRequest>();
  @$core.pragma('dart2js:noInline')
  static SignMessageWithPrivateKeyRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<SignMessageWithPrivateKeyRequest>(create);
  static SignMessageWithPrivateKeyRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get privateKey => $_getSZ(0);
  @$pb.TagNumber(1)
  set privateKey($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasPrivateKey() => $_has(0);
  @$pb.TagNumber(1)
  void clearPrivateKey() => clearField(1);

  @$pb.TagNumber(2)
  $core.String get message => $_getSZ(1);
  @$pb.TagNumber(2)
  set message($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasMessage() => $_has(1);
  @$pb.TagNumber(2)
  void clearMessage() => clearField(2);
}

class SignMessageWithPrivateKeyResponse extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'SignMessageWithPrivateKeyResponse', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'signature')
    ..hasRequiredFields = false
  ;

  SignMessageWithPrivateKeyResponse._() : super();
  factory SignMessageWithPrivateKeyResponse({
    $core.String? signature,
  }) {
    final _result = create();
    if (signature != null) {
      _result.signature = signature;
    }
    return _result;
  }
  factory SignMessageWithPrivateKeyResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory SignMessageWithPrivateKeyResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  SignMessageWithPrivateKeyResponse clone() => SignMessageWithPrivateKeyResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  SignMessageWithPrivateKeyResponse copyWith(void Function(SignMessageWithPrivateKeyResponse) updates) => super.copyWith((message) => updates(message as SignMessageWithPrivateKeyResponse)) as SignMessageWithPrivateKeyResponse; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static SignMessageWithPrivateKeyResponse create() => SignMessageWithPrivateKeyResponse._();
  SignMessageWithPrivateKeyResponse createEmptyInstance() => create();
  static $pb.PbList<SignMessageWithPrivateKeyResponse> createRepeated() => $pb.PbList<SignMessageWithPrivateKeyResponse>();
  @$core.pragma('dart2js:noInline')
  static SignMessageWithPrivateKeyResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<SignMessageWithPrivateKeyResponse>(create);
  static SignMessageWithPrivateKeyResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get signature => $_getSZ(0);
  @$pb.TagNumber(1)
  set signature($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasSignature() => $_has(0);
  @$pb.TagNumber(1)
  void clearSignature() => clearField(1);
}

class VerifyMessageRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'VerifyMessageRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'message')
    ..aOS(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'signature')
    ..aOS(3, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'publicKey')
    ..hasRequiredFields = false
  ;

  VerifyMessageRequest._() : super();
  factory VerifyMessageRequest({
    $core.String? message,
    $core.String? signature,
    $core.String? publicKey,
  }) {
    final _result = create();
    if (message != null) {
      _result.message = message;
    }
    if (signature != null) {
      _result.signature = signature;
    }
    if (publicKey != null) {
      _result.publicKey = publicKey;
    }
    return _result;
  }
  factory VerifyMessageRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory VerifyMessageRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  VerifyMessageRequest clone() => VerifyMessageRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  VerifyMessageRequest copyWith(void Function(VerifyMessageRequest) updates) => super.copyWith((message) => updates(message as VerifyMessageRequest)) as VerifyMessageRequest; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static VerifyMessageRequest create() => VerifyMessageRequest._();
  VerifyMessageRequest createEmptyInstance() => create();
  static $pb.PbList<VerifyMessageRequest> createRepeated() => $pb.PbList<VerifyMessageRequest>();
  @$core.pragma('dart2js:noInline')
  static VerifyMessageRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<VerifyMessageRequest>(create);
  static VerifyMessageRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(1)
  set message($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(1)
  void clearMessage() => clearField(1);

  @$pb.TagNumber(2)
  $core.String get signature => $_getSZ(1);
  @$pb.TagNumber(2)
  set signature($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasSignature() => $_has(1);
  @$pb.TagNumber(2)
  void clearSignature() => clearField(2);

  @$pb.TagNumber(3)
  $core.String get publicKey => $_getSZ(2);
  @$pb.TagNumber(3)
  set publicKey($core.String v) { $_setString(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasPublicKey() => $_has(2);
  @$pb.TagNumber(3)
  void clearPublicKey() => clearField(3);
}

class VerifyMessageResponse extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'VerifyMessageResponse', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOB(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'isValid')
    ..hasRequiredFields = false
  ;

  VerifyMessageResponse._() : super();
  factory VerifyMessageResponse({
    $core.bool? isValid,
  }) {
    final _result = create();
    if (isValid != null) {
      _result.isValid = isValid;
    }
    return _result;
  }
  factory VerifyMessageResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory VerifyMessageResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  VerifyMessageResponse clone() => VerifyMessageResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  VerifyMessageResponse copyWith(void Function(VerifyMessageResponse) updates) => super.copyWith((message) => updates(message as VerifyMessageResponse)) as VerifyMessageResponse; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static VerifyMessageResponse create() => VerifyMessageResponse._();
  VerifyMessageResponse createEmptyInstance() => create();
  static $pb.PbList<VerifyMessageResponse> createRepeated() => $pb.PbList<VerifyMessageResponse>();
  @$core.pragma('dart2js:noInline')
  static VerifyMessageResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<VerifyMessageResponse>(create);
  static VerifyMessageResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.bool get isValid => $_getBF(0);
  @$pb.TagNumber(1)
  set isValid($core.bool v) { $_setBool(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasIsValid() => $_has(0);
  @$pb.TagNumber(1)
  void clearIsValid() => clearField(1);
}

class BLSPublicKeyAggregateRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'BLSPublicKeyAggregateRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..pPS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'publicKeys')
    ..hasRequiredFields = false
  ;

  BLSPublicKeyAggregateRequest._() : super();
  factory BLSPublicKeyAggregateRequest({
    $core.Iterable<$core.String>? publicKeys,
  }) {
    final _result = create();
    if (publicKeys != null) {
      _result.publicKeys.addAll(publicKeys);
    }
    return _result;
  }
  factory BLSPublicKeyAggregateRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory BLSPublicKeyAggregateRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  BLSPublicKeyAggregateRequest clone() => BLSPublicKeyAggregateRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  BLSPublicKeyAggregateRequest copyWith(void Function(BLSPublicKeyAggregateRequest) updates) => super.copyWith((message) => updates(message as BLSPublicKeyAggregateRequest)) as BLSPublicKeyAggregateRequest; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static BLSPublicKeyAggregateRequest create() => BLSPublicKeyAggregateRequest._();
  BLSPublicKeyAggregateRequest createEmptyInstance() => create();
  static $pb.PbList<BLSPublicKeyAggregateRequest> createRepeated() => $pb.PbList<BLSPublicKeyAggregateRequest>();
  @$core.pragma('dart2js:noInline')
  static BLSPublicKeyAggregateRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<BLSPublicKeyAggregateRequest>(create);
  static BLSPublicKeyAggregateRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.List<$core.String> get publicKeys => $_getList(0);
}

class BLSPublicKeyAggregateResponse extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'BLSPublicKeyAggregateResponse', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'publicKey')
    ..hasRequiredFields = false
  ;

  BLSPublicKeyAggregateResponse._() : super();
  factory BLSPublicKeyAggregateResponse({
    $core.String? publicKey,
  }) {
    final _result = create();
    if (publicKey != null) {
      _result.publicKey = publicKey;
    }
    return _result;
  }
  factory BLSPublicKeyAggregateResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory BLSPublicKeyAggregateResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  BLSPublicKeyAggregateResponse clone() => BLSPublicKeyAggregateResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  BLSPublicKeyAggregateResponse copyWith(void Function(BLSPublicKeyAggregateResponse) updates) => super.copyWith((message) => updates(message as BLSPublicKeyAggregateResponse)) as BLSPublicKeyAggregateResponse; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static BLSPublicKeyAggregateResponse create() => BLSPublicKeyAggregateResponse._();
  BLSPublicKeyAggregateResponse createEmptyInstance() => create();
  static $pb.PbList<BLSPublicKeyAggregateResponse> createRepeated() => $pb.PbList<BLSPublicKeyAggregateResponse>();
  @$core.pragma('dart2js:noInline')
  static BLSPublicKeyAggregateResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<BLSPublicKeyAggregateResponse>(create);
  static BLSPublicKeyAggregateResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get publicKey => $_getSZ(0);
  @$pb.TagNumber(1)
  set publicKey($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasPublicKey() => $_has(0);
  @$pb.TagNumber(1)
  void clearPublicKey() => clearField(1);
}

class BLSSignatureAggregateRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'BLSSignatureAggregateRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..pPS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'signatures')
    ..hasRequiredFields = false
  ;

  BLSSignatureAggregateRequest._() : super();
  factory BLSSignatureAggregateRequest({
    $core.Iterable<$core.String>? signatures,
  }) {
    final _result = create();
    if (signatures != null) {
      _result.signatures.addAll(signatures);
    }
    return _result;
  }
  factory BLSSignatureAggregateRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory BLSSignatureAggregateRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  BLSSignatureAggregateRequest clone() => BLSSignatureAggregateRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  BLSSignatureAggregateRequest copyWith(void Function(BLSSignatureAggregateRequest) updates) => super.copyWith((message) => updates(message as BLSSignatureAggregateRequest)) as BLSSignatureAggregateRequest; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static BLSSignatureAggregateRequest create() => BLSSignatureAggregateRequest._();
  BLSSignatureAggregateRequest createEmptyInstance() => create();
  static $pb.PbList<BLSSignatureAggregateRequest> createRepeated() => $pb.PbList<BLSSignatureAggregateRequest>();
  @$core.pragma('dart2js:noInline')
  static BLSSignatureAggregateRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<BLSSignatureAggregateRequest>(create);
  static BLSSignatureAggregateRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.List<$core.String> get signatures => $_getList(0);
}

class BLSSignatureAggregateResponse extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'BLSSignatureAggregateResponse', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'signature')
    ..hasRequiredFields = false
  ;

  BLSSignatureAggregateResponse._() : super();
  factory BLSSignatureAggregateResponse({
    $core.String? signature,
  }) {
    final _result = create();
    if (signature != null) {
      _result.signature = signature;
    }
    return _result;
  }
  factory BLSSignatureAggregateResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory BLSSignatureAggregateResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  BLSSignatureAggregateResponse clone() => BLSSignatureAggregateResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  BLSSignatureAggregateResponse copyWith(void Function(BLSSignatureAggregateResponse) updates) => super.copyWith((message) => updates(message as BLSSignatureAggregateResponse)) as BLSSignatureAggregateResponse; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static BLSSignatureAggregateResponse create() => BLSSignatureAggregateResponse._();
  BLSSignatureAggregateResponse createEmptyInstance() => create();
  static $pb.PbList<BLSSignatureAggregateResponse> createRepeated() => $pb.PbList<BLSSignatureAggregateResponse>();
  @$core.pragma('dart2js:noInline')
  static BLSSignatureAggregateResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<BLSSignatureAggregateResponse>(create);
  static BLSSignatureAggregateResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get signature => $_getSZ(0);
  @$pb.TagNumber(1)
  set signature($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasSignature() => $_has(0);
  @$pb.TagNumber(1)
  void clearSignature() => clearField(1);
}

class UtilsApi {
  $pb.RpcClient _client;
  UtilsApi(this._client);

  $async.Future<SignMessageWithPrivateKeyResponse> signMessageWithPrivateKey($pb.ClientContext? ctx, SignMessageWithPrivateKeyRequest request) {
    var emptyResponse = SignMessageWithPrivateKeyResponse();
    return _client.invoke<SignMessageWithPrivateKeyResponse>(ctx, 'Utils', 'SignMessageWithPrivateKey', request, emptyResponse);
  }
  $async.Future<VerifyMessageResponse> verifyMessage($pb.ClientContext? ctx, VerifyMessageRequest request) {
    var emptyResponse = VerifyMessageResponse();
    return _client.invoke<VerifyMessageResponse>(ctx, 'Utils', 'VerifyMessage', request, emptyResponse);
  }
  $async.Future<BLSPublicKeyAggregateResponse> bLSPublicKeyAggregate($pb.ClientContext? ctx, BLSPublicKeyAggregateRequest request) {
    var emptyResponse = BLSPublicKeyAggregateResponse();
    return _client.invoke<BLSPublicKeyAggregateResponse>(ctx, 'Utils', 'BLSPublicKeyAggregate', request, emptyResponse);
  }
  $async.Future<BLSSignatureAggregateResponse> bLSSignatureAggregate($pb.ClientContext? ctx, BLSSignatureAggregateRequest request) {
    var emptyResponse = BLSSignatureAggregateResponse();
    return _client.invoke<BLSSignatureAggregateResponse>(ctx, 'Utils', 'BLSSignatureAggregate', request, emptyResponse);
  }
}

