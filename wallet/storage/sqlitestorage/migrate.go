package sqlitestorage

import (
	"context"
	"fmt"
	"os"

	"github.com/pactus-project/pactus/wallet/storage/jsonstorage"
)

// Migrate converts a legacy JSON wallet to the SQLite format and saves it at
// the given path (a directory). The wallet metadata (network, default fee,
// UUID and creation time), the vault and all addresses are preserved.
// The legacy JSON wallet file is left untouched; it is up to the caller to
// decide whether to remove it.
func Migrate(ctx context.Context, jsonPath, path string, opts ...Option) (*Storage, error) {
	jsonStrg, err := jsonstorage.Open(jsonPath)
	if err != nil {
		return nil, err
	}

	info := jsonStrg.WalletInfo()

	strg, err := Create(ctx, path, info.Network, jsonStrg.Vault(), opts...)
	if err != nil {
		return nil, err
	}

	// Remove the partially created storage on failure.
	cleanup := func() {
		_ = strg.Close()
		_ = os.RemoveAll(path)
	}

	// Preserve the wallet metadata from the legacy wallet.
	if err := strg.SetDefaultFee(info.DefaultFee); err != nil {
		cleanup()

		return nil, err
	}

	if err := strg.updateWalletEntry(keyUUID, info.UUID); err != nil {
		cleanup()

		return nil, err
	}
	strg.info.UUID = info.UUID

	if err := strg.updateWalletEntry(keyCreatedAt, fmt.Sprintf("%d", info.CreatedAt.Unix())); err != nil {
		cleanup()

		return nil, err
	}
	strg.info.CreatedAt = info.CreatedAt

	// Copy all addresses from the legacy wallet.
	for _, addrInfo := range jsonStrg.AllAddresses() {
		if err := strg.InsertAddress(addrInfo); err != nil {
			cleanup()

			return nil, err
		}
	}

	return strg, nil
}
