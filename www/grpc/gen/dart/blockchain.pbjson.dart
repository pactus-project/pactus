// This is a generated file - do not edit.
//
// Generated from blockchain.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports
// ignore_for_file: unused_import

import 'dart:convert' as $convert;
import 'dart:core' as $core;
import 'dart:typed_data' as $typed_data;

import 'transaction.pbjson.dart' as $0;

@$core.Deprecated('Use blockVerbosityDescriptor instead')
const BlockVerbosity$json = {
  '1': 'BlockVerbosity',
  '2': [
    {'1': 'BLOCK_VERBOSITY_DATA', '2': 0},
    {'1': 'BLOCK_VERBOSITY_INFO', '2': 1},
    {'1': 'BLOCK_VERBOSITY_TRANSACTIONS', '2': 2},
  ],
};

/// Descriptor for `BlockVerbosity`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List blockVerbosityDescriptor = $convert.base64Decode(
    'Cg5CbG9ja1ZlcmJvc2l0eRIYChRCTE9DS19WRVJCT1NJVFlfREFUQRAAEhgKFEJMT0NLX1ZFUk'
    'JPU0lUWV9JTkZPEAESIAocQkxPQ0tfVkVSQk9TSVRZX1RSQU5TQUNUSU9OUxAC');

@$core.Deprecated('Use voteTypeDescriptor instead')
const VoteType$json = {
  '1': 'VoteType',
  '2': [
    {'1': 'VOTE_TYPE_UNSPECIFIED', '2': 0},
    {'1': 'VOTE_TYPE_PREPARE', '2': 1},
    {'1': 'VOTE_TYPE_PRECOMMIT', '2': 2},
    {'1': 'VOTE_TYPE_CP_PRE_VOTE', '2': 3},
    {'1': 'VOTE_TYPE_CP_MAIN_VOTE', '2': 4},
    {'1': 'VOTE_TYPE_CP_DECIDED', '2': 5},
  ],
};

/// Descriptor for `VoteType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List voteTypeDescriptor = $convert.base64Decode(
    'CghWb3RlVHlwZRIZChVWT1RFX1RZUEVfVU5TUEVDSUZJRUQQABIVChFWT1RFX1RZUEVfUFJFUE'
    'FSRRABEhcKE1ZPVEVfVFlQRV9QUkVDT01NSVQQAhIZChVWT1RFX1RZUEVfQ1BfUFJFX1ZPVEUQ'
    'AxIaChZWT1RFX1RZUEVfQ1BfTUFJTl9WT1RFEAQSGAoUVk9URV9UWVBFX0NQX0RFQ0lERUQQBQ'
    '==');

@$core.Deprecated('Use getAccountRequestDescriptor instead')
const GetAccountRequest$json = {
  '1': 'GetAccountRequest',
  '2': [
    {'1': 'address', '3': 1, '4': 1, '5': 9, '10': 'address'},
  ],
};

/// Descriptor for `GetAccountRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getAccountRequestDescriptor = $convert.base64Decode(
    'ChFHZXRBY2NvdW50UmVxdWVzdBIYCgdhZGRyZXNzGAEgASgJUgdhZGRyZXNz');

@$core.Deprecated('Use getAccountResponseDescriptor instead')
const GetAccountResponse$json = {
  '1': 'GetAccountResponse',
  '2': [
    {
      '1': 'account',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.pactus.AccountInfo',
      '10': 'account'
    },
  ],
};

/// Descriptor for `GetAccountResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getAccountResponseDescriptor = $convert.base64Decode(
    'ChJHZXRBY2NvdW50UmVzcG9uc2USLQoHYWNjb3VudBgBIAEoCzITLnBhY3R1cy5BY2NvdW50SW'
    '5mb1IHYWNjb3VudA==');

@$core.Deprecated('Use getValidatorAddressesRequestDescriptor instead')
const GetValidatorAddressesRequest$json = {
  '1': 'GetValidatorAddressesRequest',
};

/// Descriptor for `GetValidatorAddressesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getValidatorAddressesRequestDescriptor =
    $convert.base64Decode('ChxHZXRWYWxpZGF0b3JBZGRyZXNzZXNSZXF1ZXN0');

@$core.Deprecated('Use getValidatorAddressesResponseDescriptor instead')
const GetValidatorAddressesResponse$json = {
  '1': 'GetValidatorAddressesResponse',
  '2': [
    {'1': 'addresses', '3': 1, '4': 3, '5': 9, '10': 'addresses'},
  ],
};

/// Descriptor for `GetValidatorAddressesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getValidatorAddressesResponseDescriptor =
    $convert.base64Decode(
        'Ch1HZXRWYWxpZGF0b3JBZGRyZXNzZXNSZXNwb25zZRIcCglhZGRyZXNzZXMYASADKAlSCWFkZH'
        'Jlc3Nlcw==');

@$core.Deprecated('Use getValidatorRequestDescriptor instead')
const GetValidatorRequest$json = {
  '1': 'GetValidatorRequest',
  '2': [
    {'1': 'address', '3': 1, '4': 1, '5': 9, '10': 'address'},
  ],
};

