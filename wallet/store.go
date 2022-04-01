package wallet

import (
	"encoding/json"
	"hash/crc32"
	"time"

	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
)

type store struct {
	Version   int       `json:"version"`
	CreatedAt time.Time `json:"created_at"`
	Network   int       `json:"network"`
	Encrypted bool      `json:"encrypted"`
	VaultCRC  uint32    `json:"crc"`
	Vault     *vault    `json:"vault"`
}

func recoverStore(mnemonic string, net int) (*store, error) {
	w := &store{
		Version:   1,
		CreatedAt: time.Now(),
		Network:   net,
		Encrypted: false,
		Vault:     recoverVault(mnemonic),
	}
	return w, nil
}

func newStore(passphrase string, net int) (*store, error) {
	w := &store{
		Version:   1,
		CreatedAt: time.Now(),
		Network:   net,
		Encrypted: len(passphrase) > 0,
		Vault:     newVault(passphrase),
	}
	return w, nil
}

func (s *store) calcVaultCRC() uint32 {
	d, err := json.Marshal(s.Vault)
	exitOnErr(err)
	return crc32.ChecksumIEEE(d)
}

func (s *store) Addresses(passphrase string) []crypto.Address {

	addrs := make([]crypto.Address, len(s.Vault.Addresses))
	for i, a := range s.Vault.Addresses {
		addr, err := crypto.AddressFromString(a.Address)
		exitOnErr(err)
		addrs[i] = addr
	}

	return addrs
}

func (s *store) PrivateKey(passphrase, addr string) (*bls.PrivateKey, error) {

	return s.Vault.PrivateKey(passphrase, addr)
}
