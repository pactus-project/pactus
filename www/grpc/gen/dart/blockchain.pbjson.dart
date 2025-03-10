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
    const {'1': 'VOTE_CP_PRE_VOTE', '2': 3},
    const {'1': 'VOTE_CP_MAIN_VOTE', '2': 4},
    const {'1': 'VOTE_CP_DECIDED', '2': 5},
  ],
};

/// Descriptor for `VoteType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List voteTypeDescriptor = $convert.base64Decode('CghWb3RlVHlwZRIQCgxWT1RFX1VOS05PV04QABIQCgxWT1RFX1BSRVBBUkUQARISCg5WT1RFX1BSRUNPTU1JVBACEhQKEFZPVEVfQ1BfUFJFX1ZPVEUQAxIVChFWT1RFX0NQX01BSU5fVk9URRAEEhMKD1ZPVEVfQ1BfREVDSURFRBAF');
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
@$core.Deprecated('Use getValidatorAddressesRequestDescriptor instead')
const GetValidatorAddressesRequest$json = const {
  '1': 'GetValidatorAddressesRequest',
};

/// Descriptor for `GetValidatorAddressesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getValidatorAddressesRequestDescriptor = $convert.base64Decode('ChxHZXRWYWxpZGF0b3JBZGRyZXNzZXNSZXF1ZXN0');
@$core.Deprecated('Use getValidatorAddressesResponseDescriptor instead')
const GetValidatorAddressesResponse$json = const {
  '1': 'GetValidatorAddressesResponse',
  '2': const [
    const {'1': 'addresses', '3': 1, '4': 3, '5': 9, '10': 'addresses'},
  ],
};

/// Descriptor for `GetValidatorAddressesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getValidatorAddressesResponseDescriptor = $convert.base64Decode('Ch1HZXRWYWxpZGF0b3JBZGRyZXNzZXNSZXNwb25zZRIcCglhZGRyZXNzZXMYASADKAlSCWFkZHJlc3Nlcw==');
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
@$core.Deprecated('Use getValidatorResponseDescriptor instead')
const GetValidatorResponse$json = const {
  '1': 'GetValidatorResponse',
  '2': const [
    const {'1': 'validator', '3': 1, '4': 1, '5': 11, '6': '.pactus.ValidatorInfo', '10': 'validator'},
  ],
};

/// Descriptor for `GetValidatorResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getValidatorResponseDescriptor = $convert.base64Decode('ChRHZXRWYWxpZGF0b3JSZXNwb25zZRIzCgl2YWxpZGF0b3IYASABKAsyFS5wYWN0dXMuVmFsaWRhdG9ySW5mb1IJdmFsaWRhdG9y');
@$core.Deprecated('Use getPublicKeyRequestDescriptor instead')
const GetPublicKeyRequest$json = const {
  '1': 'GetPublicKeyRequest',
  '2': const [
    const {'1': 'address', '3': 1, '4': 1, '5': 9, '10': 'address'},
  ],
};

/// Descriptor for `GetPublicKeyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getPublicKeyRequestDescriptor = $convert.base64Decode('ChNHZXRQdWJsaWNLZXlSZXF1ZXN0EhgKB2FkZHJlc3MYASABKAlSB2FkZHJlc3M=');
@$core.Deprecated('Use getPublicKeyResponseDescriptor instead')
const GetPublicKeyResponse$json = const {
  '1': 'GetPublicKeyResponse',
  '2': const [
    const {'1': 'public_key', '3': 1, '4': 1, '5': 9, '10': 'publicKey'},
  ],
};

