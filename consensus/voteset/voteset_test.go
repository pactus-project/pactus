package voteset

import (
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/pactus-project/pactus/util/errors"
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

	valsMap, valKeys, totalPower := setupCommittee(ts, 1000, 1500, 2500, 2000)

	hash1 := ts.RandHash()
	hash2 := ts.RandHash()
	height := ts.RandHeight()
	round := ts.RandRound()
	invKey := ts.RandValKey()
	valKey := valKeys[0]
	vs := NewPrecommitVoteSet(round, totalPower, valsMap)

	v1 := vote.NewPrecommitVote(hash1, height, round, invKey.Address())
	v2 := vote.NewPrecommitVote(hash1, height, round, valKey.Address())
	v3 := vote.NewPrecommitVote(hash2, height, round, valKey.Address())

	ts.HelperSignVote(invKey, v1)
	added, err := vs.AddVote(v1)
	assert.Equal(t, errors.Code(err), errors.ErrInvalidAddress) // unknown validator
	assert.False(t, added)

	ts.HelperSignVote(invKey, v2)
	added, err = vs.AddVote(v2)
	assert.Equal(t, errors.Code(err), errors.ErrInvalidSignature) // invalid signature
	assert.False(t, added)

	ts.HelperSignVote(valKey, v2)
	added, err = vs.AddVote(v2)
	assert.NoError(t, err) // ok
	assert.True(t, added)

	added, err = vs.AddVote(v2) // Adding again
	assert.False(t, added)
	assert.NoError(t, err)

	ts.HelperSignVote(valKey, v3)
	added, err = vs.AddVote(v3)
	assert.Equal(t, errors.Code(err), errors.ErrDuplicateVote)
	assert.True(t, added)
}

func TestAddBinaryVote(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	valsMap, valKeys, totalPower := setupCommittee(ts, 1000, 1500, 2500, 2000)

	hash1 := ts.RandHash()
	hash2 := ts.RandHash()
	height := ts.RandHeight()
	round := ts.RandRound()
	cpRound := ts.RandRound()
	cpVal := ts.RandInt(2)
	just := &vote.JustInitOne{}
	invKey := ts.RandValKey()
	valKey := valKeys[0]
	vs := NewCPPreVoteVoteSet(round, totalPower, valsMap)

	v1 := vote.NewCPPreVote(hash1, height, round, cpRound, vote.CPValue(cpVal), just, invKey.Address())
	v2 := vote.NewCPPreVote(hash1, height, round, cpRound, vote.CPValue(cpVal), just, valKey.Address())
	v3 := vote.NewCPPreVote(hash2, height, round, cpRound, vote.CPValue(cpVal), just, valKey.Address())

	ts.HelperSignVote(invKey, v1)
	added, err := vs.AddVote(v1)
	assert.Equal(t, errors.Code(err), errors.ErrInvalidAddress) // unknown validator
	assert.False(t, added)

	ts.HelperSignVote(invKey, v2)
	added, err = vs.AddVote(v2)
	assert.Equal(t, errors.Code(err), errors.ErrInvalidSignature) // invalid signature
	assert.False(t, added)

	ts.HelperSignVote(valKey, v2)
	added, err = vs.AddVote(v2)
	assert.NoError(t, err) // ok
	assert.True(t, added)

	added, err = vs.AddVote(v2) // Adding again
	assert.False(t, added)
	assert.NoError(t, err)

	ts.HelperSignVote(valKey, v3)
	added, err = vs.AddVote(v3)
	assert.Equal(t, errors.Code(err), errors.ErrDuplicateVote)
	assert.True(t, added)
}

