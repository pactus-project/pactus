package wallet

import (
	"errors"
	"fmt"
)

var (
	// ErrInvalidNetwork describes an error in which the network is invalid.
	ErrInvalidNetwork = errors.New("invalid network")

	// ErrOffline describes an error in which the wallet is offline.
	ErrOffline = errors.New("wallet is in offline mode")

	// ErrInvalidAddressType describes an error in which the address type is invalid.
	ErrInvalidAddressType = errors.New("invalid address type")

	// ErrAddressExists describes an error in which the address already exist
	// in wallet.
	ErrAddressExists = errors.New("address already exists")
)

// ExitsError describes an error in which a wallet exists in the
// given path.
type ExitsError struct {
	Path string
}

func (e ExitsError) Error() string {
	return fmt.Sprintf("a wallet exists at: %s", e.Path)
}

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