/// Descriptor for `GetValidatorRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getValidatorRequestDescriptor =
    $convert.base64Decode(
        'ChNHZXRWYWxpZGF0b3JSZXF1ZXN0EhgKB2FkZHJlc3MYASABKAlSB2FkZHJlc3M=');

@$core.Deprecated('Use getValidatorByNumberRequestDescriptor instead')
const GetValidatorByNumberRequest$json = {
  '1': 'GetValidatorByNumberRequest',
  '2': [
    {'1': 'number', '3': 1, '4': 1, '5': 5, '10': 'number'},
  ],
};

/// Descriptor for `GetValidatorByNumberRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getValidatorByNumberRequestDescriptor =
    $convert.base64Decode(
        'ChtHZXRWYWxpZGF0b3JCeU51bWJlclJlcXVlc3QSFgoGbnVtYmVyGAEgASgFUgZudW1iZXI=');

@$core.Deprecated('Use getValidatorResponseDescriptor instead')
const GetValidatorResponse$json = {
  '1': 'GetValidatorResponse',
  '2': [
    {
      '1': 'validator',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.pactus.ValidatorInfo',
      '10': 'validator'
    },
  ],
};

/// Descriptor for `GetValidatorResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getValidatorResponseDescriptor = $convert.base64Decode(
    'ChRHZXRWYWxpZGF0b3JSZXNwb25zZRIzCgl2YWxpZGF0b3IYASABKAsyFS5wYWN0dXMuVmFsaW'
    'RhdG9ySW5mb1IJdmFsaWRhdG9y');

@$core.Deprecated('Use getPublicKeyRequestDescriptor instead')
const GetPublicKeyRequest$json = {
  '1': 'GetPublicKeyRequest',
  '2': [
    {'1': 'address', '3': 1, '4': 1, '5': 9, '10': 'address'},
  ],
};

/// Descriptor for `GetPublicKeyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getPublicKeyRequestDescriptor =
    $convert.base64Decode(
        'ChNHZXRQdWJsaWNLZXlSZXF1ZXN0EhgKB2FkZHJlc3MYASABKAlSB2FkZHJlc3M=');

@$core.Deprecated('Use getPublicKeyResponseDescriptor instead')
const GetPublicKeyResponse$json = {
  '1': 'GetPublicKeyResponse',
  '2': [
    {'1': 'public_key', '3': 1, '4': 1, '5': 9, '10': 'publicKey'},
  ],
};

/// Descriptor for `GetPublicKeyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getPublicKeyResponseDescriptor = $convert.base64Decode(
    'ChRHZXRQdWJsaWNLZXlSZXNwb25zZRIdCgpwdWJsaWNfa2V5GAEgASgJUglwdWJsaWNLZXk=');

@$core.Deprecated('Use getBlockRequestDescriptor instead')
const GetBlockRequest$json = {
  '1': 'GetBlockRequest',
  '2': [
    {'1': 'height', '3': 1, '4': 1, '5': 13, '10': 'height'},
    {
      '1': 'verbosity',
      '3': 2,
      '4': 1,
      '5': 14,
      '6': '.pactus.BlockVerbosity',
      '10': 'verbosity'
    },
  ],
};

/// Descriptor for `GetBlockRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBlockRequestDescriptor = $convert.base64Decode(
    'Cg9HZXRCbG9ja1JlcXVlc3QSFgoGaGVpZ2h0GAEgASgNUgZoZWlnaHQSNAoJdmVyYm9zaXR5GA'
    'IgASgOMhYucGFjdHVzLkJsb2NrVmVyYm9zaXR5Ugl2ZXJib3NpdHk=');

@$core.Deprecated('Use getBlockResponseDescriptor instead')
const GetBlockResponse$json = {
  '1': 'GetBlockResponse',
  '2': [
    {'1': 'height', '3': 1, '4': 1, '5': 13, '10': 'height'},
    {'1': 'hash', '3': 2, '4': 1, '5': 9, '10': 'hash'},
    {'1': 'data', '3': 3, '4': 1, '5': 9, '10': 'data'},
    {'1': 'block_time', '3': 4, '4': 1, '5': 13, '10': 'blockTime'},
    {
      '1': 'header',
      '3': 5,
      '4': 1,
      '5': 11,
      '6': '.pactus.BlockHeaderInfo',
      '10': 'header'
    },
    {
      '1': 'prev_cert',
      '3': 6,
      '4': 1,
      '5': 11,
      '6': '.pactus.CertificateInfo',
      '10': 'prevCert'
    },
    {
      '1': 'txs',
      '3': 7,
      '4': 3,
      '5': 11,
      '6': '.pactus.TransactionInfo',
      '10': 'txs'
    },
  ],
};

