//
//  Generated code. Do not modify.
//  source: blockchain.proto
//
// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_final_fields
// ignore_for_file: unnecessary_import, unnecessary_this, unused_import

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

/// Enumeration for verbosity levels when requesting block information.
class BlockVerbosity extends $pb.ProtobufEnum {
  /// Request only block data.
  /// buf:lint:ignore ENUM_ZERO_VALUE_SUFFIX
  static const BlockVerbosity BLOCK_VERBOSITY_DATA = BlockVerbosity._(0, _omitEnumNames ? '' : 'BLOCK_VERBOSITY_DATA');
  /// Request block information and transaction IDs.
  static const BlockVerbosity BLOCK_VERBOSITY_INFO = BlockVerbosity._(1, _omitEnumNames ? '' : 'BLOCK_VERBOSITY_INFO');
  /// Request block information and detailed transaction data.
  static const BlockVerbosity BLOCK_VERBOSITY_TRANSACTIONS = BlockVerbosity._(2, _omitEnumNames ? '' : 'BLOCK_VERBOSITY_TRANSACTIONS');

  static const $core.List<BlockVerbosity> values = <BlockVerbosity> [
    BLOCK_VERBOSITY_DATA,
    BLOCK_VERBOSITY_INFO,
    BLOCK_VERBOSITY_TRANSACTIONS,
  ];

  static final $core.Map<$core.int, BlockVerbosity> _byValue = $pb.ProtobufEnum.initByValue(values);
  static BlockVerbosity? valueOf($core.int value) => _byValue[value];

  const BlockVerbosity._(super.v, super.n);
}

/// Enumeration for types of votes.
class VoteType extends $pb.ProtobufEnum {
  /// Unspecified vote type.
  static const VoteType VOTE_TYPE_UNSPECIFIED = VoteType._(0, _omitEnumNames ? '' : 'VOTE_TYPE_UNSPECIFIED');
  /// Prepare vote type.
  static const VoteType VOTE_TYPE_PREPARE = VoteType._(1, _omitEnumNames ? '' : 'VOTE_TYPE_PREPARE');
  /// Precommit vote type.
  static const VoteType VOTE_TYPE_PRECOMMIT = VoteType._(2, _omitEnumNames ? '' : 'VOTE_TYPE_PRECOMMIT');
  /// Change-proposer:pre-vote vote type.
  static const VoteType VOTE_TYPE_CP_PRE_VOTE = VoteType._(3, _omitEnumNames ? '' : 'VOTE_TYPE_CP_PRE_VOTE');
  /// Change-proposer:main-vote vote type.
  static const VoteType VOTE_TYPE_CP_MAIN_VOTE = VoteType._(4, _omitEnumNames ? '' : 'VOTE_TYPE_CP_MAIN_VOTE');
  /// Change-proposer:decided vote type.
  static const VoteType VOTE_TYPE_CP_DECIDED = VoteType._(5, _omitEnumNames ? '' : 'VOTE_TYPE_CP_DECIDED');

  static const $core.List<VoteType> values = <VoteType> [
    VOTE_TYPE_UNSPECIFIED,
    VOTE_TYPE_PREPARE,
    VOTE_TYPE_PRECOMMIT,
    VOTE_TYPE_CP_PRE_VOTE,
    VOTE_TYPE_CP_MAIN_VOTE,
    VOTE_TYPE_CP_DECIDED,
  ];

  static final $core.Map<$core.int, VoteType> _byValue = $pb.ProtobufEnum.initByValue(values);
  static VoteType? valueOf($core.int value) => _byValue[value];

  const VoteType._(super.v, super.n);
}


const _omitEnumNames = $core.bool.fromEnvironment('protobuf.omit_enum_names');
