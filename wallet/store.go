package wallet

import (
	"encoding/json"
	"hash/crc32"
	"time"

	"github.com/google/uuid"
	"github.com/zarbchain/zarb-go/wallet/vault"
)

type store struct {
	Version   int          `json:"version"`
	UUID      uuid.UUID    `json:"uuid"`
	CreatedAt time.Time    `json:"created_at"`
	Network   Network    `json:"network"`
	VaultCRC  uint32       `json:"crc"`
	Vault     *vault.Vault `json:"vault"`
}

func (s *store) calcVaultCRC() uint32 {
	d, _ := json.Marshal(s.Vault)
	return crc32.ChecksumIEEE(d)
}

func (s *store) UpdatePassword(oldPassword, newPassword string) error {
	return s.Vault.UpdatePassword(oldPassword, newPassword)
}

func (s *store) IsEncrypted() bool {
	return s.Vault.IsEncrypted()
}

func (s *store) AddressInfos() []vault.AddressInfo {
	return s.Vault.AddressInfos()
}

// AddressCount returns the number of addresses inside the wallet
func (s *store) AddressCount() int {
	return s.Vault.AddressCount()
}

func (s *store) ImportPrivateKey(password string, prvStr string) error {
	return s.Vault.ImportPrivateKey(password, prvStr)
}

func (s *store) PrivateKey(password, addr string) (string, error) {
	return s.Vault.PrivateKey(password, addr)
}

func (s *store) PublicKey(password, addr string) (string, error) {
	return s.Vault.PublicKey(password, addr)
}

func (s *store) MakeNewAddress(password, label string) (string, error) {
	return s.Vault.MakeNewAddress(password, label)
}

func (s *store) Contains(addr string) bool {
	return s.Vault.Contains(addr)
}

func (s *store) Mnemonic(password string) (string, error) {
	return s.Vault.Mnemonic(password)
}
