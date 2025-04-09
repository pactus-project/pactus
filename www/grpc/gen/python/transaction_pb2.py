# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: transaction.proto
# Protobuf Python Version: 6.30.2
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import runtime_version as _runtime_version
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
_runtime_version.ValidateProtobufRuntimeVersion(
    _runtime_version.Domain.PUBLIC,
    6,
    30,
    2,
    '',
    'transaction.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x11transaction.proto\x12\x06pactus\"c\n\x15GetTransactionRequest\x12\x0e\n\x02id\x18\x01 \x01(\tR\x02id\x12:\n\tverbosity\x18\x02 \x01(\x0e\x32\x1c.pactus.TransactionVerbosityR\tverbosity\"\x95\x01\n\x16GetTransactionResponse\x12!\n\x0c\x62lock_height\x18\x01 \x01(\rR\x0b\x62lockHeight\x12\x1d\n\nblock_time\x18\x02 \x01(\rR\tblockTime\x12\x39\n\x0btransaction\x18\x03 \x01(\x0b\x32\x17.pactus.TransactionInfoR\x0btransaction\"\x88\x01\n\x13\x43\x61lculateFeeRequest\x12\x16\n\x06\x61mount\x18\x01 \x01(\x03R\x06\x61mount\x12\x36\n\x0cpayload_type\x18\x02 \x01(\x0e\x32\x13.pactus.PayloadTypeR\x0bpayloadType\x12!\n\x0c\x66ixed_amount\x18\x03 \x01(\x08R\x0b\x66ixedAmount\"@\n\x14\x43\x61lculateFeeResponse\x12\x16\n\x06\x61mount\x18\x01 \x01(\x03R\x06\x61mount\x12\x10\n\x03\x66\x65\x65\x18\x02 \x01(\x03R\x03\x66\x65\x65\"S\n\x1b\x42roadcastTransactionRequest\x12\x34\n\x16signed_raw_transaction\x18\x01 \x01(\tR\x14signedRawTransaction\".\n\x1c\x42roadcastTransactionResponse\x12\x0e\n\x02id\x18\x01 \x01(\tR\x02id\"\xb1\x01\n GetRawTransferTransactionRequest\x12\x1b\n\tlock_time\x18\x01 \x01(\rR\x08lockTime\x12\x16\n\x06sender\x18\x02 \x01(\tR\x06sender\x12\x1a\n\x08receiver\x18\x03 \x01(\tR\x08receiver\x12\x16\n\x06\x61mount\x18\x04 \x01(\x03R\x06\x61mount\x12\x10\n\x03\x66\x65\x65\x18\x05 \x01(\x03R\x03\x66\x65\x65\x12\x12\n\x04memo\x18\x06 \x01(\tR\x04memo\"\xca\x01\n\x1cGetRawBondTransactionRequest\x12\x1b\n\tlock_time\x18\x01 \x01(\rR\x08lockTime\x12\x16\n\x06sender\x18\x02 \x01(\tR\x06sender\x12\x1a\n\x08receiver\x18\x03 \x01(\tR\x08receiver\x12\x14\n\x05stake\x18\x04 \x01(\x03R\x05stake\x12\x1d\n\npublic_key\x18\x05 \x01(\tR\tpublicKey\x12\x10\n\x03\x66\x65\x65\x18\x06 \x01(\x03R\x03\x66\x65\x65\x12\x12\n\x04memo\x18\x07 \x01(\tR\x04memo\"~\n\x1eGetRawUnbondTransactionRequest\x12\x1b\n\tlock_time\x18\x01 \x01(\rR\x08lockTime\x12+\n\x11validator_address\x18\x03 \x01(\tR\x10validatorAddress\x12\x12\n\x04memo\x18\x04 \x01(\tR\x04memo\"\xd3\x01\n GetRawWithdrawTransactionRequest\x12\x1b\n\tlock_time\x18\x01 \x01(\rR\x08lockTime\x12+\n\x11validator_address\x18\x02 \x01(\tR\x10validatorAddress\x12\'\n\x0f\x61\x63\x63ount_address\x18\x03 \x01(\tR\x0e\x61\x63\x63ountAddress\x12\x16\n\x06\x61mount\x18\x04 \x01(\x03R\x06\x61mount\x12\x10\n\x03\x66\x65\x65\x18\x05 \x01(\x03R\x03\x66\x65\x65\x12\x12\n\x04memo\x18\x06 \x01(\tR\x04memo\"T\n\x19GetRawTransactionResponse\x12\'\n\x0fraw_transaction\x18\x01 \x01(\tR\x0erawTransaction\x12\x0e\n\x02id\x18\x02 \x01(\tR\x02id\"]\n\x0fPayloadTransfer\x12\x16\n\x06sender\x18\x01 \x01(\tR\x06sender\x12\x1a\n\x08receiver\x18\x02 \x01(\tR\x08receiver\x12\x16\n\x06\x61mount\x18\x03 \x01(\x03R\x06\x61mount\"v\n\x0bPayloadBond\x12\x16\n\x06sender\x18\x01 \x01(\tR\x06sender\x12\x1a\n\x08receiver\x18\x02 \x01(\tR\x08receiver\x12\x14\n\x05stake\x18\x03 \x01(\x03R\x05stake\x12\x1d\n\npublic_key\x18\x04 \x01(\tR\tpublicKey\"B\n\x10PayloadSortition\x12\x18\n\x07\x61\x64\x64ress\x18\x01 \x01(\tR\x07\x61\x64\x64ress\x12\x14\n\x05proof\x18\x02 \x01(\tR\x05proof\"-\n\rPayloadUnbond\x12\x1c\n\tvalidator\x18\x01 \x01(\tR\tvalidator\"\x7f\n\x0fPayloadWithdraw\x12+\n\x11validator_address\x18\x01 \x01(\tR\x10validatorAddress\x12\'\n\x0f\x61\x63\x63ount_address\x18\x02 \x01(\tR\x0e\x61\x63\x63ountAddress\x12\x16\n\x06\x61mount\x18\x03 \x01(\x03R\x06\x61mount\"\xac\x04\n\x0fTransactionInfo\x12\x0e\n\x02id\x18\x01 \x01(\tR\x02id\x12\x12\n\x04\x64\x61ta\x18\x02 \x01(\tR\x04\x64\x61ta\x12\x18\n\x07version\x18\x03 \x01(\x05R\x07version\x12\x1b\n\tlock_time\x18\x04 \x01(\rR\x08lockTime\x12\x14\n\x05value\x18\x05 \x01(\x03R\x05value\x12\x10\n\x03\x66\x65\x65\x18\x06 \x01(\x03R\x03\x66\x65\x65\x12\x36\n\x0cpayload_type\x18\x07 \x01(\x0e\x32\x13.pactus.PayloadTypeR\x0bpayloadType\x12\x35\n\x08transfer\x18\x1e \x01(\x0b\x32\x17.pactus.PayloadTransferH\x00R\x08transfer\x12)\n\x04\x62ond\x18\x1f \x01(\x0b\x32\x13.pactus.PayloadBondH\x00R\x04\x62ond\x12\x38\n\tsortition\x18  \x01(\x0b\x32\x18.pactus.PayloadSortitionH\x00R\tsortition\x12/\n\x06unbond\x18! \x01(\x0b\x32\x15.pactus.PayloadUnbondH\x00R\x06unbond\x12\x35\n\x08withdraw\x18\" \x01(\x0b\x32\x17.pactus.PayloadWithdrawH\x00R\x08withdraw\x12\x12\n\x04memo\x18\x08 \x01(\tR\x04memo\x12\x1d\n\npublic_key\x18\t \x01(\tR\tpublicKey\x12\x1c\n\tsignature\x18\n \x01(\tR\tsignatureB\t\n\x07payload\"F\n\x1b\x44\x65\x63odeRawTransactionRequest\x12\'\n\x0fraw_transaction\x18\x01 \x01(\tR\x0erawTransaction\"Y\n\x1c\x44\x65\x63odeRawTransactionResponse\x12\x39\n\x0btransaction\x18\x01 \x01(\x0b\x32\x17.pactus.TransactionInfoR\x0btransaction*\xad\x01\n\x0bPayloadType\x12\x1c\n\x18PAYLOAD_TYPE_UNSPECIFIED\x10\x00\x12\x19\n\x15PAYLOAD_TYPE_TRANSFER\x10\x01\x12\x15\n\x11PAYLOAD_TYPE_BOND\x10\x02\x12\x1a\n\x16PAYLOAD_TYPE_SORTITION\x10\x03\x12\x17\n\x13PAYLOAD_TYPE_UNBOND\x10\x04\x12\x19\n\x15PAYLOAD_TYPE_WITHDRAW\x10\x05*V\n\x14TransactionVerbosity\x12\x1e\n\x1aTRANSACTION_VERBOSITY_DATA\x10\x00\x12\x1e\n\x1aTRANSACTION_VERBOSITY_INFO\x10\x01\x32\x92\x06\n\x12TransactionService\x12O\n\x0eGetTransaction\x12\x1d.pactus.GetTransactionRequest\x1a\x1e.pactus.GetTransactionResponse\x12I\n\x0c\x43\x61lculateFee\x12\x1b.pactus.CalculateFeeRequest\x1a\x1c.pactus.CalculateFeeResponse\x12\x61\n\x14\x42roadcastTransaction\x12#.pactus.BroadcastTransactionRequest\x1a$.pactus.BroadcastTransactionResponse\x12h\n\x19GetRawTransferTransaction\x12(.pactus.GetRawTransferTransactionRequest\x1a!.pactus.GetRawTransactionResponse\x12`\n\x15GetRawBondTransaction\x12$.pactus.GetRawBondTransactionRequest\x1a!.pactus.GetRawTransactionResponse\x12\x64\n\x17GetRawUnbondTransaction\x12&.pactus.GetRawUnbondTransactionRequest\x1a!.pactus.GetRawTransactionResponse\x12h\n\x19GetRawWithdrawTransaction\x12(.pactus.GetRawWithdrawTransactionRequest\x1a!.pactus.GetRawTransactionResponse\x12\x61\n\x14\x44\x65\x63odeRawTransaction\x12#.pactus.DecodeRawTransactionRequest\x1a$.pactus.DecodeRawTransactionResponseB:\n\x06pactusZ0github.com/pactus-project/pactus/www/grpc/pactusb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'transaction_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'\n\006pactusZ0github.com/pactus-project/pactus/www/grpc/pactus'
  _globals['_PAYLOADTYPE']._serialized_start=2615
  _globals['_PAYLOADTYPE']._serialized_end=2788
  _globals['_TRANSACTIONVERBOSITY']._serialized_start=2790
  _globals['_TRANSACTIONVERBOSITY']._serialized_end=2876
  _globals['_GETTRANSACTIONREQUEST']._serialized_start=29
  _globals['_GETTRANSACTIONREQUEST']._serialized_end=128
  _globals['_GETTRANSACTIONRESPONSE']._serialized_start=131
  _globals['_GETTRANSACTIONRESPONSE']._serialized_end=280
  _globals['_CALCULATEFEEREQUEST']._serialized_start=283
  _globals['_CALCULATEFEEREQUEST']._serialized_end=419
  _globals['_CALCULATEFEERESPONSE']._serialized_start=421
  _globals['_CALCULATEFEERESPONSE']._serialized_end=485
  _globals['_BROADCASTTRANSACTIONREQUEST']._serialized_start=487
  _globals['_BROADCASTTRANSACTIONREQUEST']._serialized_end=570
  _globals['_BROADCASTTRANSACTIONRESPONSE']._serialized_start=572
  _globals['_BROADCASTTRANSACTIONRESPONSE']._serialized_end=618
  _globals['_GETRAWTRANSFERTRANSACTIONREQUEST']._serialized_start=621
  _globals['_GETRAWTRANSFERTRANSACTIONREQUEST']._serialized_end=798
  _globals['_GETRAWBONDTRANSACTIONREQUEST']._serialized_start=801
  _globals['_GETRAWBONDTRANSACTIONREQUEST']._serialized_end=1003
  _globals['_GETRAWUNBONDTRANSACTIONREQUEST']._serialized_start=1005
  _globals['_GETRAWUNBONDTRANSACTIONREQUEST']._serialized_end=1131
  _globals['_GETRAWWITHDRAWTRANSACTIONREQUEST']._serialized_start=1134
  _globals['_GETRAWWITHDRAWTRANSACTIONREQUEST']._serialized_end=1345
  _globals['_GETRAWTRANSACTIONRESPONSE']._serialized_start=1347
  _globals['_GETRAWTRANSACTIONRESPONSE']._serialized_end=1431
  _globals['_PAYLOADTRANSFER']._serialized_start=1433
  _globals['_PAYLOADTRANSFER']._serialized_end=1526
  _globals['_PAYLOADBOND']._serialized_start=1528
  _globals['_PAYLOADBOND']._serialized_end=1646
  _globals['_PAYLOADSORTITION']._serialized_start=1648
  _globals['_PAYLOADSORTITION']._serialized_end=1714
  _globals['_PAYLOADUNBOND']._serialized_start=1716
  _globals['_PAYLOADUNBOND']._serialized_end=1761
  _globals['_PAYLOADWITHDRAW']._serialized_start=1763
  _globals['_PAYLOADWITHDRAW']._serialized_end=1890
  _globals['_TRANSACTIONINFO']._serialized_start=1893
  _globals['_TRANSACTIONINFO']._serialized_end=2449
  _globals['_DECODERAWTRANSACTIONREQUEST']._serialized_start=2451
  _globals['_DECODERAWTRANSACTIONREQUEST']._serialized_end=2521
  _globals['_DECODERAWTRANSACTIONRESPONSE']._serialized_start=2523
  _globals['_DECODERAWTRANSACTIONRESPONSE']._serialized_end=2612
  _globals['_TRANSACTIONSERVICE']._serialized_start=2879
  _globals['_TRANSACTIONSERVICE']._serialized_end=3665
# @@protoc_insertion_point(module_scope)
