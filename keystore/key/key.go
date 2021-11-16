package key

import (
	"fmt"

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

func GenerateRandomKey() *Key {
	addr, pk, pv := bls.RandomKeyPair()
	return &Key{
		data: keyData{
			PrivateKey: pv,
			PublicKey:  pk,
			Address:    addr,
		},
	}
}

func FromSeed(seed []byte) (*Key, error) {
	priv, err := bls.PrivateKeyFromSeed(seed)
	if err != nil {
		return nil, err
	}
	return &Key{
		data: keyData{
			PrivateKey: priv,
			PublicKey:  priv.PublicKey(),
			Address:    priv.PublicKey().Address(),
		},
	}, nil
}

// NewKey Checks if the address is derived from the given private key
func NewKey(addr crypto.Address, pv crypto.PrivateKey) (*Key, error) {
	if !addr.Verify(pv.PublicKey()) {
		return nil, fmt.Errorf("this address doesn't belong to this privatekey")
	}

	return &Key{
		data: keyData{
			PrivateKey: pv,
			PublicKey:  pv.PublicKey(),
			Address:    addr,
		},
	}, nil
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
