package voteset

import (
	"testing"

	"github.com/pactus-project/pactus/committee"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/pactus-project/pactus/util/errors"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func setupCommittee(t *testing.T, ts *testsuite.TestSuite, stakes ...int64) (committee.Committee, []crypto.Signer) {
	t.Helper()

	signers := []crypto.Signer{}
	vals := []*validator.Validator{}
	for i, s := range stakes {
		pub, pv := ts.RandomBLSKeyPair()
		val := validator.NewValidator(pub, int32(i))
		val.AddToStake(s)
		vals = append(vals, val)
		signers = append(signers, crypto.NewSigner(pv))
	}
	committee, err := committee.NewCommittee(vals, len(stakes), signers[0].Address())
	assert.NoError(t, err)
	return committee, signers
}

func TestAddVote(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	committee, signers := setupCommittee(t, ts, 1000, 1500, 2500, 2000)

	h1 := ts.RandomHash()
	invSigner := ts.RandomSigner()
	vs := NewVoteSet(5, vote.VoteTypePrecommit, committee.Validators())

	v1 := vote.NewVote(vote.VoteTypePrecommit, 100, 5, h1, invSigner.Address())
	v2 := vote.NewVote(vote.VoteTypePrecommit, 100, 5, h1, signers[0].Address())
	v3 := vote.NewVote(vote.VoteTypePrecommit, 100, 6, h1, signers[2].Address())

	invSigner.SignMsg(v1)
	err := vs.AddVote(v1)
	assert.Error(t, err) // not in committee
	assert.Nil(t, vs.ToCertificate())

	invSigner.SignMsg(v2)
	err = vs.AddVote(v2)
	assert.Error(t, err) // invalid signature

	signers[1].SignMsg(v2)
	err = vs.AddVote(v2)
	assert.Error(t, err) // wrong signer

	signers[0].SignMsg(v2)
	err = vs.AddVote(v2)
	assert.NoError(t, err) // ok

	signers[2].SignMsg(v3)
	err = vs.AddVote(v3)
	assert.Error(t, err) // invalid round
}

func TestDuplicateVote(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	committee, signers := setupCommittee(t, ts, 1000, 1500, 2500, 2000)

	h1 := ts.RandomHash()
	h2 := ts.RandomHash()
	h3 := ts.RandomHash()
	vs := NewVoteSet(0, vote.VoteTypePrepare, committee.Validators())

	correctVote := vote.NewVote(vote.VoteTypePrepare, 1, 0, h1, signers[0].Address())
	duplicatedVote1 := vote.NewVote(vote.VoteTypePrepare, 1, 0, h2, signers[0].Address())
	duplicatedVote2 := vote.NewVote(vote.VoteTypePrepare, 1, 0, h3, signers[0].Address())

	// sign the votes
	signers[0].SignMsg(correctVote)
	signers[0].SignMsg(duplicatedVote1)
	signers[0].SignMsg(duplicatedVote2)

	assert.NoError(t, vs.AddVote(correctVote)) // ok
	assert.Equal(t, vs.Len(), 1)               // correctVote

	err := vs.AddVote(duplicatedVote1)
	assert.Equal(t, errors.Code(err), errors.ErrDuplicateVote) // duplicated vote
	assert.Equal(t, vs.Len(), 2)                               // correctVote + duplicatedVote1

	err = vs.AddVote(duplicatedVote2)
	assert.Equal(t, errors.Code(err), errors.ErrDuplicateVote) // duplicated vote
	assert.Equal(t, vs.Len(), 3)                               // correctVote + duplicatedVote1 + duplicatedVote2

	assert.Error(t, vs.AddVote(correctVote)) // added before
	assert.Equal(t, vs.Len(), 3)             // correctVote + duplicatedVote1 + duplicatedVote2

	bv1 := vs.blockVotes[h1]
	bv2 := vs.blockVotes[h2]
	bv3 := vs.blockVotes[h3]
	assert.Equal(t, vs.Len(), 3)            // correctVote + duplicatedVote1 + duplicatedVote2
	assert.Equal(t, bv1.power, int64(1000)) //
	assert.Nil(t, bv2)                      //
	assert.Nil(t, bv3)                      //
}

