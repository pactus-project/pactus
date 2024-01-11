package bls

import (
	"github.com/pactus-project/pactus/crypto"
)

// ValidatorKey wraps a BLS private key, caching its public key and validator address.
type ValidatorKey struct {
	address    crypto.Address
	publicKey  *PublicKey
	privateKey *PrivateKey
}

func NewValidatorKey(prv *PrivateKey) *ValidatorKey {
	pub := prv.PublicKeyNative()

	return &ValidatorKey{
		privateKey: prv,
		publicKey:  pub,
		address:    pub.ValidatorAddress(),
	}
}

func (key *ValidatorKey) Address() crypto.Address {
	return key.address
}

func (key *ValidatorKey) PublicKey() *PublicKey {
	return key.publicKey
}

func (key *ValidatorKey) PrivateKey() *PrivateKey {
	return key.privateKey
}

func (key *ValidatorKey) Sign(data []byte) *Signature {
	return key.privateKey.SignNative(data)
}
