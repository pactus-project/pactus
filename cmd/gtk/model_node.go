//go:build gtk

package main

import (
	"github.com/pactus-project/pactus/node"
)

type nodeModel struct {
	node *node.Node
}

func newNodeModel(node *node.Node) *nodeModel {
	return &nodeModel{
		node: node,
	}
}
