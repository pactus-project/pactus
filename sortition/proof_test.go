package sortition

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProofFromString(t *testing.T) {
	_, err := ProofFromString("inv")
	require.Error(t, err)
	_, err = ProofFromBytes([]byte{0})
	require.Error(t, err)
}
