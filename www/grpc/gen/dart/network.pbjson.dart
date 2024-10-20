///
//  Generated code. Do not modify.
//  source: network.proto
//
// @dart = 2.12
// ignore_for_file: annotate_overrides,camel_case_types,constant_identifier_names,deprecated_member_use_from_same_package,directives_ordering,library_prefixes,non_constant_identifier_names,prefer_final_fields,return_of_invalid_type,unnecessary_const,unnecessary_import,unnecessary_this,unused_import,unused_shown_name

import 'dart:core' as $core;
import 'dart:convert' as $convert;
import 'dart:typed_data' as $typed_data;
@$core.Deprecated('Use getNetworkInfoRequestDescriptor instead')
const GetNetworkInfoRequest$json = const {
  '1': 'GetNetworkInfoRequest',
  '2': const [
    const {'1': 'only_connected', '3': 1, '4': 1, '5': 8, '10': 'onlyConnected'},
  ],
};

/// Descriptor for `GetNetworkInfoRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getNetworkInfoRequestDescriptor = $convert.base64Decode('ChVHZXROZXR3b3JrSW5mb1JlcXVlc3QSJQoOb25seV9jb25uZWN0ZWQYASABKAhSDW9ubHlDb25uZWN0ZWQ=');
@$core.Deprecated('Use getNetworkInfoResponseDescriptor instead')
const GetNetworkInfoResponse$json = const {
  '1': 'GetNetworkInfoResponse',
  '2': const [
    const {'1': 'network_name', '3': 1, '4': 1, '5': 9, '10': 'networkName'},
    const {'1': 'connected_peers_count', '3': 2, '4': 1, '5': 13, '10': 'connectedPeersCount'},
    const {'1': 'connected_peers', '3': 3, '4': 3, '5': 11, '6': '.pactus.PeerInfo', '10': 'connectedPeers'},
    const {'1': 'metric_info', '3': 4, '4': 1, '5': 11, '6': '.pactus.MetricInfo', '10': 'metricInfo'},
  ],
};

/// Descriptor for `GetNetworkInfoResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getNetworkInfoResponseDescriptor = $convert.base64Decode('ChZHZXROZXR3b3JrSW5mb1Jlc3BvbnNlEiEKDG5ldHdvcmtfbmFtZRgBIAEoCVILbmV0d29ya05hbWUSMgoVY29ubmVjdGVkX3BlZXJzX2NvdW50GAIgASgNUhNjb25uZWN0ZWRQZWVyc0NvdW50EjkKD2Nvbm5lY3RlZF9wZWVycxgDIAMoCzIQLnBhY3R1cy5QZWVySW5mb1IOY29ubmVjdGVkUGVlcnMSMwoLbWV0cmljX2luZm8YBCABKAsyEi5wYWN0dXMuTWV0cmljSW5mb1IKbWV0cmljSW5mbw==');
@$core.Deprecated('Use getNodeInfoRequestDescriptor instead')
const GetNodeInfoRequest$json = const {
  '1': 'GetNodeInfoRequest',
};

/// Descriptor for `GetNodeInfoRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getNodeInfoRequestDescriptor = $convert.base64Decode('ChJHZXROb2RlSW5mb1JlcXVlc3Q=');
@$core.Deprecated('Use getNodeInfoResponseDescriptor instead')
const GetNodeInfoResponse$json = const {
  '1': 'GetNodeInfoResponse',
  '2': const [
    const {'1': 'moniker', '3': 1, '4': 1, '5': 9, '10': 'moniker'},
    const {'1': 'agent', '3': 2, '4': 1, '5': 9, '10': 'agent'},
    const {'1': 'peer_id', '3': 3, '4': 1, '5': 9, '10': 'peerId'},
    const {'1': 'started_at', '3': 4, '4': 1, '5': 4, '10': 'startedAt'},
    const {'1': 'reachability', '3': 5, '4': 1, '5': 9, '10': 'reachability'},
    const {'1': 'services', '3': 6, '4': 1, '5': 5, '10': 'services'},
    const {'1': 'services_names', '3': 7, '4': 1, '5': 9, '10': 'servicesNames'},
    const {'1': 'local_addrs', '3': 8, '4': 3, '5': 9, '10': 'localAddrs'},
    const {'1': 'protocols', '3': 9, '4': 3, '5': 9, '10': 'protocols'},
    const {'1': 'clock_offset', '3': 13, '4': 1, '5': 1, '10': 'clockOffset'},
    const {'1': 'connection_info', '3': 14, '4': 1, '5': 11, '6': '.pactus.ConnectionInfo', '10': 'connectionInfo'},
    const {'1': 'fee', '3': 15, '4': 1, '5': 11, '6': '.pactus.FeeConfig', '10': 'fee'},
  ],
};

