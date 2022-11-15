# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: blockchain.proto
"""Generated protocol buffer code."""
from google.protobuf.internal import builder as _builder
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


import transaction_pb2 as transaction__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x10\x62lockchain.proto\x12\x06pactus\x1a\x11transaction.proto\"-\n\x11GetAccountRequest\x12\x18\n\x07\x61\x64\x64ress\x18\x01 \x01(\tR\x07\x61\x64\x64ress\"C\n\x12GetAccountResponse\x12-\n\x07\x61\x63\x63ount\x18\x01 \x01(\x0b\x32\x13.pactus.AccountInfoR\x07\x61\x63\x63ount\"\x16\n\x14GetValidatorsRequest\"/\n\x13GetValidatorRequest\x12\x18\n\x07\x61\x64\x64ress\x18\x01 \x01(\tR\x07\x61\x64\x64ress\"5\n\x1bGetValidatorByNumberRequest\x12\x16\n\x06number\x18\x01 \x01(\x05R\x06number\"N\n\x15GetValidatorsResponse\x12\x35\n\nvalidators\x18\x01 \x03(\x0b\x32\x15.pactus.ValidatorInfoR\nvalidators\"K\n\x14GetValidatorResponse\x12\x33\n\tvalidator\x18\x01 \x01(\x0b\x32\x15.pactus.ValidatorInfoR\tvalidator\"_\n\x0fGetBlockRequest\x12\x16\n\x06height\x18\x01 \x01(\rR\x06height\x12\x34\n\tverbosity\x18\x02 \x01(\x0e\x32\x16.pactus.BlockVerbosityR\tverbosity\"\x83\x02\n\x10GetBlockResponse\x12\x16\n\x06height\x18\x01 \x01(\rR\x06height\x12\x12\n\x04hash\x18\x02 \x01(\x0cR\x04hash\x12\x12\n\x04\x64\x61ta\x18\x03 \x01(\x0cR\x04\x64\x61ta\x12\x1d\n\nblock_time\x18\x04 \x01(\rR\tblockTime\x12/\n\x06header\x18\x05 \x01(\x0b\x32\x17.pactus.BlockHeaderInfoR\x06header\x12\x34\n\tprev_cert\x18\x06 \x01(\x0b\x32\x17.pactus.CertificateInfoR\x08prevCert\x12)\n\x03txs\x18\x07 \x03(\x0b\x32\x17.pactus.TransactionInfoR\x03txs\"-\n\x13GetBlockHashRequest\x12\x16\n\x06height\x18\x01 \x01(\rR\x06height\"*\n\x14GetBlockHashResponse\x12\x12\n\x04hash\x18\x01 \x01(\x0cR\x04hash\"+\n\x15GetBlockHeightRequest\x12\x12\n\x04hash\x18\x01 \x01(\x0cR\x04hash\"0\n\x16GetBlockHeightResponse\x12\x16\n\x06height\x18\x01 \x01(\rR\x06height\"\x1a\n\x18GetBlockchainInfoRequest\"\x83\x02\n\x19GetBlockchainInfoResponse\x12*\n\x11last_block_height\x18\x01 \x01(\rR\x0flastBlockHeight\x12&\n\x0flast_block_hash\x18\x02 \x01(\x0cR\rlastBlockHash\x12\x1f\n\x0btotal_power\x18\x03 \x01(\x03R\ntotalPower\x12\'\n\x0f\x63ommittee_power\x18\x04 \x01(\x03R\x0e\x63ommitteePower\x12H\n\x14\x63ommittee_validators\x18\x05 \x03(\x0b\x32\x15.pactus.ValidatorInfoR\x13\x63ommitteeValidators\"\x19\n\x17GetConsensusInfoRequest\"p\n\x18GetConsensusInfoResponse\x12\x16\n\x06height\x18\x01 \x01(\rR\x06height\x12\x14\n\x05round\x18\x02 \x01(\x05R\x05round\x12&\n\x05votes\x18\x03 \x03(\x0b\x32\x10.pactus.VoteInfoR\x05votes\"\x9b\x02\n\rValidatorInfo\x12\x1d\n\npublic_key\x18\x01 \x01(\tR\tpublicKey\x12\x16\n\x06number\x18\x02 \x01(\x05R\x06number\x12\x1a\n\x08sequence\x18\x03 \x01(\x05R\x08sequence\x12\x14\n\x05stake\x18\x04 \x01(\x03R\x05stake\x12.\n\x13last_bonding_height\x18\x05 \x01(\rR\x11lastBondingHeight\x12,\n\x12last_joined_height\x18\x06 \x01(\rR\x10lastJoinedHeight\x12)\n\x10unbonding_height\x18\x07 \x01(\rR\x0funbondingHeight\x12\x18\n\x07\x61\x64\x64ress\x18\x08 \x01(\tR\x07\x61\x64\x64ress\"u\n\x0b\x41\x63\x63ountInfo\x12\x18\n\x07\x61\x64\x64ress\x18\x01 \x01(\tR\x07\x61\x64\x64ress\x12\x16\n\x06number\x18\x02 \x01(\x05R\x06number\x12\x1a\n\x08sequence\x18\x03 \x01(\x05R\x08sequence\x12\x18\n\x07\x42\x61lance\x18\x04 \x01(\x03R\x07\x42\x61lance\"\xc4\x01\n\x0f\x42lockHeaderInfo\x12\x18\n\x07version\x18\x01 \x01(\x05R\x07version\x12&\n\x0fprev_block_hash\x18\x02 \x01(\x0cR\rprevBlockHash\x12\x1d\n\nstate_root\x18\x03 \x01(\x0cR\tstateRoot\x12%\n\x0esortition_seed\x18\x04 \x01(\x0cR\rsortitionSeed\x12)\n\x10proposer_address\x18\x05 \x01(\tR\x0fproposerAddress\"\x83\x01\n\x0f\x43\x65rtificateInfo\x12\x14\n\x05round\x18\x01 \x01(\x05R\x05round\x12\x1e\n\ncommitters\x18\x02 \x03(\x05R\ncommitters\x12\x1c\n\tabsentees\x18\x03 \x03(\x05R\tabsentees\x12\x1c\n\tsignature\x18\x04 \x01(\x0cR\tsignature\"{\n\x08VoteInfo\x12$\n\x04type\x18\x01 \x01(\x0e\x32\x10.pactus.VoteTypeR\x04type\x12\x14\n\x05voter\x18\x02 \x01(\tR\x05voter\x12\x1d\n\nblock_hash\x18\x03 \x01(\x0cR\tblockHash\x12\x14\n\x05round\x18\x04 \x01(\x05R\x05round*H\n\x0e\x42lockVerbosity\x12\x0e\n\nBLOCK_DATA\x10\x00\x12\x0e\n\nBLOCK_INFO\x10\x01\x12\x16\n\x12\x42LOCK_TRANSACTIONS\x10\x02*J\n\x08VoteType\x12\x10\n\x0cVOTE_PREPARE\x10\x00\x12\x12\n\x0eVOTE_PRECOMMIT\x10\x01\x12\x18\n\x14VOTE_CHANGE_PROPOSER\x10\x02\x32\xd1\x05\n\nBlockchain\x12=\n\x08GetBlock\x12\x17.pactus.GetBlockRequest\x1a\x18.pactus.GetBlockResponse\x12I\n\x0cGetBlockHash\x12\x1b.pactus.GetBlockHashRequest\x1a\x1c.pactus.GetBlockHashResponse\x12O\n\x0eGetBlockHeight\x12\x1d.pactus.GetBlockHeightRequest\x1a\x1e.pactus.GetBlockHeightResponse\x12X\n\x11GetBlockchainInfo\x12 .pactus.GetBlockchainInfoRequest\x1a!.pactus.GetBlockchainInfoResponse\x12U\n\x10GetConsensusInfo\x12\x1f.pactus.GetConsensusInfoRequest\x1a .pactus.GetConsensusInfoResponse\x12\x43\n\nGetAccount\x12\x19.pactus.GetAccountRequest\x1a\x1a.pactus.GetAccountResponse\x12I\n\x0cGetValidator\x12\x1b.pactus.GetValidatorRequest\x1a\x1c.pactus.GetValidatorResponse\x12Y\n\x14GetValidatorByNumber\x12#.pactus.GetValidatorByNumberRequest\x1a\x1c.pactus.GetValidatorResponse\x12L\n\rGetValidators\x12\x1c.pactus.GetValidatorsRequest\x1a\x1d.pactus.GetValidatorsResponseBE\n\x11pactus.blockchainZ0github.com/pactus-project/pactus/www/grpc/pactusb\x06proto3')

