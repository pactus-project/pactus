package block

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomBlock(t *testing.T) {
	b, _ := GenerateTestBlock(nil)
	assert.NoError(t, b.SanityCheck())
}

func TestMarshaling(t *testing.T) {
	b1, _ := GenerateTestBlock(nil)

	bz1, err := b1.MarshalCBOR()
	assert.NoError(t, err)
	var b2 Block
	err = b2.UnmarshalCBOR(bz1)
	assert.NoError(t, err)
	assert.NoError(t, b2.SanityCheck())
	assert.Equal(t, b1.Hash(), b2.Hash())

	bz2, _ := b1.MarshalCBOR()
	assert.Equal(t, bz1, bz2)
}
