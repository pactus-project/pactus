package amount_test

import (
	"strings"
	"testing"

	"github.com/pactus-project/pactus/types/amount"
	"github.com/stretchr/testify/require"
)

// An amount comes from the users and from the API: parsing any string never panics,
// and an amount that is parsed is written back to a string that gives the same amount.
func FuzzAmountFromString(f *testing.F) {
	for _, seed := range []string{
		"", "0", "1", "1.5", "0.000000001", "42000000", "-1", "1 PAC", "abc", "1.0000000001",
		"7000000.001", "41999999.999999999", "999999999.999999999", "-33006791.947779278",
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, text string) {
		amt, err := amount.FromString(text)
		if err != nil {
			return
		}

		again, err := amount.FromString(amt.Format(amount.WithUnit(amount.UnitPAC)))
		require.NoError(t, err)
		require.Equal(t, amt, again)

		// The delimiters are for the display only: without them, it is the same number.
		withDelimiters := amt.Format(amount.WithDelimiters())
		require.Equal(t, amt.Format(), strings.ReplaceAll(withDelimiters, ",", ""))
	})
}
