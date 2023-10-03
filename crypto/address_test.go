package crypto_test

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/bech32m"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestAddressKeyType(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	pub, _ := ts.RandBLSKeyPair()
	accAddr := pub.AccountAddress()
	valAddr := pub.ValidatorAddress()
	treasury := crypto.TreasuryAddress

	assert.True(t, accAddr.IsAccountAddress())
	assert.False(t, accAddr.IsValidatorAddress())
	assert.False(t, accAddr.IsTreasuryAddress())
	assert.False(t, valAddr.IsAccountAddress())
	assert.True(t, valAddr.IsValidatorAddress())
	assert.False(t, treasury.IsValidatorAddress())
	assert.True(t, treasury.IsAccountAddress())
	assert.True(t, treasury.IsTreasuryAddress())
	assert.NotEqual(t, accAddr, valAddr)
}

func TestString(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	a, _ := crypto.AddressFromString("pc1p0hrct7eflrpw4ccrttxzs4qud2axex4dcdzdfr")
	fmt.Println(a.String())

	addr1 := ts.RandAccAddress()
	assert.Contains(t, addr1.String(), addr1.ShortString())
}

func TestToString(t *testing.T) {
	tests := []struct {
		encoded string
		err     error
		result  *crypto.Address
	}{
		{
			"000000000000000000000000000000000000000000",
			nil,
			&crypto.Address{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			"",
			bech32m.InvalidLengthError(0),
			nil,
		},
		{
			"not_proper_encoded",
			bech32m.InvalidSeparatorIndexError(-1),
			nil,
		},
		{
			"pc1ioiooi",
			bech32m.NonCharsetCharError(105),
			nil,
		},
		{
			"pc19p72rf",
			bech32m.InvalidLengthError(0),
			nil,
		},
		{
			"qc1z0hrct7eflrpw4ccrttxzs4qud2axex4dh8zz75",
			crypto.InvalidHRPError("qc"),
			nil,
		},
		{
			"pc1p0hrct7eflrpw4ccrttxzs4qud2axex4dg8xaf5",
			bech32m.InvalidChecksumError{Expected: "cdzdfr", Actual: "g8xaf5"},
			nil,
		},
		{
			"pc1p0hrct7eflrpw4ccrttxzs4qud2axexs2dhdk8",
			crypto.InvalidLengthError(20),
			nil,
		},
		{
			"pc1r0hrct7eflrpw4ccrttxzs4qud2axex4dwc9mn4",
			crypto.InvalidAddressTypeError(3),
			nil,
		},
		{
			"PC1P0HRCT7EFLRPW4CCRTTXZS4QUD2AXEX4DCDZDFR", // UPPERCASE
			nil,
			&crypto.Address{
				0x1, 0x7d, 0xc7, 0x85, 0xfb, 0x29, 0xf8, 0xc2, 0xea, 0xe3,
				0x3, 0x5a, 0xcc, 0x28, 0x54, 0x1c, 0x6a, 0xba, 0x6c, 0x9a, 0xad,
			},
		},
		{
			"pc1p0hrct7eflrpw4ccrttxzs4qud2axex4dcdzdfr",
			nil,
			&crypto.Address{
				0x1, 0x7d, 0xc7, 0x85, 0xfb, 0x29, 0xf8, 0xc2, 0xea, 0xe3,
				0x3, 0x5a, 0xcc, 0x28, 0x54, 0x1c, 0x6a, 0xba, 0x6c, 0x9a, 0xad,
			},
		},
	}
	for no, test := range tests {
		addr, err := crypto.AddressFromString(test.encoded)
		if test.err == nil {
			assert.NoError(t, err, "test %v: unexpected error", no)
			assert.Equal(t, addr, *test.result, "test %v: invalid result", no)
			assert.Equal(t, addr.String(), strings.ToLower(test.encoded), "test %v: invalid encode", no)
		} else {
			assert.ErrorIs(t, err, test.err, "test %v: invalid error", no)
		}
	}
}

func TestAddressEncoding(t *testing.T) {
	tests := []struct {
		size int
		hex  string
		err  error
	}{
		{
			1,
			"00",
			nil,
		},
		{
			0,
			"030000000000000000000000000000000000000000",
			crypto.InvalidAddressTypeError(3),
		},
		{
			0,
			"03000102030405060708090a0b0c0d0e0f0001020304",
			crypto.InvalidAddressTypeError(3),
		},
		{
			21,
			"0100",
			io.ErrUnexpectedEOF,
		},
		{
			21,
			"01000102030405060708090a0b0c0d0e0f000102",
			io.ErrUnexpectedEOF,
		},
		{
			21,
			"01000102030405060708090a0b0c0d0e0f00010203",
			nil,
		},
		{
			21,
			"02000102030405060708090a0b0c0d0e0f00010203",
			nil,
		},
	}
	for no, test := range tests {
		data, _ := hex.DecodeString(test.hex)
		r := bytes.NewBuffer(data)
		addr := new(crypto.Address)

		err := addr.Decode(r)
		if test.err != nil {
			assert.ErrorIs(t, test.err, err, "test %v: error not matched", no)
			assert.Equal(t, addr.SerializeSize(), test.size, "test %v invalid size", no)
		} else {
			assert.NoError(t, err, "test %v expected no error", no)
			assert.Equal(t, addr.SerializeSize(), test.size, "test %v invalid size", no)

			length := addr.SerializeSize()
			for i := 0; i < length; i++ {
				w := util.NewFixedWriter(i)
				assert.Error(t, addr.Encode(w), "encode test %v failed", i)
			}
			w := util.NewFixedWriter(length)
			assert.NoError(t, addr.Encode(w))
			assert.Equal(t, data, w.Bytes())
		}
	}
}
