//
//  Generated code. Do not modify.
//  source: network.proto
//
// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_final_fields
// ignore_for_file: unnecessary_import, unnecessary_this, unused_import

import 'dart:async' as $async;
import 'dart:core' as $core;

import 'package:fixnum/fixnum.dart' as $fixnum;
import 'package:protobuf/protobuf.dart' as $pb;

import 'network.pbenum.dart';

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

export 'network.pbenum.dart';

/// Request message for retrieving overall network information.
class GetNetworkInfoRequest extends $pb.GeneratedMessage {
  factory GetNetworkInfoRequest({
    $core.bool? onlyConnected,
  }) {
    final $result = create();
    if (onlyConnected != null) {
      $result.onlyConnected = onlyConnected;
    }
    return $result;
  }
  GetNetworkInfoRequest._() : super();
  factory GetNetworkInfoRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetNetworkInfoRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetNetworkInfoRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..aOB(1, _omitFieldNames ? '' : 'onlyConnected')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetNetworkInfoRequest clone() => GetNetworkInfoRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetNetworkInfoRequest copyWith(void Function(GetNetworkInfoRequest) updates) => super.copyWith((message) => updates(message as GetNetworkInfoRequest)) as GetNetworkInfoRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetNetworkInfoRequest create() => GetNetworkInfoRequest._();
  GetNetworkInfoRequest createEmptyInstance() => create();
  static $pb.PbList<GetNetworkInfoRequest> createRepeated() => $pb.PbList<GetNetworkInfoRequest>();
  @$core.pragma('dart2js:noInline')
  static GetNetworkInfoRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetNetworkInfoRequest>(create);
  static GetNetworkInfoRequest? _defaultInstance;

