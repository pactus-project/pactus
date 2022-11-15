///
//  Generated code. Do not modify.
//  source: blockchain.proto
//
// @dart = 2.12
// ignore_for_file: annotate_overrides,camel_case_types,constant_identifier_names,deprecated_member_use_from_same_package,directives_ordering,library_prefixes,non_constant_identifier_names,prefer_final_fields,return_of_invalid_type,unnecessary_const,unnecessary_import,unnecessary_this,unused_import,unused_shown_name

import 'dart:core' as $core;
import 'dart:convert' as $convert;
import 'dart:typed_data' as $typed_data;
import 'transaction.pbjson.dart' as $0;

@$core.Deprecated('Use blockVerbosityDescriptor instead')
const BlockVerbosity$json = const {
  '1': 'BlockVerbosity',
  '2': const [
    const {'1': 'BLOCK_DATA', '2': 0},
    const {'1': 'BLOCK_INFO', '2': 1},
    const {'1': 'BLOCK_TRANSACTIONS', '2': 2},
  ],
};

/// Descriptor for `BlockVerbosity`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List blockVerbosityDescriptor = $convert.base64Decode('Cg5CbG9ja1ZlcmJvc2l0eRIOCgpCTE9DS19EQVRBEAASDgoKQkxPQ0tfSU5GTxABEhYKEkJMT0NLX1RSQU5TQUNUSU9OUxAC');
@$core.Deprecated('Use voteTypeDescriptor instead')
const VoteType$json = const {
  '1': 'VoteType',
  '2': const [
    const {'1': 'VOTE_UNKNOWN', '2': 0},
    const {'1': 'VOTE_PREPARE', '2': 1},
    const {'1': 'VOTE_PRECOMMIT', '2': 2},
    const {'1': 'VOTE_CHANGE_PROPOSER', '2': 3},
  ],
};

/// Descriptor for `VoteType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List voteTypeDescriptor = $convert.base64Decode('CghWb3RlVHlwZRIQCgxWT1RFX1VOS05PV04QABIQCgxWT1RFX1BSRVBBUkUQARISCg5WT1RFX1BSRUNPTU1JVBACEhgKFFZPVEVfQ0hBTkdFX1BST1BPU0VSEAM=');
@$core.Deprecated('Use getAccountRequestDescriptor instead')
const GetAccountRequest$json = const {
  '1': 'GetAccountRequest',
  '2': const [
    const {'1': 'address', '3': 1, '4': 1, '5': 9, '10': 'address'},
  ],
};

/// Descriptor for `GetAccountRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getAccountRequestDescriptor = $convert.base64Decode('ChFHZXRBY2NvdW50UmVxdWVzdBIYCgdhZGRyZXNzGAEgASgJUgdhZGRyZXNz');
@$core.Deprecated('Use getAccountResponseDescriptor instead')
const GetAccountResponse$json = const {
  '1': 'GetAccountResponse',
  '2': const [
    const {'1': 'account', '3': 1, '4': 1, '5': 11, '6': '.pactus.AccountInfo', '10': 'account'},
  ],
};

/// Descriptor for `GetAccountResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getAccountResponseDescriptor = $convert.base64Decode('ChJHZXRBY2NvdW50UmVzcG9uc2USLQoHYWNjb3VudBgBIAEoCzITLnBhY3R1cy5BY2NvdW50SW5mb1IHYWNjb3VudA==');
@$core.Deprecated('Use getValidatorsRequestDescriptor instead')
const GetValidatorsRequest$json = const {
  '1': 'GetValidatorsRequest',
};

/// Descriptor for `GetValidatorsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getValidatorsRequestDescriptor = $convert.base64Decode('ChRHZXRWYWxpZGF0b3JzUmVxdWVzdA==');
@$core.Deprecated('Use getValidatorRequestDescriptor instead')
const GetValidatorRequest$json = const {
  '1': 'GetValidatorRequest',
  '2': const [
    const {'1': 'address', '3': 1, '4': 1, '5': 9, '10': 'address'},
  ],
};

