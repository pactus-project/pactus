package manager

import (
	"testing"
	"time"

	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/pactus-project/pactus/util/ratelimit"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testData struct {
	*testsuite.TestSuite

	keys       []*bls.ValidatorKey
	committee  []*validator.Validator
	height     types.Height
	observer   *observer
	proposerOf func(types.Round) *validator.Validator
}

// committeeSize is the size of the test committee.
const committeeSize = 4

func setupObserver(t *testing.T) *testData {
	t.Helper()

	ts := testsuite.NewTestSuite(t)
	td := &testData{TestSuite: ts, height: ts.RandHeight()}
	for i := 0; i < committeeSize; i++ {
		key := ts.RandValKey()
		td.keys = append(td.keys, key)
		td.committee = append(td.committee, ts.GenerateTestValidator(
			testsuite.ValidatorWithPublicKey(key.PublicKey()),
			testsuite.ValidatorWithNumber(int32(i)),
			testsuite.ValidatorWithStake(amount.Amount(1000e9))))
	}
	// The proposer of a round is the member at that index, as a simple rule.
	td.proposerOf = func(round types.Round) *validator.Validator {
		return td.committee[int(round)%len(td.committee)]
	}
	td.observer = newObserver(true)
	td.observer.moveToHeight(td.height, td.committee, td.proposerOf, false)

	return td
}

func (td *testData) precommit(key *bls.ValidatorKey, round types.Round) *vote.Vote {
	vte := vote.NewPrecommitVote(td.RandHash(), td.height, round, key.Address())
	td.HelperSignVote(key, vte)

	return vte
}

func (td *testData) cpPreVote(key *bls.ValidatorKey, round types.Round, cpRound int16) *vote.Vote {
	vte := vote.NewCPPreVote(td.RandHash(), td.height, round, cpRound,
		vote.CPValueYes, &vote.JustInitYes{}, key.Address())
	td.HelperSignVote(key, vte)

	return vte
}

func (td *testData) proposal(key *bls.ValidatorKey, round types.Round) *proposal.Proposal {
	return td.GenerateTestProposal(td.height, round, testsuite.ProposalWithKey(key))
}

func TestObserverDisabled(t *testing.T) {
	td := setupObserver(t)
	obs := newObserver(false)
	obs.moveToHeight(td.height, td.committee, td.proposerOf, false)

	obs.addVote(td.precommit(td.keys[0], 0))
	obs.setProposal(td.proposal(td.keys[0], 0))

	assert.Nil(t, obs.observation())
}

func TestObserverIdleWithActiveInstance(t *testing.T) {
	td := setupObserver(t)
	td.observer.moveToHeight(td.height, td.committee, td.proposerOf, true)

	td.observer.addVote(td.precommit(td.keys[0], 0))

	assert.Nil(t, td.observer.observation(),
		"an active instance keeps the votes, so the observer must stay idle")
}

func TestObserverKeepsVerifiedVotes(t *testing.T) {
	td := setupObserver(t)

	precommit := td.precommit(td.keys[0], 0)
	cpVote := td.cpPreVote(td.keys[1], 0, 1)
	td.observer.addVote(precommit)
	td.observer.addVote(cpVote)

	obs := td.observer.observation()
	require.NotNil(t, obs)
	assert.Equal(t, td.height, obs.Height)
	assert.Equal(t, int64(4*1000e9), obs.TotalPower)
	require.Len(t, obs.Votes, 2)
	assert.Equal(t, vote.VoteTypePrecommit, obs.Votes[0].Type)
	assert.Equal(t, precommit.Signer(), obs.Votes[0].Signer)
	assert.Equal(t, precommit.BlockHash(), obs.Votes[0].BlockHash)
	assert.Equal(t, int16(0), obs.Votes[0].CPRound, "a precommit has no change-proposer round")
	assert.Equal(t, vote.VoteTypeCPPreVote, obs.Votes[1].Type)
	assert.Equal(t, int16(1), obs.Votes[1].CPRound)
	assert.Equal(t, vote.CPValueYes, obs.Votes[1].CPValue)
}

