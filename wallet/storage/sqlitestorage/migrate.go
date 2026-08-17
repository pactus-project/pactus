package sqlitestorage

import (
	"context"

	"github.com/pactus-project/pactus/wallet/storage/jsonstorage"
)

// Migrate converts a legacy JSON wallet to the SQLite format.
func Migrate(ctx context.Context, path string, jsonStrg *jsonstorage.Storage) error {
	info := jsonStrg.WalletInfo()

	strg, err := Create(ctx, path, info.Network, jsonStrg.Vault())
	if err != nil {
		return err
	}

	defer func() {
		_ = strg.Close()
	}()

	// Preserve the wallet metadata from the legacy wallet.
	if err := strg.SetDefaultFee(info.DefaultFee); err != nil {
		return err
	}

	// Copy all addresses from the legacy wallet.
	for _, addrInfo := range jsonStrg.AllAddresses() {
		if err := strg.InsertAddress(addrInfo); err != nil {
			return err
		}
	}

	return nil
}
