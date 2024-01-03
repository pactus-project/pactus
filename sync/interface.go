package sync

import (
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/sync/peerset"
	"github.com/pactus-project/pactus/sync/peerset/service"
)

type Synchronizer interface {
	Start() error
	Stop()
	Moniker() string
	SelfID() peer.ID
	PeerSet() *peerset.PeerSet
	Services() service.Services
}
