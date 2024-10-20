from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class PayloadType(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    UNKNOWN: _ClassVar[PayloadType]
    TRANSFER_PAYLOAD: _ClassVar[PayloadType]
    BOND_PAYLOAD: _ClassVar[PayloadType]
    SORTITION_PAYLOAD: _ClassVar[PayloadType]
    UNBOND_PAYLOAD: _ClassVar[PayloadType]
    WITHDRAW_PAYLOAD: _ClassVar[PayloadType]

class TransactionVerbosity(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    TRANSACTION_DATA: _ClassVar[TransactionVerbosity]
    TRANSACTION_INFO: _ClassVar[TransactionVerbosity]
UNKNOWN: PayloadType
TRANSFER_PAYLOAD: PayloadType
BOND_PAYLOAD: PayloadType
SORTITION_PAYLOAD: PayloadType
UNBOND_PAYLOAD: PayloadType
WITHDRAW_PAYLOAD: PayloadType
TRANSACTION_DATA: TransactionVerbosity
TRANSACTION_INFO: TransactionVerbosity

class GetTransactionRequest(_message.Message):
    __slots__ = ("id", "verbosity")
    ID_FIELD_NUMBER: _ClassVar[int]
    VERBOSITY_FIELD_NUMBER: _ClassVar[int]
    id: str
    verbosity: TransactionVerbosity
    def __init__(self, id: _Optional[str] = ..., verbosity: _Optional[_Union[TransactionVerbosity, str]] = ...) -> None: ...

class GetTransactionResponse(_message.Message):
    __slots__ = ("block_height", "block_time", "transaction")
    BLOCK_HEIGHT_FIELD_NUMBER: _ClassVar[int]
    BLOCK_TIME_FIELD_NUMBER: _ClassVar[int]
    TRANSACTION_FIELD_NUMBER: _ClassVar[int]
    block_height: int
    block_time: int
    transaction: TransactionInfo
    def __init__(self, block_height: _Optional[int] = ..., block_time: _Optional[int] = ..., transaction: _Optional[_Union[TransactionInfo, _Mapping]] = ...) -> None: ...

class CalculateFeeRequest(_message.Message):
    __slots__ = ("amount", "payload_type", "fixed_amount")
    AMOUNT_FIELD_NUMBER: _ClassVar[int]
    PAYLOAD_TYPE_FIELD_NUMBER: _ClassVar[int]
    FIXED_AMOUNT_FIELD_NUMBER: _ClassVar[int]
    amount: int
    payload_type: PayloadType
    fixed_amount: bool
    def __init__(self, amount: _Optional[int] = ..., payload_type: _Optional[_Union[PayloadType, str]] = ..., fixed_amount: bool = ...) -> None: ...

class CalculateFeeResponse(_message.Message):
    __slots__ = ("amount", "fee")
    AMOUNT_FIELD_NUMBER: _ClassVar[int]
    FEE_FIELD_NUMBER: _ClassVar[int]
    amount: int
    fee: int
    def __init__(self, amount: _Optional[int] = ..., fee: _Optional[int] = ...) -> None: ...

class BroadcastTransactionRequest(_message.Message):
    __slots__ = ("signed_raw_transaction",)
    SIGNED_RAW_TRANSACTION_FIELD_NUMBER: _ClassVar[int]
    signed_raw_transaction: str
    def __init__(self, signed_raw_transaction: _Optional[str] = ...) -> None: ...

class BroadcastTransactionResponse(_message.Message):
    __slots__ = ("id",)
    ID_FIELD_NUMBER: _ClassVar[int]
    id: str
    def __init__(self, id: _Optional[str] = ...) -> None: ...

class GetRawTransferTransactionRequest(_message.Message):
    __slots__ = ("lock_time", "sender", "receiver", "amount", "fee", "memo")
    LOCK_TIME_FIELD_NUMBER: _ClassVar[int]
    SENDER_FIELD_NUMBER: _ClassVar[int]
    RECEIVER_FIELD_NUMBER: _ClassVar[int]
    AMOUNT_FIELD_NUMBER: _ClassVar[int]
    FEE_FIELD_NUMBER: _ClassVar[int]
    MEMO_FIELD_NUMBER: _ClassVar[int]
    lock_time: int
    sender: str
    receiver: str
    amount: int
    fee: int
    memo: str
    def __init__(self, lock_time: _Optional[int] = ..., sender: _Optional[str] = ..., receiver: _Optional[str] = ..., amount: _Optional[int] = ..., fee: _Optional[int] = ..., memo: _Optional[str] = ...) -> None: ...

class GetRawBondTransactionRequest(_message.Message):
    __slots__ = ("lock_time", "sender", "receiver", "stake", "public_key", "fee", "memo")
    LOCK_TIME_FIELD_NUMBER: _ClassVar[int]
    SENDER_FIELD_NUMBER: _ClassVar[int]
    RECEIVER_FIELD_NUMBER: _ClassVar[int]
    STAKE_FIELD_NUMBER: _ClassVar[int]
    PUBLIC_KEY_FIELD_NUMBER: _ClassVar[int]
    FEE_FIELD_NUMBER: _ClassVar[int]
    MEMO_FIELD_NUMBER: _ClassVar[int]
    lock_time: int
    sender: str
    receiver: str
    stake: int
    public_key: str
    fee: int
    memo: str
    def __init__(self, lock_time: _Optional[int] = ..., sender: _Optional[str] = ..., receiver: _Optional[str] = ..., stake: _Optional[int] = ..., public_key: _Optional[str] = ..., fee: _Optional[int] = ..., memo: _Optional[str] = ...) -> None: ...

class GetRawUnbondTransactionRequest(_message.Message):
    __slots__ = ("lock_time", "validator_address", "memo")
    LOCK_TIME_FIELD_NUMBER: _ClassVar[int]
    VALIDATOR_ADDRESS_FIELD_NUMBER: _ClassVar[int]
    MEMO_FIELD_NUMBER: _ClassVar[int]
    lock_time: int
    validator_address: str
    memo: str
    def __init__(self, lock_time: _Optional[int] = ..., validator_address: _Optional[str] = ..., memo: _Optional[str] = ...) -> None: ...

class GetRawWithdrawTransactionRequest(_message.Message):
    __slots__ = ("lock_time", "validator_address", "account_address", "amount", "fee", "memo")
    LOCK_TIME_FIELD_NUMBER: _ClassVar[int]
    VALIDATOR_ADDRESS_FIELD_NUMBER: _ClassVar[int]
    ACCOUNT_ADDRESS_FIELD_NUMBER: _ClassVar[int]
    AMOUNT_FIELD_NUMBER: _ClassVar[int]
    FEE_FIELD_NUMBER: _ClassVar[int]
    MEMO_FIELD_NUMBER: _ClassVar[int]
    lock_time: int
    validator_address: str
    account_address: str
    amount: int
    fee: int
    memo: str
    def __init__(self, lock_time: _Optional[int] = ..., validator_address: _Optional[str] = ..., account_address: _Optional[str] = ..., amount: _Optional[int] = ..., fee: _Optional[int] = ..., memo: _Optional[str] = ...) -> None: ...

class GetRawTransactionResponse(_message.Message):
    __slots__ = ("raw_transaction", "id")
    RAW_TRANSACTION_FIELD_NUMBER: _ClassVar[int]
    ID_FIELD_NUMBER: _ClassVar[int]
    raw_transaction: str
    id: str
    def __init__(self, raw_transaction: _Optional[str] = ..., id: _Optional[str] = ...) -> None: ...

class PayloadTransfer(_message.Message):
    __slots__ = ("sender", "receiver", "amount")
    SENDER_FIELD_NUMBER: _ClassVar[int]
    RECEIVER_FIELD_NUMBER: _ClassVar[int]
    AMOUNT_FIELD_NUMBER: _ClassVar[int]
    sender: str
    receiver: str
    amount: int
    def __init__(self, sender: _Optional[str] = ..., receiver: _Optional[str] = ..., amount: _Optional[int] = ...) -> None: ...

class PayloadBond(_message.Message):
    __slots__ = ("sender", "receiver", "stake", "public_key")
    SENDER_FIELD_NUMBER: _ClassVar[int]
    RECEIVER_FIELD_NUMBER: _ClassVar[int]
    STAKE_FIELD_NUMBER: _ClassVar[int]
    PUBLIC_KEY_FIELD_NUMBER: _ClassVar[int]
    sender: str
    receiver: str
    stake: int
    public_key: str
    def __init__(self, sender: _Optional[str] = ..., receiver: _Optional[str] = ..., stake: _Optional[int] = ..., public_key: _Optional[str] = ...) -> None: ...

class PayloadSortition(_message.Message):
    __slots__ = ("address", "proof")
    ADDRESS_FIELD_NUMBER: _ClassVar[int]
    PROOF_FIELD_NUMBER: _ClassVar[int]
    address: str
    proof: str
    def __init__(self, address: _Optional[str] = ..., proof: _Optional[str] = ...) -> None: ...

class PayloadUnbond(_message.Message):
    __slots__ = ("validator",)
    VALIDATOR_FIELD_NUMBER: _ClassVar[int]
    validator: str
    def __init__(self, validator: _Optional[str] = ...) -> None: ...

class PayloadWithdraw(_message.Message):
    __slots__ = ("validator_address", "account_address", "amount")
    VALIDATOR_ADDRESS_FIELD_NUMBER: _ClassVar[int]
    ACCOUNT_ADDRESS_FIELD_NUMBER: _ClassVar[int]
    AMOUNT_FIELD_NUMBER: _ClassVar[int]
    validator_address: str
    account_address: str
    amount: int
    def __init__(self, validator_address: _Optional[str] = ..., account_address: _Optional[str] = ..., amount: _Optional[int] = ...) -> None: ...

class TransactionInfo(_message.Message):
    __slots__ = ("id", "data", "version", "lock_time", "value", "fee", "payload_type", "transfer", "bond", "sortition", "unbond", "withdraw", "memo", "public_key", "signature")
    ID_FIELD_NUMBER: _ClassVar[int]
    DATA_FIELD_NUMBER: _ClassVar[int]
    VERSION_FIELD_NUMBER: _ClassVar[int]
    LOCK_TIME_FIELD_NUMBER: _ClassVar[int]
    VALUE_FIELD_NUMBER: _ClassVar[int]
    FEE_FIELD_NUMBER: _ClassVar[int]
    PAYLOAD_TYPE_FIELD_NUMBER: _ClassVar[int]
    TRANSFER_FIELD_NUMBER: _ClassVar[int]
    BOND_FIELD_NUMBER: _ClassVar[int]
    SORTITION_FIELD_NUMBER: _ClassVar[int]
    UNBOND_FIELD_NUMBER: _ClassVar[int]
    WITHDRAW_FIELD_NUMBER: _ClassVar[int]
    MEMO_FIELD_NUMBER: _ClassVar[int]
    PUBLIC_KEY_FIELD_NUMBER: _ClassVar[int]
    SIGNATURE_FIELD_NUMBER: _ClassVar[int]
    id: str
    data: str
    version: int
    lock_time: int
    value: int
    fee: int
    payload_type: PayloadType
    transfer: PayloadTransfer
    bond: PayloadBond
    sortition: PayloadSortition
    unbond: PayloadUnbond
    withdraw: PayloadWithdraw
    memo: str
    public_key: str
    signature: str
    def __init__(self, id: _Optional[str] = ..., data: _Optional[str] = ..., version: _Optional[int] = ..., lock_time: _Optional[int] = ..., value: _Optional[int] = ..., fee: _Optional[int] = ..., payload_type: _Optional[_Union[PayloadType, str]] = ..., transfer: _Optional[_Union[PayloadTransfer, _Mapping]] = ..., bond: _Optional[_Union[PayloadBond, _Mapping]] = ..., sortition: _Optional[_Union[PayloadSortition, _Mapping]] = ..., unbond: _Optional[_Union[PayloadUnbond, _Mapping]] = ..., withdraw: _Optional[_Union[PayloadWithdraw, _Mapping]] = ..., memo: _Optional[str] = ..., public_key: _Optional[str] = ..., signature: _Optional[str] = ...) -> None: ...
