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
    const {'1': 'self_id', '3': 1, '4': 1, '5': 12, '10': 'selfId'},
    const {'1': 'peers', '3': 2, '4': 3, '5': 11, '6': '.pactus.PeerInfo', '10': 'peers'},
  ],
};

/// Descriptor for `GetNetworkInfoResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getNetworkInfoResponseDescriptor = $convert.base64Decode('ChZHZXROZXR3b3JrSW5mb1Jlc3BvbnNlEhcKB3NlbGZfaWQYASABKAxSBnNlbGZJZBImCgVwZWVycxgCIAMoCzIQLnBhY3R1cy5QZWVySW5mb1IFcGVlcnM=');
@$core.Deprecated('Use getPeerInfoRequestDescriptor instead')
const GetPeerInfoRequest$json = const {
  '1': 'GetPeerInfoRequest',
};

/// Descriptor for `GetPeerInfoRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getPeerInfoRequestDescriptor = $convert.base64Decode('ChJHZXRQZWVySW5mb1JlcXVlc3Q=');
@$core.Deprecated('Use getPeerInfoResponseDescriptor instead')
const GetPeerInfoResponse$json = const {
  '1': 'GetPeerInfoResponse',
  '2': const [
    const {'1': 'peer', '3': 1, '4': 1, '5': 11, '6': '.pactus.PeerInfo', '10': 'peer'},
  ],
};

/// Descriptor for `GetPeerInfoResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getPeerInfoResponseDescriptor = $convert.base64Decode('ChNHZXRQZWVySW5mb1Jlc3BvbnNlEiQKBHBlZXIYASABKAsyEC5wYWN0dXMuUGVlckluZm9SBHBlZXI=');
@$core.Deprecated('Use peerInfoDescriptor instead')
const PeerInfo$json = const {
  '1': 'PeerInfo',
  '2': const [
    const {'1': 'moniker', '3': 1, '4': 1, '5': 9, '10': 'moniker'},
    const {'1': 'agent', '3': 2, '4': 1, '5': 9, '10': 'agent'},
    const {'1': 'peer_id', '3': 3, '4': 1, '5': 12, '10': 'peerId'},
    const {'1': 'keys', '3': 4, '4': 3, '5': 9, '10': 'keys'},
    const {'1': 'flags', '3': 5, '4': 1, '5': 5, '10': 'flags'},
    const {'1': 'height', '3': 6, '4': 1, '5': 13, '10': 'height'},
    const {'1': 'received_messages', '3': 7, '4': 1, '5': 5, '10': 'receivedMessages'},
    const {'1': 'invalid_messages', '3': 8, '4': 1, '5': 5, '10': 'invalidMessages'},
    const {'1': 'received_bytes', '3': 9, '4': 1, '5': 5, '10': 'receivedBytes'},
    const {'1': 'status', '3': 10, '4': 1, '5': 5, '10': 'status'},
    const {'1': 'last_seen', '3': 11, '4': 1, '5': 3, '10': 'lastSeen'},
  ],
};

/// Descriptor for `PeerInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List peerInfoDescriptor = $convert.base64Decode('CghQZWVySW5mbxIYCgdtb25pa2VyGAEgASgJUgdtb25pa2VyEhQKBWFnZW50GAIgASgJUgVhZ2VudBIXCgdwZWVyX2lkGAMgASgMUgZwZWVySWQSEgoEa2V5cxgEIAMoCVIEa2V5cxIUCgVmbGFncxgFIAEoBVIFZmxhZ3MSFgoGaGVpZ2h0GAYgASgNUgZoZWlnaHQSKwoRcmVjZWl2ZWRfbWVzc2FnZXMYByABKAVSEHJlY2VpdmVkTWVzc2FnZXMSKQoQaW52YWxpZF9tZXNzYWdlcxgIIAEoBVIPaW52YWxpZE1lc3NhZ2VzEiUKDnJlY2VpdmVkX2J5dGVzGAkgASgFUg1yZWNlaXZlZEJ5dGVzEhYKBnN0YXR1cxgKIAEoBVIGc3RhdHVzEhsKCWxhc3Rfc2VlbhgLIAEoA1IIbGFzdFNlZW4=');
const $core.Map<$core.String, $core.dynamic> NetworkServiceBase$json = const {
  '1': 'Network',
  '2': const [
    const {'1': 'GetNetworkInfo', '2': '.pactus.GetNetworkInfoRequest', '3': '.pactus.GetNetworkInfoResponse'},
    const {'1': 'GetPeerInfo', '2': '.pactus.GetPeerInfoRequest', '3': '.pactus.GetPeerInfoResponse'},
  ],
};

@$core.Deprecated('Use networkServiceDescriptor instead')
const $core.Map<$core.String, $core.Map<$core.String, $core.dynamic>> NetworkServiceBase$messageJson = const {
  '.pactus.GetNetworkInfoRequest': GetNetworkInfoRequest$json,
  '.pactus.GetNetworkInfoResponse': GetNetworkInfoResponse$json,
  '.pactus.PeerInfo': PeerInfo$json,
  '.pactus.GetPeerInfoRequest': GetPeerInfoRequest$json,
  '.pactus.GetPeerInfoResponse': GetPeerInfoResponse$json,
};

/// Descriptor for `Network`. Decode as a `google.protobuf.ServiceDescriptorProto`.
final $typed_data.Uint8List networkServiceDescriptor = $convert.base64Decode('CgdOZXR3b3JrEk8KDkdldE5ldHdvcmtJbmZvEh0ucGFjdHVzLkdldE5ldHdvcmtJbmZvUmVxdWVzdBoeLnBhY3R1cy5HZXROZXR3b3JrSW5mb1Jlc3BvbnNlEkYKC0dldFBlZXJJbmZvEhoucGFjdHVzLkdldFBlZXJJbmZvUmVxdWVzdBobLnBhY3R1cy5HZXRQZWVySW5mb1Jlc3BvbnNl');
