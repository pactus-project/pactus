package log

import (
	"encoding/hex"
	"testing"

	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMustGetRound(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	log := NewLog()
	log.MoveToNewHeight(nil)
	assert.NotNil(t, log.RoundMessages(ts.RandRound()))
}

func TestAddValidVote(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	valKey := ts.RandValKey()
	val := ts.GenerateTestValidator(
		testsuite.ValidatorWithPublicKey(valKey.PublicKey()),
	)
	log := NewLog()
	log.MoveToNewHeight([]*validator.Validator{val})
	height := ts.RandHeight()
	round := ts.RandRound()

	prepares := log.PrepareVoteSet(round)
	precommits := log.PrecommitVoteSet(round)
	preVotes := log.CPPreVoteVoteSet(round)
	mainVotes := log.CPMainVoteVoteSet(round)

	vote1 := vote.NewPrepareVote(ts.RandHash(), height, round, val.Address())
	vote2 := vote.NewPrecommitVote(ts.RandHash(), height, round, val.Address())
	vote3 := vote.NewCPPreVote(ts.RandHash(), height, round, 0, vote.CPValueYes, &vote.JustInitYes{}, val.Address())
	vote4 := vote.NewCPMainVote(ts.RandHash(), height, round, 0, vote.CPValueNo, &vote.JustInitYes{}, val.Address())

	for _, v := range []*vote.Vote{vote1, vote2, vote3, vote4} {
		ts.HelperSignVote(valKey, v)

		added, err := log.AddVote(v)
		require.NoError(t, err)
		assert.True(t, added)
	}

	assert.True(t, log.HasVote(vote1.Hash()))
	assert.True(t, log.HasVote(vote2.Hash()))
	assert.True(t, log.HasVote(vote3.Hash()))
	assert.True(t, log.HasVote(vote4.Hash()))
	assert.False(t, log.HasVote(ts.RandHash()))

	assert.Contains(t, prepares.AllVotes(), vote1)
	assert.Contains(t, precommits.AllVotes(), vote2)
	assert.Contains(t, preVotes.AllVotes(), vote3)
	assert.Contains(t, mainVotes.AllVotes(), vote4)
}

func TestAddInvalidVoteType(t *testing.T) {
	log := NewLog()
	log.MoveToNewHeight(nil)

	data, _ := hex.DecodeString("A7010F0218320301045820BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB" +
		"055501AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA06f607f6")
	invVote := new(vote.Vote)
	err := invVote.UnmarshalCBOR(data)
	require.NoError(t, err)

	added, err := log.AddVote(invVote)
	require.ErrorContains(t, err, "unexpected vote type: 15")
	assert.False(t, added)
	assert.False(t, log.HasVote(invVote.Hash()))
}

func TestSetRoundProposal(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	height := ts.RandHeight()
	round := ts.RandRound()

	prop := ts.GenerateTestProposal(height, round)
	log := NewLog()
	log.MoveToNewHeight(nil)
	log.SetRoundProposal(round, prop)
	assert.False(t, log.HasRoundProposal(round+1))
	assert.True(t, log.HasRoundProposal(round))
	assert.Nil(t, log.RoundProposal(round+1))
	assert.Equal(t, prop.Hash(), log.RoundProposal(round).Hash())
}

func TestCanVote(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	val := ts.GenerateTestValidator()
	log := NewLog()
	log.MoveToNewHeight([]*validator.Validator{val})

	assert.True(t, log.CanVote(val.Address()))
	assert.False(t, log.CanVote(ts.RandValAddress()))
}
