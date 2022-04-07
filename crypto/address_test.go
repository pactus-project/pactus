package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/errors"
)

func TestFingerprint(t *testing.T) {
	addr1 := GenerateTestAddress()
	assert.Contains(t, addr1.String(), addr1.Fingerprint())
}

func TestTreasuryAddress(t *testing.T) {
	assert.Equal(t, TreasuryAddress.String(), treasuryAddressString)
	expected, err := AddressFromString(treasuryAddressString)
	assert.NoError(t, err)
	assert.Equal(t, TreasuryAddress.Bytes(), expected.Bytes())
}

func TestInvalidStrings(t *testing.T) {
	randomAddr := GenerateTestAddress()
	tests := []struct {
		name      string
		encoded   string
		decodable bool
		valid     bool
		result    *Address
	}{
		{
			"empty address",
			"",
			false, false,
			nil,
		},
		{
			"invalid character",
			"inv",
			false, false,
			nil,
		},
		{
			"invalid hrp",
			"sc1qxlq6p2aqdnqxcads964v062uvr3yv4gju4xjq38",
			false, false,
			nil,
		},
		{
			"invalid checksum",
			"zc1q97u0p0m98uv96hrqddvc2z5r34t5my6454m4svw",
			false, false,
			nil,
		},
		{
			"invalid length",
			"zc1q97u0p0m98uv96hrqddvc2z5r34t5my65pamyq",
			false, false,
			nil,
		},
		{
			"invalid type",
			"zc1qf7u0p0m98uv96hrqddvc2z5r34t5my645678lzn",
			true, false,
			&Address{0x2, 0x7d, 0xc7, 0x85, 0xfb, 0x29, 0xf8, 0xc2, 0xea, 0xe3, 0x3, 0x5a, 0xcc, 0x28, 0x54, 0x1c, 0x6a, 0xba, 0x6c, 0x9a, 0xad},
		},
		{
			"valid address in uppercase format",
			"ZC1Q97U0P0M98UV96HRQDDVC2Z5R34T5MY6454M4SVU",
			true, true,
			&Address{0x1, 0x7d, 0xc7, 0x85, 0xfb, 0x29, 0xf8, 0xc2, 0xea, 0xe3, 0x3, 0x5a, 0xcc, 0x28, 0x54, 0x1c, 0x6a, 0xba, 0x6c, 0x9a, 0xad},
		},
		{
			"valid address",
			"zc1q97u0p0m98uv96hrqddvc2z5r34t5my6454m4svu",
			true, true,
			&Address{0x1, 0x7d, 0xc7, 0x85, 0xfb, 0x29, 0xf8, 0xc2, 0xea, 0xe3, 0x3, 0x5a, 0xcc, 0x28, 0x54, 0x1c, 0x6a, 0xba, 0x6c, 0x9a, 0xad},
		},
		{
			"random address",
			randomAddr.String(),
			true, true,
			&randomAddr,
		},
	}
	for _, test := range tests {
		addr, err := AddressFromString(test.encoded)
		if test.decodable {
			assert.NoError(t, err, "test %v. unexpected error", test.name)
			assert.Equal(t, addr, *test.result, "test %v. unexpected result", test.name)
			assert.Equal(t, addr.SanityCheck() == nil, test.valid, "test %v. sanity check failed", test.name)

		} else {
			assert.Error(t, err, "test %v. should failed", test.name)
			assert.Equal(t, errors.Code(err), errors.ErrInvalidAddress)
		}
	}
}
