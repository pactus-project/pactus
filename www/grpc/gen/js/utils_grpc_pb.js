// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('grpc');
var utils_pb = require('./utils_pb.js');

function serialize_pactus_PublicKeyAggregationRequest(arg) {
  if (!(arg instanceof utils_pb.PublicKeyAggregationRequest)) {
    throw new Error('Expected argument of type pactus.PublicKeyAggregationRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_PublicKeyAggregationRequest(buffer_arg) {
  return utils_pb.PublicKeyAggregationRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_PublicKeyAggregationResponse(arg) {
  if (!(arg instanceof utils_pb.PublicKeyAggregationResponse)) {
    throw new Error('Expected argument of type pactus.PublicKeyAggregationResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_PublicKeyAggregationResponse(buffer_arg) {
  return utils_pb.PublicKeyAggregationResponse.deserializeBinary(new Uint8Array(buffer_arg));
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

function serialize_pactus_SignatureAggregationRequest(arg) {
  if (!(arg instanceof utils_pb.SignatureAggregationRequest)) {
    throw new Error('Expected argument of type pactus.SignatureAggregationRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_SignatureAggregationRequest(buffer_arg) {
  return utils_pb.SignatureAggregationRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_SignatureAggregationResponse(arg) {
  if (!(arg instanceof utils_pb.SignatureAggregationResponse)) {
    throw new Error('Expected argument of type pactus.SignatureAggregationResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_SignatureAggregationResponse(buffer_arg) {
  return utils_pb.SignatureAggregationResponse.deserializeBinary(new Uint8Array(buffer_arg));
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
// signing, verification, and etc.
var UtilsService = exports.UtilsService = {
  // SignMessageWithPrivateKey signs a message with the provided private key.
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
  // VerifyMessage verifies a signature against the public key and message.
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
  // PublicKeyAggregation aggregates multiple BLS public keys into a single key.
publicKeyAggregation: {
    path: '/pactus.Utils/PublicKeyAggregation',
    requestStream: false,
    responseStream: false,
    requestType: utils_pb.PublicKeyAggregationRequest,
    responseType: utils_pb.PublicKeyAggregationResponse,
    requestSerialize: serialize_pactus_PublicKeyAggregationRequest,
    requestDeserialize: deserialize_pactus_PublicKeyAggregationRequest,
    responseSerialize: serialize_pactus_PublicKeyAggregationResponse,
    responseDeserialize: deserialize_pactus_PublicKeyAggregationResponse,
  },
  // SignatureAggregation aggregates multiple BLS signatures into a single signature.
signatureAggregation: {
    path: '/pactus.Utils/SignatureAggregation',
    requestStream: false,
    responseStream: false,
    requestType: utils_pb.SignatureAggregationRequest,
    responseType: utils_pb.SignatureAggregationResponse,
    requestSerialize: serialize_pactus_SignatureAggregationRequest,
    requestDeserialize: deserialize_pactus_SignatureAggregationRequest,
    responseSerialize: serialize_pactus_SignatureAggregationResponse,
    responseDeserialize: deserialize_pactus_SignatureAggregationResponse,
  },
};

exports.UtilsClient = grpc.makeGenericClientConstructor(UtilsService, 'Utils');
