from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
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
