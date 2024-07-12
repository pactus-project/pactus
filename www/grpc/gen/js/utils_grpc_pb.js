// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('grpc');
var utils_pb = require('./utils_pb.js');

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
  // SignMessageWithPrivateKey sign message with provided private key.
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
  // VerifyMessage verify signature with public key and message
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
};

exports.UtilsClient = grpc.makeGenericClientConstructor(UtilsService);
