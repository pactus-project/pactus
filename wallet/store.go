package wallet

import (
	"encoding/json"
	"hash/crc32"
	"time"
)

///"use_encryption": true,

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
