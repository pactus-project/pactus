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
	pv.MoveToNewHeight(101, valSet.Validators())
	pv.MustGetRoundVotes(4)
	assert.Nil(t, pv.GetRoundVotes(5))
	assert.NotNil(t, pv.GetRoundVotes(1))
	assert.NotNil(t, pv.GetRoundVotes(4))
	assert.Equal(t, pv.GetRoundVotes(3).Prepares.Height(), 101)
	assert.Equal(t, pv.GetRoundVotes(3).Prepares.Round(), 3)
	assert.Equal(t, len(pv.roundVotes), 5)
}

func TestPendingVotesTest(t *testing.T) {
	valSet, signers := validator.GenerateTestValidatorSet()

	pv := NewPendingVotes()
	pv.MoveToNewHeight(101, valSet.Validators())
	invalidVote, _ := vote.GenerateTestPrecommitVote(55, 5)
	ok, err := pv.AddVote(invalidVote) // invalid height
	assert.False(t, ok)
	assert.Error(t, err)

	v1, _ := vote.GenerateTestPrecommitVote(101, 5)
	ok, err = pv.AddVote(v1) // invalid signer
	assert.False(t, ok)
	assert.Error(t, err)

	undefVote := vote.NewVote(vote.VoteTypePrepare, 101, 1, crypto.UndefHash, signers[0].Address())
	signers[0].SignMsg(undefVote)

	validVote := vote.NewVote(vote.VoteTypePrepare, 101, 1, crypto.GenerateTestHash(), signers[0].Address())
	signers[0].SignMsg(validVote)

	duplicateVote := vote.NewVote(vote.VoteTypePrepare, 101, 1, crypto.GenerateTestHash(), signers[0].Address())
	signers[0].SignMsg(duplicateVote)

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
	assert.Equal(t, prepares.Len(), 0) // duplicated vote has removed
	assert.Equal(t, len(pv.GetRoundVotes(1).AllVotes()), 0)
	assert.False(t, pv.HasVote(duplicateVote.Hash()))
	assert.False(t, pv.HasVote(validVote.Hash()))
	assert.False(t, pv.HasVote(invalidVote.Hash()))
}

func TestSetRoundProposal(t *testing.T) {
	valSet, _ := validator.GenerateTestValidatorSet()
	prop, _ := vote.GenerateTestProposal(101, 0)
	pv := NewPendingVotes()
	pv.MoveToNewHeight(101, valSet.Validators())
	pv.SetRoundProposal(4, prop)
	assert.False(t, pv.HasRoundProposal(0))
	assert.True(t, pv.HasRoundProposal(4))
	assert.True(t, pv.HasRoundProposal(4))
	assert.Nil(t, pv.RoundProposal(0))
	assert.Equal(t, pv.RoundProposal(4).Hash(), prop.Hash())
}
