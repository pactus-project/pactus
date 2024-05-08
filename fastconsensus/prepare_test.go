package fastconsensus

import (
	"testing"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/stretchr/testify/assert"
)

func TestChangeProposerTimeout(t *testing.T) {
	td := setup(t)

	td.enterNewHeight(td.consP)
	td.changeProposerTimeout(td.consP)

	td.shouldPublishVote(t, td.consP, vote.VoteTypeCPPreVote, hash.UndefHash)
}

func TestQueryProposal(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t)
	h := uint32(2)

	td.enterNewHeight(td.consP)
	td.enterNextRound(td.consP)
	td.queryProposalTimeout(td.consP)

	td.shouldPublishQueryProposal(t, td.consP, h)
	td.shouldNotPublish(t, td.consP, message.TypeQueryVotes)
}

func TestQueryVotes(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t)
	td.commitBlockForAllStates(t)
	h := uint32(3)
	r := int16(1)

	td.enterNewHeight(td.consP)
	td.enterNextRound(td.consP)

	// consP is the proposer for this round, but there are not enough votes.
	td.queryProposalTimeout(td.consP)
	td.shouldPublishProposal(t, td.consP, h, r)
	td.shouldPublishQueryVote(t, td.consP, h, r)
	td.shouldNotPublish(t, td.consP, message.TypeQueryProposal)
}

func TestGoToChangeProposerFromPrepare(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t)

	td.enterNewHeight(td.consP)

	td.addCPPreVote(td.consP, hash.UndefHash, 2, 0, vote.CPValueYes, &vote.JustInitYes{}, tIndexX)
	td.addCPPreVote(td.consP, hash.UndefHash, 2, 0, vote.CPValueYes, &vote.JustInitYes{}, tIndexY)

	// should move to the change proposer phase, even if it has the proposal and
	// its timer has not expired, if it has received 1/3 of the change-proposer votes.
	prop := td.makeProposal(t, 2, 0)
	td.consP.SetProposal(prop)
	td.shouldPublishVote(t, td.consP, vote.VoteTypeCPPreVote, hash.UndefHash)
}

func TestByzantineProposal(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t)
	td.commitBlockForAllStates(t)
	h := uint32(3)
	r := int16(0)
	prop := td.makeProposal(t, h, r)
	propBlockHash := prop.Block().Hash()

	td.enterNewHeight(td.consP)

	td.addPrepareVote(td.consP, propBlockHash, h, r, tIndexX)
	td.addPrepareVote(td.consP, propBlockHash, h, r, tIndexY)
	td.addPrepareVote(td.consP, propBlockHash, h, r, tIndexB)
	td.addPrepareVote(td.consP, propBlockHash, h, r, tIndexM)
	td.addPrepareVote(td.consP, propBlockHash, h, r, tIndexN)

	assert.Nil(t, td.consP.Proposal())
	td.shouldPublishQueryProposal(t, td.consP, h)

	// Byzantine node sends second proposal to Partitioned node.
	trx := tx.NewTransferTx(h, td.consX.rewardAddr,
		td.RandAccAddress(), 1000, 1000, "invalid proposal")
	td.HelperSignTransaction(td.consX.valKey.PrivateKey(), trx)
	assert.NoError(t, td.txPool.AppendTx(trx))
	byzProp := td.makeProposal(t, h, r)
	assert.NotEqual(t, prop.Hash(), byzProp.Hash())

	td.consP.SetProposal(byzProp)
	assert.Nil(t, td.consP.Proposal())
	td.shouldPublishQueryProposal(t, td.consP, h)
	td.checkHeightRound(t, td.consP, h, r)
}
