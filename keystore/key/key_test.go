package key

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyFromSeed(t *testing.T) {
	_, err := FromSeed([]byte{0})
	assert.NoError(t, err)
}