/// Descriptor for `GetPublicKeyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getPublicKeyResponseDescriptor = $convert.base64Decode('ChRHZXRQdWJsaWNLZXlSZXNwb25zZRIdCgpwdWJsaWNfa2V5GAEgASgJUglwdWJsaWNLZXk=');
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
    const {'1': 'hash', '3': 2, '4': 1, '5': 9, '10': 'hash'},
    const {'1': 'data', '3': 3, '4': 1, '5': 9, '10': 'data'},
    const {'1': 'block_time', '3': 4, '4': 1, '5': 13, '10': 'blockTime'},
    const {'1': 'header', '3': 5, '4': 1, '5': 11, '6': '.pactus.BlockHeaderInfo', '10': 'header'},
    const {'1': 'prev_cert', '3': 6, '4': 1, '5': 11, '6': '.pactus.CertificateInfo', '10': 'prevCert'},
    const {'1': 'txs', '3': 7, '4': 3, '5': 11, '6': '.pactus.TransactionInfo', '10': 'txs'},
  ],
};

/// Descriptor for `GetBlockResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBlockResponseDescriptor = $convert.base64Decode('ChBHZXRCbG9ja1Jlc3BvbnNlEhYKBmhlaWdodBgBIAEoDVIGaGVpZ2h0EhIKBGhhc2gYAiABKAlSBGhhc2gSEgoEZGF0YRgDIAEoCVIEZGF0YRIdCgpibG9ja190aW1lGAQgASgNUglibG9ja1RpbWUSLwoGaGVhZGVyGAUgASgLMhcucGFjdHVzLkJsb2NrSGVhZGVySW5mb1IGaGVhZGVyEjQKCXByZXZfY2VydBgGIAEoCzIXLnBhY3R1cy5DZXJ0aWZpY2F0ZUluZm9SCHByZXZDZXJ0EikKA3R4cxgHIAMoCzIXLnBhY3R1cy5UcmFuc2FjdGlvbkluZm9SA3R4cw==');
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
    const {'1': 'hash', '3': 1, '4': 1, '5': 9, '10': 'hash'},
  ],
};

/// Descriptor for `GetBlockHashResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBlockHashResponseDescriptor = $convert.base64Decode('ChRHZXRCbG9ja0hhc2hSZXNwb25zZRISCgRoYXNoGAEgASgJUgRoYXNo');
@$core.Deprecated('Use getBlockHeightRequestDescriptor instead')
const GetBlockHeightRequest$json = const {
  '1': 'GetBlockHeightRequest',
  '2': const [
    const {'1': 'hash', '3': 1, '4': 1, '5': 9, '10': 'hash'},
  ],
};

/// Descriptor for `GetBlockHeightRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBlockHeightRequestDescriptor = $convert.base64Decode('ChVHZXRCbG9ja0hlaWdodFJlcXVlc3QSEgoEaGFzaBgBIAEoCVIEaGFzaA==');
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
    const {'1': 'last_block_hash', '3': 2, '4': 1, '5': 9, '10': 'lastBlockHash'},
    const {'1': 'total_accounts', '3': 3, '4': 1, '5': 5, '10': 'totalAccounts'},
    const {'1': 'total_validators', '3': 4, '4': 1, '5': 5, '10': 'totalValidators'},
    const {'1': 'total_power', '3': 5, '4': 1, '5': 3, '10': 'totalPower'},
    const {'1': 'committee_power', '3': 6, '4': 1, '5': 3, '10': 'committeePower'},
    const {'1': 'committee_validators', '3': 7, '4': 3, '5': 11, '6': '.pactus.ValidatorInfo', '10': 'committeeValidators'},
    const {'1': 'is_pruned', '3': 8, '4': 1, '5': 8, '10': 'isPruned'},
    const {'1': 'pruning_height', '3': 9, '4': 1, '5': 13, '10': 'pruningHeight'},
    const {'1': 'last_block_time', '3': 10, '4': 1, '5': 3, '10': 'lastBlockTime'},
  ],
};