/// Descriptor for `GetValidatorRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getValidatorRequestDescriptor = $convert.base64Decode('ChNHZXRWYWxpZGF0b3JSZXF1ZXN0EhgKB2FkZHJlc3MYASABKAlSB2FkZHJlc3M=');
@$core.Deprecated('Use getValidatorByNumberRequestDescriptor instead')
const GetValidatorByNumberRequest$json = const {
  '1': 'GetValidatorByNumberRequest',
  '2': const [
    const {'1': 'number', '3': 1, '4': 1, '5': 5, '10': 'number'},
  ],
};

/// Descriptor for `GetValidatorByNumberRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getValidatorByNumberRequestDescriptor = $convert.base64Decode('ChtHZXRWYWxpZGF0b3JCeU51bWJlclJlcXVlc3QSFgoGbnVtYmVyGAEgASgFUgZudW1iZXI=');
@$core.Deprecated('Use getValidatorsResponseDescriptor instead')
const GetValidatorsResponse$json = const {
  '1': 'GetValidatorsResponse',
  '2': const [
    const {'1': 'validators', '3': 1, '4': 3, '5': 11, '6': '.pactus.ValidatorInfo', '10': 'validators'},
  ],
};

/// Descriptor for `GetValidatorsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getValidatorsResponseDescriptor = $convert.base64Decode('ChVHZXRWYWxpZGF0b3JzUmVzcG9uc2USNQoKdmFsaWRhdG9ycxgBIAMoCzIVLnBhY3R1cy5WYWxpZGF0b3JJbmZvUgp2YWxpZGF0b3Jz');
@$core.Deprecated('Use getValidatorResponseDescriptor instead')
const GetValidatorResponse$json = const {
  '1': 'GetValidatorResponse',
  '2': const [
    const {'1': 'validator', '3': 1, '4': 1, '5': 11, '6': '.pactus.ValidatorInfo', '10': 'validator'},
  ],
};

/// Descriptor for `GetValidatorResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getValidatorResponseDescriptor = $convert.base64Decode('ChRHZXRWYWxpZGF0b3JSZXNwb25zZRIzCgl2YWxpZGF0b3IYASABKAsyFS5wYWN0dXMuVmFsaWRhdG9ySW5mb1IJdmFsaWRhdG9y');
@$core.Deprecated('Use getBlockRequestDescriptor instead')
const GetBlockRequest$json = const {
  '1': 'GetBlockRequest',
  '2': const [
    const {'1': 'height', '3': 1, '4': 1, '5': 13, '10': 'height'},
    const {'1': 'verbosity', '3': 2, '4': 1, '5': 14, '6': '.pactus.BlockVerbosity', '10': 'verbosity'},
  ],
};

/// Descriptor for `GetBlockRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBlockRequestDescriptor = $convert.base64Decode('Cg9HZXRCbG9ja1JlcXVlc3QSFgoGaGVpZ2h0GAEgASgNUgZoZWlnaHQSNAoJdmVyYm9zaXR5GAIgASgOMhYucGFjdHVzLkJsb2NrVmVyYm9zaXR5Ugl2ZXJib3NpdHk=');
@$core.Deprecated('Use getBlockResponseDescriptor instead')
const GetBlockResponse$json = const {
  '1': 'GetBlockResponse',
  '2': const [
    const {'1': 'height', '3': 1, '4': 1, '5': 13, '10': 'height'},
    const {'1': 'hash', '3': 2, '4': 1, '5': 12, '10': 'hash'},
    const {'1': 'data', '3': 3, '4': 1, '5': 12, '10': 'data'},
    const {'1': 'block_time', '3': 4, '4': 1, '5': 13, '10': 'blockTime'},
    const {'1': 'header', '3': 5, '4': 1, '5': 11, '6': '.pactus.BlockHeaderInfo', '10': 'header'},
    const {'1': 'prev_cert', '3': 6, '4': 1, '5': 11, '6': '.pactus.CertificateInfo', '10': 'prevCert'},
    const {'1': 'txs', '3': 7, '4': 3, '5': 11, '6': '.pactus.TransactionInfo', '10': 'txs'},
  ],
};

