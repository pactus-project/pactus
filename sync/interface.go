package sync

import (
	"github.com/pactus-project/pactus/sync/peerset"
	"github.com/pactus-project/pactus/sync/peerset/peer"
	"github.com/pactus-project/pactus/sync/peerset/peer/service"
)

type Synchronizer interface {
	Start() error
	Stop()
	Moniker() string
	SelfID() peer.ID
	PeerSet() *peerset.PeerSet
	Services() service.Services
}
