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
	Version   int               `json:"version"`
	UUID      uuid.UUID         `json:"uuid"`
	CreatedAt time.Time         `json:"created_at"`
	Network   genesis.ChainType `json:"network"`
	VaultCRC  uint32            `json:"crc"`
	Vault     *vault.Vault      `json:"vault"`
	History   history           `json:"history"`
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
