package jsonstorage

import (
	"encoding/json"
	"fmt"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	blshdkeychain "github.com/pactus-project/pactus/crypto/bls/hdkeychain"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/pactus-project/pactus/wallet/addresspath"
	"github.com/pactus-project/pactus/wallet/encrypter"
	"github.com/pactus-project/pactus/wallet/types"
	"github.com/pactus-project/pactus/wallet/vault"
)

type upgrader struct {
	path string
	data []byte
}

func Upgrade(path string) error {
	data, err := util.ReadFile(path)
	if err != nil {
		return err
	}

	u := upgrader{
		path: path,
		data: data,
	}

	return u.upgrade()
}

type legacyVault struct {
	Encrypter  encrypter.Encrypter          `json:"encrypter"`
	Purposes   vault.Purposes               `json:"purposes"` // Contains Purposes of the vault
	DefaultFee amount.Amount                `json:"default_fee"`
	Addresses  map[string]types.AddressInfo `json:"addresses"`
}

type legacyStore struct {
	Vault legacyVault `json:"vault"`
}

func (u *upgrader) upgrade() error {
	store := new(store)
	err := json.Unmarshal(u.data, store)
	if err != nil {
		return err
	}

	legacyStore := new(legacyStore)
	err = json.Unmarshal(u.data, &legacyStore)
	if err != nil {
		return err
	}

	if !store.Network.IsMainnet() {
		crypto.ToTestnetHRP()
	}

	switch store.Version {
	case Version1:
		if err := u.setPublicKeys(legacyStore); err != nil {
			return err
		}

		logger.Info(fmt.Sprintf("wallet upgraded from version %d to version %d",
			Version1, Version2))

		fallthrough

	case Version2:
		if legacyStore.Vault.Encrypter.IsEncrypted() {
			store.Vault.Encrypter.Params.SetUint32("keylen", 32)
		}

		logger.Info(fmt.Sprintf("wallet upgraded from version %d to version %d",
			Version2, Version3))

		fallthrough

	case Version3:
		legacyStore.Vault.DefaultFee = amount.Amount(10_000_000) // Set default fee to 0.01 PAC

		logger.Info(fmt.Sprintf("wallet upgraded from version %d to version %d",
			Version3, Version4))

		fallthrough

	case Version4:
		store.DefaultFee = legacyStore.Vault.DefaultFee
		store.Addresses = make(map[string]types.AddressInfo)
		store.Version = Version5
		store.VaultCRC = store.calcVaultCRC()

		for addr, ai := range legacyStore.Vault.Addresses {
			store.Addresses[addr] = types.AddressInfo{
				Address:   ai.Address,
				PublicKey: ai.PublicKey,
				Label:     ai.Label,
				Path:      ai.Path,
			}
		}

		return store.Save(u.path)

	case Version5:
		// Latest version, no need to upgrade.
		return nil

	default:
		return UnsupportedVersionError{
			WalletVersion:    store.Version,
			SupportedVersion: VersionLatest,
		}
	}
}

func (*upgrader) setPublicKeys(store *legacyStore) error {
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
