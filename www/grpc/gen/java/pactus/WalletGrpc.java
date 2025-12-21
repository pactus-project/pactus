package pactus;

import static io.grpc.MethodDescriptor.generateFullMethodName;

/**
 * <pre>
 * Wallet service provides RPC methods for wallet management operations.
 * </pre>
 */
@io.grpc.stub.annotations.GrpcGenerated
public final class WalletGrpc {

  private WalletGrpc() {}

  public static final java.lang.String SERVICE_NAME = "pactus.Wallet";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<pactus.WalletOuterClass.CreateWalletRequest,
      pactus.WalletOuterClass.CreateWalletResponse> getCreateWalletMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "CreateWallet",
      requestType = pactus.WalletOuterClass.CreateWalletRequest.class,
      responseType = pactus.WalletOuterClass.CreateWalletResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.WalletOuterClass.CreateWalletRequest,
      pactus.WalletOuterClass.CreateWalletResponse> getCreateWalletMethod() {
    io.grpc.MethodDescriptor<pactus.WalletOuterClass.CreateWalletRequest, pactus.WalletOuterClass.CreateWalletResponse> getCreateWalletMethod;
    if ((getCreateWalletMethod = WalletGrpc.getCreateWalletMethod) == null) {
      synchronized (WalletGrpc.class) {
        if ((getCreateWalletMethod = WalletGrpc.getCreateWalletMethod) == null) {
          WalletGrpc.getCreateWalletMethod = getCreateWalletMethod =
              io.grpc.MethodDescriptor.<pactus.WalletOuterClass.CreateWalletRequest, pactus.WalletOuterClass.CreateWalletResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "CreateWallet"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.WalletOuterClass.CreateWalletRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.WalletOuterClass.CreateWalletResponse.getDefaultInstance()))
              .setSchemaDescriptor(new WalletMethodDescriptorSupplier("CreateWallet"))
              .build();
        }
      }
    }
    return getCreateWalletMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.WalletOuterClass.RestoreWalletRequest,
      pactus.WalletOuterClass.RestoreWalletResponse> getRestoreWalletMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "RestoreWallet",
      requestType = pactus.WalletOuterClass.RestoreWalletRequest.class,
      responseType = pactus.WalletOuterClass.RestoreWalletResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.WalletOuterClass.RestoreWalletRequest,
      pactus.WalletOuterClass.RestoreWalletResponse> getRestoreWalletMethod() {
    io.grpc.MethodDescriptor<pactus.WalletOuterClass.RestoreWalletRequest, pactus.WalletOuterClass.RestoreWalletResponse> getRestoreWalletMethod;
    if ((getRestoreWalletMethod = WalletGrpc.getRestoreWalletMethod) == null) {
      synchronized (WalletGrpc.class) {
        if ((getRestoreWalletMethod = WalletGrpc.getRestoreWalletMethod) == null) {
          WalletGrpc.getRestoreWalletMethod = getRestoreWalletMethod =
              io.grpc.MethodDescriptor.<pactus.WalletOuterClass.RestoreWalletRequest, pactus.WalletOuterClass.RestoreWalletResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "RestoreWallet"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.WalletOuterClass.RestoreWalletRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.WalletOuterClass.RestoreWalletResponse.getDefaultInstance()))
              .setSchemaDescriptor(new WalletMethodDescriptorSupplier("RestoreWallet"))
              .build();
        }
      }
    }
    return getRestoreWalletMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.WalletOuterClass.LoadWalletRequest,
      pactus.WalletOuterClass.LoadWalletResponse> getLoadWalletMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "LoadWallet",
      requestType = pactus.WalletOuterClass.LoadWalletRequest.class,
      responseType = pactus.WalletOuterClass.LoadWalletResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.WalletOuterClass.LoadWalletRequest,
      pactus.WalletOuterClass.LoadWalletResponse> getLoadWalletMethod() {
    io.grpc.MethodDescriptor<pactus.WalletOuterClass.LoadWalletRequest, pactus.WalletOuterClass.LoadWalletResponse> getLoadWalletMethod;
    if ((getLoadWalletMethod = WalletGrpc.getLoadWalletMethod) == null) {
      synchronized (WalletGrpc.class) {
        if ((getLoadWalletMethod = WalletGrpc.getLoadWalletMethod) == null) {
          WalletGrpc.getLoadWalletMethod = getLoadWalletMethod =
              io.grpc.MethodDescriptor.<pactus.WalletOuterClass.LoadWalletRequest, pactus.WalletOuterClass.LoadWalletResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "LoadWallet"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.WalletOuterClass.LoadWalletRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.WalletOuterClass.LoadWalletResponse.getDefaultInstance()))
              .setSchemaDescriptor(new WalletMethodDescriptorSupplier("LoadWallet"))
              .build();
        }
      }
    }
    return getLoadWalletMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.WalletOuterClass.UnloadWalletRequest,
      pactus.WalletOuterClass.UnloadWalletResponse> getUnloadWalletMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "UnloadWallet",
      requestType = pactus.WalletOuterClass.UnloadWalletRequest.class,
      responseType = pactus.WalletOuterClass.UnloadWalletResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.WalletOuterClass.UnloadWalletRequest,
      pactus.WalletOuterClass.UnloadWalletResponse> getUnloadWalletMethod() {
    io.grpc.MethodDescriptor<pactus.WalletOuterClass.UnloadWalletRequest, pactus.WalletOuterClass.UnloadWalletResponse> getUnloadWalletMethod;
    if ((getUnloadWalletMethod = WalletGrpc.getUnloadWalletMethod) == null) {
      synchronized (WalletGrpc.class) {
        if ((getUnloadWalletMethod = WalletGrpc.getUnloadWalletMethod) == null) {
          WalletGrpc.getUnloadWalletMethod = getUnloadWalletMethod =
              io.grpc.MethodDescriptor.<pactus.WalletOuterClass.UnloadWalletRequest, pactus.WalletOuterClass.UnloadWalletResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "UnloadWallet"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.WalletOuterClass.UnloadWalletRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.WalletOuterClass.UnloadWalletResponse.getDefaultInstance()))
              .setSchemaDescriptor(new WalletMethodDescriptorSupplier("UnloadWallet"))
              .build();
        }
      }
    }
    return getUnloadWalletMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.WalletOuterClass.GetTotalBalanceRequest,
      pactus.WalletOuterClass.GetTotalBalanceResponse> getGetTotalBalanceMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetTotalBalance",
      requestType = pactus.WalletOuterClass.GetTotalBalanceRequest.class,
      responseType = pactus.WalletOuterClass.GetTotalBalanceResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.WalletOuterClass.GetTotalBalanceRequest,
      pactus.WalletOuterClass.GetTotalBalanceResponse> getGetTotalBalanceMethod() {
    io.grpc.MethodDescriptor<pactus.WalletOuterClass.GetTotalBalanceRequest, pactus.WalletOuterClass.GetTotalBalanceResponse> getGetTotalBalanceMethod;
    if ((getGetTotalBalanceMethod = WalletGrpc.getGetTotalBalanceMethod) == null) {
      synchronized (WalletGrpc.class) {
        if ((getGetTotalBalanceMethod = WalletGrpc.getGetTotalBalanceMethod) == null) {
          WalletGrpc.getGetTotalBalanceMethod = getGetTotalBalanceMethod =
              io.grpc.MethodDescriptor.<pactus.WalletOuterClass.GetTotalBalanceRequest, pactus.WalletOuterClass.GetTotalBalanceResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetTotalBalance"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.WalletOuterClass.GetTotalBalanceRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.WalletOuterClass.GetTotalBalanceResponse.getDefaultInstance()))
              .setSchemaDescriptor(new WalletMethodDescriptorSupplier("GetTotalBalance"))
              .build();
        }
      }
    }
    return getGetTotalBalanceMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.WalletOuterClass.SignRawTransactionRequest,
      pactus.WalletOuterClass.SignRawTransactionResponse> getSignRawTransactionMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "SignRawTransaction",
      requestType = pactus.WalletOuterClass.SignRawTransactionRequest.class,
      responseType = pactus.WalletOuterClass.SignRawTransactionResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.WalletOuterClass.SignRawTransactionRequest,
      pactus.WalletOuterClass.SignRawTransactionResponse> getSignRawTransactionMethod() {
    io.grpc.MethodDescriptor<pactus.WalletOuterClass.SignRawTransactionRequest, pactus.WalletOuterClass.SignRawTransactionResponse> getSignRawTransactionMethod;
    if ((getSignRawTransactionMethod = WalletGrpc.getSignRawTransactionMethod) == null) {
      synchronized (WalletGrpc.class) {
        if ((getSignRawTransactionMethod = WalletGrpc.getSignRawTransactionMethod) == null) {
          WalletGrpc.getSignRawTransactionMethod = getSignRawTransactionMethod =
              io.grpc.MethodDescriptor.<pactus.WalletOuterClass.SignRawTransactionRequest, pactus.WalletOuterClass.SignRawTransactionResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "SignRawTransaction"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.WalletOuterClass.SignRawTransactionRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.WalletOuterClass.SignRawTransactionResponse.getDefaultInstance()))
              .setSchemaDescriptor(new WalletMethodDescriptorSupplier("SignRawTransaction"))
              .build();
        }
      }
    }
    return getSignRawTransactionMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.WalletOuterClass.GetValidatorAddressRequest,
      pactus.WalletOuterClass.GetValidatorAddressResponse> getGetValidatorAddressMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetValidatorAddress",
      requestType = pactus.WalletOuterClass.GetValidatorAddressRequest.class,
      responseType = pactus.WalletOuterClass.GetValidatorAddressResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.WalletOuterClass.GetValidatorAddressRequest,
      pactus.WalletOuterClass.GetValidatorAddressResponse> getGetValidatorAddressMethod() {
    io.grpc.MethodDescriptor<pactus.WalletOuterClass.GetValidatorAddressRequest, pactus.WalletOuterClass.GetValidatorAddressResponse> getGetValidatorAddressMethod;
    if ((getGetValidatorAddressMethod = WalletGrpc.getGetValidatorAddressMethod) == null) {
      synchronized (WalletGrpc.class) {
        if ((getGetValidatorAddressMethod = WalletGrpc.getGetValidatorAddressMethod) == null) {
          WalletGrpc.getGetValidatorAddressMethod = getGetValidatorAddressMethod =
              io.grpc.MethodDescriptor.<pactus.WalletOuterClass.GetValidatorAddressRequest, pactus.WalletOuterClass.GetValidatorAddressResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetValidatorAddress"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.WalletOuterClass.GetValidatorAddressRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.WalletOuterClass.GetValidatorAddressResponse.getDefaultInstance()))
              .setSchemaDescriptor(new WalletMethodDescriptorSupplier("GetValidatorAddress"))
              .build();
        }
      }
    }
    return getGetValidatorAddressMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.WalletOuterClass.GetNewAddressRequest,
      pactus.WalletOuterClass.GetNewAddressResponse> getGetNewAddressMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetNewAddress",
      requestType = pactus.WalletOuterClass.GetNewAddressRequest.class,
      responseType = pactus.WalletOuterClass.GetNewAddressResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.WalletOuterClass.GetNewAddressRequest,
      pactus.WalletOuterClass.GetNewAddressResponse> getGetNewAddressMethod() {
    io.grpc.MethodDescriptor<pactus.WalletOuterClass.GetNewAddressRequest, pactus.WalletOuterClass.GetNewAddressResponse> getGetNewAddressMethod;
    if ((getGetNewAddressMethod = WalletGrpc.getGetNewAddressMethod) == null) {
      synchronized (WalletGrpc.class) {
        if ((getGetNewAddressMethod = WalletGrpc.getGetNewAddressMethod) == null) {
          WalletGrpc.getGetNewAddressMethod = getGetNewAddressMethod =
              io.grpc.MethodDescriptor.<pactus.WalletOuterClass.GetNewAddressRequest, pactus.WalletOuterClass.GetNewAddressResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetNewAddress"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.WalletOuterClass.GetNewAddressRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.WalletOuterClass.GetNewAddressResponse.getDefaultInstance()))
              .setSchemaDescriptor(new WalletMethodDescriptorSupplier("GetNewAddress"))
              .build();
        }
      }
    }
    return getGetNewAddressMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.WalletOuterClass.GetAddressHistoryRequest,
      pactus.WalletOuterClass.GetAddressHistoryResponse> getGetAddressHistoryMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetAddressHistory",
      requestType = pactus.WalletOuterClass.GetAddressHistoryRequest.class,
      responseType = pactus.WalletOuterClass.GetAddressHistoryResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.WalletOuterClass.GetAddressHistoryRequest,
      pactus.WalletOuterClass.GetAddressHistoryResponse> getGetAddressHistoryMethod() {
    io.grpc.MethodDescriptor<pactus.WalletOuterClass.GetAddressHistoryRequest, pactus.WalletOuterClass.GetAddressHistoryResponse> getGetAddressHistoryMethod;
    if ((getGetAddressHistoryMethod = WalletGrpc.getGetAddressHistoryMethod) == null) {
      synchronized (WalletGrpc.class) {
        if ((getGetAddressHistoryMethod = WalletGrpc.getGetAddressHistoryMethod) == null) {
          WalletGrpc.getGetAddressHistoryMethod = getGetAddressHistoryMethod =
              io.grpc.MethodDescriptor.<pactus.WalletOuterClass.GetAddressHistoryRequest, pactus.WalletOuterClass.GetAddressHistoryResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetAddressHistory"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.WalletOuterClass.GetAddressHistoryRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.WalletOuterClass.GetAddressHistoryResponse.getDefaultInstance()))
              .setSchemaDescriptor(new WalletMethodDescriptorSupplier("GetAddressHistory"))
              .build();
        }
      }
    }
    return getGetAddressHistoryMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.WalletOuterClass.SignMessageRequest,
      pactus.WalletOuterClass.SignMessageResponse> getSignMessageMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "SignMessage",
      requestType = pactus.WalletOuterClass.SignMessageRequest.class,
      responseType = pactus.WalletOuterClass.SignMessageResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.WalletOuterClass.SignMessageRequest,
      pactus.WalletOuterClass.SignMessageResponse> getSignMessageMethod() {
    io.grpc.MethodDescriptor<pactus.WalletOuterClass.SignMessageRequest, pactus.WalletOuterClass.SignMessageResponse> getSignMessageMethod;
    if ((getSignMessageMethod = WalletGrpc.getSignMessageMethod) == null) {
      synchronized (WalletGrpc.class) {
        if ((getSignMessageMethod = WalletGrpc.getSignMessageMethod) == null) {
          WalletGrpc.getSignMessageMethod = getSignMessageMethod =
              io.grpc.MethodDescriptor.<pactus.WalletOuterClass.SignMessageRequest, pactus.WalletOuterClass.SignMessageResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "SignMessage"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.WalletOuterClass.SignMessageRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.WalletOuterClass.SignMessageResponse.getDefaultInstance()))
              .setSchemaDescriptor(new WalletMethodDescriptorSupplier("SignMessage"))
              .build();
        }
      }
    }
    return getSignMessageMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.WalletOuterClass.GetTotalStakeRequest,
      pactus.WalletOuterClass.GetTotalStakeResponse> getGetTotalStakeMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetTotalStake",
      requestType = pactus.WalletOuterClass.GetTotalStakeRequest.class,
      responseType = pactus.WalletOuterClass.GetTotalStakeResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.WalletOuterClass.GetTotalStakeRequest,
      pactus.WalletOuterClass.GetTotalStakeResponse> getGetTotalStakeMethod() {
    io.grpc.MethodDescriptor<pactus.WalletOuterClass.GetTotalStakeRequest, pactus.WalletOuterClass.GetTotalStakeResponse> getGetTotalStakeMethod;
    if ((getGetTotalStakeMethod = WalletGrpc.getGetTotalStakeMethod) == null) {
      synchronized (WalletGrpc.class) {
        if ((getGetTotalStakeMethod = WalletGrpc.getGetTotalStakeMethod) == null) {
          WalletGrpc.getGetTotalStakeMethod = getGetTotalStakeMethod =
              io.grpc.MethodDescriptor.<pactus.WalletOuterClass.GetTotalStakeRequest, pactus.WalletOuterClass.GetTotalStakeResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetTotalStake"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.WalletOuterClass.GetTotalStakeRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.WalletOuterClass.GetTotalStakeResponse.getDefaultInstance()))
              .setSchemaDescriptor(new WalletMethodDescriptorSupplier("GetTotalStake"))
              .build();
        }
      }
    }
    return getGetTotalStakeMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.WalletOuterClass.GetAddressInfoRequest,
      pactus.WalletOuterClass.GetAddressInfoResponse> getGetAddressInfoMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetAddressInfo",
      requestType = pactus.WalletOuterClass.GetAddressInfoRequest.class,
      responseType = pactus.WalletOuterClass.GetAddressInfoResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.WalletOuterClass.GetAddressInfoRequest,
      pactus.WalletOuterClass.GetAddressInfoResponse> getGetAddressInfoMethod() {
    io.grpc.MethodDescriptor<pactus.WalletOuterClass.GetAddressInfoRequest, pactus.WalletOuterClass.GetAddressInfoResponse> getGetAddressInfoMethod;
    if ((getGetAddressInfoMethod = WalletGrpc.getGetAddressInfoMethod) == null) {
      synchronized (WalletGrpc.class) {
        if ((getGetAddressInfoMethod = WalletGrpc.getGetAddressInfoMethod) == null) {
          WalletGrpc.getGetAddressInfoMethod = getGetAddressInfoMethod =
              io.grpc.MethodDescriptor.<pactus.WalletOuterClass.GetAddressInfoRequest, pactus.WalletOuterClass.GetAddressInfoResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetAddressInfo"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.WalletOuterClass.GetAddressInfoRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.WalletOuterClass.GetAddressInfoResponse.getDefaultInstance()))
              .setSchemaDescriptor(new WalletMethodDescriptorSupplier("GetAddressInfo"))
              .build();
        }
      }
    }
    return getGetAddressInfoMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.WalletOuterClass.SetAddressLabelRequest,
      pactus.WalletOuterClass.SetAddressLabelResponse> getSetAddressLabelMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "SetAddressLabel",
      requestType = pactus.WalletOuterClass.SetAddressLabelRequest.class,
      responseType = pactus.WalletOuterClass.SetAddressLabelResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.WalletOuterClass.SetAddressLabelRequest,
      pactus.WalletOuterClass.SetAddressLabelResponse> getSetAddressLabelMethod() {
    io.grpc.MethodDescriptor<pactus.WalletOuterClass.SetAddressLabelRequest, pactus.WalletOuterClass.SetAddressLabelResponse> getSetAddressLabelMethod;
    if ((getSetAddressLabelMethod = WalletGrpc.getSetAddressLabelMethod) == null) {
      synchronized (WalletGrpc.class) {
        if ((getSetAddressLabelMethod = WalletGrpc.getSetAddressLabelMethod) == null) {
          WalletGrpc.getSetAddressLabelMethod = getSetAddressLabelMethod =
              io.grpc.MethodDescriptor.<pactus.WalletOuterClass.SetAddressLabelRequest, pactus.WalletOuterClass.SetAddressLabelResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "SetAddressLabel"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.WalletOuterClass.SetAddressLabelRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.WalletOuterClass.SetAddressLabelResponse.getDefaultInstance()))
              .setSchemaDescriptor(new WalletMethodDescriptorSupplier("SetAddressLabel"))
              .build();
        }
      }
    }
    return getSetAddressLabelMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.WalletOuterClass.ListWalletsRequest,
      pactus.WalletOuterClass.ListWalletsResponse> getListWalletsMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "ListWallets",
      requestType = pactus.WalletOuterClass.ListWalletsRequest.class,
      responseType = pactus.WalletOuterClass.ListWalletsResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.WalletOuterClass.ListWalletsRequest,
      pactus.WalletOuterClass.ListWalletsResponse> getListWalletsMethod() {
    io.grpc.MethodDescriptor<pactus.WalletOuterClass.ListWalletsRequest, pactus.WalletOuterClass.ListWalletsResponse> getListWalletsMethod;
    if ((getListWalletsMethod = WalletGrpc.getListWalletsMethod) == null) {
      synchronized (WalletGrpc.class) {
        if ((getListWalletsMethod = WalletGrpc.getListWalletsMethod) == null) {
          WalletGrpc.getListWalletsMethod = getListWalletsMethod =
              io.grpc.MethodDescriptor.<pactus.WalletOuterClass.ListWalletsRequest, pactus.WalletOuterClass.ListWalletsResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "ListWallets"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.WalletOuterClass.ListWalletsRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.WalletOuterClass.ListWalletsResponse.getDefaultInstance()))
              .setSchemaDescriptor(new WalletMethodDescriptorSupplier("ListWallets"))
              .build();
        }
      }
    }
    return getListWalletsMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.WalletOuterClass.GetWalletInfoRequest,
      pactus.WalletOuterClass.GetWalletInfoResponse> getGetWalletInfoMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetWalletInfo",
      requestType = pactus.WalletOuterClass.GetWalletInfoRequest.class,
      responseType = pactus.WalletOuterClass.GetWalletInfoResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.WalletOuterClass.GetWalletInfoRequest,
      pactus.WalletOuterClass.GetWalletInfoResponse> getGetWalletInfoMethod() {
    io.grpc.MethodDescriptor<pactus.WalletOuterClass.GetWalletInfoRequest, pactus.WalletOuterClass.GetWalletInfoResponse> getGetWalletInfoMethod;
    if ((getGetWalletInfoMethod = WalletGrpc.getGetWalletInfoMethod) == null) {
      synchronized (WalletGrpc.class) {
        if ((getGetWalletInfoMethod = WalletGrpc.getGetWalletInfoMethod) == null) {
          WalletGrpc.getGetWalletInfoMethod = getGetWalletInfoMethod =
              io.grpc.MethodDescriptor.<pactus.WalletOuterClass.GetWalletInfoRequest, pactus.WalletOuterClass.GetWalletInfoResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetWalletInfo"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.WalletOuterClass.GetWalletInfoRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.WalletOuterClass.GetWalletInfoResponse.getDefaultInstance()))
              .setSchemaDescriptor(new WalletMethodDescriptorSupplier("GetWalletInfo"))
              .build();
        }
      }
    }
    return getGetWalletInfoMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.WalletOuterClass.ListAddressesRequest,
      pactus.WalletOuterClass.ListAddressesResponse> getListAddressesMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "ListAddresses",
      requestType = pactus.WalletOuterClass.ListAddressesRequest.class,
      responseType = pactus.WalletOuterClass.ListAddressesResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.WalletOuterClass.ListAddressesRequest,
      pactus.WalletOuterClass.ListAddressesResponse> getListAddressesMethod() {
    io.grpc.MethodDescriptor<pactus.WalletOuterClass.ListAddressesRequest, pactus.WalletOuterClass.ListAddressesResponse> getListAddressesMethod;
    if ((getListAddressesMethod = WalletGrpc.getListAddressesMethod) == null) {
      synchronized (WalletGrpc.class) {
        if ((getListAddressesMethod = WalletGrpc.getListAddressesMethod) == null) {
          WalletGrpc.getListAddressesMethod = getListAddressesMethod =
              io.grpc.MethodDescriptor.<pactus.WalletOuterClass.ListAddressesRequest, pactus.WalletOuterClass.ListAddressesResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "ListAddresses"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.WalletOuterClass.ListAddressesRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.WalletOuterClass.ListAddressesResponse.getDefaultInstance()))
              .setSchemaDescriptor(new WalletMethodDescriptorSupplier("ListAddresses"))
              .build();
        }
      }
    }
    return getListAddressesMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static WalletStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<WalletStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<WalletStub>() {
        @java.lang.Override
        public WalletStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new WalletStub(channel, callOptions);
        }
      };
    return WalletStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports all types of calls on the service
   */
  public static WalletBlockingV2Stub newBlockingV2Stub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<WalletBlockingV2Stub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<WalletBlockingV2Stub>() {
        @java.lang.Override
        public WalletBlockingV2Stub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new WalletBlockingV2Stub(channel, callOptions);
        }
      };
    return WalletBlockingV2Stub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static WalletBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<WalletBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<WalletBlockingStub>() {
        @java.lang.Override
        public WalletBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new WalletBlockingStub(channel, callOptions);
        }
      };
    return WalletBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static WalletFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<WalletFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<WalletFutureStub>() {
        @java.lang.Override
        public WalletFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new WalletFutureStub(channel, callOptions);
        }
      };
    return WalletFutureStub.newStub(factory, channel);
  }

  /**
   * <pre>
   * Wallet service provides RPC methods for wallet management operations.
   * </pre>
   */
  public interface AsyncService {

    /**
     * <pre>
     * CreateWallet creates a new wallet with the specified parameters.
     * </pre>
     */
    default void createWallet(pactus.WalletOuterClass.CreateWalletRequest request,
        io.grpc.stub.StreamObserver<pactus.WalletOuterClass.CreateWalletResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getCreateWalletMethod(), responseObserver);
    }

    /**
     * <pre>
     * RestoreWallet restores an existing wallet with the given mnemonic.
     * </pre>
     */
    default void restoreWallet(pactus.WalletOuterClass.RestoreWalletRequest request,
        io.grpc.stub.StreamObserver<pactus.WalletOuterClass.RestoreWalletResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getRestoreWalletMethod(), responseObserver);
    }

    /**
     * <pre>
     * LoadWallet loads an existing wallet with the given name.
     * </pre>
     */
    default void loadWallet(pactus.WalletOuterClass.LoadWalletRequest request,
        io.grpc.stub.StreamObserver<pactus.WalletOuterClass.LoadWalletResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getLoadWalletMethod(), responseObserver);
    }

    /**
     * <pre>
     * UnloadWallet unloads a currently loaded wallet with the specified name.
     * </pre>
     */
    default void unloadWallet(pactus.WalletOuterClass.UnloadWalletRequest request,
        io.grpc.stub.StreamObserver<pactus.WalletOuterClass.UnloadWalletResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getUnloadWalletMethod(), responseObserver);
    }

    /**
     * <pre>
     * GetTotalBalance returns the total available balance of the wallet.
     * </pre>
     */
    default void getTotalBalance(pactus.WalletOuterClass.GetTotalBalanceRequest request,
        io.grpc.stub.StreamObserver<pactus.WalletOuterClass.GetTotalBalanceResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetTotalBalanceMethod(), responseObserver);
    }

    /**
     * <pre>
     * SignRawTransaction signs a raw transaction for a specified wallet.
     * </pre>
     */
    default void signRawTransaction(pactus.WalletOuterClass.SignRawTransactionRequest request,
        io.grpc.stub.StreamObserver<pactus.WalletOuterClass.SignRawTransactionResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getSignRawTransactionMethod(), responseObserver);
    }

    /**
     * <pre>
     * GetValidatorAddress retrieves the validator address associated with a public key.
     * Deprecated: Will move into utils.
     * </pre>
     */
    default void getValidatorAddress(pactus.WalletOuterClass.GetValidatorAddressRequest request,
        io.grpc.stub.StreamObserver<pactus.WalletOuterClass.GetValidatorAddressResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetValidatorAddressMethod(), responseObserver);
    }

    /**
     * <pre>
     * GetNewAddress generates a new address for the specified wallet.
     * </pre>
     */
    default void getNewAddress(pactus.WalletOuterClass.GetNewAddressRequest request,
        io.grpc.stub.StreamObserver<pactus.WalletOuterClass.GetNewAddressResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetNewAddressMethod(), responseObserver);
    }

    /**
     * <pre>
     * GetAddressHistory retrieves the transaction history of an address.
     * </pre>
     */
    default void getAddressHistory(pactus.WalletOuterClass.GetAddressHistoryRequest request,
        io.grpc.stub.StreamObserver<pactus.WalletOuterClass.GetAddressHistoryResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetAddressHistoryMethod(), responseObserver);
    }

    /**
     * <pre>
     * SignMessage signs an arbitrary message using a wallet's private key.
     * </pre>
     */
    default void signMessage(pactus.WalletOuterClass.SignMessageRequest request,
        io.grpc.stub.StreamObserver<pactus.WalletOuterClass.SignMessageResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getSignMessageMethod(), responseObserver);
    }

    /**
     * <pre>
     * GetTotalStake returns the total stake amount in the wallet.
     * </pre>
     */
    default void getTotalStake(pactus.WalletOuterClass.GetTotalStakeRequest request,
        io.grpc.stub.StreamObserver<pactus.WalletOuterClass.GetTotalStakeResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetTotalStakeMethod(), responseObserver);
    }

    /**
     * <pre>
     * GetAddressInfo returns detailed information about a specific address.
     * </pre>
     */
    default void getAddressInfo(pactus.WalletOuterClass.GetAddressInfoRequest request,
        io.grpc.stub.StreamObserver<pactus.WalletOuterClass.GetAddressInfoResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetAddressInfoMethod(), responseObserver);
    }

    /**
     * <pre>
     * SetAddressLabel sets or updates the label for a given address.
     * </pre>
     */
    default void setAddressLabel(pactus.WalletOuterClass.SetAddressLabelRequest request,
        io.grpc.stub.StreamObserver<pactus.WalletOuterClass.SetAddressLabelResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getSetAddressLabelMethod(), responseObserver);
    }

    /**
     * <pre>
     * ListWallets returns a list of all available wallets.
     * </pre>
     */
    default void listWallets(pactus.WalletOuterClass.ListWalletsRequest request,
        io.grpc.stub.StreamObserver<pactus.WalletOuterClass.ListWalletsResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getListWalletsMethod(), responseObserver);
    }

    /**
     * <pre>
     * GetWalletInfo returns detailed information about a specific wallet.
     * </pre>
     */
    default void getWalletInfo(pactus.WalletOuterClass.GetWalletInfoRequest request,
        io.grpc.stub.StreamObserver<pactus.WalletOuterClass.GetWalletInfoResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetWalletInfoMethod(), responseObserver);
    }

    /**
     * <pre>
     * ListAddresses returns all addresses in the specified wallet.
     * </pre>
     */
    default void listAddresses(pactus.WalletOuterClass.ListAddressesRequest request,
        io.grpc.stub.StreamObserver<pactus.WalletOuterClass.ListAddressesResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getListAddressesMethod(), responseObserver);
    }
  }

  /**
   * Base class for the server implementation of the service Wallet.
   * <pre>
   * Wallet service provides RPC methods for wallet management operations.
   * </pre>
   */
  public static abstract class WalletImplBase
      implements io.grpc.BindableService, AsyncService {

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return WalletGrpc.bindService(this);
    }
  }

  /**
   * A stub to allow clients to do asynchronous rpc calls to service Wallet.
   * <pre>
   * Wallet service provides RPC methods for wallet management operations.
   * </pre>
   */
  public static final class WalletStub
      extends io.grpc.stub.AbstractAsyncStub<WalletStub> {
    private WalletStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected WalletStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new WalletStub(channel, callOptions);
    }

    /**
     * <pre>
     * CreateWallet creates a new wallet with the specified parameters.
     * </pre>
     */
    public void createWallet(pactus.WalletOuterClass.CreateWalletRequest request,
        io.grpc.stub.StreamObserver<pactus.WalletOuterClass.CreateWalletResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getCreateWalletMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * RestoreWallet restores an existing wallet with the given mnemonic.
     * </pre>
     */
    public void restoreWallet(pactus.WalletOuterClass.RestoreWalletRequest request,
        io.grpc.stub.StreamObserver<pactus.WalletOuterClass.RestoreWalletResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getRestoreWalletMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * LoadWallet loads an existing wallet with the given name.
     * </pre>
     */
    public void loadWallet(pactus.WalletOuterClass.LoadWalletRequest request,
        io.grpc.stub.StreamObserver<pactus.WalletOuterClass.LoadWalletResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getLoadWalletMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * UnloadWallet unloads a currently loaded wallet with the specified name.
     * </pre>
     */
    public void unloadWallet(pactus.WalletOuterClass.UnloadWalletRequest request,
        io.grpc.stub.StreamObserver<pactus.WalletOuterClass.UnloadWalletResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getUnloadWalletMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * GetTotalBalance returns the total available balance of the wallet.
     * </pre>
     */
    public void getTotalBalance(pactus.WalletOuterClass.GetTotalBalanceRequest request,
        io.grpc.stub.StreamObserver<pactus.WalletOuterClass.GetTotalBalanceResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetTotalBalanceMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * SignRawTransaction signs a raw transaction for a specified wallet.
     * </pre>
     */
    public void signRawTransaction(pactus.WalletOuterClass.SignRawTransactionRequest request,
        io.grpc.stub.StreamObserver<pactus.WalletOuterClass.SignRawTransactionResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getSignRawTransactionMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * GetValidatorAddress retrieves the validator address associated with a public key.
     * Deprecated: Will move into utils.
     * </pre>
     */
    public void getValidatorAddress(pactus.WalletOuterClass.GetValidatorAddressRequest request,
        io.grpc.stub.StreamObserver<pactus.WalletOuterClass.GetValidatorAddressResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetValidatorAddressMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * GetNewAddress generates a new address for the specified wallet.
     * </pre>
     */
    public void getNewAddress(pactus.WalletOuterClass.GetNewAddressRequest request,
        io.grpc.stub.StreamObserver<pactus.WalletOuterClass.GetNewAddressResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetNewAddressMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * GetAddressHistory retrieves the transaction history of an address.
     * </pre>
     */
    public void getAddressHistory(pactus.WalletOuterClass.GetAddressHistoryRequest request,
        io.grpc.stub.StreamObserver<pactus.WalletOuterClass.GetAddressHistoryResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetAddressHistoryMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * SignMessage signs an arbitrary message using a wallet's private key.
     * </pre>
     */
    public void signMessage(pactus.WalletOuterClass.SignMessageRequest request,
        io.grpc.stub.StreamObserver<pactus.WalletOuterClass.SignMessageResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getSignMessageMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * GetTotalStake returns the total stake amount in the wallet.
     * </pre>
     */
    public void getTotalStake(pactus.WalletOuterClass.GetTotalStakeRequest request,
        io.grpc.stub.StreamObserver<pactus.WalletOuterClass.GetTotalStakeResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetTotalStakeMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * GetAddressInfo returns detailed information about a specific address.
     * </pre>
     */
    public void getAddressInfo(pactus.WalletOuterClass.GetAddressInfoRequest request,
        io.grpc.stub.StreamObserver<pactus.WalletOuterClass.GetAddressInfoResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetAddressInfoMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * SetAddressLabel sets or updates the label for a given address.
     * </pre>
     */
    public void setAddressLabel(pactus.WalletOuterClass.SetAddressLabelRequest request,
        io.grpc.stub.StreamObserver<pactus.WalletOuterClass.SetAddressLabelResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getSetAddressLabelMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * ListWallets returns a list of all available wallets.
     * </pre>
     */
    public void listWallets(pactus.WalletOuterClass.ListWalletsRequest request,
        io.grpc.stub.StreamObserver<pactus.WalletOuterClass.ListWalletsResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getListWalletsMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * GetWalletInfo returns detailed information about a specific wallet.
     * </pre>
     */
    public void getWalletInfo(pactus.WalletOuterClass.GetWalletInfoRequest request,
        io.grpc.stub.StreamObserver<pactus.WalletOuterClass.GetWalletInfoResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetWalletInfoMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * ListAddresses returns all addresses in the specified wallet.
     * </pre>
     */
    public void listAddresses(pactus.WalletOuterClass.ListAddressesRequest request,
        io.grpc.stub.StreamObserver<pactus.WalletOuterClass.ListAddressesResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getListAddressesMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   * A stub to allow clients to do synchronous rpc calls to service Wallet.
   * <pre>
   * Wallet service provides RPC methods for wallet management operations.
   * </pre>
   */
  public static final class WalletBlockingV2Stub
      extends io.grpc.stub.AbstractBlockingStub<WalletBlockingV2Stub> {
    private WalletBlockingV2Stub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected WalletBlockingV2Stub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new WalletBlockingV2Stub(channel, callOptions);
    }

    /**
     * <pre>
     * CreateWallet creates a new wallet with the specified parameters.
     * </pre>
     */
    public pactus.WalletOuterClass.CreateWalletResponse createWallet(pactus.WalletOuterClass.CreateWalletRequest request) throws io.grpc.StatusException {
      return io.grpc.stub.ClientCalls.blockingV2UnaryCall(
          getChannel(), getCreateWalletMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * RestoreWallet restores an existing wallet with the given mnemonic.
     * </pre>
     */
    public pactus.WalletOuterClass.RestoreWalletResponse restoreWallet(pactus.WalletOuterClass.RestoreWalletRequest request) throws io.grpc.StatusException {
      return io.grpc.stub.ClientCalls.blockingV2UnaryCall(
          getChannel(), getRestoreWalletMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * LoadWallet loads an existing wallet with the given name.
     * </pre>
     */
    public pactus.WalletOuterClass.LoadWalletResponse loadWallet(pactus.WalletOuterClass.LoadWalletRequest request) throws io.grpc.StatusException {
      return io.grpc.stub.ClientCalls.blockingV2UnaryCall(
          getChannel(), getLoadWalletMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * UnloadWallet unloads a currently loaded wallet with the specified name.
     * </pre>
     */
    public pactus.WalletOuterClass.UnloadWalletResponse unloadWallet(pactus.WalletOuterClass.UnloadWalletRequest request) throws io.grpc.StatusException {
      return io.grpc.stub.ClientCalls.blockingV2UnaryCall(
          getChannel(), getUnloadWalletMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetTotalBalance returns the total available balance of the wallet.
     * </pre>
     */
    public pactus.WalletOuterClass.GetTotalBalanceResponse getTotalBalance(pactus.WalletOuterClass.GetTotalBalanceRequest request) throws io.grpc.StatusException {
      return io.grpc.stub.ClientCalls.blockingV2UnaryCall(
          getChannel(), getGetTotalBalanceMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * SignRawTransaction signs a raw transaction for a specified wallet.
     * </pre>
     */
    public pactus.WalletOuterClass.SignRawTransactionResponse signRawTransaction(pactus.WalletOuterClass.SignRawTransactionRequest request) throws io.grpc.StatusException {
      return io.grpc.stub.ClientCalls.blockingV2UnaryCall(
          getChannel(), getSignRawTransactionMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetValidatorAddress retrieves the validator address associated with a public key.
     * Deprecated: Will move into utils.
     * </pre>
     */
    public pactus.WalletOuterClass.GetValidatorAddressResponse getValidatorAddress(pactus.WalletOuterClass.GetValidatorAddressRequest request) throws io.grpc.StatusException {
      return io.grpc.stub.ClientCalls.blockingV2UnaryCall(
          getChannel(), getGetValidatorAddressMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetNewAddress generates a new address for the specified wallet.
     * </pre>
     */
    public pactus.WalletOuterClass.GetNewAddressResponse getNewAddress(pactus.WalletOuterClass.GetNewAddressRequest request) throws io.grpc.StatusException {
      return io.grpc.stub.ClientCalls.blockingV2UnaryCall(
          getChannel(), getGetNewAddressMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetAddressHistory retrieves the transaction history of an address.
     * </pre>
     */
    public pactus.WalletOuterClass.GetAddressHistoryResponse getAddressHistory(pactus.WalletOuterClass.GetAddressHistoryRequest request) throws io.grpc.StatusException {
      return io.grpc.stub.ClientCalls.blockingV2UnaryCall(
          getChannel(), getGetAddressHistoryMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * SignMessage signs an arbitrary message using a wallet's private key.
     * </pre>
     */
    public pactus.WalletOuterClass.SignMessageResponse signMessage(pactus.WalletOuterClass.SignMessageRequest request) throws io.grpc.StatusException {
      return io.grpc.stub.ClientCalls.blockingV2UnaryCall(
          getChannel(), getSignMessageMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetTotalStake returns the total stake amount in the wallet.
     * </pre>
     */
    public pactus.WalletOuterClass.GetTotalStakeResponse getTotalStake(pactus.WalletOuterClass.GetTotalStakeRequest request) throws io.grpc.StatusException {
      return io.grpc.stub.ClientCalls.blockingV2UnaryCall(
          getChannel(), getGetTotalStakeMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetAddressInfo returns detailed information about a specific address.
     * </pre>
     */
    public pactus.WalletOuterClass.GetAddressInfoResponse getAddressInfo(pactus.WalletOuterClass.GetAddressInfoRequest request) throws io.grpc.StatusException {
      return io.grpc.stub.ClientCalls.blockingV2UnaryCall(
          getChannel(), getGetAddressInfoMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * SetAddressLabel sets or updates the label for a given address.
     * </pre>
     */
    public pactus.WalletOuterClass.SetAddressLabelResponse setAddressLabel(pactus.WalletOuterClass.SetAddressLabelRequest request) throws io.grpc.StatusException {
      return io.grpc.stub.ClientCalls.blockingV2UnaryCall(
          getChannel(), getSetAddressLabelMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * ListWallets returns a list of all available wallets.
     * </pre>
     */
    public pactus.WalletOuterClass.ListWalletsResponse listWallets(pactus.WalletOuterClass.ListWalletsRequest request) throws io.grpc.StatusException {
      return io.grpc.stub.ClientCalls.blockingV2UnaryCall(
          getChannel(), getListWalletsMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetWalletInfo returns detailed information about a specific wallet.
     * </pre>
     */
    public pactus.WalletOuterClass.GetWalletInfoResponse getWalletInfo(pactus.WalletOuterClass.GetWalletInfoRequest request) throws io.grpc.StatusException {
      return io.grpc.stub.ClientCalls.blockingV2UnaryCall(
          getChannel(), getGetWalletInfoMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * ListAddresses returns all addresses in the specified wallet.
     * </pre>
     */
    public pactus.WalletOuterClass.ListAddressesResponse listAddresses(pactus.WalletOuterClass.ListAddressesRequest request) throws io.grpc.StatusException {
      return io.grpc.stub.ClientCalls.blockingV2UnaryCall(
          getChannel(), getListAddressesMethod(), getCallOptions(), request);
    }
  }

  /**
   * A stub to allow clients to do limited synchronous rpc calls to service Wallet.
   * <pre>
   * Wallet service provides RPC methods for wallet management operations.
   * </pre>
   */
  public static final class WalletBlockingStub
      extends io.grpc.stub.AbstractBlockingStub<WalletBlockingStub> {
    private WalletBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected WalletBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new WalletBlockingStub(channel, callOptions);
    }

    /**
     * <pre>
     * CreateWallet creates a new wallet with the specified parameters.
     * </pre>
     */
    public pactus.WalletOuterClass.CreateWalletResponse createWallet(pactus.WalletOuterClass.CreateWalletRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getCreateWalletMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * RestoreWallet restores an existing wallet with the given mnemonic.
     * </pre>
     */
    public pactus.WalletOuterClass.RestoreWalletResponse restoreWallet(pactus.WalletOuterClass.RestoreWalletRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getRestoreWalletMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * LoadWallet loads an existing wallet with the given name.
     * </pre>
     */
    public pactus.WalletOuterClass.LoadWalletResponse loadWallet(pactus.WalletOuterClass.LoadWalletRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getLoadWalletMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * UnloadWallet unloads a currently loaded wallet with the specified name.
     * </pre>
     */
    public pactus.WalletOuterClass.UnloadWalletResponse unloadWallet(pactus.WalletOuterClass.UnloadWalletRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getUnloadWalletMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetTotalBalance returns the total available balance of the wallet.
     * </pre>
     */
    public pactus.WalletOuterClass.GetTotalBalanceResponse getTotalBalance(pactus.WalletOuterClass.GetTotalBalanceRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetTotalBalanceMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * SignRawTransaction signs a raw transaction for a specified wallet.
     * </pre>
     */
    public pactus.WalletOuterClass.SignRawTransactionResponse signRawTransaction(pactus.WalletOuterClass.SignRawTransactionRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getSignRawTransactionMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetValidatorAddress retrieves the validator address associated with a public key.
     * Deprecated: Will move into utils.
     * </pre>
     */
    public pactus.WalletOuterClass.GetValidatorAddressResponse getValidatorAddress(pactus.WalletOuterClass.GetValidatorAddressRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetValidatorAddressMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetNewAddress generates a new address for the specified wallet.
     * </pre>
     */
    public pactus.WalletOuterClass.GetNewAddressResponse getNewAddress(pactus.WalletOuterClass.GetNewAddressRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetNewAddressMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetAddressHistory retrieves the transaction history of an address.
     * </pre>
     */
    public pactus.WalletOuterClass.GetAddressHistoryResponse getAddressHistory(pactus.WalletOuterClass.GetAddressHistoryRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetAddressHistoryMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * SignMessage signs an arbitrary message using a wallet's private key.
     * </pre>
     */
    public pactus.WalletOuterClass.SignMessageResponse signMessage(pactus.WalletOuterClass.SignMessageRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getSignMessageMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetTotalStake returns the total stake amount in the wallet.
     * </pre>
     */
    public pactus.WalletOuterClass.GetTotalStakeResponse getTotalStake(pactus.WalletOuterClass.GetTotalStakeRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetTotalStakeMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetAddressInfo returns detailed information about a specific address.
     * </pre>
     */
    public pactus.WalletOuterClass.GetAddressInfoResponse getAddressInfo(pactus.WalletOuterClass.GetAddressInfoRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetAddressInfoMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * SetAddressLabel sets or updates the label for a given address.
     * </pre>
     */
    public pactus.WalletOuterClass.SetAddressLabelResponse setAddressLabel(pactus.WalletOuterClass.SetAddressLabelRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getSetAddressLabelMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * ListWallets returns a list of all available wallets.
     * </pre>
     */
    public pactus.WalletOuterClass.ListWalletsResponse listWallets(pactus.WalletOuterClass.ListWalletsRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getListWalletsMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetWalletInfo returns detailed information about a specific wallet.
     * </pre>
     */
    public pactus.WalletOuterClass.GetWalletInfoResponse getWalletInfo(pactus.WalletOuterClass.GetWalletInfoRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetWalletInfoMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * ListAddresses returns all addresses in the specified wallet.
     * </pre>
     */
    public pactus.WalletOuterClass.ListAddressesResponse listAddresses(pactus.WalletOuterClass.ListAddressesRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getListAddressesMethod(), getCallOptions(), request);
    }
  }

  /**
   * A stub to allow clients to do ListenableFuture-style rpc calls to service Wallet.
   * <pre>
   * Wallet service provides RPC methods for wallet management operations.
   * </pre>
   */
  public static final class WalletFutureStub
      extends io.grpc.stub.AbstractFutureStub<WalletFutureStub> {
    private WalletFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected WalletFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new WalletFutureStub(channel, callOptions);
    }

    /**
     * <pre>
     * CreateWallet creates a new wallet with the specified parameters.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.WalletOuterClass.CreateWalletResponse> createWallet(
        pactus.WalletOuterClass.CreateWalletRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getCreateWalletMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * RestoreWallet restores an existing wallet with the given mnemonic.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.WalletOuterClass.RestoreWalletResponse> restoreWallet(
        pactus.WalletOuterClass.RestoreWalletRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getRestoreWalletMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * LoadWallet loads an existing wallet with the given name.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.WalletOuterClass.LoadWalletResponse> loadWallet(
        pactus.WalletOuterClass.LoadWalletRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getLoadWalletMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * UnloadWallet unloads a currently loaded wallet with the specified name.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.WalletOuterClass.UnloadWalletResponse> unloadWallet(
        pactus.WalletOuterClass.UnloadWalletRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getUnloadWalletMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * GetTotalBalance returns the total available balance of the wallet.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.WalletOuterClass.GetTotalBalanceResponse> getTotalBalance(
        pactus.WalletOuterClass.GetTotalBalanceRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetTotalBalanceMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * SignRawTransaction signs a raw transaction for a specified wallet.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.WalletOuterClass.SignRawTransactionResponse> signRawTransaction(
        pactus.WalletOuterClass.SignRawTransactionRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getSignRawTransactionMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * GetValidatorAddress retrieves the validator address associated with a public key.
     * Deprecated: Will move into utils.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.WalletOuterClass.GetValidatorAddressResponse> getValidatorAddress(
        pactus.WalletOuterClass.GetValidatorAddressRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetValidatorAddressMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * GetNewAddress generates a new address for the specified wallet.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.WalletOuterClass.GetNewAddressResponse> getNewAddress(
        pactus.WalletOuterClass.GetNewAddressRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetNewAddressMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * GetAddressHistory retrieves the transaction history of an address.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.WalletOuterClass.GetAddressHistoryResponse> getAddressHistory(
        pactus.WalletOuterClass.GetAddressHistoryRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetAddressHistoryMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * SignMessage signs an arbitrary message using a wallet's private key.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.WalletOuterClass.SignMessageResponse> signMessage(
        pactus.WalletOuterClass.SignMessageRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getSignMessageMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * GetTotalStake returns the total stake amount in the wallet.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.WalletOuterClass.GetTotalStakeResponse> getTotalStake(
        pactus.WalletOuterClass.GetTotalStakeRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetTotalStakeMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * GetAddressInfo returns detailed information about a specific address.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.WalletOuterClass.GetAddressInfoResponse> getAddressInfo(
        pactus.WalletOuterClass.GetAddressInfoRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetAddressInfoMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * SetAddressLabel sets or updates the label for a given address.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.WalletOuterClass.SetAddressLabelResponse> setAddressLabel(
        pactus.WalletOuterClass.SetAddressLabelRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getSetAddressLabelMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * ListWallets returns a list of all available wallets.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.WalletOuterClass.ListWalletsResponse> listWallets(
        pactus.WalletOuterClass.ListWalletsRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getListWalletsMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * GetWalletInfo returns detailed information about a specific wallet.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.WalletOuterClass.GetWalletInfoResponse> getWalletInfo(
        pactus.WalletOuterClass.GetWalletInfoRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetWalletInfoMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * ListAddresses returns all addresses in the specified wallet.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.WalletOuterClass.ListAddressesResponse> listAddresses(
        pactus.WalletOuterClass.ListAddressesRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getListAddressesMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_CREATE_WALLET = 0;
  private static final int METHODID_RESTORE_WALLET = 1;
  private static final int METHODID_LOAD_WALLET = 2;
  private static final int METHODID_UNLOAD_WALLET = 3;
  private static final int METHODID_GET_TOTAL_BALANCE = 4;
  private static final int METHODID_SIGN_RAW_TRANSACTION = 5;
  private static final int METHODID_GET_VALIDATOR_ADDRESS = 6;
  private static final int METHODID_GET_NEW_ADDRESS = 7;
  private static final int METHODID_GET_ADDRESS_HISTORY = 8;
  private static final int METHODID_SIGN_MESSAGE = 9;
  private static final int METHODID_GET_TOTAL_STAKE = 10;
  private static final int METHODID_GET_ADDRESS_INFO = 11;
  private static final int METHODID_SET_ADDRESS_LABEL = 12;
  private static final int METHODID_LIST_WALLETS = 13;
  private static final int METHODID_GET_WALLET_INFO = 14;
  private static final int METHODID_LIST_ADDRESSES = 15;

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
        case METHODID_CREATE_WALLET:
          serviceImpl.createWallet((pactus.WalletOuterClass.CreateWalletRequest) request,
              (io.grpc.stub.StreamObserver<pactus.WalletOuterClass.CreateWalletResponse>) responseObserver);
          break;
        case METHODID_RESTORE_WALLET:
          serviceImpl.restoreWallet((pactus.WalletOuterClass.RestoreWalletRequest) request,
              (io.grpc.stub.StreamObserver<pactus.WalletOuterClass.RestoreWalletResponse>) responseObserver);
          break;
        case METHODID_LOAD_WALLET:
          serviceImpl.loadWallet((pactus.WalletOuterClass.LoadWalletRequest) request,
              (io.grpc.stub.StreamObserver<pactus.WalletOuterClass.LoadWalletResponse>) responseObserver);
          break;
        case METHODID_UNLOAD_WALLET:
          serviceImpl.unloadWallet((pactus.WalletOuterClass.UnloadWalletRequest) request,
              (io.grpc.stub.StreamObserver<pactus.WalletOuterClass.UnloadWalletResponse>) responseObserver);
          break;
        case METHODID_GET_TOTAL_BALANCE:
          serviceImpl.getTotalBalance((pactus.WalletOuterClass.GetTotalBalanceRequest) request,
              (io.grpc.stub.StreamObserver<pactus.WalletOuterClass.GetTotalBalanceResponse>) responseObserver);
          break;
        case METHODID_SIGN_RAW_TRANSACTION:
          serviceImpl.signRawTransaction((pactus.WalletOuterClass.SignRawTransactionRequest) request,
              (io.grpc.stub.StreamObserver<pactus.WalletOuterClass.SignRawTransactionResponse>) responseObserver);
          break;
        case METHODID_GET_VALIDATOR_ADDRESS:
          serviceImpl.getValidatorAddress((pactus.WalletOuterClass.GetValidatorAddressRequest) request,
              (io.grpc.stub.StreamObserver<pactus.WalletOuterClass.GetValidatorAddressResponse>) responseObserver);
          break;
        case METHODID_GET_NEW_ADDRESS:
          serviceImpl.getNewAddress((pactus.WalletOuterClass.GetNewAddressRequest) request,
              (io.grpc.stub.StreamObserver<pactus.WalletOuterClass.GetNewAddressResponse>) responseObserver);
          break;
        case METHODID_GET_ADDRESS_HISTORY:
          serviceImpl.getAddressHistory((pactus.WalletOuterClass.GetAddressHistoryRequest) request,
              (io.grpc.stub.StreamObserver<pactus.WalletOuterClass.GetAddressHistoryResponse>) responseObserver);
          break;
        case METHODID_SIGN_MESSAGE:
          serviceImpl.signMessage((pactus.WalletOuterClass.SignMessageRequest) request,
              (io.grpc.stub.StreamObserver<pactus.WalletOuterClass.SignMessageResponse>) responseObserver);
          break;
        case METHODID_GET_TOTAL_STAKE:
          serviceImpl.getTotalStake((pactus.WalletOuterClass.GetTotalStakeRequest) request,
              (io.grpc.stub.StreamObserver<pactus.WalletOuterClass.GetTotalStakeResponse>) responseObserver);
          break;
        case METHODID_GET_ADDRESS_INFO:
          serviceImpl.getAddressInfo((pactus.WalletOuterClass.GetAddressInfoRequest) request,
              (io.grpc.stub.StreamObserver<pactus.WalletOuterClass.GetAddressInfoResponse>) responseObserver);
          break;
        case METHODID_SET_ADDRESS_LABEL:
          serviceImpl.setAddressLabel((pactus.WalletOuterClass.SetAddressLabelRequest) request,
              (io.grpc.stub.StreamObserver<pactus.WalletOuterClass.SetAddressLabelResponse>) responseObserver);
          break;
        case METHODID_LIST_WALLETS:
          serviceImpl.listWallets((pactus.WalletOuterClass.ListWalletsRequest) request,
              (io.grpc.stub.StreamObserver<pactus.WalletOuterClass.ListWalletsResponse>) responseObserver);
          break;
        case METHODID_GET_WALLET_INFO:
          serviceImpl.getWalletInfo((pactus.WalletOuterClass.GetWalletInfoRequest) request,
              (io.grpc.stub.StreamObserver<pactus.WalletOuterClass.GetWalletInfoResponse>) responseObserver);
          break;
        case METHODID_LIST_ADDRESSES:
          serviceImpl.listAddresses((pactus.WalletOuterClass.ListAddressesRequest) request,
              (io.grpc.stub.StreamObserver<pactus.WalletOuterClass.ListAddressesResponse>) responseObserver);
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
          getCreateWalletMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.WalletOuterClass.CreateWalletRequest,
              pactus.WalletOuterClass.CreateWalletResponse>(
                service, METHODID_CREATE_WALLET)))
        .addMethod(
          getRestoreWalletMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.WalletOuterClass.RestoreWalletRequest,
              pactus.WalletOuterClass.RestoreWalletResponse>(
                service, METHODID_RESTORE_WALLET)))
        .addMethod(
          getLoadWalletMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.WalletOuterClass.LoadWalletRequest,
              pactus.WalletOuterClass.LoadWalletResponse>(
                service, METHODID_LOAD_WALLET)))
        .addMethod(
          getUnloadWalletMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.WalletOuterClass.UnloadWalletRequest,
              pactus.WalletOuterClass.UnloadWalletResponse>(
                service, METHODID_UNLOAD_WALLET)))
        .addMethod(
          getGetTotalBalanceMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.WalletOuterClass.GetTotalBalanceRequest,
              pactus.WalletOuterClass.GetTotalBalanceResponse>(
                service, METHODID_GET_TOTAL_BALANCE)))
        .addMethod(
          getSignRawTransactionMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.WalletOuterClass.SignRawTransactionRequest,
              pactus.WalletOuterClass.SignRawTransactionResponse>(
                service, METHODID_SIGN_RAW_TRANSACTION)))
        .addMethod(
          getGetValidatorAddressMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.WalletOuterClass.GetValidatorAddressRequest,
              pactus.WalletOuterClass.GetValidatorAddressResponse>(
                service, METHODID_GET_VALIDATOR_ADDRESS)))
        .addMethod(
          getGetNewAddressMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.WalletOuterClass.GetNewAddressRequest,
              pactus.WalletOuterClass.GetNewAddressResponse>(
                service, METHODID_GET_NEW_ADDRESS)))
        .addMethod(
          getGetAddressHistoryMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.WalletOuterClass.GetAddressHistoryRequest,
              pactus.WalletOuterClass.GetAddressHistoryResponse>(
                service, METHODID_GET_ADDRESS_HISTORY)))
        .addMethod(
          getSignMessageMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.WalletOuterClass.SignMessageRequest,
              pactus.WalletOuterClass.SignMessageResponse>(
                service, METHODID_SIGN_MESSAGE)))
        .addMethod(
          getGetTotalStakeMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.WalletOuterClass.GetTotalStakeRequest,
              pactus.WalletOuterClass.GetTotalStakeResponse>(
                service, METHODID_GET_TOTAL_STAKE)))
        .addMethod(
          getGetAddressInfoMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.WalletOuterClass.GetAddressInfoRequest,
              pactus.WalletOuterClass.GetAddressInfoResponse>(
                service, METHODID_GET_ADDRESS_INFO)))
        .addMethod(
          getSetAddressLabelMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.WalletOuterClass.SetAddressLabelRequest,
              pactus.WalletOuterClass.SetAddressLabelResponse>(
                service, METHODID_SET_ADDRESS_LABEL)))
        .addMethod(
          getListWalletsMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.WalletOuterClass.ListWalletsRequest,
              pactus.WalletOuterClass.ListWalletsResponse>(
                service, METHODID_LIST_WALLETS)))
        .addMethod(
          getGetWalletInfoMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.WalletOuterClass.GetWalletInfoRequest,
              pactus.WalletOuterClass.GetWalletInfoResponse>(
                service, METHODID_GET_WALLET_INFO)))
        .addMethod(
          getListAddressesMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.WalletOuterClass.ListAddressesRequest,
              pactus.WalletOuterClass.ListAddressesResponse>(
                service, METHODID_LIST_ADDRESSES)))
        .build();
  }

  private static abstract class WalletBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    WalletBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return pactus.WalletOuterClass.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("Wallet");
    }
  }

  private static final class WalletFileDescriptorSupplier
      extends WalletBaseDescriptorSupplier {
    WalletFileDescriptorSupplier() {}
  }

  private static final class WalletMethodDescriptorSupplier
      extends WalletBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final java.lang.String methodName;

    WalletMethodDescriptorSupplier(java.lang.String methodName) {
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
      synchronized (WalletGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new WalletFileDescriptorSupplier())
              .addMethod(getCreateWalletMethod())
              .addMethod(getRestoreWalletMethod())
              .addMethod(getLoadWalletMethod())
              .addMethod(getUnloadWalletMethod())
              .addMethod(getGetTotalBalanceMethod())
              .addMethod(getSignRawTransactionMethod())
              .addMethod(getGetValidatorAddressMethod())
              .addMethod(getGetNewAddressMethod())
              .addMethod(getGetAddressHistoryMethod())
              .addMethod(getSignMessageMethod())
              .addMethod(getGetTotalStakeMethod())
              .addMethod(getGetAddressInfoMethod())
              .addMethod(getSetAddressLabelMethod())
              .addMethod(getListWalletsMethod())
              .addMethod(getGetWalletInfoMethod())
              .addMethod(getListAddressesMethod())
              .build();
        }
      }
    }
    return result;
  }
}
