package sortition

import (
	"testing"

	"github.com/fxamacker/cbor/v2"
	"github.com/stretchr/testify/assert"
)

func TestProofMarshaling(t *testing.T) {
	p1 := GenerateRandomProof()
	bz, err := cbor.Marshal(p1)
	assert.NoError(t, err)
	var p2 Proof
	assert.NoError(t, cbor.Unmarshal(bz, &p2))
	assert.Equal(t, p1, p2)
}
