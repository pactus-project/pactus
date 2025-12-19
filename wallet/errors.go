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

	// ErrInvalidAddressType describes an error in which the address type is invalid.
	ErrInvalidAddressType = errors.New("invalid address type")
)

// CRCNotMatchError describes an error in which the wallet CRC is not matched.
type CRCNotMatchError struct {
	Expected uint32
	Got      uint32
}

func (e CRCNotMatchError) Error() string {
	return fmt.Sprintf("crc not matched, expected: %d, got: %d", e.Expected, e.Got)
}

// ExitsError describes an error in which a wallet exists in the
// given path.
type ExitsError struct {
	Path string
}

func (e ExitsError) Error() string {
	return fmt.Sprintf("a wallet exists at: %s", e.Path)
}

// UnsupportedVersionError indicates the wallet version is incompatible with the software's supported version.
type UnsupportedVersionError struct {
	WalletVersion    int
	SupportedVersion int
}

func (e UnsupportedVersionError) Error() string {
	return fmt.Sprintf("wallet version %d is not supported, latest supported version is %d",
		e.WalletVersion, e.SupportedVersion)
}
