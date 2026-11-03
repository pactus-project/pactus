package consensus

import (
	"testing"

	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// A late proposal, for a block that is committed already, is accepted only when it is signed
// by the proposer of that block, and for that block.
func TestLateProposalChecks(t *testing.T) {
	setupLate := func(t *testing.T) (*testData, *block.Block, int) {
		t.Helper()

		td := setup(t)
		td.commitBlockForAllStates(t) // height 1
		td.enterNewHeight(td.consP)

		round := types.Round(0)
		proposerIndex := td.stateX.ProposerIndex(round)
		blk, cert := td.GenerateTestBlock(td.stateX.FakeHeight+1,
			testsuite.BlockWithProposer(td.valKeys[proposerIndex].Address()))
		_ = td.stateX.CommitBlock(blk, cert)
		_ = td.stateY.CommitBlock(blk, cert)
		_ = td.stateB.CommitBlock(blk, cert)
		_ = td.stateP.CommitBlock(blk, cert)

		return td, blk, proposerIndex
	}

	t.Run("signed by another validator", func(t *testing.T) {
		td, blk, proposerIndex := setupLate(t)
		other := (proposerIndex + 1) % len(td.valKeys)

		prop := td.makeProposal(t, 2, 0,
			testsuite.ProposalWithBlock(blk), testsuite.ProposalWithKey(td.valKeys[other]))
		td.consP.SetProposal(prop)

		assert.Nil(t, td.consP.Proposal())
	})

	t.Run("for another block", func(t *testing.T) {
		td, _, proposerIndex := setupLate(t)

		otherBlk, _ := td.GenerateTestBlock(2,
			testsuite.BlockWithProposer(td.valKeys[proposerIndex].Address()))
		prop := td.makeProposal(t, 2, 0,
			testsuite.ProposalWithBlock(otherBlk), testsuite.ProposalWithKey(td.valKeys[proposerIndex]))
		td.consP.SetProposal(prop)

		assert.Nil(t, td.consP.Proposal())
	})
}

// A proposal is signed by its proposer, but its block holds a transaction with a forged signature:
// the block hash doesn't cover the signatures of the transactions, and the proposal is ignored.
func TestProposalWithForgedTransaction(t *testing.T) {
	td := setup(t)
	td.enterNewHeight(td.consP)

	trx := td.GenerateTestTransferTx()
	trx.SetSignature(td.RandEd25519Signature())
	blk, _ := td.GenerateTestBlock(td.consP.height,
		testsuite.BlockWithTransactions(block.Txs{td.GenerateTestSubsidyTx(), trx}))
	require.Error(t, blk.BasicCheck())

	prop := td.makeProposal(t, td.consP.height, td.consP.round, testsuite.ProposalWithBlock(blk))
	td.consP.SetProposal(prop)

	assert.Nil(t, td.consP.Proposal())
}