_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, globals())
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'blockchain_pb2', globals())
if _descriptor._USE_C_DESCRIPTORS == False:

  DESCRIPTOR._options = None
  DESCRIPTOR._serialized_options = b'\n\021pactus.blockchainZ0github.com/pactus-project/pactus/www/grpc/pactus'
  _BLOCKVERBOSITY._serialized_start=2287
  _BLOCKVERBOSITY._serialized_end=2359
  _VOTETYPE._serialized_start=2361
  _VOTETYPE._serialized_end=2435
  _GETACCOUNTREQUEST._serialized_start=47
  _GETACCOUNTREQUEST._serialized_end=92
  _GETACCOUNTRESPONSE._serialized_start=94
  _GETACCOUNTRESPONSE._serialized_end=161
  _GETVALIDATORSREQUEST._serialized_start=163
  _GETVALIDATORSREQUEST._serialized_end=185
  _GETVALIDATORREQUEST._serialized_start=187
  _GETVALIDATORREQUEST._serialized_end=234
  _GETVALIDATORBYNUMBERREQUEST._serialized_start=236
  _GETVALIDATORBYNUMBERREQUEST._serialized_end=289
  _GETVALIDATORSRESPONSE._serialized_start=291
  _GETVALIDATORSRESPONSE._serialized_end=369
  _GETVALIDATORRESPONSE._serialized_start=371
  _GETVALIDATORRESPONSE._serialized_end=446
  _GETBLOCKREQUEST._serialized_start=448
  _GETBLOCKREQUEST._serialized_end=543
  _GETBLOCKRESPONSE._serialized_start=546
  _GETBLOCKRESPONSE._serialized_end=805
  _GETBLOCKHASHREQUEST._serialized_start=807
  _GETBLOCKHASHREQUEST._serialized_end=852
  _GETBLOCKHASHRESPONSE._serialized_start=854
  _GETBLOCKHASHRESPONSE._serialized_end=896
  _GETBLOCKHEIGHTREQUEST._serialized_start=898
  _GETBLOCKHEIGHTREQUEST._serialized_end=941
  _GETBLOCKHEIGHTRESPONSE._serialized_start=943
  _GETBLOCKHEIGHTRESPONSE._serialized_end=991
  _GETBLOCKCHAININFOREQUEST._serialized_start=993
  _GETBLOCKCHAININFOREQUEST._serialized_end=1019
  _GETBLOCKCHAININFORESPONSE._serialized_start=1022
  _GETBLOCKCHAININFORESPONSE._serialized_end=1281
  _GETCONSENSUSINFOREQUEST._serialized_start=1283
  _GETCONSENSUSINFOREQUEST._serialized_end=1308
  _GETCONSENSUSINFORESPONSE._serialized_start=1310
  _GETCONSENSUSINFORESPONSE._serialized_end=1422
  _VALIDATORINFO._serialized_start=1425
  _VALIDATORINFO._serialized_end=1708
  _ACCOUNTINFO._serialized_start=1710
  _ACCOUNTINFO._serialized_end=1827
  _BLOCKHEADERINFO._serialized_start=1830
  _BLOCKHEADERINFO._serialized_end=2026
  _CERTIFICATEINFO._serialized_start=2029
  _CERTIFICATEINFO._serialized_end=2160
  _VOTEINFO._serialized_start=2162
  _VOTEINFO._serialized_end=2285
  _BLOCKCHAIN._serialized_start=2438
  _BLOCKCHAIN._serialized_end=3159
# @@protoc_insertion_point(module_scope)
