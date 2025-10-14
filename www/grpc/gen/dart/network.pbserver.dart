//
//  Generated code. Do not modify.
//  source: network.proto
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

import 'network.pb.dart' as $2;
import 'network.pbjson.dart';

export 'network.pb.dart';

abstract class NetworkServiceBase extends $pb.GeneratedService {
  $async.Future<$2.GetNetworkInfoResponse> getNetworkInfo($pb.ServerContext ctx, $2.GetNetworkInfoRequest request);
  $async.Future<$2.GetNodeInfoResponse> getNodeInfo($pb.ServerContext ctx, $2.GetNodeInfoRequest request);
  $async.Future<$2.PingResponse> ping($pb.ServerContext ctx, $2.PingRequest request);

  $pb.GeneratedMessage createRequest($core.String methodName) {
    switch (methodName) {
      case 'GetNetworkInfo': return $2.GetNetworkInfoRequest();
      case 'GetNodeInfo': return $2.GetNodeInfoRequest();
      case 'Ping': return $2.PingRequest();
      default: throw $core.ArgumentError('Unknown method: $methodName');
    }
  }

  $async.Future<$pb.GeneratedMessage> handleCall($pb.ServerContext ctx, $core.String methodName, $pb.GeneratedMessage request) {
    switch (methodName) {
      case 'GetNetworkInfo': return this.getNetworkInfo(ctx, request as $2.GetNetworkInfoRequest);
      case 'GetNodeInfo': return this.getNodeInfo(ctx, request as $2.GetNodeInfoRequest);
      case 'Ping': return this.ping(ctx, request as $2.PingRequest);
      default: throw $core.ArgumentError('Unknown method: $methodName');
    }
  }

  $core.Map<$core.String, $core.dynamic> get $json => NetworkServiceBase$json;
  $core.Map<$core.String, $core.Map<$core.String, $core.dynamic>> get $messageJson => NetworkServiceBase$messageJson;
}

