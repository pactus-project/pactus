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
	depth := 1
	offset := 0
	indent := ""
	j := 1
	for i := len(nodes) - 1; i >= 0; i-- {
		if j == (1 << depth) {
			lines[offset] += "\n"
			indent += "   "
			depth++
		}
		if nodes[i] != nil {
			lines[offset] = indent + nodes[i].String()
		} else {
			lines[offset] = indent + "<EMPTY>"
		}

		j++
		offset++
	}

	return strings.Join(lines, "\n")
}
