package wallet

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
)

type Wallet struct {
	path  string
	store *store
}

/// NewWallet generates an empty wallet and save the seed string
func NewWallet(path, passphrase string) (*Wallet, error) {
	store, err := newStore(passphrase, 0)
	if err != nil {
		return nil, err
	}

	w := &Wallet{
		store: store,
		path:  path,
	}

	err = w.SaveToFile()
	if err != nil {
		return nil, err
	}

	return w, nil
}

func (w *Wallet) SaveToFile() error {
	w.store.VaultCRC = w.store.calcVaultCRC()

	bs, err := json.Marshal(w.store)
	exitOnErr(err)

	fmt.Printf("%s", bs)

	return ioutil.WriteFile(w.path, bs, 0600)
}

func (w *Wallet) PrivateKey(passphrase, addr string) (*bls.PrivateKey, error) {
	return w.store.PrivateKey(passphrase, addr)
}

func (w *Wallet) Addresses(passphrase string) []crypto.Address {
	return w.store.Addresses(passphrase)
}

func (w *Wallet) ReadFromFile(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, w)
	exitOnErr(err)

	if w.store.VaultCRC != w.store.calcVaultCRC() {
		exitOnErr(errors.New("invalid CRC"))
	}

	return nil
}
