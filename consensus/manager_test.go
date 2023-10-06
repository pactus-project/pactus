package consensus

import (
	"fmt"
	"testing"
	"time"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/logger"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func shouldPublishProposal(t *testing.T, broadcastCh chan message.Message,
	height uint32, round int16,
) *proposal.Proposal {
	t.Helper()

	timeout := time.NewTimer(1 * time.Second)

	for {
		select {
		case <-timeout.C:
			require.NoError(t, fmt.Errorf("Timeout"))
			return nil
		case msg := <-broadcastCh:
			logger.Info("shouldPublishProposal", "message", msg)

			if msg.Type() == message.TypeProposal {
				m := msg.(*message.ProposalMessage)
				require.Equal(t, m.Proposal.Height(), height)
				require.Equal(t, m.Proposal.Round(), round)
				return m.Proposal
			}
		}
	}
}

func TestManager(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	state := state.MockingState(ts)
	state.TestCommittee.Validators()

	committeeValKeys := state.TestValKeys

	rewardAddrs := []crypto.Address{
		ts.RandAccAddress(), ts.RandAccAddress(),
		ts.RandAccAddress(), ts.RandAccAddress(),
		ts.RandAccAddress(),
	}
	valKeys := make([]*bls.ValidatorKey, 5)
	valKeys[0] = committeeValKeys[0]
	valKeys[1] = ts.RandValKey()
	valKeys[2] = committeeValKeys[1]
	valKeys[3] = ts.RandValKey()
	valKeys[4] = ts.RandValKey()
	broadcastCh := make(chan message.Message, 500)

	height := ts.RandHeight()
	blk, cert := ts.GenerateTestBlock(height)
	state.TestStore.SaveBlock(blk, cert)

	Mgr := NewManager(testConfig(), state, valKeys, rewardAddrs, broadcastCh)
	mgr := Mgr.(*manager)

	consA := mgr.instances[0].(*consensus) // active
	consB := mgr.instances[1].(*consensus) // inactive
	consC := mgr.instances[2].(*consensus) // active
	consD := mgr.instances[3].(*consensus) // inactive
	consE := mgr.instances[4].(*consensus) // inactive

	assert.False(t, mgr.HasActiveInstance())
	mgr.MoveToNewHeight()
	curHeight, _ := mgr.HeightRound()

	newHeightTimeout(consA)

	t.Run("Check if keys are assigned properly", func(t *testing.T) {
		assert.Equal(t, valKeys[0].PublicKey(), consA.ConsensusKey())
		assert.Equal(t, valKeys[1].PublicKey(), consB.ConsensusKey())
		assert.Equal(t, valKeys[2].PublicKey(), consC.ConsensusKey())
		assert.Equal(t, valKeys[3].PublicKey(), consD.ConsensusKey())
		assert.Equal(t, valKeys[4].PublicKey(), consE.ConsensusKey())
	})

	t.Run("Check if all instances move to new height", func(t *testing.T) {
		assert.True(t, mgr.HasActiveInstance())
	})

	t.Run("Check if all instances move to new height", func(t *testing.T) {
		h, r := mgr.HeightRound()
		assert.Equal(t, h, curHeight)
		assert.Equal(t, r, int16(0))
		assert.True(t, mgr.HasActiveInstance())
	})

	t.Run("Testing add vote", func(t *testing.T) {
		v := vote.NewPrepareVote(ts.RandHash(), curHeight, 0, committeeValKeys[2].Address())
		ts.HelperSignVote(committeeValKeys[2], v)

		mgr.AddVote(v)

		assert.True(t, consA.HasVote(v.Hash()))
		assert.False(t, consB.HasVote(v.Hash()))
		assert.True(t, consC.HasVote(v.Hash()))
		assert.False(t, consD.HasVote(v.Hash()))
	})

	t.Run("Testing set proposal", func(t *testing.T) {
		b, _ := state.ProposeBlock(committeeValKeys[1], committeeValKeys[1].Address(), 1)
		p := proposal.NewProposal(curHeight, 1, b)
		ts.HelperSignProposal(committeeValKeys[1], p)

		mgr.SetProposal(p)

		assert.Equal(t, consA.RoundProposal(1), p)
		assert.Nil(t, consB.RoundProposal(1))
		assert.Equal(t, consC.RoundProposal(1), p)
		assert.Nil(t, consD.RoundProposal(1))
	})

	t.Run("Check if one instance publishes a proposal, the other instances receive it", func(t *testing.T) {
		p := shouldPublishProposal(t, broadcastCh, curHeight, 0)

		assert.Equal(t, mgr.RoundProposal(0), p)
		assert.Equal(t, consA.RoundProposal(0), p)
		assert.Nil(t, consB.RoundProposal(0))
	})

	t.Run("Check discarding old votes", func(t *testing.T) {
		v := vote.NewPrepareVote(ts.RandHash(), curHeight-1, 0, committeeValKeys[2].Address())
		ts.HelperSignVote(committeeValKeys[2], v)

		mgr.AddVote(v)
		assert.Empty(t, mgr.upcomingVotes)
	})

	t.Run("Check discarding old proposals", func(t *testing.T) {
		b, _ := state.ProposeBlock(committeeValKeys[1], committeeValKeys[1].Address(), 1)
		p := proposal.NewProposal(curHeight-1, 1, b)
		ts.HelperSignProposal(committeeValKeys[1], p)

		mgr.SetProposal(p)
		assert.Empty(t, mgr.upcomingProposals)
	})

	t.Run("Processing upcoming votes", func(t *testing.T) {
		v1 := vote.NewPrepareVote(ts.RandHash(), curHeight+1, 0, committeeValKeys[2].Address())
		v2 := vote.NewPrepareVote(ts.RandHash(), curHeight+2, 0, committeeValKeys[2].Address())
		v3 := vote.NewPrepareVote(ts.RandHash(), curHeight+3, 0, committeeValKeys[2].Address())

		ts.HelperSignVote(committeeValKeys[2], v1)
		ts.HelperSignVote(committeeValKeys[2], v2)
		ts.HelperSignVote(committeeValKeys[2], v3)

		mgr.AddVote(v1)
		mgr.AddVote(v2)
		mgr.AddVote(v3)

		assert.Len(t, mgr.upcomingVotes, 3)

		blk, cert := ts.GenerateTestBlockWithTime(curHeight, util.Now().Add(1*time.Hour))
		err := state.CommitBlock(blk, cert)
		assert.NoError(t, err)
		curHeight++

		blk, cert = ts.GenerateTestBlockWithTime(curHeight, util.Now().Add(2*time.Hour))
		err = state.CommitBlock(blk, cert)
		assert.NoError(t, err)
		curHeight++

		mgr.MoveToNewHeight()

		assert.Len(t, mgr.upcomingVotes, 1)
	})

	t.Run("Processing upcoming proposal", func(t *testing.T) {
		b1, _ := state.ProposeBlock(committeeValKeys[1], committeeValKeys[1].Address(), 1)
		p1 := proposal.NewProposal(curHeight+1, 1, b1)

		b2, _ := state.ProposeBlock(committeeValKeys[1], committeeValKeys[1].Address(), 1)
		p2 := proposal.NewProposal(curHeight+2, 1, b2)

		b3, _ := state.ProposeBlock(committeeValKeys[1], committeeValKeys[1].Address(), 1)
		p3 := proposal.NewProposal(curHeight+3, 1, b3)

		ts.HelperSignProposal(committeeValKeys[1], p1)
		ts.HelperSignProposal(committeeValKeys[1], p2)
		ts.HelperSignProposal(committeeValKeys[1], p3)

		mgr.SetProposal(p1)
		mgr.SetProposal(p2)
		mgr.SetProposal(p3)

		assert.Len(t, mgr.upcomingProposals, 3)

		blk, cert := ts.GenerateTestBlockWithTime(curHeight, util.Now().Add(1*time.Hour))
		err := state.CommitBlock(blk, cert)
		assert.NoError(t, err)
		curHeight++

		blk, cert = ts.GenerateTestBlockWithTime(curHeight, util.Now().Add(2*time.Hour))
		err = state.CommitBlock(blk, cert)
		assert.NoError(t, err)

		mgr.MoveToNewHeight()

		assert.Len(t, mgr.upcomingProposals, 1)
	})
}
