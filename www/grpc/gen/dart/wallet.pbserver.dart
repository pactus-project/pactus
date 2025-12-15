// This is a generated file - do not edit.
//
// Generated from wallet.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:async' as $async;
import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

import 'wallet.pb.dart' as $0;
import 'wallet.pbjson.dart';

export 'wallet.pb.dart';

abstract class WalletServiceBase extends $pb.GeneratedService {
  $async.Future<$0.CreateWalletResponse> createWallet(
      $pb.ServerContext ctx, $0.CreateWalletRequest request);
  $async.Future<$0.RestoreWalletResponse> restoreWallet(
      $pb.ServerContext ctx, $0.RestoreWalletRequest request);
  $async.Future<$0.LoadWalletResponse> loadWallet(
      $pb.ServerContext ctx, $0.LoadWalletRequest request);
  $async.Future<$0.UnloadWalletResponse> unloadWallet(
      $pb.ServerContext ctx, $0.UnloadWalletRequest request);
  $async.Future<$0.GetTotalBalanceResponse> getTotalBalance(
      $pb.ServerContext ctx, $0.GetTotalBalanceRequest request);
  $async.Future<$0.SignRawTransactionResponse> signRawTransaction(
      $pb.ServerContext ctx, $0.SignRawTransactionRequest request);
  $async.Future<$0.GetValidatorAddressResponse> getValidatorAddress(
      $pb.ServerContext ctx, $0.GetValidatorAddressRequest request);
  $async.Future<$0.GetNewAddressResponse> getNewAddress(
      $pb.ServerContext ctx, $0.GetNewAddressRequest request);
  $async.Future<$0.GetAddressHistoryResponse> getAddressHistory(
      $pb.ServerContext ctx, $0.GetAddressHistoryRequest request);
  $async.Future<$0.SignMessageResponse> signMessage(
      $pb.ServerContext ctx, $0.SignMessageRequest request);
  $async.Future<$0.GetTotalStakeResponse> getTotalStake(
      $pb.ServerContext ctx, $0.GetTotalStakeRequest request);
  $async.Future<$0.GetAddressInfoResponse> getAddressInfo(
      $pb.ServerContext ctx, $0.GetAddressInfoRequest request);
  $async.Future<$0.SetAddressLabelResponse> setAddressLabel(
      $pb.ServerContext ctx, $0.SetAddressLabelRequest request);
  $async.Future<$0.ListWalletResponse> listWallet(
      $pb.ServerContext ctx, $0.ListWalletRequest request);
  $async.Future<$0.GetWalletInfoResponse> getWalletInfo(
      $pb.ServerContext ctx, $0.GetWalletInfoRequest request);
  $async.Future<$0.ListAddressResponse> listAddress(
      $pb.ServerContext ctx, $0.ListAddressRequest request);

  $pb.GeneratedMessage createRequest($core.String methodName) {
    switch (methodName) {
      case 'CreateWallet':
        return $0.CreateWalletRequest();
      case 'RestoreWallet':
        return $0.RestoreWalletRequest();
      case 'LoadWallet':
        return $0.LoadWalletRequest();
      case 'UnloadWallet':
        return $0.UnloadWalletRequest();
      case 'GetTotalBalance':
        return $0.GetTotalBalanceRequest();
      case 'SignRawTransaction':
        return $0.SignRawTransactionRequest();
      case 'GetValidatorAddress':
        return $0.GetValidatorAddressRequest();
      case 'GetNewAddress':
        return $0.GetNewAddressRequest();
      case 'GetAddressHistory':
        return $0.GetAddressHistoryRequest();
      case 'SignMessage':
        return $0.SignMessageRequest();
      case 'GetTotalStake':
        return $0.GetTotalStakeRequest();
      case 'GetAddressInfo':
        return $0.GetAddressInfoRequest();
      case 'SetAddressLabel':
        return $0.SetAddressLabelRequest();
      case 'ListWallet':
        return $0.ListWalletRequest();
      case 'GetWalletInfo':
        return $0.GetWalletInfoRequest();
      case 'ListAddress':
        return $0.ListAddressRequest();
      default:
        throw $core.ArgumentError('Unknown method: $methodName');
    }
  }

  $async.Future<$pb.GeneratedMessage> handleCall($pb.ServerContext ctx,
      $core.String methodName, $pb.GeneratedMessage request) {
    switch (methodName) {
      case 'CreateWallet':
        return createWallet(ctx, request as $0.CreateWalletRequest);
      case 'RestoreWallet':
        return restoreWallet(ctx, request as $0.RestoreWalletRequest);
      case 'LoadWallet':
        return loadWallet(ctx, request as $0.LoadWalletRequest);
      case 'UnloadWallet':
        return unloadWallet(ctx, request as $0.UnloadWalletRequest);
      case 'GetTotalBalance':
        return getTotalBalance(ctx, request as $0.GetTotalBalanceRequest);
      case 'SignRawTransaction':
        return signRawTransaction(ctx, request as $0.SignRawTransactionRequest);
      case 'GetValidatorAddress':
        return getValidatorAddress(
            ctx, request as $0.GetValidatorAddressRequest);
      case 'GetNewAddress':
        return getNewAddress(ctx, request as $0.GetNewAddressRequest);
      case 'GetAddressHistory':
        return getAddressHistory(ctx, request as $0.GetAddressHistoryRequest);
      case 'SignMessage':
        return signMessage(ctx, request as $0.SignMessageRequest);
      case 'GetTotalStake':
        return getTotalStake(ctx, request as $0.GetTotalStakeRequest);
      case 'GetAddressInfo':
        return getAddressInfo(ctx, request as $0.GetAddressInfoRequest);
      case 'SetAddressLabel':
        return setAddressLabel(ctx, request as $0.SetAddressLabelRequest);
      case 'ListWallet':
        return listWallet(ctx, request as $0.ListWalletRequest);
      case 'GetWalletInfo':
        return getWalletInfo(ctx, request as $0.GetWalletInfoRequest);
      case 'ListAddress':
        return listAddress(ctx, request as $0.ListAddressRequest);
      default:
        throw $core.ArgumentError('Unknown method: $methodName');
    }
  }

  $core.Map<$core.String, $core.dynamic> get $json => WalletServiceBase$json;
  $core.Map<$core.String, $core.Map<$core.String, $core.dynamic>>
      get $messageJson => WalletServiceBase$messageJson;
}
