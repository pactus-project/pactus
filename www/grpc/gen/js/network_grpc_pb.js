// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('grpc');
var network_pb = require('./network_pb.js');

function serialize_pactus_GetNetworkInfoRequest(arg) {
  if (!(arg instanceof network_pb.GetNetworkInfoRequest)) {
    throw new Error('Expected argument of type pactus.GetNetworkInfoRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_GetNetworkInfoRequest(buffer_arg) {
  return network_pb.GetNetworkInfoRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_GetNetworkInfoResponse(arg) {
  if (!(arg instanceof network_pb.GetNetworkInfoResponse)) {
    throw new Error('Expected argument of type pactus.GetNetworkInfoResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_GetNetworkInfoResponse(buffer_arg) {
  return network_pb.GetNetworkInfoResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_GetPeerInfoRequest(arg) {
  if (!(arg instanceof network_pb.GetPeerInfoRequest)) {
    throw new Error('Expected argument of type pactus.GetPeerInfoRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_GetPeerInfoRequest(buffer_arg) {
  return network_pb.GetPeerInfoRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_GetPeerInfoResponse(arg) {
  if (!(arg instanceof network_pb.GetPeerInfoResponse)) {
    throw new Error('Expected argument of type pactus.GetPeerInfoResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_GetPeerInfoResponse(buffer_arg) {
  return network_pb.GetPeerInfoResponse.deserializeBinary(new Uint8Array(buffer_arg));
}


var NetworkService = exports.NetworkService = {
  getNetworkInfo: {
    path: '/pactus.Network/GetNetworkInfo',
    requestStream: false,
    responseStream: false,
    requestType: network_pb.GetNetworkInfoRequest,
    responseType: network_pb.GetNetworkInfoResponse,
    requestSerialize: serialize_pactus_GetNetworkInfoRequest,
    requestDeserialize: deserialize_pactus_GetNetworkInfoRequest,
    responseSerialize: serialize_pactus_GetNetworkInfoResponse,
    responseDeserialize: deserialize_pactus_GetNetworkInfoResponse,
  },
  getPeerInfo: {
    path: '/pactus.Network/GetPeerInfo',
    requestStream: false,
    responseStream: false,
    requestType: network_pb.GetPeerInfoRequest,
    responseType: network_pb.GetPeerInfoResponse,
    requestSerialize: serialize_pactus_GetPeerInfoRequest,
    requestDeserialize: deserialize_pactus_GetPeerInfoRequest,
    responseSerialize: serialize_pactus_GetPeerInfoResponse,
    responseDeserialize: deserialize_pactus_GetPeerInfoResponse,
  },
};

exports.NetworkClient = grpc.makeGenericClientConstructor(NetworkService);
