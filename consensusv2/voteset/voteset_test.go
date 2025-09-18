package voteset

import (
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func setupCommittee(ts *testsuite.TestSuite, stakes ...amount.Amount) (
	map[crypto.Address]*validator.Validator, []*bls.ValidatorKey, int64,
) {
	valKeys := []*bls.ValidatorKey{}
	valsMap := map[crypto.Address]*validator.Validator{}
	totalPower := int64(0)
	for i, s := range stakes {
		pub, prv := ts.RandBLSKeyPair()
		val := validator.NewValidator(pub, int32(i))
		val.AddToStake(s)
		valsMap[val.Address()] = val
		totalPower += val.Power()
		valKeys = append(valKeys, bls.NewValidatorKey(prv))
	}

	return valsMap, valKeys, totalPower
}

func TestAddBlockVote(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	valsMap, valKeys, totalPower := setupCommittee(ts, 5, 6, 8, 11)

	hash1 := ts.RandHash()
	hash2 := ts.RandHash()
	height := ts.RandHeight()
	round := ts.RandRound()
	invKey := ts.RandValKey()
	valKey := valKeys[0]
	voteSet := NewPrecommitVoteSet(round, totalPower, valsMap)
	assert.Equal(t, round, voteSet.Round())

	vote1 := vote.NewPrecommitVote(hash1, height, round, invKey.Address())
	vote2 := vote.NewPrecommitVote(hash1, height, round, valKey.Address())
	vote3 := vote.NewPrecommitVote(hash2, height, round, valKey.Address())

	ts.HelperSignVote(invKey, vote1)
	added, err := voteSet.AddVote(vote1)
	assert.ErrorIs(t, err, IneligibleVoterError{Address: vote1.Signer()}) // unknown validator
	assert.False(t, added)

	ts.HelperSignVote(invKey, vote2)
	added, err = voteSet.AddVote(vote2)
	assert.ErrorIs(t, err, crypto.ErrInvalidSignature)
	assert.False(t, added)

	ts.HelperSignVote(valKey, vote2)
	added, err = voteSet.AddVote(vote2)
	assert.NoError(t, err) // ok
	assert.True(t, added)

	added, err = voteSet.AddVote(vote2) // Adding again
	assert.False(t, added)
	assert.NoError(t, err)

	ts.HelperSignVote(valKey, vote3)
	added, err = voteSet.AddVote(vote3)
	assert.ErrorIs(t, err, ErrDoubleVote)
	assert.True(t, added)
}

func TestAddBinaryVote(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	valsMap, valKeys, totalPower := setupCommittee(ts, 5, 6, 8, 11)

	hash1 := ts.RandHash()
	hash2 := ts.RandHash()
	height := ts.RandHeight()
	round := ts.RandRound()
	cpRound := ts.RandRound()
	cpVal := ts.RandInt(2)
	just := &vote.JustInitYes{}
	invKey := ts.RandValKey()
	valKey := valKeys[0]
	voteSet := NewCPPreVoteVoteSet(round, totalPower, valsMap)

	vote1 := vote.NewCPPreVote(hash1, height, round, cpRound, vote.CPValue(cpVal), just, invKey.Address())
	vote2 := vote.NewCPPreVote(hash1, height, round, cpRound, vote.CPValue(cpVal), just, valKey.Address())
	vote3 := vote.NewCPPreVote(hash2, height, round, cpRound, vote.CPValue(cpVal), just, valKey.Address())

	ts.HelperSignVote(invKey, vote1)
	added, err := voteSet.AddVote(vote1)
	assert.ErrorIs(t, err, IneligibleVoterError{Address: vote1.Signer()}) // unknown validator
	assert.False(t, added)

	ts.HelperSignVote(invKey, vote2)
	added, err = voteSet.AddVote(vote2)
	assert.ErrorIs(t, err, crypto.ErrInvalidSignature)
	assert.False(t, added)

	ts.HelperSignVote(valKey, vote2)
	added, err = voteSet.AddVote(vote2)
	assert.NoError(t, err) // ok
	assert.True(t, added)

	added, err = voteSet.AddVote(vote2) // Adding again
	assert.False(t, added)
	assert.NoError(t, err)

	ts.HelperSignVote(valKey, vote3)
	added, err = voteSet.AddVote(vote3)
	assert.ErrorIs(t, err, ErrDoubleVote)
	assert.True(t, added)
}

func TestDoubleBlockVote(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	valsMap, valKeys, totalPower := setupCommittee(ts, 5, 6, 8, 11)

	hash1 := ts.RandHash()
	hash2 := ts.RandHash()
	hash3 := ts.RandHash()
	addr := valKeys[0].Address()
	voteSet := NewPrecommitVoteSet(0, totalPower, valsMap)

	correctVote := vote.NewPrecommitVote(hash1, 1, 0, addr)
	doubleVote1 := vote.NewPrecommitVote(hash2, 1, 0, addr)
	doubleVote2 := vote.NewPrecommitVote(hash3, 1, 0, addr)

	// sign the votes
	ts.HelperSignVote(valKeys[0], correctVote)
	ts.HelperSignVote(valKeys[0], doubleVote1)
	ts.HelperSignVote(valKeys[0], doubleVote2)

	added, err := voteSet.AddVote(correctVote)
	assert.NoError(t, err)
	assert.True(t, added)

	added, err = voteSet.AddVote(doubleVote1)
	assert.ErrorIs(t, err, ErrDoubleVote)
	assert.True(t, added)

	added, err = voteSet.AddVote(doubleVote2)
	assert.ErrorIs(t, err, ErrDoubleVote)
	assert.True(t, added)

	assert.Contains(t, voteSet.AllVotes(), correctVote)
	assert.NotContains(t, voteSet.AllVotes(), doubleVote1)
	assert.NotContains(t, voteSet.AllVotes(), doubleVote2)
}

func TestDoubleBinaryVote(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	valsMap, valKeys, totalPower := setupCommittee(ts, 5, 6, 8, 11)

	hash1 := ts.RandHash()
	hash2 := ts.RandHash()
	hash3 := ts.RandHash()
	addr := valKeys[0].Address()
	voteSet := NewCPPreVoteVoteSet(0, totalPower, valsMap)

	correctVote := vote.NewCPPreVote(hash1, 1, 0, 0, vote.CPValueYes, &vote.JustInitYes{}, addr)
	doubleVote1 := vote.NewCPPreVote(hash2, 1, 0, 0, vote.CPValueYes, &vote.JustInitYes{}, addr)
	doubleVote2 := vote.NewCPPreVote(hash3, 1, 0, 0, vote.CPValueYes, &vote.JustInitYes{}, addr)

	// sign the votes
	ts.HelperSignVote(valKeys[0], correctVote)
	ts.HelperSignVote(valKeys[0], doubleVote1)
	ts.HelperSignVote(valKeys[0], doubleVote2)

	added, err := voteSet.AddVote(correctVote)
	assert.NoError(t, err)
	assert.True(t, added)

	added, err = voteSet.AddVote(doubleVote1)
	assert.ErrorIs(t, err, ErrDoubleVote)
	assert.True(t, added)

	added, err = voteSet.AddVote(doubleVote2)
	assert.ErrorIs(t, err, ErrDoubleVote)
	assert.True(t, added)

	assert.Contains(t, voteSet.AllVotes(), correctVote)
	assert.NotContains(t, voteSet.AllVotes(), doubleVote1)
	assert.NotContains(t, voteSet.AllVotes(), doubleVote2)
}

func TestAllBlockVotes(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	valsMap, valKeys, totalPower := setupCommittee(ts, 5, 6, 8, 11)

	voteSet := NewPrecommitVoteSet(1, totalPower, valsMap)

	vote1 := vote.NewPrecommitVote(ts.RandHash(), 1, 1, valKeys[0].Address())
	vote2 := vote.NewPrecommitVote(ts.RandHash(), 1, 1, valKeys[1].Address())
	vote3 := vote.NewPrecommitVote(ts.RandHash(), 1, 1, valKeys[2].Address())

	ts.HelperSignVote(valKeys[0], vote1)
	ts.HelperSignVote(valKeys[1], vote2)
	ts.HelperSignVote(valKeys[2], vote3)

	assert.Empty(t, voteSet.AllVotes())

	_, _ = voteSet.AddVote(vote1)
	_, _ = voteSet.AddVote(vote2)
	_, _ = voteSet.AddVote(vote3)

	assert.Contains(t, voteSet.AllVotes(), vote1)
	assert.Contains(t, voteSet.AllVotes(), vote2)
	assert.Contains(t, voteSet.AllVotes(), vote3)
}

func TestAllBinaryVotes(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	valsMap, valKeys, totalPower := setupCommittee(ts, 5, 6, 8, 11)

	voteSet := NewCPMainVoteVoteSet(1, totalPower, valsMap)

	vote1 := vote.NewCPMainVote(hash.UndefHash, 1, 1, 0, vote.CPValueNo, &vote.JustInitYes{}, valKeys[0].Address())
	vote2 := vote.NewCPMainVote(hash.UndefHash, 1, 1, 1, vote.CPValueYes, &vote.JustInitYes{}, valKeys[1].Address())
	vote3 := vote.NewCPMainVote(hash.UndefHash, 1, 1, 2, vote.CPValueAbstain, &vote.JustInitYes{}, valKeys[2].Address())

	ts.HelperSignVote(valKeys[0], vote1)
	ts.HelperSignVote(valKeys[1], vote2)
	ts.HelperSignVote(valKeys[2], vote3)

	assert.Empty(t, voteSet.AllVotes())

	_, _ = voteSet.AddVote(vote1)
	_, _ = voteSet.AddVote(vote2)
	_, _ = voteSet.AddVote(vote3)

	assert.Contains(t, voteSet.AllVotes(), vote1)
	assert.Contains(t, voteSet.AllVotes(), vote2)
	assert.Contains(t, voteSet.AllVotes(), vote3)

	ranVote1 := voteSet.GetRandomVote(1, vote.CPValueNo)
	assert.Nil(t, ranVote1)

	ranVote2 := voteSet.GetRandomVote(1, vote.CPValueYes)
	assert.Equal(t, vote2, ranVote2)
}

func TestBlockQuorumVotes(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	// N = 30 (5+6+8+11)
	// f = 9
	// 1f+1 = 10
	// 2f+1 = 19
	// 3f+1 = 28
	valsMap, valKeys, totalPower := setupCommittee(ts, 5, 6, 8, 11)

	voteSet := NewPrecommitVoteSet(0, totalPower, valsMap)
	blockHash := ts.RandHash()
	vote1 := vote.NewPrecommitVote(blockHash, 1, 0, valKeys[0].Address())
	vote2 := vote.NewPrecommitVote(blockHash, 1, 0, valKeys[1].Address())
	vote3 := vote.NewPrecommitVote(blockHash, 1, 0, valKeys[2].Address())
	vote4 := vote.NewPrecommitVote(blockHash, 1, 0, valKeys[3].Address())

	ts.HelperSignVote(valKeys[0], vote1)
	ts.HelperSignVote(valKeys[1], vote2)
	ts.HelperSignVote(valKeys[2], vote3)
	ts.HelperSignVote(valKeys[3], vote4)

	_, _ = voteSet.AddVote(vote1)
	assert.False(t, voteSet.Has1FP1VotesFor(blockHash))

	// Add more votes
	_, _ = voteSet.AddVote(vote2)
	assert.True(t, voteSet.Has1FP1VotesFor(blockHash))
	assert.False(t, voteSet.Has2FP1VotesFor(blockHash))

	// Add more votes
	_, _ = voteSet.AddVote(vote3)
	assert.True(t, voteSet.Has2FP1VotesFor(blockHash))
	assert.False(t, voteSet.Has3FP1VotesFor(blockHash))

	// Add more votes
	_, _ = voteSet.AddVote(vote4)
	assert.True(t, voteSet.Has3FP1VotesFor(blockHash))
}

func TestBinaryQuorumVotes(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	// N = 30 (5+6+8+11)
	// f = 9
	// 1f+1 = 10
	// 2f+1 = 19
	// 3f+1 = 28
	valsMap, valKeys, totalPower := setupCommittee(ts, 5, 6, 8, 11)

	hash := ts.RandHash()
	height := ts.RandHeight()
	round := ts.RandRound()
	just := &vote.JustInitYes{}
	voteSet := NewCPPreVoteVoteSet(round, totalPower, valsMap)

	vote1 := vote.NewCPPreVote(hash, height, round, 0, vote.CPValueYes, just, valKeys[0].Address())
	vote2 := vote.NewCPPreVote(hash, height, round, 0, vote.CPValueYes, just, valKeys[1].Address())
	vote3 := vote.NewCPPreVote(hash, height, round, 0, vote.CPValueYes, just, valKeys[2].Address())
	vote4 := vote.NewCPPreVote(hash, height, round, 0, vote.CPValueYes, just, valKeys[3].Address())

	ts.HelperSignVote(valKeys[0], vote1)
	ts.HelperSignVote(valKeys[1], vote2)
	ts.HelperSignVote(valKeys[2], vote3)
	ts.HelperSignVote(valKeys[3], vote4)

	_, _ = voteSet.AddVote(vote1)

	assert.False(t, voteSet.Has1FP1VotesFor(0, vote.CPValueNo))
	assert.False(t, voteSet.Has1FP1VotesFor(0, vote.CPValueYes))
	assert.True(t, voteSet.HasAnyVoteFor(0, vote.CPValueYes))
	assert.False(t, voteSet.HasAnyVoteFor(0, vote.CPValueNo))
	assert.False(t, voteSet.HasAnyVoteFor(0, vote.CPValueAbstain))

	// Add more votes
	_, _ = voteSet.AddVote(vote2)

	assert.True(t, voteSet.Has1FP1VotesFor(0, vote.CPValueYes))
	assert.False(t, voteSet.Has2FP1VotesFor(0, vote.CPValueYes))

	// Add more votes
	_, _ = voteSet.AddVote(vote3)

	assert.True(t, voteSet.Has2FP1VotesFor(0, vote.CPValueYes))
	assert.False(t, voteSet.HasAllVotesFor(0, vote.CPValueNo))
	assert.True(t, voteSet.HasAllVotesFor(0, vote.CPValueYes))

	// Add more votes
	_, _ = voteSet.AddVote(vote4)

	assert.True(t, voteSet.HasAllVotesFor(0, vote.CPValueYes))
}

func TestBlockVotes(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	valsMap, valKeys, totalPower := setupCommittee(ts, 5, 6, 8, 11)

	hash := ts.RandHash()
	height := ts.RandHeight()
	round := ts.RandRound()
	voteSet := NewPrecommitVoteSet(round, totalPower, valsMap)

	vote1 := vote.NewPrecommitVote(hash, height, round, valKeys[0].Address())
	vote2 := vote.NewPrecommitVote(hash, height, round, valKeys[1].Address())

	ts.HelperSignVote(valKeys[0], vote1)
	ts.HelperSignVote(valKeys[1], vote2)

	_, _ = voteSet.AddVote(vote1)
	_, _ = voteSet.AddVote(vote2)

	bv := voteSet.BlockVotes(hash)
	assert.Contains(t, bv, vote1.Signer())
	assert.Contains(t, bv, vote2.Signer())
}

func TestBinaryVotes(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	valsMap, valKeys, totalPower := setupCommittee(ts, 5, 6, 8, 11)

	hash := ts.RandHash()
	height := ts.RandHeight()
	round := ts.RandRound()
	just := &vote.JustInitYes{}
	voteSet := NewCPPreVoteVoteSet(round, totalPower, valsMap)

	vote1 := vote.NewCPPreVote(hash, height, round, 0, vote.CPValueYes, just, valKeys[0].Address())
	vote2 := vote.NewCPPreVote(hash, height, round, 0, vote.CPValueNo, just, valKeys[1].Address())

	ts.HelperSignVote(valKeys[0], vote1)
	ts.HelperSignVote(valKeys[1], vote2)

	_, _ = voteSet.AddVote(vote1)
	_, _ = voteSet.AddVote(vote2)

	bv1 := voteSet.BinaryVotes(0, vote.CPValueYes)
	assert.Contains(t, bv1, vote1.Signer())
	assert.NotContains(t, bv1, vote2.Signer())

	bv2 := voteSet.BinaryVotes(0, vote.CPValueNo)
	assert.NotContains(t, bv2, vote1.Signer())
	assert.Contains(t, bv2, vote2.Signer())
}

func TestDecidedVoteset(t *testing.T) {
	ts := testsuite.NewTestSuite(t)
	valsMap, valKeys, totalPower := setupCommittee(ts, 5, 6, 8, 11)

	hash := ts.RandHash()
	height := ts.RandHeight()
	round := ts.RandRound()
	just := &vote.JustInitYes{}
	voteSet := NewCPDecidedVoteSet(round, totalPower, valsMap)

	vte := vote.NewCPDecidedVote(hash, height, round, 0, vote.CPValueYes, just, valKeys[0].Address())

	ts.HelperSignVote(valKeys[0], vte)

	_, err := voteSet.AddVote(vte)
	assert.NoError(t, err)
	assert.True(t, voteSet.HasAnyVoteFor(0, vote.CPValueYes))
	assert.False(t, voteSet.HasAnyVoteFor(0, vote.CPValueNo))
}

func TestBlockVotedPower(t *testing.T) {
	ts := testsuite.NewTestSuite(t)
	valsMap, valKeys, totalPower := setupCommittee(ts, 5, 6, 8, 11)

	height := ts.RandHeight()
	round := ts.RandRound()
	voteSet := NewPrecommitVoteSet(round, totalPower, valsMap)

	hash1 := ts.RandHash()
	hash2 := ts.RandHash()
	vote1a := vote.NewPrecommitVote(hash1, height, round, valKeys[0].Address())
	vote1b := vote.NewPrecommitVote(hash2, height, round, valKeys[0].Address())
	vote2 := vote.NewPrecommitVote(hash2, height, round, valKeys[1].Address())

	ts.HelperSignVote(valKeys[0], vote1a)
	ts.HelperSignVote(valKeys[0], vote1b)
	ts.HelperSignVote(valKeys[1], vote2)

	_, err := voteSet.AddVote(vote1a)
	assert.NoError(t, err)

	_, err = voteSet.AddVote(vote1b) // Double vote
	assert.ErrorIs(t, err, ErrDoubleVote)

	_, err = voteSet.AddVote(vote2)
	assert.NoError(t, err)

	assert.Equal(t, int64(5+6), voteSet.VotedPower())
}

func TestBinaryVotedPower(t *testing.T) {
	ts := testsuite.NewTestSuite(t)
	valsMap, valKeys, totalPower := setupCommittee(ts, 5, 6, 8, 11)

	height := ts.RandHeight()
	round := ts.RandRound()
	voteSet := NewCPPreVoteVoteSet(1, totalPower, valsMap)

	just := &vote.JustInitYes{}
	hash1 := ts.RandHash()
	hash2 := ts.RandHash()
	vote1a := vote.NewCPPreVote(hash1, height, round, 0, vote.CPValueYes, just, valKeys[0].Address())
	vote1b := vote.NewCPPreVote(hash2, height, round, 0, vote.CPValueYes, just, valKeys[0].Address())
	vote1c := vote.NewCPPreVote(hash1, height, round, 1, vote.CPValueYes, just, valKeys[0].Address())
	vote2 := vote.NewCPPreVote(hash1, height, round, 0, vote.CPValueYes, just, valKeys[1].Address())

	ts.HelperSignVote(valKeys[0], vote1a)
	ts.HelperSignVote(valKeys[0], vote1b)
	ts.HelperSignVote(valKeys[0], vote1c)
	ts.HelperSignVote(valKeys[1], vote2)

	_, _err := voteSet.AddVote(vote1a)
	assert.NoError(t, _err)

	_, _err = voteSet.AddVote(vote1b) // Double vote
	assert.ErrorIs(t, _err, ErrDoubleVote)

	_, _err = voteSet.AddVote(vote1c) // Next CP:Round
	assert.NoError(t, _err)

	_, _err = voteSet.AddVote(vote2)
	assert.NoError(t, _err)

	assert.Equal(t, int64(5+6), voteSet.VotedPower(0))
	assert.Equal(t, int64(5), voteSet.VotedPower(1))
}

func TestBlockHas2FP1Votes(t *testing.T) {
	ts := testsuite.NewTestSuite(t)
	valsMap, valKeys, totalPower := setupCommittee(ts, 5, 6, 8, 11)

	voteSet := NewPrecommitVoteSet(0, totalPower, valsMap)
	blockHash := ts.RandHash()
	vote0 := vote.NewPrecommitVote(ts.RandHash(), 1, 0, valKeys[0].Address()) // Byzantine vote
	vote1 := vote.NewPrecommitVote(blockHash, 1, 0, valKeys[0].Address())
	vote2 := vote.NewPrecommitVote(blockHash, 1, 0, valKeys[1].Address())
	vote3 := vote.NewPrecommitVote(blockHash, 1, 0, valKeys[2].Address())

	ts.HelperSignVote(valKeys[0], vote0)
	ts.HelperSignVote(valKeys[0], vote1)
	ts.HelperSignVote(valKeys[1], vote2)
	ts.HelperSignVote(valKeys[2], vote3)

	_, _ = voteSet.AddVote(vote0)
	_, _ = voteSet.AddVote(vote1)
	_, _ = voteSet.AddVote(vote2)
	assert.False(t, voteSet.Has2FP1Votes())

	_, _ = voteSet.AddVote(vote3)
	assert.True(t, voteSet.Has2FP1Votes())
}

func TestBinaryHas3FP1Votes(t *testing.T) {
	ts := testsuite.NewTestSuite(t)
	valsMap, valKeys, totalPower := setupCommittee(ts, 5, 6, 8, 11)

	hash := ts.RandHash()
	height := ts.RandHeight()
	round := ts.RandRound()
	just := &vote.JustInitYes{}
	voteSet := NewCPPreVoteVoteSet(round, totalPower, valsMap)

	vote0 := vote.NewCPPreVote(hash, height, round, 0, vote.CPValueNo, just, valKeys[0].Address()) // Byzantine vote
	vote1 := vote.NewCPPreVote(hash, height, round, 0, vote.CPValueYes, just, valKeys[0].Address())
	vote2 := vote.NewCPPreVote(hash, height, round, 0, vote.CPValueYes, just, valKeys[1].Address())
	vote3 := vote.NewCPPreVote(hash, height, round, 0, vote.CPValueYes, just, valKeys[2].Address())

	ts.HelperSignVote(valKeys[0], vote0)
	ts.HelperSignVote(valKeys[0], vote1)
	ts.HelperSignVote(valKeys[1], vote2)
	ts.HelperSignVote(valKeys[2], vote3)

	_, _ = voteSet.AddVote(vote0)
	_, _ = voteSet.AddVote(vote1)
	_, _ = voteSet.AddVote(vote2)
	assert.False(t, voteSet.Has2FP1Votes(0))

	_, _ = voteSet.AddVote(vote3)
	assert.True(t, voteSet.Has2FP1Votes(0))
}
