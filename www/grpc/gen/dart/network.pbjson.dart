//
//  Generated code. Do not modify.
//  source: network.proto
//
// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_final_fields
// ignore_for_file: unnecessary_import, unnecessary_this, unused_import

import 'dart:convert' as $convert;
import 'dart:core' as $core;
import 'dart:typed_data' as $typed_data;

@$core.Deprecated('Use getNetworkInfoRequestDescriptor instead')
const GetNetworkInfoRequest$json = {
  '1': 'GetNetworkInfoRequest',
  '2': [
    {'1': 'only_connected', '3': 1, '4': 1, '5': 8, '10': 'onlyConnected'},
  ],
};

/// Descriptor for `GetNetworkInfoRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getNetworkInfoRequestDescriptor = $convert.base64Decode(
    'ChVHZXROZXR3b3JrSW5mb1JlcXVlc3QSJQoOb25seV9jb25uZWN0ZWQYASABKAhSDW9ubHlDb2'
    '5uZWN0ZWQ=');

@$core.Deprecated('Use getNetworkInfoResponseDescriptor instead')
const GetNetworkInfoResponse$json = {
  '1': 'GetNetworkInfoResponse',
  '2': [
    {'1': 'network_name', '3': 1, '4': 1, '5': 9, '10': 'networkName'},
    {'1': 'connected_peers_count', '3': 2, '4': 1, '5': 13, '10': 'connectedPeersCount'},
    {'1': 'connected_peers', '3': 3, '4': 3, '5': 11, '6': '.pactus.PeerInfo', '10': 'connectedPeers'},
    {'1': 'metric_info', '3': 4, '4': 1, '5': 11, '6': '.pactus.MetricInfo', '10': 'metricInfo'},
  ],
};

/// Descriptor for `GetNetworkInfoResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getNetworkInfoResponseDescriptor = $convert.base64Decode(
    'ChZHZXROZXR3b3JrSW5mb1Jlc3BvbnNlEiEKDG5ldHdvcmtfbmFtZRgBIAEoCVILbmV0d29ya0'
    '5hbWUSMgoVY29ubmVjdGVkX3BlZXJzX2NvdW50GAIgASgNUhNjb25uZWN0ZWRQZWVyc0NvdW50'
    'EjkKD2Nvbm5lY3RlZF9wZWVycxgDIAMoCzIQLnBhY3R1cy5QZWVySW5mb1IOY29ubmVjdGVkUG'
    'VlcnMSMwoLbWV0cmljX2luZm8YBCABKAsyEi5wYWN0dXMuTWV0cmljSW5mb1IKbWV0cmljSW5m'
    'bw==');

@$core.Deprecated('Use getNodeInfoRequestDescriptor instead')
const GetNodeInfoRequest$json = {
  '1': 'GetNodeInfoRequest',
};

/// Descriptor for `GetNodeInfoRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getNodeInfoRequestDescriptor = $convert.base64Decode(
    'ChJHZXROb2RlSW5mb1JlcXVlc3Q=');

@$core.Deprecated('Use getNodeInfoResponseDescriptor instead')
const GetNodeInfoResponse$json = {
  '1': 'GetNodeInfoResponse',
  '2': [
    {'1': 'moniker', '3': 1, '4': 1, '5': 9, '10': 'moniker'},
    {'1': 'agent', '3': 2, '4': 1, '5': 9, '10': 'agent'},
    {'1': 'peer_id', '3': 3, '4': 1, '5': 9, '10': 'peerId'},
    {'1': 'started_at', '3': 4, '4': 1, '5': 4, '10': 'startedAt'},
    {'1': 'reachability', '3': 5, '4': 1, '5': 9, '10': 'reachability'},
    {'1': 'services', '3': 6, '4': 1, '5': 5, '10': 'services'},
    {'1': 'services_names', '3': 7, '4': 1, '5': 9, '10': 'servicesNames'},
    {'1': 'local_addrs', '3': 8, '4': 3, '5': 9, '10': 'localAddrs'},
    {'1': 'protocols', '3': 9, '4': 3, '5': 9, '10': 'protocols'},
    {'1': 'clock_offset', '3': 13, '4': 1, '5': 1, '10': 'clockOffset'},
    {'1': 'connection_info', '3': 14, '4': 1, '5': 11, '6': '.pactus.ConnectionInfo', '10': 'connectionInfo'},
    {'1': 'zmq_publishers', '3': 15, '4': 3, '5': 11, '6': '.pactus.ZMQPublisherInfo', '10': 'zmqPublishers'},
  ],
};