func TestQuorum(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	committee, signers := setupCommittee(t, ts, 1000, 1500, 2500, 2000)

	vs := NewVoteSet(0, vote.VoteTypePrecommit, committee.Validators())
	h1 := ts.RandomHash()
	v1 := vote.NewVote(vote.VoteTypePrecommit, 1, 0, h1, signers[0].Address())
	v2 := vote.NewVote(vote.VoteTypePrecommit, 1, 0, h1, signers[1].Address())
	v3 := vote.NewVote(vote.VoteTypePrecommit, 1, 0, h1, signers[2].Address())
	v4 := vote.NewVote(vote.VoteTypePrecommit, 1, 0, h1, signers[3].Address())

	signers[0].SignMsg(v1)
	signers[1].SignMsg(v2)
	signers[2].SignMsg(v3)
	signers[3].SignMsg(v4)

	assert.NoError(t, vs.AddVote(v1))
	assert.NoError(t, vs.AddVote(v2))
	assert.Nil(t, vs.ToCertificate())
	assert.Nil(t, vs.QuorumHash())

	assert.NoError(t, vs.AddVote(v3))
	cert1 := vs.ToCertificate()
	assert.NotNil(t, cert1)
	assert.Equal(t, cert1.Committers(), []int32{0, 1, 2, 3})
	assert.Equal(t, cert1.Absentees(), []int32{3})

	// Add one more vote
	assert.NoError(t, vs.AddVote(v4))
	assert.NotNil(t, vs.QuorumHash())
	assert.Equal(t, vs.QuorumHash(), &h1)
	assert.Equal(t, vs.Len(), 4)

	cert2 := vs.ToCertificate()
	assert.NotNil(t, cert2)
	assert.Equal(t, cert2.Committers(), []int32{0, 1, 2, 3})
	assert.Equal(t, cert2.Absentees(), []int32{})
}

func TestPower(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	committee, signers := setupCommittee(t, ts, 1000, 1500, 2500, 2000)

	vs := NewVoteSet(0, vote.VoteTypePrecommit, committee.Validators())

	h1 := ts.RandomHash()
	h2 := ts.RandomHash()
	v1 := vote.NewVote(vote.VoteTypePrecommit, 1, 0, h1, signers[0].Address())
	v2 := vote.NewVote(vote.VoteTypePrecommit, 1, 0, h1, signers[1].Address())
	v3 := vote.NewVote(vote.VoteTypePrecommit, 1, 0, h1, signers[2].Address())
	v4 := vote.NewVote(vote.VoteTypePrecommit, 1, 0, h2, signers[0].Address())

	signers[0].SignMsg(v1)
	signers[1].SignMsg(v2)
	signers[2].SignMsg(v3)
	signers[0].SignMsg(v4)

	assert.NoError(t, vs.AddVote(v1))
	assert.NoError(t, vs.AddVote(v2))
	assert.NoError(t, vs.AddVote(v3))

	assert.True(t, vs.QuorumHash().EqualsTo(h1))
	assert.Equal(t, vs.Len(), 3)

	assert.Error(t, vs.AddVote(v4)) // duplicated

	// Check accumulated power
	assert.True(t, vs.QuorumHash().EqualsTo(h1))
	assert.Equal(t, vs.Len(), 4)

	// Check previous votes
	assert.Contains(t, vs.AllVotes(), v1)
	assert.Contains(t, vs.AllVotes(), v2)
	assert.Contains(t, vs.AllVotes(), v3)
	assert.Contains(t, vs.AllVotes(), v4)
}

func TestAllVotes(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	committee, signers := setupCommittee(t, ts, 1000, 1500, 2500, 2000)

	vs := NewVoteSet(0, vote.VoteTypeChangeProposer, committee.Validators())

	v1 := vote.NewVote(vote.VoteTypeChangeProposer, 1, 0, hash.UndefHash, signers[0].Address())
	v2 := vote.NewVote(vote.VoteTypeChangeProposer, 1, 0, hash.UndefHash, signers[1].Address())
	v3 := vote.NewVote(vote.VoteTypeChangeProposer, 1, 0, hash.UndefHash, signers[2].Address())

	signers[0].SignMsg(v1)
	signers[1].SignMsg(v2)
	signers[2].SignMsg(v3)

	assert.Equal(t, vs.Len(), 0)
	assert.Empty(t, vs.AllVotes())

	assert.NoError(t, vs.AddVote(v1))
	assert.NoError(t, vs.AddVote(v2))
	assert.NoError(t, vs.AddVote(v3))

	assert.Equal(t, vs.Len(), 3)
	assert.Contains(t, vs.AllVotes(), v1)
	assert.Contains(t, vs.AllVotes(), v2)
	assert.Contains(t, vs.AllVotes(), v3)
	assert.NotNil(t, vs.QuorumHash())
}

func TestOneThirdPower(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	committee, signers := setupCommittee(t, ts, 1000, 1000, 1500, 1500)

	vs := NewVoteSet(0, vote.VoteTypeChangeProposer, committee.Validators())

	v1 := vote.NewVote(vote.VoteTypeChangeProposer, 1, 0, hash.UndefHash, signers[0].Address())
	v2 := vote.NewVote(vote.VoteTypeChangeProposer, 1, 0, hash.UndefHash, signers[1].Address())

	signers[0].SignMsg(v1)
	signers[1].SignMsg(v2)

	assert.NoError(t, vs.AddVote(v1))
	assert.False(t, vs.BlockHashHasOneThirdOfTotalPower(hash.UndefHash))

	assert.NoError(t, vs.AddVote(v2))
	assert.True(t, vs.BlockHashHasOneThirdOfTotalPower(hash.UndefHash))
}