/// Descriptor for `GetBlockchainInfoResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBlockchainInfoResponseDescriptor = $convert.base64Decode('ChlHZXRCbG9ja2NoYWluSW5mb1Jlc3BvbnNlEioKEWxhc3RfYmxvY2tfaGVpZ2h0GAEgASgNUg9sYXN0QmxvY2tIZWlnaHQSJgoPbGFzdF9ibG9ja19oYXNoGAIgASgJUg1sYXN0QmxvY2tIYXNoEiUKDnRvdGFsX2FjY291bnRzGAMgASgFUg10b3RhbEFjY291bnRzEikKEHRvdGFsX3ZhbGlkYXRvcnMYBCABKAVSD3RvdGFsVmFsaWRhdG9ycxIfCgt0b3RhbF9wb3dlchgFIAEoA1IKdG90YWxQb3dlchInCg9jb21taXR0ZWVfcG93ZXIYBiABKANSDmNvbW1pdHRlZVBvd2VyEkgKFGNvbW1pdHRlZV92YWxpZGF0b3JzGAcgAygLMhUucGFjdHVzLlZhbGlkYXRvckluZm9SE2NvbW1pdHRlZVZhbGlkYXRvcnMSGwoJaXNfcHJ1bmVkGAggASgIUghpc1BydW5lZBIlCg5wcnVuaW5nX2hlaWdodBgJIAEoDVINcHJ1bmluZ0hlaWdodBImCg9sYXN0X2Jsb2NrX3RpbWUYCiABKANSDWxhc3RCbG9ja1RpbWU=');
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
    const {'1': 'proposal', '3': 1, '4': 1, '5': 11, '6': '.pactus.ProposalInfo', '10': 'proposal'},
    const {'1': 'instances', '3': 2, '4': 3, '5': 11, '6': '.pactus.ConsensusInfo', '10': 'instances'},
  ],
};

/// Descriptor for `GetConsensusInfoResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getConsensusInfoResponseDescriptor = $convert.base64Decode('ChhHZXRDb25zZW5zdXNJbmZvUmVzcG9uc2USMAoIcHJvcG9zYWwYASABKAsyFC5wYWN0dXMuUHJvcG9zYWxJbmZvUghwcm9wb3NhbBIzCglpbnN0YW5jZXMYAiADKAsyFS5wYWN0dXMuQ29uc2Vuc3VzSW5mb1IJaW5zdGFuY2Vz');
@$core.Deprecated('Use getTxPoolContentRequestDescriptor instead')
const GetTxPoolContentRequest$json = const {
  '1': 'GetTxPoolContentRequest',
  '2': const [
    const {'1': 'payload_type', '3': 1, '4': 1, '5': 14, '6': '.pactus.PayloadType', '10': 'payloadType'},
  ],
};

/// Descriptor for `GetTxPoolContentRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getTxPoolContentRequestDescriptor = $convert.base64Decode('ChdHZXRUeFBvb2xDb250ZW50UmVxdWVzdBI2CgxwYXlsb2FkX3R5cGUYASABKA4yEy5wYWN0dXMuUGF5bG9hZFR5cGVSC3BheWxvYWRUeXBl');
@$core.Deprecated('Use getTxPoolContentResponseDescriptor instead')
const GetTxPoolContentResponse$json = const {
  '1': 'GetTxPoolContentResponse',
  '2': const [
    const {'1': 'txs', '3': 1, '4': 3, '5': 11, '6': '.pactus.TransactionInfo', '10': 'txs'},
  ],
};

