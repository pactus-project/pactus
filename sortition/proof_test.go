package sortition

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProofFromString(t *testing.T) {
	_, err := ProofFromString("inv")
	assert.Error(t, err)
	_, err = ProofFromRawBytes([]byte{0})
	assert.Error(t, err)
}
