package state

import (
	"testing"

	"github.com/pactus-project/gopkg/pipeline"
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// A node that restarts must compute the same committee and the same proposers
// as the node that kept running, whatever the round of the last committed block was.
// Otherwise, it would reject valid proposals or propose at the wrong round.
func TestRestartKeepsCommitteeAndProposers(t *testing.T) {
	td := setup(t)
	td.fakeTxPool.EXPECT().PrepareBlockTransactions().Return(block.Txs{}).AnyTimes()
	// The validators of the test can win the sortition once the bonding period is over.
	td.fakeTxPool.EXPECT().AppendTxAndBroadcast(gomock.Any()).Return(nil).AnyTimes()

	identity := func(vals []*validator.Validator) [][2]any {
		ids := make([][2]any, len(vals))
		for i, val := range vals {
			ids[i] = [2]any{val.Number(), val.Address()}
		}

		return ids
	}

	rounds := []types.Round{0, 1, 3, 2, 0, 5, 7, 1, 0, 4, 2, 9}
	for no, round := range rounds {
		blk, cert := td.makeBlockAndCertificate(t, round)
		require.NoError(t, td.state.CommitBlock(blk, cert))

		restarted, err := LoadOrNewState(t.Context(), td.state.genDoc, td.state.valKeys,
			td.fakeStore, td.fakeTxPool, pipeline.New[any](t.Context()))
		require.NoError(t, err)
		newState, _ := restarted.(*state)

		for next := types.Round(0); next < 10; next++ {
			require.Equal(t, td.state.Proposer(next).Address(), newState.Proposer(next).Address(),
				"block %d (committed at round %d): different proposer at round %d", no, round, next)
		}

		require.Equal(t, td.state.committee.Committers(), newState.committee.Committers(),
			"block %d (committed at round %d)", no, round)
		assert.Equal(t, td.state.LastBlockHash(), newState.LastBlockHash())
		assert.Equal(t, td.state.LastBlockTime(), newState.LastBlockTime())
		assert.Equal(t, td.state.stateRoot(), newState.stateRoot())
		assert.Equal(t, td.state.totalPower, newState.totalPower)
		assert.Equal(t, td.state.Params().BlockVersion, newState.Params().BlockVersion)
		assert.Equal(t, td.state.lastInfo.SortitionSeed(), newState.lastInfo.SortitionSeed())
		assert.Equal(t, td.state.lastInfo.Certificate().Hash(), newState.LastCertificate().Hash())
		assert.Equal(t, identity(td.state.lastInfo.Validators()), identity(newState.lastInfo.Validators()))
	}
}
