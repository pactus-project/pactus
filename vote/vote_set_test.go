package vote

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/validator"
)

func TestAddVote(t *testing.T) {
	h1 := crypto.GenerateTestHash()
	addr, _, pv := crypto.GenerateTestKeyPair()
	valSet, keys := validator.GenerateTestValidatorSet()
	voteSet := NewVoteSet(100, 5, VoteTypePrecommit, valSet)

	v1 := NewVote(VoteTypePrecommit, 100, 5, h1, addr)
	v2 := NewVote(VoteTypePrecommit, 100, 5, h1, keys[0].PublicKey().Address())
	v3 := NewVote(VoteTypePrecommit, 101, 5, h1, keys[1].PublicKey().Address())
	v4 := NewVote(VoteTypePrecommit, 100, 6, h1, keys[2].PublicKey().Address())

	added, err := voteSet.AddVote(v1)
	assert.False(t, added) // not in val set
	assert.Error(t, err)
	assert.Nil(t, voteSet.ToCommit())

	sig := pv.Sign(v2.SignBytes())
	v2.SetSignature(sig)
	added, err = voteSet.AddVote(v2)
	assert.False(t, added) // invalid signature
	assert.Error(t, err)

	sig = keys[1].Sign(v2.SignBytes())
	v2.SetSignature(sig)
	added, err = voteSet.AddVote(v2)
	assert.False(t, added) // wrong signer
	assert.Error(t, err)

	sig = keys[0].Sign(v2.SignBytes())
	v2.SetSignature(sig)
	added, err = voteSet.AddVote(v2)
	assert.True(t, added) // ok
	assert.NoError(t, err)

	sig = keys[1].Sign(v2.SignBytes())
	v3.SetSignature(sig)
	added, err = voteSet.AddVote(v3)
	assert.False(t, added) // invalid height
	assert.Error(t, err)

	sig = keys[2].Sign(v2.SignBytes())
	v4.SetSignature(sig)
	added, err = voteSet.AddVote(v4)
	assert.False(t, added) // invalid round
	assert.Error(t, err)
}

func TestDuplicateVote(t *testing.T) {
	h1 := crypto.GenerateTestHash()
	h2 := crypto.GenerateTestHash()
	valSet, keys := validator.GenerateTestValidatorSet()
	voteSet := NewVoteSet(1, 0, VoteTypePrevote, valSet)

	undefVote := NewVote(VoteTypePrevote, 1, 0, crypto.UndefHash, keys[0].PublicKey().Address())
	correctVote := NewVote(VoteTypePrevote, 1, 0, h1, keys[0].PublicKey().Address())
	duplicatedVote := NewVote(VoteTypePrevote, 1, 0, h2, keys[0].PublicKey().Address())

	// sign the votes
	sig := keys[0].Sign(undefVote.SignBytes())
	undefVote.SetSignature(sig)

	sig = keys[0].Sign(undefVote.SignBytes())
	undefVote.SetSignature(sig)

	sig = keys[0].Sign(correctVote.SignBytes())
	correctVote.SetSignature(sig)

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

	sig = keys[0].Sign(duplicatedVote.SignBytes())
	duplicatedVote.SetSignature(sig)
	added, err = voteSet.AddVote(duplicatedVote)
	assert.False(t, added) // ok, replace UndefHash
	assert.Error(t, err)
	assert.Equal(t, err, errors.Error(errors.ErrDuplicateVote))
}
