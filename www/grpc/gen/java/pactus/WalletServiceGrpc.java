package pactus;

import static io.grpc.MethodDescriptor.generateFullMethodName;

/**
 * <pre>
 * WalletService provides RPC methods for wallet management operations.
 * </pre>
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.71.0)",
    comments = "Source: wallet.proto")
@io.grpc.stub.annotations.GrpcGenerated
public final class WalletServiceGrpc {

  private WalletServiceGrpc() {}

  public static final java.lang.String SERVICE_NAME = "pactus.WalletService";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<pactus.Wallet.CreateWalletRequest,
      pactus.Wallet.CreateWalletResponse> getCreateWalletMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "CreateWallet",
      requestType = pactus.Wallet.CreateWalletRequest.class,
      responseType = pactus.Wallet.CreateWalletResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.Wallet.CreateWalletRequest,
      pactus.Wallet.CreateWalletResponse> getCreateWalletMethod() {
    io.grpc.MethodDescriptor<pactus.Wallet.CreateWalletRequest, pactus.Wallet.CreateWalletResponse> getCreateWalletMethod;
    if ((getCreateWalletMethod = WalletServiceGrpc.getCreateWalletMethod) == null) {
      synchronized (WalletServiceGrpc.class) {
        if ((getCreateWalletMethod = WalletServiceGrpc.getCreateWalletMethod) == null) {
          WalletServiceGrpc.getCreateWalletMethod = getCreateWalletMethod =
              io.grpc.MethodDescriptor.<pactus.Wallet.CreateWalletRequest, pactus.Wallet.CreateWalletResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "CreateWallet"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Wallet.CreateWalletRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Wallet.CreateWalletResponse.getDefaultInstance()))
              .setSchemaDescriptor(new WalletServiceMethodDescriptorSupplier("CreateWallet"))
              .build();
        }
      }
    }
    return getCreateWalletMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.Wallet.RestoreWalletRequest,
      pactus.Wallet.RestoreWalletResponse> getRestoreWalletMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "RestoreWallet",
      requestType = pactus.Wallet.RestoreWalletRequest.class,
      responseType = pactus.Wallet.RestoreWalletResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.Wallet.RestoreWalletRequest,
      pactus.Wallet.RestoreWalletResponse> getRestoreWalletMethod() {
    io.grpc.MethodDescriptor<pactus.Wallet.RestoreWalletRequest, pactus.Wallet.RestoreWalletResponse> getRestoreWalletMethod;
    if ((getRestoreWalletMethod = WalletServiceGrpc.getRestoreWalletMethod) == null) {
      synchronized (WalletServiceGrpc.class) {
        if ((getRestoreWalletMethod = WalletServiceGrpc.getRestoreWalletMethod) == null) {
          WalletServiceGrpc.getRestoreWalletMethod = getRestoreWalletMethod =
              io.grpc.MethodDescriptor.<pactus.Wallet.RestoreWalletRequest, pactus.Wallet.RestoreWalletResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "RestoreWallet"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Wallet.RestoreWalletRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Wallet.RestoreWalletResponse.getDefaultInstance()))
              .setSchemaDescriptor(new WalletServiceMethodDescriptorSupplier("RestoreWallet"))
              .build();
        }
      }
    }
    return getRestoreWalletMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.Wallet.LoadWalletRequest,
      pactus.Wallet.LoadWalletResponse> getLoadWalletMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "LoadWallet",
      requestType = pactus.Wallet.LoadWalletRequest.class,
      responseType = pactus.Wallet.LoadWalletResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.Wallet.LoadWalletRequest,
      pactus.Wallet.LoadWalletResponse> getLoadWalletMethod() {
    io.grpc.MethodDescriptor<pactus.Wallet.LoadWalletRequest, pactus.Wallet.LoadWalletResponse> getLoadWalletMethod;
    if ((getLoadWalletMethod = WalletServiceGrpc.getLoadWalletMethod) == null) {
      synchronized (WalletServiceGrpc.class) {
        if ((getLoadWalletMethod = WalletServiceGrpc.getLoadWalletMethod) == null) {
          WalletServiceGrpc.getLoadWalletMethod = getLoadWalletMethod =
              io.grpc.MethodDescriptor.<pactus.Wallet.LoadWalletRequest, pactus.Wallet.LoadWalletResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "LoadWallet"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Wallet.LoadWalletRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Wallet.LoadWalletResponse.getDefaultInstance()))
              .setSchemaDescriptor(new WalletServiceMethodDescriptorSupplier("LoadWallet"))
              .build();
        }
      }
    }
    return getLoadWalletMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.Wallet.UnloadWalletRequest,
      pactus.Wallet.UnloadWalletResponse> getUnloadWalletMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "UnloadWallet",
      requestType = pactus.Wallet.UnloadWalletRequest.class,
      responseType = pactus.Wallet.UnloadWalletResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.Wallet.UnloadWalletRequest,
      pactus.Wallet.UnloadWalletResponse> getUnloadWalletMethod() {
    io.grpc.MethodDescriptor<pactus.Wallet.UnloadWalletRequest, pactus.Wallet.UnloadWalletResponse> getUnloadWalletMethod;
    if ((getUnloadWalletMethod = WalletServiceGrpc.getUnloadWalletMethod) == null) {
      synchronized (WalletServiceGrpc.class) {
        if ((getUnloadWalletMethod = WalletServiceGrpc.getUnloadWalletMethod) == null) {
          WalletServiceGrpc.getUnloadWalletMethod = getUnloadWalletMethod =
              io.grpc.MethodDescriptor.<pactus.Wallet.UnloadWalletRequest, pactus.Wallet.UnloadWalletResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "UnloadWallet"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Wallet.UnloadWalletRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Wallet.UnloadWalletResponse.getDefaultInstance()))
              .setSchemaDescriptor(new WalletServiceMethodDescriptorSupplier("UnloadWallet"))
              .build();
        }
      }
    }
    return getUnloadWalletMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.Wallet.GetTotalBalanceRequest,
      pactus.Wallet.GetTotalBalanceResponse> getGetTotalBalanceMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetTotalBalance",
      requestType = pactus.Wallet.GetTotalBalanceRequest.class,
      responseType = pactus.Wallet.GetTotalBalanceResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.Wallet.GetTotalBalanceRequest,
      pactus.Wallet.GetTotalBalanceResponse> getGetTotalBalanceMethod() {
    io.grpc.MethodDescriptor<pactus.Wallet.GetTotalBalanceRequest, pactus.Wallet.GetTotalBalanceResponse> getGetTotalBalanceMethod;
    if ((getGetTotalBalanceMethod = WalletServiceGrpc.getGetTotalBalanceMethod) == null) {
      synchronized (WalletServiceGrpc.class) {
        if ((getGetTotalBalanceMethod = WalletServiceGrpc.getGetTotalBalanceMethod) == null) {
          WalletServiceGrpc.getGetTotalBalanceMethod = getGetTotalBalanceMethod =
              io.grpc.MethodDescriptor.<pactus.Wallet.GetTotalBalanceRequest, pactus.Wallet.GetTotalBalanceResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetTotalBalance"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Wallet.GetTotalBalanceRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Wallet.GetTotalBalanceResponse.getDefaultInstance()))
              .setSchemaDescriptor(new WalletServiceMethodDescriptorSupplier("GetTotalBalance"))
              .build();
        }
      }
    }
    return getGetTotalBalanceMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.Wallet.SignRawTransactionRequest,
      pactus.Wallet.SignRawTransactionResponse> getSignRawTransactionMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "SignRawTransaction",
      requestType = pactus.Wallet.SignRawTransactionRequest.class,
      responseType = pactus.Wallet.SignRawTransactionResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.Wallet.SignRawTransactionRequest,
      pactus.Wallet.SignRawTransactionResponse> getSignRawTransactionMethod() {
    io.grpc.MethodDescriptor<pactus.Wallet.SignRawTransactionRequest, pactus.Wallet.SignRawTransactionResponse> getSignRawTransactionMethod;
    if ((getSignRawTransactionMethod = WalletServiceGrpc.getSignRawTransactionMethod) == null) {
      synchronized (WalletServiceGrpc.class) {
        if ((getSignRawTransactionMethod = WalletServiceGrpc.getSignRawTransactionMethod) == null) {
          WalletServiceGrpc.getSignRawTransactionMethod = getSignRawTransactionMethod =
              io.grpc.MethodDescriptor.<pactus.Wallet.SignRawTransactionRequest, pactus.Wallet.SignRawTransactionResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "SignRawTransaction"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Wallet.SignRawTransactionRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Wallet.SignRawTransactionResponse.getDefaultInstance()))
              .setSchemaDescriptor(new WalletServiceMethodDescriptorSupplier("SignRawTransaction"))
              .build();
        }
      }
    }
    return getSignRawTransactionMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.Wallet.GetValidatorAddressRequest,
      pactus.Wallet.GetValidatorAddressResponse> getGetValidatorAddressMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetValidatorAddress",
      requestType = pactus.Wallet.GetValidatorAddressRequest.class,
      responseType = pactus.Wallet.GetValidatorAddressResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.Wallet.GetValidatorAddressRequest,
      pactus.Wallet.GetValidatorAddressResponse> getGetValidatorAddressMethod() {
    io.grpc.MethodDescriptor<pactus.Wallet.GetValidatorAddressRequest, pactus.Wallet.GetValidatorAddressResponse> getGetValidatorAddressMethod;
    if ((getGetValidatorAddressMethod = WalletServiceGrpc.getGetValidatorAddressMethod) == null) {
      synchronized (WalletServiceGrpc.class) {
        if ((getGetValidatorAddressMethod = WalletServiceGrpc.getGetValidatorAddressMethod) == null) {
          WalletServiceGrpc.getGetValidatorAddressMethod = getGetValidatorAddressMethod =
              io.grpc.MethodDescriptor.<pactus.Wallet.GetValidatorAddressRequest, pactus.Wallet.GetValidatorAddressResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetValidatorAddress"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Wallet.GetValidatorAddressRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Wallet.GetValidatorAddressResponse.getDefaultInstance()))
              .setSchemaDescriptor(new WalletServiceMethodDescriptorSupplier("GetValidatorAddress"))
              .build();
        }
      }
    }
    return getGetValidatorAddressMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.Wallet.GetNewAddressRequest,
      pactus.Wallet.GetNewAddressResponse> getGetNewAddressMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetNewAddress",
      requestType = pactus.Wallet.GetNewAddressRequest.class,
      responseType = pactus.Wallet.GetNewAddressResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.Wallet.GetNewAddressRequest,
      pactus.Wallet.GetNewAddressResponse> getGetNewAddressMethod() {
    io.grpc.MethodDescriptor<pactus.Wallet.GetNewAddressRequest, pactus.Wallet.GetNewAddressResponse> getGetNewAddressMethod;
    if ((getGetNewAddressMethod = WalletServiceGrpc.getGetNewAddressMethod) == null) {
      synchronized (WalletServiceGrpc.class) {
        if ((getGetNewAddressMethod = WalletServiceGrpc.getGetNewAddressMethod) == null) {
          WalletServiceGrpc.getGetNewAddressMethod = getGetNewAddressMethod =
              io.grpc.MethodDescriptor.<pactus.Wallet.GetNewAddressRequest, pactus.Wallet.GetNewAddressResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetNewAddress"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Wallet.GetNewAddressRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Wallet.GetNewAddressResponse.getDefaultInstance()))
              .setSchemaDescriptor(new WalletServiceMethodDescriptorSupplier("GetNewAddress"))
              .build();
        }
      }
    }
    return getGetNewAddressMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.Wallet.GetAddressHistoryRequest,
      pactus.Wallet.GetAddressHistoryResponse> getGetAddressHistoryMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetAddressHistory",
      requestType = pactus.Wallet.GetAddressHistoryRequest.class,
      responseType = pactus.Wallet.GetAddressHistoryResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.Wallet.GetAddressHistoryRequest,
      pactus.Wallet.GetAddressHistoryResponse> getGetAddressHistoryMethod() {
    io.grpc.MethodDescriptor<pactus.Wallet.GetAddressHistoryRequest, pactus.Wallet.GetAddressHistoryResponse> getGetAddressHistoryMethod;
    if ((getGetAddressHistoryMethod = WalletServiceGrpc.getGetAddressHistoryMethod) == null) {
      synchronized (WalletServiceGrpc.class) {
        if ((getGetAddressHistoryMethod = WalletServiceGrpc.getGetAddressHistoryMethod) == null) {
          WalletServiceGrpc.getGetAddressHistoryMethod = getGetAddressHistoryMethod =
              io.grpc.MethodDescriptor.<pactus.Wallet.GetAddressHistoryRequest, pactus.Wallet.GetAddressHistoryResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetAddressHistory"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Wallet.GetAddressHistoryRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Wallet.GetAddressHistoryResponse.getDefaultInstance()))
              .setSchemaDescriptor(new WalletServiceMethodDescriptorSupplier("GetAddressHistory"))
              .build();
        }
      }
    }
    return getGetAddressHistoryMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.Wallet.SignMessageRequest,
      pactus.Wallet.SignMessageResponse> getSignMessageMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "SignMessage",
      requestType = pactus.Wallet.SignMessageRequest.class,
      responseType = pactus.Wallet.SignMessageResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.Wallet.SignMessageRequest,
      pactus.Wallet.SignMessageResponse> getSignMessageMethod() {
    io.grpc.MethodDescriptor<pactus.Wallet.SignMessageRequest, pactus.Wallet.SignMessageResponse> getSignMessageMethod;
    if ((getSignMessageMethod = WalletServiceGrpc.getSignMessageMethod) == null) {
      synchronized (WalletServiceGrpc.class) {
        if ((getSignMessageMethod = WalletServiceGrpc.getSignMessageMethod) == null) {
          WalletServiceGrpc.getSignMessageMethod = getSignMessageMethod =
              io.grpc.MethodDescriptor.<pactus.Wallet.SignMessageRequest, pactus.Wallet.SignMessageResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "SignMessage"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Wallet.SignMessageRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Wallet.SignMessageResponse.getDefaultInstance()))
              .setSchemaDescriptor(new WalletServiceMethodDescriptorSupplier("SignMessage"))
              .build();
        }
      }
    }
    return getSignMessageMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.Wallet.GetTotalStakeRequest,
      pactus.Wallet.GetTotalStakeResponse> getGetTotalStakeMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetTotalStake",
      requestType = pactus.Wallet.GetTotalStakeRequest.class,
      responseType = pactus.Wallet.GetTotalStakeResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.Wallet.GetTotalStakeRequest,
      pactus.Wallet.GetTotalStakeResponse> getGetTotalStakeMethod() {
    io.grpc.MethodDescriptor<pactus.Wallet.GetTotalStakeRequest, pactus.Wallet.GetTotalStakeResponse> getGetTotalStakeMethod;
    if ((getGetTotalStakeMethod = WalletServiceGrpc.getGetTotalStakeMethod) == null) {
      synchronized (WalletServiceGrpc.class) {
        if ((getGetTotalStakeMethod = WalletServiceGrpc.getGetTotalStakeMethod) == null) {
          WalletServiceGrpc.getGetTotalStakeMethod = getGetTotalStakeMethod =
              io.grpc.MethodDescriptor.<pactus.Wallet.GetTotalStakeRequest, pactus.Wallet.GetTotalStakeResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetTotalStake"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Wallet.GetTotalStakeRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Wallet.GetTotalStakeResponse.getDefaultInstance()))
              .setSchemaDescriptor(new WalletServiceMethodDescriptorSupplier("GetTotalStake"))
              .build();
        }
      }
    }
    return getGetTotalStakeMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.Wallet.GetAddressInfoRequest,
      pactus.Wallet.GetAddressInfoResponse> getGetAddressInfoMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetAddressInfo",
      requestType = pactus.Wallet.GetAddressInfoRequest.class,
      responseType = pactus.Wallet.GetAddressInfoResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.Wallet.GetAddressInfoRequest,
      pactus.Wallet.GetAddressInfoResponse> getGetAddressInfoMethod() {
    io.grpc.MethodDescriptor<pactus.Wallet.GetAddressInfoRequest, pactus.Wallet.GetAddressInfoResponse> getGetAddressInfoMethod;
    if ((getGetAddressInfoMethod = WalletServiceGrpc.getGetAddressInfoMethod) == null) {
      synchronized (WalletServiceGrpc.class) {
        if ((getGetAddressInfoMethod = WalletServiceGrpc.getGetAddressInfoMethod) == null) {
          WalletServiceGrpc.getGetAddressInfoMethod = getGetAddressInfoMethod =
              io.grpc.MethodDescriptor.<pactus.Wallet.GetAddressInfoRequest, pactus.Wallet.GetAddressInfoResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetAddressInfo"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Wallet.GetAddressInfoRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Wallet.GetAddressInfoResponse.getDefaultInstance()))
              .setSchemaDescriptor(new WalletServiceMethodDescriptorSupplier("GetAddressInfo"))
              .build();
        }
      }
    }
    return getGetAddressInfoMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.Wallet.SetAddressLabelRequest,
      pactus.Wallet.SetAddressLabelResponse> getSetAddressLabelMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "SetAddressLabel",
      requestType = pactus.Wallet.SetAddressLabelRequest.class,
      responseType = pactus.Wallet.SetAddressLabelResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.Wallet.SetAddressLabelRequest,
      pactus.Wallet.SetAddressLabelResponse> getSetAddressLabelMethod() {
    io.grpc.MethodDescriptor<pactus.Wallet.SetAddressLabelRequest, pactus.Wallet.SetAddressLabelResponse> getSetAddressLabelMethod;
    if ((getSetAddressLabelMethod = WalletServiceGrpc.getSetAddressLabelMethod) == null) {
      synchronized (WalletServiceGrpc.class) {
        if ((getSetAddressLabelMethod = WalletServiceGrpc.getSetAddressLabelMethod) == null) {
          WalletServiceGrpc.getSetAddressLabelMethod = getSetAddressLabelMethod =
              io.grpc.MethodDescriptor.<pactus.Wallet.SetAddressLabelRequest, pactus.Wallet.SetAddressLabelResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "SetAddressLabel"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Wallet.SetAddressLabelRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Wallet.SetAddressLabelResponse.getDefaultInstance()))
              .setSchemaDescriptor(new WalletServiceMethodDescriptorSupplier("SetAddressLabel"))
              .build();
        }
      }
    }
    return getSetAddressLabelMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.Wallet.ListWalletRequest,
      pactus.Wallet.ListWalletResponse> getListWalletMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "ListWallet",
      requestType = pactus.Wallet.ListWalletRequest.class,
      responseType = pactus.Wallet.ListWalletResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.Wallet.ListWalletRequest,
      pactus.Wallet.ListWalletResponse> getListWalletMethod() {
    io.grpc.MethodDescriptor<pactus.Wallet.ListWalletRequest, pactus.Wallet.ListWalletResponse> getListWalletMethod;
    if ((getListWalletMethod = WalletServiceGrpc.getListWalletMethod) == null) {
      synchronized (WalletServiceGrpc.class) {
        if ((getListWalletMethod = WalletServiceGrpc.getListWalletMethod) == null) {
          WalletServiceGrpc.getListWalletMethod = getListWalletMethod =
              io.grpc.MethodDescriptor.<pactus.Wallet.ListWalletRequest, pactus.Wallet.ListWalletResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "ListWallet"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Wallet.ListWalletRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Wallet.ListWalletResponse.getDefaultInstance()))
              .setSchemaDescriptor(new WalletServiceMethodDescriptorSupplier("ListWallet"))
              .build();
        }
      }
    }
    return getListWalletMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.Wallet.GetWalletInfoRequest,
      pactus.Wallet.GetWalletInfoResponse> getGetWalletInfoMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetWalletInfo",
      requestType = pactus.Wallet.GetWalletInfoRequest.class,
      responseType = pactus.Wallet.GetWalletInfoResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.Wallet.GetWalletInfoRequest,
      pactus.Wallet.GetWalletInfoResponse> getGetWalletInfoMethod() {
    io.grpc.MethodDescriptor<pactus.Wallet.GetWalletInfoRequest, pactus.Wallet.GetWalletInfoResponse> getGetWalletInfoMethod;
    if ((getGetWalletInfoMethod = WalletServiceGrpc.getGetWalletInfoMethod) == null) {
      synchronized (WalletServiceGrpc.class) {
        if ((getGetWalletInfoMethod = WalletServiceGrpc.getGetWalletInfoMethod) == null) {
          WalletServiceGrpc.getGetWalletInfoMethod = getGetWalletInfoMethod =
              io.grpc.MethodDescriptor.<pactus.Wallet.GetWalletInfoRequest, pactus.Wallet.GetWalletInfoResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetWalletInfo"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Wallet.GetWalletInfoRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Wallet.GetWalletInfoResponse.getDefaultInstance()))
              .setSchemaDescriptor(new WalletServiceMethodDescriptorSupplier("GetWalletInfo"))
              .build();
        }
      }
    }
    return getGetWalletInfoMethod;
  }

  private static volatile io.grpc.MethodDescriptor<pactus.Wallet.ListAddressRequest,
      pactus.Wallet.ListAddressResponse> getListAddressMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "ListAddress",
      requestType = pactus.Wallet.ListAddressRequest.class,
      responseType = pactus.Wallet.ListAddressResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<pactus.Wallet.ListAddressRequest,
      pactus.Wallet.ListAddressResponse> getListAddressMethod() {
    io.grpc.MethodDescriptor<pactus.Wallet.ListAddressRequest, pactus.Wallet.ListAddressResponse> getListAddressMethod;
    if ((getListAddressMethod = WalletServiceGrpc.getListAddressMethod) == null) {
      synchronized (WalletServiceGrpc.class) {
        if ((getListAddressMethod = WalletServiceGrpc.getListAddressMethod) == null) {
          WalletServiceGrpc.getListAddressMethod = getListAddressMethod =
              io.grpc.MethodDescriptor.<pactus.Wallet.ListAddressRequest, pactus.Wallet.ListAddressResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "ListAddress"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Wallet.ListAddressRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  pactus.Wallet.ListAddressResponse.getDefaultInstance()))
              .setSchemaDescriptor(new WalletServiceMethodDescriptorSupplier("ListAddress"))
              .build();
        }
      }
    }
    return getListAddressMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static WalletServiceStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<WalletServiceStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<WalletServiceStub>() {
        @java.lang.Override
        public WalletServiceStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new WalletServiceStub(channel, callOptions);
        }
      };
    return WalletServiceStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports all types of calls on the service
   */
  public static WalletServiceBlockingV2Stub newBlockingV2Stub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<WalletServiceBlockingV2Stub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<WalletServiceBlockingV2Stub>() {
        @java.lang.Override
        public WalletServiceBlockingV2Stub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new WalletServiceBlockingV2Stub(channel, callOptions);
        }
      };
    return WalletServiceBlockingV2Stub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static WalletServiceBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<WalletServiceBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<WalletServiceBlockingStub>() {
        @java.lang.Override
        public WalletServiceBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new WalletServiceBlockingStub(channel, callOptions);
        }
      };
    return WalletServiceBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static WalletServiceFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<WalletServiceFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<WalletServiceFutureStub>() {
        @java.lang.Override
        public WalletServiceFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new WalletServiceFutureStub(channel, callOptions);
        }
      };
    return WalletServiceFutureStub.newStub(factory, channel);
  }

  /**
   * <pre>
   * WalletService provides RPC methods for wallet management operations.
   * </pre>
   */
  public interface AsyncService {

    /**
     * <pre>
     * CreateWallet creates a new wallet with the specified parameters.
     * </pre>
     */
    default void createWallet(pactus.Wallet.CreateWalletRequest request,
        io.grpc.stub.StreamObserver<pactus.Wallet.CreateWalletResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getCreateWalletMethod(), responseObserver);
    }

    /**
     * <pre>
     * RestoreWallet restores an existing wallet with the given mnemonic.
     * </pre>
     */
    default void restoreWallet(pactus.Wallet.RestoreWalletRequest request,
        io.grpc.stub.StreamObserver<pactus.Wallet.RestoreWalletResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getRestoreWalletMethod(), responseObserver);
    }

    /**
     * <pre>
     * LoadWallet loads an existing wallet with the given name.
     * </pre>
     */
    default void loadWallet(pactus.Wallet.LoadWalletRequest request,
        io.grpc.stub.StreamObserver<pactus.Wallet.LoadWalletResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getLoadWalletMethod(), responseObserver);
    }

    /**
     * <pre>
     * UnloadWallet unloads a currently loaded wallet with the specified name.
     * </pre>
     */
    default void unloadWallet(pactus.Wallet.UnloadWalletRequest request,
        io.grpc.stub.StreamObserver<pactus.Wallet.UnloadWalletResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getUnloadWalletMethod(), responseObserver);
    }

    /**
     * <pre>
     * GetTotalBalance returns the total available balance of the wallet.
     * </pre>
     */
    default void getTotalBalance(pactus.Wallet.GetTotalBalanceRequest request,
        io.grpc.stub.StreamObserver<pactus.Wallet.GetTotalBalanceResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetTotalBalanceMethod(), responseObserver);
    }

    /**
     * <pre>
     * SignRawTransaction signs a raw transaction for a specified wallet.
     * </pre>
     */
    default void signRawTransaction(pactus.Wallet.SignRawTransactionRequest request,
        io.grpc.stub.StreamObserver<pactus.Wallet.SignRawTransactionResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getSignRawTransactionMethod(), responseObserver);
    }

    /**
     * <pre>
     * GetValidatorAddress retrieves the validator address associated with a public key.
     * </pre>
     */
    default void getValidatorAddress(pactus.Wallet.GetValidatorAddressRequest request,
        io.grpc.stub.StreamObserver<pactus.Wallet.GetValidatorAddressResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetValidatorAddressMethod(), responseObserver);
    }

    /**
     * <pre>
     * GetNewAddress generates a new address for the specified wallet.
     * </pre>
     */
    default void getNewAddress(pactus.Wallet.GetNewAddressRequest request,
        io.grpc.stub.StreamObserver<pactus.Wallet.GetNewAddressResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetNewAddressMethod(), responseObserver);
    }

    /**
     * <pre>
     * GetAddressHistory retrieves the transaction history of an address.
     * </pre>
     */
    default void getAddressHistory(pactus.Wallet.GetAddressHistoryRequest request,
        io.grpc.stub.StreamObserver<pactus.Wallet.GetAddressHistoryResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetAddressHistoryMethod(), responseObserver);
    }

    /**
     * <pre>
     * SignMessage signs an arbitrary message using a wallet's private key.
     * </pre>
     */
    default void signMessage(pactus.Wallet.SignMessageRequest request,
        io.grpc.stub.StreamObserver<pactus.Wallet.SignMessageResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getSignMessageMethod(), responseObserver);
    }

    /**
     * <pre>
     * GetTotalStake returns the total stake amount in the wallet.
     * </pre>
     */
    default void getTotalStake(pactus.Wallet.GetTotalStakeRequest request,
        io.grpc.stub.StreamObserver<pactus.Wallet.GetTotalStakeResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetTotalStakeMethod(), responseObserver);
    }

    /**
     * <pre>
     * GetAddressInfo returns detailed information about a specific address.
     * </pre>
     */
    default void getAddressInfo(pactus.Wallet.GetAddressInfoRequest request,
        io.grpc.stub.StreamObserver<pactus.Wallet.GetAddressInfoResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetAddressInfoMethod(), responseObserver);
    }

    /**
     * <pre>
     * SetAddressLabel sets or updates the label for a given address.
     * </pre>
     */
    default void setAddressLabel(pactus.Wallet.SetAddressLabelRequest request,
        io.grpc.stub.StreamObserver<pactus.Wallet.SetAddressLabelResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getSetAddressLabelMethod(), responseObserver);
    }

    /**
     * <pre>
     * ListWallet returns list of all available wallets.
     * </pre>
     */
    default void listWallet(pactus.Wallet.ListWalletRequest request,
        io.grpc.stub.StreamObserver<pactus.Wallet.ListWalletResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getListWalletMethod(), responseObserver);
    }

    /**
     * <pre>
     * GetWalletInfo returns detailed information about a specific wallet.
     * </pre>
     */
    default void getWalletInfo(pactus.Wallet.GetWalletInfoRequest request,
        io.grpc.stub.StreamObserver<pactus.Wallet.GetWalletInfoResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetWalletInfoMethod(), responseObserver);
    }

    /**
     * <pre>
     * ListAddress returns all addresses in the specified wallet.
     * </pre>
     */
    default void listAddress(pactus.Wallet.ListAddressRequest request,
        io.grpc.stub.StreamObserver<pactus.Wallet.ListAddressResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getListAddressMethod(), responseObserver);
    }
  }

  /**
   * Base class for the server implementation of the service WalletService.
   * <pre>
   * WalletService provides RPC methods for wallet management operations.
   * </pre>
   */
  public static abstract class WalletServiceImplBase
      implements io.grpc.BindableService, AsyncService {

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return WalletServiceGrpc.bindService(this);
    }
  }

  /**
   * A stub to allow clients to do asynchronous rpc calls to service WalletService.
   * <pre>
   * WalletService provides RPC methods for wallet management operations.
   * </pre>
   */
  public static final class WalletServiceStub
      extends io.grpc.stub.AbstractAsyncStub<WalletServiceStub> {
    private WalletServiceStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected WalletServiceStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new WalletServiceStub(channel, callOptions);
    }

    /**
     * <pre>
     * CreateWallet creates a new wallet with the specified parameters.
     * </pre>
     */
    public void createWallet(pactus.Wallet.CreateWalletRequest request,
        io.grpc.stub.StreamObserver<pactus.Wallet.CreateWalletResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getCreateWalletMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * RestoreWallet restores an existing wallet with the given mnemonic.
     * </pre>
     */
    public void restoreWallet(pactus.Wallet.RestoreWalletRequest request,
        io.grpc.stub.StreamObserver<pactus.Wallet.RestoreWalletResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getRestoreWalletMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * LoadWallet loads an existing wallet with the given name.
     * </pre>
     */
    public void loadWallet(pactus.Wallet.LoadWalletRequest request,
        io.grpc.stub.StreamObserver<pactus.Wallet.LoadWalletResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getLoadWalletMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * UnloadWallet unloads a currently loaded wallet with the specified name.
     * </pre>
     */
    public void unloadWallet(pactus.Wallet.UnloadWalletRequest request,
        io.grpc.stub.StreamObserver<pactus.Wallet.UnloadWalletResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getUnloadWalletMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * GetTotalBalance returns the total available balance of the wallet.
     * </pre>
     */
    public void getTotalBalance(pactus.Wallet.GetTotalBalanceRequest request,
        io.grpc.stub.StreamObserver<pactus.Wallet.GetTotalBalanceResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetTotalBalanceMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * SignRawTransaction signs a raw transaction for a specified wallet.
     * </pre>
     */
    public void signRawTransaction(pactus.Wallet.SignRawTransactionRequest request,
        io.grpc.stub.StreamObserver<pactus.Wallet.SignRawTransactionResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getSignRawTransactionMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * GetValidatorAddress retrieves the validator address associated with a public key.
     * </pre>
     */
    public void getValidatorAddress(pactus.Wallet.GetValidatorAddressRequest request,
        io.grpc.stub.StreamObserver<pactus.Wallet.GetValidatorAddressResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetValidatorAddressMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * GetNewAddress generates a new address for the specified wallet.
     * </pre>
     */
    public void getNewAddress(pactus.Wallet.GetNewAddressRequest request,
        io.grpc.stub.StreamObserver<pactus.Wallet.GetNewAddressResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetNewAddressMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * GetAddressHistory retrieves the transaction history of an address.
     * </pre>
     */
    public void getAddressHistory(pactus.Wallet.GetAddressHistoryRequest request,
        io.grpc.stub.StreamObserver<pactus.Wallet.GetAddressHistoryResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetAddressHistoryMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * SignMessage signs an arbitrary message using a wallet's private key.
     * </pre>
     */
    public void signMessage(pactus.Wallet.SignMessageRequest request,
        io.grpc.stub.StreamObserver<pactus.Wallet.SignMessageResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getSignMessageMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * GetTotalStake returns the total stake amount in the wallet.
     * </pre>
     */
    public void getTotalStake(pactus.Wallet.GetTotalStakeRequest request,
        io.grpc.stub.StreamObserver<pactus.Wallet.GetTotalStakeResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetTotalStakeMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * GetAddressInfo returns detailed information about a specific address.
     * </pre>
     */
    public void getAddressInfo(pactus.Wallet.GetAddressInfoRequest request,
        io.grpc.stub.StreamObserver<pactus.Wallet.GetAddressInfoResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetAddressInfoMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * SetAddressLabel sets or updates the label for a given address.
     * </pre>
     */
    public void setAddressLabel(pactus.Wallet.SetAddressLabelRequest request,
        io.grpc.stub.StreamObserver<pactus.Wallet.SetAddressLabelResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getSetAddressLabelMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * ListWallet returns list of all available wallets.
     * </pre>
     */
    public void listWallet(pactus.Wallet.ListWalletRequest request,
        io.grpc.stub.StreamObserver<pactus.Wallet.ListWalletResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getListWalletMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * GetWalletInfo returns detailed information about a specific wallet.
     * </pre>
     */
    public void getWalletInfo(pactus.Wallet.GetWalletInfoRequest request,
        io.grpc.stub.StreamObserver<pactus.Wallet.GetWalletInfoResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetWalletInfoMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * ListAddress returns all addresses in the specified wallet.
     * </pre>
     */
    public void listAddress(pactus.Wallet.ListAddressRequest request,
        io.grpc.stub.StreamObserver<pactus.Wallet.ListAddressResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getListAddressMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   * A stub to allow clients to do synchronous rpc calls to service WalletService.
   * <pre>
   * WalletService provides RPC methods for wallet management operations.
   * </pre>
   */
  public static final class WalletServiceBlockingV2Stub
      extends io.grpc.stub.AbstractBlockingStub<WalletServiceBlockingV2Stub> {
    private WalletServiceBlockingV2Stub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected WalletServiceBlockingV2Stub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new WalletServiceBlockingV2Stub(channel, callOptions);
    }

    /**
     * <pre>
     * CreateWallet creates a new wallet with the specified parameters.
     * </pre>
     */
    public pactus.Wallet.CreateWalletResponse createWallet(pactus.Wallet.CreateWalletRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getCreateWalletMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * RestoreWallet restores an existing wallet with the given mnemonic.
     * </pre>
     */
    public pactus.Wallet.RestoreWalletResponse restoreWallet(pactus.Wallet.RestoreWalletRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getRestoreWalletMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * LoadWallet loads an existing wallet with the given name.
     * </pre>
     */
    public pactus.Wallet.LoadWalletResponse loadWallet(pactus.Wallet.LoadWalletRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getLoadWalletMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * UnloadWallet unloads a currently loaded wallet with the specified name.
     * </pre>
     */
    public pactus.Wallet.UnloadWalletResponse unloadWallet(pactus.Wallet.UnloadWalletRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getUnloadWalletMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetTotalBalance returns the total available balance of the wallet.
     * </pre>
     */
    public pactus.Wallet.GetTotalBalanceResponse getTotalBalance(pactus.Wallet.GetTotalBalanceRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetTotalBalanceMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * SignRawTransaction signs a raw transaction for a specified wallet.
     * </pre>
     */
    public pactus.Wallet.SignRawTransactionResponse signRawTransaction(pactus.Wallet.SignRawTransactionRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getSignRawTransactionMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetValidatorAddress retrieves the validator address associated with a public key.
     * </pre>
     */
    public pactus.Wallet.GetValidatorAddressResponse getValidatorAddress(pactus.Wallet.GetValidatorAddressRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetValidatorAddressMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetNewAddress generates a new address for the specified wallet.
     * </pre>
     */
    public pactus.Wallet.GetNewAddressResponse getNewAddress(pactus.Wallet.GetNewAddressRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetNewAddressMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetAddressHistory retrieves the transaction history of an address.
     * </pre>
     */
    public pactus.Wallet.GetAddressHistoryResponse getAddressHistory(pactus.Wallet.GetAddressHistoryRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetAddressHistoryMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * SignMessage signs an arbitrary message using a wallet's private key.
     * </pre>
     */
    public pactus.Wallet.SignMessageResponse signMessage(pactus.Wallet.SignMessageRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getSignMessageMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetTotalStake returns the total stake amount in the wallet.
     * </pre>
     */
    public pactus.Wallet.GetTotalStakeResponse getTotalStake(pactus.Wallet.GetTotalStakeRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetTotalStakeMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetAddressInfo returns detailed information about a specific address.
     * </pre>
     */
    public pactus.Wallet.GetAddressInfoResponse getAddressInfo(pactus.Wallet.GetAddressInfoRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetAddressInfoMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * SetAddressLabel sets or updates the label for a given address.
     * </pre>
     */
    public pactus.Wallet.SetAddressLabelResponse setAddressLabel(pactus.Wallet.SetAddressLabelRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getSetAddressLabelMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * ListWallet returns list of all available wallets.
     * </pre>
     */
    public pactus.Wallet.ListWalletResponse listWallet(pactus.Wallet.ListWalletRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getListWalletMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetWalletInfo returns detailed information about a specific wallet.
     * </pre>
     */
    public pactus.Wallet.GetWalletInfoResponse getWalletInfo(pactus.Wallet.GetWalletInfoRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetWalletInfoMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * ListAddress returns all addresses in the specified wallet.
     * </pre>
     */
    public pactus.Wallet.ListAddressResponse listAddress(pactus.Wallet.ListAddressRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getListAddressMethod(), getCallOptions(), request);
    }
  }

  /**
   * A stub to allow clients to do limited synchronous rpc calls to service WalletService.
   * <pre>
   * WalletService provides RPC methods for wallet management operations.
   * </pre>
   */
  public static final class WalletServiceBlockingStub
      extends io.grpc.stub.AbstractBlockingStub<WalletServiceBlockingStub> {
    private WalletServiceBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected WalletServiceBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new WalletServiceBlockingStub(channel, callOptions);
    }

    /**
     * <pre>
     * CreateWallet creates a new wallet with the specified parameters.
     * </pre>
     */
    public pactus.Wallet.CreateWalletResponse createWallet(pactus.Wallet.CreateWalletRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getCreateWalletMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * RestoreWallet restores an existing wallet with the given mnemonic.
     * </pre>
     */
    public pactus.Wallet.RestoreWalletResponse restoreWallet(pactus.Wallet.RestoreWalletRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getRestoreWalletMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * LoadWallet loads an existing wallet with the given name.
     * </pre>
     */
    public pactus.Wallet.LoadWalletResponse loadWallet(pactus.Wallet.LoadWalletRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getLoadWalletMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * UnloadWallet unloads a currently loaded wallet with the specified name.
     * </pre>
     */
    public pactus.Wallet.UnloadWalletResponse unloadWallet(pactus.Wallet.UnloadWalletRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getUnloadWalletMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetTotalBalance returns the total available balance of the wallet.
     * </pre>
     */
    public pactus.Wallet.GetTotalBalanceResponse getTotalBalance(pactus.Wallet.GetTotalBalanceRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetTotalBalanceMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * SignRawTransaction signs a raw transaction for a specified wallet.
     * </pre>
     */
    public pactus.Wallet.SignRawTransactionResponse signRawTransaction(pactus.Wallet.SignRawTransactionRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getSignRawTransactionMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetValidatorAddress retrieves the validator address associated with a public key.
     * </pre>
     */
    public pactus.Wallet.GetValidatorAddressResponse getValidatorAddress(pactus.Wallet.GetValidatorAddressRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetValidatorAddressMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetNewAddress generates a new address for the specified wallet.
     * </pre>
     */
    public pactus.Wallet.GetNewAddressResponse getNewAddress(pactus.Wallet.GetNewAddressRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetNewAddressMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetAddressHistory retrieves the transaction history of an address.
     * </pre>
     */
    public pactus.Wallet.GetAddressHistoryResponse getAddressHistory(pactus.Wallet.GetAddressHistoryRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetAddressHistoryMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * SignMessage signs an arbitrary message using a wallet's private key.
     * </pre>
     */
    public pactus.Wallet.SignMessageResponse signMessage(pactus.Wallet.SignMessageRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getSignMessageMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetTotalStake returns the total stake amount in the wallet.
     * </pre>
     */
    public pactus.Wallet.GetTotalStakeResponse getTotalStake(pactus.Wallet.GetTotalStakeRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetTotalStakeMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetAddressInfo returns detailed information about a specific address.
     * </pre>
     */
    public pactus.Wallet.GetAddressInfoResponse getAddressInfo(pactus.Wallet.GetAddressInfoRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetAddressInfoMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * SetAddressLabel sets or updates the label for a given address.
     * </pre>
     */
    public pactus.Wallet.SetAddressLabelResponse setAddressLabel(pactus.Wallet.SetAddressLabelRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getSetAddressLabelMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * ListWallet returns list of all available wallets.
     * </pre>
     */
    public pactus.Wallet.ListWalletResponse listWallet(pactus.Wallet.ListWalletRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getListWalletMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetWalletInfo returns detailed information about a specific wallet.
     * </pre>
     */
    public pactus.Wallet.GetWalletInfoResponse getWalletInfo(pactus.Wallet.GetWalletInfoRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetWalletInfoMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * ListAddress returns all addresses in the specified wallet.
     * </pre>
     */
    public pactus.Wallet.ListAddressResponse listAddress(pactus.Wallet.ListAddressRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getListAddressMethod(), getCallOptions(), request);
    }
  }

  /**
   * A stub to allow clients to do ListenableFuture-style rpc calls to service WalletService.
   * <pre>
   * WalletService provides RPC methods for wallet management operations.
   * </pre>
   */
  public static final class WalletServiceFutureStub
      extends io.grpc.stub.AbstractFutureStub<WalletServiceFutureStub> {
    private WalletServiceFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected WalletServiceFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new WalletServiceFutureStub(channel, callOptions);
    }

    /**
     * <pre>
     * CreateWallet creates a new wallet with the specified parameters.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.Wallet.CreateWalletResponse> createWallet(
        pactus.Wallet.CreateWalletRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getCreateWalletMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * RestoreWallet restores an existing wallet with the given mnemonic.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.Wallet.RestoreWalletResponse> restoreWallet(
        pactus.Wallet.RestoreWalletRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getRestoreWalletMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * LoadWallet loads an existing wallet with the given name.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.Wallet.LoadWalletResponse> loadWallet(
        pactus.Wallet.LoadWalletRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getLoadWalletMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * UnloadWallet unloads a currently loaded wallet with the specified name.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.Wallet.UnloadWalletResponse> unloadWallet(
        pactus.Wallet.UnloadWalletRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getUnloadWalletMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * GetTotalBalance returns the total available balance of the wallet.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.Wallet.GetTotalBalanceResponse> getTotalBalance(
        pactus.Wallet.GetTotalBalanceRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetTotalBalanceMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * SignRawTransaction signs a raw transaction for a specified wallet.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.Wallet.SignRawTransactionResponse> signRawTransaction(
        pactus.Wallet.SignRawTransactionRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getSignRawTransactionMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * GetValidatorAddress retrieves the validator address associated with a public key.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.Wallet.GetValidatorAddressResponse> getValidatorAddress(
        pactus.Wallet.GetValidatorAddressRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetValidatorAddressMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * GetNewAddress generates a new address for the specified wallet.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.Wallet.GetNewAddressResponse> getNewAddress(
        pactus.Wallet.GetNewAddressRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetNewAddressMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * GetAddressHistory retrieves the transaction history of an address.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.Wallet.GetAddressHistoryResponse> getAddressHistory(
        pactus.Wallet.GetAddressHistoryRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetAddressHistoryMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * SignMessage signs an arbitrary message using a wallet's private key.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.Wallet.SignMessageResponse> signMessage(
        pactus.Wallet.SignMessageRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getSignMessageMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * GetTotalStake returns the total stake amount in the wallet.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.Wallet.GetTotalStakeResponse> getTotalStake(
        pactus.Wallet.GetTotalStakeRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetTotalStakeMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * GetAddressInfo returns detailed information about a specific address.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.Wallet.GetAddressInfoResponse> getAddressInfo(
        pactus.Wallet.GetAddressInfoRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetAddressInfoMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * SetAddressLabel sets or updates the label for a given address.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.Wallet.SetAddressLabelResponse> setAddressLabel(
        pactus.Wallet.SetAddressLabelRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getSetAddressLabelMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * ListWallet returns list of all available wallets.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.Wallet.ListWalletResponse> listWallet(
        pactus.Wallet.ListWalletRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getListWalletMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * GetWalletInfo returns detailed information about a specific wallet.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.Wallet.GetWalletInfoResponse> getWalletInfo(
        pactus.Wallet.GetWalletInfoRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetWalletInfoMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * ListAddress returns all addresses in the specified wallet.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<pactus.Wallet.ListAddressResponse> listAddress(
        pactus.Wallet.ListAddressRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getListAddressMethod(), getCallOptions()), request);
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
  private static final int METHODID_LIST_WALLET = 13;
  private static final int METHODID_GET_WALLET_INFO = 14;
  private static final int METHODID_LIST_ADDRESS = 15;

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
          serviceImpl.createWallet((pactus.Wallet.CreateWalletRequest) request,
              (io.grpc.stub.StreamObserver<pactus.Wallet.CreateWalletResponse>) responseObserver);
          break;
        case METHODID_RESTORE_WALLET:
          serviceImpl.restoreWallet((pactus.Wallet.RestoreWalletRequest) request,
              (io.grpc.stub.StreamObserver<pactus.Wallet.RestoreWalletResponse>) responseObserver);
          break;
        case METHODID_LOAD_WALLET:
          serviceImpl.loadWallet((pactus.Wallet.LoadWalletRequest) request,
              (io.grpc.stub.StreamObserver<pactus.Wallet.LoadWalletResponse>) responseObserver);
          break;
        case METHODID_UNLOAD_WALLET:
          serviceImpl.unloadWallet((pactus.Wallet.UnloadWalletRequest) request,
              (io.grpc.stub.StreamObserver<pactus.Wallet.UnloadWalletResponse>) responseObserver);
          break;
        case METHODID_GET_TOTAL_BALANCE:
          serviceImpl.getTotalBalance((pactus.Wallet.GetTotalBalanceRequest) request,
              (io.grpc.stub.StreamObserver<pactus.Wallet.GetTotalBalanceResponse>) responseObserver);
          break;
        case METHODID_SIGN_RAW_TRANSACTION:
          serviceImpl.signRawTransaction((pactus.Wallet.SignRawTransactionRequest) request,
              (io.grpc.stub.StreamObserver<pactus.Wallet.SignRawTransactionResponse>) responseObserver);
          break;
        case METHODID_GET_VALIDATOR_ADDRESS:
          serviceImpl.getValidatorAddress((pactus.Wallet.GetValidatorAddressRequest) request,
              (io.grpc.stub.StreamObserver<pactus.Wallet.GetValidatorAddressResponse>) responseObserver);
          break;
        case METHODID_GET_NEW_ADDRESS:
          serviceImpl.getNewAddress((pactus.Wallet.GetNewAddressRequest) request,
              (io.grpc.stub.StreamObserver<pactus.Wallet.GetNewAddressResponse>) responseObserver);
          break;
        case METHODID_GET_ADDRESS_HISTORY:
          serviceImpl.getAddressHistory((pactus.Wallet.GetAddressHistoryRequest) request,
              (io.grpc.stub.StreamObserver<pactus.Wallet.GetAddressHistoryResponse>) responseObserver);
          break;
        case METHODID_SIGN_MESSAGE:
          serviceImpl.signMessage((pactus.Wallet.SignMessageRequest) request,
              (io.grpc.stub.StreamObserver<pactus.Wallet.SignMessageResponse>) responseObserver);
          break;
        case METHODID_GET_TOTAL_STAKE:
          serviceImpl.getTotalStake((pactus.Wallet.GetTotalStakeRequest) request,
              (io.grpc.stub.StreamObserver<pactus.Wallet.GetTotalStakeResponse>) responseObserver);
          break;
        case METHODID_GET_ADDRESS_INFO:
          serviceImpl.getAddressInfo((pactus.Wallet.GetAddressInfoRequest) request,
              (io.grpc.stub.StreamObserver<pactus.Wallet.GetAddressInfoResponse>) responseObserver);
          break;
        case METHODID_SET_ADDRESS_LABEL:
          serviceImpl.setAddressLabel((pactus.Wallet.SetAddressLabelRequest) request,
              (io.grpc.stub.StreamObserver<pactus.Wallet.SetAddressLabelResponse>) responseObserver);
          break;
        case METHODID_LIST_WALLET:
          serviceImpl.listWallet((pactus.Wallet.ListWalletRequest) request,
              (io.grpc.stub.StreamObserver<pactus.Wallet.ListWalletResponse>) responseObserver);
          break;
        case METHODID_GET_WALLET_INFO:
          serviceImpl.getWalletInfo((pactus.Wallet.GetWalletInfoRequest) request,
              (io.grpc.stub.StreamObserver<pactus.Wallet.GetWalletInfoResponse>) responseObserver);
          break;
        case METHODID_LIST_ADDRESS:
          serviceImpl.listAddress((pactus.Wallet.ListAddressRequest) request,
              (io.grpc.stub.StreamObserver<pactus.Wallet.ListAddressResponse>) responseObserver);
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
              pactus.Wallet.CreateWalletRequest,
              pactus.Wallet.CreateWalletResponse>(
                service, METHODID_CREATE_WALLET)))
        .addMethod(
          getRestoreWalletMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.Wallet.RestoreWalletRequest,
              pactus.Wallet.RestoreWalletResponse>(
                service, METHODID_RESTORE_WALLET)))
        .addMethod(
          getLoadWalletMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.Wallet.LoadWalletRequest,
              pactus.Wallet.LoadWalletResponse>(
                service, METHODID_LOAD_WALLET)))
        .addMethod(
          getUnloadWalletMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.Wallet.UnloadWalletRequest,
              pactus.Wallet.UnloadWalletResponse>(
                service, METHODID_UNLOAD_WALLET)))
        .addMethod(
          getGetTotalBalanceMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.Wallet.GetTotalBalanceRequest,
              pactus.Wallet.GetTotalBalanceResponse>(
                service, METHODID_GET_TOTAL_BALANCE)))
        .addMethod(
          getSignRawTransactionMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.Wallet.SignRawTransactionRequest,
              pactus.Wallet.SignRawTransactionResponse>(
                service, METHODID_SIGN_RAW_TRANSACTION)))
        .addMethod(
          getGetValidatorAddressMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.Wallet.GetValidatorAddressRequest,
              pactus.Wallet.GetValidatorAddressResponse>(
                service, METHODID_GET_VALIDATOR_ADDRESS)))
        .addMethod(
          getGetNewAddressMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.Wallet.GetNewAddressRequest,
              pactus.Wallet.GetNewAddressResponse>(
                service, METHODID_GET_NEW_ADDRESS)))
        .addMethod(
          getGetAddressHistoryMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.Wallet.GetAddressHistoryRequest,
              pactus.Wallet.GetAddressHistoryResponse>(
                service, METHODID_GET_ADDRESS_HISTORY)))
        .addMethod(
          getSignMessageMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.Wallet.SignMessageRequest,
              pactus.Wallet.SignMessageResponse>(
                service, METHODID_SIGN_MESSAGE)))
        .addMethod(
          getGetTotalStakeMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.Wallet.GetTotalStakeRequest,
              pactus.Wallet.GetTotalStakeResponse>(
                service, METHODID_GET_TOTAL_STAKE)))
        .addMethod(
          getGetAddressInfoMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.Wallet.GetAddressInfoRequest,
              pactus.Wallet.GetAddressInfoResponse>(
                service, METHODID_GET_ADDRESS_INFO)))
        .addMethod(
          getSetAddressLabelMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.Wallet.SetAddressLabelRequest,
              pactus.Wallet.SetAddressLabelResponse>(
                service, METHODID_SET_ADDRESS_LABEL)))
        .addMethod(
          getListWalletMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.Wallet.ListWalletRequest,
              pactus.Wallet.ListWalletResponse>(
                service, METHODID_LIST_WALLET)))
        .addMethod(
          getGetWalletInfoMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.Wallet.GetWalletInfoRequest,
              pactus.Wallet.GetWalletInfoResponse>(
                service, METHODID_GET_WALLET_INFO)))
        .addMethod(
          getListAddressMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              pactus.Wallet.ListAddressRequest,
              pactus.Wallet.ListAddressResponse>(
                service, METHODID_LIST_ADDRESS)))
        .build();
  }

  private static abstract class WalletServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    WalletServiceBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return pactus.Wallet.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("WalletService");
    }
  }

  private static final class WalletServiceFileDescriptorSupplier
      extends WalletServiceBaseDescriptorSupplier {
    WalletServiceFileDescriptorSupplier() {}
  }

  private static final class WalletServiceMethodDescriptorSupplier
      extends WalletServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final java.lang.String methodName;

    WalletServiceMethodDescriptorSupplier(java.lang.String methodName) {
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
      synchronized (WalletServiceGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new WalletServiceFileDescriptorSupplier())
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
              .addMethod(getListWalletMethod())
              .addMethod(getGetWalletInfoMethod())
              .addMethod(getListAddressMethod())
              .build();
        }
      }
    }
    return result;
  }
}
