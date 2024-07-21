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

	st := state.MockingState(ts)
	st.TestCommittee.Validators()

	rewardAddrs := []crypto.Address{ts.RandAccAddress(), ts.RandAccAddress()}
	valKeys := []*bls.ValidatorKey{st.TestValKeys[0], ts.RandValKey()}
	broadcastCh := make(chan message.Message, 500)

	randomHeight := ts.RandHeight()
	blk, cert := ts.GenerateTestBlock(randomHeight)
	st.TestStore.SaveBlock(blk, cert)

	mgrInst := NewManager(testConfig(), st, valKeys, rewardAddrs, broadcastCh)
	mgr := mgrInst.(*manager)

	consA := mgr.instances[0].(*consensus) // active
	consB := mgr.instances[1].(*consensus) // inactive

	t.Run("Check if keys are assigned properly", func(t *testing.T) {
		assert.Equal(t, consA.ConsensusKey(), valKeys[0].PublicKey())
		assert.Equal(t, consB.ConsensusKey(), valKeys[1].PublicKey())
	})

	t.Run("Check if all instances move to new height", func(t *testing.T) {
		stateHeight := mgr.state.LastBlockHeight()
		assert.False(t, mgr.HasActiveInstance())

		mgr.MoveToNewHeight()
		consHeight, consRound := mgr.HeightRound()

		assert.True(t, mgr.HasActiveInstance())
		assert.Equal(t, stateHeight+1, consHeight)
		assert.Zero(t, consRound)
	})

	t.Run("Testing add vote", func(t *testing.T) {
		consHeight, _ := mgr.HeightRound()
		v := vote.NewPrepareVote(ts.RandHash(), consHeight, 0, valKeys[0].Address())
		ts.HelperSignVote(valKeys[0], v)

		mgr.AddVote(v)

		assert.True(t, consA.HasVote(v.Hash()))
		assert.False(t, consB.HasVote(v.Hash()))
	})

	t.Run("Testing set proposal", func(t *testing.T) {
		consHeight, _ := mgr.HeightRound()
		b, _ := st.ProposeBlock(valKeys[0], valKeys[0].Address())
		p := proposal.NewProposal(consHeight, 0, b)
		ts.HelperSignProposal(valKeys[0], p)

		mgr.SetProposal(p)

		assert.Equal(t, p, consA.Proposal())
		assert.Nil(t, consB.Proposal())
	})

	t.Run("Check discarding old votes", func(t *testing.T) {
		consHeight, _ := mgr.HeightRound()
		v := vote.NewPrepareVote(ts.RandHash(), consHeight-1, 0, st.TestValKeys[2].Address())
		ts.HelperSignVote(st.TestValKeys[2], v)

		mgr.AddVote(v)
		assert.Empty(t, mgr.upcomingVotes)
	})

	t.Run("Check discarding old proposals", func(t *testing.T) {
		consHeight, _ := mgr.HeightRound()
		b, _ := st.ProposeBlock(valKeys[0], valKeys[0].Address())
		p := proposal.NewProposal(consHeight-1, 1, b)
		ts.HelperSignProposal(valKeys[0], p)

		mgr.SetProposal(p)
		assert.Empty(t, mgr.upcomingProposals)
	})

	t.Run("Processing upcoming votes", func(t *testing.T) {
		consHeight, _ := mgr.HeightRound()
		v1 := vote.NewPrepareVote(ts.RandHash(), consHeight+1, 0, valKeys[0].Address())
		v2 := vote.NewPrepareVote(ts.RandHash(), consHeight+2, 0, valKeys[0].Address())
		v3 := vote.NewPrepareVote(ts.RandHash(), consHeight+3, 0, valKeys[0].Address())

		ts.HelperSignVote(valKeys[0], v1)
		ts.HelperSignVote(valKeys[0], v2)
		ts.HelperSignVote(valKeys[0], v3)

		mgr.AddVote(v1)
		mgr.AddVote(v2)
		mgr.AddVote(v3)

		assert.Len(t, mgr.upcomingVotes, 3)

		blk1, cert1 := ts.GenerateTestBlock(consHeight)
		err := st.CommitBlock(blk1, cert1)
		assert.NoError(t, err)

		blk2, cert2 := ts.GenerateTestBlock(consHeight + 1)
		err = st.CommitBlock(blk2, cert2)
		assert.NoError(t, err)

		mgr.MoveToNewHeight()

		assert.Len(t, mgr.upcomingVotes, 1)
	})

	t.Run("Processing upcoming proposal", func(t *testing.T) {
		consHeight, _ := mgr.HeightRound()
		b1, _ := st.ProposeBlock(valKeys[0], valKeys[0].Address())
		p1 := proposal.NewProposal(consHeight+1, 0, b1)

		b2, _ := st.ProposeBlock(valKeys[0], valKeys[0].Address())
		p2 := proposal.NewProposal(consHeight+2, 0, b2)

		b3, _ := st.ProposeBlock(valKeys[0], valKeys[0].Address())
		p3 := proposal.NewProposal(consHeight+3, 0, b3)

		ts.HelperSignProposal(valKeys[0], p1)
		ts.HelperSignProposal(valKeys[0], p2)
		ts.HelperSignProposal(valKeys[0], p3)

		mgr.SetProposal(p1)
		mgr.SetProposal(p2)
		mgr.SetProposal(p3)

		assert.Len(t, mgr.upcomingProposals, 3)

		blk1, cert1 := ts.GenerateTestBlock(consHeight)
		err := st.CommitBlock(blk1, cert1)
		assert.NoError(t, err)

		blk2, cert2 := ts.GenerateTestBlock(consHeight + 1)
		err = st.CommitBlock(blk2, cert2)
		assert.NoError(t, err)

		mgr.MoveToNewHeight()

		assert.Len(t, mgr.upcomingProposals, 1)
	})
}

func TestMediator(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	st := state.MockingState(ts)
	cmt, valKeys := ts.GenerateTestCommittee(4)
	st.TestCommittee = cmt
	st.TestParams.BlockIntervalInSecond = 1

	rewardAddrs := []crypto.Address{
		ts.RandAccAddress(), ts.RandAccAddress(),
		ts.RandAccAddress(), ts.RandAccAddress(),
	}
	broadcastCh := make(chan message.Message, 500)

	stateHeight := ts.RandHeight()
	blk, cert := ts.GenerateTestBlock(stateHeight)
	st.TestStore.SaveBlock(blk, cert)

	mgrInst := NewManager(testConfig(), st, valKeys, rewardAddrs, broadcastCh)
	mgr := mgrInst.(*manager)

	mgr.MoveToNewHeight()

	for {
		msg := <-broadcastCh
		logger.Info("shouldPublishProposal", "msg", msg)

		m, ok := msg.(*message.BlockAnnounceMessage)
		if ok {
			require.Equal(t, stateHeight+1, m.Height())

			return
		}
	}
}
