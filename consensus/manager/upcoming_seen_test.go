package manager

import (
	"testing"

	"github.com/pactus-project/pactus/consensus"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// The hashes of the cached messages are forgotten once the messages are replayed:
// they would otherwise add up, height after height.
func TestUpcomingSeenIsClearedOnReplay(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	mockState := state.NewMockState(ts.MockController())
	consA := consensus.NewMockConsensus(ts.MockController())

	height := ts.RandHeight()
	consA.EXPECT().IsActive().Return(true).AnyTimes()
	consA.EXPECT().HeightRound().Return(height, ts.RandRound()).AnyTimes()
	consA.EXPECT().MoveToNewHeight().AnyTimes()
	consA.EXPECT().AddVote(gomock.Any()).AnyTimes()
	consA.EXPECT().SetProposal(gomock.Any()).AnyTimes()
	mockState.EXPECT().LastBlockHeight().Return(height).AnyTimes()

	mgr, ok := NewManager(mockState, []consensus.Consensus{consA}).(*manager)
	require.True(t, ok)

	vte, _ := ts.GenerateTestPrecommitVote(height+1, 0)
	mgr.AddVote(vte)
	mgr.SetProposal(ts.GenerateTestProposal(height+1, 0))
	require.Len(t, mgr.upcomingSeen, 2)

	mgr.MoveToNewHeight()
	assert.Empty(t, mgr.upcomingSeen)
	assert.Empty(t, mgr.upcomingVotes)
	assert.Empty(t, mgr.upcomingProposals)
}
