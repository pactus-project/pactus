package consensus

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/vote"
)

func TestConsensusSetProposalAfterCommit(t *testing.T) {
	setup(t)

	tConsX.enterNewHeight()
	tConsP.enterNewHeight()
	p := tConsX.LastProposal()

	commitBlockForAllStates(t)
	tConsP.SetProposal(p)
	assert.Nil(t, tConsP.LastProposal())
}

func TestGotoNextRoundWithoutProposal(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)
	commitBlockForAllStates(t)

	tConsP.enterNewHeight()

	testAddVote(t, tConsP, vote.VoteTypePrecommit, 3, 0, crypto.UndefHash, tIndexX, false)
	testAddVote(t, tConsP, vote.VoteTypePrecommit, 3, 0, crypto.UndefHash, tIndexY, false)
	testAddVote(t, tConsP, vote.VoteTypePrecommit, 3, 0, crypto.UndefHash, tIndexB, false)

	checkHRSWait(t, tConsP, 3, 1, hrs.StepTypePrepare)
}

func TestSecondProposalCommitted(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)
	commitBlockForAllStates(t)

	tConsX.enterNewHeight()
	tConsB.enterNewHeight()
	tConsP.enterNewHeight()
	tConsP.enterNewRound(1)

	// Now it's turn for Byzantine node to propose a block
	// Other nodes are going to not accept its proposal, even it is valid
	p1 := tConsB.LastProposal() // valid proposal for first round
	p2 := tConsP.LastProposal() // valid proposal for second round

	// Probably we have blocked Byzantine node
	//tConsX.SetProposal(p1)

	shouldPublishVote(t, tConsX, vote.VoteTypePrepare, crypto.UndefHash)
	testAddVote(t, tConsX, vote.VoteTypePrepare, 3, 0, crypto.UndefHash, tIndexY, false)
	testAddVote(t, tConsX, vote.VoteTypePrepare, 3, 0, p1.Block().Hash(), tIndexB, false)
	testAddVote(t, tConsX, vote.VoteTypePrepare, 3, 0, crypto.UndefHash, tIndexP, false)

	shouldPublishVote(t, tConsX, vote.VoteTypePrecommit, crypto.UndefHash)
	testAddVote(t, tConsX, vote.VoteTypePrecommit, 3, 0, crypto.UndefHash, tIndexY, false)
	testAddVote(t, tConsX, vote.VoteTypePrecommit, 3, 0, p1.Block().Hash(), tIndexB, false)
	testAddVote(t, tConsX, vote.VoteTypePrecommit, 3, 0, crypto.UndefHash, tIndexP, false)

	tConsX.SetProposal(p2)

	testAddVote(t, tConsX, vote.VoteTypePrepare, 3, 1, p2.Block().Hash(), tIndexY, false)
	testAddVote(t, tConsX, vote.VoteTypePrepare, 3, 1, crypto.UndefHash, tIndexB, false)
	testAddVote(t, tConsX, vote.VoteTypePrepare, 3, 1, p2.Block().Hash(), tIndexP, false)
	shouldPublishVote(t, tConsX, vote.VoteTypePrepare, p2.Block().Hash())

	testAddVote(t, tConsX, vote.VoteTypePrecommit, 3, 1, p2.Block().Hash(), tIndexY, false)
	testAddVote(t, tConsX, vote.VoteTypePrecommit, 3, 1, crypto.UndefHash, tIndexB, false)
	testAddVote(t, tConsX, vote.VoteTypePrecommit, 3, 1, p2.Block().Hash(), tIndexP, false)

	shouldPublishVote(t, tConsX, vote.VoteTypePrecommit, p2.Block().Hash())
	shouldPublishBlockAnnounce(t, tConsX, p2.Block().Hash())
}

func TestNetworkLagging1(t *testing.T) {
	setup(t)

	tConsP.enterNewHeight()

	h := 1
	r := 0
	p := makeProposal(t, h, r)
	// We don't set proposal for second validator here
	// tConsP.SetProposal(p)

	checkHRSWait(t, tConsP, h, r, hrs.StepTypePrepare)
	shouldPublishQueryProposal(t, tConsP, h, r)
	shouldPublishVote(t, tConsP, vote.VoteTypePrepare, crypto.UndefHash)

	testAddVote(t, tConsP, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexX, false)
	testAddVote(t, tConsP, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexY, false)
	testAddVote(t, tConsP, vote.VoteTypePrepare, h, r, crypto.UndefHash, tIndexB, false)

	// Proposal received now, set it
	tConsP.SetProposal(p)
	checkHRSWait(t, tConsP, h, r, hrs.StepTypePrecommit)
	shouldPublishVote(t, tConsP, vote.VoteTypePrecommit, p.Block().Hash())
}

func TestNetworkLagging2(t *testing.T) {
	setup(t)

	tConsP.enterNewHeight()
	tConsX.enterNewHeight()

	p1 := tConsX.LastProposal()
	// We don't set proposal for second validator here
	// tConsP.SetProposal(p1)

	// Networks lags and we don't receive prepare from val_1 and pre-commit from val_4
	testAddVote(t, tConsP, vote.VoteTypePrecommit, 1, 0, p1.Block().Hash(), tIndexX, false)
	testAddVote(t, tConsP, vote.VoteTypePrecommit, 1, 0, p1.Block().Hash(), tIndexY, false)

	checkHRS(t, tConsP, 1, 0, hrs.StepTypePropose)
	assert.Nil(t, tConsP.pendingVotes.roundVotes[0].Precommits.QuorumBlock())

	shouldPublishQueryProposal(t, tConsP, 1, 0)
	shouldPublishVote(t, tConsP, vote.VoteTypePrepare, crypto.UndefHash)

	// Proposal received now, set it
	tConsP.SetProposal(p1)

	shouldPublishVote(t, tConsP, vote.VoteTypePrepare, p1.Block().Hash())
	checkHRSWait(t, tConsP, 1, 0, hrs.StepTypePrepare)

	// We can't go to precommit stage, because we haven't prepared yet
	// But if we receive another vote we go to commit phase directly
	// Let's do it
	testAddVote(t, tConsP, vote.VoteTypePrecommit, 1, 0, p1.Block().Hash(), tIndexB, false)
	shouldPublishBlockAnnounce(t, tConsP, p1.Block().Hash())
}

