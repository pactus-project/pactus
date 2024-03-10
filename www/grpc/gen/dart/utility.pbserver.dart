///
//  Generated code. Do not modify.
//  source: utility.proto
//
// @dart = 2.12
// ignore_for_file: annotate_overrides,camel_case_types,constant_identifier_names,deprecated_member_use_from_same_package,directives_ordering,library_prefixes,non_constant_identifier_names,prefer_final_fields,return_of_invalid_type,unnecessary_const,unnecessary_import,unnecessary_this,unused_import,unused_shown_name

import 'dart:async' as $async;

import 'package:protobuf/protobuf.dart' as $pb;

import 'dart:core' as $core;
import 'utility.pb.dart' as $3;
import 'utility.pbjson.dart';

export 'utility.pb.dart';

abstract class UtilityServiceBase extends $pb.GeneratedService {
  $async.Future<$3.CalculateFeeResponse> calculateFee($pb.ServerContext ctx, $3.CalculateFeeRequest request);

  $pb.GeneratedMessage createRequest($core.String method) {
    switch (method) {
      case 'CalculateFee': return $3.CalculateFeeRequest();
      default: throw $core.ArgumentError('Unknown method: $method');
    }
  }

  $async.Future<$pb.GeneratedMessage> handleCall($pb.ServerContext ctx, $core.String method, $pb.GeneratedMessage request) {
    switch (method) {
      case 'CalculateFee': return this.calculateFee(ctx, request as $3.CalculateFeeRequest);
      default: throw $core.ArgumentError('Unknown method: $method');
    }
  }

  $core.Map<$core.String, $core.dynamic> get $json => UtilityServiceBase$json;
  $core.Map<$core.String, $core.Map<$core.String, $core.dynamic>> get $messageJson => UtilityServiceBase$messageJson;
}