/// Descriptor for `GetTxPoolContentResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getTxPoolContentResponseDescriptor = $convert.base64Decode('ChhHZXRUeFBvb2xDb250ZW50UmVzcG9uc2USKQoDdHhzGAEgAygLMhcucGFjdHVzLlRyYW5zYWN0aW9uSW5mb1IDdHhz');
@$core.Deprecated('Use validatorInfoDescriptor instead')
const ValidatorInfo$json = const {
  '1': 'ValidatorInfo',
  '2': const [
    const {'1': 'hash', '3': 1, '4': 1, '5': 9, '10': 'hash'},
    const {'1': 'data', '3': 2, '4': 1, '5': 9, '10': 'data'},
    const {'1': 'public_key', '3': 3, '4': 1, '5': 9, '10': 'publicKey'},
    const {'1': 'number', '3': 4, '4': 1, '5': 5, '10': 'number'},
    const {'1': 'stake', '3': 5, '4': 1, '5': 3, '10': 'stake'},
    const {'1': 'last_bonding_height', '3': 6, '4': 1, '5': 13, '10': 'lastBondingHeight'},
    const {'1': 'last_sortition_height', '3': 7, '4': 1, '5': 13, '10': 'lastSortitionHeight'},
    const {'1': 'unbonding_height', '3': 8, '4': 1, '5': 13, '10': 'unbondingHeight'},
    const {'1': 'address', '3': 9, '4': 1, '5': 9, '10': 'address'},
    const {'1': 'availability_score', '3': 10, '4': 1, '5': 1, '10': 'availabilityScore'},
  ],
};

/// Descriptor for `ValidatorInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List validatorInfoDescriptor = $convert.base64Decode('Cg1WYWxpZGF0b3JJbmZvEhIKBGhhc2gYASABKAlSBGhhc2gSEgoEZGF0YRgCIAEoCVIEZGF0YRIdCgpwdWJsaWNfa2V5GAMgASgJUglwdWJsaWNLZXkSFgoGbnVtYmVyGAQgASgFUgZudW1iZXISFAoFc3Rha2UYBSABKANSBXN0YWtlEi4KE2xhc3RfYm9uZGluZ19oZWlnaHQYBiABKA1SEWxhc3RCb25kaW5nSGVpZ2h0EjIKFWxhc3Rfc29ydGl0aW9uX2hlaWdodBgHIAEoDVITbGFzdFNvcnRpdGlvbkhlaWdodBIpChB1bmJvbmRpbmdfaGVpZ2h0GAggASgNUg91bmJvbmRpbmdIZWlnaHQSGAoHYWRkcmVzcxgJIAEoCVIHYWRkcmVzcxItChJhdmFpbGFiaWxpdHlfc2NvcmUYCiABKAFSEWF2YWlsYWJpbGl0eVNjb3Jl');
@$core.Deprecated('Use accountInfoDescriptor instead')
const AccountInfo$json = const {
  '1': 'AccountInfo',
  '2': const [
    const {'1': 'hash', '3': 1, '4': 1, '5': 9, '10': 'hash'},
    const {'1': 'data', '3': 2, '4': 1, '5': 9, '10': 'data'},
    const {'1': 'number', '3': 3, '4': 1, '5': 5, '10': 'number'},
    const {'1': 'balance', '3': 4, '4': 1, '5': 3, '10': 'balance'},
    const {'1': 'address', '3': 5, '4': 1, '5': 9, '10': 'address'},
  ],
};

/// Descriptor for `AccountInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List accountInfoDescriptor = $convert.base64Decode('CgtBY2NvdW50SW5mbxISCgRoYXNoGAEgASgJUgRoYXNoEhIKBGRhdGEYAiABKAlSBGRhdGESFgoGbnVtYmVyGAMgASgFUgZudW1iZXISGAoHYmFsYW5jZRgEIAEoA1IHYmFsYW5jZRIYCgdhZGRyZXNzGAUgASgJUgdhZGRyZXNz');
@$core.Deprecated('Use blockHeaderInfoDescriptor instead')
const BlockHeaderInfo$json = const {
  '1': 'BlockHeaderInfo',
  '2': const [
    const {'1': 'version', '3': 1, '4': 1, '5': 5, '10': 'version'},
    const {'1': 'prev_block_hash', '3': 2, '4': 1, '5': 9, '10': 'prevBlockHash'},
    const {'1': 'state_root', '3': 3, '4': 1, '5': 9, '10': 'stateRoot'},
    const {'1': 'sortition_seed', '3': 4, '4': 1, '5': 9, '10': 'sortitionSeed'},
    const {'1': 'proposer_address', '3': 5, '4': 1, '5': 9, '10': 'proposerAddress'},
  ],
};

