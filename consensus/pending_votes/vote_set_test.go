package pending_votes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/committee"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/validator"
	"github.com/zarbchain/zarb-go/vote"
)

func setupCommittee(t *testing.T, stakes ...int64) (*committee.Committee, []crypto.Signer) {

	signers := []crypto.Signer{}
	vals := []*validator.Validator{}
	for i, s := range stakes {
		signer := crypto.GenerateTestSigner()
		val := validator.NewValidator(signer.PublicKey(), i, 0)
		val.AddToStake(s)
		vals = append(vals, val)
		signers = append(signers, signer)
	}
	committee, _ := committee.NewCommittee(vals, len(stakes), signers[0].Address())
	return committee, signers
}

func TestAddVote(t *testing.T) {
	committee, signers := setupCommittee(t, 1000, 1500, 2500, 2000)

	h1 := crypto.GenerateTestHash()
	invSigner := crypto.GenerateTestSigner()
	vs := NewVoteSet(100, 5, vote.VoteTypePrecommit, committee.Validators())

	v1 := vote.NewVote(vote.VoteTypePrecommit, 100, 5, h1, invSigner.Address())
	v2 := vote.NewVote(vote.VoteTypePrecommit, 100, 5, h1, signers[0].Address())
	v3 := vote.NewVote(vote.VoteTypePrecommit, 101, 5, h1, signers[1].Address())
	v4 := vote.NewVote(vote.VoteTypePrecommit, 100, 6, h1, signers[2].Address())

	invSigner.SignMsg(v1)
	added, err := vs.AddVote(v1)
	assert.False(t, added) // not in committee
	assert.Error(t, err)
	assert.Nil(t, vs.ToCertificate())

	invSigner.SignMsg(v2)
	added, err = vs.AddVote(v2)
	assert.False(t, added) // invalid signature
	assert.Error(t, err)

	signers[1].SignMsg(v2)
	added, err = vs.AddVote(v2)
	assert.False(t, added) // wrong signer
	assert.Error(t, err)

	signers[0].SignMsg(v2)
	added, err = vs.AddVote(v2)
	assert.True(t, added) // ok
	assert.NoError(t, err)

	signers[1].SignMsg(v3)
	added, err = vs.AddVote(v3)
	assert.False(t, added) // invalid height
	assert.Error(t, err)

	signers[2].SignMsg(v4)
	added, err = vs.AddVote(v4)
	assert.False(t, added) // invalid round
	assert.Error(t, err)
}

