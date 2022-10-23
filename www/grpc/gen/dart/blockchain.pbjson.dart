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
@$core.Deprecated('Use accountRequestDescriptor instead')
const AccountRequest$json = const {
  '1': 'AccountRequest',
  '2': const [
    const {'1': 'address', '3': 1, '4': 1, '5': 9, '10': 'address'},
  ],
};

/// Descriptor for `AccountRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List accountRequestDescriptor = $convert.base64Decode('Cg5BY2NvdW50UmVxdWVzdBIYCgdhZGRyZXNzGAEgASgJUgdhZGRyZXNz');
@$core.Deprecated('Use accountResponseDescriptor instead')
const AccountResponse$json = const {
  '1': 'AccountResponse',
  '2': const [
    const {'1': 'account', '3': 1, '4': 1, '5': 11, '6': '.pactus.AccountInfo', '10': 'account'},
  ],
};

/// Descriptor for `AccountResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List accountResponseDescriptor = $convert.base64Decode('Cg9BY2NvdW50UmVzcG9uc2USLQoHYWNjb3VudBgBIAEoCzITLnBhY3R1cy5BY2NvdW50SW5mb1IHYWNjb3VudA==');
@$core.Deprecated('Use validatorsRequestDescriptor instead')
const ValidatorsRequest$json = const {
  '1': 'ValidatorsRequest',
};

/// Descriptor for `ValidatorsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List validatorsRequestDescriptor = $convert.base64Decode('ChFWYWxpZGF0b3JzUmVxdWVzdA==');
@$core.Deprecated('Use validatorRequestDescriptor instead')
const ValidatorRequest$json = const {
  '1': 'ValidatorRequest',
  '2': const [
    const {'1': 'address', '3': 1, '4': 1, '5': 9, '10': 'address'},
  ],
};

/// Descriptor for `ValidatorRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List validatorRequestDescriptor = $convert.base64Decode('ChBWYWxpZGF0b3JSZXF1ZXN0EhgKB2FkZHJlc3MYASABKAlSB2FkZHJlc3M=');
@$core.Deprecated('Use validatorByNumberRequestDescriptor instead')
const ValidatorByNumberRequest$json = const {
  '1': 'ValidatorByNumberRequest',
  '2': const [
    const {'1': 'number', '3': 1, '4': 1, '5': 5, '10': 'number'},
  ],
};

/// Descriptor for `ValidatorByNumberRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List validatorByNumberRequestDescriptor = $convert.base64Decode('ChhWYWxpZGF0b3JCeU51bWJlclJlcXVlc3QSFgoGbnVtYmVyGAEgASgFUgZudW1iZXI=');
@$core.Deprecated('Use validatorsResponseDescriptor instead')
const ValidatorsResponse$json = const {
  '1': 'ValidatorsResponse',
  '2': const [
    const {'1': 'validators', '3': 1, '4': 3, '5': 11, '6': '.pactus.ValidatorInfo', '10': 'validators'},
  ],
};

/// Descriptor for `ValidatorsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List validatorsResponseDescriptor = $convert.base64Decode('ChJWYWxpZGF0b3JzUmVzcG9uc2USNQoKdmFsaWRhdG9ycxgBIAMoCzIVLnBhY3R1cy5WYWxpZGF0b3JJbmZvUgp2YWxpZGF0b3Jz');
@$core.Deprecated('Use validatorResponseDescriptor instead')
const ValidatorResponse$json = const {
  '1': 'ValidatorResponse',
  '2': const [
    const {'1': 'validator', '3': 1, '4': 1, '5': 11, '6': '.pactus.ValidatorInfo', '10': 'validator'},
  ],
};

/// Descriptor for `ValidatorResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List validatorResponseDescriptor = $convert.base64Decode('ChFWYWxpZGF0b3JSZXNwb25zZRIzCgl2YWxpZGF0b3IYASABKAsyFS5wYWN0dXMuVmFsaWRhdG9ySW5mb1IJdmFsaWRhdG9y');
@$core.Deprecated('Use blockRequestDescriptor instead')
const BlockRequest$json = const {
  '1': 'BlockRequest',
  '2': const [
    const {'1': 'height', '3': 1, '4': 1, '5': 13, '10': 'height'},
    const {'1': 'verbosity', '3': 2, '4': 1, '5': 14, '6': '.pactus.BlockVerbosity', '10': 'verbosity'},
  ],
};

