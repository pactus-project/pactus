package manager

import (
	"errors"
)

var (
	// ErrWalletAlreadyExists indicates a wallet already exists on disk.
	ErrWalletAlreadyExists = errors.New("wallet already exists")

	// ErrWalletAlreadyLoaded indicates a wallet is already loaded in memory.
	ErrWalletAlreadyLoaded = errors.New("wallet already loaded")

	// ErrWalletNotLoaded indicates a wallet is not loaded in memory.
	ErrWalletNotLoaded = errors.New("wallet is not loaded")
)
