// This is a generated file - do not edit.
//
// Generated from utils.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:async' as $async;
import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

import 'utils.pb.dart' as $0;
import 'utils.pbjson.dart';

export 'utils.pb.dart';

abstract class UtilsServiceBase extends $pb.GeneratedService {
  $async.Future<$0.SignMessageWithPrivateKeyResponse> signMessageWithPrivateKey(
      $pb.ServerContext ctx, $0.SignMessageWithPrivateKeyRequest request);
  $async.Future<$0.VerifyMessageResponse> verifyMessage(
      $pb.ServerContext ctx, $0.VerifyMessageRequest request);
  $async.Future<$0.PublicKeyAggregationResponse> publicKeyAggregation(
      $pb.ServerContext ctx, $0.PublicKeyAggregationRequest request);
  $async.Future<$0.SignatureAggregationResponse> signatureAggregation(
      $pb.ServerContext ctx, $0.SignatureAggregationRequest request);

  $pb.GeneratedMessage createRequest($core.String methodName) {
    switch (methodName) {
      case 'SignMessageWithPrivateKey':
        return $0.SignMessageWithPrivateKeyRequest();
      case 'VerifyMessage':
        return $0.VerifyMessageRequest();
      case 'PublicKeyAggregation':
        return $0.PublicKeyAggregationRequest();
      case 'SignatureAggregation':
        return $0.SignatureAggregationRequest();
      default:
        throw $core.ArgumentError('Unknown method: $methodName');
    }
  }

  $async.Future<$pb.GeneratedMessage> handleCall($pb.ServerContext ctx,
      $core.String methodName, $pb.GeneratedMessage request) {
    switch (methodName) {
      case 'SignMessageWithPrivateKey':
        return signMessageWithPrivateKey(
            ctx, request as $0.SignMessageWithPrivateKeyRequest);
      case 'VerifyMessage':
        return verifyMessage(ctx, request as $0.VerifyMessageRequest);
      case 'PublicKeyAggregation':
        return publicKeyAggregation(
            ctx, request as $0.PublicKeyAggregationRequest);
      case 'SignatureAggregation':
        return signatureAggregation(
            ctx, request as $0.SignatureAggregationRequest);
      default:
        throw $core.ArgumentError('Unknown method: $methodName');
    }
  }

  $core.Map<$core.String, $core.dynamic> get $json => UtilsServiceBase$json;
  $core.Map<$core.String, $core.Map<$core.String, $core.dynamic>>
      get $messageJson => UtilsServiceBase$messageJson;
}
