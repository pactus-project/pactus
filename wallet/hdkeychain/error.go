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

	// ErrDeriveBeyondMaxDepth describes an error in which the caller
	// has attempted to derive more than 255 keys from a root key.
	ErrDeriveBeyondMaxDepth = errors.New("cannot derive a key with more than " +
		"255 indices in its path")

	// ErrNotPrivExtKey describes an error in which the caller attempted
	// to extract a private key from a public extended key.
	ErrNotPrivExtKey = errors.New("unable to create private keys from a " +
		"public extended key")

	// ErrInvalidChild describes an error in which the child at a specific
	// index is invalid due to the derived key falling outside of the valid
	// range for BLS private keys.  This error indicates the caller
	// should simply ignore the invalid child extended key at this index and
	// increment to the next index.
	ErrInvalidChild = errors.New("the extended key at this index is invalid")

	// ErrUnusableSeed describes an error in which the provided seed is not
	// usable due to the derived key falling outside of the valid range for
	// BLS private keys.  This error indicates the caller must choose
	// another seed.
	ErrUnusableSeed = errors.New("unusable seed")

	// ErrInvalidSeedLen describes an error in which the provided seed or
	// seed length is not in the allowed range.
	ErrInvalidSeedLen = fmt.Errorf("seed length must be between %d and %d "+
		"bits", MinSeedBytes*8, MaxSeedBytes*8)

	//
	ErrInvalidKey = errors.New("key is invalid")
)
