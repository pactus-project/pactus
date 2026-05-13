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

	state := state.NewMockState(ts.Ctrl)
	consA := consensus.NewMockConsensus(ts.Ctrl)
	consB := consensus.NewMockConsensus(ts.Ctrl)
	instances := []consensus.Consensus{consA, consB}

	consA.EXPECT().MoveToNewHeight().Return().AnyTimes()
	consB.EXPECT().MoveToNewHeight().Return().AnyTimes()

	mgr := manager.NewManager(state, instances)

	height := ts.RandHeight()

	t.Run("Has no active instances", func(t *testing.T) {
		consA.EXPECT().IsActive().Return(false).Times(1)
		consB.EXPECT().IsActive().Return(false).Times(1)

		assert.False(t, mgr.HasActiveInstance())
	})

	t.Run("Has an active instances", func(t *testing.T) {
		consA.EXPECT().IsActive().Return(false).AnyTimes()
		consB.EXPECT().IsActive().Return(true).AnyTimes()

		assert.True(t, mgr.HasActiveInstance())

		t.Run("Testing AddVote", func(t *testing.T) {
			t.Run("Add votes for previous height", func(t *testing.T) {
				consB.EXPECT().HeightRound().Return(height, types.Round(0)).Times(1)
				vote, _ := ts.GenerateTestPrecommitVote(height-1, ts.RandRound())

				state.EXPECT().UpdateLastCertificate(vote)
				mgr.AddVote(vote)

				consB.EXPECT().HeightRound().Return(height+1, types.Round(0)).Times(1)
				mgr.MoveToNewHeight()
			})

			t.Run("Add votes for current height", func(t *testing.T) {
				consB.EXPECT().HeightRound().Return(height, types.Round(0)).Times(1)
				vote, _ := ts.GenerateTestPrecommitVote(height, ts.RandRound())

				consA.EXPECT().AddVote(vote).Return().Times(1)
				consB.EXPECT().AddVote(vote).Return().Times(1)
				mgr.AddVote(vote)

				consB.EXPECT().HeightRound().Return(height+1, types.Round(0)).Times(1)
				mgr.MoveToNewHeight()
			})

			t.Run("Add votes for next height", func(t *testing.T) {
				consB.EXPECT().HeightRound().Return(height, types.Round(0)).Times(1)
				vote, _ := ts.GenerateTestPrecommitVote(height+1, ts.RandRound())

				mgr.AddVote(vote)

				// Moving too the next height, votes should be added.
				consB.EXPECT().HeightRound().Return(height+1, types.Round(0)).Times(1)
				consA.EXPECT().AddVote(vote).Return().Times(1)
				consB.EXPECT().AddVote(vote).Return().Times(1)
				mgr.MoveToNewHeight()
			})
		})

		t.Run("Testing SetProposal", func(t *testing.T) {
			t.Run("Set proposal for previous height", func(t *testing.T) {
				consB.EXPECT().HeightRound().Return(height, types.Round(0)).Times(1)
				prop := ts.GenerateTestProposal(height-1, ts.RandRound())

				mgr.SetProposal(prop)

				consB.EXPECT().HeightRound().Return(height+1, types.Round(0)).Times(1)
				mgr.MoveToNewHeight()
			})

			t.Run("Set proposal for current height", func(t *testing.T) {
				consB.EXPECT().HeightRound().Return(height, types.Round(0)).Times(1)
				prop := ts.GenerateTestProposal(height, ts.RandRound())

				consA.EXPECT().SetProposal(prop).Return().Times(1)
				consB.EXPECT().SetProposal(prop).Return().Times(1)

				mgr.SetProposal(prop)

				consB.EXPECT().HeightRound().Return(height+1, types.Round(0)).Times(1)

				mgr.MoveToNewHeight()
			})

			t.Run("Set proposal for next height", func(t *testing.T) {
				consB.EXPECT().HeightRound().Return(height, types.Round(0)).Times(1)
				prop := ts.GenerateTestProposal(height+1, ts.RandRound())

				mgr.SetProposal(prop)

				// Moving too the next height, votes should be added.
				consB.EXPECT().HeightRound().Return(height+1, types.Round(0)).Times(1)
				consA.EXPECT().SetProposal(prop).Return().Times(1)
				consB.EXPECT().SetProposal(prop).Return().Times(1)

				mgr.MoveToNewHeight()
			})
		})

		// t.Run("Testing set proposal", func(t *testing.T) {
		// 	consHeight, _ := mgr.HeightRound()
		// 	blk, _ := state.ProposeBlock(valKeys[0], valKeys[0].Address())
		// 	prop := proposal.NewProposal(consHeight, 0, blk)
		// 	ts.HelperSignProposal(valKeys[0], prop)

		// 	mgr.SetProposal(prop)

		// 	assert.Equal(t, prop, consA.Proposal())
		// 	assert.Nil(t, consB.Proposal())
		// })

		// // t.Run("Check discarding old votes", func(t *testing.T) {
		// // 	consHeight, _ := mgr.HeightRound()
		// // 	v := vote.NewPrepareVote(ts.RandHash(), consHeight-1, 0, state.TestValKeys[2].Address())
		// // 	ts.HelperSignVote(state.TestValKeys[2], v)

		// // 	mgr.AddVote(v)
		// // 	assert.Empty(t, mgr.upcomingVotes)
		// // })

		// t.Run("Check discarding old proposals", func(t *testing.T) {
		// 	consHeight, _ := mgr.HeightRound()
		// 	blk, _ := state.ProposeBlock(valKeys[0], valKeys[0].Address())
		// 	prop := proposal.NewProposal(consHeight-1, 1, blk)
		// 	ts.HelperSignProposal(valKeys[0], prop)

		// 	mgr.SetProposal(prop)
		// 	assert.Empty(t, mgr.upcomingProposals)
		// })

		// t.Run("Processing upcoming votes", func(t *testing.T) {
		// 	consHeight, _ := mgr.HeightRound()
		// 	vote1 := vote.NewPrepareVote(ts.RandHash(), consHeight+1, 0, valKeys[0].Address())
		// 	vote2 := vote.NewPrepareVote(ts.RandHash(), consHeight+2, 0, valKeys[0].Address())
		// 	vote3 := vote.NewPrepareVote(ts.RandHash(), consHeight+3, 0, valKeys[0].Address())

		// 	ts.HelperSignVote(valKeys[0], vote1)
		// 	ts.HelperSignVote(valKeys[0], vote2)
		// 	ts.HelperSignVote(valKeys[0], vote3)

		// 	mgr.AddVote(vote1)
		// 	mgr.AddVote(vote2)
		// 	mgr.AddVote(vote3)

		// 	assert.Len(t, mgr.upcomingVotes, 3)

		// 	blk1, cert1 := ts.GenerateTestBlock(consHeight)
		// 	err := state.CommitBlock(blk1, cert1)
		// 	require.NoError(t, err)

		// 	blk2, cert2 := ts.GenerateTestBlock(consHeight + 1)
		// 	err = state.CommitBlock(blk2, cert2)
		// 	require.NoError(t, err)

		// 	mgr.MoveToNewHeight()

		// 	assert.Len(t, mgr.upcomingVotes, 1)
		// })

		// t.Run("Processing upcoming proposal", func(t *testing.T) {
		// 	consHeight, _ := mgr.HeightRound()
		// 	prop1 := ts.GenerateTestProposal(consHeight+1, 0)
		// 	prop2 := ts.GenerateTestProposal(consHeight+2, 0)
		// 	prop3 := ts.GenerateTestProposal(consHeight+3, 0)

		// 	mgr.SetProposal(prop1)
		// 	mgr.SetProposal(prop2)
		// 	mgr.SetProposal(prop3)

		// 	assert.Len(t, mgr.upcomingProposals, 3)

		// 	blk1, cert1 := ts.GenerateTestBlock(consHeight)
		// 	err := state.CommitBlock(blk1, cert1)
		// 	require.NoError(t, err)

		// 	blk2, cert2 := ts.GenerateTestBlock(consHeight + 1)
		// 	err = state.CommitBlock(blk2, cert2)
		// 	require.NoError(t, err)

		// 	mgr.MoveToNewHeight()

		// 	assert.Len(t, mgr.upcomingProposals, 1)
		// })

	})
}

