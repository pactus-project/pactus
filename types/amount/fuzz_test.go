package amount_test

import (
	"testing"

	"github.com/pactus-project/pactus/types/amount"
	"github.com/stretchr/testify/require"
)

// An amount comes from the users and from the API: parsing any string never panics,
// and an amount that is parsed is written back to a string that gives the same amount.
func FuzzAmountFromString(f *testing.F) {
	for _, seed := range []string{"", "0", "1", "1.5", "0.000000001", "42000000", "-1", "1 PAC", "abc", "1.0000000001"} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, text string) {
		amt, err := amount.FromString(text)
		if err != nil {
			return
		}

		// An amount is parsed and formatted as a float, which is exact for the small amounts only.
		const exactLimit = 1_000_000 * amount.NanoPACPerPAC
		if amt > exactLimit || amt < -exactLimit {
			return
		}

		again, err := amount.FromString(amt.Format(amount.WithUnit(amount.UnitPAC)))
		require.NoError(t, err)
		require.Equal(t, amt, again)
	})
}
