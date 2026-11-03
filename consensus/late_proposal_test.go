package consensus

import (
	"testing"

	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLateProposalIsNotAuthenticated checks that the "block already committed" path of
// SetProposal() (prop.Height() == LastBlockHeight) still verifies the signature of the proposer.
// That branch only compares the block hash, so anybody can make a validator store and
// precommit a proposal with an invalid signature, for any round >= the current one.
// v1 counterpart of consensusv2's test of the same name.
func TestLateProposalIsNotAuthenticated(t *testing.T) {
	td := setup(t)

	td.commitBlockForAllStates(t) // height 1
	td.enterNewHeight(td.consP)

	height := types.Height(2)
	blk, _ := td.commitBlockForAllStates(t) // height 2: P is late

	prop := td.makeProposal(t, height, 0, testsuite.ProposalWithBlock(blk))
	prop.SetSignature(td.valKeys[tIndexX].Sign(prop.SignBytes())) // X is not the proposer of round 0
	require.Error(t, prop.Verify(td.stateP.Proposer(0).PublicKey()))

	td.consP.SetProposal(prop)

	assert.Nil(t, td.consP.Proposal(), "a proposal with an invalid signature must be ignored")
}
