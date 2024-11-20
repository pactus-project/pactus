package persistentmerkle

import (
	"math"

	"github.com/pactus-project/pactus/crypto/hash"
)

type Tree struct {
	nodes     map[uint32]*node
	maxWidth  int32
	maxHeight int32
}

type node struct {
	width  int32
	height int32
	hash   *hash.Hash
}

// nodeID return the node ID (four bytes):
// +-+---+
// |h| w |
// +-+---+
// h: height
// w: width
func nodeID(width, height int32) uint32 {
	return (uint32(height&0xff) << 24) | uint32(width&0xffffff)
}

func New() *Tree {
	return &Tree{
		nodes: make(map[uint32]*node),
	}
}

func (*Tree) createNode(width, height int32) *node {
	return &node{
		width:  width,
		height: height,
	}
}

func (t *Tree) getNode(width, height int32) *node {
	id := nodeID(width, height)

	return t.nodes[id]
}

func (t *Tree) getOrCreateNode(width, height int32) *node {
	id := nodeID(width, height)
	node, ok := t.nodes[id]
	if !ok {
		node = t.createNode(width, height)
		t.nodes[id] = node
	}

	return node
}

func (t *Tree) invalidateNode(width, height int32) {
	n := t.getOrCreateNode(width, height)
	n.hash = nil
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

	node := t.getOrCreateNode(leaf, 0)
	node.hash = &h

	w := leaf / 2
	for h := int32(1); h < t.maxHeight; h++ {
		t.invalidateNode(w, h)
		w /= 2
	}
}

func (t *Tree) Root() hash.Hash {
	return t.nodeHash(0, t.maxHeight-1)
}

func (t *Tree) nodeHash(width, height int32) hash.Hash {
	node := t.getNode(width, height)
	if node == nil {
		node = t.getNode(width-1, height)
		if node == nil {
			panic("invalid merkle tree")
		}
	}
	if node.hash != nil {
		return *node.hash
	}

	left := t.nodeHash(width*2, height-1)
	right := t.nodeHash(width*2+1, height-1)

	data := make([]byte, len(left)+len(right))
	copy(data[:hash.HashSize], left.Bytes())
	copy(data[hash.HashSize:], right.Bytes())

	h := hash.CalcHash(data)
	node.hash = &h

	return h
}
