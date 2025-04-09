package pactus;

import static io.grpc.MethodDescriptor.generateFullMethodName;

/**
 * <pre>
 * Utils service defines RPC methods for utility functions such as message
 * signing, verification, and etc.
 * </pre>
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.71.0)",
    comments = "Source: utils.proto")
@io.grpc.stub.annotations.GrpcGenerated
public final class UtilsServiceGrpc {

  private UtilsServiceGrpc() {}

  public static final java.lang.String SERVICE_NAME = "pactus.UtilsService";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<pactus.Utils.SignMessageWithPrivateKeyRequest,
      pactus.Utils.SignMessageWithPrivateKeyResponse> getSignMessageWithPrivateKeyMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "SignMessageWithPrivateKey",
      requestType = pactus.Utils.SignMessageWithPrivateKeyRequest.class,
      responseType = pactus.Utils.SignMessageWithPrivateKeyResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.Utils.SignMessageWithPrivateKeyRequest,
      pactus.Utils.SignMessageWithPrivateKeyResponse> getSignMessageWithPrivateKeyMethod() {
    io.grpc.MethodDescriptor<pactus.Utils.SignMessageWithPrivateKeyRequest, pactus.Utils.SignMessageWithPrivateKeyResponse> getSignMessageWithPrivateKeyMethod;
    if ((getSignMessageWithPrivateKeyMethod = UtilsServiceGrpc.getSignMessageWithPrivateKeyMethod) == null) {
      synchronized (UtilsServiceGrpc.class) {
        if ((getSignMessageWithPrivateKeyMethod = UtilsServiceGrpc.getSignMessageWithPrivateKeyMethod) == null) {
          UtilsServiceGrpc.getSignMessageWithPrivateKeyMethod = getSignMessageWithPrivateKeyMethod =
              io.grpc.MethodDescriptor.<pactus.Utils.SignMessageWithPrivateKeyRequest, pactus.Utils.SignMessageWithPrivateKeyResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "SignMessageWithPrivateKey"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Utils.SignMessageWithPrivateKeyRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Utils.SignMessageWithPrivateKeyResponse.getDefaultInstance()))
              .setSchemaDescriptor(new UtilsServiceMethodDescriptorSupplier("SignMessageWithPrivateKey"))
              .build();
        }
      }
    }
    return getSignMessageWithPrivateKeyMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.Utils.VerifyMessageRequest,
      pactus.Utils.VerifyMessageResponse> getVerifyMessageMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "VerifyMessage",
      requestType = pactus.Utils.VerifyMessageRequest.class,
      responseType = pactus.Utils.VerifyMessageResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.Utils.VerifyMessageRequest,
      pactus.Utils.VerifyMessageResponse> getVerifyMessageMethod() {
    io.grpc.MethodDescriptor<pactus.Utils.VerifyMessageRequest, pactus.Utils.VerifyMessageResponse> getVerifyMessageMethod;
    if ((getVerifyMessageMethod = UtilsServiceGrpc.getVerifyMessageMethod) == null) {
      synchronized (UtilsServiceGrpc.class) {
        if ((getVerifyMessageMethod = UtilsServiceGrpc.getVerifyMessageMethod) == null) {
          UtilsServiceGrpc.getVerifyMessageMethod = getVerifyMessageMethod =
              io.grpc.MethodDescriptor.<pactus.Utils.VerifyMessageRequest, pactus.Utils.VerifyMessageResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "VerifyMessage"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Utils.VerifyMessageRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Utils.VerifyMessageResponse.getDefaultInstance()))
              .setSchemaDescriptor(new UtilsServiceMethodDescriptorSupplier("VerifyMessage"))
              .build();
        }
      }
    }
    return getVerifyMessageMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.Utils.PublicKeyAggregationRequest,
      pactus.Utils.PublicKeyAggregationResponse> getPublicKeyAggregationMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "PublicKeyAggregation",
      requestType = pactus.Utils.PublicKeyAggregationRequest.class,
      responseType = pactus.Utils.PublicKeyAggregationResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.Utils.PublicKeyAggregationRequest,
      pactus.Utils.PublicKeyAggregationResponse> getPublicKeyAggregationMethod() {
    io.grpc.MethodDescriptor<pactus.Utils.PublicKeyAggregationRequest, pactus.Utils.PublicKeyAggregationResponse> getPublicKeyAggregationMethod;
    if ((getPublicKeyAggregationMethod = UtilsServiceGrpc.getPublicKeyAggregationMethod) == null) {
      synchronized (UtilsServiceGrpc.class) {
        if ((getPublicKeyAggregationMethod = UtilsServiceGrpc.getPublicKeyAggregationMethod) == null) {
          UtilsServiceGrpc.getPublicKeyAggregationMethod = getPublicKeyAggregationMethod =
              io.grpc.MethodDescriptor.<pactus.Utils.PublicKeyAggregationRequest, pactus.Utils.PublicKeyAggregationResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "PublicKeyAggregation"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Utils.PublicKeyAggregationRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Utils.PublicKeyAggregationResponse.getDefaultInstance()))
              .setSchemaDescriptor(new UtilsServiceMethodDescriptorSupplier("PublicKeyAggregation"))
              .build();
        }
      }
    }
    return getPublicKeyAggregationMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.Utils.SignatureAggregationRequest,
      pactus.Utils.SignatureAggregationResponse> getSignatureAggregationMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "SignatureAggregation",
      requestType = pactus.Utils.SignatureAggregationRequest.class,
      responseType = pactus.Utils.SignatureAggregationResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.Utils.SignatureAggregationRequest,
      pactus.Utils.SignatureAggregationResponse> getSignatureAggregationMethod() {
    io.grpc.MethodDescriptor<pactus.Utils.SignatureAggregationRequest, pactus.Utils.SignatureAggregationResponse> getSignatureAggregationMethod;
    if ((getSignatureAggregationMethod = UtilsServiceGrpc.getSignatureAggregationMethod) == null) {
      synchronized (UtilsServiceGrpc.class) {
        if ((getSignatureAggregationMethod = UtilsServiceGrpc.getSignatureAggregationMethod) == null) {
          UtilsServiceGrpc.getSignatureAggregationMethod = getSignatureAggregationMethod =
              io.grpc.MethodDescriptor.<pactus.Utils.SignatureAggregationRequest, pactus.Utils.SignatureAggregationResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "SignatureAggregation"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Utils.SignatureAggregationRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Utils.SignatureAggregationResponse.getDefaultInstance()))
              .setSchemaDescriptor(new UtilsServiceMethodDescriptorSupplier("SignatureAggregation"))
              .build();
        }
      }
    }
    return getSignatureAggregationMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static UtilsServiceStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<UtilsServiceStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<UtilsServiceStub>() {
        @java.lang.Override
        public UtilsServiceStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new UtilsServiceStub(channel, callOptions);
        }
      };
    return UtilsServiceStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports all types of calls on the service
   */
  public static UtilsServiceBlockingV2Stub newBlockingV2Stub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<UtilsServiceBlockingV2Stub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<UtilsServiceBlockingV2Stub>() {
        @java.lang.Override
        public UtilsServiceBlockingV2Stub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new UtilsServiceBlockingV2Stub(channel, callOptions);
        }
      };
    return UtilsServiceBlockingV2Stub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static UtilsServiceBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<UtilsServiceBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<UtilsServiceBlockingStub>() {
        @java.lang.Override
        public UtilsServiceBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new UtilsServiceBlockingStub(channel, callOptions);
        }
      };
    return UtilsServiceBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static UtilsServiceFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<UtilsServiceFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<UtilsServiceFutureStub>() {
        @java.lang.Override
        public UtilsServiceFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new UtilsServiceFutureStub(channel, callOptions);
        }
      };
    return UtilsServiceFutureStub.newStub(factory, channel);
  }

  /**
   * <pre>
   * Utils service defines RPC methods for utility functions such as message
   * signing, verification, and etc.
   * </pre>
   */
  public interface AsyncService {

    /**
     * <pre>
     * SignMessageWithPrivateKey signs a message with the provided private key.
     * </pre>
     */
    default void signMessageWithPrivateKey(pactus.Utils.SignMessageWithPrivateKeyRequest request,
        io.grpc.stub.StreamObserver<pactus.Utils.SignMessageWithPrivateKeyResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getSignMessageWithPrivateKeyMethod(), responseObserver);
    }

    /**
     * <pre>
     * VerifyMessage verifies a signature against the public key and message.
     * </pre>
     */
    default void verifyMessage(pactus.Utils.VerifyMessageRequest request,
        io.grpc.stub.StreamObserver<pactus.Utils.VerifyMessageResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getVerifyMessageMethod(), responseObserver);
    }

    /**
     * <pre>
     * PublicKeyAggregation aggregates multiple BLS public keys into a single key.
     * </pre>
     */
    default void publicKeyAggregation(pactus.Utils.PublicKeyAggregationRequest request,
        io.grpc.stub.StreamObserver<pactus.Utils.PublicKeyAggregationResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getPublicKeyAggregationMethod(), responseObserver);
    }

    /**
     * <pre>
     * SignatureAggregation aggregates multiple BLS signatures into a single signature.
     * </pre>
     */
    default void signatureAggregation(pactus.Utils.SignatureAggregationRequest request,
        io.grpc.stub.StreamObserver<pactus.Utils.SignatureAggregationResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getSignatureAggregationMethod(), responseObserver);
    }
  }

  /**
   * Base class for the server implementation of the service UtilsService.
   * <pre>
   * Utils service defines RPC methods for utility functions such as message
   * signing, verification, and etc.
   * </pre>
   */
  public static abstract class UtilsServiceImplBase
      implements io.grpc.BindableService, AsyncService {

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return UtilsServiceGrpc.bindService(this);
    }
  }

  /**
   * A stub to allow clients to do asynchronous rpc calls to service UtilsService.
   * <pre>
   * Utils service defines RPC methods for utility functions such as message
   * signing, verification, and etc.
   * </pre>
   */
  public static final class UtilsServiceStub
      extends io.grpc.stub.AbstractAsyncStub<UtilsServiceStub> {
    private UtilsServiceStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected UtilsServiceStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new UtilsServiceStub(channel, callOptions);
    }

    /**
     * <pre>
     * SignMessageWithPrivateKey signs a message with the provided private key.
     * </pre>
     */
    public void signMessageWithPrivateKey(pactus.Utils.SignMessageWithPrivateKeyRequest request,
        io.grpc.stub.StreamObserver<pactus.Utils.SignMessageWithPrivateKeyResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getSignMessageWithPrivateKeyMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * VerifyMessage verifies a signature against the public key and message.
     * </pre>
     */
    public void verifyMessage(pactus.Utils.VerifyMessageRequest request,
        io.grpc.stub.StreamObserver<pactus.Utils.VerifyMessageResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getVerifyMessageMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * PublicKeyAggregation aggregates multiple BLS public keys into a single key.
     * </pre>
     */
    public void publicKeyAggregation(pactus.Utils.PublicKeyAggregationRequest request,
        io.grpc.stub.StreamObserver<pactus.Utils.PublicKeyAggregationResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getPublicKeyAggregationMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * SignatureAggregation aggregates multiple BLS signatures into a single signature.
     * </pre>
     */
    public void signatureAggregation(pactus.Utils.SignatureAggregationRequest request,
        io.grpc.stub.StreamObserver<pactus.Utils.SignatureAggregationResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getSignatureAggregationMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   * A stub to allow clients to do synchronous rpc calls to service UtilsService.
   * <pre>
   * Utils service defines RPC methods for utility functions such as message
   * signing, verification, and etc.
   * </pre>
   */
  public static final class UtilsServiceBlockingV2Stub
      extends io.grpc.stub.AbstractBlockingStub<UtilsServiceBlockingV2Stub> {
    private UtilsServiceBlockingV2Stub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected UtilsServiceBlockingV2Stub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new UtilsServiceBlockingV2Stub(channel, callOptions);
    }

    /**
     * <pre>
     * SignMessageWithPrivateKey signs a message with the provided private key.
     * </pre>
     */
    public pactus.Utils.SignMessageWithPrivateKeyResponse signMessageWithPrivateKey(pactus.Utils.SignMessageWithPrivateKeyRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getSignMessageWithPrivateKeyMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * VerifyMessage verifies a signature against the public key and message.
     * </pre>
     */
    public pactus.Utils.VerifyMessageResponse verifyMessage(pactus.Utils.VerifyMessageRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getVerifyMessageMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * PublicKeyAggregation aggregates multiple BLS public keys into a single key.
     * </pre>
     */
    public pactus.Utils.PublicKeyAggregationResponse publicKeyAggregation(pactus.Utils.PublicKeyAggregationRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getPublicKeyAggregationMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * SignatureAggregation aggregates multiple BLS signatures into a single signature.
     * </pre>
     */
    public pactus.Utils.SignatureAggregationResponse signatureAggregation(pactus.Utils.SignatureAggregationRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getSignatureAggregationMethod(), getCallOptions(), request);
    }
  }

  /**
   * A stub to allow clients to do limited synchronous rpc calls to service UtilsService.
   * <pre>
   * Utils service defines RPC methods for utility functions such as message
   * signing, verification, and etc.
   * </pre>
   */
  public static final class UtilsServiceBlockingStub
      extends io.grpc.stub.AbstractBlockingStub<UtilsServiceBlockingStub> {
    private UtilsServiceBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected UtilsServiceBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new UtilsServiceBlockingStub(channel, callOptions);
    }

    /**
     * <pre>
     * SignMessageWithPrivateKey signs a message with the provided private key.
     * </pre>
     */
    public pactus.Utils.SignMessageWithPrivateKeyResponse signMessageWithPrivateKey(pactus.Utils.SignMessageWithPrivateKeyRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getSignMessageWithPrivateKeyMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * VerifyMessage verifies a signature against the public key and message.
     * </pre>
     */
    public pactus.Utils.VerifyMessageResponse verifyMessage(pactus.Utils.VerifyMessageRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getVerifyMessageMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * PublicKeyAggregation aggregates multiple BLS public keys into a single key.
     * </pre>
     */
    public pactus.Utils.PublicKeyAggregationResponse publicKeyAggregation(pactus.Utils.PublicKeyAggregationRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getPublicKeyAggregationMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * SignatureAggregation aggregates multiple BLS signatures into a single signature.
     * </pre>
     */
    public pactus.Utils.SignatureAggregationResponse signatureAggregation(pactus.Utils.SignatureAggregationRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getSignatureAggregationMethod(), getCallOptions(), request);
    }
  }

  /**
   * A stub to allow clients to do ListenableFuture-style rpc calls to service UtilsService.
   * <pre>
   * Utils service defines RPC methods for utility functions such as message
   * signing, verification, and etc.
   * </pre>
   */
  public static final class UtilsServiceFutureStub
      extends io.grpc.stub.AbstractFutureStub<UtilsServiceFutureStub> {
    private UtilsServiceFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected UtilsServiceFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new UtilsServiceFutureStub(channel, callOptions);
    }

    /**
     * <pre>
     * SignMessageWithPrivateKey signs a message with the provided private key.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.Utils.SignMessageWithPrivateKeyResponse> signMessageWithPrivateKey(
        pactus.Utils.SignMessageWithPrivateKeyRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getSignMessageWithPrivateKeyMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * VerifyMessage verifies a signature against the public key and message.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.Utils.VerifyMessageResponse> verifyMessage(
        pactus.Utils.VerifyMessageRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getVerifyMessageMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * PublicKeyAggregation aggregates multiple BLS public keys into a single key.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.Utils.PublicKeyAggregationResponse> publicKeyAggregation(
        pactus.Utils.PublicKeyAggregationRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getPublicKeyAggregationMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * SignatureAggregation aggregates multiple BLS signatures into a single signature.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.Utils.SignatureAggregationResponse> signatureAggregation(
        pactus.Utils.SignatureAggregationRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getSignatureAggregationMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_SIGN_MESSAGE_WITH_PRIVATE_KEY = 0;
  private static final int METHODID_VERIFY_MESSAGE = 1;
  private static final int METHODID_PUBLIC_KEY_AGGREGATION = 2;
  private static final int METHODID_SIGNATURE_AGGREGATION = 3;

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
        case METHODID_SIGN_MESSAGE_WITH_PRIVATE_KEY:
          serviceImpl.signMessageWithPrivateKey((pactus.Utils.SignMessageWithPrivateKeyRequest) request,
              (io.grpc.stub.StreamObserver<pactus.Utils.SignMessageWithPrivateKeyResponse>) responseObserver);
          break;
        case METHODID_VERIFY_MESSAGE:
          serviceImpl.verifyMessage((pactus.Utils.VerifyMessageRequest) request,
              (io.grpc.stub.StreamObserver<pactus.Utils.VerifyMessageResponse>) responseObserver);
          break;
        case METHODID_PUBLIC_KEY_AGGREGATION:
          serviceImpl.publicKeyAggregation((pactus.Utils.PublicKeyAggregationRequest) request,
              (io.grpc.stub.StreamObserver<pactus.Utils.PublicKeyAggregationResponse>) responseObserver);
          break;
        case METHODID_SIGNATURE_AGGREGATION:
          serviceImpl.signatureAggregation((pactus.Utils.SignatureAggregationRequest) request,
              (io.grpc.stub.StreamObserver<pactus.Utils.SignatureAggregationResponse>) responseObserver);
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
          getSignMessageWithPrivateKeyMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.Utils.SignMessageWithPrivateKeyRequest,
              pactus.Utils.SignMessageWithPrivateKeyResponse>(
                service, METHODID_SIGN_MESSAGE_WITH_PRIVATE_KEY)))
        .addMethod(
          getVerifyMessageMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.Utils.VerifyMessageRequest,
              pactus.Utils.VerifyMessageResponse>(
                service, METHODID_VERIFY_MESSAGE)))
        .addMethod(
          getPublicKeyAggregationMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.Utils.PublicKeyAggregationRequest,
              pactus.Utils.PublicKeyAggregationResponse>(
                service, METHODID_PUBLIC_KEY_AGGREGATION)))
        .addMethod(
          getSignatureAggregationMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.Utils.SignatureAggregationRequest,
              pactus.Utils.SignatureAggregationResponse>(
                service, METHODID_SIGNATURE_AGGREGATION)))
        .build();
  }

  private static abstract class UtilsServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    UtilsServiceBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return pactus.Utils.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("UtilsService");
    }
  }

  private static final class UtilsServiceFileDescriptorSupplier
      extends UtilsServiceBaseDescriptorSupplier {
    UtilsServiceFileDescriptorSupplier() {}
  }

  private static final class UtilsServiceMethodDescriptorSupplier
      extends UtilsServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final java.lang.String methodName;

    UtilsServiceMethodDescriptorSupplier(java.lang.String methodName) {
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
      synchronized (UtilsServiceGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new UtilsServiceFileDescriptorSupplier())
              .addMethod(getSignMessageWithPrivateKeyMethod())
              .addMethod(getVerifyMessageMethod())
              .addMethod(getPublicKeyAggregationMethod())
              .addMethod(getSignatureAggregationMethod())
              .build();
        }
      }
    }
    return result;
  }
}
