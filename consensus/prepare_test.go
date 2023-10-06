package consensus

import (
	"testing"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/vote"
)

func TestChangeProposerTimeout(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consP)
	td.changeProposerTimeout(td.consP)

	td.shouldPublishVote(t, td.consP, vote.VoteTypeCPPreVote, hash.UndefHash)
}

func TestQueryProposal(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t)
	h := uint32(2)
	r := int16(1)

	td.enterNewHeight(td.consP)
	td.enterNextRound(td.consP)
	td.queryProposalTimeout(td.consP)

	td.shouldPublishQueryProposal(t, td.consP, h)
	td.shouldPublishQueryVote(t, td.consP, h, r)
}

func TestQueryVotes(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t)
	h := uint32(2)
	r := int16(1)

	td.enterNewHeight(td.consP)
	td.enterNextRound(td.consP)

	p := td.makeProposal(t, h, r)
	td.consP.SetProposal(p)

	// consP has a valid proposal but not enough votes.
	td.queryProposalTimeout(td.consP)
	td.shouldPublishQueryVote(t, td.consP, h, r)
}

func TestGoToChangeProposerFromPrepare(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t)

	td.enterNewHeight(td.consP)

	td.addCPPreVote(td.consP, hash.UndefHash, 2, 0, 0, vote.CPValueOne, &vote.JustInitOne{}, tIndexX)
	td.addCPPreVote(td.consP, hash.UndefHash, 2, 0, 0, vote.CPValueOne, &vote.JustInitOne{}, tIndexY)

	// should move to the change proposer phase, even if it has the proposal and
	// its timer has not expired, if it has received 1/3 of the change-proposer votes.
	p := td.makeProposal(t, 2, 0)
	td.consP.SetProposal(p)
	td.shouldPublishVote(t, td.consP, vote.VoteTypeCPPreVote, hash.UndefHash)
}
