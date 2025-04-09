package pactus;

import static io.grpc.MethodDescriptor.generateFullMethodName;

/**
 * <pre>
 * Transaction service defines various RPC methods for interacting with
 * transactions.
 * </pre>
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.71.0)",
    comments = "Source: transaction.proto")
@io.grpc.stub.annotations.GrpcGenerated
public final class TransactionServiceGrpc {

  private TransactionServiceGrpc() {}

  public static final java.lang.String SERVICE_NAME = "pactus.TransactionService";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<pactus.Transaction.GetTransactionRequest,
      pactus.Transaction.GetTransactionResponse> getGetTransactionMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetTransaction",
      requestType = pactus.Transaction.GetTransactionRequest.class,
      responseType = pactus.Transaction.GetTransactionResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.Transaction.GetTransactionRequest,
      pactus.Transaction.GetTransactionResponse> getGetTransactionMethod() {
    io.grpc.MethodDescriptor<pactus.Transaction.GetTransactionRequest, pactus.Transaction.GetTransactionResponse> getGetTransactionMethod;
    if ((getGetTransactionMethod = TransactionServiceGrpc.getGetTransactionMethod) == null) {
      synchronized (TransactionServiceGrpc.class) {
        if ((getGetTransactionMethod = TransactionServiceGrpc.getGetTransactionMethod) == null) {
          TransactionServiceGrpc.getGetTransactionMethod = getGetTransactionMethod =
              io.grpc.MethodDescriptor.<pactus.Transaction.GetTransactionRequest, pactus.Transaction.GetTransactionResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetTransaction"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Transaction.GetTransactionRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Transaction.GetTransactionResponse.getDefaultInstance()))
              .setSchemaDescriptor(new TransactionServiceMethodDescriptorSupplier("GetTransaction"))
              .build();
        }
      }
    }
    return getGetTransactionMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.Transaction.CalculateFeeRequest,
      pactus.Transaction.CalculateFeeResponse> getCalculateFeeMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "CalculateFee",
      requestType = pactus.Transaction.CalculateFeeRequest.class,
      responseType = pactus.Transaction.CalculateFeeResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.Transaction.CalculateFeeRequest,
      pactus.Transaction.CalculateFeeResponse> getCalculateFeeMethod() {
    io.grpc.MethodDescriptor<pactus.Transaction.CalculateFeeRequest, pactus.Transaction.CalculateFeeResponse> getCalculateFeeMethod;
    if ((getCalculateFeeMethod = TransactionServiceGrpc.getCalculateFeeMethod) == null) {
      synchronized (TransactionServiceGrpc.class) {
        if ((getCalculateFeeMethod = TransactionServiceGrpc.getCalculateFeeMethod) == null) {
          TransactionServiceGrpc.getCalculateFeeMethod = getCalculateFeeMethod =
              io.grpc.MethodDescriptor.<pactus.Transaction.CalculateFeeRequest, pactus.Transaction.CalculateFeeResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "CalculateFee"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Transaction.CalculateFeeRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Transaction.CalculateFeeResponse.getDefaultInstance()))
              .setSchemaDescriptor(new TransactionServiceMethodDescriptorSupplier("CalculateFee"))
              .build();
        }
      }
    }
    return getCalculateFeeMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.Transaction.BroadcastTransactionRequest,
      pactus.Transaction.BroadcastTransactionResponse> getBroadcastTransactionMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "BroadcastTransaction",
      requestType = pactus.Transaction.BroadcastTransactionRequest.class,
      responseType = pactus.Transaction.BroadcastTransactionResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.Transaction.BroadcastTransactionRequest,
      pactus.Transaction.BroadcastTransactionResponse> getBroadcastTransactionMethod() {
    io.grpc.MethodDescriptor<pactus.Transaction.BroadcastTransactionRequest, pactus.Transaction.BroadcastTransactionResponse> getBroadcastTransactionMethod;
    if ((getBroadcastTransactionMethod = TransactionServiceGrpc.getBroadcastTransactionMethod) == null) {
      synchronized (TransactionServiceGrpc.class) {
        if ((getBroadcastTransactionMethod = TransactionServiceGrpc.getBroadcastTransactionMethod) == null) {
          TransactionServiceGrpc.getBroadcastTransactionMethod = getBroadcastTransactionMethod =
              io.grpc.MethodDescriptor.<pactus.Transaction.BroadcastTransactionRequest, pactus.Transaction.BroadcastTransactionResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "BroadcastTransaction"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Transaction.BroadcastTransactionRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Transaction.BroadcastTransactionResponse.getDefaultInstance()))
              .setSchemaDescriptor(new TransactionServiceMethodDescriptorSupplier("BroadcastTransaction"))
              .build();
        }
      }
    }
    return getBroadcastTransactionMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.Transaction.GetRawTransferTransactionRequest,
      pactus.Transaction.GetRawTransactionResponse> getGetRawTransferTransactionMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetRawTransferTransaction",
      requestType = pactus.Transaction.GetRawTransferTransactionRequest.class,
      responseType = pactus.Transaction.GetRawTransactionResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.Transaction.GetRawTransferTransactionRequest,
      pactus.Transaction.GetRawTransactionResponse> getGetRawTransferTransactionMethod() {
    io.grpc.MethodDescriptor<pactus.Transaction.GetRawTransferTransactionRequest, pactus.Transaction.GetRawTransactionResponse> getGetRawTransferTransactionMethod;
    if ((getGetRawTransferTransactionMethod = TransactionServiceGrpc.getGetRawTransferTransactionMethod) == null) {
      synchronized (TransactionServiceGrpc.class) {
        if ((getGetRawTransferTransactionMethod = TransactionServiceGrpc.getGetRawTransferTransactionMethod) == null) {
          TransactionServiceGrpc.getGetRawTransferTransactionMethod = getGetRawTransferTransactionMethod =
              io.grpc.MethodDescriptor.<pactus.Transaction.GetRawTransferTransactionRequest, pactus.Transaction.GetRawTransactionResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetRawTransferTransaction"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Transaction.GetRawTransferTransactionRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Transaction.GetRawTransactionResponse.getDefaultInstance()))
              .setSchemaDescriptor(new TransactionServiceMethodDescriptorSupplier("GetRawTransferTransaction"))
              .build();
        }
      }
    }
    return getGetRawTransferTransactionMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.Transaction.GetRawBondTransactionRequest,
      pactus.Transaction.GetRawTransactionResponse> getGetRawBondTransactionMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetRawBondTransaction",
      requestType = pactus.Transaction.GetRawBondTransactionRequest.class,
      responseType = pactus.Transaction.GetRawTransactionResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.Transaction.GetRawBondTransactionRequest,
      pactus.Transaction.GetRawTransactionResponse> getGetRawBondTransactionMethod() {
    io.grpc.MethodDescriptor<pactus.Transaction.GetRawBondTransactionRequest, pactus.Transaction.GetRawTransactionResponse> getGetRawBondTransactionMethod;
    if ((getGetRawBondTransactionMethod = TransactionServiceGrpc.getGetRawBondTransactionMethod) == null) {
      synchronized (TransactionServiceGrpc.class) {
        if ((getGetRawBondTransactionMethod = TransactionServiceGrpc.getGetRawBondTransactionMethod) == null) {
          TransactionServiceGrpc.getGetRawBondTransactionMethod = getGetRawBondTransactionMethod =
              io.grpc.MethodDescriptor.<pactus.Transaction.GetRawBondTransactionRequest, pactus.Transaction.GetRawTransactionResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetRawBondTransaction"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Transaction.GetRawBondTransactionRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Transaction.GetRawTransactionResponse.getDefaultInstance()))
              .setSchemaDescriptor(new TransactionServiceMethodDescriptorSupplier("GetRawBondTransaction"))
              .build();
        }
      }
    }
    return getGetRawBondTransactionMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.Transaction.GetRawUnbondTransactionRequest,
      pactus.Transaction.GetRawTransactionResponse> getGetRawUnbondTransactionMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetRawUnbondTransaction",
      requestType = pactus.Transaction.GetRawUnbondTransactionRequest.class,
      responseType = pactus.Transaction.GetRawTransactionResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.Transaction.GetRawUnbondTransactionRequest,
      pactus.Transaction.GetRawTransactionResponse> getGetRawUnbondTransactionMethod() {
    io.grpc.MethodDescriptor<pactus.Transaction.GetRawUnbondTransactionRequest, pactus.Transaction.GetRawTransactionResponse> getGetRawUnbondTransactionMethod;
    if ((getGetRawUnbondTransactionMethod = TransactionServiceGrpc.getGetRawUnbondTransactionMethod) == null) {
      synchronized (TransactionServiceGrpc.class) {
        if ((getGetRawUnbondTransactionMethod = TransactionServiceGrpc.getGetRawUnbondTransactionMethod) == null) {
          TransactionServiceGrpc.getGetRawUnbondTransactionMethod = getGetRawUnbondTransactionMethod =
              io.grpc.MethodDescriptor.<pactus.Transaction.GetRawUnbondTransactionRequest, pactus.Transaction.GetRawTransactionResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetRawUnbondTransaction"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Transaction.GetRawUnbondTransactionRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Transaction.GetRawTransactionResponse.getDefaultInstance()))
              .setSchemaDescriptor(new TransactionServiceMethodDescriptorSupplier("GetRawUnbondTransaction"))
              .build();
        }
      }
    }
    return getGetRawUnbondTransactionMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.Transaction.GetRawWithdrawTransactionRequest,
      pactus.Transaction.GetRawTransactionResponse> getGetRawWithdrawTransactionMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetRawWithdrawTransaction",
      requestType = pactus.Transaction.GetRawWithdrawTransactionRequest.class,
      responseType = pactus.Transaction.GetRawTransactionResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.Transaction.GetRawWithdrawTransactionRequest,
      pactus.Transaction.GetRawTransactionResponse> getGetRawWithdrawTransactionMethod() {
    io.grpc.MethodDescriptor<pactus.Transaction.GetRawWithdrawTransactionRequest, pactus.Transaction.GetRawTransactionResponse> getGetRawWithdrawTransactionMethod;
    if ((getGetRawWithdrawTransactionMethod = TransactionServiceGrpc.getGetRawWithdrawTransactionMethod) == null) {
      synchronized (TransactionServiceGrpc.class) {
        if ((getGetRawWithdrawTransactionMethod = TransactionServiceGrpc.getGetRawWithdrawTransactionMethod) == null) {
          TransactionServiceGrpc.getGetRawWithdrawTransactionMethod = getGetRawWithdrawTransactionMethod =
              io.grpc.MethodDescriptor.<pactus.Transaction.GetRawWithdrawTransactionRequest, pactus.Transaction.GetRawTransactionResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetRawWithdrawTransaction"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Transaction.GetRawWithdrawTransactionRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Transaction.GetRawTransactionResponse.getDefaultInstance()))
              .setSchemaDescriptor(new TransactionServiceMethodDescriptorSupplier("GetRawWithdrawTransaction"))
              .build();
        }
      }
    }
    return getGetRawWithdrawTransactionMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.Transaction.DecodeRawTransactionRequest,
      pactus.Transaction.DecodeRawTransactionResponse> getDecodeRawTransactionMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "DecodeRawTransaction",
      requestType = pactus.Transaction.DecodeRawTransactionRequest.class,
      responseType = pactus.Transaction.DecodeRawTransactionResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.Transaction.DecodeRawTransactionRequest,
      pactus.Transaction.DecodeRawTransactionResponse> getDecodeRawTransactionMethod() {
    io.grpc.MethodDescriptor<pactus.Transaction.DecodeRawTransactionRequest, pactus.Transaction.DecodeRawTransactionResponse> getDecodeRawTransactionMethod;
    if ((getDecodeRawTransactionMethod = TransactionServiceGrpc.getDecodeRawTransactionMethod) == null) {
      synchronized (TransactionServiceGrpc.class) {
        if ((getDecodeRawTransactionMethod = TransactionServiceGrpc.getDecodeRawTransactionMethod) == null) {
          TransactionServiceGrpc.getDecodeRawTransactionMethod = getDecodeRawTransactionMethod =
              io.grpc.MethodDescriptor.<pactus.Transaction.DecodeRawTransactionRequest, pactus.Transaction.DecodeRawTransactionResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "DecodeRawTransaction"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Transaction.DecodeRawTransactionRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Transaction.DecodeRawTransactionResponse.getDefaultInstance()))
              .setSchemaDescriptor(new TransactionServiceMethodDescriptorSupplier("DecodeRawTransaction"))
              .build();
        }
      }
    }
    return getDecodeRawTransactionMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static TransactionServiceStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<TransactionServiceStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<TransactionServiceStub>() {
        @java.lang.Override
        public TransactionServiceStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new TransactionServiceStub(channel, callOptions);
        }
      };
    return TransactionServiceStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports all types of calls on the service
   */
  public static TransactionServiceBlockingV2Stub newBlockingV2Stub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<TransactionServiceBlockingV2Stub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<TransactionServiceBlockingV2Stub>() {
        @java.lang.Override
        public TransactionServiceBlockingV2Stub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new TransactionServiceBlockingV2Stub(channel, callOptions);
        }
      };
    return TransactionServiceBlockingV2Stub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static TransactionServiceBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<TransactionServiceBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<TransactionServiceBlockingStub>() {
        @java.lang.Override
        public TransactionServiceBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new TransactionServiceBlockingStub(channel, callOptions);
        }
      };
    return TransactionServiceBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static TransactionServiceFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<TransactionServiceFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<TransactionServiceFutureStub>() {
        @java.lang.Override
        public TransactionServiceFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new TransactionServiceFutureStub(channel, callOptions);
        }
      };
    return TransactionServiceFutureStub.newStub(factory, channel);
  }

  /**
   * <pre>
   * Transaction service defines various RPC methods for interacting with
   * transactions.
   * </pre>
   */
  public interface AsyncService {

    /**
     * <pre>
     * GetTransaction retrieves transaction details based on the provided request
     * parameters.
     * </pre>
     */
    default void getTransaction(pactus.Transaction.GetTransactionRequest request,
        io.grpc.stub.StreamObserver<pactus.Transaction.GetTransactionResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetTransactionMethod(), responseObserver);
    }

    /**
     * <pre>
     * CalculateFee calculates the transaction fee based on the specified amount
     * and payload type.
     * </pre>
     */
    default void calculateFee(pactus.Transaction.CalculateFeeRequest request,
        io.grpc.stub.StreamObserver<pactus.Transaction.CalculateFeeResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getCalculateFeeMethod(), responseObserver);
    }

    /**
     * <pre>
     * BroadcastTransaction broadcasts a signed transaction to the network.
     * </pre>
     */
    default void broadcastTransaction(pactus.Transaction.BroadcastTransactionRequest request,
        io.grpc.stub.StreamObserver<pactus.Transaction.BroadcastTransactionResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getBroadcastTransactionMethod(), responseObserver);
    }

    /**
     * <pre>
     * GetRawTransferTransaction retrieves raw details of a transfer transaction.
     * buf:lint:ignore RPC_REQUEST_RESPONSE_UNIQUE
     * buf:lint:ignore RPC_RESPONSE_STANDARD_NAME
     * </pre>
     */
    default void getRawTransferTransaction(pactus.Transaction.GetRawTransferTransactionRequest request,
        io.grpc.stub.StreamObserver<pactus.Transaction.GetRawTransactionResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetRawTransferTransactionMethod(), responseObserver);
    }

    /**
     * <pre>
     * GetRawBondTransaction retrieves raw details of a bond transaction.
     * buf:lint:ignore RPC_REQUEST_RESPONSE_UNIQUE
     * buf:lint:ignore RPC_RESPONSE_STANDARD_NAME
     * </pre>
     */
    default void getRawBondTransaction(pactus.Transaction.GetRawBondTransactionRequest request,
        io.grpc.stub.StreamObserver<pactus.Transaction.GetRawTransactionResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetRawBondTransactionMethod(), responseObserver);
    }

    /**
     * <pre>
     * GetRawUnbondTransaction retrieves raw details of an unbond transaction.
     * buf:lint:ignore RPC_REQUEST_RESPONSE_UNIQUE
     * buf:lint:ignore RPC_RESPONSE_STANDARD_NAME
     * </pre>
     */
    default void getRawUnbondTransaction(pactus.Transaction.GetRawUnbondTransactionRequest request,
        io.grpc.stub.StreamObserver<pactus.Transaction.GetRawTransactionResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetRawUnbondTransactionMethod(), responseObserver);
    }

    /**
     * <pre>
     * GetRawWithdrawTransaction retrieves raw details of a withdraw transaction.
     * buf:lint:ignore RPC_REQUEST_RESPONSE_UNIQUE
     * buf:lint:ignore RPC_RESPONSE_STANDARD_NAME
     * </pre>
     */
    default void getRawWithdrawTransaction(pactus.Transaction.GetRawWithdrawTransactionRequest request,
        io.grpc.stub.StreamObserver<pactus.Transaction.GetRawTransactionResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetRawWithdrawTransactionMethod(), responseObserver);
    }

    /**
     * <pre>
     * DecodeRawTransaction accepts raw transaction and returns decoded transaction.
     * </pre>
     */
    default void decodeRawTransaction(pactus.Transaction.DecodeRawTransactionRequest request,
        io.grpc.stub.StreamObserver<pactus.Transaction.DecodeRawTransactionResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getDecodeRawTransactionMethod(), responseObserver);
    }
  }

  /**
   * Base class for the server implementation of the service TransactionService.
   * <pre>
   * Transaction service defines various RPC methods for interacting with
   * transactions.
   * </pre>
   */
  public static abstract class TransactionServiceImplBase
      implements io.grpc.BindableService, AsyncService {

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return TransactionServiceGrpc.bindService(this);
    }
  }

  /**
   * A stub to allow clients to do asynchronous rpc calls to service TransactionService.
   * <pre>
   * Transaction service defines various RPC methods for interacting with
   * transactions.
   * </pre>
   */
  public static final class TransactionServiceStub
      extends io.grpc.stub.AbstractAsyncStub<TransactionServiceStub> {
    private TransactionServiceStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected TransactionServiceStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new TransactionServiceStub(channel, callOptions);
    }

    /**
     * <pre>
     * GetTransaction retrieves transaction details based on the provided request
     * parameters.
     * </pre>
     */
    public void getTransaction(pactus.Transaction.GetTransactionRequest request,
        io.grpc.stub.StreamObserver<pactus.Transaction.GetTransactionResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetTransactionMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * CalculateFee calculates the transaction fee based on the specified amount
     * and payload type.
     * </pre>
     */
    public void calculateFee(pactus.Transaction.CalculateFeeRequest request,
        io.grpc.stub.StreamObserver<pactus.Transaction.CalculateFeeResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getCalculateFeeMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * BroadcastTransaction broadcasts a signed transaction to the network.
     * </pre>
     */
    public void broadcastTransaction(pactus.Transaction.BroadcastTransactionRequest request,
        io.grpc.stub.StreamObserver<pactus.Transaction.BroadcastTransactionResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getBroadcastTransactionMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * GetRawTransferTransaction retrieves raw details of a transfer transaction.
     * buf:lint:ignore RPC_REQUEST_RESPONSE_UNIQUE
     * buf:lint:ignore RPC_RESPONSE_STANDARD_NAME
     * </pre>
     */
    public void getRawTransferTransaction(pactus.Transaction.GetRawTransferTransactionRequest request,
        io.grpc.stub.StreamObserver<pactus.Transaction.GetRawTransactionResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetRawTransferTransactionMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * GetRawBondTransaction retrieves raw details of a bond transaction.
     * buf:lint:ignore RPC_REQUEST_RESPONSE_UNIQUE
     * buf:lint:ignore RPC_RESPONSE_STANDARD_NAME
     * </pre>
     */
    public void getRawBondTransaction(pactus.Transaction.GetRawBondTransactionRequest request,
        io.grpc.stub.StreamObserver<pactus.Transaction.GetRawTransactionResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetRawBondTransactionMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * GetRawUnbondTransaction retrieves raw details of an unbond transaction.
     * buf:lint:ignore RPC_REQUEST_RESPONSE_UNIQUE
     * buf:lint:ignore RPC_RESPONSE_STANDARD_NAME
     * </pre>
     */
    public void getRawUnbondTransaction(pactus.Transaction.GetRawUnbondTransactionRequest request,
        io.grpc.stub.StreamObserver<pactus.Transaction.GetRawTransactionResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetRawUnbondTransactionMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * GetRawWithdrawTransaction retrieves raw details of a withdraw transaction.
     * buf:lint:ignore RPC_REQUEST_RESPONSE_UNIQUE
     * buf:lint:ignore RPC_RESPONSE_STANDARD_NAME
     * </pre>
     */
    public void getRawWithdrawTransaction(pactus.Transaction.GetRawWithdrawTransactionRequest request,
        io.grpc.stub.StreamObserver<pactus.Transaction.GetRawTransactionResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetRawWithdrawTransactionMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * DecodeRawTransaction accepts raw transaction and returns decoded transaction.
     * </pre>
     */
    public void decodeRawTransaction(pactus.Transaction.DecodeRawTransactionRequest request,
        io.grpc.stub.StreamObserver<pactus.Transaction.DecodeRawTransactionResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getDecodeRawTransactionMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   * A stub to allow clients to do synchronous rpc calls to service TransactionService.
   * <pre>
   * Transaction service defines various RPC methods for interacting with
   * transactions.
   * </pre>
   */
  public static final class TransactionServiceBlockingV2Stub
      extends io.grpc.stub.AbstractBlockingStub<TransactionServiceBlockingV2Stub> {
    private TransactionServiceBlockingV2Stub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected TransactionServiceBlockingV2Stub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new TransactionServiceBlockingV2Stub(channel, callOptions);
    }

    /**
     * <pre>
     * GetTransaction retrieves transaction details based on the provided request
     * parameters.
     * </pre>
     */
    public pactus.Transaction.GetTransactionResponse getTransaction(pactus.Transaction.GetTransactionRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetTransactionMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * CalculateFee calculates the transaction fee based on the specified amount
     * and payload type.
     * </pre>
     */
    public pactus.Transaction.CalculateFeeResponse calculateFee(pactus.Transaction.CalculateFeeRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getCalculateFeeMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * BroadcastTransaction broadcasts a signed transaction to the network.
     * </pre>
     */
    public pactus.Transaction.BroadcastTransactionResponse broadcastTransaction(pactus.Transaction.BroadcastTransactionRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getBroadcastTransactionMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetRawTransferTransaction retrieves raw details of a transfer transaction.
     * buf:lint:ignore RPC_REQUEST_RESPONSE_UNIQUE
     * buf:lint:ignore RPC_RESPONSE_STANDARD_NAME
     * </pre>
     */
    public pactus.Transaction.GetRawTransactionResponse getRawTransferTransaction(pactus.Transaction.GetRawTransferTransactionRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetRawTransferTransactionMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetRawBondTransaction retrieves raw details of a bond transaction.
     * buf:lint:ignore RPC_REQUEST_RESPONSE_UNIQUE
     * buf:lint:ignore RPC_RESPONSE_STANDARD_NAME
     * </pre>
     */
    public pactus.Transaction.GetRawTransactionResponse getRawBondTransaction(pactus.Transaction.GetRawBondTransactionRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetRawBondTransactionMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetRawUnbondTransaction retrieves raw details of an unbond transaction.
     * buf:lint:ignore RPC_REQUEST_RESPONSE_UNIQUE
     * buf:lint:ignore RPC_RESPONSE_STANDARD_NAME
     * </pre>
     */
    public pactus.Transaction.GetRawTransactionResponse getRawUnbondTransaction(pactus.Transaction.GetRawUnbondTransactionRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetRawUnbondTransactionMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetRawWithdrawTransaction retrieves raw details of a withdraw transaction.
     * buf:lint:ignore RPC_REQUEST_RESPONSE_UNIQUE
     * buf:lint:ignore RPC_RESPONSE_STANDARD_NAME
     * </pre>
     */
    public pactus.Transaction.GetRawTransactionResponse getRawWithdrawTransaction(pactus.Transaction.GetRawWithdrawTransactionRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetRawWithdrawTransactionMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * DecodeRawTransaction accepts raw transaction and returns decoded transaction.
     * </pre>
     */
    public pactus.Transaction.DecodeRawTransactionResponse decodeRawTransaction(pactus.Transaction.DecodeRawTransactionRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getDecodeRawTransactionMethod(), getCallOptions(), request);
    }
  }

  /**
   * A stub to allow clients to do limited synchronous rpc calls to service TransactionService.
   * <pre>
   * Transaction service defines various RPC methods for interacting with
   * transactions.
   * </pre>
   */
  public static final class TransactionServiceBlockingStub
      extends io.grpc.stub.AbstractBlockingStub<TransactionServiceBlockingStub> {
    private TransactionServiceBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected TransactionServiceBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new TransactionServiceBlockingStub(channel, callOptions);
    }

    /**
     * <pre>
     * GetTransaction retrieves transaction details based on the provided request
     * parameters.
     * </pre>
     */
    public pactus.Transaction.GetTransactionResponse getTransaction(pactus.Transaction.GetTransactionRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetTransactionMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * CalculateFee calculates the transaction fee based on the specified amount
     * and payload type.
     * </pre>
     */
    public pactus.Transaction.CalculateFeeResponse calculateFee(pactus.Transaction.CalculateFeeRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getCalculateFeeMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * BroadcastTransaction broadcasts a signed transaction to the network.
     * </pre>
     */
    public pactus.Transaction.BroadcastTransactionResponse broadcastTransaction(pactus.Transaction.BroadcastTransactionRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getBroadcastTransactionMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetRawTransferTransaction retrieves raw details of a transfer transaction.
     * buf:lint:ignore RPC_REQUEST_RESPONSE_UNIQUE
     * buf:lint:ignore RPC_RESPONSE_STANDARD_NAME
     * </pre>
     */
    public pactus.Transaction.GetRawTransactionResponse getRawTransferTransaction(pactus.Transaction.GetRawTransferTransactionRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetRawTransferTransactionMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetRawBondTransaction retrieves raw details of a bond transaction.
     * buf:lint:ignore RPC_REQUEST_RESPONSE_UNIQUE
     * buf:lint:ignore RPC_RESPONSE_STANDARD_NAME
     * </pre>
     */
    public pactus.Transaction.GetRawTransactionResponse getRawBondTransaction(pactus.Transaction.GetRawBondTransactionRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetRawBondTransactionMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetRawUnbondTransaction retrieves raw details of an unbond transaction.
     * buf:lint:ignore RPC_REQUEST_RESPONSE_UNIQUE
     * buf:lint:ignore RPC_RESPONSE_STANDARD_NAME
     * </pre>
     */
    public pactus.Transaction.GetRawTransactionResponse getRawUnbondTransaction(pactus.Transaction.GetRawUnbondTransactionRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetRawUnbondTransactionMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetRawWithdrawTransaction retrieves raw details of a withdraw transaction.
     * buf:lint:ignore RPC_REQUEST_RESPONSE_UNIQUE
     * buf:lint:ignore RPC_RESPONSE_STANDARD_NAME
     * </pre>
     */
    public pactus.Transaction.GetRawTransactionResponse getRawWithdrawTransaction(pactus.Transaction.GetRawWithdrawTransactionRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetRawWithdrawTransactionMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * DecodeRawTransaction accepts raw transaction and returns decoded transaction.
     * </pre>
     */
    public pactus.Transaction.DecodeRawTransactionResponse decodeRawTransaction(pactus.Transaction.DecodeRawTransactionRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getDecodeRawTransactionMethod(), getCallOptions(), request);
    }
  }

  /**
   * A stub to allow clients to do ListenableFuture-style rpc calls to service TransactionService.
   * <pre>
   * Transaction service defines various RPC methods for interacting with
   * transactions.
   * </pre>
   */
  public static final class TransactionServiceFutureStub
      extends io.grpc.stub.AbstractFutureStub<TransactionServiceFutureStub> {
    private TransactionServiceFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected TransactionServiceFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new TransactionServiceFutureStub(channel, callOptions);
    }

    /**
     * <pre>
     * GetTransaction retrieves transaction details based on the provided request
     * parameters.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.Transaction.GetTransactionResponse> getTransaction(
        pactus.Transaction.GetTransactionRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetTransactionMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * CalculateFee calculates the transaction fee based on the specified amount
     * and payload type.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.Transaction.CalculateFeeResponse> calculateFee(
        pactus.Transaction.CalculateFeeRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getCalculateFeeMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * BroadcastTransaction broadcasts a signed transaction to the network.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.Transaction.BroadcastTransactionResponse> broadcastTransaction(
        pactus.Transaction.BroadcastTransactionRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getBroadcastTransactionMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * GetRawTransferTransaction retrieves raw details of a transfer transaction.
     * buf:lint:ignore RPC_REQUEST_RESPONSE_UNIQUE
     * buf:lint:ignore RPC_RESPONSE_STANDARD_NAME
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.Transaction.GetRawTransactionResponse> getRawTransferTransaction(
        pactus.Transaction.GetRawTransferTransactionRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetRawTransferTransactionMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * GetRawBondTransaction retrieves raw details of a bond transaction.
     * buf:lint:ignore RPC_REQUEST_RESPONSE_UNIQUE
     * buf:lint:ignore RPC_RESPONSE_STANDARD_NAME
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.Transaction.GetRawTransactionResponse> getRawBondTransaction(
        pactus.Transaction.GetRawBondTransactionRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetRawBondTransactionMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * GetRawUnbondTransaction retrieves raw details of an unbond transaction.
     * buf:lint:ignore RPC_REQUEST_RESPONSE_UNIQUE
     * buf:lint:ignore RPC_RESPONSE_STANDARD_NAME
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.Transaction.GetRawTransactionResponse> getRawUnbondTransaction(
        pactus.Transaction.GetRawUnbondTransactionRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetRawUnbondTransactionMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * GetRawWithdrawTransaction retrieves raw details of a withdraw transaction.
     * buf:lint:ignore RPC_REQUEST_RESPONSE_UNIQUE
     * buf:lint:ignore RPC_RESPONSE_STANDARD_NAME
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.Transaction.GetRawTransactionResponse> getRawWithdrawTransaction(
        pactus.Transaction.GetRawWithdrawTransactionRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetRawWithdrawTransactionMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * DecodeRawTransaction accepts raw transaction and returns decoded transaction.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.Transaction.DecodeRawTransactionResponse> decodeRawTransaction(
        pactus.Transaction.DecodeRawTransactionRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getDecodeRawTransactionMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_GET_TRANSACTION = 0;
  private static final int METHODID_CALCULATE_FEE = 1;
  private static final int METHODID_BROADCAST_TRANSACTION = 2;
  private static final int METHODID_GET_RAW_TRANSFER_TRANSACTION = 3;
  private static final int METHODID_GET_RAW_BOND_TRANSACTION = 4;
  private static final int METHODID_GET_RAW_UNBOND_TRANSACTION = 5;
  private static final int METHODID_GET_RAW_WITHDRAW_TRANSACTION = 6;
  private static final int METHODID_DECODE_RAW_TRANSACTION = 7;

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
        case METHODID_GET_TRANSACTION:
          serviceImpl.getTransaction((pactus.Transaction.GetTransactionRequest) request,
              (io.grpc.stub.StreamObserver<pactus.Transaction.GetTransactionResponse>) responseObserver);
          break;
        case METHODID_CALCULATE_FEE:
          serviceImpl.calculateFee((pactus.Transaction.CalculateFeeRequest) request,
              (io.grpc.stub.StreamObserver<pactus.Transaction.CalculateFeeResponse>) responseObserver);
          break;
        case METHODID_BROADCAST_TRANSACTION:
          serviceImpl.broadcastTransaction((pactus.Transaction.BroadcastTransactionRequest) request,
              (io.grpc.stub.StreamObserver<pactus.Transaction.BroadcastTransactionResponse>) responseObserver);
          break;
        case METHODID_GET_RAW_TRANSFER_TRANSACTION:
          serviceImpl.getRawTransferTransaction((pactus.Transaction.GetRawTransferTransactionRequest) request,
              (io.grpc.stub.StreamObserver<pactus.Transaction.GetRawTransactionResponse>) responseObserver);
          break;
        case METHODID_GET_RAW_BOND_TRANSACTION:
          serviceImpl.getRawBondTransaction((pactus.Transaction.GetRawBondTransactionRequest) request,
              (io.grpc.stub.StreamObserver<pactus.Transaction.GetRawTransactionResponse>) responseObserver);
          break;
        case METHODID_GET_RAW_UNBOND_TRANSACTION:
          serviceImpl.getRawUnbondTransaction((pactus.Transaction.GetRawUnbondTransactionRequest) request,
              (io.grpc.stub.StreamObserver<pactus.Transaction.GetRawTransactionResponse>) responseObserver);
          break;
        case METHODID_GET_RAW_WITHDRAW_TRANSACTION:
          serviceImpl.getRawWithdrawTransaction((pactus.Transaction.GetRawWithdrawTransactionRequest) request,
              (io.grpc.stub.StreamObserver<pactus.Transaction.GetRawTransactionResponse>) responseObserver);
          break;
        case METHODID_DECODE_RAW_TRANSACTION:
          serviceImpl.decodeRawTransaction((pactus.Transaction.DecodeRawTransactionRequest) request,
              (io.grpc.stub.StreamObserver<pactus.Transaction.DecodeRawTransactionResponse>) responseObserver);
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
          getGetTransactionMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.Transaction.GetTransactionRequest,
              pactus.Transaction.GetTransactionResponse>(
                service, METHODID_GET_TRANSACTION)))
        .addMethod(
          getCalculateFeeMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.Transaction.CalculateFeeRequest,
              pactus.Transaction.CalculateFeeResponse>(
                service, METHODID_CALCULATE_FEE)))
        .addMethod(
          getBroadcastTransactionMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.Transaction.BroadcastTransactionRequest,
              pactus.Transaction.BroadcastTransactionResponse>(
                service, METHODID_BROADCAST_TRANSACTION)))
        .addMethod(
          getGetRawTransferTransactionMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.Transaction.GetRawTransferTransactionRequest,
              pactus.Transaction.GetRawTransactionResponse>(
                service, METHODID_GET_RAW_TRANSFER_TRANSACTION)))
        .addMethod(
          getGetRawBondTransactionMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.Transaction.GetRawBondTransactionRequest,
              pactus.Transaction.GetRawTransactionResponse>(
                service, METHODID_GET_RAW_BOND_TRANSACTION)))
        .addMethod(
          getGetRawUnbondTransactionMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.Transaction.GetRawUnbondTransactionRequest,
              pactus.Transaction.GetRawTransactionResponse>(
                service, METHODID_GET_RAW_UNBOND_TRANSACTION)))
        .addMethod(
          getGetRawWithdrawTransactionMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.Transaction.GetRawWithdrawTransactionRequest,
              pactus.Transaction.GetRawTransactionResponse>(
                service, METHODID_GET_RAW_WITHDRAW_TRANSACTION)))
        .addMethod(
          getDecodeRawTransactionMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.Transaction.DecodeRawTransactionRequest,
              pactus.Transaction.DecodeRawTransactionResponse>(
                service, METHODID_DECODE_RAW_TRANSACTION)))
        .build();
  }

  private static abstract class TransactionServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    TransactionServiceBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return pactus.Transaction.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("TransactionService");
    }
  }

  private static final class TransactionServiceFileDescriptorSupplier
      extends TransactionServiceBaseDescriptorSupplier {
    TransactionServiceFileDescriptorSupplier() {}
  }

  private static final class TransactionServiceMethodDescriptorSupplier
      extends TransactionServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final java.lang.String methodName;

    TransactionServiceMethodDescriptorSupplier(java.lang.String methodName) {
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
      synchronized (TransactionServiceGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new TransactionServiceFileDescriptorSupplier())
              .addMethod(getGetTransactionMethod())
              .addMethod(getCalculateFeeMethod())
              .addMethod(getBroadcastTransactionMethod())
              .addMethod(getGetRawTransferTransactionMethod())
              .addMethod(getGetRawBondTransactionMethod())
              .addMethod(getGetRawUnbondTransactionMethod())
              .addMethod(getGetRawWithdrawTransactionMethod())
              .addMethod(getDecodeRawTransactionMethod())
              .build();
        }
      }
    }
    return result;
  }
}
