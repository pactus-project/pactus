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
    const {'1': 'total_sent_bytes', '3': 2, '4': 1, '5': 3, '10': 'totalSentBytes'},
    const {'1': 'total_received_bytes', '3': 3, '4': 1, '5': 3, '10': 'totalReceivedBytes'},
    const {'1': 'connected_peers_count', '3': 4, '4': 1, '5': 13, '10': 'connectedPeersCount'},
    const {'1': 'connected_peers', '3': 5, '4': 3, '5': 11, '6': '.pactus.PeerInfo', '10': 'connectedPeers'},
    const {'1': 'sent_bytes', '3': 6, '4': 3, '5': 11, '6': '.pactus.GetNetworkInfoResponse.SentBytesEntry', '10': 'sentBytes'},
    const {'1': 'received_bytes', '3': 7, '4': 3, '5': 11, '6': '.pactus.GetNetworkInfoResponse.ReceivedBytesEntry', '10': 'receivedBytes'},
  ],
  '3': const [GetNetworkInfoResponse_SentBytesEntry$json, GetNetworkInfoResponse_ReceivedBytesEntry$json],
};

@$core.Deprecated('Use getNetworkInfoResponseDescriptor instead')
const GetNetworkInfoResponse_SentBytesEntry$json = const {
  '1': 'SentBytesEntry',
  '2': const [
    const {'1': 'key', '3': 1, '4': 1, '5': 5, '10': 'key'},
    const {'1': 'value', '3': 2, '4': 1, '5': 3, '10': 'value'},
  ],
  '7': const {'7': true},
};

@$core.Deprecated('Use getNetworkInfoResponseDescriptor instead')
const GetNetworkInfoResponse_ReceivedBytesEntry$json = const {
  '1': 'ReceivedBytesEntry',
  '2': const [
    const {'1': 'key', '3': 1, '4': 1, '5': 5, '10': 'key'},
    const {'1': 'value', '3': 2, '4': 1, '5': 3, '10': 'value'},
  ],
  '7': const {'7': true},
};

/// Descriptor for `GetNetworkInfoResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getNetworkInfoResponseDescriptor = $convert.base64Decode('ChZHZXROZXR3b3JrSW5mb1Jlc3BvbnNlEiEKDG5ldHdvcmtfbmFtZRgBIAEoCVILbmV0d29ya05hbWUSKAoQdG90YWxfc2VudF9ieXRlcxgCIAEoA1IOdG90YWxTZW50Qnl0ZXMSMAoUdG90YWxfcmVjZWl2ZWRfYnl0ZXMYAyABKANSEnRvdGFsUmVjZWl2ZWRCeXRlcxIyChVjb25uZWN0ZWRfcGVlcnNfY291bnQYBCABKA1SE2Nvbm5lY3RlZFBlZXJzQ291bnQSOQoPY29ubmVjdGVkX3BlZXJzGAUgAygLMhAucGFjdHVzLlBlZXJJbmZvUg5jb25uZWN0ZWRQZWVycxJMCgpzZW50X2J5dGVzGAYgAygLMi0ucGFjdHVzLkdldE5ldHdvcmtJbmZvUmVzcG9uc2UuU2VudEJ5dGVzRW50cnlSCXNlbnRCeXRlcxJYCg5yZWNlaXZlZF9ieXRlcxgHIAMoCzIxLnBhY3R1cy5HZXROZXR3b3JrSW5mb1Jlc3BvbnNlLlJlY2VpdmVkQnl0ZXNFbnRyeVINcmVjZWl2ZWRCeXRlcxo8Cg5TZW50Qnl0ZXNFbnRyeRIQCgNrZXkYASABKAVSA2tleRIUCgV2YWx1ZRgCIAEoA1IFdmFsdWU6AjgBGkAKElJlY2VpdmVkQnl0ZXNFbnRyeRIQCgNrZXkYASABKAVSA2tleRIUCgV2YWx1ZRgCIAEoA1IFdmFsdWU6AjgB');
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
  ],
};

