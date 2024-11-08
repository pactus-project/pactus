///
//  Generated code. Do not modify.
//  source: utils.proto
//
// @dart = 2.12
// ignore_for_file: annotate_overrides,camel_case_types,constant_identifier_names,deprecated_member_use_from_same_package,directives_ordering,library_prefixes,non_constant_identifier_names,prefer_final_fields,return_of_invalid_type,unnecessary_const,unnecessary_import,unnecessary_this,unused_import,unused_shown_name

import 'dart:async' as $async;

import 'package:protobuf/protobuf.dart' as $pb;

import 'dart:core' as $core;
import 'utils.pb.dart' as $3;
import 'utils.pbjson.dart';

export 'utils.pb.dart';

abstract class UtilsServiceBase extends $pb.GeneratedService {
  $async.Future<$3.SignMessageWithPrivateKeyResponse> signMessageWithPrivateKey($pb.ServerContext ctx, $3.SignMessageWithPrivateKeyRequest request);
  $async.Future<$3.VerifyMessageResponse> verifyMessage($pb.ServerContext ctx, $3.VerifyMessageRequest request);
  $async.Future<$3.BLSPublicKeyAggregationResponse> bLSPublicKeyAggregation($pb.ServerContext ctx, $3.BLSPublicKeyAggregationRequest request);
  $async.Future<$3.BLSSignatureAggregationResponse> bLSSignatureAggregation($pb.ServerContext ctx, $3.BLSSignatureAggregationRequest request);

  $pb.GeneratedMessage createRequest($core.String method) {
    switch (method) {
      case 'SignMessageWithPrivateKey': return $3.SignMessageWithPrivateKeyRequest();
      case 'VerifyMessage': return $3.VerifyMessageRequest();
      case 'BLSPublicKeyAggregation': return $3.BLSPublicKeyAggregationRequest();
      case 'BLSSignatureAggregation': return $3.BLSSignatureAggregationRequest();
      default: throw $core.ArgumentError('Unknown method: $method');
    }
  }

  $async.Future<$pb.GeneratedMessage> handleCall($pb.ServerContext ctx, $core.String method, $pb.GeneratedMessage request) {
    switch (method) {
      case 'SignMessageWithPrivateKey': return this.signMessageWithPrivateKey(ctx, request as $3.SignMessageWithPrivateKeyRequest);
      case 'VerifyMessage': return this.verifyMessage(ctx, request as $3.VerifyMessageRequest);
      case 'BLSPublicKeyAggregation': return this.bLSPublicKeyAggregation(ctx, request as $3.BLSPublicKeyAggregationRequest);
      case 'BLSSignatureAggregation': return this.bLSSignatureAggregation(ctx, request as $3.BLSSignatureAggregationRequest);
      default: throw $core.ArgumentError('Unknown method: $method');
    }
  }

  $core.Map<$core.String, $core.dynamic> get $json => UtilsServiceBase$json;
  $core.Map<$core.String, $core.Map<$core.String, $core.dynamic>> get $messageJson => UtilsServiceBase$messageJson;
}

