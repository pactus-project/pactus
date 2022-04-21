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
	Name      string       `json:"name"`
	UUID      uuid.UUID    `json:"uuid"`
	CreatedAt time.Time    `json:"created_at"`
	Network   int          `json:"network"`
	VaultCRC  uint32       `json:"crc"`
	Vault     *vault.Vault `json:"vault"`
}

func (s *store) calcVaultCRC() uint32 {
	d, _ := json.Marshal(s.Vault)
	return crc32.ChecksumIEEE(d)
}

func (s *store) UpdatePassword(oldPassphrase, newPassphrase string) error {
	return s.Vault.UpdatePassword(oldPassphrase, newPassphrase)
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

func (s *store) ImportPrivateKey(passphrase string, prvStr string) error {
	return s.Vault.ImportPrivateKey(passphrase, prvStr)
}

func (s *store) PrivateKey(passphrase, addr string) (string, error) {
	return s.Vault.PrivateKey(passphrase, addr)
}

func (s *store) PublicKey(passphrase, addr string) (string, error) {
	return s.Vault.PublicKey(passphrase, addr)
}

func (s *store) MakeNewAddress(passphrase, label string) (string, error) {
	return s.Vault.MakeNewAddress(passphrase, label)
}

func (s *store) Contains(addr string) bool {
	return s.Vault.Contains(addr)
}

func (s *store) Mnemonic(passphrase string) (string, error) {
	return s.Vault.Mnemonic(passphrase)
}
