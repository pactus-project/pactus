package pactus.blockchain;

import static io.grpc.MethodDescriptor.generateFullMethodName;

/**
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.50.2)",
    comments = "Source: blockchain.proto")
@io.grpc.stub.annotations.GrpcGenerated
public final class BlockchainGrpc {

  private BlockchainGrpc() {}

  public static final String SERVICE_NAME = "pactus.Blockchain";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<pactus.blockchain.BlockchainOuterClass.BlockRequest,
      pactus.blockchain.BlockchainOuterClass.BlockResponse> getGetBlockMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetBlock",
      requestType = pactus.blockchain.BlockchainOuterClass.BlockRequest.class,
      responseType = pactus.blockchain.BlockchainOuterClass.BlockResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.blockchain.BlockchainOuterClass.BlockRequest,
      pactus.blockchain.BlockchainOuterClass.BlockResponse> getGetBlockMethod() {
    io.grpc.MethodDescriptor<pactus.blockchain.BlockchainOuterClass.BlockRequest, pactus.blockchain.BlockchainOuterClass.BlockResponse> getGetBlockMethod;
    if ((getGetBlockMethod = BlockchainGrpc.getGetBlockMethod) == null) {
      synchronized (BlockchainGrpc.class) {
        if ((getGetBlockMethod = BlockchainGrpc.getGetBlockMethod) == null) {
          BlockchainGrpc.getGetBlockMethod = getGetBlockMethod =
              io.grpc.MethodDescriptor.<pactus.blockchain.BlockchainOuterClass.BlockRequest, pactus.blockchain.BlockchainOuterClass.BlockResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetBlock"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.blockchain.BlockchainOuterClass.BlockRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.blockchain.BlockchainOuterClass.BlockResponse.getDefaultInstance()))
              .setSchemaDescriptor(new BlockchainMethodDescriptorSupplier("GetBlock"))
              .build();
        }
      }
    }
    return getGetBlockMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.blockchain.BlockchainOuterClass.BlockHashRequest,
      pactus.blockchain.BlockchainOuterClass.BlockHashResponse> getGetBlockHashMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetBlockHash",
      requestType = pactus.blockchain.BlockchainOuterClass.BlockHashRequest.class,
      responseType = pactus.blockchain.BlockchainOuterClass.BlockHashResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.blockchain.BlockchainOuterClass.BlockHashRequest,
      pactus.blockchain.BlockchainOuterClass.BlockHashResponse> getGetBlockHashMethod() {
    io.grpc.MethodDescriptor<pactus.blockchain.BlockchainOuterClass.BlockHashRequest, pactus.blockchain.BlockchainOuterClass.BlockHashResponse> getGetBlockHashMethod;
    if ((getGetBlockHashMethod = BlockchainGrpc.getGetBlockHashMethod) == null) {
      synchronized (BlockchainGrpc.class) {
        if ((getGetBlockHashMethod = BlockchainGrpc.getGetBlockHashMethod) == null) {
          BlockchainGrpc.getGetBlockHashMethod = getGetBlockHashMethod =
              io.grpc.MethodDescriptor.<pactus.blockchain.BlockchainOuterClass.BlockHashRequest, pactus.blockchain.BlockchainOuterClass.BlockHashResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetBlockHash"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.blockchain.BlockchainOuterClass.BlockHashRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.blockchain.BlockchainOuterClass.BlockHashResponse.getDefaultInstance()))
              .setSchemaDescriptor(new BlockchainMethodDescriptorSupplier("GetBlockHash"))
              .build();
        }
      }
    }
    return getGetBlockHashMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.blockchain.BlockchainOuterClass.BlockHeightRequest,
      pactus.blockchain.BlockchainOuterClass.BlockHeightResponse> getGetBlockHeightMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetBlockHeight",
      requestType = pactus.blockchain.BlockchainOuterClass.BlockHeightRequest.class,
      responseType = pactus.blockchain.BlockchainOuterClass.BlockHeightResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.blockchain.BlockchainOuterClass.BlockHeightRequest,
      pactus.blockchain.BlockchainOuterClass.BlockHeightResponse> getGetBlockHeightMethod() {
    io.grpc.MethodDescriptor<pactus.blockchain.BlockchainOuterClass.BlockHeightRequest, pactus.blockchain.BlockchainOuterClass.BlockHeightResponse> getGetBlockHeightMethod;
    if ((getGetBlockHeightMethod = BlockchainGrpc.getGetBlockHeightMethod) == null) {
      synchronized (BlockchainGrpc.class) {
        if ((getGetBlockHeightMethod = BlockchainGrpc.getGetBlockHeightMethod) == null) {
          BlockchainGrpc.getGetBlockHeightMethod = getGetBlockHeightMethod =
              io.grpc.MethodDescriptor.<pactus.blockchain.BlockchainOuterClass.BlockHeightRequest, pactus.blockchain.BlockchainOuterClass.BlockHeightResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetBlockHeight"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.blockchain.BlockchainOuterClass.BlockHeightRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.blockchain.BlockchainOuterClass.BlockHeightResponse.getDefaultInstance()))
              .setSchemaDescriptor(new BlockchainMethodDescriptorSupplier("GetBlockHeight"))
              .build();
        }
      }
    }
    return getGetBlockHeightMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.blockchain.BlockchainOuterClass.AccountRequest,
      pactus.blockchain.BlockchainOuterClass.AccountResponse> getGetAccountMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetAccount",
      requestType = pactus.blockchain.BlockchainOuterClass.AccountRequest.class,
      responseType = pactus.blockchain.BlockchainOuterClass.AccountResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.blockchain.BlockchainOuterClass.AccountRequest,
      pactus.blockchain.BlockchainOuterClass.AccountResponse> getGetAccountMethod() {
    io.grpc.MethodDescriptor<pactus.blockchain.BlockchainOuterClass.AccountRequest, pactus.blockchain.BlockchainOuterClass.AccountResponse> getGetAccountMethod;
    if ((getGetAccountMethod = BlockchainGrpc.getGetAccountMethod) == null) {
      synchronized (BlockchainGrpc.class) {
        if ((getGetAccountMethod = BlockchainGrpc.getGetAccountMethod) == null) {
          BlockchainGrpc.getGetAccountMethod = getGetAccountMethod =
              io.grpc.MethodDescriptor.<pactus.blockchain.BlockchainOuterClass.AccountRequest, pactus.blockchain.BlockchainOuterClass.AccountResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetAccount"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.blockchain.BlockchainOuterClass.AccountRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.blockchain.BlockchainOuterClass.AccountResponse.getDefaultInstance()))
              .setSchemaDescriptor(new BlockchainMethodDescriptorSupplier("GetAccount"))
              .build();
        }
      }
    }
    return getGetAccountMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.blockchain.BlockchainOuterClass.ValidatorsRequest,
      pactus.blockchain.BlockchainOuterClass.ValidatorsResponse> getGetValidatorsMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetValidators",
      requestType = pactus.blockchain.BlockchainOuterClass.ValidatorsRequest.class,
      responseType = pactus.blockchain.BlockchainOuterClass.ValidatorsResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.blockchain.BlockchainOuterClass.ValidatorsRequest,
      pactus.blockchain.BlockchainOuterClass.ValidatorsResponse> getGetValidatorsMethod() {
    io.grpc.MethodDescriptor<pactus.blockchain.BlockchainOuterClass.ValidatorsRequest, pactus.blockchain.BlockchainOuterClass.ValidatorsResponse> getGetValidatorsMethod;
    if ((getGetValidatorsMethod = BlockchainGrpc.getGetValidatorsMethod) == null) {
      synchronized (BlockchainGrpc.class) {
        if ((getGetValidatorsMethod = BlockchainGrpc.getGetValidatorsMethod) == null) {
          BlockchainGrpc.getGetValidatorsMethod = getGetValidatorsMethod =
              io.grpc.MethodDescriptor.<pactus.blockchain.BlockchainOuterClass.ValidatorsRequest, pactus.blockchain.BlockchainOuterClass.ValidatorsResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetValidators"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.blockchain.BlockchainOuterClass.ValidatorsRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.blockchain.BlockchainOuterClass.ValidatorsResponse.getDefaultInstance()))
              .setSchemaDescriptor(new BlockchainMethodDescriptorSupplier("GetValidators"))
              .build();
        }
      }
    }
    return getGetValidatorsMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.blockchain.BlockchainOuterClass.ValidatorRequest,
      pactus.blockchain.BlockchainOuterClass.ValidatorResponse> getGetValidatorMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetValidator",
      requestType = pactus.blockchain.BlockchainOuterClass.ValidatorRequest.class,
      responseType = pactus.blockchain.BlockchainOuterClass.ValidatorResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.blockchain.BlockchainOuterClass.ValidatorRequest,
      pactus.blockchain.BlockchainOuterClass.ValidatorResponse> getGetValidatorMethod() {
    io.grpc.MethodDescriptor<pactus.blockchain.BlockchainOuterClass.ValidatorRequest, pactus.blockchain.BlockchainOuterClass.ValidatorResponse> getGetValidatorMethod;
    if ((getGetValidatorMethod = BlockchainGrpc.getGetValidatorMethod) == null) {
      synchronized (BlockchainGrpc.class) {
        if ((getGetValidatorMethod = BlockchainGrpc.getGetValidatorMethod) == null) {
          BlockchainGrpc.getGetValidatorMethod = getGetValidatorMethod =
              io.grpc.MethodDescriptor.<pactus.blockchain.BlockchainOuterClass.ValidatorRequest, pactus.blockchain.BlockchainOuterClass.ValidatorResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetValidator"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.blockchain.BlockchainOuterClass.ValidatorRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.blockchain.BlockchainOuterClass.ValidatorResponse.getDefaultInstance()))
              .setSchemaDescriptor(new BlockchainMethodDescriptorSupplier("GetValidator"))
              .build();
        }
      }
    }
    return getGetValidatorMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.blockchain.BlockchainOuterClass.ValidatorByNumberRequest,
      pactus.blockchain.BlockchainOuterClass.ValidatorResponse> getGetValidatorByNumberMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetValidatorByNumber",
      requestType = pactus.blockchain.BlockchainOuterClass.ValidatorByNumberRequest.class,
      responseType = pactus.blockchain.BlockchainOuterClass.ValidatorResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.blockchain.BlockchainOuterClass.ValidatorByNumberRequest,
      pactus.blockchain.BlockchainOuterClass.ValidatorResponse> getGetValidatorByNumberMethod() {
    io.grpc.MethodDescriptor<pactus.blockchain.BlockchainOuterClass.ValidatorByNumberRequest, pactus.blockchain.BlockchainOuterClass.ValidatorResponse> getGetValidatorByNumberMethod;
    if ((getGetValidatorByNumberMethod = BlockchainGrpc.getGetValidatorByNumberMethod) == null) {
      synchronized (BlockchainGrpc.class) {
        if ((getGetValidatorByNumberMethod = BlockchainGrpc.getGetValidatorByNumberMethod) == null) {
          BlockchainGrpc.getGetValidatorByNumberMethod = getGetValidatorByNumberMethod =
              io.grpc.MethodDescriptor.<pactus.blockchain.BlockchainOuterClass.ValidatorByNumberRequest, pactus.blockchain.BlockchainOuterClass.ValidatorResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetValidatorByNumber"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.blockchain.BlockchainOuterClass.ValidatorByNumberRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.blockchain.BlockchainOuterClass.ValidatorResponse.getDefaultInstance()))
              .setSchemaDescriptor(new BlockchainMethodDescriptorSupplier("GetValidatorByNumber"))
              .build();
        }
      }
    }
    return getGetValidatorByNumberMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.blockchain.BlockchainOuterClass.BlockchainInfoRequest,
      pactus.blockchain.BlockchainOuterClass.BlockchainInfoResponse> getGetBlockchainInfoMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetBlockchainInfo",
      requestType = pactus.blockchain.BlockchainOuterClass.BlockchainInfoRequest.class,
      responseType = pactus.blockchain.BlockchainOuterClass.BlockchainInfoResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.blockchain.BlockchainOuterClass.BlockchainInfoRequest,
      pactus.blockchain.BlockchainOuterClass.BlockchainInfoResponse> getGetBlockchainInfoMethod() {
    io.grpc.MethodDescriptor<pactus.blockchain.BlockchainOuterClass.BlockchainInfoRequest, pactus.blockchain.BlockchainOuterClass.BlockchainInfoResponse> getGetBlockchainInfoMethod;
    if ((getGetBlockchainInfoMethod = BlockchainGrpc.getGetBlockchainInfoMethod) == null) {
      synchronized (BlockchainGrpc.class) {
        if ((getGetBlockchainInfoMethod = BlockchainGrpc.getGetBlockchainInfoMethod) == null) {
          BlockchainGrpc.getGetBlockchainInfoMethod = getGetBlockchainInfoMethod =
              io.grpc.MethodDescriptor.<pactus.blockchain.BlockchainOuterClass.BlockchainInfoRequest, pactus.blockchain.BlockchainOuterClass.BlockchainInfoResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetBlockchainInfo"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.blockchain.BlockchainOuterClass.BlockchainInfoRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.blockchain.BlockchainOuterClass.BlockchainInfoResponse.getDefaultInstance()))
              .setSchemaDescriptor(new BlockchainMethodDescriptorSupplier("GetBlockchainInfo"))
              .build();
        }
      }
    }
    return getGetBlockchainInfoMethod;
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
   */
  public static abstract class BlockchainImplBase implements io.grpc.BindableService {

    /**
     */
    public void getBlock(pactus.blockchain.BlockchainOuterClass.BlockRequest request,
        io.grpc.stub.StreamObserver<pactus.blockchain.BlockchainOuterClass.BlockResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetBlockMethod(), responseObserver);
    }

    /**
     */
    public void getBlockHash(pactus.blockchain.BlockchainOuterClass.BlockHashRequest request,
        io.grpc.stub.StreamObserver<pactus.blockchain.BlockchainOuterClass.BlockHashResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetBlockHashMethod(), responseObserver);
    }

    /**
     */
    public void getBlockHeight(pactus.blockchain.BlockchainOuterClass.BlockHeightRequest request,
        io.grpc.stub.StreamObserver<pactus.blockchain.BlockchainOuterClass.BlockHeightResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetBlockHeightMethod(), responseObserver);
    }

    /**
     */
    public void getAccount(pactus.blockchain.BlockchainOuterClass.AccountRequest request,
        io.grpc.stub.StreamObserver<pactus.blockchain.BlockchainOuterClass.AccountResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetAccountMethod(), responseObserver);
    }

    /**
     */
    public void getValidators(pactus.blockchain.BlockchainOuterClass.ValidatorsRequest request,
        io.grpc.stub.StreamObserver<pactus.blockchain.BlockchainOuterClass.ValidatorsResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetValidatorsMethod(), responseObserver);
    }

    /**
     */
    public void getValidator(pactus.blockchain.BlockchainOuterClass.ValidatorRequest request,
        io.grpc.stub.StreamObserver<pactus.blockchain.BlockchainOuterClass.ValidatorResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetValidatorMethod(), responseObserver);
    }

    /**
     */
    public void getValidatorByNumber(pactus.blockchain.BlockchainOuterClass.ValidatorByNumberRequest request,
        io.grpc.stub.StreamObserver<pactus.blockchain.BlockchainOuterClass.ValidatorResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetValidatorByNumberMethod(), responseObserver);
    }

    /**
     */
    public void getBlockchainInfo(pactus.blockchain.BlockchainOuterClass.BlockchainInfoRequest request,
        io.grpc.stub.StreamObserver<pactus.blockchain.BlockchainOuterClass.BlockchainInfoResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetBlockchainInfoMethod(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getGetBlockMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                pactus.blockchain.BlockchainOuterClass.BlockRequest,
                pactus.blockchain.BlockchainOuterClass.BlockResponse>(
                  this, METHODID_GET_BLOCK)))
          .addMethod(
            getGetBlockHashMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                pactus.blockchain.BlockchainOuterClass.BlockHashRequest,
                pactus.blockchain.BlockchainOuterClass.BlockHashResponse>(
                  this, METHODID_GET_BLOCK_HASH)))
          .addMethod(
            getGetBlockHeightMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                pactus.blockchain.BlockchainOuterClass.BlockHeightRequest,
                pactus.blockchain.BlockchainOuterClass.BlockHeightResponse>(
                  this, METHODID_GET_BLOCK_HEIGHT)))
          .addMethod(
            getGetAccountMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                pactus.blockchain.BlockchainOuterClass.AccountRequest,
                pactus.blockchain.BlockchainOuterClass.AccountResponse>(
                  this, METHODID_GET_ACCOUNT)))
          .addMethod(
            getGetValidatorsMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                pactus.blockchain.BlockchainOuterClass.ValidatorsRequest,
                pactus.blockchain.BlockchainOuterClass.ValidatorsResponse>(
                  this, METHODID_GET_VALIDATORS)))
          .addMethod(
            getGetValidatorMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                pactus.blockchain.BlockchainOuterClass.ValidatorRequest,
                pactus.blockchain.BlockchainOuterClass.ValidatorResponse>(
                  this, METHODID_GET_VALIDATOR)))
          .addMethod(
            getGetValidatorByNumberMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                pactus.blockchain.BlockchainOuterClass.ValidatorByNumberRequest,
                pactus.blockchain.BlockchainOuterClass.ValidatorResponse>(
                  this, METHODID_GET_VALIDATOR_BY_NUMBER)))
          .addMethod(
            getGetBlockchainInfoMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                pactus.blockchain.BlockchainOuterClass.BlockchainInfoRequest,
                pactus.blockchain.BlockchainOuterClass.BlockchainInfoResponse>(
                  this, METHODID_GET_BLOCKCHAIN_INFO)))
          .build();
    }
  }

  /**
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
     */
    public void getBlock(pactus.blockchain.BlockchainOuterClass.BlockRequest request,
        io.grpc.stub.StreamObserver<pactus.blockchain.BlockchainOuterClass.BlockResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetBlockMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void getBlockHash(pactus.blockchain.BlockchainOuterClass.BlockHashRequest request,
        io.grpc.stub.StreamObserver<pactus.blockchain.BlockchainOuterClass.BlockHashResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetBlockHashMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void getBlockHeight(pactus.blockchain.BlockchainOuterClass.BlockHeightRequest request,
        io.grpc.stub.StreamObserver<pactus.blockchain.BlockchainOuterClass.BlockHeightResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetBlockHeightMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void getAccount(pactus.blockchain.BlockchainOuterClass.AccountRequest request,
        io.grpc.stub.StreamObserver<pactus.blockchain.BlockchainOuterClass.AccountResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetAccountMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void getValidators(pactus.blockchain.BlockchainOuterClass.ValidatorsRequest request,
        io.grpc.stub.StreamObserver<pactus.blockchain.BlockchainOuterClass.ValidatorsResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetValidatorsMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void getValidator(pactus.blockchain.BlockchainOuterClass.ValidatorRequest request,
        io.grpc.stub.StreamObserver<pactus.blockchain.BlockchainOuterClass.ValidatorResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetValidatorMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void getValidatorByNumber(pactus.blockchain.BlockchainOuterClass.ValidatorByNumberRequest request,
        io.grpc.stub.StreamObserver<pactus.blockchain.BlockchainOuterClass.ValidatorResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetValidatorByNumberMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void getBlockchainInfo(pactus.blockchain.BlockchainOuterClass.BlockchainInfoRequest request,
        io.grpc.stub.StreamObserver<pactus.blockchain.BlockchainOuterClass.BlockchainInfoResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetBlockchainInfoMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
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
     */
    public pactus.blockchain.BlockchainOuterClass.BlockResponse getBlock(pactus.blockchain.BlockchainOuterClass.BlockRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetBlockMethod(), getCallOptions(), request);
    }

    /**
     */
    public pactus.blockchain.BlockchainOuterClass.BlockHashResponse getBlockHash(pactus.blockchain.BlockchainOuterClass.BlockHashRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetBlockHashMethod(), getCallOptions(), request);
    }

    /**
     */
    public pactus.blockchain.BlockchainOuterClass.BlockHeightResponse getBlockHeight(pactus.blockchain.BlockchainOuterClass.BlockHeightRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetBlockHeightMethod(), getCallOptions(), request);
    }

    /**
     */
    public pactus.blockchain.BlockchainOuterClass.AccountResponse getAccount(pactus.blockchain.BlockchainOuterClass.AccountRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetAccountMethod(), getCallOptions(), request);
    }

    /**
     */
    public pactus.blockchain.BlockchainOuterClass.ValidatorsResponse getValidators(pactus.blockchain.BlockchainOuterClass.ValidatorsRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetValidatorsMethod(), getCallOptions(), request);
    }

    /**
     */
    public pactus.blockchain.BlockchainOuterClass.ValidatorResponse getValidator(pactus.blockchain.BlockchainOuterClass.ValidatorRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetValidatorMethod(), getCallOptions(), request);
    }

    /**
     */
    public pactus.blockchain.BlockchainOuterClass.ValidatorResponse getValidatorByNumber(pactus.blockchain.BlockchainOuterClass.ValidatorByNumberRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetValidatorByNumberMethod(), getCallOptions(), request);
    }

    /**
     */
    public pactus.blockchain.BlockchainOuterClass.BlockchainInfoResponse getBlockchainInfo(pactus.blockchain.BlockchainOuterClass.BlockchainInfoRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetBlockchainInfoMethod(), getCallOptions(), request);
    }
  }

  /**
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
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.blockchain.BlockchainOuterClass.BlockResponse> getBlock(
        pactus.blockchain.BlockchainOuterClass.BlockRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetBlockMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.blockchain.BlockchainOuterClass.BlockHashResponse> getBlockHash(
        pactus.blockchain.BlockchainOuterClass.BlockHashRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetBlockHashMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.blockchain.BlockchainOuterClass.BlockHeightResponse> getBlockHeight(
        pactus.blockchain.BlockchainOuterClass.BlockHeightRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetBlockHeightMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.blockchain.BlockchainOuterClass.AccountResponse> getAccount(
        pactus.blockchain.BlockchainOuterClass.AccountRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetAccountMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.blockchain.BlockchainOuterClass.ValidatorsResponse> getValidators(
        pactus.blockchain.BlockchainOuterClass.ValidatorsRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetValidatorsMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.blockchain.BlockchainOuterClass.ValidatorResponse> getValidator(
        pactus.blockchain.BlockchainOuterClass.ValidatorRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetValidatorMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.blockchain.BlockchainOuterClass.ValidatorResponse> getValidatorByNumber(
        pactus.blockchain.BlockchainOuterClass.ValidatorByNumberRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetValidatorByNumberMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.blockchain.BlockchainOuterClass.BlockchainInfoResponse> getBlockchainInfo(
        pactus.blockchain.BlockchainOuterClass.BlockchainInfoRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetBlockchainInfoMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_GET_BLOCK = 0;
  private static final int METHODID_GET_BLOCK_HASH = 1;
  private static final int METHODID_GET_BLOCK_HEIGHT = 2;
  private static final int METHODID_GET_ACCOUNT = 3;
  private static final int METHODID_GET_VALIDATORS = 4;
  private static final int METHODID_GET_VALIDATOR = 5;
  private static final int METHODID_GET_VALIDATOR_BY_NUMBER = 6;
  private static final int METHODID_GET_BLOCKCHAIN_INFO = 7;

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
          serviceImpl.getBlock((pactus.blockchain.BlockchainOuterClass.BlockRequest) request,
              (io.grpc.stub.StreamObserver<pactus.blockchain.BlockchainOuterClass.BlockResponse>) responseObserver);
          break;
        case METHODID_GET_BLOCK_HASH:
          serviceImpl.getBlockHash((pactus.blockchain.BlockchainOuterClass.BlockHashRequest) request,
              (io.grpc.stub.StreamObserver<pactus.blockchain.BlockchainOuterClass.BlockHashResponse>) responseObserver);
          break;
        case METHODID_GET_BLOCK_HEIGHT:
          serviceImpl.getBlockHeight((pactus.blockchain.BlockchainOuterClass.BlockHeightRequest) request,
              (io.grpc.stub.StreamObserver<pactus.blockchain.BlockchainOuterClass.BlockHeightResponse>) responseObserver);
          break;
        case METHODID_GET_ACCOUNT:
          serviceImpl.getAccount((pactus.blockchain.BlockchainOuterClass.AccountRequest) request,
              (io.grpc.stub.StreamObserver<pactus.blockchain.BlockchainOuterClass.AccountResponse>) responseObserver);
          break;
        case METHODID_GET_VALIDATORS:
          serviceImpl.getValidators((pactus.blockchain.BlockchainOuterClass.ValidatorsRequest) request,
              (io.grpc.stub.StreamObserver<pactus.blockchain.BlockchainOuterClass.ValidatorsResponse>) responseObserver);
          break;
        case METHODID_GET_VALIDATOR:
          serviceImpl.getValidator((pactus.blockchain.BlockchainOuterClass.ValidatorRequest) request,
              (io.grpc.stub.StreamObserver<pactus.blockchain.BlockchainOuterClass.ValidatorResponse>) responseObserver);
          break;
        case METHODID_GET_VALIDATOR_BY_NUMBER:
          serviceImpl.getValidatorByNumber((pactus.blockchain.BlockchainOuterClass.ValidatorByNumberRequest) request,
              (io.grpc.stub.StreamObserver<pactus.blockchain.BlockchainOuterClass.ValidatorResponse>) responseObserver);
          break;
        case METHODID_GET_BLOCKCHAIN_INFO:
          serviceImpl.getBlockchainInfo((pactus.blockchain.BlockchainOuterClass.BlockchainInfoRequest) request,
              (io.grpc.stub.StreamObserver<pactus.blockchain.BlockchainOuterClass.BlockchainInfoResponse>) responseObserver);
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
              .addMethod(getGetAccountMethod())
              .addMethod(getGetValidatorsMethod())
              .addMethod(getGetValidatorMethod())
              .addMethod(getGetValidatorByNumberMethod())
              .addMethod(getGetBlockchainInfoMethod())
              .build();
        }
      }
    }
    return result;
  }
}
