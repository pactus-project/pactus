package wallet

import (
	"encoding/json"
	"errors"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"time"

	"github.com/tyler-smith/go-bip39"
)

type Wallet interface {
	SaveToFile(path string) error
}

///"use_encryption": true,

type wallet struct {
	Version   uint32    `json:"version"`
	CreatedAt time.Time `json:"created_at"`
	VaultCRC  uint32    `json:"checksum"`
	Vault     encrypted `json:"vault"`
}

/// cipher text
type encrypted struct {
	Method     string `json:"method,omitempty"`
	Params     params `json:"params,omitempty"`
	CipherText string `json:"ct"`
}

/// GenerateWallet generates an empty wallet and save the seed string
func GenerateWallet(passphrase string) (Wallet, error) {
	w := &wallet{
		Version:   1,
		CreatedAt: time.Now(),
	}

	e := w.encrypter(passphrase)

	entropy, err := bip39.NewEntropy(128)
	exitOnErr(err)
	mnemonic, err := bip39.NewMnemonic(entropy)
	exitOnErr(err)

	v := &vault{
		Seed: seed{
			Function: "BIP_39",
			Seed:     e.encrypt(mnemonic, passphrase),
		},
	}

	// v.Addresses = make
	// ikm, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	// exitOnErr(err)
	// for i := 0; i < 21; i++ {
	// 	prv, err := bls.PrivateKeyFromSeed(ikm, util.Uint32ToSlice(uint32(i)))
	// 	exitOnErr(err)
	// }

	//prv, err := bls.PrivateKeyFromSeed(seed, nil)

	w.Vault = encryptInterface(e, v, passphrase)

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

	fmt.Printf("%s", bs)

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
