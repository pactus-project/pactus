package wallet

import (
	"encoding/json"
	"hash/crc32"
	"time"

	"github.com/pactus-project/pactus/crypto"

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
	for addrKey, info := range s.Vault.Addresses {
		if info.PublicKey == "" {
			addr, err := crypto.AddressFromString(info.Address)
			if err != nil {
				return nil
			}

			var xPub string
			if addr.IsAccountAddress() {
				xPub = s.Vault.Purposes.PurposeBLS.XPubAccount
			} else if addr.IsValidatorAddress() {
				xPub = s.Vault.Purposes.PurposeBLS.XPubValidator
			}

			pubKey, err := blshdkeychain.NewKeyFromString(xPub)
			if err != nil {
				return err
			}
			info.PublicKey = pubKey.String()
			s.Vault.Addresses[addrKey] = info
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
