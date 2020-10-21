package stats

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/logger"
)

// Stats hold statistic data about peers' behaviors
type Stats struct {
	peers map[peer.ID]Peer
}

func NewStats(logger *logger.Logger) *Stats {
	return &Stats{
		peers: make(map[peer.ID]Peer),
	}
}
