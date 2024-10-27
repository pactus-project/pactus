// This file contains code modified from the btcd project,
// which is licensed under the ISC License.
//
// Original license: https://github.com/btcsuite/btcd/blob/master/LICENSE
//

package bech32m

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"testing"
)

// TestBech32M tests that the following set of strings, based on the test
// vectors in BIP-350 are either valid or invalid using the new bech32m
// checksum algo. Some of these strings are similar to the set of above test
// vectors, but end up with different checksums.
func TestBech32M(t *testing.T) {
	tests := []struct {
		str           string
		expectedError error
	}{
		{"A1LQFN3A", nil},
		{"a1lqfn3a", nil},
		{"an83characterlonghumanreadablepartthatcontainsthetheexcludedcharactersbioandnumber11sg7hg6", nil},
		{"abcdef1l7aum6echk45nj3s0wdvt2fg8x9yrzpqzd3ryx", nil},
		{"11llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllludsr8", nil},
		{"split1checkupstagehandshakeupstreamerranterredcaperredlc445v", nil},
		{"?1v759aa", nil},

		// Additional test vectors used in bitcoin core
		{"\x201xj0phk", InvalidCharacterError('\x20')},
		{"\x7f1g6xzxy", InvalidCharacterError('\x7f')},
		{"\x801vctc34", InvalidCharacterError('\x80')},
		{
			"an84characterslonghumanreadablepartthatcontainsthetheexcludedcharactersbioandnumber11d6pts4",
			InvalidLengthError(91),
		},
		{"qyrz8wqd2c9m", InvalidSeparatorIndexError(-1)},
		{"1qyrz8wqd2c9m", InvalidSeparatorIndexError(0)},
		{"y1b0jsk6g", NonCharsetCharError(98)},
		{"lt1igcx5c0", NonCharsetCharError(105)},
		{"in1muywd", InvalidSeparatorIndexError(2)},
		{"mm1crxm3i", NonCharsetCharError(105)},
		{"au1s5cgom", NonCharsetCharError(111)},
		{"16plkw9", InvalidLengthError(7)},
		{"1p2gdwpf", InvalidSeparatorIndexError(0)},

		{" 1nwldj5", InvalidCharacterError(' ')},
		{"\x7f" + "1axkwrx", InvalidCharacterError(0x7f)},
		{"\x801eym55h", InvalidCharacterError(0x80)},
	}

	for no, tt := range tests {
		str := tt.str
		hrp, decoded, err := Decode(str)
		if !errors.Is(err, tt.expectedError) {
			t.Errorf("%d: (%v) expected decoding error %v "+
				"instead got %v", no, str, tt.expectedError,
				err)

			continue
		}

		if err != nil {
			// End test case here if a decoding error was expected.
			continue
		}

		// Check that it encodes to the same string, using bech32 m.
		encoded, err := Encode(hrp, decoded)
		if err != nil {
			t.Errorf("encoding failed: %v", err)
		}

		if !strings.EqualFold(encoded, str) {
			t.Errorf("expected data to encode to %v, but got %v",
				str, encoded)
		}

		// Flip a bit in the string an make sure it is caught.
		pos := strings.LastIndexAny(str, "1")
		flipped := str[:pos+1] + string((str[pos+1] ^ 1)) + str[pos+2:]
		_, _, err = Decode(flipped)
		if err == nil {
			t.Error("expected decoding to fail")
		}
	}
}

