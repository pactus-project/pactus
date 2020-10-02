package simpleMerkle

import (
	"strings"
)

func (tree *SimpleMerkleTree) ToString() string {
	nodes := tree.merkles
	if len(nodes) == 0 {
		return ""
	}

	lines := make([]string, len(tree.merkles))
	depth := tree.Depth() - 1
	offset := 0
	indent := ""
	for i := len(nodes) - 1; i >= 0; i-- {
		if i == (1 << depth) {
			indent += "   "
		}
		if nodes[i] != nil {
			lines[offset] = indent + nodes[i].String()
		}
		offset++
	}

	return strings.Join(lines, "\n")
}