func TestLateProposal(t *testing.T) {
	setup(t)

	tConsP.enterNewHeight()

	h := 1
	r := 0
	p := makeProposal(t, h, r)

	// tConsP is partitioned, so tConsP doesn't have the proposal
	testAddVote(t, tConsP, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexX, false)
	testAddVote(t, tConsP, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexY, false)
	testAddVote(t, tConsP, vote.VoteTypePrepare, h, r, crypto.UndefHash, tIndexB, false)

	testAddVote(t, tConsP, vote.VoteTypePrecommit, h, r, p.Block().Hash(), tIndexX, false)
	testAddVote(t, tConsP, vote.VoteTypePrecommit, h, r, p.Block().Hash(), tIndexY, false)
	testAddVote(t, tConsP, vote.VoteTypePrecommit, h, r, crypto.UndefHash, tIndexB, false)

	// Now partition healed.
	tConsP.SetProposal(p)
	shouldPublishVote(t, tConsP, vote.VoteTypePrecommit, p.Block().Hash())
	shouldPublishVote(t, tConsP, vote.VoteTypePrepare, p.Block().Hash())

	assert.True(t, tConsP.isCommitted)
}

func TestLateUndefVote(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)
	commitBlockForAllStates(t)

	tConsX.enterNewHeight()

	h := 3
	r := 0
	p := makeProposal(t, h, r) // Other nodes doesn't accept byzantine proposal

	// tConsP is partitioned, so tConsX doesn't have the proposal
	testAddVote(t, tConsX, vote.VoteTypePrepare, h, r, crypto.UndefHash, tIndexY, false)
	testAddVote(t, tConsX, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexB, false)

	testAddVote(t, tConsX, vote.VoteTypePrecommit, h, r, crypto.UndefHash, tIndexY, false)
	testAddVote(t, tConsX, vote.VoteTypePrecommit, h, r, p.Block().Hash(), tIndexB, false)

	// Now partition healed.
	testAddVote(t, tConsX, vote.VoteTypePrecommit, h, r, crypto.UndefHash, tIndexP, false)
	testAddVote(t, tConsX, vote.VoteTypePrecommit, h, r, crypto.UndefHash, tIndexP, false)

	checkHRSWait(t, tConsX, h, r+1, hrs.StepTypePropose)
}

func TestLateProposal2(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)
	commitBlockForAllStates(t)
	commitBlockForAllStates(t)

	tConsX.enterNewHeight()
	tConsP.enterNewHeight()

	p := tConsP.LastProposal()

	// tConsP is partitioned, so tConsX doesn't have the proposal
	testAddVote(t, tConsX, vote.VoteTypePrepare, 4, 0, crypto.UndefHash, tIndexY, false)
	testAddVote(t, tConsX, vote.VoteTypePrepare, 4, 0, crypto.UndefHash, tIndexB, false)

	testAddVote(t, tConsX, vote.VoteTypePrecommit, 4, 0, crypto.UndefHash, tIndexY, false)
	testAddVote(t, tConsX, vote.VoteTypePrecommit, 4, 0, crypto.UndefHash, tIndexB, false)

	checkHRSWait(t, tConsX, 4, 1, hrs.StepTypePrepare)

	// Now partition healed, but it's too late, We already moved to the next round
	tConsX.SetProposal(p)

	checkHRSWait(t, tConsX, 4, 1, hrs.StepTypePrepare)
}

func TestSetProposalForNextRoundWithoutFinishingTheFirstRound(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)

	tConsX.enterNewHeight()

	// Byzantine node sends proposal for second round (his turn)
	b, err := tConsB.state.ProposeBlock(1)
	assert.NoError(t, err)
	p := vote.NewProposal(2, 1, *b)
	tSigners[tIndexB].SignMsg(p)

	tConsX.SetProposal(p)
	assert.Nil(t, tConsX.LastProposal())
}

func TestEnterPrepareAfterPrecommit(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)
	commitBlockForAllStates(t)
	commitBlockForAllStates(t)

	tConsX.enterNewHeight()

	p := makeProposal(t, 4, 0)

	// tConsP is partitioned, so tConsX doesn't have the proposal
	testAddVote(t, tConsX, vote.VoteTypePrepare, 4, 0, crypto.UndefHash, tIndexY, false)
	testAddVote(t, tConsX, vote.VoteTypePrepare, 4, 0, crypto.UndefHash, tIndexP, false)
	checkHRSWait(t, tConsX, 4, 0, hrs.StepTypePrecommit)

	tConsX.SetProposal(p)
	tConsX.enterPrepare(0)

	votes := tConsX.RoundVotes(0)
	for _, v := range votes {
		if v.BlockHash() == p.Block().Hash() {
			assert.NoError(t, fmt.Errorf("The proposal should be ignored, because we have voted for undef before"))
		}
	}
}