/// Descriptor for `GetNodeInfoResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getNodeInfoResponseDescriptor = $convert.base64Decode(
    'ChNHZXROb2RlSW5mb1Jlc3BvbnNlEhgKB21vbmlrZXIYASABKAlSB21vbmlrZXISFAoFYWdlbn'
    'QYAiABKAlSBWFnZW50EhcKB3BlZXJfaWQYAyABKAlSBnBlZXJJZBIdCgpzdGFydGVkX2F0GAQg'
    'ASgEUglzdGFydGVkQXQSIgoMcmVhY2hhYmlsaXR5GAUgASgJUgxyZWFjaGFiaWxpdHkSGgoIc2'
    'VydmljZXMYBiABKAVSCHNlcnZpY2VzEiUKDnNlcnZpY2VzX25hbWVzGAcgASgJUg1zZXJ2aWNl'
    'c05hbWVzEh8KC2xvY2FsX2FkZHJzGAggAygJUgpsb2NhbEFkZHJzEhwKCXByb3RvY29scxgJIA'
    'MoCVIJcHJvdG9jb2xzEiEKDGNsb2NrX29mZnNldBgNIAEoAVILY2xvY2tPZmZzZXQSPwoPY29u'
    'bmVjdGlvbl9pbmZvGA4gASgLMhYucGFjdHVzLkNvbm5lY3Rpb25JbmZvUg5jb25uZWN0aW9uSW'
    '5mbxI/Cg56bXFfcHVibGlzaGVycxgPIAMoCzIYLnBhY3R1cy5aTVFQdWJsaXNoZXJJbmZvUg16'
    'bXFQdWJsaXNoZXJz');

@$core.Deprecated('Use zMQPublisherInfoDescriptor instead')
const ZMQPublisherInfo$json = {
  '1': 'ZMQPublisherInfo',
  '2': [
    {'1': 'topic', '3': 1, '4': 1, '5': 9, '10': 'topic'},
    {'1': 'address', '3': 2, '4': 1, '5': 9, '10': 'address'},
    {'1': 'hwm', '3': 3, '4': 1, '5': 5, '10': 'hwm'},
  ],
};

/// Descriptor for `ZMQPublisherInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List zMQPublisherInfoDescriptor = $convert.base64Decode(
    'ChBaTVFQdWJsaXNoZXJJbmZvEhQKBXRvcGljGAEgASgJUgV0b3BpYxIYCgdhZGRyZXNzGAIgAS'
    'gJUgdhZGRyZXNzEhAKA2h3bRgDIAEoBVIDaHdt');

@$core.Deprecated('Use peerInfoDescriptor instead')
const PeerInfo$json = {
  '1': 'PeerInfo',
  '2': [
    {'1': 'status', '3': 1, '4': 1, '5': 5, '10': 'status'},
    {'1': 'moniker', '3': 2, '4': 1, '5': 9, '10': 'moniker'},
    {'1': 'agent', '3': 3, '4': 1, '5': 9, '10': 'agent'},
    {'1': 'peer_id', '3': 4, '4': 1, '5': 9, '10': 'peerId'},
    {'1': 'consensus_keys', '3': 5, '4': 3, '5': 9, '10': 'consensusKeys'},
    {'1': 'consensus_addresses', '3': 6, '4': 3, '5': 9, '10': 'consensusAddresses'},
    {'1': 'services', '3': 7, '4': 1, '5': 13, '10': 'services'},
    {'1': 'last_block_hash', '3': 8, '4': 1, '5': 9, '10': 'lastBlockHash'},
    {'1': 'height', '3': 9, '4': 1, '5': 13, '10': 'height'},
    {'1': 'last_sent', '3': 10, '4': 1, '5': 3, '10': 'lastSent'},
    {'1': 'last_received', '3': 11, '4': 1, '5': 3, '10': 'lastReceived'},
    {'1': 'address', '3': 12, '4': 1, '5': 9, '10': 'address'},
    {'1': 'direction', '3': 13, '4': 1, '5': 9, '10': 'direction'},
    {'1': 'protocols', '3': 14, '4': 3, '5': 9, '10': 'protocols'},
    {'1': 'total_sessions', '3': 15, '4': 1, '5': 5, '10': 'totalSessions'},
    {'1': 'completed_sessions', '3': 16, '4': 1, '5': 5, '10': 'completedSessions'},
    {'1': 'metric_info', '3': 17, '4': 1, '5': 11, '6': '.pactus.MetricInfo', '10': 'metricInfo'},
  ],
};

