package crypto_test

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/util/errors"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestAddressKeyEqualsTo(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	addr1 := ts.RandomAddress()
	addr2 := ts.RandomAddress()

	assert.True(t, addr1.EqualsTo(addr1))
	assert.False(t, addr1.EqualsTo(addr2))
	assert.Equal(t, addr1, addr1)
	assert.NotEqual(t, addr1, addr2)
}

func TestString(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	addr1 := ts.RandomAddress()
	assert.Contains(t, addr1.String(), addr1.ShortString())
}

func TestToString(t *testing.T) {
	tests := []struct {
		errMsg    string
		encoded   string
		decodable bool
		result    *crypto.Address
	}{
		{
			"",
			"000000000000000000000000000000000000000000",
			true,
			&crypto.Address{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			"invalid bech32 string length 0",
			"",
			false,
			nil,
		},
		{
			"invalid separator index -1",
			"not_proper_encoded",
			false,
			nil,
		},
		{
			"invalid character not part of charset: 105",
			"pc1ioiooi",
			false,
			nil,
		},
		{
			"invalid bech32 string length 0",
			"pc19p72rf",
			false,
			nil,
		},
		{
			"invalid hrp: qc",
			"qc1z0hrct7eflrpw4ccrttxzs4qud2axex4dh8zz75",
			false,
			nil,
		},
		{
			"invalid checksum (expected cdzdfr got g8xaf5)",
			"pc1p0hrct7eflrpw4ccrttxzs4qud2axex4dg8xaf5",
			false,
			nil,
		},
		{
			"address should be 21 bytes, but it is 20 bytes",
			"pc1p0hrct7eflrpw4ccrttxzs4qud2axexs2dhdk8",
			false,
			nil,
		},
		{
			"invalid address key type: 2",
			"pc1z0hrct7eflrpw4ccrttxzs4qud2axex4d9xjs77",
			false,
			nil,
		},
		{
			"",
			"PC1P0HRCT7EFLRPW4CCRTTXZS4QUD2AXEX4DCDZDFR", // UPPERCASE
			true,
			&crypto.Address{
				0x1, 0x7d, 0xc7, 0x85, 0xfb, 0x29, 0xf8, 0xc2, 0xea, 0xe3,
				0x3, 0x5a, 0xcc, 0x28, 0x54, 0x1c, 0x6a, 0xba, 0x6c, 0x9a, 0xad,
			},
		},
		{
			"",
			"pc1p0hrct7eflrpw4ccrttxzs4qud2axex4dcdzdfr",
			true,
			&crypto.Address{
				0x1, 0x7d, 0xc7, 0x85, 0xfb, 0x29, 0xf8, 0xc2, 0xea, 0xe3,
				0x3, 0x5a, 0xcc, 0x28, 0x54, 0x1c, 0x6a, 0xba, 0x6c, 0x9a, 0xad,
			},
		},
	}
	for no, test := range tests {
		addr, err := crypto.AddressFromString(test.encoded)
		if test.decodable {
			assert.NoError(t, err, "test %v: unexpected error", no)
			assert.Equal(t, addr, *test.result, "test %v: invalid result", no)
			assert.Equal(t, addr.String(), strings.ToLower(test.encoded), "test %v: invalid encode", no)
		} else {
			assert.Equal(t, errors.Code(err), errors.ErrInvalidAddress, "test %v: invalid error code", no)
			assert.Contains(t, err.Error(), test.errMsg, "test %v: error not matched", no)
		}
	}
}

func TestAddressBasicCheck(t *testing.T) {
	tests := []struct {
		errMsg  string
		hex     string
		invalid bool
	}{
		{
			"invalid address data",
			"00ffffffffffffffffffffffffffffffffffffffff",
			true,
		},
		{
			"invalid address type",
			"020000000000000000000000000000000000000000",
			true,
		},
		{
			"",
			"000000000000000000000000000000000000000000",
			false,
		},
		{
			"",
			"010000000000000000000000000000000000000000",
			false,
		},
	}
	for no, test := range tests {
		data, _ := hex.DecodeString(test.hex)
		addr := crypto.Address{}
		copy(addr[:], data)

		err := addr.BasicCheck()
		if !test.invalid {
			assert.NoError(t, err, "test %v unexpected error", no)
		} else {
			assert.Error(t, err, "test %v expected error", no)
			assert.Contains(t, err.Error(), test.errMsg, "test %v: error not matched", no)
		}
	}
}
