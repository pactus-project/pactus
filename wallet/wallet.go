package wallet

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
)

type Wallet struct {
	path  string
	store *Store
}

/// OpenWallet generates an empty wallet and save the seed string
func OpenWallet(path string) (*Wallet, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	s := new(Store)
	err = json.Unmarshal(data, s)
	exitOnErr(err)

	if s.VaultCRC != s.calcVaultCRC() {
		exitOnErr(errors.New("invalid CRC"))
	}

	return &Wallet{
		store: s,
		path:  path,
	}, nil
}

/// Recover recovers a wallet from mnemonic (seed phrase)
func RecoverWallet(path, mnemonic string, net int) (*Wallet, error) {
	store := RecoverStore(mnemonic, net)

	w := &Wallet{
		store: store,
		path:  path,
	}

	err := w.SaveToFile()
	if err != nil {
		return nil, err
	}

	return w, nil
}

/// NewWallet generates an empty wallet and save the seed string
func NewWallet(path, passphrase string, net int) (*Wallet, error) {
	store := NewStore(passphrase, net)
	w := &Wallet{
		store: store,
		path:  path,
	}

	err := w.SaveToFile()
	if err != nil {
		return nil, err
	}

	return w, nil
}

func (w *Wallet) IsEncrypted() bool {
	return w.store.Encrypted
}

func (w *Wallet) SaveToFile() error {
	w.store.VaultCRC = w.store.calcVaultCRC()

	bs, err := json.Marshal(w.store)
	exitOnErr(err)

	return ioutil.WriteFile(w.path, bs, 0600)
}

func (w *Wallet) ImportPrivateKey(passphrase string, prv *bls.PrivateKey) error {
	err := w.store.ImportPrivateKey(passphrase, prv)
	if err != nil {
		return err
	}
	return w.SaveToFile()
}

func (w *Wallet) PrivateKey(passphrase, addr string) (*bls.PrivateKey, error) {
	return w.store.PrivateKey(passphrase, addr)
}

func (w *Wallet) Mnemonic(passphrase string) string {
	return w.store.Mnemonic(passphrase)
}

func (w *Wallet) Addresses() []crypto.Address {
	return w.store.Addresses()
}
