package wallet

import (
	"github.com/tyler-smith/go-bip39"
	"github.com/zarbchain/zarb-go/crypto/bls"
)

type seed struct {
	Method     string    `json:"Method"`
	ParentSeed encrypted `json:"seed"`
	ParentKey  encrypted `json:"xprv"`
}

func recoverSeed(mnemonic string) seed {
	ikm, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	exitOnErr(err)

	parentKey, err := bls.PrivateKeyFromSeed(ikm, nil)

	s := seed{
		Method:     "BIP-39",
		ParentSeed: newNopeEncrypter().encrypt(mnemonic),
		ParentKey:  newNopeEncrypter().encrypt(parentKey.String()),
	}

	return s
}

func newSeed(passphrase string) seed {
	entropy, err := bip39.NewEntropy(128)
	exitOnErr(err)
	mnemonic, err := bip39.NewMnemonic(entropy)
	exitOnErr(err)
	ikm, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	exitOnErr(err)

	parentKey, err := bls.PrivateKeyFromSeed(ikm, nil)

	s := seed{
		Method:     "BIP-39",
		ParentSeed: newEncrypter(passphrase).encrypt(mnemonic),
		ParentKey:  newEncrypter(passphrase).encrypt(parentKey.String()),
	}

	return s
}

func (s *seed) parentSeed(passphrase string) []byte {
	h, err := bip39.NewSeedWithErrorChecking(s.mnemonic(passphrase), "")
	exitOnErr(err)

	return h
}

func (s *seed) mnemonic(passphrase string) string {
	m, err := newEncrypter(passphrase).decrypt(s.ParentSeed)
	exitOnErr(err)

	return m
}

func (s *seed) parentKey(passphrase string) *bls.PrivateKey {
	m, err := newEncrypter(passphrase).decrypt(s.ParentKey)
	prv, err := bls.PrivateKeyFromString(m)
	exitOnErr(err)

	return prv
}
