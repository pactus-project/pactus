package sortition

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProofFromString(t *testing.T) {
	_, err := ProofFromString("inv")
	assert.Error(t, err)
	_, err = ProofFromBytes([]byte{0})
	assert.Error(t, err)
}

func TestProofMarshaling(t *testing.T) {
	p1 := GenerateRandomProof()
	bz, err := json.Marshal(p1)
	assert.NoError(t, err)
	var p2 Proof
	assert.NoError(t, json.Unmarshal(bz, &p2))
	assert.Equal(t, p1, p2)
}
