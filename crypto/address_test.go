package crypto_test

import (
	"bytes"
	"encoding/hex"
	"io"
	"strings"
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/bech32m"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestTreasuryAddressType(t *testing.T) {
	treasury := crypto.TreasuryAddress

	assert.False(t, treasury.IsValidatorAddress())
	assert.True(t, treasury.IsAccountAddress())
	assert.True(t, treasury.IsTreasuryAddress())
}

func TestAddressType(t *testing.T) {
	tests := []struct {
		address   string
		account   bool
		validator bool
	}{
		{address: "pc1pjneygutecly9gtandrdt8j36v8g4fl42k4y5xp", account: false, validator: true},
		{address: "pc1z0m0vw8sjfgv7f2zgq2hfxutg8rwn7gpffhe8tf", account: true, validator: false},
		{address: "pc1rcx9x55nfme5juwdgxd2ksjdcmhvmvkrygmxpa3", account: true, validator: false},
	}

	for _, tt := range tests {
		addr, _ := crypto.AddressFromString(tt.address)

		assert.Equal(t, tt.account, addr.IsAccountAddress())
		assert.Equal(t, tt.validator, addr.IsValidatorAddress())
	}
}

func TestShortString(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	addr1 := ts.RandAccAddress()
	assert.Contains(t, addr1.String(), addr1.LogString())
}

func TestFromString(t *testing.T) {
	tests := []struct {
		encoded  string
		err      error
		bytes    []byte
		addrType crypto.AddressType
	}{
		{
			"000000000000000000000000000000000000000000",
			nil,
			[]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			crypto.AddressTypeTreasury,
		},
		{
			"00",
			bech32m.InvalidLengthError(2),
			nil,
			0,
		},
		{
			"",
			bech32m.InvalidLengthError(0),
			nil,
			0,
		},
		{
			"not_proper_encoded",
			bech32m.InvalidSeparatorIndexError(-1),
			nil,
			0,
		},
		{
			"pc1ioiooi",
			bech32m.NonCharsetCharError(105),
			nil,
			0,
		},
		{
			"pc19p72rf",
			bech32m.InvalidLengthError(0),
			nil,
			0,
		},
		{
			"qc1z0hrct7eflrpw4ccrttxzs4qud2axex4dh8zz75",
			crypto.InvalidHRPError("qc"),
			nil,
			0,
		},
		{
			"pc1p0hrct7eflrpw4ccrttxzs4qud2axex4dg8xaf5",
			bech32m.InvalidChecksumError{Expected: "cdzdfr", Actual: "g8xaf5"},
			nil,
			0,
		},
		{
			"pc1p0hrct7eflrpw4ccrttxzs4qud2axexs2dhdk8",
			crypto.InvalidLengthError(20),
			nil,
			0,
		},
		{
			"pc1y0hrct7eflrpw4ccrttxzs4qud2axex4dksmred",
			crypto.InvalidAddressTypeError(4),
			nil,
			0,
		},
		{
			"PC1P0HRCT7EFLRPW4CCRTTXZS4QUD2AXEX4DCDZDFR", // UPPERCASE
			nil,
			[]byte{
				0x01, 0x7d, 0xc7, 0x85, 0xfb, 0x29, 0xf8, 0xc2, 0xea, 0xe3,
				0x03, 0x5a, 0xcc, 0x28, 0x54, 0x1c, 0x6a, 0xba, 0x6c, 0x9a, 0xad,
			},
			crypto.AddressTypeValidator,
		},
		{
			"pc1p0hrct7eflrpw4ccrttxzs4qud2axex4dcdzdfr",
			nil,
			[]byte{
				0x01, 0x7d, 0xc7, 0x85, 0xfb, 0x29, 0xf8, 0xc2, 0xea, 0xe3,
				0x03, 0x5a, 0xcc, 0x28, 0x54, 0x1c, 0x6a, 0xba, 0x6c, 0x9a, 0xad,
			},
			crypto.AddressTypeValidator,
		},
		{
			"pc1zzqkzzu4vyddss052as6c37qrdcfptegquw826x",
			nil,
			[]byte{
				0x02, 0x10, 0x2c, 0x21, 0x72, 0xac, 0x23, 0x5b, 0x08, 0x3e, 0x8a,
				0xec, 0x35, 0x88, 0xf8, 0x03, 0x6e, 0x12, 0x15, 0xe5, 0x00,
			},
			crypto.AddressTypeBLSAccount,
		},
		{
			"pc1rspm7ps49gar9ft5g0tkl6lhxs8ygeakq87quh3",
			nil,
			[]byte{
				0x03, 0x80, 0x77, 0xe0, 0xc2, 0xa5, 0x47, 0x46, 0x54, 0xae,
				0x88, 0x7a, 0xed, 0xfd, 0x7e, 0xe6, 0x81, 0xc8, 0x8c, 0xf6, 0xc0,
			},
			crypto.AddressTypeEd25519Account,
		},
	}
	for no, tt := range tests {
		addr, err := crypto.AddressFromString(tt.encoded)
		if tt.err == nil {
			assert.NoError(t, err, "test %v: unexpected error", no)
			assert.Equal(t, tt.bytes, addr.Bytes(), "test %v: invalid result", no)
			assert.Equal(t, strings.ToLower(tt.encoded), addr.String(), "test %v: invalid encode", no)
			assert.Equal(t, tt.addrType, addr.Type(), "test %v: invalid type", no)
		} else {
			assert.ErrorIs(t, err, tt.err, "test %v: invalid error", no)
		}
	}
}

func TestAddressDecoding(t *testing.T) {
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
			"040000000000000000000000000000000000000000",
			crypto.InvalidAddressTypeError(4),
		},
		{
			0,
			"04000102030405060708090a0b0c0d0e0f0001020304",
			crypto.InvalidAddressTypeError(4),
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
		{
			21,
			"03000102030405060708090a0b0c0d0e0f00010203",
			nil,
		},
	}
	for no, tt := range tests {
		data, _ := hex.DecodeString(tt.hex)
		buf := bytes.NewBuffer(data)
		addr := new(crypto.Address)

		err := addr.Decode(buf)
		if tt.err != nil {
			assert.ErrorIs(t, err, tt.err, "test %v: error not matched", no)
			assert.Equal(t, tt.size, addr.SerializeSize(), "test %v invalid size", no)
		} else {
			assert.NoError(t, err, "test %v expected no error", no)
			assert.Equal(t, tt.size, addr.SerializeSize(), "test %v invalid size", no)

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
