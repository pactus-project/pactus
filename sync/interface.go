package sync

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/sync/peerset"
)

type Synchronizer interface {
	Start() error
	Stop()
	PeerID() peer.ID
	Peers() []*peerset.Peer
	Fingerprint() string
}
