package consensus

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/validator"
	"github.com/zarbchain/zarb-go/vote"
)

func TestMustGetRound(t *testing.T) {
	valSet, _ := validator.GenerateTestValidatorSet()
	pv := NewPendingVotes()
	pv.MoveToNewHeight(101, valSet.CopyValidators())
	pv.MustGetRoundVotes(4)
	assert.Nil(t, pv.GetRoundVotes(5))
	assert.NotNil(t, pv.GetRoundVotes(1))
	assert.NotNil(t, pv.GetRoundVotes(4))
	assert.Equal(t, pv.roundVotes[3].Prepares.Height(), 101)
	assert.Equal(t, pv.roundVotes[3].Prepares.Round(), 3)
	assert.Equal(t, len(pv.roundVotes), 5)
}

func TestPendingVotesTest(t *testing.T) {
	valSet, keys := validator.GenerateTestValidatorSet()

	pv := NewPendingVotes()
	pv.MoveToNewHeight(101, valSet.CopyValidators())
	invalidVote, _ := vote.GenerateTestPrecommitVote(55, 5)
	ok, err := pv.AddVote(invalidVote) // invalid height
	assert.False(t, ok)
	assert.Error(t, err)

	v1, _ := vote.GenerateTestPrecommitVote(101, 5)
	ok, err = pv.AddVote(v1) // invalid signer
	assert.False(t, ok)
	assert.Error(t, err)

	undefVote := vote.NewVote(vote.VoteTypePrepare, 101, 1, crypto.UndefHash, keys[0].PublicKey().Address())
	undefVote.SetSignature(keys[0].Sign(undefVote.SignBytes()))

	validVote := vote.NewVote(vote.VoteTypePrepare, 101, 1, crypto.GenerateTestHash(), keys[0].PublicKey().Address())
	validVote.SetSignature(keys[0].Sign(validVote.SignBytes()))

	duplicateVote := vote.NewVote(vote.VoteTypePrepare, 101, 1, crypto.GenerateTestHash(), keys[0].PublicKey().Address())
	duplicateVote.SetSignature(keys[0].Sign(duplicateVote.SignBytes()))

	ok, err = pv.AddVote(undefVote)
	assert.True(t, ok)
	assert.NoError(t, err)

	ok, err = pv.AddVote(validVote)
	assert.True(t, ok)
	assert.NoError(t, err)

	// Because of network lagging we might receive nil-vote after block-vote
	// We don't add this vote and we don't report it as duplicated
	ok, err = pv.AddVote(undefVote)
	assert.False(t, ok)
	assert.NoError(t, err)

	// Definitely it is a duplicated error
	ok, err = pv.AddVote(duplicateVote)
	assert.False(t, ok) // duplicated vote
	assert.Error(t, err)

	prepares := pv.PrepareVoteSet(1)
	assert.Equal(t, prepares.Len(), 1)              // validVote
	assert.Equal(t, len(pv.roundVotes[1].votes), 2) // validVote + duplicateVote
	assert.True(t, pv.HasVote(validVote.Hash()))
	assert.False(t, pv.HasVote(invalidVote.Hash()))
}

func TestSetProposal(t *testing.T) {
	valSet, _ := validator.GenerateTestValidatorSet()
	prop, _ := vote.GenerateTestProposal(101, 0)
	pv := NewPendingVotes()
	pv.MoveToNewHeight(101, valSet.CopyValidators())
	pv.SetRoundProposal(4, prop)
	assert.False(t, pv.HasRoundProposal(0))
	assert.True(t, pv.HasRoundProposal(4))
	assert.Equal(t, pv.RoundProposal(4).Hash(), prop.Hash())
}
