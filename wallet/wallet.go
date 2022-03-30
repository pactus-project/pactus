package wallet

import (
	"time"
)

type Wallet interface {
}

type wallet struct {
	Version   uint32    `json:"version"`
	CreatedAt time.Time `json:"created_at"`
	VaultCRC  uint32    `json:"checksum"`
	Vault     encrypted `json:"vault"`
}

/// cipher text
type encrypted struct {
	Method  string `json:"method,omitempty"`
	Params  params `json:"params,omitempty"`
	Message string `json:"message"`
}

/// GenerateWallet generates an empty wallet and save the seed string
func GenerateWallet(seed string, passphrase string) (Wallet, error) {
	w := &wallet{
		Version:   1,
		CreatedAt: time.Now(),
	}

	e := w.encrypter(passphrase)
	w.Vault = e.encrypt("test", passphrase)

	return w, nil
}

func (w *wallet) encrypter(passphrase string) encrypter {
	if len(passphrase) == 0 {
		return &nopeEncrypter{}
	}
	return &argon2Encrypter{}
}
