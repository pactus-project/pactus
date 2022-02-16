package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/consensus/vote"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/util"
)

func TestParsingQueryVotesMessages(t *testing.T) {
	setup(t)

	consensusHeight := tState.LastBlockHeight() + 1
	v1, _ := vote.GenerateTestPrecommitVote(consensusHeight, 0)
	tConsensus.AddVote(v1)
	pid := util.RandomPeerID()
	pld := payload.NewQueryVotesPayload(consensusHeight, 1)

	t.Run("Not in the committee, should not respond to the query vote message", func(t *testing.T) {
		assert.Error(t, testReceiveingNewMessage(t, tSync, pld, pid))
	})

	pub, _ := bls.GenerateTestKeyPair()
	testAddPeerToCommittee(t, pub, pid)

	t.Run("In the committee, should respond to the query vote message", func(t *testing.T) {
		assert.NoError(t, testReceiveingNewMessage(t, tSync, pld, pid))

		msg := shouldPublishPayloadWithThisType(t, tNetwork, payload.PayloadTypeVote)
		assert.Equal(t, msg.Payload.(*payload.VotePayload).Vote.Hash(), v1.Hash())
	})

	t.Run("In the committee, but doesn't have the vote", func(t *testing.T) {
		pld := payload.NewQueryVotesPayload(consensusHeight+1, 1)
		assert.NoError(t, testReceiveingNewMessage(t, tSync, pld, pid))

		shouldNotPublishPayloadWithThisType(t, tNetwork, payload.PayloadTypeVote)
	})
}

func TestBroadcastingQueryVotesMessages(t *testing.T) {
	setup(t)

	consensusHeight := tState.LastBlockHeight() + 1
	pld := payload.NewQueryVotesPayload(consensusHeight, 1)

	t.Run("Not in the committee, should not send query vote message", func(t *testing.T) {
		tSync.broadcast(pld)

		shouldNotPublishPayloadWithThisType(t, tNetwork, payload.PayloadTypeQueryVotes)
	})

	testAddPeerToCommittee(t, tSync.signer.PublicKey(), tSync.SelfID())
	t.Run("In the committee, should send query vote message", func(t *testing.T) {
		tSync.broadcast(pld)

		shouldPublishPayloadWithThisType(t, tNetwork, payload.PayloadTypeQueryVotes)
	})
}
