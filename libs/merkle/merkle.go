package simpleMerkle

import (
	"math"

	"github.com/zarbchain/zarb-go/crypto"
)

type SimpleMerkleTree struct {
	merkles []*crypto.Hash
}

// nextPowerOfTwo returns the next highest power of two from a given number if
// it is not already a power of two.  This is a helper function used during the
// calculation of a merkle tree.
func nextPowerOfTwo(n int) int {
	// Return the number if it's already a power of 2.
	if n&(n-1) == 0 {
		return n
	}

	// Figure out and return the next power of two.
	exponent := uint(math.Log2(float64(n))) + 1
	return 1 << exponent // 2^exponent
}

// HashMerkleBranches takes two hashes, treated as the left and right tree
// nodes, and returns the hash of their concatenation.  This is a helper
// function used to aid in the generation of a merkle tree.
func HashMerkleBranches(left *crypto.Hash, right *crypto.Hash) *crypto.Hash {
	// Concatenate the left and right nodes.
	var hash [crypto.HashSize * 2]byte
	copy(hash[:crypto.HashSize], left.RawBytes())
	copy(hash[crypto.HashSize:], right.RawBytes())

	newHash := crypto.HashH(hash[:])
	return &newHash
}

func NewTreeFromSlices(slices [][]byte) *SimpleMerkleTree {
	hashes := make([]crypto.Hash, len(slices))
	for i, b := range slices {
		hashes[i] = crypto.HashH(b)
	}

	return NewTreeFromHashes(hashes)
}

func NewTreeFromHashes(hashes []crypto.Hash) *SimpleMerkleTree {
	if len(hashes) == 0 {
		return nil
	}
	// abcdww Calculate how many entries are required to hold the binary merkle
	// tree as a linear array and create an array of that size.
	nextPoT := nextPowerOfTwo(len(hashes))
	arraySize := nextPoT*2 - 1
	merkles := make([]*crypto.Hash, arraySize)

	for i, _ := range hashes {
		merkles[i] = &hashes[i]
	}

	// Start the array offset after the last transaction and adjusted to the
	// next power of two.
	offset := nextPoT
	for i := 0; i < arraySize-1; i += 2 {
		switch {
		// When there is no left child node, the parent is nil too.
		case merkles[i] == nil:
			merkles[offset] = nil

		// When there is no right child, the parent is generated by
		// hashing the concatenation of the left child with itself.
		case merkles[i+1] == nil:
			newHash := HashMerkleBranches(merkles[i], merkles[i])
			merkles[offset] = newHash

		// The normal case sets the parent node to the double sha256
		// of the concatentation of the left and right children.
		default:
			newHash := HashMerkleBranches(merkles[i], merkles[i+1])
			merkles[offset] = newHash
		}
		offset++
	}

	return &SimpleMerkleTree{merkles: merkles}
}

func (tree *SimpleMerkleTree) Root() *crypto.Hash {
	if tree == nil {
		return &crypto.UndefHash
	}
	return tree.merkles[len(tree.merkles)-1]
}

func (tree *SimpleMerkleTree) Depth() int {
	if tree == nil {
		return 0
	}
	return int(math.Log2(float64(len(tree.merkles))))
}
