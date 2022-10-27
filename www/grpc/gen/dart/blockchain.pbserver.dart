///
//  Generated code. Do not modify.
//  source: blockchain.proto
//
// @dart = 2.12
// ignore_for_file: annotate_overrides,camel_case_types,constant_identifier_names,deprecated_member_use_from_same_package,directives_ordering,library_prefixes,non_constant_identifier_names,prefer_final_fields,return_of_invalid_type,unnecessary_const,unnecessary_import,unnecessary_this,unused_import,unused_shown_name

import 'dart:async' as $async;

import 'package:protobuf/protobuf.dart' as $pb;

import 'dart:core' as $core;
import 'blockchain.pb.dart' as $1;
import 'blockchain.pbjson.dart';

export 'blockchain.pb.dart';

abstract class BlockchainServiceBase extends $pb.GeneratedService {
  $async.Future<$1.BlockResponse> getBlock($pb.ServerContext ctx, $1.BlockRequest request);
  $async.Future<$1.BlockHashResponse> getBlockHash($pb.ServerContext ctx, $1.BlockHashRequest request);
  $async.Future<$1.BlockHeightResponse> getBlockHeight($pb.ServerContext ctx, $1.BlockHeightRequest request);
  $async.Future<$1.AccountResponse> getAccount($pb.ServerContext ctx, $1.AccountRequest request);
  $async.Future<$1.ValidatorsResponse> getValidators($pb.ServerContext ctx, $1.ValidatorsRequest request);
  $async.Future<$1.ValidatorResponse> getValidator($pb.ServerContext ctx, $1.ValidatorRequest request);
  $async.Future<$1.ValidatorResponse> getValidatorByNumber($pb.ServerContext ctx, $1.ValidatorByNumberRequest request);
  $async.Future<$1.BlockchainInfoResponse> getBlockchainInfo($pb.ServerContext ctx, $1.BlockchainInfoRequest request);

  $pb.GeneratedMessage createRequest($core.String method) {
    switch (method) {
      case 'GetBlock': return $1.BlockRequest();
      case 'GetBlockHash': return $1.BlockHashRequest();
      case 'GetBlockHeight': return $1.BlockHeightRequest();
      case 'GetAccount': return $1.AccountRequest();
      case 'GetValidators': return $1.ValidatorsRequest();
      case 'GetValidator': return $1.ValidatorRequest();
      case 'GetValidatorByNumber': return $1.ValidatorByNumberRequest();
      case 'GetBlockchainInfo': return $1.BlockchainInfoRequest();
      default: throw $core.ArgumentError('Unknown method: $method');
    }
  }

  $async.Future<$pb.GeneratedMessage> handleCall($pb.ServerContext ctx, $core.String method, $pb.GeneratedMessage request) {
    switch (method) {
      case 'GetBlock': return this.getBlock(ctx, request as $1.BlockRequest);
      case 'GetBlockHash': return this.getBlockHash(ctx, request as $1.BlockHashRequest);
      case 'GetBlockHeight': return this.getBlockHeight(ctx, request as $1.BlockHeightRequest);
      case 'GetAccount': return this.getAccount(ctx, request as $1.AccountRequest);
      case 'GetValidators': return this.getValidators(ctx, request as $1.ValidatorsRequest);
      case 'GetValidator': return this.getValidator(ctx, request as $1.ValidatorRequest);
      case 'GetValidatorByNumber': return this.getValidatorByNumber(ctx, request as $1.ValidatorByNumberRequest);
      case 'GetBlockchainInfo': return this.getBlockchainInfo(ctx, request as $1.BlockchainInfoRequest);
      default: throw $core.ArgumentError('Unknown method: $method');
    }
  }

  $core.Map<$core.String, $core.dynamic> get $json => BlockchainServiceBase$json;
  $core.Map<$core.String, $core.Map<$core.String, $core.dynamic>> get $messageJson => BlockchainServiceBase$messageJson;
}

