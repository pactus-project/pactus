package stats

import (
	"github.com/zarbchain/zarb-go/consensus/hrs"
)

type Node struct {
	hrs hrs.HRS
}

func NewNode() *Node {
	return &Node{}
}
