package pactus.transaction;

import static io.grpc.MethodDescriptor.generateFullMethodName;

/**
 * <pre>
 * Transaction service defines various RPC methods for interacting with
 * transactions.
 * </pre>
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.50.2)",
    comments = "Source: transaction.proto")
@io.grpc.stub.annotations.GrpcGenerated
public final class TransactionGrpc {

  private TransactionGrpc() {}

  public static final String SERVICE_NAME = "pactus.Transaction";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<pactus.transaction.TransactionOuterClass.GetTransactionRequest,
      pactus.transaction.TransactionOuterClass.GetTransactionResponse> getGetTransactionMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetTransaction",
      requestType = pactus.transaction.TransactionOuterClass.GetTransactionRequest.class,
      responseType = pactus.transaction.TransactionOuterClass.GetTransactionResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.transaction.TransactionOuterClass.GetTransactionRequest,
      pactus.transaction.TransactionOuterClass.GetTransactionResponse> getGetTransactionMethod() {
    io.grpc.MethodDescriptor<pactus.transaction.TransactionOuterClass.GetTransactionRequest, pactus.transaction.TransactionOuterClass.GetTransactionResponse> getGetTransactionMethod;
    if ((getGetTransactionMethod = TransactionGrpc.getGetTransactionMethod) == null) {
      synchronized (TransactionGrpc.class) {
        if ((getGetTransactionMethod = TransactionGrpc.getGetTransactionMethod) == null) {
          TransactionGrpc.getGetTransactionMethod = getGetTransactionMethod =
              io.grpc.MethodDescriptor.<pactus.transaction.TransactionOuterClass.GetTransactionRequest, pactus.transaction.TransactionOuterClass.GetTransactionResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetTransaction"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.transaction.TransactionOuterClass.GetTransactionRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.transaction.TransactionOuterClass.GetTransactionResponse.getDefaultInstance()))
              .setSchemaDescriptor(new TransactionMethodDescriptorSupplier("GetTransaction"))
              .build();
        }
      }
    }
    return getGetTransactionMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.transaction.TransactionOuterClass.CalculateFeeRequest,
      pactus.transaction.TransactionOuterClass.CalculateFeeResponse> getCalculateFeeMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "CalculateFee",
      requestType = pactus.transaction.TransactionOuterClass.CalculateFeeRequest.class,
      responseType = pactus.transaction.TransactionOuterClass.CalculateFeeResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.transaction.TransactionOuterClass.CalculateFeeRequest,
      pactus.transaction.TransactionOuterClass.CalculateFeeResponse> getCalculateFeeMethod() {
    io.grpc.MethodDescriptor<pactus.transaction.TransactionOuterClass.CalculateFeeRequest, pactus.transaction.TransactionOuterClass.CalculateFeeResponse> getCalculateFeeMethod;
    if ((getCalculateFeeMethod = TransactionGrpc.getCalculateFeeMethod) == null) {
      synchronized (TransactionGrpc.class) {
        if ((getCalculateFeeMethod = TransactionGrpc.getCalculateFeeMethod) == null) {
          TransactionGrpc.getCalculateFeeMethod = getCalculateFeeMethod =
              io.grpc.MethodDescriptor.<pactus.transaction.TransactionOuterClass.CalculateFeeRequest, pactus.transaction.TransactionOuterClass.CalculateFeeResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "CalculateFee"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.transaction.TransactionOuterClass.CalculateFeeRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.transaction.TransactionOuterClass.CalculateFeeResponse.getDefaultInstance()))
              .setSchemaDescriptor(new TransactionMethodDescriptorSupplier("CalculateFee"))
              .build();
        }
      }
    }
    return getCalculateFeeMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.transaction.TransactionOuterClass.BroadcastTransactionRequest,
      pactus.transaction.TransactionOuterClass.BroadcastTransactionResponse> getBroadcastTransactionMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "BroadcastTransaction",
      requestType = pactus.transaction.TransactionOuterClass.BroadcastTransactionRequest.class,
      responseType = pactus.transaction.TransactionOuterClass.BroadcastTransactionResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.transaction.TransactionOuterClass.BroadcastTransactionRequest,
      pactus.transaction.TransactionOuterClass.BroadcastTransactionResponse> getBroadcastTransactionMethod() {
    io.grpc.MethodDescriptor<pactus.transaction.TransactionOuterClass.BroadcastTransactionRequest, pactus.transaction.TransactionOuterClass.BroadcastTransactionResponse> getBroadcastTransactionMethod;
    if ((getBroadcastTransactionMethod = TransactionGrpc.getBroadcastTransactionMethod) == null) {
      synchronized (TransactionGrpc.class) {
        if ((getBroadcastTransactionMethod = TransactionGrpc.getBroadcastTransactionMethod) == null) {
          TransactionGrpc.getBroadcastTransactionMethod = getBroadcastTransactionMethod =
              io.grpc.MethodDescriptor.<pactus.transaction.TransactionOuterClass.BroadcastTransactionRequest, pactus.transaction.TransactionOuterClass.BroadcastTransactionResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "BroadcastTransaction"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.transaction.TransactionOuterClass.BroadcastTransactionRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.transaction.TransactionOuterClass.BroadcastTransactionResponse.getDefaultInstance()))
              .setSchemaDescriptor(new TransactionMethodDescriptorSupplier("BroadcastTransaction"))
              .build();
        }
      }
    }
    return getBroadcastTransactionMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.transaction.TransactionOuterClass.GetRawTransactionRequest,
      pactus.transaction.TransactionOuterClass.GetRawTransactionResponse> getGetRawTransactionMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetRawTransaction",
      requestType = pactus.transaction.TransactionOuterClass.GetRawTransactionRequest.class,
      responseType = pactus.transaction.TransactionOuterClass.GetRawTransactionResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.transaction.TransactionOuterClass.GetRawTransactionRequest,
      pactus.transaction.TransactionOuterClass.GetRawTransactionResponse> getGetRawTransactionMethod() {
    io.grpc.MethodDescriptor<pactus.transaction.TransactionOuterClass.GetRawTransactionRequest, pactus.transaction.TransactionOuterClass.GetRawTransactionResponse> getGetRawTransactionMethod;
    if ((getGetRawTransactionMethod = TransactionGrpc.getGetRawTransactionMethod) == null) {
      synchronized (TransactionGrpc.class) {
        if ((getGetRawTransactionMethod = TransactionGrpc.getGetRawTransactionMethod) == null) {
          TransactionGrpc.getGetRawTransactionMethod = getGetRawTransactionMethod =
              io.grpc.MethodDescriptor.<pactus.transaction.TransactionOuterClass.GetRawTransactionRequest, pactus.transaction.TransactionOuterClass.GetRawTransactionResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetRawTransaction"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.transaction.TransactionOuterClass.GetRawTransactionRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.transaction.TransactionOuterClass.GetRawTransactionResponse.getDefaultInstance()))
              .setSchemaDescriptor(new TransactionMethodDescriptorSupplier("GetRawTransaction"))
              .build();
        }
      }
    }
    return getGetRawTransactionMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static TransactionStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<TransactionStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<TransactionStub>() {
        @java.lang.Override
        public TransactionStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new TransactionStub(channel, callOptions);
        }
      };
    return TransactionStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static TransactionBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<TransactionBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<TransactionBlockingStub>() {
        @java.lang.Override
        public TransactionBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new TransactionBlockingStub(channel, callOptions);
        }
      };
    return TransactionBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static TransactionFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<TransactionFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<TransactionFutureStub>() {
        @java.lang.Override
        public TransactionFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new TransactionFutureStub(channel, callOptions);
        }
      };
    return TransactionFutureStub.newStub(factory, channel);
  }

  /**
   * <pre>
   * Transaction service defines various RPC methods for interacting with
   * transactions.
   * </pre>
   */
  public static abstract class TransactionImplBase implements io.grpc.BindableService {

    /**
     * <pre>
     * GetTransaction retrieves transaction details based on the provided request
     * parameters.
     * </pre>
     */
    public void getTransaction(pactus.transaction.TransactionOuterClass.GetTransactionRequest request,
        io.grpc.stub.StreamObserver<pactus.transaction.TransactionOuterClass.GetTransactionResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetTransactionMethod(), responseObserver);
    }

    /**
     * <pre>
     * CalculateFee calculates the transaction fee based on the specified amount
     * and payload type.
     * </pre>
     */
    public void calculateFee(pactus.transaction.TransactionOuterClass.CalculateFeeRequest request,
        io.grpc.stub.StreamObserver<pactus.transaction.TransactionOuterClass.CalculateFeeResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getCalculateFeeMethod(), responseObserver);
    }

    /**
     * <pre>
     * BroadcastTransaction broadcasts a signed transaction to the network.
     * </pre>
     */
    public void broadcastTransaction(pactus.transaction.TransactionOuterClass.BroadcastTransactionRequest request,
        io.grpc.stub.StreamObserver<pactus.transaction.TransactionOuterClass.BroadcastTransactionResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getBroadcastTransactionMethod(), responseObserver);
    }

    /**
     * <pre>
     * GetRawTransaction retrieves raw details of transfer, bond, unbond or withdraw transaction.
     * </pre>
     */
    public void getRawTransaction(pactus.transaction.TransactionOuterClass.GetRawTransactionRequest request,
        io.grpc.stub.StreamObserver<pactus.transaction.TransactionOuterClass.GetRawTransactionResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetRawTransactionMethod(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getGetTransactionMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                pactus.transaction.TransactionOuterClass.GetTransactionRequest,
                pactus.transaction.TransactionOuterClass.GetTransactionResponse>(
                  this, METHODID_GET_TRANSACTION)))
          .addMethod(
            getCalculateFeeMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                pactus.transaction.TransactionOuterClass.CalculateFeeRequest,
                pactus.transaction.TransactionOuterClass.CalculateFeeResponse>(
                  this, METHODID_CALCULATE_FEE)))
          .addMethod(
            getBroadcastTransactionMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                pactus.transaction.TransactionOuterClass.BroadcastTransactionRequest,
                pactus.transaction.TransactionOuterClass.BroadcastTransactionResponse>(
                  this, METHODID_BROADCAST_TRANSACTION)))
          .addMethod(
            getGetRawTransactionMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                pactus.transaction.TransactionOuterClass.GetRawTransactionRequest,
                pactus.transaction.TransactionOuterClass.GetRawTransactionResponse>(
                  this, METHODID_GET_RAW_TRANSACTION)))
          .build();
    }
  }

  /**
   * <pre>
   * Transaction service defines various RPC methods for interacting with
   * transactions.
   * </pre>
   */
  public static final class TransactionStub extends io.grpc.stub.AbstractAsyncStub<TransactionStub> {
    private TransactionStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected TransactionStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new TransactionStub(channel, callOptions);
    }

    /**
     * <pre>
     * GetTransaction retrieves transaction details based on the provided request
     * parameters.
     * </pre>
     */
    public void getTransaction(pactus.transaction.TransactionOuterClass.GetTransactionRequest request,
        io.grpc.stub.StreamObserver<pactus.transaction.TransactionOuterClass.GetTransactionResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetTransactionMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * CalculateFee calculates the transaction fee based on the specified amount
     * and payload type.
     * </pre>
     */
    public void calculateFee(pactus.transaction.TransactionOuterClass.CalculateFeeRequest request,
        io.grpc.stub.StreamObserver<pactus.transaction.TransactionOuterClass.CalculateFeeResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getCalculateFeeMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * BroadcastTransaction broadcasts a signed transaction to the network.
     * </pre>
     */
    public void broadcastTransaction(pactus.transaction.TransactionOuterClass.BroadcastTransactionRequest request,
        io.grpc.stub.StreamObserver<pactus.transaction.TransactionOuterClass.BroadcastTransactionResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getBroadcastTransactionMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * GetRawTransaction retrieves raw details of transfer, bond, unbond or withdraw transaction.
     * </pre>
     */
    public void getRawTransaction(pactus.transaction.TransactionOuterClass.GetRawTransactionRequest request,
        io.grpc.stub.StreamObserver<pactus.transaction.TransactionOuterClass.GetRawTransactionResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetRawTransactionMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   * <pre>
   * Transaction service defines various RPC methods for interacting with
   * transactions.
   * </pre>
   */
  public static final class TransactionBlockingStub extends io.grpc.stub.AbstractBlockingStub<TransactionBlockingStub> {
    private TransactionBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected TransactionBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new TransactionBlockingStub(channel, callOptions);
    }

    /**
     * <pre>
     * GetTransaction retrieves transaction details based on the provided request
     * parameters.
     * </pre>
     */
    public pactus.transaction.TransactionOuterClass.GetTransactionResponse getTransaction(pactus.transaction.TransactionOuterClass.GetTransactionRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetTransactionMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * CalculateFee calculates the transaction fee based on the specified amount
     * and payload type.
     * </pre>
     */
    public pactus.transaction.TransactionOuterClass.CalculateFeeResponse calculateFee(pactus.transaction.TransactionOuterClass.CalculateFeeRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getCalculateFeeMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * BroadcastTransaction broadcasts a signed transaction to the network.
     * </pre>
     */
    public pactus.transaction.TransactionOuterClass.BroadcastTransactionResponse broadcastTransaction(pactus.transaction.TransactionOuterClass.BroadcastTransactionRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getBroadcastTransactionMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetRawTransaction retrieves raw details of transfer, bond, unbond or withdraw transaction.
     * </pre>
     */
    public pactus.transaction.TransactionOuterClass.GetRawTransactionResponse getRawTransaction(pactus.transaction.TransactionOuterClass.GetRawTransactionRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetRawTransactionMethod(), getCallOptions(), request);
    }
  }

  /**
   * <pre>
   * Transaction service defines various RPC methods for interacting with
   * transactions.
   * </pre>
   */
  public static final class TransactionFutureStub extends io.grpc.stub.AbstractFutureStub<TransactionFutureStub> {
    private TransactionFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected TransactionFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new TransactionFutureStub(channel, callOptions);
    }

    /**
     * <pre>
     * GetTransaction retrieves transaction details based on the provided request
     * parameters.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.transaction.TransactionOuterClass.GetTransactionResponse> getTransaction(
        pactus.transaction.TransactionOuterClass.GetTransactionRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetTransactionMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * CalculateFee calculates the transaction fee based on the specified amount
     * and payload type.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.transaction.TransactionOuterClass.CalculateFeeResponse> calculateFee(
        pactus.transaction.TransactionOuterClass.CalculateFeeRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getCalculateFeeMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * BroadcastTransaction broadcasts a signed transaction to the network.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.transaction.TransactionOuterClass.BroadcastTransactionResponse> broadcastTransaction(
        pactus.transaction.TransactionOuterClass.BroadcastTransactionRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getBroadcastTransactionMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * GetRawTransaction retrieves raw details of transfer, bond, unbond or withdraw transaction.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.transaction.TransactionOuterClass.GetRawTransactionResponse> getRawTransaction(
        pactus.transaction.TransactionOuterClass.GetRawTransactionRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetRawTransactionMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_GET_TRANSACTION = 0;
  private static final int METHODID_CALCULATE_FEE = 1;
  private static final int METHODID_BROADCAST_TRANSACTION = 2;
  private static final int METHODID_GET_RAW_TRANSACTION = 3;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final TransactionImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(TransactionImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_GET_TRANSACTION:
          serviceImpl.getTransaction((pactus.transaction.TransactionOuterClass.GetTransactionRequest) request,
              (io.grpc.stub.StreamObserver<pactus.transaction.TransactionOuterClass.GetTransactionResponse>) responseObserver);
          break;
        case METHODID_CALCULATE_FEE:
          serviceImpl.calculateFee((pactus.transaction.TransactionOuterClass.CalculateFeeRequest) request,
              (io.grpc.stub.StreamObserver<pactus.transaction.TransactionOuterClass.CalculateFeeResponse>) responseObserver);
          break;
        case METHODID_BROADCAST_TRANSACTION:
          serviceImpl.broadcastTransaction((pactus.transaction.TransactionOuterClass.BroadcastTransactionRequest) request,
              (io.grpc.stub.StreamObserver<pactus.transaction.TransactionOuterClass.BroadcastTransactionResponse>) responseObserver);
          break;
        case METHODID_GET_RAW_TRANSACTION:
          serviceImpl.getRawTransaction((pactus.transaction.TransactionOuterClass.GetRawTransactionRequest) request,
              (io.grpc.stub.StreamObserver<pactus.transaction.TransactionOuterClass.GetRawTransactionResponse>) responseObserver);
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

  private static abstract class TransactionBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    TransactionBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return pactus.transaction.TransactionOuterClass.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("Transaction");
    }
  }

  private static final class TransactionFileDescriptorSupplier
      extends TransactionBaseDescriptorSupplier {
    TransactionFileDescriptorSupplier() {}
  }

  private static final class TransactionMethodDescriptorSupplier
      extends TransactionBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    TransactionMethodDescriptorSupplier(String methodName) {
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
      synchronized (TransactionGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new TransactionFileDescriptorSupplier())
              .addMethod(getGetTransactionMethod())
              .addMethod(getCalculateFeeMethod())
              .addMethod(getBroadcastTransactionMethod())
              .addMethod(getGetRawTransactionMethod())
              .build();
        }
      }
    }
    return result;
  }
}
