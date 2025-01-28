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
    ..aOB(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'onlyConnected')
    ..hasRequiredFields = false
  ;

  GetNetworkInfoRequest._() : super();
  factory GetNetworkInfoRequest({
    $core.bool? onlyConnected,
  }) {
    final _result = create();
    if (onlyConnected != null) {
      _result.onlyConnected = onlyConnected;
    }
    return _result;
  }
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

  @$pb.TagNumber(1)
  $core.bool get onlyConnected => $_getBF(0);
  @$pb.TagNumber(1)
  set onlyConnected($core.bool v) { $_setBool(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasOnlyConnected() => $_has(0);
  @$pb.TagNumber(1)
  void clearOnlyConnected() => clearField(1);
}

class GetNetworkInfoResponse extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'GetNetworkInfoResponse', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'networkName')
    ..a<$core.int>(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'connectedPeersCount', $pb.PbFieldType.OU3)
    ..pc<PeerInfo>(3, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'connectedPeers', $pb.PbFieldType.PM, subBuilder: PeerInfo.create)
    ..aOM<MetricInfo>(4, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'metricInfo', subBuilder: MetricInfo.create)
    ..hasRequiredFields = false
  ;

  GetNetworkInfoResponse._() : super();
  factory GetNetworkInfoResponse({
    $core.String? networkName,
    $core.int? connectedPeersCount,
    $core.Iterable<PeerInfo>? connectedPeers,
    MetricInfo? metricInfo,
  }) {
    final _result = create();
    if (networkName != null) {
      _result.networkName = networkName;
    }
    if (connectedPeersCount != null) {
      _result.connectedPeersCount = connectedPeersCount;
    }
    if (connectedPeers != null) {
      _result.connectedPeers.addAll(connectedPeers);
    }
    if (metricInfo != null) {
      _result.metricInfo = metricInfo;
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
  $core.String get networkName => $_getSZ(0);
  @$pb.TagNumber(1)
  set networkName($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasNetworkName() => $_has(0);
  @$pb.TagNumber(1)
  void clearNetworkName() => clearField(1);

  @$pb.TagNumber(2)
  $core.int get connectedPeersCount => $_getIZ(1);
  @$pb.TagNumber(2)
  set connectedPeersCount($core.int v) { $_setUnsignedInt32(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasConnectedPeersCount() => $_has(1);
  @$pb.TagNumber(2)
  void clearConnectedPeersCount() => clearField(2);

  @$pb.TagNumber(3)
  $core.List<PeerInfo> get connectedPeers => $_getList(2);

  @$pb.TagNumber(4)
  MetricInfo get metricInfo => $_getN(3);
  @$pb.TagNumber(4)
  set metricInfo(MetricInfo v) { setField(4, v); }
  @$pb.TagNumber(4)
  $core.bool hasMetricInfo() => $_has(3);
  @$pb.TagNumber(4)
  void clearMetricInfo() => clearField(4);
  @$pb.TagNumber(4)
  MetricInfo ensureMetricInfo() => $_ensure(3);
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
    ..aOS(3, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'peerId')
    ..a<$fixnum.Int64>(4, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'startedAt', $pb.PbFieldType.OU6, defaultOrMaker: $fixnum.Int64.ZERO)
    ..aOS(5, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'reachability')
    ..a<$core.int>(6, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'services', $pb.PbFieldType.O3)
    ..aOS(7, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'servicesNames')
    ..pPS(8, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'localAddrs')
    ..pPS(9, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'protocols')
    ..a<$core.double>(13, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'clockOffset', $pb.PbFieldType.OD)
    ..aOM<ConnectionInfo>(14, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'connectionInfo', subBuilder: ConnectionInfo.create)
    ..pc<ZMQPublisherInfo>(15, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'zmqPublishers', $pb.PbFieldType.PM, subBuilder: ZMQPublisherInfo.create)
    ..hasRequiredFields = false
  ;

  GetNodeInfoResponse._() : super();
  factory GetNodeInfoResponse({
    $core.String? moniker,
    $core.String? agent,
    $core.String? peerId,
    $fixnum.Int64? startedAt,
    $core.String? reachability,
    $core.int? services,
    $core.String? servicesNames,
    $core.Iterable<$core.String>? localAddrs,
    $core.Iterable<$core.String>? protocols,
    $core.double? clockOffset,
    ConnectionInfo? connectionInfo,
    $core.Iterable<ZMQPublisherInfo>? zmqPublishers,
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
    if (startedAt != null) {
      _result.startedAt = startedAt;
    }
    if (reachability != null) {
      _result.reachability = reachability;
    }
    if (services != null) {
      _result.services = services;
    }
    if (servicesNames != null) {
      _result.servicesNames = servicesNames;
    }
    if (localAddrs != null) {
      _result.localAddrs.addAll(localAddrs);
    }
    if (protocols != null) {
      _result.protocols.addAll(protocols);
    }
    if (clockOffset != null) {
      _result.clockOffset = clockOffset;
    }
    if (connectionInfo != null) {
      _result.connectionInfo = connectionInfo;
    }
    if (zmqPublishers != null) {
      _result.zmqPublishers.addAll(zmqPublishers);
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
  $core.String get peerId => $_getSZ(2);
  @$pb.TagNumber(3)
  set peerId($core.String v) { $_setString(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasPeerId() => $_has(2);
  @$pb.TagNumber(3)
  void clearPeerId() => clearField(3);

  @$pb.TagNumber(4)
  $fixnum.Int64 get startedAt => $_getI64(3);
  @$pb.TagNumber(4)
  set startedAt($fixnum.Int64 v) { $_setInt64(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasStartedAt() => $_has(3);
  @$pb.TagNumber(4)
  void clearStartedAt() => clearField(4);

  @$pb.TagNumber(5)
  $core.String get reachability => $_getSZ(4);
  @$pb.TagNumber(5)
  set reachability($core.String v) { $_setString(4, v); }
  @$pb.TagNumber(5)
  $core.bool hasReachability() => $_has(4);
  @$pb.TagNumber(5)
  void clearReachability() => clearField(5);

  @$pb.TagNumber(6)
  $core.int get services => $_getIZ(5);
  @$pb.TagNumber(6)
  set services($core.int v) { $_setSignedInt32(5, v); }
  @$pb.TagNumber(6)
  $core.bool hasServices() => $_has(5);
  @$pb.TagNumber(6)
  void clearServices() => clearField(6);

  @$pb.TagNumber(7)
  $core.String get servicesNames => $_getSZ(6);
  @$pb.TagNumber(7)
  set servicesNames($core.String v) { $_setString(6, v); }
  @$pb.TagNumber(7)
  $core.bool hasServicesNames() => $_has(6);
  @$pb.TagNumber(7)
  void clearServicesNames() => clearField(7);

  @$pb.TagNumber(8)
  $core.List<$core.String> get localAddrs => $_getList(7);

  @$pb.TagNumber(9)
  $core.List<$core.String> get protocols => $_getList(8);

  @$pb.TagNumber(13)
  $core.double get clockOffset => $_getN(9);
  @$pb.TagNumber(13)
  set clockOffset($core.double v) { $_setDouble(9, v); }
  @$pb.TagNumber(13)
  $core.bool hasClockOffset() => $_has(9);
  @$pb.TagNumber(13)
  void clearClockOffset() => clearField(13);

  @$pb.TagNumber(14)
  ConnectionInfo get connectionInfo => $_getN(10);
  @$pb.TagNumber(14)
  set connectionInfo(ConnectionInfo v) { setField(14, v); }
  @$pb.TagNumber(14)
  $core.bool hasConnectionInfo() => $_has(10);
  @$pb.TagNumber(14)
  void clearConnectionInfo() => clearField(14);
  @$pb.TagNumber(14)
  ConnectionInfo ensureConnectionInfo() => $_ensure(10);

  @$pb.TagNumber(15)
  $core.List<ZMQPublisherInfo> get zmqPublishers => $_getList(11);
}

class ZMQPublisherInfo extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'ZMQPublisherInfo', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'topic')
    ..aOS(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'address')
    ..a<$core.int>(3, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'hwm', $pb.PbFieldType.O3)
    ..hasRequiredFields = false
  ;

  ZMQPublisherInfo._() : super();
  factory ZMQPublisherInfo({
    $core.String? topic,
    $core.String? address,
    $core.int? hwm,
  }) {
    final _result = create();
    if (topic != null) {
      _result.topic = topic;
    }
    if (address != null) {
      _result.address = address;
    }
    if (hwm != null) {
      _result.hwm = hwm;
    }
    return _result;
  }
  factory ZMQPublisherInfo.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory ZMQPublisherInfo.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  ZMQPublisherInfo clone() => ZMQPublisherInfo()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  ZMQPublisherInfo copyWith(void Function(ZMQPublisherInfo) updates) => super.copyWith((message) => updates(message as ZMQPublisherInfo)) as ZMQPublisherInfo; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static ZMQPublisherInfo create() => ZMQPublisherInfo._();
  ZMQPublisherInfo createEmptyInstance() => create();
  static $pb.PbList<ZMQPublisherInfo> createRepeated() => $pb.PbList<ZMQPublisherInfo>();
  @$core.pragma('dart2js:noInline')
  static ZMQPublisherInfo getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ZMQPublisherInfo>(create);
  static ZMQPublisherInfo? _defaultInstance;

  @$pb.TagNumber(1)
  $core.String get topic => $_getSZ(0);
  @$pb.TagNumber(1)
  set topic($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasTopic() => $_has(0);
  @$pb.TagNumber(1)
  void clearTopic() => clearField(1);

  @$pb.TagNumber(2)
  $core.String get address => $_getSZ(1);
  @$pb.TagNumber(2)
  set address($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasAddress() => $_has(1);
  @$pb.TagNumber(2)
  void clearAddress() => clearField(2);

  @$pb.TagNumber(3)
  $core.int get hwm => $_getIZ(2);
  @$pb.TagNumber(3)
  set hwm($core.int v) { $_setSignedInt32(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasHwm() => $_has(2);
  @$pb.TagNumber(3)
  void clearHwm() => clearField(3);
}

class PeerInfo extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'PeerInfo', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..a<$core.int>(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'status', $pb.PbFieldType.O3)
    ..aOS(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'moniker')
    ..aOS(3, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'agent')
    ..aOS(4, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'peerId')
    ..pPS(5, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'consensusKeys')
    ..pPS(6, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'consensusAddresses')
    ..a<$core.int>(7, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'services', $pb.PbFieldType.OU3)
    ..aOS(8, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'lastBlockHash')
    ..a<$core.int>(9, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'height', $pb.PbFieldType.OU3)
    ..aInt64(10, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'lastSent')
    ..aInt64(11, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'lastReceived')
    ..aOS(12, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'address')
    ..aOS(13, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'direction')
    ..pPS(14, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'protocols')
    ..a<$core.int>(15, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'totalSessions', $pb.PbFieldType.O3)
    ..a<$core.int>(16, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'completedSessions', $pb.PbFieldType.O3)
    ..aOM<MetricInfo>(17, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'metricInfo', subBuilder: MetricInfo.create)
    ..hasRequiredFields = false
  ;

  PeerInfo._() : super();
  factory PeerInfo({
    $core.int? status,
    $core.String? moniker,
    $core.String? agent,
    $core.String? peerId,
    $core.Iterable<$core.String>? consensusKeys,
    $core.Iterable<$core.String>? consensusAddresses,
    $core.int? services,
    $core.String? lastBlockHash,
    $core.int? height,
    $fixnum.Int64? lastSent,
    $fixnum.Int64? lastReceived,
    $core.String? address,
    $core.String? direction,
    $core.Iterable<$core.String>? protocols,
    $core.int? totalSessions,
    $core.int? completedSessions,
    MetricInfo? metricInfo,
  }) {
    final _result = create();
    if (status != null) {
      _result.status = status;
    }
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
    if (consensusAddresses != null) {
      _result.consensusAddresses.addAll(consensusAddresses);
    }
    if (services != null) {
      _result.services = services;
    }
    if (lastBlockHash != null) {
      _result.lastBlockHash = lastBlockHash;
    }
    if (height != null) {
      _result.height = height;
    }
    if (lastSent != null) {
      _result.lastSent = lastSent;
    }
    if (lastReceived != null) {
      _result.lastReceived = lastReceived;
    }
    if (address != null) {
      _result.address = address;
    }
    if (direction != null) {
      _result.direction = direction;
    }
    if (protocols != null) {
      _result.protocols.addAll(protocols);
    }
    if (totalSessions != null) {
      _result.totalSessions = totalSessions;
    }
    if (completedSessions != null) {
      _result.completedSessions = completedSessions;
    }
    if (metricInfo != null) {
      _result.metricInfo = metricInfo;
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
  $core.int get status => $_getIZ(0);
  @$pb.TagNumber(1)
  set status($core.int v) { $_setSignedInt32(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasStatus() => $_has(0);
  @$pb.TagNumber(1)
  void clearStatus() => clearField(1);

  @$pb.TagNumber(2)
  $core.String get moniker => $_getSZ(1);
  @$pb.TagNumber(2)
  set moniker($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasMoniker() => $_has(1);
  @$pb.TagNumber(2)
  void clearMoniker() => clearField(2);

  @$pb.TagNumber(3)
  $core.String get agent => $_getSZ(2);
  @$pb.TagNumber(3)
  set agent($core.String v) { $_setString(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasAgent() => $_has(2);
  @$pb.TagNumber(3)
  void clearAgent() => clearField(3);

  @$pb.TagNumber(4)
  $core.String get peerId => $_getSZ(3);
  @$pb.TagNumber(4)
  set peerId($core.String v) { $_setString(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasPeerId() => $_has(3);
  @$pb.TagNumber(4)
  void clearPeerId() => clearField(4);

  @$pb.TagNumber(5)
  $core.List<$core.String> get consensusKeys => $_getList(4);

  @$pb.TagNumber(6)
  $core.List<$core.String> get consensusAddresses => $_getList(5);

  @$pb.TagNumber(7)
  $core.int get services => $_getIZ(6);
  @$pb.TagNumber(7)
  set services($core.int v) { $_setUnsignedInt32(6, v); }
  @$pb.TagNumber(7)
  $core.bool hasServices() => $_has(6);
  @$pb.TagNumber(7)
  void clearServices() => clearField(7);

  @$pb.TagNumber(8)
  $core.String get lastBlockHash => $_getSZ(7);
  @$pb.TagNumber(8)
  set lastBlockHash($core.String v) { $_setString(7, v); }
  @$pb.TagNumber(8)
  $core.bool hasLastBlockHash() => $_has(7);
  @$pb.TagNumber(8)
  void clearLastBlockHash() => clearField(8);

  @$pb.TagNumber(9)
  $core.int get height => $_getIZ(8);
  @$pb.TagNumber(9)
  set height($core.int v) { $_setUnsignedInt32(8, v); }
  @$pb.TagNumber(9)
  $core.bool hasHeight() => $_has(8);
  @$pb.TagNumber(9)
  void clearHeight() => clearField(9);

  @$pb.TagNumber(10)
  $fixnum.Int64 get lastSent => $_getI64(9);
  @$pb.TagNumber(10)
  set lastSent($fixnum.Int64 v) { $_setInt64(9, v); }
  @$pb.TagNumber(10)
  $core.bool hasLastSent() => $_has(9);
  @$pb.TagNumber(10)
  void clearLastSent() => clearField(10);

  @$pb.TagNumber(11)
  $fixnum.Int64 get lastReceived => $_getI64(10);
  @$pb.TagNumber(11)
  set lastReceived($fixnum.Int64 v) { $_setInt64(10, v); }
  @$pb.TagNumber(11)
  $core.bool hasLastReceived() => $_has(10);
  @$pb.TagNumber(11)
  void clearLastReceived() => clearField(11);

  @$pb.TagNumber(12)
  $core.String get address => $_getSZ(11);
  @$pb.TagNumber(12)
  set address($core.String v) { $_setString(11, v); }
  @$pb.TagNumber(12)
  $core.bool hasAddress() => $_has(11);
  @$pb.TagNumber(12)
  void clearAddress() => clearField(12);

  @$pb.TagNumber(13)
  $core.String get direction => $_getSZ(12);
  @$pb.TagNumber(13)
  set direction($core.String v) { $_setString(12, v); }
  @$pb.TagNumber(13)
  $core.bool hasDirection() => $_has(12);
  @$pb.TagNumber(13)
  void clearDirection() => clearField(13);

  @$pb.TagNumber(14)
  $core.List<$core.String> get protocols => $_getList(13);

  @$pb.TagNumber(15)
  $core.int get totalSessions => $_getIZ(14);
  @$pb.TagNumber(15)
  set totalSessions($core.int v) { $_setSignedInt32(14, v); }
  @$pb.TagNumber(15)
  $core.bool hasTotalSessions() => $_has(14);
  @$pb.TagNumber(15)
  void clearTotalSessions() => clearField(15);

  @$pb.TagNumber(16)
  $core.int get completedSessions => $_getIZ(15);
  @$pb.TagNumber(16)
  set completedSessions($core.int v) { $_setSignedInt32(15, v); }
  @$pb.TagNumber(16)
  $core.bool hasCompletedSessions() => $_has(15);
  @$pb.TagNumber(16)
  void clearCompletedSessions() => clearField(16);

  @$pb.TagNumber(17)
  MetricInfo get metricInfo => $_getN(16);
  @$pb.TagNumber(17)
  set metricInfo(MetricInfo v) { setField(17, v); }
  @$pb.TagNumber(17)
  $core.bool hasMetricInfo() => $_has(16);
  @$pb.TagNumber(17)
  void clearMetricInfo() => clearField(17);
  @$pb.TagNumber(17)
  MetricInfo ensureMetricInfo() => $_ensure(16);
}

class ConnectionInfo extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'ConnectionInfo', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..a<$fixnum.Int64>(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'connections', $pb.PbFieldType.OU6, defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'inboundConnections', $pb.PbFieldType.OU6, defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(3, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'outboundConnections', $pb.PbFieldType.OU6, defaultOrMaker: $fixnum.Int64.ZERO)
    ..hasRequiredFields = false
  ;

  ConnectionInfo._() : super();
  factory ConnectionInfo({
    $fixnum.Int64? connections,
    $fixnum.Int64? inboundConnections,
    $fixnum.Int64? outboundConnections,
  }) {
    final _result = create();
    if (connections != null) {
      _result.connections = connections;
    }
    if (inboundConnections != null) {
      _result.inboundConnections = inboundConnections;
    }
    if (outboundConnections != null) {
      _result.outboundConnections = outboundConnections;
    }
    return _result;
  }
  factory ConnectionInfo.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory ConnectionInfo.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  ConnectionInfo clone() => ConnectionInfo()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  ConnectionInfo copyWith(void Function(ConnectionInfo) updates) => super.copyWith((message) => updates(message as ConnectionInfo)) as ConnectionInfo; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static ConnectionInfo create() => ConnectionInfo._();
  ConnectionInfo createEmptyInstance() => create();
  static $pb.PbList<ConnectionInfo> createRepeated() => $pb.PbList<ConnectionInfo>();
  @$core.pragma('dart2js:noInline')
  static ConnectionInfo getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ConnectionInfo>(create);
  static ConnectionInfo? _defaultInstance;

  @$pb.TagNumber(1)
  $fixnum.Int64 get connections => $_getI64(0);
  @$pb.TagNumber(1)
  set connections($fixnum.Int64 v) { $_setInt64(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasConnections() => $_has(0);
  @$pb.TagNumber(1)
  void clearConnections() => clearField(1);

  @$pb.TagNumber(2)
  $fixnum.Int64 get inboundConnections => $_getI64(1);
  @$pb.TagNumber(2)
  set inboundConnections($fixnum.Int64 v) { $_setInt64(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasInboundConnections() => $_has(1);
  @$pb.TagNumber(2)
  void clearInboundConnections() => clearField(2);

  @$pb.TagNumber(3)
  $fixnum.Int64 get outboundConnections => $_getI64(2);
  @$pb.TagNumber(3)
  set outboundConnections($fixnum.Int64 v) { $_setInt64(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasOutboundConnections() => $_has(2);
  @$pb.TagNumber(3)
  void clearOutboundConnections() => clearField(3);
}

class MetricInfo extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'MetricInfo', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOM<CounterInfo>(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'TotalInvalid', protoName: 'TotalInvalid', subBuilder: CounterInfo.create)
    ..aOM<CounterInfo>(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'TotalSent', protoName: 'TotalSent', subBuilder: CounterInfo.create)
    ..aOM<CounterInfo>(3, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'TotalReceived', protoName: 'TotalReceived', subBuilder: CounterInfo.create)
    ..m<$core.int, CounterInfo>(4, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'MessageSent', protoName: 'MessageSent', entryClassName: 'MetricInfo.MessageSentEntry', keyFieldType: $pb.PbFieldType.O3, valueFieldType: $pb.PbFieldType.OM, valueCreator: CounterInfo.create, packageName: const $pb.PackageName('pactus'))
    ..m<$core.int, CounterInfo>(5, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'MessageReceived', protoName: 'MessageReceived', entryClassName: 'MetricInfo.MessageReceivedEntry', keyFieldType: $pb.PbFieldType.O3, valueFieldType: $pb.PbFieldType.OM, valueCreator: CounterInfo.create, packageName: const $pb.PackageName('pactus'))
    ..hasRequiredFields = false
  ;

  MetricInfo._() : super();
  factory MetricInfo({
    CounterInfo? totalInvalid,
    CounterInfo? totalSent,
    CounterInfo? totalReceived,
    $core.Map<$core.int, CounterInfo>? messageSent,
    $core.Map<$core.int, CounterInfo>? messageReceived,
  }) {
    final _result = create();
    if (totalInvalid != null) {
      _result.totalInvalid = totalInvalid;
    }
    if (totalSent != null) {
      _result.totalSent = totalSent;
    }
    if (totalReceived != null) {
      _result.totalReceived = totalReceived;
    }
    if (messageSent != null) {
      _result.messageSent.addAll(messageSent);
    }
    if (messageReceived != null) {
      _result.messageReceived.addAll(messageReceived);
    }
    return _result;
  }
  factory MetricInfo.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory MetricInfo.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  MetricInfo clone() => MetricInfo()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  MetricInfo copyWith(void Function(MetricInfo) updates) => super.copyWith((message) => updates(message as MetricInfo)) as MetricInfo; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static MetricInfo create() => MetricInfo._();
  MetricInfo createEmptyInstance() => create();
  static $pb.PbList<MetricInfo> createRepeated() => $pb.PbList<MetricInfo>();
  @$core.pragma('dart2js:noInline')
  static MetricInfo getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<MetricInfo>(create);
  static MetricInfo? _defaultInstance;

  @$pb.TagNumber(1)
  CounterInfo get totalInvalid => $_getN(0);
  @$pb.TagNumber(1)
  set totalInvalid(CounterInfo v) { setField(1, v); }
  @$pb.TagNumber(1)
  $core.bool hasTotalInvalid() => $_has(0);
  @$pb.TagNumber(1)
  void clearTotalInvalid() => clearField(1);
  @$pb.TagNumber(1)
  CounterInfo ensureTotalInvalid() => $_ensure(0);

  @$pb.TagNumber(2)
  CounterInfo get totalSent => $_getN(1);
  @$pb.TagNumber(2)
  set totalSent(CounterInfo v) { setField(2, v); }
  @$pb.TagNumber(2)
  $core.bool hasTotalSent() => $_has(1);
  @$pb.TagNumber(2)
  void clearTotalSent() => clearField(2);
  @$pb.TagNumber(2)
  CounterInfo ensureTotalSent() => $_ensure(1);

  @$pb.TagNumber(3)
  CounterInfo get totalReceived => $_getN(2);
  @$pb.TagNumber(3)
  set totalReceived(CounterInfo v) { setField(3, v); }
  @$pb.TagNumber(3)
  $core.bool hasTotalReceived() => $_has(2);
  @$pb.TagNumber(3)
  void clearTotalReceived() => clearField(3);
  @$pb.TagNumber(3)
  CounterInfo ensureTotalReceived() => $_ensure(2);

  @$pb.TagNumber(4)
  $core.Map<$core.int, CounterInfo> get messageSent => $_getMap(3);

  @$pb.TagNumber(5)
  $core.Map<$core.int, CounterInfo> get messageReceived => $_getMap(4);
}

class CounterInfo extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'CounterInfo', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..a<$fixnum.Int64>(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'Bytes', $pb.PbFieldType.OU6, protoName: 'Bytes', defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'Bundles', $pb.PbFieldType.OU6, protoName: 'Bundles', defaultOrMaker: $fixnum.Int64.ZERO)
    ..hasRequiredFields = false
  ;

  CounterInfo._() : super();
  factory CounterInfo({
    $fixnum.Int64? bytes,
    $fixnum.Int64? bundles,
  }) {
    final _result = create();
    if (bytes != null) {
      _result.bytes = bytes;
    }
    if (bundles != null) {
      _result.bundles = bundles;
    }
    return _result;
  }
  factory CounterInfo.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory CounterInfo.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  CounterInfo clone() => CounterInfo()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  CounterInfo copyWith(void Function(CounterInfo) updates) => super.copyWith((message) => updates(message as CounterInfo)) as CounterInfo; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static CounterInfo create() => CounterInfo._();
  CounterInfo createEmptyInstance() => create();
  static $pb.PbList<CounterInfo> createRepeated() => $pb.PbList<CounterInfo>();
  @$core.pragma('dart2js:noInline')
  static CounterInfo getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<CounterInfo>(create);
  static CounterInfo? _defaultInstance;

  @$pb.TagNumber(1)
  $fixnum.Int64 get bytes => $_getI64(0);
  @$pb.TagNumber(1)
  set bytes($fixnum.Int64 v) { $_setInt64(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasBytes() => $_has(0);
  @$pb.TagNumber(1)
  void clearBytes() => clearField(1);

  @$pb.TagNumber(2)
  $fixnum.Int64 get bundles => $_getI64(1);
  @$pb.TagNumber(2)
  set bundles($fixnum.Int64 v) { $_setInt64(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasBundles() => $_has(1);
  @$pb.TagNumber(2)
  void clearBundles() => clearField(2);
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