// TestMixedCaseEncode ensures mixed case HRPs are converted to lowercase as
// expected when encoding and that decoding the produced encoding when converted
// to all uppercase produces the lowercase HRP and original data.
func TestMixedCaseEncode(t *testing.T) {
	tests := []struct {
		name    string
		hrp     string
		data    string
		encoded string
	}{{
		name:    "all uppercase HRP with no data",
		hrp:     "A",
		data:    "",
		encoded: "a1lqfn3a",
	}, {
		name:    "all uppercase HRP with data",
		hrp:     "UPPERCASE",
		data:    "787878",
		encoded: "uppercase10pu8s9vw67r",
	}, {
		name:    "mixed case HRP even offsets uppercase",
		hrp:     "AbCdEf",
		data:    "00443214c74254b635cf84653a56d7c675be77df",
		encoded: "abcdef1qpzry9x8gf2tvdw0s3jn54khce6mua7lwusvrv",
	}, {
		name:    "mixed case HRP odd offsets uppercase ",
		hrp:     "aBcDeF",
		data:    "00443214c74254b635cf84653a56d7c675be77df",
		encoded: "abcdef1qpzry9x8gf2tvdw0s3jn54khce6mua7lwusvrv",
	}, {
		name:    "all lowercase HRP",
		hrp:     "abcdef",
		data:    "00443214c74254b635cf84653a56d7c675be77df",
		encoded: "abcdef1qpzry9x8gf2tvdw0s3jn54khce6mua7lwusvrv",
	}}

	for _, tt := range tests {
		// Convert the text hex to bytes, convert those bytes from base256 to
		// base32, then ensure the encoded result with the HRP provided in the
		// test data is as expected.
		data, err := hex.DecodeString(tt.data)
		if err != nil {
			t.Errorf("%q: invalid hex %q: %v", tt.name, tt.data, err)

			continue
		}
		convertedData, err := ConvertBits(data, 8, 5, true)
		if err != nil {
			t.Errorf("%q: unexpected convert bits error: %v", tt.name,
				err)

			continue
		}
		gotEncoded, err := Encode(tt.hrp, convertedData)
		if err != nil {
			t.Errorf("%q: unexpected encode error: %v", tt.name, err)

			continue
		}
		if gotEncoded != tt.encoded {
			t.Errorf("%q: mismatched encoding -- got %q, want %q", tt.name,
				gotEncoded, tt.encoded)

			continue
		}

		// Ensure the decoding the expected lowercase encoding converted to all
		// uppercase produces the lowercase HRP and original data.
		gotHRP, gotData, err := Decode(strings.ToUpper(tt.encoded))
		if err != nil {
			t.Errorf("%q: unexpected decode error: %v", tt.name, err)

			continue
		}
		wantHRP := strings.ToLower(tt.hrp)
		if gotHRP != wantHRP {
			t.Errorf("%q: mismatched decoded HRP -- got %q, want %q", tt.name,
				gotHRP, wantHRP)

			continue
		}
		convertedGotData, err := ConvertBits(gotData, 5, 8, false)
		if err != nil {
			t.Errorf("%q: unexpected convert bits error: %v", tt.name,
				err)

			continue
		}
		if !bytes.Equal(convertedGotData, data) {
			t.Errorf("%q: mismatched data -- got %x, want %x", tt.name,
				convertedGotData, data)

			continue
		}
	}
}

// TestCanDecodeUnlimtedBech32 tests whether decoding a large bech32 string works
// when using the DecodeNoLimit version.
func TestCanDecodeUnlimtedBech32(t *testing.T) {
	input := "11qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqsqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqp2krp0"

	// basic check that an input of this length errors on regular Decode()
	_, _, err := Decode(input)
	if err == nil {
		t.Fatalf("Test vector not appropriate")
	}

	// Try and decode it.
	hrp, data, err := DecodeNoLimit(input)
	if err != nil {
		t.Fatalf("Expected decoding of large string to work. Got error: %v", err)
	}

	// Verify data for correctness.
	if hrp != "1" {
		t.Fatalf("Unexpected hrp: %v", hrp)
	}
	decodedHex := fmt.Sprintf("%x", data)
	expected := "00000000000000000000000000000000000000000000000000000000000000000000000000000000000" +
		"00000000010000000000000000000000000000000000000000000000000000000000000000000000000"
	if decodedHex != expected {
		t.Fatalf("Unexpected decoded data: %s", decodedHex)
	}
}

