package wallet

import (
	"encoding/json"
	"fmt"
	"hash/crc32"
	"time"

	"github.com/google/uuid"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	blshdkeychain "github.com/pactus-project/pactus/crypto/bls/hdkeychain"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/pactus-project/pactus/wallet/addresspath"
	"github.com/pactus-project/pactus/wallet/vault"
)

const (
	Version1 = 1 // initial version
	Version2 = 2 // supporting Ed25519

	VersionLatest = Version2
)

type Store struct {
	Version   int               `json:"version"`
	UUID      uuid.UUID         `json:"uuid"`
	CreatedAt time.Time         `json:"created_at"`
	Network   genesis.ChainType `json:"network"`
	VaultCRC  uint32            `json:"crc"`
	Vault     *vault.Vault      `json:"vault"`
	History   history           `json:"history"`
}

func FromBytes(data []byte) (*Store, error) {
	s := new(Store)
	if err := json.Unmarshal(data, s); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Store) ToBytes() ([]byte, error) {
	s.VaultCRC = s.calcVaultCRC()

	return json.MarshalIndent(s, "  ", "  ")
}

func (s *Store) ValidateCRC() error {
	crc := s.calcVaultCRC()
	if s.VaultCRC != crc {
		return CRCNotMatchError{
			Expected: crc,
			Got:      s.VaultCRC,
		}
	}

	return nil
}

func (s *Store) UpgradeWallet(walletPath string) error {
	oldVersion := s.Version
	switch oldVersion {
	case Version1:
		if err := s.setPublicKeys(); err != nil {
			return err
		}

	case Version2:
		// Current version
		return nil

	default:
		return UnsupportedVersionError{
			WalletVersion:    s.Version,
			SupportedVersion: VersionLatest,
		}
	}

	s.VaultCRC = s.calcVaultCRC()
	s.Version = Version2

	bs, err := s.ToBytes()
	if err != nil {
		return err
	}

	err = util.WriteFile(walletPath, bs)
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("wallet upgraded from version %d to version %d",
		oldVersion, VersionLatest))

	return nil
}

func (s *Store) setPublicKeys() error {
	for addrKey, info := range s.Vault.Addresses {
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
			xPub = s.Vault.Purposes.PurposeBLS.XPubAccount
		} else if addr.IsValidatorAddress() {
			xPub = s.Vault.Purposes.PurposeBLS.XPubValidator
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
		s.Vault.Addresses[addrKey] = info
	}

	return nil
}

func (s *Store) calcVaultCRC() uint32 {
	d, err := json.Marshal(s.Vault)
	if err != nil {
		return 0
	}

	return crc32.ChecksumIEEE(d)
}
