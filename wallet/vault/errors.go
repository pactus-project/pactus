package vault

import (
	"errors"
)

var (
	// ErrInvalidPath describes an error in which the key path is invalid.
	ErrInvalidPath = errors.New("the key path is invalid")

	// ErrNeutered describes an error in which the wallet is neutered.
	ErrNeutered = errors.New("wallet is neutered")

	// ErrUnsupportedPurpose describes an error in which the purpose is not supported.
	ErrUnsupportedPurpose = errors.New("unsupported purpose")
)
