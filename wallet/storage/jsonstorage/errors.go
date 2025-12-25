package jsonstorage

import (
	"errors"
	"fmt"
)

var ErrUnsupported = errors.New("operation not supported")

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
