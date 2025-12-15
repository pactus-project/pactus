package pactus;

import static io.grpc.MethodDescriptor.generateFullMethodName;

/**
 * <pre>
 * Transaction service defines various RPC methods for interacting with transactions.
 * </pre>
 */
@io.grpc.stub.annotations.GrpcGenerated
public final class TransactionGrpc {

  private TransactionGrpc() {}

  public static final java.lang.String SERVICE_NAME = "pactus.Transaction";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<pactus.TransactionOuterClass.GetTransactionRequest,
      pactus.TransactionOuterClass.GetTransactionResponse> getGetTransactionMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetTransaction",
      requestType = pactus.TransactionOuterClass.GetTransactionRequest.class,
      responseType = pactus.TransactionOuterClass.GetTransactionResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.TransactionOuterClass.GetTransactionRequest,
      pactus.TransactionOuterClass.GetTransactionResponse> getGetTransactionMethod() {
    io.grpc.MethodDescriptor<pactus.TransactionOuterClass.GetTransactionRequest, pactus.TransactionOuterClass.GetTransactionResponse> getGetTransactionMethod;
    if ((getGetTransactionMethod = TransactionGrpc.getGetTransactionMethod) == null) {
      synchronized (TransactionGrpc.class) {
        if ((getGetTransactionMethod = TransactionGrpc.getGetTransactionMethod) == null) {
          TransactionGrpc.getGetTransactionMethod = getGetTransactionMethod =
              io.grpc.MethodDescriptor.<pactus.TransactionOuterClass.GetTransactionRequest, pactus.TransactionOuterClass.GetTransactionResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetTransaction"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.TransactionOuterClass.GetTransactionRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.TransactionOuterClass.GetTransactionResponse.getDefaultInstance()))
              .setSchemaDescriptor(new TransactionMethodDescriptorSupplier("GetTransaction"))
              .build();
        }
      }
    }
    return getGetTransactionMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.TransactionOuterClass.CalculateFeeRequest,
      pactus.TransactionOuterClass.CalculateFeeResponse> getCalculateFeeMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "CalculateFee",
      requestType = pactus.TransactionOuterClass.CalculateFeeRequest.class,
      responseType = pactus.TransactionOuterClass.CalculateFeeResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.TransactionOuterClass.CalculateFeeRequest,
      pactus.TransactionOuterClass.CalculateFeeResponse> getCalculateFeeMethod() {
    io.grpc.MethodDescriptor<pactus.TransactionOuterClass.CalculateFeeRequest, pactus.TransactionOuterClass.CalculateFeeResponse> getCalculateFeeMethod;
    if ((getCalculateFeeMethod = TransactionGrpc.getCalculateFeeMethod) == null) {
      synchronized (TransactionGrpc.class) {
        if ((getCalculateFeeMethod = TransactionGrpc.getCalculateFeeMethod) == null) {
          TransactionGrpc.getCalculateFeeMethod = getCalculateFeeMethod =
              io.grpc.MethodDescriptor.<pactus.TransactionOuterClass.CalculateFeeRequest, pactus.TransactionOuterClass.CalculateFeeResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "CalculateFee"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.TransactionOuterClass.CalculateFeeRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.TransactionOuterClass.CalculateFeeResponse.getDefaultInstance()))
              .setSchemaDescriptor(new TransactionMethodDescriptorSupplier("CalculateFee"))
              .build();
        }
      }
    }
    return getCalculateFeeMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.TransactionOuterClass.BroadcastTransactionRequest,
      pactus.TransactionOuterClass.BroadcastTransactionResponse> getBroadcastTransactionMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "BroadcastTransaction",
      requestType = pactus.TransactionOuterClass.BroadcastTransactionRequest.class,
      responseType = pactus.TransactionOuterClass.BroadcastTransactionResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.TransactionOuterClass.BroadcastTransactionRequest,
      pactus.TransactionOuterClass.BroadcastTransactionResponse> getBroadcastTransactionMethod() {
    io.grpc.MethodDescriptor<pactus.TransactionOuterClass.BroadcastTransactionRequest, pactus.TransactionOuterClass.BroadcastTransactionResponse> getBroadcastTransactionMethod;
    if ((getBroadcastTransactionMethod = TransactionGrpc.getBroadcastTransactionMethod) == null) {
      synchronized (TransactionGrpc.class) {
        if ((getBroadcastTransactionMethod = TransactionGrpc.getBroadcastTransactionMethod) == null) {
          TransactionGrpc.getBroadcastTransactionMethod = getBroadcastTransactionMethod =
              io.grpc.MethodDescriptor.<pactus.TransactionOuterClass.BroadcastTransactionRequest, pactus.TransactionOuterClass.BroadcastTransactionResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "BroadcastTransaction"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.TransactionOuterClass.BroadcastTransactionRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.TransactionOuterClass.BroadcastTransactionResponse.getDefaultInstance()))
              .setSchemaDescriptor(new TransactionMethodDescriptorSupplier("BroadcastTransaction"))
              .build();
        }
      }
    }
    return getBroadcastTransactionMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.TransactionOuterClass.GetRawTransferTransactionRequest,
      pactus.TransactionOuterClass.GetRawTransactionResponse> getGetRawTransferTransactionMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetRawTransferTransaction",
      requestType = pactus.TransactionOuterClass.GetRawTransferTransactionRequest.class,
      responseType = pactus.TransactionOuterClass.GetRawTransactionResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.TransactionOuterClass.GetRawTransferTransactionRequest,
      pactus.TransactionOuterClass.GetRawTransactionResponse> getGetRawTransferTransactionMethod() {
    io.grpc.MethodDescriptor<pactus.TransactionOuterClass.GetRawTransferTransactionRequest, pactus.TransactionOuterClass.GetRawTransactionResponse> getGetRawTransferTransactionMethod;
    if ((getGetRawTransferTransactionMethod = TransactionGrpc.getGetRawTransferTransactionMethod) == null) {
      synchronized (TransactionGrpc.class) {
        if ((getGetRawTransferTransactionMethod = TransactionGrpc.getGetRawTransferTransactionMethod) == null) {
          TransactionGrpc.getGetRawTransferTransactionMethod = getGetRawTransferTransactionMethod =
              io.grpc.MethodDescriptor.<pactus.TransactionOuterClass.GetRawTransferTransactionRequest, pactus.TransactionOuterClass.GetRawTransactionResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetRawTransferTransaction"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.TransactionOuterClass.GetRawTransferTransactionRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.TransactionOuterClass.GetRawTransactionResponse.getDefaultInstance()))
              .setSchemaDescriptor(new TransactionMethodDescriptorSupplier("GetRawTransferTransaction"))
              .build();
        }
      }
    }
    return getGetRawTransferTransactionMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.TransactionOuterClass.GetRawBondTransactionRequest,
      pactus.TransactionOuterClass.GetRawTransactionResponse> getGetRawBondTransactionMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetRawBondTransaction",
      requestType = pactus.TransactionOuterClass.GetRawBondTransactionRequest.class,
      responseType = pactus.TransactionOuterClass.GetRawTransactionResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.TransactionOuterClass.GetRawBondTransactionRequest,
      pactus.TransactionOuterClass.GetRawTransactionResponse> getGetRawBondTransactionMethod() {
    io.grpc.MethodDescriptor<pactus.TransactionOuterClass.GetRawBondTransactionRequest, pactus.TransactionOuterClass.GetRawTransactionResponse> getGetRawBondTransactionMethod;
    if ((getGetRawBondTransactionMethod = TransactionGrpc.getGetRawBondTransactionMethod) == null) {
      synchronized (TransactionGrpc.class) {
        if ((getGetRawBondTransactionMethod = TransactionGrpc.getGetRawBondTransactionMethod) == null) {
          TransactionGrpc.getGetRawBondTransactionMethod = getGetRawBondTransactionMethod =
              io.grpc.MethodDescriptor.<pactus.TransactionOuterClass.GetRawBondTransactionRequest, pactus.TransactionOuterClass.GetRawTransactionResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetRawBondTransaction"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.TransactionOuterClass.GetRawBondTransactionRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.TransactionOuterClass.GetRawTransactionResponse.getDefaultInstance()))
              .setSchemaDescriptor(new TransactionMethodDescriptorSupplier("GetRawBondTransaction"))
              .build();
        }
      }
    }
    return getGetRawBondTransactionMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.TransactionOuterClass.GetRawUnbondTransactionRequest,
      pactus.TransactionOuterClass.GetRawTransactionResponse> getGetRawUnbondTransactionMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetRawUnbondTransaction",
      requestType = pactus.TransactionOuterClass.GetRawUnbondTransactionRequest.class,
      responseType = pactus.TransactionOuterClass.GetRawTransactionResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.TransactionOuterClass.GetRawUnbondTransactionRequest,
      pactus.TransactionOuterClass.GetRawTransactionResponse> getGetRawUnbondTransactionMethod() {
    io.grpc.MethodDescriptor<pactus.TransactionOuterClass.GetRawUnbondTransactionRequest, pactus.TransactionOuterClass.GetRawTransactionResponse> getGetRawUnbondTransactionMethod;
    if ((getGetRawUnbondTransactionMethod = TransactionGrpc.getGetRawUnbondTransactionMethod) == null) {
      synchronized (TransactionGrpc.class) {
        if ((getGetRawUnbondTransactionMethod = TransactionGrpc.getGetRawUnbondTransactionMethod) == null) {
          TransactionGrpc.getGetRawUnbondTransactionMethod = getGetRawUnbondTransactionMethod =
              io.grpc.MethodDescriptor.<pactus.TransactionOuterClass.GetRawUnbondTransactionRequest, pactus.TransactionOuterClass.GetRawTransactionResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetRawUnbondTransaction"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.TransactionOuterClass.GetRawUnbondTransactionRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.TransactionOuterClass.GetRawTransactionResponse.getDefaultInstance()))
              .setSchemaDescriptor(new TransactionMethodDescriptorSupplier("GetRawUnbondTransaction"))
              .build();
        }
      }
    }
    return getGetRawUnbondTransactionMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.TransactionOuterClass.GetRawWithdrawTransactionRequest,
      pactus.TransactionOuterClass.GetRawTransactionResponse> getGetRawWithdrawTransactionMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetRawWithdrawTransaction",
      requestType = pactus.TransactionOuterClass.GetRawWithdrawTransactionRequest.class,
      responseType = pactus.TransactionOuterClass.GetRawTransactionResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.TransactionOuterClass.GetRawWithdrawTransactionRequest,
      pactus.TransactionOuterClass.GetRawTransactionResponse> getGetRawWithdrawTransactionMethod() {
    io.grpc.MethodDescriptor<pactus.TransactionOuterClass.GetRawWithdrawTransactionRequest, pactus.TransactionOuterClass.GetRawTransactionResponse> getGetRawWithdrawTransactionMethod;
    if ((getGetRawWithdrawTransactionMethod = TransactionGrpc.getGetRawWithdrawTransactionMethod) == null) {
      synchronized (TransactionGrpc.class) {
        if ((getGetRawWithdrawTransactionMethod = TransactionGrpc.getGetRawWithdrawTransactionMethod) == null) {
          TransactionGrpc.getGetRawWithdrawTransactionMethod = getGetRawWithdrawTransactionMethod =
              io.grpc.MethodDescriptor.<pactus.TransactionOuterClass.GetRawWithdrawTransactionRequest, pactus.TransactionOuterClass.GetRawTransactionResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetRawWithdrawTransaction"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.TransactionOuterClass.GetRawWithdrawTransactionRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.TransactionOuterClass.GetRawTransactionResponse.getDefaultInstance()))
              .setSchemaDescriptor(new TransactionMethodDescriptorSupplier("GetRawWithdrawTransaction"))
              .build();
        }
      }
    }
    return getGetRawWithdrawTransactionMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.TransactionOuterClass.GetRawBatchTransferTransactionRequest,
      pactus.TransactionOuterClass.GetRawTransactionResponse> getGetRawBatchTransferTransactionMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetRawBatchTransferTransaction",
      requestType = pactus.TransactionOuterClass.GetRawBatchTransferTransactionRequest.class,
      responseType = pactus.TransactionOuterClass.GetRawTransactionResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.TransactionOuterClass.GetRawBatchTransferTransactionRequest,
      pactus.TransactionOuterClass.GetRawTransactionResponse> getGetRawBatchTransferTransactionMethod() {
    io.grpc.MethodDescriptor<pactus.TransactionOuterClass.GetRawBatchTransferTransactionRequest, pactus.TransactionOuterClass.GetRawTransactionResponse> getGetRawBatchTransferTransactionMethod;
    if ((getGetRawBatchTransferTransactionMethod = TransactionGrpc.getGetRawBatchTransferTransactionMethod) == null) {
      synchronized (TransactionGrpc.class) {
        if ((getGetRawBatchTransferTransactionMethod = TransactionGrpc.getGetRawBatchTransferTransactionMethod) == null) {
          TransactionGrpc.getGetRawBatchTransferTransactionMethod = getGetRawBatchTransferTransactionMethod =
              io.grpc.MethodDescriptor.<pactus.TransactionOuterClass.GetRawBatchTransferTransactionRequest, pactus.TransactionOuterClass.GetRawTransactionResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetRawBatchTransferTransaction"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.TransactionOuterClass.GetRawBatchTransferTransactionRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.TransactionOuterClass.GetRawTransactionResponse.getDefaultInstance()))
              .setSchemaDescriptor(new TransactionMethodDescriptorSupplier("GetRawBatchTransferTransaction"))
              .build();
        }
      }
    }
    return getGetRawBatchTransferTransactionMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.TransactionOuterClass.DecodeRawTransactionRequest,
      pactus.TransactionOuterClass.DecodeRawTransactionResponse> getDecodeRawTransactionMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "DecodeRawTransaction",
      requestType = pactus.TransactionOuterClass.DecodeRawTransactionRequest.class,
      responseType = pactus.TransactionOuterClass.DecodeRawTransactionResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.TransactionOuterClass.DecodeRawTransactionRequest,
      pactus.TransactionOuterClass.DecodeRawTransactionResponse> getDecodeRawTransactionMethod() {
    io.grpc.MethodDescriptor<pactus.TransactionOuterClass.DecodeRawTransactionRequest, pactus.TransactionOuterClass.DecodeRawTransactionResponse> getDecodeRawTransactionMethod;
    if ((getDecodeRawTransactionMethod = TransactionGrpc.getDecodeRawTransactionMethod) == null) {
      synchronized (TransactionGrpc.class) {
        if ((getDecodeRawTransactionMethod = TransactionGrpc.getDecodeRawTransactionMethod) == null) {
          TransactionGrpc.getDecodeRawTransactionMethod = getDecodeRawTransactionMethod =
              io.grpc.MethodDescriptor.<pactus.TransactionOuterClass.DecodeRawTransactionRequest, pactus.TransactionOuterClass.DecodeRawTransactionResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "DecodeRawTransaction"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.TransactionOuterClass.DecodeRawTransactionRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.TransactionOuterClass.DecodeRawTransactionResponse.getDefaultInstance()))
              .setSchemaDescriptor(new TransactionMethodDescriptorSupplier("DecodeRawTransaction"))
              .build();
        }
      }
    }
    return getDecodeRawTransactionMethod;
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
   * Creates a new blocking-style stub that supports all types of calls on the service
   */
  public static TransactionBlockingV2Stub newBlockingV2Stub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<TransactionBlockingV2Stub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<TransactionBlockingV2Stub>() {
        @java.lang.Override
        public TransactionBlockingV2Stub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new TransactionBlockingV2Stub(channel, callOptions);
        }
      };
    return TransactionBlockingV2Stub.newStub(factory, channel);
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
   * Transaction service defines various RPC methods for interacting with transactions.
   * </pre>
   */
  public interface AsyncService {

    /**
     * <pre>
     * GetTransaction retrieves transaction details based on the provided request parameters.
     * </pre>
     */
    default void getTransaction(pactus.TransactionOuterClass.GetTransactionRequest request,
        io.grpc.stub.StreamObserver<pactus.TransactionOuterClass.GetTransactionResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetTransactionMethod(), responseObserver);
    }

    /**
     * <pre>
     * CalculateFee calculates the transaction fee based on the specified amount and payload type.
     * </pre>
     */
    default void calculateFee(pactus.TransactionOuterClass.CalculateFeeRequest request,
        io.grpc.stub.StreamObserver<pactus.TransactionOuterClass.CalculateFeeResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getCalculateFeeMethod(), responseObserver);
    }

    /**
     * <pre>
     * BroadcastTransaction broadcasts a signed transaction to the network.
     * </pre>
     */
    default void broadcastTransaction(pactus.TransactionOuterClass.BroadcastTransactionRequest request,
        io.grpc.stub.StreamObserver<pactus.TransactionOuterClass.BroadcastTransactionResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getBroadcastTransactionMethod(), responseObserver);
    }

    /**
     * <pre>
     * GetRawTransferTransaction retrieves raw details of a transfer transaction.
     * </pre>
     */
    default void getRawTransferTransaction(pactus.TransactionOuterClass.GetRawTransferTransactionRequest request,
        io.grpc.stub.StreamObserver<pactus.TransactionOuterClass.GetRawTransactionResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetRawTransferTransactionMethod(), responseObserver);
    }

    /**
     * <pre>
     * GetRawBondTransaction retrieves raw details of a bond transaction.
     * </pre>
     */
    default void getRawBondTransaction(pactus.TransactionOuterClass.GetRawBondTransactionRequest request,
        io.grpc.stub.StreamObserver<pactus.TransactionOuterClass.GetRawTransactionResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetRawBondTransactionMethod(), responseObserver);
    }

    /**
     * <pre>
     * GetRawUnbondTransaction retrieves raw details of an unbond transaction.
     * </pre>
     */
    default void getRawUnbondTransaction(pactus.TransactionOuterClass.GetRawUnbondTransactionRequest request,
        io.grpc.stub.StreamObserver<pactus.TransactionOuterClass.GetRawTransactionResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetRawUnbondTransactionMethod(), responseObserver);
    }

    /**
     * <pre>
     * GetRawWithdrawTransaction retrieves raw details of a withdraw transaction.
     * </pre>
     */
    default void getRawWithdrawTransaction(pactus.TransactionOuterClass.GetRawWithdrawTransactionRequest request,
        io.grpc.stub.StreamObserver<pactus.TransactionOuterClass.GetRawTransactionResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetRawWithdrawTransactionMethod(), responseObserver);
    }

    /**
     * <pre>
     * GetRawBatchTransferTransaction retrieves raw details of batch transfer transaction.
     * </pre>
     */
    default void getRawBatchTransferTransaction(pactus.TransactionOuterClass.GetRawBatchTransferTransactionRequest request,
        io.grpc.stub.StreamObserver<pactus.TransactionOuterClass.GetRawTransactionResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetRawBatchTransferTransactionMethod(), responseObserver);
    }

    /**
     * <pre>
     * DecodeRawTransaction accepts raw transaction and returns decoded transaction.
     * </pre>
     */
    default void decodeRawTransaction(pactus.TransactionOuterClass.DecodeRawTransactionRequest request,
        io.grpc.stub.StreamObserver<pactus.TransactionOuterClass.DecodeRawTransactionResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getDecodeRawTransactionMethod(), responseObserver);
    }
  }

  /**
   * Base class for the server implementation of the service Transaction.
   * <pre>
   * Transaction service defines various RPC methods for interacting with transactions.
   * </pre>
   */
  public static abstract class TransactionImplBase
      implements io.grpc.BindableService, AsyncService {

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return TransactionGrpc.bindService(this);
    }
  }

  /**
   * A stub to allow clients to do asynchronous rpc calls to service Transaction.
   * <pre>
   * Transaction service defines various RPC methods for interacting with transactions.
   * </pre>
   */
  public static final class TransactionStub
      extends io.grpc.stub.AbstractAsyncStub<TransactionStub> {
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
     * GetTransaction retrieves transaction details based on the provided request parameters.
     * </pre>
     */
    public void getTransaction(pactus.TransactionOuterClass.GetTransactionRequest request,
        io.grpc.stub.StreamObserver<pactus.TransactionOuterClass.GetTransactionResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetTransactionMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * CalculateFee calculates the transaction fee based on the specified amount and payload type.
     * </pre>
     */
    public void calculateFee(pactus.TransactionOuterClass.CalculateFeeRequest request,
        io.grpc.stub.StreamObserver<pactus.TransactionOuterClass.CalculateFeeResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getCalculateFeeMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * BroadcastTransaction broadcasts a signed transaction to the network.
     * </pre>
     */
    public void broadcastTransaction(pactus.TransactionOuterClass.BroadcastTransactionRequest request,
        io.grpc.stub.StreamObserver<pactus.TransactionOuterClass.BroadcastTransactionResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getBroadcastTransactionMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * GetRawTransferTransaction retrieves raw details of a transfer transaction.
     * </pre>
     */
    public void getRawTransferTransaction(pactus.TransactionOuterClass.GetRawTransferTransactionRequest request,
        io.grpc.stub.StreamObserver<pactus.TransactionOuterClass.GetRawTransactionResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetRawTransferTransactionMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * GetRawBondTransaction retrieves raw details of a bond transaction.
     * </pre>
     */
    public void getRawBondTransaction(pactus.TransactionOuterClass.GetRawBondTransactionRequest request,
        io.grpc.stub.StreamObserver<pactus.TransactionOuterClass.GetRawTransactionResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetRawBondTransactionMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * GetRawUnbondTransaction retrieves raw details of an unbond transaction.
     * </pre>
     */
    public void getRawUnbondTransaction(pactus.TransactionOuterClass.GetRawUnbondTransactionRequest request,
        io.grpc.stub.StreamObserver<pactus.TransactionOuterClass.GetRawTransactionResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetRawUnbondTransactionMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * GetRawWithdrawTransaction retrieves raw details of a withdraw transaction.
     * </pre>
     */
    public void getRawWithdrawTransaction(pactus.TransactionOuterClass.GetRawWithdrawTransactionRequest request,
        io.grpc.stub.StreamObserver<pactus.TransactionOuterClass.GetRawTransactionResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetRawWithdrawTransactionMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * GetRawBatchTransferTransaction retrieves raw details of batch transfer transaction.
     * </pre>
     */
    public void getRawBatchTransferTransaction(pactus.TransactionOuterClass.GetRawBatchTransferTransactionRequest request,
        io.grpc.stub.StreamObserver<pactus.TransactionOuterClass.GetRawTransactionResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetRawBatchTransferTransactionMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * DecodeRawTransaction accepts raw transaction and returns decoded transaction.
     * </pre>
     */
    public void decodeRawTransaction(pactus.TransactionOuterClass.DecodeRawTransactionRequest request,
        io.grpc.stub.StreamObserver<pactus.TransactionOuterClass.DecodeRawTransactionResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getDecodeRawTransactionMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   * A stub to allow clients to do synchronous rpc calls to service Transaction.
   * <pre>
   * Transaction service defines various RPC methods for interacting with transactions.
   * </pre>
   */
  public static final class TransactionBlockingV2Stub
      extends io.grpc.stub.AbstractBlockingStub<TransactionBlockingV2Stub> {
    private TransactionBlockingV2Stub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected TransactionBlockingV2Stub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new TransactionBlockingV2Stub(channel, callOptions);
    }

    /**
     * <pre>
     * GetTransaction retrieves transaction details based on the provided request parameters.
     * </pre>
     */
    public pactus.TransactionOuterClass.GetTransactionResponse getTransaction(pactus.TransactionOuterClass.GetTransactionRequest request) throws io.grpc.StatusException {
      return io.grpc.stub.ClientCalls.blockingV2UnaryCall(
          getChannel(), getGetTransactionMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * CalculateFee calculates the transaction fee based on the specified amount and payload type.
     * </pre>
     */
    public pactus.TransactionOuterClass.CalculateFeeResponse calculateFee(pactus.TransactionOuterClass.CalculateFeeRequest request) throws io.grpc.StatusException {
      return io.grpc.stub.ClientCalls.blockingV2UnaryCall(
          getChannel(), getCalculateFeeMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * BroadcastTransaction broadcasts a signed transaction to the network.
     * </pre>
     */
    public pactus.TransactionOuterClass.BroadcastTransactionResponse broadcastTransaction(pactus.TransactionOuterClass.BroadcastTransactionRequest request) throws io.grpc.StatusException {
      return io.grpc.stub.ClientCalls.blockingV2UnaryCall(
          getChannel(), getBroadcastTransactionMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetRawTransferTransaction retrieves raw details of a transfer transaction.
     * </pre>
     */
    public pactus.TransactionOuterClass.GetRawTransactionResponse getRawTransferTransaction(pactus.TransactionOuterClass.GetRawTransferTransactionRequest request) throws io.grpc.StatusException {
      return io.grpc.stub.ClientCalls.blockingV2UnaryCall(
          getChannel(), getGetRawTransferTransactionMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetRawBondTransaction retrieves raw details of a bond transaction.
     * </pre>
     */
    public pactus.TransactionOuterClass.GetRawTransactionResponse getRawBondTransaction(pactus.TransactionOuterClass.GetRawBondTransactionRequest request) throws io.grpc.StatusException {
      return io.grpc.stub.ClientCalls.blockingV2UnaryCall(
          getChannel(), getGetRawBondTransactionMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetRawUnbondTransaction retrieves raw details of an unbond transaction.
     * </pre>
     */
    public pactus.TransactionOuterClass.GetRawTransactionResponse getRawUnbondTransaction(pactus.TransactionOuterClass.GetRawUnbondTransactionRequest request) throws io.grpc.StatusException {
      return io.grpc.stub.ClientCalls.blockingV2UnaryCall(
          getChannel(), getGetRawUnbondTransactionMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetRawWithdrawTransaction retrieves raw details of a withdraw transaction.
     * </pre>
     */
    public pactus.TransactionOuterClass.GetRawTransactionResponse getRawWithdrawTransaction(pactus.TransactionOuterClass.GetRawWithdrawTransactionRequest request) throws io.grpc.StatusException {
      return io.grpc.stub.ClientCalls.blockingV2UnaryCall(
          getChannel(), getGetRawWithdrawTransactionMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetRawBatchTransferTransaction retrieves raw details of batch transfer transaction.
     * </pre>
     */
    public pactus.TransactionOuterClass.GetRawTransactionResponse getRawBatchTransferTransaction(pactus.TransactionOuterClass.GetRawBatchTransferTransactionRequest request) throws io.grpc.StatusException {
      return io.grpc.stub.ClientCalls.blockingV2UnaryCall(
          getChannel(), getGetRawBatchTransferTransactionMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * DecodeRawTransaction accepts raw transaction and returns decoded transaction.
     * </pre>
     */
    public pactus.TransactionOuterClass.DecodeRawTransactionResponse decodeRawTransaction(pactus.TransactionOuterClass.DecodeRawTransactionRequest request) throws io.grpc.StatusException {
      return io.grpc.stub.ClientCalls.blockingV2UnaryCall(
          getChannel(), getDecodeRawTransactionMethod(), getCallOptions(), request);
    }
  }

  /**
   * A stub to allow clients to do limited synchronous rpc calls to service Transaction.
   * <pre>
   * Transaction service defines various RPC methods for interacting with transactions.
   * </pre>
   */
  public static final class TransactionBlockingStub
      extends io.grpc.stub.AbstractBlockingStub<TransactionBlockingStub> {
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
     * GetTransaction retrieves transaction details based on the provided request parameters.
     * </pre>
     */
    public pactus.TransactionOuterClass.GetTransactionResponse getTransaction(pactus.TransactionOuterClass.GetTransactionRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetTransactionMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * CalculateFee calculates the transaction fee based on the specified amount and payload type.
     * </pre>
     */
    public pactus.TransactionOuterClass.CalculateFeeResponse calculateFee(pactus.TransactionOuterClass.CalculateFeeRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getCalculateFeeMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * BroadcastTransaction broadcasts a signed transaction to the network.
     * </pre>
     */
    public pactus.TransactionOuterClass.BroadcastTransactionResponse broadcastTransaction(pactus.TransactionOuterClass.BroadcastTransactionRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getBroadcastTransactionMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetRawTransferTransaction retrieves raw details of a transfer transaction.
     * </pre>
     */
    public pactus.TransactionOuterClass.GetRawTransactionResponse getRawTransferTransaction(pactus.TransactionOuterClass.GetRawTransferTransactionRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetRawTransferTransactionMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetRawBondTransaction retrieves raw details of a bond transaction.
     * </pre>
     */
    public pactus.TransactionOuterClass.GetRawTransactionResponse getRawBondTransaction(pactus.TransactionOuterClass.GetRawBondTransactionRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetRawBondTransactionMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetRawUnbondTransaction retrieves raw details of an unbond transaction.
     * </pre>
     */
    public pactus.TransactionOuterClass.GetRawTransactionResponse getRawUnbondTransaction(pactus.TransactionOuterClass.GetRawUnbondTransactionRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetRawUnbondTransactionMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetRawWithdrawTransaction retrieves raw details of a withdraw transaction.
     * </pre>
     */
    public pactus.TransactionOuterClass.GetRawTransactionResponse getRawWithdrawTransaction(pactus.TransactionOuterClass.GetRawWithdrawTransactionRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetRawWithdrawTransactionMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetRawBatchTransferTransaction retrieves raw details of batch transfer transaction.
     * </pre>
     */
    public pactus.TransactionOuterClass.GetRawTransactionResponse getRawBatchTransferTransaction(pactus.TransactionOuterClass.GetRawBatchTransferTransactionRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetRawBatchTransferTransactionMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * DecodeRawTransaction accepts raw transaction and returns decoded transaction.
     * </pre>
     */
    public pactus.TransactionOuterClass.DecodeRawTransactionResponse decodeRawTransaction(pactus.TransactionOuterClass.DecodeRawTransactionRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getDecodeRawTransactionMethod(), getCallOptions(), request);
    }
  }

  /**
   * A stub to allow clients to do ListenableFuture-style rpc calls to service Transaction.
   * <pre>
   * Transaction service defines various RPC methods for interacting with transactions.
   * </pre>
   */
  public static final class TransactionFutureStub
      extends io.grpc.stub.AbstractFutureStub<TransactionFutureStub> {
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
     * GetTransaction retrieves transaction details based on the provided request parameters.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.TransactionOuterClass.GetTransactionResponse> getTransaction(
        pactus.TransactionOuterClass.GetTransactionRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetTransactionMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * CalculateFee calculates the transaction fee based on the specified amount and payload type.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.TransactionOuterClass.CalculateFeeResponse> calculateFee(
        pactus.TransactionOuterClass.CalculateFeeRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getCalculateFeeMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * BroadcastTransaction broadcasts a signed transaction to the network.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.TransactionOuterClass.BroadcastTransactionResponse> broadcastTransaction(
        pactus.TransactionOuterClass.BroadcastTransactionRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getBroadcastTransactionMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * GetRawTransferTransaction retrieves raw details of a transfer transaction.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.TransactionOuterClass.GetRawTransactionResponse> getRawTransferTransaction(
        pactus.TransactionOuterClass.GetRawTransferTransactionRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetRawTransferTransactionMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * GetRawBondTransaction retrieves raw details of a bond transaction.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.TransactionOuterClass.GetRawTransactionResponse> getRawBondTransaction(
        pactus.TransactionOuterClass.GetRawBondTransactionRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetRawBondTransactionMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * GetRawUnbondTransaction retrieves raw details of an unbond transaction.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.TransactionOuterClass.GetRawTransactionResponse> getRawUnbondTransaction(
        pactus.TransactionOuterClass.GetRawUnbondTransactionRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetRawUnbondTransactionMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * GetRawWithdrawTransaction retrieves raw details of a withdraw transaction.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.TransactionOuterClass.GetRawTransactionResponse> getRawWithdrawTransaction(
        pactus.TransactionOuterClass.GetRawWithdrawTransactionRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetRawWithdrawTransactionMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * GetRawBatchTransferTransaction retrieves raw details of batch transfer transaction.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.TransactionOuterClass.GetRawTransactionResponse> getRawBatchTransferTransaction(
        pactus.TransactionOuterClass.GetRawBatchTransferTransactionRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetRawBatchTransferTransactionMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * DecodeRawTransaction accepts raw transaction and returns decoded transaction.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.TransactionOuterClass.DecodeRawTransactionResponse> decodeRawTransaction(
        pactus.TransactionOuterClass.DecodeRawTransactionRequest request) {
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
  private static final int METHODID_GET_RAW_BATCH_TRANSFER_TRANSACTION = 7;
  private static final int METHODID_DECODE_RAW_TRANSACTION = 8;

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
          serviceImpl.getTransaction((pactus.TransactionOuterClass.GetTransactionRequest) request,
              (io.grpc.stub.StreamObserver<pactus.TransactionOuterClass.GetTransactionResponse>) responseObserver);
          break;
        case METHODID_CALCULATE_FEE:
          serviceImpl.calculateFee((pactus.TransactionOuterClass.CalculateFeeRequest) request,
              (io.grpc.stub.StreamObserver<pactus.TransactionOuterClass.CalculateFeeResponse>) responseObserver);
          break;
        case METHODID_BROADCAST_TRANSACTION:
          serviceImpl.broadcastTransaction((pactus.TransactionOuterClass.BroadcastTransactionRequest) request,
              (io.grpc.stub.StreamObserver<pactus.TransactionOuterClass.BroadcastTransactionResponse>) responseObserver);
          break;
        case METHODID_GET_RAW_TRANSFER_TRANSACTION:
          serviceImpl.getRawTransferTransaction((pactus.TransactionOuterClass.GetRawTransferTransactionRequest) request,
              (io.grpc.stub.StreamObserver<pactus.TransactionOuterClass.GetRawTransactionResponse>) responseObserver);
          break;
        case METHODID_GET_RAW_BOND_TRANSACTION:
          serviceImpl.getRawBondTransaction((pactus.TransactionOuterClass.GetRawBondTransactionRequest) request,
              (io.grpc.stub.StreamObserver<pactus.TransactionOuterClass.GetRawTransactionResponse>) responseObserver);
          break;
        case METHODID_GET_RAW_UNBOND_TRANSACTION:
          serviceImpl.getRawUnbondTransaction((pactus.TransactionOuterClass.GetRawUnbondTransactionRequest) request,
              (io.grpc.stub.StreamObserver<pactus.TransactionOuterClass.GetRawTransactionResponse>) responseObserver);
          break;
        case METHODID_GET_RAW_WITHDRAW_TRANSACTION:
          serviceImpl.getRawWithdrawTransaction((pactus.TransactionOuterClass.GetRawWithdrawTransactionRequest) request,
              (io.grpc.stub.StreamObserver<pactus.TransactionOuterClass.GetRawTransactionResponse>) responseObserver);
          break;
        case METHODID_GET_RAW_BATCH_TRANSFER_TRANSACTION:
          serviceImpl.getRawBatchTransferTransaction((pactus.TransactionOuterClass.GetRawBatchTransferTransactionRequest) request,
              (io.grpc.stub.StreamObserver<pactus.TransactionOuterClass.GetRawTransactionResponse>) responseObserver);
          break;
        case METHODID_DECODE_RAW_TRANSACTION:
          serviceImpl.decodeRawTransaction((pactus.TransactionOuterClass.DecodeRawTransactionRequest) request,
              (io.grpc.stub.StreamObserver<pactus.TransactionOuterClass.DecodeRawTransactionResponse>) responseObserver);
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
              pactus.TransactionOuterClass.GetTransactionRequest,
              pactus.TransactionOuterClass.GetTransactionResponse>(
                service, METHODID_GET_TRANSACTION)))
        .addMethod(
          getCalculateFeeMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.TransactionOuterClass.CalculateFeeRequest,
              pactus.TransactionOuterClass.CalculateFeeResponse>(
                service, METHODID_CALCULATE_FEE)))
        .addMethod(
          getBroadcastTransactionMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.TransactionOuterClass.BroadcastTransactionRequest,
              pactus.TransactionOuterClass.BroadcastTransactionResponse>(
                service, METHODID_BROADCAST_TRANSACTION)))
        .addMethod(
          getGetRawTransferTransactionMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.TransactionOuterClass.GetRawTransferTransactionRequest,
              pactus.TransactionOuterClass.GetRawTransactionResponse>(
                service, METHODID_GET_RAW_TRANSFER_TRANSACTION)))
        .addMethod(
          getGetRawBondTransactionMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.TransactionOuterClass.GetRawBondTransactionRequest,
              pactus.TransactionOuterClass.GetRawTransactionResponse>(
                service, METHODID_GET_RAW_BOND_TRANSACTION)))
        .addMethod(
          getGetRawUnbondTransactionMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.TransactionOuterClass.GetRawUnbondTransactionRequest,
              pactus.TransactionOuterClass.GetRawTransactionResponse>(
                service, METHODID_GET_RAW_UNBOND_TRANSACTION)))
        .addMethod(
          getGetRawWithdrawTransactionMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.TransactionOuterClass.GetRawWithdrawTransactionRequest,
              pactus.TransactionOuterClass.GetRawTransactionResponse>(
                service, METHODID_GET_RAW_WITHDRAW_TRANSACTION)))
        .addMethod(
          getGetRawBatchTransferTransactionMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.TransactionOuterClass.GetRawBatchTransferTransactionRequest,
              pactus.TransactionOuterClass.GetRawTransactionResponse>(
                service, METHODID_GET_RAW_BATCH_TRANSFER_TRANSACTION)))
        .addMethod(
          getDecodeRawTransactionMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.TransactionOuterClass.DecodeRawTransactionRequest,
              pactus.TransactionOuterClass.DecodeRawTransactionResponse>(
                service, METHODID_DECODE_RAW_TRANSACTION)))
        .build();
  }

  private static abstract class TransactionBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    TransactionBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return pactus.TransactionOuterClass.getDescriptor();
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
    private final java.lang.String methodName;

    TransactionMethodDescriptorSupplier(java.lang.String methodName) {
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
              .addMethod(getGetRawUnbondTransactionMethod())
              .addMethod(getGetRawWithdrawTransactionMethod())
              .addMethod(getGetRawBatchTransferTransactionMethod())
              .addMethod(getDecodeRawTransactionMethod())
              .build();
        }
      }
    }
    return result;
  }
}
