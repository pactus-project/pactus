// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('@grpc/grpc-js');
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

function serialize_pactus_ListPeersRequest(arg) {
  if (!(arg instanceof network_pb.ListPeersRequest)) {
    throw new Error('Expected argument of type pactus.ListPeersRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_ListPeersRequest(buffer_arg) {
  return network_pb.ListPeersRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_ListPeersResponse(arg) {
  if (!(arg instanceof network_pb.ListPeersResponse)) {
    throw new Error('Expected argument of type pactus.ListPeersResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_ListPeersResponse(buffer_arg) {
  return network_pb.ListPeersResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_PingRequest(arg) {
  if (!(arg instanceof network_pb.PingRequest)) {
    throw new Error('Expected argument of type pactus.PingRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_PingRequest(buffer_arg) {
  return network_pb.PingRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_PingResponse(arg) {
  if (!(arg instanceof network_pb.PingResponse)) {
    throw new Error('Expected argument of type pactus.PingResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_PingResponse(buffer_arg) {
  return network_pb.PingResponse.deserializeBinary(new Uint8Array(buffer_arg));
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
  // ListPeers lists all peers in the network.
listPeers: {
    path: '/pactus.Network/ListPeers',
    requestStream: false,
    responseStream: false,
    requestType: network_pb.ListPeersRequest,
    responseType: network_pb.ListPeersResponse,
    requestSerialize: serialize_pactus_ListPeersRequest,
    requestDeserialize: deserialize_pactus_ListPeersRequest,
    responseSerialize: serialize_pactus_ListPeersResponse,
    responseDeserialize: deserialize_pactus_ListPeersResponse,
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
  // Ping provides a simple connectivity test and latency measurement.
ping: {
    path: '/pactus.Network/Ping',
    requestStream: false,
    responseStream: false,
    requestType: network_pb.PingRequest,
    responseType: network_pb.PingResponse,
    requestSerialize: serialize_pactus_PingRequest,
    requestDeserialize: deserialize_pactus_PingRequest,
    responseSerialize: serialize_pactus_PingResponse,
    responseDeserialize: deserialize_pactus_PingResponse,
  },
};

exports.NetworkClient = grpc.makeGenericClientConstructor(NetworkService, 'Network');
