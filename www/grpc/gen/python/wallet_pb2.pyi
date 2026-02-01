import transaction_pb2 as _transaction_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable, Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class AddressType(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    ADDRESS_TYPE_TREASURY: _ClassVar[AddressType]
    ADDRESS_TYPE_VALIDATOR: _ClassVar[AddressType]
    ADDRESS_TYPE_BLS_ACCOUNT: _ClassVar[AddressType]
    ADDRESS_TYPE_ED25519_ACCOUNT: _ClassVar[AddressType]

class TxDirection(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    TX_DIRECTION_ANY: _ClassVar[TxDirection]
    TX_DIRECTION_INCOMING: _ClassVar[TxDirection]
    TX_DIRECTION_OUTGOING: _ClassVar[TxDirection]

class TransactionStatus(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    TRANSACTION_STATUS_PENDING: _ClassVar[TransactionStatus]
    TRANSACTION_STATUS_CONFIRMED: _ClassVar[TransactionStatus]
    TRANSACTION_STATUS_FAILED: _ClassVar[TransactionStatus]
ADDRESS_TYPE_TREASURY: AddressType
ADDRESS_TYPE_VALIDATOR: AddressType
ADDRESS_TYPE_BLS_ACCOUNT: AddressType
ADDRESS_TYPE_ED25519_ACCOUNT: AddressType
TX_DIRECTION_ANY: TxDirection
TX_DIRECTION_INCOMING: TxDirection
TX_DIRECTION_OUTGOING: TxDirection
TRANSACTION_STATUS_PENDING: TransactionStatus
TRANSACTION_STATUS_CONFIRMED: TransactionStatus
TRANSACTION_STATUS_FAILED: TransactionStatus

class AddressInfo(_message.Message):
    __slots__ = ()
    ADDRESS_FIELD_NUMBER: _ClassVar[int]
    PUBLIC_KEY_FIELD_NUMBER: _ClassVar[int]
    LABEL_FIELD_NUMBER: _ClassVar[int]
    PATH_FIELD_NUMBER: _ClassVar[int]
    address: str
    public_key: str
    label: str
    path: str
    def __init__(self, address: _Optional[str] = ..., public_key: _Optional[str] = ..., label: _Optional[str] = ..., path: _Optional[str] = ...) -> None: ...

class GetNewAddressRequest(_message.Message):
    __slots__ = ()
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
    __slots__ = ()
    WALLET_NAME_FIELD_NUMBER: _ClassVar[int]
    ADDR_FIELD_NUMBER: _ClassVar[int]
    wallet_name: str
    addr: AddressInfo
    def __init__(self, wallet_name: _Optional[str] = ..., addr: _Optional[_Union[AddressInfo, _Mapping]] = ...) -> None: ...

class RestoreWalletRequest(_message.Message):
    __slots__ = ()
    WALLET_NAME_FIELD_NUMBER: _ClassVar[int]
    MNEMONIC_FIELD_NUMBER: _ClassVar[int]
    PASSWORD_FIELD_NUMBER: _ClassVar[int]
    wallet_name: str
    mnemonic: str
    password: str
    def __init__(self, wallet_name: _Optional[str] = ..., mnemonic: _Optional[str] = ..., password: _Optional[str] = ...) -> None: ...

class RestoreWalletResponse(_message.Message):
    __slots__ = ()
    WALLET_NAME_FIELD_NUMBER: _ClassVar[int]
    wallet_name: str
    def __init__(self, wallet_name: _Optional[str] = ...) -> None: ...

class CreateWalletRequest(_message.Message):
    __slots__ = ()
    WALLET_NAME_FIELD_NUMBER: _ClassVar[int]
    PASSWORD_FIELD_NUMBER: _ClassVar[int]
    wallet_name: str
    password: str
    def __init__(self, wallet_name: _Optional[str] = ..., password: _Optional[str] = ...) -> None: ...

class CreateWalletResponse(_message.Message):
    __slots__ = ()
    WALLET_NAME_FIELD_NUMBER: _ClassVar[int]
    MNEMONIC_FIELD_NUMBER: _ClassVar[int]
    wallet_name: str
    mnemonic: str
    def __init__(self, wallet_name: _Optional[str] = ..., mnemonic: _Optional[str] = ...) -> None: ...

class LoadWalletRequest(_message.Message):
    __slots__ = ()
    WALLET_NAME_FIELD_NUMBER: _ClassVar[int]
    wallet_name: str
    def __init__(self, wallet_name: _Optional[str] = ...) -> None: ...

class LoadWalletResponse(_message.Message):
    __slots__ = ()
    WALLET_NAME_FIELD_NUMBER: _ClassVar[int]
    wallet_name: str
    def __init__(self, wallet_name: _Optional[str] = ...) -> None: ...

class UnloadWalletRequest(_message.Message):
    __slots__ = ()
    WALLET_NAME_FIELD_NUMBER: _ClassVar[int]
    wallet_name: str
    def __init__(self, wallet_name: _Optional[str] = ...) -> None: ...

class UnloadWalletResponse(_message.Message):
    __slots__ = ()
    WALLET_NAME_FIELD_NUMBER: _ClassVar[int]
    wallet_name: str
    def __init__(self, wallet_name: _Optional[str] = ...) -> None: ...

class GetValidatorAddressRequest(_message.Message):
    __slots__ = ()
    PUBLIC_KEY_FIELD_NUMBER: _ClassVar[int]
    public_key: str
    def __init__(self, public_key: _Optional[str] = ...) -> None: ...

class GetValidatorAddressResponse(_message.Message):
    __slots__ = ()
    ADDRESS_FIELD_NUMBER: _ClassVar[int]
    address: str
    def __init__(self, address: _Optional[str] = ...) -> None: ...

class SignRawTransactionRequest(_message.Message):
    __slots__ = ()
    WALLET_NAME_FIELD_NUMBER: _ClassVar[int]
    RAW_TRANSACTION_FIELD_NUMBER: _ClassVar[int]
    PASSWORD_FIELD_NUMBER: _ClassVar[int]
    wallet_name: str
    raw_transaction: str
    password: str
    def __init__(self, wallet_name: _Optional[str] = ..., raw_transaction: _Optional[str] = ..., password: _Optional[str] = ...) -> None: ...

class SignRawTransactionResponse(_message.Message):
    __slots__ = ()
    TRANSACTION_ID_FIELD_NUMBER: _ClassVar[int]
    SIGNED_RAW_TRANSACTION_FIELD_NUMBER: _ClassVar[int]
    transaction_id: str
    signed_raw_transaction: str
    def __init__(self, transaction_id: _Optional[str] = ..., signed_raw_transaction: _Optional[str] = ...) -> None: ...

class GetTotalBalanceRequest(_message.Message):
    __slots__ = ()
    WALLET_NAME_FIELD_NUMBER: _ClassVar[int]
    wallet_name: str
    def __init__(self, wallet_name: _Optional[str] = ...) -> None: ...

class GetTotalBalanceResponse(_message.Message):
    __slots__ = ()
    WALLET_NAME_FIELD_NUMBER: _ClassVar[int]
    TOTAL_BALANCE_FIELD_NUMBER: _ClassVar[int]
    wallet_name: str
    total_balance: int
    def __init__(self, wallet_name: _Optional[str] = ..., total_balance: _Optional[int] = ...) -> None: ...

class SignMessageRequest(_message.Message):
    __slots__ = ()
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
    __slots__ = ()
    SIGNATURE_FIELD_NUMBER: _ClassVar[int]
    signature: str
    def __init__(self, signature: _Optional[str] = ...) -> None: ...

class GetTotalStakeRequest(_message.Message):
    __slots__ = ()
    WALLET_NAME_FIELD_NUMBER: _ClassVar[int]
    wallet_name: str
    def __init__(self, wallet_name: _Optional[str] = ...) -> None: ...

class GetTotalStakeResponse(_message.Message):
    __slots__ = ()
    WALLET_NAME_FIELD_NUMBER: _ClassVar[int]
    TOTAL_STAKE_FIELD_NUMBER: _ClassVar[int]
    wallet_name: str
    total_stake: int
    def __init__(self, wallet_name: _Optional[str] = ..., total_stake: _Optional[int] = ...) -> None: ...

class GetAddressInfoRequest(_message.Message):
    __slots__ = ()
    WALLET_NAME_FIELD_NUMBER: _ClassVar[int]
    ADDRESS_FIELD_NUMBER: _ClassVar[int]
    wallet_name: str
    address: str
    def __init__(self, wallet_name: _Optional[str] = ..., address: _Optional[str] = ...) -> None: ...

class GetAddressInfoResponse(_message.Message):
    __slots__ = ()
    WALLET_NAME_FIELD_NUMBER: _ClassVar[int]
    ADDR_FIELD_NUMBER: _ClassVar[int]
    wallet_name: str
    addr: AddressInfo
    def __init__(self, wallet_name: _Optional[str] = ..., addr: _Optional[_Union[AddressInfo, _Mapping]] = ...) -> None: ...

class SetAddressLabelRequest(_message.Message):
    __slots__ = ()
    WALLET_NAME_FIELD_NUMBER: _ClassVar[int]
    PASSWORD_FIELD_NUMBER: _ClassVar[int]
    ADDRESS_FIELD_NUMBER: _ClassVar[int]
    LABEL_FIELD_NUMBER: _ClassVar[int]
    wallet_name: str
    password: str
    address: str
    label: str
    def __init__(self, wallet_name: _Optional[str] = ..., password: _Optional[str] = ..., address: _Optional[str] = ..., label: _Optional[str] = ...) -> None: ...

class SetAddressLabelResponse(_message.Message):
    __slots__ = ()
    WALLET_NAME_FIELD_NUMBER: _ClassVar[int]
    ADDRESS_FIELD_NUMBER: _ClassVar[int]
    LABEL_FIELD_NUMBER: _ClassVar[int]
    wallet_name: str
    address: str
    label: str
    def __init__(self, wallet_name: _Optional[str] = ..., address: _Optional[str] = ..., label: _Optional[str] = ...) -> None: ...

class ListWalletsRequest(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class ListWalletsResponse(_message.Message):
    __slots__ = ()
    WALLETS_FIELD_NUMBER: _ClassVar[int]
    wallets: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, wallets: _Optional[_Iterable[str]] = ...) -> None: ...

class GetWalletInfoRequest(_message.Message):
    __slots__ = ()
    WALLET_NAME_FIELD_NUMBER: _ClassVar[int]
    wallet_name: str
    def __init__(self, wallet_name: _Optional[str] = ...) -> None: ...

class GetWalletInfoResponse(_message.Message):
    __slots__ = ()
    WALLET_NAME_FIELD_NUMBER: _ClassVar[int]
    VERSION_FIELD_NUMBER: _ClassVar[int]
    NETWORK_FIELD_NUMBER: _ClassVar[int]
    ENCRYPTED_FIELD_NUMBER: _ClassVar[int]
    UUID_FIELD_NUMBER: _ClassVar[int]
    CREATED_AT_FIELD_NUMBER: _ClassVar[int]
    DEFAULT_FEE_FIELD_NUMBER: _ClassVar[int]
    DRIVER_FIELD_NUMBER: _ClassVar[int]
    PATH_FIELD_NUMBER: _ClassVar[int]
    wallet_name: str
    version: int
    network: str
    encrypted: bool
    uuid: str
    created_at: int
    default_fee: int
    driver: str
    path: str
    def __init__(self, wallet_name: _Optional[str] = ..., version: _Optional[int] = ..., network: _Optional[str] = ..., encrypted: _Optional[bool] = ..., uuid: _Optional[str] = ..., created_at: _Optional[int] = ..., default_fee: _Optional[int] = ..., driver: _Optional[str] = ..., path: _Optional[str] = ...) -> None: ...

class ListAddressesRequest(_message.Message):
    __slots__ = ()
    WALLET_NAME_FIELD_NUMBER: _ClassVar[int]
    ADDRESS_TYPES_FIELD_NUMBER: _ClassVar[int]
    wallet_name: str
    address_types: _containers.RepeatedScalarFieldContainer[AddressType]
    def __init__(self, wallet_name: _Optional[str] = ..., address_types: _Optional[_Iterable[_Union[AddressType, str]]] = ...) -> None: ...

class ListAddressesResponse(_message.Message):
    __slots__ = ()
    WALLET_NAME_FIELD_NUMBER: _ClassVar[int]
    ADDRS_FIELD_NUMBER: _ClassVar[int]
    wallet_name: str
    addrs: _containers.RepeatedCompositeFieldContainer[AddressInfo]
    def __init__(self, wallet_name: _Optional[str] = ..., addrs: _Optional[_Iterable[_Union[AddressInfo, _Mapping]]] = ...) -> None: ...

class UpdatePasswordRequest(_message.Message):
    __slots__ = ()
    WALLET_NAME_FIELD_NUMBER: _ClassVar[int]
    OLD_PASSWORD_FIELD_NUMBER: _ClassVar[int]
    NEW_PASSWORD_FIELD_NUMBER: _ClassVar[int]
    wallet_name: str
    old_password: str
    new_password: str
    def __init__(self, wallet_name: _Optional[str] = ..., old_password: _Optional[str] = ..., new_password: _Optional[str] = ...) -> None: ...

class UpdatePasswordResponse(_message.Message):
    __slots__ = ()
    WALLET_NAME_FIELD_NUMBER: _ClassVar[int]
    wallet_name: str
    def __init__(self, wallet_name: _Optional[str] = ...) -> None: ...

class WalletTransactionInfo(_message.Message):
    __slots__ = ()
    NO_FIELD_NUMBER: _ClassVar[int]
    TX_ID_FIELD_NUMBER: _ClassVar[int]
    SENDER_FIELD_NUMBER: _ClassVar[int]
    RECEIVER_FIELD_NUMBER: _ClassVar[int]
    DIRECTION_FIELD_NUMBER: _ClassVar[int]
    AMOUNT_FIELD_NUMBER: _ClassVar[int]
    FEE_FIELD_NUMBER: _ClassVar[int]
    MEMO_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    BLOCK_HEIGHT_FIELD_NUMBER: _ClassVar[int]
    PAYLOAD_TYPE_FIELD_NUMBER: _ClassVar[int]
    DATA_FIELD_NUMBER: _ClassVar[int]
    COMMENT_FIELD_NUMBER: _ClassVar[int]
    CREATED_AT_FIELD_NUMBER: _ClassVar[int]
    UPDATED_AT_FIELD_NUMBER: _ClassVar[int]
    no: int
    tx_id: str
    sender: str
    receiver: str
    direction: TxDirection
    amount: int
    fee: int
    memo: str
    status: TransactionStatus
    block_height: int
    payload_type: _transaction_pb2.PayloadType
    data: bytes
    comment: str
    created_at: int
    updated_at: int
    def __init__(self, no: _Optional[int] = ..., tx_id: _Optional[str] = ..., sender: _Optional[str] = ..., receiver: _Optional[str] = ..., direction: _Optional[_Union[TxDirection, str]] = ..., amount: _Optional[int] = ..., fee: _Optional[int] = ..., memo: _Optional[str] = ..., status: _Optional[_Union[TransactionStatus, str]] = ..., block_height: _Optional[int] = ..., payload_type: _Optional[_Union[_transaction_pb2.PayloadType, str]] = ..., data: _Optional[bytes] = ..., comment: _Optional[str] = ..., created_at: _Optional[int] = ..., updated_at: _Optional[int] = ...) -> None: ...

class ListTransactionsRequest(_message.Message):
    __slots__ = ()
    WALLET_NAME_FIELD_NUMBER: _ClassVar[int]
    ADDRESS_FIELD_NUMBER: _ClassVar[int]
    DIRECTION_FIELD_NUMBER: _ClassVar[int]
    COUNT_FIELD_NUMBER: _ClassVar[int]
    SKIP_FIELD_NUMBER: _ClassVar[int]
    wallet_name: str
    address: str
    direction: TxDirection
    count: int
    skip: int
    def __init__(self, wallet_name: _Optional[str] = ..., address: _Optional[str] = ..., direction: _Optional[_Union[TxDirection, str]] = ..., count: _Optional[int] = ..., skip: _Optional[int] = ...) -> None: ...

class ListTransactionsResponse(_message.Message):
    __slots__ = ()
    WALLET_NAME_FIELD_NUMBER: _ClassVar[int]
    TXS_FIELD_NUMBER: _ClassVar[int]
    wallet_name: str
    txs: _containers.RepeatedCompositeFieldContainer[WalletTransactionInfo]
    def __init__(self, wallet_name: _Optional[str] = ..., txs: _Optional[_Iterable[_Union[WalletTransactionInfo, _Mapping]]] = ...) -> None: ...
