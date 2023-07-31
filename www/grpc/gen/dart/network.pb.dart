///
//  Generated code. Do not modify.
//  source: network.proto
//
// @dart = 2.12
// ignore_for_file: annotate_overrides,camel_case_types,constant_identifier_names,directives_ordering,library_prefixes,non_constant_identifier_names,prefer_final_fields,return_of_invalid_type,unnecessary_const,unnecessary_import,unnecessary_this,unused_import,unused_shown_name

import 'dart:async' as $async;
import 'dart:core' as $core;

import 'package:fixnum/fixnum.dart' as $fixnum;
import 'package:protobuf/protobuf.dart' as $pb;

class GetNetworkInfoRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'GetNetworkInfoRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..hasRequiredFields = false
  ;

  GetNetworkInfoRequest._() : super();
  factory GetNetworkInfoRequest() => create();
  factory GetNetworkInfoRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetNetworkInfoRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetNetworkInfoRequest clone() => GetNetworkInfoRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetNetworkInfoRequest copyWith(void Function(GetNetworkInfoRequest) updates) => super.copyWith((message) => updates(message as GetNetworkInfoRequest)) as GetNetworkInfoRequest; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static GetNetworkInfoRequest create() => GetNetworkInfoRequest._();
  GetNetworkInfoRequest createEmptyInstance() => create();
  static $pb.PbList<GetNetworkInfoRequest> createRepeated() => $pb.PbList<GetNetworkInfoRequest>();
  @$core.pragma('dart2js:noInline')
  static GetNetworkInfoRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetNetworkInfoRequest>(create);
  static GetNetworkInfoRequest? _defaultInstance;
}

