///
//  Generated code. Do not modify.
//  source: utility.proto
//
// @dart = 2.12
// ignore_for_file: annotate_overrides,camel_case_types,constant_identifier_names,directives_ordering,library_prefixes,non_constant_identifier_names,prefer_final_fields,return_of_invalid_type,unnecessary_const,unnecessary_import,unnecessary_this,unused_import,unused_shown_name

import 'dart:async' as $async;
import 'dart:core' as $core;

import 'package:fixnum/fixnum.dart' as $fixnum;
import 'package:protobuf/protobuf.dart' as $pb;

import 'transaction.pbenum.dart' as $0;

class CalculateFeeRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'CalculateFeeRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aInt64(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'amount')
    ..e<$0.PayloadType>(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'payloadType', $pb.PbFieldType.OE, defaultOrMaker: $0.PayloadType.UNKNOWN, valueOf: $0.PayloadType.valueOf, enumValues: $0.PayloadType.values)
    ..aOB(3, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'fixedAmount')
    ..hasRequiredFields = false
  ;

  CalculateFeeRequest._() : super();
  factory CalculateFeeRequest({
    $fixnum.Int64? amount,
    $0.PayloadType? payloadType,
    $core.bool? fixedAmount,
  }) {
    final _result = create();
    if (amount != null) {
      _result.amount = amount;
    }
    if (payloadType != null) {
      _result.payloadType = payloadType;
    }
    if (fixedAmount != null) {
      _result.fixedAmount = fixedAmount;
    }
    return _result;
  }
  factory CalculateFeeRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory CalculateFeeRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  CalculateFeeRequest clone() => CalculateFeeRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  CalculateFeeRequest copyWith(void Function(CalculateFeeRequest) updates) => super.copyWith((message) => updates(message as CalculateFeeRequest)) as CalculateFeeRequest; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static CalculateFeeRequest create() => CalculateFeeRequest._();
  CalculateFeeRequest createEmptyInstance() => create();
  static $pb.PbList<CalculateFeeRequest> createRepeated() => $pb.PbList<CalculateFeeRequest>();
  @$core.pragma('dart2js:noInline')
  static CalculateFeeRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<CalculateFeeRequest>(create);
  static CalculateFeeRequest? _defaultInstance;

  @$pb.TagNumber(1)
  $fixnum.Int64 get amount => $_getI64(0);
  @$pb.TagNumber(1)
  set amount($fixnum.Int64 v) { $_setInt64(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasAmount() => $_has(0);
  @$pb.TagNumber(1)
  void clearAmount() => clearField(1);

  @$pb.TagNumber(2)
  $0.PayloadType get payloadType => $_getN(1);
  @$pb.TagNumber(2)
  set payloadType($0.PayloadType v) { setField(2, v); }
  @$pb.TagNumber(2)
  $core.bool hasPayloadType() => $_has(1);
  @$pb.TagNumber(2)
  void clearPayloadType() => clearField(2);

  @$pb.TagNumber(3)
  $core.bool get fixedAmount => $_getBF(2);
  @$pb.TagNumber(3)
  set fixedAmount($core.bool v) { $_setBool(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasFixedAmount() => $_has(2);
  @$pb.TagNumber(3)
  void clearFixedAmount() => clearField(3);
}

class CalculateFeeResponse extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'CalculateFeeResponse', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aInt64(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'amount')
    ..aInt64(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'fee')
    ..hasRequiredFields = false
  ;

  CalculateFeeResponse._() : super();
  factory CalculateFeeResponse({
    $fixnum.Int64? amount,
    $fixnum.Int64? fee,
  }) {
    final _result = create();
    if (amount != null) {
      _result.amount = amount;
    }
    if (fee != null) {
      _result.fee = fee;
    }
    return _result;
  }
  factory CalculateFeeResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory CalculateFeeResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  CalculateFeeResponse clone() => CalculateFeeResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  CalculateFeeResponse copyWith(void Function(CalculateFeeResponse) updates) => super.copyWith((message) => updates(message as CalculateFeeResponse)) as CalculateFeeResponse; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static CalculateFeeResponse create() => CalculateFeeResponse._();
  CalculateFeeResponse createEmptyInstance() => create();
  static $pb.PbList<CalculateFeeResponse> createRepeated() => $pb.PbList<CalculateFeeResponse>();
  @$core.pragma('dart2js:noInline')
  static CalculateFeeResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<CalculateFeeResponse>(create);
  static CalculateFeeResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $fixnum.Int64 get amount => $_getI64(0);
  @$pb.TagNumber(1)
  set amount($fixnum.Int64 v) { $_setInt64(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasAmount() => $_has(0);
  @$pb.TagNumber(1)
  void clearAmount() => clearField(1);

  @$pb.TagNumber(2)
  $fixnum.Int64 get fee => $_getI64(1);
  @$pb.TagNumber(2)
  set fee($fixnum.Int64 v) { $_setInt64(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasFee() => $_has(1);
  @$pb.TagNumber(2)
  void clearFee() => clearField(2);
}

class UtilityApi {
  $pb.RpcClient _client;
  UtilityApi(this._client);

  $async.Future<CalculateFeeResponse> calculateFee($pb.ClientContext? ctx, CalculateFeeRequest request) {
    var emptyResponse = CalculateFeeResponse();
    return _client.invoke<CalculateFeeResponse>(ctx, 'Utility', 'CalculateFee', request, emptyResponse);
  }
}