// TestBech32Base256 ensures decoding and encoding various bech32, HRPs, and
// data produces the expected results when using EncodeFromBase256 and
// DecodeToBase256.  It includes tests for proper handling of case
// manipulations.
func TestBech32Base256(t *testing.T) {
	tests := []struct {
		name    string // test name
		encoded string // bech32 string to decode
		hrp     string // expected human-readable part
		data    string // expected hex-encoded data
		err     error  // expected error
	}{{
		name:    "all uppercase, no data",
		encoded: "A1LQFN3A",
		hrp:     "a",
		data:    "",
	}, {
		name:    "long hrp with separator and excluded chars, no data",
		encoded: "an83characterlonghumanreadablepartthatcontainsthenumber1andtheexcludedcharactersbio17hy8dj",
		hrp:     "an83characterlonghumanreadablepartthatcontainsthenumber1andtheexcludedcharactersbio",
		data:    "",
	}, {
		name:    "6 char hrp with data with leading zero",
		encoded: "abcdef1qpzry9x8gf2tvdw0s3jn54khce6mua7lwusvrv",
		hrp:     "abcdef",
		data:    "00443214c74254b635cf84653a56d7c675be77df",
	}, {
		name:    "hrp same as separator and max length encoded string",
		encoded: "11qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqdm6ems",
		hrp:     "1",
		data:    "000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
	}, {
		name:    "5 char hrp with data chosen to produce human-readable data part",
		encoded: "split1checkupstagehandshakeupstreamerranterredcaperredlc445v",
		hrp:     "split",
		data:    "c5f38b70305f519bf66d85fb6cf03058f3dde463ecd7918f2dc743918f2d",
	}, {
		name:    "same as previous but with checksum invalidated",
		encoded: "split1checkupstagehandshakeupstreamerranterredcaperred2y9e2w",
		err:     InvalidChecksumError{"lc445v", "2y9e2w"},
	}, {
		name:    "hrp with invalid character (space)",
		encoded: "s lit1checkupstagehandshakeupstreamerranterredcaperredp8hs2p",
		err:     InvalidCharacterError(' '),
	}, {
		name:    "hrp with invalid character (DEL)",
		encoded: "spl\x7ft1checkupstagehandshakeupstreamerranterredcaperredlc445v",
		err:     InvalidCharacterError(127),
	}, {
		name:    "data part with invalid character (o)",
		encoded: "split1cheo2y9e2w",
		err:     NonCharsetCharError('o'),
	}, {
		name:    "data part too short",
		encoded: "split1a2y9w",
		err:     InvalidSeparatorIndexError(5),
	}, {
		name:    "empty hrp",
		encoded: "1checkupstagehandshakeupstreamerranterredcaperredlc445v",
		err:     InvalidSeparatorIndexError(0),
	}, {
		name:    "no separator",
		encoded: "pzry9x0s0muk",
		err:     InvalidSeparatorIndexError(-1),
	}, {
		name:    "too long by one char",
		encoded: "11qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqsqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqc8247j",
		err:     InvalidLengthError(91),
	}, {
		name:    "invalid due to mixed case in hrp",
		encoded: "aBcdef1qpzry9x8gf2tvdw0s3jn54khce6mua7lmqqqxw",
		err:     MixedCaseError{},
	}, {
		name:    "invalid due to mixed case in data part",
		encoded: "abcdef1Qpzry9x8gf2tvdw0s3jn54khce6mua7lmqqqxw",
		err:     MixedCaseError{},
	}}

	for _, tt := range tests {
		// Ensure the decode either produces an error or not as expected.
		str := tt.encoded
		gotHRP, gotData, err := DecodeToBase256(str)
		if !errors.Is(tt.err, err) {
			t.Errorf("%q: unexpected decode error -- got %v, want %v",
				tt.name, err, tt.err)

			continue
		}
		if err != nil {
			// End test case here if a decoding error was expected.
			continue
		}

		// Ensure the expected HRP and original data are as expected.
		if gotHRP != tt.hrp {
			t.Errorf("%q: mismatched decoded HRP -- got %q, want %q", tt.name,
				gotHRP, tt.hrp)

			continue
		}
		data, err := hex.DecodeString(tt.data)
		if err != nil {
			t.Errorf("%q: invalid hex %q: %v", tt.name, tt.data, err)

			continue
		}
		if !bytes.Equal(gotData, data) {
			t.Errorf("%q: mismatched data -- got %x, want %x", tt.name,
				gotData, data)

			continue
		}

		// Encode the same data with the HRP converted to all uppercase and
		// ensure the result is the lowercase version of the original encoded
		// bech32 string.
		gotEncoded, err := EncodeFromBase256(strings.ToUpper(tt.hrp), data)
		if err != nil {
			t.Errorf("%q: unexpected uppercase HRP encode error: %v", tt.name,
				err)
		}
		wantEncoded := strings.ToLower(str)
		if gotEncoded != wantEncoded {
			t.Errorf("%q: mismatched encoding -- got %q, want %q", tt.name,
				gotEncoded, wantEncoded)
		}

		// Encode the same data with the HRP converted to all lowercase and
		// ensure the result is the lowercase version of the original encoded
		// bech32 string.
		gotEncoded, err = EncodeFromBase256(strings.ToLower(tt.hrp), data)
		if err != nil {
			t.Errorf("%q: unexpected lowercase HRP encode error: %v", tt.name,
				err)
		}
		if gotEncoded != wantEncoded {
			t.Errorf("%q: mismatched encoding -- got %q, want %q", tt.name,
				gotEncoded, wantEncoded)
		}

		// Encode the same data with the HRP converted to mixed upper and
		// lowercase and ensure the result is the lowercase version of the
		// original encoded bech32 string.
		var mixedHRPBuilder strings.Builder
		for index, chr := range tt.hrp {
			if index%2 == 0 {
				mixedHRPBuilder.WriteString(strings.ToUpper(string(chr)))

				continue
			}
			mixedHRPBuilder.WriteRune(chr)
		}
		gotEncoded, err = EncodeFromBase256(mixedHRPBuilder.String(), data)
		if err != nil {
			t.Errorf("%q: unexpected lowercase HRP encode error: %v", tt.name,
				err)
		}
		if gotEncoded != wantEncoded {
			t.Errorf("%q: mismatched encoding -- got %q, want %q", tt.name,
				gotEncoded, wantEncoded)
		}

		// Ensure a bit flip in the string is caught.
		pos := strings.LastIndexAny(tt.encoded, "1")
		flipped := str[:pos+1] + string((str[pos+1] ^ 1)) + str[pos+2:]
		_, _, err = DecodeToBase256(flipped)
		if err == nil {
			t.Error("expected decoding to fail")
		}
	}
}

