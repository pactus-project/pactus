// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('grpc');
var network_pb = require('./network_pb.js');

function serialize_pactus_NetworkInfoRequest(arg) {
  if (!(arg instanceof network_pb.NetworkInfoRequest)) {
    throw new Error('Expected argument of type pactus.NetworkInfoRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_NetworkInfoRequest(buffer_arg) {
  return network_pb.NetworkInfoRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_NetworkInfoResponse(arg) {
  if (!(arg instanceof network_pb.NetworkInfoResponse)) {
    throw new Error('Expected argument of type pactus.NetworkInfoResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_NetworkInfoResponse(buffer_arg) {
  return network_pb.NetworkInfoResponse.deserializeBinary(new Uint8Array(buffer_arg));
}


var NetworkService = exports.NetworkService = {
  getNetworkInfo: {
    path: '/pactus.Network/GetNetworkInfo',
    requestStream: false,
    responseStream: false,
    requestType: network_pb.NetworkInfoRequest,
    responseType: network_pb.NetworkInfoResponse,
    requestSerialize: serialize_pactus_NetworkInfoRequest,
    requestDeserialize: deserialize_pactus_NetworkInfoRequest,
    responseSerialize: serialize_pactus_NetworkInfoResponse,
    responseDeserialize: deserialize_pactus_NetworkInfoResponse,
  },
};

exports.NetworkClient = grpc.makeGenericClientConstructor(NetworkService);
