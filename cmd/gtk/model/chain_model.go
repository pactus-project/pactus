//go:build gtk

package model

import (
	"github.com/pactus-project/pactus/node"
)

type ChainModel struct {
	node *node.Node
}

func NewNodeModel(n *node.Node) *ChainModel {
	return &ChainModel{node: n}
}
