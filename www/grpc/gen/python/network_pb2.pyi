from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class GetNetworkInfoRequest(_message.Message):
    __slots__ = ("only_connected",)
    ONLY_CONNECTED_FIELD_NUMBER: _ClassVar[int]
    only_connected: bool
    def __init__(self, only_connected: bool = ...) -> None: ...

class GetNetworkInfoResponse(_message.Message):
    __slots__ = ("network_name", "total_sent_bytes", "total_received_bytes", "connected_peers_count", "connected_peers", "sent_bytes", "received_bytes")
    class SentBytesEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: int
        value: int
        def __init__(self, key: _Optional[int] = ..., value: _Optional[int] = ...) -> None: ...
    class ReceivedBytesEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: int
        value: int
        def __init__(self, key: _Optional[int] = ..., value: _Optional[int] = ...) -> None: ...
    NETWORK_NAME_FIELD_NUMBER: _ClassVar[int]
    TOTAL_SENT_BYTES_FIELD_NUMBER: _ClassVar[int]
    TOTAL_RECEIVED_BYTES_FIELD_NUMBER: _ClassVar[int]
    CONNECTED_PEERS_COUNT_FIELD_NUMBER: _ClassVar[int]
    CONNECTED_PEERS_FIELD_NUMBER: _ClassVar[int]
    SENT_BYTES_FIELD_NUMBER: _ClassVar[int]
    RECEIVED_BYTES_FIELD_NUMBER: _ClassVar[int]
    network_name: str
    total_sent_bytes: int
    total_received_bytes: int
    connected_peers_count: int
    connected_peers: _containers.RepeatedCompositeFieldContainer[PeerInfo]
    sent_bytes: _containers.ScalarMap[int, int]
    received_bytes: _containers.ScalarMap[int, int]
    def __init__(self, network_name: _Optional[str] = ..., total_sent_bytes: _Optional[int] = ..., total_received_bytes: _Optional[int] = ..., connected_peers_count: _Optional[int] = ..., connected_peers: _Optional[_Iterable[_Union[PeerInfo, _Mapping]]] = ..., sent_bytes: _Optional[_Mapping[int, int]] = ..., received_bytes: _Optional[_Mapping[int, int]] = ...) -> None: ...

