package wallet

import (
	"encoding/json"
	"errors"
	"hash/crc32"
	"io/ioutil"
	"time"
)

type Wallet interface {
	SaveToFile(path string) error
	
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

func (w *wallet) SaveToFile(path string) error {
	w.VaultCRC = w.calcVaultCRC()

	bs, err := json.Marshal(w)
	exitOnErr(err)

	return ioutil.WriteFile(path, bs, 0600)
}

func (w *wallet) ReadFromFile(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, w)
	exitOnErr(err)

	if w.VaultCRC != w.calcVaultCRC() {
		exitOnErr(errors.New("invalid CRC"))
	}

	return nil
}

func (w *wallet) calcVaultCRC() uint32 {
	d, err := json.Marshal(w)
	exitOnErr(err)
	return crc32.ChecksumIEEE(d)
}
