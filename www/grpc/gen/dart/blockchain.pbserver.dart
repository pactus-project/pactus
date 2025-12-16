// This is a generated file - do not edit.
//
// Generated from blockchain.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:async' as $async;
import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

import 'blockchain.pb.dart' as $1;
import 'blockchain.pbjson.dart';

export 'blockchain.pb.dart';

abstract class BlockchainServiceBase extends $pb.GeneratedService {
  $async.Future<$1.GetBlockResponse> getBlock(
      $pb.ServerContext ctx, $1.GetBlockRequest request);
  $async.Future<$1.GetBlockHashResponse> getBlockHash(
      $pb.ServerContext ctx, $1.GetBlockHashRequest request);
  $async.Future<$1.GetBlockHeightResponse> getBlockHeight(
      $pb.ServerContext ctx, $1.GetBlockHeightRequest request);
  $async.Future<$1.GetBlockchainInfoResponse> getBlockchainInfo(
      $pb.ServerContext ctx, $1.GetBlockchainInfoRequest request);
  $async.Future<$1.GetConsensusInfoResponse> getConsensusInfo(
      $pb.ServerContext ctx, $1.GetConsensusInfoRequest request);
  $async.Future<$1.GetAccountResponse> getAccount(
      $pb.ServerContext ctx, $1.GetAccountRequest request);
  $async.Future<$1.GetValidatorResponse> getValidator(
      $pb.ServerContext ctx, $1.GetValidatorRequest request);
  $async.Future<$1.GetValidatorResponse> getValidatorByNumber(
      $pb.ServerContext ctx, $1.GetValidatorByNumberRequest request);
  $async.Future<$1.GetValidatorAddressesResponse> getValidatorAddresses(
      $pb.ServerContext ctx, $1.GetValidatorAddressesRequest request);
  $async.Future<$1.GetPublicKeyResponse> getPublicKey(
      $pb.ServerContext ctx, $1.GetPublicKeyRequest request);
  $async.Future<$1.GetTxPoolContentResponse> getTxPoolContent(
      $pb.ServerContext ctx, $1.GetTxPoolContentRequest request);

  $pb.GeneratedMessage createRequest($core.String methodName) {
    switch (methodName) {
      case 'GetBlock':
        return $1.GetBlockRequest();
      case 'GetBlockHash':
        return $1.GetBlockHashRequest();
      case 'GetBlockHeight':
        return $1.GetBlockHeightRequest();
      case 'GetBlockchainInfo':
        return $1.GetBlockchainInfoRequest();
      case 'GetConsensusInfo':
        return $1.GetConsensusInfoRequest();
      case 'GetAccount':
        return $1.GetAccountRequest();
      case 'GetValidator':
        return $1.GetValidatorRequest();
      case 'GetValidatorByNumber':
        return $1.GetValidatorByNumberRequest();
      case 'GetValidatorAddresses':
        return $1.GetValidatorAddressesRequest();
      case 'GetPublicKey':
        return $1.GetPublicKeyRequest();
      case 'GetTxPoolContent':
        return $1.GetTxPoolContentRequest();
      default:
        throw $core.ArgumentError('Unknown method: $methodName');
    }
  }

  $async.Future<$pb.GeneratedMessage> handleCall($pb.ServerContext ctx,
      $core.String methodName, $pb.GeneratedMessage request) {
    switch (methodName) {
      case 'GetBlock':
        return getBlock(ctx, request as $1.GetBlockRequest);
      case 'GetBlockHash':
        return getBlockHash(ctx, request as $1.GetBlockHashRequest);
      case 'GetBlockHeight':
        return getBlockHeight(ctx, request as $1.GetBlockHeightRequest);
      case 'GetBlockchainInfo':
        return getBlockchainInfo(ctx, request as $1.GetBlockchainInfoRequest);
      case 'GetConsensusInfo':
        return getConsensusInfo(ctx, request as $1.GetConsensusInfoRequest);
      case 'GetAccount':
        return getAccount(ctx, request as $1.GetAccountRequest);
      case 'GetValidator':
        return getValidator(ctx, request as $1.GetValidatorRequest);
      case 'GetValidatorByNumber':
        return getValidatorByNumber(
            ctx, request as $1.GetValidatorByNumberRequest);
      case 'GetValidatorAddresses':
        return getValidatorAddresses(
            ctx, request as $1.GetValidatorAddressesRequest);
      case 'GetPublicKey':
        return getPublicKey(ctx, request as $1.GetPublicKeyRequest);
      case 'GetTxPoolContent':
        return getTxPoolContent(ctx, request as $1.GetTxPoolContentRequest);
      default:
        throw $core.ArgumentError('Unknown method: $methodName');
    }
  }

  $core.Map<$core.String, $core.dynamic> get $json =>
      BlockchainServiceBase$json;
  $core.Map<$core.String, $core.Map<$core.String, $core.dynamic>>
      get $messageJson => BlockchainServiceBase$messageJson;
}
