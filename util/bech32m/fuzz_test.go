package bech32m_test

import (
	"strings"
	"testing"

	"github.com/pactus-project/pactus/util/bech32m"
	"github.com/stretchr/testify/require"
)

// Decoding any string never panics, and a string that is decoded
// is encoded back to the same string, in lower case.
func FuzzDecode(f *testing.F) {
	for _, seed := range []string{
		"", "A1LQFN3A", "a1lqfn3a", "abcdef1l7aum6echk45nj3s0wdvt2fg8x9yrzpqzd3ryx",
		"split1checkupstagehandshakeupstreamerranterredcaperredlc445v", "?1v759aa", "1qzzfhee",
		"pc1p0hrct7eflrpw4ccrttxzs4qud2axex4dcdzdfr",
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, text string) {
		hrp, data, err := bech32m.Decode(text)
		if err != nil {
			return
		}

		encoded, err := bech32m.Encode(hrp, data)
		require.NoError(t, err)
		require.Equal(t, strings.ToLower(text), encoded)

		_, _, _ = bech32m.DecodeToBase256(text)
	})
}
