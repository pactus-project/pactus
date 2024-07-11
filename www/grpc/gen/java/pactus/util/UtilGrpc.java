package pactus.util;

import static io.grpc.MethodDescriptor.generateFullMethodName;

/**
 * <pre>
 * Util service defines various RPC methods for interacting with
 * Utils.
 * </pre>
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.50.2)",
    comments = "Source: util.proto")
@io.grpc.stub.annotations.GrpcGenerated
public final class UtilGrpc {

  private UtilGrpc() {}

  public static final String SERVICE_NAME = "pactus.Util";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<pactus.util.UtilOuterClass.SignMessageWithPrivateKeyRequest,
      pactus.util.UtilOuterClass.SignMessageWithPrivateKeyResponse> getSignMessageWithPrivateKeyMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "SignMessageWithPrivateKey",
      requestType = pactus.util.UtilOuterClass.SignMessageWithPrivateKeyRequest.class,
      responseType = pactus.util.UtilOuterClass.SignMessageWithPrivateKeyResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.util.UtilOuterClass.SignMessageWithPrivateKeyRequest,
      pactus.util.UtilOuterClass.SignMessageWithPrivateKeyResponse> getSignMessageWithPrivateKeyMethod() {
    io.grpc.MethodDescriptor<pactus.util.UtilOuterClass.SignMessageWithPrivateKeyRequest, pactus.util.UtilOuterClass.SignMessageWithPrivateKeyResponse> getSignMessageWithPrivateKeyMethod;
    if ((getSignMessageWithPrivateKeyMethod = UtilGrpc.getSignMessageWithPrivateKeyMethod) == null) {
      synchronized (UtilGrpc.class) {
        if ((getSignMessageWithPrivateKeyMethod = UtilGrpc.getSignMessageWithPrivateKeyMethod) == null) {
          UtilGrpc.getSignMessageWithPrivateKeyMethod = getSignMessageWithPrivateKeyMethod =
              io.grpc.MethodDescriptor.<pactus.util.UtilOuterClass.SignMessageWithPrivateKeyRequest, pactus.util.UtilOuterClass.SignMessageWithPrivateKeyResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "SignMessageWithPrivateKey"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.util.UtilOuterClass.SignMessageWithPrivateKeyRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.util.UtilOuterClass.SignMessageWithPrivateKeyResponse.getDefaultInstance()))
              .setSchemaDescriptor(new UtilMethodDescriptorSupplier("SignMessageWithPrivateKey"))
              .build();
        }
      }
    }
    return getSignMessageWithPrivateKeyMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.util.UtilOuterClass.VerifyMessageRequest,
      pactus.util.UtilOuterClass.VerifyMessageResponse> getVerifyMessageMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "VerifyMessage",
      requestType = pactus.util.UtilOuterClass.VerifyMessageRequest.class,
      responseType = pactus.util.UtilOuterClass.VerifyMessageResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.util.UtilOuterClass.VerifyMessageRequest,
      pactus.util.UtilOuterClass.VerifyMessageResponse> getVerifyMessageMethod() {
    io.grpc.MethodDescriptor<pactus.util.UtilOuterClass.VerifyMessageRequest, pactus.util.UtilOuterClass.VerifyMessageResponse> getVerifyMessageMethod;
    if ((getVerifyMessageMethod = UtilGrpc.getVerifyMessageMethod) == null) {
      synchronized (UtilGrpc.class) {
        if ((getVerifyMessageMethod = UtilGrpc.getVerifyMessageMethod) == null) {
          UtilGrpc.getVerifyMessageMethod = getVerifyMessageMethod =
              io.grpc.MethodDescriptor.<pactus.util.UtilOuterClass.VerifyMessageRequest, pactus.util.UtilOuterClass.VerifyMessageResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "VerifyMessage"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.util.UtilOuterClass.VerifyMessageRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.util.UtilOuterClass.VerifyMessageResponse.getDefaultInstance()))
              .setSchemaDescriptor(new UtilMethodDescriptorSupplier("VerifyMessage"))
              .build();
        }
      }
    }
    return getVerifyMessageMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static UtilStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<UtilStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<UtilStub>() {
        @java.lang.Override
        public UtilStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new UtilStub(channel, callOptions);
        }
      };
    return UtilStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static UtilBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<UtilBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<UtilBlockingStub>() {
        @java.lang.Override
        public UtilBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new UtilBlockingStub(channel, callOptions);
        }
      };
    return UtilBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static UtilFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<UtilFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<UtilFutureStub>() {
        @java.lang.Override
        public UtilFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new UtilFutureStub(channel, callOptions);
        }
      };
    return UtilFutureStub.newStub(factory, channel);
  }

  /**
   * <pre>
   * Util service defines various RPC methods for interacting with
   * Utils.
   * </pre>
   */
  public static abstract class UtilImplBase implements io.grpc.BindableService {

    /**
     * <pre>
     * SignMessageWithPrivateKey
     * </pre>
     */
    public void signMessageWithPrivateKey(pactus.util.UtilOuterClass.SignMessageWithPrivateKeyRequest request,
        io.grpc.stub.StreamObserver<pactus.util.UtilOuterClass.SignMessageWithPrivateKeyResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getSignMessageWithPrivateKeyMethod(), responseObserver);
    }

    /**
     * <pre>
     * VerifyMessage
     * </pre>
     */
    public void verifyMessage(pactus.util.UtilOuterClass.VerifyMessageRequest request,
        io.grpc.stub.StreamObserver<pactus.util.UtilOuterClass.VerifyMessageResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getVerifyMessageMethod(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getSignMessageWithPrivateKeyMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                pactus.util.UtilOuterClass.SignMessageWithPrivateKeyRequest,
                pactus.util.UtilOuterClass.SignMessageWithPrivateKeyResponse>(
                  this, METHODID_SIGN_MESSAGE_WITH_PRIVATE_KEY)))
          .addMethod(
            getVerifyMessageMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                pactus.util.UtilOuterClass.VerifyMessageRequest,
                pactus.util.UtilOuterClass.VerifyMessageResponse>(
                  this, METHODID_VERIFY_MESSAGE)))
          .build();
    }
  }

  /**
   * <pre>
   * Util service defines various RPC methods for interacting with
   * Utils.
   * </pre>
   */
  public static final class UtilStub extends io.grpc.stub.AbstractAsyncStub<UtilStub> {
    private UtilStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected UtilStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new UtilStub(channel, callOptions);
    }

    /**
     * <pre>
     * SignMessageWithPrivateKey
     * </pre>
     */
    public void signMessageWithPrivateKey(pactus.util.UtilOuterClass.SignMessageWithPrivateKeyRequest request,
        io.grpc.stub.StreamObserver<pactus.util.UtilOuterClass.SignMessageWithPrivateKeyResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getSignMessageWithPrivateKeyMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * VerifyMessage
     * </pre>
     */
    public void verifyMessage(pactus.util.UtilOuterClass.VerifyMessageRequest request,
        io.grpc.stub.StreamObserver<pactus.util.UtilOuterClass.VerifyMessageResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getVerifyMessageMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   * <pre>
   * Util service defines various RPC methods for interacting with
   * Utils.
   * </pre>
   */
  public static final class UtilBlockingStub extends io.grpc.stub.AbstractBlockingStub<UtilBlockingStub> {
    private UtilBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected UtilBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new UtilBlockingStub(channel, callOptions);
    }

    /**
     * <pre>
     * SignMessageWithPrivateKey
     * </pre>
     */
    public pactus.util.UtilOuterClass.SignMessageWithPrivateKeyResponse signMessageWithPrivateKey(pactus.util.UtilOuterClass.SignMessageWithPrivateKeyRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getSignMessageWithPrivateKeyMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * VerifyMessage
     * </pre>
     */
    public pactus.util.UtilOuterClass.VerifyMessageResponse verifyMessage(pactus.util.UtilOuterClass.VerifyMessageRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getVerifyMessageMethod(), getCallOptions(), request);
    }
  }

  /**
   * <pre>
   * Util service defines various RPC methods for interacting with
   * Utils.
   * </pre>
   */
  public static final class UtilFutureStub extends io.grpc.stub.AbstractFutureStub<UtilFutureStub> {
    private UtilFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected UtilFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new UtilFutureStub(channel, callOptions);
    }

    /**
     * <pre>
     * SignMessageWithPrivateKey
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.util.UtilOuterClass.SignMessageWithPrivateKeyResponse> signMessageWithPrivateKey(
        pactus.util.UtilOuterClass.SignMessageWithPrivateKeyRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getSignMessageWithPrivateKeyMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * VerifyMessage
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.util.UtilOuterClass.VerifyMessageResponse> verifyMessage(
        pactus.util.UtilOuterClass.VerifyMessageRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getVerifyMessageMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_SIGN_MESSAGE_WITH_PRIVATE_KEY = 0;
  private static final int METHODID_VERIFY_MESSAGE = 1;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final UtilImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(UtilImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_SIGN_MESSAGE_WITH_PRIVATE_KEY:
          serviceImpl.signMessageWithPrivateKey((pactus.util.UtilOuterClass.SignMessageWithPrivateKeyRequest) request,
              (io.grpc.stub.StreamObserver<pactus.util.UtilOuterClass.SignMessageWithPrivateKeyResponse>) responseObserver);
          break;
        case METHODID_VERIFY_MESSAGE:
          serviceImpl.verifyMessage((pactus.util.UtilOuterClass.VerifyMessageRequest) request,
              (io.grpc.stub.StreamObserver<pactus.util.UtilOuterClass.VerifyMessageResponse>) responseObserver);
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

  private static abstract class UtilBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    UtilBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return pactus.util.UtilOuterClass.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("Util");
    }
  }

  private static final class UtilFileDescriptorSupplier
      extends UtilBaseDescriptorSupplier {
    UtilFileDescriptorSupplier() {}
  }

  private static final class UtilMethodDescriptorSupplier
      extends UtilBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    UtilMethodDescriptorSupplier(String methodName) {
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
      synchronized (UtilGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new UtilFileDescriptorSupplier())
              .addMethod(getSignMessageWithPrivateKeyMethod())
              .addMethod(getVerifyMessageMethod())
              .build();
        }
      }
    }
    return result;
  }
}
