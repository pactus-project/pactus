package pactus;

import static io.grpc.MethodDescriptor.generateFullMethodName;

/**
 * <pre>
 * Blockchain service defines RPC methods for interacting with the blockchain.
 * </pre>
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.71.0)",
    comments = "Source: blockchain.proto")
@io.grpc.stub.annotations.GrpcGenerated
public final class BlockchainServiceGrpc {

  private BlockchainServiceGrpc() {}

  public static final java.lang.String SERVICE_NAME = "pactus.BlockchainService";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<pactus.Blockchain.GetBlockRequest,
      pactus.Blockchain.GetBlockResponse> getGetBlockMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetBlock",
      requestType = pactus.Blockchain.GetBlockRequest.class,
      responseType = pactus.Blockchain.GetBlockResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.Blockchain.GetBlockRequest,
      pactus.Blockchain.GetBlockResponse> getGetBlockMethod() {
    io.grpc.MethodDescriptor<pactus.Blockchain.GetBlockRequest, pactus.Blockchain.GetBlockResponse> getGetBlockMethod;
    if ((getGetBlockMethod = BlockchainServiceGrpc.getGetBlockMethod) == null) {
      synchronized (BlockchainServiceGrpc.class) {
        if ((getGetBlockMethod = BlockchainServiceGrpc.getGetBlockMethod) == null) {
          BlockchainServiceGrpc.getGetBlockMethod = getGetBlockMethod =
              io.grpc.MethodDescriptor.<pactus.Blockchain.GetBlockRequest, pactus.Blockchain.GetBlockResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetBlock"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Blockchain.GetBlockRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Blockchain.GetBlockResponse.getDefaultInstance()))
              .setSchemaDescriptor(new BlockchainServiceMethodDescriptorSupplier("GetBlock"))
              .build();
        }
      }
    }
    return getGetBlockMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.Blockchain.GetBlockHashRequest,
      pactus.Blockchain.GetBlockHashResponse> getGetBlockHashMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetBlockHash",
      requestType = pactus.Blockchain.GetBlockHashRequest.class,
      responseType = pactus.Blockchain.GetBlockHashResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.Blockchain.GetBlockHashRequest,
      pactus.Blockchain.GetBlockHashResponse> getGetBlockHashMethod() {
    io.grpc.MethodDescriptor<pactus.Blockchain.GetBlockHashRequest, pactus.Blockchain.GetBlockHashResponse> getGetBlockHashMethod;
    if ((getGetBlockHashMethod = BlockchainServiceGrpc.getGetBlockHashMethod) == null) {
      synchronized (BlockchainServiceGrpc.class) {
        if ((getGetBlockHashMethod = BlockchainServiceGrpc.getGetBlockHashMethod) == null) {
          BlockchainServiceGrpc.getGetBlockHashMethod = getGetBlockHashMethod =
              io.grpc.MethodDescriptor.<pactus.Blockchain.GetBlockHashRequest, pactus.Blockchain.GetBlockHashResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetBlockHash"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Blockchain.GetBlockHashRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Blockchain.GetBlockHashResponse.getDefaultInstance()))
              .setSchemaDescriptor(new BlockchainServiceMethodDescriptorSupplier("GetBlockHash"))
              .build();
        }
      }
    }
    return getGetBlockHashMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.Blockchain.GetBlockHeightRequest,
      pactus.Blockchain.GetBlockHeightResponse> getGetBlockHeightMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetBlockHeight",
      requestType = pactus.Blockchain.GetBlockHeightRequest.class,
      responseType = pactus.Blockchain.GetBlockHeightResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.Blockchain.GetBlockHeightRequest,
      pactus.Blockchain.GetBlockHeightResponse> getGetBlockHeightMethod() {
    io.grpc.MethodDescriptor<pactus.Blockchain.GetBlockHeightRequest, pactus.Blockchain.GetBlockHeightResponse> getGetBlockHeightMethod;
    if ((getGetBlockHeightMethod = BlockchainServiceGrpc.getGetBlockHeightMethod) == null) {
      synchronized (BlockchainServiceGrpc.class) {
        if ((getGetBlockHeightMethod = BlockchainServiceGrpc.getGetBlockHeightMethod) == null) {
          BlockchainServiceGrpc.getGetBlockHeightMethod = getGetBlockHeightMethod =
              io.grpc.MethodDescriptor.<pactus.Blockchain.GetBlockHeightRequest, pactus.Blockchain.GetBlockHeightResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetBlockHeight"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Blockchain.GetBlockHeightRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Blockchain.GetBlockHeightResponse.getDefaultInstance()))
              .setSchemaDescriptor(new BlockchainServiceMethodDescriptorSupplier("GetBlockHeight"))
              .build();
        }
      }
    }
    return getGetBlockHeightMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.Blockchain.GetBlockchainInfoRequest,
      pactus.Blockchain.GetBlockchainInfoResponse> getGetBlockchainInfoMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetBlockchainInfo",
      requestType = pactus.Blockchain.GetBlockchainInfoRequest.class,
      responseType = pactus.Blockchain.GetBlockchainInfoResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.Blockchain.GetBlockchainInfoRequest,
      pactus.Blockchain.GetBlockchainInfoResponse> getGetBlockchainInfoMethod() {
    io.grpc.MethodDescriptor<pactus.Blockchain.GetBlockchainInfoRequest, pactus.Blockchain.GetBlockchainInfoResponse> getGetBlockchainInfoMethod;
    if ((getGetBlockchainInfoMethod = BlockchainServiceGrpc.getGetBlockchainInfoMethod) == null) {
      synchronized (BlockchainServiceGrpc.class) {
        if ((getGetBlockchainInfoMethod = BlockchainServiceGrpc.getGetBlockchainInfoMethod) == null) {
          BlockchainServiceGrpc.getGetBlockchainInfoMethod = getGetBlockchainInfoMethod =
              io.grpc.MethodDescriptor.<pactus.Blockchain.GetBlockchainInfoRequest, pactus.Blockchain.GetBlockchainInfoResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetBlockchainInfo"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Blockchain.GetBlockchainInfoRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Blockchain.GetBlockchainInfoResponse.getDefaultInstance()))
              .setSchemaDescriptor(new BlockchainServiceMethodDescriptorSupplier("GetBlockchainInfo"))
              .build();
        }
      }
    }
    return getGetBlockchainInfoMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.Blockchain.GetConsensusInfoRequest,
      pactus.Blockchain.GetConsensusInfoResponse> getGetConsensusInfoMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetConsensusInfo",
      requestType = pactus.Blockchain.GetConsensusInfoRequest.class,
      responseType = pactus.Blockchain.GetConsensusInfoResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.Blockchain.GetConsensusInfoRequest,
      pactus.Blockchain.GetConsensusInfoResponse> getGetConsensusInfoMethod() {
    io.grpc.MethodDescriptor<pactus.Blockchain.GetConsensusInfoRequest, pactus.Blockchain.GetConsensusInfoResponse> getGetConsensusInfoMethod;
    if ((getGetConsensusInfoMethod = BlockchainServiceGrpc.getGetConsensusInfoMethod) == null) {
      synchronized (BlockchainServiceGrpc.class) {
        if ((getGetConsensusInfoMethod = BlockchainServiceGrpc.getGetConsensusInfoMethod) == null) {
          BlockchainServiceGrpc.getGetConsensusInfoMethod = getGetConsensusInfoMethod =
              io.grpc.MethodDescriptor.<pactus.Blockchain.GetConsensusInfoRequest, pactus.Blockchain.GetConsensusInfoResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetConsensusInfo"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Blockchain.GetConsensusInfoRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Blockchain.GetConsensusInfoResponse.getDefaultInstance()))
              .setSchemaDescriptor(new BlockchainServiceMethodDescriptorSupplier("GetConsensusInfo"))
              .build();
        }
      }
    }
    return getGetConsensusInfoMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.Blockchain.GetAccountRequest,
      pactus.Blockchain.GetAccountResponse> getGetAccountMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetAccount",
      requestType = pactus.Blockchain.GetAccountRequest.class,
      responseType = pactus.Blockchain.GetAccountResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.Blockchain.GetAccountRequest,
      pactus.Blockchain.GetAccountResponse> getGetAccountMethod() {
    io.grpc.MethodDescriptor<pactus.Blockchain.GetAccountRequest, pactus.Blockchain.GetAccountResponse> getGetAccountMethod;
    if ((getGetAccountMethod = BlockchainServiceGrpc.getGetAccountMethod) == null) {
      synchronized (BlockchainServiceGrpc.class) {
        if ((getGetAccountMethod = BlockchainServiceGrpc.getGetAccountMethod) == null) {
          BlockchainServiceGrpc.getGetAccountMethod = getGetAccountMethod =
              io.grpc.MethodDescriptor.<pactus.Blockchain.GetAccountRequest, pactus.Blockchain.GetAccountResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetAccount"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Blockchain.GetAccountRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Blockchain.GetAccountResponse.getDefaultInstance()))
              .setSchemaDescriptor(new BlockchainServiceMethodDescriptorSupplier("GetAccount"))
              .build();
        }
      }
    }
    return getGetAccountMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.Blockchain.GetValidatorRequest,
      pactus.Blockchain.GetValidatorResponse> getGetValidatorMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetValidator",
      requestType = pactus.Blockchain.GetValidatorRequest.class,
      responseType = pactus.Blockchain.GetValidatorResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.Blockchain.GetValidatorRequest,
      pactus.Blockchain.GetValidatorResponse> getGetValidatorMethod() {
    io.grpc.MethodDescriptor<pactus.Blockchain.GetValidatorRequest, pactus.Blockchain.GetValidatorResponse> getGetValidatorMethod;
    if ((getGetValidatorMethod = BlockchainServiceGrpc.getGetValidatorMethod) == null) {
      synchronized (BlockchainServiceGrpc.class) {
        if ((getGetValidatorMethod = BlockchainServiceGrpc.getGetValidatorMethod) == null) {
          BlockchainServiceGrpc.getGetValidatorMethod = getGetValidatorMethod =
              io.grpc.MethodDescriptor.<pactus.Blockchain.GetValidatorRequest, pactus.Blockchain.GetValidatorResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetValidator"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Blockchain.GetValidatorRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Blockchain.GetValidatorResponse.getDefaultInstance()))
              .setSchemaDescriptor(new BlockchainServiceMethodDescriptorSupplier("GetValidator"))
              .build();
        }
      }
    }
    return getGetValidatorMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.Blockchain.GetValidatorByNumberRequest,
      pactus.Blockchain.GetValidatorResponse> getGetValidatorByNumberMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetValidatorByNumber",
      requestType = pactus.Blockchain.GetValidatorByNumberRequest.class,
      responseType = pactus.Blockchain.GetValidatorResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.Blockchain.GetValidatorByNumberRequest,
      pactus.Blockchain.GetValidatorResponse> getGetValidatorByNumberMethod() {
    io.grpc.MethodDescriptor<pactus.Blockchain.GetValidatorByNumberRequest, pactus.Blockchain.GetValidatorResponse> getGetValidatorByNumberMethod;
    if ((getGetValidatorByNumberMethod = BlockchainServiceGrpc.getGetValidatorByNumberMethod) == null) {
      synchronized (BlockchainServiceGrpc.class) {
        if ((getGetValidatorByNumberMethod = BlockchainServiceGrpc.getGetValidatorByNumberMethod) == null) {
          BlockchainServiceGrpc.getGetValidatorByNumberMethod = getGetValidatorByNumberMethod =
              io.grpc.MethodDescriptor.<pactus.Blockchain.GetValidatorByNumberRequest, pactus.Blockchain.GetValidatorResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetValidatorByNumber"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Blockchain.GetValidatorByNumberRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Blockchain.GetValidatorResponse.getDefaultInstance()))
              .setSchemaDescriptor(new BlockchainServiceMethodDescriptorSupplier("GetValidatorByNumber"))
              .build();
        }
      }
    }
    return getGetValidatorByNumberMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.Blockchain.GetValidatorAddressesRequest,
      pactus.Blockchain.GetValidatorAddressesResponse> getGetValidatorAddressesMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetValidatorAddresses",
      requestType = pactus.Blockchain.GetValidatorAddressesRequest.class,
      responseType = pactus.Blockchain.GetValidatorAddressesResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.Blockchain.GetValidatorAddressesRequest,
      pactus.Blockchain.GetValidatorAddressesResponse> getGetValidatorAddressesMethod() {
    io.grpc.MethodDescriptor<pactus.Blockchain.GetValidatorAddressesRequest, pactus.Blockchain.GetValidatorAddressesResponse> getGetValidatorAddressesMethod;
    if ((getGetValidatorAddressesMethod = BlockchainServiceGrpc.getGetValidatorAddressesMethod) == null) {
      synchronized (BlockchainServiceGrpc.class) {
        if ((getGetValidatorAddressesMethod = BlockchainServiceGrpc.getGetValidatorAddressesMethod) == null) {
          BlockchainServiceGrpc.getGetValidatorAddressesMethod = getGetValidatorAddressesMethod =
              io.grpc.MethodDescriptor.<pactus.Blockchain.GetValidatorAddressesRequest, pactus.Blockchain.GetValidatorAddressesResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetValidatorAddresses"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Blockchain.GetValidatorAddressesRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Blockchain.GetValidatorAddressesResponse.getDefaultInstance()))
              .setSchemaDescriptor(new BlockchainServiceMethodDescriptorSupplier("GetValidatorAddresses"))
              .build();
        }
      }
    }
    return getGetValidatorAddressesMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.Blockchain.GetPublicKeyRequest,
      pactus.Blockchain.GetPublicKeyResponse> getGetPublicKeyMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetPublicKey",
      requestType = pactus.Blockchain.GetPublicKeyRequest.class,
      responseType = pactus.Blockchain.GetPublicKeyResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.Blockchain.GetPublicKeyRequest,
      pactus.Blockchain.GetPublicKeyResponse> getGetPublicKeyMethod() {
    io.grpc.MethodDescriptor<pactus.Blockchain.GetPublicKeyRequest, pactus.Blockchain.GetPublicKeyResponse> getGetPublicKeyMethod;
    if ((getGetPublicKeyMethod = BlockchainServiceGrpc.getGetPublicKeyMethod) == null) {
      synchronized (BlockchainServiceGrpc.class) {
        if ((getGetPublicKeyMethod = BlockchainServiceGrpc.getGetPublicKeyMethod) == null) {
          BlockchainServiceGrpc.getGetPublicKeyMethod = getGetPublicKeyMethod =
              io.grpc.MethodDescriptor.<pactus.Blockchain.GetPublicKeyRequest, pactus.Blockchain.GetPublicKeyResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetPublicKey"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Blockchain.GetPublicKeyRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Blockchain.GetPublicKeyResponse.getDefaultInstance()))
              .setSchemaDescriptor(new BlockchainServiceMethodDescriptorSupplier("GetPublicKey"))
              .build();
        }
      }
    }
    return getGetPublicKeyMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.Blockchain.GetTxPoolContentRequest,
      pactus.Blockchain.GetTxPoolContentResponse> getGetTxPoolContentMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetTxPoolContent",
      requestType = pactus.Blockchain.GetTxPoolContentRequest.class,
      responseType = pactus.Blockchain.GetTxPoolContentResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.Blockchain.GetTxPoolContentRequest,
      pactus.Blockchain.GetTxPoolContentResponse> getGetTxPoolContentMethod() {
    io.grpc.MethodDescriptor<pactus.Blockchain.GetTxPoolContentRequest, pactus.Blockchain.GetTxPoolContentResponse> getGetTxPoolContentMethod;
    if ((getGetTxPoolContentMethod = BlockchainServiceGrpc.getGetTxPoolContentMethod) == null) {
      synchronized (BlockchainServiceGrpc.class) {
        if ((getGetTxPoolContentMethod = BlockchainServiceGrpc.getGetTxPoolContentMethod) == null) {
          BlockchainServiceGrpc.getGetTxPoolContentMethod = getGetTxPoolContentMethod =
              io.grpc.MethodDescriptor.<pactus.Blockchain.GetTxPoolContentRequest, pactus.Blockchain.GetTxPoolContentResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetTxPoolContent"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Blockchain.GetTxPoolContentRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Blockchain.GetTxPoolContentResponse.getDefaultInstance()))
              .setSchemaDescriptor(new BlockchainServiceMethodDescriptorSupplier("GetTxPoolContent"))
              .build();
        }
      }
    }
    return getGetTxPoolContentMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static BlockchainServiceStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<BlockchainServiceStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<BlockchainServiceStub>() {
        @java.lang.Override
        public BlockchainServiceStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new BlockchainServiceStub(channel, callOptions);
        }
      };
    return BlockchainServiceStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports all types of calls on the service
   */
  public static BlockchainServiceBlockingV2Stub newBlockingV2Stub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<BlockchainServiceBlockingV2Stub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<BlockchainServiceBlockingV2Stub>() {
        @java.lang.Override
        public BlockchainServiceBlockingV2Stub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new BlockchainServiceBlockingV2Stub(channel, callOptions);
        }
      };
    return BlockchainServiceBlockingV2Stub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static BlockchainServiceBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<BlockchainServiceBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<BlockchainServiceBlockingStub>() {
        @java.lang.Override
        public BlockchainServiceBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new BlockchainServiceBlockingStub(channel, callOptions);
        }
      };
    return BlockchainServiceBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static BlockchainServiceFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<BlockchainServiceFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<BlockchainServiceFutureStub>() {
        @java.lang.Override
        public BlockchainServiceFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new BlockchainServiceFutureStub(channel, callOptions);
        }
      };
    return BlockchainServiceFutureStub.newStub(factory, channel);
  }

  /**
   * <pre>
   * Blockchain service defines RPC methods for interacting with the blockchain.
   * </pre>
   */
  public interface AsyncService {

    /**
     * <pre>
     * GetBlock retrieves information about a block based on the provided request parameters.
     * </pre>
     */
    default void getBlock(pactus.Blockchain.GetBlockRequest request,
        io.grpc.stub.StreamObserver<pactus.Blockchain.GetBlockResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetBlockMethod(), responseObserver);
    }

    /**
     * <pre>
     * GetBlockHash retrieves the hash of a block at the specified height.
     * </pre>
     */
    default void getBlockHash(pactus.Blockchain.GetBlockHashRequest request,
        io.grpc.stub.StreamObserver<pactus.Blockchain.GetBlockHashResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetBlockHashMethod(), responseObserver);
    }

    /**
     * <pre>
     * GetBlockHeight retrieves the height of a block with the specified hash.
     * </pre>
     */
    default void getBlockHeight(pactus.Blockchain.GetBlockHeightRequest request,
        io.grpc.stub.StreamObserver<pactus.Blockchain.GetBlockHeightResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetBlockHeightMethod(), responseObserver);
    }

    /**
     * <pre>
     * GetBlockchainInfo retrieves general information about the blockchain.
     * </pre>
     */
    default void getBlockchainInfo(pactus.Blockchain.GetBlockchainInfoRequest request,
        io.grpc.stub.StreamObserver<pactus.Blockchain.GetBlockchainInfoResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetBlockchainInfoMethod(), responseObserver);
    }

    /**
     * <pre>
     * GetConsensusInfo retrieves information about the consensus instances.
     * </pre>
     */
    default void getConsensusInfo(pactus.Blockchain.GetConsensusInfoRequest request,
        io.grpc.stub.StreamObserver<pactus.Blockchain.GetConsensusInfoResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetConsensusInfoMethod(), responseObserver);
    }

    /**
     * <pre>
     * GetAccount retrieves information about an account based on the provided address.
     * </pre>
     */
    default void getAccount(pactus.Blockchain.GetAccountRequest request,
        io.grpc.stub.StreamObserver<pactus.Blockchain.GetAccountResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetAccountMethod(), responseObserver);
    }

    /**
     * <pre>
     * GetValidator retrieves information about a validator based on the provided address.
     * buf:lint:ignore RPC_REQUEST_RESPONSE_UNIQUE
     * buf:lint:ignore RPC_RESPONSE_STANDARD_NAME
     * </pre>
     */
    default void getValidator(pactus.Blockchain.GetValidatorRequest request,
        io.grpc.stub.StreamObserver<pactus.Blockchain.GetValidatorResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetValidatorMethod(), responseObserver);
    }

    /**
     * <pre>
     * GetValidatorByNumber retrieves information about a validator based on the provided number.
     * buf:lint:ignore RPC_REQUEST_RESPONSE_UNIQUE
     * buf:lint:ignore RPC_RESPONSE_STANDARD_NAME
     * </pre>
     */
    default void getValidatorByNumber(pactus.Blockchain.GetValidatorByNumberRequest request,
        io.grpc.stub.StreamObserver<pactus.Blockchain.GetValidatorResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetValidatorByNumberMethod(), responseObserver);
    }

    /**
     * <pre>
     * GetValidatorAddresses retrieves a list of all validator addresses.
     * </pre>
     */
    default void getValidatorAddresses(pactus.Blockchain.GetValidatorAddressesRequest request,
        io.grpc.stub.StreamObserver<pactus.Blockchain.GetValidatorAddressesResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetValidatorAddressesMethod(), responseObserver);
    }

    /**
     * <pre>
     * GetPublicKey retrieves the public key of an account based on the provided address.
     * </pre>
     */
    default void getPublicKey(pactus.Blockchain.GetPublicKeyRequest request,
        io.grpc.stub.StreamObserver<pactus.Blockchain.GetPublicKeyResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetPublicKeyMethod(), responseObserver);
    }

    /**
     * <pre>
     * GetTxPoolContent retrieves current transactions in the transaction pool.
     * </pre>
     */
    default void getTxPoolContent(pactus.Blockchain.GetTxPoolContentRequest request,
        io.grpc.stub.StreamObserver<pactus.Blockchain.GetTxPoolContentResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetTxPoolContentMethod(), responseObserver);
    }
  }

  /**
   * Base class for the server implementation of the service BlockchainService.
   * <pre>
   * Blockchain service defines RPC methods for interacting with the blockchain.
   * </pre>
   */
  public static abstract class BlockchainServiceImplBase
      implements io.grpc.BindableService, AsyncService {

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return BlockchainServiceGrpc.bindService(this);
    }
  }

  /**
   * A stub to allow clients to do asynchronous rpc calls to service BlockchainService.
   * <pre>
   * Blockchain service defines RPC methods for interacting with the blockchain.
   * </pre>
   */
  public static final class BlockchainServiceStub
      extends io.grpc.stub.AbstractAsyncStub<BlockchainServiceStub> {
    private BlockchainServiceStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected BlockchainServiceStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new BlockchainServiceStub(channel, callOptions);
    }

    /**
     * <pre>
     * GetBlock retrieves information about a block based on the provided request parameters.
     * </pre>
     */
    public void getBlock(pactus.Blockchain.GetBlockRequest request,
        io.grpc.stub.StreamObserver<pactus.Blockchain.GetBlockResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetBlockMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * GetBlockHash retrieves the hash of a block at the specified height.
     * </pre>
     */
    public void getBlockHash(pactus.Blockchain.GetBlockHashRequest request,
        io.grpc.stub.StreamObserver<pactus.Blockchain.GetBlockHashResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetBlockHashMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * GetBlockHeight retrieves the height of a block with the specified hash.
     * </pre>
     */
    public void getBlockHeight(pactus.Blockchain.GetBlockHeightRequest request,
        io.grpc.stub.StreamObserver<pactus.Blockchain.GetBlockHeightResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetBlockHeightMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * GetBlockchainInfo retrieves general information about the blockchain.
     * </pre>
     */
    public void getBlockchainInfo(pactus.Blockchain.GetBlockchainInfoRequest request,
        io.grpc.stub.StreamObserver<pactus.Blockchain.GetBlockchainInfoResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetBlockchainInfoMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * GetConsensusInfo retrieves information about the consensus instances.
     * </pre>
     */
    public void getConsensusInfo(pactus.Blockchain.GetConsensusInfoRequest request,
        io.grpc.stub.StreamObserver<pactus.Blockchain.GetConsensusInfoResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetConsensusInfoMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * GetAccount retrieves information about an account based on the provided address.
     * </pre>
     */
    public void getAccount(pactus.Blockchain.GetAccountRequest request,
        io.grpc.stub.StreamObserver<pactus.Blockchain.GetAccountResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetAccountMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * GetValidator retrieves information about a validator based on the provided address.
     * buf:lint:ignore RPC_REQUEST_RESPONSE_UNIQUE
     * buf:lint:ignore RPC_RESPONSE_STANDARD_NAME
     * </pre>
     */
    public void getValidator(pactus.Blockchain.GetValidatorRequest request,
        io.grpc.stub.StreamObserver<pactus.Blockchain.GetValidatorResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetValidatorMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * GetValidatorByNumber retrieves information about a validator based on the provided number.
     * buf:lint:ignore RPC_REQUEST_RESPONSE_UNIQUE
     * buf:lint:ignore RPC_RESPONSE_STANDARD_NAME
     * </pre>
     */
    public void getValidatorByNumber(pactus.Blockchain.GetValidatorByNumberRequest request,
        io.grpc.stub.StreamObserver<pactus.Blockchain.GetValidatorResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetValidatorByNumberMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * GetValidatorAddresses retrieves a list of all validator addresses.
     * </pre>
     */
    public void getValidatorAddresses(pactus.Blockchain.GetValidatorAddressesRequest request,
        io.grpc.stub.StreamObserver<pactus.Blockchain.GetValidatorAddressesResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetValidatorAddressesMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * GetPublicKey retrieves the public key of an account based on the provided address.
     * </pre>
     */
    public void getPublicKey(pactus.Blockchain.GetPublicKeyRequest request,
        io.grpc.stub.StreamObserver<pactus.Blockchain.GetPublicKeyResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetPublicKeyMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * GetTxPoolContent retrieves current transactions in the transaction pool.
     * </pre>
     */
    public void getTxPoolContent(pactus.Blockchain.GetTxPoolContentRequest request,
        io.grpc.stub.StreamObserver<pactus.Blockchain.GetTxPoolContentResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetTxPoolContentMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   * A stub to allow clients to do synchronous rpc calls to service BlockchainService.
   * <pre>
   * Blockchain service defines RPC methods for interacting with the blockchain.
   * </pre>
   */
  public static final class BlockchainServiceBlockingV2Stub
      extends io.grpc.stub.AbstractBlockingStub<BlockchainServiceBlockingV2Stub> {
    private BlockchainServiceBlockingV2Stub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected BlockchainServiceBlockingV2Stub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new BlockchainServiceBlockingV2Stub(channel, callOptions);
    }

    /**
     * <pre>
     * GetBlock retrieves information about a block based on the provided request parameters.
     * </pre>
     */
    public pactus.Blockchain.GetBlockResponse getBlock(pactus.Blockchain.GetBlockRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetBlockMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetBlockHash retrieves the hash of a block at the specified height.
     * </pre>
     */
    public pactus.Blockchain.GetBlockHashResponse getBlockHash(pactus.Blockchain.GetBlockHashRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetBlockHashMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetBlockHeight retrieves the height of a block with the specified hash.
     * </pre>
     */
    public pactus.Blockchain.GetBlockHeightResponse getBlockHeight(pactus.Blockchain.GetBlockHeightRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetBlockHeightMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetBlockchainInfo retrieves general information about the blockchain.
     * </pre>
     */
    public pactus.Blockchain.GetBlockchainInfoResponse getBlockchainInfo(pactus.Blockchain.GetBlockchainInfoRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetBlockchainInfoMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetConsensusInfo retrieves information about the consensus instances.
     * </pre>
     */
    public pactus.Blockchain.GetConsensusInfoResponse getConsensusInfo(pactus.Blockchain.GetConsensusInfoRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetConsensusInfoMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetAccount retrieves information about an account based on the provided address.
     * </pre>
     */
    public pactus.Blockchain.GetAccountResponse getAccount(pactus.Blockchain.GetAccountRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetAccountMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetValidator retrieves information about a validator based on the provided address.
     * buf:lint:ignore RPC_REQUEST_RESPONSE_UNIQUE
     * buf:lint:ignore RPC_RESPONSE_STANDARD_NAME
     * </pre>
     */
    public pactus.Blockchain.GetValidatorResponse getValidator(pactus.Blockchain.GetValidatorRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetValidatorMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetValidatorByNumber retrieves information about a validator based on the provided number.
     * buf:lint:ignore RPC_REQUEST_RESPONSE_UNIQUE
     * buf:lint:ignore RPC_RESPONSE_STANDARD_NAME
     * </pre>
     */
    public pactus.Blockchain.GetValidatorResponse getValidatorByNumber(pactus.Blockchain.GetValidatorByNumberRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetValidatorByNumberMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetValidatorAddresses retrieves a list of all validator addresses.
     * </pre>
     */
    public pactus.Blockchain.GetValidatorAddressesResponse getValidatorAddresses(pactus.Blockchain.GetValidatorAddressesRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetValidatorAddressesMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetPublicKey retrieves the public key of an account based on the provided address.
     * </pre>
     */
    public pactus.Blockchain.GetPublicKeyResponse getPublicKey(pactus.Blockchain.GetPublicKeyRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetPublicKeyMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetTxPoolContent retrieves current transactions in the transaction pool.
     * </pre>
     */
    public pactus.Blockchain.GetTxPoolContentResponse getTxPoolContent(pactus.Blockchain.GetTxPoolContentRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetTxPoolContentMethod(), getCallOptions(), request);
    }
  }

  /**
   * A stub to allow clients to do limited synchronous rpc calls to service BlockchainService.
   * <pre>
   * Blockchain service defines RPC methods for interacting with the blockchain.
   * </pre>
   */
  public static final class BlockchainServiceBlockingStub
      extends io.grpc.stub.AbstractBlockingStub<BlockchainServiceBlockingStub> {
    private BlockchainServiceBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected BlockchainServiceBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new BlockchainServiceBlockingStub(channel, callOptions);
    }

    /**
     * <pre>
     * GetBlock retrieves information about a block based on the provided request parameters.
     * </pre>
     */
    public pactus.Blockchain.GetBlockResponse getBlock(pactus.Blockchain.GetBlockRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetBlockMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetBlockHash retrieves the hash of a block at the specified height.
     * </pre>
     */
    public pactus.Blockchain.GetBlockHashResponse getBlockHash(pactus.Blockchain.GetBlockHashRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetBlockHashMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetBlockHeight retrieves the height of a block with the specified hash.
     * </pre>
     */
    public pactus.Blockchain.GetBlockHeightResponse getBlockHeight(pactus.Blockchain.GetBlockHeightRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetBlockHeightMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetBlockchainInfo retrieves general information about the blockchain.
     * </pre>
     */
    public pactus.Blockchain.GetBlockchainInfoResponse getBlockchainInfo(pactus.Blockchain.GetBlockchainInfoRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetBlockchainInfoMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetConsensusInfo retrieves information about the consensus instances.
     * </pre>
     */
    public pactus.Blockchain.GetConsensusInfoResponse getConsensusInfo(pactus.Blockchain.GetConsensusInfoRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetConsensusInfoMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetAccount retrieves information about an account based on the provided address.
     * </pre>
     */
    public pactus.Blockchain.GetAccountResponse getAccount(pactus.Blockchain.GetAccountRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetAccountMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetValidator retrieves information about a validator based on the provided address.
     * buf:lint:ignore RPC_REQUEST_RESPONSE_UNIQUE
     * buf:lint:ignore RPC_RESPONSE_STANDARD_NAME
     * </pre>
     */
    public pactus.Blockchain.GetValidatorResponse getValidator(pactus.Blockchain.GetValidatorRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetValidatorMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetValidatorByNumber retrieves information about a validator based on the provided number.
     * buf:lint:ignore RPC_REQUEST_RESPONSE_UNIQUE
     * buf:lint:ignore RPC_RESPONSE_STANDARD_NAME
     * </pre>
     */
    public pactus.Blockchain.GetValidatorResponse getValidatorByNumber(pactus.Blockchain.GetValidatorByNumberRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetValidatorByNumberMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetValidatorAddresses retrieves a list of all validator addresses.
     * </pre>
     */
    public pactus.Blockchain.GetValidatorAddressesResponse getValidatorAddresses(pactus.Blockchain.GetValidatorAddressesRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetValidatorAddressesMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetPublicKey retrieves the public key of an account based on the provided address.
     * </pre>
     */
    public pactus.Blockchain.GetPublicKeyResponse getPublicKey(pactus.Blockchain.GetPublicKeyRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetPublicKeyMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetTxPoolContent retrieves current transactions in the transaction pool.
     * </pre>
     */
    public pactus.Blockchain.GetTxPoolContentResponse getTxPoolContent(pactus.Blockchain.GetTxPoolContentRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetTxPoolContentMethod(), getCallOptions(), request);
    }
  }

  /**
   * A stub to allow clients to do ListenableFuture-style rpc calls to service BlockchainService.
   * <pre>
   * Blockchain service defines RPC methods for interacting with the blockchain.
   * </pre>
   */
  public static final class BlockchainServiceFutureStub
      extends io.grpc.stub.AbstractFutureStub<BlockchainServiceFutureStub> {
    private BlockchainServiceFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected BlockchainServiceFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new BlockchainServiceFutureStub(channel, callOptions);
    }

    /**
     * <pre>
     * GetBlock retrieves information about a block based on the provided request parameters.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.Blockchain.GetBlockResponse> getBlock(
        pactus.Blockchain.GetBlockRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetBlockMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * GetBlockHash retrieves the hash of a block at the specified height.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.Blockchain.GetBlockHashResponse> getBlockHash(
        pactus.Blockchain.GetBlockHashRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetBlockHashMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * GetBlockHeight retrieves the height of a block with the specified hash.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.Blockchain.GetBlockHeightResponse> getBlockHeight(
        pactus.Blockchain.GetBlockHeightRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetBlockHeightMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * GetBlockchainInfo retrieves general information about the blockchain.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.Blockchain.GetBlockchainInfoResponse> getBlockchainInfo(
        pactus.Blockchain.GetBlockchainInfoRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetBlockchainInfoMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * GetConsensusInfo retrieves information about the consensus instances.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.Blockchain.GetConsensusInfoResponse> getConsensusInfo(
        pactus.Blockchain.GetConsensusInfoRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetConsensusInfoMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * GetAccount retrieves information about an account based on the provided address.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.Blockchain.GetAccountResponse> getAccount(
        pactus.Blockchain.GetAccountRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetAccountMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * GetValidator retrieves information about a validator based on the provided address.
     * buf:lint:ignore RPC_REQUEST_RESPONSE_UNIQUE
     * buf:lint:ignore RPC_RESPONSE_STANDARD_NAME
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.Blockchain.GetValidatorResponse> getValidator(
        pactus.Blockchain.GetValidatorRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetValidatorMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * GetValidatorByNumber retrieves information about a validator based on the provided number.
     * buf:lint:ignore RPC_REQUEST_RESPONSE_UNIQUE
     * buf:lint:ignore RPC_RESPONSE_STANDARD_NAME
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.Blockchain.GetValidatorResponse> getValidatorByNumber(
        pactus.Blockchain.GetValidatorByNumberRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetValidatorByNumberMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * GetValidatorAddresses retrieves a list of all validator addresses.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.Blockchain.GetValidatorAddressesResponse> getValidatorAddresses(
        pactus.Blockchain.GetValidatorAddressesRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetValidatorAddressesMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * GetPublicKey retrieves the public key of an account based on the provided address.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.Blockchain.GetPublicKeyResponse> getPublicKey(
        pactus.Blockchain.GetPublicKeyRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetPublicKeyMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * GetTxPoolContent retrieves current transactions in the transaction pool.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.Blockchain.GetTxPoolContentResponse> getTxPoolContent(
        pactus.Blockchain.GetTxPoolContentRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetTxPoolContentMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_GET_BLOCK = 0;
  private static final int METHODID_GET_BLOCK_HASH = 1;
  private static final int METHODID_GET_BLOCK_HEIGHT = 2;
  private static final int METHODID_GET_BLOCKCHAIN_INFO = 3;
  private static final int METHODID_GET_CONSENSUS_INFO = 4;
  private static final int METHODID_GET_ACCOUNT = 5;
  private static final int METHODID_GET_VALIDATOR = 6;
  private static final int METHODID_GET_VALIDATOR_BY_NUMBER = 7;
  private static final int METHODID_GET_VALIDATOR_ADDRESSES = 8;
  private static final int METHODID_GET_PUBLIC_KEY = 9;
  private static final int METHODID_GET_TX_POOL_CONTENT = 10;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final AsyncService serviceImpl;
    private final int methodId;

    MethodHandlers(AsyncService serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_GET_BLOCK:
          serviceImpl.getBlock((pactus.Blockchain.GetBlockRequest) request,
              (io.grpc.stub.StreamObserver<pactus.Blockchain.GetBlockResponse>) responseObserver);
          break;
        case METHODID_GET_BLOCK_HASH:
          serviceImpl.getBlockHash((pactus.Blockchain.GetBlockHashRequest) request,
              (io.grpc.stub.StreamObserver<pactus.Blockchain.GetBlockHashResponse>) responseObserver);
          break;
        case METHODID_GET_BLOCK_HEIGHT:
          serviceImpl.getBlockHeight((pactus.Blockchain.GetBlockHeightRequest) request,
              (io.grpc.stub.StreamObserver<pactus.Blockchain.GetBlockHeightResponse>) responseObserver);
          break;
        case METHODID_GET_BLOCKCHAIN_INFO:
          serviceImpl.getBlockchainInfo((pactus.Blockchain.GetBlockchainInfoRequest) request,
              (io.grpc.stub.StreamObserver<pactus.Blockchain.GetBlockchainInfoResponse>) responseObserver);
          break;
        case METHODID_GET_CONSENSUS_INFO:
          serviceImpl.getConsensusInfo((pactus.Blockchain.GetConsensusInfoRequest) request,
              (io.grpc.stub.StreamObserver<pactus.Blockchain.GetConsensusInfoResponse>) responseObserver);
          break;
        case METHODID_GET_ACCOUNT:
          serviceImpl.getAccount((pactus.Blockchain.GetAccountRequest) request,
              (io.grpc.stub.StreamObserver<pactus.Blockchain.GetAccountResponse>) responseObserver);
          break;
        case METHODID_GET_VALIDATOR:
          serviceImpl.getValidator((pactus.Blockchain.GetValidatorRequest) request,
              (io.grpc.stub.StreamObserver<pactus.Blockchain.GetValidatorResponse>) responseObserver);
          break;
        case METHODID_GET_VALIDATOR_BY_NUMBER:
          serviceImpl.getValidatorByNumber((pactus.Blockchain.GetValidatorByNumberRequest) request,
              (io.grpc.stub.StreamObserver<pactus.Blockchain.GetValidatorResponse>) responseObserver);
          break;
        case METHODID_GET_VALIDATOR_ADDRESSES:
          serviceImpl.getValidatorAddresses((pactus.Blockchain.GetValidatorAddressesRequest) request,
              (io.grpc.stub.StreamObserver<pactus.Blockchain.GetValidatorAddressesResponse>) responseObserver);
          break;
        case METHODID_GET_PUBLIC_KEY:
          serviceImpl.getPublicKey((pactus.Blockchain.GetPublicKeyRequest) request,
              (io.grpc.stub.StreamObserver<pactus.Blockchain.GetPublicKeyResponse>) responseObserver);
          break;
        case METHODID_GET_TX_POOL_CONTENT:
          serviceImpl.getTxPoolContent((pactus.Blockchain.GetTxPoolContentRequest) request,
              (io.grpc.stub.StreamObserver<pactus.Blockchain.GetTxPoolContentResponse>) responseObserver);
          break;
        default:
          throw new AssertionError();
      }
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public io.grpc.stub.StreamObserver<Req> invoke(
        io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        default:
          throw new AssertionError();
      }
    }
  }

  public static final io.grpc.ServerServiceDefinition bindService(AsyncService service) {
    return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
        .addMethod(
          getGetBlockMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.Blockchain.GetBlockRequest,
              pactus.Blockchain.GetBlockResponse>(
                service, METHODID_GET_BLOCK)))
        .addMethod(
          getGetBlockHashMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.Blockchain.GetBlockHashRequest,
              pactus.Blockchain.GetBlockHashResponse>(
                service, METHODID_GET_BLOCK_HASH)))
        .addMethod(
          getGetBlockHeightMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.Blockchain.GetBlockHeightRequest,
              pactus.Blockchain.GetBlockHeightResponse>(
                service, METHODID_GET_BLOCK_HEIGHT)))
        .addMethod(
          getGetBlockchainInfoMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.Blockchain.GetBlockchainInfoRequest,
              pactus.Blockchain.GetBlockchainInfoResponse>(
                service, METHODID_GET_BLOCKCHAIN_INFO)))
        .addMethod(
          getGetConsensusInfoMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.Blockchain.GetConsensusInfoRequest,
              pactus.Blockchain.GetConsensusInfoResponse>(
                service, METHODID_GET_CONSENSUS_INFO)))
        .addMethod(
          getGetAccountMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.Blockchain.GetAccountRequest,
              pactus.Blockchain.GetAccountResponse>(
                service, METHODID_GET_ACCOUNT)))
        .addMethod(
          getGetValidatorMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.Blockchain.GetValidatorRequest,
              pactus.Blockchain.GetValidatorResponse>(
                service, METHODID_GET_VALIDATOR)))
        .addMethod(
          getGetValidatorByNumberMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.Blockchain.GetValidatorByNumberRequest,
              pactus.Blockchain.GetValidatorResponse>(
                service, METHODID_GET_VALIDATOR_BY_NUMBER)))
        .addMethod(
          getGetValidatorAddressesMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.Blockchain.GetValidatorAddressesRequest,
              pactus.Blockchain.GetValidatorAddressesResponse>(
                service, METHODID_GET_VALIDATOR_ADDRESSES)))
        .addMethod(
          getGetPublicKeyMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.Blockchain.GetPublicKeyRequest,
              pactus.Blockchain.GetPublicKeyResponse>(
                service, METHODID_GET_PUBLIC_KEY)))
        .addMethod(
          getGetTxPoolContentMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.Blockchain.GetTxPoolContentRequest,
              pactus.Blockchain.GetTxPoolContentResponse>(
                service, METHODID_GET_TX_POOL_CONTENT)))
        .build();
  }

  private static abstract class BlockchainServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    BlockchainServiceBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return pactus.Blockchain.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("BlockchainService");
    }
  }

  private static final class BlockchainServiceFileDescriptorSupplier
      extends BlockchainServiceBaseDescriptorSupplier {
    BlockchainServiceFileDescriptorSupplier() {}
  }

  private static final class BlockchainServiceMethodDescriptorSupplier
      extends BlockchainServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final java.lang.String methodName;

    BlockchainServiceMethodDescriptorSupplier(java.lang.String methodName) {
      this.methodName = methodName;
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.MethodDescriptor getMethodDescriptor() {
      return getServiceDescriptor().findMethodByName(methodName);
    }
  }

  private static volatile io.grpc.ServiceDescriptor serviceDescriptor;

  public static io.grpc.ServiceDescriptor getServiceDescriptor() {
    io.grpc.ServiceDescriptor result = serviceDescriptor;
    if (result == null) {
      synchronized (BlockchainServiceGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new BlockchainServiceFileDescriptorSupplier())
              .addMethod(getGetBlockMethod())
              .addMethod(getGetBlockHashMethod())
              .addMethod(getGetBlockHeightMethod())
              .addMethod(getGetBlockchainInfoMethod())
              .addMethod(getGetConsensusInfoMethod())
              .addMethod(getGetAccountMethod())
              .addMethod(getGetValidatorMethod())
              .addMethod(getGetValidatorByNumberMethod())
              .addMethod(getGetValidatorAddressesMethod())
              .addMethod(getGetPublicKeyMethod())
              .addMethod(getGetTxPoolContentMethod())
              .build();
        }
      }
    }
    return result;
  }
}
