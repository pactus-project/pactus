package vote

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/validator"
)

func setupValidatorSet(t *testing.T, stakes ...int64) (*validator.ValidatorSet, []crypto.Signer) {

	signers := []crypto.Signer{}
	vals := []*validator.Validator{}
	for i, s := range stakes {
		signer := crypto.GenerateTestSigner()
		val := validator.NewValidator(signer.PublicKey(), i, 0)
		val.AddToStake(s)
		vals = append(vals, val)
		signers = append(signers, signer)
	}
	valset, _ := validator.NewValidatorSet(vals, len(stakes), signers[0].Address())
	return valset, signers
}

func TestAddVote(t *testing.T) {
	valSet, signers := setupValidatorSet(t, 1000, 1500, 2500, 2000)

	h1 := crypto.GenerateTestHash()
	invSigner := crypto.GenerateTestSigner()
	vs := NewVoteSet(100, 5, VoteTypePrecommit, valSet.Validators())

	v1 := NewVote(VoteTypePrecommit, 100, 5, h1, invSigner.Address())
	v2 := NewVote(VoteTypePrecommit, 100, 5, h1, signers[0].Address())
	v3 := NewVote(VoteTypePrecommit, 101, 5, h1, signers[1].Address())
	v4 := NewVote(VoteTypePrecommit, 100, 6, h1, signers[2].Address())

	invSigner.SignMsg(v1)
	added, err := vs.AddVote(v1)
	assert.False(t, added) // not in val set
	assert.Error(t, err)
	assert.Nil(t, vs.ToCommit())

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
	valSet, signers := setupValidatorSet(t, 1000, 1500, 2500, 2000)

	h1 := crypto.GenerateTestHash()
	h2 := crypto.GenerateTestHash()
	vs := NewVoteSet(1, 0, VoteTypePrepare, valSet.Validators())

	undefVote := NewVote(VoteTypePrepare, 1, 0, crypto.UndefHash, signers[0].Address())
	correctVote := NewVote(VoteTypePrepare, 1, 0, h1, signers[0].Address())
	duplicatedVote := NewVote(VoteTypePrepare, 1, 0, h2, signers[0].Address())

	// sign the votes
	signers[0].SignMsg(undefVote)
	signers[0].SignMsg(correctVote)
	signers[0].SignMsg(duplicatedVote)

	added, err := vs.AddVote(undefVote)
	assert.True(t, added) // ok
	assert.NoError(t, err)

	added, err = vs.AddVote(undefVote)
	assert.False(t, added) // added before
	assert.NoError(t, err)
	assert.Equal(t, vs.Len(), 1)
	assert.Equal(t, vs.Power(), int64(1000)) // First validator's stake

	added, err = vs.AddVote(correctVote)
	assert.True(t, added) // ok, replace UndefHash
	assert.NoError(t, err)
	assert.Equal(t, vs.Len(), 1)             // The vote has replaced
	assert.Equal(t, vs.Power(), int64(1000)) // First validator's stake

	// Again add undef vote
	added, err = vs.AddVote(undefVote)
	assert.False(t, added) // ok
	assert.NoError(t, err)
	assert.Equal(t, vs.Len(), 1)

	bv := vs.blockVotes[h1]
	added, err = vs.AddVote(duplicatedVote)
	assert.False(t, added)
	assert.Error(t, err)
	assert.Equal(t, err, errors.Error(errors.ErrDuplicateVote))
	assert.Equal(t, bv.power, int64(0))   // vote has removed
	assert.Equal(t, vs.Len(), 0)          //
	assert.Equal(t, vs.Power(), int64(0)) //
}

func TestQuorum(t *testing.T) {
	valSet, signers := setupValidatorSet(t, 1000, 1500, 2500, 2000)

	vs := NewVoteSet(1, 0, VoteTypePrecommit, valSet.Validators())
	h1 := crypto.GenerateTestHash()
	v1 := NewVote(VoteTypePrecommit, 1, 0, h1, signers[0].Address())
	v2 := NewVote(VoteTypePrecommit, 1, 0, h1, signers[1].Address())
	v3 := NewVote(VoteTypePrecommit, 1, 0, h1, signers[2].Address())
	v4 := NewVote(VoteTypePrecommit, 1, 0, h1, signers[3].Address())

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
	ok, _ = vs.AddVote(v3)
	assert.True(t, ok)
	assert.True(t, vs.HasQuorum())
	ok, _ = vs.AddVote(v4)
	assert.True(t, ok)
	assert.True(t, vs.HasQuorum())
	assert.NotNil(t, vs.QuorumBlock())
	assert.Equal(t, vs.QuorumBlock(), &h1)
	assert.Equal(t, vs.Len(), 4)

	c := vs.ToCommit()
	assert.NotNil(t, c)
	assert.Equal(t, c.Committers(), []block.Committer{
		{Number: 0, Status: 1},
		{Number: 1, Status: 1},
		{Number: 2, Status: 1},
		{Number: 3, Status: 1},
	})
}

