package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/message"
	"github.com/zarbchain/zarb-go/message/payload"
	"github.com/zarbchain/zarb-go/version"
)

func TestSendSalamBadGenesisHash(t *testing.T) {
	setup(t)

	invGenHash := crypto.GenerateTestHash()
	_, pub, _ := crypto.GenerateTestKeyPair()

	msg := message.NewSalamMessage("bad-genesis", pub, tAnotherPeerID, invGenHash, 0, 0)
	d, _ := msg.Encode()
	tAliceNetAPI.CheckAndParsMessage(d, tAnotherPeerID)
	msg2 := tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeAleyk)
	pld := msg2.Payload.(*payload.AleykPayload)

	assert.Equal(t, pld.Response.Status, payload.SalamResponseCodeRejected)
}

func TestSendSalamPeerBehind(t *testing.T) {
	setup(t)
	_, pub, _ := crypto.GenerateTestKeyPair()

	msg := message.NewSalamMessage("kitty", pub, tAnotherPeerID, tAliceState.GenHash, 3, 0)
	d, _ := msg.Encode()
	tAliceNetAPI.CheckAndParsMessage(d, tAnotherPeerID)
	msg2 := tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeAleyk)
	pld := msg2.Payload.(*payload.AleykPayload)

	assert.Equal(t, pld.Response.Status, payload.SalamResponseCodeOK)
	assert.Equal(t, tBobSync.peerSet.MaxClaimedHeight(), tAliceState.LastBlockHeight())

	p := tAliceSync.peerSet.GetPeer(tAnotherPeerID)
	assert.Equal(t, p.Version(), version.NodeVersion)
	assert.Equal(t, p.Moniker(), "kitty")
	assert.True(t, pub.EqualsTo(p.PublicKey()))
	assert.Equal(t, p.PeerID(), tAnotherPeerID)
	assert.Equal(t, p.Address(), pub.Address())
	assert.Equal(t, p.Height(), 3)
}

func TestSendSalamPeerAhead(t *testing.T) {
	setup(t)

	_, pub, _ := crypto.GenerateTestKeyPair()

	msg := message.NewSalamMessage("kitty", pub, tAnotherPeerID, tAliceState.GenHash, 111, 0)
	d, _ := msg.Encode()
	tAliceNetAPI.CheckAndParsMessage(d, tAnotherPeerID)
	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeAleyk)
	tAliceNetAPI.ShouldPublishThisMessage(t, message.NewLatestBlocksRequestMessage(tAliceState.LastBlockHeight()+1, tAliceState.LastBlockHash()))

	assert.Equal(t, tAliceSync.peerSet.MaxClaimedHeight(), 111)
}

func TestSendAleykPeerBehind(t *testing.T) {
	setup(t)
	_, pub, _ := crypto.GenerateTestKeyPair()

	msg := message.NewAleykMessage("kitty", pub, tAnotherPeerID, 1, 0, 0, "Welcome!")
	d, _ := msg.Encode()
	tAliceNetAPI.CheckAndParsMessage(d, tAnotherPeerID)
	tAliceNetAPI.ShouldNotPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocksRequest)
}

func TestSendAleykPeerAhead(t *testing.T) {
	setup(t)
	_, pub, _ := crypto.GenerateTestKeyPair()

	msg := message.NewAleykMessage("kitty", pub, tAnotherPeerID, 111, 0, 0, "Welcome!")
	d, _ := msg.Encode()
	tAliceNetAPI.CheckAndParsMessage(d, tAnotherPeerID)
	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocksRequest)
	assert.Equal(t, tAliceSync.peerSet.MaxClaimedHeight(), 111)
}

func TestSendAleykPeerSameHeight(t *testing.T) {
	setup(t)
	_, pub, _ := crypto.GenerateTestKeyPair()

	msg := message.NewAleykMessage("kitty", pub, tAnotherPeerID, tAliceState.LastBlockHeight(), 0, 0, "Welcome!")
	d, _ := msg.Encode()
	tAliceNetAPI.CheckAndParsMessage(d, tAnotherPeerID)
	tAliceNetAPI.ShouldNotPublishMessageWithThisType(t, payload.PayloadTypeLatestBlocksRequest)
}
