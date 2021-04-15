package pending_votes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/committee"
	"github.com/zarbchain/zarb-go/consensus/proposal"
	"github.com/zarbchain/zarb-go/consensus/vote"
	"github.com/zarbchain/zarb-go/crypto"
)

func TestMustGetRound(t *testing.T) {
	committee, _ := committee.GenerateTestCommittee()
	pv := NewPendingVotes()
	pv.MoveToNewHeight(101, committee.Validators())
	pv.MustGetRoundVotes(4)
	assert.Nil(t, pv.GetRoundVotes(5))
	assert.NotNil(t, pv.GetRoundVotes(1))
	assert.NotNil(t, pv.GetRoundVotes(4))
	assert.Equal(t, len(pv.roundVotes), 5)
}

func TestAddVotes(t *testing.T) {
	committee, signers := committee.GenerateTestCommittee()

	pv := NewPendingVotes()
	pv.MoveToNewHeight(101, committee.Validators())
	invalidVote, _ := vote.GenerateTestPrecommitVote(55, 5)
	err := pv.AddVote(invalidVote) // invalid height
	assert.Error(t, err)

	v1, _ := vote.GenerateTestPrecommitVote(101, 5)
	err = pv.AddVote(v1) // invalid signer
	assert.Error(t, err)

	validVote := vote.NewVote(vote.VoteTypePrepare, 101, 1, crypto.GenerateTestHash(), signers[0].Address())
	signers[0].SignMsg(validVote)

	duplicateVote := vote.NewVote(vote.VoteTypePrepare, 101, 1, crypto.GenerateTestHash(), signers[0].Address())
	signers[0].SignMsg(duplicateVote)

	err = pv.AddVote(validVote)
	assert.NoError(t, err)

	// Definitely it is a duplicated error
	err = pv.AddVote(duplicateVote)
	assert.Error(t, err) // duplicated vote error

	prepares := pv.PrepareVoteSet(1)
	precommits := pv.PrecommitVoteSet(1)
	assert.Equal(t, prepares.Len(), 2)   //  Vote + Duplicated
	assert.Equal(t, precommits.Len(), 0) // no precommit votes
	assert.Equal(t, len(pv.GetRoundVotes(1).AllVotes()), 2)
	assert.True(t, pv.HasVote(duplicateVote.Hash()))
	assert.True(t, pv.HasVote(validVote.Hash()))
	assert.False(t, pv.HasVote(invalidVote.Hash()))
}

func TestSetRoundProposal(t *testing.T) {
	committee, _ := committee.GenerateTestCommittee()
	prop, _ := proposal.GenerateTestProposal(101, 0)
	pv := NewPendingVotes()
	pv.MoveToNewHeight(101, committee.Validators())
	pv.SetRoundProposal(4, prop)
	assert.False(t, pv.HasRoundProposal(0))
	assert.True(t, pv.HasRoundProposal(4))
	assert.True(t, pv.HasRoundProposal(4))
	assert.Nil(t, pv.RoundProposal(0))
	assert.Nil(t, pv.RoundProposal(5))
	assert.Equal(t, pv.RoundProposal(4).Hash(), prop.Hash())
}

func TestCanVote(t *testing.T) {
	committee, signers := committee.GenerateTestCommittee()
	pv := NewPendingVotes()
	pv.MoveToNewHeight(101, committee.Validators())

	addr, _, _ := crypto.GenerateTestKeyPair()
	assert.True(t, pv.CanVote(signers[0].Address()))
	assert.False(t, pv.CanVote(addr))
}