/// Descriptor for `GetNodeInfoResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getNodeInfoResponseDescriptor = $convert.base64Decode('ChNHZXROb2RlSW5mb1Jlc3BvbnNlEhgKB21vbmlrZXIYASABKAlSB21vbmlrZXISFAoFYWdlbnQYAiABKAlSBWFnZW50EhcKB3BlZXJfaWQYAyABKAlSBnBlZXJJZBIdCgpzdGFydGVkX2F0GAQgASgEUglzdGFydGVkQXQSIgoMcmVhY2hhYmlsaXR5GAUgASgJUgxyZWFjaGFiaWxpdHkSGgoIc2VydmljZXMYBiABKAVSCHNlcnZpY2VzEiUKDnNlcnZpY2VzX25hbWVzGAcgASgJUg1zZXJ2aWNlc05hbWVzEh8KC2xvY2FsX2FkZHJzGAggAygJUgpsb2NhbEFkZHJzEhwKCXByb3RvY29scxgJIAMoCVIJcHJvdG9jb2xzEiEKDGNsb2NrX29mZnNldBgNIAEoAVILY2xvY2tPZmZzZXQSPwoPY29ubmVjdGlvbl9pbmZvGA4gASgLMhYucGFjdHVzLkNvbm5lY3Rpb25JbmZvUg5jb25uZWN0aW9uSW5mbw==');
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
    const {'1': 'received_bundles', '3': 10, '4': 1, '5': 5, '10': 'receivedBundles'},
    const {'1': 'invalid_bundles', '3': 11, '4': 1, '5': 5, '10': 'invalidBundles'},
    const {'1': 'last_sent', '3': 12, '4': 1, '5': 3, '10': 'lastSent'},
    const {'1': 'last_received', '3': 13, '4': 1, '5': 3, '10': 'lastReceived'},
    const {'1': 'sent_bytes', '3': 14, '4': 3, '5': 11, '6': '.pactus.PeerInfo.SentBytesEntry', '10': 'sentBytes'},
    const {'1': 'received_bytes', '3': 15, '4': 3, '5': 11, '6': '.pactus.PeerInfo.ReceivedBytesEntry', '10': 'receivedBytes'},
    const {'1': 'address', '3': 16, '4': 1, '5': 9, '10': 'address'},
    const {'1': 'direction', '3': 17, '4': 1, '5': 9, '10': 'direction'},
    const {'1': 'protocols', '3': 18, '4': 3, '5': 9, '10': 'protocols'},
    const {'1': 'total_sessions', '3': 19, '4': 1, '5': 5, '10': 'totalSessions'},
    const {'1': 'completed_sessions', '3': 20, '4': 1, '5': 5, '10': 'completedSessions'},
  ],
  '3': const [PeerInfo_SentBytesEntry$json, PeerInfo_ReceivedBytesEntry$json],
};

@$core.Deprecated('Use peerInfoDescriptor instead')
const PeerInfo_SentBytesEntry$json = const {
  '1': 'SentBytesEntry',
  '2': const [
    const {'1': 'key', '3': 1, '4': 1, '5': 5, '10': 'key'},
    const {'1': 'value', '3': 2, '4': 1, '5': 3, '10': 'value'},
  ],
  '7': const {'7': true},
};

@$core.Deprecated('Use peerInfoDescriptor instead')
const PeerInfo_ReceivedBytesEntry$json = const {
  '1': 'ReceivedBytesEntry',
  '2': const [
    const {'1': 'key', '3': 1, '4': 1, '5': 5, '10': 'key'},
    const {'1': 'value', '3': 2, '4': 1, '5': 3, '10': 'value'},
  ],
  '7': const {'7': true},
};