// BenchmarkEncodeDecodeCycle performs a benchmark for a full encode/decode
// cycle of a bech32 string. It also  reports the allocation count, which we
// expect to be 2 for a fully optimized cycle.
func BenchmarkEncodeDecodeCycle(b *testing.B) {
	// Use a fixed, 49-byte raw data for testing.
	inputData, err := hex.DecodeString(
		"cbe6365ddbcda9a9915422c3f091c13f8c7b2f263b8d34067bd12c274408473fa764871c9dd51b1bb34873b3473b633ed1")
	if err != nil {
		b.Fatalf("failed to initialize input data: %v", err)
	}

	// Convert this into a 79-byte, base 32 byte slice.
	base32Input, err := ConvertBits(inputData, 8, 5, true)
	if err != nil {
		b.Fatalf("failed to convert input to 32 bits-per-element: %v", err)
	}

	// Use a fixed hrp for the tests. This should generate an encoded bech32
	// string of size 90 (the maximum allowed by BIP-173).
	hrp := "bc"

	// Begin the benchmark. Given that we test one roundtrip per iteration
	// (that is, one Encode() and one Decode() operation), we expect at most
	// 2 allocations per reported test op.
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		str, err := Encode(hrp, base32Input)
		if err != nil {
			b.Fatalf("failed to encode input: %v", err)
		}

		_, _, err = Decode(str)
		if err != nil {
			b.Fatalf("failed to decode string: %v", err)
		}
	}
}

