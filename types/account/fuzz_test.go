package account_test

import (
	"testing"

	"github.com/pactus-project/pactus/types/account"
	"github.com/stretchr/testify/require"
)

// An account that decodes is encoded back to the same bytes: it has one encoding.
func FuzzAccountRoundTrip(f *testing.F) {
	// The seeds are in testdata/fuzz/FuzzAccountRoundTrip.
	f.Add([]byte{})

	f.Fuzz(func(t *testing.T, data []byte) {
		acc, err := account.FromBytes(data)
		if err != nil {
			return
		}

		encoded, err := acc.Bytes()
		require.NoError(t, err)
		require.Equal(t, data, encoded)
	})
}
