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
  $async.Future<$1.GetBlockResponse> getBlock($pb.ServerContext ctx, $1.GetBlockRequest request);
  $async.Future<$1.GetBlockHashResponse> getBlockHash($pb.ServerContext ctx, $1.GetBlockHashRequest request);
  $async.Future<$1.GetBlockHeightResponse> getBlockHeight($pb.ServerContext ctx, $1.GetBlockHeightRequest request);
  $async.Future<$1.GetBlockchainInfoResponse> getBlockchainInfo($pb.ServerContext ctx, $1.GetBlockchainInfoRequest request);
  $async.Future<$1.GetConsensusInfoResponse> getConsensusInfo($pb.ServerContext ctx, $1.GetConsensusInfoRequest request);
  $async.Future<$1.GetAccountResponse> getAccount($pb.ServerContext ctx, $1.GetAccountRequest request);
  $async.Future<$1.GetAccountResponse> getAccountByNumber($pb.ServerContext ctx, $1.GetAccountByNumberRequest request);
  $async.Future<$1.GetValidatorResponse> getValidator($pb.ServerContext ctx, $1.GetValidatorRequest request);
  $async.Future<$1.GetValidatorResponse> getValidatorByNumber($pb.ServerContext ctx, $1.GetValidatorByNumberRequest request);
  $async.Future<$1.GetValidatorAddressesResponse> getValidatorAddresses($pb.ServerContext ctx, $1.GetValidatorAddressesRequest request);

  $pb.GeneratedMessage createRequest($core.String method) {
    switch (method) {
      case 'GetBlock': return $1.GetBlockRequest();
      case 'GetBlockHash': return $1.GetBlockHashRequest();
      case 'GetBlockHeight': return $1.GetBlockHeightRequest();
      case 'GetBlockchainInfo': return $1.GetBlockchainInfoRequest();
      case 'GetConsensusInfo': return $1.GetConsensusInfoRequest();
      case 'GetAccount': return $1.GetAccountRequest();
      case 'GetAccountByNumber': return $1.GetAccountByNumberRequest();
      case 'GetValidator': return $1.GetValidatorRequest();
      case 'GetValidatorByNumber': return $1.GetValidatorByNumberRequest();
      case 'GetValidatorAddresses': return $1.GetValidatorAddressesRequest();
      default: throw $core.ArgumentError('Unknown method: $method');
    }
  }

  $async.Future<$pb.GeneratedMessage> handleCall($pb.ServerContext ctx, $core.String method, $pb.GeneratedMessage request) {
    switch (method) {
      case 'GetBlock': return this.getBlock(ctx, request as $1.GetBlockRequest);
      case 'GetBlockHash': return this.getBlockHash(ctx, request as $1.GetBlockHashRequest);
      case 'GetBlockHeight': return this.getBlockHeight(ctx, request as $1.GetBlockHeightRequest);
      case 'GetBlockchainInfo': return this.getBlockchainInfo(ctx, request as $1.GetBlockchainInfoRequest);
      case 'GetConsensusInfo': return this.getConsensusInfo(ctx, request as $1.GetConsensusInfoRequest);
      case 'GetAccount': return this.getAccount(ctx, request as $1.GetAccountRequest);
      case 'GetAccountByNumber': return this.getAccountByNumber(ctx, request as $1.GetAccountByNumberRequest);
      case 'GetValidator': return this.getValidator(ctx, request as $1.GetValidatorRequest);
      case 'GetValidatorByNumber': return this.getValidatorByNumber(ctx, request as $1.GetValidatorByNumberRequest);
      case 'GetValidatorAddresses': return this.getValidatorAddresses(ctx, request as $1.GetValidatorAddressesRequest);
      default: throw $core.ArgumentError('Unknown method: $method');
    }
  }

  $core.Map<$core.String, $core.dynamic> get $json => BlockchainServiceBase$json;
  $core.Map<$core.String, $core.Map<$core.String, $core.dynamic>> get $messageJson => BlockchainServiceBase$messageJson;
}

