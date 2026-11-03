package tx_test

import (
	"bytes"
	"testing"

	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTxDecodeIsTransactional(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	trxA := ts.GenerateTestTransferTx()
	trxB := ts.GenerateTestBondTx()
	dataA, err := trxA.Bytes()
	require.NoError(t, err)

	t.Run("a failed decoding leaves the transaction untouched", func(t *testing.T) {
		trx := ts.GenerateTestTransferTx()
		before := trx.ID()

		require.Error(t, trx.Decode(bytes.NewReader(dataA[:len(dataA)-5])))
		assert.Equal(t, before, trx.ID())
		assert.NotNil(t, trx.Payload())
	})

	t.Run("the caches of the previous content are dropped", func(t *testing.T) {
		trx := trxB
		require.NoError(t, trx.BasicCheck())
		_ = trx.ID()

		require.NoError(t, trx.Decode(bytes.NewReader(dataA)))
		assert.Equal(t, trxA.ID(), trx.ID())
		assert.Equal(t, trxA.Payload().Type(), trx.Payload().Type())
	})
}

func TestTxWithPublicKeyButNoSignature(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	trx := ts.GenerateTestTransferTx()
	trx.SetSignature(nil)
	require.NotNil(t, trx.PublicKey())

	data, err := trx.Bytes()
	require.NoError(t, err)
	assert.Len(t, data, trx.SerializeSize(), "the size matches the encoding")

	// The following transaction of a block must not be read from the bytes of the public key.
	buf := bytes.NewBuffer(data)
	trxNext := ts.GenerateTestTransferTx()
	require.NoError(t, trxNext.Encode(buf))

	first, second := new(tx.Tx), new(tx.Tx)
	require.NoError(t, first.Decode(buf))
	require.NoError(t, second.Decode(buf))
	assert.Equal(t, trxNext.ID(), second.ID())
	assert.Zero(t, buf.Len())
}
