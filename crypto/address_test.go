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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		{address: "pc1pcs9ezsmn6n7fc8jxzxks4n2lyw4ltzsdc9v8qn", account: false, validator: true},
		{address: "pc1zcs9ezsmn6n7fc8jxzxks4n2lyw4ltzsd9wu6hw", account: true, validator: false},
		{address: "pc1rj65g93q7lpdq0366vst22l7va9d26j3l2vr0em", account: true, validator: false},
		{address: "pc1y90qakls8jlz9hyvdcsqsj0yj2lrqz26vqu7l0z", account: true, validator: false},
	}

	for _, tt := range tests {
		addr, _ := crypto.AddressFromString(tt.address)

		assert.Equal(t, tt.account, addr.IsAccountAddress())
		assert.Equal(t, tt.validator, addr.IsValidatorAddress())
		assert.False(t, addr.IsTreasuryAddress())
	}
}

func TestFromString(t *testing.T) {
	tests := []struct {
		encoded  string
		err      error
		hex      string
		addrType crypto.AddressType
	}{
		{
			"000000000000000000000000000000000000000000",
			nil,
			"000000000000000000000000000000000000000000",
			crypto.AddressTypeTreasury,
		},
		{
			"00",
			bech32m.InvalidLengthError(2),
			"",
			0,
		},
		{
			"",
			bech32m.InvalidLengthError(0),
			"",
			0,
		},
		{
			"not_proper_encoded",
			bech32m.InvalidSeparatorIndexError(-1),
			"",
			0,
		},
		{
			"pc1ioiooi",
			bech32m.NonCharsetCharError(105),
			"",
			0,
		},
		{
			"pc19p72rf",
			bech32m.InvalidLengthError(0),
			"",
			0,
		},
		{
			"qc1z0hrct7eflrpw4ccrttxzs4qud2axex4dh8zz75",
			crypto.InvalidHRPError("qc"),
			"",
			0,
		},
		{
			"pc1p0hrct7eflrpw4ccrttxzs4qud2axexs2dhdk8",
			crypto.InvalidLengthError(20),
			"",
			0,
		},
		{
			"pc1p0hrct7eflrpw4ccrttxzs4qud2axex4dg8xaf5",
			bech32m.InvalidChecksumError{Expected: "cdzdfr", Actual: "g8xaf5"},
			"",
			0,
		},
		{
			"pc19q5qqzqsrqszsvpcgpy9qkrqdpc8sqqgz24gkaq",
			crypto.InvalidAddressTypeError(5),
			"",
			0,
		},
		{
			"PC1PCS9EZSMN6N7FC8JXZXKS4N2LYW4LTZSDC9V8QN", // UPPERCASE
			nil,
			"01c40b914373d4fc9c1e4611ad0acd5f23abf58a0d",
			crypto.AddressTypeValidator,
		},
		{
			"pc1pcs9ezsmn6n7fc8jxzxks4n2lyw4ltzsdc9v8qn",
			nil,
			"01c40b914373d4fc9c1e4611ad0acd5f23abf58a0d",
			crypto.AddressTypeValidator,
		},
		{
			"pc1zcs9ezsmn6n7fc8jxzxks4n2lyw4ltzsd9wu6hw",
			nil,
			"02c40b914373d4fc9c1e4611ad0acd5f23abf58a0d",
			crypto.AddressTypeBLSAccount,
		},
		{
			"pc1rj65g93q7lpdq0366vst22l7va9d26j3l2vr0em",
			nil,
			"0396a882c41ef85a07c75a6416a57fcce95aad4a3f",
			crypto.AddressTypeEd25519Account,
		},
		{
			"pc1y90qakls8jlz9hyvdcsqsj0yj2lrqz26vqu7l0z",
			nil,
			"042bc1db7e0797c45b918dc401093c9257c6012b4c",
			crypto.AddressTypeSecp256k1Account,
		},
	}
	for no, tt := range tests {
		addr, err := crypto.AddressFromString(tt.encoded)
		if tt.err == nil {
			data, _ := hex.DecodeString(tt.hex)
			require.NoError(t, err, "test %v: unexpected error", no)
			assert.Equal(t, data, addr.Bytes(), "test %v: invalid result", no)
			assert.Equal(t, strings.ToLower(tt.encoded), addr.String(), "test %v: invalid encode", no)
			assert.Equal(t, tt.addrType, addr.Type(), "test %v: invalid type", no)
		} else {
			require.ErrorIs(t, err, tt.err, "test %v: invalid error", no)
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
			"",
			io.EOF,
		},
		{
			1,
			"00",
			nil,
		},
		{
			21,
			"040000000000000000000000000000000000000000",
			nil,
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
		{
			21,
			"04000102030405060708090a0b0c0d0e0f00010203",
			nil,
		},
		{
			0,
			"05000102030405060708090a0b0c0d0e0f00010203",
			crypto.InvalidAddressTypeError(5),
		},
	}
	for no, tt := range tests {
		data, _ := hex.DecodeString(tt.hex)
		buf := bytes.NewBuffer(data)
		addr := new(crypto.Address)

		err := addr.Decode(buf)
		if tt.err != nil {
			require.ErrorIs(t, err, tt.err, "test %v: error not matched", no)
			assert.Equal(t, tt.size, addr.SerializeSize(), "test %v invalid size", no)
		} else {
			require.NoError(t, err, "test %v expected no error", no)
			assert.Equal(t, tt.size, addr.SerializeSize(), "test %v invalid size", no)

			length := addr.SerializeSize()
			for i := 0; i < length; i++ {
				w := util.NewFixedWriter(i)
				require.Error(t, addr.Encode(w), "encode test %v failed", no)
			}
			w := util.NewFixedWriter(length)
			require.NoError(t, addr.Encode(w))
			assert.Equal(t, data, w.Bytes())
		}
	}
}

func TestShortString(t *testing.T) {
	h, err := crypto.AddressFromString("pc1pcs9ezsmn6n7fc8jxzxks4n2lyw4ltzsdc9v8qn")
	require.NoError(t, err)

	assert.Equal(t, "pc1pcs9e-9v8qn", h.ShortString())
	assert.Equal(t, h.ShortString(), h.LogString())
}

func TestAddressTypeString(t *testing.T) {
	tests := []struct {
		str string
		typ crypto.AddressType
		err error
	}{
		{
			str: "invalid_type",
			typ: crypto.AddressType(255),
			err: crypto.ErrInvalidAddressType,
		},
		{
			str: "validator",
			typ: crypto.AddressType(1),
			err: nil,
		},
		{
			str: "bls",
			typ: crypto.AddressType(2),
			err: nil,
		},
		{
			str: "ed25519",
			typ: crypto.AddressType(3),
			err: nil,
		},
		{
			str: "secp256k1",
			typ: crypto.AddressType(4),
			err: nil,
		},
	}

	for no, tt := range tests {
		typ, err := crypto.AddressTypeFromString(tt.str)
		if tt.err != nil {
			require.ErrorIs(t, err, tt.err, "test %v: error not matched", no)
			assert.Equal(t, tt.typ, typ)
		} else {
			require.NoError(t, err, "test %v expected no error", no)
			assert.Equal(t, tt.typ, typ)
			assert.Equal(t, tt.str, typ.String())
		}
	}

	t.Run("backward compatibility with old _account suffix", func(t *testing.T) {
		oldFormats := []struct {
			str string
			typ crypto.AddressType
		}{
			{"bls_account", crypto.AddressTypeBLSAccount},
			{"ed25519_account", crypto.AddressTypeEd25519Account},
			{"secp256k1_account", crypto.AddressTypeSecp256k1Account},
		}
		for _, tc := range oldFormats {
			typ, err := crypto.AddressTypeFromString(tc.str)
			require.NoError(t, err)
			assert.Equal(t, tc.typ, typ)
		}
	})
}