/// Descriptor for `GetBlockResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBlockResponseDescriptor = $convert.base64Decode(
    'ChBHZXRCbG9ja1Jlc3BvbnNlEhYKBmhlaWdodBgBIAEoDVIGaGVpZ2h0EhIKBGhhc2gYAiABKA'
    'lSBGhhc2gSEgoEZGF0YRgDIAEoCVIEZGF0YRIdCgpibG9ja190aW1lGAQgASgNUglibG9ja1Rp'
    'bWUSLwoGaGVhZGVyGAUgASgLMhcucGFjdHVzLkJsb2NrSGVhZGVySW5mb1IGaGVhZGVyEjQKCX'
    'ByZXZfY2VydBgGIAEoCzIXLnBhY3R1cy5DZXJ0aWZpY2F0ZUluZm9SCHByZXZDZXJ0EikKA3R4'
    'cxgHIAMoCzIXLnBhY3R1cy5UcmFuc2FjdGlvbkluZm9SA3R4cw==');

@$core.Deprecated('Use getBlockHashRequestDescriptor instead')
const GetBlockHashRequest$json = {
  '1': 'GetBlockHashRequest',
  '2': [
    {'1': 'height', '3': 1, '4': 1, '5': 13, '10': 'height'},
  ],
};

/// Descriptor for `GetBlockHashRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBlockHashRequestDescriptor =
    $convert.base64Decode(
        'ChNHZXRCbG9ja0hhc2hSZXF1ZXN0EhYKBmhlaWdodBgBIAEoDVIGaGVpZ2h0');

@$core.Deprecated('Use getBlockHashResponseDescriptor instead')
const GetBlockHashResponse$json = {
  '1': 'GetBlockHashResponse',
  '2': [
    {'1': 'hash', '3': 1, '4': 1, '5': 9, '10': 'hash'},
  ],
};

/// Descriptor for `GetBlockHashResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBlockHashResponseDescriptor = $convert
    .base64Decode('ChRHZXRCbG9ja0hhc2hSZXNwb25zZRISCgRoYXNoGAEgASgJUgRoYXNo');

@$core.Deprecated('Use getBlockHeightRequestDescriptor instead')
const GetBlockHeightRequest$json = {
  '1': 'GetBlockHeightRequest',
  '2': [
    {'1': 'hash', '3': 1, '4': 1, '5': 9, '10': 'hash'},
  ],
};

/// Descriptor for `GetBlockHeightRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBlockHeightRequestDescriptor =
    $convert.base64Decode(
        'ChVHZXRCbG9ja0hlaWdodFJlcXVlc3QSEgoEaGFzaBgBIAEoCVIEaGFzaA==');

@$core.Deprecated('Use getBlockHeightResponseDescriptor instead')
const GetBlockHeightResponse$json = {
  '1': 'GetBlockHeightResponse',
  '2': [
    {'1': 'height', '3': 1, '4': 1, '5': 13, '10': 'height'},
  ],
};

/// Descriptor for `GetBlockHeightResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBlockHeightResponseDescriptor =
    $convert.base64Decode(
        'ChZHZXRCbG9ja0hlaWdodFJlc3BvbnNlEhYKBmhlaWdodBgBIAEoDVIGaGVpZ2h0');

@$core.Deprecated('Use getBlockchainInfoRequestDescriptor instead')
const GetBlockchainInfoRequest$json = {
  '1': 'GetBlockchainInfoRequest',
};

/// Descriptor for `GetBlockchainInfoRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBlockchainInfoRequestDescriptor =
    $convert.base64Decode('ChhHZXRCbG9ja2NoYWluSW5mb1JlcXVlc3Q=');

@$core.Deprecated('Use getBlockchainInfoResponseDescriptor instead')
const GetBlockchainInfoResponse$json = {
  '1': 'GetBlockchainInfoResponse',
  '2': [
    {
      '1': 'last_block_height',
      '3': 1,
      '4': 1,
      '5': 13,
      '10': 'lastBlockHeight'
    },
    {'1': 'last_block_hash', '3': 2, '4': 1, '5': 9, '10': 'lastBlockHash'},
    {'1': 'last_block_time', '3': 10, '4': 1, '5': 3, '10': 'lastBlockTime'},
    {'1': 'total_accounts', '3': 3, '4': 1, '5': 5, '10': 'totalAccounts'},
    {'1': 'total_validators', '3': 4, '4': 1, '5': 5, '10': 'totalValidators'},
    {
      '1': 'active_validators',
      '3': 12,
      '4': 1,
      '5': 5,
      '10': 'activeValidators'
    },
    {'1': 'total_power', '3': 5, '4': 1, '5': 3, '10': 'totalPower'},
    {'1': 'committee_power', '3': 6, '4': 1, '5': 3, '10': 'committeePower'},
    {'1': 'is_pruned', '3': 8, '4': 1, '5': 8, '10': 'isPruned'},
    {'1': 'pruning_height', '3': 9, '4': 1, '5': 13, '10': 'pruningHeight'},
    {'1': 'in_committee', '3': 13, '4': 1, '5': 8, '10': 'inCommittee'},
  ],
};

