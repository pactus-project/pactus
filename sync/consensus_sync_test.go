package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/sync/message"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/vote"
)

func TestRequestForProposal(t *testing.T) {

	t.Run("Alice and bob are in same height. Alice has proposal. Bob ask for the proposal", func(t *testing.T) {
		setup(t)

		hrs := hrs.NewHRS(100, 1, 6)
		p, _ := vote.GenerateTestProposal(hrs.Height(), hrs.Round())
		tAliceConsensus.SetProposal(p)
		tAliceConsensus.HRS_ = hrs

		tBobBroadcastCh <- message.NewQueryProposalMessage(hrs.Height(), hrs.Round())
		tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeQueryProposal)
		tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeProposal)

		assert.Equal(t, tBobConsensus.Proposal.Hash(), tBobConsensus.Proposal.Hash())
	})

	t.Run("Alice and bob are in same height. Alice doesn't have have proposal. Bob ask for the proposal", func(t *testing.T) {
		setup(t)

		hrs := hrs.NewHRS(101, 1, 6)
		tAliceConsensus.HRS_ = hrs

		tBobBroadcastCh <- message.NewQueryProposalMessage(hrs.Height(), hrs.Round())
		tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeQueryProposal)

		// Alice doesn't respond
		tAliceNetAPI.ShouldNotPublishMessageWithThisType(t, payload.PayloadTypeProposal)
	})

	t.Run("Alice and bob are in same height. Alice is in next round. Alice has proposal. Bob ask for the proposal", func(t *testing.T) {
		setup(t)

		hrs := hrs.NewHRS(102, 1, 6)
		p, _ := vote.GenerateTestProposal(hrs.Height(), hrs.Round())
		tAliceConsensus.SetProposal(p)
		tAliceConsensus.HRS_ = hrs

		tBobBroadcastCh <- message.NewQueryProposalMessage(hrs.Height(), hrs.Round()-1)
		tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeQueryProposal)
		tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeProposal)

		assert.Equal(t, tBobConsensus.Proposal.Hash(), tBobConsensus.Proposal.Hash())
	})

	t.Run("Alice and bob are in same height. Alice is in previous round. Alice has proposal. Bob ask for the proposal", func(t *testing.T) {
		setup(t)

		hrs := hrs.NewHRS(103, 1, 6)
		p, _ := vote.GenerateTestProposal(hrs.Height(), hrs.Round())
		tAliceConsensus.SetProposal(p)
		tAliceConsensus.HRS_ = hrs

		tBobBroadcastCh <- message.NewQueryProposalMessage(hrs.Height(), hrs.Round()+1)
		tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeQueryProposal)

		// Alice doesn't respond
		tAliceNetAPI.ShouldNotPublishMessageWithThisType(t, payload.PayloadTypeProposal)
	})

}

func TestUpdateConsensus(t *testing.T) {
	setup(t)

	v, _ := vote.GenerateTestPrecommitVote(1, 1)
	p, _ := vote.GenerateTestProposal(1, 1)

	tAliceSync.consensusSync.BroadcastVote(v)
	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeVote)

	tAliceSync.consensusSync.BroadcastProposal(p)
	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeProposal)

	assert.Equal(t, tBobConsensus.Votes[0].Hash(), v.Hash())
	assert.Equal(t, tBobConsensus.Proposal.Hash(), p.Hash())
}

func TestSendVoteSet(t *testing.T) {
	setup(t)

	tAliceConsensus.HRS_ = hrs.NewHRS(100, 1, 1)
	tBobConsensus.HRS_ = hrs.NewHRS(100, 1, 1)
	v1, _ := vote.GenerateTestPrepareVote(100, 0)
	v2, _ := vote.GenerateTestPrepareVote(100, 1)
	v3, _ := vote.GenerateTestPrepareVote(100, 1)
	v4, _ := vote.GenerateTestPrepareVote(100, 1)
	v5, _ := vote.GenerateTestPrepareVote(101, 1)

	tAliceConsensus.Votes = []*vote.Vote{v1, v2, v3}
	tBobConsensus.Votes = []*vote.Vote{v2, v4}
	tBobSync.consensusSync.BroadcastVoteSet()
	tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeVoteSet)
	tAliceNetAPI.ShouldPublishThisMessage(t, message.NewVoteMessage(v3))

	tBobBroadcastCh <- message.NewVoteSetMessage(101, 1, []crypto.Hash{v5.Hash()})
	tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeVoteSet)
	tAliceNetAPI.ShouldNotPublishMessageWithThisType(t, payload.PayloadTypeVote)
}
