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
	invalid_vote, _ := vote.GenerateTestPrecommitVote(55, 5)
	ok, err := hvs.AddVote(invalid_vote) // invalid height
	assert.False(t, ok)
	assert.Error(t, err)

	v1, _ := vote.GenerateTestPrecommitVote(101, 5)
	ok, err = hvs.AddVote(v1) // invalid signer
	assert.False(t, ok)
	assert.Error(t, err)

	v2 := vote.NewVote(vote.VoteTypePrevote, 101, 1, crypto.UndefHash, keys[0].PublicKey().Address())
	v2.SetSignature(keys[0].Sign(v2.SignBytes()))
	ok, err = hvs.AddVote(v2)
	assert.True(t, ok)
	assert.NoError(t, err)

	v3 := vote.NewVote(vote.VoteTypePrevote, 101, 1, crypto.GenerateTestHash(), keys[0].PublicKey().Address())
	v3.SetSignature(keys[0].Sign(v3.SignBytes()))
	ok, err = hvs.AddVote(v3)
	assert.True(t, ok)
	assert.NoError(t, err)

	v4 := vote.NewVote(vote.VoteTypePrevote, 101, 1, crypto.GenerateTestHash(), keys[0].PublicKey().Address())
	v4.SetSignature(keys[0].Sign(v4.SignBytes()))
	ok, err = hvs.AddVote(v4)
	assert.False(t, ok) // duplicated vote
	assert.Error(t, err)

	prevotes := hvs.Prevotes(1)
	assert.Equal(t, prevotes.Len(), 1) // only v3
	assert.Equal(t, len(hvs.votes), 2) // v2 + v3

}
