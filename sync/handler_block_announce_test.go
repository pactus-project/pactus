package sync

import (
	"testing"

	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/services"
	"github.com/pactus-project/pactus/types/certificate"
	"github.com/stretchr/testify/assert"
)

func TestParsingBlockAnnounceMessages(t *testing.T) {
	td := setup(t, nil)

	pid := td.RandPeerID()
	lastHeight := td.state.LastBlockHeight()
	blk1, cert1 := td.GenerateTestBlock(lastHeight + 1)
	msg1 := message.NewBlockAnnounceMessage(blk1, cert1)

	blk2, cert2 := td.GenerateTestBlock(lastHeight + 2)
	msg2 := message.NewBlockAnnounceMessage(blk2, cert2)

	pub, _ := td.RandBLSKeyPair()
	td.addPeer(t, pub, pid, services.New(services.Network))

	t.Run("Receiving new block announce message, without committing previous block", func(t *testing.T) {
		assert.NoError(t, td.receivingNewMessage(td.sync, msg2, pid))

		msg1 := td.shouldPublishMessageWithThisType(t, td.network, message.TypeBlocksRequest)
		assert.Equal(t, msg1.Message.(*message.BlocksRequestMessage).From, lastHeight+1)

		peer := td.sync.peerSet.GetPeer(pid)
		assert.Equal(t, peer.Height, lastHeight+2)
		assert.Equal(t, td.sync.state.LastBlockHeight(), lastHeight)
	})

	t.Run("Receiving missed block, should commit both blocks", func(t *testing.T) {
		assert.NoError(t, td.receivingNewMessage(td.sync, msg1, pid))

		assert.Equal(t, td.sync.state.LastBlockHeight(), lastHeight+2)
	})
}

func TestInvalidBlockAnnounce(t *testing.T) {
	td := setup(t, nil)

	pid := td.RandPeerID()
	lastHeight := td.state.LastBlockHeight()
	blk, _ := td.GenerateTestBlock(lastHeight + 1)
	invCert := certificate.NewCertificate(lastHeight+1, 0, nil, nil, nil)
	msg := message.NewBlockAnnounceMessage(blk, invCert)

	err := td.receivingNewMessage(td.sync, msg, pid)
	assert.Error(t, err)
}

func TestBroadcastingBlockAnnounceMessages(t *testing.T) {
	td := setup(t, nil)

	td.state.CommitTestBlocks(21)
	blk, _ := td.state.CommittedBlock(td.state.LastBlockHeight()).ToBlock()
	msg := message.NewBlockAnnounceMessage(
		blk, td.state.LastCertificate())

	t.Run("Not in the committee, should not broadcast block announce message", func(t *testing.T) {
		td.sync.broadcast(msg)

		td.shouldNotPublishMessageWithThisType(t, td.network, message.TypeBlockAnnounce)
	})

	td.addPeerToCommittee(t, td.sync.SelfID(), td.sync.valKeys[0].PublicKey())

	t.Run("In the committee, should broadcast block announce message", func(t *testing.T) {
		td.sync.broadcast(msg)

		msg1 := td.shouldPublishMessageWithThisType(t, td.network, message.TypeBlockAnnounce)
		assert.Equal(t, msg1.Message.(*message.BlockAnnounceMessage).Certificate.Height(), msg.Certificate.Height())
	})
}
