///
//  Generated code. Do not modify.
//  source: network.proto
//
// @dart = 2.12
// ignore_for_file: annotate_overrides,camel_case_types,constant_identifier_names,directives_ordering,library_prefixes,non_constant_identifier_names,prefer_final_fields,return_of_invalid_type,unnecessary_const,unnecessary_import,unnecessary_this,unused_import,unused_shown_name

import 'dart:async' as $async;
import 'dart:core' as $core;

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
    ..a<$core.List<$core.int>>(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'selfId', $pb.PbFieldType.OY)
    ..pc<PeerInfo>(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'peers', $pb.PbFieldType.PM, subBuilder: PeerInfo.create)
    ..hasRequiredFields = false
  ;

  GetNetworkInfoResponse._() : super();
  factory GetNetworkInfoResponse({
    $core.List<$core.int>? selfId,
    $core.Iterable<PeerInfo>? peers,
  }) {
    final _result = create();
    if (selfId != null) {
      _result.selfId = selfId;
    }
    if (peers != null) {
      _result.peers.addAll(peers);
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
  $core.List<$core.int> get selfId => $_getN(0);
  @$pb.TagNumber(1)
  set selfId($core.List<$core.int> v) { $_setBytes(0, v); }
  @$pb.TagNumber(1)
  $core.bool hasSelfId() => $_has(0);
  @$pb.TagNumber(1)
  void clearSelfId() => clearField(1);

  @$pb.TagNumber(2)
  $core.List<PeerInfo> get peers => $_getList(1);
}

class GetPeerInfoRequest extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'GetPeerInfoRequest', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..hasRequiredFields = false
  ;

  GetPeerInfoRequest._() : super();
  factory GetPeerInfoRequest() => create();
  factory GetPeerInfoRequest.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetPeerInfoRequest.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetPeerInfoRequest clone() => GetPeerInfoRequest()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetPeerInfoRequest copyWith(void Function(GetPeerInfoRequest) updates) => super.copyWith((message) => updates(message as GetPeerInfoRequest)) as GetPeerInfoRequest; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static GetPeerInfoRequest create() => GetPeerInfoRequest._();
  GetPeerInfoRequest createEmptyInstance() => create();
  static $pb.PbList<GetPeerInfoRequest> createRepeated() => $pb.PbList<GetPeerInfoRequest>();
  @$core.pragma('dart2js:noInline')
  static GetPeerInfoRequest getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetPeerInfoRequest>(create);
  static GetPeerInfoRequest? _defaultInstance;
}

