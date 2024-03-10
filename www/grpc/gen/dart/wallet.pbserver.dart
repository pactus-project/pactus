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
  $async.Future<$4.LoadWalletResponse> loadWallet($pb.ServerContext ctx, $4.LoadWalletRequest request);
  $async.Future<$4.UnloadWalletResponse> unloadWallet($pb.ServerContext ctx, $4.UnloadWalletRequest request);
  $async.Future<$4.LockWalletResponse> lockWallet($pb.ServerContext ctx, $4.LockWalletRequest request);
  $async.Future<$4.UnlockWalletResponse> unlockWallet($pb.ServerContext ctx, $4.UnlockWalletRequest request);
  $async.Future<$4.SignRawTransactionResponse> signRawTransaction($pb.ServerContext ctx, $4.SignRawTransactionRequest request);
  $async.Future<$4.GetValidatorAddressResponse> getValidatorAddress($pb.ServerContext ctx, $4.GetValidatorAddressRequest request);

  $pb.GeneratedMessage createRequest($core.String method) {
    switch (method) {
      case 'CreateWallet': return $4.CreateWalletRequest();
      case 'LoadWallet': return $4.LoadWalletRequest();
      case 'UnloadWallet': return $4.UnloadWalletRequest();
      case 'LockWallet': return $4.LockWalletRequest();
      case 'UnlockWallet': return $4.UnlockWalletRequest();
      case 'SignRawTransaction': return $4.SignRawTransactionRequest();
      case 'GetValidatorAddress': return $4.GetValidatorAddressRequest();
      default: throw $core.ArgumentError('Unknown method: $method');
    }
  }

  $async.Future<$pb.GeneratedMessage> handleCall($pb.ServerContext ctx, $core.String method, $pb.GeneratedMessage request) {
    switch (method) {
      case 'CreateWallet': return this.createWallet(ctx, request as $4.CreateWalletRequest);
      case 'LoadWallet': return this.loadWallet(ctx, request as $4.LoadWalletRequest);
      case 'UnloadWallet': return this.unloadWallet(ctx, request as $4.UnloadWalletRequest);
      case 'LockWallet': return this.lockWallet(ctx, request as $4.LockWalletRequest);
      case 'UnlockWallet': return this.unlockWallet(ctx, request as $4.UnlockWalletRequest);
      case 'SignRawTransaction': return this.signRawTransaction(ctx, request as $4.SignRawTransactionRequest);
      case 'GetValidatorAddress': return this.getValidatorAddress(ctx, request as $4.GetValidatorAddressRequest);
      default: throw $core.ArgumentError('Unknown method: $method');
    }
  }

  $core.Map<$core.String, $core.dynamic> get $json => WalletServiceBase$json;
  $core.Map<$core.String, $core.Map<$core.String, $core.dynamic>> get $messageJson => WalletServiceBase$messageJson;
}

