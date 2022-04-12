package wallet

import (
	"errors"
	"fmt"
)

var (
	// ErrInvalidCRC describes an error in which the wallet CRC is
	// invalid
	ErrInvalidCRC = errors.New("invalid CRC")

	// ErrInvalidNetwork describes an error in which the network is not
	// valid
	ErrInvalidNetwork = errors.New("invalid network")

	// ErrInvalidPassword describes an error in which the password is
	// invalid
	ErrInvalidPassword = errors.New("invalid password")

	// ErrAddressExists describes an error in which the address already
	// exist in wallet
	ErrAddressExists = errors.New("address already exists")
)

// ErrWalletExits describes an error in which a wallet exists in the
// given path
type ErrWalletExits struct {
	Path string
}

func NewErrWalletExits(path string) *ErrWalletExits {
	return &ErrWalletExits{Path: path}
}

func (e *ErrWalletExits) Error() string {
	return fmt.Sprintf("a wallet exists in: %s", e.Path)
}

// ErrAddressNotFound describes an error in which the address doesn't
// exist in wallet
type ErrAddressNotFound struct {
	addr string
}

func NewErrAddressNotFound(addr string) *ErrAddressNotFound {
	return &ErrAddressNotFound{addr: addr}
}

func (e *ErrAddressNotFound) Error() string {
	return fmt.Sprintf("address not found: %s", e.addr)
}
