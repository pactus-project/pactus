package wallet

import (
	"encoding/json"
	"hash/crc32"
	"time"

	"github.com/google/uuid"
	"github.com/zarbchain/zarb-go/types/crypto"
	"github.com/zarbchain/zarb-go/wallet/vault"
)

type store struct {
	data storeData
}

type storeData struct {
	Version   int          `json:"version"`
	UUID      uuid.UUID    `json:"uuid"`
	CreatedAt time.Time    `json:"created_at"`
	Network   Network      `json:"network"`
	VaultCRC  uint32       `json:"crc"`
	Vault     *vault.Vault `json:"vault"`
}

func (s *store) MarshalJSON() ([]byte, error) {
	s.data.VaultCRC = s.calcVaultCRC()
	return json.MarshalIndent(&s.data, "  ", "  ")
}

func (s *store) UnmarshalJSON(bs []byte) error {
	err := json.Unmarshal(bs, &s.data)
	if err != nil {
		return err
	}

	if s.data.VaultCRC != s.calcVaultCRC() {
		return ErrInvalidCRC
	}

	return nil
}

func (s *store) calcVaultCRC() uint32 {
	d, _ := json.Marshal(s.data.Vault)
	return crc32.ChecksumIEEE(d)
}

func (s *store) UpdatePassword(oldPassword, newPassword string) error {
	return s.data.Vault.UpdatePassword(oldPassword, newPassword)
}

func (s *store) IsEncrypted() bool {
	return s.data.Vault.IsEncrypted()
}

func (s *store) AddressInfo(addr string) *vault.AddressInfo {
	return s.data.Vault.AddressInfo(addr)
}

func (s *store) AddressLabels() []vault.AddressInfo {
	return s.data.Vault.AddressLabels()
}

// AddressCount returns the number of addresses inside the wallet.
func (s *store) AddressCount() int {
	return s.data.Vault.AddressCount()
}

func (s *store) ImportPrivateKey(password string, prv crypto.PrivateKey) error {
	return s.data.Vault.ImportPrivateKey(password, prv)
}

func (s *store) PrivateKey(password, addr string) (crypto.PrivateKey, error) {
	return s.data.Vault.PrivateKey(password, addr)
}

func (s *store) DeriveNewAddress(label string) (string, error) {
	return s.data.Vault.DeriveNewAddress(label, vault.PurposeBLS12381)
}

func (s *store) Contains(addr string) bool {
	return s.data.Vault.Contains(addr)
}

func (s *store) Mnemonic(password string) (string, error) {
	return s.data.Vault.Mnemonic(password)
}

// Label returns label of addr.
func (s *store) Label(addr string) string {
	return s.data.Vault.Label(addr)
}

// SetLabel sets label for addr.
func (s *store) SetLabel(addr, label string) error {
	return s.data.Vault.SetLabel(addr, label)
}
