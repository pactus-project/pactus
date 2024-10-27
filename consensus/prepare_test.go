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
	height := uint32(2)

	td.enterNewHeight(td.consP)
	td.enterNextRound(td.consP)
	td.queryProposalTimeout(td.consP)

	td.shouldPublishQueryProposal(t, td.consP, height)
	td.shouldNotPublish(t, td.consP, message.TypeQueryVote)
}

func TestQueryVote(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t)
	td.commitBlockForAllStates(t)
	height := uint32(3)
	round := int16(1)

	td.enterNewHeight(td.consP)
	td.enterNextRound(td.consP)

	// consP is the proposer for this round, but there are not enough votes.
	td.queryProposalTimeout(td.consP)
	td.shouldPublishProposal(t, td.consP, height, round)
	td.shouldPublishQueryVote(t, td.consP, height, round)
	td.shouldNotPublish(t, td.consP, message.TypeQueryProposal)
}

func TestGoToChangeProposerFromPrepare(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t)

	td.enterNewHeight(td.consP)

	td.addCPPreVote(td.consP, hash.UndefHash, 2, 0, vote.CPValueYes, &vote.JustInitYes{}, tIndexX)
	td.addCPPreVote(td.consP, hash.UndefHash, 2, 0, vote.CPValueYes, &vote.JustInitYes{}, tIndexY)

	// should move to the change proposer phase, even if it has the proposal and
	// its timer has not expired, if it has received 1/3 of the change-proposer votes.
	p := td.makeProposal(t, 2, 0)
	td.consP.SetProposal(p)
	td.shouldPublishVote(t, td.consP, vote.VoteTypeCPPreVote, hash.UndefHash)
}