class GetNodeInfoRequest(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class GetNodeInfoResponse(_message.Message):
    __slots__ = ("moniker", "agent", "peer_id", "started_at", "reachability", "services", "services_names", "local_addrs", "protocols", "clock_offset", "connection_info")
    MONIKER_FIELD_NUMBER: _ClassVar[int]
    AGENT_FIELD_NUMBER: _ClassVar[int]
    PEER_ID_FIELD_NUMBER: _ClassVar[int]
    STARTED_AT_FIELD_NUMBER: _ClassVar[int]
    REACHABILITY_FIELD_NUMBER: _ClassVar[int]
    SERVICES_FIELD_NUMBER: _ClassVar[int]
    SERVICES_NAMES_FIELD_NUMBER: _ClassVar[int]
    LOCAL_ADDRS_FIELD_NUMBER: _ClassVar[int]
    PROTOCOLS_FIELD_NUMBER: _ClassVar[int]
    CLOCK_OFFSET_FIELD_NUMBER: _ClassVar[int]
    CONNECTION_INFO_FIELD_NUMBER: _ClassVar[int]
    moniker: str
    agent: str
    peer_id: str
    started_at: int
    reachability: str
    services: int
    services_names: str
    local_addrs: _containers.RepeatedScalarFieldContainer[str]
    protocols: _containers.RepeatedScalarFieldContainer[str]
    clock_offset: float
    connection_info: ConnectionInfo
    def __init__(self, moniker: _Optional[str] = ..., agent: _Optional[str] = ..., peer_id: _Optional[str] = ..., started_at: _Optional[int] = ..., reachability: _Optional[str] = ..., services: _Optional[int] = ..., services_names: _Optional[str] = ..., local_addrs: _Optional[_Iterable[str]] = ..., protocols: _Optional[_Iterable[str]] = ..., clock_offset: _Optional[float] = ..., connection_info: _Optional[_Union[ConnectionInfo, _Mapping]] = ...) -> None: ...

class PeerInfo(_message.Message):
    __slots__ = ("status", "moniker", "agent", "peer_id", "consensus_keys", "consensus_addresses", "services", "last_block_hash", "height", "received_bundles", "invalid_bundles", "last_sent", "last_received", "sent_bytes", "received_bytes", "address", "direction", "protocols", "total_sessions", "completed_sessions")
    class SentBytesEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: int
        value: int
        def __init__(self, key: _Optional[int] = ..., value: _Optional[int] = ...) -> None: ...
    class ReceivedBytesEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: int
        value: int
        def __init__(self, key: _Optional[int] = ..., value: _Optional[int] = ...) -> None: ...
    STATUS_FIELD_NUMBER: _ClassVar[int]
    MONIKER_FIELD_NUMBER: _ClassVar[int]
    AGENT_FIELD_NUMBER: _ClassVar[int]
    PEER_ID_FIELD_NUMBER: _ClassVar[int]
    CONSENSUS_KEYS_FIELD_NUMBER: _ClassVar[int]
    CONSENSUS_ADDRESSES_FIELD_NUMBER: _ClassVar[int]
    SERVICES_FIELD_NUMBER: _ClassVar[int]
    LAST_BLOCK_HASH_FIELD_NUMBER: _ClassVar[int]
    HEIGHT_FIELD_NUMBER: _ClassVar[int]
    RECEIVED_BUNDLES_FIELD_NUMBER: _ClassVar[int]
    INVALID_BUNDLES_FIELD_NUMBER: _ClassVar[int]
    LAST_SENT_FIELD_NUMBER: _ClassVar[int]
    LAST_RECEIVED_FIELD_NUMBER: _ClassVar[int]
    SENT_BYTES_FIELD_NUMBER: _ClassVar[int]
    RECEIVED_BYTES_FIELD_NUMBER: _ClassVar[int]
    ADDRESS_FIELD_NUMBER: _ClassVar[int]
    DIRECTION_FIELD_NUMBER: _ClassVar[int]
    PROTOCOLS_FIELD_NUMBER: _ClassVar[int]
    TOTAL_SESSIONS_FIELD_NUMBER: _ClassVar[int]
    COMPLETED_SESSIONS_FIELD_NUMBER: _ClassVar[int]
    status: int
    moniker: str
    agent: str
    peer_id: str
    consensus_keys: _containers.RepeatedScalarFieldContainer[str]
    consensus_addresses: _containers.RepeatedScalarFieldContainer[str]
    services: int
    last_block_hash: str
    height: int
    received_bundles: int
    invalid_bundles: int
    last_sent: int
    last_received: int
    sent_bytes: _containers.ScalarMap[int, int]
    received_bytes: _containers.ScalarMap[int, int]
    address: str
    direction: str
    protocols: _containers.RepeatedScalarFieldContainer[str]
    total_sessions: int
    completed_sessions: int
    def __init__(self, status: _Optional[int] = ..., moniker: _Optional[str] = ..., agent: _Optional[str] = ..., peer_id: _Optional[str] = ..., consensus_keys: _Optional[_Iterable[str]] = ..., consensus_addresses: _Optional[_Iterable[str]] = ..., services: _Optional[int] = ..., last_block_hash: _Optional[str] = ..., height: _Optional[int] = ..., received_bundles: _Optional[int] = ..., invalid_bundles: _Optional[int] = ..., last_sent: _Optional[int] = ..., last_received: _Optional[int] = ..., sent_bytes: _Optional[_Mapping[int, int]] = ..., received_bytes: _Optional[_Mapping[int, int]] = ..., address: _Optional[str] = ..., direction: _Optional[str] = ..., protocols: _Optional[_Iterable[str]] = ..., total_sessions: _Optional[int] = ..., completed_sessions: _Optional[int] = ...) -> None: ...

class ConnectionInfo(_message.Message):
    __slots__ = ("connections", "inbound_connections", "outbound_connections")
    CONNECTIONS_FIELD_NUMBER: _ClassVar[int]
    INBOUND_CONNECTIONS_FIELD_NUMBER: _ClassVar[int]
    OUTBOUND_CONNECTIONS_FIELD_NUMBER: _ClassVar[int]
    connections: int
    inbound_connections: int
    outbound_connections: int
    def __init__(self, connections: _Optional[int] = ..., inbound_connections: _Optional[int] = ..., outbound_connections: _Optional[int] = ...) -> None: ...
