package pactus.utils;

import static io.grpc.MethodDescriptor.generateFullMethodName;

/**
 * <pre>
 * Utils service defines RPC methods for utility functions such as message
 * signing and verification.
 * </pre>
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.50.2)",
    comments = "Source: utils.proto")
@io.grpc.stub.annotations.GrpcGenerated
public final class UtilsGrpc {

  private UtilsGrpc() {}

  public static final String SERVICE_NAME = "pactus.Utils";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<pactus.utils.UtilsOuterClass.SignMessageWithPrivateKeyRequest,
      pactus.utils.UtilsOuterClass.SignMessageWithPrivateKeyResponse> getSignMessageWithPrivateKeyMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "SignMessageWithPrivateKey",
      requestType = pactus.utils.UtilsOuterClass.SignMessageWithPrivateKeyRequest.class,
      responseType = pactus.utils.UtilsOuterClass.SignMessageWithPrivateKeyResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.utils.UtilsOuterClass.SignMessageWithPrivateKeyRequest,
      pactus.utils.UtilsOuterClass.SignMessageWithPrivateKeyResponse> getSignMessageWithPrivateKeyMethod() {
    io.grpc.MethodDescriptor<pactus.utils.UtilsOuterClass.SignMessageWithPrivateKeyRequest, pactus.utils.UtilsOuterClass.SignMessageWithPrivateKeyResponse> getSignMessageWithPrivateKeyMethod;
    if ((getSignMessageWithPrivateKeyMethod = UtilsGrpc.getSignMessageWithPrivateKeyMethod) == null) {
      synchronized (UtilsGrpc.class) {
        if ((getSignMessageWithPrivateKeyMethod = UtilsGrpc.getSignMessageWithPrivateKeyMethod) == null) {
          UtilsGrpc.getSignMessageWithPrivateKeyMethod = getSignMessageWithPrivateKeyMethod =
              io.grpc.MethodDescriptor.<pactus.utils.UtilsOuterClass.SignMessageWithPrivateKeyRequest, pactus.utils.UtilsOuterClass.SignMessageWithPrivateKeyResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "SignMessageWithPrivateKey"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.utils.UtilsOuterClass.SignMessageWithPrivateKeyRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.utils.UtilsOuterClass.SignMessageWithPrivateKeyResponse.getDefaultInstance()))
              .setSchemaDescriptor(new UtilsMethodDescriptorSupplier("SignMessageWithPrivateKey"))
              .build();
        }
      }
    }
    return getSignMessageWithPrivateKeyMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.utils.UtilsOuterClass.VerifyMessageRequest,
      pactus.utils.UtilsOuterClass.VerifyMessageResponse> getVerifyMessageMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "VerifyMessage",
      requestType = pactus.utils.UtilsOuterClass.VerifyMessageRequest.class,
      responseType = pactus.utils.UtilsOuterClass.VerifyMessageResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.utils.UtilsOuterClass.VerifyMessageRequest,
      pactus.utils.UtilsOuterClass.VerifyMessageResponse> getVerifyMessageMethod() {
    io.grpc.MethodDescriptor<pactus.utils.UtilsOuterClass.VerifyMessageRequest, pactus.utils.UtilsOuterClass.VerifyMessageResponse> getVerifyMessageMethod;
    if ((getVerifyMessageMethod = UtilsGrpc.getVerifyMessageMethod) == null) {
      synchronized (UtilsGrpc.class) {
        if ((getVerifyMessageMethod = UtilsGrpc.getVerifyMessageMethod) == null) {
          UtilsGrpc.getVerifyMessageMethod = getVerifyMessageMethod =
              io.grpc.MethodDescriptor.<pactus.utils.UtilsOuterClass.VerifyMessageRequest, pactus.utils.UtilsOuterClass.VerifyMessageResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "VerifyMessage"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.utils.UtilsOuterClass.VerifyMessageRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.utils.UtilsOuterClass.VerifyMessageResponse.getDefaultInstance()))
              .setSchemaDescriptor(new UtilsMethodDescriptorSupplier("VerifyMessage"))
              .build();
        }
      }
    }
    return getVerifyMessageMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.utils.UtilsOuterClass.BLSPublicKeyAggregateRequest,
      pactus.utils.UtilsOuterClass.BLSPublicKeyAggregateResponse> getBLSPublicKeyAggregateMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "BLSPublicKeyAggregate",
      requestType = pactus.utils.UtilsOuterClass.BLSPublicKeyAggregateRequest.class,
      responseType = pactus.utils.UtilsOuterClass.BLSPublicKeyAggregateResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.utils.UtilsOuterClass.BLSPublicKeyAggregateRequest,
      pactus.utils.UtilsOuterClass.BLSPublicKeyAggregateResponse> getBLSPublicKeyAggregateMethod() {
    io.grpc.MethodDescriptor<pactus.utils.UtilsOuterClass.BLSPublicKeyAggregateRequest, pactus.utils.UtilsOuterClass.BLSPublicKeyAggregateResponse> getBLSPublicKeyAggregateMethod;
    if ((getBLSPublicKeyAggregateMethod = UtilsGrpc.getBLSPublicKeyAggregateMethod) == null) {
      synchronized (UtilsGrpc.class) {
        if ((getBLSPublicKeyAggregateMethod = UtilsGrpc.getBLSPublicKeyAggregateMethod) == null) {
          UtilsGrpc.getBLSPublicKeyAggregateMethod = getBLSPublicKeyAggregateMethod =
              io.grpc.MethodDescriptor.<pactus.utils.UtilsOuterClass.BLSPublicKeyAggregateRequest, pactus.utils.UtilsOuterClass.BLSPublicKeyAggregateResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "BLSPublicKeyAggregate"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.utils.UtilsOuterClass.BLSPublicKeyAggregateRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.utils.UtilsOuterClass.BLSPublicKeyAggregateResponse.getDefaultInstance()))
              .setSchemaDescriptor(new UtilsMethodDescriptorSupplier("BLSPublicKeyAggregate"))
              .build();
        }
      }
    }
    return getBLSPublicKeyAggregateMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.utils.UtilsOuterClass.BLSSignatureAggregateRequest,
      pactus.utils.UtilsOuterClass.BLSSignatureAggregateResponse> getBLSSignatureAggregateMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "BLSSignatureAggregate",
      requestType = pactus.utils.UtilsOuterClass.BLSSignatureAggregateRequest.class,
      responseType = pactus.utils.UtilsOuterClass.BLSSignatureAggregateResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.utils.UtilsOuterClass.BLSSignatureAggregateRequest,
      pactus.utils.UtilsOuterClass.BLSSignatureAggregateResponse> getBLSSignatureAggregateMethod() {
    io.grpc.MethodDescriptor<pactus.utils.UtilsOuterClass.BLSSignatureAggregateRequest, pactus.utils.UtilsOuterClass.BLSSignatureAggregateResponse> getBLSSignatureAggregateMethod;
    if ((getBLSSignatureAggregateMethod = UtilsGrpc.getBLSSignatureAggregateMethod) == null) {
      synchronized (UtilsGrpc.class) {
        if ((getBLSSignatureAggregateMethod = UtilsGrpc.getBLSSignatureAggregateMethod) == null) {
          UtilsGrpc.getBLSSignatureAggregateMethod = getBLSSignatureAggregateMethod =
              io.grpc.MethodDescriptor.<pactus.utils.UtilsOuterClass.BLSSignatureAggregateRequest, pactus.utils.UtilsOuterClass.BLSSignatureAggregateResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "BLSSignatureAggregate"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.utils.UtilsOuterClass.BLSSignatureAggregateRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.utils.UtilsOuterClass.BLSSignatureAggregateResponse.getDefaultInstance()))
              .setSchemaDescriptor(new UtilsMethodDescriptorSupplier("BLSSignatureAggregate"))
              .build();
        }
      }
    }
    return getBLSSignatureAggregateMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static UtilsStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<UtilsStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<UtilsStub>() {
        @java.lang.Override
        public UtilsStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new UtilsStub(channel, callOptions);
        }
      };
    return UtilsStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static UtilsBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<UtilsBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<UtilsBlockingStub>() {
        @java.lang.Override
        public UtilsBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new UtilsBlockingStub(channel, callOptions);
        }
      };
    return UtilsBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static UtilsFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<UtilsFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<UtilsFutureStub>() {
        @java.lang.Override
        public UtilsFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new UtilsFutureStub(channel, callOptions);
        }
      };
    return UtilsFutureStub.newStub(factory, channel);
  }

  /**
   * <pre>
   * Utils service defines RPC methods for utility functions such as message
   * signing and verification.
   * </pre>
   */
  public static abstract class UtilsImplBase implements io.grpc.BindableService {

    /**
     * <pre>
     * SignMessageWithPrivateKey signs message with provided private key.
     * </pre>
     */
    public void signMessageWithPrivateKey(pactus.utils.UtilsOuterClass.SignMessageWithPrivateKeyRequest request,
        io.grpc.stub.StreamObserver<pactus.utils.UtilsOuterClass.SignMessageWithPrivateKeyResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getSignMessageWithPrivateKeyMethod(), responseObserver);
    }

    /**
     * <pre>
     * VerifyMessage verifies signature with public key and message.
     * </pre>
     */
    public void verifyMessage(pactus.utils.UtilsOuterClass.VerifyMessageRequest request,
        io.grpc.stub.StreamObserver<pactus.utils.UtilsOuterClass.VerifyMessageResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getVerifyMessageMethod(), responseObserver);
    }

    /**
     * <pre>
     * BLSPublicKeyAggregate aggregates bls public keys.
     * </pre>
     */
    public void bLSPublicKeyAggregate(pactus.utils.UtilsOuterClass.BLSPublicKeyAggregateRequest request,
        io.grpc.stub.StreamObserver<pactus.utils.UtilsOuterClass.BLSPublicKeyAggregateResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getBLSPublicKeyAggregateMethod(), responseObserver);
    }

    /**
     * <pre>
     * BLSSignatureAggregate aggregates bls signatures.
     * </pre>
     */
    public void bLSSignatureAggregate(pactus.utils.UtilsOuterClass.BLSSignatureAggregateRequest request,
        io.grpc.stub.StreamObserver<pactus.utils.UtilsOuterClass.BLSSignatureAggregateResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getBLSSignatureAggregateMethod(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getSignMessageWithPrivateKeyMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                pactus.utils.UtilsOuterClass.SignMessageWithPrivateKeyRequest,
                pactus.utils.UtilsOuterClass.SignMessageWithPrivateKeyResponse>(
                  this, METHODID_SIGN_MESSAGE_WITH_PRIVATE_KEY)))
          .addMethod(
            getVerifyMessageMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                pactus.utils.UtilsOuterClass.VerifyMessageRequest,
                pactus.utils.UtilsOuterClass.VerifyMessageResponse>(
                  this, METHODID_VERIFY_MESSAGE)))
          .addMethod(
            getBLSPublicKeyAggregateMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                pactus.utils.UtilsOuterClass.BLSPublicKeyAggregateRequest,
                pactus.utils.UtilsOuterClass.BLSPublicKeyAggregateResponse>(
                  this, METHODID_BLSPUBLIC_KEY_AGGREGATE)))
          .addMethod(
            getBLSSignatureAggregateMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                pactus.utils.UtilsOuterClass.BLSSignatureAggregateRequest,
                pactus.utils.UtilsOuterClass.BLSSignatureAggregateResponse>(
                  this, METHODID_BLSSIGNATURE_AGGREGATE)))
          .build();
    }
  }

  /**
   * <pre>
   * Utils service defines RPC methods for utility functions such as message
   * signing and verification.
   * </pre>
   */
  public static final class UtilsStub extends io.grpc.stub.AbstractAsyncStub<UtilsStub> {
    private UtilsStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected UtilsStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new UtilsStub(channel, callOptions);
    }

    /**
     * <pre>
     * SignMessageWithPrivateKey signs message with provided private key.
     * </pre>
     */
    public void signMessageWithPrivateKey(pactus.utils.UtilsOuterClass.SignMessageWithPrivateKeyRequest request,
        io.grpc.stub.StreamObserver<pactus.utils.UtilsOuterClass.SignMessageWithPrivateKeyResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getSignMessageWithPrivateKeyMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * VerifyMessage verifies signature with public key and message.
     * </pre>
     */
    public void verifyMessage(pactus.utils.UtilsOuterClass.VerifyMessageRequest request,
        io.grpc.stub.StreamObserver<pactus.utils.UtilsOuterClass.VerifyMessageResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getVerifyMessageMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * BLSPublicKeyAggregate aggregates bls public keys.
     * </pre>
     */
    public void bLSPublicKeyAggregate(pactus.utils.UtilsOuterClass.BLSPublicKeyAggregateRequest request,
        io.grpc.stub.StreamObserver<pactus.utils.UtilsOuterClass.BLSPublicKeyAggregateResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getBLSPublicKeyAggregateMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * BLSSignatureAggregate aggregates bls signatures.
     * </pre>
     */
    public void bLSSignatureAggregate(pactus.utils.UtilsOuterClass.BLSSignatureAggregateRequest request,
        io.grpc.stub.StreamObserver<pactus.utils.UtilsOuterClass.BLSSignatureAggregateResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getBLSSignatureAggregateMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   * <pre>
   * Utils service defines RPC methods for utility functions such as message
   * signing and verification.
   * </pre>
   */
  public static final class UtilsBlockingStub extends io.grpc.stub.AbstractBlockingStub<UtilsBlockingStub> {
    private UtilsBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected UtilsBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new UtilsBlockingStub(channel, callOptions);
    }

    /**
     * <pre>
     * SignMessageWithPrivateKey signs message with provided private key.
     * </pre>
     */
    public pactus.utils.UtilsOuterClass.SignMessageWithPrivateKeyResponse signMessageWithPrivateKey(pactus.utils.UtilsOuterClass.SignMessageWithPrivateKeyRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getSignMessageWithPrivateKeyMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * VerifyMessage verifies signature with public key and message.
     * </pre>
     */
    public pactus.utils.UtilsOuterClass.VerifyMessageResponse verifyMessage(pactus.utils.UtilsOuterClass.VerifyMessageRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getVerifyMessageMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * BLSPublicKeyAggregate aggregates bls public keys.
     * </pre>
     */
    public pactus.utils.UtilsOuterClass.BLSPublicKeyAggregateResponse bLSPublicKeyAggregate(pactus.utils.UtilsOuterClass.BLSPublicKeyAggregateRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getBLSPublicKeyAggregateMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * BLSSignatureAggregate aggregates bls signatures.
     * </pre>
     */
    public pactus.utils.UtilsOuterClass.BLSSignatureAggregateResponse bLSSignatureAggregate(pactus.utils.UtilsOuterClass.BLSSignatureAggregateRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getBLSSignatureAggregateMethod(), getCallOptions(), request);
    }
  }

  /**
   * <pre>
   * Utils service defines RPC methods for utility functions such as message
   * signing and verification.
   * </pre>
   */
  public static final class UtilsFutureStub extends io.grpc.stub.AbstractFutureStub<UtilsFutureStub> {
    private UtilsFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected UtilsFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new UtilsFutureStub(channel, callOptions);
    }

    /**
     * <pre>
     * SignMessageWithPrivateKey signs message with provided private key.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.utils.UtilsOuterClass.SignMessageWithPrivateKeyResponse> signMessageWithPrivateKey(
        pactus.utils.UtilsOuterClass.SignMessageWithPrivateKeyRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getSignMessageWithPrivateKeyMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * VerifyMessage verifies signature with public key and message.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.utils.UtilsOuterClass.VerifyMessageResponse> verifyMessage(
        pactus.utils.UtilsOuterClass.VerifyMessageRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getVerifyMessageMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * BLSPublicKeyAggregate aggregates bls public keys.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.utils.UtilsOuterClass.BLSPublicKeyAggregateResponse> bLSPublicKeyAggregate(
        pactus.utils.UtilsOuterClass.BLSPublicKeyAggregateRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getBLSPublicKeyAggregateMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * BLSSignatureAggregate aggregates bls signatures.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.utils.UtilsOuterClass.BLSSignatureAggregateResponse> bLSSignatureAggregate(
        pactus.utils.UtilsOuterClass.BLSSignatureAggregateRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getBLSSignatureAggregateMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_SIGN_MESSAGE_WITH_PRIVATE_KEY = 0;
  private static final int METHODID_VERIFY_MESSAGE = 1;
  private static final int METHODID_BLSPUBLIC_KEY_AGGREGATE = 2;
  private static final int METHODID_BLSSIGNATURE_AGGREGATE = 3;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final UtilsImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(UtilsImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_SIGN_MESSAGE_WITH_PRIVATE_KEY:
          serviceImpl.signMessageWithPrivateKey((pactus.utils.UtilsOuterClass.SignMessageWithPrivateKeyRequest) request,
              (io.grpc.stub.StreamObserver<pactus.utils.UtilsOuterClass.SignMessageWithPrivateKeyResponse>) responseObserver);
          break;
        case METHODID_VERIFY_MESSAGE:
          serviceImpl.verifyMessage((pactus.utils.UtilsOuterClass.VerifyMessageRequest) request,
              (io.grpc.stub.StreamObserver<pactus.utils.UtilsOuterClass.VerifyMessageResponse>) responseObserver);
          break;
        case METHODID_BLSPUBLIC_KEY_AGGREGATE:
          serviceImpl.bLSPublicKeyAggregate((pactus.utils.UtilsOuterClass.BLSPublicKeyAggregateRequest) request,
              (io.grpc.stub.StreamObserver<pactus.utils.UtilsOuterClass.BLSPublicKeyAggregateResponse>) responseObserver);
          break;
        case METHODID_BLSSIGNATURE_AGGREGATE:
          serviceImpl.bLSSignatureAggregate((pactus.utils.UtilsOuterClass.BLSSignatureAggregateRequest) request,
              (io.grpc.stub.StreamObserver<pactus.utils.UtilsOuterClass.BLSSignatureAggregateResponse>) responseObserver);
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

  private static abstract class UtilsBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    UtilsBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return pactus.utils.UtilsOuterClass.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("Utils");
    }
  }

  private static final class UtilsFileDescriptorSupplier
      extends UtilsBaseDescriptorSupplier {
    UtilsFileDescriptorSupplier() {}
  }

  private static final class UtilsMethodDescriptorSupplier
      extends UtilsBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    UtilsMethodDescriptorSupplier(String methodName) {
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
      synchronized (UtilsGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new UtilsFileDescriptorSupplier())
              .addMethod(getSignMessageWithPrivateKeyMethod())
              .addMethod(getVerifyMessageMethod())
              .addMethod(getBLSPublicKeyAggregateMethod())
              .addMethod(getBLSSignatureAggregateMethod())
              .build();
        }
      }
    }
    return result;
  }
}
