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
};

/// Descriptor for `GetNetworkInfoRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getNetworkInfoRequestDescriptor = $convert.base64Decode('ChVHZXROZXR3b3JrSW5mb1JlcXVlc3Q=');
@$core.Deprecated('Use getNetworkInfoResponseDescriptor instead')
const GetNetworkInfoResponse$json = const {
  '1': 'GetNetworkInfoResponse',
  '2': const [
    const {'1': 'total_sent_bytes', '3': 1, '4': 1, '5': 5, '10': 'totalSentBytes'},
    const {'1': 'total_received_bytes', '3': 2, '4': 1, '5': 5, '10': 'totalReceivedBytes'},
    const {'1': 'started_at', '3': 3, '4': 1, '5': 3, '10': 'startedAt'},
    const {'1': 'peers', '3': 4, '4': 3, '5': 11, '6': '.pactus.PeerInfo', '10': 'peers'},
    const {'1': 'sent_bytes', '3': 5, '4': 3, '5': 11, '6': '.pactus.GetNetworkInfoResponse.SentBytesEntry', '10': 'sentBytes'},
    const {'1': 'received_bytes', '3': 6, '4': 3, '5': 11, '6': '.pactus.GetNetworkInfoResponse.ReceivedBytesEntry', '10': 'receivedBytes'},
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
final $typed_data.Uint8List getNetworkInfoResponseDescriptor = $convert.base64Decode('ChZHZXROZXR3b3JrSW5mb1Jlc3BvbnNlEigKEHRvdGFsX3NlbnRfYnl0ZXMYASABKAVSDnRvdGFsU2VudEJ5dGVzEjAKFHRvdGFsX3JlY2VpdmVkX2J5dGVzGAIgASgFUhJ0b3RhbFJlY2VpdmVkQnl0ZXMSHQoKc3RhcnRlZF9hdBgDIAEoA1IJc3RhcnRlZEF0EiYKBXBlZXJzGAQgAygLMhAucGFjdHVzLlBlZXJJbmZvUgVwZWVycxJMCgpzZW50X2J5dGVzGAUgAygLMi0ucGFjdHVzLkdldE5ldHdvcmtJbmZvUmVzcG9uc2UuU2VudEJ5dGVzRW50cnlSCXNlbnRCeXRlcxJYCg5yZWNlaXZlZF9ieXRlcxgGIAMoCzIxLnBhY3R1cy5HZXROZXR3b3JrSW5mb1Jlc3BvbnNlLlJlY2VpdmVkQnl0ZXNFbnRyeVINcmVjZWl2ZWRCeXRlcxo8Cg5TZW50Qnl0ZXNFbnRyeRIQCgNrZXkYASABKAVSA2tleRIUCgV2YWx1ZRgCIAEoA1IFdmFsdWU6AjgBGkAKElJlY2VpdmVkQnl0ZXNFbnRyeRIQCgNrZXkYASABKAVSA2tleRIUCgV2YWx1ZRgCIAEoA1IFdmFsdWU6AjgB');
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
    const {'1': 'peer_id', '3': 3, '4': 1, '5': 12, '10': 'peerId'},
  ],
};

