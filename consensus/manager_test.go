package consensus

import (
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManager(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	state := state.MockingState(ts)
	state.TestCommittee.Validators()

	rewardAddrs := []crypto.Address{ts.RandAccAddress(), ts.RandAccAddress()}
	valKeys := []*bls.ValidatorKey{state.TestValKeys[0], ts.RandValKey()}
	broadcastCh := make(chan message.Message, 500)

	stateHeight := ts.RandHeight()
	blk, cert := ts.GenerateTestBlock(stateHeight)
	state.TestStore.SaveBlock(blk, cert)

	Mgr := NewManager(testConfig(), state, valKeys, rewardAddrs, broadcastCh)
	mgr := Mgr.(*manager)

	consA := mgr.instances[0].(*consensus) // active
	consB := mgr.instances[1].(*consensus) // inactive

	t.Run("Check if keys are assigned properly", func(t *testing.T) {
		assert.Equal(t, valKeys[0].PublicKey(), consA.ConsensusKey())
		assert.Equal(t, valKeys[1].PublicKey(), consB.ConsensusKey())
	})

	t.Run("Check if all instances move to new height", func(t *testing.T) {
		assert.False(t, mgr.HasActiveInstance())
		mgr.MoveToNewHeight()
		h, r := mgr.HeightRound()

		assert.True(t, mgr.HasActiveInstance())
		assert.Equal(t, stateHeight+1, h)
		assert.Zero(t, r)
	})

	t.Run("Testing add vote", func(t *testing.T) {
		v := vote.NewPrepareVote(ts.RandHash(), stateHeight+1, 0, valKeys[0].Address())
		ts.HelperSignVote(valKeys[0], v)

		mgr.AddVote(v)

		assert.True(t, consA.HasVote(v.Hash()))
		assert.False(t, consB.HasVote(v.Hash()))
	})

	t.Run("Testing set proposal", func(t *testing.T) {
		b, _ := state.ProposeBlock(valKeys[0], valKeys[0].Address(), 1)
		p := proposal.NewProposal(stateHeight+1, 0, b)
		ts.HelperSignProposal(valKeys[0], p)

		mgr.SetProposal(p)

		// assert.Equal(t, p, mgr.Proposal())
		assert.Equal(t, p, consA.Proposal())
		assert.Nil(t, consB.Proposal())
	})

	t.Run("Check discarding old votes", func(t *testing.T) {
		v := vote.NewPrepareVote(ts.RandHash(), stateHeight-1, 0, state.TestValKeys[2].Address())
		ts.HelperSignVote(state.TestValKeys[2], v)

		mgr.AddVote(v)
		assert.Empty(t, mgr.upcomingVotes)
	})

	t.Run("Check discarding old proposals", func(t *testing.T) {
		b, _ := state.ProposeBlock(valKeys[0], valKeys[0].Address(), 1)
		p := proposal.NewProposal(stateHeight-1, 1, b)
		ts.HelperSignProposal(valKeys[0], p)

		mgr.SetProposal(p)
		assert.Empty(t, mgr.upcomingProposals)
	})

	t.Run("Processing upcoming votes", func(t *testing.T) {
		v1 := vote.NewPrepareVote(ts.RandHash(), stateHeight+2, 0, valKeys[0].Address())
		v2 := vote.NewPrepareVote(ts.RandHash(), stateHeight+3, 0, valKeys[0].Address())
		v3 := vote.NewPrepareVote(ts.RandHash(), stateHeight+4, 0, valKeys[0].Address())

		ts.HelperSignVote(valKeys[0], v1)
		ts.HelperSignVote(valKeys[0], v2)
		ts.HelperSignVote(valKeys[0], v3)

		mgr.AddVote(v1)
		mgr.AddVote(v2)
		mgr.AddVote(v3)

		assert.Len(t, mgr.upcomingVotes, 3)

		blk, cert := ts.GenerateTestBlock(stateHeight + 1)
		err := state.CommitBlock(blk, cert)
		assert.NoError(t, err)
		stateHeight++

		blk, cert = ts.GenerateTestBlock(stateHeight + 1)
		err = state.CommitBlock(blk, cert)
		assert.NoError(t, err)
		stateHeight++

		mgr.MoveToNewHeight()

		assert.Len(t, mgr.upcomingVotes, 1)
	})

	t.Run("Processing upcoming proposal", func(t *testing.T) {
		b1, _ := state.ProposeBlock(valKeys[0], valKeys[0].Address(), 1)
		p1 := proposal.NewProposal(stateHeight+2, 0, b1)

		b2, _ := state.ProposeBlock(valKeys[0], valKeys[0].Address(), 1)
		p2 := proposal.NewProposal(stateHeight+3, 0, b2)

		b3, _ := state.ProposeBlock(valKeys[0], valKeys[0].Address(), 1)
		p3 := proposal.NewProposal(stateHeight+4, 0, b3)

		ts.HelperSignProposal(valKeys[0], p1)
		ts.HelperSignProposal(valKeys[0], p2)
		ts.HelperSignProposal(valKeys[0], p3)

		mgr.SetProposal(p1)
		mgr.SetProposal(p2)
		mgr.SetProposal(p3)

		assert.Len(t, mgr.upcomingProposals, 3)

		blk, cert := ts.GenerateTestBlock(stateHeight + 1)
		err := state.CommitBlock(blk, cert)
		assert.NoError(t, err)
		stateHeight++

		blk, cert = ts.GenerateTestBlock(stateHeight + 1)
		err = state.CommitBlock(blk, cert)
		assert.NoError(t, err)
		stateHeight++

		mgr.MoveToNewHeight()

		assert.Len(t, mgr.upcomingProposals, 1)
	})
}

func TestMediator(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	state := state.MockingState(ts)
	committee, valKeys := ts.GenerateTestCommittee(4)
	state.TestCommittee = committee
	state.TestParams.BlockIntervalInSecond = 1

	rewardAddrs := []crypto.Address{
		ts.RandAccAddress(), ts.RandAccAddress(),
		ts.RandAccAddress(), ts.RandAccAddress(),
	}
	broadcastCh := make(chan message.Message, 500)

	stateHeight := ts.RandHeight()
	blk, cert := ts.GenerateTestBlock(stateHeight)
	state.TestStore.SaveBlock(blk, cert)

	Mgr := NewManager(testConfig(), state, valKeys, rewardAddrs, broadcastCh)
	mgr := Mgr.(*manager)

	mgr.MoveToNewHeight()

	for {
		msg := <-broadcastCh
		logger.Info("shouldPublishProposal", "message", msg)

		m, ok := msg.(*message.BlockAnnounceMessage)
		if ok {
			require.Equal(t, m.Height(), stateHeight+1)
			return
		}
	}
}
