package stats

import (
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/version"
)

type Peer struct {
	Version       version.Version
	Moniker       string
	PublicKey     crypto.PublicKey
	GenesisHash   crypto.Hash
	Height        int
	ReceivedMsg   int
	InvalidMsg    int
	ReceivedBytes int
}

func NewPeer() *Peer {
	return &Peer{}
}

func (p *Peer) BelongsToSameNetwork(genesisHash crypto.Hash) bool {
	if p.GenesisHash.IsUndef() {
		return true
	}
	return p.GenesisHash.EqualsTo(genesisHash)
}
