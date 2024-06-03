package consensus

import (
	"testing"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/sync/bundle/message"
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

	td.enterNewHeight(td.consP)
	td.enterNextRound(td.consP)
	td.queryProposalTimeout(td.consP)

	td.shouldPublishQueryProposal(t, td.consP, h)
	td.shouldNotPublish(t, td.consP, message.TypeQueryVote)
}

func TestQueryVotes(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t)
	td.commitBlockForAllStates(t)
	h := uint32(3)
	r := int16(1)

	td.enterNewHeight(td.consP)
	td.enterNextRound(td.consP)

	// consP is the proposer for this round, but there are not enough votes.
	td.queryProposalTimeout(td.consP)
	td.shouldPublishProposal(t, td.consP, h, r)
	td.shouldPublishQueryVote(t, td.consP, h, r)
	td.shouldNotPublish(t, td.consP, message.TypeQueryProposal)
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
