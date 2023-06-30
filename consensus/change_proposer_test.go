package consensus

import (
	"testing"
	"time"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/stretchr/testify/assert"
)

func TestChangeProposer(t *testing.T) {
	td := setup(t)

	td.consP.config.ChangeProposerTimeout = 100 * time.Millisecond
	td.enterNewHeight(td.consP)

	td.shouldPublishVote(t, td.consP, vote.VoteTypeChangeProposer, hash.UndefHash)
}

func TestGotoNewRound(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t)

	td.consP.config.ChangeProposerTimeout = 100 * time.Millisecond
	td.enterNewHeight(td.consP)

	td.shouldPublishVote(t, td.consP, vote.VoteTypeChangeProposer, hash.UndefHash)
	td.addVote(td.consP, vote.VoteTypeChangeProposer, 2, 0, hash.UndefHash, tIndexX)
	td.addVote(td.consP, vote.VoteTypeChangeProposer, 2, 0, hash.UndefHash, tIndexY)

	td.checkHeightRound(t, td.consP, 2, 1)
}

func TestSetProposalAfterChangeProposer(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t)

	td.enterNewHeight(td.consP)

	td.addVote(td.consP, vote.VoteTypeChangeProposer, 2, 0, hash.UndefHash, tIndexX)
	td.addVote(td.consP, vote.VoteTypeChangeProposer, 2, 0, hash.UndefHash, tIndexY)

	p := td.makeProposal(t, 2, 0)
	td.consP.SetProposal(p)
	assert.Nil(t, td.consP.RoundProposal(0))

	td.shouldPublishVote(t, td.consP, vote.VoteTypeChangeProposer, hash.UndefHash)
}

func TestRemoveProposalAfterChangeProposer(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t)

	td.enterNewHeight(td.consP)
	p := td.makeProposal(t, 2, 0)
	td.consP.SetProposal(p)
	assert.NotNil(t, td.consP.RoundProposal(0))

	td.addVote(td.consP, vote.VoteTypeChangeProposer, 2, 0, hash.UndefHash, tIndexX)
	td.addVote(td.consP, vote.VoteTypeChangeProposer, 2, 0, hash.UndefHash, tIndexY)

	td.shouldPublishVote(t, td.consP, vote.VoteTypeChangeProposer, hash.UndefHash)
	td.checkHeightRound(t, td.consP, 2, 1)
	assert.Nil(t, td.consP.RoundProposal(0))
}
