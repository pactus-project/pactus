package sortition_test

import (
	"testing"

	"github.com/pactus-project/pactus/sortition"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestSeedFromString(t *testing.T) {
	_, err := sortition.VerifiableSeedFromString("inv")
	assert.Error(t, err)
	_, err = sortition.VerifiableSeedFromBytes([]byte{0})
	assert.Error(t, err)
}

func TestValidate(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	signer := ts.RandSigner()
	seed1 := ts.RandSeed()
	seed2 := seed1.GenerateNext(signer)
	seed3 := sortition.VerifiableSeed{}
	seed4, _ := sortition.VerifiableSeedFromString(
		"C00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")

	assert.True(t, seed2.Verify(signer.PublicKey(), seed1))
	assert.False(t, seed1.Verify(signer.PublicKey(), seed2))
	assert.False(t, seed2.Verify(signer.PublicKey(), ts.RandSeed()))
	assert.False(t, seed3.Verify(signer.PublicKey(), seed1))
	assert.False(t, seed4.Verify(signer.PublicKey(), seed1))
}
