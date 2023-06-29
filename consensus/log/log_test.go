package log

import (
	"testing"

	"github.com/pactus-project/pactus/types/vote"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestMustGetRound(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	committee, _ := ts.GenerateTestCommittee(4)
	log := NewLog()
	log.MoveToNewHeight(committee.Validators())
	log.MustGetRoundMessages(4)
	assert.Nil(t, log.RoundMessages(5))
	assert.NotNil(t, log.RoundMessages(1))
	assert.NotNil(t, log.RoundMessages(4))
	assert.Equal(t, len(log.roundMessages), 5)
}

func addVotes(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	committee, signers := ts.GenerateTestCommittee(4)

	log := NewLog()
	log.MoveToNewHeight(committee.Validators())
	invalidVote, _ := ts.GenerateTestPrecommitVote(55, 5)
	err := log.AddVote(invalidVote) // invalid height
	assert.Error(t, err)

	v1, _ := ts.GenerateTestPrecommitVote(101, 5)
	err = log.AddVote(v1) // invalid signer
	assert.Error(t, err)

	validVote := vote.NewVote(vote.VoteTypePrepare, 101, 1, ts.RandomHash(), signers[0].Address())
	signers[0].SignMsg(validVote)

	duplicateVote := vote.NewVote(vote.VoteTypePrepare, 101, 1, ts.RandomHash(), signers[0].Address())
	signers[0].SignMsg(duplicateVote)

	err = log.AddVote(validVote)
	assert.NoError(t, err)

	// Definitely it is a duplicated error
	err = log.AddVote(duplicateVote)
	assert.Error(t, err) // duplicated vote error

	prepares := log.PrepareVoteSet(1)
	precommits := log.PrecommitVoteSet(1)
	assert.Equal(t, prepares.Len(), 2)   //  Vote + Duplicated
	assert.Equal(t, precommits.Len(), 0) // no precommit votes
	assert.Equal(t, len(log.RoundMessages(1).AllVotes()), 2)
	assert.True(t, log.HasVote(duplicateVote.Hash()))
	assert.True(t, log.HasVote(validVote.Hash()))
	assert.False(t, log.HasVote(invalidVote.Hash()))
}

func TestSetRoundProposal(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	committee, _ := ts.GenerateTestCommittee(4)
	prop, _ := ts.GenerateTestProposal(101, 0)
	log := NewLog()
	log.MoveToNewHeight(committee.Validators())
	log.SetRoundProposal(4, prop)
	assert.False(t, log.HasRoundProposal(0))
	assert.True(t, log.HasRoundProposal(4))
	assert.True(t, log.HasRoundProposal(4))
	assert.Nil(t, log.RoundProposal(0))
	assert.Nil(t, log.RoundProposal(5))
	assert.Equal(t, log.RoundProposal(4).Hash(), prop.Hash())
}

func TestCanVote(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	committee, signers := ts.GenerateTestCommittee(4)
	log := NewLog()
	log.MoveToNewHeight(committee.Validators())

	addr := ts.RandomAddress()
	assert.True(t, log.CanVote(signers[0].Address()))
	assert.False(t, log.CanVote(addr))
}
