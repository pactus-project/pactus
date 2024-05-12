//go:build gtk

package main

import (
	"github.com/pactus-project/pactus/node"
	"github.com/pactus-project/pactus/util/ntp"
)

type nodeModel struct {
	node *node.Node
	ntp  *ntp.Server
}

func newNodeModel(nde *node.Node) *nodeModel {
	return &nodeModel{
		node: nde,
		ntp:  ntp.NewNtpServer(),
	}
}
