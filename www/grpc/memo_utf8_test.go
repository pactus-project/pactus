package grpc_test

import (
	"testing"

	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/util/testsuite"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Tx.BasicCheck only checks the memo's length (types/tx/tx.go),
// never that it is valid UTF-8. The TransactionInfo.Memo proto field is a
// proto3 string, and the generated Go runtime (google.golang.org/protobuf)
// rejects marshaling a string field that is not valid UTF-8. A single paid
// transaction with a non-UTF-8 memo therefore passes BasicCheck, gets
// committed, and then makes GetTransaction fail for every caller trying to
// read it back, because transactionToProto copies trx.Memo() verbatim into
// the proto message before it is marshaled.
//
// This test expresses the desired behavior (GetTransaction succeeds, with a
// sanitized memo) and failed before the fix, because the server call itself fails.
func TestNonUTF8MemoBreaksGetTransaction(t *testing.T) {
	td := setup(t, nil)
	client := td.transactionClient(t)

	_, signer := td.RandBLSKeyPair()
	sender := signer.PublicKey().AccountAddress()
	receiver := td.RandAccAddress()

	badMemo := string([]byte{0xff, 0xfe, 0x80})
	trx := tx.NewTransferTx(td.RandHeight(), sender, receiver, td.RandAmount(), td.RandFee(), tx.WithMemo(badMemo))
	td.HelperSignTransaction(signer, trx)

	require.NoError(t, trx.BasicCheck(),
		"a transaction with a non-UTF-8 memo must still pass BasicCheck")

	// FakeState.CommittedTx (AnyTimes default stub) delegates to
	// FakeStore.Transaction, which scans FakeBlocks for a matching tx ID.
	blk, _ := td.GenerateTestBlock(1, testsuite.BlockWithTransactions(block.Txs{trx}))
	td.FakeState.FakeStore.FakeBlocks[1] = blk

	res, err := client.GetTransaction(t.Context(), &pactus.GetTransactionRequest{
		Id:        trx.ID().String(),
		Verbosity: pactus.TransactionVerbosity_TRANSACTION_VERBOSITY_INFO,
	})

	require.NoError(t, err, "a committed transaction with a non-UTF-8 memo must still be readable through GetTransaction")
	require.NotNil(t, res)
	assert.NotNil(t, res.Transaction)
}
