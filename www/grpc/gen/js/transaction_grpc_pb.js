// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('grpc');
var transaction_pb = require('./transaction_pb.js');

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

function serialize_pactus_TransactionRequest(arg) {
  if (!(arg instanceof transaction_pb.TransactionRequest)) {
    throw new Error('Expected argument of type pactus.TransactionRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_TransactionRequest(buffer_arg) {
  return transaction_pb.TransactionRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_TransactionResponse(arg) {
  if (!(arg instanceof transaction_pb.TransactionResponse)) {
    throw new Error('Expected argument of type pactus.TransactionResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_TransactionResponse(buffer_arg) {
  return transaction_pb.TransactionResponse.deserializeBinary(new Uint8Array(buffer_arg));
}


var TransactionService = exports.TransactionService = {
  getTransaction: {
    path: '/pactus.Transaction/GetTransaction',
    requestStream: false,
    responseStream: false,
    requestType: transaction_pb.TransactionRequest,
    responseType: transaction_pb.TransactionResponse,
    requestSerialize: serialize_pactus_TransactionRequest,
    requestDeserialize: deserialize_pactus_TransactionRequest,
    responseSerialize: serialize_pactus_TransactionResponse,
    responseDeserialize: deserialize_pactus_TransactionResponse,
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