/// Descriptor for `GetNodeInfoResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getNodeInfoResponseDescriptor = $convert.base64Decode('ChNHZXROb2RlSW5mb1Jlc3BvbnNlEhgKB21vbmlrZXIYASABKAlSB21vbmlrZXISFAoFYWdlbnQYAiABKAlSBWFnZW50EhcKB3BlZXJfaWQYAyABKAlSBnBlZXJJZBIdCgpzdGFydGVkX2F0GAQgASgEUglzdGFydGVkQXQSIgoMcmVhY2hhYmlsaXR5GAUgASgJUgxyZWFjaGFiaWxpdHkSGgoIc2VydmljZXMYBiABKAVSCHNlcnZpY2VzEiUKDnNlcnZpY2VzX25hbWVzGAcgASgJUg1zZXJ2aWNlc05hbWVzEh8KC2xvY2FsX2FkZHJzGAggAygJUgpsb2NhbEFkZHJzEhwKCXByb3RvY29scxgJIAMoCVIJcHJvdG9jb2xzEiEKDGNsb2NrX29mZnNldBgNIAEoAVILY2xvY2tPZmZzZXQSPwoPY29ubmVjdGlvbl9pbmZvGA4gASgLMhYucGFjdHVzLkNvbm5lY3Rpb25JbmZvUg5jb25uZWN0aW9uSW5mbxIjCgNmZWUYDyABKAsyES5wYWN0dXMuRmVlQ29uZmlnUgNmZWU=');
@$core.Deprecated('Use peerInfoDescriptor instead')
const PeerInfo$json = const {
  '1': 'PeerInfo',
  '2': const [
    const {'1': 'status', '3': 1, '4': 1, '5': 5, '10': 'status'},
    const {'1': 'moniker', '3': 2, '4': 1, '5': 9, '10': 'moniker'},
    const {'1': 'agent', '3': 3, '4': 1, '5': 9, '10': 'agent'},
    const {'1': 'peer_id', '3': 4, '4': 1, '5': 9, '10': 'peerId'},
    const {'1': 'consensus_keys', '3': 5, '4': 3, '5': 9, '10': 'consensusKeys'},
    const {'1': 'consensus_addresses', '3': 6, '4': 3, '5': 9, '10': 'consensusAddresses'},
    const {'1': 'services', '3': 7, '4': 1, '5': 13, '10': 'services'},
    const {'1': 'last_block_hash', '3': 8, '4': 1, '5': 9, '10': 'lastBlockHash'},
    const {'1': 'height', '3': 9, '4': 1, '5': 13, '10': 'height'},
    const {'1': 'last_sent', '3': 10, '4': 1, '5': 3, '10': 'lastSent'},
    const {'1': 'last_received', '3': 11, '4': 1, '5': 3, '10': 'lastReceived'},
    const {'1': 'address', '3': 12, '4': 1, '5': 9, '10': 'address'},
    const {'1': 'direction', '3': 13, '4': 1, '5': 9, '10': 'direction'},
    const {'1': 'protocols', '3': 14, '4': 3, '5': 9, '10': 'protocols'},
    const {'1': 'total_sessions', '3': 15, '4': 1, '5': 5, '10': 'totalSessions'},
    const {'1': 'completed_sessions', '3': 16, '4': 1, '5': 5, '10': 'completedSessions'},
    const {'1': 'metric_info', '3': 17, '4': 1, '5': 11, '6': '.pactus.MetricInfo', '10': 'metricInfo'},
  ],
};

