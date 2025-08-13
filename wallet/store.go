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
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/pactus-project/pactus/wallet/addresspath"
	"github.com/pactus-project/pactus/wallet/vault"
)

const (
	Version1 = 1 // Initial version
	Version2 = 2 // Supporting Ed25519
	Version3 = 3 // Supporting AEC-256-CBC encryption method
	Version4 = 4 // Set Default Fee for the Wallet

	VersionLatest = Version4
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

func (s *Store) Clone() *Store {
	clonedVault := *s.Vault // Assuming Vault has proper pointer handling internally
	clonedHistory := s.History

	return &Store{
		Version:   s.Version,
		UUID:      s.UUID,
		CreatedAt: s.CreatedAt,
		Network:   s.Network,
		VaultCRC:  s.VaultCRC,
		Vault:     &clonedVault,
		History:   clonedHistory,
	}
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
	if !s.Network.IsMainnet() {
		crypto.ToTestnetHRP()
	}

	oldVersion := s.Version
	switch oldVersion {
	case Version1:
		if err := s.setPublicKeys(); err != nil {
			return err
		}
		s.Version = Version2

		logger.Info(fmt.Sprintf("wallet upgraded from version %d to version %d",
			Version1, Version2))

		fallthrough

	case Version2:
		if s.Vault.IsEncrypted() {
			s.Vault.Encrypter.Params.SetUint32("keylen", 32)
		}
		s.Version = Version3

		logger.Info(fmt.Sprintf("wallet upgraded from version %d to version %d",
			Version2, Version3))
		fallthrough

	case Version3:
		s.Vault.DefaultFee = amount.Amount(10_000_000) // Set default fee to 0.01 PAC
		s.Version = Version4

		logger.Info(fmt.Sprintf("wallet upgraded from version %d to version %d",
			Version3, Version4))

	case Version4:
		return nil

	default:
		return UnsupportedVersionError{
			WalletVersion:    s.Version,
			SupportedVersion: VersionLatest,
		}
	}

	// Write wallet data.
	s.VaultCRC = s.calcVaultCRC()

	bs, err := s.ToBytes()
	if err != nil {
		return err
	}

	err = util.WriteFile(walletPath, bs)
	if err != nil {
		return err
	}

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