/// Descriptor for `GetBlockchainInfoResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBlockchainInfoResponseDescriptor = $convert.base64Decode(
    'ChlHZXRCbG9ja2NoYWluSW5mb1Jlc3BvbnNlEioKEWxhc3RfYmxvY2tfaGVpZ2h0GAEgASgNUg'
    '9sYXN0QmxvY2tIZWlnaHQSJgoPbGFzdF9ibG9ja19oYXNoGAIgASgJUg1sYXN0QmxvY2tIYXNo'
    'EiYKD2xhc3RfYmxvY2tfdGltZRgKIAEoA1INbGFzdEJsb2NrVGltZRIlCg50b3RhbF9hY2NvdW'
    '50cxgDIAEoBVINdG90YWxBY2NvdW50cxIpChB0b3RhbF92YWxpZGF0b3JzGAQgASgFUg90b3Rh'
    'bFZhbGlkYXRvcnMSKwoRYWN0aXZlX3ZhbGlkYXRvcnMYDCABKAVSEGFjdGl2ZVZhbGlkYXRvcn'
    'MSHwoLdG90YWxfcG93ZXIYBSABKANSCnRvdGFsUG93ZXISJwoPY29tbWl0dGVlX3Bvd2VyGAYg'
    'ASgDUg5jb21taXR0ZWVQb3dlchIbCglpc19wcnVuZWQYCCABKAhSCGlzUHJ1bmVkEiUKDnBydW'
    '5pbmdfaGVpZ2h0GAkgASgNUg1wcnVuaW5nSGVpZ2h0EiEKDGluX2NvbW1pdHRlZRgNIAEoCFIL'
    'aW5Db21taXR0ZWU=');

@$core.Deprecated('Use getCommitteeInfoRequestDescriptor instead')
const GetCommitteeInfoRequest$json = {
  '1': 'GetCommitteeInfoRequest',
};

/// Descriptor for `GetCommitteeInfoRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getCommitteeInfoRequestDescriptor =
    $convert.base64Decode('ChdHZXRDb21taXR0ZWVJbmZvUmVxdWVzdA==');

@$core.Deprecated('Use getCommitteeInfoResponseDescriptor instead')
const GetCommitteeInfoResponse$json = {
  '1': 'GetCommitteeInfoResponse',
  '2': [
    {'1': 'committee_power', '3': 1, '4': 1, '5': 3, '10': 'committeePower'},
    {
      '1': 'validators',
      '3': 2,
      '4': 3,
      '5': 11,
      '6': '.pactus.ValidatorInfo',
      '10': 'validators'
    },
    {
      '1': 'protocol_versions',
      '3': 3,
      '4': 3,
      '5': 11,
      '6': '.pactus.GetCommitteeInfoResponse.ProtocolVersionsEntry',
      '10': 'protocolVersions'
    },
  ],
  '3': [GetCommitteeInfoResponse_ProtocolVersionsEntry$json],
};

@$core.Deprecated('Use getCommitteeInfoResponseDescriptor instead')
const GetCommitteeInfoResponse_ProtocolVersionsEntry$json = {
  '1': 'ProtocolVersionsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 5, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 1, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `GetCommitteeInfoResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getCommitteeInfoResponseDescriptor = $convert.base64Decode(
    'ChhHZXRDb21taXR0ZWVJbmZvUmVzcG9uc2USJwoPY29tbWl0dGVlX3Bvd2VyGAEgASgDUg5jb2'
    '1taXR0ZWVQb3dlchI1Cgp2YWxpZGF0b3JzGAIgAygLMhUucGFjdHVzLlZhbGlkYXRvckluZm9S'
    'CnZhbGlkYXRvcnMSYwoRcHJvdG9jb2xfdmVyc2lvbnMYAyADKAsyNi5wYWN0dXMuR2V0Q29tbW'
    'l0dGVlSW5mb1Jlc3BvbnNlLlByb3RvY29sVmVyc2lvbnNFbnRyeVIQcHJvdG9jb2xWZXJzaW9u'
    'cxpDChVQcm90b2NvbFZlcnNpb25zRW50cnkSEAoDa2V5GAEgASgFUgNrZXkSFAoFdmFsdWUYAi'
    'ABKAFSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use getConsensusInfoRequestDescriptor instead')
const GetConsensusInfoRequest$json = {
  '1': 'GetConsensusInfoRequest',
};

/// Descriptor for `GetConsensusInfoRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getConsensusInfoRequestDescriptor =
    $convert.base64Decode('ChdHZXRDb25zZW5zdXNJbmZvUmVxdWVzdA==');

@$core.Deprecated('Use getConsensusInfoResponseDescriptor instead')
const GetConsensusInfoResponse$json = {
  '1': 'GetConsensusInfoResponse',
  '2': [
    {
      '1': 'proposal',
      '3': 1,
      '4': 1,
      '5': 11,
      '6': '.pactus.ProposalInfo',
      '10': 'proposal'
    },
    {
      '1': 'instances',
      '3': 2,
      '4': 3,
      '5': 11,
      '6': '.pactus.ConsensusInfo',
      '10': 'instances'
    },
  ],
};

