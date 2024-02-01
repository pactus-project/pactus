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

function serialize_pactus_GetNodeInfoRequest(arg) {
  if (!(arg instanceof network_pb.GetNodeInfoRequest)) {
    throw new Error('Expected argument of type pactus.GetNodeInfoRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_GetNodeInfoRequest(buffer_arg) {
  return network_pb.GetNodeInfoRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_GetNodeInfoResponse(arg) {
  if (!(arg instanceof network_pb.GetNodeInfoResponse)) {
    throw new Error('Expected argument of type pactus.GetNodeInfoResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_GetNodeInfoResponse(buffer_arg) {
  return network_pb.GetNodeInfoResponse.deserializeBinary(new Uint8Array(buffer_arg));
}


// Network service provides RPCs for retrieving information about the network.
var NetworkService = exports.NetworkService = {
  // GetNetworkInfo retrieves information about the overall network.
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
  // GetNodeInfo retrieves information about a specific node in the network.
getNodeInfo: {
    path: '/pactus.Network/GetNodeInfo',
    requestStream: false,
    responseStream: false,
    requestType: network_pb.GetNodeInfoRequest,
    responseType: network_pb.GetNodeInfoResponse,
    requestSerialize: serialize_pactus_GetNodeInfoRequest,
    requestDeserialize: deserialize_pactus_GetNodeInfoRequest,
    responseSerialize: serialize_pactus_GetNodeInfoResponse,
    responseDeserialize: deserialize_pactus_GetNodeInfoResponse,
  },
};

exports.NetworkClient = grpc.makeGenericClientConstructor(NetworkService);