// This test is very important. Change it with cautious
func TestUpdateVote(t *testing.T) {
	valSet, signers := setupValidatorSet(t, 1000, 1500, 2500, 2000)

	vs := NewVoteSet(1, 0, VoteTypePrecommit, valSet.Validators())

	h1 := crypto.GenerateTestHash()
	v1 := NewVote(VoteTypePrecommit, 1, 0, crypto.UndefHash, signers[0].Address())
	v2 := NewVote(VoteTypePrecommit, 1, 0, crypto.UndefHash, signers[1].Address())
	v3 := NewVote(VoteTypePrecommit, 1, 0, crypto.UndefHash, signers[2].Address())
	v4 := NewVote(VoteTypePrecommit, 1, 0, h1, signers[0].Address())
	v5 := NewVote(VoteTypePrecommit, 1, 0, h1, signers[1].Address())
	v6 := NewVote(VoteTypePrecommit, 1, 0, h1, signers[2].Address())

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
	assert.True(t, vs.QuorumBlock().EqualsTo(crypto.UndefHash))
	assert.Equal(t, vs.Len(), 3)
	assert.Equal(t, vs.Power(), int64(1000+1500+2500))

	// Update vote
	ok, _ = vs.AddVote(v4)
	assert.True(t, ok)

	// Check block votes power
	bv1 := vs.blockVotes[crypto.UndefHash]
	assert.Equal(t, bv1.power, int64(1500+2500))
	bv2 := vs.blockVotes[h1]
	assert.Equal(t, bv2.power, int64(1000))

	// Check previous votes
	_, exists1 := bv1.votes[v1.Signer()]
	_, exists2 := bv1.votes[v2.Signer()]
	assert.False(t, exists1)
	assert.True(t, exists2)

	// Check accumulated power
	assert.True(t, vs.HasQuorum())
	assert.Nil(t, vs.QuorumBlock())
	assert.Equal(t, vs.Power(), int64(1000+1500+2500))
	assert.Equal(t, vs.Len(), 3)

	// Update more votes
	ok, _ = vs.AddVote(v5)
	assert.True(t, ok)
	ok, _ = vs.AddVote(v6)
	assert.True(t, ok)

	// Check block votes power
	bv1 = vs.blockVotes[crypto.UndefHash]
	assert.Equal(t, bv1.power, int64(0))
	bv2 = vs.blockVotes[h1]
	assert.Equal(t, bv2.power, int64(1000+1500+2500))

	assert.True(t, vs.HasQuorum())
	assert.Equal(t, vs.QuorumBlock(), &h1)
	assert.Equal(t, vs.Power(), int64(1000+1500+2500))
	assert.Equal(t, vs.Len(), 3)

	// Check previous votes
	_, exists1 = bv1.votes[v1.Signer()]
	_, exists2 = bv1.votes[v2.Signer()]
	assert.False(t, exists1)
	assert.False(t, exists2)
}

func TestAllVotes(t *testing.T) {
	valSet, signers := setupValidatorSet(t, 1000, 1500, 2500, 2000)

	vs := NewVoteSet(1, 0, VoteTypePrecommit, valSet.Validators())

	h1 := crypto.GenerateTestHash()
	v1 := NewVote(VoteTypePrecommit, 1, 0, crypto.UndefHash, signers[0].Address())
	v2 := NewVote(VoteTypePrecommit, 1, 0, h1, signers[1].Address())

	signers[0].SignMsg(v1)
	signers[1].SignMsg(v2)

	assert.Equal(t, vs.Len(), 0)
	assert.Empty(t, vs.AllVotes())

	ok, _ := vs.AddVote(v1)
	assert.True(t, ok)
	ok, _ = vs.AddVote(v2)
	assert.True(t, ok)

	assert.Equal(t, vs.Len(), 2)
	assert.Contains(t, vs.AllVotes(), v1)
	assert.Contains(t, vs.AllVotes(), v2)
}