/// Descriptor for `GetConsensusInfoResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getConsensusInfoResponseDescriptor = $convert.base64Decode(
    'ChhHZXRDb25zZW5zdXNJbmZvUmVzcG9uc2USMAoIcHJvcG9zYWwYASABKAsyFC5wYWN0dXMuUH'
    'JvcG9zYWxJbmZvUghwcm9wb3NhbBIzCglpbnN0YW5jZXMYAiADKAsyFS5wYWN0dXMuQ29uc2Vu'
    'c3VzSW5mb1IJaW5zdGFuY2Vz');

@$core.Deprecated('Use getTxPoolContentRequestDescriptor instead')
const GetTxPoolContentRequest$json = {
  '1': 'GetTxPoolContentRequest',
  '2': [
    {
      '1': 'payload_type',
      '3': 1,
      '4': 1,
      '5': 14,
      '6': '.pactus.PayloadType',
      '10': 'payloadType'
    },
  ],
};

/// Descriptor for `GetTxPoolContentRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getTxPoolContentRequestDescriptor =
    $convert.base64Decode(
        'ChdHZXRUeFBvb2xDb250ZW50UmVxdWVzdBI2CgxwYXlsb2FkX3R5cGUYASABKA4yEy5wYWN0dX'
        'MuUGF5bG9hZFR5cGVSC3BheWxvYWRUeXBl');

@$core.Deprecated('Use getTxPoolContentResponseDescriptor instead')
const GetTxPoolContentResponse$json = {
  '1': 'GetTxPoolContentResponse',
  '2': [
    {
      '1': 'txs',
      '3': 1,
      '4': 3,
      '5': 11,
      '6': '.pactus.TransactionInfo',
      '10': 'txs'
    },
  ],
};

/// Descriptor for `GetTxPoolContentResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getTxPoolContentResponseDescriptor =
    $convert.base64Decode(
        'ChhHZXRUeFBvb2xDb250ZW50UmVzcG9uc2USKQoDdHhzGAEgAygLMhcucGFjdHVzLlRyYW5zYW'
        'N0aW9uSW5mb1IDdHhz');

@$core.Deprecated('Use validatorInfoDescriptor instead')
const ValidatorInfo$json = {
  '1': 'ValidatorInfo',
  '2': [
    {'1': 'hash', '3': 1, '4': 1, '5': 9, '10': 'hash'},
    {'1': 'data', '3': 2, '4': 1, '5': 9, '10': 'data'},
    {'1': 'public_key', '3': 3, '4': 1, '5': 9, '10': 'publicKey'},
    {'1': 'number', '3': 4, '4': 1, '5': 5, '10': 'number'},
    {'1': 'stake', '3': 5, '4': 1, '5': 3, '10': 'stake'},
    {
      '1': 'last_bonding_height',
      '3': 6,
      '4': 1,
      '5': 13,
      '10': 'lastBondingHeight'
    },
    {
      '1': 'last_sortition_height',
      '3': 7,
      '4': 1,
      '5': 13,
      '10': 'lastSortitionHeight'
    },
    {'1': 'unbonding_height', '3': 8, '4': 1, '5': 13, '10': 'unbondingHeight'},
    {'1': 'address', '3': 9, '4': 1, '5': 9, '10': 'address'},
    {
      '1': 'availability_score',
      '3': 10,
      '4': 1,
      '5': 1,
      '10': 'availabilityScore'
    },
    {'1': 'protocol_version', '3': 11, '4': 1, '5': 5, '10': 'protocolVersion'},
  ],
};

/// Descriptor for `ValidatorInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List validatorInfoDescriptor = $convert.base64Decode(
    'Cg1WYWxpZGF0b3JJbmZvEhIKBGhhc2gYASABKAlSBGhhc2gSEgoEZGF0YRgCIAEoCVIEZGF0YR'
    'IdCgpwdWJsaWNfa2V5GAMgASgJUglwdWJsaWNLZXkSFgoGbnVtYmVyGAQgASgFUgZudW1iZXIS'
    'FAoFc3Rha2UYBSABKANSBXN0YWtlEi4KE2xhc3RfYm9uZGluZ19oZWlnaHQYBiABKA1SEWxhc3'
    'RCb25kaW5nSGVpZ2h0EjIKFWxhc3Rfc29ydGl0aW9uX2hlaWdodBgHIAEoDVITbGFzdFNvcnRp'
    'dGlvbkhlaWdodBIpChB1bmJvbmRpbmdfaGVpZ2h0GAggASgNUg91bmJvbmRpbmdIZWlnaHQSGA'
    'oHYWRkcmVzcxgJIAEoCVIHYWRkcmVzcxItChJhdmFpbGFiaWxpdHlfc2NvcmUYCiABKAFSEWF2'
    'YWlsYWJpbGl0eVNjb3JlEikKEHByb3RvY29sX3ZlcnNpb24YCyABKAVSD3Byb3RvY29sVmVyc2'
    'lvbg==');

