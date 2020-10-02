package simpleMerkle

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
)

func strToHash(str string) crypto.Hash {
	h := crypto.HashH([]byte(str))
	return h
}

func TestMerkleTree(t *testing.T) {
	slices := [][]byte{[]byte("a"), []byte("b"), []byte("c")}
	hashes := []crypto.Hash{
		strToHash("a"),
		strToHash("b"),
		strToHash("c"),
	}

	tree1 := NewTreeFromHashes(hashes)
	if tree1.Root().String() != "905b17edcf8b6fb1415b32cdbab3e02c2c93f80a345de80ea2bbf9feba9f5a55" {
		t.Errorf("Invalid merkle root")
	}

	tree2 := NewTreeFromSlices(slices)
	assert.Equal(t, tree1.Root(), tree2.Root())
}
