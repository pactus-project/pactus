package hdkeychain

import (
	"errors"
	"fmt"
)

var (
	// ErrDeriveHardFromPublic describes an error in which the caller
	// attempted to derive a hardened extended key from a public key.
	ErrDeriveHardFromPublic = errors.New("cannot derive a hardened key " +
		"from a public key")

	// ErrNotPrivExtKey describes an error in which the caller attempted
	// to extract a private key from a public extended key.
	ErrNotPrivExtKey = errors.New("unable to create private keys from a " +
		"public extended key")

	// ErrUnusableSeed describes an error in which the provided seed is not
	// usable due to the derived key falling outside of the valid range for
	// BLS private keys.  This error indicates the caller must choose
	// another seed.
	ErrUnusableSeed = errors.New("unusable seed")

	// ErrInvalidSeedLen describes an error in which the provided seed or
	// seed length is not in the allowed range.
	ErrInvalidSeedLen = fmt.Errorf("seed length must be between %d and %d "+
		"bits", MinSeedBytes*8, MaxSeedBytes*8)

	// ErrInvalidKeyData describes an error in which the provided key is
	// not valid.
	ErrInvalidKeyData = errors.New("key data is invalid")

	// ErrInvalidHRP describes an error in which the HRP is not valid.
	ErrInvalidHRP = errors.New("HRP is invalid")
)
