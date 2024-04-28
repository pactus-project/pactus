package consensus

import (
	"testing"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/stretchr/testify/assert"
)

func TestPrecommitQueryProposal(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t)
	h := uint32(2)
	r := int16(0)

	td.enterNewHeight(td.consP)

	p := td.makeProposal(t, h, r)

	td.addPrepareVote(td.consP, p.Block().Hash(), h, r, tIndexX)
	td.addPrepareVote(td.consP, p.Block().Hash(), h, r, tIndexY)
	td.addPrepareVote(td.consP, p.Block().Hash(), h, r, tIndexB)

	td.addPrecommitVote(td.consP, p.Block().Hash(), h, r, tIndexX)
	td.addPrecommitVote(td.consP, p.Block().Hash(), h, r, tIndexY)
	td.addPrecommitVote(td.consP, p.Block().Hash(), h, r, tIndexB)

	td.shouldPublishQueryProposal(t, td.consP, h)
}

func TestPrecommitDuplicatedProposal(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t)
	h := uint32(2)
	r := int16(0)

	p1 := td.makeProposal(t, h, r)
	trx := tx.NewTransferTx(h, td.consX.rewardAddr,
		td.RandAccAddress(), 1000, 1000, "invalid proposal")
	td.HelperSignTransaction(td.consX.valKey.PrivateKey(), trx)

	assert.NoError(t, td.txPool.AppendTx(trx))
	p2 := td.makeProposal(t, h, r)
	assert.NotEqual(t, p1.Hash(), p2.Hash())

	td.enterNewHeight(td.consP)

	// Byzantine node sends second proposal to Partitioned node
	// in prepare step
	td.consP.SetProposal(p2)
	assert.NotNil(t, td.consP.Proposal())

	td.addPrepareVote(td.consP, p1.Block().Hash(), h, r, tIndexX)
	td.addPrepareVote(td.consP, p1.Block().Hash(), h, r, tIndexY)
	td.addPrepareVote(td.consP, p1.Block().Hash(), h, r, tIndexB)

	assert.Nil(t, td.consP.Proposal())
	td.shouldPublishQueryProposal(t, td.consP, h)

	// Byzantine node sends second proposal to Partitioned node,
	// in precommit step
	td.consP.SetProposal(p2)
	assert.Nil(t, td.consP.Proposal())
	td.shouldPublishQueryProposal(t, td.consP, h)

	td.consP.SetProposal(p1)
	assert.NotNil(t, td.consP.Proposal())
}

func TestGoToChangeProposerFromPrecommit(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t)
	h := uint32(2)
	r := int16(0)

	td.enterNewHeight(td.consP)
	blockHash := td.RandHash()

	td.addPrepareVote(td.consP, blockHash, h, r, tIndexX)
	td.addPrepareVote(td.consP, blockHash, h, r, tIndexY)
	td.addPrepareVote(td.consP, blockHash, h, r, tIndexB)

	td.addCPPreVote(td.consP, hash.UndefHash, h, r, 0, vote.CPValueYes, &vote.JustInitYes{}, tIndexX)
	td.addCPPreVote(td.consP, hash.UndefHash, h, r, 0, vote.CPValueYes, &vote.JustInitYes{}, tIndexY)

	td.shouldPublishVote(t, td.consP, vote.VoteTypeCPPreVote, blockHash)
}
