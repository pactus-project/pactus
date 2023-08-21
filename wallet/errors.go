package wallet

import (
	"errors"
	"fmt"
)

var (
	// ErrInvalidCRC describes an error in which the wallet CRC is invalid.
	ErrInvalidCRC = errors.New("invalid CRC")

	// ErrInvalidNetwork describes an error in which the network is invalid.
	ErrInvalidNetwork = errors.New("invalid network")

	// ErrOffline describes an error in which the wallet is offline.
	ErrOffline = errors.New("wallet is in offline mode")

	// ErrHistoryExists describes an error in which the transaction already exists
	// in history.
	ErrHistoryExists = errors.New("transaction already exists")
)

// WalletExitsError describes an error in which a wallet exists in the
// given path.
type WalletExitsError struct { //nolint
	Path string
}

func NewWalletExitsError(path string) error {
	return WalletExitsError{Path: path}
}

func (e WalletExitsError) Error() string {
	return fmt.Sprintf("a wallet exists at: %s", e.Path)
}