func TestDuplicateBlockVote(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	valsMap, valKeys, totalPower := setupCommittee(ts, 1, 1, 1, 1)

	h1 := ts.RandHash()
	h2 := ts.RandHash()
	h3 := ts.RandHash()
	addr := valKeys[0].Address()
	vs := NewPrepareVoteSet(0, totalPower, valsMap)

	correctVote := vote.NewPrepareVote(h1, 1, 0, addr)
	duplicatedVote1 := vote.NewPrepareVote(h2, 1, 0, addr)
	duplicatedVote2 := vote.NewPrepareVote(h3, 1, 0, addr)

	// sign the votes
	ts.HelperSignVote(valKeys[0], correctVote)
	ts.HelperSignVote(valKeys[0], duplicatedVote1)
	ts.HelperSignVote(valKeys[0], duplicatedVote2)

	added, err := vs.AddVote(correctVote)
	assert.NoError(t, err)
	assert.True(t, added)

	added, err = vs.AddVote(duplicatedVote1)
	assert.Equal(t, errors.Code(err), errors.ErrDuplicateVote)
	assert.True(t, added)

	added, err = vs.AddVote(duplicatedVote2)
	assert.Equal(t, errors.Code(err), errors.ErrDuplicateVote)
	assert.True(t, added)

	bv1 := vs.BlockVotes(h1)
	bv2 := vs.BlockVotes(h2)
	bv3 := vs.BlockVotes(h3)
	assert.Equal(t, bv1[addr], correctVote)
	assert.Equal(t, bv2[addr], duplicatedVote1)
	assert.Equal(t, bv3[addr], duplicatedVote2)
	assert.False(t, vs.HasQuorumHash())
}

func TestDuplicateBinaryVote(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	valsMap, valKeys, totalPower := setupCommittee(ts, 1, 1, 1, 1)

	h1 := ts.RandHash()
	h2 := ts.RandHash()
	h3 := ts.RandHash()
	addr := valKeys[0].Address()
	vs := NewCPPreVoteVoteSet(0, totalPower, valsMap)

	correctVote := vote.NewCPPreVote(h1, 1, 0, 0, vote.CPValueOne, &vote.JustInitOne{}, addr)
	duplicatedVote1 := vote.NewCPPreVote(h2, 1, 0, 0, vote.CPValueOne, &vote.JustInitOne{}, addr)
	duplicatedVote2 := vote.NewCPPreVote(h3, 1, 0, 0, vote.CPValueOne, &vote.JustInitOne{}, addr)

	// sign the votes
	ts.HelperSignVote(valKeys[0], correctVote)
	ts.HelperSignVote(valKeys[0], duplicatedVote1)
	ts.HelperSignVote(valKeys[0], duplicatedVote2)

	added, err := vs.AddVote(correctVote)
	assert.NoError(t, err)
	assert.True(t, added)

	added, err = vs.AddVote(duplicatedVote1)
	assert.Equal(t, errors.Code(err), errors.ErrDuplicateVote)
	assert.True(t, added)

	added, err = vs.AddVote(duplicatedVote2)
	assert.Equal(t, errors.Code(err), errors.ErrDuplicateVote)
	assert.True(t, added)

	assert.False(t, vs.HasOneThirdOfTotalPower(0))
}

func TestQuorum(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	valsMap, valKeys, totalPower := setupCommittee(ts, 1000, 1500, 2500, 2000)

	vs := NewPrecommitVoteSet(0, totalPower, valsMap)
	blockHash := ts.RandHash()
	v1 := vote.NewPrecommitVote(blockHash, 1, 0, valKeys[0].Address())
	v2 := vote.NewPrecommitVote(blockHash, 1, 0, valKeys[1].Address())
	v3 := vote.NewPrecommitVote(blockHash, 1, 0, valKeys[2].Address())
	v4 := vote.NewPrecommitVote(blockHash, 1, 0, valKeys[3].Address())

	ts.HelperSignVote(valKeys[0], v1)
	ts.HelperSignVote(valKeys[1], v2)
	ts.HelperSignVote(valKeys[2], v3)
	ts.HelperSignVote(valKeys[3], v4)

	_, err := vs.AddVote(v1)
	assert.NoError(t, err)

	_, err = vs.AddVote(v2)
	assert.NoError(t, err)

	assert.Nil(t, vs.QuorumHash())
	assert.False(t, vs.HasQuorumHash())
	assert.Contains(t, vs.BlockVotes(blockHash), v1.Signer())
	assert.Contains(t, vs.BlockVotes(blockHash), v2.Signer())

	_, err = vs.AddVote(v3)
	assert.NoError(t, err)

	assert.True(t, vs.HasQuorumHash())
	assert.Contains(t, vs.BlockVotes(blockHash), v3.Signer())
	assert.NotContains(t, vs.BlockVotes(blockHash), v4.Signer())

	// Add one more vote
	_, err = vs.AddVote(v4)
	assert.NoError(t, err)

	assert.NotNil(t, vs.QuorumHash())
	assert.Equal(t, vs.QuorumHash(), &blockHash)
	assert.True(t, vs.HasQuorumHash())
	assert.Contains(t, vs.BlockVotes(blockHash), v4.Signer())
}