class GetPeerInfoResponse extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'GetPeerInfoResponse', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'moniker')
    ..aOS(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'agent')
    ..a<$core.List<$core.int>>(3, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'peerId', $pb.PbFieldType.OY)
    ..aOS(4, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'publicKey')
    ..a<$core.int>(5, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'flags', $pb.PbFieldType.O3)
    ..a<$core.int>(6, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'height', $pb.PbFieldType.OU3)
    ..a<$core.int>(7, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'receivedMessages', $pb.PbFieldType.O3)
    ..a<$core.int>(8, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'invalidMessages', $pb.PbFieldType.O3)
    ..a<$core.int>(9, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'receivedBytes', $pb.PbFieldType.O3)
    ..hasRequiredFields = false
  ;

  GetPeerInfoResponse._() : super();
  factory GetPeerInfoResponse({
    $core.String? moniker,
    $core.String? agent,
    $core.List<$core.int>? peerId,
    $core.String? publicKey,
    $core.int? flags,
    $core.int? height,
    $core.int? receivedMessages,
    $core.int? invalidMessages,
    $core.int? receivedBytes,
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
    if (publicKey != null) {
      _result.publicKey = publicKey;
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
    return _result;
  }
  factory GetPeerInfoResponse.fromBuffer($core.List<$core.int> i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromBuffer(i, r);
  factory GetPeerInfoResponse.fromJson($core.String i, [$pb.ExtensionRegistry r = $pb.ExtensionRegistry.EMPTY]) => create()..mergeFromJson(i, r);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.deepCopy] instead. '
  'Will be removed in next major version')
  GetPeerInfoResponse clone() => GetPeerInfoResponse()..mergeFromMessage(this);
  @$core.Deprecated(
  'Using this can add significant overhead to your binary. '
  'Use [GeneratedMessageGenericExtensions.rebuild] instead. '
  'Will be removed in next major version')
  GetPeerInfoResponse copyWith(void Function(GetPeerInfoResponse) updates) => super.copyWith((message) => updates(message as GetPeerInfoResponse)) as GetPeerInfoResponse; // ignore: deprecated_member_use
  $pb.BuilderInfo get info_ => _i;
  @$core.pragma('dart2js:noInline')
  static GetPeerInfoResponse create() => GetPeerInfoResponse._();
  GetPeerInfoResponse createEmptyInstance() => create();
  static $pb.PbList<GetPeerInfoResponse> createRepeated() => $pb.PbList<GetPeerInfoResponse>();
  @$core.pragma('dart2js:noInline')
  static GetPeerInfoResponse getDefault() => _defaultInstance ??= $pb.GeneratedMessage.$_defaultFor<GetPeerInfoResponse>(create);
  static GetPeerInfoResponse? _defaultInstance;

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
  $core.String get publicKey => $_getSZ(3);
  @$pb.TagNumber(4)
  set publicKey($core.String v) { $_setString(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasPublicKey() => $_has(3);
  @$pb.TagNumber(4)
  void clearPublicKey() => clearField(4);

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
}

class PeerInfo extends $pb.GeneratedMessage {
  static final $pb.BuilderInfo _i = $pb.BuilderInfo(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'PeerInfo', package: const $pb.PackageName(const $core.bool.fromEnvironment('protobuf.omit_message_names') ? '' : 'pactus'), createEmptyInstance: create)
    ..aOS(1, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'moniker')
    ..aOS(2, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'agent')
    ..a<$core.List<$core.int>>(3, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'peerId', $pb.PbFieldType.OY)
    ..aOS(4, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'publicKey')
    ..a<$core.int>(5, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'flags', $pb.PbFieldType.O3)
    ..a<$core.int>(6, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'height', $pb.PbFieldType.OU3)
    ..a<$core.int>(7, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'receivedMessages', $pb.PbFieldType.O3)
    ..a<$core.int>(8, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'invalidMessages', $pb.PbFieldType.O3)
    ..a<$core.int>(9, const $core.bool.fromEnvironment('protobuf.omit_field_names') ? '' : 'receivedBytes', $pb.PbFieldType.O3)
    ..hasRequiredFields = false
  ;

  PeerInfo._() : super();
  factory PeerInfo({
    $core.String? moniker,
    $core.String? agent,
    $core.List<$core.int>? peerId,
    $core.String? publicKey,
    $core.int? flags,
    $core.int? height,
    $core.int? receivedMessages,
    $core.int? invalidMessages,
    $core.int? receivedBytes,
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
    if (publicKey != null) {
      _result.publicKey = publicKey;
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
  $core.String get publicKey => $_getSZ(3);
  @$pb.TagNumber(4)
  set publicKey($core.String v) { $_setString(3, v); }
  @$pb.TagNumber(4)
  $core.bool hasPublicKey() => $_has(3);
  @$pb.TagNumber(4)
  void clearPublicKey() => clearField(4);

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
}

class NetworkApi {
  $pb.RpcClient _client;
  NetworkApi(this._client);

  $async.Future<GetNetworkInfoResponse> getNetworkInfo($pb.ClientContext? ctx, GetNetworkInfoRequest request) {
    var emptyResponse = GetNetworkInfoResponse();
    return _client.invoke<GetNetworkInfoResponse>(ctx, 'Network', 'GetNetworkInfo', request, emptyResponse);
  }
  $async.Future<GetPeerInfoResponse> getPeerInfo($pb.ClientContext? ctx, GetPeerInfoRequest request) {
    var emptyResponse = GetPeerInfoResponse();
    return _client.invoke<GetPeerInfoResponse>(ctx, 'Network', 'GetPeerInfo', request, emptyResponse);
  }
}

