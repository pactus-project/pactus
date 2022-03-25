package sortition

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto/bls"
)

func TestSeedFromString(t *testing.T) {
	_, err := VerifiableSeedFromString("inv")
	assert.Error(t, err)
	_, err = VerifiableSeedFromBytes([]byte{0})
	assert.Error(t, err)
}

func TestValidate(t *testing.T) {
	signer := bls.GenerateTestSigner()
	seed1 := GenerateRandomSeed()
	seed2 := seed1.Generate(signer)
	seed3 := VerifiableSeed{}
	seed4, _ := VerifiableSeedFromString("C00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")

	assert.True(t, seed2.Verify(signer.PublicKey(), seed1))
	assert.False(t, seed1.Verify(signer.PublicKey(), seed2))
	assert.False(t, seed2.Verify(signer.PublicKey(), GenerateRandomSeed()))
	assert.False(t, seed3.Verify(signer.PublicKey(), seed1))
	assert.False(t, seed4.Verify(signer.PublicKey(), seed1))
}

func TestSeedMarshaling(t *testing.T) {
	s1 := GenerateRandomSeed()
	bz, err := json.Marshal(s1)
	assert.NoError(t, err)
	var s2 VerifiableSeed
	assert.NoError(t, json.Unmarshal(bz, &s2))
	assert.Equal(t, s1, s2)
}
