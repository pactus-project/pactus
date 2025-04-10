import transaction_pb2 as _transaction_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable, Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class BlockVerbosity(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    BLOCK_VERBOSITY_DATA: _ClassVar[BlockVerbosity]
    BLOCK_VERBOSITY_INFO: _ClassVar[BlockVerbosity]
    BLOCK_VERBOSITY_TRANSACTIONS: _ClassVar[BlockVerbosity]

class VoteType(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    VOTE_TYPE_UNSPECIFIED: _ClassVar[VoteType]
    VOTE_TYPE_PREPARE: _ClassVar[VoteType]
    VOTE_TYPE_PRECOMMIT: _ClassVar[VoteType]
    VOTE_TYPE_CP_PRE_VOTE: _ClassVar[VoteType]
    VOTE_TYPE_CP_MAIN_VOTE: _ClassVar[VoteType]
    VOTE_TYPE_CP_DECIDED: _ClassVar[VoteType]
BLOCK_VERBOSITY_DATA: BlockVerbosity
BLOCK_VERBOSITY_INFO: BlockVerbosity
BLOCK_VERBOSITY_TRANSACTIONS: BlockVerbosity
VOTE_TYPE_UNSPECIFIED: VoteType
VOTE_TYPE_PREPARE: VoteType
VOTE_TYPE_PRECOMMIT: VoteType
VOTE_TYPE_CP_PRE_VOTE: VoteType
VOTE_TYPE_CP_MAIN_VOTE: VoteType
VOTE_TYPE_CP_DECIDED: VoteType

class GetAccountRequest(_message.Message):
    __slots__ = ("address",)
    ADDRESS_FIELD_NUMBER: _ClassVar[int]
    address: str
    def __init__(self, address: _Optional[str] = ...) -> None: ...

class GetAccountResponse(_message.Message):
    __slots__ = ("account",)
    ACCOUNT_FIELD_NUMBER: _ClassVar[int]
    account: AccountInfo
    def __init__(self, account: _Optional[_Union[AccountInfo, _Mapping]] = ...) -> None: ...

class GetValidatorAddressesRequest(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class GetValidatorAddressesResponse(_message.Message):
    __slots__ = ("addresses",)
    ADDRESSES_FIELD_NUMBER: _ClassVar[int]
    addresses: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, addresses: _Optional[_Iterable[str]] = ...) -> None: ...

class GetValidatorRequest(_message.Message):
    __slots__ = ("address",)
    ADDRESS_FIELD_NUMBER: _ClassVar[int]
    address: str
    def __init__(self, address: _Optional[str] = ...) -> None: ...

class GetValidatorByNumberRequest(_message.Message):
    __slots__ = ("number",)
    NUMBER_FIELD_NUMBER: _ClassVar[int]
    number: int
    def __init__(self, number: _Optional[int] = ...) -> None: ...

class GetValidatorResponse(_message.Message):
    __slots__ = ("validator",)
    VALIDATOR_FIELD_NUMBER: _ClassVar[int]
    validator: ValidatorInfo
    def __init__(self, validator: _Optional[_Union[ValidatorInfo, _Mapping]] = ...) -> None: ...

class GetPublicKeyRequest(_message.Message):
    __slots__ = ("address",)
    ADDRESS_FIELD_NUMBER: _ClassVar[int]
    address: str
    def __init__(self, address: _Optional[str] = ...) -> None: ...

class GetPublicKeyResponse(_message.Message):
    __slots__ = ("public_key",)
    PUBLIC_KEY_FIELD_NUMBER: _ClassVar[int]
    public_key: str
    def __init__(self, public_key: _Optional[str] = ...) -> None: ...

class GetBlockRequest(_message.Message):
    __slots__ = ("height", "verbosity")
    HEIGHT_FIELD_NUMBER: _ClassVar[int]
    VERBOSITY_FIELD_NUMBER: _ClassVar[int]
    height: int
    verbosity: BlockVerbosity
    def __init__(self, height: _Optional[int] = ..., verbosity: _Optional[_Union[BlockVerbosity, str]] = ...) -> None: ...

class GetBlockResponse(_message.Message):
    __slots__ = ("height", "hash", "data", "block_time", "header", "prev_cert", "txs")
    HEIGHT_FIELD_NUMBER: _ClassVar[int]
    HASH_FIELD_NUMBER: _ClassVar[int]
    DATA_FIELD_NUMBER: _ClassVar[int]
    BLOCK_TIME_FIELD_NUMBER: _ClassVar[int]
    HEADER_FIELD_NUMBER: _ClassVar[int]
    PREV_CERT_FIELD_NUMBER: _ClassVar[int]
    TXS_FIELD_NUMBER: _ClassVar[int]
    height: int
    hash: str
    data: str
    block_time: int
    header: BlockHeaderInfo
    prev_cert: CertificateInfo
    txs: _containers.RepeatedCompositeFieldContainer[_transaction_pb2.TransactionInfo]
    def __init__(self, height: _Optional[int] = ..., hash: _Optional[str] = ..., data: _Optional[str] = ..., block_time: _Optional[int] = ..., header: _Optional[_Union[BlockHeaderInfo, _Mapping]] = ..., prev_cert: _Optional[_Union[CertificateInfo, _Mapping]] = ..., txs: _Optional[_Iterable[_Union[_transaction_pb2.TransactionInfo, _Mapping]]] = ...) -> None: ...

class GetBlockHashRequest(_message.Message):
    __slots__ = ("height",)
    HEIGHT_FIELD_NUMBER: _ClassVar[int]
    height: int
    def __init__(self, height: _Optional[int] = ...) -> None: ...

class GetBlockHashResponse(_message.Message):
    __slots__ = ("hash",)
    HASH_FIELD_NUMBER: _ClassVar[int]
    hash: str
    def __init__(self, hash: _Optional[str] = ...) -> None: ...

class GetBlockHeightRequest(_message.Message):
    __slots__ = ("hash",)
    HASH_FIELD_NUMBER: _ClassVar[int]
    hash: str
    def __init__(self, hash: _Optional[str] = ...) -> None: ...

class GetBlockHeightResponse(_message.Message):
    __slots__ = ("height",)
    HEIGHT_FIELD_NUMBER: _ClassVar[int]
    height: int
    def __init__(self, height: _Optional[int] = ...) -> None: ...

class GetBlockchainInfoRequest(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class GetBlockchainInfoResponse(_message.Message):
    __slots__ = ("last_block_height", "last_block_hash", "total_accounts", "total_validators", "total_power", "committee_power", "committee_validators", "is_pruned", "pruning_height", "last_block_time")
    LAST_BLOCK_HEIGHT_FIELD_NUMBER: _ClassVar[int]
    LAST_BLOCK_HASH_FIELD_NUMBER: _ClassVar[int]
    TOTAL_ACCOUNTS_FIELD_NUMBER: _ClassVar[int]
    TOTAL_VALIDATORS_FIELD_NUMBER: _ClassVar[int]
    TOTAL_POWER_FIELD_NUMBER: _ClassVar[int]
    COMMITTEE_POWER_FIELD_NUMBER: _ClassVar[int]
    COMMITTEE_VALIDATORS_FIELD_NUMBER: _ClassVar[int]
    IS_PRUNED_FIELD_NUMBER: _ClassVar[int]
    PRUNING_HEIGHT_FIELD_NUMBER: _ClassVar[int]
    LAST_BLOCK_TIME_FIELD_NUMBER: _ClassVar[int]
    last_block_height: int
    last_block_hash: str
    total_accounts: int
    total_validators: int
    total_power: int
    committee_power: int
    committee_validators: _containers.RepeatedCompositeFieldContainer[ValidatorInfo]
    is_pruned: bool
    pruning_height: int
    last_block_time: int
    def __init__(self, last_block_height: _Optional[int] = ..., last_block_hash: _Optional[str] = ..., total_accounts: _Optional[int] = ..., total_validators: _Optional[int] = ..., total_power: _Optional[int] = ..., committee_power: _Optional[int] = ..., committee_validators: _Optional[_Iterable[_Union[ValidatorInfo, _Mapping]]] = ..., is_pruned: bool = ..., pruning_height: _Optional[int] = ..., last_block_time: _Optional[int] = ...) -> None: ...

class GetConsensusInfoRequest(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class GetConsensusInfoResponse(_message.Message):
    __slots__ = ("proposal", "instances")
    PROPOSAL_FIELD_NUMBER: _ClassVar[int]
    INSTANCES_FIELD_NUMBER: _ClassVar[int]
    proposal: ProposalInfo
    instances: _containers.RepeatedCompositeFieldContainer[ConsensusInfo]
    def __init__(self, proposal: _Optional[_Union[ProposalInfo, _Mapping]] = ..., instances: _Optional[_Iterable[_Union[ConsensusInfo, _Mapping]]] = ...) -> None: ...

class GetTxPoolContentRequest(_message.Message):
    __slots__ = ("payload_type",)
    PAYLOAD_TYPE_FIELD_NUMBER: _ClassVar[int]
    payload_type: _transaction_pb2.PayloadType
    def __init__(self, payload_type: _Optional[_Union[_transaction_pb2.PayloadType, str]] = ...) -> None: ...

class GetTxPoolContentResponse(_message.Message):
    __slots__ = ("txs",)
    TXS_FIELD_NUMBER: _ClassVar[int]
    txs: _containers.RepeatedCompositeFieldContainer[_transaction_pb2.TransactionInfo]
    def __init__(self, txs: _Optional[_Iterable[_Union[_transaction_pb2.TransactionInfo, _Mapping]]] = ...) -> None: ...

class ValidatorInfo(_message.Message):
    __slots__ = ("hash", "data", "public_key", "number", "stake", "last_bonding_height", "last_sortition_height", "unbonding_height", "address", "availability_score")
    HASH_FIELD_NUMBER: _ClassVar[int]
    DATA_FIELD_NUMBER: _ClassVar[int]
    PUBLIC_KEY_FIELD_NUMBER: _ClassVar[int]
    NUMBER_FIELD_NUMBER: _ClassVar[int]
    STAKE_FIELD_NUMBER: _ClassVar[int]
    LAST_BONDING_HEIGHT_FIELD_NUMBER: _ClassVar[int]
    LAST_SORTITION_HEIGHT_FIELD_NUMBER: _ClassVar[int]
    UNBONDING_HEIGHT_FIELD_NUMBER: _ClassVar[int]
    ADDRESS_FIELD_NUMBER: _ClassVar[int]
    AVAILABILITY_SCORE_FIELD_NUMBER: _ClassVar[int]
    hash: str
    data: str
    public_key: str
    number: int
    stake: int
    last_bonding_height: int
    last_sortition_height: int
    unbonding_height: int
    address: str
    availability_score: float
    def __init__(self, hash: _Optional[str] = ..., data: _Optional[str] = ..., public_key: _Optional[str] = ..., number: _Optional[int] = ..., stake: _Optional[int] = ..., last_bonding_height: _Optional[int] = ..., last_sortition_height: _Optional[int] = ..., unbonding_height: _Optional[int] = ..., address: _Optional[str] = ..., availability_score: _Optional[float] = ...) -> None: ...

class AccountInfo(_message.Message):
    __slots__ = ("hash", "data", "number", "balance", "address")
    HASH_FIELD_NUMBER: _ClassVar[int]
    DATA_FIELD_NUMBER: _ClassVar[int]
    NUMBER_FIELD_NUMBER: _ClassVar[int]
    BALANCE_FIELD_NUMBER: _ClassVar[int]
    ADDRESS_FIELD_NUMBER: _ClassVar[int]
    hash: str
    data: str
    number: int
    balance: int
    address: str
    def __init__(self, hash: _Optional[str] = ..., data: _Optional[str] = ..., number: _Optional[int] = ..., balance: _Optional[int] = ..., address: _Optional[str] = ...) -> None: ...

class BlockHeaderInfo(_message.Message):
    __slots__ = ("version", "prev_block_hash", "state_root", "sortition_seed", "proposer_address")
    VERSION_FIELD_NUMBER: _ClassVar[int]
    PREV_BLOCK_HASH_FIELD_NUMBER: _ClassVar[int]
    STATE_ROOT_FIELD_NUMBER: _ClassVar[int]
    SORTITION_SEED_FIELD_NUMBER: _ClassVar[int]
    PROPOSER_ADDRESS_FIELD_NUMBER: _ClassVar[int]
    version: int
    prev_block_hash: str
    state_root: str
    sortition_seed: str
    proposer_address: str
    def __init__(self, version: _Optional[int] = ..., prev_block_hash: _Optional[str] = ..., state_root: _Optional[str] = ..., sortition_seed: _Optional[str] = ..., proposer_address: _Optional[str] = ...) -> None: ...

class CertificateInfo(_message.Message):
    __slots__ = ("hash", "round", "committers", "absentees", "signature")
    HASH_FIELD_NUMBER: _ClassVar[int]
    ROUND_FIELD_NUMBER: _ClassVar[int]
    COMMITTERS_FIELD_NUMBER: _ClassVar[int]
    ABSENTEES_FIELD_NUMBER: _ClassVar[int]
    SIGNATURE_FIELD_NUMBER: _ClassVar[int]
    hash: str
    round: int
    committers: _containers.RepeatedScalarFieldContainer[int]
    absentees: _containers.RepeatedScalarFieldContainer[int]
    signature: str
    def __init__(self, hash: _Optional[str] = ..., round: _Optional[int] = ..., committers: _Optional[_Iterable[int]] = ..., absentees: _Optional[_Iterable[int]] = ..., signature: _Optional[str] = ...) -> None: ...

class VoteInfo(_message.Message):
    __slots__ = ("type", "voter", "block_hash", "round", "cp_round", "cp_value")
    TYPE_FIELD_NUMBER: _ClassVar[int]
    VOTER_FIELD_NUMBER: _ClassVar[int]
    BLOCK_HASH_FIELD_NUMBER: _ClassVar[int]
    ROUND_FIELD_NUMBER: _ClassVar[int]
    CP_ROUND_FIELD_NUMBER: _ClassVar[int]
    CP_VALUE_FIELD_NUMBER: _ClassVar[int]
    type: VoteType
    voter: str
    block_hash: str
    round: int
    cp_round: int
    cp_value: int
    def __init__(self, type: _Optional[_Union[VoteType, str]] = ..., voter: _Optional[str] = ..., block_hash: _Optional[str] = ..., round: _Optional[int] = ..., cp_round: _Optional[int] = ..., cp_value: _Optional[int] = ...) -> None: ...

class ConsensusInfo(_message.Message):
    __slots__ = ("address", "active", "height", "round", "votes")
    ADDRESS_FIELD_NUMBER: _ClassVar[int]
    ACTIVE_FIELD_NUMBER: _ClassVar[int]
    HEIGHT_FIELD_NUMBER: _ClassVar[int]
    ROUND_FIELD_NUMBER: _ClassVar[int]
    VOTES_FIELD_NUMBER: _ClassVar[int]
    address: str
    active: bool
    height: int
    round: int
    votes: _containers.RepeatedCompositeFieldContainer[VoteInfo]
    def __init__(self, address: _Optional[str] = ..., active: bool = ..., height: _Optional[int] = ..., round: _Optional[int] = ..., votes: _Optional[_Iterable[_Union[VoteInfo, _Mapping]]] = ...) -> None: ...

class ProposalInfo(_message.Message):
    __slots__ = ("height", "round", "block_data", "signature")
    HEIGHT_FIELD_NUMBER: _ClassVar[int]
    ROUND_FIELD_NUMBER: _ClassVar[int]
    BLOCK_DATA_FIELD_NUMBER: _ClassVar[int]
    SIGNATURE_FIELD_NUMBER: _ClassVar[int]
    height: int
    round: int
    block_data: str
    signature: str
    def __init__(self, height: _Optional[int] = ..., round: _Optional[int] = ..., block_data: _Optional[str] = ..., signature: _Optional[str] = ...) -> None: ...
