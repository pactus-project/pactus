package wallet

import (
	"encoding/json"
	"hash/crc32"
	"time"

	"github.com/google/uuid"
	blshdkeychain "github.com/pactus-project/pactus/crypto/bls/hdkeychain"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/wallet/vault"
)

type store struct {
	Version   int               `json:"version"`
	UUID      uuid.UUID         `json:"uuid"`
	CreatedAt time.Time         `json:"created_at"`
	Network   genesis.ChainType `json:"network"`
	VaultCRC  uint32            `json:"crc"`
	Vault     *vault.Vault      `json:"vault"`
	History   history           `json:"history"`
}

func (s *store) ToBytes() ([]byte, error) {
	s.VaultCRC = s.calcVaultCRC()

	return json.MarshalIndent(s, "  ", "  ")
}

func (s *store) ValidateCRC() error {
	crc := s.calcVaultCRC()
	if s.VaultCRC != crc {
		return CRCNotMatchError{
			Expected: crc,
			Got:      s.VaultCRC,
		}
	}

	return nil
}

func (s *store) UpgradeWallet(walletPath string) error {
	if s.Version != Version2 {
		if err := s.setPublicKeys(); err != nil {
			return err
		}

		s.VaultCRC = s.calcVaultCRC()
		s.Version = Version2
		bs, err := s.ToBytes()
		if err != nil {
			return err
		}

		return util.WriteFile(walletPath, bs)
	}

	return nil
}

func (s *store) setPublicKeys() error {
	for addrKey, addrInfo := range s.Vault.Addresses {
		if addrInfo.PublicKey == "" {
			pubKey, err := blshdkeychain.NewKeyFromString(s.Vault.Purposes.PurposeBLS.XPubAccount)
			if err != nil {
				return err
			}
			addrInfo.PublicKey = pubKey.String()
			s.Vault.Addresses[addrKey] = addrInfo
		}
	}

	return nil
}

func (s *store) calcVaultCRC() uint32 {
	d, err := json.Marshal(s.Vault)
	if err != nil {
		return 0
	}

	return crc32.ChecksumIEEE(d)
}
