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

	valKey := ts.RandValKey()
	seed1 := ts.RandSeed()
	seed2 := seed1.GenerateNext(valKey.PrivateKey())
	seed3 := sortition.VerifiableSeed{}
	seed4, _ := sortition.VerifiableSeedFromString(
		"C00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")

	assert.True(t, seed2.Verify(valKey.PublicKey(), seed1))
	assert.False(t, seed1.Verify(valKey.PublicKey(), seed2))
	assert.False(t, seed2.Verify(valKey.PublicKey(), ts.RandSeed()))
	assert.False(t, seed3.Verify(valKey.PublicKey(), seed1))
	assert.False(t, seed4.Verify(valKey.PublicKey(), seed1))
}
