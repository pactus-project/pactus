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

import 'wallet.pb.dart' as $1;
import 'wallet.pbjson.dart';

export 'wallet.pb.dart';

abstract class WalletServiceBase extends $pb.GeneratedService {
  $async.Future<$1.CreateWalletResponse> createWallet(
      $pb.ServerContext ctx, $1.CreateWalletRequest request);
  $async.Future<$1.RestoreWalletResponse> restoreWallet(
      $pb.ServerContext ctx, $1.RestoreWalletRequest request);
  $async.Future<$1.LoadWalletResponse> loadWallet(
      $pb.ServerContext ctx, $1.LoadWalletRequest request);
  $async.Future<$1.UnloadWalletResponse> unloadWallet(
      $pb.ServerContext ctx, $1.UnloadWalletRequest request);
  $async.Future<$1.ListWalletsResponse> listWallets(
      $pb.ServerContext ctx, $1.ListWalletsRequest request);
  $async.Future<$1.GetWalletInfoResponse> getWalletInfo(
      $pb.ServerContext ctx, $1.GetWalletInfoRequest request);
  $async.Future<$1.UpdatePasswordResponse> updatePassword(
      $pb.ServerContext ctx, $1.UpdatePasswordRequest request);
  $async.Future<$1.GetTotalBalanceResponse> getTotalBalance(
      $pb.ServerContext ctx, $1.GetTotalBalanceRequest request);
  $async.Future<$1.GetTotalStakeResponse> getTotalStake(
      $pb.ServerContext ctx, $1.GetTotalStakeRequest request);
  $async.Future<$1.GetValidatorAddressResponse> getValidatorAddress(
      $pb.ServerContext ctx, $1.GetValidatorAddressRequest request);
  $async.Future<$1.GetAddressInfoResponse> getAddressInfo(
      $pb.ServerContext ctx, $1.GetAddressInfoRequest request);
  $async.Future<$1.SetAddressLabelResponse> setAddressLabel(
      $pb.ServerContext ctx, $1.SetAddressLabelRequest request);
  $async.Future<$1.GetNewAddressResponse> getNewAddress(
      $pb.ServerContext ctx, $1.GetNewAddressRequest request);
  $async.Future<$1.ListAddressesResponse> listAddresses(
      $pb.ServerContext ctx, $1.ListAddressesRequest request);
  $async.Future<$1.SignMessageResponse> signMessage(
      $pb.ServerContext ctx, $1.SignMessageRequest request);
  $async.Future<$1.SignRawTransactionResponse> signRawTransaction(
      $pb.ServerContext ctx, $1.SignRawTransactionRequest request);
  $async.Future<$1.ListTransactionsResponse> listTransactions(
      $pb.ServerContext ctx, $1.ListTransactionsRequest request);

  $pb.GeneratedMessage createRequest($core.String methodName) {
    switch (methodName) {
      case 'CreateWallet':
        return $1.CreateWalletRequest();
      case 'RestoreWallet':
        return $1.RestoreWalletRequest();
      case 'LoadWallet':
        return $1.LoadWalletRequest();
      case 'UnloadWallet':
        return $1.UnloadWalletRequest();
      case 'ListWallets':
        return $1.ListWalletsRequest();
      case 'GetWalletInfo':
        return $1.GetWalletInfoRequest();
      case 'UpdatePassword':
        return $1.UpdatePasswordRequest();
      case 'GetTotalBalance':
        return $1.GetTotalBalanceRequest();
      case 'GetTotalStake':
        return $1.GetTotalStakeRequest();
      case 'GetValidatorAddress':
        return $1.GetValidatorAddressRequest();
      case 'GetAddressInfo':
        return $1.GetAddressInfoRequest();
      case 'SetAddressLabel':
        return $1.SetAddressLabelRequest();
      case 'GetNewAddress':
        return $1.GetNewAddressRequest();
      case 'ListAddresses':
        return $1.ListAddressesRequest();
      case 'SignMessage':
        return $1.SignMessageRequest();
      case 'SignRawTransaction':
        return $1.SignRawTransactionRequest();
      case 'ListTransactions':
        return $1.ListTransactionsRequest();
      default:
        throw $core.ArgumentError('Unknown method: $methodName');
    }
  }

  $async.Future<$pb.GeneratedMessage> handleCall($pb.ServerContext ctx,
      $core.String methodName, $pb.GeneratedMessage request) {
    switch (methodName) {
      case 'CreateWallet':
        return createWallet(ctx, request as $1.CreateWalletRequest);
      case 'RestoreWallet':
        return restoreWallet(ctx, request as $1.RestoreWalletRequest);
      case 'LoadWallet':
        return loadWallet(ctx, request as $1.LoadWalletRequest);
      case 'UnloadWallet':
        return unloadWallet(ctx, request as $1.UnloadWalletRequest);
      case 'ListWallets':
        return listWallets(ctx, request as $1.ListWalletsRequest);
      case 'GetWalletInfo':
        return getWalletInfo(ctx, request as $1.GetWalletInfoRequest);
      case 'UpdatePassword':
        return updatePassword(ctx, request as $1.UpdatePasswordRequest);
      case 'GetTotalBalance':
        return getTotalBalance(ctx, request as $1.GetTotalBalanceRequest);
      case 'GetTotalStake':
        return getTotalStake(ctx, request as $1.GetTotalStakeRequest);
      case 'GetValidatorAddress':
        return getValidatorAddress(
            ctx, request as $1.GetValidatorAddressRequest);
      case 'GetAddressInfo':
        return getAddressInfo(ctx, request as $1.GetAddressInfoRequest);
      case 'SetAddressLabel':
        return setAddressLabel(ctx, request as $1.SetAddressLabelRequest);
      case 'GetNewAddress':
        return getNewAddress(ctx, request as $1.GetNewAddressRequest);
      case 'ListAddresses':
        return listAddresses(ctx, request as $1.ListAddressesRequest);
      case 'SignMessage':
        return signMessage(ctx, request as $1.SignMessageRequest);
      case 'SignRawTransaction':
        return signRawTransaction(ctx, request as $1.SignRawTransactionRequest);
      case 'ListTransactions':
        return listTransactions(ctx, request as $1.ListTransactionsRequest);
      default:
        throw $core.ArgumentError('Unknown method: $methodName');
    }
  }

  $core.Map<$core.String, $core.dynamic> get $json => WalletServiceBase$json;
  $core.Map<$core.String, $core.Map<$core.String, $core.dynamic>>
      get $messageJson => WalletServiceBase$messageJson;
}
