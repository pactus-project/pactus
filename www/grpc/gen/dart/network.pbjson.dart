///
//  Generated code. Do not modify.
//  source: network.proto
//
// @dart = 2.12
// ignore_for_file: annotate_overrides,camel_case_types,constant_identifier_names,deprecated_member_use_from_same_package,directives_ordering,library_prefixes,non_constant_identifier_names,prefer_final_fields,return_of_invalid_type,unnecessary_const,unnecessary_import,unnecessary_this,unused_import,unused_shown_name

import 'dart:core' as $core;
import 'dart:convert' as $convert;
import 'dart:typed_data' as $typed_data;
@$core.Deprecated('Use networkInfoRequestDescriptor instead')
const NetworkInfoRequest$json = const {
  '1': 'NetworkInfoRequest',
};

/// Descriptor for `NetworkInfoRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List networkInfoRequestDescriptor = $convert.base64Decode('ChJOZXR3b3JrSW5mb1JlcXVlc3Q=');
@$core.Deprecated('Use networkInfoResponseDescriptor instead')
const NetworkInfoResponse$json = const {
  '1': 'NetworkInfoResponse',
  '2': const [
    const {'1': 'self_id', '3': 1, '4': 1, '5': 12, '10': 'selfId'},
    const {'1': 'peers', '3': 2, '4': 3, '5': 11, '6': '.pactus.PeerInfo', '10': 'peers'},
  ],
};

/// Descriptor for `NetworkInfoResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List networkInfoResponseDescriptor = $convert.base64Decode('ChNOZXR3b3JrSW5mb1Jlc3BvbnNlEhcKB3NlbGZfaWQYASABKAxSBnNlbGZJZBImCgVwZWVycxgCIAMoCzIQLnBhY3R1cy5QZWVySW5mb1IFcGVlcnM=');
@$core.Deprecated('Use peerInfoDescriptor instead')
const PeerInfo$json = const {
  '1': 'PeerInfo',
  '2': const [
    const {'1': 'moniker', '3': 1, '4': 1, '5': 9, '10': 'moniker'},
    const {'1': 'agent', '3': 2, '4': 1, '5': 9, '10': 'agent'},
    const {'1': 'peer_id', '3': 3, '4': 1, '5': 12, '10': 'peerId'},
    const {'1': 'public_key', '3': 4, '4': 1, '5': 9, '10': 'publicKey'},
    const {'1': 'flags', '3': 5, '4': 1, '5': 5, '10': 'flags'},
    const {'1': 'height', '3': 6, '4': 1, '5': 13, '10': 'height'},
    const {'1': 'received_messages', '3': 7, '4': 1, '5': 5, '10': 'receivedMessages'},
    const {'1': 'invalid_messages', '3': 8, '4': 1, '5': 5, '10': 'invalidMessages'},
    const {'1': 'received_bytes', '3': 9, '4': 1, '5': 5, '10': 'receivedBytes'},
  ],
};

/// Descriptor for `PeerInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List peerInfoDescriptor = $convert.base64Decode('CghQZWVySW5mbxIYCgdtb25pa2VyGAEgASgJUgdtb25pa2VyEhQKBWFnZW50GAIgASgJUgVhZ2VudBIXCgdwZWVyX2lkGAMgASgMUgZwZWVySWQSHQoKcHVibGljX2tleRgEIAEoCVIJcHVibGljS2V5EhQKBWZsYWdzGAUgASgFUgVmbGFncxIWCgZoZWlnaHQYBiABKA1SBmhlaWdodBIrChFyZWNlaXZlZF9tZXNzYWdlcxgHIAEoBVIQcmVjZWl2ZWRNZXNzYWdlcxIpChBpbnZhbGlkX21lc3NhZ2VzGAggASgFUg9pbnZhbGlkTWVzc2FnZXMSJQoOcmVjZWl2ZWRfYnl0ZXMYCSABKAVSDXJlY2VpdmVkQnl0ZXM=');
@$core.Deprecated('Use peerInfoRequestDescriptor instead')
const PeerInfoRequest$json = const {
  '1': 'PeerInfoRequest',
};

