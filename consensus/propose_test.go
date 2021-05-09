package consensus

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/consensus/proposal"
	"github.com/zarbchain/zarb-go/consensus/vote"
)

func TestProposeBlock(t *testing.T) {
	setup(t)
	testEnterNewHeight(tConsX)
	shouldPublishProposal(t, tConsX, 1, 0)
}

func TestSetProposalInvalidProposer(t *testing.T) {
	setup(t)

	testEnterNewHeight(tConsY)
	assert.Nil(t, tConsY.RoundProposal(0))

	addr := tSigners[tIndexB].Address()
	b, _ := block.GenerateTestBlock(&addr, nil)
	p := proposal.NewProposal(1, 0, b)

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
	p := proposal.NewProposal(1, 2, invBlock)
	tSigners[tIndexB].SignMsg(p)

	testEnterNewHeight(tConsP)
	testEnterNextRound(tConsP)
	testEnterNextRound(tConsP)

	tConsP.SetProposal(p)
	assert.Nil(t, tConsP.RoundProposal(2))
}

func TestSetProposalInvalidHeight(t *testing.T) {
	setup(t)

	a := tSigners[tIndexB].Address()
	invBlock, _ := block.GenerateTestBlock(&a, nil)
	p := proposal.NewProposal(2, 0, invBlock)
	tSigners[tIndexB].SignMsg(p)

	testEnterNewHeight(tConsY)
	tConsY.SetProposal(p)
	assert.Nil(t, tConsY.RoundProposal(2))
}

func TestConsensusSetProposalAfterCommit(t *testing.T) {
	setup(t)

	p := makeProposal(t, 1, 0)

	testEnterNewHeight(tConsP)
	commitBlockForAllStates(t)
	tConsP.SetProposal(p)
	assert.Nil(t, tConsP.RoundProposal(0))
}

func TestNetworkLagging(t *testing.T) {
	setup(t)

	testEnterNewHeight(tConsP)

	h := 1
	r := 0
	p := makeProposal(t, h, r)
	// We don't receive proposal on time
	// tConsP.SetProposal(p)

	testAddVote(tConsP, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexX)
	testAddVote(tConsP, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexY)
	shouldPublishQueryProposal(t, tConsP, h, r)

	// Proposal receives now
	tConsP.SetProposal(p)

	shouldPublishVote(t, tConsP, vote.VoteTypePrepare, p.Block().Hash())
	shouldPublishVote(t, tConsP, vote.VoteTypePrecommit, p.Block().Hash())
}

func TestProposalNextRound(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)

	testEnterNewHeight(tConsX)

	// Byzantine node sends proposal for second round (his turn)
	b, err := tConsB.state.ProposeBlock(1)
	assert.NoError(t, err)
	p := proposal.NewProposal(2, 1, b)
	tSigners[tIndexB].SignMsg(p)

	tConsX.SetProposal(p)

	// tConsX doesn't accept the proposal for next rounds
	assert.NotNil(t, tConsX.RoundProposal(1))
}
