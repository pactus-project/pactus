package pactus.blockchain;

import static io.grpc.MethodDescriptor.generateFullMethodName;

/**
 * <pre>
 * Blockchain service defines RPC methods for interacting with the blockchain.
 * </pre>
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.50.2)",
    comments = "Source: blockchain.proto")
@io.grpc.stub.annotations.GrpcGenerated
public final class BlockchainGrpc {

  private BlockchainGrpc() {}

  public static final String SERVICE_NAME = "pactus.Blockchain";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<pactus.blockchain.BlockchainOuterClass.GetBlockRequest,
      pactus.blockchain.BlockchainOuterClass.GetBlockResponse> getGetBlockMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetBlock",
      requestType = pactus.blockchain.BlockchainOuterClass.GetBlockRequest.class,
      responseType = pactus.blockchain.BlockchainOuterClass.GetBlockResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.blockchain.BlockchainOuterClass.GetBlockRequest,
      pactus.blockchain.BlockchainOuterClass.GetBlockResponse> getGetBlockMethod() {
    io.grpc.MethodDescriptor<pactus.blockchain.BlockchainOuterClass.GetBlockRequest, pactus.blockchain.BlockchainOuterClass.GetBlockResponse> getGetBlockMethod;
    if ((getGetBlockMethod = BlockchainGrpc.getGetBlockMethod) == null) {
      synchronized (BlockchainGrpc.class) {
        if ((getGetBlockMethod = BlockchainGrpc.getGetBlockMethod) == null) {
          BlockchainGrpc.getGetBlockMethod = getGetBlockMethod =
              io.grpc.MethodDescriptor.<pactus.blockchain.BlockchainOuterClass.GetBlockRequest, pactus.blockchain.BlockchainOuterClass.GetBlockResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetBlock"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.blockchain.BlockchainOuterClass.GetBlockRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.blockchain.BlockchainOuterClass.GetBlockResponse.getDefaultInstance()))
              .setSchemaDescriptor(new BlockchainMethodDescriptorSupplier("GetBlock"))
              .build();
        }
      }
    }
    return getGetBlockMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.blockchain.BlockchainOuterClass.GetBlockHashRequest,
      pactus.blockchain.BlockchainOuterClass.GetBlockHashResponse> getGetBlockHashMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetBlockHash",
      requestType = pactus.blockchain.BlockchainOuterClass.GetBlockHashRequest.class,
      responseType = pactus.blockchain.BlockchainOuterClass.GetBlockHashResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.blockchain.BlockchainOuterClass.GetBlockHashRequest,
      pactus.blockchain.BlockchainOuterClass.GetBlockHashResponse> getGetBlockHashMethod() {
    io.grpc.MethodDescriptor<pactus.blockchain.BlockchainOuterClass.GetBlockHashRequest, pactus.blockchain.BlockchainOuterClass.GetBlockHashResponse> getGetBlockHashMethod;
    if ((getGetBlockHashMethod = BlockchainGrpc.getGetBlockHashMethod) == null) {
      synchronized (BlockchainGrpc.class) {
        if ((getGetBlockHashMethod = BlockchainGrpc.getGetBlockHashMethod) == null) {
          BlockchainGrpc.getGetBlockHashMethod = getGetBlockHashMethod =
              io.grpc.MethodDescriptor.<pactus.blockchain.BlockchainOuterClass.GetBlockHashRequest, pactus.blockchain.BlockchainOuterClass.GetBlockHashResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetBlockHash"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.blockchain.BlockchainOuterClass.GetBlockHashRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.blockchain.BlockchainOuterClass.GetBlockHashResponse.getDefaultInstance()))
              .setSchemaDescriptor(new BlockchainMethodDescriptorSupplier("GetBlockHash"))
              .build();
        }
      }
    }
    return getGetBlockHashMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.blockchain.BlockchainOuterClass.GetBlockHeightRequest,
      pactus.blockchain.BlockchainOuterClass.GetBlockHeightResponse> getGetBlockHeightMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetBlockHeight",
      requestType = pactus.blockchain.BlockchainOuterClass.GetBlockHeightRequest.class,
      responseType = pactus.blockchain.BlockchainOuterClass.GetBlockHeightResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.blockchain.BlockchainOuterClass.GetBlockHeightRequest,
      pactus.blockchain.BlockchainOuterClass.GetBlockHeightResponse> getGetBlockHeightMethod() {
    io.grpc.MethodDescriptor<pactus.blockchain.BlockchainOuterClass.GetBlockHeightRequest, pactus.blockchain.BlockchainOuterClass.GetBlockHeightResponse> getGetBlockHeightMethod;
    if ((getGetBlockHeightMethod = BlockchainGrpc.getGetBlockHeightMethod) == null) {
      synchronized (BlockchainGrpc.class) {
        if ((getGetBlockHeightMethod = BlockchainGrpc.getGetBlockHeightMethod) == null) {
          BlockchainGrpc.getGetBlockHeightMethod = getGetBlockHeightMethod =
              io.grpc.MethodDescriptor.<pactus.blockchain.BlockchainOuterClass.GetBlockHeightRequest, pactus.blockchain.BlockchainOuterClass.GetBlockHeightResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetBlockHeight"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.blockchain.BlockchainOuterClass.GetBlockHeightRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.blockchain.BlockchainOuterClass.GetBlockHeightResponse.getDefaultInstance()))
              .setSchemaDescriptor(new BlockchainMethodDescriptorSupplier("GetBlockHeight"))
              .build();
        }
      }
    }
    return getGetBlockHeightMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.blockchain.BlockchainOuterClass.GetBlockchainInfoRequest,
      pactus.blockchain.BlockchainOuterClass.GetBlockchainInfoResponse> getGetBlockchainInfoMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetBlockchainInfo",
      requestType = pactus.blockchain.BlockchainOuterClass.GetBlockchainInfoRequest.class,
      responseType = pactus.blockchain.BlockchainOuterClass.GetBlockchainInfoResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.blockchain.BlockchainOuterClass.GetBlockchainInfoRequest,
      pactus.blockchain.BlockchainOuterClass.GetBlockchainInfoResponse> getGetBlockchainInfoMethod() {
    io.grpc.MethodDescriptor<pactus.blockchain.BlockchainOuterClass.GetBlockchainInfoRequest, pactus.blockchain.BlockchainOuterClass.GetBlockchainInfoResponse> getGetBlockchainInfoMethod;
    if ((getGetBlockchainInfoMethod = BlockchainGrpc.getGetBlockchainInfoMethod) == null) {
      synchronized (BlockchainGrpc.class) {
        if ((getGetBlockchainInfoMethod = BlockchainGrpc.getGetBlockchainInfoMethod) == null) {
          BlockchainGrpc.getGetBlockchainInfoMethod = getGetBlockchainInfoMethod =
              io.grpc.MethodDescriptor.<pactus.blockchain.BlockchainOuterClass.GetBlockchainInfoRequest, pactus.blockchain.BlockchainOuterClass.GetBlockchainInfoResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetBlockchainInfo"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.blockchain.BlockchainOuterClass.GetBlockchainInfoRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.blockchain.BlockchainOuterClass.GetBlockchainInfoResponse.getDefaultInstance()))
              .setSchemaDescriptor(new BlockchainMethodDescriptorSupplier("GetBlockchainInfo"))
              .build();
        }
      }
    }
    return getGetBlockchainInfoMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.blockchain.BlockchainOuterClass.GetConsensusInfoRequest,
      pactus.blockchain.BlockchainOuterClass.GetConsensusInfoResponse> getGetConsensusInfoMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetConsensusInfo",
      requestType = pactus.blockchain.BlockchainOuterClass.GetConsensusInfoRequest.class,
      responseType = pactus.blockchain.BlockchainOuterClass.GetConsensusInfoResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.blockchain.BlockchainOuterClass.GetConsensusInfoRequest,
      pactus.blockchain.BlockchainOuterClass.GetConsensusInfoResponse> getGetConsensusInfoMethod() {
    io.grpc.MethodDescriptor<pactus.blockchain.BlockchainOuterClass.GetConsensusInfoRequest, pactus.blockchain.BlockchainOuterClass.GetConsensusInfoResponse> getGetConsensusInfoMethod;
    if ((getGetConsensusInfoMethod = BlockchainGrpc.getGetConsensusInfoMethod) == null) {
      synchronized (BlockchainGrpc.class) {
        if ((getGetConsensusInfoMethod = BlockchainGrpc.getGetConsensusInfoMethod) == null) {
          BlockchainGrpc.getGetConsensusInfoMethod = getGetConsensusInfoMethod =
              io.grpc.MethodDescriptor.<pactus.blockchain.BlockchainOuterClass.GetConsensusInfoRequest, pactus.blockchain.BlockchainOuterClass.GetConsensusInfoResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetConsensusInfo"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.blockchain.BlockchainOuterClass.GetConsensusInfoRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.blockchain.BlockchainOuterClass.GetConsensusInfoResponse.getDefaultInstance()))
              .setSchemaDescriptor(new BlockchainMethodDescriptorSupplier("GetConsensusInfo"))
              .build();
        }
      }
    }
    return getGetConsensusInfoMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.blockchain.BlockchainOuterClass.GetAccountRequest,
      pactus.blockchain.BlockchainOuterClass.GetAccountResponse> getGetAccountMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetAccount",
      requestType = pactus.blockchain.BlockchainOuterClass.GetAccountRequest.class,
      responseType = pactus.blockchain.BlockchainOuterClass.GetAccountResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.blockchain.BlockchainOuterClass.GetAccountRequest,
      pactus.blockchain.BlockchainOuterClass.GetAccountResponse> getGetAccountMethod() {
    io.grpc.MethodDescriptor<pactus.blockchain.BlockchainOuterClass.GetAccountRequest, pactus.blockchain.BlockchainOuterClass.GetAccountResponse> getGetAccountMethod;
    if ((getGetAccountMethod = BlockchainGrpc.getGetAccountMethod) == null) {
      synchronized (BlockchainGrpc.class) {
        if ((getGetAccountMethod = BlockchainGrpc.getGetAccountMethod) == null) {
          BlockchainGrpc.getGetAccountMethod = getGetAccountMethod =
              io.grpc.MethodDescriptor.<pactus.blockchain.BlockchainOuterClass.GetAccountRequest, pactus.blockchain.BlockchainOuterClass.GetAccountResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetAccount"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.blockchain.BlockchainOuterClass.GetAccountRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.blockchain.BlockchainOuterClass.GetAccountResponse.getDefaultInstance()))
              .setSchemaDescriptor(new BlockchainMethodDescriptorSupplier("GetAccount"))
              .build();
        }
      }
    }
    return getGetAccountMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.blockchain.BlockchainOuterClass.GetValidatorRequest,
      pactus.blockchain.BlockchainOuterClass.GetValidatorResponse> getGetValidatorMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetValidator",
      requestType = pactus.blockchain.BlockchainOuterClass.GetValidatorRequest.class,
      responseType = pactus.blockchain.BlockchainOuterClass.GetValidatorResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.blockchain.BlockchainOuterClass.GetValidatorRequest,
      pactus.blockchain.BlockchainOuterClass.GetValidatorResponse> getGetValidatorMethod() {
    io.grpc.MethodDescriptor<pactus.blockchain.BlockchainOuterClass.GetValidatorRequest, pactus.blockchain.BlockchainOuterClass.GetValidatorResponse> getGetValidatorMethod;
    if ((getGetValidatorMethod = BlockchainGrpc.getGetValidatorMethod) == null) {
      synchronized (BlockchainGrpc.class) {
        if ((getGetValidatorMethod = BlockchainGrpc.getGetValidatorMethod) == null) {
          BlockchainGrpc.getGetValidatorMethod = getGetValidatorMethod =
              io.grpc.MethodDescriptor.<pactus.blockchain.BlockchainOuterClass.GetValidatorRequest, pactus.blockchain.BlockchainOuterClass.GetValidatorResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetValidator"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.blockchain.BlockchainOuterClass.GetValidatorRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.blockchain.BlockchainOuterClass.GetValidatorResponse.getDefaultInstance()))
              .setSchemaDescriptor(new BlockchainMethodDescriptorSupplier("GetValidator"))
              .build();
        }
      }
    }
    return getGetValidatorMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.blockchain.BlockchainOuterClass.GetValidatorByNumberRequest,
      pactus.blockchain.BlockchainOuterClass.GetValidatorResponse> getGetValidatorByNumberMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetValidatorByNumber",
      requestType = pactus.blockchain.BlockchainOuterClass.GetValidatorByNumberRequest.class,
      responseType = pactus.blockchain.BlockchainOuterClass.GetValidatorResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.blockchain.BlockchainOuterClass.GetValidatorByNumberRequest,
      pactus.blockchain.BlockchainOuterClass.GetValidatorResponse> getGetValidatorByNumberMethod() {
    io.grpc.MethodDescriptor<pactus.blockchain.BlockchainOuterClass.GetValidatorByNumberRequest, pactus.blockchain.BlockchainOuterClass.GetValidatorResponse> getGetValidatorByNumberMethod;
    if ((getGetValidatorByNumberMethod = BlockchainGrpc.getGetValidatorByNumberMethod) == null) {
      synchronized (BlockchainGrpc.class) {
        if ((getGetValidatorByNumberMethod = BlockchainGrpc.getGetValidatorByNumberMethod) == null) {
          BlockchainGrpc.getGetValidatorByNumberMethod = getGetValidatorByNumberMethod =
              io.grpc.MethodDescriptor.<pactus.blockchain.BlockchainOuterClass.GetValidatorByNumberRequest, pactus.blockchain.BlockchainOuterClass.GetValidatorResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetValidatorByNumber"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.blockchain.BlockchainOuterClass.GetValidatorByNumberRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.blockchain.BlockchainOuterClass.GetValidatorResponse.getDefaultInstance()))
              .setSchemaDescriptor(new BlockchainMethodDescriptorSupplier("GetValidatorByNumber"))
              .build();
        }
      }
    }
    return getGetValidatorByNumberMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.blockchain.BlockchainOuterClass.GetValidatorAddressesRequest,
      pactus.blockchain.BlockchainOuterClass.GetValidatorAddressesResponse> getGetValidatorAddressesMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetValidatorAddresses",
      requestType = pactus.blockchain.BlockchainOuterClass.GetValidatorAddressesRequest.class,
      responseType = pactus.blockchain.BlockchainOuterClass.GetValidatorAddressesResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.blockchain.BlockchainOuterClass.GetValidatorAddressesRequest,
      pactus.blockchain.BlockchainOuterClass.GetValidatorAddressesResponse> getGetValidatorAddressesMethod() {
    io.grpc.MethodDescriptor<pactus.blockchain.BlockchainOuterClass.GetValidatorAddressesRequest, pactus.blockchain.BlockchainOuterClass.GetValidatorAddressesResponse> getGetValidatorAddressesMethod;
    if ((getGetValidatorAddressesMethod = BlockchainGrpc.getGetValidatorAddressesMethod) == null) {
      synchronized (BlockchainGrpc.class) {
        if ((getGetValidatorAddressesMethod = BlockchainGrpc.getGetValidatorAddressesMethod) == null) {
          BlockchainGrpc.getGetValidatorAddressesMethod = getGetValidatorAddressesMethod =
              io.grpc.MethodDescriptor.<pactus.blockchain.BlockchainOuterClass.GetValidatorAddressesRequest, pactus.blockchain.BlockchainOuterClass.GetValidatorAddressesResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetValidatorAddresses"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.blockchain.BlockchainOuterClass.GetValidatorAddressesRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.blockchain.BlockchainOuterClass.GetValidatorAddressesResponse.getDefaultInstance()))
              .setSchemaDescriptor(new BlockchainMethodDescriptorSupplier("GetValidatorAddresses"))
              .build();
        }
      }
    }
    return getGetValidatorAddressesMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.blockchain.BlockchainOuterClass.GetPublicKeyRequest,
      pactus.blockchain.BlockchainOuterClass.GetPublicKeyResponse> getGetPublicKeyMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetPublicKey",
      requestType = pactus.blockchain.BlockchainOuterClass.GetPublicKeyRequest.class,
      responseType = pactus.blockchain.BlockchainOuterClass.GetPublicKeyResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.blockchain.BlockchainOuterClass.GetPublicKeyRequest,
      pactus.blockchain.BlockchainOuterClass.GetPublicKeyResponse> getGetPublicKeyMethod() {
    io.grpc.MethodDescriptor<pactus.blockchain.BlockchainOuterClass.GetPublicKeyRequest, pactus.blockchain.BlockchainOuterClass.GetPublicKeyResponse> getGetPublicKeyMethod;
    if ((getGetPublicKeyMethod = BlockchainGrpc.getGetPublicKeyMethod) == null) {
      synchronized (BlockchainGrpc.class) {
        if ((getGetPublicKeyMethod = BlockchainGrpc.getGetPublicKeyMethod) == null) {
          BlockchainGrpc.getGetPublicKeyMethod = getGetPublicKeyMethod =
              io.grpc.MethodDescriptor.<pactus.blockchain.BlockchainOuterClass.GetPublicKeyRequest, pactus.blockchain.BlockchainOuterClass.GetPublicKeyResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetPublicKey"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.blockchain.BlockchainOuterClass.GetPublicKeyRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.blockchain.BlockchainOuterClass.GetPublicKeyResponse.getDefaultInstance()))
              .setSchemaDescriptor(new BlockchainMethodDescriptorSupplier("GetPublicKey"))
              .build();
        }
      }
    }
    return getGetPublicKeyMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static BlockchainStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<BlockchainStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<BlockchainStub>() {
        @java.lang.Override
        public BlockchainStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new BlockchainStub(channel, callOptions);
        }
      };
    return BlockchainStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static BlockchainBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<BlockchainBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<BlockchainBlockingStub>() {
        @java.lang.Override
        public BlockchainBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new BlockchainBlockingStub(channel, callOptions);
        }
      };
    return BlockchainBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static BlockchainFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<BlockchainFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<BlockchainFutureStub>() {
        @java.lang.Override
        public BlockchainFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new BlockchainFutureStub(channel, callOptions);
        }
      };
    return BlockchainFutureStub.newStub(factory, channel);
  }

  /**
   * <pre>
   * Blockchain service defines RPC methods for interacting with the blockchain.
   * </pre>
   */
  public static abstract class BlockchainImplBase implements io.grpc.BindableService {

    /**
     * <pre>
     * GetBlock retrieves information about a block based on the provided request
     * parameters.
     * </pre>
     */
    public void getBlock(pactus.blockchain.BlockchainOuterClass.GetBlockRequest request,
        io.grpc.stub.StreamObserver<pactus.blockchain.BlockchainOuterClass.GetBlockResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetBlockMethod(), responseObserver);
    }

    /**
     * <pre>
     * GetBlockHash retrieves the hash of a block at the specified height.
     * </pre>
     */
    public void getBlockHash(pactus.blockchain.BlockchainOuterClass.GetBlockHashRequest request,
        io.grpc.stub.StreamObserver<pactus.blockchain.BlockchainOuterClass.GetBlockHashResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetBlockHashMethod(), responseObserver);
    }

    /**
     * <pre>
     * GetBlockHeight retrieves the height of a block with the specified hash.
     * </pre>
     */
    public void getBlockHeight(pactus.blockchain.BlockchainOuterClass.GetBlockHeightRequest request,
        io.grpc.stub.StreamObserver<pactus.blockchain.BlockchainOuterClass.GetBlockHeightResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetBlockHeightMethod(), responseObserver);
    }

    /**
     * <pre>
     * GetBlockchainInfo retrieves general information about the blockchain.
     * </pre>
     */
    public void getBlockchainInfo(pactus.blockchain.BlockchainOuterClass.GetBlockchainInfoRequest request,
        io.grpc.stub.StreamObserver<pactus.blockchain.BlockchainOuterClass.GetBlockchainInfoResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetBlockchainInfoMethod(), responseObserver);
    }

    /**
     * <pre>
     * GetConsensusInfo retrieves information about the consensus instances.
     * </pre>
     */
    public void getConsensusInfo(pactus.blockchain.BlockchainOuterClass.GetConsensusInfoRequest request,
        io.grpc.stub.StreamObserver<pactus.blockchain.BlockchainOuterClass.GetConsensusInfoResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetConsensusInfoMethod(), responseObserver);
    }

    /**
     * <pre>
     * GetAccount retrieves information about an account based on the provided
     * address.
     * </pre>
     */
    public void getAccount(pactus.blockchain.BlockchainOuterClass.GetAccountRequest request,
        io.grpc.stub.StreamObserver<pactus.blockchain.BlockchainOuterClass.GetAccountResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetAccountMethod(), responseObserver);
    }

    /**
     * <pre>
     * GetValidator retrieves information about a validator based on the provided
     * address.
     * </pre>
     */
    public void getValidator(pactus.blockchain.BlockchainOuterClass.GetValidatorRequest request,
        io.grpc.stub.StreamObserver<pactus.blockchain.BlockchainOuterClass.GetValidatorResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetValidatorMethod(), responseObserver);
    }

    /**
     * <pre>
     * GetValidatorByNumber retrieves information about a validator based on the
     * provided number.
     * </pre>
     */
    public void getValidatorByNumber(pactus.blockchain.BlockchainOuterClass.GetValidatorByNumberRequest request,
        io.grpc.stub.StreamObserver<pactus.blockchain.BlockchainOuterClass.GetValidatorResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetValidatorByNumberMethod(), responseObserver);
    }

    /**
     * <pre>
     * GetValidatorAddresses retrieves a list of all validator addresses.
     * </pre>
     */
    public void getValidatorAddresses(pactus.blockchain.BlockchainOuterClass.GetValidatorAddressesRequest request,
        io.grpc.stub.StreamObserver<pactus.blockchain.BlockchainOuterClass.GetValidatorAddressesResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetValidatorAddressesMethod(), responseObserver);
    }

    /**
     * <pre>
     * GetPublicKey retrieves the public key of an account based on the provided
     * address.
     * </pre>
     */
    public void getPublicKey(pactus.blockchain.BlockchainOuterClass.GetPublicKeyRequest request,
        io.grpc.stub.StreamObserver<pactus.blockchain.BlockchainOuterClass.GetPublicKeyResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetPublicKeyMethod(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getGetBlockMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                pactus.blockchain.BlockchainOuterClass.GetBlockRequest,
                pactus.blockchain.BlockchainOuterClass.GetBlockResponse>(
                  this, METHODID_GET_BLOCK)))
          .addMethod(
            getGetBlockHashMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                pactus.blockchain.BlockchainOuterClass.GetBlockHashRequest,
                pactus.blockchain.BlockchainOuterClass.GetBlockHashResponse>(
                  this, METHODID_GET_BLOCK_HASH)))
          .addMethod(
            getGetBlockHeightMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                pactus.blockchain.BlockchainOuterClass.GetBlockHeightRequest,
                pactus.blockchain.BlockchainOuterClass.GetBlockHeightResponse>(
                  this, METHODID_GET_BLOCK_HEIGHT)))
          .addMethod(
            getGetBlockchainInfoMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                pactus.blockchain.BlockchainOuterClass.GetBlockchainInfoRequest,
                pactus.blockchain.BlockchainOuterClass.GetBlockchainInfoResponse>(
                  this, METHODID_GET_BLOCKCHAIN_INFO)))
          .addMethod(
            getGetConsensusInfoMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                pactus.blockchain.BlockchainOuterClass.GetConsensusInfoRequest,
                pactus.blockchain.BlockchainOuterClass.GetConsensusInfoResponse>(
                  this, METHODID_GET_CONSENSUS_INFO)))
          .addMethod(
            getGetAccountMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                pactus.blockchain.BlockchainOuterClass.GetAccountRequest,
                pactus.blockchain.BlockchainOuterClass.GetAccountResponse>(
                  this, METHODID_GET_ACCOUNT)))
          .addMethod(
            getGetValidatorMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                pactus.blockchain.BlockchainOuterClass.GetValidatorRequest,
                pactus.blockchain.BlockchainOuterClass.GetValidatorResponse>(
                  this, METHODID_GET_VALIDATOR)))
          .addMethod(
            getGetValidatorByNumberMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                pactus.blockchain.BlockchainOuterClass.GetValidatorByNumberRequest,
                pactus.blockchain.BlockchainOuterClass.GetValidatorResponse>(
                  this, METHODID_GET_VALIDATOR_BY_NUMBER)))
          .addMethod(
            getGetValidatorAddressesMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                pactus.blockchain.BlockchainOuterClass.GetValidatorAddressesRequest,
                pactus.blockchain.BlockchainOuterClass.GetValidatorAddressesResponse>(
                  this, METHODID_GET_VALIDATOR_ADDRESSES)))
          .addMethod(
            getGetPublicKeyMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                pactus.blockchain.BlockchainOuterClass.GetPublicKeyRequest,
                pactus.blockchain.BlockchainOuterClass.GetPublicKeyResponse>(
                  this, METHODID_GET_PUBLIC_KEY)))
          .build();
    }
  }

  /**
   * <pre>
   * Blockchain service defines RPC methods for interacting with the blockchain.
   * </pre>
   */
  public static final class BlockchainStub extends io.grpc.stub.AbstractAsyncStub<BlockchainStub> {
    private BlockchainStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected BlockchainStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new BlockchainStub(channel, callOptions);
    }

    /**
     * <pre>
     * GetBlock retrieves information about a block based on the provided request
     * parameters.
     * </pre>
     */
    public void getBlock(pactus.blockchain.BlockchainOuterClass.GetBlockRequest request,
        io.grpc.stub.StreamObserver<pactus.blockchain.BlockchainOuterClass.GetBlockResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetBlockMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * GetBlockHash retrieves the hash of a block at the specified height.
     * </pre>
     */
    public void getBlockHash(pactus.blockchain.BlockchainOuterClass.GetBlockHashRequest request,
        io.grpc.stub.StreamObserver<pactus.blockchain.BlockchainOuterClass.GetBlockHashResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetBlockHashMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * GetBlockHeight retrieves the height of a block with the specified hash.
     * </pre>
     */
    public void getBlockHeight(pactus.blockchain.BlockchainOuterClass.GetBlockHeightRequest request,
        io.grpc.stub.StreamObserver<pactus.blockchain.BlockchainOuterClass.GetBlockHeightResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetBlockHeightMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * GetBlockchainInfo retrieves general information about the blockchain.
     * </pre>
     */
    public void getBlockchainInfo(pactus.blockchain.BlockchainOuterClass.GetBlockchainInfoRequest request,
        io.grpc.stub.StreamObserver<pactus.blockchain.BlockchainOuterClass.GetBlockchainInfoResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetBlockchainInfoMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * GetConsensusInfo retrieves information about the consensus instances.
     * </pre>
     */
    public void getConsensusInfo(pactus.blockchain.BlockchainOuterClass.GetConsensusInfoRequest request,
        io.grpc.stub.StreamObserver<pactus.blockchain.BlockchainOuterClass.GetConsensusInfoResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetConsensusInfoMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * GetAccount retrieves information about an account based on the provided
     * address.
     * </pre>
     */
    public void getAccount(pactus.blockchain.BlockchainOuterClass.GetAccountRequest request,
        io.grpc.stub.StreamObserver<pactus.blockchain.BlockchainOuterClass.GetAccountResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetAccountMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * GetValidator retrieves information about a validator based on the provided
     * address.
     * </pre>
     */
    public void getValidator(pactus.blockchain.BlockchainOuterClass.GetValidatorRequest request,
        io.grpc.stub.StreamObserver<pactus.blockchain.BlockchainOuterClass.GetValidatorResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetValidatorMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * GetValidatorByNumber retrieves information about a validator based on the
     * provided number.
     * </pre>
     */
    public void getValidatorByNumber(pactus.blockchain.BlockchainOuterClass.GetValidatorByNumberRequest request,
        io.grpc.stub.StreamObserver<pactus.blockchain.BlockchainOuterClass.GetValidatorResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetValidatorByNumberMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * GetValidatorAddresses retrieves a list of all validator addresses.
     * </pre>
     */
    public void getValidatorAddresses(pactus.blockchain.BlockchainOuterClass.GetValidatorAddressesRequest request,
        io.grpc.stub.StreamObserver<pactus.blockchain.BlockchainOuterClass.GetValidatorAddressesResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetValidatorAddressesMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * GetPublicKey retrieves the public key of an account based on the provided
     * address.
     * </pre>
     */
    public void getPublicKey(pactus.blockchain.BlockchainOuterClass.GetPublicKeyRequest request,
        io.grpc.stub.StreamObserver<pactus.blockchain.BlockchainOuterClass.GetPublicKeyResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetPublicKeyMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   * <pre>
   * Blockchain service defines RPC methods for interacting with the blockchain.
   * </pre>
   */
  public static final class BlockchainBlockingStub extends io.grpc.stub.AbstractBlockingStub<BlockchainBlockingStub> {
    private BlockchainBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected BlockchainBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new BlockchainBlockingStub(channel, callOptions);
    }

    /**
     * <pre>
     * GetBlock retrieves information about a block based on the provided request
     * parameters.
     * </pre>
     */
    public pactus.blockchain.BlockchainOuterClass.GetBlockResponse getBlock(pactus.blockchain.BlockchainOuterClass.GetBlockRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetBlockMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetBlockHash retrieves the hash of a block at the specified height.
     * </pre>
     */
    public pactus.blockchain.BlockchainOuterClass.GetBlockHashResponse getBlockHash(pactus.blockchain.BlockchainOuterClass.GetBlockHashRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetBlockHashMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetBlockHeight retrieves the height of a block with the specified hash.
     * </pre>
     */
    public pactus.blockchain.BlockchainOuterClass.GetBlockHeightResponse getBlockHeight(pactus.blockchain.BlockchainOuterClass.GetBlockHeightRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetBlockHeightMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetBlockchainInfo retrieves general information about the blockchain.
     * </pre>
     */
    public pactus.blockchain.BlockchainOuterClass.GetBlockchainInfoResponse getBlockchainInfo(pactus.blockchain.BlockchainOuterClass.GetBlockchainInfoRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetBlockchainInfoMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetConsensusInfo retrieves information about the consensus instances.
     * </pre>
     */
    public pactus.blockchain.BlockchainOuterClass.GetConsensusInfoResponse getConsensusInfo(pactus.blockchain.BlockchainOuterClass.GetConsensusInfoRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetConsensusInfoMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetAccount retrieves information about an account based on the provided
     * address.
     * </pre>
     */
    public pactus.blockchain.BlockchainOuterClass.GetAccountResponse getAccount(pactus.blockchain.BlockchainOuterClass.GetAccountRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetAccountMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetValidator retrieves information about a validator based on the provided
     * address.
     * </pre>
     */
    public pactus.blockchain.BlockchainOuterClass.GetValidatorResponse getValidator(pactus.blockchain.BlockchainOuterClass.GetValidatorRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetValidatorMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetValidatorByNumber retrieves information about a validator based on the
     * provided number.
     * </pre>
     */
    public pactus.blockchain.BlockchainOuterClass.GetValidatorResponse getValidatorByNumber(pactus.blockchain.BlockchainOuterClass.GetValidatorByNumberRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetValidatorByNumberMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetValidatorAddresses retrieves a list of all validator addresses.
     * </pre>
     */
    public pactus.blockchain.BlockchainOuterClass.GetValidatorAddressesResponse getValidatorAddresses(pactus.blockchain.BlockchainOuterClass.GetValidatorAddressesRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetValidatorAddressesMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetPublicKey retrieves the public key of an account based on the provided
     * address.
     * </pre>
     */
    public pactus.blockchain.BlockchainOuterClass.GetPublicKeyResponse getPublicKey(pactus.blockchain.BlockchainOuterClass.GetPublicKeyRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetPublicKeyMethod(), getCallOptions(), request);
    }
  }

  /**
   * <pre>
   * Blockchain service defines RPC methods for interacting with the blockchain.
   * </pre>
   */
  public static final class BlockchainFutureStub extends io.grpc.stub.AbstractFutureStub<BlockchainFutureStub> {
    private BlockchainFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected BlockchainFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new BlockchainFutureStub(channel, callOptions);
    }

    /**
     * <pre>
     * GetBlock retrieves information about a block based on the provided request
     * parameters.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.blockchain.BlockchainOuterClass.GetBlockResponse> getBlock(
        pactus.blockchain.BlockchainOuterClass.GetBlockRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetBlockMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * GetBlockHash retrieves the hash of a block at the specified height.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.blockchain.BlockchainOuterClass.GetBlockHashResponse> getBlockHash(
        pactus.blockchain.BlockchainOuterClass.GetBlockHashRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetBlockHashMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * GetBlockHeight retrieves the height of a block with the specified hash.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.blockchain.BlockchainOuterClass.GetBlockHeightResponse> getBlockHeight(
        pactus.blockchain.BlockchainOuterClass.GetBlockHeightRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetBlockHeightMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * GetBlockchainInfo retrieves general information about the blockchain.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.blockchain.BlockchainOuterClass.GetBlockchainInfoResponse> getBlockchainInfo(
        pactus.blockchain.BlockchainOuterClass.GetBlockchainInfoRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetBlockchainInfoMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * GetConsensusInfo retrieves information about the consensus instances.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.blockchain.BlockchainOuterClass.GetConsensusInfoResponse> getConsensusInfo(
        pactus.blockchain.BlockchainOuterClass.GetConsensusInfoRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetConsensusInfoMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * GetAccount retrieves information about an account based on the provided
     * address.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.blockchain.BlockchainOuterClass.GetAccountResponse> getAccount(
        pactus.blockchain.BlockchainOuterClass.GetAccountRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetAccountMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * GetValidator retrieves information about a validator based on the provided
     * address.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.blockchain.BlockchainOuterClass.GetValidatorResponse> getValidator(
        pactus.blockchain.BlockchainOuterClass.GetValidatorRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetValidatorMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * GetValidatorByNumber retrieves information about a validator based on the
     * provided number.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.blockchain.BlockchainOuterClass.GetValidatorResponse> getValidatorByNumber(
        pactus.blockchain.BlockchainOuterClass.GetValidatorByNumberRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetValidatorByNumberMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * GetValidatorAddresses retrieves a list of all validator addresses.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.blockchain.BlockchainOuterClass.GetValidatorAddressesResponse> getValidatorAddresses(
        pactus.blockchain.BlockchainOuterClass.GetValidatorAddressesRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetValidatorAddressesMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * GetPublicKey retrieves the public key of an account based on the provided
     * address.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.blockchain.BlockchainOuterClass.GetPublicKeyResponse> getPublicKey(
        pactus.blockchain.BlockchainOuterClass.GetPublicKeyRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetPublicKeyMethod(), getCallOptions()), request);
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

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final BlockchainImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(BlockchainImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_GET_BLOCK:
          serviceImpl.getBlock((pactus.blockchain.BlockchainOuterClass.GetBlockRequest) request,
              (io.grpc.stub.StreamObserver<pactus.blockchain.BlockchainOuterClass.GetBlockResponse>) responseObserver);
          break;
        case METHODID_GET_BLOCK_HASH:
          serviceImpl.getBlockHash((pactus.blockchain.BlockchainOuterClass.GetBlockHashRequest) request,
              (io.grpc.stub.StreamObserver<pactus.blockchain.BlockchainOuterClass.GetBlockHashResponse>) responseObserver);
          break;
        case METHODID_GET_BLOCK_HEIGHT:
          serviceImpl.getBlockHeight((pactus.blockchain.BlockchainOuterClass.GetBlockHeightRequest) request,
              (io.grpc.stub.StreamObserver<pactus.blockchain.BlockchainOuterClass.GetBlockHeightResponse>) responseObserver);
          break;
        case METHODID_GET_BLOCKCHAIN_INFO:
          serviceImpl.getBlockchainInfo((pactus.blockchain.BlockchainOuterClass.GetBlockchainInfoRequest) request,
              (io.grpc.stub.StreamObserver<pactus.blockchain.BlockchainOuterClass.GetBlockchainInfoResponse>) responseObserver);
          break;
        case METHODID_GET_CONSENSUS_INFO:
          serviceImpl.getConsensusInfo((pactus.blockchain.BlockchainOuterClass.GetConsensusInfoRequest) request,
              (io.grpc.stub.StreamObserver<pactus.blockchain.BlockchainOuterClass.GetConsensusInfoResponse>) responseObserver);
          break;
        case METHODID_GET_ACCOUNT:
          serviceImpl.getAccount((pactus.blockchain.BlockchainOuterClass.GetAccountRequest) request,
              (io.grpc.stub.StreamObserver<pactus.blockchain.BlockchainOuterClass.GetAccountResponse>) responseObserver);
          break;
        case METHODID_GET_VALIDATOR:
          serviceImpl.getValidator((pactus.blockchain.BlockchainOuterClass.GetValidatorRequest) request,
              (io.grpc.stub.StreamObserver<pactus.blockchain.BlockchainOuterClass.GetValidatorResponse>) responseObserver);
          break;
        case METHODID_GET_VALIDATOR_BY_NUMBER:
          serviceImpl.getValidatorByNumber((pactus.blockchain.BlockchainOuterClass.GetValidatorByNumberRequest) request,
              (io.grpc.stub.StreamObserver<pactus.blockchain.BlockchainOuterClass.GetValidatorResponse>) responseObserver);
          break;
        case METHODID_GET_VALIDATOR_ADDRESSES:
          serviceImpl.getValidatorAddresses((pactus.blockchain.BlockchainOuterClass.GetValidatorAddressesRequest) request,
              (io.grpc.stub.StreamObserver<pactus.blockchain.BlockchainOuterClass.GetValidatorAddressesResponse>) responseObserver);
          break;
        case METHODID_GET_PUBLIC_KEY:
          serviceImpl.getPublicKey((pactus.blockchain.BlockchainOuterClass.GetPublicKeyRequest) request,
              (io.grpc.stub.StreamObserver<pactus.blockchain.BlockchainOuterClass.GetPublicKeyResponse>) responseObserver);
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

  private static abstract class BlockchainBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    BlockchainBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return pactus.blockchain.BlockchainOuterClass.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("Blockchain");
    }
  }

  private static final class BlockchainFileDescriptorSupplier
      extends BlockchainBaseDescriptorSupplier {
    BlockchainFileDescriptorSupplier() {}
  }

  private static final class BlockchainMethodDescriptorSupplier
      extends BlockchainBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    BlockchainMethodDescriptorSupplier(String methodName) {
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
      synchronized (BlockchainGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new BlockchainFileDescriptorSupplier())
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
              .build();
        }
      }
    }
    return result;
  }
}
