package key

import (
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
)

type Key struct {
	data keyData
}

type keyData struct {
	Address    crypto.Address
	PublicKey  crypto.PublicKey
	PrivateKey crypto.PrivateKey
}

func FromSeed(seed []byte) (*Key, error) {
	prv, err := bls.PrivateKeyFromSeed(seed, nil)
	if err != nil {
		return nil, err
	}
	return &Key{
		data: keyData{
			PrivateKey: prv,
			PublicKey:  prv.PublicKey(),
			Address:    prv.PublicKey().Address(),
		},
	}, nil
}

// NewKey Checks if the address is derived from the given private key
func NewKey(pv crypto.PrivateKey) *Key {
	return &Key{
		data: keyData{
			PrivateKey: pv,
			PublicKey:  pv.PublicKey(),
			Address:    pv.PublicKey().Address(),
		},
	}
}

func (k *Key) Address() crypto.Address {
	return k.data.Address
}

func (k *Key) PublicKey() crypto.PublicKey {
	return k.data.PublicKey
}

func (k *Key) PrivateKey() crypto.PrivateKey {
	return k.data.PrivateKey
}

func (k *Key) ToSigner() crypto.Signer {
	return crypto.NewSigner(k.data.PrivateKey)
}
