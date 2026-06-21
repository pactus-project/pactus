package manager_test

import (
	"testing"

	"github.com/pactus-project/pactus/consensus"
	"github.com/pactus-project/pactus/consensus/manager"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestManager(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	state := state.NewMockState(ts.MockController())
	consA := consensus.NewMockConsensus(ts.MockController())
	consB := consensus.NewMockConsensus(ts.MockController())
	instances := []consensus.Consensus{consA, consB}

	consA.EXPECT().MoveToNewHeight().Return().AnyTimes()
	consB.EXPECT().MoveToNewHeight().Return().AnyTimes()

	mgr := manager.NewManager(state, instances)
	height := ts.RandHeight()
	round := ts.RandRound()

	t.Run("Test Instances", func(t *testing.T) {
		instances := mgr.Instances()
		assert.Equal(t, []consensus.ConsensusReader{consA, consB}, instances)
	})

	t.Run("Has no active instances", func(t *testing.T) {
		consA.EXPECT().IsActive().Return(false).Times(2)
		consB.EXPECT().IsActive().Return(false).Times(2)

		assert.False(t, mgr.HasActiveInstance())

		t.Run("Test Get Proposal", func(t *testing.T) {
			prop := ts.GenerateTestProposal(height-1, round)
			consA.EXPECT().Proposal().Return(prop).Times(1)

			assert.Equal(t, prop, mgr.Proposal())
		})
	})

	t.Run("Has an active instances", func(t *testing.T) {
		consA.EXPECT().IsActive().Return(false).AnyTimes()
		consB.EXPECT().IsActive().Return(true).AnyTimes()

		assert.True(t, mgr.HasActiveInstance())

		t.Run("Test Get Proposal", func(t *testing.T) {
			prop := ts.GenerateTestProposal(height, round)
			consB.EXPECT().Proposal().Return(prop).Times(1)

			assert.Equal(t, prop, mgr.Proposal())
		})

		t.Run("Testing HeightRound", func(t *testing.T) {
			consB.EXPECT().HeightRound().Return(height, round).Times(1)

			h, r := mgr.HeightRound()
			assert.Equal(t, height, h)
			assert.Equal(t, round, r)
		})

		t.Run("Testing HandleQueryProposal", func(*testing.T) {
			consB.EXPECT().HandleQueryProposal(height, round).Times(1)
			mgr.HandleQueryProposal(height, round)
		})

		t.Run("Testing HandleQueryVote", func(*testing.T) {
			consB.EXPECT().HandleQueryVote(height, round).Times(1)
			mgr.HandleQueryVote(height, round)
		})

		t.Run("Testing AddVote", func(t *testing.T) {
			t.Run("Discard old votes", func(*testing.T) {
				consB.EXPECT().HeightRound().Return(height, types.Round(0)).Times(1)
				vote, _ := ts.GenerateTestPrecommitVote(height-2, round)

				mgr.AddVote(vote)
			})

			t.Run("Add votes for previous height", func(*testing.T) {
				consB.EXPECT().HeightRound().Return(height, types.Round(0)).Times(1)
				vote, _ := ts.GenerateTestPrecommitVote(height-1, round)

				state.EXPECT().UpdateLastCertificate(vote)
				mgr.AddVote(vote)
			})

			t.Run("Add votes for current height", func(*testing.T) {
				consB.EXPECT().HeightRound().Return(height, types.Round(0)).Times(1)
				vote, _ := ts.GenerateTestPrecommitVote(height, round)

				consA.EXPECT().AddVote(vote).Return().Times(1)
				consB.EXPECT().AddVote(vote).Return().Times(1)
				mgr.AddVote(vote)
			})

			t.Run("Add votes for next height", func(*testing.T) {
				consB.EXPECT().HeightRound().Return(height, types.Round(0)).Times(1)
				vote, _ := ts.GenerateTestPrecommitVote(height+1, round)

				mgr.AddVote(vote)

				// Moving too the next height, votes should be added.
				state.EXPECT().LastBlockHeight().Return(height + 1).Times(1)
				consB.EXPECT().HeightRound().Return(height+1, types.Round(0)).Times(1)
				consA.EXPECT().AddVote(vote).Return().Times(1)
				consB.EXPECT().AddVote(vote).Return().Times(1)
				mgr.MoveToNewHeight()
			})

			t.Run("Discard future votes", func(*testing.T) {
				consB.EXPECT().HeightRound().Return(height, types.Round(0)).Times(1)
				vote, _ := ts.GenerateTestPrecommitVote(height+2, round)

				mgr.AddVote(vote)
			})
		})

		t.Run("Testing SetProposal", func(t *testing.T) {
			t.Run("Discard old proposals", func(*testing.T) {
				consB.EXPECT().HeightRound().Return(height, types.Round(0)).Times(1)
				prop := ts.GenerateTestProposal(height-1, round)

				mgr.SetProposal(prop)
			})

			t.Run("Set proposal for current height", func(*testing.T) {
				consB.EXPECT().HeightRound().Return(height, types.Round(0)).Times(1)
				prop := ts.GenerateTestProposal(height, round)

				consA.EXPECT().SetProposal(prop).Return().Times(1)
				consB.EXPECT().SetProposal(prop).Return().Times(1)

				mgr.SetProposal(prop)
			})

			t.Run("Set proposal for next height", func(*testing.T) {
				consB.EXPECT().HeightRound().Return(height, types.Round(0)).Times(1)
				prop := ts.GenerateTestProposal(height+1, round)

				mgr.SetProposal(prop)

				// Moving too the next height, votes should be added.
				state.EXPECT().LastBlockHeight().Return(height + 1).Times(1)
				consB.EXPECT().HeightRound().Return(height+1, types.Round(0)).Times(1)
				consA.EXPECT().SetProposal(prop).Return().Times(1)
				consB.EXPECT().SetProposal(prop).Return().Times(1)

				mgr.MoveToNewHeight()
			})
		})
	})
}
