import transaction_pb2 as _transaction_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class AddressType(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    ADDRESS_TYPE_TREASURY: _ClassVar[AddressType]
    ADDRESS_TYPE_VALIDATOR: _ClassVar[AddressType]
    ADDRESS_TYPE_BLS_ACCOUNT: _ClassVar[AddressType]
    ADDRESS_TYPE_ED25519_ACCOUNT: _ClassVar[AddressType]
ADDRESS_TYPE_TREASURY: AddressType
ADDRESS_TYPE_VALIDATOR: AddressType
ADDRESS_TYPE_BLS_ACCOUNT: AddressType
ADDRESS_TYPE_ED25519_ACCOUNT: AddressType

class AddressInfo(_message.Message):
    __slots__ = ("address", "public_key", "label", "path")
    ADDRESS_FIELD_NUMBER: _ClassVar[int]
    PUBLIC_KEY_FIELD_NUMBER: _ClassVar[int]
    LABEL_FIELD_NUMBER: _ClassVar[int]
    PATH_FIELD_NUMBER: _ClassVar[int]
    address: str
    public_key: str
    label: str
    path: str
    def __init__(self, address: _Optional[str] = ..., public_key: _Optional[str] = ..., label: _Optional[str] = ..., path: _Optional[str] = ...) -> None: ...

class HistoryInfo(_message.Message):
    __slots__ = ("transaction_id", "time", "payload_type", "description", "amount")
    TRANSACTION_ID_FIELD_NUMBER: _ClassVar[int]
    TIME_FIELD_NUMBER: _ClassVar[int]
    PAYLOAD_TYPE_FIELD_NUMBER: _ClassVar[int]
    DESCRIPTION_FIELD_NUMBER: _ClassVar[int]
    AMOUNT_FIELD_NUMBER: _ClassVar[int]
    transaction_id: str
    time: int
    payload_type: str
    description: str
    amount: int
    def __init__(self, transaction_id: _Optional[str] = ..., time: _Optional[int] = ..., payload_type: _Optional[str] = ..., description: _Optional[str] = ..., amount: _Optional[int] = ...) -> None: ...

class GetAddressHistoryRequest(_message.Message):
    __slots__ = ("wallet_name", "address")
    WALLET_NAME_FIELD_NUMBER: _ClassVar[int]
    ADDRESS_FIELD_NUMBER: _ClassVar[int]
    wallet_name: str
    address: str
    def __init__(self, wallet_name: _Optional[str] = ..., address: _Optional[str] = ...) -> None: ...

class GetAddressHistoryResponse(_message.Message):
    __slots__ = ("history_info",)
    HISTORY_INFO_FIELD_NUMBER: _ClassVar[int]
    history_info: _containers.RepeatedCompositeFieldContainer[HistoryInfo]
    def __init__(self, history_info: _Optional[_Iterable[_Union[HistoryInfo, _Mapping]]] = ...) -> None: ...

class GetNewAddressRequest(_message.Message):
    __slots__ = ("wallet_name", "address_type", "label", "password")
    WALLET_NAME_FIELD_NUMBER: _ClassVar[int]
    ADDRESS_TYPE_FIELD_NUMBER: _ClassVar[int]
    LABEL_FIELD_NUMBER: _ClassVar[int]
    PASSWORD_FIELD_NUMBER: _ClassVar[int]
    wallet_name: str
    address_type: AddressType
    label: str
    password: str
    def __init__(self, wallet_name: _Optional[str] = ..., address_type: _Optional[_Union[AddressType, str]] = ..., label: _Optional[str] = ..., password: _Optional[str] = ...) -> None: ...

class GetNewAddressResponse(_message.Message):
    __slots__ = ("wallet_name", "address_info")
    WALLET_NAME_FIELD_NUMBER: _ClassVar[int]
    ADDRESS_INFO_FIELD_NUMBER: _ClassVar[int]
    wallet_name: str
    address_info: AddressInfo
    def __init__(self, wallet_name: _Optional[str] = ..., address_info: _Optional[_Union[AddressInfo, _Mapping]] = ...) -> None: ...

class RestoreWalletRequest(_message.Message):
    __slots__ = ("wallet_name", "mnemonic", "password")
    WALLET_NAME_FIELD_NUMBER: _ClassVar[int]
    MNEMONIC_FIELD_NUMBER: _ClassVar[int]
    PASSWORD_FIELD_NUMBER: _ClassVar[int]
    wallet_name: str
    mnemonic: str
    password: str
    def __init__(self, wallet_name: _Optional[str] = ..., mnemonic: _Optional[str] = ..., password: _Optional[str] = ...) -> None: ...

class RestoreWalletResponse(_message.Message):
    __slots__ = ("wallet_name",)
    WALLET_NAME_FIELD_NUMBER: _ClassVar[int]
    wallet_name: str
    def __init__(self, wallet_name: _Optional[str] = ...) -> None: ...

class CreateWalletRequest(_message.Message):
    __slots__ = ("wallet_name", "password")
    WALLET_NAME_FIELD_NUMBER: _ClassVar[int]
    PASSWORD_FIELD_NUMBER: _ClassVar[int]
    wallet_name: str
    password: str
    def __init__(self, wallet_name: _Optional[str] = ..., password: _Optional[str] = ...) -> None: ...

class CreateWalletResponse(_message.Message):
    __slots__ = ("mnemonic",)
    MNEMONIC_FIELD_NUMBER: _ClassVar[int]
    mnemonic: str
    def __init__(self, mnemonic: _Optional[str] = ...) -> None: ...

class LoadWalletRequest(_message.Message):
    __slots__ = ("wallet_name",)
    WALLET_NAME_FIELD_NUMBER: _ClassVar[int]
    wallet_name: str
    def __init__(self, wallet_name: _Optional[str] = ...) -> None: ...

class LoadWalletResponse(_message.Message):
    __slots__ = ("wallet_name",)
    WALLET_NAME_FIELD_NUMBER: _ClassVar[int]
    wallet_name: str
    def __init__(self, wallet_name: _Optional[str] = ...) -> None: ...

class UnloadWalletRequest(_message.Message):
    __slots__ = ("wallet_name",)
    WALLET_NAME_FIELD_NUMBER: _ClassVar[int]
    wallet_name: str
    def __init__(self, wallet_name: _Optional[str] = ...) -> None: ...

class UnloadWalletResponse(_message.Message):
    __slots__ = ("wallet_name",)
    WALLET_NAME_FIELD_NUMBER: _ClassVar[int]
    wallet_name: str
    def __init__(self, wallet_name: _Optional[str] = ...) -> None: ...

class GetValidatorAddressRequest(_message.Message):
    __slots__ = ("public_key",)
    PUBLIC_KEY_FIELD_NUMBER: _ClassVar[int]
    public_key: str
    def __init__(self, public_key: _Optional[str] = ...) -> None: ...

class GetValidatorAddressResponse(_message.Message):
    __slots__ = ("address",)
    ADDRESS_FIELD_NUMBER: _ClassVar[int]
    address: str
    def __init__(self, address: _Optional[str] = ...) -> None: ...

class SignRawTransactionRequest(_message.Message):
    __slots__ = ("wallet_name", "raw_transaction", "password")
    WALLET_NAME_FIELD_NUMBER: _ClassVar[int]
    RAW_TRANSACTION_FIELD_NUMBER: _ClassVar[int]
    PASSWORD_FIELD_NUMBER: _ClassVar[int]
    wallet_name: str
    raw_transaction: str
    password: str
    def __init__(self, wallet_name: _Optional[str] = ..., raw_transaction: _Optional[str] = ..., password: _Optional[str] = ...) -> None: ...

class SignRawTransactionResponse(_message.Message):
    __slots__ = ("transaction_id", "signed_raw_transaction")
    TRANSACTION_ID_FIELD_NUMBER: _ClassVar[int]
    SIGNED_RAW_TRANSACTION_FIELD_NUMBER: _ClassVar[int]
    transaction_id: str
    signed_raw_transaction: str
    def __init__(self, transaction_id: _Optional[str] = ..., signed_raw_transaction: _Optional[str] = ...) -> None: ...

class GetTotalBalanceRequest(_message.Message):
    __slots__ = ("wallet_name",)
    WALLET_NAME_FIELD_NUMBER: _ClassVar[int]
    wallet_name: str
    def __init__(self, wallet_name: _Optional[str] = ...) -> None: ...

class GetTotalBalanceResponse(_message.Message):
    __slots__ = ("wallet_name", "total_balance")
    WALLET_NAME_FIELD_NUMBER: _ClassVar[int]
    TOTAL_BALANCE_FIELD_NUMBER: _ClassVar[int]
    wallet_name: str
    total_balance: int
    def __init__(self, wallet_name: _Optional[str] = ..., total_balance: _Optional[int] = ...) -> None: ...

class SignMessageRequest(_message.Message):
    __slots__ = ("wallet_name", "password", "address", "message")
    WALLET_NAME_FIELD_NUMBER: _ClassVar[int]
    PASSWORD_FIELD_NUMBER: _ClassVar[int]
    ADDRESS_FIELD_NUMBER: _ClassVar[int]
    MESSAGE_FIELD_NUMBER: _ClassVar[int]
    wallet_name: str
    password: str
    address: str
    message: str
    def __init__(self, wallet_name: _Optional[str] = ..., password: _Optional[str] = ..., address: _Optional[str] = ..., message: _Optional[str] = ...) -> None: ...

class SignMessageResponse(_message.Message):
    __slots__ = ("signature",)
    SIGNATURE_FIELD_NUMBER: _ClassVar[int]
    signature: str
    def __init__(self, signature: _Optional[str] = ...) -> None: ...