/// Descriptor for `GetBlockResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBlockResponseDescriptor = $convert.base64Decode('ChBHZXRCbG9ja1Jlc3BvbnNlEhYKBmhlaWdodBgBIAEoDVIGaGVpZ2h0EhIKBGhhc2gYAiABKAxSBGhhc2gSEgoEZGF0YRgDIAEoDFIEZGF0YRIdCgpibG9ja190aW1lGAQgASgNUglibG9ja1RpbWUSLwoGaGVhZGVyGAUgASgLMhcucGFjdHVzLkJsb2NrSGVhZGVySW5mb1IGaGVhZGVyEjQKCXByZXZfY2VydBgGIAEoCzIXLnBhY3R1cy5DZXJ0aWZpY2F0ZUluZm9SCHByZXZDZXJ0EikKA3R4cxgHIAMoCzIXLnBhY3R1cy5UcmFuc2FjdGlvbkluZm9SA3R4cw==');
@$core.Deprecated('Use getBlockHashRequestDescriptor instead')
const GetBlockHashRequest$json = const {
  '1': 'GetBlockHashRequest',
  '2': const [
    const {'1': 'height', '3': 1, '4': 1, '5': 13, '10': 'height'},
  ],
};

/// Descriptor for `GetBlockHashRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBlockHashRequestDescriptor = $convert.base64Decode('ChNHZXRCbG9ja0hhc2hSZXF1ZXN0EhYKBmhlaWdodBgBIAEoDVIGaGVpZ2h0');
@$core.Deprecated('Use getBlockHashResponseDescriptor instead')
const GetBlockHashResponse$json = const {
  '1': 'GetBlockHashResponse',
  '2': const [
    const {'1': 'hash', '3': 1, '4': 1, '5': 12, '10': 'hash'},
  ],
};

/// Descriptor for `GetBlockHashResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBlockHashResponseDescriptor = $convert.base64Decode('ChRHZXRCbG9ja0hhc2hSZXNwb25zZRISCgRoYXNoGAEgASgMUgRoYXNo');
@$core.Deprecated('Use getBlockHeightRequestDescriptor instead')
const GetBlockHeightRequest$json = const {
  '1': 'GetBlockHeightRequest',
  '2': const [
    const {'1': 'hash', '3': 1, '4': 1, '5': 12, '10': 'hash'},
  ],
};

/// Descriptor for `GetBlockHeightRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBlockHeightRequestDescriptor = $convert.base64Decode('ChVHZXRCbG9ja0hlaWdodFJlcXVlc3QSEgoEaGFzaBgBIAEoDFIEaGFzaA==');
@$core.Deprecated('Use getBlockHeightResponseDescriptor instead')
const GetBlockHeightResponse$json = const {
  '1': 'GetBlockHeightResponse',
  '2': const [
    const {'1': 'height', '3': 1, '4': 1, '5': 13, '10': 'height'},
  ],
};

/// Descriptor for `GetBlockHeightResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBlockHeightResponseDescriptor = $convert.base64Decode('ChZHZXRCbG9ja0hlaWdodFJlc3BvbnNlEhYKBmhlaWdodBgBIAEoDVIGaGVpZ2h0');
@$core.Deprecated('Use getBlockchainInfoRequestDescriptor instead')
const GetBlockchainInfoRequest$json = const {
  '1': 'GetBlockchainInfoRequest',
};

/// Descriptor for `GetBlockchainInfoRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBlockchainInfoRequestDescriptor = $convert.base64Decode('ChhHZXRCbG9ja2NoYWluSW5mb1JlcXVlc3Q=');
@$core.Deprecated('Use getBlockchainInfoResponseDescriptor instead')
const GetBlockchainInfoResponse$json = const {
  '1': 'GetBlockchainInfoResponse',
  '2': const [
    const {'1': 'last_block_height', '3': 1, '4': 1, '5': 13, '10': 'lastBlockHeight'},
    const {'1': 'last_block_hash', '3': 2, '4': 1, '5': 12, '10': 'lastBlockHash'},
    const {'1': 'total_power', '3': 3, '4': 1, '5': 3, '10': 'totalPower'},
    const {'1': 'committee_power', '3': 4, '4': 1, '5': 3, '10': 'committeePower'},
    const {'1': 'committee_validators', '3': 5, '4': 3, '5': 11, '6': '.pactus.ValidatorInfo', '10': 'committeeValidators'},
  ],
};

