package crypto

import (
	"encoding/hex"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddressJSONMarshaling(t *testing.T) {
	addr1 := GenerateTestAddress()

	js, err := json.Marshal(addr1)
	assert.NoError(t, err)
	assert.Contains(t, string(js), addr1.String())
	assert.Contains(t, addr1.String(), addr1.Fingerprint())
}

func TestAddressFromString(t *testing.T) {
	addr1 := GenerateTestAddress()
	prv2, err := AddressFromString(addr1.String())
	assert.NoError(t, err)
	assert.True(t, addr1.EqualsTo(prv2))

	_, err = AddressFromString("")
	assert.Error(t, err)

	_, err = AddressFromString("inv")
	assert.Error(t, err)
}

func TestAddressEmpty(t *testing.T) {
	_, err := AddressFromBytes(nil)
	assert.Error(t, err)

	_, err = AddressFromBytes([]byte{1})
	assert.Error(t, err)
}

func TestTreasuryAddress(t *testing.T) {
	assert.Equal(t, TreasuryAddress.String(), treasuryAddressString)
	expected, err := AddressFromString(treasuryAddressString)
	assert.NoError(t, err)
	assert.Equal(t, TreasuryAddress.RawBytes(), expected.RawBytes())
}

func TestInvalidBech32(t *testing.T) {
	// ok
	addr, err := AddressFromString("zc17mka0cw484es5whq638xkm89msgzczmrmf7p27")
	assert.NoError(t, err)
	assert.Equal(t, addr.Fingerprint(), "zc17mka0cw48")

	// Invalid hrp
	_, err = AddressFromString("sc17mka0cw484es5whq638xkm89msgzczmr75t2kv")
	assert.Error(t, err)

	// Invalid type
	_, err = AddressFromString("zc27mka0cw484es5whq638xkm89msgzczmrpd86dv")
	assert.Error(t, err)

	// Invalid checksum
	_, err = AddressFromString("zc17mka0cw484es5whq638xkm89msgzczmrwy64dz")
	assert.Error(t, err)
}

func TestAddressSanityCheck(t *testing.T) {
	inv, _ := hex.DecodeString(strings.Repeat("ff", AddressSize))
	addr1, err := AddressFromBytes(inv)
	assert.NoError(t, err)
	assert.Error(t, addr1.SanityCheck(), "invalid address type")
	addr1[0] = 1
	assert.NoError(t, addr1.SanityCheck())
}
