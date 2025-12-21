package jsonstorage

import (
	"fmt"
	"time"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	blshdkeychain "github.com/pactus-project/pactus/crypto/bls/hdkeychain"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/pactus-project/pactus/wallet/addresspath"
	"github.com/pactus-project/pactus/wallet/types"
	"github.com/pactus-project/pactus/wallet/version"
)

type upgrader struct {
	store *Store
}

func NewUpgrader(store *Store) *upgrader {
	return &upgrader{store: store}
}

func (u *upgrader) Upgrade() error {
	if !u.store.Network.IsMainnet() {
		crypto.ToTestnetHRP()
	}

	oldVersion := u.store.Version
	switch oldVersion {
	case version.Version1:
		if err := u.setPublicKeys(); err != nil {
			return err
		}
		u.store.Version = version.Version2

		logger.Info(fmt.Sprintf("wallet upgraded from version %d to version %d",
			version.Version1, version.Version2))

		fallthrough

	case version.Version2:
		if u.store.Vault.IsEncrypted() {
			u.store.Vault.Encrypter.Params.SetUint32("keylen", 32)
		}
		u.store.Version = version.Version3

		logger.Info(fmt.Sprintf("wallet upgraded from version %d to version %d",
			version.Version2, version.Version3))

		fallthrough

	case version.Version3:
		u.store.Vault.DefaultFeeDeprecated = amount.Amount(10_000_000) // Set default fee to 0.01 PAC
		u.store.Version = version.Version4

		logger.Info(fmt.Sprintf("wallet upgraded from version %d to version %d",
			version.Version3, version.Version4))

	case version.Version4:
		u.store.DefaultFee = u.store.Vault.DefaultFeeDeprecated

		now := time.Now()
		for addr, ai := range u.store.Vault.AddressesDeprecated {
			u.store.Addresses[addr] = types.AddressInfo{
				Address:   ai.Address,
				PublicKey: ai.PublicKey,
				Label:     ai.Label,
				Path:      ai.Path,
				CreatedAt: now,
				UpdatedAt: now,
			}
		}

		u.store.Vault.DefaultFeeDeprecated = 0
		u.store.Vault.AddressesDeprecated = nil

		return nil

	default:
		return UnsupportedVersionError{
			WalletVersion:    u.store.Version,
			SupportedVersion: version.VersionLatest,
		}
	}

	return nil
}

func (u *upgrader) setPublicKeys() error {
	for addrKey, info := range u.store.Vault.AddressesDeprecated {
		if info.PublicKey != "" {
			continue
		}

		// Some old wallet doesn't have public key for all addresses.
		addr, err := crypto.AddressFromString(info.Address)
		if err != nil {
			return err
		}

		var xPub string
		if addr.IsAccountAddress() {
			xPub = u.store.Vault.Purposes.PurposeBLS.XPubAccount
		} else if addr.IsValidatorAddress() {
			xPub = u.store.Vault.Purposes.PurposeBLS.XPubValidator
		}

		ext, err := blshdkeychain.NewKeyFromString(xPub)
		if err != nil {
			return err
		}

		p, err := addresspath.FromString(info.Path)
		if err != nil {
			return err
		}

		extendedKey, err := ext.Derive(p.AddressIndex())
		if err != nil {
			return err
		}

		blsPubKey, err := bls.PublicKeyFromBytes(extendedKey.RawPublicKey())
		if err != nil {
			return err
		}

		info.PublicKey = blsPubKey.String()
		u.store.Vault.AddressesDeprecated[addrKey] = info
	}

	return nil
}