/// Descriptor for `BlockHeaderInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List blockHeaderInfoDescriptor = $convert.base64Decode('Cg9CbG9ja0hlYWRlckluZm8SGAoHdmVyc2lvbhgBIAEoBVIHdmVyc2lvbhImCg9wcmV2X2Jsb2NrX2hhc2gYAiABKAlSDXByZXZCbG9ja0hhc2gSHQoKc3RhdGVfcm9vdBgDIAEoCVIJc3RhdGVSb290EiUKDnNvcnRpdGlvbl9zZWVkGAQgASgJUg1zb3J0aXRpb25TZWVkEikKEHByb3Bvc2VyX2FkZHJlc3MYBSABKAlSD3Byb3Bvc2VyQWRkcmVzcw==');
@$core.Deprecated('Use certificateInfoDescriptor instead')
const CertificateInfo$json = const {
  '1': 'CertificateInfo',
  '2': const [
    const {'1': 'hash', '3': 1, '4': 1, '5': 9, '10': 'hash'},
    const {'1': 'round', '3': 2, '4': 1, '5': 5, '10': 'round'},
    const {'1': 'committers', '3': 3, '4': 3, '5': 5, '10': 'committers'},
    const {'1': 'absentees', '3': 4, '4': 3, '5': 5, '10': 'absentees'},
    const {'1': 'signature', '3': 5, '4': 1, '5': 9, '10': 'signature'},
  ],
};

/// Descriptor for `CertificateInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List certificateInfoDescriptor = $convert.base64Decode('Cg9DZXJ0aWZpY2F0ZUluZm8SEgoEaGFzaBgBIAEoCVIEaGFzaBIUCgVyb3VuZBgCIAEoBVIFcm91bmQSHgoKY29tbWl0dGVycxgDIAMoBVIKY29tbWl0dGVycxIcCglhYnNlbnRlZXMYBCADKAVSCWFic2VudGVlcxIcCglzaWduYXR1cmUYBSABKAlSCXNpZ25hdHVyZQ==');
@$core.Deprecated('Use voteInfoDescriptor instead')
const VoteInfo$json = const {
  '1': 'VoteInfo',
  '2': const [
    const {'1': 'type', '3': 1, '4': 1, '5': 14, '6': '.pactus.VoteType', '10': 'type'},
    const {'1': 'voter', '3': 2, '4': 1, '5': 9, '10': 'voter'},
    const {'1': 'block_hash', '3': 3, '4': 1, '5': 9, '10': 'blockHash'},
    const {'1': 'round', '3': 4, '4': 1, '5': 5, '10': 'round'},
    const {'1': 'cp_round', '3': 5, '4': 1, '5': 5, '10': 'cpRound'},
    const {'1': 'cp_value', '3': 6, '4': 1, '5': 5, '10': 'cpValue'},
  ],
};

/// Descriptor for `VoteInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List voteInfoDescriptor = $convert.base64Decode('CghWb3RlSW5mbxIkCgR0eXBlGAEgASgOMhAucGFjdHVzLlZvdGVUeXBlUgR0eXBlEhQKBXZvdGVyGAIgASgJUgV2b3RlchIdCgpibG9ja19oYXNoGAMgASgJUglibG9ja0hhc2gSFAoFcm91bmQYBCABKAVSBXJvdW5kEhkKCGNwX3JvdW5kGAUgASgFUgdjcFJvdW5kEhkKCGNwX3ZhbHVlGAYgASgFUgdjcFZhbHVl');
@$core.Deprecated('Use consensusInfoDescriptor instead')
const ConsensusInfo$json = const {
  '1': 'ConsensusInfo',
  '2': const [
    const {'1': 'address', '3': 1, '4': 1, '5': 9, '10': 'address'},
    const {'1': 'active', '3': 2, '4': 1, '5': 8, '10': 'active'},
    const {'1': 'height', '3': 3, '4': 1, '5': 13, '10': 'height'},
    const {'1': 'round', '3': 4, '4': 1, '5': 5, '10': 'round'},
    const {'1': 'votes', '3': 5, '4': 3, '5': 11, '6': '.pactus.VoteInfo', '10': 'votes'},
  ],
};