// func TestMediator(t *testing.T) {
// 	ts := testsuite.NewTestSuite(t)

// 	state := state.NewMockState(ts.Ctrl)
// 	_, valKeys := ts.GenerateTestCommittee(4)

// 	rewardAddrs := []crypto.Address{
// 		ts.RandAccAddress(), ts.RandAccAddress(),
// 		ts.RandAccAddress(), ts.RandAccAddress(),
// 	}
// 	stateHeight := ts.RandHeight()
// 	blk, cert := ts.GenerateTestBlock(stateHeight)
// 	state.TestStore.SaveBlock(blk, cert)
// 	pipe := pipeline.New[message.Message](t.Context())
// 	conf := consensusv2.DefaultConfig()

// 	mgr := manager.NewManagerV2(t.Context(), conf, state, valKeys, rewardAddrs, pipe)
// 	// mgr := mgrInt.(*manager)

// 	mgr.MoveToNewHeight()

// 	for {
// 		msg := <-pipe.UnsafeGetChannel()
// 		logger.Info("Published Vote", "msg", msg, "type", msg.Type())

// 		m, ok := msg.(*message.BlockAnnounceMessage)
// 		if ok {
// 			require.Equal(t, stateHeight+1, m.Height())

// 			return
// 		}
// 	}
// }
