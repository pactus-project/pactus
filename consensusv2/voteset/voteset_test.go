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

	valsMap, valKeys, totalPower := setupCommittee(ts, 1, 1, 1, 1, 1, 1)

	hash1 := ts.RandHash()
	hash2 := ts.RandHash()
	height := ts.RandHeight()
	round := ts.RandRound()
	invKey := ts.RandValKey()
	valKey := valKeys[0]
	voteSet := NewPrepareVoteSet(round, totalPower, valsMap)
	assert.Equal(t, round, voteSet.Round())

	vote1 := vote.NewPrepareVote(hash1, height, round, invKey.Address())
	vote2 := vote.NewPrepareVote(hash1, height, round, valKey.Address())
	vote3 := vote.NewPrepareVote(hash2, height, round, valKey.Address())

	ts.HelperSignVote(invKey, vote1)
	added, err := voteSet.AddVote(vote1)
	assert.ErrorIs(t, err, IneligibleVoterError{Address: vote1.Signer()}) // unknown validator
	assert.False(t, added)
	assert.False(t, voteSet.HasVoted(vote1.Signer()))

	ts.HelperSignVote(invKey, vote2)
	added, err = voteSet.AddVote(vote2)
	assert.ErrorIs(t, err, crypto.ErrInvalidSignature)
	assert.False(t, added)
	assert.False(t, voteSet.HasVoted(vote2.Signer()))

	ts.HelperSignVote(valKey, vote2)
	added, err = voteSet.AddVote(vote2)
	assert.NoError(t, err) // ok
	assert.True(t, added)
	assert.True(t, voteSet.HasVoted(vote2.Signer()))

	added, err = voteSet.AddVote(vote2) // Adding again
	assert.False(t, added)
	assert.NoError(t, err)
	assert.True(t, voteSet.HasVoted(vote2.Signer()))

	ts.HelperSignVote(valKey, vote3)
	added, err = voteSet.AddVote(vote3)
	assert.ErrorIs(t, err, ErrDuplicatedVote)
	assert.True(t, added)
	assert.True(t, voteSet.HasVoted(vote3.Signer()))
}

func TestAddBinaryVote(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	valsMap, valKeys, totalPower := setupCommittee(ts, 1, 1, 1, 1)

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
	assert.ErrorIs(t, err, ErrDuplicatedVote)
	assert.True(t, added)
}

func TestDuplicateBlockVote(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	valsMap, valKeys, totalPower := setupCommittee(ts, 1, 1, 1, 1, 1, 1)

	hash1 := ts.RandHash()
	hash2 := ts.RandHash()
	hash3 := ts.RandHash()
	addr := valKeys[0].Address()
	voteSet := NewPrepareVoteSet(0, totalPower, valsMap)

	correctVote := vote.NewPrepareVote(hash1, 1, 0, addr)
	duplicatedVote1 := vote.NewPrepareVote(hash2, 1, 0, addr)
	duplicatedVote2 := vote.NewPrepareVote(hash3, 1, 0, addr)

	// sign the votes
	ts.HelperSignVote(valKeys[0], correctVote)
	ts.HelperSignVote(valKeys[0], duplicatedVote1)
	ts.HelperSignVote(valKeys[0], duplicatedVote2)

	added, err := voteSet.AddVote(correctVote)
	assert.NoError(t, err)
	assert.True(t, added)

	added, err = voteSet.AddVote(duplicatedVote1)
	assert.ErrorIs(t, err, ErrDuplicatedVote)
	assert.True(t, added)

	added, err = voteSet.AddVote(duplicatedVote2)
	assert.ErrorIs(t, err, ErrDuplicatedVote)
	assert.True(t, added)

	bv1 := voteSet.BlockVotes(hash1)
	bv2 := voteSet.BlockVotes(hash2)
	bv3 := voteSet.BlockVotes(hash3)
	assert.Equal(t, correctVote, bv1[addr])
	assert.Equal(t, duplicatedVote1, bv2[addr])
	assert.Equal(t, duplicatedVote2, bv3[addr])
	assert.False(t, voteSet.HasQuorumHash())

	assert.Contains(t, voteSet.AllVotes(), correctVote)
	assert.NotContains(t, voteSet.AllVotes(), duplicatedVote1)
	assert.NotContains(t, voteSet.AllVotes(), duplicatedVote2)
}

