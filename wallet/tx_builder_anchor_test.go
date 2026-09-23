package wallet

import (
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types/tx/payload"
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

	builder.typ = payload.TypeTransfer
	trx, err := builder.build()
	require.NoError(t, err)
	require.Equal(t, payload.TypeTransfer, trx.Payload().Type())
}
