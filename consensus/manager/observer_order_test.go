package manager

import (
	"math/rand"
	"slices"
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

// The observation is sorted, so that two nodes that saw the same votes report the same list.
// The votes are ordered by round, type, change-proposer round, then signer.
func TestObservedVotesOrder(t *testing.T) {
	ts := testsuite.NewTestSuite(t)
	low := crypto.NewAddress(crypto.AddressTypeValidator, make([]byte, 20))
	high := ts.RandValAddress()
	if slices.Compare(low.Bytes(), high.Bytes()) > 0 {
		low, high = high, low
	}

	ordered := []ObservedVote{
		{Round: 0, Type: vote.VoteTypePrepare, CPRound: 0, Signer: low},
		{Round: 0, Type: vote.VoteTypePrepare, CPRound: 0, Signer: high},
		{Round: 0, Type: vote.VoteTypePrepare, CPRound: 1, Signer: low},
		{Round: 0, Type: vote.VoteTypePrecommit, CPRound: 0, Signer: low},
		{Round: 1, Type: vote.VoteTypePrepare, CPRound: 0, Signer: low},
	}

	for i := range ordered {
		assert.Zero(t, compareObservedVotes(ordered[i], ordered[i]))
		for j := i + 1; j < len(ordered); j++ {
			assert.Negative(t, compareObservedVotes(ordered[i], ordered[j]), "%d before %d", i, j)
			assert.Positive(t, compareObservedVotes(ordered[j], ordered[i]), "%d after %d", j, i)
		}
	}

	shuffled := slices.Clone(ordered)
	rand.New(rand.NewSource(1)).Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})
	slices.SortStableFunc(shuffled, compareObservedVotes)
	assert.Equal(t, ordered, shuffled)
}

// The tallies are ordered by round, type, change-proposer round, value, then block hash.
func TestObservedTalliesOrder(t *testing.T) {
	lowHash := hash.UndefHash
	highHash := hash.CalcHash([]byte("block"))
	if slices.Compare(lowHash.Bytes(), highHash.Bytes()) > 0 {
		lowHash, highHash = highHash, lowHash
	}

	ordered := []ObservedTally{
		{Round: 0, Type: vote.VoteTypePrepare, CPRound: 0, CPValue: vote.CPValueNo, BlockHash: lowHash},
		{Round: 0, Type: vote.VoteTypePrepare, CPRound: 0, CPValue: vote.CPValueNo, BlockHash: highHash},
		{Round: 0, Type: vote.VoteTypePrepare, CPRound: 0, CPValue: vote.CPValueYes, BlockHash: lowHash},
		{Round: 0, Type: vote.VoteTypePrepare, CPRound: 1, CPValue: vote.CPValueNo, BlockHash: lowHash},
		{Round: 0, Type: vote.VoteTypePrecommit, CPRound: 0, CPValue: vote.CPValueNo, BlockHash: lowHash},
		{Round: 1, Type: vote.VoteTypePrepare, CPRound: 0, CPValue: vote.CPValueNo, BlockHash: lowHash},
	}

	for i := range ordered {
		assert.Zero(t, compareObservedTallies(ordered[i], ordered[i]))
		for j := i + 1; j < len(ordered); j++ {
			assert.Negative(t, compareObservedTallies(ordered[i], ordered[j]), "%d before %d", i, j)
			assert.Positive(t, compareObservedTallies(ordered[j], ordered[i]), "%d after %d", j, i)
		}
	}
}
