package executor

import (
	"bytes"
	"testing"

	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/stretchr/testify/require"
)

func TestAnchorRejectedByExecutor(t *testing.T) {
	td := setup(t)

	acc, addr := td.addTestAccount(t)
	before, err := acc.Bytes()
	require.NoError(t, err)
	require.Len(t, before, 12)
	require.False(t, acc.HasAnchor())

	pub, prv := td.RandBLSKeyPair()
	from := pub.AccountAddress()
	fee := td.RandFee()
	setTx := tx.NewAnchorTx(td.sbx.CurrentHeight(), from, payload.AnchorActionSet,
		bytes.Repeat([]byte{0x11}, 32), "abc", 0x04, 1, fee)
	td.HelperSignTransaction(prv, setTx)
	_, err = MakeExecutor(setTx, td.sbx)
	require.ErrorIs(t, err, InvalidPayloadTypeError{PayloadType: payload.TypeAnchor})

	delTx := tx.NewAnchorTx(td.sbx.CurrentHeight(), from, payload.AnchorActionDelete, nil, "", 0, 0, fee)
	td.HelperSignTransaction(prv, delTx)
	_, err = MakeExecutor(delTx, td.sbx)
	require.ErrorIs(t, err, InvalidPayloadTypeError{PayloadType: payload.TypeAnchor})

	after, err := acc.Bytes()
	require.NoError(t, err)
	require.Equal(t, before, after)
	require.False(t, acc.HasAnchor())

	receiver := td.RandAccAddress()
	amt, transferFee := td.randAmountFee(acc.Balance())
	transfer := tx.NewTransferTx(td.sbx.CurrentHeight(), addr, receiver, amt, transferFee)
	td.check(t, transfer, true, nil)
	td.execute(t, transfer)
	require.Equal(t, amt, td.sbx.Account(receiver).Balance())

	valPub, _ := td.RandBLSKeyPair()
	bondTx := tx.NewBondTx(td.sbx.CurrentHeight(), addr, valPub.ValidatorAddress(), valPub, 1, 1)
	_, err = MakeExecutor(bondTx, td.sbx)
	require.NoError(t, err)
}
