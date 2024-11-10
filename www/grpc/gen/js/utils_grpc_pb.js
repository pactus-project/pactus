// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('grpc');
var utils_pb = require('./utils_pb.js');

function serialize_pactus_BLSPublicKeyAggregationRequest(arg) {
  if (!(arg instanceof utils_pb.BLSPublicKeyAggregationRequest)) {
    throw new Error('Expected argument of type pactus.BLSPublicKeyAggregationRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_BLSPublicKeyAggregationRequest(buffer_arg) {
  return utils_pb.BLSPublicKeyAggregationRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_BLSPublicKeyAggregationResponse(arg) {
  if (!(arg instanceof utils_pb.BLSPublicKeyAggregationResponse)) {
    throw new Error('Expected argument of type pactus.BLSPublicKeyAggregationResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_BLSPublicKeyAggregationResponse(buffer_arg) {
  return utils_pb.BLSPublicKeyAggregationResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_BLSSignatureAggregationRequest(arg) {
  if (!(arg instanceof utils_pb.BLSSignatureAggregationRequest)) {
    throw new Error('Expected argument of type pactus.BLSSignatureAggregationRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_BLSSignatureAggregationRequest(buffer_arg) {
  return utils_pb.BLSSignatureAggregationRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_BLSSignatureAggregationResponse(arg) {
  if (!(arg instanceof utils_pb.BLSSignatureAggregationResponse)) {
    throw new Error('Expected argument of type pactus.BLSSignatureAggregationResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_BLSSignatureAggregationResponse(buffer_arg) {
  return utils_pb.BLSSignatureAggregationResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_SignMessageWithPrivateKeyRequest(arg) {
  if (!(arg instanceof utils_pb.SignMessageWithPrivateKeyRequest)) {
    throw new Error('Expected argument of type pactus.SignMessageWithPrivateKeyRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_SignMessageWithPrivateKeyRequest(buffer_arg) {
  return utils_pb.SignMessageWithPrivateKeyRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_SignMessageWithPrivateKeyResponse(arg) {
  if (!(arg instanceof utils_pb.SignMessageWithPrivateKeyResponse)) {
    throw new Error('Expected argument of type pactus.SignMessageWithPrivateKeyResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_SignMessageWithPrivateKeyResponse(buffer_arg) {
  return utils_pb.SignMessageWithPrivateKeyResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_VerifyMessageRequest(arg) {
  if (!(arg instanceof utils_pb.VerifyMessageRequest)) {
    throw new Error('Expected argument of type pactus.VerifyMessageRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_VerifyMessageRequest(buffer_arg) {
  return utils_pb.VerifyMessageRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_VerifyMessageResponse(arg) {
  if (!(arg instanceof utils_pb.VerifyMessageResponse)) {
    throw new Error('Expected argument of type pactus.VerifyMessageResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_VerifyMessageResponse(buffer_arg) {
  return utils_pb.VerifyMessageResponse.deserializeBinary(new Uint8Array(buffer_arg));
}


// Utils service defines RPC methods for utility functions such as message
// signing and verification.
var UtilsService = exports.UtilsService = {
  // SignMessageWithPrivateKey signs message with provided private key.
signMessageWithPrivateKey: {
    path: '/pactus.Utils/SignMessageWithPrivateKey',
    requestStream: false,
    responseStream: false,
    requestType: utils_pb.SignMessageWithPrivateKeyRequest,
    responseType: utils_pb.SignMessageWithPrivateKeyResponse,
    requestSerialize: serialize_pactus_SignMessageWithPrivateKeyRequest,
    requestDeserialize: deserialize_pactus_SignMessageWithPrivateKeyRequest,
    responseSerialize: serialize_pactus_SignMessageWithPrivateKeyResponse,
    responseDeserialize: deserialize_pactus_SignMessageWithPrivateKeyResponse,
  },
  // VerifyMessage verifies signature with public key and message.
verifyMessage: {
    path: '/pactus.Utils/VerifyMessage',
    requestStream: false,
    responseStream: false,
    requestType: utils_pb.VerifyMessageRequest,
    responseType: utils_pb.VerifyMessageResponse,
    requestSerialize: serialize_pactus_VerifyMessageRequest,
    requestDeserialize: deserialize_pactus_VerifyMessageRequest,
    responseSerialize: serialize_pactus_VerifyMessageResponse,
    responseDeserialize: deserialize_pactus_VerifyMessageResponse,
  },
  // BLSPublicKeyAggregation aggregates bls public keys.
bLSPublicKeyAggregation: {
    path: '/pactus.Utils/BLSPublicKeyAggregation',
    requestStream: false,
    responseStream: false,
    requestType: utils_pb.BLSPublicKeyAggregationRequest,
    responseType: utils_pb.BLSPublicKeyAggregationResponse,
    requestSerialize: serialize_pactus_BLSPublicKeyAggregationRequest,
    requestDeserialize: deserialize_pactus_BLSPublicKeyAggregationRequest,
    responseSerialize: serialize_pactus_BLSPublicKeyAggregationResponse,
    responseDeserialize: deserialize_pactus_BLSPublicKeyAggregationResponse,
  },
  // BLSSignatureAggregation aggregates bls signatures.
bLSSignatureAggregation: {
    path: '/pactus.Utils/BLSSignatureAggregation',
    requestStream: false,
    responseStream: false,
    requestType: utils_pb.BLSSignatureAggregationRequest,
    responseType: utils_pb.BLSSignatureAggregationResponse,
    requestSerialize: serialize_pactus_BLSSignatureAggregationRequest,
    requestDeserialize: deserialize_pactus_BLSSignatureAggregationRequest,
    responseSerialize: serialize_pactus_BLSSignatureAggregationResponse,
    responseDeserialize: deserialize_pactus_BLSSignatureAggregationResponse,
  },
};

exports.UtilsClient = grpc.makeGenericClientConstructor(UtilsService);
