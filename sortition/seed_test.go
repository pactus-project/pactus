package sortition

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto/bls"
)

func TestValidate(t *testing.T) {
	signer := bls.GenerateTestSigner()
	seed1 := GenerateRandomSeed()
	seed2 := seed1.Generate(signer)

	assert.True(t, seed2.Verify(signer.PublicKey(), seed1))
	assert.False(t, seed1.Verify(signer.PublicKey(), seed2))
	assert.False(t, seed2.Verify(signer.PublicKey(), GenerateRandomSeed()))
}