// TestConvertBits tests whether base conversion works using TestConvertBits().
func TestConvertBits(t *testing.T) {
	tests := []struct {
		input    string
		output   string
		fromBits uint8
		toBits   uint8
		pad      bool
	}{
		// Trivial empty conversions.
		{"", "", 8, 5, false},
		{"", "", 8, 5, true},
		{"", "", 5, 8, false},
		{"", "", 5, 8, true},

		// Conversions of 0 value with/without padding.
		{"00", "00", 8, 5, false},
		{"00", "0000", 8, 5, true},
		{"0000", "00", 5, 8, false},
		{"0000", "0000", 5, 8, true},

		// Testing when conversion ends exactly at the byte edge. This makes
		// both padded and unpadded versions the same.
		{"0000000000", "0000000000000000", 8, 5, false},
		{"0000000000", "0000000000000000", 8, 5, true},
		{"0000000000000000", "0000000000", 5, 8, false},
		{"0000000000000000", "0000000000", 5, 8, true},

		// Conversions of full byte sequences.
		{"ffffff", "1f1f1f1f1e", 8, 5, true},
		{"1f1f1f1f1e", "ffffff", 5, 8, false},
		{"1f1f1f1f1e", "ffffff00", 5, 8, true},

		// Sample random conversions.
		{"c9ca", "190705", 8, 5, false},
		{"c9ca", "19070500", 8, 5, true},
		{"19070500", "c9ca", 5, 8, false},
		{"19070500", "c9ca00", 5, 8, true},

		// Test cases tested on TestConvertBitsFailures with their corresponding
		// fixes.
		{"ff", "1f1c", 8, 5, true},
		{"1f1c10", "ff20", 5, 8, true},

		// Large conversions.
		{
			"cbe6365ddbcda9a9915422c3f091c13f8c7b2f263b8d34067bd12c274408473fa764871c9dd51b1bb34873b3473b633ed1",
			"190f13030c170e1b1916141a13040a14040b011f01040e01071e0607160b1906070e06130801131" +
				"b1a0416020e110008081c1f1a0e19040703120e1d0a06181b160d0407070c1a07070d11131d1408",
			8, 5, true,
		},
		{
			"190f13030c170e1b1916141a13040a14040b011f01040e01071e0607160b1906070e06130801131" +
				"b1a0416020e110008081c1f1a0e19040703120e1d0a06181b160d0407070c1a07070d11131d1408",
			"cbe6365ddbcda9a9915422c3f091c13f8c7b2f263b8d34067bd12c274408473fa764871c9dd51b1bb34873b3473b633ed100",
			5, 8, true,
		},
	}

	for no, tt := range tests {
		input, err := hex.DecodeString(tt.input)
		if err != nil {
			t.Fatalf("invalid test input data: %v", err)
		}

		expected, err := hex.DecodeString(tt.output)
		if err != nil {
			t.Fatalf("invalid test output data: %v", err)
		}

		actual, err := ConvertBits(input, tt.fromBits, tt.toBits, tt.pad)
		if err != nil {
			t.Fatalf("test case %d failed: %v", no, err)
		}

		if !bytes.Equal(actual, expected) {
			t.Fatalf("test case %d has wrong output; expected=%x actual=%x",
				no, expected, actual)
		}
	}
}

// TestConvertBitsFailures tests for the expected conversion failures of
// ConvertBits().
func TestConvertBitsFailures(t *testing.T) {
	tests := []struct {
		input    string
		fromBits uint8
		toBits   uint8
		pad      bool
		err      error
	}{
		// Not enough output bytes when not using padding.
		{"ff", 8, 5, false, InvalidIncompleteGroupError{}},
		{"1f1c10", 5, 8, false, InvalidIncompleteGroupError{}},

		// Unsupported bit conversions.
		{"", 0, 5, false, InvalidBitGroupsError{}},
		{"", 10, 5, false, InvalidBitGroupsError{}},
		{"", 5, 0, false, InvalidBitGroupsError{}},
		{"", 5, 10, false, InvalidBitGroupsError{}},
	}

	for no, tt := range tests {
		input, err := hex.DecodeString(tt.input)
		if err != nil {
			t.Fatalf("invalid test input data: %v", err)
		}

		_, err = ConvertBits(input, tt.fromBits, tt.toBits, tt.pad)
		if !errors.Is(err, tt.err) {
			t.Fatalf("test case %d failure: expected '%v' got '%v'", no,
				tt.err, err)
		}
	}
}