/// Descriptor for `PeerInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List peerInfoDescriptor = $convert.base64Decode('CghQZWVySW5mbxIWCgZzdGF0dXMYASABKAVSBnN0YXR1cxIYCgdtb25pa2VyGAIgASgJUgdtb25pa2VyEhQKBWFnZW50GAMgASgJUgVhZ2VudBIXCgdwZWVyX2lkGAQgASgJUgZwZWVySWQSJQoOY29uc2Vuc3VzX2tleXMYBSADKAlSDWNvbnNlbnN1c0tleXMSLwoTY29uc2Vuc3VzX2FkZHJlc3NlcxgGIAMoCVISY29uc2Vuc3VzQWRkcmVzc2VzEhoKCHNlcnZpY2VzGAcgASgNUghzZXJ2aWNlcxImCg9sYXN0X2Jsb2NrX2hhc2gYCCABKAlSDWxhc3RCbG9ja0hhc2gSFgoGaGVpZ2h0GAkgASgNUgZoZWlnaHQSGwoJbGFzdF9zZW50GAogASgDUghsYXN0U2VudBIjCg1sYXN0X3JlY2VpdmVkGAsgASgDUgxsYXN0UmVjZWl2ZWQSGAoHYWRkcmVzcxgMIAEoCVIHYWRkcmVzcxIcCglkaXJlY3Rpb24YDSABKAlSCWRpcmVjdGlvbhIcCglwcm90b2NvbHMYDiADKAlSCXByb3RvY29scxIlCg50b3RhbF9zZXNzaW9ucxgPIAEoBVINdG90YWxTZXNzaW9ucxItChJjb21wbGV0ZWRfc2Vzc2lvbnMYECABKAVSEWNvbXBsZXRlZFNlc3Npb25zEjMKC21ldHJpY19pbmZvGBEgASgLMhIucGFjdHVzLk1ldHJpY0luZm9SCm1ldHJpY0luZm8=');
@$core.Deprecated('Use connectionInfoDescriptor instead')
const ConnectionInfo$json = const {
  '1': 'ConnectionInfo',
  '2': const [
    const {'1': 'connections', '3': 1, '4': 1, '5': 4, '10': 'connections'},
    const {'1': 'inbound_connections', '3': 2, '4': 1, '5': 4, '10': 'inboundConnections'},
    const {'1': 'outbound_connections', '3': 3, '4': 1, '5': 4, '10': 'outboundConnections'},
  ],
};

/// Descriptor for `ConnectionInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List connectionInfoDescriptor = $convert.base64Decode('Cg5Db25uZWN0aW9uSW5mbxIgCgtjb25uZWN0aW9ucxgBIAEoBFILY29ubmVjdGlvbnMSLwoTaW5ib3VuZF9jb25uZWN0aW9ucxgCIAEoBFISaW5ib3VuZENvbm5lY3Rpb25zEjEKFG91dGJvdW5kX2Nvbm5lY3Rpb25zGAMgASgEUhNvdXRib3VuZENvbm5lY3Rpb25z');
@$core.Deprecated('Use metricInfoDescriptor instead')
const MetricInfo$json = const {
  '1': 'MetricInfo',
  '2': const [
    const {'1': 'TotalInvalid', '3': 1, '4': 1, '5': 11, '6': '.pactus.CounterInfo', '10': 'TotalInvalid'},
    const {'1': 'TotalSent', '3': 2, '4': 1, '5': 11, '6': '.pactus.CounterInfo', '10': 'TotalSent'},
    const {'1': 'TotalReceived', '3': 3, '4': 1, '5': 11, '6': '.pactus.CounterInfo', '10': 'TotalReceived'},
    const {'1': 'MessageSent', '3': 4, '4': 3, '5': 11, '6': '.pactus.MetricInfo.MessageSentEntry', '10': 'MessageSent'},
    const {'1': 'MessageReceived', '3': 5, '4': 3, '5': 11, '6': '.pactus.MetricInfo.MessageReceivedEntry', '10': 'MessageReceived'},
  ],
  '3': const [MetricInfo_MessageSentEntry$json, MetricInfo_MessageReceivedEntry$json],
};

@$core.Deprecated('Use metricInfoDescriptor instead')
const MetricInfo_MessageSentEntry$json = const {
  '1': 'MessageSentEntry',
  '2': const [
    const {'1': 'key', '3': 1, '4': 1, '5': 5, '10': 'key'},
    const {'1': 'value', '3': 2, '4': 1, '5': 11, '6': '.pactus.CounterInfo', '10': 'value'},
  ],
  '7': const {'7': true},
};

@$core.Deprecated('Use metricInfoDescriptor instead')
const MetricInfo_MessageReceivedEntry$json = const {
  '1': 'MessageReceivedEntry',
  '2': const [
    const {'1': 'key', '3': 1, '4': 1, '5': 5, '10': 'key'},
    const {'1': 'value', '3': 2, '4': 1, '5': 11, '6': '.pactus.CounterInfo', '10': 'value'},
  ],
  '7': const {'7': true},
};

