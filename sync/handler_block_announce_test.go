package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/util"
)

func TestParsingBlockAnnounceMessages(t *testing.T) {
	setup(t)

	lastBlockHash := tState.LastBlockHash()
	lastBlockheight := tState.LastBlockHeight()
	b1, _ := block.GenerateTestBlock(nil, &lastBlockHash)
	lastBlockHash = b1.Hash()
	b2, _ := block.GenerateTestBlock(nil, &lastBlockHash)
	c2 := block.GenerateTestCertificate(b2.Hash())

	pid := util.RandomPeerID()
	pld := payload.NewBlockAnnouncePayload(lastBlockheight+2, b2, c2)

	pub, _ := bls.GenerateTestKeyPair()
	testAddPeer(t, pub, pid)

	t.Run("Receiving new block announce message, without committing previous block", func(t *testing.T) {
		assert.NoError(t, testReceiveingNewMessage(tSync, pld, pid))

		msg1 := shouldPublishPayloadWithThisType(t, tNetwork, payload.PayloadTypeBlocksRequest)
		assert.Equal(t, msg1.Payload.(*payload.BlocksRequestPayload).From, lastBlockheight+1)
	})

}

func TestBroadcastingBlockAnnounceMessages(t *testing.T) {
	setup(t)

	pld := payload.NewBlockAnnouncePayload(
		tState.LastBlockHeight(),
		tState.Block(tState.LastBlockHeight()),
		tState.LastCertificate())

	t.Run("Not in the committee, should not broadcast block announce message", func(t *testing.T) {
		tSync.broadcast(pld)

		shouldNotPublishPayloadWithThisType(t, tNetwork, payload.PayloadTypeBlockAnnounce)
	})

	testAddPeerToCommittee(t, tSync.SelfID(), tSync.signer.PublicKey())

	t.Run("In the committee, should broadcast block announce message", func(t *testing.T) {
		tSync.broadcast(pld)

		msg1 := shouldPublishPayloadWithThisType(t, tNetwork, payload.PayloadTypeBlockAnnounce)
		assert.Equal(t, msg1.Payload.(*payload.BlockAnnouncePayload).Height, pld.Height)
	})
}
