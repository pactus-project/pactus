// This is a generated file - do not edit.
//
// Generated from network.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:async' as $async;
import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

import 'network.pb.dart' as $0;
import 'network.pbjson.dart';

export 'network.pb.dart';

abstract class NetworkServiceBase extends $pb.GeneratedService {
  $async.Future<$0.GetNetworkInfoResponse> getNetworkInfo(
      $pb.ServerContext ctx, $0.GetNetworkInfoRequest request);
  $async.Future<$0.ListPeersResponse> listPeers(
      $pb.ServerContext ctx, $0.ListPeersRequest request);
  $async.Future<$0.GetNodeInfoResponse> getNodeInfo(
      $pb.ServerContext ctx, $0.GetNodeInfoRequest request);
  $async.Future<$0.PingResponse> ping(
      $pb.ServerContext ctx, $0.PingRequest request);

  $pb.GeneratedMessage createRequest($core.String methodName) {
    switch (methodName) {
      case 'GetNetworkInfo':
        return $0.GetNetworkInfoRequest();
      case 'ListPeers':
        return $0.ListPeersRequest();
      case 'GetNodeInfo':
        return $0.GetNodeInfoRequest();
      case 'Ping':
        return $0.PingRequest();
      default:
        throw $core.ArgumentError('Unknown method: $methodName');
    }
  }

  $async.Future<$pb.GeneratedMessage> handleCall($pb.ServerContext ctx,
      $core.String methodName, $pb.GeneratedMessage request) {
    switch (methodName) {
      case 'GetNetworkInfo':
        return getNetworkInfo(ctx, request as $0.GetNetworkInfoRequest);
      case 'ListPeers':
        return listPeers(ctx, request as $0.ListPeersRequest);
      case 'GetNodeInfo':
        return getNodeInfo(ctx, request as $0.GetNodeInfoRequest);
      case 'Ping':
        return ping(ctx, request as $0.PingRequest);
      default:
        throw $core.ArgumentError('Unknown method: $methodName');
    }
  }

  $core.Map<$core.String, $core.dynamic> get $json => NetworkServiceBase$json;
  $core.Map<$core.String, $core.Map<$core.String, $core.dynamic>>
      get $messageJson => NetworkServiceBase$messageJson;
}