func TestDuplicateBinaryVote(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	valsMap, valKeys, totalPower := setupCommittee(ts, 1, 1, 1, 1)

	hash1 := ts.RandHash()
	hash2 := ts.RandHash()
	hash3 := ts.RandHash()
	addr := valKeys[0].Address()
	voteSet := NewCPPreVoteVoteSet(0, totalPower, valsMap)

	correctVote := vote.NewCPPreVote(hash1, 1, 0, 0, vote.CPValueYes, &vote.JustInitYes{}, addr)
	duplicatedVote1 := vote.NewCPPreVote(hash2, 1, 0, 0, vote.CPValueYes, &vote.JustInitYes{}, addr)
	duplicatedVote2 := vote.NewCPPreVote(hash3, 1, 0, 0, vote.CPValueYes, &vote.JustInitYes{}, addr)

	// sign the votes
	ts.HelperSignVote(valKeys[0], correctVote)
	ts.HelperSignVote(valKeys[0], duplicatedVote1)
	ts.HelperSignVote(valKeys[0], duplicatedVote2)

	added, err := voteSet.AddVote(correctVote)
	assert.NoError(t, err)
	assert.True(t, added)

	added, err = voteSet.AddVote(duplicatedVote1)
	assert.ErrorIs(t, err, ErrDuplicatedVote)
	assert.True(t, added)

	added, err = voteSet.AddVote(duplicatedVote2)
	assert.ErrorIs(t, err, ErrDuplicatedVote)
	assert.True(t, added)

	assert.False(t, voteSet.HasFPlusOneVotesFor(0, vote.CPValueNo))

	assert.Contains(t, voteSet.AllVotes(), correctVote)
	assert.NotContains(t, voteSet.AllVotes(), duplicatedVote1)
	assert.NotContains(t, voteSet.AllVotes(), duplicatedVote2)
}

func TestQuorum(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	// N = 4501
	// f = 1500
	// 2f+1 = 3001
	// 3f+1 = 4501
	valsMap, valKeys, totalPower := setupCommittee(ts, 1000, 900, 801, 700, 600, 500)

	voteSet := NewPrepareVoteSet(0, totalPower, valsMap)
	blockHash := ts.RandHash()
	vote1 := vote.NewPrepareVote(blockHash, 1, 0, valKeys[0].Address())
	vote2 := vote.NewPrepareVote(blockHash, 1, 0, valKeys[1].Address())
	vote3 := vote.NewPrepareVote(blockHash, 1, 0, valKeys[2].Address())
	vote4 := vote.NewPrepareVote(blockHash, 1, 0, valKeys[3].Address())
	vote5 := vote.NewPrepareVote(blockHash, 1, 0, valKeys[4].Address())
	vote6 := vote.NewPrepareVote(blockHash, 1, 0, valKeys[5].Address())

	ts.HelperSignVote(valKeys[0], vote1)
	ts.HelperSignVote(valKeys[1], vote2)
	ts.HelperSignVote(valKeys[2], vote3)
	ts.HelperSignVote(valKeys[3], vote4)
	ts.HelperSignVote(valKeys[4], vote5)
	ts.HelperSignVote(valKeys[5], vote6)

	_, err := voteSet.AddVote(vote1)
	assert.NoError(t, err)
	_, err = voteSet.AddVote(vote2)
	assert.NoError(t, err)

	assert.True(t, voteSet.HasQuorumHash())
	assert.Equal(t, &blockHash, voteSet.QuorumHash())

	// Add more votes
	_, err = voteSet.AddVote(vote3)
	assert.NoError(t, err)

	assert.False(t, voteSet.HasMajorityQuorum(blockHash))
	assert.False(t, voteSet.HasAbsoluteQuorum())

	// Add more votes
	_, err = voteSet.AddVote(vote4)
	assert.NoError(t, err)
	_, err = voteSet.AddVote(vote5)
	assert.NoError(t, err)

	assert.True(t, voteSet.HasMajorityQuorum(blockHash))
	assert.False(t, voteSet.HasAbsoluteQuorum())

	// Add more votes
	_, err = voteSet.AddVote(vote6)
	assert.NoError(t, err)
	assert.True(t, voteSet.HasMajorityQuorum(blockHash))
	assert.True(t, voteSet.HasAbsoluteQuorum())
}