/// Descriptor for `GetBlockchainInfoResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBlockchainInfoResponseDescriptor = $convert.base64Decode('ChlHZXRCbG9ja2NoYWluSW5mb1Jlc3BvbnNlEioKEWxhc3RfYmxvY2tfaGVpZ2h0GAEgASgNUg9sYXN0QmxvY2tIZWlnaHQSJgoPbGFzdF9ibG9ja19oYXNoGAIgASgMUg1sYXN0QmxvY2tIYXNoEh8KC3RvdGFsX3Bvd2VyGAMgASgDUgp0b3RhbFBvd2VyEicKD2NvbW1pdHRlZV9wb3dlchgEIAEoA1IOY29tbWl0dGVlUG93ZXISSAoUY29tbWl0dGVlX3ZhbGlkYXRvcnMYBSADKAsyFS5wYWN0dXMuVmFsaWRhdG9ySW5mb1ITY29tbWl0dGVlVmFsaWRhdG9ycw==');
@$core.Deprecated('Use getConsensusInfoRequestDescriptor instead')
const GetConsensusInfoRequest$json = const {
  '1': 'GetConsensusInfoRequest',
};

/// Descriptor for `GetConsensusInfoRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getConsensusInfoRequestDescriptor = $convert.base64Decode('ChdHZXRDb25zZW5zdXNJbmZvUmVxdWVzdA==');
@$core.Deprecated('Use getConsensusInfoResponseDescriptor instead')
const GetConsensusInfoResponse$json = const {
  '1': 'GetConsensusInfoResponse',
  '2': const [
    const {'1': 'height', '3': 1, '4': 1, '5': 13, '10': 'height'},
    const {'1': 'round', '3': 2, '4': 1, '5': 5, '10': 'round'},
    const {'1': 'votes', '3': 3, '4': 3, '5': 11, '6': '.pactus.VoteInfo', '10': 'votes'},
  ],
};

/// Descriptor for `GetConsensusInfoResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getConsensusInfoResponseDescriptor = $convert.base64Decode('ChhHZXRDb25zZW5zdXNJbmZvUmVzcG9uc2USFgoGaGVpZ2h0GAEgASgNUgZoZWlnaHQSFAoFcm91bmQYAiABKAVSBXJvdW5kEiYKBXZvdGVzGAMgAygLMhAucGFjdHVzLlZvdGVJbmZvUgV2b3Rlcw==');
@$core.Deprecated('Use validatorInfoDescriptor instead')
const ValidatorInfo$json = const {
  '1': 'ValidatorInfo',
  '2': const [
    const {'1': 'hash', '3': 1, '4': 1, '5': 12, '10': 'hash'},
    const {'1': 'data', '3': 2, '4': 1, '5': 12, '10': 'data'},
    const {'1': 'public_key', '3': 3, '4': 1, '5': 9, '10': 'publicKey'},
    const {'1': 'number', '3': 4, '4': 1, '5': 5, '10': 'number'},
    const {'1': 'sequence', '3': 5, '4': 1, '5': 5, '10': 'sequence'},
    const {'1': 'stake', '3': 6, '4': 1, '5': 3, '10': 'stake'},
    const {'1': 'last_bonding_height', '3': 7, '4': 1, '5': 13, '10': 'lastBondingHeight'},
    const {'1': 'last_joined_height', '3': 8, '4': 1, '5': 13, '10': 'lastJoinedHeight'},
    const {'1': 'unbonding_height', '3': 9, '4': 1, '5': 13, '10': 'unbondingHeight'},
    const {'1': 'address', '3': 10, '4': 1, '5': 9, '10': 'address'},
  ],
};

