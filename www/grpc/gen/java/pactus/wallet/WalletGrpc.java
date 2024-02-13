package pactus.wallet;

import static io.grpc.MethodDescriptor.generateFullMethodName;

/**
 * <pre>
 * Define the Wallet service with various RPC methods for wallet management.
 * </pre>
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.50.2)",
    comments = "Source: wallet.proto")
@io.grpc.stub.annotations.GrpcGenerated
public final class WalletGrpc {

  private WalletGrpc() {}

  public static final String SERVICE_NAME = "pactus.Wallet";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<pactus.wallet.WalletOuterClass.CreateWalletRequest,
      pactus.wallet.WalletOuterClass.CreateWalletResponse> getCreateWalletMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "CreateWallet",
      requestType = pactus.wallet.WalletOuterClass.CreateWalletRequest.class,
      responseType = pactus.wallet.WalletOuterClass.CreateWalletResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.wallet.WalletOuterClass.CreateWalletRequest,
      pactus.wallet.WalletOuterClass.CreateWalletResponse> getCreateWalletMethod() {
    io.grpc.MethodDescriptor<pactus.wallet.WalletOuterClass.CreateWalletRequest, pactus.wallet.WalletOuterClass.CreateWalletResponse> getCreateWalletMethod;
    if ((getCreateWalletMethod = WalletGrpc.getCreateWalletMethod) == null) {
      synchronized (WalletGrpc.class) {
        if ((getCreateWalletMethod = WalletGrpc.getCreateWalletMethod) == null) {
          WalletGrpc.getCreateWalletMethod = getCreateWalletMethod =
              io.grpc.MethodDescriptor.<pactus.wallet.WalletOuterClass.CreateWalletRequest, pactus.wallet.WalletOuterClass.CreateWalletResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "CreateWallet"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.wallet.WalletOuterClass.CreateWalletRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.wallet.WalletOuterClass.CreateWalletResponse.getDefaultInstance()))
              .setSchemaDescriptor(new WalletMethodDescriptorSupplier("CreateWallet"))
              .build();
        }
      }
    }
    return getCreateWalletMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.wallet.WalletOuterClass.LoadWalletRequest,
      pactus.wallet.WalletOuterClass.LoadWalletResponse> getLoadWalletMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "LoadWallet",
      requestType = pactus.wallet.WalletOuterClass.LoadWalletRequest.class,
      responseType = pactus.wallet.WalletOuterClass.LoadWalletResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.wallet.WalletOuterClass.LoadWalletRequest,
      pactus.wallet.WalletOuterClass.LoadWalletResponse> getLoadWalletMethod() {
    io.grpc.MethodDescriptor<pactus.wallet.WalletOuterClass.LoadWalletRequest, pactus.wallet.WalletOuterClass.LoadWalletResponse> getLoadWalletMethod;
    if ((getLoadWalletMethod = WalletGrpc.getLoadWalletMethod) == null) {
      synchronized (WalletGrpc.class) {
        if ((getLoadWalletMethod = WalletGrpc.getLoadWalletMethod) == null) {
          WalletGrpc.getLoadWalletMethod = getLoadWalletMethod =
              io.grpc.MethodDescriptor.<pactus.wallet.WalletOuterClass.LoadWalletRequest, pactus.wallet.WalletOuterClass.LoadWalletResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "LoadWallet"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.wallet.WalletOuterClass.LoadWalletRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.wallet.WalletOuterClass.LoadWalletResponse.getDefaultInstance()))
              .setSchemaDescriptor(new WalletMethodDescriptorSupplier("LoadWallet"))
              .build();
        }
      }
    }
    return getLoadWalletMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.wallet.WalletOuterClass.UnloadWalletRequest,
      pactus.wallet.WalletOuterClass.UnloadWalletResponse> getUnloadWalletMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "UnloadWallet",
      requestType = pactus.wallet.WalletOuterClass.UnloadWalletRequest.class,
      responseType = pactus.wallet.WalletOuterClass.UnloadWalletResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.wallet.WalletOuterClass.UnloadWalletRequest,
      pactus.wallet.WalletOuterClass.UnloadWalletResponse> getUnloadWalletMethod() {
    io.grpc.MethodDescriptor<pactus.wallet.WalletOuterClass.UnloadWalletRequest, pactus.wallet.WalletOuterClass.UnloadWalletResponse> getUnloadWalletMethod;
    if ((getUnloadWalletMethod = WalletGrpc.getUnloadWalletMethod) == null) {
      synchronized (WalletGrpc.class) {
        if ((getUnloadWalletMethod = WalletGrpc.getUnloadWalletMethod) == null) {
          WalletGrpc.getUnloadWalletMethod = getUnloadWalletMethod =
              io.grpc.MethodDescriptor.<pactus.wallet.WalletOuterClass.UnloadWalletRequest, pactus.wallet.WalletOuterClass.UnloadWalletResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "UnloadWallet"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.wallet.WalletOuterClass.UnloadWalletRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.wallet.WalletOuterClass.UnloadWalletResponse.getDefaultInstance()))
              .setSchemaDescriptor(new WalletMethodDescriptorSupplier("UnloadWallet"))
              .build();
        }
      }
    }
    return getUnloadWalletMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.wallet.WalletOuterClass.LockWalletRequest,
      pactus.wallet.WalletOuterClass.LockWalletResponse> getLockWalletMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "LockWallet",
      requestType = pactus.wallet.WalletOuterClass.LockWalletRequest.class,
      responseType = pactus.wallet.WalletOuterClass.LockWalletResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.wallet.WalletOuterClass.LockWalletRequest,
      pactus.wallet.WalletOuterClass.LockWalletResponse> getLockWalletMethod() {
    io.grpc.MethodDescriptor<pactus.wallet.WalletOuterClass.LockWalletRequest, pactus.wallet.WalletOuterClass.LockWalletResponse> getLockWalletMethod;
    if ((getLockWalletMethod = WalletGrpc.getLockWalletMethod) == null) {
      synchronized (WalletGrpc.class) {
        if ((getLockWalletMethod = WalletGrpc.getLockWalletMethod) == null) {
          WalletGrpc.getLockWalletMethod = getLockWalletMethod =
              io.grpc.MethodDescriptor.<pactus.wallet.WalletOuterClass.LockWalletRequest, pactus.wallet.WalletOuterClass.LockWalletResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "LockWallet"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.wallet.WalletOuterClass.LockWalletRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.wallet.WalletOuterClass.LockWalletResponse.getDefaultInstance()))
              .setSchemaDescriptor(new WalletMethodDescriptorSupplier("LockWallet"))
              .build();
        }
      }
    }
    return getLockWalletMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.wallet.WalletOuterClass.UnlockWalletRequest,
      pactus.wallet.WalletOuterClass.UnlockWalletResponse> getUnlockWalletMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "UnlockWallet",
      requestType = pactus.wallet.WalletOuterClass.UnlockWalletRequest.class,
      responseType = pactus.wallet.WalletOuterClass.UnlockWalletResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.wallet.WalletOuterClass.UnlockWalletRequest,
      pactus.wallet.WalletOuterClass.UnlockWalletResponse> getUnlockWalletMethod() {
    io.grpc.MethodDescriptor<pactus.wallet.WalletOuterClass.UnlockWalletRequest, pactus.wallet.WalletOuterClass.UnlockWalletResponse> getUnlockWalletMethod;
    if ((getUnlockWalletMethod = WalletGrpc.getUnlockWalletMethod) == null) {
      synchronized (WalletGrpc.class) {
        if ((getUnlockWalletMethod = WalletGrpc.getUnlockWalletMethod) == null) {
          WalletGrpc.getUnlockWalletMethod = getUnlockWalletMethod =
              io.grpc.MethodDescriptor.<pactus.wallet.WalletOuterClass.UnlockWalletRequest, pactus.wallet.WalletOuterClass.UnlockWalletResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "UnlockWallet"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.wallet.WalletOuterClass.UnlockWalletRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.wallet.WalletOuterClass.UnlockWalletResponse.getDefaultInstance()))
              .setSchemaDescriptor(new WalletMethodDescriptorSupplier("UnlockWallet"))
              .build();
        }
      }
    }
    return getUnlockWalletMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.wallet.WalletOuterClass.SignRawTransactionRequest,
      pactus.wallet.WalletOuterClass.SignRawTransactionResponse> getSignRawTransactionMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "SignRawTransaction",
      requestType = pactus.wallet.WalletOuterClass.SignRawTransactionRequest.class,
      responseType = pactus.wallet.WalletOuterClass.SignRawTransactionResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.wallet.WalletOuterClass.SignRawTransactionRequest,
      pactus.wallet.WalletOuterClass.SignRawTransactionResponse> getSignRawTransactionMethod() {
    io.grpc.MethodDescriptor<pactus.wallet.WalletOuterClass.SignRawTransactionRequest, pactus.wallet.WalletOuterClass.SignRawTransactionResponse> getSignRawTransactionMethod;
    if ((getSignRawTransactionMethod = WalletGrpc.getSignRawTransactionMethod) == null) {
      synchronized (WalletGrpc.class) {
        if ((getSignRawTransactionMethod = WalletGrpc.getSignRawTransactionMethod) == null) {
          WalletGrpc.getSignRawTransactionMethod = getSignRawTransactionMethod =
              io.grpc.MethodDescriptor.<pactus.wallet.WalletOuterClass.SignRawTransactionRequest, pactus.wallet.WalletOuterClass.SignRawTransactionResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "SignRawTransaction"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.wallet.WalletOuterClass.SignRawTransactionRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.wallet.WalletOuterClass.SignRawTransactionResponse.getDefaultInstance()))
              .setSchemaDescriptor(new WalletMethodDescriptorSupplier("SignRawTransaction"))
              .build();
        }
      }
    }
    return getSignRawTransactionMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.wallet.WalletOuterClass.GetValidatorAddressRequest,
      pactus.wallet.WalletOuterClass.GetValidatorAddressResponse> getGetValidatorAddressMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetValidatorAddress",
      requestType = pactus.wallet.WalletOuterClass.GetValidatorAddressRequest.class,
      responseType = pactus.wallet.WalletOuterClass.GetValidatorAddressResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.wallet.WalletOuterClass.GetValidatorAddressRequest,
      pactus.wallet.WalletOuterClass.GetValidatorAddressResponse> getGetValidatorAddressMethod() {
    io.grpc.MethodDescriptor<pactus.wallet.WalletOuterClass.GetValidatorAddressRequest, pactus.wallet.WalletOuterClass.GetValidatorAddressResponse> getGetValidatorAddressMethod;
    if ((getGetValidatorAddressMethod = WalletGrpc.getGetValidatorAddressMethod) == null) {
      synchronized (WalletGrpc.class) {
        if ((getGetValidatorAddressMethod = WalletGrpc.getGetValidatorAddressMethod) == null) {
          WalletGrpc.getGetValidatorAddressMethod = getGetValidatorAddressMethod =
              io.grpc.MethodDescriptor.<pactus.wallet.WalletOuterClass.GetValidatorAddressRequest, pactus.wallet.WalletOuterClass.GetValidatorAddressResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetValidatorAddress"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.wallet.WalletOuterClass.GetValidatorAddressRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.wallet.WalletOuterClass.GetValidatorAddressResponse.getDefaultInstance()))
              .setSchemaDescriptor(new WalletMethodDescriptorSupplier("GetValidatorAddress"))
              .build();
        }
      }
    }
    return getGetValidatorAddressMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static WalletStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<WalletStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<WalletStub>() {
        @java.lang.Override
        public WalletStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new WalletStub(channel, callOptions);
        }
      };
    return WalletStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static WalletBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<WalletBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<WalletBlockingStub>() {
        @java.lang.Override
        public WalletBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new WalletBlockingStub(channel, callOptions);
        }
      };
    return WalletBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static WalletFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<WalletFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<WalletFutureStub>() {
        @java.lang.Override
        public WalletFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new WalletFutureStub(channel, callOptions);
        }
      };
    return WalletFutureStub.newStub(factory, channel);
  }

  /**
   * <pre>
   * Define the Wallet service with various RPC methods for wallet management.
   * </pre>
   */
  public static abstract class WalletImplBase implements io.grpc.BindableService {

    /**
     * <pre>
     * CreateWallet creates a new wallet with the specified parameters.
     * </pre>
     */
    public void createWallet(pactus.wallet.WalletOuterClass.CreateWalletRequest request,
        io.grpc.stub.StreamObserver<pactus.wallet.WalletOuterClass.CreateWalletResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getCreateWalletMethod(), responseObserver);
    }

    /**
     * <pre>
     * LoadWallet loads an existing wallet with the given name.
     * </pre>
     */
    public void loadWallet(pactus.wallet.WalletOuterClass.LoadWalletRequest request,
        io.grpc.stub.StreamObserver<pactus.wallet.WalletOuterClass.LoadWalletResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getLoadWalletMethod(), responseObserver);
    }

    /**
     * <pre>
     * UnloadWallet unloads a currently loaded wallet with the specified name.
     * </pre>
     */
    public void unloadWallet(pactus.wallet.WalletOuterClass.UnloadWalletRequest request,
        io.grpc.stub.StreamObserver<pactus.wallet.WalletOuterClass.UnloadWalletResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getUnloadWalletMethod(), responseObserver);
    }

    /**
     * <pre>
     * LockWallet locks a currently loaded wallet with the provided password and
     * timeout.
     * </pre>
     */
    public void lockWallet(pactus.wallet.WalletOuterClass.LockWalletRequest request,
        io.grpc.stub.StreamObserver<pactus.wallet.WalletOuterClass.LockWalletResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getLockWalletMethod(), responseObserver);
    }

    /**
     * <pre>
     * UnlockWallet unlocks a locked wallet with the provided password and
     * timeout.
     * </pre>
     */
    public void unlockWallet(pactus.wallet.WalletOuterClass.UnlockWalletRequest request,
        io.grpc.stub.StreamObserver<pactus.wallet.WalletOuterClass.UnlockWalletResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getUnlockWalletMethod(), responseObserver);
    }

    /**
     * <pre>
     * SignRawTransaction signs a raw transaction for a specified wallet.
     * </pre>
     */
    public void signRawTransaction(pactus.wallet.WalletOuterClass.SignRawTransactionRequest request,
        io.grpc.stub.StreamObserver<pactus.wallet.WalletOuterClass.SignRawTransactionResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getSignRawTransactionMethod(), responseObserver);
    }

    /**
     * <pre>
     * GetValidatorAddress retrieves the validator address associated with a
     * public key.
     * </pre>
     */
    public void getValidatorAddress(pactus.wallet.WalletOuterClass.GetValidatorAddressRequest request,
        io.grpc.stub.StreamObserver<pactus.wallet.WalletOuterClass.GetValidatorAddressResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetValidatorAddressMethod(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getCreateWalletMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                pactus.wallet.WalletOuterClass.CreateWalletRequest,
                pactus.wallet.WalletOuterClass.CreateWalletResponse>(
                  this, METHODID_CREATE_WALLET)))
          .addMethod(
            getLoadWalletMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                pactus.wallet.WalletOuterClass.LoadWalletRequest,
                pactus.wallet.WalletOuterClass.LoadWalletResponse>(
                  this, METHODID_LOAD_WALLET)))
          .addMethod(
            getUnloadWalletMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                pactus.wallet.WalletOuterClass.UnloadWalletRequest,
                pactus.wallet.WalletOuterClass.UnloadWalletResponse>(
                  this, METHODID_UNLOAD_WALLET)))
          .addMethod(
            getLockWalletMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                pactus.wallet.WalletOuterClass.LockWalletRequest,
                pactus.wallet.WalletOuterClass.LockWalletResponse>(
                  this, METHODID_LOCK_WALLET)))
          .addMethod(
            getUnlockWalletMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                pactus.wallet.WalletOuterClass.UnlockWalletRequest,
                pactus.wallet.WalletOuterClass.UnlockWalletResponse>(
                  this, METHODID_UNLOCK_WALLET)))
          .addMethod(
            getSignRawTransactionMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                pactus.wallet.WalletOuterClass.SignRawTransactionRequest,
                pactus.wallet.WalletOuterClass.SignRawTransactionResponse>(
                  this, METHODID_SIGN_RAW_TRANSACTION)))
          .addMethod(
            getGetValidatorAddressMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                pactus.wallet.WalletOuterClass.GetValidatorAddressRequest,
                pactus.wallet.WalletOuterClass.GetValidatorAddressResponse>(
                  this, METHODID_GET_VALIDATOR_ADDRESS)))
          .build();
    }
  }

  /**
   * <pre>
   * Define the Wallet service with various RPC methods for wallet management.
   * </pre>
   */
  public static final class WalletStub extends io.grpc.stub.AbstractAsyncStub<WalletStub> {
    private WalletStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected WalletStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new WalletStub(channel, callOptions);
    }

    /**
     * <pre>
     * CreateWallet creates a new wallet with the specified parameters.
     * </pre>
     */
    public void createWallet(pactus.wallet.WalletOuterClass.CreateWalletRequest request,
        io.grpc.stub.StreamObserver<pactus.wallet.WalletOuterClass.CreateWalletResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getCreateWalletMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * LoadWallet loads an existing wallet with the given name.
     * </pre>
     */
    public void loadWallet(pactus.wallet.WalletOuterClass.LoadWalletRequest request,
        io.grpc.stub.StreamObserver<pactus.wallet.WalletOuterClass.LoadWalletResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getLoadWalletMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * UnloadWallet unloads a currently loaded wallet with the specified name.
     * </pre>
     */
    public void unloadWallet(pactus.wallet.WalletOuterClass.UnloadWalletRequest request,
        io.grpc.stub.StreamObserver<pactus.wallet.WalletOuterClass.UnloadWalletResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getUnloadWalletMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * LockWallet locks a currently loaded wallet with the provided password and
     * timeout.
     * </pre>
     */
    public void lockWallet(pactus.wallet.WalletOuterClass.LockWalletRequest request,
        io.grpc.stub.StreamObserver<pactus.wallet.WalletOuterClass.LockWalletResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getLockWalletMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * UnlockWallet unlocks a locked wallet with the provided password and
     * timeout.
     * </pre>
     */
    public void unlockWallet(pactus.wallet.WalletOuterClass.UnlockWalletRequest request,
        io.grpc.stub.StreamObserver<pactus.wallet.WalletOuterClass.UnlockWalletResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getUnlockWalletMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * SignRawTransaction signs a raw transaction for a specified wallet.
     * </pre>
     */
    public void signRawTransaction(pactus.wallet.WalletOuterClass.SignRawTransactionRequest request,
        io.grpc.stub.StreamObserver<pactus.wallet.WalletOuterClass.SignRawTransactionResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getSignRawTransactionMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * GetValidatorAddress retrieves the validator address associated with a
     * public key.
     * </pre>
     */
    public void getValidatorAddress(pactus.wallet.WalletOuterClass.GetValidatorAddressRequest request,
        io.grpc.stub.StreamObserver<pactus.wallet.WalletOuterClass.GetValidatorAddressResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetValidatorAddressMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   * <pre>
   * Define the Wallet service with various RPC methods for wallet management.
   * </pre>
   */
  public static final class WalletBlockingStub extends io.grpc.stub.AbstractBlockingStub<WalletBlockingStub> {
    private WalletBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected WalletBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new WalletBlockingStub(channel, callOptions);
    }

    /**
     * <pre>
     * CreateWallet creates a new wallet with the specified parameters.
     * </pre>
     */
    public pactus.wallet.WalletOuterClass.CreateWalletResponse createWallet(pactus.wallet.WalletOuterClass.CreateWalletRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getCreateWalletMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * LoadWallet loads an existing wallet with the given name.
     * </pre>
     */
    public pactus.wallet.WalletOuterClass.LoadWalletResponse loadWallet(pactus.wallet.WalletOuterClass.LoadWalletRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getLoadWalletMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * UnloadWallet unloads a currently loaded wallet with the specified name.
     * </pre>
     */
    public pactus.wallet.WalletOuterClass.UnloadWalletResponse unloadWallet(pactus.wallet.WalletOuterClass.UnloadWalletRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getUnloadWalletMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * LockWallet locks a currently loaded wallet with the provided password and
     * timeout.
     * </pre>
     */
    public pactus.wallet.WalletOuterClass.LockWalletResponse lockWallet(pactus.wallet.WalletOuterClass.LockWalletRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getLockWalletMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * UnlockWallet unlocks a locked wallet with the provided password and
     * timeout.
     * </pre>
     */
    public pactus.wallet.WalletOuterClass.UnlockWalletResponse unlockWallet(pactus.wallet.WalletOuterClass.UnlockWalletRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getUnlockWalletMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * SignRawTransaction signs a raw transaction for a specified wallet.
     * </pre>
     */
    public pactus.wallet.WalletOuterClass.SignRawTransactionResponse signRawTransaction(pactus.wallet.WalletOuterClass.SignRawTransactionRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getSignRawTransactionMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetValidatorAddress retrieves the validator address associated with a
     * public key.
     * </pre>
     */
    public pactus.wallet.WalletOuterClass.GetValidatorAddressResponse getValidatorAddress(pactus.wallet.WalletOuterClass.GetValidatorAddressRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetValidatorAddressMethod(), getCallOptions(), request);
    }
  }

  /**
   * <pre>
   * Define the Wallet service with various RPC methods for wallet management.
   * </pre>
   */
  public static final class WalletFutureStub extends io.grpc.stub.AbstractFutureStub<WalletFutureStub> {
    private WalletFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected WalletFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new WalletFutureStub(channel, callOptions);
    }

    /**
     * <pre>
     * CreateWallet creates a new wallet with the specified parameters.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.wallet.WalletOuterClass.CreateWalletResponse> createWallet(
        pactus.wallet.WalletOuterClass.CreateWalletRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getCreateWalletMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * LoadWallet loads an existing wallet with the given name.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.wallet.WalletOuterClass.LoadWalletResponse> loadWallet(
        pactus.wallet.WalletOuterClass.LoadWalletRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getLoadWalletMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * UnloadWallet unloads a currently loaded wallet with the specified name.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.wallet.WalletOuterClass.UnloadWalletResponse> unloadWallet(
        pactus.wallet.WalletOuterClass.UnloadWalletRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getUnloadWalletMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * LockWallet locks a currently loaded wallet with the provided password and
     * timeout.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.wallet.WalletOuterClass.LockWalletResponse> lockWallet(
        pactus.wallet.WalletOuterClass.LockWalletRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getLockWalletMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * UnlockWallet unlocks a locked wallet with the provided password and
     * timeout.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.wallet.WalletOuterClass.UnlockWalletResponse> unlockWallet(
        pactus.wallet.WalletOuterClass.UnlockWalletRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getUnlockWalletMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * SignRawTransaction signs a raw transaction for a specified wallet.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.wallet.WalletOuterClass.SignRawTransactionResponse> signRawTransaction(
        pactus.wallet.WalletOuterClass.SignRawTransactionRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getSignRawTransactionMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * GetValidatorAddress retrieves the validator address associated with a
     * public key.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.wallet.WalletOuterClass.GetValidatorAddressResponse> getValidatorAddress(
        pactus.wallet.WalletOuterClass.GetValidatorAddressRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetValidatorAddressMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_CREATE_WALLET = 0;
  private static final int METHODID_LOAD_WALLET = 1;
  private static final int METHODID_UNLOAD_WALLET = 2;
  private static final int METHODID_LOCK_WALLET = 3;
  private static final int METHODID_UNLOCK_WALLET = 4;
  private static final int METHODID_SIGN_RAW_TRANSACTION = 5;
  private static final int METHODID_GET_VALIDATOR_ADDRESS = 6;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final WalletImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(WalletImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_CREATE_WALLET:
          serviceImpl.createWallet((pactus.wallet.WalletOuterClass.CreateWalletRequest) request,
              (io.grpc.stub.StreamObserver<pactus.wallet.WalletOuterClass.CreateWalletResponse>) responseObserver);
          break;
        case METHODID_LOAD_WALLET:
          serviceImpl.loadWallet((pactus.wallet.WalletOuterClass.LoadWalletRequest) request,
              (io.grpc.stub.StreamObserver<pactus.wallet.WalletOuterClass.LoadWalletResponse>) responseObserver);
          break;
        case METHODID_UNLOAD_WALLET:
          serviceImpl.unloadWallet((pactus.wallet.WalletOuterClass.UnloadWalletRequest) request,
              (io.grpc.stub.StreamObserver<pactus.wallet.WalletOuterClass.UnloadWalletResponse>) responseObserver);
          break;
        case METHODID_LOCK_WALLET:
          serviceImpl.lockWallet((pactus.wallet.WalletOuterClass.LockWalletRequest) request,
              (io.grpc.stub.StreamObserver<pactus.wallet.WalletOuterClass.LockWalletResponse>) responseObserver);
          break;
        case METHODID_UNLOCK_WALLET:
          serviceImpl.unlockWallet((pactus.wallet.WalletOuterClass.UnlockWalletRequest) request,
              (io.grpc.stub.StreamObserver<pactus.wallet.WalletOuterClass.UnlockWalletResponse>) responseObserver);
          break;
        case METHODID_SIGN_RAW_TRANSACTION:
          serviceImpl.signRawTransaction((pactus.wallet.WalletOuterClass.SignRawTransactionRequest) request,
              (io.grpc.stub.StreamObserver<pactus.wallet.WalletOuterClass.SignRawTransactionResponse>) responseObserver);
          break;
        case METHODID_GET_VALIDATOR_ADDRESS:
          serviceImpl.getValidatorAddress((pactus.wallet.WalletOuterClass.GetValidatorAddressRequest) request,
              (io.grpc.stub.StreamObserver<pactus.wallet.WalletOuterClass.GetValidatorAddressResponse>) responseObserver);
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

  private static abstract class WalletBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    WalletBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return pactus.wallet.WalletOuterClass.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("Wallet");
    }
  }

  private static final class WalletFileDescriptorSupplier
      extends WalletBaseDescriptorSupplier {
    WalletFileDescriptorSupplier() {}
  }

  private static final class WalletMethodDescriptorSupplier
      extends WalletBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    WalletMethodDescriptorSupplier(String methodName) {
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
      synchronized (WalletGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new WalletFileDescriptorSupplier())
              .addMethod(getCreateWalletMethod())
              .addMethod(getLoadWalletMethod())
              .addMethod(getUnloadWalletMethod())
              .addMethod(getLockWalletMethod())
              .addMethod(getUnlockWalletMethod())
              .addMethod(getSignRawTransactionMethod())
              .addMethod(getGetValidatorAddressMethod())
              .build();
        }
      }
    }
    return result;
  }
}
