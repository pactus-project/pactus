package sync

import (
	"strings"
	"testing"

	"github.com/pactus-project/pactus/network"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// prepareBlocks (sync/sync.go) used to collect up to
// config.BlockPerMessage (60) raw block byte slices with no cumulative size
// check at all. The default network.MaxStreamMessageSize is 8 MiB, and the
// receiver's stream reader is wrapped in util.LimitReaderClose at exactly
// that size (network/stream.go). 60 consecutive near-full blocks (up to
// state/param.MaxTransactionsPerBlock = 1000 transactions each) together
// exceed 8 MiB, so every honest peer used to truncate the same response and
// the requester's CBOR decode failed: a new node stalled on that exact
// height range, since prepareBlocks always returned the same oversized set
// for it.
//
// This test builds several real, near-full blocks (1 subsidy + 999
// maximal-memo transfer transactions each, matching MaxTransactionsPerBlock)
// through the real block encoding path, commits them, and calls
// prepareBlocks directly: it must never hand back a response whose total
// byte size exceeds MaxStreamMessageSize, resuming from a lower height
// instead.
func TestBlocksResponseCanExceedMaxStreamMessageSize(t *testing.T) {
	td := setup(t, nil)

	const txPerBlock = 1000            // state/param.Params.MaxTransactionsPerBlock
	const nearFullBlocks = 30          // one near-full block is ~268 KB; 30 of them exceed the 7 MiB budget.
	maxMemo := strings.Repeat("A", 64) // the maximum memo length a transaction may carry.
	maxStreamMessageSize := network.DefaultConfig().MaxStreamMessageSize

	fromHeight := td.state.FakeHeight + 1
	for i := 0; i < nearFullBlocks; i++ {
		txs := block.NewTxs()
		txs.Append(td.GenerateTestSubsidyTx())
		for j := 0; j < txPerBlock-1; j++ {
			_, signer := td.RandBLSKeyPair()
			sender := signer.PublicKey().AccountAddress()
			receiver := td.RandAccAddress()
			trx := tx.NewTransferTx(td.RandHeight(), sender, receiver, td.RandAmount(), td.RandFee(), tx.WithMemo(maxMemo))
			td.HelperSignTransaction(signer, trx)
			txs.Append(trx)
		}

		height := td.state.FakeHeight + 1
		blk, cert := td.GenerateTestBlock(height, testsuite.BlockWithTransactions(txs))
		td.state.AddTestBlock(blk, cert)
		td.state.FakeHeight = height
	}

	blocksData := td.sync.prepareBlocks(fromHeight, nearFullBlocks)

	require.NotEmpty(t, blocksData, "prepareBlocks must still make progress, even if it cannot fit every requested block")

	totalBytes := 0
	for _, data := range blocksData {
		totalBytes += len(data)
	}

	t.Logf("requested %d near-full blocks (%d tx each); prepareBlocks returned %d, %d bytes (limit: %d)",
		nearFullBlocks, txPerBlock, len(blocksData), totalBytes, maxStreamMessageSize)

	assert.LessOrEqual(t, totalBytes, maxStreamMessageSize,
		"prepareBlocks must never return more bytes than a single stream message can carry")
	assert.Less(t, len(blocksData), nearFullBlocks,
		"with near-full blocks, prepareBlocks must stop before BlockPerMessage once the byte budget is reached")
}
