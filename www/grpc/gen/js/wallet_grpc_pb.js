// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('grpc');
var wallet_pb = require('./wallet_pb.js');
var google_api_annotations_pb = require('./google/api/annotations_pb.js');

function serialize_pactus_CreateWalletRequest(arg) {
  if (!(arg instanceof wallet_pb.CreateWalletRequest)) {
    throw new Error('Expected argument of type pactus.CreateWalletRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_CreateWalletRequest(buffer_arg) {
  return wallet_pb.CreateWalletRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_CreateWalletResponse(arg) {
  if (!(arg instanceof wallet_pb.CreateWalletResponse)) {
    throw new Error('Expected argument of type pactus.CreateWalletResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_CreateWalletResponse(buffer_arg) {
  return wallet_pb.CreateWalletResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_GenerateMnemonicRequest(arg) {
  if (!(arg instanceof wallet_pb.GenerateMnemonicRequest)) {
    throw new Error('Expected argument of type pactus.GenerateMnemonicRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_GenerateMnemonicRequest(buffer_arg) {
  return wallet_pb.GenerateMnemonicRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_GenerateMnemonicResponse(arg) {
  if (!(arg instanceof wallet_pb.GenerateMnemonicResponse)) {
    throw new Error('Expected argument of type pactus.GenerateMnemonicResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_GenerateMnemonicResponse(buffer_arg) {
  return wallet_pb.GenerateMnemonicResponse.deserializeBinary(new Uint8Array(buffer_arg));
}


var WalletService = exports.WalletService = {
  generateMnemonic: {
    path: '/pactus.Wallet/GenerateMnemonic',
    requestStream: false,
    responseStream: false,
    requestType: wallet_pb.GenerateMnemonicRequest,
    responseType: wallet_pb.GenerateMnemonicResponse,
    requestSerialize: serialize_pactus_GenerateMnemonicRequest,
    requestDeserialize: deserialize_pactus_GenerateMnemonicRequest,
    responseSerialize: serialize_pactus_GenerateMnemonicResponse,
    responseDeserialize: deserialize_pactus_GenerateMnemonicResponse,
  },
  createWallet: {
    path: '/pactus.Wallet/CreateWallet',
    requestStream: false,
    responseStream: false,
    requestType: wallet_pb.CreateWalletRequest,
    responseType: wallet_pb.CreateWalletResponse,
    requestSerialize: serialize_pactus_CreateWalletRequest,
    requestDeserialize: deserialize_pactus_CreateWalletRequest,
    responseSerialize: serialize_pactus_CreateWalletResponse,
    responseDeserialize: deserialize_pactus_CreateWalletResponse,
  },
};

exports.WalletClient = grpc.makeGenericClientConstructor(WalletService);
