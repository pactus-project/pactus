package key

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyFromSeed(t *testing.T) {
	_, err := FromSeed([]byte{0})
	assert.Error(t, err)

	seed := [32]byte{}
	_, err = FromSeed(seed[:])
	assert.NoError(t, err)
}