/// Descriptor for `BlockRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List blockRequestDescriptor = $convert.base64Decode('CgxCbG9ja1JlcXVlc3QSFgoGaGVpZ2h0GAEgASgNUgZoZWlnaHQSNAoJdmVyYm9zaXR5GAIgASgOMhYucGFjdHVzLkJsb2NrVmVyYm9zaXR5Ugl2ZXJib3NpdHk=');
@$core.Deprecated('Use blockResponseDescriptor instead')
const BlockResponse$json = const {
  '1': 'BlockResponse',
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

/// Descriptor for `BlockResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List blockResponseDescriptor = $convert.base64Decode('Cg1CbG9ja1Jlc3BvbnNlEhYKBmhlaWdodBgBIAEoDVIGaGVpZ2h0EhIKBGhhc2gYAiABKAxSBGhhc2gSEgoEZGF0YRgDIAEoDFIEZGF0YRIdCgpibG9ja190aW1lGAQgASgNUglibG9ja1RpbWUSLwoGaGVhZGVyGAUgASgLMhcucGFjdHVzLkJsb2NrSGVhZGVySW5mb1IGaGVhZGVyEjQKCXByZXZfY2VydBgGIAEoCzIXLnBhY3R1cy5DZXJ0aWZpY2F0ZUluZm9SCHByZXZDZXJ0EikKA3R4cxgHIAMoCzIXLnBhY3R1cy5UcmFuc2FjdGlvbkluZm9SA3R4cw==');
@$core.Deprecated('Use blockHashRequestDescriptor instead')
const BlockHashRequest$json = const {
  '1': 'BlockHashRequest',
  '2': const [
    const {'1': 'height', '3': 1, '4': 1, '5': 13, '10': 'height'},
  ],
};

/// Descriptor for `BlockHashRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List blockHashRequestDescriptor = $convert.base64Decode('ChBCbG9ja0hhc2hSZXF1ZXN0EhYKBmhlaWdodBgBIAEoDVIGaGVpZ2h0');
@$core.Deprecated('Use blockHashResponseDescriptor instead')
const BlockHashResponse$json = const {
  '1': 'BlockHashResponse',
  '2': const [
    const {'1': 'hash', '3': 1, '4': 1, '5': 12, '10': 'hash'},
  ],
};

/// Descriptor for `BlockHashResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List blockHashResponseDescriptor = $convert.base64Decode('ChFCbG9ja0hhc2hSZXNwb25zZRISCgRoYXNoGAEgASgMUgRoYXNo');
@$core.Deprecated('Use blockHeightRequestDescriptor instead')
const BlockHeightRequest$json = const {
  '1': 'BlockHeightRequest',
  '2': const [
    const {'1': 'hash', '3': 1, '4': 1, '5': 12, '10': 'hash'},
  ],
};

/// Descriptor for `BlockHeightRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List blockHeightRequestDescriptor = $convert.base64Decode('ChJCbG9ja0hlaWdodFJlcXVlc3QSEgoEaGFzaBgBIAEoDFIEaGFzaA==');
@$core.Deprecated('Use blockHeightResponseDescriptor instead')
const BlockHeightResponse$json = const {
  '1': 'BlockHeightResponse',
  '2': const [
    const {'1': 'height', '3': 1, '4': 1, '5': 13, '10': 'height'},
  ],
};

/// Descriptor for `BlockHeightResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List blockHeightResponseDescriptor = $convert.base64Decode('ChNCbG9ja0hlaWdodFJlc3BvbnNlEhYKBmhlaWdodBgBIAEoDVIGaGVpZ2h0');
@$core.Deprecated('Use blockchainInfoRequestDescriptor instead')
const BlockchainInfoRequest$json = const {
  '1': 'BlockchainInfoRequest',
};

/// Descriptor for `BlockchainInfoRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List blockchainInfoRequestDescriptor = $convert.base64Decode('ChVCbG9ja2NoYWluSW5mb1JlcXVlc3Q=');
@$core.Deprecated('Use blockchainInfoResponseDescriptor instead')
const BlockchainInfoResponse$json = const {
  '1': 'BlockchainInfoResponse',
  '2': const [
    const {'1': 'last_block_height', '3': 1, '4': 1, '5': 13, '10': 'lastBlockHeight'},
    const {'1': 'last_block_hash', '3': 2, '4': 1, '5': 12, '10': 'lastBlockHash'},
  ],
};

/// Descriptor for `BlockchainInfoResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List blockchainInfoResponseDescriptor = $convert.base64Decode('ChZCbG9ja2NoYWluSW5mb1Jlc3BvbnNlEioKEWxhc3RfYmxvY2tfaGVpZ2h0GAEgASgNUg9sYXN0QmxvY2tIZWlnaHQSJgoPbGFzdF9ibG9ja19oYXNoGAIgASgMUg1sYXN0QmxvY2tIYXNo');
@$core.Deprecated('Use validatorInfoDescriptor instead')
const ValidatorInfo$json = const {
  '1': 'ValidatorInfo',
  '2': const [
    const {'1': 'public_key', '3': 1, '4': 1, '5': 9, '10': 'publicKey'},
    const {'1': 'number', '3': 2, '4': 1, '5': 5, '10': 'number'},
    const {'1': 'sequence', '3': 3, '4': 1, '5': 5, '10': 'sequence'},
    const {'1': 'stake', '3': 4, '4': 1, '5': 3, '10': 'stake'},
    const {'1': 'last_bonding_height', '3': 5, '4': 1, '5': 13, '10': 'lastBondingHeight'},
    const {'1': 'last_joined_height', '3': 6, '4': 1, '5': 13, '10': 'lastJoinedHeight'},
    const {'1': 'unbonding_height', '3': 7, '4': 1, '5': 13, '10': 'unbondingHeight'},
    const {'1': 'address', '3': 8, '4': 1, '5': 9, '10': 'address'},
  ],
};

/// Descriptor for `ValidatorInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List validatorInfoDescriptor = $convert.base64Decode('Cg1WYWxpZGF0b3JJbmZvEh0KCnB1YmxpY19rZXkYASABKAlSCXB1YmxpY0tleRIWCgZudW1iZXIYAiABKAVSBm51bWJlchIaCghzZXF1ZW5jZRgDIAEoBVIIc2VxdWVuY2USFAoFc3Rha2UYBCABKANSBXN0YWtlEi4KE2xhc3RfYm9uZGluZ19oZWlnaHQYBSABKA1SEWxhc3RCb25kaW5nSGVpZ2h0EiwKEmxhc3Rfam9pbmVkX2hlaWdodBgGIAEoDVIQbGFzdEpvaW5lZEhlaWdodBIpChB1bmJvbmRpbmdfaGVpZ2h0GAcgASgNUg91bmJvbmRpbmdIZWlnaHQSGAoHYWRkcmVzcxgIIAEoCVIHYWRkcmVzcw==');
@$core.Deprecated('Use accountInfoDescriptor instead')
const AccountInfo$json = const {
  '1': 'AccountInfo',
  '2': const [
    const {'1': 'address', '3': 1, '4': 1, '5': 9, '10': 'address'},
    const {'1': 'number', '3': 2, '4': 1, '5': 5, '10': 'number'},
    const {'1': 'sequence', '3': 3, '4': 1, '5': 5, '10': 'sequence'},
    const {'1': 'Balance', '3': 4, '4': 1, '5': 3, '10': 'Balance'},
  ],
};

/// Descriptor for `AccountInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List accountInfoDescriptor = $convert.base64Decode('CgtBY2NvdW50SW5mbxIYCgdhZGRyZXNzGAEgASgJUgdhZGRyZXNzEhYKBm51bWJlchgCIAEoBVIGbnVtYmVyEhoKCHNlcXVlbmNlGAMgASgFUghzZXF1ZW5jZRIYCgdCYWxhbmNlGAQgASgDUgdCYWxhbmNl');
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
    const {'1': 'round', '3': 1, '4': 1, '5': 5, '10': 'round'},
    const {'1': 'committers', '3': 2, '4': 3, '5': 5, '10': 'committers'},
    const {'1': 'absentees', '3': 3, '4': 3, '5': 5, '10': 'absentees'},
    const {'1': 'signature', '3': 4, '4': 1, '5': 12, '10': 'signature'},
  ],
};