func TestAllBlockVotes(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	valsMap, valKeys, totalPower := setupCommittee(ts, 1000, 1500, 2500, 2000)

	vs := NewPrecommitVoteSet(1, totalPower, valsMap)

	h1 := ts.RandHash()
	h2 := ts.RandHash()
	v1 := vote.NewPrecommitVote(h1, 1, 1, valKeys[0].Address())
	v2 := vote.NewPrecommitVote(h1, 1, 1, valKeys[1].Address())
	v3 := vote.NewPrecommitVote(h1, 1, 1, valKeys[2].Address())
	v4 := vote.NewPrecommitVote(h2, 1, 1, valKeys[0].Address())

	ts.HelperSignVote(valKeys[0], v1)
	ts.HelperSignVote(valKeys[1], v2)
	ts.HelperSignVote(valKeys[2], v3)
	ts.HelperSignVote(valKeys[3], v4)

	_, err := vs.AddVote(v1)
	assert.NoError(t, err)

	_, err = vs.AddVote(v2)
	assert.NoError(t, err)

	_, err = vs.AddVote(v3)
	assert.NoError(t, err)

	assert.Equal(t, vs.QuorumHash(), &h1)

	_, err = vs.AddVote(v4)
	assert.Error(t, err) // duplicated

	// Check accumulated power
	assert.Equal(t, vs.QuorumHash(), &h1)

	// Check previous votes
	assert.Contains(t, vs.AllVotes(), v1)
	assert.Contains(t, vs.AllVotes(), v2)
	assert.Contains(t, vs.AllVotes(), v3)
	assert.NotContains(t, vs.AllVotes(), v4) // Should add duplicated votes?
}

func TestAllBinaryVotes(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	valsMap, valKeys, totalPower := setupCommittee(ts, 1000, 1500, 2500, 2000)

	vs := NewCPMainVoteVoteSet(1, totalPower, valsMap)

	v1 := vote.NewCPMainVote(hash.UndefHash, 1, 1, 0, vote.CPValueZero, &vote.JustInitOne{}, valKeys[0].Address())
	v2 := vote.NewCPMainVote(hash.UndefHash, 1, 1, 1, vote.CPValueOne, &vote.JustInitOne{}, valKeys[1].Address())
	v3 := vote.NewCPMainVote(hash.UndefHash, 1, 1, 2, vote.CPValueAbstain, &vote.JustInitOne{}, valKeys[2].Address())

	ts.HelperSignVote(valKeys[0], v1)
	ts.HelperSignVote(valKeys[1], v2)
	ts.HelperSignVote(valKeys[2], v3)

	assert.Empty(t, vs.AllVotes())

	_, err := vs.AddVote(v1)
	assert.NoError(t, err)

	_, err = vs.AddVote(v2)
	assert.NoError(t, err)

	_, err = vs.AddVote(v3)
	assert.NoError(t, err)

	assert.Contains(t, vs.AllVotes(), v1)
	assert.Contains(t, vs.AllVotes(), v2)
	assert.Contains(t, vs.AllVotes(), v3)

	ranVote1 := vs.GetRandomVote(1, vote.CPValueZero)
	assert.Nil(t, ranVote1)

	ranVote2 := vs.GetRandomVote(1, vote.CPValueOne)
	assert.Equal(t, ranVote2, v2)
}

