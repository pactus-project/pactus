// This is a generated file - do not edit.
//
// Generated from network.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:async' as $async;
import 'dart:core' as $core;

import 'package:fixnum/fixnum.dart' as $fixnum;
import 'package:protobuf/protobuf.dart' as $pb;

import 'network.pbenum.dart';

export 'package:protobuf/protobuf.dart' show GeneratedMessageGenericExtensions;

export 'network.pbenum.dart';

/// Request message for retrieving overall network information.
class GetNetworkInfoRequest extends $pb.GeneratedMessage {
  factory GetNetworkInfoRequest() => create();

  GetNetworkInfoRequest._();

  factory GetNetworkInfoRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetNetworkInfoRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetNetworkInfoRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetNetworkInfoRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetNetworkInfoRequest copyWith(
          void Function(GetNetworkInfoRequest) updates) =>
      super.copyWith((message) => updates(message as GetNetworkInfoRequest))
          as GetNetworkInfoRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetNetworkInfoRequest create() => GetNetworkInfoRequest._();
  @$core.override
  GetNetworkInfoRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetNetworkInfoRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetNetworkInfoRequest>(create);
  static GetNetworkInfoRequest? _defaultInstance;
}

/// Response message contains information about the overall network.
class GetNetworkInfoResponse extends $pb.GeneratedMessage {
  factory GetNetworkInfoResponse({
    $core.String? networkName,
    $core.int? connectedPeersCount,
    MetricInfo? metricInfo,
  }) {
    final result = create();
    if (networkName != null) result.networkName = networkName;
    if (connectedPeersCount != null)
      result.connectedPeersCount = connectedPeersCount;
    if (metricInfo != null) result.metricInfo = metricInfo;
    return result;
  }

  GetNetworkInfoResponse._();

