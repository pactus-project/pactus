package consensus

import (
	"testing"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/protocol"
	"github.com/pactus-project/pactus/types/vote"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The first round of the change-proposer phase has two kinds of pre-votes: `yes` for no block,
// and `no` for the proposed block. A `yes` pre-vote with a block hash is counted in the quorum of `yes`,
// but it can't be in the certificate, which is for no block.
func TestCPPreVoteYesHasNoBlockHash(t *testing.T) {
	newVote := func(td *testData, blockHash hash.Hash) *vote.Vote {
		v := vote.NewCPPreVote(blockHash, 1, 0, 0, vote.CPValueYes, &vote.JustInitYes{}, td.valKeys[tIndexB].Address())
		td.HelperSignVote(td.valKeys[tIndexB], v)

		return v
	}

	t.Run("before the hardened rules, a block hash is accepted", func(t *testing.T) {
		td := setup(t)
		td.stateP.StateParams.BlockVersion = protocol.ProtocolVersion5
		td.enterNewHeight(td.consP)

		require.NoError(t, td.consP.changeProposer.checkJust(newVote(td, td.RandHash())), "legacy rules")
	})

	t.Run("hardened rules", func(t *testing.T) {
		td := setup(t)
		td.stateP.StateParams.BlockVersion = protocol.ProtocolVersion6
		td.enterNewHeight(td.consP)

		require.NoError(t, td.consP.changeProposer.checkJust(newVote(td, hash.UndefHash)))

		blockHash := td.RandHash()
		err := td.consP.changeProposer.checkJust(newVote(td, blockHash))
		require.ErrorAs(t, err, &InvalidJustificationError{})
		require.ErrorContains(t, err, blockHash.String())

		// And the vote is not added to the log.
		v := newVote(td, td.RandHash())
		td.consP.AddVote(v)
		assert.False(t, td.consP.log.HasVote(v.Hash()))
	})
}
