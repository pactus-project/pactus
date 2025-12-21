package jsonstorage

import (
	"errors"
	"fmt"
)

var (
	// ErrHistoryExists describes an error in which the transaction already exists
	// in history.
	ErrHistoryExists = errors.New("transaction already exists")

	// ErrInvalidAddressType describes an error in which the address type is invalid.
	ErrInvalidAddressType = errors.New("invalid address type")

	// ErrAddressExists describes an error in which the address already exist
	// in wallet.
	ErrAddressExists = errors.New("address already exists")
)

// CRCNotMatchError describes an error in which the wallet CRC is not matched.
type CRCNotMatchError struct {
	Expected uint32
	Got      uint32
}

func (e CRCNotMatchError) Error() string {
	return fmt.Sprintf("crc not matched, expected: %d, got: %d", e.Expected, e.Got)
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
