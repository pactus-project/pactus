package stats

import (
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/version"
)

type Node struct {
	Version     version.Version
	GenesisHash crypto.Hash
	HRS         hrs.HRS
}

func NewNode() *Node {
	return &Node{}
}
