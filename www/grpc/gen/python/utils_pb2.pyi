from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class SignMessageWithPrivateKeyRequest(_message.Message):
    __slots__ = ("private_key", "message")
    PRIVATE_KEY_FIELD_NUMBER: _ClassVar[int]
    MESSAGE_FIELD_NUMBER: _ClassVar[int]
    private_key: str
    message: str
    def __init__(self, private_key: _Optional[str] = ..., message: _Optional[str] = ...) -> None: ...

class SignMessageWithPrivateKeyResponse(_message.Message):
    __slots__ = ("signature",)
    SIGNATURE_FIELD_NUMBER: _ClassVar[int]
    signature: str
    def __init__(self, signature: _Optional[str] = ...) -> None: ...

class VerifyMessageRequest(_message.Message):
    __slots__ = ("message", "signature", "public_key")
    MESSAGE_FIELD_NUMBER: _ClassVar[int]
    SIGNATURE_FIELD_NUMBER: _ClassVar[int]
    PUBLIC_KEY_FIELD_NUMBER: _ClassVar[int]
    message: str
    signature: str
    public_key: str
    def __init__(self, message: _Optional[str] = ..., signature: _Optional[str] = ..., public_key: _Optional[str] = ...) -> None: ...

class VerifyMessageResponse(_message.Message):
    __slots__ = ("is_valid",)
    IS_VALID_FIELD_NUMBER: _ClassVar[int]
    is_valid: bool
    def __init__(self, is_valid: bool = ...) -> None: ...

class PublicKeyAggregationRequest(_message.Message):
    __slots__ = ("public_keys",)
    PUBLIC_KEYS_FIELD_NUMBER: _ClassVar[int]
    public_keys: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, public_keys: _Optional[_Iterable[str]] = ...) -> None: ...

class PublicKeyAggregationResponse(_message.Message):
    __slots__ = ("public_key", "address")
    PUBLIC_KEY_FIELD_NUMBER: _ClassVar[int]
    ADDRESS_FIELD_NUMBER: _ClassVar[int]
    public_key: str
    address: str
    def __init__(self, public_key: _Optional[str] = ..., address: _Optional[str] = ...) -> None: ...

class SignatureAggregationRequest(_message.Message):
    __slots__ = ("signatures",)
    SIGNATURES_FIELD_NUMBER: _ClassVar[int]
    signatures: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, signatures: _Optional[_Iterable[str]] = ...) -> None: ...

class SignatureAggregationResponse(_message.Message):
    __slots__ = ("signature",)
    SIGNATURE_FIELD_NUMBER: _ClassVar[int]
    signature: str
    def __init__(self, signature: _Optional[str] = ...) -> None: ...
