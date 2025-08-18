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
	"github.com/pactus-project/pactus/util/pipeline"
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
	pipe := pipeline.MockingPipeline[message.Message]()

	randomHeight := ts.RandHeight()
	rndBlk, rndCert := ts.GenerateTestBlock(randomHeight)
	state.TestStore.SaveBlock(rndBlk, rndCert)

	mgrInt := NewManager(testConfig(), state, valKeys, rewardAddrs, pipe)
	mgr := mgrInt.(*manager)

	consA := mgr.instances[0].(*consensus) // active
	consB := mgr.instances[1].(*consensus) // inactive

	t.Run("Check if keys are assigned properly", func(t *testing.T) {
		assert.Equal(t, consA.ConsensusKey(), valKeys[0].PublicKey())
		assert.Equal(t, consB.ConsensusKey(), valKeys[1].PublicKey())
	})

	t.Run("Check if all instances move to new height", func(t *testing.T) {
		mgr.MoveToNewHeight()

		stateHeight := mgr.state.LastBlockHeight()
		consHeight, consRound := mgr.HeightRound()

		assert.True(t, mgr.HasActiveInstance())
		assert.Equal(t, stateHeight+1, consHeight)
		assert.Zero(t, consRound)
	})

	t.Run("Testing add vote", func(t *testing.T) {
		consHeight, _ := mgr.HeightRound()
		vote := vote.NewPrepareVote(ts.RandHash(), consHeight, 0, valKeys[0].Address())
		ts.HelperSignVote(valKeys[0], vote)

		mgr.AddVote(vote)

		assert.True(t, consA.HasVote(vote.Hash()))
		assert.False(t, consB.HasVote(vote.Hash()))
	})

	t.Run("Testing set proposal", func(t *testing.T) {
		consHeight, _ := mgr.HeightRound()
		blk, _ := state.ProposeBlock(valKeys[0], valKeys[0].Address())
		prop := proposal.NewProposal(consHeight, 0, blk)
		ts.HelperSignProposal(valKeys[0], prop)

		mgr.SetProposal(prop)

		assert.Equal(t, prop, consA.Proposal())
		assert.Nil(t, consB.Proposal())
	})

	t.Run("Check discarding old votes", func(t *testing.T) {
		consHeight, _ := mgr.HeightRound()
		v := vote.NewPrepareVote(ts.RandHash(), consHeight-1, 0, state.TestValKeys[2].Address())
		ts.HelperSignVote(state.TestValKeys[2], v)

		mgr.AddVote(v)
		assert.Empty(t, mgr.upcomingVotes)
	})

	t.Run("Check discarding old proposals", func(t *testing.T) {
		consHeight, _ := mgr.HeightRound()
		blk, _ := state.ProposeBlock(valKeys[0], valKeys[0].Address())
		prop := proposal.NewProposal(consHeight-1, 1, blk)
		ts.HelperSignProposal(valKeys[0], prop)

		mgr.SetProposal(prop)
		assert.Empty(t, mgr.upcomingProposals)
	})

	t.Run("Processing upcoming votes", func(t *testing.T) {
		consHeight, _ := mgr.HeightRound()
		vote1 := vote.NewPrepareVote(ts.RandHash(), consHeight+1, 0, valKeys[0].Address())
		vote2 := vote.NewPrepareVote(ts.RandHash(), consHeight+2, 0, valKeys[0].Address())
		vote3 := vote.NewPrepareVote(ts.RandHash(), consHeight+3, 0, valKeys[0].Address())

		ts.HelperSignVote(valKeys[0], vote1)
		ts.HelperSignVote(valKeys[0], vote2)
		ts.HelperSignVote(valKeys[0], vote3)

		mgr.AddVote(vote1)
		mgr.AddVote(vote2)
		mgr.AddVote(vote3)

		assert.Len(t, mgr.upcomingVotes, 3)

		blk1, cert1 := ts.GenerateTestBlock(consHeight)
		err := state.CommitBlock(blk1, cert1)
		assert.NoError(t, err)

		blk2, cert2 := ts.GenerateTestBlock(consHeight + 1)
		err = state.CommitBlock(blk2, cert2)
		assert.NoError(t, err)

		mgr.MoveToNewHeight()

		assert.Len(t, mgr.upcomingVotes, 1)
	})

	t.Run("Processing upcoming proposal", func(t *testing.T) {
		consHeight, _ := mgr.HeightRound()
		prop1 := ts.GenerateTestProposal(consHeight+1, 0)
		prop2 := ts.GenerateTestProposal(consHeight+2, 0)
		prop3 := ts.GenerateTestProposal(consHeight+3, 0)

		mgr.SetProposal(prop1)
		mgr.SetProposal(prop2)
		mgr.SetProposal(prop3)

		assert.Len(t, mgr.upcomingProposals, 3)

		blk1, cert1 := ts.GenerateTestBlock(consHeight)
		err := state.CommitBlock(blk1, cert1)
		assert.NoError(t, err)

		blk2, cert2 := ts.GenerateTestBlock(consHeight + 1)
		err = state.CommitBlock(blk2, cert2)
		assert.NoError(t, err)

		mgr.MoveToNewHeight()

		assert.Len(t, mgr.upcomingProposals, 1)
	})
}

func TestMediator(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	state := state.MockingState(ts)
	cmt, valKeys := ts.GenerateTestCommittee(4)
	state.TestCommittee = cmt
	state.TestParams.BlockIntervalInSecond = 1

	rewardAddrs := []crypto.Address{
		ts.RandAccAddress(), ts.RandAccAddress(),
		ts.RandAccAddress(), ts.RandAccAddress(),
	}
	stateHeight := ts.RandHeight()
	blk, cert := ts.GenerateTestBlock(stateHeight)
	state.TestStore.SaveBlock(blk, cert)
	pipe := pipeline.MockingPipeline[message.Message]()
	mgrInt := NewManager(testConfig(), state, valKeys, rewardAddrs, pipe)
	mgr := mgrInt.(*manager)

	mgr.MoveToNewHeight()

	for {
		msg := <-pipe.UnsafeGetChannel()
		logger.Info("shouldPublishProposal", "msg", msg)

		m, ok := msg.(*message.BlockAnnounceMessage)
		if ok {
			require.Equal(t, stateHeight+1, m.Height())

			return
		}
	}
}