@$core.Deprecated('Use accountInfoDescriptor instead')
const AccountInfo$json = {
  '1': 'AccountInfo',
  '2': [
    {'1': 'hash', '3': 1, '4': 1, '5': 9, '10': 'hash'},
    {'1': 'data', '3': 2, '4': 1, '5': 9, '10': 'data'},
    {'1': 'number', '3': 3, '4': 1, '5': 5, '10': 'number'},
    {'1': 'balance', '3': 4, '4': 1, '5': 3, '10': 'balance'},
    {'1': 'address', '3': 5, '4': 1, '5': 9, '10': 'address'},
  ],
};

/// Descriptor for `AccountInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List accountInfoDescriptor = $convert.base64Decode(
    'CgtBY2NvdW50SW5mbxISCgRoYXNoGAEgASgJUgRoYXNoEhIKBGRhdGEYAiABKAlSBGRhdGESFg'
    'oGbnVtYmVyGAMgASgFUgZudW1iZXISGAoHYmFsYW5jZRgEIAEoA1IHYmFsYW5jZRIYCgdhZGRy'
    'ZXNzGAUgASgJUgdhZGRyZXNz');

@$core.Deprecated('Use blockHeaderInfoDescriptor instead')
const BlockHeaderInfo$json = {
  '1': 'BlockHeaderInfo',
  '2': [
    {'1': 'version', '3': 1, '4': 1, '5': 5, '10': 'version'},
    {'1': 'prev_block_hash', '3': 2, '4': 1, '5': 9, '10': 'prevBlockHash'},
    {'1': 'state_root', '3': 3, '4': 1, '5': 9, '10': 'stateRoot'},
    {'1': 'sortition_seed', '3': 4, '4': 1, '5': 9, '10': 'sortitionSeed'},
    {'1': 'proposer_address', '3': 5, '4': 1, '5': 9, '10': 'proposerAddress'},
  ],
};

/// Descriptor for `BlockHeaderInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List blockHeaderInfoDescriptor = $convert.base64Decode(
    'Cg9CbG9ja0hlYWRlckluZm8SGAoHdmVyc2lvbhgBIAEoBVIHdmVyc2lvbhImCg9wcmV2X2Jsb2'
    'NrX2hhc2gYAiABKAlSDXByZXZCbG9ja0hhc2gSHQoKc3RhdGVfcm9vdBgDIAEoCVIJc3RhdGVS'
    'b290EiUKDnNvcnRpdGlvbl9zZWVkGAQgASgJUg1zb3J0aXRpb25TZWVkEikKEHByb3Bvc2VyX2'
    'FkZHJlc3MYBSABKAlSD3Byb3Bvc2VyQWRkcmVzcw==');

@$core.Deprecated('Use certificateInfoDescriptor instead')
const CertificateInfo$json = {
  '1': 'CertificateInfo',
  '2': [
    {'1': 'hash', '3': 1, '4': 1, '5': 9, '10': 'hash'},
    {'1': 'round', '3': 2, '4': 1, '5': 5, '10': 'round'},
    {'1': 'committers', '3': 3, '4': 3, '5': 5, '10': 'committers'},
    {'1': 'absentees', '3': 4, '4': 3, '5': 5, '10': 'absentees'},
    {'1': 'signature', '3': 5, '4': 1, '5': 9, '10': 'signature'},
  ],
};

/// Descriptor for `CertificateInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List certificateInfoDescriptor = $convert.base64Decode(
    'Cg9DZXJ0aWZpY2F0ZUluZm8SEgoEaGFzaBgBIAEoCVIEaGFzaBIUCgVyb3VuZBgCIAEoBVIFcm'
    '91bmQSHgoKY29tbWl0dGVycxgDIAMoBVIKY29tbWl0dGVycxIcCglhYnNlbnRlZXMYBCADKAVS'
    'CWFic2VudGVlcxIcCglzaWduYXR1cmUYBSABKAlSCXNpZ25hdHVyZQ==');

@$core.Deprecated('Use voteInfoDescriptor instead')
const VoteInfo$json = {
  '1': 'VoteInfo',
  '2': [
    {
      '1': 'type',
      '3': 1,
      '4': 1,
      '5': 14,
      '6': '.pactus.VoteType',
      '10': 'type'
    },
    {'1': 'voter', '3': 2, '4': 1, '5': 9, '10': 'voter'},
    {'1': 'block_hash', '3': 3, '4': 1, '5': 9, '10': 'blockHash'},
    {'1': 'round', '3': 4, '4': 1, '5': 5, '10': 'round'},
    {'1': 'cp_round', '3': 5, '4': 1, '5': 5, '10': 'cpRound'},
    {'1': 'cp_value', '3': 6, '4': 1, '5': 5, '10': 'cpValue'},
  ],
};

