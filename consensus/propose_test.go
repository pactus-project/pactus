package consensus

import (
	"testing"

	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/stretchr/testify/assert"
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
	b := block.GenerateTestBlock(&addr, nil)
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
	invBlock := block.GenerateTestBlock(&a, nil)
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
	invBlock := block.GenerateTestBlock(&a, nil)
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

	h := uint32(1)
	r := int16(0)
	p := makeProposal(t, h, r)

	// tConsP doesn't have the proposal, but it has received prepared votes from other peers
	testAddVote(tConsP, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexX)
	testAddVote(tConsP, vote.VoteTypePrepare, h, r, p.Block().Hash(), tIndexY)
	shouldPublishQueryProposal(t, tConsP, h, r)

	// Proposal is received now
	tConsP.SetProposal(p)

	shouldPublishVote(t, tConsP, vote.VoteTypePrepare, p.Block().Hash())
	shouldPublishVote(t, tConsP, vote.VoteTypePrecommit, p.Block().Hash())
}

func TestProposalNextRound(t *testing.T) {
	setup(t)

	commitBlockForAllStates(t)

	testEnterNewHeight(tConsX)

	// Byzantine node sends proposal for the second round (his turn) even before the first round is started
	b, err := tConsB.state.ProposeBlock(tConsB.signer, tConsB.rewardAddr, 1)
	assert.NoError(t, err)
	p := proposal.NewProposal(2, 1, b)
	tSigners[tIndexB].SignMsg(p)

	tConsX.SetProposal(p)

	// tConsX accepts his proposal, but doesn't move to the next round
	assert.NotNil(t, tConsX.RoundProposal(1))
	assert.Equal(t, tConsX.height, uint32(2))
	assert.Equal(t, tConsX.round, int16(0))
}
