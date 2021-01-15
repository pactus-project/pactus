package consensus

import (
	"testing"

	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/vote"
)

func TestPrecommitWithNoProposal(t *testing.T) {
	setup(t)

	tConsX.enterNewHeight()
	p1 := tConsX.LastProposal()

	tConsP.enterNewHeight()
	checkHRSWait(t, tConsP, 1, 0, hrs.StepTypePropose)

	testAddVote(t, tConsP, vote.VoteTypePrepare, 1, 0, p1.Block().Hash(), tIndexX, false)
	testAddVote(t, tConsP, vote.VoteTypePrepare, 1, 0, p1.Block().Hash(), tIndexY, false)
	testAddVote(t, tConsP, vote.VoteTypePrepare, 1, 0, p1.Block().Hash(), tIndexB, false)

	checkHRSWait(t, tConsP, 1, 0, hrs.StepTypePrecommit)
	shouldPublishQueryProposal(t, tConsP, 1, 0)
	shouldPublishVote(t, tConsP, vote.VoteTypePrepare, crypto.UndefHash)

	tConsP.SetProposal(p1)

	shouldPublishVote(t, tConsP, vote.VoteTypePrecommit, p1.Block().Hash())
	shouldPublishVote(t, tConsP, vote.VoteTypePrepare, p1.Block().Hash())
}