func TestObserverRejectsUntrustedVotes(t *testing.T) {
	td := setupObserver(t)

	t.Run("signer outside the committee", func(t *testing.T) {
		outsider := td.RandValKey()
		td.observer.addVote(td.precommit(outsider, 0))

		assert.Empty(t, td.observer.observation().Votes)
	})

	t.Run("forged signature", func(t *testing.T) {
		forged := vote.NewPrecommitVote(td.RandHash(), td.height, 0, td.keys[0].Address())
		td.HelperSignVote(td.keys[1], forged)
		td.observer.addVote(forged)

		assert.Empty(t, td.observer.observation().Votes)
		assert.Contains(t, td.observer.badMessages, forged.Hash(),
			"a forged vote is remembered, so it is not verified again")
	})

	t.Run("another height", func(t *testing.T) {
		vte := vote.NewPrecommitVote(td.RandHash(), td.height+1, 0, td.keys[0].Address())
		td.HelperSignVote(td.keys[0], vte)
		td.observer.addVote(vte)

		assert.Empty(t, td.observer.observation().Votes)
	})

	t.Run("change-proposer round above the bound", func(t *testing.T) {
		dropped := td.observer.observation().Dropped
		td.observer.addVote(td.cpPreVote(td.keys[0], 0, observerMaxCPRound))

		assert.Empty(t, td.observer.observation().Votes)
		assert.Equal(t, dropped+1, td.observer.observation().Dropped)
	})
}

func TestObserverIgnoresFilledSlot(t *testing.T) {
	td := setupObserver(t)

	first := td.precommit(td.keys[0], 0)
	td.observer.addVote(first)

	// A second, validly signed vote of the same member for the same slot
	// (an equivocation, or a replay with another block hash) is not verified.
	second := td.precommit(td.keys[0], 0)
	td.observer.addVote(second)
	td.observer.addVote(first)

	obs := td.observer.observation()
	require.Len(t, obs.Votes, 1)
	assert.Equal(t, first.BlockHash(), obs.Votes[0].BlockHash)
}

func TestObserverSumsOtherRounds(t *testing.T) {
	td := setupObserver(t)

	// Round 2 has the votes of three members, rounds 3 and 4 one vote each.
	blockHash := td.RandHash()
	for _, key := range td.keys[:3] {
		vte := vote.NewPrecommitVote(blockHash, td.height, 2, key.Address())
		td.HelperSignVote(key, vte)
		td.observer.addVote(vte)
	}
	td.observer.addVote(td.precommit(td.keys[0], 3))
	td.observer.addVote(td.precommit(td.keys[0], 4))

	// Round 2 is shown in full, since more than a third of the power voted on
	// it, and round 4 as the highest of the others. Round 3 is summed.
	obs := td.observer.observation()
	require.Len(t, obs.Votes, 4)
	shown := map[types.Round]int{}
	for _, vte := range obs.Votes {
		shown[vte.Round]++
	}
	assert.Equal(t, map[types.Round]int{2: 3, 4: 1}, shown)

	require.Len(t, obs.Tallies, 1)
	assert.Equal(t, types.Round(3), obs.Tallies[0].Round)
	assert.Equal(t, 1, obs.Tallies[0].Votes)
	assert.Equal(t, int64(1000e9), obs.Tallies[0].Power)
}

func TestObserverByzantineCannotHideCommitteeRounds(t *testing.T) {
	td := setupObserver(t)

	// The committee works on round 0.
	blockHash := td.RandHash()
	for _, key := range td.keys[:3] {
		vte := vote.NewPrecommitVote(blockHash, td.height, 0, key.Address())
		td.HelperSignVote(key, vte)
		td.observer.addVote(vte)
	}

	// A byzantine member signs valid votes for many far rounds.
	byzantine := td.keys[3]
	for round := types.Round(400); round < 400+3*observerMaxRounds; round++ {
		td.observer.addVote(td.precommit(byzantine, round))
		td.observer.addVote(td.cpPreVote(byzantine, round, 1))
	}

	assert.LessOrEqual(t, len(td.observer.rounds), observerMaxRounds)
	assert.Contains(t, td.observer.rounds, types.Round(0),
		"the round with the most power is never evicted by a single member")

	obs := td.observer.observation()
	committeeVotes := 0
	for _, vte := range obs.Votes {
		if vte.Round == 0 {
			committeeVotes++
		}
	}
	assert.Equal(t, 3, committeeVotes, "the committee round is still shown vote by vote")
}

