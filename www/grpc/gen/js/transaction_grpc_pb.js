// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('grpc');
var transaction_pb = require('./transaction_pb.js');

function serialize_pactus_BroadcastTransactionRequest(arg) {
  if (!(arg instanceof transaction_pb.BroadcastTransactionRequest)) {
    throw new Error('Expected argument of type pactus.BroadcastTransactionRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_BroadcastTransactionRequest(buffer_arg) {
  return transaction_pb.BroadcastTransactionRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_BroadcastTransactionResponse(arg) {
  if (!(arg instanceof transaction_pb.BroadcastTransactionResponse)) {
    throw new Error('Expected argument of type pactus.BroadcastTransactionResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_BroadcastTransactionResponse(buffer_arg) {
  return transaction_pb.BroadcastTransactionResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

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

function serialize_pactus_GetRawUnbondTransactionRequest(arg) {
  if (!(arg instanceof transaction_pb.GetRawUnbondTransactionRequest)) {
    throw new Error('Expected argument of type pactus.GetRawUnbondTransactionRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_GetRawUnbondTransactionRequest(buffer_arg) {
  return transaction_pb.GetRawUnbondTransactionRequest.deserializeBinary(new Uint8Array(buffer_arg));
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


// Transaction service defines various RPC methods for interacting with
// transactions.
var TransactionService = exports.TransactionService = {
  // GetTransaction retrieves transaction details based on the provided request
// parameters.
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
  // CalculateFee calculates the transaction fee based on the specified amount
// and payload type.
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
  // BroadcastTransaction broadcasts a signed transaction to the network.
broadcastTransaction: {
    path: '/pactus.Transaction/BroadcastTransaction',
    requestStream: false,
    responseStream: false,
    requestType: transaction_pb.BroadcastTransactionRequest,
    responseType: transaction_pb.BroadcastTransactionResponse,
    requestSerialize: serialize_pactus_BroadcastTransactionRequest,
    requestDeserialize: deserialize_pactus_BroadcastTransactionRequest,
    responseSerialize: serialize_pactus_BroadcastTransactionResponse,
    responseDeserialize: deserialize_pactus_BroadcastTransactionResponse,
  },
  // GetRawTransferTransaction retrieves raw details of a transfer transaction.
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
  // GetRawBondTransaction retrieves raw details of a bond transaction.
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
  // GetRawUnbondTransaction retrieves raw details of an unbond transaction.
getRawUnbondTransaction: {
    path: '/pactus.Transaction/GetRawUnbondTransaction',
    requestStream: false,
    responseStream: false,
    requestType: transaction_pb.GetRawUnbondTransactionRequest,
    responseType: transaction_pb.GetRawTransactionResponse,
    requestSerialize: serialize_pactus_GetRawUnbondTransactionRequest,
    requestDeserialize: deserialize_pactus_GetRawUnbondTransactionRequest,
    responseSerialize: serialize_pactus_GetRawTransactionResponse,
    responseDeserialize: deserialize_pactus_GetRawTransactionResponse,
  },
  // GetRawWithdrawTransaction retrieves raw details of a withdraw transaction.
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
