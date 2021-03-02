package sortition

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
)

func TestSeedFromString(t *testing.T) {
	_, err := SeedFromString("inv")
	assert.Error(t, err)
	_, err = SeedFromRawBytes([]byte{0})
	assert.Error(t, err)
}

func TestValidate(t *testing.T) {
	signer := crypto.GenerateTestSigner()
	seed1 := GenerateRandomSeed()
	seed2 := seed1.Generate(signer)

	assert.True(t, seed2.Validate(signer.PublicKey(), seed1))
	assert.False(t, seed1.Validate(signer.PublicKey(), seed2))
	assert.False(t, seed2.Validate(signer.PublicKey(), GenerateRandomSeed()))
}