func TestObserverRoundPowerCountsEachMemberOnce(t *testing.T) {
	td := setupObserver(t)

	// One member casting many votes in a round does not reach a third of the power.
	for cpRound := int16(0); cpRound < 5; cpRound++ {
		td.observer.addVote(td.cpPreVote(td.keys[0], 0, cpRound))
	}

	assert.Equal(t, int64(1000e9), td.observer.rounds[0].power)
}

func TestObserverBoundsRounds(t *testing.T) {
	td := setupObserver(t)

	// A stall: the rounds keep growing and the height does not change.
	for round := types.Round(0); round < 3*observerMaxRounds; round++ {
		td.observer.addVote(td.precommit(td.keys[int(round)%4], round))
	}

	assert.LessOrEqual(t, len(td.observer.rounds), observerMaxRounds)
	obs := td.observer.observation()
	assert.LessOrEqual(t, len(obs.Votes), observerFullRounds*len(td.committee)*observerMaxCPRound)
	assert.Positive(t, obs.Dropped)
	assert.Contains(t, td.observer.rounds, types.Round(3*observerMaxRounds-1),
		"the newest round is kept")
}

func TestObserverForgedVoteCannotEvictRounds(t *testing.T) {
	td := setupObserver(t)

	for round := types.Round(0); round < observerMaxRounds; round++ {
		td.observer.addVote(td.precommit(td.keys[0], round))
	}
	before := len(td.observer.rounds)

	forged := vote.NewPrecommitVote(td.RandHash(), td.height, 1000, td.keys[0].Address())
	td.HelperSignVote(td.keys[1], forged)
	td.observer.addVote(forged)

	assert.Len(t, td.observer.rounds, before)
	assert.NotContains(t, td.observer.rounds, types.Round(1000))
	assert.Contains(t, td.observer.rounds, types.Round(0))
}

func TestObserverRateLimitsVerification(t *testing.T) {
	td := setupObserver(t)
	// The budget is per second: the test can take longer than that, for example under the race
	// detector. The same budget over a long window makes it independent of the speed of the machine.
	td.observer.verifyLimit = ratelimit.NewRateLimit(observerMaxVerifications, time.Hour)

	// A flood of distinct forged votes for open slots: each one would cost a
	// signature check, so the observer stops checking past its budget.
	for i := 0; i < 2*observerMaxVerifications; i++ {
		round := types.Round(i / observerMaxCPRound)
		cpRound := int16(i % observerMaxCPRound)
		forged := vote.NewCPPreVote(td.RandHash(), td.height, round, cpRound,
			vote.CPValueYes, &vote.JustInitYes{}, td.keys[0].Address())
		td.HelperSignVote(td.keys[2], forged) // signed by another member
		td.observer.addVote(forged)
	}

	obs := td.observer.observation()
	assert.Empty(t, obs.Votes)
	assert.GreaterOrEqual(t, obs.Dropped, uint64(observerMaxVerifications),
		"checks past the budget are skipped")
}

