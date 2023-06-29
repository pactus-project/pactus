package wallet

import (
	"encoding/json"
	"hash/crc32"
	"time"

	"github.com/google/uuid"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/wallet/vault"
)

type store struct {
	CreatedAt time.Time         `json:"created_at"`
	History   history           `json:"history"`
	Vault     *vault.Vault      `json:"vault"`
	Version   int               `json:"version"`
	VaultCRC  uint32            `json:"crc"`
	UUID      uuid.UUID         `json:"uuid"`
	Network   genesis.ChainType `json:"network"`
}

func (s *store) Load() ([]byte, error) {
	s.VaultCRC = s.calcVaultCRC()
	return json.MarshalIndent(s, "  ", "  ")
}

func (s *store) Save(bs []byte) error {
	err := json.Unmarshal(bs, s)
	if err != nil {
		return err
	}

	if s.VaultCRC != s.calcVaultCRC() {
		return ErrInvalidCRC
	}

	return nil
}

func (s *store) calcVaultCRC() uint32 {
	d, _ := json.Marshal(s.Vault)
	return crc32.ChecksumIEEE(d)
}
