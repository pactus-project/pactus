package consensusv2

import (
	"testing"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/types/vote"
)

// PASSED-2

func TestPrecommitStrongCommit(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t)
	height := uint32(2)
	round := int16(0)

	td.enterNewHeight(td.consP)
	prop := td.makeProposal(t, height, round)
	propBlockHash := prop.Block().Hash()

	td.addPrecommitVote(td.consP, propBlockHash, height, round, tIndexX)
	td.addPrecommitVote(td.consP, propBlockHash, height, round, tIndexY)
	td.addPrecommitVote(td.consP, propBlockHash, height, round, tIndexB)

	td.consP.SetProposal(prop)

	td.shouldPublishVote(t, td.consP, vote.VoteTypePrecommit, propBlockHash)
	td.shouldPublishBlockAnnounce(t, td.consP, propBlockHash)
	td.shouldNotPublish(t, td.consP, message.TypeQueryProposal)
	td.shouldNotPublish(t, td.consP, message.TypeQueryVote)
}

func TestPrecommitQueryProposal(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t)
	height := uint32(2)
	round := int16(0)

	// ConsP is not the proposer for this round.
	td.enterNewHeight(td.consP)
	td.queryProposalTimeout(td.consP)

	td.shouldPublishQueryProposal(t, td.consP, height, round)
	td.shouldNotPublish(t, td.consP, message.TypeQueryVote)
}

func TestPrecommitQueryVote(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t)
	height := uint32(2)
	round := int16(0)

	// ConsY is  the proposer for this round.
	td.enterNewHeight(td.consY)
	td.queryProposalTimeout(td.consY)

	td.shouldNotPublish(t, td.consY, message.TypeQueryProposal)
	td.shouldPublishQueryVote(t, td.consY, height, round)
}

func TestPrecommitChangeProposerTimeout(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consP)
	td.changeProposerTimeout(td.consP)

	td.shouldPublishVote(t, td.consP, vote.VoteTypeCPPreVote, hash.UndefHash)
}

func TestPrecommitChangeProposer(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t)
	height := uint32(2)
	round := int16(0)

	td.enterNewHeight(td.consP)
	prop := td.makeProposal(t, height, round)
	td.consP.SetProposal(prop)

	td.addCPPreVote(td.consP, hash.UndefHash, height, round, vote.CPValueYes, &vote.JustInitYes{}, tIndexX)
	td.addCPPreVote(td.consP, hash.UndefHash, height, round, vote.CPValueYes, &vote.JustInitYes{}, tIndexY)

	// should move to the change proposer phase, even if it has the proposal and
	// its timer has not expired, if it has received 1/3 of the change-proposer votes.
	td.shouldPublishVote(t, td.consP, vote.VoteTypeCPPreVote, hash.UndefHash)
}

func TestPrecommitQueryProposalWithCert(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t)
	height := uint32(2)
	round := int16(0)

	td.enterNewHeight(td.consP)
	td.consP.cpDecidedCert = td.GenerateTestVoteCertificate(height)

	td.consP.currentState.decide()

	td.shouldPublishQueryProposal(t, td.consP, height, round)
	td.shouldNotPublish(t, td.consP, message.TypeQueryVote)
}

func TestPrecommitQueryVoteWithCert(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t)
	height := uint32(2)
	round := int16(0)

	td.enterNewHeight(td.consP)
	prop := td.makeProposal(t, height, round)
	td.consP.SetProposal(prop)
	td.consP.cpDecidedCert = td.GenerateTestVoteCertificate(height)

	td.consP.currentState.decide()

	td.shouldNotPublish(t, td.consP, message.TypeQueryProposal)
	td.shouldPublishQueryVote(t, td.consP, height, round)
}
