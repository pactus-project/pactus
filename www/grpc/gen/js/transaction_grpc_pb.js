// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('grpc');
var transaction_pb = require('./transaction_pb.js');

function serialize_pactus_CalculateFeeRequest(arg) {
  if (!(arg instanceof transaction_pb.CalculateFeeRequest)) {
    throw new Error('Expected argument of type pactus.CalculateFeeRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_CalculateFeeRequest(buffer_arg) {
  return transaction_pb.CalculateFeeRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_CalculateFeeResponse(arg) {
  if (!(arg instanceof transaction_pb.CalculateFeeResponse)) {
    throw new Error('Expected argument of type pactus.CalculateFeeResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_CalculateFeeResponse(buffer_arg) {
  return transaction_pb.CalculateFeeResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_GetRawBondTransactionRequest(arg) {
  if (!(arg instanceof transaction_pb.GetRawBondTransactionRequest)) {
    throw new Error('Expected argument of type pactus.GetRawBondTransactionRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_GetRawBondTransactionRequest(buffer_arg) {
  return transaction_pb.GetRawBondTransactionRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_GetRawTransactionResponse(arg) {
  if (!(arg instanceof transaction_pb.GetRawTransactionResponse)) {
    throw new Error('Expected argument of type pactus.GetRawTransactionResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_GetRawTransactionResponse(buffer_arg) {
  return transaction_pb.GetRawTransactionResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_GetRawTransferTransactionRequest(arg) {
  if (!(arg instanceof transaction_pb.GetRawTransferTransactionRequest)) {
    throw new Error('Expected argument of type pactus.GetRawTransferTransactionRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_GetRawTransferTransactionRequest(buffer_arg) {
  return transaction_pb.GetRawTransferTransactionRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_GetRawUnBondTransactionRequest(arg) {
  if (!(arg instanceof transaction_pb.GetRawUnBondTransactionRequest)) {
    throw new Error('Expected argument of type pactus.GetRawUnBondTransactionRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_GetRawUnBondTransactionRequest(buffer_arg) {
  return transaction_pb.GetRawUnBondTransactionRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_GetRawWithdrawTransactionRequest(arg) {
  if (!(arg instanceof transaction_pb.GetRawWithdrawTransactionRequest)) {
    throw new Error('Expected argument of type pactus.GetRawWithdrawTransactionRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_GetRawWithdrawTransactionRequest(buffer_arg) {
  return transaction_pb.GetRawWithdrawTransactionRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

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
  calculateFee: {
    path: '/pactus.Transaction/CalculateFee',
    requestStream: false,
    responseStream: false,
    requestType: transaction_pb.CalculateFeeRequest,
    responseType: transaction_pb.CalculateFeeResponse,
    requestSerialize: serialize_pactus_CalculateFeeRequest,
    requestDeserialize: deserialize_pactus_CalculateFeeRequest,
    responseSerialize: serialize_pactus_CalculateFeeResponse,
    responseDeserialize: deserialize_pactus_CalculateFeeResponse,
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
  getRawTransferTransaction: {
    path: '/pactus.Transaction/GetRawTransferTransaction',
    requestStream: false,
    responseStream: false,
    requestType: transaction_pb.GetRawTransferTransactionRequest,
    responseType: transaction_pb.GetRawTransactionResponse,
    requestSerialize: serialize_pactus_GetRawTransferTransactionRequest,
    requestDeserialize: deserialize_pactus_GetRawTransferTransactionRequest,
    responseSerialize: serialize_pactus_GetRawTransactionResponse,
    responseDeserialize: deserialize_pactus_GetRawTransactionResponse,
  },
  getRawBondTransaction: {
    path: '/pactus.Transaction/GetRawBondTransaction',
    requestStream: false,
    responseStream: false,
    requestType: transaction_pb.GetRawBondTransactionRequest,
    responseType: transaction_pb.GetRawTransactionResponse,
    requestSerialize: serialize_pactus_GetRawBondTransactionRequest,
    requestDeserialize: deserialize_pactus_GetRawBondTransactionRequest,
    responseSerialize: serialize_pactus_GetRawTransactionResponse,
    responseDeserialize: deserialize_pactus_GetRawTransactionResponse,
  },
  getRawUnBondTransaction: {
    path: '/pactus.Transaction/GetRawUnBondTransaction',
    requestStream: false,
    responseStream: false,
    requestType: transaction_pb.GetRawUnBondTransactionRequest,
    responseType: transaction_pb.GetRawTransactionResponse,
    requestSerialize: serialize_pactus_GetRawUnBondTransactionRequest,
    requestDeserialize: deserialize_pactus_GetRawUnBondTransactionRequest,
    responseSerialize: serialize_pactus_GetRawTransactionResponse,
    responseDeserialize: deserialize_pactus_GetRawTransactionResponse,
  },
  getRawWithdrawTransaction: {
    path: '/pactus.Transaction/GetRawWithdrawTransaction',
    requestStream: false,
    responseStream: false,
    requestType: transaction_pb.GetRawWithdrawTransactionRequest,
    responseType: transaction_pb.GetRawTransactionResponse,
    requestSerialize: serialize_pactus_GetRawWithdrawTransactionRequest,
    requestDeserialize: deserialize_pactus_GetRawWithdrawTransactionRequest,
    responseSerialize: serialize_pactus_GetRawTransactionResponse,
    responseDeserialize: deserialize_pactus_GetRawTransactionResponse,
  },
};

exports.TransactionClient = grpc.makeGenericClientConstructor(TransactionService);
