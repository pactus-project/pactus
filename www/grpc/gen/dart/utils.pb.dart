//
//  Generated code. Do not modify.
//  source: utils.proto
//
// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_final_fields
// ignore_for_file: unnecessary_import, unnecessary_this, unused_import

import 'dart:async' as $async;
import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

/// Request message for signing a message with a private key.
class SignMessageWithPrivateKeyRequest extends $pb.GeneratedMessage {
  factory SignMessageWithPrivateKeyRequest({
    $core.String? privateKey,
    $core.String? message,
  }) {
    final $result = create();
    if (privateKey != null) {
      $result.privateKey = privateKey;
    }
    if (message != null) {
      $result.message = message;
    }
    return $result;
  }
  SignMessageWithPrivateKeyRequest._() : super();
  factory SignMessageWithPrivateKeyRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory SignMessageWithPrivateKeyRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'SignMessageWithPrivateKeyRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'privateKey')
    ..aOS(2, _omitFieldNames ? '' : 'message')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  SignMessageWithPrivateKeyRequest clone() => SignMessageWithPrivateKeyRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  SignMessageWithPrivateKeyRequest copyWith(void Function(SignMessageWithPrivateKeyRequest) updates) => super.copyWith((message) => updates(message as SignMessageWithPrivateKeyRequest)) as SignMessageWithPrivateKeyRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SignMessageWithPrivateKeyRequest create() => SignMessageWithPrivateKeyRequest._();
  SignMessageWithPrivateKeyRequest createEmptyInstance() => create();
  static $pb.PbList<SignMessageWithPrivateKeyRequest> createRepeated() => $pb.PbList<SignMessageWithPrivateKeyRequest>();
  @$core.pragma('dart2js:noInline')
  static SignMessageWithPrivateKeyRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<SignMessageWithPrivateKeyRequest>(create);
  static SignMessageWithPrivateKeyRequest? _defaultInstance;

