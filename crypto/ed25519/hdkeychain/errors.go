package hdkeychain

import (
	"errors"
	"fmt"
)

var (
	// ErrInvalidSeedLen describes an error in which the provided seed or
	// seed length is not in the allowed range.
	ErrInvalidSeedLen = fmt.Errorf("seed length must be between %d and %d "+
		"bits", MinSeedBytes*8, MaxSeedBytes*8)

	// ErrInvalidKeyData describes an error in which the provided key is
	// not valid.
	ErrInvalidKeyData = errors.New("key data is invalid")

	// ErrInvalidHRP describes an error in which the HRP is not valid.
	ErrInvalidHRP = errors.New("HRP is invalid")

	// ErrNonHardenedPath is returned when a non-hardened derivation path is used,
	// which is not supported by ed25519.
	ErrNonHardenedPath = errors.New("non-hardened derivation not supported")
)
