package consensus

import (
	"testing"

	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/vote"
)

func TestPrecommitWithNoProposal(t *testing.T) {

	cons := newTestConsensus(t, VAL3)

	cons.enterNewHeight(1)

	h := crypto.GenerateTestHash()
	testAddVote(t, cons, vote.VoteTypePrevote, 1, 0, h, VAL1, false) // invalid votes
	testAddVote(t, cons, vote.VoteTypePrevote, 1, 0, h, VAL2, false) // invalid votes
	testAddVote(t, cons, vote.VoteTypePrevote, 1, 0, h, VAL4, false) // invalid votes
	shouldPublishProposalReqquest(t, cons)
	shouldPublishUndefVote(t, cons)

	checkHRS(t, cons, 1, 0, hrs.StepTypePrecommit)

	cons1 := newTestConsensus(t, VAL1)
	b := cons1.state.ProposeBlock()
	p := vote.NewProposal(1, 0, b)
	tSigners[VAL1].SignMsg(p)

	// Here we have valid proposal but invalid votes.
	cons.SetProposal(p)
	shouldPublishUndefVote(t, cons)
}