// BenchmarkConvertBitsDown benchmarks the speed and memory allocation behavior
// of ConvertBits when converting from a higher base into a lower base (e.g. 8
// => 5).
//
// Only a single allocation is expected, which is used for the output array.
func BenchmarkConvertBitsDown(b *testing.B) {
	// Use a fixed, 49-byte raw data for testing.
	inputData, err := hex.DecodeString(
		"cbe6365ddbcda9a9915422c3f091c13f8c7b2f263b8d34067bd12c274408473fa764871c9dd51b1bb34873b3473b633ed1")
	if err != nil {
		b.Fatalf("failed to initialize input data: %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ConvertBits(inputData, 8, 5, true)
		if err != nil {
			b.Fatalf("error converting bits: %v", err)
		}
	}
}

// BenchmarkConvertBitsDown benchmarks the speed and memory allocation behavior
// of ConvertBits when converting from a lower base into a higher base (e.g. 5
// => 8).
//
// Only a single allocation is expected, which is used for the output array.
func BenchmarkConvertBitsUp(b *testing.B) {
	// Use a fixed, 79-byte raw data for testing.
	inputData, err := hex.DecodeString(
		"190f13030c170e1b1916141a13040a14040b011f01040e01071e0607160b1906070e06130801131" +
			"b1a0416020e110008081c1f1a0e19040703120e1d0a06181b160d0407070c1a07070d11131d1408")
	if err != nil {
		b.Fatalf("failed to initialize input data: %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ConvertBits(inputData, 8, 5, true)
		if err != nil {
			b.Fatalf("error converting bits: %v", err)
		}
	}
}

// TestEncodeFromBase256WithType tests for the expected behavior of
// EncodeFromBase256WithType function.
func TestEncodeFromBase256WithType(t *testing.T) {
	tests := []struct {
		hrp           string
		typ           byte
		input         string
		expectedBech  string
		expectedError error
	}{
		{"A", 0, "", "a1qy52hkn", nil},
		{"AbC", 1, "1234", "abc1pzg6qgtt0h8", nil},
		{"", 1, "abcd", "1p40xsjtqww4", nil},
		{"", 32, "1", "", InvalidDataByteError(32)},
	}

	for no, tt := range tests {
		data, _ := hex.DecodeString(tt.input)
		enc, err := EncodeFromBase256WithType(tt.hrp, tt.typ, data)
		if !errors.Is(err, tt.expectedError) {
			t.Errorf("%d: (%v) expected encoding error "+
				"instead got %v", no, tt.expectedError,
				err)

			continue
		}

		if enc != tt.expectedBech {
			t.Errorf("%d: mismatched encoding -- got %q, want %q", no,
				enc, tt.expectedBech)
		}
	}
}

// TestDecodeToBase256WithTypeNoLimit tests for the expected behavior of
// DecodeToBase256WithTypeNoLimit function.
func TestDecodeToBase256WithTypeNoLimit(t *testing.T) {
	tests := []struct {
		bech          string
		expectedHRP   string
		expectedTyp   byte
		expectedData  string
		expectedError error
	}{
		{"a1qy52hkn", "a", 0, "", nil},
		{"abc1pzg6qgtt0h8", "abc", 1, "1234", nil},
		{"1p40xsjtqww4", "", 0, "", InvalidSeparatorIndexError(0)},
		{"a1lqfn3a", "", 0, "", InvalidLengthError(0)},
	}

	for no, tt := range tests {
		hrp, typ, data, err := DecodeToBase256WithTypeNoLimit(tt.bech)
		if !errors.Is(err, tt.expectedError) {
			t.Errorf("%d: (%v) expected encoding error "+
				"instead got %v", no, tt.expectedError,
				err)

			continue
		}

		if hrp != tt.expectedHRP {
			t.Errorf("%d: mismatched HRP -- got %q, want %q", no,
				hrp, tt.expectedHRP)
		}

		if typ != tt.expectedTyp {
			t.Errorf("%d: mismatched Type -- got %q, want %q", no,
				typ, tt.expectedTyp)
		}

		expectedData, _ := hex.DecodeString(tt.expectedData)
		if !bytes.Equal(expectedData, data) {
			t.Errorf("%d: mismatched HRP -- got \"%x\", want %q", no,
				data, tt.expectedData)
		}
	}
}
