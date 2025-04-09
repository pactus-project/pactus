# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: utils.proto
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
    'utils.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x0butils.proto\x12\x06pactus\"]\n SignMessageWithPrivateKeyRequest\x12\x1f\n\x0bprivate_key\x18\x01 \x01(\tR\nprivateKey\x12\x18\n\x07message\x18\x02 \x01(\tR\x07message\"A\n!SignMessageWithPrivateKeyResponse\x12\x1c\n\tsignature\x18\x01 \x01(\tR\tsignature\"m\n\x14VerifyMessageRequest\x12\x18\n\x07message\x18\x01 \x01(\tR\x07message\x12\x1c\n\tsignature\x18\x02 \x01(\tR\tsignature\x12\x1d\n\npublic_key\x18\x03 \x01(\tR\tpublicKey\"2\n\x15VerifyMessageResponse\x12\x19\n\x08is_valid\x18\x01 \x01(\x08R\x07isValid\">\n\x1bPublicKeyAggregationRequest\x12\x1f\n\x0bpublic_keys\x18\x01 \x03(\tR\npublicKeys\"W\n\x1cPublicKeyAggregationResponse\x12\x1d\n\npublic_key\x18\x01 \x01(\tR\tpublicKey\x12\x18\n\x07\x61\x64\x64ress\x18\x02 \x01(\tR\x07\x61\x64\x64ress\"=\n\x1bSignatureAggregationRequest\x12\x1e\n\nsignatures\x18\x01 \x03(\tR\nsignatures\"<\n\x1cSignatureAggregationResponse\x12\x1c\n\tsignature\x18\x01 \x01(\tR\tsignature2\x94\x03\n\x0cUtilsService\x12p\n\x19SignMessageWithPrivateKey\x12(.pactus.SignMessageWithPrivateKeyRequest\x1a).pactus.SignMessageWithPrivateKeyResponse\x12L\n\rVerifyMessage\x12\x1c.pactus.VerifyMessageRequest\x1a\x1d.pactus.VerifyMessageResponse\x12\x61\n\x14PublicKeyAggregation\x12#.pactus.PublicKeyAggregationRequest\x1a$.pactus.PublicKeyAggregationResponse\x12\x61\n\x14SignatureAggregation\x12#.pactus.SignatureAggregationRequest\x1a$.pactus.SignatureAggregationResponseB:\n\x06pactusZ0github.com/pactus-project/pactus/www/grpc/pactusb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'utils_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'\n\006pactusZ0github.com/pactus-project/pactus/www/grpc/pactus'
  _globals['_SIGNMESSAGEWITHPRIVATEKEYREQUEST']._serialized_start=23
  _globals['_SIGNMESSAGEWITHPRIVATEKEYREQUEST']._serialized_end=116
  _globals['_SIGNMESSAGEWITHPRIVATEKEYRESPONSE']._serialized_start=118
  _globals['_SIGNMESSAGEWITHPRIVATEKEYRESPONSE']._serialized_end=183
  _globals['_VERIFYMESSAGEREQUEST']._serialized_start=185
  _globals['_VERIFYMESSAGEREQUEST']._serialized_end=294
  _globals['_VERIFYMESSAGERESPONSE']._serialized_start=296
  _globals['_VERIFYMESSAGERESPONSE']._serialized_end=346
  _globals['_PUBLICKEYAGGREGATIONREQUEST']._serialized_start=348
  _globals['_PUBLICKEYAGGREGATIONREQUEST']._serialized_end=410
  _globals['_PUBLICKEYAGGREGATIONRESPONSE']._serialized_start=412
  _globals['_PUBLICKEYAGGREGATIONRESPONSE']._serialized_end=499
  _globals['_SIGNATUREAGGREGATIONREQUEST']._serialized_start=501
  _globals['_SIGNATUREAGGREGATIONREQUEST']._serialized_end=562
  _globals['_SIGNATUREAGGREGATIONRESPONSE']._serialized_start=564
  _globals['_SIGNATUREAGGREGATIONRESPONSE']._serialized_end=624
  _globals['_UTILSSERVICE']._serialized_start=627
  _globals['_UTILSSERVICE']._serialized_end=1031
# @@protoc_insertion_point(module_scope)
