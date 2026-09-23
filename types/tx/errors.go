package tx

import (
	"errors"

	"github.com/pactus-project/pactus/types/tx/payload"
)

// ErrInvalidSigner is returned when the signer address is not valid.
var ErrInvalidSigner = errors.New("invalid signer address")

// ErrTrailingBytes is returned when an anchor transaction has leftover bytes.
var ErrTrailingBytes = errors.New("trailing bytes")

// BasicCheckError is returned when the basic check on the transaction fails.
type BasicCheckError struct {
	Reason string
}

func (e BasicCheckError) Error() string {
	return e.Reason
}

// InvalidPayloadTypeError is returned when the payload type is not valid.
type InvalidPayloadTypeError struct {
	PayloadType payload.Type
}

func (e InvalidPayloadTypeError) Error() string {
	return "invalid payload type: " + e.PayloadType.String()
}
