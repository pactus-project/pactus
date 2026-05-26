package sqlitestorage

import (
	"errors"
	"fmt"
)

var ErrUnsupported = errors.New("operation not supported")

// UnsupportedVersionError indicates the wallet version is incompatible with the software's supported version.
type UnsupportedVersionError struct {
	WalletVersion    int
	SupportedVersion int
}

func (e UnsupportedVersionError) Error() string {
	return fmt.Sprintf("wallet version %d is not supported, latest supported version is %d",
		e.WalletVersion, e.SupportedVersion)
}
