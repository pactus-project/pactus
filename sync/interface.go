package sync

import (
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/sync/peerset"
)

type Synchronizer interface {
	Start() error
	Stop()
	SelfID() peer.ID
	Peers() []peerset.Peer
	Fingerprint() string
}
