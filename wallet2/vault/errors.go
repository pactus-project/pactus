package vault

import (
	"errors"
	"fmt"
)

var (
	// ErrAddressExists describes an error in which the address already exist
	// in wallet.
	ErrAddressExists = errors.New("address already exists")

	// ErrInvalidPath describes an error in which the key path is invalid.
	ErrInvalidPath = errors.New("the key path is invalid")

	// ErrNeutered describes an error in which the wallet is neutered.
	ErrNeutered = errors.New("wallet is neutered")

	// ErrInvalidCoinType describes an error in which the coin type is not valid.
	ErrInvalidCoinType = errors.New("invalid coin type")

	// ErrUnsupportedPurpose describes an error in which the purpose is not supported.
	ErrUnsupportedPurpose = errors.New("unsupported purpose")
)

// AddressNotFoundError describes an error in which the address doesn't exist
// in wallet.
type AddressNotFoundError struct {
	addr string
}

func NewErrAddressNotFound(addr string) error {
	return AddressNotFoundError{addr: addr}
}

func (e AddressNotFoundError) Error() string {
	return fmt.Sprintf("address not found: %s", e.addr)
}
