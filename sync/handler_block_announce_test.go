package sync

import (
	"testing"

	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/stretchr/testify/assert"
)

func TestParsingBlockAnnounceMessages(t *testing.T) {
	td := setup(t, nil)

	lastBlockHeight := td.state.LastBlockHeight()
	b1 := td.GenerateTestBlock(nil, nil)
	c1 := td.GenerateTestCertificate(b1.Hash())
	b2 := td.GenerateTestBlock(nil, nil)
	c2 := td.GenerateTestCertificate(b2.Hash())

	pid := td.RandomPeerID()
	msg1 := message.NewBlockAnnounceMessage(lastBlockHeight+1, b1, c1)
	msg2 := message.NewBlockAnnounceMessage(lastBlockHeight+2, b2, c2)

	pub, _ := td.RandomBLSKeyPair()
	td.addPeer(t, pub, pid, false)

	t.Run("Receiving new block announce message, without committing previous block", func(t *testing.T) {
		assert.NoError(t, td.receivingNewMessage(td.sync, msg2, pid))

		msg1 := td.shouldPublishMessageWithThisType(t, td.network, message.TypeBlocksRequest)
		assert.Equal(t, msg1.Message.(*message.BlocksRequestMessage).From, lastBlockHeight+1)

		peer := td.sync.peerSet.GetPeer(pid)
		assert.Equal(t, peer.Height, lastBlockHeight+2)
		assert.Equal(t, td.sync.state.LastBlockHeight(), lastBlockHeight)
		assert.Equal(t, td.sync.peerSet.MaxClaimedHeight(), lastBlockHeight+2)
	})

	t.Run("Receiving missed block, should commit both blocks", func(t *testing.T) {
		assert.NoError(t, td.receivingNewMessage(td.sync, msg1, pid))

		assert.Equal(t, td.sync.state.LastBlockHeight(), lastBlockHeight+2)
		assert.Equal(t, td.sync.peerSet.MaxClaimedHeight(), lastBlockHeight+2)
	})
}

func TestBroadcastingBlockAnnounceMessages(t *testing.T) {
	td := setup(t, nil)

	msg := message.NewBlockAnnounceMessage(
		td.state.LastBlockHeight(),
		td.state.StoredBlock(td.state.LastBlockHeight()).ToBlock(),
		td.state.LastCertificate())

	t.Run("Not in the committee, should not broadcast block announce message", func(t *testing.T) {
		td.sync.broadcast(msg)

		td.shouldNotPublishMessageWithThisType(t, td.network, message.TypeBlockAnnounce)
	})

	td.addPeerToCommittee(t, td.sync.SelfID(), td.sync.signers[0].PublicKey())

	t.Run("In the committee, should broadcast block announce message", func(t *testing.T) {
		td.sync.broadcast(msg)

		msg1 := td.shouldPublishMessageWithThisType(t, td.network, message.TypeBlockAnnounce)
		assert.Equal(t, msg1.Message.(*message.BlockAnnounceMessage).Height, msg.Height)
	})
}
