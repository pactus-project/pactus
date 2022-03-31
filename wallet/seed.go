package wallet

import (
	"github.com/tyler-smith/go-bip39"
	"github.com/zarbchain/zarb-go/crypto/bls"
)

type seed struct {
	Method    string    `json:"Method"`
	Seed      encrypted `json:"seed"`
	ParentKey string    `json:"parent_key"` // No Master!
}

func newSeed(passphrase string) seed {
	entropy, err := bip39.NewEntropy(128)
	exitOnErr(err)
	mnemonic, err := bip39.NewMnemonic(entropy)
	exitOnErr(err)
	ikm, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	exitOnErr(err)

	parentKey, err := bls.PrivateKeyFromSeed(ikm, nil)

	e := newEncrypter(passphrase)
	s := seed{
		Method:    "BIP-39",
		Seed:      e.encrypt(mnemonic),
		ParentKey: parentKey.String(),
	}

	return s
}

func (s *seed) hashedSeed(passphrase string) []byte {
	h, err := bip39.NewSeedWithErrorChecking(s.mnemonic(passphrase), "")
	exitOnErr(err)

	return h
}

func (s *seed) mnemonic(passphrase string) string {
	e := newEncrypter(passphrase)
	m, err := e.decrypt(s.Seed)
	exitOnErr(err)

	return m
}

func (s *seed) parentKey() *bls.PrivateKey {
	prv, err := bls.PrivateKeyFromString(s.ParentKey)
	exitOnErr(err)

	return prv
}
