package vault

import (
	"errors"
)

var (
	// ErrInvalidPath describes an error in which the key path is invalid.
	ErrInvalidPath = errors.New("the key path is invalid")

	// ErrInvalidPublicKey describes an error in which the public key is invalid.
	ErrInvalidPublicKey = errors.New("invalid public key")

	// ErrInvalidPrivateKey describes an error in which the private key is invalid.
	ErrInvalidPrivateKey = errors.New("invalid private key")
	// ErrNeutered describes an error in which the wallet is neutered.
	ErrNeutered = errors.New("wallet is neutered")

	// ErrUnsupportedPurpose describes an error in which the purpose is not supported.
	ErrUnsupportedPurpose = errors.New("unsupported purpose")
)
