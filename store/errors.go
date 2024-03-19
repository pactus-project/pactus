package store

import (
	"fmt"

	"github.com/pactus-project/pactus/crypto"
)

// ConfigError is returned when the store configuration is invalid.
type ConfigError struct {
	Reason string
}

func (e ConfigError) Error() string {
	return e.Reason
}

// PublicKeyNotFoundError is returned when the public key associated with an address
// is not found in the store.
type PublicKeyNotFoundError struct {
	Address crypto.Address
}

func (e PublicKeyNotFoundError) Error() string {
	return fmt.Sprintf("public key not found for: %s",
		e.Address.String())
}
