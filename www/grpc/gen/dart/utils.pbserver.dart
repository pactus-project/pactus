//
//  Generated code. Do not modify.
//  source: utils.proto
//
// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_final_fields
// ignore_for_file: unnecessary_import, unnecessary_this, unused_import

import 'dart:async' as $async;
import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

import 'utils.pb.dart' as $3;
import 'utils.pbjson.dart';

export 'utils.pb.dart';

abstract class UtilsServiceBase extends $pb.GeneratedService {
  $async.Future<$3.SignMessageWithPrivateKeyResponse> signMessageWithPrivateKey($pb.ServerContext ctx, $3.SignMessageWithPrivateKeyRequest request);
  $async.Future<$3.VerifyMessageResponse> verifyMessage($pb.ServerContext ctx, $3.VerifyMessageRequest request);
  $async.Future<$3.PublicKeyAggregationResponse> publicKeyAggregation($pb.ServerContext ctx, $3.PublicKeyAggregationRequest request);
  $async.Future<$3.SignatureAggregationResponse> signatureAggregation($pb.ServerContext ctx, $3.SignatureAggregationRequest request);

  $pb.GeneratedMessage createRequest($core.String methodName) {
    switch (methodName) {
      case 'SignMessageWithPrivateKey': return $3.SignMessageWithPrivateKeyRequest();
      case 'VerifyMessage': return $3.VerifyMessageRequest();
      case 'PublicKeyAggregation': return $3.PublicKeyAggregationRequest();
      case 'SignatureAggregation': return $3.SignatureAggregationRequest();
      default: throw $core.ArgumentError('Unknown method: $methodName');
    }
  }

  $async.Future<$pb.GeneratedMessage> handleCall($pb.ServerContext ctx, $core.String methodName, $pb.GeneratedMessage request) {
    switch (methodName) {
      case 'SignMessageWithPrivateKey': return this.signMessageWithPrivateKey(ctx, request as $3.SignMessageWithPrivateKeyRequest);
      case 'VerifyMessage': return this.verifyMessage(ctx, request as $3.VerifyMessageRequest);
      case 'PublicKeyAggregation': return this.publicKeyAggregation(ctx, request as $3.PublicKeyAggregationRequest);
      case 'SignatureAggregation': return this.signatureAggregation(ctx, request as $3.SignatureAggregationRequest);
      default: throw $core.ArgumentError('Unknown method: $methodName');
    }
  }

  $core.Map<$core.String, $core.dynamic> get $json => UtilsServiceBase$json;
  $core.Map<$core.String, $core.Map<$core.String, $core.dynamic>> get $messageJson => UtilsServiceBase$messageJson;
}

