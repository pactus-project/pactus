// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('grpc');
var utility_pb = require('./utility_pb.js');
var transaction_pb = require('./transaction_pb.js');

function serialize_pactus_CalculateFeeRequest(arg) {
  if (!(arg instanceof utility_pb.CalculateFeeRequest)) {
    throw new Error('Expected argument of type pactus.CalculateFeeRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_CalculateFeeRequest(buffer_arg) {
  return utility_pb.CalculateFeeRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_pactus_CalculateFeeResponse(arg) {
  if (!(arg instanceof utility_pb.CalculateFeeResponse)) {
    throw new Error('Expected argument of type pactus.CalculateFeeResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_pactus_CalculateFeeResponse(buffer_arg) {
  return utility_pb.CalculateFeeResponse.deserializeBinary(new Uint8Array(buffer_arg));
}


var UtilityService = exports.UtilityService = {
  // CalculateFee calculates the transaction fee based on the specified amount
// and payload type.
calculateFee: {
    path: '/pactus.Utility/CalculateFee',
    requestStream: false,
    responseStream: false,
    requestType: utility_pb.CalculateFeeRequest,
    responseType: utility_pb.CalculateFeeResponse,
    requestSerialize: serialize_pactus_CalculateFeeRequest,
    requestDeserialize: deserialize_pactus_CalculateFeeRequest,
    responseSerialize: serialize_pactus_CalculateFeeResponse,
    responseDeserialize: deserialize_pactus_CalculateFeeResponse,
  },
};

exports.UtilityClient = grpc.makeGenericClientConstructor(UtilityService);
