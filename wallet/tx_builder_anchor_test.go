package wallet

import (
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/require"
)

func TestTxBuilderAnchor(t *testing.T) {
	sender := crypto.TreasuryAddress
	receiver := crypto.TreasuryAddress
	builder := &txBuilder{
		typ:      payload.TypeAnchor,
		lockTime: 1,
		sender:   &sender,
		receiver: &receiver,
	}

	_, err := builder.build()
	require.Error(t, err)

	ts := testsuite.NewTestSuite(t)
	pub, _ := ts.RandEd25519KeyPair()
	from := pub.AccountAddress()
	other := crypto.TreasuryAddress
	require.NoError(t, OptionMemo("note")(builder))
	require.NoError(t, OptionFee("1")(builder))
	builder.sender = &from
	builder.receiver = &other
	builder.lockTime = 9
	builder.anchorRoot = bytesOf(0x11, 32)
	builder.anchorURI = "uri"
	builder.anchorType = 0
	builder.amount = 3
	builder.typ = payload.TypeAnchor

	trx, err := builder.build()
	require.NoError(t, err)
	pld := trx.Payload().(*payload.AnchorPayload)
	require.Equal(t, from, pld.From)
	require.NotEqual(t, other, pld.From)
	require.Equal(t, payload.AnchorActionSet, pld.Action)
	require.Equal(t, bytesOf(0x11, 32), pld.RootHash)
	require.Equal(t, "uri", pld.ManifestURI)
	require.Equal(t, uint8(0), pld.AnchorType)
	require.Equal(t, amount.Amount(3), pld.Deposit)
	require.Equal(t, "note", trx.Memo())
	require.Equal(t, amount.Amount(1e9), trx.Fee())

	builder.anchorRoot = bytesOf(0x22, 31)
	_, err = builder.build()
	require.Error(t, err)

	builder.anchorAction = payload.AnchorActionDelete
	builder.anchorRoot = nil
	builder.anchorURI = "ignored"
	builder.anchorType = 0xFF
	builder.amount = 99
	trx, err = builder.build()
	require.NoError(t, err)
	pld = trx.Payload().(*payload.AnchorPayload)
	require.Equal(t, payload.AnchorActionDelete, pld.Action)
	require.Empty(t, pld.RootHash)
	require.Equal(t, amount.Amount(0), pld.Deposit)
	raw, err := trx.Bytes()
	require.NoError(t, err)
	require.NotContains(t, string(raw), "ignored")
}

func bytesOf(fill byte, n int) []byte {
	out := make([]byte, n)
	for i := range out {
		out[i] = fill
	}

	return out
}
