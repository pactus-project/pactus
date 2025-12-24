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

type store struct {
	Version    int                          `json:"version"`
	UUID       uuid.UUID                    `json:"uuid"`
	CreatedAt  time.Time                    `json:"created_at"`
	Network    genesis.ChainType            `json:"network"`
	VaultCRC   uint32                       `json:"crc"`
	DefaultFee amount.Amount                `json:"default_fee"`
	Vault      vault.Vault                  `json:"vault"`
	Addresses  map[string]types.AddressInfo `json:"addresses"`
	History    history                      `json:"history"`
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

func (s *store) calcVaultCRC() uint32 {
	d, err := json.Marshal(s.Vault)
	if err != nil {
		return 0
	}

	return crc32.ChecksumIEEE(d)
}
