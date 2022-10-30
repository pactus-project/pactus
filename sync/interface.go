package sync

import (
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/sync/peerset"
)

type Synchronizer interface {
	Start() error
	Stop()
	Moniker() string
	PublicKey() crypto.PublicKey
	SelfID() peer.ID
	Peers() []peerset.Peer
	Fingerprint() string
}
