package sync

import (
	"testing"

	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/sync/services"
	"github.com/pactus-project/pactus/types/certificate"
	"github.com/pactus-project/pactus/util"
	"github.com/stretchr/testify/assert"
)

func TestParsingBlockAnnounceMessages(t *testing.T) {
	td := setup(t, nil)

	pid := td.RandPeerID()
	lastBlockHeight := td.state.LastBlockHeight()
	blockInterval := td.state.Genesis().Params().BlockInterval()
	blockTime := util.RoundNow(int(blockInterval.Seconds()))
	b1 := td.GenerateTestBlockWithTime(blockTime)
	c1 := td.GenerateTestCertificate()
	msg1 := message.NewBlockAnnounceMessage(lastBlockHeight+1, b1, c1)

	blockTime = blockTime.Add(blockInterval)
	b2 := td.GenerateTestBlockWithTime(blockTime)
	c2 := td.GenerateTestCertificate()
	msg2 := message.NewBlockAnnounceMessage(lastBlockHeight+2, b2, c2)

	pub, _ := td.RandBLSKeyPair()
	td.addPeer(t, pub, pid, services.New(services.Network))

	t.Run("Receiving new block announce message, without committing previous block", func(t *testing.T) {
		assert.NoError(t, td.receivingNewMessage(td.sync, msg2, pid))

		msg1 := td.shouldPublishMessageWithThisType(t, td.network, message.TypeBlocksRequest)
		assert.Equal(t, msg1.Message.(*message.BlocksRequestMessage).From, lastBlockHeight+1)

		peer := td.sync.peerSet.GetPeer(pid)
		assert.Equal(t, peer.Height, lastBlockHeight+2)
		assert.Equal(t, td.sync.state.LastBlockHeight(), lastBlockHeight)
	})

	t.Run("Receiving missed block, should commit both blocks", func(t *testing.T) {
		assert.NoError(t, td.receivingNewMessage(td.sync, msg1, pid))

		assert.Equal(t, td.sync.state.LastBlockHeight(), lastBlockHeight+2)
	})
}

func TestInvalidBlockAnnounce(t *testing.T) {
	td := setup(t, nil)

	pid := td.RandPeerID()
	lastBlockHeight := td.state.LastBlockHeight()
	blk := td.GenerateTestBlock()
	invCert := certificate.NewCertificate(td.RandHeight(), 0, nil, nil, nil)
	msg := message.NewBlockAnnounceMessage(lastBlockHeight+1, blk, invCert)

	err := td.receivingNewMessage(td.sync, msg, pid)
	assert.Error(t, err)
}

func TestBroadcastingBlockAnnounceMessages(t *testing.T) {
	td := setup(t, nil)

	td.state.CommitTestBlocks(21)
	blk, _ := td.state.CommittedBlock(td.state.LastBlockHeight()).ToBlock()
	msg := message.NewBlockAnnounceMessage(
		td.state.LastBlockHeight(), blk, td.state.LastCertificate())

	t.Run("Not in the committee, should not broadcast block announce message", func(t *testing.T) {
		td.sync.broadcast(msg)

		td.shouldNotPublishMessageWithThisType(t, td.network, message.TypeBlockAnnounce)
	})

	td.addPeerToCommittee(t, td.sync.SelfID(), td.sync.valKeys[0].PublicKey())

	t.Run("In the committee, should broadcast block announce message", func(t *testing.T) {
		td.sync.broadcast(msg)

		msg1 := td.shouldPublishMessageWithThisType(t, td.network, message.TypeBlockAnnounce)
		assert.Equal(t, msg1.Message.(*message.BlockAnnounceMessage).Height, msg.Height)
	})
}
