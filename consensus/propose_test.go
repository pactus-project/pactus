package consensus

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/proposal"
	"github.com/zarbchain/zarb-go/vote"
)

func TestSetProposalInvalidProposer(t *testing.T) {
	setup(t)

	tConsY.enterNewHeight()
	assert.Nil(t, tConsY.RoundProposal(0))

	addr := tSigners[tIndexB].Address()
	b, _ := block.GenerateTestBlock(&addr, nil)
	p := proposal.NewProposal(1, 0, *b)

	tConsY.SetProposal(p)
	assert.Nil(t, tConsY.RoundProposal(0))

	tSigners[tIndexB].SignMsg(p) // Invalid signature
	tConsY.SetProposal(p)
	assert.Nil(t, tConsY.RoundProposal(0))
}

func TestSetProposalInvalidBlock(t *testing.T) {
	setup(t)

	a := tSigners[tIndexB].Address()
	invBlock, _ := block.GenerateTestBlock(&a, nil)
	p := proposal.NewProposal(1, 2, *invBlock)
	tSigners[tIndexB].SignMsg(p)

	tConsY.enterNewHeight()
	tConsY.enterNewRound(2)
	tConsY.SetProposal(p)
	assert.Nil(t, tConsY.RoundProposal(2))
}

func TestSetProposalInvalidHeight(t *testing.T) {
	setup(t)

	a := tSigners[tIndexB].Address()
	invBlock, _ := block.GenerateTestBlock(&a, nil)
	p := proposal.NewProposal(2, 0, *invBlock)
	tSigners[tIndexB].SignMsg(p)

	tConsY.enterNewHeight()
	tConsY.SetProposal(p)
	assert.Nil(t, tConsY.RoundProposal(2))
}

func TestConsensusSetProposalAfterCommit(t *testing.T) {
	setup(t)

	p := makeProposal(t, 1, 0)

	tConsP.enterNewHeight()
	commitBlockForAllStates(t)
	tConsP.SetProposal(p)
	assert.Nil(t, tConsP.RoundProposal(0))
}

func TestGotoNextRoundWithoutProposal(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)
	commitBlockForAllStates(t)

	tConsP.enterNewHeight()

	testAddVote(t, tConsP, vote.VoteTypePrecommit, 3, 0, crypto.UndefHash, tIndexX, false)
	testAddVote(t, tConsP, vote.VoteTypePrecommit, 3, 0, crypto.UndefHash, tIndexY, false)
	testAddVote(t, tConsP, vote.VoteTypePrecommit, 3, 0, crypto.UndefHash, tIndexB, false)

	checkHRS(t, tConsP, 3, 1, hrs.StepTypePrepare)
}

func TestSecondProposalCommitted(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)
	commitBlockForAllStates(t)

	tConsX.enterNewHeight()

	// Now it's turn for Byzantine node to propose a block
	// Other nodes are going to not accept its proposal, even it is valid
	p1 := makeProposal(t, 3, 0) // valid proposal for round 0, byzantine proposer
	p2 := makeProposal(t, 3, 1) // valid proposal for round 1, partitioned proposer

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

	shouldPublishVote(t, tConsX, vote.VoteTypePrepare, p2.Block().Hash())
	testAddVote(t, tConsX, vote.VoteTypePrepare, 3, 1, p2.Block().Hash(), tIndexY, false)
	testAddVote(t, tConsX, vote.VoteTypePrepare, 3, 1, crypto.UndefHash, tIndexB, false)
	testAddVote(t, tConsX, vote.VoteTypePrepare, 3, 1, p2.Block().Hash(), tIndexP, false)

	shouldPublishVote(t, tConsX, vote.VoteTypePrecommit, p2.Block().Hash())
	testAddVote(t, tConsX, vote.VoteTypePrecommit, 3, 1, p2.Block().Hash(), tIndexY, false)
	testAddVote(t, tConsX, vote.VoteTypePrecommit, 3, 1, crypto.UndefHash, tIndexB, false)
	testAddVote(t, tConsX, vote.VoteTypePrecommit, 3, 1, p2.Block().Hash(), tIndexP, false)

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

	// Now let's set the proposal
	tConsP.SetProposal(p)
	checkHRS(t, tConsP, h, r, hrs.StepTypePrecommit)
	shouldPublishVote(t, tConsP, vote.VoteTypePrecommit, p.Block().Hash())
}

