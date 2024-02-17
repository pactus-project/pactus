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

  private static volatile io.grpc.MethodDescriptor<pactus.transaction.TransactionOuterClass.GetRawTransferTransactionRequest,
      pactus.transaction.TransactionOuterClass.GetRawTransactionResponse> getGetRawTransferTransactionMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetRawTransferTransaction",
      requestType = pactus.transaction.TransactionOuterClass.GetRawTransferTransactionRequest.class,
      responseType = pactus.transaction.TransactionOuterClass.GetRawTransactionResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.transaction.TransactionOuterClass.GetRawTransferTransactionRequest,
      pactus.transaction.TransactionOuterClass.GetRawTransactionResponse> getGetRawTransferTransactionMethod() {
    io.grpc.MethodDescriptor<pactus.transaction.TransactionOuterClass.GetRawTransferTransactionRequest, pactus.transaction.TransactionOuterClass.GetRawTransactionResponse> getGetRawTransferTransactionMethod;
    if ((getGetRawTransferTransactionMethod = TransactionGrpc.getGetRawTransferTransactionMethod) == null) {
      synchronized (TransactionGrpc.class) {
        if ((getGetRawTransferTransactionMethod = TransactionGrpc.getGetRawTransferTransactionMethod) == null) {
          TransactionGrpc.getGetRawTransferTransactionMethod = getGetRawTransferTransactionMethod =
              io.grpc.MethodDescriptor.<pactus.transaction.TransactionOuterClass.GetRawTransferTransactionRequest, pactus.transaction.TransactionOuterClass.GetRawTransactionResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetRawTransferTransaction"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.transaction.TransactionOuterClass.GetRawTransferTransactionRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.transaction.TransactionOuterClass.GetRawTransactionResponse.getDefaultInstance()))
              .setSchemaDescriptor(new TransactionMethodDescriptorSupplier("GetRawTransferTransaction"))
              .build();
        }
      }
    }
    return getGetRawTransferTransactionMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.transaction.TransactionOuterClass.GetRawBondTransactionRequest,
      pactus.transaction.TransactionOuterClass.GetRawTransactionResponse> getGetRawBondTransactionMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetRawBondTransaction",
      requestType = pactus.transaction.TransactionOuterClass.GetRawBondTransactionRequest.class,
      responseType = pactus.transaction.TransactionOuterClass.GetRawTransactionResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.transaction.TransactionOuterClass.GetRawBondTransactionRequest,
      pactus.transaction.TransactionOuterClass.GetRawTransactionResponse> getGetRawBondTransactionMethod() {
    io.grpc.MethodDescriptor<pactus.transaction.TransactionOuterClass.GetRawBondTransactionRequest, pactus.transaction.TransactionOuterClass.GetRawTransactionResponse> getGetRawBondTransactionMethod;
    if ((getGetRawBondTransactionMethod = TransactionGrpc.getGetRawBondTransactionMethod) == null) {
      synchronized (TransactionGrpc.class) {
        if ((getGetRawBondTransactionMethod = TransactionGrpc.getGetRawBondTransactionMethod) == null) {
          TransactionGrpc.getGetRawBondTransactionMethod = getGetRawBondTransactionMethod =
              io.grpc.MethodDescriptor.<pactus.transaction.TransactionOuterClass.GetRawBondTransactionRequest, pactus.transaction.TransactionOuterClass.GetRawTransactionResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetRawBondTransaction"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.transaction.TransactionOuterClass.GetRawBondTransactionRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.transaction.TransactionOuterClass.GetRawTransactionResponse.getDefaultInstance()))
              .setSchemaDescriptor(new TransactionMethodDescriptorSupplier("GetRawBondTransaction"))
              .build();
        }
      }
    }
    return getGetRawBondTransactionMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.transaction.TransactionOuterClass.GetRawUnBondTransactionRequest,
      pactus.transaction.TransactionOuterClass.GetRawTransactionResponse> getGetRawUnBondTransactionMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetRawUnBondTransaction",
      requestType = pactus.transaction.TransactionOuterClass.GetRawUnBondTransactionRequest.class,
      responseType = pactus.transaction.TransactionOuterClass.GetRawTransactionResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.transaction.TransactionOuterClass.GetRawUnBondTransactionRequest,
      pactus.transaction.TransactionOuterClass.GetRawTransactionResponse> getGetRawUnBondTransactionMethod() {
    io.grpc.MethodDescriptor<pactus.transaction.TransactionOuterClass.GetRawUnBondTransactionRequest, pactus.transaction.TransactionOuterClass.GetRawTransactionResponse> getGetRawUnBondTransactionMethod;
    if ((getGetRawUnBondTransactionMethod = TransactionGrpc.getGetRawUnBondTransactionMethod) == null) {
      synchronized (TransactionGrpc.class) {
        if ((getGetRawUnBondTransactionMethod = TransactionGrpc.getGetRawUnBondTransactionMethod) == null) {
          TransactionGrpc.getGetRawUnBondTransactionMethod = getGetRawUnBondTransactionMethod =
              io.grpc.MethodDescriptor.<pactus.transaction.TransactionOuterClass.GetRawUnBondTransactionRequest, pactus.transaction.TransactionOuterClass.GetRawTransactionResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetRawUnBondTransaction"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.transaction.TransactionOuterClass.GetRawUnBondTransactionRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.transaction.TransactionOuterClass.GetRawTransactionResponse.getDefaultInstance()))
              .setSchemaDescriptor(new TransactionMethodDescriptorSupplier("GetRawUnBondTransaction"))
              .build();
        }
      }
    }
    return getGetRawUnBondTransactionMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.transaction.TransactionOuterClass.GetRawWithdrawTransactionRequest,
      pactus.transaction.TransactionOuterClass.GetRawTransactionResponse> getGetRawWithdrawTransactionMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetRawWithdrawTransaction",
      requestType = pactus.transaction.TransactionOuterClass.GetRawWithdrawTransactionRequest.class,
      responseType = pactus.transaction.TransactionOuterClass.GetRawTransactionResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.transaction.TransactionOuterClass.GetRawWithdrawTransactionRequest,
      pactus.transaction.TransactionOuterClass.GetRawTransactionResponse> getGetRawWithdrawTransactionMethod() {
    io.grpc.MethodDescriptor<pactus.transaction.TransactionOuterClass.GetRawWithdrawTransactionRequest, pactus.transaction.TransactionOuterClass.GetRawTransactionResponse> getGetRawWithdrawTransactionMethod;
    if ((getGetRawWithdrawTransactionMethod = TransactionGrpc.getGetRawWithdrawTransactionMethod) == null) {
      synchronized (TransactionGrpc.class) {
        if ((getGetRawWithdrawTransactionMethod = TransactionGrpc.getGetRawWithdrawTransactionMethod) == null) {
          TransactionGrpc.getGetRawWithdrawTransactionMethod = getGetRawWithdrawTransactionMethod =
              io.grpc.MethodDescriptor.<pactus.transaction.TransactionOuterClass.GetRawWithdrawTransactionRequest, pactus.transaction.TransactionOuterClass.GetRawTransactionResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetRawWithdrawTransaction"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.transaction.TransactionOuterClass.GetRawWithdrawTransactionRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.transaction.TransactionOuterClass.GetRawTransactionResponse.getDefaultInstance()))
              .setSchemaDescriptor(new TransactionMethodDescriptorSupplier("GetRawWithdrawTransaction"))
              .build();
        }
      }
    }
    return getGetRawWithdrawTransactionMethod;
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
     * GetRawTransferTransaction retrieves raw details of a transfer transaction.
     * </pre>
     */
    public void getRawTransferTransaction(pactus.transaction.TransactionOuterClass.GetRawTransferTransactionRequest request,
        io.grpc.stub.StreamObserver<pactus.transaction.TransactionOuterClass.GetRawTransactionResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetRawTransferTransactionMethod(), responseObserver);
    }

    /**
     * <pre>
     * GetRawBondTransaction retrieves raw details of a bond transaction.
     * </pre>
     */
    public void getRawBondTransaction(pactus.transaction.TransactionOuterClass.GetRawBondTransactionRequest request,
        io.grpc.stub.StreamObserver<pactus.transaction.TransactionOuterClass.GetRawTransactionResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetRawBondTransactionMethod(), responseObserver);
    }

    /**
     * <pre>
     * GetRawUnBondTransaction retrieves raw details of an unbond transaction.
     * </pre>
     */
    public void getRawUnBondTransaction(pactus.transaction.TransactionOuterClass.GetRawUnBondTransactionRequest request,
        io.grpc.stub.StreamObserver<pactus.transaction.TransactionOuterClass.GetRawTransactionResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetRawUnBondTransactionMethod(), responseObserver);
    }

    /**
     * <pre>
     * GetRawWithdrawTransaction retrieves raw details of a withdraw transaction.
     * </pre>
     */
    public void getRawWithdrawTransaction(pactus.transaction.TransactionOuterClass.GetRawWithdrawTransactionRequest request,
        io.grpc.stub.StreamObserver<pactus.transaction.TransactionOuterClass.GetRawTransactionResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetRawWithdrawTransactionMethod(), responseObserver);
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
            getGetRawTransferTransactionMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                pactus.transaction.TransactionOuterClass.GetRawTransferTransactionRequest,
                pactus.transaction.TransactionOuterClass.GetRawTransactionResponse>(
                  this, METHODID_GET_RAW_TRANSFER_TRANSACTION)))
          .addMethod(
            getGetRawBondTransactionMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                pactus.transaction.TransactionOuterClass.GetRawBondTransactionRequest,
                pactus.transaction.TransactionOuterClass.GetRawTransactionResponse>(
                  this, METHODID_GET_RAW_BOND_TRANSACTION)))
          .addMethod(
            getGetRawUnBondTransactionMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                pactus.transaction.TransactionOuterClass.GetRawUnBondTransactionRequest,
                pactus.transaction.TransactionOuterClass.GetRawTransactionResponse>(
                  this, METHODID_GET_RAW_UN_BOND_TRANSACTION)))
          .addMethod(
            getGetRawWithdrawTransactionMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                pactus.transaction.TransactionOuterClass.GetRawWithdrawTransactionRequest,
                pactus.transaction.TransactionOuterClass.GetRawTransactionResponse>(
                  this, METHODID_GET_RAW_WITHDRAW_TRANSACTION)))
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
     * GetRawTransferTransaction retrieves raw details of a transfer transaction.
     * </pre>
     */
    public void getRawTransferTransaction(pactus.transaction.TransactionOuterClass.GetRawTransferTransactionRequest request,
        io.grpc.stub.StreamObserver<pactus.transaction.TransactionOuterClass.GetRawTransactionResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetRawTransferTransactionMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * GetRawBondTransaction retrieves raw details of a bond transaction.
     * </pre>
     */
    public void getRawBondTransaction(pactus.transaction.TransactionOuterClass.GetRawBondTransactionRequest request,
        io.grpc.stub.StreamObserver<pactus.transaction.TransactionOuterClass.GetRawTransactionResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetRawBondTransactionMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * GetRawUnBondTransaction retrieves raw details of an unbond transaction.
     * </pre>
     */
    public void getRawUnBondTransaction(pactus.transaction.TransactionOuterClass.GetRawUnBondTransactionRequest request,
        io.grpc.stub.StreamObserver<pactus.transaction.TransactionOuterClass.GetRawTransactionResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetRawUnBondTransactionMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * GetRawWithdrawTransaction retrieves raw details of a withdraw transaction.
     * </pre>
     */
    public void getRawWithdrawTransaction(pactus.transaction.TransactionOuterClass.GetRawWithdrawTransactionRequest request,
        io.grpc.stub.StreamObserver<pactus.transaction.TransactionOuterClass.GetRawTransactionResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetRawWithdrawTransactionMethod(), getCallOptions()), request, responseObserver);
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
     * GetRawTransferTransaction retrieves raw details of a transfer transaction.
     * </pre>
     */
    public pactus.transaction.TransactionOuterClass.GetRawTransactionResponse getRawTransferTransaction(pactus.transaction.TransactionOuterClass.GetRawTransferTransactionRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetRawTransferTransactionMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetRawBondTransaction retrieves raw details of a bond transaction.
     * </pre>
     */
    public pactus.transaction.TransactionOuterClass.GetRawTransactionResponse getRawBondTransaction(pactus.transaction.TransactionOuterClass.GetRawBondTransactionRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetRawBondTransactionMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetRawUnBondTransaction retrieves raw details of an unbond transaction.
     * </pre>
     */
    public pactus.transaction.TransactionOuterClass.GetRawTransactionResponse getRawUnBondTransaction(pactus.transaction.TransactionOuterClass.GetRawUnBondTransactionRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetRawUnBondTransactionMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetRawWithdrawTransaction retrieves raw details of a withdraw transaction.
     * </pre>
     */
    public pactus.transaction.TransactionOuterClass.GetRawTransactionResponse getRawWithdrawTransaction(pactus.transaction.TransactionOuterClass.GetRawWithdrawTransactionRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetRawWithdrawTransactionMethod(), getCallOptions(), request);
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
     * GetRawTransferTransaction retrieves raw details of a transfer transaction.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.transaction.TransactionOuterClass.GetRawTransactionResponse> getRawTransferTransaction(
        pactus.transaction.TransactionOuterClass.GetRawTransferTransactionRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetRawTransferTransactionMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * GetRawBondTransaction retrieves raw details of a bond transaction.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.transaction.TransactionOuterClass.GetRawTransactionResponse> getRawBondTransaction(
        pactus.transaction.TransactionOuterClass.GetRawBondTransactionRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetRawBondTransactionMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * GetRawUnBondTransaction retrieves raw details of an unbond transaction.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.transaction.TransactionOuterClass.GetRawTransactionResponse> getRawUnBondTransaction(
        pactus.transaction.TransactionOuterClass.GetRawUnBondTransactionRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetRawUnBondTransactionMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * GetRawWithdrawTransaction retrieves raw details of a withdraw transaction.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.transaction.TransactionOuterClass.GetRawTransactionResponse> getRawWithdrawTransaction(
        pactus.transaction.TransactionOuterClass.GetRawWithdrawTransactionRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetRawWithdrawTransactionMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_GET_TRANSACTION = 0;
  private static final int METHODID_CALCULATE_FEE = 1;
  private static final int METHODID_BROADCAST_TRANSACTION = 2;
  private static final int METHODID_GET_RAW_TRANSFER_TRANSACTION = 3;
  private static final int METHODID_GET_RAW_BOND_TRANSACTION = 4;
  private static final int METHODID_GET_RAW_UN_BOND_TRANSACTION = 5;
  private static final int METHODID_GET_RAW_WITHDRAW_TRANSACTION = 6;

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
        case METHODID_GET_RAW_TRANSFER_TRANSACTION:
          serviceImpl.getRawTransferTransaction((pactus.transaction.TransactionOuterClass.GetRawTransferTransactionRequest) request,
              (io.grpc.stub.StreamObserver<pactus.transaction.TransactionOuterClass.GetRawTransactionResponse>) responseObserver);
          break;
        case METHODID_GET_RAW_BOND_TRANSACTION:
          serviceImpl.getRawBondTransaction((pactus.transaction.TransactionOuterClass.GetRawBondTransactionRequest) request,
              (io.grpc.stub.StreamObserver<pactus.transaction.TransactionOuterClass.GetRawTransactionResponse>) responseObserver);
          break;
        case METHODID_GET_RAW_UN_BOND_TRANSACTION:
          serviceImpl.getRawUnBondTransaction((pactus.transaction.TransactionOuterClass.GetRawUnBondTransactionRequest) request,
              (io.grpc.stub.StreamObserver<pactus.transaction.TransactionOuterClass.GetRawTransactionResponse>) responseObserver);
          break;
        case METHODID_GET_RAW_WITHDRAW_TRANSACTION:
          serviceImpl.getRawWithdrawTransaction((pactus.transaction.TransactionOuterClass.GetRawWithdrawTransactionRequest) request,
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
              .addMethod(getGetRawTransferTransactionMethod())
              .addMethod(getGetRawBondTransactionMethod())
              .addMethod(getGetRawUnBondTransactionMethod())
              .addMethod(getGetRawWithdrawTransactionMethod())
              .build();
        }
      }
    }
    return result;
  }
}