/// Descriptor for `PeerInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List peerInfoDescriptor = $convert.base64Decode(
    'CghQZWVySW5mbxIWCgZzdGF0dXMYASABKAVSBnN0YXR1cxIYCgdtb25pa2VyGAIgASgJUgdtb2'
    '5pa2VyEhQKBWFnZW50GAMgASgJUgVhZ2VudBIXCgdwZWVyX2lkGAQgASgJUgZwZWVySWQSJQoO'
    'Y29uc2Vuc3VzX2tleXMYBSADKAlSDWNvbnNlbnN1c0tleXMSLwoTY29uc2Vuc3VzX2FkZHJlc3'
    'NlcxgGIAMoCVISY29uc2Vuc3VzQWRkcmVzc2VzEhoKCHNlcnZpY2VzGAcgASgNUghzZXJ2aWNl'
    'cxImCg9sYXN0X2Jsb2NrX2hhc2gYCCABKAlSDWxhc3RCbG9ja0hhc2gSFgoGaGVpZ2h0GAkgAS'
    'gNUgZoZWlnaHQSGwoJbGFzdF9zZW50GAogASgDUghsYXN0U2VudBIjCg1sYXN0X3JlY2VpdmVk'
    'GAsgASgDUgxsYXN0UmVjZWl2ZWQSGAoHYWRkcmVzcxgMIAEoCVIHYWRkcmVzcxIcCglkaXJlY3'
    'Rpb24YDSABKAlSCWRpcmVjdGlvbhIcCglwcm90b2NvbHMYDiADKAlSCXByb3RvY29scxIlCg50'
    'b3RhbF9zZXNzaW9ucxgPIAEoBVINdG90YWxTZXNzaW9ucxItChJjb21wbGV0ZWRfc2Vzc2lvbn'
    'MYECABKAVSEWNvbXBsZXRlZFNlc3Npb25zEjMKC21ldHJpY19pbmZvGBEgASgLMhIucGFjdHVz'
    'Lk1ldHJpY0luZm9SCm1ldHJpY0luZm8=');

@$core.Deprecated('Use connectionInfoDescriptor instead')
const ConnectionInfo$json = {
  '1': 'ConnectionInfo',
  '2': [
    {'1': 'connections', '3': 1, '4': 1, '5': 4, '10': 'connections'},
    {'1': 'inbound_connections', '3': 2, '4': 1, '5': 4, '10': 'inboundConnections'},
    {'1': 'outbound_connections', '3': 3, '4': 1, '5': 4, '10': 'outboundConnections'},
  ],
};

/// Descriptor for `ConnectionInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List connectionInfoDescriptor = $convert.base64Decode(
    'Cg5Db25uZWN0aW9uSW5mbxIgCgtjb25uZWN0aW9ucxgBIAEoBFILY29ubmVjdGlvbnMSLwoTaW'
    '5ib3VuZF9jb25uZWN0aW9ucxgCIAEoBFISaW5ib3VuZENvbm5lY3Rpb25zEjEKFG91dGJvdW5k'
    'X2Nvbm5lY3Rpb25zGAMgASgEUhNvdXRib3VuZENvbm5lY3Rpb25z');

@$core.Deprecated('Use metricInfoDescriptor instead')
const MetricInfo$json = {
  '1': 'MetricInfo',
  '2': [
    {'1': 'total_invalid', '3': 1, '4': 1, '5': 11, '6': '.pactus.CounterInfo', '10': 'totalInvalid'},
    {'1': 'total_sent', '3': 2, '4': 1, '5': 11, '6': '.pactus.CounterInfo', '10': 'totalSent'},
    {'1': 'total_received', '3': 3, '4': 1, '5': 11, '6': '.pactus.CounterInfo', '10': 'totalReceived'},
    {'1': 'message_sent', '3': 4, '4': 3, '5': 11, '6': '.pactus.MetricInfo.MessageSentEntry', '10': 'messageSent'},
    {'1': 'message_received', '3': 5, '4': 3, '5': 11, '6': '.pactus.MetricInfo.MessageReceivedEntry', '10': 'messageReceived'},
  ],
  '3': [MetricInfo_MessageSentEntry$json, MetricInfo_MessageReceivedEntry$json],
};