func TestDuplicateVote(t *testing.T) {
	committee, signers := setupCommittee(t, 1000, 1500, 2500, 2000)

	h1 := crypto.GenerateTestHash()
	h2 := crypto.GenerateTestHash()
	h3 := crypto.GenerateTestHash()
	vs := NewVoteSet(1, 0, vote.VoteTypePrepare, committee.Validators())

	nullVote := vote.NewVote(vote.VoteTypePrepare, 1, 0, crypto.UndefHash, signers[0].Address())
	correctVote := vote.NewVote(vote.VoteTypePrepare, 1, 0, h1, signers[0].Address())
	duplicatedVote1 := vote.NewVote(vote.VoteTypePrepare, 1, 0, h2, signers[0].Address())
	duplicatedVote2 := vote.NewVote(vote.VoteTypePrepare, 1, 0, h3, signers[0].Address())

	// sign the votes
	signers[0].SignMsg(nullVote)
	signers[0].SignMsg(correctVote)
	signers[0].SignMsg(duplicatedVote1)
	signers[0].SignMsg(duplicatedVote2)

	added, err := vs.AddVote(nullVote)
	assert.True(t, added)                               // ok
	assert.NoError(t, err)                              //
	assert.Equal(t, vs.Len(), 1)                        //
	assert.Equal(t, vs.AccumulatedPower(), int64(1000)) // First validator's stake

	added, err = vs.AddVote(correctVote)
	assert.True(t, added)                               // ok
	assert.NoError(t, err)                              //
	assert.Equal(t, vs.Len(), 2)                        // null + block_vote
	assert.Equal(t, vs.AccumulatedPower(), int64(1000)) //

	added, err = vs.AddVote(duplicatedVote1)
	assert.True(t, added)                               // ok
	assert.Error(t, err)                                //
	assert.Equal(t, vs.Len(), 3)                        // null + block_vote + duplicated1
	assert.Equal(t, vs.AccumulatedPower(), int64(1000)) //
	assert.Equal(t, err, errors.Error(errors.ErrDuplicateVote))

	added, err = vs.AddVote(duplicatedVote2)
	assert.True(t, added)                               // ok
	assert.Error(t, err)                                //
	assert.Equal(t, vs.Len(), 4)                        // null + block_vote + duplicated1 + duplicated2
	assert.Equal(t, vs.AccumulatedPower(), int64(1000)) //
	assert.Equal(t, err, errors.Error(errors.ErrDuplicateVote))

	added, err = vs.AddVote(nullVote)
	assert.False(t, added)                              // added before
	assert.NoError(t, err)                              //
	assert.Equal(t, vs.Len(), 4)                        //
	assert.Equal(t, vs.AccumulatedPower(), int64(1000)) //

	// Again add null_vote
	added, err = vs.AddVote(nullVote)
	assert.False(t, added)                              // added before
	assert.NoError(t, err)                              //
	assert.Equal(t, vs.Len(), 4)                        //
	assert.Equal(t, vs.AccumulatedPower(), int64(1000)) //

	bv1 := vs.blockVotes[h1]
	bv2 := vs.blockVotes[h2]
	bv3 := vs.blockVotes[h3]
	bv4 := vs.blockVotes[crypto.UndefHash]
	assert.Equal(t, vs.Len(), 4)                        // vote + duplicated1 + duplicated2
	assert.Equal(t, vs.AccumulatedPower(), int64(1000)) //
	assert.Equal(t, bv1.power, int64(1000))             //
	assert.Equal(t, bv2.power, int64(1000))             //
	assert.Equal(t, bv3.power, int64(1000))             //
	assert.Equal(t, bv4.power, int64(0))                //

	// Add vote again
}

func TestQuorum(t *testing.T) {
	committee, signers := setupCommittee(t, 1000, 1500, 2500, 2000)

	vs := NewVoteSet(1, 0, vote.VoteTypePrecommit, committee.Validators())
	h1 := crypto.GenerateTestHash()
	v1 := vote.NewVote(vote.VoteTypePrecommit, 1, 0, h1, signers[0].Address())
	v2 := vote.NewVote(vote.VoteTypePrecommit, 1, 0, h1, signers[1].Address())
	v3 := vote.NewVote(vote.VoteTypePrecommit, 1, 0, h1, signers[2].Address())
	v4 := vote.NewVote(vote.VoteTypePrecommit, 1, 0, h1, signers[3].Address())

	signers[0].SignMsg(v1)
	signers[1].SignMsg(v2)
	signers[2].SignMsg(v3)
	signers[3].SignMsg(v4)

	ok, _ := vs.AddVote(v1)
	assert.True(t, ok)
	assert.False(t, vs.HasQuorum())
	ok, _ = vs.AddVote(v2)
	assert.True(t, ok)
	assert.False(t, vs.HasQuorum())
	assert.Nil(t, vs.ToCertificate())

	ok, _ = vs.AddVote(v3)
	assert.True(t, ok)
	assert.True(t, vs.HasQuorum())

	cert1 := vs.ToCertificate()
	assert.NotNil(t, cert1)
	assert.Equal(t, cert1.Committers(), []int{0, 1, 2, 3})
	assert.Equal(t, cert1.Absences(), []int{3})

	// Add one more vote
	ok, _ = vs.AddVote(v4)
	assert.True(t, ok)
	assert.True(t, vs.HasQuorum())
	assert.NotNil(t, vs.QuorumBlock())
	assert.Equal(t, vs.QuorumBlock(), &h1)
	assert.Equal(t, vs.Len(), 4)

	cert2 := vs.ToCertificate()
	assert.NotNil(t, cert2)
	assert.Equal(t, cert2.Committers(), []int{0, 1, 2, 3})
	assert.Equal(t, cert2.Absences(), []int{})
}

