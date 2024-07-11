///
//  Generated code. Do not modify.
//  source: wallet.proto
//
// @dart = 2.12
// ignore_for_file: annotate_overrides,camel_case_types,constant_identifier_names,deprecated_member_use_from_same_package,directives_ordering,library_prefixes,non_constant_identifier_names,prefer_final_fields,return_of_invalid_type,unnecessary_const,unnecessary_import,unnecessary_this,unused_import,unused_shown_name

import 'dart:async' as $async;

import 'package:protobuf/protobuf.dart' as $pb;

import 'dart:core' as $core;
import 'wallet.pb.dart' as $4;
import 'wallet.pbjson.dart';

export 'wallet.pb.dart';

abstract class WalletServiceBase extends $pb.GeneratedService {
  $async.Future<$4.CreateWalletResponse> createWallet($pb.ServerContext ctx, $4.CreateWalletRequest request);
  $async.Future<$4.RestoreWalletResponse> restoreWallet($pb.ServerContext ctx, $4.RestoreWalletRequest request);
  $async.Future<$4.LoadWalletResponse> loadWallet($pb.ServerContext ctx, $4.LoadWalletRequest request);
  $async.Future<$4.UnloadWalletResponse> unloadWallet($pb.ServerContext ctx, $4.UnloadWalletRequest request);
  $async.Future<$4.GetTotalBalanceResponse> getTotalBalance($pb.ServerContext ctx, $4.GetTotalBalanceRequest request);
  $async.Future<$4.SignRawTransactionResponse> signRawTransaction($pb.ServerContext ctx, $4.SignRawTransactionRequest request);
  $async.Future<$4.GetValidatorAddressResponse> getValidatorAddress($pb.ServerContext ctx, $4.GetValidatorAddressRequest request);
  $async.Future<$4.GetNewAddressResponse> getNewAddress($pb.ServerContext ctx, $4.GetNewAddressRequest request);
  $async.Future<$4.GetAddressHistoryResponse> getAddressHistory($pb.ServerContext ctx, $4.GetAddressHistoryRequest request);
  $async.Future<$4.SignMessageResponse> signMessage($pb.ServerContext ctx, $4.SignMessageRequest request);

  $pb.GeneratedMessage createRequest($core.String method) {
    switch (method) {
      case 'CreateWallet': return $4.CreateWalletRequest();
      case 'RestoreWallet': return $4.RestoreWalletRequest();
      case 'LoadWallet': return $4.LoadWalletRequest();
      case 'UnloadWallet': return $4.UnloadWalletRequest();
      case 'GetTotalBalance': return $4.GetTotalBalanceRequest();
      case 'SignRawTransaction': return $4.SignRawTransactionRequest();
      case 'GetValidatorAddress': return $4.GetValidatorAddressRequest();
      case 'GetNewAddress': return $4.GetNewAddressRequest();
      case 'GetAddressHistory': return $4.GetAddressHistoryRequest();
      case 'SignMessage': return $4.SignMessageRequest();
      default: throw $core.ArgumentError('Unknown method: $method');
    }
  }

  $async.Future<$pb.GeneratedMessage> handleCall($pb.ServerContext ctx, $core.String method, $pb.GeneratedMessage request) {
    switch (method) {
      case 'CreateWallet': return this.createWallet(ctx, request as $4.CreateWalletRequest);
      case 'RestoreWallet': return this.restoreWallet(ctx, request as $4.RestoreWalletRequest);
      case 'LoadWallet': return this.loadWallet(ctx, request as $4.LoadWalletRequest);
      case 'UnloadWallet': return this.unloadWallet(ctx, request as $4.UnloadWalletRequest);
      case 'GetTotalBalance': return this.getTotalBalance(ctx, request as $4.GetTotalBalanceRequest);
      case 'SignRawTransaction': return this.signRawTransaction(ctx, request as $4.SignRawTransactionRequest);
      case 'GetValidatorAddress': return this.getValidatorAddress(ctx, request as $4.GetValidatorAddressRequest);
      case 'GetNewAddress': return this.getNewAddress(ctx, request as $4.GetNewAddressRequest);
      case 'GetAddressHistory': return this.getAddressHistory(ctx, request as $4.GetAddressHistoryRequest);
      case 'SignMessage': return this.signMessage(ctx, request as $4.SignMessageRequest);
      default: throw $core.ArgumentError('Unknown method: $method');
    }
  }

  $core.Map<$core.String, $core.dynamic> get $json => WalletServiceBase$json;
  $core.Map<$core.String, $core.Map<$core.String, $core.dynamic>> get $messageJson => WalletServiceBase$messageJson;
}