/// Descriptor for `CertificateInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List certificateInfoDescriptor = $convert.base64Decode('Cg9DZXJ0aWZpY2F0ZUluZm8SFAoFcm91bmQYASABKAVSBXJvdW5kEh4KCmNvbW1pdHRlcnMYAiADKAVSCmNvbW1pdHRlcnMSHAoJYWJzZW50ZWVzGAMgAygFUglhYnNlbnRlZXMSHAoJc2lnbmF0dXJlGAQgASgMUglzaWduYXR1cmU=');
const $core.Map<$core.String, $core.dynamic> BlockchainServiceBase$json = const {
  '1': 'Blockchain',
  '2': const [
    const {'1': 'GetBlock', '2': '.pactus.BlockRequest', '3': '.pactus.BlockResponse'},
    const {'1': 'GetBlockHash', '2': '.pactus.BlockHashRequest', '3': '.pactus.BlockHashResponse'},
    const {'1': 'GetBlockHeight', '2': '.pactus.BlockHeightRequest', '3': '.pactus.BlockHeightResponse'},
    const {'1': 'GetAccount', '2': '.pactus.AccountRequest', '3': '.pactus.AccountResponse'},
    const {'1': 'GetValidators', '2': '.pactus.ValidatorsRequest', '3': '.pactus.ValidatorsResponse'},
    const {'1': 'GetValidator', '2': '.pactus.ValidatorRequest', '3': '.pactus.ValidatorResponse'},
    const {'1': 'GetValidatorByNumber', '2': '.pactus.ValidatorByNumberRequest', '3': '.pactus.ValidatorResponse'},
    const {'1': 'GetBlockchainInfo', '2': '.pactus.BlockchainInfoRequest', '3': '.pactus.BlockchainInfoResponse'},
  ],
};

