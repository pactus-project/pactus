package firewall

import (
	"testing"
	"time"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/sync/message"
	"github.com/zarbchain/zarb-go/sync/peerset"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/validator"
)

var tFirewall *Firewall
var tAnotherPeerID peer.ID

func setup(t *testing.T) {
	peerSet := peerset.NewPeerSet(3 * time.Second)
	valSet, _ := validator.GenerateTestValidatorSet()
	state := state.MockingState(valSet)
	tFirewall = NewFirewall(peerSet, state)
	tAnotherPeerID = util.RandomPeerID()
}

func TestIncreaseMsgCounter(t *testing.T) {
	setup(t)

	msg := message.NewQueryProposalMessage(tAnotherPeerID, 1, 0)
	d, _ := msg.Encode()
	assert.NotNil(t, tFirewall.ParsMessage(d, tAnotherPeerID))
	p := tFirewall.peerSet.GetPeer(tAnotherPeerID)
	assert.False(t, tFirewall.badPeer(p))

	tFirewall.ParsMessage([]byte("bad"), tAnotherPeerID)
	p = tFirewall.peerSet.GetPeer(tAnotherPeerID)
	assert.True(t, tFirewall.badPeer(p))
	assert.Nil(t, tFirewall.ParsMessage(d, tAnotherPeerID))
}