func TestNetworkLagging2(t *testing.T) {
	setup(t)

	h := 1
	r := 0
	p1 := makeProposal(t, h, r)

	tConsP.enterNewHeight()
	// We don't set proposal for second validator here
	// tConsP.SetProposal(p1)

	// Networks lags and we don't receive prepare from val_1 and pre-commit from val_4
	testAddVote(t, tConsP, vote.VoteTypePrecommit, h, r, p1.Block().Hash(), tIndexX, false)
	testAddVote(t, tConsP, vote.VoteTypePrecommit, h, r, p1.Block().Hash(), tIndexY, false)

	checkHRS(t, tConsP, h, r, hrs.StepTypePropose)

	shouldPublishQueryProposal(t, tConsP, h, r)
	shouldPublishVote(t, tConsP, vote.VoteTypePrepare, crypto.UndefHash)

	// Now let's set the proposal
	tConsP.SetProposal(p1)

	shouldPublishVote(t, tConsP, vote.VoteTypePrepare, p1.Block().Hash())
	checkHRS(t, tConsP, h, r, hrs.StepTypePrepare)

	// We can't go to precommit stage, because we haven't prepared yet
	// But if we receive another vote we go to commit phase directly
	// Let's do it
	testAddVote(t, tConsP, vote.VoteTypePrecommit, h, r, p1.Block().Hash(), tIndexB, false)
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

	assert.True(t, tConsP.status.IsCommitted())
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

	h := 3
	p := makeProposal(t, h, 0) // tConsP should propose for this round

	tConsX.enterNewHeight()

	// tConsP is partitioned, so tConsX doesn't have the proposal
	shouldPublishVote(t, tConsX, vote.VoteTypePrepare, crypto.UndefHash)
	testAddVote(t, tConsX, vote.VoteTypePrepare, h, 0, crypto.UndefHash, tIndexY, false)
	testAddVote(t, tConsX, vote.VoteTypePrepare, h, 0, crypto.UndefHash, tIndexB, false)

	shouldPublishVote(t, tConsX, vote.VoteTypePrecommit, crypto.UndefHash)
	testAddVote(t, tConsX, vote.VoteTypePrecommit, h, 0, crypto.UndefHash, tIndexY, false)
	testAddVote(t, tConsX, vote.VoteTypePrecommit, h, 0, crypto.UndefHash, tIndexB, false)

	checkHRSWait(t, tConsX, h, 1, hrs.StepTypePrepare)

	// Now partition healed, but it's too late, We already moved to the next round
	tConsX.SetProposal(p)

	checkHRS(t, tConsX, h, 1, hrs.StepTypePrepare)
}

func TestSetProposalForNextRoundWithoutFinishingTheFirstRound(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)

	tConsX.enterNewHeight()

	// Byzantine node sends proposal for second round (his turn)
	b, err := tConsB.state.ProposeBlock(1)
	assert.NoError(t, err)
	p := proposal.NewProposal(2, 1, *b)
	tSigners[tIndexB].SignMsg(p)

	tConsX.SetProposal(p)
	// tConsX doesn't accept the proposal for next rounds
	assert.Nil(t, tConsX.RoundProposal(1))

	// But doesn't move to prepare phase
	checkHRS(t, tConsX, 2, 0, hrs.StepTypePropose)
}

func TestEnterPrepareAfterPrecommit(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)
	commitBlockForAllStates(t)
	commitBlockForAllStates(t)

	h := 4
	r := 0
	p := makeProposal(t, h, r)

	// tConsP is partitioned, so tConsX doesn't have the proposal
	tConsX.enterNewHeight()
	shouldPublishVote(t, tConsX, vote.VoteTypePrepare, crypto.UndefHash)

	testAddVote(t, tConsX, vote.VoteTypePrepare, h, r, crypto.UndefHash, tIndexY, false)
	testAddVote(t, tConsX, vote.VoteTypePrepare, h, r, crypto.UndefHash, tIndexB, false)
	checkHRS(t, tConsX, h, r, hrs.StepTypePrecommit)

	shouldPublishVote(t, tConsX, vote.VoteTypePrecommit, crypto.UndefHash)

	testAddVote(t, tConsX, vote.VoteTypePrecommit, h, r, crypto.UndefHash, tIndexY, false)
	testAddVote(t, tConsX, vote.VoteTypePrecommit, h, r, crypto.GenerateTestHash(), tIndexB, false)

	// Now partition healed
	tConsX.SetProposal(p)
	tConsX.enterPrepare(0)
	shouldPublishVote(t, tConsX, vote.VoteTypePrepare, p.Block().Hash())

	testAddVote(t, tConsX, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexY, false)
	testAddVote(t, tConsX, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexP, false)

	shouldPublishVote(t, tConsX, vote.VoteTypePrecommit, p.Block().Hash())
}

func TestProposeIvalidArgs(t *testing.T) {
	setup(t)

	tConsP.hrs = hrs.NewHRS(1, 0, hrs.StepTypeNewHeight)
	// Invalid args for propose phase
	tConsP.enterPropose(1)
	checkHRS(t, tConsP, 1, 0, hrs.StepTypeNewHeight)
}

func TestCreateProposal(t *testing.T) {
	setup(t)

	tConsX.enterNewHeight()
	tConsY.enterNewHeight()

	tConsX.createProposal(1, 0)
	assert.NotNil(t, tConsX.RoundProposal(0))

	tConsY.createProposal(1, 0)
	assert.Nil(t, tConsY.RoundProposal(0))
}
