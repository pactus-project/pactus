package cmd

import (
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
)

func ShortHash(id string) string {
	h, err := hash.FromString(id)
	if err != nil {
		return id
	}

	return h.ShortString()
}

func ShortAddress(addr string) string {
	a, err := crypto.AddressFromString(addr)
	if err != nil {
		return addr
	}

	return a.ShortString()
}
