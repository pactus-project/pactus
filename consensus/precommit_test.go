package consensus

import (
	"testing"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/stretchr/testify/assert"
)

func TestPrecommitQueryProposal(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t)
	height := types.Height(2)
	round := types.Round(0)

	td.enterNewHeight(td.consP)

	prop := td.makeProposal(t, height, round)

	td.addPrepareVote(td.consP, prop.Block().Hash(), height, round, tIndexX)
	td.addPrepareVote(td.consP, prop.Block().Hash(), height, round, tIndexY)
	td.addPrepareVote(td.consP, prop.Block().Hash(), height, round, tIndexB)

	td.addPrecommitVote(td.consP, prop.Block().Hash(), height, round, tIndexX)
	td.addPrecommitVote(td.consP, prop.Block().Hash(), height, round, tIndexY)
	td.addPrecommitVote(td.consP, prop.Block().Hash(), height, round, tIndexB)

	td.shouldPublishQueryProposal(t, td.consP, height)
}

func TestPrecommitDuplicatedProposal(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t)
	height := types.Height(2)
	round := types.Round(0)

	prop1 := td.makeProposal(t, height, round)
	prop2 := td.makeProposal(t, height, round, td.RandAccAddress())
	assert.NotEqual(t, prop1.Hash(), prop2.Hash())

	td.enterNewHeight(td.consP)

	// Byzantine node sends second proposal to Partitioned node
	// in prepare step
	td.consP.SetProposal(prop2)
	assert.NotNil(t, td.consP.Proposal())

	td.addPrepareVote(td.consP, prop1.Block().Hash(), height, round, tIndexX)
	td.addPrepareVote(td.consP, prop1.Block().Hash(), height, round, tIndexY)
	td.addPrepareVote(td.consP, prop1.Block().Hash(), height, round, tIndexB)

	assert.Nil(t, td.consP.Proposal())
	td.shouldPublishQueryProposal(t, td.consP, height)

	// Byzantine node sends second proposal to Partitioned node,
	// in precommit step
	td.consP.SetProposal(prop2)
	assert.Nil(t, td.consP.Proposal())
	td.shouldPublishQueryProposal(t, td.consP, height)

	td.consP.SetProposal(prop1)
	assert.NotNil(t, td.consP.Proposal())
}

func TestGoToChangeProposerFromPrecommit(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t)
	height := types.Height(2)
	round := types.Round(0)

	td.enterNewHeight(td.consP)
	blockHash := td.RandHash()

	td.addPrepareVote(td.consP, blockHash, height, round, tIndexX)
	td.addPrepareVote(td.consP, blockHash, height, round, tIndexY)
	td.addPrepareVote(td.consP, blockHash, height, round, tIndexB)

	td.addCPPreVote(td.consP, hash.UndefHash, height, round, vote.CPValueYes, &vote.JustInitYes{}, tIndexX)
	td.addCPPreVote(td.consP, hash.UndefHash, height, round, vote.CPValueYes, &vote.JustInitYes{}, tIndexY)

	td.shouldPublishVote(t, td.consP, vote.VoteTypeCPPreVote, blockHash)
}
