package firewall

import (
	"testing"

	peer "github.com/libp2p/go-libp2p-peer"
	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/state"

	"github.com/zarbchain/zarb-go/sync/peerset"
)

var tFirewall *Firewall
var tAnotherPeerID peer.ID

func setup(t *testing.T) {
	peerSet := peerset.NewPeerSet()
	state := state.MockingState()
	tFirewall = NewFirewall(peerSet, state)
	tAnotherPeerID, _ = peer.IDB58Decode("12D3KooWBtNwU6PiV9KrVXqhNeoeP2vrvJs7USAXtkapgCs6TwUm")
}

func TestIncreaseMsgCounter(t *testing.T) {
	setup(t)

	tFirewall.ParsMessage([]byte("bad"), tAnotherPeerID)
	p := tFirewall.peerSet.GetPeer(tAnotherPeerID)
	assert.True(t, tFirewall.badPeer(p))
}
