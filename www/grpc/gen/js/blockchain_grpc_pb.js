// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('grpc');
var blockchain_pb = require('./blockchain_pb.js');
var transaction_pb = require('./transaction_pb.js');

function serialize_pactus_AccountRequest(arg) {
  if (!(arg instanceof blockchain_pb.AccountRequest)) {
    throw new Error('Expected argument of type pactus.AccountRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_AccountRequest(buffer_arg) {
  return blockchain_pb.AccountRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_AccountResponse(arg) {
  if (!(arg instanceof blockchain_pb.AccountResponse)) {
    throw new Error('Expected argument of type pactus.AccountResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_AccountResponse(buffer_arg) {
  return blockchain_pb.AccountResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_BlockHashRequest(arg) {
  if (!(arg instanceof blockchain_pb.BlockHashRequest)) {
    throw new Error('Expected argument of type pactus.BlockHashRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_BlockHashRequest(buffer_arg) {
  return blockchain_pb.BlockHashRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_BlockHashResponse(arg) {
  if (!(arg instanceof blockchain_pb.BlockHashResponse)) {
    throw new Error('Expected argument of type pactus.BlockHashResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_BlockHashResponse(buffer_arg) {
  return blockchain_pb.BlockHashResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_BlockHeightRequest(arg) {
  if (!(arg instanceof blockchain_pb.BlockHeightRequest)) {
    throw new Error('Expected argument of type pactus.BlockHeightRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_BlockHeightRequest(buffer_arg) {
  return blockchain_pb.BlockHeightRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_BlockHeightResponse(arg) {
  if (!(arg instanceof blockchain_pb.BlockHeightResponse)) {
    throw new Error('Expected argument of type pactus.BlockHeightResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_BlockHeightResponse(buffer_arg) {
  return blockchain_pb.BlockHeightResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_BlockRequest(arg) {
  if (!(arg instanceof blockchain_pb.BlockRequest)) {
    throw new Error('Expected argument of type pactus.BlockRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_BlockRequest(buffer_arg) {
  return blockchain_pb.BlockRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_BlockResponse(arg) {
  if (!(arg instanceof blockchain_pb.BlockResponse)) {
    throw new Error('Expected argument of type pactus.BlockResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_BlockResponse(buffer_arg) {
  return blockchain_pb.BlockResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_BlockchainInfoRequest(arg) {
  if (!(arg instanceof blockchain_pb.BlockchainInfoRequest)) {
    throw new Error('Expected argument of type pactus.BlockchainInfoRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_BlockchainInfoRequest(buffer_arg) {
  return blockchain_pb.BlockchainInfoRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_BlockchainInfoResponse(arg) {
  if (!(arg instanceof blockchain_pb.BlockchainInfoResponse)) {
    throw new Error('Expected argument of type pactus.BlockchainInfoResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_BlockchainInfoResponse(buffer_arg) {
  return blockchain_pb.BlockchainInfoResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_ValidatorByNumberRequest(arg) {
  if (!(arg instanceof blockchain_pb.ValidatorByNumberRequest)) {
    throw new Error('Expected argument of type pactus.ValidatorByNumberRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_ValidatorByNumberRequest(buffer_arg) {
  return blockchain_pb.ValidatorByNumberRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_ValidatorRequest(arg) {
  if (!(arg instanceof blockchain_pb.ValidatorRequest)) {
    throw new Error('Expected argument of type pactus.ValidatorRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_ValidatorRequest(buffer_arg) {
  return blockchain_pb.ValidatorRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_ValidatorResponse(arg) {
  if (!(arg instanceof blockchain_pb.ValidatorResponse)) {
    throw new Error('Expected argument of type pactus.ValidatorResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_ValidatorResponse(buffer_arg) {
  return blockchain_pb.ValidatorResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_ValidatorsRequest(arg) {
  if (!(arg instanceof blockchain_pb.ValidatorsRequest)) {
    throw new Error('Expected argument of type pactus.ValidatorsRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_ValidatorsRequest(buffer_arg) {
  return blockchain_pb.ValidatorsRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_ValidatorsResponse(arg) {
  if (!(arg instanceof blockchain_pb.ValidatorsResponse)) {
    throw new Error('Expected argument of type pactus.ValidatorsResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_ValidatorsResponse(buffer_arg) {
  return blockchain_pb.ValidatorsResponse.deserializeBinary(new Uint8Array(buffer_arg));
}


var BlockchainService = exports.BlockchainService = {
  getBlock: {
    path: '/pactus.Blockchain/GetBlock',
    requestStream: false,
    responseStream: false,
    requestType: blockchain_pb.BlockRequest,
    responseType: blockchain_pb.BlockResponse,
    requestSerialize: serialize_pactus_BlockRequest,
    requestDeserialize: deserialize_pactus_BlockRequest,
    responseSerialize: serialize_pactus_BlockResponse,
    responseDeserialize: deserialize_pactus_BlockResponse,
  },
  getBlockHash: {
    path: '/pactus.Blockchain/GetBlockHash',
    requestStream: false,
    responseStream: false,
    requestType: blockchain_pb.BlockHashRequest,
    responseType: blockchain_pb.BlockHashResponse,
    requestSerialize: serialize_pactus_BlockHashRequest,
    requestDeserialize: deserialize_pactus_BlockHashRequest,
    responseSerialize: serialize_pactus_BlockHashResponse,
    responseDeserialize: deserialize_pactus_BlockHashResponse,
  },
  getBlockHeight: {
    path: '/pactus.Blockchain/GetBlockHeight',
    requestStream: false,
    responseStream: false,
    requestType: blockchain_pb.BlockHeightRequest,
    responseType: blockchain_pb.BlockHeightResponse,
    requestSerialize: serialize_pactus_BlockHeightRequest,
    requestDeserialize: deserialize_pactus_BlockHeightRequest,
    responseSerialize: serialize_pactus_BlockHeightResponse,
    responseDeserialize: deserialize_pactus_BlockHeightResponse,
  },
  getAccount: {
    path: '/pactus.Blockchain/GetAccount',
    requestStream: false,
    responseStream: false,
    requestType: blockchain_pb.AccountRequest,
    responseType: blockchain_pb.AccountResponse,
    requestSerialize: serialize_pactus_AccountRequest,
    requestDeserialize: deserialize_pactus_AccountRequest,
    responseSerialize: serialize_pactus_AccountResponse,
    responseDeserialize: deserialize_pactus_AccountResponse,
  },
  getValidators: {
    path: '/pactus.Blockchain/GetValidators',
    requestStream: false,
    responseStream: false,
    requestType: blockchain_pb.ValidatorsRequest,
    responseType: blockchain_pb.ValidatorsResponse,
    requestSerialize: serialize_pactus_ValidatorsRequest,
    requestDeserialize: deserialize_pactus_ValidatorsRequest,
    responseSerialize: serialize_pactus_ValidatorsResponse,
    responseDeserialize: deserialize_pactus_ValidatorsResponse,
  },
  getValidator: {
    path: '/pactus.Blockchain/GetValidator',
    requestStream: false,
    responseStream: false,
    requestType: blockchain_pb.ValidatorRequest,
    responseType: blockchain_pb.ValidatorResponse,
    requestSerialize: serialize_pactus_ValidatorRequest,
    requestDeserialize: deserialize_pactus_ValidatorRequest,
    responseSerialize: serialize_pactus_ValidatorResponse,
    responseDeserialize: deserialize_pactus_ValidatorResponse,
  },
  getValidatorByNumber: {
    path: '/pactus.Blockchain/GetValidatorByNumber',
    requestStream: false,
    responseStream: false,
    requestType: blockchain_pb.ValidatorByNumberRequest,
    responseType: blockchain_pb.ValidatorResponse,
    requestSerialize: serialize_pactus_ValidatorByNumberRequest,
    requestDeserialize: deserialize_pactus_ValidatorByNumberRequest,
    responseSerialize: serialize_pactus_ValidatorResponse,
    responseDeserialize: deserialize_pactus_ValidatorResponse,
  },
  getBlockchainInfo: {
    path: '/pactus.Blockchain/GetBlockchainInfo',
    requestStream: false,
    responseStream: false,
    requestType: blockchain_pb.BlockchainInfoRequest,
    responseType: blockchain_pb.BlockchainInfoResponse,
    requestSerialize: serialize_pactus_BlockchainInfoRequest,
    requestDeserialize: deserialize_pactus_BlockchainInfoRequest,
    responseSerialize: serialize_pactus_BlockchainInfoResponse,
    responseDeserialize: deserialize_pactus_BlockchainInfoResponse,
  },
};

exports.BlockchainClient = grpc.makeGenericClientConstructor(BlockchainService);
