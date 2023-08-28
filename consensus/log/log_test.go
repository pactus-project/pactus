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
	assert.NotNil(t, log.RoundMessages(ts.RandRound()))
}

func TestAddValidVote(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	committee, signers := ts.GenerateTestCommittee(4)
	log := NewLog()
	log.MoveToNewHeight(committee.Validators())
	h := ts.RandHeight()
	r := ts.RandRound()

	prepares := log.PrepareVoteSet(r)
	precommits := log.PrecommitVoteSet(r)
	preVotes := log.CPPreVoteVoteSet(r)
	mainVotes := log.CPMainVoteVoteSet(r)

	v1 := vote.NewPrepareVote(ts.RandHash(), h, r, signers[0].Address())
	v2 := vote.NewPrecommitVote(ts.RandHash(), h, r, signers[0].Address())
	v3 := vote.NewCPPreVote(ts.RandHash(), h, r, 0, vote.CPValueOne, &vote.JustInitOne{}, signers[0].Address())
	v4 := vote.NewCPMainVote(ts.RandHash(), h, r, 0, vote.CPValueZero, &vote.JustInitOne{}, signers[0].Address())

	for _, v := range []*vote.Vote{v1, v2, v3, v4} {
		signers[0].SignMsg(v)

		added, err := log.AddVote(v)
		assert.NoError(t, err)
		assert.True(t, added)
	}

	assert.True(t, log.HasVote(v1.Hash()))
	assert.True(t, log.HasVote(v2.Hash()))
	assert.True(t, log.HasVote(v3.Hash()))
	assert.True(t, log.HasVote(v4.Hash()))
	assert.False(t, log.HasVote(ts.RandHash()))

	assert.Contains(t, prepares.AllVotes(), v1)
	assert.Contains(t, precommits.AllVotes(), v2)
	assert.Contains(t, preVotes.AllVotes(), v3)
	assert.Contains(t, mainVotes.AllVotes(), v4)
}

func TestAddInvalidVote(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	committee, signers := ts.GenerateTestCommittee(4)
	log := NewLog()
	log.MoveToNewHeight(committee.Validators())
	h := ts.RandHeight()
	r := ts.RandRound()

	invVote := vote.NewVote(5, ts.RandHash(), h, r, signers[0].Address())
	signers[0].SignMsg(invVote)

	added, err := log.AddVote(invVote)
	assert.Error(t, err)
	assert.False(t, added)
	assert.False(t, log.HasVote(invVote.Hash()))
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

	addr := ts.RandAddress()
	assert.True(t, log.CanVote(signers[0].Address()))
	assert.False(t, log.CanVote(addr))
}
