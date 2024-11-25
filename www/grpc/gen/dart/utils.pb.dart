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

class BLSPublicKeyAggregationRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'BLSPublicKeyAggregationRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..pPS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'publicKeys')
    ..hasRequiredFields = false
  ;

  BLSPublicKeyAggregationRequest._() : super();
  factory BLSPublicKeyAggregationRequest({
    $core.Iterable<$core.String>? publicKeys,
  }) {
    final _result = create();
    if (publicKeys != null) {
      _result.publicKeys.addAll(publicKeys);
    }
    return _result;
  }
  factory BLSPublicKeyAggregationRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory BLSPublicKeyAggregationRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  BLSPublicKeyAggregationRequest clone() => BLSPublicKeyAggregationRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  BLSPublicKeyAggregationRequest copyWith(void Function(BLSPublicKeyAggregationRequest) updates) => super.copyWith((message) => updates(message as BLSPublicKeyAggregationRequest)) as BLSPublicKeyAggregationRequest; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static BLSPublicKeyAggregationRequest create() => BLSPublicKeyAggregationRequest._();
  BLSPublicKeyAggregationRequest createEmptyInstance() => create();
  static $pb.PbList<BLSPublicKeyAggregationRequest> createRepeated() => $pb.PbList<BLSPublicKeyAggregationRequest>();
  @$core.pragma('dart2js:noInline')
  static BLSPublicKeyAggregationRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<BLSPublicKeyAggregationRequest>(create);
  static BLSPublicKeyAggregationRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.List<$core.String> get publicKeys => $_getList(0);
}

class BLSPublicKeyAggregationResponse extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'BLSPublicKeyAggregationResponse', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'publicKey')
    ..aOS(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'address')
    ..hasRequiredFields = false
  ;

  BLSPublicKeyAggregationResponse._() : super();
  factory BLSPublicKeyAggregationResponse({
    $core.String? publicKey,
    $core.String? address,
  }) {
    final _result = create();
    if (publicKey != null) {
      _result.publicKey = publicKey;
    }
    if (address != null) {
      _result.address = address;
    }
    return _result;
  }
  factory BLSPublicKeyAggregationResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory BLSPublicKeyAggregationResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  BLSPublicKeyAggregationResponse clone() => BLSPublicKeyAggregationResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  BLSPublicKeyAggregationResponse copyWith(void Function(BLSPublicKeyAggregationResponse) updates) => super.copyWith((message) => updates(message as BLSPublicKeyAggregationResponse)) as BLSPublicKeyAggregationResponse; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static BLSPublicKeyAggregationResponse create() => BLSPublicKeyAggregationResponse._();
  BLSPublicKeyAggregationResponse createEmptyInstance() => create();
  static $pb.PbList<BLSPublicKeyAggregationResponse> createRepeated() => $pb.PbList<BLSPublicKeyAggregationResponse>();
  @$core.pragma('dart2js:noInline')
  static BLSPublicKeyAggregationResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<BLSPublicKeyAggregationResponse>(create);
  static BLSPublicKeyAggregationResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get publicKey => $_getSZ(0);
  @$pb.TagNumber(1)
  set publicKey($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasPublicKey() => $_has(0);
  @$pb.TagNumber(1)
  void clearPublicKey() => clearField(1);

  @$pb.TagNumber(2)
  $core.String get address => $_getSZ(1);
  @$pb.TagNumber(2)
  set address($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasAddress() => $_has(1);
  @$pb.TagNumber(2)
  void clearAddress() => clearField(2);
}

class BLSSignatureAggregationRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'BLSSignatureAggregationRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..pPS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'signatures')
    ..hasRequiredFields = false
  ;

  BLSSignatureAggregationRequest._() : super();
  factory BLSSignatureAggregationRequest({
    $core.Iterable<$core.String>? signatures,
  }) {
    final _result = create();
    if (signatures != null) {
      _result.signatures.addAll(signatures);
    }
    return _result;
  }
  factory BLSSignatureAggregationRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory BLSSignatureAggregationRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  BLSSignatureAggregationRequest clone() => BLSSignatureAggregationRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  BLSSignatureAggregationRequest copyWith(void Function(BLSSignatureAggregationRequest) updates) => super.copyWith((message) => updates(message as BLSSignatureAggregationRequest)) as BLSSignatureAggregationRequest; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static BLSSignatureAggregationRequest create() => BLSSignatureAggregationRequest._();
  BLSSignatureAggregationRequest createEmptyInstance() => create();
  static $pb.PbList<BLSSignatureAggregationRequest> createRepeated() => $pb.PbList<BLSSignatureAggregationRequest>();
  @$core.pragma('dart2js:noInline')
  static BLSSignatureAggregationRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<BLSSignatureAggregationRequest>(create);
  static BLSSignatureAggregationRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $core.List<$core.String> get signatures => $_getList(0);
}

class BLSSignatureAggregationResponse extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'BLSSignatureAggregationResponse', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'signature')
    ..hasRequiredFields = false
  ;

  BLSSignatureAggregationResponse._() : super();
  factory BLSSignatureAggregationResponse({
    $core.String? signature,
  }) {
    final _result = create();
    if (signature != null) {
      _result.signature = signature;
    }
    return _result;
  }
  factory BLSSignatureAggregationResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory BLSSignatureAggregationResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  BLSSignatureAggregationResponse clone() => BLSSignatureAggregationResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  BLSSignatureAggregationResponse copyWith(void Function(BLSSignatureAggregationResponse) updates) => super.copyWith((message) => updates(message as BLSSignatureAggregationResponse)) as BLSSignatureAggregationResponse; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static BLSSignatureAggregationResponse create() => BLSSignatureAggregationResponse._();
  BLSSignatureAggregationResponse createEmptyInstance() => create();
  static $pb.PbList<BLSSignatureAggregationResponse> createRepeated() => $pb.PbList<BLSSignatureAggregationResponse>();
  @$core.pragma('dart2js:noInline')
  static BLSSignatureAggregationResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<BLSSignatureAggregationResponse>(create);
  static BLSSignatureAggregationResponse? _defaultInstance;

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
  $async.Future<BLSPublicKeyAggregationResponse> bLSPublicKeyAggregation($pb.ClientContext? ctx, BLSPublicKeyAggregationRequest request) {
    var emptyResponse = BLSPublicKeyAggregationResponse();
    return _client.invoke<BLSPublicKeyAggregationResponse>(ctx, 'Utils', 'BLSPublicKeyAggregation', request, emptyResponse);
  }
  $async.Future<BLSSignatureAggregationResponse> bLSSignatureAggregation($pb.ClientContext? ctx, BLSSignatureAggregationRequest request) {
    var emptyResponse = BLSSignatureAggregationResponse();
    return _client.invoke<BLSSignatureAggregationResponse>(ctx, 'Utils', 'BLSSignatureAggregation', request, emptyResponse);
  }
}

