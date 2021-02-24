package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/sync/message"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/vote"
)

func TestProposalToCache(t *testing.T) {
	setup(t)

	p, _ := vote.GenerateTestProposal(106, 0)

	tAliceSync.consensusSync.BroadcastProposal(p)
	tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeProposal)
	assert.NotNil(t, tBobSync.cache.GetProposal(p.Height(), p.Round()))
}

func TestRequestForProposal(t *testing.T) {
	setup(t)

	joinAliceToTheSet(t)
	joinBobToTheSet(t)

	hrs := tAliceConsensus.HRS()
	assert.Equal(t, hrs.Height(), tAliceState.LastBlockHeight()+1)

	t.Run("Alice and bob are in same height. Alice doesn't have have proposal. Bob ask for the proposal", func(t *testing.T) {
		tBobBroadcastCh <- message.NewQueryProposalMessage(tBobPeerID, hrs.Height(), 0)
		tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeQueryProposal)

		// Alice doesn't respond
		tAliceNetAPI.ShouldNotPublishMessageWithThisType(t, payload.PayloadTypeProposal)
	})

	p1, _ := vote.GenerateTestProposal(hrs.Height(), 0)
	tAliceConsensus.SetProposal(p1)

	t.Run("Alice and bob are in same height. Alice has proposal. Bob ask for the proposal", func(t *testing.T) {
		tBobBroadcastCh <- message.NewQueryProposalMessage(tBobPeerID, hrs.Height(), 0)

		tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeQueryProposal)
		tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeProposal)

		assert.Equal(t, tBobConsensus.Proposal.Hash(), p1.Hash())
	})

	t.Run("Alice and bob are in same height. Bob is in next round. Bob ask for the proposal", func(t *testing.T) {
		tBobBroadcastCh <- message.NewQueryProposalMessage(tBobPeerID, hrs.Height(), 1)
		tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeQueryProposal)

		// Alice doesn't respond
		tAliceNetAPI.ShouldNotPublishMessageWithThisType(t, payload.PayloadTypeProposal)
	})

	p2, _ := vote.GenerateTestProposal(hrs.Height(), 1)
	tAliceConsensus.Proposal = p2
	tAliceConsensus.Round = 1

	t.Run("Alice and bob are in same height. Alice is in next round. Alice has proposal. Bob ask for the proposal", func(t *testing.T) {
		tBobBroadcastCh <- message.NewQueryProposalMessage(tBobPeerID, hrs.Height(), 1)

		tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeQueryProposal)
		tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeProposal)

		assert.Equal(t, tBobConsensus.Proposal.Hash(), p2.Hash())
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

func TestProcessQueryVote(t *testing.T) {
	setup(t)

	disableHeartbeat(t)
	joinAliceToTheSet(t)
	joinBobToTheSet(t)

	hrs := tAliceConsensus.HRS()
	v1, _ := vote.GenerateTestPrepareVote(hrs.Height(), 0)
	v2, _ := vote.GenerateTestPrepareVote(hrs.Height(), 1)
	tAliceConsensus.Votes = []*vote.Vote{v1, v2}

	t.Run("Alice and bob are in same height. Bob queries for votes, alice sends a random vote", func(t *testing.T) {
		tBobSync.consensusSync.BroadcastQueryVotes(hrs.Height(), 1)
		tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeQueryVotes)
		tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeVote)
	})
}

func TestProcessHeartbeatForQueryProposal(t *testing.T) {
	setup(t)

	joinAliceToTheSet(t)
	joinBobToTheSet(t)

	hrs := tAliceConsensus.HRS()

	p, _ := vote.GenerateTestProposal(hrs.Height(), hrs.Round())
	tBobConsensus.SetProposal(p)

	t.Run("Alice Doesn't have proposal. She should query it.", func(t *testing.T) {
		tAliceSync.broadcastHeartBeat()
		tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeQueryProposal)
		tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeProposal)
	})

	v1, _ := vote.GenerateTestPrepareVote(hrs.Height(), 0)
	v2, _ := vote.GenerateTestPrepareVote(hrs.Height(), 1)
	tAliceConsensus.Votes = []*vote.Vote{v1, v2}

	t.Run("Alice and bob are in same HRS.", func(t *testing.T) {
		tAliceSync.broadcastHeartBeat()
		tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeVote)
		tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeHeartBeat)
	})

	tAliceConsensus.Round = 1
	t.Run("Alice is in the next round. Bob isn't", func(t *testing.T) {
		tAliceSync.broadcastHeartBeat()
		tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeVote)
		tAliceNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeHeartBeat)

		tBobNetAPI.ShouldPublishMessageWithThisType(t, payload.PayloadTypeQueryVotes)
	})
}