// This test is very important. Change it with cautious
func TestPower(t *testing.T) {
	committee, signers := setupCommittee(t, 1000, 1500, 2500, 2000)

	vs := NewVoteSet(1, 0, vote.VoteTypePrecommit, committee.Validators())

	h1 := crypto.GenerateTestHash()
	v1 := vote.NewVote(vote.VoteTypePrecommit, 1, 0, crypto.UndefHash, signers[0].Address())
	v2 := vote.NewVote(vote.VoteTypePrecommit, 1, 0, crypto.UndefHash, signers[1].Address())
	v3 := vote.NewVote(vote.VoteTypePrecommit, 1, 0, crypto.UndefHash, signers[2].Address())
	v4 := vote.NewVote(vote.VoteTypePrecommit, 1, 0, h1, signers[0].Address())
	v5 := vote.NewVote(vote.VoteTypePrecommit, 1, 0, h1, signers[1].Address())
	v6 := vote.NewVote(vote.VoteTypePrecommit, 1, 0, h1, signers[2].Address())

	signers[0].SignMsg(v1)
	signers[1].SignMsg(v2)
	signers[2].SignMsg(v3)
	signers[0].SignMsg(v4)
	signers[1].SignMsg(v5)
	signers[2].SignMsg(v6)

	ok, _ := vs.AddVote(v1)
	assert.True(t, ok)
	ok, _ = vs.AddVote(v2)
	assert.True(t, ok)
	ok, _ = vs.AddVote(v3)
	assert.True(t, ok)

	assert.True(t, vs.HasQuorum())
	assert.True(t, vs.QuorumBlock().IsUndef())
	assert.Equal(t, vs.Len(), 3)
	assert.Equal(t, vs.AccumulatedPower(), int64(1000+1500+2500))

	ok, _ = vs.AddVote(v4)
	assert.True(t, ok)

	// Check block votes power
	bv1 := vs.blockVotes[crypto.UndefHash]
	bv2 := vs.blockVotes[h1]
	assert.Equal(t, bv1.power, int64(1500+2500))
	assert.Equal(t, bv2.power, int64(1000))

	// Check previous votes
	assert.Contains(t, vs.AllVotes(), v1)
	assert.Contains(t, vs.AllVotes(), v2)

	// Check accumulated power
	assert.True(t, vs.HasQuorum())
	assert.Nil(t, vs.QuorumBlock())
	assert.Equal(t, vs.AccumulatedPower(), int64(1000+1500+2500))
	assert.Equal(t, vs.Len(), 4)

	// Add more votes
	ok, _ = vs.AddVote(v5)
	assert.True(t, ok)
	ok, _ = vs.AddVote(v6)
	assert.True(t, ok)

	// Check block votes power
	bv1 = vs.blockVotes[crypto.UndefHash]
	bv2 = vs.blockVotes[h1]
	assert.Equal(t, bv1.power, int64(0))
	assert.Equal(t, bv2.power, int64(1000+1500+2500))

	assert.True(t, vs.HasQuorum())
	assert.Equal(t, vs.QuorumBlock(), &h1)
	assert.Equal(t, vs.AccumulatedPower(), int64(1000+1500+2500))
	assert.Equal(t, vs.Len(), 6)

	// Check previous votes
	assert.Contains(t, vs.AllVotes(), v1)
	assert.Contains(t, vs.AllVotes(), v2)
	assert.Contains(t, vs.AllVotes(), v3)
	assert.Contains(t, vs.AllVotes(), v4)
	assert.Contains(t, vs.AllVotes(), v5)
	assert.Contains(t, vs.AllVotes(), v6)
}

