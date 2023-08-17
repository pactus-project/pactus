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
)

// ErrAddressNotFound describes an error in which the address doesn't exist
// in wallet.
type ErrAddressNotFound struct {
	addr string
}

func NewErrAddressNotFound(addr string) error {
	return ErrAddressNotFound{addr: addr}
}

func (e ErrAddressNotFound) Error() string {
	return fmt.Sprintf("address not found: %s", e.addr)
}