  /// The private key to sign the message.
  @$pb.TagNumber(1)
  $core.String get privateKey => $_getSZ(0);
  @$pb.TagNumber(1)
  set privateKey($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasPrivateKey() => $_has(0);
  @$pb.TagNumber(1)
  void clearPrivateKey() => $_clearField(1);

  /// The message content to be signed.
  @$pb.TagNumber(2)
  $core.String get message => $_getSZ(1);
  @$pb.TagNumber(2)
  set message($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasMessage() => $_has(1);
  @$pb.TagNumber(2)
  void clearMessage() => $_clearField(2);
}

/// Response message contains the signature generated from the message.
class SignMessageWithPrivateKeyResponse extends $pb.GeneratedMessage {
  factory SignMessageWithPrivateKeyResponse({
    $core.String? signature,
  }) {
    final $result = create();
    if (signature != null) {
      $result.signature = signature;
    }
    return $result;
  }
  SignMessageWithPrivateKeyResponse._() : super();
  factory SignMessageWithPrivateKeyResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory SignMessageWithPrivateKeyResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'SignMessageWithPrivateKeyResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'signature')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  SignMessageWithPrivateKeyResponse clone() => SignMessageWithPrivateKeyResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  SignMessageWithPrivateKeyResponse copyWith(void Function(SignMessageWithPrivateKeyResponse) updates) => super.copyWith((message) => updates(message as SignMessageWithPrivateKeyResponse)) as SignMessageWithPrivateKeyResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SignMessageWithPrivateKeyResponse create() => SignMessageWithPrivateKeyResponse._();
  SignMessageWithPrivateKeyResponse createEmptyInstance() => create();
  static $pb.PbList<SignMessageWithPrivateKeyResponse> createRepeated() => $pb.PbList<SignMessageWithPrivateKeyResponse>();
  @$core.pragma('dart2js:noInline')
  static SignMessageWithPrivateKeyResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<SignMessageWithPrivateKeyResponse>(create);
  static SignMessageWithPrivateKeyResponse? _defaultInstance;

  /// The resulting signature in hexadecimal format.
  @$pb.TagNumber(1)
  $core.String get signature => $_getSZ(0);
  @$pb.TagNumber(1)
  set signature($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasSignature() => $_has(0);
  @$pb.TagNumber(1)
  void clearSignature() => $_clearField(1);
}

/// Request message for verifying a message signature.
class VerifyMessageRequest extends $pb.GeneratedMessage {
  factory VerifyMessageRequest({
    $core.String? message,
    $core.String? signature,
    $core.String? publicKey,
  }) {
    final $result = create();
    if (message != null) {
      $result.message = message;
    }
    if (signature != null) {
      $result.signature = signature;
    }
    if (publicKey != null) {
      $result.publicKey = publicKey;
    }
    return $result;
  }
  VerifyMessageRequest._() : super();
  factory VerifyMessageRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory VerifyMessageRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'VerifyMessageRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'message')
    ..aOS(2, _omitFieldNames ? '' : 'signature')
    ..aOS(3, _omitFieldNames ? '' : 'publicKey')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  VerifyMessageRequest clone() => VerifyMessageRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  VerifyMessageRequest copyWith(void Function(VerifyMessageRequest) updates) => super.copyWith((message) => updates(message as VerifyMessageRequest)) as VerifyMessageRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static VerifyMessageRequest create() => VerifyMessageRequest._();
  VerifyMessageRequest createEmptyInstance() => create();
  static $pb.PbList<VerifyMessageRequest> createRepeated() => $pb.PbList<VerifyMessageRequest>();
  @$core.pragma('dart2js:noInline')
  static VerifyMessageRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<VerifyMessageRequest>(create);
  static VerifyMessageRequest? _defaultInstance;

  /// The original message content that was signed.
  @$pb.TagNumber(1)
  $core.String get message => $_getSZ(0);
  @$pb.TagNumber(1)
  set message($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasMessage() => $_has(0);
  @$pb.TagNumber(1)
  void clearMessage() => $_clearField(1);

  /// The signature to verify in hexadecimal format.
  @$pb.TagNumber(2)
  $core.String get signature => $_getSZ(1);
  @$pb.TagNumber(2)
  set signature($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasSignature() => $_has(1);
  @$pb.TagNumber(2)
  void clearSignature() => $_clearField(2);

  /// The public key of the signer.
  @$pb.TagNumber(3)
  $core.String get publicKey => $_getSZ(2);
  @$pb.TagNumber(3)
  set publicKey($core.String v) { $_setString(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasPublicKey() => $_has(2);
  @$pb.TagNumber(3)
  void clearPublicKey() => $_clearField(3);
}

/// Response message contains the verification result.
class VerifyMessageResponse extends $pb.GeneratedMessage {
  factory VerifyMessageResponse({
    $core.bool? isValid,
  }) {
    final $result = create();
    if (isValid != null) {
      $result.isValid = isValid;
    }
    return $result;
  }
  VerifyMessageResponse._() : super();
  factory VerifyMessageResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory VerifyMessageResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'VerifyMessageResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..aOB(1, _omitFieldNames ? '' : 'isValid')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  VerifyMessageResponse clone() => VerifyMessageResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  VerifyMessageResponse copyWith(void Function(VerifyMessageResponse) updates) => super.copyWith((message) => updates(message as VerifyMessageResponse)) as VerifyMessageResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static VerifyMessageResponse create() => VerifyMessageResponse._();
  VerifyMessageResponse createEmptyInstance() => create();
  static $pb.PbList<VerifyMessageResponse> createRepeated() => $pb.PbList<VerifyMessageResponse>();
  @$core.pragma('dart2js:noInline')
  static VerifyMessageResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<VerifyMessageResponse>(create);
  static VerifyMessageResponse? _defaultInstance;

  /// Boolean indicating whether the signature is valid for the given message and public key.
  @$pb.TagNumber(1)
  $core.bool get isValid => $_getBF(0);
  @$pb.TagNumber(1)
  set isValid($core.bool v) { $_setBool(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasIsValid() => $_has(0);
  @$pb.TagNumber(1)
  void clearIsValid() => $_clearField(1);
}

/// Request message for aggregating multiple BLS public keys.
class PublicKeyAggregationRequest extends $pb.GeneratedMessage {
  factory PublicKeyAggregationRequest({
    $core.Iterable<$core.String>? publicKeys,
  }) {
    final $result = create();
    if (publicKeys != null) {
      $result.publicKeys.addAll(publicKeys);
    }
    return $result;
  }
  PublicKeyAggregationRequest._() : super();
  factory PublicKeyAggregationRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory PublicKeyAggregationRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'PublicKeyAggregationRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..pPS(1, _omitFieldNames ? '' : 'publicKeys')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  PublicKeyAggregationRequest clone() => PublicKeyAggregationRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  PublicKeyAggregationRequest copyWith(void Function(PublicKeyAggregationRequest) updates) => super.copyWith((message) => updates(message as PublicKeyAggregationRequest)) as PublicKeyAggregationRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PublicKeyAggregationRequest create() => PublicKeyAggregationRequest._();
  PublicKeyAggregationRequest createEmptyInstance() => create();
  static $pb.PbList<PublicKeyAggregationRequest> createRepeated() => $pb.PbList<PublicKeyAggregationRequest>();
  @$core.pragma('dart2js:noInline')
  static PublicKeyAggregationRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<PublicKeyAggregationRequest>(create);
  static PublicKeyAggregationRequest? _defaultInstance;

  /// List of BLS public keys to be aggregated.
  @$pb.TagNumber(1)
  $pb.PbList<$core.String> get publicKeys => $_getList(0);
}

/// Response message contains the aggregated BLS public key result.
class PublicKeyAggregationResponse extends $pb.GeneratedMessage {
  factory PublicKeyAggregationResponse({
    $core.String? publicKey,
    $core.String? address,
  }) {
    final $result = create();
    if (publicKey != null) {
      $result.publicKey = publicKey;
    }
    if (address != null) {
      $result.address = address;
    }
    return $result;
  }
  PublicKeyAggregationResponse._() : super();
  factory PublicKeyAggregationResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory PublicKeyAggregationResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'PublicKeyAggregationResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'publicKey')
    ..aOS(2, _omitFieldNames ? '' : 'address')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  PublicKeyAggregationResponse clone() => PublicKeyAggregationResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  PublicKeyAggregationResponse copyWith(void Function(PublicKeyAggregationResponse) updates) => super.copyWith((message) => updates(message as PublicKeyAggregationResponse)) as PublicKeyAggregationResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PublicKeyAggregationResponse create() => PublicKeyAggregationResponse._();
  PublicKeyAggregationResponse createEmptyInstance() => create();
  static $pb.PbList<PublicKeyAggregationResponse> createRepeated() => $pb.PbList<PublicKeyAggregationResponse>();
  @$core.pragma('dart2js:noInline')
  static PublicKeyAggregationResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<PublicKeyAggregationResponse>(create);
  static PublicKeyAggregationResponse? _defaultInstance;

  /// The aggregated BLS public key.
  @$pb.TagNumber(1)
  $core.String get publicKey => $_getSZ(0);
  @$pb.TagNumber(1)
  set publicKey($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasPublicKey() => $_has(0);
  @$pb.TagNumber(1)
  void clearPublicKey() => $_clearField(1);

  /// The blockchain address derived from the aggregated public key.
  @$pb.TagNumber(2)
  $core.String get address => $_getSZ(1);
  @$pb.TagNumber(2)
  set address($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasAddress() => $_has(1);
  @$pb.TagNumber(2)
  void clearAddress() => $_clearField(2);
}

/// Request message for aggregating multiple BLS signatures.
class SignatureAggregationRequest extends $pb.GeneratedMessage {
  factory SignatureAggregationRequest({
    $core.Iterable<$core.String>? signatures,
  }) {
    final $result = create();
    if (signatures != null) {
      $result.signatures.addAll(signatures);
    }
    return $result;
  }
  SignatureAggregationRequest._() : super();
  factory SignatureAggregationRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory SignatureAggregationRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'SignatureAggregationRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..pPS(1, _omitFieldNames ? '' : 'signatures')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  SignatureAggregationRequest clone() => SignatureAggregationRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  SignatureAggregationRequest copyWith(void Function(SignatureAggregationRequest) updates) => super.copyWith((message) => updates(message as SignatureAggregationRequest)) as SignatureAggregationRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SignatureAggregationRequest create() => SignatureAggregationRequest._();
  SignatureAggregationRequest createEmptyInstance() => create();
  static $pb.PbList<SignatureAggregationRequest> createRepeated() => $pb.PbList<SignatureAggregationRequest>();
  @$core.pragma('dart2js:noInline')
  static SignatureAggregationRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<SignatureAggregationRequest>(create);
  static SignatureAggregationRequest? _defaultInstance;

  /// List of BLS signatures to be aggregated.
  @$pb.TagNumber(1)
  $pb.PbList<$core.String> get signatures => $_getList(0);
}

/// Response message contains the aggregated BLS signature.
class SignatureAggregationResponse extends $pb.GeneratedMessage {
  factory SignatureAggregationResponse({
    $core.String? signature,
  }) {
    final $result = create();
    if (signature != null) {
      $result.signature = signature;
    }
    return $result;
  }
  SignatureAggregationResponse._() : super();
  factory SignatureAggregationResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory SignatureAggregationResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'SignatureAggregationResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'signature')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  SignatureAggregationResponse clone() => SignatureAggregationResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  SignatureAggregationResponse copyWith(void Function(SignatureAggregationResponse) updates) => super.copyWith((message) => updates(message as SignatureAggregationResponse)) as SignatureAggregationResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static SignatureAggregationResponse create() => SignatureAggregationResponse._();
  SignatureAggregationResponse createEmptyInstance() => create();
  static $pb.PbList<SignatureAggregationResponse> createRepeated() => $pb.PbList<SignatureAggregationResponse>();
  @$core.pragma('dart2js:noInline')
  static SignatureAggregationResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<SignatureAggregationResponse>(create);
  static SignatureAggregationResponse? _defaultInstance;

  /// The aggregated BLS signature in hexadecimal format.
  @$pb.TagNumber(1)
  $core.String get signature => $_getSZ(0);
  @$pb.TagNumber(1)
  set signature($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasSignature() => $_has(0);
  @$pb.TagNumber(1)
  void clearSignature() => $_clearField(1);
}

/// Utils service defines RPC methods for utility functions such as message
/// signing, verification, and other cryptographic operations.
class UtilsApi {
  $pb.RpcClient _client;
  UtilsApi(this._client);

  /// SignMessageWithPrivateKey signs a message with the provided private key.
  $async.Future<SignMessageWithPrivateKeyResponse> signMessageWithPrivateKey($pb.ClientContext? ctx, SignMessageWithPrivateKeyRequest request) =>
    _client.invoke<SignMessageWithPrivateKeyResponse>(ctx, 'Utils', 'SignMessageWithPrivateKey', request, SignMessageWithPrivateKeyResponse())
  ;
  /// VerifyMessage verifies a signature against the public key and message.
  $async.Future<VerifyMessageResponse> verifyMessage($pb.ClientContext? ctx, VerifyMessageRequest request) =>
    _client.invoke<VerifyMessageResponse>(ctx, 'Utils', 'VerifyMessage', request, VerifyMessageResponse())
  ;
  /// PublicKeyAggregation aggregates multiple BLS public keys into a single key.
  $async.Future<PublicKeyAggregationResponse> publicKeyAggregation($pb.ClientContext? ctx, PublicKeyAggregationRequest request) =>
    _client.invoke<PublicKeyAggregationResponse>(ctx, 'Utils', 'PublicKeyAggregation', request, PublicKeyAggregationResponse())
  ;
  /// SignatureAggregation aggregates multiple BLS signatures into a single signature.
  $async.Future<SignatureAggregationResponse> signatureAggregation($pb.ClientContext? ctx, SignatureAggregationRequest request) =>
    _client.invoke<SignatureAggregationResponse>(ctx, 'Utils', 'SignatureAggregation', request, SignatureAggregationResponse())
  ;
}


const _omitFieldNames = $core.bool.fromEnvironment('protobuf.omit_field_names');
const _omitMessageNames = $core.bool.fromEnvironment('protobuf.omit_message_names');