/// Descriptor for `ConsensusInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List consensusInfoDescriptor = $convert.base64Decode('Cg1Db25zZW5zdXNJbmZvEhgKB2FkZHJlc3MYASABKAlSB2FkZHJlc3MSFgoGYWN0aXZlGAIgASgIUgZhY3RpdmUSFgoGaGVpZ2h0GAMgASgNUgZoZWlnaHQSFAoFcm91bmQYBCABKAVSBXJvdW5kEiYKBXZvdGVzGAUgAygLMhAucGFjdHVzLlZvdGVJbmZvUgV2b3Rlcw==');
@$core.Deprecated('Use proposalInfoDescriptor instead')
const ProposalInfo$json = const {
  '1': 'ProposalInfo',
  '2': const [
    const {'1': 'height', '3': 1, '4': 1, '5': 13, '10': 'height'},
    const {'1': 'round', '3': 2, '4': 1, '5': 5, '10': 'round'},
    const {'1': 'block_data', '3': 3, '4': 1, '5': 9, '10': 'blockData'},
    const {'1': 'signature', '3': 4, '4': 1, '5': 9, '10': 'signature'},
  ],
};

/// Descriptor for `ProposalInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List proposalInfoDescriptor = $convert.base64Decode('CgxQcm9wb3NhbEluZm8SFgoGaGVpZ2h0GAEgASgNUgZoZWlnaHQSFAoFcm91bmQYAiABKAVSBXJvdW5kEh0KCmJsb2NrX2RhdGEYAyABKAlSCWJsb2NrRGF0YRIcCglzaWduYXR1cmUYBCABKAlSCXNpZ25hdHVyZQ==');
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
    const {'1': 'GetValidatorAddresses', '2': '.pactus.GetValidatorAddressesRequest', '3': '.pactus.GetValidatorAddressesResponse'},
    const {'1': 'GetPublicKey', '2': '.pactus.GetPublicKeyRequest', '3': '.pactus.GetPublicKeyResponse'},
    const {'1': 'GetTxPoolContent', '2': '.pactus.GetTxPoolContentRequest', '3': '.pactus.GetTxPoolContentResponse'},
  ],
};

@$core.Deprecated('Use blockchainServiceDescriptor instead')
const $core.Map<$core.String, $core.Map<$core.String, $core.dynamic>> BlockchainServiceBase$messageJson = const {
  '.pactus.GetBlockRequest': GetBlockRequest$json,
  '.pactus.GetBlockResponse': GetBlockResponse$json,
  '.pactus.BlockHeaderInfo': BlockHeaderInfo$json,
  '.pactus.CertificateInfo': CertificateInfo$json,
  '.pactus.TransactionInfo': $0.TransactionInfo$json,
  '.pactus.PayloadTransfer': $0.PayloadTransfer$json,
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
  '.pactus.ProposalInfo': ProposalInfo$json,
  '.pactus.ConsensusInfo': ConsensusInfo$json,
  '.pactus.VoteInfo': VoteInfo$json,
  '.pactus.GetAccountRequest': GetAccountRequest$json,
  '.pactus.GetAccountResponse': GetAccountResponse$json,
  '.pactus.AccountInfo': AccountInfo$json,
  '.pactus.GetValidatorRequest': GetValidatorRequest$json,
  '.pactus.GetValidatorResponse': GetValidatorResponse$json,
  '.pactus.GetValidatorByNumberRequest': GetValidatorByNumberRequest$json,
  '.pactus.GetValidatorAddressesRequest': GetValidatorAddressesRequest$json,
  '.pactus.GetValidatorAddressesResponse': GetValidatorAddressesResponse$json,
  '.pactus.GetPublicKeyRequest': GetPublicKeyRequest$json,
  '.pactus.GetPublicKeyResponse': GetPublicKeyResponse$json,
  '.pactus.GetTxPoolContentRequest': GetTxPoolContentRequest$json,
  '.pactus.GetTxPoolContentResponse': GetTxPoolContentResponse$json,
};