/// Descriptor for `ValidatorInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List validatorInfoDescriptor = $convert.base64Decode('Cg1WYWxpZGF0b3JJbmZvEhIKBGhhc2gYASABKAxSBGhhc2gSEgoEZGF0YRgCIAEoDFIEZGF0YRIdCgpwdWJsaWNfa2V5GAMgASgJUglwdWJsaWNLZXkSFgoGbnVtYmVyGAQgASgFUgZudW1iZXISGgoIc2VxdWVuY2UYBSABKAVSCHNlcXVlbmNlEhQKBXN0YWtlGAYgASgDUgVzdGFrZRIuChNsYXN0X2JvbmRpbmdfaGVpZ2h0GAcgASgNUhFsYXN0Qm9uZGluZ0hlaWdodBIsChJsYXN0X2pvaW5lZF9oZWlnaHQYCCABKA1SEGxhc3RKb2luZWRIZWlnaHQSKQoQdW5ib25kaW5nX2hlaWdodBgJIAEoDVIPdW5ib25kaW5nSGVpZ2h0EhgKB2FkZHJlc3MYCiABKAlSB2FkZHJlc3M=');
@$core.Deprecated('Use accountInfoDescriptor instead')
const AccountInfo$json = const {
  '1': 'AccountInfo',
  '2': const [
    const {'1': 'hash', '3': 1, '4': 1, '5': 12, '10': 'hash'},
    const {'1': 'data', '3': 2, '4': 1, '5': 12, '10': 'data'},
    const {'1': 'address', '3': 3, '4': 1, '5': 9, '10': 'address'},
    const {'1': 'number', '3': 4, '4': 1, '5': 5, '10': 'number'},
    const {'1': 'sequence', '3': 5, '4': 1, '5': 5, '10': 'sequence'},
    const {'1': 'balance', '3': 6, '4': 1, '5': 3, '10': 'balance'},
  ],
};

/// Descriptor for `AccountInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List accountInfoDescriptor = $convert.base64Decode('CgtBY2NvdW50SW5mbxISCgRoYXNoGAEgASgMUgRoYXNoEhIKBGRhdGEYAiABKAxSBGRhdGESGAoHYWRkcmVzcxgDIAEoCVIHYWRkcmVzcxIWCgZudW1iZXIYBCABKAVSBm51bWJlchIaCghzZXF1ZW5jZRgFIAEoBVIIc2VxdWVuY2USGAoHYmFsYW5jZRgGIAEoA1IHYmFsYW5jZQ==');
@$core.Deprecated('Use blockHeaderInfoDescriptor instead')
const BlockHeaderInfo$json = const {
  '1': 'BlockHeaderInfo',
  '2': const [
    const {'1': 'version', '3': 1, '4': 1, '5': 5, '10': 'version'},
    const {'1': 'prev_block_hash', '3': 2, '4': 1, '5': 12, '10': 'prevBlockHash'},
    const {'1': 'state_root', '3': 3, '4': 1, '5': 12, '10': 'stateRoot'},
    const {'1': 'sortition_seed', '3': 4, '4': 1, '5': 12, '10': 'sortitionSeed'},
    const {'1': 'proposer_address', '3': 5, '4': 1, '5': 9, '10': 'proposerAddress'},
  ],
};

/// Descriptor for `BlockHeaderInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List blockHeaderInfoDescriptor = $convert.base64Decode('Cg9CbG9ja0hlYWRlckluZm8SGAoHdmVyc2lvbhgBIAEoBVIHdmVyc2lvbhImCg9wcmV2X2Jsb2NrX2hhc2gYAiABKAxSDXByZXZCbG9ja0hhc2gSHQoKc3RhdGVfcm9vdBgDIAEoDFIJc3RhdGVSb290EiUKDnNvcnRpdGlvbl9zZWVkGAQgASgMUg1zb3J0aXRpb25TZWVkEikKEHByb3Bvc2VyX2FkZHJlc3MYBSABKAlSD3Byb3Bvc2VyQWRkcmVzcw==');
@$core.Deprecated('Use certificateInfoDescriptor instead')
const CertificateInfo$json = const {
  '1': 'CertificateInfo',
  '2': const [
    const {'1': 'hash', '3': 1, '4': 1, '5': 12, '10': 'hash'},
    const {'1': 'round', '3': 2, '4': 1, '5': 5, '10': 'round'},
    const {'1': 'committers', '3': 3, '4': 3, '5': 5, '10': 'committers'},
    const {'1': 'absentees', '3': 4, '4': 3, '5': 5, '10': 'absentees'},
    const {'1': 'signature', '3': 5, '4': 1, '5': 12, '10': 'signature'},
  ],
};

