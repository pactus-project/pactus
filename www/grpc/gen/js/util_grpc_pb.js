// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('grpc');
var util_pb = require('./util_pb.js');

function serialize_pactus_SignMessageWithPrivateKeyRequest(arg) {
  if (!(arg instanceof util_pb.SignMessageWithPrivateKeyRequest)) {
    throw new Error('Expected argument of type pactus.SignMessageWithPrivateKeyRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_SignMessageWithPrivateKeyRequest(buffer_arg) {
  return util_pb.SignMessageWithPrivateKeyRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_SignMessageWithPrivateKeyResponse(arg) {
  if (!(arg instanceof util_pb.SignMessageWithPrivateKeyResponse)) {
    throw new Error('Expected argument of type pactus.SignMessageWithPrivateKeyResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_SignMessageWithPrivateKeyResponse(buffer_arg) {
  return util_pb.SignMessageWithPrivateKeyResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_VerifyMessageRequest(arg) {
  if (!(arg instanceof util_pb.VerifyMessageRequest)) {
    throw new Error('Expected argument of type pactus.VerifyMessageRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_VerifyMessageRequest(buffer_arg) {
  return util_pb.VerifyMessageRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_VerifyMessageResponse(arg) {
  if (!(arg instanceof util_pb.VerifyMessageResponse)) {
    throw new Error('Expected argument of type pactus.VerifyMessageResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_VerifyMessageResponse(buffer_arg) {
  return util_pb.VerifyMessageResponse.deserializeBinary(new Uint8Array(buffer_arg));
}


// Util service defines various RPC methods for interacting with
// Utils.
var UtilService = exports.UtilService = {
  // SignMessageWithPrivateKey
signMessageWithPrivateKey: {
    path: '/pactus.Util/SignMessageWithPrivateKey',
    requestStream: false,
    responseStream: false,
    requestType: util_pb.SignMessageWithPrivateKeyRequest,
    responseType: util_pb.SignMessageWithPrivateKeyResponse,
    requestSerialize: serialize_pactus_SignMessageWithPrivateKeyRequest,
    requestDeserialize: deserialize_pactus_SignMessageWithPrivateKeyRequest,
    responseSerialize: serialize_pactus_SignMessageWithPrivateKeyResponse,
    responseDeserialize: deserialize_pactus_SignMessageWithPrivateKeyResponse,
  },
  // VerifyMessage
verifyMessage: {
    path: '/pactus.Util/VerifyMessage',
    requestStream: false,
    responseStream: false,
    requestType: util_pb.VerifyMessageRequest,
    responseType: util_pb.VerifyMessageResponse,
    requestSerialize: serialize_pactus_VerifyMessageRequest,
    requestDeserialize: deserialize_pactus_VerifyMessageRequest,
    responseSerialize: serialize_pactus_VerifyMessageResponse,
    responseDeserialize: deserialize_pactus_VerifyMessageResponse,
  },
};

exports.UtilClient = grpc.makeGenericClientConstructor(UtilService);
