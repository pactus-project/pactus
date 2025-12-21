// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('@grpc/grpc-js');
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

function serialize_pactus_GetAddressHistoryRequest(arg) {
  if (!(arg instanceof wallet_pb.GetAddressHistoryRequest)) {
    throw new Error('Expected argument of type pactus.GetAddressHistoryRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_GetAddressHistoryRequest(buffer_arg) {
  return wallet_pb.GetAddressHistoryRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_GetAddressHistoryResponse(arg) {
  if (!(arg instanceof wallet_pb.GetAddressHistoryResponse)) {
    throw new Error('Expected argument of type pactus.GetAddressHistoryResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_GetAddressHistoryResponse(buffer_arg) {
  return wallet_pb.GetAddressHistoryResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_GetAddressInfoRequest(arg) {
  if (!(arg instanceof wallet_pb.GetAddressInfoRequest)) {
    throw new Error('Expected argument of type pactus.GetAddressInfoRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_GetAddressInfoRequest(buffer_arg) {
  return wallet_pb.GetAddressInfoRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_GetAddressInfoResponse(arg) {
  if (!(arg instanceof wallet_pb.GetAddressInfoResponse)) {
    throw new Error('Expected argument of type pactus.GetAddressInfoResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_GetAddressInfoResponse(buffer_arg) {
  return wallet_pb.GetAddressInfoResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_GetNewAddressRequest(arg) {
  if (!(arg instanceof wallet_pb.GetNewAddressRequest)) {
    throw new Error('Expected argument of type pactus.GetNewAddressRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_GetNewAddressRequest(buffer_arg) {
  return wallet_pb.GetNewAddressRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_GetNewAddressResponse(arg) {
  if (!(arg instanceof wallet_pb.GetNewAddressResponse)) {
    throw new Error('Expected argument of type pactus.GetNewAddressResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_GetNewAddressResponse(buffer_arg) {
  return wallet_pb.GetNewAddressResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_GetTotalBalanceRequest(arg) {
  if (!(arg instanceof wallet_pb.GetTotalBalanceRequest)) {
    throw new Error('Expected argument of type pactus.GetTotalBalanceRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_GetTotalBalanceRequest(buffer_arg) {
  return wallet_pb.GetTotalBalanceRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_GetTotalBalanceResponse(arg) {
  if (!(arg instanceof wallet_pb.GetTotalBalanceResponse)) {
    throw new Error('Expected argument of type pactus.GetTotalBalanceResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_GetTotalBalanceResponse(buffer_arg) {
  return wallet_pb.GetTotalBalanceResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_GetTotalStakeRequest(arg) {
  if (!(arg instanceof wallet_pb.GetTotalStakeRequest)) {
    throw new Error('Expected argument of type pactus.GetTotalStakeRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_GetTotalStakeRequest(buffer_arg) {
  return wallet_pb.GetTotalStakeRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_GetTotalStakeResponse(arg) {
  if (!(arg instanceof wallet_pb.GetTotalStakeResponse)) {
    throw new Error('Expected argument of type pactus.GetTotalStakeResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_GetTotalStakeResponse(buffer_arg) {
  return wallet_pb.GetTotalStakeResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_GetValidatorAddressRequest(arg) {
  if (!(arg instanceof wallet_pb.GetValidatorAddressRequest)) {
    throw new Error('Expected argument of type pactus.GetValidatorAddressRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_GetValidatorAddressRequest(buffer_arg) {
  return wallet_pb.GetValidatorAddressRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_GetValidatorAddressResponse(arg) {
  if (!(arg instanceof wallet_pb.GetValidatorAddressResponse)) {
    throw new Error('Expected argument of type pactus.GetValidatorAddressResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_GetValidatorAddressResponse(buffer_arg) {
  return wallet_pb.GetValidatorAddressResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_GetWalletInfoRequest(arg) {
  if (!(arg instanceof wallet_pb.GetWalletInfoRequest)) {
    throw new Error('Expected argument of type pactus.GetWalletInfoRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_GetWalletInfoRequest(buffer_arg) {
  return wallet_pb.GetWalletInfoRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_GetWalletInfoResponse(arg) {
  if (!(arg instanceof wallet_pb.GetWalletInfoResponse)) {
    throw new Error('Expected argument of type pactus.GetWalletInfoResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_GetWalletInfoResponse(buffer_arg) {
  return wallet_pb.GetWalletInfoResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_ListAddressesRequest(arg) {
  if (!(arg instanceof wallet_pb.ListAddressesRequest)) {
    throw new Error('Expected argument of type pactus.ListAddressesRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_ListAddressesRequest(buffer_arg) {
  return wallet_pb.ListAddressesRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_ListAddressesResponse(arg) {
  if (!(arg instanceof wallet_pb.ListAddressesResponse)) {
    throw new Error('Expected argument of type pactus.ListAddressesResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_ListAddressesResponse(buffer_arg) {
  return wallet_pb.ListAddressesResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_ListWalletsRequest(arg) {
  if (!(arg instanceof wallet_pb.ListWalletsRequest)) {
    throw new Error('Expected argument of type pactus.ListWalletsRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_ListWalletsRequest(buffer_arg) {
  return wallet_pb.ListWalletsRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_ListWalletsResponse(arg) {
  if (!(arg instanceof wallet_pb.ListWalletsResponse)) {
    throw new Error('Expected argument of type pactus.ListWalletsResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_ListWalletsResponse(buffer_arg) {
  return wallet_pb.ListWalletsResponse.deserializeBinary(new Uint8Array(buffer_arg));
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

function serialize_pactus_RestoreWalletRequest(arg) {
  if (!(arg instanceof wallet_pb.RestoreWalletRequest)) {
    throw new Error('Expected argument of type pactus.RestoreWalletRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_RestoreWalletRequest(buffer_arg) {
  return wallet_pb.RestoreWalletRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_RestoreWalletResponse(arg) {
  if (!(arg instanceof wallet_pb.RestoreWalletResponse)) {
    throw new Error('Expected argument of type pactus.RestoreWalletResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_RestoreWalletResponse(buffer_arg) {
  return wallet_pb.RestoreWalletResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_SetAddressLabelRequest(arg) {
  if (!(arg instanceof wallet_pb.SetAddressLabelRequest)) {
    throw new Error('Expected argument of type pactus.SetAddressLabelRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_SetAddressLabelRequest(buffer_arg) {
  return wallet_pb.SetAddressLabelRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_SetAddressLabelResponse(arg) {
  if (!(arg instanceof wallet_pb.SetAddressLabelResponse)) {
    throw new Error('Expected argument of type pactus.SetAddressLabelResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_SetAddressLabelResponse(buffer_arg) {
  return wallet_pb.SetAddressLabelResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_SignMessageRequest(arg) {
  if (!(arg instanceof wallet_pb.SignMessageRequest)) {
    throw new Error('Expected argument of type pactus.SignMessageRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_SignMessageRequest(buffer_arg) {
  return wallet_pb.SignMessageRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_SignMessageResponse(arg) {
  if (!(arg instanceof wallet_pb.SignMessageResponse)) {
    throw new Error('Expected argument of type pactus.SignMessageResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_SignMessageResponse(buffer_arg) {
  return wallet_pb.SignMessageResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_SignRawTransactionRequest(arg) {
  if (!(arg instanceof wallet_pb.SignRawTransactionRequest)) {
    throw new Error('Expected argument of type pactus.SignRawTransactionRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_SignRawTransactionRequest(buffer_arg) {
  return wallet_pb.SignRawTransactionRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_SignRawTransactionResponse(arg) {
  if (!(arg instanceof wallet_pb.SignRawTransactionResponse)) {
    throw new Error('Expected argument of type pactus.SignRawTransactionResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_SignRawTransactionResponse(buffer_arg) {
  return wallet_pb.SignRawTransactionResponse.deserializeBinary(new Uint8Array(buffer_arg));
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


// Wallet service provides RPC methods for wallet management operations.
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
  // RestoreWallet restores an existing wallet with the given mnemonic.
restoreWallet: {
    path: '/pactus.Wallet/RestoreWallet',
    requestStream: false,
    responseStream: false,
    requestType: wallet_pb.RestoreWalletRequest,
    responseType: wallet_pb.RestoreWalletResponse,
    requestSerialize: serialize_pactus_RestoreWalletRequest,
    requestDeserialize: deserialize_pactus_RestoreWalletRequest,
    responseSerialize: serialize_pactus_RestoreWalletResponse,
    responseDeserialize: deserialize_pactus_RestoreWalletResponse,
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
  // GetTotalBalance returns the total available balance of the wallet.
getTotalBalance: {
    path: '/pactus.Wallet/GetTotalBalance',
    requestStream: false,
    responseStream: false,
    requestType: wallet_pb.GetTotalBalanceRequest,
    responseType: wallet_pb.GetTotalBalanceResponse,
    requestSerialize: serialize_pactus_GetTotalBalanceRequest,
    requestDeserialize: deserialize_pactus_GetTotalBalanceRequest,
    responseSerialize: serialize_pactus_GetTotalBalanceResponse,
    responseDeserialize: deserialize_pactus_GetTotalBalanceResponse,
  },
  // SignRawTransaction signs a raw transaction for a specified wallet.
signRawTransaction: {
    path: '/pactus.Wallet/SignRawTransaction',
    requestStream: false,
    responseStream: false,
    requestType: wallet_pb.SignRawTransactionRequest,
    responseType: wallet_pb.SignRawTransactionResponse,
    requestSerialize: serialize_pactus_SignRawTransactionRequest,
    requestDeserialize: deserialize_pactus_SignRawTransactionRequest,
    responseSerialize: serialize_pactus_SignRawTransactionResponse,
    responseDeserialize: deserialize_pactus_SignRawTransactionResponse,
  },
  // GetValidatorAddress retrieves the validator address associated with a public key.
// Deprecated: Will move into utils.
getValidatorAddress: {
    path: '/pactus.Wallet/GetValidatorAddress',
    requestStream: false,
    responseStream: false,
    requestType: wallet_pb.GetValidatorAddressRequest,
    responseType: wallet_pb.GetValidatorAddressResponse,
    requestSerialize: serialize_pactus_GetValidatorAddressRequest,
    requestDeserialize: deserialize_pactus_GetValidatorAddressRequest,
    responseSerialize: serialize_pactus_GetValidatorAddressResponse,
    responseDeserialize: deserialize_pactus_GetValidatorAddressResponse,
  },
  // GetNewAddress generates a new address for the specified wallet.
getNewAddress: {
    path: '/pactus.Wallet/GetNewAddress',
    requestStream: false,
    responseStream: false,
    requestType: wallet_pb.GetNewAddressRequest,
    responseType: wallet_pb.GetNewAddressResponse,
    requestSerialize: serialize_pactus_GetNewAddressRequest,
    requestDeserialize: deserialize_pactus_GetNewAddressRequest,
    responseSerialize: serialize_pactus_GetNewAddressResponse,
    responseDeserialize: deserialize_pactus_GetNewAddressResponse,
  },
  // GetAddressHistory retrieves the transaction history of an address.
getAddressHistory: {
    path: '/pactus.Wallet/GetAddressHistory',
    requestStream: false,
    responseStream: false,
    requestType: wallet_pb.GetAddressHistoryRequest,
    responseType: wallet_pb.GetAddressHistoryResponse,
    requestSerialize: serialize_pactus_GetAddressHistoryRequest,
    requestDeserialize: deserialize_pactus_GetAddressHistoryRequest,
    responseSerialize: serialize_pactus_GetAddressHistoryResponse,
    responseDeserialize: deserialize_pactus_GetAddressHistoryResponse,
  },
  // SignMessage signs an arbitrary message using a wallet's private key.
signMessage: {
    path: '/pactus.Wallet/SignMessage',
    requestStream: false,
    responseStream: false,
    requestType: wallet_pb.SignMessageRequest,
    responseType: wallet_pb.SignMessageResponse,
    requestSerialize: serialize_pactus_SignMessageRequest,
    requestDeserialize: deserialize_pactus_SignMessageRequest,
    responseSerialize: serialize_pactus_SignMessageResponse,
    responseDeserialize: deserialize_pactus_SignMessageResponse,
  },
  // GetTotalStake returns the total stake amount in the wallet.
getTotalStake: {
    path: '/pactus.Wallet/GetTotalStake',
    requestStream: false,
    responseStream: false,
    requestType: wallet_pb.GetTotalStakeRequest,
    responseType: wallet_pb.GetTotalStakeResponse,
    requestSerialize: serialize_pactus_GetTotalStakeRequest,
    requestDeserialize: deserialize_pactus_GetTotalStakeRequest,
    responseSerialize: serialize_pactus_GetTotalStakeResponse,
    responseDeserialize: deserialize_pactus_GetTotalStakeResponse,
  },
  // GetAddressInfo returns detailed information about a specific address.
getAddressInfo: {
    path: '/pactus.Wallet/GetAddressInfo',
    requestStream: false,
    responseStream: false,
    requestType: wallet_pb.GetAddressInfoRequest,
    responseType: wallet_pb.GetAddressInfoResponse,
    requestSerialize: serialize_pactus_GetAddressInfoRequest,
    requestDeserialize: deserialize_pactus_GetAddressInfoRequest,
    responseSerialize: serialize_pactus_GetAddressInfoResponse,
    responseDeserialize: deserialize_pactus_GetAddressInfoResponse,
  },
  // SetAddressLabel sets or updates the label for a given address.
setAddressLabel: {
    path: '/pactus.Wallet/SetAddressLabel',
    requestStream: false,
    responseStream: false,
    requestType: wallet_pb.SetAddressLabelRequest,
    responseType: wallet_pb.SetAddressLabelResponse,
    requestSerialize: serialize_pactus_SetAddressLabelRequest,
    requestDeserialize: deserialize_pactus_SetAddressLabelRequest,
    responseSerialize: serialize_pactus_SetAddressLabelResponse,
    responseDeserialize: deserialize_pactus_SetAddressLabelResponse,
  },
  // ListWallets returns a list of all available wallets.
listWallets: {
    path: '/pactus.Wallet/ListWallets',
    requestStream: false,
    responseStream: false,
    requestType: wallet_pb.ListWalletsRequest,
    responseType: wallet_pb.ListWalletsResponse,
    requestSerialize: serialize_pactus_ListWalletsRequest,
    requestDeserialize: deserialize_pactus_ListWalletsRequest,
    responseSerialize: serialize_pactus_ListWalletsResponse,
    responseDeserialize: deserialize_pactus_ListWalletsResponse,
  },
  // GetWalletInfo returns detailed information about a specific wallet.
getWalletInfo: {
    path: '/pactus.Wallet/GetWalletInfo',
    requestStream: false,
    responseStream: false,
    requestType: wallet_pb.GetWalletInfoRequest,
    responseType: wallet_pb.GetWalletInfoResponse,
    requestSerialize: serialize_pactus_GetWalletInfoRequest,
    requestDeserialize: deserialize_pactus_GetWalletInfoRequest,
    responseSerialize: serialize_pactus_GetWalletInfoResponse,
    responseDeserialize: deserialize_pactus_GetWalletInfoResponse,
  },
  // ListAddresses returns all addresses in the specified wallet.
listAddresses: {
    path: '/pactus.Wallet/ListAddresses',
    requestStream: false,
    responseStream: false,
    requestType: wallet_pb.ListAddressesRequest,
    responseType: wallet_pb.ListAddressesResponse,
    requestSerialize: serialize_pactus_ListAddressesRequest,
    requestDeserialize: deserialize_pactus_ListAddressesRequest,
    responseSerialize: serialize_pactus_ListAddressesResponse,
    responseDeserialize: deserialize_pactus_ListAddressesResponse,
  },
};

exports.WalletClient = grpc.makeGenericClientConstructor(WalletService, 'Wallet');