  /// If true, returns only peers that are currently connected.
  @$pb.TagNumber(1)
  $core.bool get onlyConnected => $_getBF(0);
  @$pb.TagNumber(1)
  set onlyConnected($core.bool v) { $_setBool(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasOnlyConnected() => $_has(0);
  @$pb.TagNumber(1)
  void clearOnlyConnected() => $_clearField(1);
}

/// Response message contains information about the overall network.
class GetNetworkInfoResponse extends $pb.GeneratedMessage {
  factory GetNetworkInfoResponse({
    $core.String? networkName,
    $core.int? connectedPeersCount,
    $core.Iterable<PeerInfo>? connectedPeers,
    MetricInfo? metricInfo,
  }) {
    final $result = create();
    if (networkName != null) {
      $result.networkName = networkName;
    }
    if (connectedPeersCount != null) {
      $result.connectedPeersCount = connectedPeersCount;
    }
    if (connectedPeers != null) {
      $result.connectedPeers.addAll(connectedPeers);
    }
    if (metricInfo != null) {
      $result.metricInfo = metricInfo;
    }
    return $result;
  }
  GetNetworkInfoResponse._() : super();
  factory GetNetworkInfoResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetNetworkInfoResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetNetworkInfoResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'networkName')
    ..a<$core.int>(2, _omitFieldNames ? '' : 'connectedPeersCount', $pb.PbFieldType.OU3)
    ..pc<PeerInfo>(3, _omitFieldNames ? '' : 'connectedPeers', $pb.PbFieldType.PM, subBuilder: PeerInfo.create)
    ..aOM<MetricInfo>(4, _omitFieldNames ? '' : 'metricInfo', subBuilder: MetricInfo.create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetNetworkInfoResponse clone() => GetNetworkInfoResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetNetworkInfoResponse copyWith(void Function(GetNetworkInfoResponse) updates) => super.copyWith((message) => updates(message as GetNetworkInfoResponse)) as GetNetworkInfoResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetNetworkInfoResponse create() => GetNetworkInfoResponse._();
  GetNetworkInfoResponse createEmptyInstance() => create();
  static $pb.PbList<GetNetworkInfoResponse> createRepeated() => $pb.PbList<GetNetworkInfoResponse>();
  @$core.pragma('dart2js:noInline')
  static GetNetworkInfoResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetNetworkInfoResponse>(create);
  static GetNetworkInfoResponse? _defaultInstance;

  /// Name of the network.
  @$pb.TagNumber(1)
  $core.String get networkName => $_getSZ(0);
  @$pb.TagNumber(1)
  set networkName($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasNetworkName() => $_has(0);
  @$pb.TagNumber(1)
  void clearNetworkName() => $_clearField(1);

  /// Number of connected peers.
  @$pb.TagNumber(2)
  $core.int get connectedPeersCount => $_getIZ(1);
  @$pb.TagNumber(2)
  set connectedPeersCount($core.int v) { $_setUnsignedInt32(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasConnectedPeersCount() => $_has(1);
  @$pb.TagNumber(2)
  void clearConnectedPeersCount() => $_clearField(2);

  /// List of connected peers.
  @$pb.TagNumber(3)
  $pb.PbList<PeerInfo> get connectedPeers => $_getList(2);

  /// Metrics related to node activity.
  @$pb.TagNumber(4)
  MetricInfo get metricInfo => $_getN(3);
  @$pb.TagNumber(4)
  set metricInfo(MetricInfo v) { $_setField(4, v); }
  @$pb.TagNumber(4)
  $core.bool hasMetricInfo() => $_has(3);
  @$pb.TagNumber(4)
  void clearMetricInfo() => $_clearField(4);
  @$pb.TagNumber(4)
  MetricInfo ensureMetricInfo() => $_ensure(3);
}

/// Request message for retrieving information of the node.
class GetNodeInfoRequest extends $pb.GeneratedMessage {
  factory GetNodeInfoRequest() => create();
  GetNodeInfoRequest._() : super();
  factory GetNodeInfoRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetNodeInfoRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetNodeInfoRequest', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetNodeInfoRequest clone() => GetNodeInfoRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetNodeInfoRequest copyWith(void Function(GetNodeInfoRequest) updates) => super.copyWith((message) => updates(message as GetNodeInfoRequest)) as GetNodeInfoRequest;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetNodeInfoRequest create() => GetNodeInfoRequest._();
  GetNodeInfoRequest createEmptyInstance() => create();
  static $pb.PbList<GetNodeInfoRequest> createRepeated() => $pb.PbList<GetNodeInfoRequest>();
  @$core.pragma('dart2js:noInline')
  static GetNodeInfoRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetNodeInfoRequest>(create);
  static GetNodeInfoRequest? _defaultInstance;
}

/// Response message contains information about a specific node in the network.
class GetNodeInfoResponse extends $pb.GeneratedMessage {
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
    $fixnum.Int64? currentTime,
  }) {
    final $result = create();
    if (moniker != null) {
      $result.moniker = moniker;
    }
    if (agent != null) {
      $result.agent = agent;
    }
    if (peerId != null) {
      $result.peerId = peerId;
    }
    if (startedAt != null) {
      $result.startedAt = startedAt;
    }
    if (reachability != null) {
      $result.reachability = reachability;
    }
    if (services != null) {
      $result.services = services;
    }
    if (servicesNames != null) {
      $result.servicesNames = servicesNames;
    }
    if (localAddrs != null) {
      $result.localAddrs.addAll(localAddrs);
    }
    if (protocols != null) {
      $result.protocols.addAll(protocols);
    }
    if (clockOffset != null) {
      $result.clockOffset = clockOffset;
    }
    if (connectionInfo != null) {
      $result.connectionInfo = connectionInfo;
    }
    if (zmqPublishers != null) {
      $result.zmqPublishers.addAll(zmqPublishers);
    }
    if (currentTime != null) {
      $result.currentTime = currentTime;
    }
    return $result;
  }
  GetNodeInfoResponse._() : super();
  factory GetNodeInfoResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetNodeInfoResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'GetNodeInfoResponse', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'moniker')
    ..aOS(2, _omitFieldNames ? '' : 'agent')
    ..aOS(3, _omitFieldNames ? '' : 'peerId')
    ..a<$fixnum.Int64>(4, _omitFieldNames ? '' : 'startedAt', $pb.PbFieldType.OU6, defaultOrMaker: $fixnum.Int64.ZERO)
    ..aOS(5, _omitFieldNames ? '' : 'reachability')
    ..a<$core.int>(6, _omitFieldNames ? '' : 'services', $pb.PbFieldType.O3)
    ..aOS(7, _omitFieldNames ? '' : 'servicesNames')
    ..pPS(8, _omitFieldNames ? '' : 'localAddrs')
    ..pPS(9, _omitFieldNames ? '' : 'protocols')
    ..a<$core.double>(13, _omitFieldNames ? '' : 'clockOffset', $pb.PbFieldType.OD)
    ..aOM<ConnectionInfo>(14, _omitFieldNames ? '' : 'connectionInfo', subBuilder: ConnectionInfo.create)
    ..pc<ZMQPublisherInfo>(15, _omitFieldNames ? '' : 'zmqPublishers', $pb.PbFieldType.PM, subBuilder: ZMQPublisherInfo.create)
    ..a<$fixnum.Int64>(16, _omitFieldNames ? '' : 'currentTime', $pb.PbFieldType.OU6, defaultOrMaker: $fixnum.Int64.ZERO)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetNodeInfoResponse clone() => GetNodeInfoResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetNodeInfoResponse copyWith(void Function(GetNodeInfoResponse) updates) => super.copyWith((message) => updates(message as GetNodeInfoResponse)) as GetNodeInfoResponse;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetNodeInfoResponse create() => GetNodeInfoResponse._();
  GetNodeInfoResponse createEmptyInstance() => create();
  static $pb.PbList<GetNodeInfoResponse> createRepeated() => $pb.PbList<GetNodeInfoResponse>();
  @$core.pragma('dart2js:noInline')
  static GetNodeInfoResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetNodeInfoResponse>(create);
  static GetNodeInfoResponse? _defaultInstance;

  /// Moniker or Human-readable name identifying this node in the network.
  @$pb.TagNumber(1)
  $core.String get moniker => $_getSZ(0);
  @$pb.TagNumber(1)
  set moniker($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasMoniker() => $_has(0);
  @$pb.TagNumber(1)
  void clearMoniker() => $_clearField(1);

  /// Version and agent details of the node.
  @$pb.TagNumber(2)
  $core.String get agent => $_getSZ(1);
  @$pb.TagNumber(2)
  set agent($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasAgent() => $_has(1);
  @$pb.TagNumber(2)
  void clearAgent() => $_clearField(2);

  /// Peer ID of the node.
  @$pb.TagNumber(3)
  $core.String get peerId => $_getSZ(2);
  @$pb.TagNumber(3)
  set peerId($core.String v) { $_setString(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasPeerId() => $_has(2);
  @$pb.TagNumber(3)
  void clearPeerId() => $_clearField(3);

  /// Unix timestamp when the node was started (UTC).
  @$pb.TagNumber(4)
  $fixnum.Int64 get startedAt => $_getI64(3);
  @$pb.TagNumber(4)
  set startedAt($fixnum.Int64 v) { $_setInt64(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasStartedAt() => $_has(3);
  @$pb.TagNumber(4)
  void clearStartedAt() => $_clearField(4);

  /// Reachability status of the node.
  @$pb.TagNumber(5)
  $core.String get reachability => $_getSZ(4);
  @$pb.TagNumber(5)
  set reachability($core.String v) { $_setString(4, v); }
  @$pb.TagNumber(5)
  $core.bool hasReachability() => $_has(4);
  @$pb.TagNumber(5)
  void clearReachability() => $_clearField(5);

  /// Bitfield representing the services provided by the node.
  @$pb.TagNumber(6)
  $core.int get services => $_getIZ(5);
  @$pb.TagNumber(6)
  set services($core.int v) { $_setSignedInt32(5, v); }
  @$pb.TagNumber(6)
  $core.bool hasServices() => $_has(5);
  @$pb.TagNumber(6)
  void clearServices() => $_clearField(6);

  /// Names of services provided by the node.
  @$pb.TagNumber(7)
  $core.String get servicesNames => $_getSZ(6);
  @$pb.TagNumber(7)
  set servicesNames($core.String v) { $_setString(6, v); }
  @$pb.TagNumber(7)
  $core.bool hasServicesNames() => $_has(6);
  @$pb.TagNumber(7)
  void clearServicesNames() => $_clearField(7);

  /// List of addresses associated with the node.
  @$pb.TagNumber(8)
  $pb.PbList<$core.String> get localAddrs => $_getList(7);

  /// List of protocols supported by the node.
  @$pb.TagNumber(9)
  $pb.PbList<$core.String> get protocols => $_getList(8);

  /// Offset between the node's clock and the network's clock (in seconds).
  @$pb.TagNumber(13)
  $core.double get clockOffset => $_getN(9);
  @$pb.TagNumber(13)
  set clockOffset($core.double v) { $_setDouble(9, v); }
  @$pb.TagNumber(13)
  $core.bool hasClockOffset() => $_has(9);
  @$pb.TagNumber(13)
  void clearClockOffset() => $_clearField(13);

  /// Information about the node's connections.
  @$pb.TagNumber(14)
  ConnectionInfo get connectionInfo => $_getN(10);
  @$pb.TagNumber(14)
  set connectionInfo(ConnectionInfo v) { $_setField(14, v); }
  @$pb.TagNumber(14)
  $core.bool hasConnectionInfo() => $_has(10);
  @$pb.TagNumber(14)
  void clearConnectionInfo() => $_clearField(14);
  @$pb.TagNumber(14)
  ConnectionInfo ensureConnectionInfo() => $_ensure(10);

  /// List of active ZeroMQ publishers.
  @$pb.TagNumber(15)
  $pb.PbList<ZMQPublisherInfo> get zmqPublishers => $_getList(11);

  /// Current Unix timestamp of the node (UTC).
  @$pb.TagNumber(16)
  $fixnum.Int64 get currentTime => $_getI64(12);
  @$pb.TagNumber(16)
  set currentTime($fixnum.Int64 v) { $_setInt64(12, v); }
  @$pb.TagNumber(16)
  $core.bool hasCurrentTime() => $_has(12);
  @$pb.TagNumber(16)
  void clearCurrentTime() => $_clearField(16);
}

/// ZMQPublisherInfo contains information about a ZeroMQ publisher.
class ZMQPublisherInfo extends $pb.GeneratedMessage {
  factory ZMQPublisherInfo({
    $core.String? topic,
    $core.String? address,
    $core.int? hwm,
  }) {
    final $result = create();
    if (topic != null) {
      $result.topic = topic;
    }
    if (address != null) {
      $result.address = address;
    }
    if (hwm != null) {
      $result.hwm = hwm;
    }
    return $result;
  }
  ZMQPublisherInfo._() : super();
  factory ZMQPublisherInfo.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory ZMQPublisherInfo.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'ZMQPublisherInfo', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'topic')
    ..aOS(2, _omitFieldNames ? '' : 'address')
    ..a<$core.int>(3, _omitFieldNames ? '' : 'hwm', $pb.PbFieldType.O3)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  ZMQPublisherInfo clone() => ZMQPublisherInfo()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  ZMQPublisherInfo copyWith(void Function(ZMQPublisherInfo) updates) => super.copyWith((message) => updates(message as ZMQPublisherInfo)) as ZMQPublisherInfo;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ZMQPublisherInfo create() => ZMQPublisherInfo._();
  ZMQPublisherInfo createEmptyInstance() => create();
  static $pb.PbList<ZMQPublisherInfo> createRepeated() => $pb.PbList<ZMQPublisherInfo>();
  @$core.pragma('dart2js:noInline')
  static ZMQPublisherInfo getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ZMQPublisherInfo>(create);
  static ZMQPublisherInfo? _defaultInstance;

  /// The topic associated with the publisher.
  @$pb.TagNumber(1)
  $core.String get topic => $_getSZ(0);
  @$pb.TagNumber(1)
  set topic($core.String v) { $_setString(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasTopic() => $_has(0);
  @$pb.TagNumber(1)
  void clearTopic() => $_clearField(1);

  /// The address of the publisher.
  @$pb.TagNumber(2)
  $core.String get address => $_getSZ(1);
  @$pb.TagNumber(2)
  set address($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasAddress() => $_has(1);
  @$pb.TagNumber(2)
  void clearAddress() => $_clearField(2);

  /// The high-water mark (HWM) for the publisher, indicating the
  /// maximum number of messages to queue before dropping older ones.
  @$pb.TagNumber(3)
  $core.int get hwm => $_getIZ(2);
  @$pb.TagNumber(3)
  set hwm($core.int v) { $_setSignedInt32(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasHwm() => $_has(2);
  @$pb.TagNumber(3)
  void clearHwm() => $_clearField(3);
}

/// PeerInfo contains information about a peer in the network.
class PeerInfo extends $pb.GeneratedMessage {
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
    Direction? direction,
    $core.Iterable<$core.String>? protocols,
    $core.int? totalSessions,
    $core.int? completedSessions,
    MetricInfo? metricInfo,
    $core.bool? outboundHelloSent,
  }) {
    final $result = create();
    if (status != null) {
      $result.status = status;
    }
    if (moniker != null) {
      $result.moniker = moniker;
    }
    if (agent != null) {
      $result.agent = agent;
    }
    if (peerId != null) {
      $result.peerId = peerId;
    }
    if (consensusKeys != null) {
      $result.consensusKeys.addAll(consensusKeys);
    }
    if (consensusAddresses != null) {
      $result.consensusAddresses.addAll(consensusAddresses);
    }
    if (services != null) {
      $result.services = services;
    }
    if (lastBlockHash != null) {
      $result.lastBlockHash = lastBlockHash;
    }
    if (height != null) {
      $result.height = height;
    }
    if (lastSent != null) {
      $result.lastSent = lastSent;
    }
    if (lastReceived != null) {
      $result.lastReceived = lastReceived;
    }
    if (address != null) {
      $result.address = address;
    }
    if (direction != null) {
      $result.direction = direction;
    }
    if (protocols != null) {
      $result.protocols.addAll(protocols);
    }
    if (totalSessions != null) {
      $result.totalSessions = totalSessions;
    }
    if (completedSessions != null) {
      $result.completedSessions = completedSessions;
    }
    if (metricInfo != null) {
      $result.metricInfo = metricInfo;
    }
    if (outboundHelloSent != null) {
      $result.outboundHelloSent = outboundHelloSent;
    }
    return $result;
  }
  PeerInfo._() : super();
  factory PeerInfo.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory PeerInfo.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'PeerInfo', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..a<$core.int>(1, _omitFieldNames ? '' : 'status', $pb.PbFieldType.O3)
    ..aOS(2, _omitFieldNames ? '' : 'moniker')
    ..aOS(3, _omitFieldNames ? '' : 'agent')
    ..aOS(4, _omitFieldNames ? '' : 'peerId')
    ..pPS(5, _omitFieldNames ? '' : 'consensusKeys')
    ..pPS(6, _omitFieldNames ? '' : 'consensusAddresses')
    ..a<$core.int>(7, _omitFieldNames ? '' : 'services', $pb.PbFieldType.OU3)
    ..aOS(8, _omitFieldNames ? '' : 'lastBlockHash')
    ..a<$core.int>(9, _omitFieldNames ? '' : 'height', $pb.PbFieldType.OU3)
    ..aInt64(10, _omitFieldNames ? '' : 'lastSent')
    ..aInt64(11, _omitFieldNames ? '' : 'lastReceived')
    ..aOS(12, _omitFieldNames ? '' : 'address')
    ..e<Direction>(13, _omitFieldNames ? '' : 'direction', $pb.PbFieldType.OE, defaultOrMaker: Direction.DIRECTION_UNKNOWN, valueOf: Direction.valueOf, enumValues: Direction.values)
    ..pPS(14, _omitFieldNames ? '' : 'protocols')
    ..a<$core.int>(15, _omitFieldNames ? '' : 'totalSessions', $pb.PbFieldType.O3)
    ..a<$core.int>(16, _omitFieldNames ? '' : 'completedSessions', $pb.PbFieldType.O3)
    ..aOM<MetricInfo>(17, _omitFieldNames ? '' : 'metricInfo', subBuilder: MetricInfo.create)
    ..aOB(18, _omitFieldNames ? '' : 'outboundHelloSent')
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  PeerInfo clone() => PeerInfo()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  PeerInfo copyWith(void Function(PeerInfo) updates) => super.copyWith((message) => updates(message as PeerInfo)) as PeerInfo;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PeerInfo create() => PeerInfo._();
  PeerInfo createEmptyInstance() => create();
  static $pb.PbList<PeerInfo> createRepeated() => $pb.PbList<PeerInfo>();
  @$core.pragma('dart2js:noInline')
  static PeerInfo getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<PeerInfo>(create);
  static PeerInfo? _defaultInstance;

  /// Current status of the peer (e.g., connected, disconnected).
  @$pb.TagNumber(1)
  $core.int get status => $_getIZ(0);
  @$pb.TagNumber(1)
  set status($core.int v) { $_setSignedInt32(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasStatus() => $_has(0);
  @$pb.TagNumber(1)
  void clearStatus() => $_clearField(1);

  /// Moniker or Human-Readable name of the peer.
  @$pb.TagNumber(2)
  $core.String get moniker => $_getSZ(1);
  @$pb.TagNumber(2)
  set moniker($core.String v) { $_setString(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasMoniker() => $_has(1);
  @$pb.TagNumber(2)
  void clearMoniker() => $_clearField(2);

  /// Version and agent details of the peer.
  @$pb.TagNumber(3)
  $core.String get agent => $_getSZ(2);
  @$pb.TagNumber(3)
  set agent($core.String v) { $_setString(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasAgent() => $_has(2);
  @$pb.TagNumber(3)
  void clearAgent() => $_clearField(3);

  /// Peer ID of the peer in P2P network.
  @$pb.TagNumber(4)
  $core.String get peerId => $_getSZ(3);
  @$pb.TagNumber(4)
  set peerId($core.String v) { $_setString(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasPeerId() => $_has(3);
  @$pb.TagNumber(4)
  void clearPeerId() => $_clearField(4);

  /// List of consensus keys used by the peer.
  @$pb.TagNumber(5)
  $pb.PbList<$core.String> get consensusKeys => $_getList(4);

  /// List of consensus addresses used by the peer.
  @$pb.TagNumber(6)
  $pb.PbList<$core.String> get consensusAddresses => $_getList(5);

  /// Bitfield representing the services provided by the peer.
  @$pb.TagNumber(7)
  $core.int get services => $_getIZ(6);
  @$pb.TagNumber(7)
  set services($core.int v) { $_setUnsignedInt32(6, v); }
  @$pb.TagNumber(7)
  $core.bool hasServices() => $_has(6);
  @$pb.TagNumber(7)
  void clearServices() => $_clearField(7);

  /// Hash of the last block the peer knows.
  @$pb.TagNumber(8)
  $core.String get lastBlockHash => $_getSZ(7);
  @$pb.TagNumber(8)
  set lastBlockHash($core.String v) { $_setString(7, v); }
  @$pb.TagNumber(8)
  $core.bool hasLastBlockHash() => $_has(7);
  @$pb.TagNumber(8)
  void clearLastBlockHash() => $_clearField(8);

  /// Blockchain height of the peer.
  @$pb.TagNumber(9)
  $core.int get height => $_getIZ(8);
  @$pb.TagNumber(9)
  set height($core.int v) { $_setUnsignedInt32(8, v); }
  @$pb.TagNumber(9)
  $core.bool hasHeight() => $_has(8);
  @$pb.TagNumber(9)
  void clearHeight() => $_clearField(9);

  /// Unix timestamp of the last bundle sent to the peer (UTC).
  @$pb.TagNumber(10)
  $fixnum.Int64 get lastSent => $_getI64(9);
  @$pb.TagNumber(10)
  set lastSent($fixnum.Int64 v) { $_setInt64(9, v); }
  @$pb.TagNumber(10)
  $core.bool hasLastSent() => $_has(9);
  @$pb.TagNumber(10)
  void clearLastSent() => $_clearField(10);

  /// Unix timestamp of the last bundle received from the peer (UTC).
  @$pb.TagNumber(11)
  $fixnum.Int64 get lastReceived => $_getI64(10);
  @$pb.TagNumber(11)
  set lastReceived($fixnum.Int64 v) { $_setInt64(10, v); }
  @$pb.TagNumber(11)
  $core.bool hasLastReceived() => $_has(10);
  @$pb.TagNumber(11)
  void clearLastReceived() => $_clearField(11);

  /// Network address of the peer.
  @$pb.TagNumber(12)
  $core.String get address => $_getSZ(11);
  @$pb.TagNumber(12)
  set address($core.String v) { $_setString(11, v); }
  @$pb.TagNumber(12)
  $core.bool hasAddress() => $_has(11);
  @$pb.TagNumber(12)
  void clearAddress() => $_clearField(12);

  /// Connection direction (e.g., inbound, outbound).
  @$pb.TagNumber(13)
  Direction get direction => $_getN(12);
  @$pb.TagNumber(13)
  set direction(Direction v) { $_setField(13, v); }
  @$pb.TagNumber(13)
  $core.bool hasDirection() => $_has(12);
  @$pb.TagNumber(13)
  void clearDirection() => $_clearField(13);

  /// List of protocols supported by the peer.
  @$pb.TagNumber(14)
  $pb.PbList<$core.String> get protocols => $_getList(13);

  /// Total download sessions with the peer.
  @$pb.TagNumber(15)
  $core.int get totalSessions => $_getIZ(14);
  @$pb.TagNumber(15)
  set totalSessions($core.int v) { $_setSignedInt32(14, v); }
  @$pb.TagNumber(15)
  $core.bool hasTotalSessions() => $_has(14);
  @$pb.TagNumber(15)
  void clearTotalSessions() => $_clearField(15);

  /// Completed download sessions with the peer.
  @$pb.TagNumber(16)
  $core.int get completedSessions => $_getIZ(15);
  @$pb.TagNumber(16)
  set completedSessions($core.int v) { $_setSignedInt32(15, v); }
  @$pb.TagNumber(16)
  $core.bool hasCompletedSessions() => $_has(15);
  @$pb.TagNumber(16)
  void clearCompletedSessions() => $_clearField(16);

  /// Metrics related to peer activity.
  @$pb.TagNumber(17)
  MetricInfo get metricInfo => $_getN(16);
  @$pb.TagNumber(17)
  set metricInfo(MetricInfo v) { $_setField(17, v); }
  @$pb.TagNumber(17)
  $core.bool hasMetricInfo() => $_has(16);
  @$pb.TagNumber(17)
  void clearMetricInfo() => $_clearField(17);
  @$pb.TagNumber(17)
  MetricInfo ensureMetricInfo() => $_ensure(16);

  /// Whether we've sent the hello message for outbound connections.
  @$pb.TagNumber(18)
  $core.bool get outboundHelloSent => $_getBF(17);
  @$pb.TagNumber(18)
  set outboundHelloSent($core.bool v) { $_setBool(17, v); }
  @$pb.TagNumber(18)
  $core.bool hasOutboundHelloSent() => $_has(17);
  @$pb.TagNumber(18)
  void clearOutboundHelloSent() => $_clearField(18);
}

/// ConnectionInfo contains information about the node's connections.
class ConnectionInfo extends $pb.GeneratedMessage {
  factory ConnectionInfo({
    $fixnum.Int64? connections,
    $fixnum.Int64? inboundConnections,
    $fixnum.Int64? outboundConnections,
  }) {
    final $result = create();
    if (connections != null) {
      $result.connections = connections;
    }
    if (inboundConnections != null) {
      $result.inboundConnections = inboundConnections;
    }
    if (outboundConnections != null) {
      $result.outboundConnections = outboundConnections;
    }
    return $result;
  }
  ConnectionInfo._() : super();
  factory ConnectionInfo.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory ConnectionInfo.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'ConnectionInfo', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..a<$fixnum.Int64>(1, _omitFieldNames ? '' : 'connections', $pb.PbFieldType.OU6, defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(2, _omitFieldNames ? '' : 'inboundConnections', $pb.PbFieldType.OU6, defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(3, _omitFieldNames ? '' : 'outboundConnections', $pb.PbFieldType.OU6, defaultOrMaker: $fixnum.Int64.ZERO)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  ConnectionInfo clone() => ConnectionInfo()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  ConnectionInfo copyWith(void Function(ConnectionInfo) updates) => super.copyWith((message) => updates(message as ConnectionInfo)) as ConnectionInfo;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ConnectionInfo create() => ConnectionInfo._();
  ConnectionInfo createEmptyInstance() => create();
  static $pb.PbList<ConnectionInfo> createRepeated() => $pb.PbList<ConnectionInfo>();
  @$core.pragma('dart2js:noInline')
  static ConnectionInfo getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<ConnectionInfo>(create);
  static ConnectionInfo? _defaultInstance;

  /// Total number of connections.
  @$pb.TagNumber(1)
  $fixnum.Int64 get connections => $_getI64(0);
  @$pb.TagNumber(1)
  set connections($fixnum.Int64 v) { $_setInt64(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasConnections() => $_has(0);
  @$pb.TagNumber(1)
  void clearConnections() => $_clearField(1);

  /// Number of inbound connections.
  @$pb.TagNumber(2)
  $fixnum.Int64 get inboundConnections => $_getI64(1);
  @$pb.TagNumber(2)
  set inboundConnections($fixnum.Int64 v) { $_setInt64(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasInboundConnections() => $_has(1);
  @$pb.TagNumber(2)
  void clearInboundConnections() => $_clearField(2);

  /// Number of outbound connections.
  @$pb.TagNumber(3)
  $fixnum.Int64 get outboundConnections => $_getI64(2);
  @$pb.TagNumber(3)
  set outboundConnections($fixnum.Int64 v) { $_setInt64(2, v); }
  @$pb.TagNumber(3)
  $core.bool hasOutboundConnections() => $_has(2);
  @$pb.TagNumber(3)
  void clearOutboundConnections() => $_clearField(3);
}

/// MetricInfo contains metrics data regarding network activity.
class MetricInfo extends $pb.GeneratedMessage {
  factory MetricInfo({
    CounterInfo? totalInvalid,
    CounterInfo? totalSent,
    CounterInfo? totalReceived,
    $pb.PbMap<$core.int, CounterInfo>? messageSent,
    $pb.PbMap<$core.int, CounterInfo>? messageReceived,
  }) {
    final $result = create();
    if (totalInvalid != null) {
      $result.totalInvalid = totalInvalid;
    }
    if (totalSent != null) {
      $result.totalSent = totalSent;
    }
    if (totalReceived != null) {
      $result.totalReceived = totalReceived;
    }
    if (messageSent != null) {
      $result.messageSent.addAll(messageSent);
    }
    if (messageReceived != null) {
      $result.messageReceived.addAll(messageReceived);
    }
    return $result;
  }
  MetricInfo._() : super();
  factory MetricInfo.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory MetricInfo.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'MetricInfo', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..aOM<CounterInfo>(1, _omitFieldNames ? '' : 'totalInvalid', subBuilder: CounterInfo.create)
    ..aOM<CounterInfo>(2, _omitFieldNames ? '' : 'totalSent', subBuilder: CounterInfo.create)
    ..aOM<CounterInfo>(3, _omitFieldNames ? '' : 'totalReceived', subBuilder: CounterInfo.create)
    ..m<$core.int, CounterInfo>(4, _omitFieldNames ? '' : 'messageSent', entryClassName: 'MetricInfo.MessageSentEntry', keyFieldType: $pb.PbFieldType.O3, valueFieldType: $pb.PbFieldType.OM, valueCreator: CounterInfo.create, valueDefaultOrMaker: CounterInfo.getDefault, packageName: const $pb.PackageName('pactus'))
    ..m<$core.int, CounterInfo>(5, _omitFieldNames ? '' : 'messageReceived', entryClassName: 'MetricInfo.MessageReceivedEntry', keyFieldType: $pb.PbFieldType.O3, valueFieldType: $pb.PbFieldType.OM, valueCreator: CounterInfo.create, valueDefaultOrMaker: CounterInfo.getDefault, packageName: const $pb.PackageName('pactus'))
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  MetricInfo clone() => MetricInfo()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  MetricInfo copyWith(void Function(MetricInfo) updates) => super.copyWith((message) => updates(message as MetricInfo)) as MetricInfo;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MetricInfo create() => MetricInfo._();
  MetricInfo createEmptyInstance() => create();
  static $pb.PbList<MetricInfo> createRepeated() => $pb.PbList<MetricInfo>();
  @$core.pragma('dart2js:noInline')
  static MetricInfo getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<MetricInfo>(create);
  static MetricInfo? _defaultInstance;

  /// Total number of invalid bundles.
  @$pb.TagNumber(1)
  CounterInfo get totalInvalid => $_getN(0);
  @$pb.TagNumber(1)
  set totalInvalid(CounterInfo v) { $_setField(1, v); }
  @$pb.TagNumber(1)
  $core.bool hasTotalInvalid() => $_has(0);
  @$pb.TagNumber(1)
  void clearTotalInvalid() => $_clearField(1);
  @$pb.TagNumber(1)
  CounterInfo ensureTotalInvalid() => $_ensure(0);

  /// Total number of bundles sent.
  @$pb.TagNumber(2)
  CounterInfo get totalSent => $_getN(1);
  @$pb.TagNumber(2)
  set totalSent(CounterInfo v) { $_setField(2, v); }
  @$pb.TagNumber(2)
  $core.bool hasTotalSent() => $_has(1);
  @$pb.TagNumber(2)
  void clearTotalSent() => $_clearField(2);
  @$pb.TagNumber(2)
  CounterInfo ensureTotalSent() => $_ensure(1);

  /// Total number of bundles received.
  @$pb.TagNumber(3)
  CounterInfo get totalReceived => $_getN(2);
  @$pb.TagNumber(3)
  set totalReceived(CounterInfo v) { $_setField(3, v); }
  @$pb.TagNumber(3)
  $core.bool hasTotalReceived() => $_has(2);
  @$pb.TagNumber(3)
  void clearTotalReceived() => $_clearField(3);
  @$pb.TagNumber(3)
  CounterInfo ensureTotalReceived() => $_ensure(2);

  /// Number of sent bundles categorized by message type.
  @$pb.TagNumber(4)
  $pb.PbMap<$core.int, CounterInfo> get messageSent => $_getMap(3);

  /// Number of received bundles categorized by message type.
  @$pb.TagNumber(5)
  $pb.PbMap<$core.int, CounterInfo> get messageReceived => $_getMap(4);
}

/// CounterInfo holds counter data regarding byte and bundle counts.
class CounterInfo extends $pb.GeneratedMessage {
  factory CounterInfo({
    $fixnum.Int64? bytes,
    $fixnum.Int64? bundles,
  }) {
    final $result = create();
    if (bytes != null) {
      $result.bytes = bytes;
    }
    if (bundles != null) {
      $result.bundles = bundles;
    }
    return $result;
  }
  CounterInfo._() : super();
  factory CounterInfo.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory CounterInfo.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(_omitMessageNames ? '' : 'CounterInfo', package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'), createEmptyInstance: create)
    ..a<$fixnum.Int64>(1, _omitFieldNames ? '' : 'bytes', $pb.PbFieldType.OU6, defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(2, _omitFieldNames ? '' : 'bundles', $pb.PbFieldType.OU6, defaultOrMaker: $fixnum.Int64.ZERO)
    ..hasRequiredFields = false
  ;

  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  CounterInfo clone() => CounterInfo()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  CounterInfo copyWith(void Function(CounterInfo) updates) => super.copyWith((message) => updates(message as CounterInfo)) as CounterInfo;

  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CounterInfo create() => CounterInfo._();
  CounterInfo createEmptyInstance() => create();
  static $pb.PbList<CounterInfo> createRepeated() => $pb.PbList<CounterInfo>();
  @$core.pragma('dart2js:noInline')
  static CounterInfo getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<CounterInfo>(create);
  static CounterInfo? _defaultInstance;

  /// Total number of bytes.
  @$pb.TagNumber(1)
  $fixnum.Int64 get bytes => $_getI64(0);
  @$pb.TagNumber(1)
  set bytes($fixnum.Int64 v) { $_setInt64(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasBytes() => $_has(0);
  @$pb.TagNumber(1)
  void clearBytes() => $_clearField(1);

  /// Total number of bundles.
  @$pb.TagNumber(2)
  $fixnum.Int64 get bundles => $_getI64(1);
  @$pb.TagNumber(2)
  set bundles($fixnum.Int64 v) { $_setInt64(1, v); }
  @$pb.TagNumber(2)
  $core.bool hasBundles() => $_has(1);
  @$pb.TagNumber(2)
  void clearBundles() => $_clearField(2);
}

/// Network service provides RPCs for retrieving information about the network.
class NetworkApi {
  $pb.RpcClient _client;
  NetworkApi(this._client);

  /// GetNetworkInfo retrieves information about the overall network.
  $async.Future<GetNetworkInfoResponse> getNetworkInfo($pb.ClientContext? ctx, GetNetworkInfoRequest request) =>
    _client.invoke<GetNetworkInfoResponse>(ctx, 'Network', 'GetNetworkInfo', request, GetNetworkInfoResponse())
  ;
  /// GetNodeInfo retrieves information about a specific node in the network.
  $async.Future<GetNodeInfoResponse> getNodeInfo($pb.ClientContext? ctx, GetNodeInfoRequest request) =>
    _client.invoke<GetNodeInfoResponse>(ctx, 'Network', 'GetNodeInfo', request, GetNodeInfoResponse())
  ;
}


const _omitFieldNames = $core.bool.fromEnvironment('protobuf.omit_field_names');
const _omitMessageNames = $core.bool.fromEnvironment('protobuf.omit_message_names');
