/**
 * @fileoverview gRPC-Web generated client stub for pactus
 * @enhanceable
 * @public
 */

// Code generated by protoc-gen-grpc-web. DO NOT EDIT.
// versions:
// 	protoc-gen-grpc-web v1.5.0
// 	protoc              v0.0.0
// source: utils.proto


/* eslint-disable */
// @ts-nocheck



const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.pactus = require('./utils_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.pactus.UtilsClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options.format = 'binary';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname.replace(/\/+$/, '');

};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.pactus.UtilsPromiseClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options.format = 'binary';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname.replace(/\/+$/, '');

};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.pactus.SignMessageWithPrivateKeyRequest,
 *   !proto.pactus.SignMessageWithPrivateKeyResponse>}
 */
const methodDescriptor_Utils_SignMessageWithPrivateKey = new grpc.web.MethodDescriptor(
  '/pactus.Utils/SignMessageWithPrivateKey',
  grpc.web.MethodType.UNARY,
  proto.pactus.SignMessageWithPrivateKeyRequest,
  proto.pactus.SignMessageWithPrivateKeyResponse,
  /**
   * @param {!proto.pactus.SignMessageWithPrivateKeyRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.pactus.SignMessageWithPrivateKeyResponse.deserializeBinary
);


/**
 * @param {!proto.pactus.SignMessageWithPrivateKeyRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.pactus.SignMessageWithPrivateKeyResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.pactus.SignMessageWithPrivateKeyResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.pactus.UtilsClient.prototype.signMessageWithPrivateKey =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/pactus.Utils/SignMessageWithPrivateKey',
      request,
      metadata || {},
      methodDescriptor_Utils_SignMessageWithPrivateKey,
      callback);
};


/**
 * @param {!proto.pactus.SignMessageWithPrivateKeyRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.pactus.SignMessageWithPrivateKeyResponse>}
 *     Promise that resolves to the response
 */
proto.pactus.UtilsPromiseClient.prototype.signMessageWithPrivateKey =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/pactus.Utils/SignMessageWithPrivateKey',
      request,
      metadata || {},
      methodDescriptor_Utils_SignMessageWithPrivateKey);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.pactus.VerifyMessageRequest,
 *   !proto.pactus.VerifyMessageResponse>}
 */
const methodDescriptor_Utils_VerifyMessage = new grpc.web.MethodDescriptor(
  '/pactus.Utils/VerifyMessage',
  grpc.web.MethodType.UNARY,
  proto.pactus.VerifyMessageRequest,
  proto.pactus.VerifyMessageResponse,
  /**
   * @param {!proto.pactus.VerifyMessageRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.pactus.VerifyMessageResponse.deserializeBinary
);


/**
 * @param {!proto.pactus.VerifyMessageRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.pactus.VerifyMessageResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.pactus.VerifyMessageResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.pactus.UtilsClient.prototype.verifyMessage =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/pactus.Utils/VerifyMessage',
      request,
      metadata || {},
      methodDescriptor_Utils_VerifyMessage,
      callback);
};


/**
 * @param {!proto.pactus.VerifyMessageRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.pactus.VerifyMessageResponse>}
 *     Promise that resolves to the response
 */
proto.pactus.UtilsPromiseClient.prototype.verifyMessage =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/pactus.Utils/VerifyMessage',
      request,
      metadata || {},
      methodDescriptor_Utils_VerifyMessage);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.pactus.PublicKeyAggregationRequest,
 *   !proto.pactus.PublicKeyAggregationResponse>}
 */
const methodDescriptor_Utils_PublicKeyAggregation = new grpc.web.MethodDescriptor(
  '/pactus.Utils/PublicKeyAggregation',
  grpc.web.MethodType.UNARY,
  proto.pactus.PublicKeyAggregationRequest,
  proto.pactus.PublicKeyAggregationResponse,
  /**
   * @param {!proto.pactus.PublicKeyAggregationRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.pactus.PublicKeyAggregationResponse.deserializeBinary
);


/**
 * @param {!proto.pactus.PublicKeyAggregationRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.pactus.PublicKeyAggregationResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.pactus.PublicKeyAggregationResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.pactus.UtilsClient.prototype.publicKeyAggregation =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/pactus.Utils/PublicKeyAggregation',
      request,
      metadata || {},
      methodDescriptor_Utils_PublicKeyAggregation,
      callback);
};


/**
 * @param {!proto.pactus.PublicKeyAggregationRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.pactus.PublicKeyAggregationResponse>}
 *     Promise that resolves to the response
 */
proto.pactus.UtilsPromiseClient.prototype.publicKeyAggregation =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/pactus.Utils/PublicKeyAggregation',
      request,
      metadata || {},
      methodDescriptor_Utils_PublicKeyAggregation);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.pactus.SignatureAggregationRequest,
 *   !proto.pactus.SignatureAggregationResponse>}
 */
const methodDescriptor_Utils_SignatureAggregation = new grpc.web.MethodDescriptor(
  '/pactus.Utils/SignatureAggregation',
  grpc.web.MethodType.UNARY,
  proto.pactus.SignatureAggregationRequest,
  proto.pactus.SignatureAggregationResponse,
  /**
   * @param {!proto.pactus.SignatureAggregationRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.pactus.SignatureAggregationResponse.deserializeBinary
);


/**
 * @param {!proto.pactus.SignatureAggregationRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.pactus.SignatureAggregationResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.pactus.SignatureAggregationResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.pactus.UtilsClient.prototype.signatureAggregation =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/pactus.Utils/SignatureAggregation',
      request,
      metadata || {},
      methodDescriptor_Utils_SignatureAggregation,
      callback);
};


/**
 * @param {!proto.pactus.SignatureAggregationRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.pactus.SignatureAggregationResponse>}
 *     Promise that resolves to the response
 */
proto.pactus.UtilsPromiseClient.prototype.signatureAggregation =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/pactus.Utils/SignatureAggregation',
      request,
      metadata || {},
      methodDescriptor_Utils_SignatureAggregation);
};


module.exports = proto.pactus;

