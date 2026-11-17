package crypto_test

import (
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/stretchr/testify/require"
)

// An address comes from the users and from the API: parsing any string never panics,
// and an address that is parsed is written back to a string that gives the same address.
func FuzzAddressFromString(f *testing.F) {
	f.Add("")
	f.Add("000000000000000000000000000000000000000000")
	f.Add("pc1p0hrct7eflrpw4ccrttxzs4qud2axex4dcdzdfr")
	f.Add("pc1zzqkzzu4vyddss052as6c37qrdcfptegquw826x")
	f.Add("pc1rspm7ps49gar9ft5g0tkl6lhxs8ygeakq87quh3")
	f.Add("PC1P0HRCT7EFLRPW4CCRTTXZS4QUD2AXEX4DCDZDFR")

	f.Fuzz(func(t *testing.T, text string) {
		addr, err := crypto.AddressFromString(text)
		if err != nil {
			return
		}

		_ = addr.ShortString()
		_ = addr.LogString()
		_ = addr.Type()

		again, err := crypto.AddressFromString(addr.String())
		require.NoError(t, err)
		require.Equal(t, addr, again)
	})
}