func TestOneThirdPower(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	// total power = 3000
	// 1/3 of total power = 1000
	// 2/3 of total power = 2000
	valsMap, valKeys, totalPower := setupCommittee(ts, 999, 3, 999, 999)

	h := ts.RandHash()
	height := ts.RandHeight()
	round := ts.RandRound()
	just := &vote.JustInitOne{}
	vs := NewCPPreVoteVoteSet(round, totalPower, valsMap)

	v1 := vote.NewCPPreVote(h, height, round, 0, vote.CPValueOne, just, valKeys[0].Address())
	v2 := vote.NewCPPreVote(h, height, round, 0, vote.CPValueOne, just, valKeys[1].Address())
	v3 := vote.NewCPPreVote(h, height, round, 0, vote.CPValueOne, just, valKeys[2].Address())
	v4 := vote.NewCPPreVote(h, height, round, 0, vote.CPValueZero, just, valKeys[3].Address())

	ts.HelperSignVote(valKeys[0], v1)
	ts.HelperSignVote(valKeys[1], v2)
	ts.HelperSignVote(valKeys[2], v3)
	ts.HelperSignVote(valKeys[3], v4)

	_, err := vs.AddVote(v1)
	assert.NoError(t, err)
	assert.False(t, vs.HasOneThirdOfTotalPower(0))
	assert.True(t, vs.HasAnyVoteFor(0, vote.CPValueOne))
	assert.False(t, vs.HasAnyVoteFor(0, vote.CPValueZero))
	assert.False(t, vs.HasAnyVoteFor(0, vote.CPValueAbstain))

	_, err = vs.AddVote(v2)
	assert.NoError(t, err)
	assert.True(t, vs.HasOneThirdOfTotalPower(0))
	assert.False(t, vs.HasTwoThirdOfTotalPower(0))

	_, err = vs.AddVote(v3)
	assert.NoError(t, err)
	assert.True(t, vs.HasTwoThirdOfTotalPower(0))
	assert.False(t, vs.HasAnyVoteFor(0, vote.CPValueZero))
	assert.True(t, vs.HasAnyVoteFor(0, vote.CPValueOne))
	assert.False(t, vs.HasQuorumVotesFor(0, vote.CPValueZero))
	assert.True(t, vs.HasQuorumVotesFor(0, vote.CPValueOne))
	assert.True(t, vs.HasAllVotesFor(0, vote.CPValueOne))

	_, err = vs.AddVote(v4)
	assert.NoError(t, err)
	assert.True(t, vs.HasAnyVoteFor(0, vote.CPValueZero))
	assert.False(t, vs.HasQuorumVotesFor(0, vote.CPValueZero))
	assert.True(t, vs.HasQuorumVotesFor(0, vote.CPValueOne))
	assert.False(t, vs.HasAllVotesFor(0, vote.CPValueOne))

	bv1 := vs.BinaryVotes(0, vote.CPValueOne)
	bv2 := vs.BinaryVotes(0, vote.CPValueZero)

	assert.Contains(t, bv1, v1.Signer())
	assert.Contains(t, bv1, v2.Signer())
	assert.Contains(t, bv1, v3.Signer())
	assert.Contains(t, bv2, v4.Signer())
}

func TestDecidedVoteset(t *testing.T) {
	ts := testsuite.NewTestSuite(t)
	valsMap, valKeys, totalPower := setupCommittee(ts, 1, 1, 1, 1)

	h := ts.RandHash()
	height := ts.RandHeight()
	round := ts.RandRound()
	just := &vote.JustInitOne{}
	vs := NewCPDecidedVoteVoteSet(round, totalPower, valsMap)

	v1 := vote.NewCPDecidedVote(h, height, round, 0, vote.CPValueOne, just, valKeys[0].Address())

	ts.HelperSignVote(valKeys[0], v1)

	_, err := vs.AddVote(v1)
	assert.NoError(t, err)
	assert.True(t, vs.HasAnyVoteFor(0, vote.CPValueOne))
	assert.False(t, vs.HasAnyVoteFor(0, vote.CPValueZero))
}