func TestAllBlockVotes(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	valsMap, valKeys, totalPower := setupCommittee(ts, 1, 1, 1, 1, 1, 1)

	voteSet := NewPrecommitVoteSet(1, totalPower, valsMap)

	hash1 := ts.RandHash()
	vote1 := vote.NewPrecommitVote(hash1, 1, 1, valKeys[0].Address())
	vote2 := vote.NewPrecommitVote(hash1, 1, 1, valKeys[1].Address())
	vote3 := vote.NewPrecommitVote(hash1, 1, 1, valKeys[2].Address())
	vote4 := vote.NewPrecommitVote(hash1, 1, 1, valKeys[3].Address())

	ts.HelperSignVote(valKeys[0], vote1)
	ts.HelperSignVote(valKeys[1], vote2)
	ts.HelperSignVote(valKeys[2], vote3)
	ts.HelperSignVote(valKeys[3], vote4)

	_, err := voteSet.AddVote(vote1)
	assert.NoError(t, err)

	_, err = voteSet.AddVote(vote2)
	assert.NoError(t, err)

	_, err = voteSet.AddVote(vote3)
	assert.NoError(t, err)

	_, err = voteSet.AddVote(vote4)
	assert.NoError(t, err)

	assert.Equal(t, &hash1, voteSet.QuorumHash())

	// Check accumulated power
	assert.Equal(t, &hash1, voteSet.QuorumHash())

	// Check previous votes
	assert.Contains(t, voteSet.AllVotes(), vote1)
	assert.Contains(t, voteSet.AllVotes(), vote2)
	assert.Contains(t, voteSet.AllVotes(), vote3)
	assert.Contains(t, voteSet.AllVotes(), vote4)
}

func TestAllBinaryVotes(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	valsMap, valKeys, totalPower := setupCommittee(ts, 1, 1, 1, 1)

	voteSet := NewCPMainVoteVoteSet(1, totalPower, valsMap)

	vote1 := vote.NewCPMainVote(hash.UndefHash, 1, 1, 0, vote.CPValueNo, &vote.JustInitYes{}, valKeys[0].Address())
	vote2 := vote.NewCPMainVote(hash.UndefHash, 1, 1, 1, vote.CPValueYes, &vote.JustInitYes{}, valKeys[1].Address())
	vote3 := vote.NewCPMainVote(hash.UndefHash, 1, 1, 2, vote.CPValueAbstain, &vote.JustInitYes{}, valKeys[2].Address())

	ts.HelperSignVote(valKeys[0], vote1)
	ts.HelperSignVote(valKeys[1], vote2)
	ts.HelperSignVote(valKeys[2], vote3)

	assert.Empty(t, voteSet.AllVotes())

	_, err := voteSet.AddVote(vote1)
	assert.NoError(t, err)

	_, err = voteSet.AddVote(vote2)
	assert.NoError(t, err)

	_, err = voteSet.AddVote(vote3)
	assert.NoError(t, err)

	assert.Contains(t, voteSet.AllVotes(), vote1)
	assert.Contains(t, voteSet.AllVotes(), vote2)
	assert.Contains(t, voteSet.AllVotes(), vote3)

	ranVote1 := voteSet.GetRandomVote(1, vote.CPValueNo)
	assert.Nil(t, ranVote1)

	ranVote2 := voteSet.GetRandomVote(1, vote.CPValueYes)
	assert.Equal(t, vote2, ranVote2)
}

