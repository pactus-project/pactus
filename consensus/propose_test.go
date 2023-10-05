package consensus

import (
	"testing"

	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/stretchr/testify/assert"
)

func TestProposeBlock(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consX)
	p := td.shouldPublishProposal(t, td.consX, 1, 0)
	assert.Equal(t, td.consX.valKey.Address(), p.Block().Header().ProposerAddress())
}

func TestSetProposalInvalidProposer(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consY)
	assert.Nil(t, td.consY.RoundProposal(0))

	addr := td.valKeys[tIndexB].Address()
	blk, _ := td.GenerateTestBlockWithProposer(1, addr)
	invalidProp := proposal.NewProposal(1, 0, blk)

	td.consY.SetProposal(invalidProp)
	assert.Nil(t, td.consY.RoundProposal(0))

	td.HelperSignProposal(td.valKeys[tIndexB], invalidProp)
	td.consY.SetProposal(invalidProp)
	assert.Nil(t, td.consY.RoundProposal(0))
}

func TestSetProposalInvalidBlock(t *testing.T) {
	td := setup(t)

	addr := td.valKeys[tIndexB].Address()
	blk, _ := td.GenerateTestBlockWithProposer(1, addr)
	invProp := proposal.NewProposal(1, 2, blk)
	td.HelperSignProposal(td.valKeys[tIndexB], invProp)

	td.enterNewHeight(td.consP)
	td.enterNextRound(td.consP)
	td.enterNextRound(td.consP)

	td.consP.SetProposal(invProp)
	assert.Nil(t, td.consP.RoundProposal(2))
}

func TestSetProposalInvalidHeight(t *testing.T) {
	td := setup(t)

	addr := td.valKeys[tIndexB].Address()
	blk, _ := td.GenerateTestBlockWithProposer(2, addr)
	invProp := proposal.NewProposal(2, 0, blk)
	td.HelperSignProposal(td.valKeys[tIndexB], invProp)

	td.enterNewHeight(td.consY)
	td.consY.SetProposal(invProp)
	assert.Nil(t, td.consY.RoundProposal(2))
}

func TestSetProposalAfterCommit(t *testing.T) {
	td := setup(t)

	p0 := td.makeProposal(t, 1, 0)
	p1 := td.makeProposal(t, 1, 1)

	td.enterNewHeight(td.consP)
	td.commitBlockForAllStates(t)

	td.consP.SetProposal(p0)
	assert.NotNil(t, td.consP.RoundProposal(0))

	td.consP.SetProposal(p1)
	assert.Nil(t, td.consP.RoundProposal(1))
}

func TestNetworkLagging(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consP)

	h := uint32(1)
	r := int16(0)
	p := td.makeProposal(t, h, r)

	// consP doesn't have the proposal, but it has received prepared votes from other peers
	td.addPrepareVote(td.consP, p.Block().Hash(), h, r, tIndexX)
	td.addPrepareVote(td.consP, p.Block().Hash(), h, r, tIndexY)

	td.queryProposalTimeout(td.consP)
	td.shouldPublishQueryProposal(t, td.consP, h, r)

	// Proposal is received now
	td.consP.SetProposal(p)

	td.shouldPublishVote(t, td.consP, vote.VoteTypePrepare, p.Block().Hash())
	td.shouldPublishVote(t, td.consP, vote.VoteTypePrecommit, p.Block().Hash())
}

func TestProposalNextRound(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t)

	td.enterNewHeight(td.consX)

	// Byzantine node sends proposal for the second round (his turn) even before the first round is started
	b, err := td.consB.state.ProposeBlock(td.consB.valKey, td.consB.rewardAddr, 1)
	assert.NoError(t, err)
	p := proposal.NewProposal(2, 1, b)
	td.HelperSignProposal(td.valKeys[tIndexB], p)

	td.consX.SetProposal(p)

	// consX accepts his proposal, but doesn't move to the next round
	assert.NotNil(t, td.consX.RoundProposal(1))
	assert.Equal(t, td.consX.height, uint32(2))
	assert.Equal(t, td.consX.round, int16(0))
}
