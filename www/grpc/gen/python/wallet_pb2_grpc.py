# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc

import wallet_pb2 as wallet__pb2


class WalletStub(object):
    """Define the Wallet service with various RPC methods for wallet management.
    """

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.CreateWallet = channel.unary_unary(
                '/pactus.Wallet/CreateWallet',
                request_serializer=wallet__pb2.CreateWalletRequest.SerializeToString,
                response_deserializer=wallet__pb2.CreateWalletResponse.FromString,
                )
        self.LoadWallet = channel.unary_unary(
                '/pactus.Wallet/LoadWallet',
                request_serializer=wallet__pb2.LoadWalletRequest.SerializeToString,
                response_deserializer=wallet__pb2.LoadWalletResponse.FromString,
                )
        self.UnloadWallet = channel.unary_unary(
                '/pactus.Wallet/UnloadWallet',
                request_serializer=wallet__pb2.UnloadWalletRequest.SerializeToString,
                response_deserializer=wallet__pb2.UnloadWalletResponse.FromString,
                )
        self.LockWallet = channel.unary_unary(
                '/pactus.Wallet/LockWallet',
                request_serializer=wallet__pb2.LockWalletRequest.SerializeToString,
                response_deserializer=wallet__pb2.LockWalletResponse.FromString,
                )
        self.UnlockWallet = channel.unary_unary(
                '/pactus.Wallet/UnlockWallet',
                request_serializer=wallet__pb2.UnlockWalletRequest.SerializeToString,
                response_deserializer=wallet__pb2.UnlockWalletResponse.FromString,
                )
        self.SignRawTransaction = channel.unary_unary(
                '/pactus.Wallet/SignRawTransaction',
                request_serializer=wallet__pb2.SignRawTransactionRequest.SerializeToString,
                response_deserializer=wallet__pb2.SignRawTransactionResponse.FromString,
                )
        self.GetValidatorAddress = channel.unary_unary(
                '/pactus.Wallet/GetValidatorAddress',
                request_serializer=wallet__pb2.GetValidatorAddressRequest.SerializeToString,
                response_deserializer=wallet__pb2.GetValidatorAddressResponse.FromString,
                )


class WalletServicer(object):
    """Define the Wallet service with various RPC methods for wallet management.
    """

    def CreateWallet(self, request, context):
        """CreateWallet creates a new wallet with the specified parameters.
        """
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def LoadWallet(self, request, context):
        """LoadWallet loads an existing wallet with the given name.
        """
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def UnloadWallet(self, request, context):
        """UnloadWallet unloads a currently loaded wallet with the specified name.
        """
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def LockWallet(self, request, context):
        """LockWallet locks a currently loaded wallet with the provided password and
        timeout.
        """
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def UnlockWallet(self, request, context):
        """UnlockWallet unlocks a locked wallet with the provided password and
        timeout.
        """
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def SignRawTransaction(self, request, context):
        """SignRawTransaction signs a raw transaction for a specified wallet.
        """
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def GetValidatorAddress(self, request, context):
        """GetValidatorAddress retrieves the validator address associated with a
        public key.
        """
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_WalletServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'CreateWallet': grpc.unary_unary_rpc_method_handler(
                    servicer.CreateWallet,
                    request_deserializer=wallet__pb2.CreateWalletRequest.FromString,
                    response_serializer=wallet__pb2.CreateWalletResponse.SerializeToString,
            ),
            'LoadWallet': grpc.unary_unary_rpc_method_handler(
                    servicer.LoadWallet,
                    request_deserializer=wallet__pb2.LoadWalletRequest.FromString,
                    response_serializer=wallet__pb2.LoadWalletResponse.SerializeToString,
            ),
            'UnloadWallet': grpc.unary_unary_rpc_method_handler(
                    servicer.UnloadWallet,
                    request_deserializer=wallet__pb2.UnloadWalletRequest.FromString,
                    response_serializer=wallet__pb2.UnloadWalletResponse.SerializeToString,
            ),
            'LockWallet': grpc.unary_unary_rpc_method_handler(
                    servicer.LockWallet,
                    request_deserializer=wallet__pb2.LockWalletRequest.FromString,
                    response_serializer=wallet__pb2.LockWalletResponse.SerializeToString,
            ),
            'UnlockWallet': grpc.unary_unary_rpc_method_handler(
                    servicer.UnlockWallet,
                    request_deserializer=wallet__pb2.UnlockWalletRequest.FromString,
                    response_serializer=wallet__pb2.UnlockWalletResponse.SerializeToString,
            ),
            'SignRawTransaction': grpc.unary_unary_rpc_method_handler(
                    servicer.SignRawTransaction,
                    request_deserializer=wallet__pb2.SignRawTransactionRequest.FromString,
                    response_serializer=wallet__pb2.SignRawTransactionResponse.SerializeToString,
            ),
            'GetValidatorAddress': grpc.unary_unary_rpc_method_handler(
                    servicer.GetValidatorAddress,
                    request_deserializer=wallet__pb2.GetValidatorAddressRequest.FromString,
                    response_serializer=wallet__pb2.GetValidatorAddressResponse.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'pactus.Wallet', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))


 # This class is part of an EXPERIMENTAL API.
class Wallet(object):
    """Define the Wallet service with various RPC methods for wallet management.
    """

    @staticmethod
    def CreateWallet(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/pactus.Wallet/CreateWallet',
            wallet__pb2.CreateWalletRequest.SerializeToString,
            wallet__pb2.CreateWalletResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def LoadWallet(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/pactus.Wallet/LoadWallet',
            wallet__pb2.LoadWalletRequest.SerializeToString,
            wallet__pb2.LoadWalletResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def UnloadWallet(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/pactus.Wallet/UnloadWallet',
            wallet__pb2.UnloadWalletRequest.SerializeToString,
            wallet__pb2.UnloadWalletResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def LockWallet(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/pactus.Wallet/LockWallet',
            wallet__pb2.LockWalletRequest.SerializeToString,
            wallet__pb2.LockWalletResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def UnlockWallet(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/pactus.Wallet/UnlockWallet',
            wallet__pb2.UnlockWalletRequest.SerializeToString,
            wallet__pb2.UnlockWalletResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def SignRawTransaction(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/pactus.Wallet/SignRawTransaction',
            wallet__pb2.SignRawTransactionRequest.SerializeToString,
            wallet__pb2.SignRawTransactionResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def GetValidatorAddress(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/pactus.Wallet/GetValidatorAddress',
            wallet__pb2.GetValidatorAddressRequest.SerializeToString,
            wallet__pb2.GetValidatorAddressResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)
