package crypto

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/pactus-project/pactus/util/errors"
	"github.com/stretchr/testify/assert"
)

func TestFingerprint(t *testing.T) {
	addr1 := GenerateTestAddress()
	assert.Contains(t, addr1.String(), addr1.Fingerprint())
}

func TestToString(t *testing.T) {
	randomAddr := GenerateTestAddress()
	tests := []struct {
		name      string
		encoded   string
		decodable bool
		result    *Address
	}{
		{
			"treasury address",
			"000000000000000000000000000000000000000000",
			true,
			&Address{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			"empty address",
			"",
			false,
			nil,
		},
		{
			"invalid character",
			"inv",
			false,
			nil,
		},
		{
			"no type",
			"pc19p72rf",
			false,
			nil,
		},
		{
			"invalid hrp",
			"qc1z0hrct7eflrpw4ccrttxzs4qud2axex4dh8zz75",
			false,
			nil,
		},
		{
			"invalid checksum",
			"pc1p0hrct7eflrpw4ccrttxzs4qud2axex4dg8xaf5",
			false,
			nil,
		},
		{
			"invalid length",
			"pc1p0hrct7eflrpw4ccrttxzs4qud2axexs2dhdk8",
			false,
			nil,
		},
		{
			"invalid type",
			"pc1z0hrct7eflrpw4ccrttxzs4qud2axex4d9xjs77",
			false,
			nil,
		},
		{
			"valid address in uppercase format",
			"PC1P0HRCT7EFLRPW4CCRTTXZS4QUD2AXEX4DCDZDFR",
			true,
			&Address{0x1, 0x7d, 0xc7, 0x85, 0xfb, 0x29, 0xf8, 0xc2, 0xea, 0xe3,
				0x3, 0x5a, 0xcc, 0x28, 0x54, 0x1c, 0x6a, 0xba, 0x6c, 0x9a, 0xad},
		},
		{
			"valid address",
			"pc1p0hrct7eflrpw4ccrttxzs4qud2axex4dcdzdfr",
			true,
			&Address{0x1, 0x7d, 0xc7, 0x85, 0xfb, 0x29, 0xf8, 0xc2, 0xea, 0xe3,
				0x3, 0x5a, 0xcc, 0x28, 0x54, 0x1c, 0x6a, 0xba, 0x6c, 0x9a, 0xad},
		},
		{
			"random address",
			randomAddr.String(),
			true,
			&randomAddr,
		},
	}
	for _, test := range tests {
		addr, err := AddressFromString(test.encoded)
		if test.decodable {
			assert.NoError(t, err,
				"test '%v' failed. unexpected error", test.name)
			assert.Equal(t, addr, *test.result,
				"test '%v' failed.. unexpected result", test.name)
			assert.Equal(t, addr.String(), strings.ToLower(test.encoded),
				"test '%v' failed.. encoded failed", test.name)
		} else {
			assert.Error(t, err,
				"test '%v failed.. should return error", test.name)
			assert.Equal(t, errors.Code(err), errors.ErrInvalidAddress)
		}
	}
}

func TestAddressSanityCheck(t *testing.T) {
	tests := []struct {
		name    string
		hex     string
		invalid bool
	}{
		{
			"type zero with non zero data",
			"00ffffffffffffffffffffffffffffffffffffffff",
			true,
		},
		{
			"invalid type",
			"020000000000000000000000000000000000000000",
			true,
		},
		{
			"treasury address",
			"000000000000000000000000000000000000000000",
			false,
		},
		{
			"valid address",
			"010000000000000000000000000000000000000000",
			false,
		},
	}
	for _, test := range tests {
		data, _ := hex.DecodeString(test.hex)
		addr := Address{}
		copy(addr[:], data)

		if test.invalid {
			assert.Error(t, addr.SanityCheck(),
				"test '%v'. expected error", test.name)
		} else {
			assert.NoError(t, addr.SanityCheck(),
				"test '%v'. unexpected error", test.name)
		}
	}
}