/// Descriptor for `VoteInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List voteInfoDescriptor = $convert.base64Decode(
    'CghWb3RlSW5mbxIkCgR0eXBlGAEgASgOMhAucGFjdHVzLlZvdGVUeXBlUgR0eXBlEhQKBXZvdG'
    'VyGAIgASgJUgV2b3RlchIdCgpibG9ja19oYXNoGAMgASgJUglibG9ja0hhc2gSFAoFcm91bmQY'
    'BCABKAVSBXJvdW5kEhkKCGNwX3JvdW5kGAUgASgFUgdjcFJvdW5kEhkKCGNwX3ZhbHVlGAYgAS'
    'gFUgdjcFZhbHVl');

@$core.Deprecated('Use consensusInfoDescriptor instead')
const ConsensusInfo$json = {
  '1': 'ConsensusInfo',
  '2': [
    {'1': 'address', '3': 1, '4': 1, '5': 9, '10': 'address'},
    {'1': 'active', '3': 2, '4': 1, '5': 8, '10': 'active'},
    {'1': 'height', '3': 3, '4': 1, '5': 13, '10': 'height'},
    {'1': 'round', '3': 4, '4': 1, '5': 5, '10': 'round'},
    {
      '1': 'votes',
      '3': 5,
      '4': 3,
      '5': 11,
      '6': '.pactus.VoteInfo',
      '10': 'votes'
    },
  ],
};

/// Descriptor for `ConsensusInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List consensusInfoDescriptor = $convert.base64Decode(
    'Cg1Db25zZW5zdXNJbmZvEhgKB2FkZHJlc3MYASABKAlSB2FkZHJlc3MSFgoGYWN0aXZlGAIgAS'
    'gIUgZhY3RpdmUSFgoGaGVpZ2h0GAMgASgNUgZoZWlnaHQSFAoFcm91bmQYBCABKAVSBXJvdW5k'
    'EiYKBXZvdGVzGAUgAygLMhAucGFjdHVzLlZvdGVJbmZvUgV2b3Rlcw==');

@$core.Deprecated('Use proposalInfoDescriptor instead')
const ProposalInfo$json = {
  '1': 'ProposalInfo',
  '2': [
    {'1': 'height', '3': 1, '4': 1, '5': 13, '10': 'height'},
    {'1': 'round', '3': 2, '4': 1, '5': 5, '10': 'round'},
    {'1': 'block_data', '3': 3, '4': 1, '5': 9, '10': 'blockData'},
    {'1': 'signature', '3': 4, '4': 1, '5': 9, '10': 'signature'},
  ],
};

/// Descriptor for `ProposalInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List proposalInfoDescriptor = $convert.base64Decode(
    'CgxQcm9wb3NhbEluZm8SFgoGaGVpZ2h0GAEgASgNUgZoZWlnaHQSFAoFcm91bmQYAiABKAVSBX'
    'JvdW5kEh0KCmJsb2NrX2RhdGEYAyABKAlSCWJsb2NrRGF0YRIcCglzaWduYXR1cmUYBCABKAlS'
    'CXNpZ25hdHVyZQ==');

const $core.Map<$core.String, $core.dynamic> BlockchainServiceBase$json = {
  '1': 'Blockchain',
  '2': [
    {
      '1': 'GetBlock',
      '2': '.pactus.GetBlockRequest',
      '3': '.pactus.GetBlockResponse'
    },
    {
      '1': 'GetBlockHash',
      '2': '.pactus.GetBlockHashRequest',
      '3': '.pactus.GetBlockHashResponse'
    },
    {
      '1': 'GetBlockHeight',
      '2': '.pactus.GetBlockHeightRequest',
      '3': '.pactus.GetBlockHeightResponse'
    },
    {
      '1': 'GetBlockchainInfo',
      '2': '.pactus.GetBlockchainInfoRequest',
      '3': '.pactus.GetBlockchainInfoResponse'
    },
    {
      '1': 'GetCommitteeInfo',
      '2': '.pactus.GetCommitteeInfoRequest',
      '3': '.pactus.GetCommitteeInfoResponse'
    },
    {
      '1': 'GetConsensusInfo',
      '2': '.pactus.GetConsensusInfoRequest',
      '3': '.pactus.GetConsensusInfoResponse'
    },
    {
      '1': 'GetAccount',
      '2': '.pactus.GetAccountRequest',
      '3': '.pactus.GetAccountResponse'
    },
    {
      '1': 'GetValidator',
      '2': '.pactus.GetValidatorRequest',
      '3': '.pactus.GetValidatorResponse'
    },
    {
      '1': 'GetValidatorByNumber',
      '2': '.pactus.GetValidatorByNumberRequest',
      '3': '.pactus.GetValidatorResponse'
    },
    {
      '1': 'GetValidatorAddresses',
      '2': '.pactus.GetValidatorAddressesRequest',
      '3': '.pactus.GetValidatorAddressesResponse'
    },
    {
      '1': 'GetPublicKey',
      '2': '.pactus.GetPublicKeyRequest',
      '3': '.pactus.GetPublicKeyResponse'
    },
    {
      '1': 'GetTxPoolContent',
      '2': '.pactus.GetTxPoolContentRequest',
      '3': '.pactus.GetTxPoolContentResponse'
    },
  ],
};

