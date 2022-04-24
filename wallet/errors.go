package wallet

import (
	"errors"
	"fmt"
)

var (
	// ErrInvalidCRC describes an error in which the wallet CRC is invalid
	ErrInvalidCRC = errors.New("invalid CRC")

	// ErrInvalidNetwork describes an error in which the network is invalid
	ErrInvalidNetwork = errors.New("invalid network")
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
	return fmt.Sprintf("a wallet exists at: %s", e.Path)
}
