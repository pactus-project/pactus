package persistentmerkle

import (
	"math/rand"
	"testing"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/stretchr/testify/require"
)

// referenceRoot computes the root from all the leaves, level by level. A node without
// a right sibling is hashed with itself, like the tree does.
func referenceRoot(leaves []hash.Hash) hash.Hash {
	if len(leaves) == 0 {
		return hash.UndefHash
	}

	level := leaves
	for len(level) > 1 {
		next := make([]hash.Hash, 0, (len(level)+1)/2)
		for i := 0; i < len(level); i += 2 {
			left := level[i]
			right := left
			if i+1 < len(level) {
				right = level[i+1]
			}

			data := make([]byte, 0, 2*hash.HashSize)
			data = append(data, left.Bytes()...)
			data = append(data, right.Bytes()...)
			next = append(next, hash.CalcHash(data))
		}
		level = next
	}

	return level[0]
}

// The root is the one of the reference, for any sequence of new leaves and updates.
func TestRootMatchesReference(t *testing.T) {
	for seed := int64(1); seed <= 20; seed++ {
		rnd := rand.New(rand.NewSource(seed))
		tree := New()
		leaves := []hash.Hash{}

		for step := range 1500 {
			leaf := len(leaves)
			if leaf > 0 && rnd.Intn(3) > 0 {
				leaf = rnd.Intn(len(leaves))
			}

			data := make([]byte, 8)
			_, _ = rnd.Read(data)
			h := hash.CalcHash(data)
			if leaf == len(leaves) {
				leaves = append(leaves, h)
			} else {
				leaves[leaf] = h
			}
			tree.SetHash(int32(leaf), h)

			if rnd.Intn(5) == 0 || step == 1499 {
				require.Equal(t, referenceRoot(leaves), tree.Root(), "seed %d, step %d", seed, step)
			}
		}
	}
}
