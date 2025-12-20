//go:build gtk

package model

import "github.com/pactus-project/pactus/node"

type NodeModel struct {
	Node *node.Node
}

func NewNodeModel(n *node.Node) *NodeModel {
	return &NodeModel{Node: n}
}
