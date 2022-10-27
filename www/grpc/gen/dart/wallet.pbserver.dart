///
//  Generated code. Do not modify.
//  source: wallet.proto
//
// @dart = 2.12
// ignore_for_file: annotate_overrides,camel_case_types,constant_identifier_names,deprecated_member_use_from_same_package,directives_ordering,library_prefixes,non_constant_identifier_names,prefer_final_fields,return_of_invalid_type,unnecessary_const,unnecessary_import,unnecessary_this,unused_import,unused_shown_name

import 'dart:async' as $async;

import 'package:protobuf/protobuf.dart' as $pb;

import 'dart:core' as $core;
import 'wallet.pb.dart' as $3;
import 'wallet.pbjson.dart';

export 'wallet.pb.dart';

abstract class WalletServiceBase extends $pb.GeneratedService {
  $async.Future<$3.GenerateMnemonicResponse> generateMnemonic($pb.ServerContext ctx, $3.GenerateMnemonicRequest request);
  $async.Future<$3.CreateWalletResponse> createWallet($pb.ServerContext ctx, $3.CreateWalletRequest request);

  $pb.GeneratedMessage createRequest($core.String method) {
    switch (method) {
      case 'GenerateMnemonic': return $3.GenerateMnemonicRequest();
      case 'CreateWallet': return $3.CreateWalletRequest();
      default: throw $core.ArgumentError('Unknown method: $method');
    }
  }

  $async.Future<$pb.GeneratedMessage> handleCall($pb.ServerContext ctx, $core.String method, $pb.GeneratedMessage request) {
    switch (method) {
      case 'GenerateMnemonic': return this.generateMnemonic(ctx, request as $3.GenerateMnemonicRequest);
      case 'CreateWallet': return this.createWallet(ctx, request as $3.CreateWalletRequest);
      default: throw $core.ArgumentError('Unknown method: $method');
    }
  }

  $core.Map<$core.String, $core.dynamic> get $json => WalletServiceBase$json;
  $core.Map<$core.String, $core.Map<$core.String, $core.dynamic>> get $messageJson => WalletServiceBase$messageJson;
}

