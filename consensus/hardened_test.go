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

// In the next rounds, an honest validator only sends a `yes` pre-vote for no block too.
// A `yes` pre-vote with a block hash and a soft or hard justification must be rejected like
// the one of the first round: the main-votes that are justified with it are rejected by the others.
func TestCPPreVoteYesHasNoBlockHashInNextRounds(t *testing.T) {
	newVote := func(td *testData, blockHash hash.Hash, just vote.Just) *vote.Vote {
		v := vote.NewCPPreVote(blockHash, 1, 0, 1, vote.CPValueYes, just, td.valKeys[tIndexB].Address())
		td.HelperSignVote(td.valKeys[tIndexB], v)

		return v
	}
	justs := map[string]vote.Just{
		"soft": &vote.JustPreVoteSoft{},
		"hard": &vote.JustPreVoteHard{},
	}

	for name, just := range justs {
		t.Run(name+", before the hardened rules", func(t *testing.T) {
			td := setup(t)
			td.stateP.StateParams.BlockVersion = protocol.ProtocolVersion5
			td.enterNewHeight(td.consP)

			require.NoError(t, td.consP.changeProposer.checkYesHasNoBlockHash(newVote(td, td.RandHash(), just)))
		})

		t.Run(name+", hardened rules", func(t *testing.T) {
			td := setup(t)
			td.stateP.StateParams.BlockVersion = protocol.ProtocolVersion6
			td.enterNewHeight(td.consP)

			blockHash := td.RandHash()
			err := td.consP.changeProposer.checkJust(newVote(td, blockHash, just))
			require.ErrorAs(t, err, &InvalidJustificationError{})
			require.ErrorContains(t, err, "invalid block hash: "+blockHash.String())

			require.NoError(t, td.consP.changeProposer.checkYesHasNoBlockHash(newVote(td, hash.UndefHash, just)))
		})
	}

	t.Run("a `no` pre-vote keeps its block hash", func(t *testing.T) {
		td := setup(t)
		td.stateP.StateParams.BlockVersion = protocol.ProtocolVersion6
		td.enterNewHeight(td.consP)

		v := vote.NewCPPreVote(td.RandHash(), 1, 0, 1, vote.CPValueNo, &vote.JustPreVoteSoft{},
			td.valKeys[tIndexB].Address())
		require.NoError(t, td.consP.changeProposer.checkYesHasNoBlockHash(v))
	})
}
