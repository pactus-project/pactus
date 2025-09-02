package consensusv2

import (
	"testing"

	"github.com/pactus-project/pactus/types/proposal"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestProposePublishProposal(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consX)
	p := td.shouldPublishProposal(t, td.consX, 1, 0)
	assert.Equal(t, td.consX.valKey.Address(), p.Block().Header().ProposerAddress())
}

func TestProposeInvalidProposer(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consY)
	assert.Nil(t, td.consY.Proposal())

	addr := td.consB.valKey.Address()
	blk, _ := td.GenerateTestBlock(1, testsuite.BlockWithProposer(addr))
	invalidProp := proposal.NewProposal(1, 0, blk)

	td.consY.SetProposal(invalidProp)
	assert.Nil(t, td.consY.Proposal())

	td.HelperSignProposal(td.consB.valKey, invalidProp)
	td.consY.SetProposal(invalidProp)
	assert.Nil(t, td.consY.Proposal())
}

func TestProposeInvalidBlock(t *testing.T) {
	td := setup(t)

	addr := td.consB.valKey.Address()
	blk, _ := td.GenerateTestBlock(1, testsuite.BlockWithProposer(addr))
	invProp := proposal.NewProposal(1, 2, blk)
	td.HelperSignProposal(td.consB.valKey, invProp)

	td.enterNewHeight(td.consP)
	td.enterNextRound(td.consP)
	td.enterNextRound(td.consP)

	td.consP.SetProposal(invProp)
	assert.Nil(t, td.consP.Proposal())
}

func TestProposeInvalidHeight(t *testing.T) {
	td := setup(t)

	addr := td.consB.valKey.Address()
	blk, _ := td.GenerateTestBlock(2, testsuite.BlockWithProposer(addr))
	invProp := proposal.NewProposal(2, 0, blk)
	td.HelperSignProposal(td.consB.valKey, invProp)

	td.enterNewHeight(td.consY)
	td.consY.SetProposal(invProp)
	assert.Nil(t, td.consY.Proposal())
}

func TestProposeNetworkLagging(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consP)

	height := uint32(1)
	round := int16(0)
	prop := td.makeProposal(t, height, round)

	// consP doesn't have the proposal, but it has received prepared votes from other peers
	td.addPrecommitVote(td.consP, prop.Block().Hash(), height, round, tIndexX)
	td.addPrecommitVote(td.consP, prop.Block().Hash(), height, round, tIndexY)

	td.queryProposalTimeout(td.consP)
	td.shouldPublishQueryProposal(t, td.consP, height, round)

	// Proposal is received now
	td.consP.SetProposal(prop)

	td.shouldPublishVote(t, td.consP, vote.VoteTypePrecommit, prop.Block().Hash())
}

func TestProposeNextRound(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t)

	td.enterNewHeight(td.consX)

	// Byzantine node sends proposal for the second round (his turn) even before the first round is started
	b, err := td.consB.bcState.ProposeBlock(td.consB.valKey, td.consB.rewardAddr)
	assert.NoError(t, err)
	p := proposal.NewProposal(2, 1, b)
	td.HelperSignProposal(td.consB.valKey, p)

	td.consX.SetProposal(p)

	// consX accepts his proposal, but doesn't move to the next round
	assert.NotNil(t, td.consX.log.RoundProposal(1))
	assert.Nil(t, td.consX.Proposal())
	assert.Equal(t, uint32(2), td.consX.height)
	assert.Equal(t, int16(0), td.consX.round)
}
