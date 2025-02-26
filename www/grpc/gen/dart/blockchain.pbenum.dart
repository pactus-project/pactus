///
//  Generated code. Do not modify.
//  source: blockchain.proto
//
// @dart = 2.12
// ignore_for_file: annotate_overrides,camel_case_types,constant_identifier_names,directives_ordering,library_prefixes,non_constant_identifier_names,prefer_final_fields,return_of_invalid_type,unnecessary_const,unnecessary_import,unnecessary_this,unused_import,unused_shown_name

// ignore_for_file: UNDEFINED_SHOWN_NAME
import 'dart:core' as $core;
import 'package:protobuf/protobuf.dart' as $pb;

class BlockVerbosity extends $pb.ProtobufEnum {
  static const BlockVerbosity BLOCK_DATA = BlockVerbosity._(0, const $core.bool.fromEnvironment('protobuf.omit_enum_names') ? '' : 'BLOCK_DATA');
  static const BlockVerbosity BLOCK_INFO = BlockVerbosity._(1, const $core.bool.fromEnvironment('protobuf.omit_enum_names') ? '' : 'BLOCK_INFO');
  static const BlockVerbosity BLOCK_TRANSACTIONS = BlockVerbosity._(2, const $core.bool.fromEnvironment('protobuf.omit_enum_names') ? '' : 'BLOCK_TRANSACTIONS');

  static const $core.List<BlockVerbosity> values = <BlockVerbosity> [
    BLOCK_DATA,
    BLOCK_INFO,
    BLOCK_TRANSACTIONS,
  ];

  static final $core.Map<$core.int, BlockVerbosity> _byValue = $pb.ProtobufEnum.initByValue(values);
  static BlockVerbosity? valueOf($core.int value) => _byValue[value];

  const BlockVerbosity._($core.int v, $core.String n) : super(v, n);
}

class VoteType extends $pb.ProtobufEnum {
  static const VoteType VOTE_UNKNOWN = VoteType._(0, const $core.bool.fromEnvironment('protobuf.omit_enum_names') ? '' : 'VOTE_UNKNOWN');
  static const VoteType VOTE_PREPARE = VoteType._(1, const $core.bool.fromEnvironment('protobuf.omit_enum_names') ? '' : 'VOTE_PREPARE');
  static const VoteType VOTE_PRECOMMIT = VoteType._(2, const $core.bool.fromEnvironment('protobuf.omit_enum_names') ? '' : 'VOTE_PRECOMMIT');
  static const VoteType VOTE_CP_PRE_VOTE = VoteType._(3, const $core.bool.fromEnvironment('protobuf.omit_enum_names') ? '' : 'VOTE_CP_PRE_VOTE');
  static const VoteType VOTE_CP_MAIN_VOTE = VoteType._(4, const $core.bool.fromEnvironment('protobuf.omit_enum_names') ? '' : 'VOTE_CP_MAIN_VOTE');
  static const VoteType VOTE_CP_DECIDED = VoteType._(5, const $core.bool.fromEnvironment('protobuf.omit_enum_names') ? '' : 'VOTE_CP_DECIDED');

  static const $core.List<VoteType> values = <VoteType> [
    VOTE_UNKNOWN,
    VOTE_PREPARE,
    VOTE_PRECOMMIT,
    VOTE_CP_PRE_VOTE,
    VOTE_CP_MAIN_VOTE,
    VOTE_CP_DECIDED,
  ];

  static final $core.Map<$core.int, VoteType> _byValue = $pb.ProtobufEnum.initByValue(values);
  static VoteType? valueOf($core.int value) => _byValue[value];

  const VoteType._($core.int v, $core.String n) : super(v, n);
}

