package wallet

import (
	"encoding/json"
	"hash/crc32"
	"time"

	"github.com/google/uuid"
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

func (s *store) AddressInfos() []vault.AddressInfo {
	return s.data.Vault.AddressInfos()
}

// AddressCount returns the number of addresses inside the wallet
func (s *store) AddressCount() int {
	return s.data.Vault.AddressCount()
}

func (s *store) ImportPrivateKey(password string, prvStr string) error {
	return s.data.Vault.ImportPrivateKey(password, prvStr)
}

func (s *store) PrivateKey(password, addr string) (string, error) {
	return s.data.Vault.PrivateKey(password, addr)
}

func (s *store) PublicKey(password, addr string) (string, error) {
	return s.data.Vault.PublicKey(password, addr)
}

func (s *store) MakeNewAddress(password, label string) (string, error) {
	return s.data.Vault.MakeNewAddress(password, label)
}

func (s *store) Contains(addr string) bool {
	return s.data.Vault.Contains(addr)
}

func (s *store) Mnemonic(password string) (string, error) {
	return s.data.Vault.Mnemonic(password)
}

// SetLabel returns label of addrStr
func (s *store) Label(addrStr string) string {
	return s.data.Vault.Label(addrStr)
}

// SetLabel sets label for addr
func (s *store) SetLabel(addrStr, label string) error {
	return s.data.Vault.SetLabel(addrStr, label)
}