class GetNetworkInfoResponse extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'GetNetworkInfoResponse', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..a<$core.int>(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'totalSentBytes', $pb.PbFieldType.O3)
    ..a<$core.int>(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'totalReceivedBytes', $pb.PbFieldType.O3)
    ..aInt64(3, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'startedAt')
    ..pc<PeerInfo>(4, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'peers', $pb.PbFieldType.PM, subBuilder: PeerInfo.create)
    ..m<$core.int, $fixnum.Int64>(5, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'sentBytes', entryClassName: 'GetNetworkInfoResponse.SentBytesEntry', keyFieldType: $pb.PbFieldType.O3, valueFieldType: $pb.PbFieldType.O6, packageName: const $pb.PackageName('pactus'))
    ..m<$core.int, $fixnum.Int64>(6, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'receivedBytes', entryClassName: 'GetNetworkInfoResponse.ReceivedBytesEntry', keyFieldType: $pb.PbFieldType.O3, valueFieldType: $pb.PbFieldType.O6, packageName: const $pb.PackageName('pactus'))
    ..hasRequiredFields = false
  ;

  GetNetworkInfoResponse._() : super();
  factory GetNetworkInfoResponse({
    $core.int? totalSentBytes,
    $core.int? totalReceivedBytes,
    $fixnum.Int64? startedAt,
    $core.Iterable<PeerInfo>? peers,
    $core.Map<$core.int, $fixnum.Int64>? sentBytes,
    $core.Map<$core.int, $fixnum.Int64>? receivedBytes,
  }) {
    final _result = create();
    if (totalSentBytes != null) {
      _result.totalSentBytes = totalSentBytes;
    }
    if (totalReceivedBytes != null) {
      _result.totalReceivedBytes = totalReceivedBytes;
    }
    if (startedAt != null) {
      _result.startedAt = startedAt;
    }
    if (peers != null) {
      _result.peers.addAll(peers);
    }
    if (sentBytes != null) {
      _result.sentBytes.addAll(sentBytes);
    }
    if (receivedBytes != null) {
      _result.receivedBytes.addAll(receivedBytes);
    }
    return _result;
  }
  factory GetNetworkInfoResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetNetworkInfoResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetNetworkInfoResponse clone() => GetNetworkInfoResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetNetworkInfoResponse copyWith(void Function(GetNetworkInfoResponse) updates) => super.copyWith((message) => updates(message as GetNetworkInfoResponse)) as GetNetworkInfoResponse; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static GetNetworkInfoResponse create() => GetNetworkInfoResponse._();
  GetNetworkInfoResponse createEmptyInstance() => create();
  static $pb.PbList<GetNetworkInfoResponse> createRepeated() => $pb.PbList<GetNetworkInfoResponse>();
  @$core.pragma('dart2js:noInline')
  static GetNetworkInfoResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetNetworkInfoResponse>(create);
  static GetNetworkInfoResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.int get totalSentBytes => $_getIZ(0);
  @$pb.TagNumber(1)
  set totalSentBytes($core.int v) { $_setSignedInt32(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasTotalSentBytes() => $_has(0);
  @$pb.TagNumber(1)
  void clearTotalSentBytes() => clearField(1);

  @$pb.TagNumber(2)
  $core.int get totalReceivedBytes => $_getIZ(1);
  @$pb.TagNumber(2)
  set totalReceivedBytes($core.int v) { $_setSignedInt32(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasTotalReceivedBytes() => $_has(1);
  @$pb.TagNumber(2)
  void clearTotalReceivedBytes() => clearField(2);

  @$pb.TagNumber(3)
  $fixnum.Int64 get startedAt => $_getI64(2);
  @$pb.TagNumber(3)
  set startedAt($fixnum.Int64 v) { $_setInt64(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasStartedAt() => $_has(2);
  @$pb.TagNumber(3)
  void clearStartedAt() => clearField(3);

  @$pb.TagNumber(4)
  $core.List<PeerInfo> get peers => $_getList(3);

  @$pb.TagNumber(5)
  $core.Map<$core.int, $fixnum.Int64> get sentBytes => $_getMap(4);

  @$pb.TagNumber(6)
  $core.Map<$core.int, $fixnum.Int64> get receivedBytes => $_getMap(5);
}

class GetNodeInfoRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'GetNodeInfoRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..hasRequiredFields = false
  ;

  GetNodeInfoRequest._() : super();
  factory GetNodeInfoRequest() => create();
  factory GetNodeInfoRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetNodeInfoRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetNodeInfoRequest clone() => GetNodeInfoRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetNodeInfoRequest copyWith(void Function(GetNodeInfoRequest) updates) => super.copyWith((message) => updates(message as GetNodeInfoRequest)) as GetNodeInfoRequest; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static GetNodeInfoRequest create() => GetNodeInfoRequest._();
  GetNodeInfoRequest createEmptyInstance() => create();
  static $pb.PbList<GetNodeInfoRequest> createRepeated() => $pb.PbList<GetNodeInfoRequest>();
  @$core.pragma('dart2js:noInline')
  static GetNodeInfoRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetNodeInfoRequest>(create);
  static GetNodeInfoRequest? _defaultInstance;
}

class GetNodeInfoResponse extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'GetNodeInfoResponse', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'moniker')
    ..aOS(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'agent')
    ..a<$core.List<$core.int>>(3, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'peerId', $pb.PbFieldType.OY)
    ..hasRequiredFields = false
  ;

  GetNodeInfoResponse._() : super();
  factory GetNodeInfoResponse({
    $core.String? moniker,
    $core.String? agent,
    $core.List<$core.int>? peerId,
  }) {
    final _result = create();
    if (moniker != null) {
      _result.moniker = moniker;
    }
    if (agent != null) {
      _result.agent = agent;
    }
    if (peerId != null) {
      _result.peerId = peerId;
    }
    return _result;
  }
  factory GetNodeInfoResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetNodeInfoResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetNodeInfoResponse clone() => GetNodeInfoResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetNodeInfoResponse copyWith(void Function(GetNodeInfoResponse) updates) => super.copyWith((message) => updates(message as GetNodeInfoResponse)) as GetNodeInfoResponse; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static GetNodeInfoResponse create() => GetNodeInfoResponse._();
  GetNodeInfoResponse createEmptyInstance() => create();
  static $pb.PbList<GetNodeInfoResponse> createRepeated() => $pb.PbList<GetNodeInfoResponse>();
  @$core.pragma('dart2js:noInline')
  static GetNodeInfoResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetNodeInfoResponse>(create);
  static GetNodeInfoResponse? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get moniker => $_getSZ(0);
  @$pb.TagNumber(1)
  set moniker($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasMoniker() => $_has(0);
  @$pb.TagNumber(1)
  void clearMoniker() => clearField(1);

  @$pb.TagNumber(2)
  $core.String get agent => $_getSZ(1);
  @$pb.TagNumber(2)
  set agent($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasAgent() => $_has(1);
  @$pb.TagNumber(2)
  void clearAgent() => clearField(2);

  @$pb.TagNumber(3)
  $core.List<$core.int> get peerId => $_getN(2);
  @$pb.TagNumber(3)
  set peerId($core.List<$core.int> v) { $_setBytes(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasPeerId() => $_has(2);
  @$pb.TagNumber(3)
  void clearPeerId() => clearField(3);
}

class PeerInfo extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'PeerInfo', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'moniker')
    ..aOS(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'agent')
    ..a<$core.List<$core.int>>(3, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'peerId', $pb.PbFieldType.OY)
    ..pPS(4, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'consensusKeys')
    ..a<$core.int>(5, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'flags', $pb.PbFieldType.O3)
    ..a<$core.int>(6, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'height', $pb.PbFieldType.OU3)
    ..a<$core.int>(7, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'receivedMessages', $pb.PbFieldType.O3)
    ..a<$core.int>(8, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'invalidMessages', $pb.PbFieldType.O3)
    ..a<$core.int>(9, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'receivedBytes', $pb.PbFieldType.O3)
    ..a<$core.int>(10, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'sentBytes', $pb.PbFieldType.O3)
    ..a<$core.int>(11, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'status', $pb.PbFieldType.O3)
    ..aInt64(12, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'lastSent')
    ..aInt64(13, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'lastReceived')
    ..a<$core.int>(14, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'sendSuccess', $pb.PbFieldType.O3)
    ..a<$core.int>(15, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'sendFailed', $pb.PbFieldType.O3)
    ..a<$core.List<$core.int>>(16, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'lastBlockHash', $pb.PbFieldType.OY)
    ..hasRequiredFields = false
  ;

  PeerInfo._() : super();
  factory PeerInfo({
    $core.String? moniker,
    $core.String? agent,
    $core.List<$core.int>? peerId,
    $core.Iterable<$core.String>? consensusKeys,
    $core.int? flags,
    $core.int? height,
    $core.int? receivedMessages,
    $core.int? invalidMessages,
    $core.int? receivedBytes,
    $core.int? sentBytes,
    $core.int? status,
    $fixnum.Int64? lastSent,
    $fixnum.Int64? lastReceived,
    $core.int? sendSuccess,
    $core.int? sendFailed,
    $core.List<$core.int>? lastBlockHash,
  }) {
    final _result = create();
    if (moniker != null) {
      _result.moniker = moniker;
    }
    if (agent != null) {
      _result.agent = agent;
    }
    if (peerId != null) {
      _result.peerId = peerId;
    }
    if (consensusKeys != null) {
      _result.consensusKeys.addAll(consensusKeys);
    }
    if (flags != null) {
      _result.flags = flags;
    }
    if (height != null) {
      _result.height = height;
    }
    if (receivedMessages != null) {
      _result.receivedMessages = receivedMessages;
    }
    if (invalidMessages != null) {
      _result.invalidMessages = invalidMessages;
    }
    if (receivedBytes != null) {
      _result.receivedBytes = receivedBytes;
    }
    if (sentBytes != null) {
      _result.sentBytes = sentBytes;
    }
    if (status != null) {
      _result.status = status;
    }
    if (lastSent != null) {
      _result.lastSent = lastSent;
    }
    if (lastReceived != null) {
      _result.lastReceived = lastReceived;
    }
    if (sendSuccess != null) {
      _result.sendSuccess = sendSuccess;
    }
    if (sendFailed != null) {
      _result.sendFailed = sendFailed;
    }
    if (lastBlockHash != null) {
      _result.lastBlockHash = lastBlockHash;
    }
    return _result;
  }
  factory PeerInfo.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory PeerInfo.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  PeerInfo clone() => PeerInfo()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  PeerInfo copyWith(void Function(PeerInfo) updates) => super.copyWith((message) => updates(message as PeerInfo)) as PeerInfo; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static PeerInfo create() => PeerInfo._();
  PeerInfo createEmptyInstance() => create();
  static $pb.PbList<PeerInfo> createRepeated() => $pb.PbList<PeerInfo>();
  @$core.pragma('dart2js:noInline')
  static PeerInfo getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<PeerInfo>(create);
  static PeerInfo? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get moniker => $_getSZ(0);
  @$pb.TagNumber(1)
  set moniker($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasMoniker() => $_has(0);
  @$pb.TagNumber(1)
  void clearMoniker() => clearField(1);

  @$pb.TagNumber(2)
  $core.String get agent => $_getSZ(1);
  @$pb.TagNumber(2)
  set agent($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasAgent() => $_has(1);
  @$pb.TagNumber(2)
  void clearAgent() => clearField(2);

  @$pb.TagNumber(3)
  $core.List<$core.int> get peerId => $_getN(2);
  @$pb.TagNumber(3)
  set peerId($core.List<$core.int> v) { $_setBytes(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasPeerId() => $_has(2);
  @$pb.TagNumber(3)
  void clearPeerId() => clearField(3);

  @$pb.TagNumber(4)
  $core.List<$core.String> get consensusKeys => $_getList(3);

  @$pb.TagNumber(5)
  $core.int get flags => $_getIZ(4);
  @$pb.TagNumber(5)
  set flags($core.int v) { $_setSignedInt32(4, v); }
  @$pb.TagNumber(5)
  $core.bool hasFlags() => $_has(4);
  @$pb.TagNumber(5)
  void clearFlags() => clearField(5);

  @$pb.TagNumber(6)
  $core.int get height => $_getIZ(5);
  @$pb.TagNumber(6)
  set height($core.int v) { $_setUnsignedInt32(5, v); }
  @$pb.TagNumber(6)
  $core.bool hasHeight() => $_has(5);
  @$pb.TagNumber(6)
  void clearHeight() => clearField(6);

  @$pb.TagNumber(7)
  $core.int get receivedMessages => $_getIZ(6);
  @$pb.TagNumber(7)
  set receivedMessages($core.int v) { $_setSignedInt32(6, v); }
  @$pb.TagNumber(7)
  $core.bool hasReceivedMessages() => $_has(6);
  @$pb.TagNumber(7)
  void clearReceivedMessages() => clearField(7);

  @$pb.TagNumber(8)
  $core.int get invalidMessages => $_getIZ(7);
  @$pb.TagNumber(8)
  set invalidMessages($core.int v) { $_setSignedInt32(7, v); }
  @$pb.TagNumber(8)
  $core.bool hasInvalidMessages() => $_has(7);
  @$pb.TagNumber(8)
  void clearInvalidMessages() => clearField(8);

  @$pb.TagNumber(9)
  $core.int get receivedBytes => $_getIZ(8);
  @$pb.TagNumber(9)
  set receivedBytes($core.int v) { $_setSignedInt32(8, v); }
  @$pb.TagNumber(9)
  $core.bool hasReceivedBytes() => $_has(8);
  @$pb.TagNumber(9)
  void clearReceivedBytes() => clearField(9);

  @$pb.TagNumber(10)
  $core.int get sentBytes => $_getIZ(9);
  @$pb.TagNumber(10)
  set sentBytes($core.int v) { $_setSignedInt32(9, v); }
  @$pb.TagNumber(10)
  $core.bool hasSentBytes() => $_has(9);
  @$pb.TagNumber(10)
  void clearSentBytes() => clearField(10);

  @$pb.TagNumber(11)
  $core.int get status => $_getIZ(10);
  @$pb.TagNumber(11)
  set status($core.int v) { $_setSignedInt32(10, v); }
  @$pb.TagNumber(11)
  $core.bool hasStatus() => $_has(10);
  @$pb.TagNumber(11)
  void clearStatus() => clearField(11);

  @$pb.TagNumber(12)
  $fixnum.Int64 get lastSent => $_getI64(11);
  @$pb.TagNumber(12)
  set lastSent($fixnum.Int64 v) { $_setInt64(11, v); }
  @$pb.TagNumber(12)
  $core.bool hasLastSent() => $_has(11);
  @$pb.TagNumber(12)
  void clearLastSent() => clearField(12);

  @$pb.TagNumber(13)
  $fixnum.Int64 get lastReceived => $_getI64(12);
  @$pb.TagNumber(13)
  set lastReceived($fixnum.Int64 v) { $_setInt64(12, v); }
  @$pb.TagNumber(13)
  $core.bool hasLastReceived() => $_has(12);
  @$pb.TagNumber(13)
  void clearLastReceived() => clearField(13);

  @$pb.TagNumber(14)
  $core.int get sendSuccess => $_getIZ(13);
  @$pb.TagNumber(14)
  set sendSuccess($core.int v) { $_setSignedInt32(13, v); }
  @$pb.TagNumber(14)
  $core.bool hasSendSuccess() => $_has(13);
  @$pb.TagNumber(14)
  void clearSendSuccess() => clearField(14);

  @$pb.TagNumber(15)
  $core.int get sendFailed => $_getIZ(14);
  @$pb.TagNumber(15)
  set sendFailed($core.int v) { $_setSignedInt32(14, v); }
  @$pb.TagNumber(15)
  $core.bool hasSendFailed() => $_has(14);
  @$pb.TagNumber(15)
  void clearSendFailed() => clearField(15);

  @$pb.TagNumber(16)
  $core.List<$core.int> get lastBlockHash => $_getN(15);
  @$pb.TagNumber(16)
  set lastBlockHash($core.List<$core.int> v) { $_setBytes(15, v); }
  @$pb.TagNumber(16)
  $core.bool hasLastBlockHash() => $_has(15);
  @$pb.TagNumber(16)
  void clearLastBlockHash() => clearField(16);
}

class NetworkApi {
  $pb.RpcClient _client;
  NetworkApi(this._client);

  $async.Future<GetNetworkInfoResponse> getNetworkInfo($pb.ClientContext? ctx, GetNetworkInfoRequest request) {
    var emptyResponse = GetNetworkInfoResponse();
    return _client.invoke<GetNetworkInfoResponse>(ctx, 'Network', 'GetNetworkInfo', request, emptyResponse);
  }
  $async.Future<GetNodeInfoResponse> getNodeInfo($pb.ClientContext? ctx, GetNodeInfoRequest request) {
    var emptyResponse = GetNodeInfoResponse();
    return _client.invoke<GetNodeInfoResponse>(ctx, 'Network', 'GetNodeInfo', request, emptyResponse);
  }
}

