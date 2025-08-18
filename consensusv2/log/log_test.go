package log

import (
	"encoding/hex"
	"testing"

	"github.com/pactus-project/pactus/types/vote"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestMustGetRound(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	cmt, _ := ts.GenerateTestCommittee(6)
	log := NewLog()
	log.MoveToNewHeight(cmt.Validators())
	assert.NotNil(t, log.RoundMessages(ts.RandRound()))
}

func TestAddValidVote(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	cmt, valKeys := ts.GenerateTestCommittee(6)
	log := NewLog()
	log.MoveToNewHeight(cmt.Validators())
	h := ts.RandHeight()
	r := ts.RandRound()

	precommits := log.PrecommitVoteSet(r)
	preVotes := log.CPPreVoteVoteSet(r)
	mainVotes := log.CPMainVoteVoteSet(r)

	vote1 := vote.NewPrepareVote(ts.RandHash(), h, r, valKeys[0].Address())
	vote2 := vote.NewPrecommitVote(ts.RandHash(), h, r, valKeys[0].Address())
	vote3 := vote.NewCPPreVote(ts.RandHash(), h, r, 0, vote.CPValueYes, &vote.JustInitYes{}, valKeys[0].Address())
	vote4 := vote.NewCPMainVote(ts.RandHash(), h, r, 0, vote.CPValueNo, &vote.JustInitYes{}, valKeys[0].Address())

	for _, v := range []*vote.Vote{vote1, vote2, vote3, vote4} {
		ts.HelperSignVote(valKeys[0], v)

		added, err := log.AddVote(v)
		assert.NoError(t, err)
		assert.True(t, added)
	}

	assert.True(t, log.HasVote(vote1.Hash()))
	assert.True(t, log.HasVote(vote2.Hash()))
	assert.True(t, log.HasVote(vote3.Hash()))
	assert.True(t, log.HasVote(vote4.Hash()))
	assert.False(t, log.HasVote(ts.RandHash()))

	assert.Contains(t, precommits.AllVotes(), vote2)
	assert.Contains(t, preVotes.AllVotes(), vote3)
	assert.Contains(t, mainVotes.AllVotes(), vote4)
}

func TestAddInvalidVoteType(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	cmt, _ := ts.GenerateTestCommittee(6)
	log := NewLog()
	log.MoveToNewHeight(cmt.Validators())

	data, _ := hex.DecodeString("A701050218320301045820BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB" +
		"055501AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA06f607f6")
	invVote := new(vote.Vote)
	err := invVote.UnmarshalCBOR(data)
	assert.NoError(t, err)

	added, err := log.AddVote(invVote)
	assert.Error(t, err)
	assert.False(t, added)
	assert.False(t, log.HasVote(invVote.Hash()))
}

func TestSetRoundProposal(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	cmt, _ := ts.GenerateTestCommittee(6)
	prop := ts.GenerateTestProposal(101, 0)
	log := NewLog()
	log.MoveToNewHeight(cmt.Validators())
	log.SetRoundProposal(4, prop)
	assert.False(t, log.HasRoundProposal(0))
	assert.True(t, log.HasRoundProposal(4))
	assert.True(t, log.HasRoundProposal(4))
	assert.Nil(t, log.RoundProposal(0))
	assert.Nil(t, log.RoundProposal(5))
	assert.Equal(t, prop.Hash(), log.RoundProposal(4).Hash())
}

func TestCanVote(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	cmt, valKeys := ts.GenerateTestCommittee(6)
	log := NewLog()
	log.MoveToNewHeight(cmt.Validators())

	addr := ts.RandAccAddress()
	assert.True(t, log.CanVote(valKeys[0].Address()))
	assert.False(t, log.CanVote(addr))
}
