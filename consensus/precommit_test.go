package consensus

import (
	"testing"

	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/vote"
)

func TestPrecommitNoProposal(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)

	h := 2
	r := 0
	p := makeProposal(t, h, r)

	tConsP.enterNewHeight()
	checkHRSWait(t, tConsP, h, r, hrs.StepTypePropose) // We can't prepared, because we don't have proposal
	shouldPublishQueryProposal(t, tConsP, h, r)
	shouldPublishVote(t, tConsP, vote.VoteTypePrepare, crypto.UndefHash)

	// Still no proposal
	testAddVote(t, tConsP, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexX, false)
	testAddVote(t, tConsP, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexY, false)
	testAddVote(t, tConsP, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexB, false)

	checkHRSWait(t, tConsP, h, r, hrs.StepTypePrecommit)
	shouldPublishQueryProposal(t, tConsP, h, r)

	// Set proposal now
	tConsP.SetProposal(p)
	shouldPublishVote(t, tConsP, vote.VoteTypePrecommit, p.Block().Hash())
}

// This is a worse case scenario
func TestPrecommitNoProposalWithPrecommitQuorom(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)

	h := 2
	r := 0
	p := makeProposal(t, h, r)

	tConsP.enterNewHeight()
	checkHRSWait(t, tConsP, h, r, hrs.StepTypePropose)
	shouldPublishQueryProposal(t, tConsP, h, r)
	shouldPublishVote(t, tConsP, vote.VoteTypePrepare, crypto.UndefHash)

	// Still no proposal
	testAddVote(t, tConsP, vote.VoteTypePrecommit, h, r, p.Block().Hash(), tIndexX, false)
	testAddVote(t, tConsP, vote.VoteTypePrecommit, h, r, p.Block().Hash(), tIndexY, false)
	testAddVote(t, tConsP, vote.VoteTypePrecommit, h, r, p.Block().Hash(), tIndexB, false)

	checkHRSWait(t, tConsP, h, r, hrs.StepTypeCommit)

	// Set proposal now
	tConsP.SetProposal(p)
	shouldPublishVote(t, tConsP, vote.VoteTypePrepare, p.Block().Hash())

	testAddVote(t, tConsP, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexX, false)
	testAddVote(t, tConsP, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexY, false)
	testAddVote(t, tConsP, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexB, false)

	shouldPublishBlockAnnounce(t, tConsP, p.Block().Hash())
	shouldPublishVote(t, tConsP, vote.VoteTypePrecommit, p.Block().Hash())
}
