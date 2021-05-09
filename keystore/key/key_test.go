package key

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyGeneration(t *testing.T) {
	k1 := GenerateRandomKey()
	k2 := GenerateRandomKey()
	k3, err := NewKey(k1.Address(), k2.PrivateKey())

	assert.NotNil(t, k1)
	assert.NotNil(t, k2)
	assert.Nil(t, k3)
	assert.Error(t, err)
}

func TestKeyFromSeed(t *testing.T) {
	_, err := FromSeed([]byte{0})
	assert.NoError(t, err)
}
