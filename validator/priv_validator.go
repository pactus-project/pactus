package validator

import (
	"gitlab.com/zarb-chain/zarb-go/crypto"
)

type Signable interface {
	SignBytes() []byte
	SetSignature(sig crypto.Signature)
}
type PrivValidator struct {
	address    crypto.Address
	publicKey  crypto.PublicKey
	privateKey crypto.PrivateKey
}

func NewPrivValidator(pv crypto.PrivateKey) *PrivValidator {
	return &PrivValidator{
		privateKey: pv,
		publicKey:  pv.PublicKey(),
		address:    pv.PublicKey().Address(),
	}
}

func (pv *PrivValidator) Address() crypto.Address {
	return pv.address
}

func (pv *PrivValidator) PublicKey() crypto.PublicKey {
	return pv.publicKey
}

func (pv *PrivValidator) SignMsg(msg Signable) {
	bz := msg.SignBytes()
	sig := pv.privateKey.Sign(bz)
	msg.SetSignature(sig)
}