/// Descriptor for `MetricInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List metricInfoDescriptor = $convert.base64Decode('CgpNZXRyaWNJbmZvEjcKDFRvdGFsSW52YWxpZBgBIAEoCzITLnBhY3R1cy5Db3VudGVySW5mb1IMVG90YWxJbnZhbGlkEjEKCVRvdGFsU2VudBgCIAEoCzITLnBhY3R1cy5Db3VudGVySW5mb1IJVG90YWxTZW50EjkKDVRvdGFsUmVjZWl2ZWQYAyABKAsyEy5wYWN0dXMuQ291bnRlckluZm9SDVRvdGFsUmVjZWl2ZWQSRQoLTWVzc2FnZVNlbnQYBCADKAsyIy5wYWN0dXMuTWV0cmljSW5mby5NZXNzYWdlU2VudEVudHJ5UgtNZXNzYWdlU2VudBJRCg9NZXNzYWdlUmVjZWl2ZWQYBSADKAsyJy5wYWN0dXMuTWV0cmljSW5mby5NZXNzYWdlUmVjZWl2ZWRFbnRyeVIPTWVzc2FnZVJlY2VpdmVkGlMKEE1lc3NhZ2VTZW50RW50cnkSEAoDa2V5GAEgASgFUgNrZXkSKQoFdmFsdWUYAiABKAsyEy5wYWN0dXMuQ291bnRlckluZm9SBXZhbHVlOgI4ARpXChRNZXNzYWdlUmVjZWl2ZWRFbnRyeRIQCgNrZXkYASABKAVSA2tleRIpCgV2YWx1ZRgCIAEoCzITLnBhY3R1cy5Db3VudGVySW5mb1IFdmFsdWU6AjgB');
@$core.Deprecated('Use counterInfoDescriptor instead')
const CounterInfo$json = const {
  '1': 'CounterInfo',
  '2': const [
    const {'1': 'Bytes', '3': 1, '4': 1, '5': 4, '10': 'Bytes'},
    const {'1': 'Bundles', '3': 2, '4': 1, '5': 4, '10': 'Bundles'},
  ],
};

/// Descriptor for `CounterInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List counterInfoDescriptor = $convert.base64Decode('CgtDb3VudGVySW5mbxIUCgVCeXRlcxgBIAEoBFIFQnl0ZXMSGAoHQnVuZGxlcxgCIAEoBFIHQnVuZGxlcw==');
@$core.Deprecated('Use feeConfigDescriptor instead')
const FeeConfig$json = const {
  '1': 'FeeConfig',
  '2': const [
    const {'1': 'fixed_fee', '3': 1, '4': 1, '5': 1, '10': 'fixedFee'},
    const {'1': 'daily_limit', '3': 2, '4': 1, '5': 13, '10': 'dailyLimit'},
    const {'1': 'unit_price', '3': 3, '4': 1, '5': 1, '10': 'unitPrice'},
  ],
};

/// Descriptor for `FeeConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List feeConfigDescriptor = $convert.base64Decode('CglGZWVDb25maWcSGwoJZml4ZWRfZmVlGAEgASgBUghmaXhlZEZlZRIfCgtkYWlseV9saW1pdBgCIAEoDVIKZGFpbHlMaW1pdBIdCgp1bml0X3ByaWNlGAMgASgBUgl1bml0UHJpY2U=');
const $core.Map<$core.String, $core.dynamic> NetworkServiceBase$json = const {
  '1': 'Network',
  '2': const [
    const {'1': 'GetNetworkInfo', '2': '.pactus.GetNetworkInfoRequest', '3': '.pactus.GetNetworkInfoResponse'},
    const {'1': 'GetNodeInfo', '2': '.pactus.GetNodeInfoRequest', '3': '.pactus.GetNodeInfoResponse'},
  ],
};

@$core.Deprecated('Use networkServiceDescriptor instead')
const $core.Map<$core.String, $core.Map<$core.String, $core.dynamic>> NetworkServiceBase$messageJson = const {
  '.pactus.GetNetworkInfoRequest': GetNetworkInfoRequest$json,
  '.pactus.GetNetworkInfoResponse': GetNetworkInfoResponse$json,
  '.pactus.PeerInfo': PeerInfo$json,
  '.pactus.MetricInfo': MetricInfo$json,
  '.pactus.CounterInfo': CounterInfo$json,
  '.pactus.MetricInfo.MessageSentEntry': MetricInfo_MessageSentEntry$json,
  '.pactus.MetricInfo.MessageReceivedEntry': MetricInfo_MessageReceivedEntry$json,
  '.pactus.GetNodeInfoRequest': GetNodeInfoRequest$json,
  '.pactus.GetNodeInfoResponse': GetNodeInfoResponse$json,
  '.pactus.ConnectionInfo': ConnectionInfo$json,
  '.pactus.FeeConfig': FeeConfig$json,
};

/// Descriptor for `Network`. Decode as a `google.protobuf.ServiceDescriptorProto`.
final $typed_data.Uint8List networkServiceDescriptor = $convert.base64Decode('CgdOZXR3b3JrEk8KDkdldE5ldHdvcmtJbmZvEh0ucGFjdHVzLkdldE5ldHdvcmtJbmZvUmVxdWVzdBoeLnBhY3R1cy5HZXROZXR3b3JrSW5mb1Jlc3BvbnNlEkYKC0dldE5vZGVJbmZvEhoucGFjdHVzLkdldE5vZGVJbmZvUmVxdWVzdBobLnBhY3R1cy5HZXROb2RlSW5mb1Jlc3BvbnNl');