@$core.Deprecated('Use metricInfoDescriptor instead')
const MetricInfo_MessageSentEntry$json = {
  '1': 'MessageSentEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 5, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 11, '6': '.pactus.CounterInfo', '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use metricInfoDescriptor instead')
const MetricInfo_MessageReceivedEntry$json = {
  '1': 'MessageReceivedEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 5, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 11, '6': '.pactus.CounterInfo', '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `MetricInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List metricInfoDescriptor = $convert.base64Decode(
    'CgpNZXRyaWNJbmZvEjgKDXRvdGFsX2ludmFsaWQYASABKAsyEy5wYWN0dXMuQ291bnRlckluZm'
    '9SDHRvdGFsSW52YWxpZBIyCgp0b3RhbF9zZW50GAIgASgLMhMucGFjdHVzLkNvdW50ZXJJbmZv'
    'Ugl0b3RhbFNlbnQSOgoOdG90YWxfcmVjZWl2ZWQYAyABKAsyEy5wYWN0dXMuQ291bnRlckluZm'
    '9SDXRvdGFsUmVjZWl2ZWQSRgoMbWVzc2FnZV9zZW50GAQgAygLMiMucGFjdHVzLk1ldHJpY0lu'
    'Zm8uTWVzc2FnZVNlbnRFbnRyeVILbWVzc2FnZVNlbnQSUgoQbWVzc2FnZV9yZWNlaXZlZBgFIA'
    'MoCzInLnBhY3R1cy5NZXRyaWNJbmZvLk1lc3NhZ2VSZWNlaXZlZEVudHJ5Ug9tZXNzYWdlUmVj'
    'ZWl2ZWQaUwoQTWVzc2FnZVNlbnRFbnRyeRIQCgNrZXkYASABKAVSA2tleRIpCgV2YWx1ZRgCIA'
    'EoCzITLnBhY3R1cy5Db3VudGVySW5mb1IFdmFsdWU6AjgBGlcKFE1lc3NhZ2VSZWNlaXZlZEVu'
    'dHJ5EhAKA2tleRgBIAEoBVIDa2V5EikKBXZhbHVlGAIgASgLMhMucGFjdHVzLkNvdW50ZXJJbm'
    'ZvUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use counterInfoDescriptor instead')
const CounterInfo$json = {
  '1': 'CounterInfo',
  '2': [
    {'1': 'bytes', '3': 1, '4': 1, '5': 4, '10': 'bytes'},
    {'1': 'bundles', '3': 2, '4': 1, '5': 4, '10': 'bundles'},
  ],
};

/// Descriptor for `CounterInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List counterInfoDescriptor = $convert.base64Decode(
    'CgtDb3VudGVySW5mbxIUCgVieXRlcxgBIAEoBFIFYnl0ZXMSGAoHYnVuZGxlcxgCIAEoBFIHYn'
    'VuZGxlcw==');

const $core.Map<$core.String, $core.dynamic> NetworkServiceBase$json = {
  '1': 'Network',
  '2': [
    {'1': 'GetNetworkInfo', '2': '.pactus.GetNetworkInfoRequest', '3': '.pactus.GetNetworkInfoResponse'},
    {'1': 'GetNodeInfo', '2': '.pactus.GetNodeInfoRequest', '3': '.pactus.GetNodeInfoResponse'},
  ],
};

@$core.Deprecated('Use networkServiceDescriptor instead')
const $core.Map<$core.String, $core.Map<$core.String, $core.dynamic>> NetworkServiceBase$messageJson = {
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
  '.pactus.ZMQPublisherInfo': ZMQPublisherInfo$json,
};

/// Descriptor for `Network`. Decode as a `google.protobuf.ServiceDescriptorProto`.
final $typed_data.Uint8List networkServiceDescriptor = $convert.base64Decode(
    'CgdOZXR3b3JrEk8KDkdldE5ldHdvcmtJbmZvEh0ucGFjdHVzLkdldE5ldHdvcmtJbmZvUmVxdW'
    'VzdBoeLnBhY3R1cy5HZXROZXR3b3JrSW5mb1Jlc3BvbnNlEkYKC0dldE5vZGVJbmZvEhoucGFj'
    'dHVzLkdldE5vZGVJbmZvUmVxdWVzdBobLnBhY3R1cy5HZXROb2RlSW5mb1Jlc3BvbnNl');

