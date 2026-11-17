package persistentmerkle

import (
	"math"

	"github.com/pactus-project/pactus/crypto/hash"
)

// Tree is a merkle tree that keeps the hashes of its nodes. Setting a leaf invalidates
// the nodes above it, and Root only hashes again the invalidated nodes.
//
// The nodes are kept in one flat array per height, indexed by their width: a node costs
// its hash and one byte, where a node of a map cost about 170 bytes per leaf.
type Tree struct {
	levels    []level
	maxWidth  int32
	maxHeight int32
}

// level holds the nodes of one height of the tree.
type level struct {
	hashes []hash.Hash
	states []nodeState
}

type nodeState uint8

const (
	// nodeAbsent is a position without node. The node on its left is used instead.
	nodeAbsent nodeState = iota
	// nodeInvalid is a node whose hash must be computed from its children.
	nodeInvalid
	// nodeValid is a node whose hash is computed, or set for a leaf.
	nodeValid
)

func New() *Tree {
	return &Tree{}
}

func (t *Tree) nodeState(width, height int32) nodeState {
	if width < 0 || int(height) >= len(t.levels) {
		return nodeAbsent
	}

	lvl := &t.levels[height]
	if int(width) >= len(lvl.states) {
		return nodeAbsent
	}

	return lvl.states[width]
}

// setNode creates the node if it doesn't exist, and sets its state and its hash.
func (t *Tree) setNode(width, height int32, state nodeState, h hash.Hash) {
	for len(t.levels) <= int(height) {
		t.levels = append(t.levels, level{})
	}

	lvl := &t.levels[height]
	if missing := int(width) + 1 - len(lvl.states); missing > 0 {
		lvl.hashes = append(lvl.hashes, make([]hash.Hash, missing)...)
		lvl.states = append(lvl.states, make([]nodeState, missing)...)
	}

	lvl.states[width] = state
	lvl.hashes[width] = h
}

func (t *Tree) invalidateNode(width, height int32) {
	t.setNode(width, height, nodeInvalid, hash.UndefHash)
}

func (t *Tree) recalculateHeight(maxWidth int32) {
	if maxWidth > t.maxWidth {
		t.maxWidth = maxWidth

		maxHeight := math.Log2(float64(maxWidth))
		if math.Remainder(maxHeight, 1.0) != 0 {
			t.maxHeight = int32(math.Trunc(maxHeight)) + 2
		} else {
			t.maxHeight = int32(math.Trunc(maxHeight)) + 1
		}
	}
}

func (t *Tree) SetData(leaf int32, data []byte) {
	t.SetHash(leaf, hash.CalcHash(data))
}

func (t *Tree) SetHash(leaf int32, h hash.Hash) {
	t.recalculateHeight(leaf + 1)

	t.setNode(leaf, 0, nodeValid, h)

	w := leaf / 2
	for h := int32(1); h < t.maxHeight; h++ {
		t.invalidateNode(w, h)
		w /= 2
	}
}

// Root returns the root of the tree. The root of an empty tree is the undefined hash.
func (t *Tree) Root() hash.Hash {
	if t.maxHeight == 0 {
		return hash.UndefHash
	}

	return t.nodeHash(0, t.maxHeight-1)
}

func (t *Tree) nodeHash(width, height int32) hash.Hash {
	switch t.nodeState(width, height) {
	case nodeAbsent:
		// There is no node on the right: the node on the left is used instead.
		// It is hashed from its own children, never from the ones of this position.
		if t.nodeState(width-1, height) == nodeAbsent {
			panic("invalid merkle tree")
		}

		return t.nodeHash(width-1, height)

	case nodeValid:
		return t.levels[height].hashes[width]

	case nodeInvalid:
		left := t.nodeHash(width*2, height-1)
		right := t.nodeHash(width*2+1, height-1)

		data := make([]byte, len(left)+len(right))
		copy(data[:hash.HashSize], left.Bytes())
		copy(data[hash.HashSize:], right.Bytes())

		h := hash.CalcHash(data)
		t.setNode(width, height, nodeValid, h)

		return h
	}

	panic("invalid node state")
}