/// Descriptor for `PeerInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List peerInfoDescriptor = $convert.base64Decode('CghQZWVySW5mbxIWCgZzdGF0dXMYASABKAVSBnN0YXR1cxIYCgdtb25pa2VyGAIgASgJUgdtb25pa2VyEhQKBWFnZW50GAMgASgJUgVhZ2VudBIXCgdwZWVyX2lkGAQgASgJUgZwZWVySWQSJQoOY29uc2Vuc3VzX2tleXMYBSADKAlSDWNvbnNlbnN1c0tleXMSLwoTY29uc2Vuc3VzX2FkZHJlc3NlcxgGIAMoCVISY29uc2Vuc3VzQWRkcmVzc2VzEhoKCHNlcnZpY2VzGAcgASgNUghzZXJ2aWNlcxImCg9sYXN0X2Jsb2NrX2hhc2gYCCABKAlSDWxhc3RCbG9ja0hhc2gSFgoGaGVpZ2h0GAkgASgNUgZoZWlnaHQSKQoQcmVjZWl2ZWRfYnVuZGxlcxgKIAEoBVIPcmVjZWl2ZWRCdW5kbGVzEicKD2ludmFsaWRfYnVuZGxlcxgLIAEoBVIOaW52YWxpZEJ1bmRsZXMSGwoJbGFzdF9zZW50GAwgASgDUghsYXN0U2VudBIjCg1sYXN0X3JlY2VpdmVkGA0gASgDUgxsYXN0UmVjZWl2ZWQSPgoKc2VudF9ieXRlcxgOIAMoCzIfLnBhY3R1cy5QZWVySW5mby5TZW50Qnl0ZXNFbnRyeVIJc2VudEJ5dGVzEkoKDnJlY2VpdmVkX2J5dGVzGA8gAygLMiMucGFjdHVzLlBlZXJJbmZvLlJlY2VpdmVkQnl0ZXNFbnRyeVINcmVjZWl2ZWRCeXRlcxIYCgdhZGRyZXNzGBAgASgJUgdhZGRyZXNzEhwKCWRpcmVjdGlvbhgRIAEoCVIJZGlyZWN0aW9uEhwKCXByb3RvY29scxgSIAMoCVIJcHJvdG9jb2xzEiUKDnRvdGFsX3Nlc3Npb25zGBMgASgFUg10b3RhbFNlc3Npb25zEi0KEmNvbXBsZXRlZF9zZXNzaW9ucxgUIAEoBVIRY29tcGxldGVkU2Vzc2lvbnMaPAoOU2VudEJ5dGVzRW50cnkSEAoDa2V5GAEgASgFUgNrZXkSFAoFdmFsdWUYAiABKANSBXZhbHVlOgI4ARpAChJSZWNlaXZlZEJ5dGVzRW50cnkSEAoDa2V5GAEgASgFUgNrZXkSFAoFdmFsdWUYAiABKANSBXZhbHVlOgI4AQ==');
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
  '.pactus.PeerInfo.SentBytesEntry': PeerInfo_SentBytesEntry$json,
  '.pactus.PeerInfo.ReceivedBytesEntry': PeerInfo_ReceivedBytesEntry$json,
  '.pactus.GetNetworkInfoResponse.SentBytesEntry': GetNetworkInfoResponse_SentBytesEntry$json,
  '.pactus.GetNetworkInfoResponse.ReceivedBytesEntry': GetNetworkInfoResponse_ReceivedBytesEntry$json,
  '.pactus.GetNodeInfoRequest': GetNodeInfoRequest$json,
  '.pactus.GetNodeInfoResponse': GetNodeInfoResponse$json,
  '.pactus.ConnectionInfo': ConnectionInfo$json,
};

/// Descriptor for `Network`. Decode as a `google.protobuf.ServiceDescriptorProto`.
final $typed_data.Uint8List networkServiceDescriptor = $convert.base64Decode('CgdOZXR3b3JrEk8KDkdldE5ldHdvcmtJbmZvEh0ucGFjdHVzLkdldE5ldHdvcmtJbmZvUmVxdWVzdBoeLnBhY3R1cy5HZXROZXR3b3JrSW5mb1Jlc3BvbnNlEkYKC0dldE5vZGVJbmZvEhoucGFjdHVzLkdldE5vZGVJbmZvUmVxdWVzdBobLnBhY3R1cy5HZXROb2RlSW5mb1Jlc3BvbnNl');