func TestAllVotes(t *testing.T) {
	committee, signers := setupCommittee(t, 1000, 1500, 2500, 2000)

	vs := NewVoteSet(1, 0, vote.VoteTypePrecommit, committee.Validators())

	v1 := vote.NewVote(vote.VoteTypePrecommit, 1, 0, crypto.GenerateTestHash(), signers[0].Address())
	v2 := vote.NewVote(vote.VoteTypePrecommit, 1, 0, crypto.UndefHash, signers[0].Address())
	v3 := vote.NewVote(vote.VoteTypePrecommit, 1, 0, crypto.GenerateTestHash(), signers[1].Address())

	signers[0].SignMsg(v1)
	signers[0].SignMsg(v2)
	signers[1].SignMsg(v3)

	assert.Equal(t, vs.Len(), 0)
	assert.Empty(t, vs.AllVotes())

	ok, _ := vs.AddVote(v1)
	assert.True(t, ok)
	ok, _ = vs.AddVote(v2)
	assert.False(t, ok) // Ignore null after block vote
	ok, _ = vs.AddVote(v3)
	assert.True(t, ok)

	assert.Equal(t, vs.Len(), 3)
	assert.Contains(t, vs.AllVotes(), v1)
	assert.Contains(t, vs.AllVotes(), v2)
	assert.Contains(t, vs.AllVotes(), v3)
}

func TestUpdateQuoromBlock(t *testing.T) {
	committee, signers := setupCommittee(t, 1000, 1000, 2500, 2000)

	vs := NewVoteSet(1, 0, vote.VoteTypePrepare, committee.Validators())

	h1 := crypto.GenerateTestHash()
	v1 := vote.NewVote(vote.VoteTypePrepare, 1, 0, crypto.UndefHash, signers[3].Address())
	v2 := vote.NewVote(vote.VoteTypePrepare, 1, 0, crypto.UndefHash, signers[2].Address())
	v3 := vote.NewVote(vote.VoteTypePrepare, 1, 0, h1, signers[3].Address())
	v4 := vote.NewVote(vote.VoteTypePrepare, 1, 0, h1, signers[2].Address())
	v5 := vote.NewVote(vote.VoteTypePrepare, 1, 0, crypto.UndefHash, signers[1].Address())

	signers[3].SignMsg(v1)
	signers[2].SignMsg(v2)
	signers[3].SignMsg(v3)
	signers[2].SignMsg(v4)
	signers[1].SignMsg(v5)

	ok, _ := vs.AddVote(v1)
	assert.True(t, ok)
	assert.False(t, vs.HasQuorum())
	assert.Nil(t, vs.QuorumBlock())

	ok, _ = vs.AddVote(v2)
	assert.True(t, ok)
	assert.True(t, vs.HasQuorum())
	assert.Equal(t, vs.QuorumBlock(), &crypto.UndefHash)

	ok, _ = vs.AddVote(v3)
	assert.True(t, ok)
	assert.True(t, vs.HasQuorum())
	assert.Nil(t, vs.QuorumBlock())

	ok, _ = vs.AddVote(v4)
	assert.True(t, ok)
	assert.True(t, vs.HasQuorum())
	assert.Equal(t, vs.QuorumBlock(), &h1)

	ok, _ = vs.AddVote(v5)
	assert.True(t, ok)
	assert.True(t, vs.HasQuorum())
	assert.Equal(t, vs.QuorumBlock(), &h1)
}

func TestHasOneThirdOfTotalPower(t *testing.T) {
	committee, signers := setupCommittee(t, 1000, 1500, 2500, 2000)

	vs := NewVoteSet(1, 0, vote.VoteTypePrepare, committee.Validators())

	v1 := vote.NewVote(vote.VoteTypePrepare, 1, 0, crypto.UndefHash, signers[0].Address())
	v2 := vote.NewVote(vote.VoteTypePrepare, 1, 0, crypto.UndefHash, signers[1].Address())

	signers[0].SignMsg(v1)
	signers[1].SignMsg(v2)

	assert.False(t, vs.HasOneThirdOfTotalPower(crypto.UndefHash))

	ok, _ := vs.AddVote(v1)
	assert.True(t, ok)
	assert.False(t, vs.HasOneThirdOfTotalPower(crypto.UndefHash))

	ok, _ = vs.AddVote(v2)
	assert.True(t, ok)
	assert.True(t, vs.HasOneThirdOfTotalPower(crypto.UndefHash))
}
