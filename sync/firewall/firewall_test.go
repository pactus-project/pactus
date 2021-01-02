package firewall

import (
	"testing"

	"github.com/zarbchain/zarb-go/vote"

	peer "github.com/libp2p/go-libp2p-peer"
	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/message"
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

	msg := message.NewProposalRequestMessage(1, 1)
	d, _ := msg.Encode()
	assert.Nil(t, tFirewall.ParsMessage(d, tAnotherPeerID))
}

func TestIncreaseHeight(t *testing.T) {
	setup(t)
	_, pub, _ := crypto.GenerateTestKeyPair()

	msg1 := message.NewSalamMessage("kitty", pub, tAnotherPeerID, crypto.GenerateTestHash(), 3, 0)
	d, _ := msg1.Encode()
	assert.NotNil(t, tFirewall.ParsMessage(d, tAnotherPeerID))
	assert.Equal(t, tFirewall.peerSet.MaxClaimedHeight(), 3)

	msg2 := message.NewAleykMessage("kitty-2", pub, tAnotherPeerID, 4, 0, 0, "Welcome!")
	d, _ = msg2.Encode()
	assert.NotNil(t, tFirewall.ParsMessage(d, tAnotherPeerID))
	assert.Equal(t, tFirewall.peerSet.MaxClaimedHeight(), 4)

	msg3 := message.NewHeartBeatMessage(crypto.GenerateTestHash(), hrs.NewHRS(6, 0, 1))
	d, _ = msg3.Encode()
	assert.NotNil(t, tFirewall.ParsMessage(d, tAnotherPeerID))
	assert.Equal(t, tFirewall.peerSet.MaxClaimedHeight(), 5)

	p, _ := vote.GenerateTestProposal(7, 1)
	msg4 := message.NewProposalMessage(p)
	d, _ = msg4.Encode()
	assert.NotNil(t, tFirewall.ParsMessage(d, tAnotherPeerID))
	assert.Equal(t, tFirewall.peerSet.MaxClaimedHeight(), 6)

	v, _ := vote.GenerateTestPrepareVote(8, 0)
	msg5 := message.NewVoteMessage(v)
	d, _ = msg5.Encode()
	assert.NotNil(t, tFirewall.ParsMessage(d, tAnotherPeerID))
	assert.Equal(t, tFirewall.peerSet.MaxClaimedHeight(), 7)
}
