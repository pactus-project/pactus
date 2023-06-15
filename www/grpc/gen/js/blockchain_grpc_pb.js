// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('grpc');
var blockchain_pb = require('./blockchain_pb.js');
var transaction_pb = require('./transaction_pb.js');

function serialize_pactus_GetAccountByNumberRequest(arg) {
  if (!(arg instanceof blockchain_pb.GetAccountByNumberRequest)) {
    throw new Error('Expected argument of type pactus.GetAccountByNumberRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_GetAccountByNumberRequest(buffer_arg) {
  return blockchain_pb.GetAccountByNumberRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_GetAccountRequest(arg) {
  if (!(arg instanceof blockchain_pb.GetAccountRequest)) {
    throw new Error('Expected argument of type pactus.GetAccountRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_GetAccountRequest(buffer_arg) {
  return blockchain_pb.GetAccountRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_GetAccountResponse(arg) {
  if (!(arg instanceof blockchain_pb.GetAccountResponse)) {
    throw new Error('Expected argument of type pactus.GetAccountResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_GetAccountResponse(buffer_arg) {
  return blockchain_pb.GetAccountResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_GetBlockHashRequest(arg) {
  if (!(arg instanceof blockchain_pb.GetBlockHashRequest)) {
    throw new Error('Expected argument of type pactus.GetBlockHashRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_GetBlockHashRequest(buffer_arg) {
  return blockchain_pb.GetBlockHashRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_GetBlockHashResponse(arg) {
  if (!(arg instanceof blockchain_pb.GetBlockHashResponse)) {
    throw new Error('Expected argument of type pactus.GetBlockHashResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_GetBlockHashResponse(buffer_arg) {
  return blockchain_pb.GetBlockHashResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_GetBlockHeightRequest(arg) {
  if (!(arg instanceof blockchain_pb.GetBlockHeightRequest)) {
    throw new Error('Expected argument of type pactus.GetBlockHeightRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_GetBlockHeightRequest(buffer_arg) {
  return blockchain_pb.GetBlockHeightRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_GetBlockHeightResponse(arg) {
  if (!(arg instanceof blockchain_pb.GetBlockHeightResponse)) {
    throw new Error('Expected argument of type pactus.GetBlockHeightResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_GetBlockHeightResponse(buffer_arg) {
  return blockchain_pb.GetBlockHeightResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_GetBlockRequest(arg) {
  if (!(arg instanceof blockchain_pb.GetBlockRequest)) {
    throw new Error('Expected argument of type pactus.GetBlockRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_GetBlockRequest(buffer_arg) {
  return blockchain_pb.GetBlockRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_GetBlockResponse(arg) {
  if (!(arg instanceof blockchain_pb.GetBlockResponse)) {
    throw new Error('Expected argument of type pactus.GetBlockResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_GetBlockResponse(buffer_arg) {
  return blockchain_pb.GetBlockResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_GetBlockchainInfoRequest(arg) {
  if (!(arg instanceof blockchain_pb.GetBlockchainInfoRequest)) {
    throw new Error('Expected argument of type pactus.GetBlockchainInfoRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_GetBlockchainInfoRequest(buffer_arg) {
  return blockchain_pb.GetBlockchainInfoRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_GetBlockchainInfoResponse(arg) {
  if (!(arg instanceof blockchain_pb.GetBlockchainInfoResponse)) {
    throw new Error('Expected argument of type pactus.GetBlockchainInfoResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_GetBlockchainInfoResponse(buffer_arg) {
  return blockchain_pb.GetBlockchainInfoResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_GetConsensusInfoRequest(arg) {
  if (!(arg instanceof blockchain_pb.GetConsensusInfoRequest)) {
    throw new Error('Expected argument of type pactus.GetConsensusInfoRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_GetConsensusInfoRequest(buffer_arg) {
  return blockchain_pb.GetConsensusInfoRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_GetConsensusInfoResponse(arg) {
  if (!(arg instanceof blockchain_pb.GetConsensusInfoResponse)) {
    throw new Error('Expected argument of type pactus.GetConsensusInfoResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_GetConsensusInfoResponse(buffer_arg) {
  return blockchain_pb.GetConsensusInfoResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_GetValidatorByNumberRequest(arg) {
  if (!(arg instanceof blockchain_pb.GetValidatorByNumberRequest)) {
    throw new Error('Expected argument of type pactus.GetValidatorByNumberRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_GetValidatorByNumberRequest(buffer_arg) {
  return blockchain_pb.GetValidatorByNumberRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_GetValidatorRequest(arg) {
  if (!(arg instanceof blockchain_pb.GetValidatorRequest)) {
    throw new Error('Expected argument of type pactus.GetValidatorRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_GetValidatorRequest(buffer_arg) {
  return blockchain_pb.GetValidatorRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_GetValidatorResponse(arg) {
  if (!(arg instanceof blockchain_pb.GetValidatorResponse)) {
    throw new Error('Expected argument of type pactus.GetValidatorResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_GetValidatorResponse(buffer_arg) {
  return blockchain_pb.GetValidatorResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_GetValidatorsRequest(arg) {
  if (!(arg instanceof blockchain_pb.GetValidatorsRequest)) {
    throw new Error('Expected argument of type pactus.GetValidatorsRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_GetValidatorsRequest(buffer_arg) {
  return blockchain_pb.GetValidatorsRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_GetValidatorsResponse(arg) {
  if (!(arg instanceof blockchain_pb.GetValidatorsResponse)) {
    throw new Error('Expected argument of type pactus.GetValidatorsResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_GetValidatorsResponse(buffer_arg) {
  return blockchain_pb.GetValidatorsResponse.deserializeBinary(new Uint8Array(buffer_arg));
}


var BlockchainService = exports.BlockchainService = {
  getBlock: {
    path: '/pactus.Blockchain/GetBlock',
    requestStream: false,
    responseStream: false,
    requestType: blockchain_pb.GetBlockRequest,
    responseType: blockchain_pb.GetBlockResponse,
    requestSerialize: serialize_pactus_GetBlockRequest,
    requestDeserialize: deserialize_pactus_GetBlockRequest,
    responseSerialize: serialize_pactus_GetBlockResponse,
    responseDeserialize: deserialize_pactus_GetBlockResponse,
  },
  getBlockHash: {
    path: '/pactus.Blockchain/GetBlockHash',
    requestStream: false,
    responseStream: false,
    requestType: blockchain_pb.GetBlockHashRequest,
    responseType: blockchain_pb.GetBlockHashResponse,
    requestSerialize: serialize_pactus_GetBlockHashRequest,
    requestDeserialize: deserialize_pactus_GetBlockHashRequest,
    responseSerialize: serialize_pactus_GetBlockHashResponse,
    responseDeserialize: deserialize_pactus_GetBlockHashResponse,
  },
  getBlockHeight: {
    path: '/pactus.Blockchain/GetBlockHeight',
    requestStream: false,
    responseStream: false,
    requestType: blockchain_pb.GetBlockHeightRequest,
    responseType: blockchain_pb.GetBlockHeightResponse,
    requestSerialize: serialize_pactus_GetBlockHeightRequest,
    requestDeserialize: deserialize_pactus_GetBlockHeightRequest,
    responseSerialize: serialize_pactus_GetBlockHeightResponse,
    responseDeserialize: deserialize_pactus_GetBlockHeightResponse,
  },
  getBlockchainInfo: {
    path: '/pactus.Blockchain/GetBlockchainInfo',
    requestStream: false,
    responseStream: false,
    requestType: blockchain_pb.GetBlockchainInfoRequest,
    responseType: blockchain_pb.GetBlockchainInfoResponse,
    requestSerialize: serialize_pactus_GetBlockchainInfoRequest,
    requestDeserialize: deserialize_pactus_GetBlockchainInfoRequest,
    responseSerialize: serialize_pactus_GetBlockchainInfoResponse,
    responseDeserialize: deserialize_pactus_GetBlockchainInfoResponse,
  },
  getConsensusInfo: {
    path: '/pactus.Blockchain/GetConsensusInfo',
    requestStream: false,
    responseStream: false,
    requestType: blockchain_pb.GetConsensusInfoRequest,
    responseType: blockchain_pb.GetConsensusInfoResponse,
    requestSerialize: serialize_pactus_GetConsensusInfoRequest,
    requestDeserialize: deserialize_pactus_GetConsensusInfoRequest,
    responseSerialize: serialize_pactus_GetConsensusInfoResponse,
    responseDeserialize: deserialize_pactus_GetConsensusInfoResponse,
  },
  getAccount: {
    path: '/pactus.Blockchain/GetAccount',
    requestStream: false,
    responseStream: false,
    requestType: blockchain_pb.GetAccountRequest,
    responseType: blockchain_pb.GetAccountResponse,
    requestSerialize: serialize_pactus_GetAccountRequest,
    requestDeserialize: deserialize_pactus_GetAccountRequest,
    responseSerialize: serialize_pactus_GetAccountResponse,
    responseDeserialize: deserialize_pactus_GetAccountResponse,
  },
  getAccountByNumber: {
    path: '/pactus.Blockchain/GetAccountByNumber',
    requestStream: false,
    responseStream: false,
    requestType: blockchain_pb.GetAccountByNumberRequest,
    responseType: blockchain_pb.GetAccountResponse,
    requestSerialize: serialize_pactus_GetAccountByNumberRequest,
    requestDeserialize: deserialize_pactus_GetAccountByNumberRequest,
    responseSerialize: serialize_pactus_GetAccountResponse,
    responseDeserialize: deserialize_pactus_GetAccountResponse,
  },
  getValidator: {
    path: '/pactus.Blockchain/GetValidator',
    requestStream: false,
    responseStream: false,
    requestType: blockchain_pb.GetValidatorRequest,
    responseType: blockchain_pb.GetValidatorResponse,
    requestSerialize: serialize_pactus_GetValidatorRequest,
    requestDeserialize: deserialize_pactus_GetValidatorRequest,
    responseSerialize: serialize_pactus_GetValidatorResponse,
    responseDeserialize: deserialize_pactus_GetValidatorResponse,
  },
  getValidatorByNumber: {
    path: '/pactus.Blockchain/GetValidatorByNumber',
    requestStream: false,
    responseStream: false,
    requestType: blockchain_pb.GetValidatorByNumberRequest,
    responseType: blockchain_pb.GetValidatorResponse,
    requestSerialize: serialize_pactus_GetValidatorByNumberRequest,
    requestDeserialize: deserialize_pactus_GetValidatorByNumberRequest,
    responseSerialize: serialize_pactus_GetValidatorResponse,
    responseDeserialize: deserialize_pactus_GetValidatorResponse,
  },
  getValidators: {
    path: '/pactus.Blockchain/GetValidators',
    requestStream: false,
    responseStream: false,
    requestType: blockchain_pb.GetValidatorsRequest,
    responseType: blockchain_pb.GetValidatorsResponse,
    requestSerialize: serialize_pactus_GetValidatorsRequest,
    requestDeserialize: deserialize_pactus_GetValidatorsRequest,
    responseSerialize: serialize_pactus_GetValidatorsResponse,
    responseDeserialize: deserialize_pactus_GetValidatorsResponse,
  },
};

exports.BlockchainClient = grpc.makeGenericClientConstructor(BlockchainService);
