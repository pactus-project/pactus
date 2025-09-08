//
//  Generated code. Do not modify.
//  source: network.proto
//
// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_final_fields
// ignore_for_file: unnecessary_import, unnecessary_this, unused_import

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

/// Direction represents the connection direction between peers.
class Direction extends $pb.ProtobufEnum {
  /// Unknown direction (default value).
  static const Direction DIRECTION_UNKNOWN = Direction._(0, _omitEnumNames ? '' : 'DIRECTION_UNKNOWN');
  /// Inbound connection - peer connected to us.
  static const Direction DIRECTION_INBOUND = Direction._(1, _omitEnumNames ? '' : 'DIRECTION_INBOUND');
  /// Outbound connection - we connected to peer.
  static const Direction DIRECTION_OUTBOUND = Direction._(2, _omitEnumNames ? '' : 'DIRECTION_OUTBOUND');

  static const $core.List<Direction> values = <Direction> [
    DIRECTION_UNKNOWN,
    DIRECTION_INBOUND,
    DIRECTION_OUTBOUND,
  ];

  static final $core.Map<$core.int, Direction> _byValue = $pb.ProtobufEnum.initByValue(values);
  static Direction? valueOf($core.int value) => _byValue[value];

  const Direction._(super.v, super.n);
}


const _omitEnumNames = $core.bool.fromEnvironment('protobuf.omit_enum_names');
