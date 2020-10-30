package test

import (
	"crypto/rand"

	"github.com/zarbchain/zarb-go/crypto"
)

// ---------
// For tests
func GenerateTestHash() crypto.Hash {
	p := make([]byte, 10)
	random := rand.Reader
	random.Read(p)
	return crypto.HashH(p)
}
