package pactus;

import static io.grpc.MethodDescriptor.generateFullMethodName;

/**
 * <pre>
 * Network service provides RPCs for retrieving information about the network.
 * </pre>
 */
@io.grpc.stub.annotations.GrpcGenerated
public final class NetworkGrpc {

  private NetworkGrpc() {}

  public static final java.lang.String SERVICE_NAME = "pactus.Network";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<pactus.NetworkOuterClass.GetNetworkInfoRequest,
      pactus.NetworkOuterClass.GetNetworkInfoResponse> getGetNetworkInfoMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetNetworkInfo",
      requestType = pactus.NetworkOuterClass.GetNetworkInfoRequest.class,
      responseType = pactus.NetworkOuterClass.GetNetworkInfoResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.NetworkOuterClass.GetNetworkInfoRequest,
      pactus.NetworkOuterClass.GetNetworkInfoResponse> getGetNetworkInfoMethod() {
    io.grpc.MethodDescriptor<pactus.NetworkOuterClass.GetNetworkInfoRequest, pactus.NetworkOuterClass.GetNetworkInfoResponse> getGetNetworkInfoMethod;
    if ((getGetNetworkInfoMethod = NetworkGrpc.getGetNetworkInfoMethod) == null) {
      synchronized (NetworkGrpc.class) {
        if ((getGetNetworkInfoMethod = NetworkGrpc.getGetNetworkInfoMethod) == null) {
          NetworkGrpc.getGetNetworkInfoMethod = getGetNetworkInfoMethod =
              io.grpc.MethodDescriptor.<pactus.NetworkOuterClass.GetNetworkInfoRequest, pactus.NetworkOuterClass.GetNetworkInfoResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetNetworkInfo"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.NetworkOuterClass.GetNetworkInfoRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.NetworkOuterClass.GetNetworkInfoResponse.getDefaultInstance()))
              .setSchemaDescriptor(new NetworkMethodDescriptorSupplier("GetNetworkInfo"))
              .build();
        }
      }
    }
    return getGetNetworkInfoMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.NetworkOuterClass.GetNodeInfoRequest,
      pactus.NetworkOuterClass.GetNodeInfoResponse> getGetNodeInfoMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetNodeInfo",
      requestType = pactus.NetworkOuterClass.GetNodeInfoRequest.class,
      responseType = pactus.NetworkOuterClass.GetNodeInfoResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.NetworkOuterClass.GetNodeInfoRequest,
      pactus.NetworkOuterClass.GetNodeInfoResponse> getGetNodeInfoMethod() {
    io.grpc.MethodDescriptor<pactus.NetworkOuterClass.GetNodeInfoRequest, pactus.NetworkOuterClass.GetNodeInfoResponse> getGetNodeInfoMethod;
    if ((getGetNodeInfoMethod = NetworkGrpc.getGetNodeInfoMethod) == null) {
      synchronized (NetworkGrpc.class) {
        if ((getGetNodeInfoMethod = NetworkGrpc.getGetNodeInfoMethod) == null) {
          NetworkGrpc.getGetNodeInfoMethod = getGetNodeInfoMethod =
              io.grpc.MethodDescriptor.<pactus.NetworkOuterClass.GetNodeInfoRequest, pactus.NetworkOuterClass.GetNodeInfoResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetNodeInfo"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.NetworkOuterClass.GetNodeInfoRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.NetworkOuterClass.GetNodeInfoResponse.getDefaultInstance()))
              .setSchemaDescriptor(new NetworkMethodDescriptorSupplier("GetNodeInfo"))
              .build();
        }
      }
    }
    return getGetNodeInfoMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.NetworkOuterClass.PingRequest,
      pactus.NetworkOuterClass.PingResponse> getPingMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Ping",
      requestType = pactus.NetworkOuterClass.PingRequest.class,
      responseType = pactus.NetworkOuterClass.PingResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.NetworkOuterClass.PingRequest,
      pactus.NetworkOuterClass.PingResponse> getPingMethod() {
    io.grpc.MethodDescriptor<pactus.NetworkOuterClass.PingRequest, pactus.NetworkOuterClass.PingResponse> getPingMethod;
    if ((getPingMethod = NetworkGrpc.getPingMethod) == null) {
      synchronized (NetworkGrpc.class) {
        if ((getPingMethod = NetworkGrpc.getPingMethod) == null) {
          NetworkGrpc.getPingMethod = getPingMethod =
              io.grpc.MethodDescriptor.<pactus.NetworkOuterClass.PingRequest, pactus.NetworkOuterClass.PingResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Ping"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.NetworkOuterClass.PingRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.NetworkOuterClass.PingResponse.getDefaultInstance()))
              .setSchemaDescriptor(new NetworkMethodDescriptorSupplier("Ping"))
              .build();
        }
      }
    }
    return getPingMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static NetworkStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<NetworkStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<NetworkStub>() {
        @java.lang.Override
        public NetworkStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new NetworkStub(channel, callOptions);
        }
      };
    return NetworkStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports all types of calls on the service
   */
  public static NetworkBlockingV2Stub newBlockingV2Stub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<NetworkBlockingV2Stub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<NetworkBlockingV2Stub>() {
        @java.lang.Override
        public NetworkBlockingV2Stub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new NetworkBlockingV2Stub(channel, callOptions);
        }
      };
    return NetworkBlockingV2Stub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static NetworkBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<NetworkBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<NetworkBlockingStub>() {
        @java.lang.Override
        public NetworkBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new NetworkBlockingStub(channel, callOptions);
        }
      };
    return NetworkBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static NetworkFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<NetworkFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<NetworkFutureStub>() {
        @java.lang.Override
        public NetworkFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new NetworkFutureStub(channel, callOptions);
        }
      };
    return NetworkFutureStub.newStub(factory, channel);
  }

  /**
   * <pre>
   * Network service provides RPCs for retrieving information about the network.
   * </pre>
   */
  public interface AsyncService {

    /**
     * <pre>
     * GetNetworkInfo retrieves information about the overall network.
     * </pre>
     */
    default void getNetworkInfo(pactus.NetworkOuterClass.GetNetworkInfoRequest request,
        io.grpc.stub.StreamObserver<pactus.NetworkOuterClass.GetNetworkInfoResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetNetworkInfoMethod(), responseObserver);
    }

    /**
     * <pre>
     * GetNodeInfo retrieves information about a specific node in the network.
     * </pre>
     */
    default void getNodeInfo(pactus.NetworkOuterClass.GetNodeInfoRequest request,
        io.grpc.stub.StreamObserver<pactus.NetworkOuterClass.GetNodeInfoResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetNodeInfoMethod(), responseObserver);
    }

    /**
     * <pre>
     * Ping provides a simple connectivity test and latency measurement.
     * </pre>
     */
    default void ping(pactus.NetworkOuterClass.PingRequest request,
        io.grpc.stub.StreamObserver<pactus.NetworkOuterClass.PingResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getPingMethod(), responseObserver);
    }
  }

  /**
   * Base class for the server implementation of the service Network.
   * <pre>
   * Network service provides RPCs for retrieving information about the network.
   * </pre>
   */
  public static abstract class NetworkImplBase
      implements io.grpc.BindableService, AsyncService {

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return NetworkGrpc.bindService(this);
    }
  }

  /**
   * A stub to allow clients to do asynchronous rpc calls to service Network.
   * <pre>
   * Network service provides RPCs for retrieving information about the network.
   * </pre>
   */
  public static final class NetworkStub
      extends io.grpc.stub.AbstractAsyncStub<NetworkStub> {
    private NetworkStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected NetworkStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new NetworkStub(channel, callOptions);
    }

    /**
     * <pre>
     * GetNetworkInfo retrieves information about the overall network.
     * </pre>
     */
    public void getNetworkInfo(pactus.NetworkOuterClass.GetNetworkInfoRequest request,
        io.grpc.stub.StreamObserver<pactus.NetworkOuterClass.GetNetworkInfoResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetNetworkInfoMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * GetNodeInfo retrieves information about a specific node in the network.
     * </pre>
     */
    public void getNodeInfo(pactus.NetworkOuterClass.GetNodeInfoRequest request,
        io.grpc.stub.StreamObserver<pactus.NetworkOuterClass.GetNodeInfoResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetNodeInfoMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * Ping provides a simple connectivity test and latency measurement.
     * </pre>
     */
    public void ping(pactus.NetworkOuterClass.PingRequest request,
        io.grpc.stub.StreamObserver<pactus.NetworkOuterClass.PingResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getPingMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   * A stub to allow clients to do synchronous rpc calls to service Network.
   * <pre>
   * Network service provides RPCs for retrieving information about the network.
   * </pre>
   */
  public static final class NetworkBlockingV2Stub
      extends io.grpc.stub.AbstractBlockingStub<NetworkBlockingV2Stub> {
    private NetworkBlockingV2Stub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected NetworkBlockingV2Stub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new NetworkBlockingV2Stub(channel, callOptions);
    }

    /**
     * <pre>
     * GetNetworkInfo retrieves information about the overall network.
     * </pre>
     */
    public pactus.NetworkOuterClass.GetNetworkInfoResponse getNetworkInfo(pactus.NetworkOuterClass.GetNetworkInfoRequest request) throws io.grpc.StatusException {
      return io.grpc.stub.ClientCalls.blockingV2UnaryCall(
          getChannel(), getGetNetworkInfoMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetNodeInfo retrieves information about a specific node in the network.
     * </pre>
     */
    public pactus.NetworkOuterClass.GetNodeInfoResponse getNodeInfo(pactus.NetworkOuterClass.GetNodeInfoRequest request) throws io.grpc.StatusException {
      return io.grpc.stub.ClientCalls.blockingV2UnaryCall(
          getChannel(), getGetNodeInfoMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * Ping provides a simple connectivity test and latency measurement.
     * </pre>
     */
    public pactus.NetworkOuterClass.PingResponse ping(pactus.NetworkOuterClass.PingRequest request) throws io.grpc.StatusException {
      return io.grpc.stub.ClientCalls.blockingV2UnaryCall(
          getChannel(), getPingMethod(), getCallOptions(), request);
    }
  }

  /**
   * A stub to allow clients to do limited synchronous rpc calls to service Network.
   * <pre>
   * Network service provides RPCs for retrieving information about the network.
   * </pre>
   */
  public static final class NetworkBlockingStub
      extends io.grpc.stub.AbstractBlockingStub<NetworkBlockingStub> {
    private NetworkBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected NetworkBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new NetworkBlockingStub(channel, callOptions);
    }

    /**
     * <pre>
     * GetNetworkInfo retrieves information about the overall network.
     * </pre>
     */
    public pactus.NetworkOuterClass.GetNetworkInfoResponse getNetworkInfo(pactus.NetworkOuterClass.GetNetworkInfoRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetNetworkInfoMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetNodeInfo retrieves information about a specific node in the network.
     * </pre>
     */
    public pactus.NetworkOuterClass.GetNodeInfoResponse getNodeInfo(pactus.NetworkOuterClass.GetNodeInfoRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetNodeInfoMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * Ping provides a simple connectivity test and latency measurement.
     * </pre>
     */
    public pactus.NetworkOuterClass.PingResponse ping(pactus.NetworkOuterClass.PingRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getPingMethod(), getCallOptions(), request);
    }
  }

  /**
   * A stub to allow clients to do ListenableFuture-style rpc calls to service Network.
   * <pre>
   * Network service provides RPCs for retrieving information about the network.
   * </pre>
   */
  public static final class NetworkFutureStub
      extends io.grpc.stub.AbstractFutureStub<NetworkFutureStub> {
    private NetworkFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected NetworkFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new NetworkFutureStub(channel, callOptions);
    }

    /**
     * <pre>
     * GetNetworkInfo retrieves information about the overall network.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.NetworkOuterClass.GetNetworkInfoResponse> getNetworkInfo(
        pactus.NetworkOuterClass.GetNetworkInfoRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetNetworkInfoMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * GetNodeInfo retrieves information about a specific node in the network.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.NetworkOuterClass.GetNodeInfoResponse> getNodeInfo(
        pactus.NetworkOuterClass.GetNodeInfoRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetNodeInfoMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * Ping provides a simple connectivity test and latency measurement.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.NetworkOuterClass.PingResponse> ping(
        pactus.NetworkOuterClass.PingRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getPingMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_GET_NETWORK_INFO = 0;
  private static final int METHODID_GET_NODE_INFO = 1;
  private static final int METHODID_PING = 2;

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
        case METHODID_GET_NETWORK_INFO:
          serviceImpl.getNetworkInfo((pactus.NetworkOuterClass.GetNetworkInfoRequest) request,
              (io.grpc.stub.StreamObserver<pactus.NetworkOuterClass.GetNetworkInfoResponse>) responseObserver);
          break;
        case METHODID_GET_NODE_INFO:
          serviceImpl.getNodeInfo((pactus.NetworkOuterClass.GetNodeInfoRequest) request,
              (io.grpc.stub.StreamObserver<pactus.NetworkOuterClass.GetNodeInfoResponse>) responseObserver);
          break;
        case METHODID_PING:
          serviceImpl.ping((pactus.NetworkOuterClass.PingRequest) request,
              (io.grpc.stub.StreamObserver<pactus.NetworkOuterClass.PingResponse>) responseObserver);
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
          getGetNetworkInfoMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.NetworkOuterClass.GetNetworkInfoRequest,
              pactus.NetworkOuterClass.GetNetworkInfoResponse>(
                service, METHODID_GET_NETWORK_INFO)))
        .addMethod(
          getGetNodeInfoMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.NetworkOuterClass.GetNodeInfoRequest,
              pactus.NetworkOuterClass.GetNodeInfoResponse>(
                service, METHODID_GET_NODE_INFO)))
        .addMethod(
          getPingMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.NetworkOuterClass.PingRequest,
              pactus.NetworkOuterClass.PingResponse>(
                service, METHODID_PING)))
        .build();
  }

  private static abstract class NetworkBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    NetworkBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return pactus.NetworkOuterClass.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("Network");
    }
  }

  private static final class NetworkFileDescriptorSupplier
      extends NetworkBaseDescriptorSupplier {
    NetworkFileDescriptorSupplier() {}
  }

  private static final class NetworkMethodDescriptorSupplier
      extends NetworkBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final java.lang.String methodName;

    NetworkMethodDescriptorSupplier(java.lang.String methodName) {
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
      synchronized (NetworkGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new NetworkFileDescriptorSupplier())
              .addMethod(getGetNetworkInfoMethod())
              .addMethod(getGetNodeInfoMethod())
              .addMethod(getPingMethod())
              .build();
        }
      }
    }
    return result;
  }
}