@$core.Deprecated('Use blockchainServiceDescriptor instead')
const $core.Map<$core.String, $core.Map<$core.String, $core.dynamic>> BlockchainServiceBase$messageJson = const {
  '.pactus.BlockRequest': BlockRequest$json,
  '.pactus.BlockResponse': BlockResponse$json,
  '.pactus.BlockHeaderInfo': BlockHeaderInfo$json,
  '.pactus.CertificateInfo': CertificateInfo$json,
  '.pactus.TransactionInfo': $0.TransactionInfo$json,
  '.pactus.PayloadSend': $0.PayloadSend$json,
  '.pactus.PayloadBond': $0.PayloadBond$json,
  '.pactus.PayloadSortition': $0.PayloadSortition$json,
  '.pactus.BlockHashRequest': BlockHashRequest$json,
  '.pactus.BlockHashResponse': BlockHashResponse$json,
  '.pactus.BlockHeightRequest': BlockHeightRequest$json,
  '.pactus.BlockHeightResponse': BlockHeightResponse$json,
  '.pactus.AccountRequest': AccountRequest$json,
  '.pactus.AccountResponse': AccountResponse$json,
  '.pactus.AccountInfo': AccountInfo$json,
  '.pactus.ValidatorsRequest': ValidatorsRequest$json,
  '.pactus.ValidatorsResponse': ValidatorsResponse$json,
  '.pactus.ValidatorInfo': ValidatorInfo$json,
  '.pactus.ValidatorRequest': ValidatorRequest$json,
  '.pactus.ValidatorResponse': ValidatorResponse$json,
  '.pactus.ValidatorByNumberRequest': ValidatorByNumberRequest$json,
  '.pactus.BlockchainInfoRequest': BlockchainInfoRequest$json,
  '.pactus.BlockchainInfoResponse': BlockchainInfoResponse$json,
};

/// Descriptor for `Blockchain`. Decode as a `google.protobuf.ServiceDescriptorProto`.
final $typed_data.Uint8List blockchainServiceDescriptor = $convert.base64Decode('CgpCbG9ja2NoYWluEjcKCEdldEJsb2NrEhQucGFjdHVzLkJsb2NrUmVxdWVzdBoVLnBhY3R1cy5CbG9ja1Jlc3BvbnNlEkMKDEdldEJsb2NrSGFzaBIYLnBhY3R1cy5CbG9ja0hhc2hSZXF1ZXN0GhkucGFjdHVzLkJsb2NrSGFzaFJlc3BvbnNlEkkKDkdldEJsb2NrSGVpZ2h0EhoucGFjdHVzLkJsb2NrSGVpZ2h0UmVxdWVzdBobLnBhY3R1cy5CbG9ja0hlaWdodFJlc3BvbnNlEj0KCkdldEFjY291bnQSFi5wYWN0dXMuQWNjb3VudFJlcXVlc3QaFy5wYWN0dXMuQWNjb3VudFJlc3BvbnNlEkYKDUdldFZhbGlkYXRvcnMSGS5wYWN0dXMuVmFsaWRhdG9yc1JlcXVlc3QaGi5wYWN0dXMuVmFsaWRhdG9yc1Jlc3BvbnNlEkMKDEdldFZhbGlkYXRvchIYLnBhY3R1cy5WYWxpZGF0b3JSZXF1ZXN0GhkucGFjdHVzLlZhbGlkYXRvclJlc3BvbnNlElMKFEdldFZhbGlkYXRvckJ5TnVtYmVyEiAucGFjdHVzLlZhbGlkYXRvckJ5TnVtYmVyUmVxdWVzdBoZLnBhY3R1cy5WYWxpZGF0b3JSZXNwb25zZRJSChFHZXRCbG9ja2NoYWluSW5mbxIdLnBhY3R1cy5CbG9ja2NoYWluSW5mb1JlcXVlc3QaHi5wYWN0dXMuQmxvY2tjaGFpbkluZm9SZXNwb25zZQ==');