  factory GetNetworkInfoResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetNetworkInfoResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetNetworkInfoResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'networkName')
    ..aI(2, _omitFieldNames ? '' : 'connectedPeersCount',
        fieldType: $pb.PbFieldType.OU3)
    ..aOM<MetricInfo>(4, _omitFieldNames ? '' : 'metricInfo',
        subBuilder: MetricInfo.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetNetworkInfoResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetNetworkInfoResponse copyWith(
          void Function(GetNetworkInfoResponse) updates) =>
      super.copyWith((message) => updates(message as GetNetworkInfoResponse))
          as GetNetworkInfoResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetNetworkInfoResponse create() => GetNetworkInfoResponse._();
  @$core.override
  GetNetworkInfoResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetNetworkInfoResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetNetworkInfoResponse>(create);
  static GetNetworkInfoResponse? _defaultInstance;

  /// Name of the network.
  @$pb.TagNumber(1)
  $core.String get networkName => $_getSZ(0);
  @$pb.TagNumber(1)
  set networkName($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasNetworkName() => $_has(0);
  @$pb.TagNumber(1)
  void clearNetworkName() => $_clearField(1);

  /// Number of connected peers.
  @$pb.TagNumber(2)
  $core.int get connectedPeersCount => $_getIZ(1);
  @$pb.TagNumber(2)
  set connectedPeersCount($core.int value) => $_setUnsignedInt32(1, value);
  @$pb.TagNumber(2)
  $core.bool hasConnectedPeersCount() => $_has(1);
  @$pb.TagNumber(2)
  void clearConnectedPeersCount() => $_clearField(2);

  /// Metrics related to node activity.
  @$pb.TagNumber(4)
  MetricInfo get metricInfo => $_getN(2);
  @$pb.TagNumber(4)
  set metricInfo(MetricInfo value) => $_setField(4, value);
  @$pb.TagNumber(4)
  $core.bool hasMetricInfo() => $_has(2);
  @$pb.TagNumber(4)
  void clearMetricInfo() => $_clearField(4);
  @$pb.TagNumber(4)
  MetricInfo ensureMetricInfo() => $_ensure(2);
}

/// Request message for listing peers.
class ListPeersRequest extends $pb.GeneratedMessage {
  factory ListPeersRequest({
    $core.bool? includeDisconnected,
  }) {
    final result = create();
    if (includeDisconnected != null)
      result.includeDisconnected = includeDisconnected;
    return result;
  }

  ListPeersRequest._();

  factory ListPeersRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListPeersRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListPeersRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOB(1, _omitFieldNames ? '' : 'includeDisconnected')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListPeersRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListPeersRequest copyWith(void Function(ListPeersRequest) updates) =>
      super.copyWith((message) => updates(message as ListPeersRequest))
          as ListPeersRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListPeersRequest create() => ListPeersRequest._();
  @$core.override
  ListPeersRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListPeersRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListPeersRequest>(create);
  static ListPeersRequest? _defaultInstance;

  /// If true, includes disconnected peers (default: connected peers only).
  @$pb.TagNumber(1)
  $core.bool get includeDisconnected => $_getBF(0);
  @$pb.TagNumber(1)
  set includeDisconnected($core.bool value) => $_setBool(0, value);
  @$pb.TagNumber(1)
  $core.bool hasIncludeDisconnected() => $_has(0);
  @$pb.TagNumber(1)
  void clearIncludeDisconnected() => $_clearField(1);
}

/// Response message for listing peers.
class ListPeersResponse extends $pb.GeneratedMessage {
  factory ListPeersResponse({
    $core.Iterable<PeerInfo>? peers,
  }) {
    final result = create();
    if (peers != null) result.peers.addAll(peers);
    return result;
  }

  ListPeersResponse._();

  factory ListPeersResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ListPeersResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ListPeersResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..pPM<PeerInfo>(1, _omitFieldNames ? '' : 'peers',
        subBuilder: PeerInfo.create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListPeersResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ListPeersResponse copyWith(void Function(ListPeersResponse) updates) =>
      super.copyWith((message) => updates(message as ListPeersResponse))
          as ListPeersResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ListPeersResponse create() => ListPeersResponse._();
  @$core.override
  ListPeersResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ListPeersResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ListPeersResponse>(create);
  static ListPeersResponse? _defaultInstance;

  /// List of peers.
  @$pb.TagNumber(1)
  $pb.PbList<PeerInfo> get peers => $_getList(0);
}

/// Request message for retrieving information of the node.
class GetNodeInfoRequest extends $pb.GeneratedMessage {
  factory GetNodeInfoRequest() => create();

  GetNodeInfoRequest._();

  factory GetNodeInfoRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetNodeInfoRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetNodeInfoRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetNodeInfoRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetNodeInfoRequest copyWith(void Function(GetNodeInfoRequest) updates) =>
      super.copyWith((message) => updates(message as GetNodeInfoRequest))
          as GetNodeInfoRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetNodeInfoRequest create() => GetNodeInfoRequest._();
  @$core.override
  GetNodeInfoRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetNodeInfoRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetNodeInfoRequest>(create);
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
    final result = create();
    if (moniker != null) result.moniker = moniker;
    if (agent != null) result.agent = agent;
    if (peerId != null) result.peerId = peerId;
    if (startedAt != null) result.startedAt = startedAt;
    if (reachability != null) result.reachability = reachability;
    if (services != null) result.services = services;
    if (servicesNames != null) result.servicesNames = servicesNames;
    if (localAddrs != null) result.localAddrs.addAll(localAddrs);
    if (protocols != null) result.protocols.addAll(protocols);
    if (clockOffset != null) result.clockOffset = clockOffset;
    if (connectionInfo != null) result.connectionInfo = connectionInfo;
    if (zmqPublishers != null) result.zmqPublishers.addAll(zmqPublishers);
    if (currentTime != null) result.currentTime = currentTime;
    return result;
  }

  GetNodeInfoResponse._();

  factory GetNodeInfoResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory GetNodeInfoResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'GetNodeInfoResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'moniker')
    ..aOS(2, _omitFieldNames ? '' : 'agent')
    ..aOS(3, _omitFieldNames ? '' : 'peerId')
    ..a<$fixnum.Int64>(
        4, _omitFieldNames ? '' : 'startedAt', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..aOS(5, _omitFieldNames ? '' : 'reachability')
    ..aI(6, _omitFieldNames ? '' : 'services')
    ..aOS(7, _omitFieldNames ? '' : 'servicesNames')
    ..pPS(8, _omitFieldNames ? '' : 'localAddrs')
    ..pPS(9, _omitFieldNames ? '' : 'protocols')
    ..aD(13, _omitFieldNames ? '' : 'clockOffset')
    ..aOM<ConnectionInfo>(14, _omitFieldNames ? '' : 'connectionInfo',
        subBuilder: ConnectionInfo.create)
    ..pPM<ZMQPublisherInfo>(15, _omitFieldNames ? '' : 'zmqPublishers',
        subBuilder: ZMQPublisherInfo.create)
    ..a<$fixnum.Int64>(
        16, _omitFieldNames ? '' : 'currentTime', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetNodeInfoResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  GetNodeInfoResponse copyWith(void Function(GetNodeInfoResponse) updates) =>
      super.copyWith((message) => updates(message as GetNodeInfoResponse))
          as GetNodeInfoResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static GetNodeInfoResponse create() => GetNodeInfoResponse._();
  @$core.override
  GetNodeInfoResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static GetNodeInfoResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<GetNodeInfoResponse>(create);
  static GetNodeInfoResponse? _defaultInstance;

  /// Moniker or Human-readable name identifying this node in the network.
  @$pb.TagNumber(1)
  $core.String get moniker => $_getSZ(0);
  @$pb.TagNumber(1)
  set moniker($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasMoniker() => $_has(0);
  @$pb.TagNumber(1)
  void clearMoniker() => $_clearField(1);

  /// Version and agent details of the node.
  @$pb.TagNumber(2)
  $core.String get agent => $_getSZ(1);
  @$pb.TagNumber(2)
  set agent($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasAgent() => $_has(1);
  @$pb.TagNumber(2)
  void clearAgent() => $_clearField(2);

  /// Peer ID of the node.
  @$pb.TagNumber(3)
  $core.String get peerId => $_getSZ(2);
  @$pb.TagNumber(3)
  set peerId($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasPeerId() => $_has(2);
  @$pb.TagNumber(3)
  void clearPeerId() => $_clearField(3);

  /// Unix timestamp when the node was started (UTC).
  @$pb.TagNumber(4)
  $fixnum.Int64 get startedAt => $_getI64(3);
  @$pb.TagNumber(4)
  set startedAt($fixnum.Int64 value) => $_setInt64(3, value);
  @$pb.TagNumber(4)
  $core.bool hasStartedAt() => $_has(3);
  @$pb.TagNumber(4)
  void clearStartedAt() => $_clearField(4);

  /// Reachability status of the node.
  @$pb.TagNumber(5)
  $core.String get reachability => $_getSZ(4);
  @$pb.TagNumber(5)
  set reachability($core.String value) => $_setString(4, value);
  @$pb.TagNumber(5)
  $core.bool hasReachability() => $_has(4);
  @$pb.TagNumber(5)
  void clearReachability() => $_clearField(5);

  /// Bitfield representing the services provided by the node.
  @$pb.TagNumber(6)
  $core.int get services => $_getIZ(5);
  @$pb.TagNumber(6)
  set services($core.int value) => $_setSignedInt32(5, value);
  @$pb.TagNumber(6)
  $core.bool hasServices() => $_has(5);
  @$pb.TagNumber(6)
  void clearServices() => $_clearField(6);

  /// Names of services provided by the node.
  @$pb.TagNumber(7)
  $core.String get servicesNames => $_getSZ(6);
  @$pb.TagNumber(7)
  set servicesNames($core.String value) => $_setString(6, value);
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
  set clockOffset($core.double value) => $_setDouble(9, value);
  @$pb.TagNumber(13)
  $core.bool hasClockOffset() => $_has(9);
  @$pb.TagNumber(13)
  void clearClockOffset() => $_clearField(13);

  /// Information about the node's connections.
  @$pb.TagNumber(14)
  ConnectionInfo get connectionInfo => $_getN(10);
  @$pb.TagNumber(14)
  set connectionInfo(ConnectionInfo value) => $_setField(14, value);
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
  set currentTime($fixnum.Int64 value) => $_setInt64(12, value);
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
    final result = create();
    if (topic != null) result.topic = topic;
    if (address != null) result.address = address;
    if (hwm != null) result.hwm = hwm;
    return result;
  }

  ZMQPublisherInfo._();

  factory ZMQPublisherInfo.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ZMQPublisherInfo.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ZMQPublisherInfo',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOS(1, _omitFieldNames ? '' : 'topic')
    ..aOS(2, _omitFieldNames ? '' : 'address')
    ..aI(3, _omitFieldNames ? '' : 'hwm')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ZMQPublisherInfo clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ZMQPublisherInfo copyWith(void Function(ZMQPublisherInfo) updates) =>
      super.copyWith((message) => updates(message as ZMQPublisherInfo))
          as ZMQPublisherInfo;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ZMQPublisherInfo create() => ZMQPublisherInfo._();
  @$core.override
  ZMQPublisherInfo createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ZMQPublisherInfo getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ZMQPublisherInfo>(create);
  static ZMQPublisherInfo? _defaultInstance;

  /// The topic associated with the publisher.
  @$pb.TagNumber(1)
  $core.String get topic => $_getSZ(0);
  @$pb.TagNumber(1)
  set topic($core.String value) => $_setString(0, value);
  @$pb.TagNumber(1)
  $core.bool hasTopic() => $_has(0);
  @$pb.TagNumber(1)
  void clearTopic() => $_clearField(1);

  /// The address of the publisher.
  @$pb.TagNumber(2)
  $core.String get address => $_getSZ(1);
  @$pb.TagNumber(2)
  set address($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasAddress() => $_has(1);
  @$pb.TagNumber(2)
  void clearAddress() => $_clearField(2);

  /// The high-water mark (HWM) for the publisher, indicating the
  /// maximum number of messages to queue before dropping older ones.
  @$pb.TagNumber(3)
  $core.int get hwm => $_getIZ(2);
  @$pb.TagNumber(3)
  set hwm($core.int value) => $_setSignedInt32(2, value);
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
    final result = create();
    if (status != null) result.status = status;
    if (moniker != null) result.moniker = moniker;
    if (agent != null) result.agent = agent;
    if (peerId != null) result.peerId = peerId;
    if (consensusKeys != null) result.consensusKeys.addAll(consensusKeys);
    if (consensusAddresses != null)
      result.consensusAddresses.addAll(consensusAddresses);
    if (services != null) result.services = services;
    if (lastBlockHash != null) result.lastBlockHash = lastBlockHash;
    if (height != null) result.height = height;
    if (lastSent != null) result.lastSent = lastSent;
    if (lastReceived != null) result.lastReceived = lastReceived;
    if (address != null) result.address = address;
    if (direction != null) result.direction = direction;
    if (protocols != null) result.protocols.addAll(protocols);
    if (totalSessions != null) result.totalSessions = totalSessions;
    if (completedSessions != null) result.completedSessions = completedSessions;
    if (metricInfo != null) result.metricInfo = metricInfo;
    if (outboundHelloSent != null) result.outboundHelloSent = outboundHelloSent;
    return result;
  }

  PeerInfo._();

  factory PeerInfo.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PeerInfo.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PeerInfo',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aI(1, _omitFieldNames ? '' : 'status')
    ..aOS(2, _omitFieldNames ? '' : 'moniker')
    ..aOS(3, _omitFieldNames ? '' : 'agent')
    ..aOS(4, _omitFieldNames ? '' : 'peerId')
    ..pPS(5, _omitFieldNames ? '' : 'consensusKeys')
    ..pPS(6, _omitFieldNames ? '' : 'consensusAddresses')
    ..aI(7, _omitFieldNames ? '' : 'services', fieldType: $pb.PbFieldType.OU3)
    ..aOS(8, _omitFieldNames ? '' : 'lastBlockHash')
    ..aI(9, _omitFieldNames ? '' : 'height', fieldType: $pb.PbFieldType.OU3)
    ..aInt64(10, _omitFieldNames ? '' : 'lastSent')
    ..aInt64(11, _omitFieldNames ? '' : 'lastReceived')
    ..aOS(12, _omitFieldNames ? '' : 'address')
    ..aE<Direction>(13, _omitFieldNames ? '' : 'direction',
        enumValues: Direction.values)
    ..pPS(14, _omitFieldNames ? '' : 'protocols')
    ..aI(15, _omitFieldNames ? '' : 'totalSessions')
    ..aI(16, _omitFieldNames ? '' : 'completedSessions')
    ..aOM<MetricInfo>(17, _omitFieldNames ? '' : 'metricInfo',
        subBuilder: MetricInfo.create)
    ..aOB(18, _omitFieldNames ? '' : 'outboundHelloSent')
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PeerInfo clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PeerInfo copyWith(void Function(PeerInfo) updates) =>
      super.copyWith((message) => updates(message as PeerInfo)) as PeerInfo;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PeerInfo create() => PeerInfo._();
  @$core.override
  PeerInfo createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PeerInfo getDefault() =>
      _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<PeerInfo>(create);
  static PeerInfo? _defaultInstance;

  /// Current status of the peer (e.g., connected, disconnected).
  @$pb.TagNumber(1)
  $core.int get status => $_getIZ(0);
  @$pb.TagNumber(1)
  set status($core.int value) => $_setSignedInt32(0, value);
  @$pb.TagNumber(1)
  $core.bool hasStatus() => $_has(0);
  @$pb.TagNumber(1)
  void clearStatus() => $_clearField(1);

  /// Moniker or Human-Readable name of the peer.
  @$pb.TagNumber(2)
  $core.String get moniker => $_getSZ(1);
  @$pb.TagNumber(2)
  set moniker($core.String value) => $_setString(1, value);
  @$pb.TagNumber(2)
  $core.bool hasMoniker() => $_has(1);
  @$pb.TagNumber(2)
  void clearMoniker() => $_clearField(2);

  /// Version and agent details of the peer.
  @$pb.TagNumber(3)
  $core.String get agent => $_getSZ(2);
  @$pb.TagNumber(3)
  set agent($core.String value) => $_setString(2, value);
  @$pb.TagNumber(3)
  $core.bool hasAgent() => $_has(2);
  @$pb.TagNumber(3)
  void clearAgent() => $_clearField(3);

  /// Peer ID of the peer in P2P network.
  @$pb.TagNumber(4)
  $core.String get peerId => $_getSZ(3);
  @$pb.TagNumber(4)
  set peerId($core.String value) => $_setString(3, value);
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
  set services($core.int value) => $_setUnsignedInt32(6, value);
  @$pb.TagNumber(7)
  $core.bool hasServices() => $_has(6);
  @$pb.TagNumber(7)
  void clearServices() => $_clearField(7);

  /// Hash of the last block the peer knows.
  @$pb.TagNumber(8)
  $core.String get lastBlockHash => $_getSZ(7);
  @$pb.TagNumber(8)
  set lastBlockHash($core.String value) => $_setString(7, value);
  @$pb.TagNumber(8)
  $core.bool hasLastBlockHash() => $_has(7);
  @$pb.TagNumber(8)
  void clearLastBlockHash() => $_clearField(8);

  /// Blockchain height of the peer.
  @$pb.TagNumber(9)
  $core.int get height => $_getIZ(8);
  @$pb.TagNumber(9)
  set height($core.int value) => $_setUnsignedInt32(8, value);
  @$pb.TagNumber(9)
  $core.bool hasHeight() => $_has(8);
  @$pb.TagNumber(9)
  void clearHeight() => $_clearField(9);

  /// Unix timestamp of the last bundle sent to the peer (UTC).
  @$pb.TagNumber(10)
  $fixnum.Int64 get lastSent => $_getI64(9);
  @$pb.TagNumber(10)
  set lastSent($fixnum.Int64 value) => $_setInt64(9, value);
  @$pb.TagNumber(10)
  $core.bool hasLastSent() => $_has(9);
  @$pb.TagNumber(10)
  void clearLastSent() => $_clearField(10);

  /// Unix timestamp of the last bundle received from the peer (UTC).
  @$pb.TagNumber(11)
  $fixnum.Int64 get lastReceived => $_getI64(10);
  @$pb.TagNumber(11)
  set lastReceived($fixnum.Int64 value) => $_setInt64(10, value);
  @$pb.TagNumber(11)
  $core.bool hasLastReceived() => $_has(10);
  @$pb.TagNumber(11)
  void clearLastReceived() => $_clearField(11);

  /// Network address of the peer.
  @$pb.TagNumber(12)
  $core.String get address => $_getSZ(11);
  @$pb.TagNumber(12)
  set address($core.String value) => $_setString(11, value);
  @$pb.TagNumber(12)
  $core.bool hasAddress() => $_has(11);
  @$pb.TagNumber(12)
  void clearAddress() => $_clearField(12);

  /// Connection direction (e.g., inbound, outbound).
  @$pb.TagNumber(13)
  Direction get direction => $_getN(12);
  @$pb.TagNumber(13)
  set direction(Direction value) => $_setField(13, value);
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
  set totalSessions($core.int value) => $_setSignedInt32(14, value);
  @$pb.TagNumber(15)
  $core.bool hasTotalSessions() => $_has(14);
  @$pb.TagNumber(15)
  void clearTotalSessions() => $_clearField(15);

  /// Completed download sessions with the peer.
  @$pb.TagNumber(16)
  $core.int get completedSessions => $_getIZ(15);
  @$pb.TagNumber(16)
  set completedSessions($core.int value) => $_setSignedInt32(15, value);
  @$pb.TagNumber(16)
  $core.bool hasCompletedSessions() => $_has(15);
  @$pb.TagNumber(16)
  void clearCompletedSessions() => $_clearField(16);

  /// Metrics related to peer activity.
  @$pb.TagNumber(17)
  MetricInfo get metricInfo => $_getN(16);
  @$pb.TagNumber(17)
  set metricInfo(MetricInfo value) => $_setField(17, value);
  @$pb.TagNumber(17)
  $core.bool hasMetricInfo() => $_has(16);
  @$pb.TagNumber(17)
  void clearMetricInfo() => $_clearField(17);
  @$pb.TagNumber(17)
  MetricInfo ensureMetricInfo() => $_ensure(16);

  /// Whether the hello message was sent from the outbound connection.
  @$pb.TagNumber(18)
  $core.bool get outboundHelloSent => $_getBF(17);
  @$pb.TagNumber(18)
  set outboundHelloSent($core.bool value) => $_setBool(17, value);
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
    final result = create();
    if (connections != null) result.connections = connections;
    if (inboundConnections != null)
      result.inboundConnections = inboundConnections;
    if (outboundConnections != null)
      result.outboundConnections = outboundConnections;
    return result;
  }

  ConnectionInfo._();

  factory ConnectionInfo.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory ConnectionInfo.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'ConnectionInfo',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..a<$fixnum.Int64>(
        1, _omitFieldNames ? '' : 'connections', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(
        2, _omitFieldNames ? '' : 'inboundConnections', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(
        3, _omitFieldNames ? '' : 'outboundConnections', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ConnectionInfo clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  ConnectionInfo copyWith(void Function(ConnectionInfo) updates) =>
      super.copyWith((message) => updates(message as ConnectionInfo))
          as ConnectionInfo;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static ConnectionInfo create() => ConnectionInfo._();
  @$core.override
  ConnectionInfo createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static ConnectionInfo getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<ConnectionInfo>(create);
  static ConnectionInfo? _defaultInstance;

  /// Total number of connections.
  @$pb.TagNumber(1)
  $fixnum.Int64 get connections => $_getI64(0);
  @$pb.TagNumber(1)
  set connections($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(1)
  $core.bool hasConnections() => $_has(0);
  @$pb.TagNumber(1)
  void clearConnections() => $_clearField(1);

  /// Number of inbound connections.
  @$pb.TagNumber(2)
  $fixnum.Int64 get inboundConnections => $_getI64(1);
  @$pb.TagNumber(2)
  set inboundConnections($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(2)
  $core.bool hasInboundConnections() => $_has(1);
  @$pb.TagNumber(2)
  void clearInboundConnections() => $_clearField(2);

  /// Number of outbound connections.
  @$pb.TagNumber(3)
  $fixnum.Int64 get outboundConnections => $_getI64(2);
  @$pb.TagNumber(3)
  set outboundConnections($fixnum.Int64 value) => $_setInt64(2, value);
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
    $core.Iterable<$core.MapEntry<$core.int, CounterInfo>>? messageSent,
    $core.Iterable<$core.MapEntry<$core.int, CounterInfo>>? messageReceived,
  }) {
    final result = create();
    if (totalInvalid != null) result.totalInvalid = totalInvalid;
    if (totalSent != null) result.totalSent = totalSent;
    if (totalReceived != null) result.totalReceived = totalReceived;
    if (messageSent != null) result.messageSent.addEntries(messageSent);
    if (messageReceived != null)
      result.messageReceived.addEntries(messageReceived);
    return result;
  }

  MetricInfo._();

  factory MetricInfo.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory MetricInfo.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'MetricInfo',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..aOM<CounterInfo>(1, _omitFieldNames ? '' : 'totalInvalid',
        subBuilder: CounterInfo.create)
    ..aOM<CounterInfo>(2, _omitFieldNames ? '' : 'totalSent',
        subBuilder: CounterInfo.create)
    ..aOM<CounterInfo>(3, _omitFieldNames ? '' : 'totalReceived',
        subBuilder: CounterInfo.create)
    ..m<$core.int, CounterInfo>(4, _omitFieldNames ? '' : 'messageSent',
        entryClassName: 'MetricInfo.MessageSentEntry',
        keyFieldType: $pb.PbFieldType.O3,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: CounterInfo.create,
        valueDefaultOrMaker: CounterInfo.getDefault,
        packageName: const $pb.PackageName('pactus'))
    ..m<$core.int, CounterInfo>(5, _omitFieldNames ? '' : 'messageReceived',
        entryClassName: 'MetricInfo.MessageReceivedEntry',
        keyFieldType: $pb.PbFieldType.O3,
        valueFieldType: $pb.PbFieldType.OM,
        valueCreator: CounterInfo.create,
        valueDefaultOrMaker: CounterInfo.getDefault,
        packageName: const $pb.PackageName('pactus'))
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MetricInfo clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  MetricInfo copyWith(void Function(MetricInfo) updates) =>
      super.copyWith((message) => updates(message as MetricInfo)) as MetricInfo;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static MetricInfo create() => MetricInfo._();
  @$core.override
  MetricInfo createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static MetricInfo getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<MetricInfo>(create);
  static MetricInfo? _defaultInstance;

  /// Total number of invalid bundles.
  @$pb.TagNumber(1)
  CounterInfo get totalInvalid => $_getN(0);
  @$pb.TagNumber(1)
  set totalInvalid(CounterInfo value) => $_setField(1, value);
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
  set totalSent(CounterInfo value) => $_setField(2, value);
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
  set totalReceived(CounterInfo value) => $_setField(3, value);
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
    final result = create();
    if (bytes != null) result.bytes = bytes;
    if (bundles != null) result.bundles = bundles;
    return result;
  }

  CounterInfo._();

  factory CounterInfo.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory CounterInfo.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'CounterInfo',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..a<$fixnum.Int64>(1, _omitFieldNames ? '' : 'bytes', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..a<$fixnum.Int64>(2, _omitFieldNames ? '' : 'bundles', $pb.PbFieldType.OU6,
        defaultOrMaker: $fixnum.Int64.ZERO)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CounterInfo clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  CounterInfo copyWith(void Function(CounterInfo) updates) =>
      super.copyWith((message) => updates(message as CounterInfo))
          as CounterInfo;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static CounterInfo create() => CounterInfo._();
  @$core.override
  CounterInfo createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static CounterInfo getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<CounterInfo>(create);
  static CounterInfo? _defaultInstance;

  /// Total number of bytes.
  @$pb.TagNumber(1)
  $fixnum.Int64 get bytes => $_getI64(0);
  @$pb.TagNumber(1)
  set bytes($fixnum.Int64 value) => $_setInt64(0, value);
  @$pb.TagNumber(1)
  $core.bool hasBytes() => $_has(0);
  @$pb.TagNumber(1)
  void clearBytes() => $_clearField(1);

  /// Total number of bundles.
  @$pb.TagNumber(2)
  $fixnum.Int64 get bundles => $_getI64(1);
  @$pb.TagNumber(2)
  set bundles($fixnum.Int64 value) => $_setInt64(1, value);
  @$pb.TagNumber(2)
  $core.bool hasBundles() => $_has(1);
  @$pb.TagNumber(2)
  void clearBundles() => $_clearField(2);
}

/// Request message for ping - intentionally empty for measuring round-trip time.
class PingRequest extends $pb.GeneratedMessage {
  factory PingRequest() => create();

  PingRequest._();

  factory PingRequest.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PingRequest.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PingRequest',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PingRequest clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PingRequest copyWith(void Function(PingRequest) updates) =>
      super.copyWith((message) => updates(message as PingRequest))
          as PingRequest;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PingRequest create() => PingRequest._();
  @$core.override
  PingRequest createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PingRequest getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PingRequest>(create);
  static PingRequest? _defaultInstance;
}

/// Response message for ping - intentionally empty for measuring round-trip time.
class PingResponse extends $pb.GeneratedMessage {
  factory PingResponse() => create();

  PingResponse._();

  factory PingResponse.fromBuffer($core.List<$core.int> data,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromBuffer(data, registry);
  factory PingResponse.fromJson($core.String json,
          [$pb.ExtensionRegistry registry = $pb.ExtensionRegistry.EMPTY]) =>
      create()..mergeFromJson(json, registry);

  static final $pb.BuilderInfo _i = $pb.BuilderInfo(
      _omitMessageNames ? '' : 'PingResponse',
      package: const $pb.PackageName(_omitMessageNames ? '' : 'pactus'),
      createEmptyInstance: create)
    ..hasRequiredFields = false;

  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PingResponse clone() => deepCopy();
  @$core.Deprecated('See https://github.com/google/protobuf.dart/issues/998.')
  PingResponse copyWith(void Function(PingResponse) updates) =>
      super.copyWith((message) => updates(message as PingResponse))
          as PingResponse;

  @$core.override
  $pb.BuilderInfo get info_ => _i;

  @$core.pragma('dart2js:noInline')
  static PingResponse create() => PingResponse._();
  @$core.override
  PingResponse createEmptyInstance() => create();
  @$core.pragma('dart2js:noInline')
  static PingResponse getDefault() => _defaultInstance ??=
      $pb.GeneratedMessage.$_defaultFor<PingResponse>(create);
  static PingResponse? _defaultInstance;
}

/// Network service provides RPCs for retrieving information about the network.
class NetworkApi {
  final $pb.RpcClient _client;

  NetworkApi(this._client);

  /// GetNetworkInfo retrieves information about the overall network.
  $async.Future<GetNetworkInfoResponse> getNetworkInfo(
          $pb.ClientContext? ctx, GetNetworkInfoRequest request) =>
      _client.invoke<GetNetworkInfoResponse>(
          ctx, 'Network', 'GetNetworkInfo', request, GetNetworkInfoResponse());

  /// ListPeers lists all peers in the network.
  $async.Future<ListPeersResponse> listPeers(
          $pb.ClientContext? ctx, ListPeersRequest request) =>
      _client.invoke<ListPeersResponse>(
          ctx, 'Network', 'ListPeers', request, ListPeersResponse());

  /// GetNodeInfo retrieves information about a specific node in the network.
  $async.Future<GetNodeInfoResponse> getNodeInfo(
          $pb.ClientContext? ctx, GetNodeInfoRequest request) =>
      _client.invoke<GetNodeInfoResponse>(
          ctx, 'Network', 'GetNodeInfo', request, GetNodeInfoResponse());

  /// Ping provides a simple connectivity test and latency measurement.
  $async.Future<PingResponse> ping(
          $pb.ClientContext? ctx, PingRequest request) =>
      _client.invoke<PingResponse>(
          ctx, 'Network', 'Ping', request, PingResponse());
}

const $core.bool _omitFieldNames =
    $core.bool.fromEnvironment('protobuf.omit_field_names');
const $core.bool _omitMessageNames =
    $core.bool.fromEnvironment('protobuf.omit_message_names');