/// Descriptor for `CertificateInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List certificateInfoDescriptor = $convert.base64Decode('Cg9DZXJ0aWZpY2F0ZUluZm8SEgoEaGFzaBgBIAEoDFIEaGFzaBIUCgVyb3VuZBgCIAEoBVIFcm91bmQSHgoKY29tbWl0dGVycxgDIAMoBVIKY29tbWl0dGVycxIcCglhYnNlbnRlZXMYBCADKAVSCWFic2VudGVlcxIcCglzaWduYXR1cmUYBSABKAxSCXNpZ25hdHVyZQ==');
@$core.Deprecated('Use voteInfoDescriptor instead')
const VoteInfo$json = const {
  '1': 'VoteInfo',
  '2': const [
    const {'1': 'type', '3': 1, '4': 1, '5': 14, '6': '.pactus.VoteType', '10': 'type'},
    const {'1': 'voter', '3': 2, '4': 1, '5': 9, '10': 'voter'},
    const {'1': 'block_hash', '3': 3, '4': 1, '5': 12, '10': 'blockHash'},
    const {'1': 'round', '3': 4, '4': 1, '5': 5, '10': 'round'},
  ],
};

/// Descriptor for `VoteInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List voteInfoDescriptor = $convert.base64Decode('CghWb3RlSW5mbxIkCgR0eXBlGAEgASgOMhAucGFjdHVzLlZvdGVUeXBlUgR0eXBlEhQKBXZvdGVyGAIgASgJUgV2b3RlchIdCgpibG9ja19oYXNoGAMgASgMUglibG9ja0hhc2gSFAoFcm91bmQYBCABKAVSBXJvdW5k');
const $core.Map<$core.String, $core.dynamic> BlockchainServiceBase$json = const {
  '1': 'Blockchain',
  '2': const [
    const {'1': 'GetBlock', '2': '.pactus.GetBlockRequest', '3': '.pactus.GetBlockResponse'},
    const {'1': 'GetBlockHash', '2': '.pactus.GetBlockHashRequest', '3': '.pactus.GetBlockHashResponse'},
    const {'1': 'GetBlockHeight', '2': '.pactus.GetBlockHeightRequest', '3': '.pactus.GetBlockHeightResponse'},
    const {'1': 'GetBlockchainInfo', '2': '.pactus.GetBlockchainInfoRequest', '3': '.pactus.GetBlockchainInfoResponse'},
    const {'1': 'GetConsensusInfo', '2': '.pactus.GetConsensusInfoRequest', '3': '.pactus.GetConsensusInfoResponse'},
    const {'1': 'GetAccount', '2': '.pactus.GetAccountRequest', '3': '.pactus.GetAccountResponse'},
    const {'1': 'GetValidator', '2': '.pactus.GetValidatorRequest', '3': '.pactus.GetValidatorResponse'},
    const {'1': 'GetValidatorByNumber', '2': '.pactus.GetValidatorByNumberRequest', '3': '.pactus.GetValidatorResponse'},
    const {'1': 'GetValidators', '2': '.pactus.GetValidatorsRequest', '3': '.pactus.GetValidatorsResponse'},
  ],
};

