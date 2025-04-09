package pactus;

import static io.grpc.MethodDescriptor.generateFullMethodName;

/**
 * <pre>
 * NetworkService provides RPCs for retrieving information about the network.
 * </pre>
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.71.0)",
    comments = "Source: network.proto")
@io.grpc.stub.annotations.GrpcGenerated
public final class NetworkServiceGrpc {

  private NetworkServiceGrpc() {}

  public static final java.lang.String SERVICE_NAME = "pactus.NetworkService";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<pactus.Network.GetNetworkInfoRequest,
      pactus.Network.GetNetworkInfoResponse> getGetNetworkInfoMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetNetworkInfo",
      requestType = pactus.Network.GetNetworkInfoRequest.class,
      responseType = pactus.Network.GetNetworkInfoResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.Network.GetNetworkInfoRequest,
      pactus.Network.GetNetworkInfoResponse> getGetNetworkInfoMethod() {
    io.grpc.MethodDescriptor<pactus.Network.GetNetworkInfoRequest, pactus.Network.GetNetworkInfoResponse> getGetNetworkInfoMethod;
    if ((getGetNetworkInfoMethod = NetworkServiceGrpc.getGetNetworkInfoMethod) == null) {
      synchronized (NetworkServiceGrpc.class) {
        if ((getGetNetworkInfoMethod = NetworkServiceGrpc.getGetNetworkInfoMethod) == null) {
          NetworkServiceGrpc.getGetNetworkInfoMethod = getGetNetworkInfoMethod =
              io.grpc.MethodDescriptor.<pactus.Network.GetNetworkInfoRequest, pactus.Network.GetNetworkInfoResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetNetworkInfo"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Network.GetNetworkInfoRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Network.GetNetworkInfoResponse.getDefaultInstance()))
              .setSchemaDescriptor(new NetworkServiceMethodDescriptorSupplier("GetNetworkInfo"))
              .build();
        }
      }
    }
    return getGetNetworkInfoMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.Network.GetNodeInfoRequest,
      pactus.Network.GetNodeInfoResponse> getGetNodeInfoMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetNodeInfo",
      requestType = pactus.Network.GetNodeInfoRequest.class,
      responseType = pactus.Network.GetNodeInfoResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.Network.GetNodeInfoRequest,
      pactus.Network.GetNodeInfoResponse> getGetNodeInfoMethod() {
    io.grpc.MethodDescriptor<pactus.Network.GetNodeInfoRequest, pactus.Network.GetNodeInfoResponse> getGetNodeInfoMethod;
    if ((getGetNodeInfoMethod = NetworkServiceGrpc.getGetNodeInfoMethod) == null) {
      synchronized (NetworkServiceGrpc.class) {
        if ((getGetNodeInfoMethod = NetworkServiceGrpc.getGetNodeInfoMethod) == null) {
          NetworkServiceGrpc.getGetNodeInfoMethod = getGetNodeInfoMethod =
              io.grpc.MethodDescriptor.<pactus.Network.GetNodeInfoRequest, pactus.Network.GetNodeInfoResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetNodeInfo"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Network.GetNodeInfoRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Network.GetNodeInfoResponse.getDefaultInstance()))
              .setSchemaDescriptor(new NetworkServiceMethodDescriptorSupplier("GetNodeInfo"))
              .build();
        }
      }
    }
    return getGetNodeInfoMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static NetworkServiceStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<NetworkServiceStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<NetworkServiceStub>() {
        @java.lang.Override
        public NetworkServiceStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new NetworkServiceStub(channel, callOptions);
        }
      };
    return NetworkServiceStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports all types of calls on the service
   */
  public static NetworkServiceBlockingV2Stub newBlockingV2Stub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<NetworkServiceBlockingV2Stub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<NetworkServiceBlockingV2Stub>() {
        @java.lang.Override
        public NetworkServiceBlockingV2Stub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new NetworkServiceBlockingV2Stub(channel, callOptions);
        }
      };
    return NetworkServiceBlockingV2Stub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static NetworkServiceBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<NetworkServiceBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<NetworkServiceBlockingStub>() {
        @java.lang.Override
        public NetworkServiceBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new NetworkServiceBlockingStub(channel, callOptions);
        }
      };
    return NetworkServiceBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static NetworkServiceFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<NetworkServiceFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<NetworkServiceFutureStub>() {
        @java.lang.Override
        public NetworkServiceFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new NetworkServiceFutureStub(channel, callOptions);
        }
      };
    return NetworkServiceFutureStub.newStub(factory, channel);
  }

  /**
   * <pre>
   * NetworkService provides RPCs for retrieving information about the network.
   * </pre>
   */
  public interface AsyncService {

    /**
     * <pre>
     * GetNetworkInfo retrieves information about the overall network.
     * </pre>
     */
    default void getNetworkInfo(pactus.Network.GetNetworkInfoRequest request,
        io.grpc.stub.StreamObserver<pactus.Network.GetNetworkInfoResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetNetworkInfoMethod(), responseObserver);
    }

    /**
     * <pre>
     * GetNodeInfo retrieves information about a specific node in the network.
     * </pre>
     */
    default void getNodeInfo(pactus.Network.GetNodeInfoRequest request,
        io.grpc.stub.StreamObserver<pactus.Network.GetNodeInfoResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetNodeInfoMethod(), responseObserver);
    }
  }

  /**
   * Base class for the server implementation of the service NetworkService.
   * <pre>
   * NetworkService provides RPCs for retrieving information about the network.
   * </pre>
   */
  public static abstract class NetworkServiceImplBase
      implements io.grpc.BindableService, AsyncService {

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return NetworkServiceGrpc.bindService(this);
    }
  }

  /**
   * A stub to allow clients to do asynchronous rpc calls to service NetworkService.
   * <pre>
   * NetworkService provides RPCs for retrieving information about the network.
   * </pre>
   */
  public static final class NetworkServiceStub
      extends io.grpc.stub.AbstractAsyncStub<NetworkServiceStub> {
    private NetworkServiceStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected NetworkServiceStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new NetworkServiceStub(channel, callOptions);
    }

    /**
     * <pre>
     * GetNetworkInfo retrieves information about the overall network.
     * </pre>
     */
    public void getNetworkInfo(pactus.Network.GetNetworkInfoRequest request,
        io.grpc.stub.StreamObserver<pactus.Network.GetNetworkInfoResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetNetworkInfoMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * GetNodeInfo retrieves information about a specific node in the network.
     * </pre>
     */
    public void getNodeInfo(pactus.Network.GetNodeInfoRequest request,
        io.grpc.stub.StreamObserver<pactus.Network.GetNodeInfoResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetNodeInfoMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   * A stub to allow clients to do synchronous rpc calls to service NetworkService.
   * <pre>
   * NetworkService provides RPCs for retrieving information about the network.
   * </pre>
   */
  public static final class NetworkServiceBlockingV2Stub
      extends io.grpc.stub.AbstractBlockingStub<NetworkServiceBlockingV2Stub> {
    private NetworkServiceBlockingV2Stub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected NetworkServiceBlockingV2Stub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new NetworkServiceBlockingV2Stub(channel, callOptions);
    }

    /**
     * <pre>
     * GetNetworkInfo retrieves information about the overall network.
     * </pre>
     */
    public pactus.Network.GetNetworkInfoResponse getNetworkInfo(pactus.Network.GetNetworkInfoRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetNetworkInfoMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetNodeInfo retrieves information about a specific node in the network.
     * </pre>
     */
    public pactus.Network.GetNodeInfoResponse getNodeInfo(pactus.Network.GetNodeInfoRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetNodeInfoMethod(), getCallOptions(), request);
    }
  }

  /**
   * A stub to allow clients to do limited synchronous rpc calls to service NetworkService.
   * <pre>
   * NetworkService provides RPCs for retrieving information about the network.
   * </pre>
   */
  public static final class NetworkServiceBlockingStub
      extends io.grpc.stub.AbstractBlockingStub<NetworkServiceBlockingStub> {
    private NetworkServiceBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected NetworkServiceBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new NetworkServiceBlockingStub(channel, callOptions);
    }

    /**
     * <pre>
     * GetNetworkInfo retrieves information about the overall network.
     * </pre>
     */
    public pactus.Network.GetNetworkInfoResponse getNetworkInfo(pactus.Network.GetNetworkInfoRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetNetworkInfoMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetNodeInfo retrieves information about a specific node in the network.
     * </pre>
     */
    public pactus.Network.GetNodeInfoResponse getNodeInfo(pactus.Network.GetNodeInfoRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetNodeInfoMethod(), getCallOptions(), request);
    }
  }

  /**
   * A stub to allow clients to do ListenableFuture-style rpc calls to service NetworkService.
   * <pre>
   * NetworkService provides RPCs for retrieving information about the network.
   * </pre>
   */
  public static final class NetworkServiceFutureStub
      extends io.grpc.stub.AbstractFutureStub<NetworkServiceFutureStub> {
    private NetworkServiceFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected NetworkServiceFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new NetworkServiceFutureStub(channel, callOptions);
    }

    /**
     * <pre>
     * GetNetworkInfo retrieves information about the overall network.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.Network.GetNetworkInfoResponse> getNetworkInfo(
        pactus.Network.GetNetworkInfoRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetNetworkInfoMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * GetNodeInfo retrieves information about a specific node in the network.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.Network.GetNodeInfoResponse> getNodeInfo(
        pactus.Network.GetNodeInfoRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetNodeInfoMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_GET_NETWORK_INFO = 0;
  private static final int METHODID_GET_NODE_INFO = 1;

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
          serviceImpl.getNetworkInfo((pactus.Network.GetNetworkInfoRequest) request,
              (io.grpc.stub.StreamObserver<pactus.Network.GetNetworkInfoResponse>) responseObserver);
          break;
        case METHODID_GET_NODE_INFO:
          serviceImpl.getNodeInfo((pactus.Network.GetNodeInfoRequest) request,
              (io.grpc.stub.StreamObserver<pactus.Network.GetNodeInfoResponse>) responseObserver);
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
              pactus.Network.GetNetworkInfoRequest,
              pactus.Network.GetNetworkInfoResponse>(
                service, METHODID_GET_NETWORK_INFO)))
        .addMethod(
          getGetNodeInfoMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.Network.GetNodeInfoRequest,
              pactus.Network.GetNodeInfoResponse>(
                service, METHODID_GET_NODE_INFO)))
        .build();
  }

  private static abstract class NetworkServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    NetworkServiceBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return pactus.Network.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("NetworkService");
    }
  }

  private static final class NetworkServiceFileDescriptorSupplier
      extends NetworkServiceBaseDescriptorSupplier {
    NetworkServiceFileDescriptorSupplier() {}
  }

  private static final class NetworkServiceMethodDescriptorSupplier
      extends NetworkServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final java.lang.String methodName;

    NetworkServiceMethodDescriptorSupplier(java.lang.String methodName) {
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
      synchronized (NetworkServiceGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new NetworkServiceFileDescriptorSupplier())
              .addMethod(getGetNetworkInfoMethod())
              .addMethod(getGetNodeInfoMethod())
              .build();
        }
      }
    }
    return result;
  }
}