@$core.Deprecated('Use blockchainServiceDescriptor instead')
const $core.Map<$core.String, $core.Map<$core.String, $core.dynamic>>
    BlockchainServiceBase$messageJson = {
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
  '.pactus.PayloadBatchTransfer': $0.PayloadBatchTransfer$json,
  '.pactus.Recipient': $0.Recipient$json,
  '.pactus.GetBlockHashRequest': GetBlockHashRequest$json,
  '.pactus.GetBlockHashResponse': GetBlockHashResponse$json,
  '.pactus.GetBlockHeightRequest': GetBlockHeightRequest$json,
  '.pactus.GetBlockHeightResponse': GetBlockHeightResponse$json,
  '.pactus.GetBlockchainInfoRequest': GetBlockchainInfoRequest$json,
  '.pactus.GetBlockchainInfoResponse': GetBlockchainInfoResponse$json,
  '.pactus.GetCommitteeInfoRequest': GetCommitteeInfoRequest$json,
  '.pactus.GetCommitteeInfoResponse': GetCommitteeInfoResponse$json,
  '.pactus.ValidatorInfo': ValidatorInfo$json,
  '.pactus.GetCommitteeInfoResponse.ProtocolVersionsEntry':
      GetCommitteeInfoResponse_ProtocolVersionsEntry$json,
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
final $typed_data.Uint8List blockchainServiceDescriptor = $convert.base64Decode(
    'CgpCbG9ja2NoYWluEj0KCEdldEJsb2NrEhcucGFjdHVzLkdldEJsb2NrUmVxdWVzdBoYLnBhY3'
    'R1cy5HZXRCbG9ja1Jlc3BvbnNlEkkKDEdldEJsb2NrSGFzaBIbLnBhY3R1cy5HZXRCbG9ja0hh'
    'c2hSZXF1ZXN0GhwucGFjdHVzLkdldEJsb2NrSGFzaFJlc3BvbnNlEk8KDkdldEJsb2NrSGVpZ2'
    'h0Eh0ucGFjdHVzLkdldEJsb2NrSGVpZ2h0UmVxdWVzdBoeLnBhY3R1cy5HZXRCbG9ja0hlaWdo'
    'dFJlc3BvbnNlElgKEUdldEJsb2NrY2hhaW5JbmZvEiAucGFjdHVzLkdldEJsb2NrY2hhaW5Jbm'
    'ZvUmVxdWVzdBohLnBhY3R1cy5HZXRCbG9ja2NoYWluSW5mb1Jlc3BvbnNlElUKEEdldENvbW1p'
    'dHRlZUluZm8SHy5wYWN0dXMuR2V0Q29tbWl0dGVlSW5mb1JlcXVlc3QaIC5wYWN0dXMuR2V0Q2'
    '9tbWl0dGVlSW5mb1Jlc3BvbnNlElUKEEdldENvbnNlbnN1c0luZm8SHy5wYWN0dXMuR2V0Q29u'
    'c2Vuc3VzSW5mb1JlcXVlc3QaIC5wYWN0dXMuR2V0Q29uc2Vuc3VzSW5mb1Jlc3BvbnNlEkMKCk'
    'dldEFjY291bnQSGS5wYWN0dXMuR2V0QWNjb3VudFJlcXVlc3QaGi5wYWN0dXMuR2V0QWNjb3Vu'
    'dFJlc3BvbnNlEkkKDEdldFZhbGlkYXRvchIbLnBhY3R1cy5HZXRWYWxpZGF0b3JSZXF1ZXN0Gh'
    'wucGFjdHVzLkdldFZhbGlkYXRvclJlc3BvbnNlElkKFEdldFZhbGlkYXRvckJ5TnVtYmVyEiMu'
    'cGFjdHVzLkdldFZhbGlkYXRvckJ5TnVtYmVyUmVxdWVzdBocLnBhY3R1cy5HZXRWYWxpZGF0b3'
    'JSZXNwb25zZRJkChVHZXRWYWxpZGF0b3JBZGRyZXNzZXMSJC5wYWN0dXMuR2V0VmFsaWRhdG9y'
    'QWRkcmVzc2VzUmVxdWVzdBolLnBhY3R1cy5HZXRWYWxpZGF0b3JBZGRyZXNzZXNSZXNwb25zZR'
    'JJCgxHZXRQdWJsaWNLZXkSGy5wYWN0dXMuR2V0UHVibGljS2V5UmVxdWVzdBocLnBhY3R1cy5H'
    'ZXRQdWJsaWNLZXlSZXNwb25zZRJVChBHZXRUeFBvb2xDb250ZW50Eh8ucGFjdHVzLkdldFR4UG'
    '9vbENvbnRlbnRSZXF1ZXN0GiAucGFjdHVzLkdldFR4UG9vbENvbnRlbnRSZXNwb25zZQ==');
