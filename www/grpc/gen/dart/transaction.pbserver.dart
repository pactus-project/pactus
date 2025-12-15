// This is a generated file - do not edit.
//
// Generated from transaction.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:async' as $async;
import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

import 'transaction.pb.dart' as $0;
import 'transaction.pbjson.dart';

export 'transaction.pb.dart';

abstract class TransactionServiceBase extends $pb.GeneratedService {
  $async.Future<$0.GetTransactionResponse> getTransaction(
      $pb.ServerContext ctx, $0.GetTransactionRequest request);
  $async.Future<$0.CalculateFeeResponse> calculateFee(
      $pb.ServerContext ctx, $0.CalculateFeeRequest request);
  $async.Future<$0.BroadcastTransactionResponse> broadcastTransaction(
      $pb.ServerContext ctx, $0.BroadcastTransactionRequest request);
  $async.Future<$0.GetRawTransactionResponse> getRawTransferTransaction(
      $pb.ServerContext ctx, $0.GetRawTransferTransactionRequest request);
  $async.Future<$0.GetRawTransactionResponse> getRawBondTransaction(
      $pb.ServerContext ctx, $0.GetRawBondTransactionRequest request);
  $async.Future<$0.GetRawTransactionResponse> getRawUnbondTransaction(
      $pb.ServerContext ctx, $0.GetRawUnbondTransactionRequest request);
  $async.Future<$0.GetRawTransactionResponse> getRawWithdrawTransaction(
      $pb.ServerContext ctx, $0.GetRawWithdrawTransactionRequest request);
  $async.Future<$0.GetRawTransactionResponse> getRawBatchTransferTransaction(
      $pb.ServerContext ctx, $0.GetRawBatchTransferTransactionRequest request);
  $async.Future<$0.DecodeRawTransactionResponse> decodeRawTransaction(
      $pb.ServerContext ctx, $0.DecodeRawTransactionRequest request);

  $pb.GeneratedMessage createRequest($core.String methodName) {
    switch (methodName) {
      case 'GetTransaction':
        return $0.GetTransactionRequest();
      case 'CalculateFee':
        return $0.CalculateFeeRequest();
      case 'BroadcastTransaction':
        return $0.BroadcastTransactionRequest();
      case 'GetRawTransferTransaction':
        return $0.GetRawTransferTransactionRequest();
      case 'GetRawBondTransaction':
        return $0.GetRawBondTransactionRequest();
      case 'GetRawUnbondTransaction':
        return $0.GetRawUnbondTransactionRequest();
      case 'GetRawWithdrawTransaction':
        return $0.GetRawWithdrawTransactionRequest();
      case 'GetRawBatchTransferTransaction':
        return $0.GetRawBatchTransferTransactionRequest();
      case 'DecodeRawTransaction':
        return $0.DecodeRawTransactionRequest();
      default:
        throw $core.ArgumentError('Unknown method: $methodName');
    }
  }

  $async.Future<$pb.GeneratedMessage> handleCall($pb.ServerContext ctx,
      $core.String methodName, $pb.GeneratedMessage request) {
    switch (methodName) {
      case 'GetTransaction':
        return getTransaction(ctx, request as $0.GetTransactionRequest);
      case 'CalculateFee':
        return calculateFee(ctx, request as $0.CalculateFeeRequest);
      case 'BroadcastTransaction':
        return broadcastTransaction(
            ctx, request as $0.BroadcastTransactionRequest);
      case 'GetRawTransferTransaction':
        return getRawTransferTransaction(
            ctx, request as $0.GetRawTransferTransactionRequest);
      case 'GetRawBondTransaction':
        return getRawBondTransaction(
            ctx, request as $0.GetRawBondTransactionRequest);
      case 'GetRawUnbondTransaction':
        return getRawUnbondTransaction(
            ctx, request as $0.GetRawUnbondTransactionRequest);
      case 'GetRawWithdrawTransaction':
        return getRawWithdrawTransaction(
            ctx, request as $0.GetRawWithdrawTransactionRequest);
      case 'GetRawBatchTransferTransaction':
        return getRawBatchTransferTransaction(
            ctx, request as $0.GetRawBatchTransferTransactionRequest);
      case 'DecodeRawTransaction':
        return decodeRawTransaction(
            ctx, request as $0.DecodeRawTransactionRequest);
      default:
        throw $core.ArgumentError('Unknown method: $methodName');
    }
  }

  $core.Map<$core.String, $core.dynamic> get $json =>
      TransactionServiceBase$json;
  $core.Map<$core.String, $core.Map<$core.String, $core.dynamic>>
      get $messageJson => TransactionServiceBase$messageJson;
}
