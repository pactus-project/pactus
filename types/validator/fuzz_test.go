package validator_test

import (
	"testing"

	"github.com/pactus-project/pactus/types/validator"
	"github.com/stretchr/testify/require"
)

// A validator that decodes is encoded back to the same bytes: it has one encoding.
func FuzzValidatorRoundTrip(f *testing.F) {
	// The seeds are in testdata/fuzz/FuzzValidatorRoundTrip.
	f.Add([]byte{})

	f.Fuzz(func(t *testing.T, data []byte) {
		val, err := validator.FromBytes(data)
		if err != nil {
			return
		}

		encoded, err := val.Bytes()
		require.NoError(t, err)
		require.Equal(t, data, encoded)
	})
}