/// Descriptor for `Blockchain`. Decode as a `google.protobuf.ServiceDescriptorProto`.
final $typed_data.Uint8List blockchainServiceDescriptor = $convert.base64Decode('CgpCbG9ja2NoYWluEj0KCEdldEJsb2NrEhcucGFjdHVzLkdldEJsb2NrUmVxdWVzdBoYLnBhY3R1cy5HZXRCbG9ja1Jlc3BvbnNlEkkKDEdldEJsb2NrSGFzaBIbLnBhY3R1cy5HZXRCbG9ja0hhc2hSZXF1ZXN0GhwucGFjdHVzLkdldEJsb2NrSGFzaFJlc3BvbnNlEk8KDkdldEJsb2NrSGVpZ2h0Eh0ucGFjdHVzLkdldEJsb2NrSGVpZ2h0UmVxdWVzdBoeLnBhY3R1cy5HZXRCbG9ja0hlaWdodFJlc3BvbnNlElgKEUdldEJsb2NrY2hhaW5JbmZvEiAucGFjdHVzLkdldEJsb2NrY2hhaW5JbmZvUmVxdWVzdBohLnBhY3R1cy5HZXRCbG9ja2NoYWluSW5mb1Jlc3BvbnNlElUKEEdldENvbnNlbnN1c0luZm8SHy5wYWN0dXMuR2V0Q29uc2Vuc3VzSW5mb1JlcXVlc3QaIC5wYWN0dXMuR2V0Q29uc2Vuc3VzSW5mb1Jlc3BvbnNlEkMKCkdldEFjY291bnQSGS5wYWN0dXMuR2V0QWNjb3VudFJlcXVlc3QaGi5wYWN0dXMuR2V0QWNjb3VudFJlc3BvbnNlEkkKDEdldFZhbGlkYXRvchIbLnBhY3R1cy5HZXRWYWxpZGF0b3JSZXF1ZXN0GhwucGFjdHVzLkdldFZhbGlkYXRvclJlc3BvbnNlElkKFEdldFZhbGlkYXRvckJ5TnVtYmVyEiMucGFjdHVzLkdldFZhbGlkYXRvckJ5TnVtYmVyUmVxdWVzdBocLnBhY3R1cy5HZXRWYWxpZGF0b3JSZXNwb25zZRJkChVHZXRWYWxpZGF0b3JBZGRyZXNzZXMSJC5wYWN0dXMuR2V0VmFsaWRhdG9yQWRkcmVzc2VzUmVxdWVzdBolLnBhY3R1cy5HZXRWYWxpZGF0b3JBZGRyZXNzZXNSZXNwb25zZRJJCgxHZXRQdWJsaWNLZXkSGy5wYWN0dXMuR2V0UHVibGljS2V5UmVxdWVzdBocLnBhY3R1cy5HZXRQdWJsaWNLZXlSZXNwb25zZRJVChBHZXRUeFBvb2xDb250ZW50Eh8ucGFjdHVzLkdldFR4UG9vbENvbnRlbnRSZXF1ZXN0GiAucGFjdHVzLkdldFR4UG9vbENvbnRlbnRSZXNwb25zZQ==');
