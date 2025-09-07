package peer

import (
	"time"

	lp2pnetwork "github.com/libp2p/go-libp2p/core/network"
	lp2ppeer "github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/sync/peerset/peer/metric"
	"github.com/pactus-project/pactus/sync/peerset/peer/service"
	"github.com/pactus-project/pactus/sync/peerset/peer/status"
)

type ID = lp2ppeer.ID

type Peer struct {
	Status            status.Status
	Moniker           string
	Agent             string
	Address           string
	Direction         lp2pnetwork.Direction
	Protocols         []string
	PeerID            ID
	ConsensusKeys     []*bls.PublicKey
	Services          service.Services
	LastSent          time.Time
	LastReceived      time.Time
	LastBlockHash     hash.Hash
	Height            uint32
	TotalSessions     int
	CompletedSessions int
	Metric            metric.Metric
}

func NewPeer(peerID ID) *Peer {
	return &Peer{
		ConsensusKeys: make([]*bls.PublicKey, 0),
		Status:        status.StatusUnknown,
		PeerID:        peerID,
		Metric:        metric.NewMetric(),
		Protocols:     make([]string, 0),
	}
}

func (p *Peer) IsFullNode() bool {
	return p.Services.IsFullNode()
}

func (p *Peer) DownloadScore() int {
	return (p.CompletedSessions + 1) * 100 / (p.TotalSessions + 1)
}