/// Descriptor for `GetNodeInfoResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getNodeInfoResponseDescriptor = $convert.base64Decode('ChNHZXROb2RlSW5mb1Jlc3BvbnNlEhgKB21vbmlrZXIYASABKAlSB21vbmlrZXISFAoFYWdlbnQYAiABKAlSBWFnZW50EhcKB3BlZXJfaWQYAyABKAxSBnBlZXJJZA==');
@$core.Deprecated('Use peerInfoDescriptor instead')
const PeerInfo$json = const {
  '1': 'PeerInfo',
  '2': const [
    const {'1': 'status', '3': 1, '4': 1, '5': 5, '10': 'status'},
    const {'1': 'moniker', '3': 2, '4': 1, '5': 9, '10': 'moniker'},
    const {'1': 'agent', '3': 3, '4': 1, '5': 9, '10': 'agent'},
    const {'1': 'peer_id', '3': 4, '4': 1, '5': 12, '10': 'peerId'},
    const {'1': 'consensus_keys', '3': 5, '4': 3, '5': 9, '10': 'consensusKeys'},
    const {'1': 'services', '3': 6, '4': 1, '5': 13, '10': 'services'},
    const {'1': 'last_block_hash', '3': 7, '4': 1, '5': 12, '10': 'lastBlockHash'},
    const {'1': 'height', '3': 8, '4': 1, '5': 13, '10': 'height'},
    const {'1': 'received_messages', '3': 9, '4': 1, '5': 5, '10': 'receivedMessages'},
    const {'1': 'invalid_messages', '3': 10, '4': 1, '5': 5, '10': 'invalidMessages'},
    const {'1': 'last_sent', '3': 11, '4': 1, '5': 3, '10': 'lastSent'},
    const {'1': 'last_received', '3': 12, '4': 1, '5': 3, '10': 'lastReceived'},
    const {'1': 'sent_bytes', '3': 13, '4': 3, '5': 11, '6': '.pactus.PeerInfo.SentBytesEntry', '10': 'sentBytes'},
    const {'1': 'received_bytes', '3': 14, '4': 3, '5': 11, '6': '.pactus.PeerInfo.ReceivedBytesEntry', '10': 'receivedBytes'},
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
final $typed_data.Uint8List peerInfoDescriptor = $convert.base64Decode('CghQZWVySW5mbxIWCgZzdGF0dXMYASABKAVSBnN0YXR1cxIYCgdtb25pa2VyGAIgASgJUgdtb25pa2VyEhQKBWFnZW50GAMgASgJUgVhZ2VudBIXCgdwZWVyX2lkGAQgASgMUgZwZWVySWQSJQoOY29uc2Vuc3VzX2tleXMYBSADKAlSDWNvbnNlbnN1c0tleXMSGgoIc2VydmljZXMYBiABKA1SCHNlcnZpY2VzEiYKD2xhc3RfYmxvY2tfaGFzaBgHIAEoDFINbGFzdEJsb2NrSGFzaBIWCgZoZWlnaHQYCCABKA1SBmhlaWdodBIrChFyZWNlaXZlZF9tZXNzYWdlcxgJIAEoBVIQcmVjZWl2ZWRNZXNzYWdlcxIpChBpbnZhbGlkX21lc3NhZ2VzGAogASgFUg9pbnZhbGlkTWVzc2FnZXMSGwoJbGFzdF9zZW50GAsgASgDUghsYXN0U2VudBIjCg1sYXN0X3JlY2VpdmVkGAwgASgDUgxsYXN0UmVjZWl2ZWQSPgoKc2VudF9ieXRlcxgNIAMoCzIfLnBhY3R1cy5QZWVySW5mby5TZW50Qnl0ZXNFbnRyeVIJc2VudEJ5dGVzEkoKDnJlY2VpdmVkX2J5dGVzGA4gAygLMiMucGFjdHVzLlBlZXJJbmZvLlJlY2VpdmVkQnl0ZXNFbnRyeVINcmVjZWl2ZWRCeXRlcxo8Cg5TZW50Qnl0ZXNFbnRyeRIQCgNrZXkYASABKAVSA2tleRIUCgV2YWx1ZRgCIAEoA1IFdmFsdWU6AjgBGkAKElJlY2VpdmVkQnl0ZXNFbnRyeRIQCgNrZXkYASABKAVSA2tleRIUCgV2YWx1ZRgCIAEoA1IFdmFsdWU6AjgB');
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
};

/// Descriptor for `Network`. Decode as a `google.protobuf.ServiceDescriptorProto`.
final $typed_data.Uint8List networkServiceDescriptor = $convert.base64Decode('CgdOZXR3b3JrEk8KDkdldE5ldHdvcmtJbmZvEh0ucGFjdHVzLkdldE5ldHdvcmtJbmZvUmVxdWVzdBoeLnBhY3R1cy5HZXROZXR3b3JrSW5mb1Jlc3BvbnNlEkYKC0dldE5vZGVJbmZvEhoucGFjdHVzLkdldE5vZGVJbmZvUmVxdWVzdBobLnBhY3R1cy5HZXROb2RlSW5mb1Jlc3BvbnNl');
