// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('grpc');
var transaction_pb = require('./transaction_pb.js');

function serialize_pactus_GetTransactionRequest(arg) {
  if (!(arg instanceof transaction_pb.GetTransactionRequest)) {
    throw new Error('Expected argument of type pactus.GetTransactionRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_GetTransactionRequest(buffer_arg) {
  return transaction_pb.GetTransactionRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_GetTransactionResponse(arg) {
  if (!(arg instanceof transaction_pb.GetTransactionResponse)) {
    throw new Error('Expected argument of type pactus.GetTransactionResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_GetTransactionResponse(buffer_arg) {
  return transaction_pb.GetTransactionResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_SendRawTransactionRequest(arg) {
  if (!(arg instanceof transaction_pb.SendRawTransactionRequest)) {
    throw new Error('Expected argument of type pactus.SendRawTransactionRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_SendRawTransactionRequest(buffer_arg) {
  return transaction_pb.SendRawTransactionRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_SendRawTransactionResponse(arg) {
  if (!(arg instanceof transaction_pb.SendRawTransactionResponse)) {
    throw new Error('Expected argument of type pactus.SendRawTransactionResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_SendRawTransactionResponse(buffer_arg) {
  return transaction_pb.SendRawTransactionResponse.deserializeBinary(new Uint8Array(buffer_arg));
}


var TransactionService = exports.TransactionService = {
  getTransaction: {
    path: '/pactus.Transaction/GetTransaction',
    requestStream: false,
    responseStream: false,
    requestType: transaction_pb.GetTransactionRequest,
    responseType: transaction_pb.GetTransactionResponse,
    requestSerialize: serialize_pactus_GetTransactionRequest,
    requestDeserialize: deserialize_pactus_GetTransactionRequest,
    responseSerialize: serialize_pactus_GetTransactionResponse,
    responseDeserialize: deserialize_pactus_GetTransactionResponse,
  },
  sendRawTransaction: {
    path: '/pactus.Transaction/SendRawTransaction',
    requestStream: false,
    responseStream: false,
    requestType: transaction_pb.SendRawTransactionRequest,
    responseType: transaction_pb.SendRawTransactionResponse,
    requestSerialize: serialize_pactus_SendRawTransactionRequest,
    requestDeserialize: deserialize_pactus_SendRawTransactionRequest,
    responseSerialize: serialize_pactus_SendRawTransactionResponse,
    responseDeserialize: deserialize_pactus_SendRawTransactionResponse,
  },
};

exports.TransactionClient = grpc.makeGenericClientConstructor(TransactionService);
