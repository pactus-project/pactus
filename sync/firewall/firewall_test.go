package firewall

import (
	"testing"
	"time"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/committee"
	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/sync/message"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/sync/peerset"
	"github.com/zarbchain/zarb-go/util"
)

var tFirewall *Firewall
var tAnotherPeerID peer.ID

func setup(t *testing.T) {
	peerSet := peerset.NewPeerSet(3 * time.Second)
	committee, _ := committee.GenerateTestCommittee()
	state := state.MockingState(committee)
	tFirewall = NewFirewall(peerSet, state)
	tAnotherPeerID = util.RandomPeerID()
}

func TestIncreaseMsgCounter(t *testing.T) {
	setup(t)

	msg := message.NewMessage(tAnotherPeerID, payload.NewQueryProposalPayload(1, 0))
	d, _ := msg.Encode()
	assert.NotNil(t, tFirewall.OpenMessage(d, tAnotherPeerID))
	p := tFirewall.peerSet.GetPeer(tAnotherPeerID)
	assert.False(t, tFirewall.badPeer(p))

	tFirewall.OpenMessage([]byte("bad"), tAnotherPeerID)
	p = tFirewall.peerSet.GetPeer(tAnotherPeerID)
	assert.True(t, tFirewall.badPeer(p))
	assert.Nil(t, tFirewall.OpenMessage(d, tAnotherPeerID))
}
