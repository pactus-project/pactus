package crypto

import (
	"errors"
	"fmt"
)

// ErrInvalidSignature is returned when a signature is invalid.
var ErrInvalidSignature = errors.New("invalid signature")

// InvalidLengthError is returned when the length of the data
// does not match the expected length.
type InvalidLengthError int

func (e InvalidLengthError) Error() string {
	return fmt.Sprintf("invalid length: %d", int(e))
}

// InvalidHRPError is returned when the provided HRP code
// does not match the expected value.
type InvalidHRPError string

func (e InvalidHRPError) Error() string {
	return fmt.Sprintf("invalid HRP: %s", string(e))
}

// InvalidAddressTypeError is returned when the address type is not recognized or supported.
type InvalidAddressTypeError int

func (e InvalidAddressTypeError) Error() string {
	return fmt.Sprintf("invalid address type: %d", int(e))
}

// AddressMismatchError is returned when the provided address is not derived
// from the corresponding public key.
type AddressMismatchError struct {
	Expected Address
	Got      Address
}

func (e AddressMismatchError) Error() string {
	return fmt.Sprintf("address mismatch: expected %s, got %s",
		e.Expected.String(), e.Got.String())
}
