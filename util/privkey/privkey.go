package privkey

import (
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/ed25519"
)

type PrivateKeyType uint8

const (
	// PrivateKeyUnknown defines the unknown PrivateKey type.
	PrivateKeyTypeUnknown PrivateKeyType = iota
	// PrivateKeyEd25519 defines the PrivateKey type for the Ed25519 signature algorithm.
	PrivateKeyTypeEd25519
	// PrivateKeyBLS defines the PrivateKey type for the BLS signature algorithm.
	PrivateKeyTypeBLS
)

func PrivateKeyTypeFromString(prvStr string) (PrivateKeyType, error) {
	_, err := bls.PrivateKeyFromString(prvStr)
	if err == nil {
		return PrivateKeyTypeBLS, nil
	}
	_, err = ed25519.PrivateKeyFromString(prvStr)
	if err == nil {
		return PrivateKeyTypeEd25519, nil
	}

	return PrivateKeyTypeUnknown, err
}