@$core.Deprecated('Use blockchainServiceDescriptor instead')
const $core.Map<$core.String, $core.Map<$core.String, $core.dynamic>> BlockchainServiceBase$messageJson = const {
  '.pactus.GetBlockRequest': GetBlockRequest$json,
  '.pactus.GetBlockResponse': GetBlockResponse$json,
  '.pactus.BlockHeaderInfo': BlockHeaderInfo$json,
  '.pactus.CertificateInfo': CertificateInfo$json,
  '.pactus.TransactionInfo': $0.TransactionInfo$json,
  '.pactus.PayloadSend': $0.PayloadSend$json,
  '.pactus.PayloadBond': $0.PayloadBond$json,
  '.pactus.PayloadSortition': $0.PayloadSortition$json,
  '.pactus.PayloadUnbond': $0.PayloadUnbond$json,
  '.pactus.PayloadWithdraw': $0.PayloadWithdraw$json,
  '.pactus.GetBlockHashRequest': GetBlockHashRequest$json,
  '.pactus.GetBlockHashResponse': GetBlockHashResponse$json,
  '.pactus.GetBlockHeightRequest': GetBlockHeightRequest$json,
  '.pactus.GetBlockHeightResponse': GetBlockHeightResponse$json,
  '.pactus.GetBlockchainInfoRequest': GetBlockchainInfoRequest$json,
  '.pactus.GetBlockchainInfoResponse': GetBlockchainInfoResponse$json,
  '.pactus.ValidatorInfo': ValidatorInfo$json,
  '.pactus.GetConsensusInfoRequest': GetConsensusInfoRequest$json,
  '.pactus.GetConsensusInfoResponse': GetConsensusInfoResponse$json,
  '.pactus.VoteInfo': VoteInfo$json,
  '.pactus.GetAccountRequest': GetAccountRequest$json,
  '.pactus.GetAccountResponse': GetAccountResponse$json,
  '.pactus.AccountInfo': AccountInfo$json,
  '.pactus.GetValidatorRequest': GetValidatorRequest$json,
  '.pactus.GetValidatorResponse': GetValidatorResponse$json,
  '.pactus.GetValidatorByNumberRequest': GetValidatorByNumberRequest$json,
  '.pactus.GetValidatorsRequest': GetValidatorsRequest$json,
  '.pactus.GetValidatorsResponse': GetValidatorsResponse$json,
};

/// Descriptor for `Blockchain`. Decode as a `google.protobuf.ServiceDescriptorProto`.
final $typed_data.Uint8List blockchainServiceDescriptor = $convert.base64Decode('CgpCbG9ja2NoYWluEj0KCEdldEJsb2NrEhcucGFjdHVzLkdldEJsb2NrUmVxdWVzdBoYLnBhY3R1cy5HZXRCbG9ja1Jlc3BvbnNlEkkKDEdldEJsb2NrSGFzaBIbLnBhY3R1cy5HZXRCbG9ja0hhc2hSZXF1ZXN0GhwucGFjdHVzLkdldEJsb2NrSGFzaFJlc3BvbnNlEk8KDkdldEJsb2NrSGVpZ2h0Eh0ucGFjdHVzLkdldEJsb2NrSGVpZ2h0UmVxdWVzdBoeLnBhY3R1cy5HZXRCbG9ja0hlaWdodFJlc3BvbnNlElgKEUdldEJsb2NrY2hhaW5JbmZvEiAucGFjdHVzLkdldEJsb2NrY2hhaW5JbmZvUmVxdWVzdBohLnBhY3R1cy5HZXRCbG9ja2NoYWluSW5mb1Jlc3BvbnNlElUKEEdldENvbnNlbnN1c0luZm8SHy5wYWN0dXMuR2V0Q29uc2Vuc3VzSW5mb1JlcXVlc3QaIC5wYWN0dXMuR2V0Q29uc2Vuc3VzSW5mb1Jlc3BvbnNlEkMKCkdldEFjY291bnQSGS5wYWN0dXMuR2V0QWNjb3VudFJlcXVlc3QaGi5wYWN0dXMuR2V0QWNjb3VudFJlc3BvbnNlEkkKDEdldFZhbGlkYXRvchIbLnBhY3R1cy5HZXRWYWxpZGF0b3JSZXF1ZXN0GhwucGFjdHVzLkdldFZhbGlkYXRvclJlc3BvbnNlElkKFEdldFZhbGlkYXRvckJ5TnVtYmVyEiMucGFjdHVzLkdldFZhbGlkYXRvckJ5TnVtYmVyUmVxdWVzdBocLnBhY3R1cy5HZXRWYWxpZGF0b3JSZXNwb25zZRJMCg1HZXRWYWxpZGF0b3JzEhwucGFjdHVzLkdldFZhbGlkYXRvcnNSZXF1ZXN0Gh0ucGFjdHVzLkdldFZhbGlkYXRvcnNSZXNwb25zZQ==');
