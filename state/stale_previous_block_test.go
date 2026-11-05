package state

import (
	"testing"

	"github.com/pactus-project/pactus/types/block"
	"github.com/stretchr/testify/require"
)

// A valid proposer builds a block on top of an older block of the chain: the
// previous hash and the previous certificate are the ones of the block at
// height tip-2 (same committee, so the certificate is valid for the validators
// of the last block), while the state root and the sortition seed are the
// current ones. ValidateBlock must reject a block that doesn't extend the tip.
func TestValidateBlockWithStalePrevBlock(t *testing.T) {
	td := setup(t)

	td.fakeTxPool.EXPECT().PrepareBlockTransactions().Return(block.Txs{}).AnyTimes()

	lastHeight := td.state.LastBlockHeight()
	oldHeight := lastHeight - 2
	oldHash := td.fakeStore.BlockHash(oldHeight)
	oldCert := td.fakeStore.FakeCertificates[oldHeight]
	require.NotNil(t, oldCert)
	require.NotEqual(t, td.state.LastBlockHash(), oldHash)

	valKey := td.proposerKey(t, 0)
	good, err := td.state.ProposeBlock(valKey, td.RandAccAddress())
	require.NoError(t, err)
	require.NoError(t, td.state.ValidateBlock(good, 0))

	stale := block.MakeBlock(good.Header().Version(), good.Header().Time(), good.Transactions(),
		oldHash, good.Header().StateRoot(), oldCert, good.Header().SortitionSeed(),
		good.Header().ProposerAddress())

	err = td.state.ValidateBlock(stale, 0)
	require.Error(t, err, "a block whose previous hash is not the last block hash must be rejected")
}
