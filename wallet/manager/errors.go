package manager

import (
	"errors"
)

var (
	// ErrWalletAlreadyExists indicates a wallet already exists on disk.
	ErrWalletAlreadyExists = errors.New("wallet already exists")

	// TODO: rename me
	// ErrWalletNotLoaded indicates a wallet is not loaded in memory.
	ErrWalletNotLoaded = errors.New("wallet is not loaded")
)

// ConfigError is returned when the config is not valid with a descriptive Reason message.
type ConfigError struct {
	Reason string
}

func (e ConfigError) Error() string {
	return e.Reason
}