/// Descriptor for `PeerInfoRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List peerInfoRequestDescriptor = $convert.base64Decode('Cg9QZWVySW5mb1JlcXVlc3Q=');
@$core.Deprecated('Use peerInfoResponseDescriptor instead')
const PeerInfoResponse$json = const {
  '1': 'PeerInfoResponse',
  '2': const [
    const {'1': 'moniker', '3': 1, '4': 1, '5': 9, '10': 'moniker'},
    const {'1': 'agent', '3': 2, '4': 1, '5': 9, '10': 'agent'},
    const {'1': 'peer_id', '3': 3, '4': 1, '5': 12, '10': 'peerId'},
    const {'1': 'public_key', '3': 4, '4': 1, '5': 9, '10': 'publicKey'},
    const {'1': 'flags', '3': 5, '4': 1, '5': 5, '10': 'flags'},
    const {'1': 'height', '3': 6, '4': 1, '5': 13, '10': 'height'},
    const {'1': 'received_messages', '3': 7, '4': 1, '5': 5, '10': 'receivedMessages'},
    const {'1': 'invalid_messages', '3': 8, '4': 1, '5': 5, '10': 'invalidMessages'},
    const {'1': 'received_bytes', '3': 9, '4': 1, '5': 5, '10': 'receivedBytes'},
  ],
};

/// Descriptor for `PeerInfoResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List peerInfoResponseDescriptor = $convert.base64Decode('ChBQZWVySW5mb1Jlc3BvbnNlEhgKB21vbmlrZXIYASABKAlSB21vbmlrZXISFAoFYWdlbnQYAiABKAlSBWFnZW50EhcKB3BlZXJfaWQYAyABKAxSBnBlZXJJZBIdCgpwdWJsaWNfa2V5GAQgASgJUglwdWJsaWNLZXkSFAoFZmxhZ3MYBSABKAVSBWZsYWdzEhYKBmhlaWdodBgGIAEoDVIGaGVpZ2h0EisKEXJlY2VpdmVkX21lc3NhZ2VzGAcgASgFUhByZWNlaXZlZE1lc3NhZ2VzEikKEGludmFsaWRfbWVzc2FnZXMYCCABKAVSD2ludmFsaWRNZXNzYWdlcxIlCg5yZWNlaXZlZF9ieXRlcxgJIAEoBVINcmVjZWl2ZWRCeXRlcw==');
const $core.Map<$core.String, $core.dynamic> NetworkServiceBase$json = const {
  '1': 'Network',
  '2': const [
    const {'1': 'GetNetworkInfo', '2': '.pactus.NetworkInfoRequest', '3': '.pactus.NetworkInfoResponse'},
    const {'1': 'GetPeerInfo', '2': '.pactus.PeerInfoRequest', '3': '.pactus.PeerInfoResponse'},
  ],
};

@$core.Deprecated('Use networkServiceDescriptor instead')
const $core.Map<$core.String, $core.Map<$core.String, $core.dynamic>> NetworkServiceBase$messageJson = const {
  '.pactus.NetworkInfoRequest': NetworkInfoRequest$json,
  '.pactus.NetworkInfoResponse': NetworkInfoResponse$json,
  '.pactus.PeerInfo': PeerInfo$json,
  '.pactus.PeerInfoRequest': PeerInfoRequest$json,
  '.pactus.PeerInfoResponse': PeerInfoResponse$json,
};

/// Descriptor for `Network`. Decode as a `google.protobuf.ServiceDescriptorProto`.
final $typed_data.Uint8List networkServiceDescriptor = $convert.base64Decode('CgdOZXR3b3JrEkkKDkdldE5ldHdvcmtJbmZvEhoucGFjdHVzLk5ldHdvcmtJbmZvUmVxdWVzdBobLnBhY3R1cy5OZXR3b3JrSW5mb1Jlc3BvbnNlEkAKC0dldFBlZXJJbmZvEhcucGFjdHVzLlBlZXJJbmZvUmVxdWVzdBoYLnBhY3R1cy5QZWVySW5mb1Jlc3BvbnNl');
