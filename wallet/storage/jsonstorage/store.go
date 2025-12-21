package jsonstorage

import (
	"encoding/json"
	"hash/crc32"
	"time"

	"github.com/google/uuid"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/wallet/types"
	"github.com/pactus-project/pactus/wallet/vault"
)

type Store struct {
	Version    int                          `json:"version"`
	UUID       uuid.UUID                    `json:"uuid"`
	CreatedAt  time.Time                    `json:"created_at"`
	Network    genesis.ChainType            `json:"network"`
	VaultCRC   uint32                       `json:"crc"`
	Vault      vault.Vault                  `json:"vault"`
	DefaultFee amount.Amount                `json:"default_fee"`
	Addresses  map[string]types.AddressInfo `json:"addresses"`
	History    history                      `json:"history"`
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

func (s *Store) Upgrade() error {
	u := NewUpgrader(s)
	if err := u.Upgrade(); err != nil {
		return err
	}

	// Write wallet data.
	s.VaultCRC = s.calcVaultCRC()

	return nil
}

func (s *Store) calcVaultCRC() uint32 {
	d, err := json.Marshal(s.Vault)
	if err != nil {
		return 0
	}

	return crc32.ChecksumIEEE(d)
}
