package manager

import (
	"runtime"
	"testing"

	"github.com/pactus-project/pactus/consensus"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// manager.AddVote caches ANY vote for height+1 into
// upcomingVotes with no size bound and no validation whatsoever (no
// signature check, no check that the signer is even a validator) — unlike
// the vote.Height() == curHeight branch, which forwards to each consensus
// instance's own AddVote (itself checked there). A peer that is never asked
// to authenticate can grow this slice without limit simply by sending votes
// one height ahead.
//
// This test expresses the desired behavior (the cache must stay bounded, or
// at least reject unsigned/unverifiable votes) and fails today: 10,000
// completely unsigned, non-committee votes are all cached.
func TestUpcomingVotesCacheIsUnbounded(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	mockState := state.NewMockState(ts.MockController())
	consA := consensus.NewMockConsensus(ts.MockController())

	height := ts.RandHeight()
	round := ts.RandRound()
	consA.EXPECT().IsActive().Return(true).AnyTimes()
	consA.EXPECT().HeightRound().Return(height, round).AnyTimes()

	mgrIface := NewManager(mockState, []consensus.Consensus{consA})
	mgr, ok := mgrIface.(*manager)
	require.True(t, ok)

	const acceptableBound = 1000
	const injectedCount = 10_000

	var before, after runtime.MemStats
	runtime.GC() //nolint:revive // the test measures the memory held by the caches.
	runtime.ReadMemStats(&before)

	for i := 0; i < injectedCount; i++ {
		// A vote for height+1, unsigned, from an address that is not (and
		// need not be) a validator. Nothing checks either fact before caching it.
		v := vote.NewPrecommitVote(hash.UndefHash, height+1, round, ts.RandAccAddress())
		mgr.AddVote(v)
	}

	runtime.ReadMemStats(&after)
	t.Logf("%d unsigned, non-committee votes cached; heap grew by %d bytes; upcomingVotes length = %d",
		injectedCount, after.TotalAlloc-before.TotalAlloc, len(mgr.upcomingVotes))

	assert.Less(t, len(mgr.upcomingVotes), acceptableBound,
		"upcomingVotes must not grow without bound from unauthenticated, unsigned votes")
}

// Same claim for proposals: SetProposal caches any proposal for height+1 with
// no bound and no signature check.
func TestUpcomingProposalsCacheIsUnbounded(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	mockState := state.NewMockState(ts.MockController())
	consA := consensus.NewMockConsensus(ts.MockController())

	height := ts.RandHeight()
	round := ts.RandRound()
	consA.EXPECT().IsActive().Return(true).AnyTimes()
	consA.EXPECT().HeightRound().Return(height, round).AnyTimes()

	mgrIface := NewManager(mockState, []consensus.Consensus{consA})
	mgr, ok := mgrIface.(*manager)
	require.True(t, ok)

	const acceptableBound = 200
	const injectedCount = 500 // GenerateTestProposal is comparatively expensive; keep the default run fast.

	for i := 0; i < injectedCount; i++ {
		prop := ts.GenerateTestProposal(height+1, types.Round(0))
		prop.SetSignature(nil) // unsigned: nothing checks it before caching.
		mgr.SetProposal(prop)
	}

	t.Logf("upcomingProposals length = %d", len(mgr.upcomingProposals))

	assert.Less(t, len(mgr.upcomingProposals), acceptableBound,
		"upcomingProposals must not grow without bound from unauthenticated, unsigned proposals")
}
