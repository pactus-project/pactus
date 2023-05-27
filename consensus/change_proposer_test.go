package consensus

import (
	"testing"
	"time"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/stretchr/testify/assert"
)

func TestChangeProposer(t *testing.T) {
	setup(t)

	tConsP.config.ChangeProposerTimeout = 100 * time.Millisecond
	testEnterNewHeight(tConsP)

	shouldPublishVote(t, tConsP, vote.VoteTypeChangeProposer, hash.UndefHash)
}

func TestGotoNewRound(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)

	tConsP.config.ChangeProposerTimeout = 100 * time.Millisecond
	testEnterNewHeight(tConsP)

	shouldPublishVote(t, tConsP, vote.VoteTypeChangeProposer, hash.UndefHash)
	testAddVote(tConsP, vote.VoteTypeChangeProposer, 2, 0, hash.UndefHash, tIndexX)
	testAddVote(tConsP, vote.VoteTypeChangeProposer, 2, 0, hash.UndefHash, tIndexY)

	checkHeightRound(t, tConsP, 2, 1)
}

func TestSetProposalAfterChangeProposer(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)

	testEnterNewHeight(tConsP)

	testAddVote(tConsP, vote.VoteTypeChangeProposer, 2, 0, hash.UndefHash, tIndexX)
	testAddVote(tConsP, vote.VoteTypeChangeProposer, 2, 0, hash.UndefHash, tIndexY)

	p := makeProposal(t, 2, 0)
	tConsP.SetProposal(p)
	assert.Nil(t, tConsP.RoundProposal(0))

	shouldPublishVote(t, tConsP, vote.VoteTypeChangeProposer, hash.UndefHash)
}

func TestRemoveProposalAfterChangeProposer(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)

	testEnterNewHeight(tConsP)
	p := makeProposal(t, 2, 0)
	tConsP.SetProposal(p)
	assert.NotNil(t, tConsP.RoundProposal(0))

	testAddVote(tConsP, vote.VoteTypeChangeProposer, 2, 0, hash.UndefHash, tIndexX)
	testAddVote(tConsP, vote.VoteTypeChangeProposer, 2, 0, hash.UndefHash, tIndexY)

	shouldPublishVote(t, tConsP, vote.VoteTypeChangeProposer, hash.UndefHash)
	checkHeightRound(t, tConsP, 2, 1)
	assert.Nil(t, tConsP.RoundProposal(0))
}
