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
public final class UtilsGrpc {

  private UtilsGrpc() {}

  public static final java.lang.String SERVICE_NAME = "pactus.Utils";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<pactus.UtilsOuterClass.SignMessageWithPrivateKeyRequest,
      pactus.UtilsOuterClass.SignMessageWithPrivateKeyResponse> getSignMessageWithPrivateKeyMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "SignMessageWithPrivateKey",
      requestType = pactus.UtilsOuterClass.SignMessageWithPrivateKeyRequest.class,
      responseType = pactus.UtilsOuterClass.SignMessageWithPrivateKeyResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.UtilsOuterClass.SignMessageWithPrivateKeyRequest,
      pactus.UtilsOuterClass.SignMessageWithPrivateKeyResponse> getSignMessageWithPrivateKeyMethod() {
    io.grpc.MethodDescriptor<pactus.UtilsOuterClass.SignMessageWithPrivateKeyRequest, pactus.UtilsOuterClass.SignMessageWithPrivateKeyResponse> getSignMessageWithPrivateKeyMethod;
    if ((getSignMessageWithPrivateKeyMethod = UtilsGrpc.getSignMessageWithPrivateKeyMethod) == null) {
      synchronized (UtilsGrpc.class) {
        if ((getSignMessageWithPrivateKeyMethod = UtilsGrpc.getSignMessageWithPrivateKeyMethod) == null) {
          UtilsGrpc.getSignMessageWithPrivateKeyMethod = getSignMessageWithPrivateKeyMethod =
              io.grpc.MethodDescriptor.<pactus.UtilsOuterClass.SignMessageWithPrivateKeyRequest, pactus.UtilsOuterClass.SignMessageWithPrivateKeyResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "SignMessageWithPrivateKey"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.UtilsOuterClass.SignMessageWithPrivateKeyRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.UtilsOuterClass.SignMessageWithPrivateKeyResponse.getDefaultInstance()))
              .setSchemaDescriptor(new UtilsMethodDescriptorSupplier("SignMessageWithPrivateKey"))
              .build();
        }
      }
    }
    return getSignMessageWithPrivateKeyMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.UtilsOuterClass.VerifyMessageRequest,
      pactus.UtilsOuterClass.VerifyMessageResponse> getVerifyMessageMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "VerifyMessage",
      requestType = pactus.UtilsOuterClass.VerifyMessageRequest.class,
      responseType = pactus.UtilsOuterClass.VerifyMessageResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.UtilsOuterClass.VerifyMessageRequest,
      pactus.UtilsOuterClass.VerifyMessageResponse> getVerifyMessageMethod() {
    io.grpc.MethodDescriptor<pactus.UtilsOuterClass.VerifyMessageRequest, pactus.UtilsOuterClass.VerifyMessageResponse> getVerifyMessageMethod;
    if ((getVerifyMessageMethod = UtilsGrpc.getVerifyMessageMethod) == null) {
      synchronized (UtilsGrpc.class) {
        if ((getVerifyMessageMethod = UtilsGrpc.getVerifyMessageMethod) == null) {
          UtilsGrpc.getVerifyMessageMethod = getVerifyMessageMethod =
              io.grpc.MethodDescriptor.<pactus.UtilsOuterClass.VerifyMessageRequest, pactus.UtilsOuterClass.VerifyMessageResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "VerifyMessage"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.UtilsOuterClass.VerifyMessageRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.UtilsOuterClass.VerifyMessageResponse.getDefaultInstance()))
              .setSchemaDescriptor(new UtilsMethodDescriptorSupplier("VerifyMessage"))
              .build();
        }
      }
    }
    return getVerifyMessageMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.UtilsOuterClass.PublicKeyAggregationRequest,
      pactus.UtilsOuterClass.PublicKeyAggregationResponse> getPublicKeyAggregationMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "PublicKeyAggregation",
      requestType = pactus.UtilsOuterClass.PublicKeyAggregationRequest.class,
      responseType = pactus.UtilsOuterClass.PublicKeyAggregationResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.UtilsOuterClass.PublicKeyAggregationRequest,
      pactus.UtilsOuterClass.PublicKeyAggregationResponse> getPublicKeyAggregationMethod() {
    io.grpc.MethodDescriptor<pactus.UtilsOuterClass.PublicKeyAggregationRequest, pactus.UtilsOuterClass.PublicKeyAggregationResponse> getPublicKeyAggregationMethod;
    if ((getPublicKeyAggregationMethod = UtilsGrpc.getPublicKeyAggregationMethod) == null) {
      synchronized (UtilsGrpc.class) {
        if ((getPublicKeyAggregationMethod = UtilsGrpc.getPublicKeyAggregationMethod) == null) {
          UtilsGrpc.getPublicKeyAggregationMethod = getPublicKeyAggregationMethod =
              io.grpc.MethodDescriptor.<pactus.UtilsOuterClass.PublicKeyAggregationRequest, pactus.UtilsOuterClass.PublicKeyAggregationResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "PublicKeyAggregation"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.UtilsOuterClass.PublicKeyAggregationRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.UtilsOuterClass.PublicKeyAggregationResponse.getDefaultInstance()))
              .setSchemaDescriptor(new UtilsMethodDescriptorSupplier("PublicKeyAggregation"))
              .build();
        }
      }
    }
    return getPublicKeyAggregationMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.UtilsOuterClass.SignatureAggregationRequest,
      pactus.UtilsOuterClass.SignatureAggregationResponse> getSignatureAggregationMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "SignatureAggregation",
      requestType = pactus.UtilsOuterClass.SignatureAggregationRequest.class,
      responseType = pactus.UtilsOuterClass.SignatureAggregationResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.UtilsOuterClass.SignatureAggregationRequest,
      pactus.UtilsOuterClass.SignatureAggregationResponse> getSignatureAggregationMethod() {
    io.grpc.MethodDescriptor<pactus.UtilsOuterClass.SignatureAggregationRequest, pactus.UtilsOuterClass.SignatureAggregationResponse> getSignatureAggregationMethod;
    if ((getSignatureAggregationMethod = UtilsGrpc.getSignatureAggregationMethod) == null) {
      synchronized (UtilsGrpc.class) {
        if ((getSignatureAggregationMethod = UtilsGrpc.getSignatureAggregationMethod) == null) {
          UtilsGrpc.getSignatureAggregationMethod = getSignatureAggregationMethod =
              io.grpc.MethodDescriptor.<pactus.UtilsOuterClass.SignatureAggregationRequest, pactus.UtilsOuterClass.SignatureAggregationResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "SignatureAggregation"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.UtilsOuterClass.SignatureAggregationRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.UtilsOuterClass.SignatureAggregationResponse.getDefaultInstance()))
              .setSchemaDescriptor(new UtilsMethodDescriptorSupplier("SignatureAggregation"))
              .build();
        }
      }
    }
    return getSignatureAggregationMethod;
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
   * Creates a new blocking-style stub that supports all types of calls on the service
   */
  public static UtilsBlockingV2Stub newBlockingV2Stub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<UtilsBlockingV2Stub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<UtilsBlockingV2Stub>() {
        @java.lang.Override
        public UtilsBlockingV2Stub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new UtilsBlockingV2Stub(channel, callOptions);
        }
      };
    return UtilsBlockingV2Stub.newStub(factory, channel);
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
   * signing, verification, and etc.
   * </pre>
   */
  public interface AsyncService {

    /**
     * <pre>
     * SignMessageWithPrivateKey signs a message with the provided private key.
     * </pre>
     */
    default void signMessageWithPrivateKey(pactus.UtilsOuterClass.SignMessageWithPrivateKeyRequest request,
        io.grpc.stub.StreamObserver<pactus.UtilsOuterClass.SignMessageWithPrivateKeyResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getSignMessageWithPrivateKeyMethod(), responseObserver);
    }

    /**
     * <pre>
     * VerifyMessage verifies a signature against the public key and message.
     * </pre>
     */
    default void verifyMessage(pactus.UtilsOuterClass.VerifyMessageRequest request,
        io.grpc.stub.StreamObserver<pactus.UtilsOuterClass.VerifyMessageResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getVerifyMessageMethod(), responseObserver);
    }

    /**
     * <pre>
     * PublicKeyAggregation aggregates multiple BLS public keys into a single key.
     * </pre>
     */
    default void publicKeyAggregation(pactus.UtilsOuterClass.PublicKeyAggregationRequest request,
        io.grpc.stub.StreamObserver<pactus.UtilsOuterClass.PublicKeyAggregationResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getPublicKeyAggregationMethod(), responseObserver);
    }

    /**
     * <pre>
     * SignatureAggregation aggregates multiple BLS signatures into a single signature.
     * </pre>
     */
    default void signatureAggregation(pactus.UtilsOuterClass.SignatureAggregationRequest request,
        io.grpc.stub.StreamObserver<pactus.UtilsOuterClass.SignatureAggregationResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getSignatureAggregationMethod(), responseObserver);
    }
  }

  /**
   * Base class for the server implementation of the service Utils.
   * <pre>
   * Utils service defines RPC methods for utility functions such as message
   * signing, verification, and etc.
   * </pre>
   */
  public static abstract class UtilsImplBase
      implements io.grpc.BindableService, AsyncService {

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return UtilsGrpc.bindService(this);
    }
  }

  /**
   * A stub to allow clients to do asynchronous rpc calls to service Utils.
   * <pre>
   * Utils service defines RPC methods for utility functions such as message
   * signing, verification, and etc.
   * </pre>
   */
  public static final class UtilsStub
      extends io.grpc.stub.AbstractAsyncStub<UtilsStub> {
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
     * SignMessageWithPrivateKey signs a message with the provided private key.
     * </pre>
     */
    public void signMessageWithPrivateKey(pactus.UtilsOuterClass.SignMessageWithPrivateKeyRequest request,
        io.grpc.stub.StreamObserver<pactus.UtilsOuterClass.SignMessageWithPrivateKeyResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getSignMessageWithPrivateKeyMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * VerifyMessage verifies a signature against the public key and message.
     * </pre>
     */
    public void verifyMessage(pactus.UtilsOuterClass.VerifyMessageRequest request,
        io.grpc.stub.StreamObserver<pactus.UtilsOuterClass.VerifyMessageResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getVerifyMessageMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * PublicKeyAggregation aggregates multiple BLS public keys into a single key.
     * </pre>
     */
    public void publicKeyAggregation(pactus.UtilsOuterClass.PublicKeyAggregationRequest request,
        io.grpc.stub.StreamObserver<pactus.UtilsOuterClass.PublicKeyAggregationResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getPublicKeyAggregationMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * SignatureAggregation aggregates multiple BLS signatures into a single signature.
     * </pre>
     */
    public void signatureAggregation(pactus.UtilsOuterClass.SignatureAggregationRequest request,
        io.grpc.stub.StreamObserver<pactus.UtilsOuterClass.SignatureAggregationResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getSignatureAggregationMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   * A stub to allow clients to do synchronous rpc calls to service Utils.
   * <pre>
   * Utils service defines RPC methods for utility functions such as message
   * signing, verification, and etc.
   * </pre>
   */
  public static final class UtilsBlockingV2Stub
      extends io.grpc.stub.AbstractBlockingStub<UtilsBlockingV2Stub> {
    private UtilsBlockingV2Stub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected UtilsBlockingV2Stub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new UtilsBlockingV2Stub(channel, callOptions);
    }

    /**
     * <pre>
     * SignMessageWithPrivateKey signs a message with the provided private key.
     * </pre>
     */
    public pactus.UtilsOuterClass.SignMessageWithPrivateKeyResponse signMessageWithPrivateKey(pactus.UtilsOuterClass.SignMessageWithPrivateKeyRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getSignMessageWithPrivateKeyMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * VerifyMessage verifies a signature against the public key and message.
     * </pre>
     */
    public pactus.UtilsOuterClass.VerifyMessageResponse verifyMessage(pactus.UtilsOuterClass.VerifyMessageRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getVerifyMessageMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * PublicKeyAggregation aggregates multiple BLS public keys into a single key.
     * </pre>
     */
    public pactus.UtilsOuterClass.PublicKeyAggregationResponse publicKeyAggregation(pactus.UtilsOuterClass.PublicKeyAggregationRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getPublicKeyAggregationMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * SignatureAggregation aggregates multiple BLS signatures into a single signature.
     * </pre>
     */
    public pactus.UtilsOuterClass.SignatureAggregationResponse signatureAggregation(pactus.UtilsOuterClass.SignatureAggregationRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getSignatureAggregationMethod(), getCallOptions(), request);
    }
  }

  /**
   * A stub to allow clients to do limited synchronous rpc calls to service Utils.
   * <pre>
   * Utils service defines RPC methods for utility functions such as message
   * signing, verification, and etc.
   * </pre>
   */
  public static final class UtilsBlockingStub
      extends io.grpc.stub.AbstractBlockingStub<UtilsBlockingStub> {
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
     * SignMessageWithPrivateKey signs a message with the provided private key.
     * </pre>
     */
    public pactus.UtilsOuterClass.SignMessageWithPrivateKeyResponse signMessageWithPrivateKey(pactus.UtilsOuterClass.SignMessageWithPrivateKeyRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getSignMessageWithPrivateKeyMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * VerifyMessage verifies a signature against the public key and message.
     * </pre>
     */
    public pactus.UtilsOuterClass.VerifyMessageResponse verifyMessage(pactus.UtilsOuterClass.VerifyMessageRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getVerifyMessageMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * PublicKeyAggregation aggregates multiple BLS public keys into a single key.
     * </pre>
     */
    public pactus.UtilsOuterClass.PublicKeyAggregationResponse publicKeyAggregation(pactus.UtilsOuterClass.PublicKeyAggregationRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getPublicKeyAggregationMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * SignatureAggregation aggregates multiple BLS signatures into a single signature.
     * </pre>
     */
    public pactus.UtilsOuterClass.SignatureAggregationResponse signatureAggregation(pactus.UtilsOuterClass.SignatureAggregationRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getSignatureAggregationMethod(), getCallOptions(), request);
    }
  }

  /**
   * A stub to allow clients to do ListenableFuture-style rpc calls to service Utils.
   * <pre>
   * Utils service defines RPC methods for utility functions such as message
   * signing, verification, and etc.
   * </pre>
   */
  public static final class UtilsFutureStub
      extends io.grpc.stub.AbstractFutureStub<UtilsFutureStub> {
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
     * SignMessageWithPrivateKey signs a message with the provided private key.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.UtilsOuterClass.SignMessageWithPrivateKeyResponse> signMessageWithPrivateKey(
        pactus.UtilsOuterClass.SignMessageWithPrivateKeyRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getSignMessageWithPrivateKeyMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * VerifyMessage verifies a signature against the public key and message.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.UtilsOuterClass.VerifyMessageResponse> verifyMessage(
        pactus.UtilsOuterClass.VerifyMessageRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getVerifyMessageMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * PublicKeyAggregation aggregates multiple BLS public keys into a single key.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.UtilsOuterClass.PublicKeyAggregationResponse> publicKeyAggregation(
        pactus.UtilsOuterClass.PublicKeyAggregationRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getPublicKeyAggregationMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * SignatureAggregation aggregates multiple BLS signatures into a single signature.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.UtilsOuterClass.SignatureAggregationResponse> signatureAggregation(
        pactus.UtilsOuterClass.SignatureAggregationRequest request) {
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
          serviceImpl.signMessageWithPrivateKey((pactus.UtilsOuterClass.SignMessageWithPrivateKeyRequest) request,
              (io.grpc.stub.StreamObserver<pactus.UtilsOuterClass.SignMessageWithPrivateKeyResponse>) responseObserver);
          break;
        case METHODID_VERIFY_MESSAGE:
          serviceImpl.verifyMessage((pactus.UtilsOuterClass.VerifyMessageRequest) request,
              (io.grpc.stub.StreamObserver<pactus.UtilsOuterClass.VerifyMessageResponse>) responseObserver);
          break;
        case METHODID_PUBLIC_KEY_AGGREGATION:
          serviceImpl.publicKeyAggregation((pactus.UtilsOuterClass.PublicKeyAggregationRequest) request,
              (io.grpc.stub.StreamObserver<pactus.UtilsOuterClass.PublicKeyAggregationResponse>) responseObserver);
          break;
        case METHODID_SIGNATURE_AGGREGATION:
          serviceImpl.signatureAggregation((pactus.UtilsOuterClass.SignatureAggregationRequest) request,
              (io.grpc.stub.StreamObserver<pactus.UtilsOuterClass.SignatureAggregationResponse>) responseObserver);
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
              pactus.UtilsOuterClass.SignMessageWithPrivateKeyRequest,
              pactus.UtilsOuterClass.SignMessageWithPrivateKeyResponse>(
                service, METHODID_SIGN_MESSAGE_WITH_PRIVATE_KEY)))
        .addMethod(
          getVerifyMessageMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.UtilsOuterClass.VerifyMessageRequest,
              pactus.UtilsOuterClass.VerifyMessageResponse>(
                service, METHODID_VERIFY_MESSAGE)))
        .addMethod(
          getPublicKeyAggregationMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.UtilsOuterClass.PublicKeyAggregationRequest,
              pactus.UtilsOuterClass.PublicKeyAggregationResponse>(
                service, METHODID_PUBLIC_KEY_AGGREGATION)))
        .addMethod(
          getSignatureAggregationMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.UtilsOuterClass.SignatureAggregationRequest,
              pactus.UtilsOuterClass.SignatureAggregationResponse>(
                service, METHODID_SIGNATURE_AGGREGATION)))
        .build();
  }

  private static abstract class UtilsBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    UtilsBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return pactus.UtilsOuterClass.getDescriptor();
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
    private final java.lang.String methodName;

    UtilsMethodDescriptorSupplier(java.lang.String methodName) {
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
              .addMethod(getPublicKeyAggregationMethod())
              .addMethod(getSignatureAggregationMethod())
              .build();
        }
      }
    }
    return result;
  }
}
