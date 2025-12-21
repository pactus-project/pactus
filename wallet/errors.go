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
)

// ExitsError describes an error in which a wallet exists in the
// given path.
type ExitsError struct {
	Path string
}

func (e ExitsError) Error() string {
	return fmt.Sprintf("a wallet exists at: %s", e.Path)
}
