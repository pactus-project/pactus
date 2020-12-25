package consensus

import (
	"testing"
	"time"

	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/vote"
)

func TestPrecommitWithNoProposal(t *testing.T) {

	cons := newTestConsensus(t, VAL3)

	cons.enterNewHeight(1)

	p1 := makeTestProposal(t, VAL1, 1, 0)

	testAddVote(t, cons, vote.VoteTypePrepare, 1, 0, p1.Block().Hash(), VAL1, false)
	testAddVote(t, cons, vote.VoteTypePrepare, 1, 0, p1.Block().Hash(), VAL2, false)
	testAddVote(t, cons, vote.VoteTypePrepare, 1, 0, p1.Block().Hash(), VAL4, false)
	shouldPublishProposalReqquest(t, cons)

	checkHRS(t, cons, 1, 0, hrs.StepTypePrecommit)

	time.Sleep(1 * time.Second)           // This will change block timestamp
	p2 := makeTestProposal(t, VAL1, 1, 0) // Invalid proposal
	cons.SetProposal(p2)
	shouldPublishProposalReqquest(t, cons)

	cons.SetProposal(p1)
	shouldPublishVote(t, cons, vote.VoteTypePrecommit, p1.Block().Hash())
}
