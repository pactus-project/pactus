package manager_test

import (
	"testing"

	"github.com/pactus-project/pactus/consensus"
	"github.com/pactus-project/pactus/consensus/manager"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/state"
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// The state and the instances are strict mocks: the observer must not call
// anything else on them, so it can neither write to the state nor make an
// instance sign or vote.
func setupObservingManager(t *testing.T, active bool) (
	*testsuite.TestSuite, manager.ConsensusManager, *bls.ValidatorKey, types.Height,
) {
	t.Helper()

	ts := testsuite.NewTestSuite(t)
	mockState := state.NewMockState(ts.MockController())
	cons := consensus.NewMockConsensus(ts.MockController())
	mgr := manager.NewManager(mockState, []consensus.Consensus{cons}, manager.WithObserver(true))

	key := ts.RandValKey()
	member := ts.GenerateTestValidator(testsuite.ValidatorWithPublicKey(key.PublicKey()))
	height := ts.RandHeight()

	cons.EXPECT().IsActive().Return(active).AnyTimes()
	cons.EXPECT().HeightRound().Return(height, types.Round(0)).AnyTimes()
	cons.EXPECT().MoveToNewHeight().Times(1)
	cons.EXPECT().AddVote(gomock.Any()).AnyTimes()
	// The height of the consensus is committed, so the manager moves on.
	mockState.EXPECT().LastBlockHeight().Return(height).Times(1)
	mockState.EXPECT().CommitteeValidators().Return([]*validator.Validator{member}).Times(1)

	mgr.MoveToNewHeight()

	return ts, mgr, key, height
}

func TestManagerObserverKeepsVotesWhenNotInCommittee(t *testing.T) {
	ts, mgr, key, height := setupObservingManager(t, false)

	vte := vote.NewPrecommitVote(ts.RandHash(), height, 0, key.Address())
	ts.HelperSignVote(key, vte)
	mgr.AddVote(vte)

	obs := mgr.Observation()
	require.NotNil(t, obs)
	require.Len(t, obs.Votes, 1)
	assert.Equal(t, key.Address(), obs.Votes[0].Signer)
}

func TestManagerObserverIdleWhenInCommittee(t *testing.T) {
	ts, mgr, key, height := setupObservingManager(t, true)

	vte := vote.NewPrecommitVote(ts.RandHash(), height, 0, key.Address())
	ts.HelperSignVote(key, vte)
	mgr.AddVote(vte)

	assert.Nil(t, mgr.Observation())
}

func TestManagerObserverOffByDefault(t *testing.T) {
	ts := testsuite.NewTestSuite(t)
	mockState := state.NewMockState(ts.MockController())
	cons := consensus.NewMockConsensus(ts.MockController())
	mgr := manager.NewManager(mockState, []consensus.Consensus{cons})

	// With the observer off, moving to a new height does not read the committee.
	cons.EXPECT().HeightRound().Return(types.Height(10), types.Round(0)).AnyTimes()
	cons.EXPECT().IsActive().Return(false).AnyTimes()
	cons.EXPECT().MoveToNewHeight().Times(1)
	mockState.EXPECT().LastBlockHeight().Return(types.Height(10)).Times(1)

	mgr.MoveToNewHeight()

	assert.Nil(t, mgr.Observation())
}
