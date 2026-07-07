package jsonstorage

import (
	"encoding/json"
	"fmt"

	"github.com/pactus-project/gopkg/logger"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	blshdkeychain "github.com/pactus-project/pactus/crypto/bls/hdkeychain"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/wallet/addresspath"
	"github.com/pactus-project/pactus/wallet/encrypter"
	"github.com/pactus-project/pactus/wallet/types"
	"github.com/pactus-project/pactus/wallet/vault"
)

type legacyVault struct {
	Encrypter  encrypter.Encrypter           `json:"encrypter"`
	Purposes   vault.Purposes                `json:"purposes"` // Contains Purposes of the vault
	DefaultFee amount.Amount                 `json:"default_fee"`
	Addresses  map[string]*types.AddressInfo `json:"addresses"`
}

type legacyStore struct {
	Vault legacyVault `json:"vault"`
}

func upgrade(path string) error {
	data, err := util.ReadFile(path)
	if err != nil {
		return err
	}

	store := new(store)
	err = json.Unmarshal(data, store)
	if err != nil {
		return err
	}

	legacyStore := new(legacyStore)
	err = json.Unmarshal(data, &legacyStore)
	if err != nil {
		return err
	}

	if !store.Network.IsMainnet() {
		crypto.ToTestnetHRP()
	}

	switch store.Version {
	case Version1:
		if err := setPublicKeys(legacyStore); err != nil {
			return err
		}

		logger.Info(fmt.Sprintf("wallet upgraded from version %d to version %d", Version1, Version2))

		fallthrough

	case Version2:
		if legacyStore.Vault.Encrypter.IsEncrypted() {
			store.Vault.Encrypter.Params.SetUint32("keylen", 32)
		}

		logger.Info(fmt.Sprintf("wallet upgraded from version %d to version %d", Version2, Version3))

		fallthrough

	case Version3:
		legacyStore.Vault.DefaultFee = amount.Amount(10_000_000) // Set default fee to 0.01 PAC

		logger.Info(fmt.Sprintf("wallet upgraded from version %d to version %d", Version3, Version4))

		fallthrough

	case Version4:
		store.DefaultFee = legacyStore.Vault.DefaultFee
		store.Addresses = make(map[string]*types.AddressInfo)
		store.VaultCRC = store.calcVaultCRC()

		for addr, ai := range legacyStore.Vault.Addresses {
			store.Addresses[addr] = &types.AddressInfo{
				Address:   ai.Address,
				PublicKey: ai.PublicKey,
				Label:     ai.Label,
				Path:      ai.Path,
			}
		}

		logger.Info(fmt.Sprintf("wallet upgraded from version %d to version %d", Version4, Version5))

		fallthrough

	case Version5:
		store.Version = Version6
		store.Vault.Purposes.PurposeBIP44.NextSecp256k1Index = 0

		logger.Info(fmt.Sprintf("wallet upgraded from version %d to version %d", Version5, Version6))

		return store.Save(path)

	case Version6:
		// Latest version, no need to upgrade.
		return nil

	default:
		return UnsupportedVersionError{
			WalletVersion:    store.Version,
			SupportedVersion: VersionLatest,
		}
	}
}

func setPublicKeys(store *legacyStore) error {
	for addrKey, info := range store.Vault.Addresses {
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
			xPub = store.Vault.Purposes.PurposeBLS.XPubAccount
		} else if addr.IsValidatorAddress() {
			xPub = store.Vault.Purposes.PurposeBLS.XPubValidator
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
		store.Vault.Addresses[addrKey] = info
	}

	return nil
}
