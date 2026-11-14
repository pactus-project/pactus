package tx_test

import (
	"testing"

	"github.com/pactus-project/pactus/types/tx"
	"github.com/stretchr/testify/require"
)

// A transaction that decodes is encoded back to the same bytes: one transaction has one encoding.
func FuzzTxRoundTrip(f *testing.F) {
	// The seeds are in testdata/fuzz/FuzzTxRoundTrip: one of each payload type.
	f.Add([]byte{})

	f.Fuzz(func(t *testing.T, data []byte) {
		trx, err := tx.FromBytes(data)
		if err != nil {
			return
		}

		encoded, err := trx.Bytes()
		require.NoError(t, err)
		require.Equal(t, data, encoded)
	})
}