func TestOneThirdPower(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	// N = 3001
	// f = 1000
	// f+1 = 1001
	// 2f+1 = 2001
	// 3f+1 = 3001
	valsMap, valKeys, totalPower := setupCommittee(ts, 1000, 1, 1000, 1000)

	hash := ts.RandHash()
	height := ts.RandHeight()
	round := ts.RandRound()
	just := &vote.JustInitYes{}
	voteSet := NewCPPreVoteVoteSet(round, totalPower, valsMap)

	vote1 := vote.NewCPPreVote(hash, height, round, 0, vote.CPValueYes, just, valKeys[0].Address())
	vote2 := vote.NewCPPreVote(hash, height, round, 0, vote.CPValueYes, just, valKeys[1].Address())
	vote3 := vote.NewCPPreVote(hash, height, round, 0, vote.CPValueYes, just, valKeys[2].Address())
	vote4 := vote.NewCPPreVote(hash, height, round, 0, vote.CPValueNo, just, valKeys[3].Address())

	ts.HelperSignVote(valKeys[0], vote1)
	ts.HelperSignVote(valKeys[1], vote2)
	ts.HelperSignVote(valKeys[2], vote3)
	ts.HelperSignVote(valKeys[3], vote4)

	_, err := voteSet.AddVote(vote1)
	assert.NoError(t, err)
	assert.False(t, voteSet.HasFPlusOneVotesFor(0, vote.CPValueNo))
	assert.True(t, voteSet.HasAnyVoteFor(0, vote.CPValueYes))
	assert.False(t, voteSet.HasAnyVoteFor(0, vote.CPValueNo))
	assert.False(t, voteSet.HasAnyVoteFor(0, vote.CPValueAbstain))

	_, err = voteSet.AddVote(vote2)
	assert.NoError(t, err)
	assert.True(t, voteSet.HasFPlusOneVotesFor(0, vote.CPValueYes))
	assert.False(t, voteSet.HasTwoFPlusOneVotes(0))

	_, err = voteSet.AddVote(vote3)
	assert.NoError(t, err)
	assert.True(t, voteSet.HasTwoFPlusOneVotes(0))
	assert.False(t, voteSet.HasAnyVoteFor(0, vote.CPValueNo))
	assert.True(t, voteSet.HasAnyVoteFor(0, vote.CPValueYes))
	assert.False(t, voteSet.HasTwoFPlusOneVotesFor(0, vote.CPValueNo))
	assert.True(t, voteSet.HasTwoFPlusOneVotesFor(0, vote.CPValueYes))
	assert.True(t, voteSet.HasAllVotesFor(0, vote.CPValueYes))

	_, err = voteSet.AddVote(vote4)
	assert.NoError(t, err)
	assert.True(t, voteSet.HasAnyVoteFor(0, vote.CPValueNo))
	assert.False(t, voteSet.HasTwoFPlusOneVotesFor(0, vote.CPValueNo))
	assert.True(t, voteSet.HasTwoFPlusOneVotesFor(0, vote.CPValueYes))
	assert.False(t, voteSet.HasAllVotesFor(0, vote.CPValueYes))

	bv1 := voteSet.BinaryVotes(0, vote.CPValueYes)
	bv2 := voteSet.BinaryVotes(0, vote.CPValueNo)

	assert.Contains(t, bv1, vote1.Signer())
	assert.Contains(t, bv1, vote2.Signer())
	assert.Contains(t, bv1, vote3.Signer())
	assert.Contains(t, bv2, vote4.Signer())
}

func TestDecidedVoteset(t *testing.T) {
	ts := testsuite.NewTestSuite(t)
	valsMap, valKeys, totalPower := setupCommittee(ts, 1, 1, 1, 1)

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
