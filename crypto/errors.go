package crypto

import (
	"fmt"
)

// InvalidLengthError is returned when the length of the data
// does not match the expected length.
type InvalidLengthError struct {
	Expected int
	Got      int
}

func (e InvalidLengthError) Error() string {
	return fmt.Sprintf("invalid length: expected %d, got %d", e.Expected, e.Got)
}

// InvalidHRPError is returned when the provided HRP code
// does not match the expected value.
type InvalidHRPError struct {
	Expected string
	Got      string
}

func (e InvalidHRPError) Error() string {
	return fmt.Sprintf("the HRP code is invalid: expected %s, got %s", e.Expected, e.Got)
}

// InvalidAddressTypeError is returned when the address type is not recognized or supported.
type InvalidAddressTypeError struct {
	Type AddressType
}

func (e InvalidAddressTypeError) Error() string {
	return fmt.Sprintf("invalid address type: got %d", e.Type)
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
