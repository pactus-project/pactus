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

	// ErrHistoryExists describes an error in which the transaction already exists
	// in history.
	ErrHistoryExists = errors.New("transaction already exists")
)

// CRCNotMatchError describes an error in which the wallet CRC is not macthed.
type CRCNotMatchError struct {
	Expected uint32
	Got      uint32
}

func (e CRCNotMatchError) Error() string {
	return fmt.Sprintf("crc not matched, expected: %d, got: %d", e.Expected, e.Got)
}

// WalletExitsError describes an error in which a wallet exists in the
// given path.
type WalletExitsError struct { //nolint
	Path string
}

func (e WalletExitsError) Error() string {
	return fmt.Sprintf("a wallet exists at: %s", e.Path)
}