func TestObserverProposals(t *testing.T) {
	td := setupObserver(t)

	t.Run("proposal of the round's proposer", func(t *testing.T) {
		prop := td.proposal(td.keys[1], 1)
		td.observer.setProposal(prop)

		obs := td.observer.observation()
		require.Len(t, obs.Proposals, 1)
		assert.Equal(t, types.Round(1), obs.Proposals[0].Round)
		assert.Equal(t, prop.Block().Hash(), obs.Proposals[0].BlockHash)
		assert.Equal(t, td.keys[1].Address(), obs.Proposals[0].Proposer)
	})

	t.Run("second proposal for the same round", func(t *testing.T) {
		td.observer.setProposal(td.proposal(td.keys[1], 1))

		assert.Len(t, td.observer.observation().Proposals, 1)
	})

	t.Run("proposal of another member", func(t *testing.T) {
		td.observer.setProposal(td.proposal(td.keys[3], 2))

		assert.Len(t, td.observer.observation().Proposals, 1)
	})

	t.Run("proposal with a forged signature", func(t *testing.T) {
		prop := td.proposal(td.keys[2], 2)
		prop.SetSignature(td.keys[3].Sign(prop.SignBytes()))
		td.observer.setProposal(prop)

		assert.Len(t, td.observer.observation().Proposals, 1)
	})
}

func TestObserverProposalsFollowVotedRounds(t *testing.T) {
	td := setupObserver(t)

	// A member fills every proposal slot with rounds nobody votes on.
	for i := 0; i < observerMaxRounds; i++ {
		round := types.Round(committeeSize * (i + 10))
		td.observer.setProposal(td.proposal(td.keys[0], round))
	}
	require.Len(t, td.observer.proposals, observerMaxRounds)

	// The committee votes on round 1, so its proposal still gets in.
	td.observer.addVote(td.precommit(td.keys[2], 1))
	td.observer.setProposal(td.proposal(td.keys[1], 1))

	obs := td.observer.observation()
	require.Len(t, obs.Proposals, 1)
	assert.Equal(t, types.Round(1), obs.Proposals[0].Round)
}

func TestObserverMoveToHeightResets(t *testing.T) {
	td := setupObserver(t)

	td.observer.addVote(td.precommit(td.keys[0], 0))
	td.observer.setProposal(td.proposal(td.keys[0], 0))
	td.height++
	td.observer.moveToHeight(td.height, td.committee, td.proposerOf, false)

	obs := td.observer.observation()
	assert.Equal(t, td.height, obs.Height)
	assert.Empty(t, obs.Votes)
	assert.Empty(t, obs.Tallies)
	assert.Empty(t, obs.Proposals)
	assert.Zero(t, obs.Dropped)
}

func TestObserverObservationIsACopy(t *testing.T) {
	td := setupObserver(t)

	td.observer.addVote(td.precommit(td.keys[0], 0))
	obs := td.observer.observation()
	obs.Votes[0].Round = 99

	assert.Equal(t, types.Round(0), td.observer.observation().Votes[0].Round)
}

// FuzzObserverAddVote feeds the observer votes decoded from arbitrary bytes, as
// a peer can send them. It must never panic: a panic would stop the node.
func FuzzObserverAddVote(f *testing.F) {
	ts := testsuite.NewTestSuite(&testing.T{})
	key := ts.RandValKey()
	member := ts.GenerateTestValidator(testsuite.ValidatorWithPublicKey(key.PublicKey()))
	height := types.Height(100)

	seeds := []*vote.Vote{
		vote.NewPrecommitVote(ts.RandHash(), height, 0, key.Address()),
		vote.NewPrepareVote(ts.RandHash(), height, 1, key.Address()),
		vote.NewCPPreVote(ts.RandHash(), height, 0, 1, vote.CPValueYes, &vote.JustInitYes{}, key.Address()),
		vote.NewCPMainVote(ts.RandHash(), height, 0, 2, vote.CPValueAbstain, &vote.JustInitYes{}, key.Address()),
	}
	for _, seed := range seeds {
		ts.HelperSignVote(key, seed)
		data, err := seed.MarshalCBOR()
		require.NoError(f, err)
		f.Add(data)
	}

	f.Fuzz(func(_ *testing.T, data []byte) {
		vte := new(vote.Vote)
		if err := vte.UnmarshalCBOR(data); err != nil {
			return
		}

		obs := newObserver(true)
		obs.moveToHeight(height, []*validator.Validator{member},
			func(types.Round) *validator.Validator { return member }, false)
		obs.addVote(vte)
		_ = obs.observation()
	})
}
