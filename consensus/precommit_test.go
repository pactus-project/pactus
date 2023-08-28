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

	td.enterNewHeight(td.consP)

	p := td.makeProposal(t, 2, 0)

	td.addPrepareVote(td.consP, p.Block().Hash(), 2, 0, tIndexX)
	td.addPrepareVote(td.consP, p.Block().Hash(), 2, 0, tIndexY)
	td.addPrepareVote(td.consP, p.Block().Hash(), 2, 0, tIndexB)

	td.addPrecommitVote(td.consP, p.Block().Hash(), 2, 0, tIndexX)
	td.addPrecommitVote(td.consP, p.Block().Hash(), 2, 0, tIndexY)
	td.addPrecommitVote(td.consP, p.Block().Hash(), 2, 0, tIndexB)

	td.shouldPublishQueryProposal(t, td.consP, 2, 0)
}

func TestPrecommitInvalidProposal(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t)

	p1 := td.makeProposal(t, 2, 0)
	trx := tx.NewTransferTx(hash.UndefHash.Stamp(), 1, td.signers[0].Address(),
		td.signers[1].Address(), 1000, 1000, "invalid proposal")
	td.signers[0].SignMsg(trx)
	assert.NoError(t, td.txPool.AppendTx(trx))
	p2 := td.makeProposal(t, 2, 0)
	assert.NotEqual(t, p1.Hash(), p2.Hash())

	td.enterNewHeight(td.consP)

	td.addPrepareVote(td.consP, p1.Block().Hash(), 2, 0, tIndexX)
	td.addPrepareVote(td.consP, p1.Block().Hash(), 2, 0, tIndexY)
	td.addPrepareVote(td.consP, p1.Block().Hash(), 2, 0, tIndexB)

	td.consP.SetProposal(p2)
	assert.Nil(t, td.consP.RoundProposal(0))

	td.consP.SetProposal(p1)
	assert.NotNil(t, td.consP.RoundProposal(0))
}

func TestGoToChangeProposerFromPrecommit(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t)

	td.enterNewHeight(td.consP)
	h := td.RandHash()

	td.addPrepareVote(td.consP, h, 2, 0, tIndexX)
	td.addPrepareVote(td.consP, h, 2, 0, tIndexY)
	td.addPrepareVote(td.consP, h, 2, 0, tIndexB)

	td.addCPPreVote(td.consP, hash.UndefHash, 2, 0, 0, vote.CPValueOne, &vote.JustInitOne{}, tIndexX)
	td.addCPPreVote(td.consP, hash.UndefHash, 2, 0, 0, vote.CPValueOne, &vote.JustInitOne{}, tIndexY)

	td.shouldPublishVote(t, td.consP, vote.VoteTypeCPPreVote, h)
}
