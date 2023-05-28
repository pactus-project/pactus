package sync

import (
	"testing"

	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/network"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/types/block"
	"github.com/stretchr/testify/assert"
)

func TestParsingBlockAnnounceMessages(t *testing.T) {
	setup(t)

	lastBlockHeight := tState.LastBlockHeight()
	b1 := block.GenerateTestBlock(nil, nil)
	c1 := block.GenerateTestCertificate(b1.Hash())
	b2 := block.GenerateTestBlock(nil, nil)
	c2 := block.GenerateTestCertificate(b2.Hash())

	pid := network.TestRandomPeerID()
	msg1 := message.NewBlockAnnounceMessage(lastBlockHeight+1, b1, c1)
	msg2 := message.NewBlockAnnounceMessage(lastBlockHeight+2, b2, c2)

	pub, _ := bls.GenerateTestKeyPair()
	testAddPeer(t, pub, pid, false)

	t.Run("Receiving new block announce message, without committing previous block", func(t *testing.T) {
		assert.NoError(t, testReceivingNewMessage(tSync, msg2, pid))

		msg1 := shouldPublishMessageWithThisType(t, tNetwork, message.MessageTypeBlocksRequest)
		assert.Equal(t, msg1.Message.(*message.BlocksRequestMessage).From, lastBlockHeight+1)

		peer := tSync.peerSet.GetPeer(pid)
		assert.Equal(t, peer.Height, lastBlockHeight+2)
		assert.Equal(t, tSync.state.LastBlockHeight(), lastBlockHeight)
		assert.Equal(t, tSync.peerSet.MaxClaimedHeight(), lastBlockHeight+2)
	})

	t.Run("Receiving missed block, should commit both blocks", func(t *testing.T) {
		assert.NoError(t, testReceivingNewMessage(tSync, msg1, pid))

		peer := tSync.peerSet.GetPeer(pid)
		assert.Equal(t, peer.Height, lastBlockHeight+2)
		assert.Equal(t, tSync.state.LastBlockHeight(), lastBlockHeight+2)
		assert.Equal(t, tSync.peerSet.MaxClaimedHeight(), lastBlockHeight+2)
	})
}

func TestBroadcastingBlockAnnounceMessages(t *testing.T) {
	setup(t)

	msg := message.NewBlockAnnounceMessage(
		tState.LastBlockHeight(),
		tState.StoredBlock(tState.LastBlockHeight()).ToBlock(),
		tState.LastCertificate())

	t.Run("Not in the committee, should not broadcast block announce message", func(t *testing.T) {
		tSync.broadcast(msg)

		shouldNotPublishMessageWithThisType(t, tNetwork, message.MessageTypeBlockAnnounce)
	})

	testAddPeerToCommittee(t, tSync.SelfID(), tSync.signers[0].PublicKey())

	t.Run("In the committee, should broadcast block announce message", func(t *testing.T) {
		tSync.broadcast(msg)

		msg1 := shouldPublishMessageWithThisType(t, tNetwork, message.MessageTypeBlockAnnounce)
		assert.Equal(t, msg1.Message.(*message.BlockAnnounceMessage).Height, msg.Height)
	})
}
