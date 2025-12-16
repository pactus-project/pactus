// This is a generated file - do not edit.
//
// Generated from network.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

/// Direction represents the connection direction between peers.
class Direction extends $pb.ProtobufEnum {
  /// Unknown direction (default value).
  static const Direction DIRECTION_UNKNOWN =
      Direction._(0, _omitEnumNames ? '' : 'DIRECTION_UNKNOWN');

  /// Inbound connection - peer connected to us.
  static const Direction DIRECTION_INBOUND =
      Direction._(1, _omitEnumNames ? '' : 'DIRECTION_INBOUND');

  /// Outbound connection - we connected to peer.
  static const Direction DIRECTION_OUTBOUND =
      Direction._(2, _omitEnumNames ? '' : 'DIRECTION_OUTBOUND');

  static const $core.List<Direction> values = <Direction>[
    DIRECTION_UNKNOWN,
    DIRECTION_INBOUND,
    DIRECTION_OUTBOUND,
  ];

  static final $core.List<Direction?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static Direction? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const Direction._(super.value, super.name);
}

const $core.bool _omitEnumNames =
    $core.bool.fromEnvironment('protobuf.omit_enum_names');
