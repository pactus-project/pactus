// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('grpc');
var wallet_pb = require('./wallet_pb.js');

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

function serialize_pactus_LoadWalletRequest(arg) {
  if (!(arg instanceof wallet_pb.LoadWalletRequest)) {
    throw new Error('Expected argument of type pactus.LoadWalletRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_LoadWalletRequest(buffer_arg) {
  return wallet_pb.LoadWalletRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_LoadWalletResponse(arg) {
  if (!(arg instanceof wallet_pb.LoadWalletResponse)) {
    throw new Error('Expected argument of type pactus.LoadWalletResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_LoadWalletResponse(buffer_arg) {
  return wallet_pb.LoadWalletResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_LockWalletRequest(arg) {
  if (!(arg instanceof wallet_pb.LockWalletRequest)) {
    throw new Error('Expected argument of type pactus.LockWalletRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_LockWalletRequest(buffer_arg) {
  return wallet_pb.LockWalletRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_LockWalletResponse(arg) {
  if (!(arg instanceof wallet_pb.LockWalletResponse)) {
    throw new Error('Expected argument of type pactus.LockWalletResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_LockWalletResponse(buffer_arg) {
  return wallet_pb.LockWalletResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_UnloadWalletRequest(arg) {
  if (!(arg instanceof wallet_pb.UnloadWalletRequest)) {
    throw new Error('Expected argument of type pactus.UnloadWalletRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_UnloadWalletRequest(buffer_arg) {
  return wallet_pb.UnloadWalletRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_UnloadWalletResponse(arg) {
  if (!(arg instanceof wallet_pb.UnloadWalletResponse)) {
    throw new Error('Expected argument of type pactus.UnloadWalletResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_UnloadWalletResponse(buffer_arg) {
  return wallet_pb.UnloadWalletResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_UnlockWalletRequest(arg) {
  if (!(arg instanceof wallet_pb.UnlockWalletRequest)) {
    throw new Error('Expected argument of type pactus.UnlockWalletRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_UnlockWalletRequest(buffer_arg) {
  return wallet_pb.UnlockWalletRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_UnlockWalletResponse(arg) {
  if (!(arg instanceof wallet_pb.UnlockWalletResponse)) {
    throw new Error('Expected argument of type pactus.UnlockWalletResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_UnlockWalletResponse(buffer_arg) {
  return wallet_pb.UnlockWalletResponse.deserializeBinary(new Uint8Array(buffer_arg));
}


// Wallet service defines RPC methods for managing wallet operations.
var WalletService = exports.WalletService = {
  // CreateWallet creates a new wallet with the specified parameters.
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
  // LoadWallet loads an existing wallet with the given name.
loadWallet: {
    path: '/pactus.Wallet/LoadWallet',
    requestStream: false,
    responseStream: false,
    requestType: wallet_pb.LoadWalletRequest,
    responseType: wallet_pb.LoadWalletResponse,
    requestSerialize: serialize_pactus_LoadWalletRequest,
    requestDeserialize: deserialize_pactus_LoadWalletRequest,
    responseSerialize: serialize_pactus_LoadWalletResponse,
    responseDeserialize: deserialize_pactus_LoadWalletResponse,
  },
  // UnloadWallet unloads a currently loaded wallet with the specified name.
unloadWallet: {
    path: '/pactus.Wallet/UnloadWallet',
    requestStream: false,
    responseStream: false,
    requestType: wallet_pb.UnloadWalletRequest,
    responseType: wallet_pb.UnloadWalletResponse,
    requestSerialize: serialize_pactus_UnloadWalletRequest,
    requestDeserialize: deserialize_pactus_UnloadWalletRequest,
    responseSerialize: serialize_pactus_UnloadWalletResponse,
    responseDeserialize: deserialize_pactus_UnloadWalletResponse,
  },
  // LockWallet locks a currently loaded wallet with the provided password and timeout.
lockWallet: {
    path: '/pactus.Wallet/LockWallet',
    requestStream: false,
    responseStream: false,
    requestType: wallet_pb.LockWalletRequest,
    responseType: wallet_pb.LockWalletResponse,
    requestSerialize: serialize_pactus_LockWalletRequest,
    requestDeserialize: deserialize_pactus_LockWalletRequest,
    responseSerialize: serialize_pactus_LockWalletResponse,
    responseDeserialize: deserialize_pactus_LockWalletResponse,
  },
  // UnlockWallet unlocks a locked wallet with the provided password and timeout.
unlockWallet: {
    path: '/pactus.Wallet/UnlockWallet',
    requestStream: false,
    responseStream: false,
    requestType: wallet_pb.UnlockWalletRequest,
    responseType: wallet_pb.UnlockWalletResponse,
    requestSerialize: serialize_pactus_UnlockWalletRequest,
    requestDeserialize: deserialize_pactus_UnlockWalletRequest,
    responseSerialize: serialize_pactus_UnlockWalletResponse,
    responseDeserialize: deserialize_pactus_UnlockWalletResponse,
  },
};

exports.WalletClient = grpc.makeGenericClientConstructor(WalletService);
