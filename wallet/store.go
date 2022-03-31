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
	Vault     encrypted `json:"vault"`
}

func newStore(passphrase string, net int) (*store, error) {
	w := &store{
		Version:   1,
		CreatedAt: time.Now(),
		Network:   net,
		Encrypted: len(passphrase) > 0,
	}

	v := newVault(passphrase)

	e := newEncrypter(passphrase)
	w.Vault = encryptInterface(e, v)

	return w, nil
}

func (s *store) calcVaultCRC() uint32 {
	d, err := json.Marshal(s.Vault)
	exitOnErr(err)
	return crc32.ChecksumIEEE(d)
}

func (s *store) Addresses(passphrase string) []crypto.Address {
	e := newEncrypter(passphrase)
	js, err := e.decrypt(s.Vault)
	exitOnErr(err)
	v := new(vault)

	err = json.Unmarshal([]byte(js), v)
	exitOnErr(err)

	addrs := make([]crypto.Address, len(v.Addresses))
	for i, a := range v.Addresses {
		addr, err := crypto.AddressFromString(a.Address)
		exitOnErr(err)
		addrs[i] = addr
	}

	return addrs
}

func (s *store) PrivateKey(passphrase, addr string) (*bls.PrivateKey, error) {
	e := newEncrypter(passphrase)
	js, err := e.decrypt(s.Vault)
	exitOnErr(err)
	v := new(vault)

	err = json.Unmarshal([]byte(js), v)
	exitOnErr(err)

	return v.PrivateKey(passphrase, addr)
}
