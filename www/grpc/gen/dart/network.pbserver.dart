///
//  Generated code. Do not modify.
//  source: network.proto
//
// @dart = 2.12
// ignore_for_file: annotate_overrides,camel_case_types,constant_identifier_names,deprecated_member_use_from_same_package,directives_ordering,library_prefixes,non_constant_identifier_names,prefer_final_fields,return_of_invalid_type,unnecessary_const,unnecessary_import,unnecessary_this,unused_import,unused_shown_name

import 'dart:async' as $async;

import 'package:protobuf/protobuf.dart' as $pb;

import 'dart:core' as $core;
import 'network.pb.dart' as $2;
import 'network.pbjson.dart';

export 'network.pb.dart';

abstract class NetworkServiceBase extends $pb.GeneratedService {
  $async.Future<$2.GetNetworkInfoResponse> getNetworkInfo($pb.ServerContext ctx, $2.GetNetworkInfoRequest request);
  $async.Future<$2.GetNodeInfoResponse> getNodeInfo($pb.ServerContext ctx, $2.GetNodeInfoRequest request);
  $async.Future<$2.GetPeersInfoResponse> getPeersInfo($pb.ServerContext ctx, $2.GetPeersInfoRequest request);

  $pb.GeneratedMessage createRequest($core.String method) {
    switch (method) {
      case 'GetNetworkInfo': return $2.GetNetworkInfoRequest();
      case 'GetNodeInfo': return $2.GetNodeInfoRequest();
      case 'GetPeersInfo': return $2.GetPeersInfoRequest();
      default: throw $core.ArgumentError('Unknown method: $method');
    }
  }

  $async.Future<$pb.GeneratedMessage> handleCall($pb.ServerContext ctx, $core.String method, $pb.GeneratedMessage request) {
    switch (method) {
      case 'GetNetworkInfo': return this.getNetworkInfo(ctx, request as $2.GetNetworkInfoRequest);
      case 'GetNodeInfo': return this.getNodeInfo(ctx, request as $2.GetNodeInfoRequest);
      case 'GetPeersInfo': return this.getPeersInfo(ctx, request as $2.GetPeersInfoRequest);
      default: throw $core.ArgumentError('Unknown method: $method');
    }
  }

  $core.Map<$core.String, $core.dynamic> get $json => NetworkServiceBase$json;
  $core.Map<$core.String, $core.Map<$core.String, $core.dynamic>> get $messageJson => NetworkServiceBase$messageJson;
}

