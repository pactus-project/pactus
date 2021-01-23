package vote

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/validator"
)

func TestAddVote(t *testing.T) {
	h1 := crypto.GenerateTestHash()
	invSigner := crypto.GenerateTestSigner()
	valSet, signers := validator.GenerateTestValidatorSet()
	voteSet := NewVoteSet(100, 5, VoteTypePrecommit, valSet.CopyValidators())

	v1 := NewVote(VoteTypePrecommit, 100, 5, h1, invSigner.Address())
	v2 := NewVote(VoteTypePrecommit, 100, 5, h1, signers[0].Address())
	v3 := NewVote(VoteTypePrecommit, 101, 5, h1, signers[1].Address())
	v4 := NewVote(VoteTypePrecommit, 100, 6, h1, signers[2].Address())

	invSigner.SignMsg(v1)
	added, err := voteSet.AddVote(v1)
	assert.False(t, added) // not in val set
	assert.Error(t, err)
	assert.Nil(t, voteSet.ToCommit())

	invSigner.SignMsg(v2)
	added, err = voteSet.AddVote(v2)
	assert.False(t, added) // invalid signature
	assert.Error(t, err)

	signers[1].SignMsg(v2)
	added, err = voteSet.AddVote(v2)
	assert.False(t, added) // wrong signer
	assert.Error(t, err)

	signers[0].SignMsg(v2)
	added, err = voteSet.AddVote(v2)
	assert.True(t, added) // ok
	assert.NoError(t, err)

	signers[1].SignMsg(v3)
	added, err = voteSet.AddVote(v3)
	assert.False(t, added) // invalid height
	assert.Error(t, err)

	signers[2].SignMsg(v4)
	added, err = voteSet.AddVote(v4)
	assert.False(t, added) // invalid round
	assert.Error(t, err)
}

func TestDuplicateVote(t *testing.T) {
	h1 := crypto.GenerateTestHash()
	h2 := crypto.GenerateTestHash()
	valSet, signers := validator.GenerateTestValidatorSet()
	voteSet := NewVoteSet(1, 0, VoteTypePrepare, valSet.CopyValidators())

	undefVote := NewVote(VoteTypePrepare, 1, 0, crypto.UndefHash, signers[0].Address())
	correctVote := NewVote(VoteTypePrepare, 1, 0, h1, signers[0].Address())
	duplicatedVote := NewVote(VoteTypePrepare, 1, 0, h2, signers[0].Address())

	// sign the votes
	signers[0].SignMsg(undefVote)
	signers[0].SignMsg(correctVote)
	signers[0].SignMsg(duplicatedVote)

	added, err := voteSet.AddVote(undefVote)
	assert.True(t, added) // ok
	assert.NoError(t, err)

	added, err = voteSet.AddVote(undefVote)
	assert.False(t, added) // added before
	assert.NoError(t, err)

	added, err = voteSet.AddVote(correctVote)
	assert.True(t, added) // ok, replace UndefHash
	assert.NoError(t, err)
	assert.Equal(t, len(voteSet.AllVotes()), 1)

	// Again add undef vote
	added, err = voteSet.AddVote(undefVote)
	assert.False(t, added) // ok
	assert.NoError(t, err)
	assert.Equal(t, len(voteSet.AllVotes()), 1)

	added, err = voteSet.AddVote(duplicatedVote)
	assert.False(t, added) // ok, replace UndefHash
	assert.Error(t, err)
	assert.Equal(t, err, errors.Error(errors.ErrDuplicateVote))
}

func TestQuorum(t *testing.T) {
	valSet, signers := validator.GenerateTestValidatorSet()
	voteSet := NewVoteSet(1, 0, VoteTypePrecommit, valSet.CopyValidators())
	h1 := crypto.GenerateTestHash()
	h2 := crypto.GenerateTestHash()
	v1 := NewVote(VoteTypePrecommit, 1, 0, h1, signers[0].Address())
	v2 := NewVote(VoteTypePrecommit, 1, 0, h1, signers[1].Address())
	v3 := NewVote(VoteTypePrecommit, 1, 0, h1, signers[2].Address())
	v4 := NewVote(VoteTypePrecommit, 1, 0, h1, signers[3].Address())

	signers[0].SignMsg(v1)
	signers[1].SignMsg(v2)
	signers[2].SignMsg(v3)
	signers[3].SignMsg(v4)

	ok, _ := voteSet.AddVote(v1)
	assert.True(t, ok)
	assert.False(t, voteSet.HasQuorum())
	ok, _ = voteSet.AddVote(v2)
	assert.True(t, ok)
	assert.False(t, voteSet.HasQuorum())
	ok, _ = voteSet.AddVote(v3)
	assert.True(t, ok)
	assert.True(t, voteSet.HasQuorum())
	ok, _ = voteSet.AddVote(v4)
	assert.True(t, ok)
	assert.True(t, voteSet.HasQuorum())
	assert.True(t, voteSet.HasQuorumBlock(h1))
	assert.False(t, voteSet.HasQuorumBlock(h2))
	assert.NotNil(t, voteSet.QuorumBlock())
	assert.Equal(t, voteSet.QuorumBlock(), &h1)

	c := voteSet.ToCommit()
	assert.NotNil(t, c)
	assert.Equal(t, c.Committers(), []block.Committer{
		{Number: 0, Status: 1},
		{Number: 1, Status: 1},
		{Number: 2, Status: 1},
		{Number: 3, Status: 1},
	})
}

func TestUpdateVote(t *testing.T) {
	valSet, signers := validator.GenerateTestValidatorSet()
	voteSet := NewVoteSet(1, 0, VoteTypePrecommit, valSet.CopyValidators())

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

	ok, _ := voteSet.AddVote(v1)
	assert.True(t, ok)
	ok, _ = voteSet.AddVote(v2)
	assert.True(t, ok)
	ok, _ = voteSet.AddVote(v3)
	assert.True(t, ok)

	assert.True(t, voteSet.HasQuorum())
	assert.True(t, voteSet.HasQuorumBlock(crypto.UndefHash))
	assert.True(t, voteSet.QuorumBlock().EqualsTo(crypto.UndefHash))

	// Update vote
	ok, _ = voteSet.AddVote(v4)
	assert.True(t, ok)

	assert.True(t, voteSet.HasQuorum())
	assert.False(t, voteSet.HasQuorumBlock(crypto.UndefHash))
	assert.Nil(t, voteSet.QuorumBlock())
	assert.Equal(t, voteSet.sum, 3)

	ok, _ = voteSet.AddVote(v5)
	assert.True(t, ok)
	ok, _ = voteSet.AddVote(v6)
	assert.True(t, ok)
	assert.True(t, voteSet.QuorumBlock().EqualsTo(h1))
	assert.Equal(t, voteSet.sum, 3)

}
