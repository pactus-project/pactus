package consensus

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/vote"

	"github.com/zarbchain/zarb-go/validator"
)

func TestHeightVoteSetTest(t *testing.T) {
	vset, keys := validator.GenerateTestValidatorSet()

	hvs := NewHeightVoteSet(101, vset)
	invalidVote, _ := vote.GenerateTestPrecommitVote(55, 5)
	ok, err := hvs.AddVote(invalidVote) // invalid height
	assert.False(t, ok)
	assert.Error(t, err)

	v1, _ := vote.GenerateTestPrecommitVote(101, 5)
	ok, err = hvs.AddVote(v1) // invalid signer
	assert.False(t, ok)
	assert.Error(t, err)

	undefVote := vote.NewVote(vote.VoteTypePrevote, 101, 1, crypto.UndefHash, keys[0].PublicKey().Address())
	undefVote.SetSignature(keys[0].Sign(undefVote.SignBytes()))

	validVote := vote.NewVote(vote.VoteTypePrevote, 101, 1, crypto.GenerateTestHash(), keys[0].PublicKey().Address())
	validVote.SetSignature(keys[0].Sign(validVote.SignBytes()))

	duplicateVote := vote.NewVote(vote.VoteTypePrevote, 101, 1, crypto.GenerateTestHash(), keys[0].PublicKey().Address())
	duplicateVote.SetSignature(keys[0].Sign(duplicateVote.SignBytes()))

	ok, err = hvs.AddVote(undefVote)
	assert.True(t, ok)
	assert.NoError(t, err)

	ok, err = hvs.AddVote(validVote)
	assert.True(t, ok)
	assert.NoError(t, err)

	ok, err = hvs.AddVote(undefVote)
	assert.False(t, ok)
	assert.NoError(t, err)

	ok, err = hvs.AddVote(duplicateVote)
	assert.False(t, ok) // duplicated vote
	assert.Error(t, err)

	prevotes := hvs.Prevotes(1)
	assert.Equal(t, prevotes.Len(), 1) // validVote
	assert.Equal(t, len(hvs.votes), 2) // validVote + duplicateVote

}
