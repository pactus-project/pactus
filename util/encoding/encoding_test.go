// This file contains code modified from the btcd project,
// which is licensed under the ISC License.
//
// Original license: https://github.com/btcsuite/btcd/blob/master/LICENSE
//

package encoding

import (
	"bytes"
	"crypto/rand"
	"io"
	"math/big"
	"reflect"
	"strings"
	"testing"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestElementEncoding tests encode and decode for various element types.  This
// is mainly to test the "fast" paths in readElement and writeElement which use
// type assertions to avoid reflection when possible.
func TestElementEncoding(t *testing.T) {
	type writeElementReflect int32

	tests := []struct {
		in  any    // Value to encode
		buf []byte // encoding bytes
	}{
		{int8(-128), []byte{0x80}},
		{int8(127), []byte{0x7f}},
		{uint8(1), []byte{0x01}},
		{int16(-32256), []byte{0x00, 0x82}},
		{int16(127), []byte{0x7f, 0x00}},
		{uint16(65535), []byte{0xff, 0xff}},
		{int32(-1), []byte{0xff, 0xff, 0xff, 0xff}},
		{int32(1), []byte{0x01, 0x00, 0x00, 0x00}},
		{uint32(256), []byte{0x00, 0x01, 0x00, 0x00}},
		{int64(-65536), []byte{0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}},
		{int64(65536), []byte{0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{uint64(4294967296), []byte{0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00}},
		{true, []byte{0x01}},
		{false, []byte{0x00}},
		{
			&hash.Hash{
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
				0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
				0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
				0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20,
			},
			[]byte{
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
				0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
				0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
				0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20,
			},
		},

		// Type not supported by the "fast" path and requires reflection.
		{
			writeElementReflect(1),
			[]byte{0x01, 0x00, 0x00, 0x00},
		},
	}

	t.Logf("Running %d tests", len(tests))
	for no, tt := range tests {
		var buf bytes.Buffer
		err := WriteElement(&buf, tt.in)
		assert.NoError(t, err, "writeElement #%d", no)
		assert.Equal(t, buf.Bytes(), tt.buf, "writeElement #%d", no)

		rbuf := bytes.NewReader(tt.buf)
		val := tt.in
		if reflect.ValueOf(tt.in).Kind() != reflect.Ptr {
			val = reflect.New(reflect.TypeOf(tt.in)).Interface()
		}
		err = ReadElement(rbuf, val)
		assert.NoError(t, err, "readElement #%d", no)

		ival := val
		if reflect.ValueOf(tt.in).Kind() != reflect.Ptr {
			ival = reflect.Indirect(reflect.ValueOf(val)).Interface()
		}
		assert.Equal(t, ival, tt.in, "readElement #%d", no)
	}
}

// TestElementEncodingErrors performs negative tests against encode and decode
// of various element types to confirm error paths work correctly.
func TestElementEncodingErrors(t *testing.T) {
	tests := []struct {
		in       any   // Value to encode
		max      int   // Max size of fixed buffer to induce errors
		writeErr error // Expected write error
		readErr  error // Expected read error
	}{
		{int8(127), 0, io.ErrShortWrite, io.EOF},
		{uint8(1), 0, io.ErrShortWrite, io.EOF},
		{int16(127), 1, io.ErrShortWrite, io.ErrUnexpectedEOF},
		{uint16(256), 1, io.ErrShortWrite, io.ErrUnexpectedEOF},
		{int32(256), 3, io.ErrShortWrite, io.ErrUnexpectedEOF},
		{uint32(256), 3, io.ErrShortWrite, io.ErrUnexpectedEOF},
		{int64(65536), 7, io.ErrShortWrite, io.ErrUnexpectedEOF},
		{uint64(4294967296), 7, io.ErrShortWrite, io.ErrUnexpectedEOF},
		{true, 0, io.ErrShortWrite, io.EOF},
		{false, 0, io.ErrShortWrite, io.EOF},
		{
			&hash.Hash{
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
				0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
				0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
				0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20,
			},
			0, io.ErrShortWrite, io.EOF,
		},
	}

	t.Logf("Running %d tests", len(tests))
	for no, tt := range tests {
		w := util.NewFixedWriter(tt.max)
		err := WriteElement(w, tt.in)
		assert.ErrorIs(t, err, tt.writeErr, "writeElement #%d", no)

		r := util.NewFixedReader(tt.max, nil)
		val := tt.in
		if reflect.ValueOf(tt.in).Kind() != reflect.Ptr {
			val = reflect.New(reflect.TypeOf(tt.in)).Interface()
		}
		err = ReadElement(r, val)
		assert.ErrorIs(t, err, tt.readErr, "readElement #%d", no)
	}
}

// TestVarStringEncoding tests encode and decode for variable length strings.
func TestVarStringEncoding(t *testing.T) {
	// str256 is a string that takes a 2-byte varint to encode.
	str256 := strings.Repeat("test", 64)

	tests := []struct {
		in  string // String to encode
		out string // String to decoded value
		buf []byte // Encoding bytes
	}{
		// Latest protocol version.
		// Empty string
		{"", "", []byte{0x00}},
		// Single byte varint + string
		{"Test", "Test", append([]byte{0x04}, []byte("Test")...)},
		// 2-byte varint + string
		{str256, str256, append([]byte{0x80, 0x02}, []byte(str256)...)},
	}

	t.Logf("Running %d tests", len(tests))
	for no, tt := range tests {
		var buf bytes.Buffer
		err := WriteVarString(&buf, tt.in)
		assert.NoError(t, err, "WriteVarString #%d ", no)
		assert.Equal(t, buf.Bytes(), tt.buf, "WriteVarString #%d", no)

		rbuf := bytes.NewReader(tt.buf)
		val, err := ReadVarString(rbuf)
		assert.NoError(t, err, "ReadVarString #%d", no)
		assert.Equal(t, val, tt.out, "ReadVarString #%d", no)
		assert.Equal(t, VarStringSerializeSize(tt.in), len(tt.buf))
	}
}

// TestVarStringEncodingErrors performs negative tests against encode and
// decode of variable length strings to confirm error paths work correctly.
func TestVarStringEncodingErrors(t *testing.T) {
	// str256 is a string that takes a 2-byte varint to encode.
	str256 := strings.Repeat("test", 64)

	tests := []struct {
		in       string // Value to encode
		buf      []byte // Encoding bytes
		max      int    // Max size of fixed buffer to induce errors
		writeErr error  // Expected write error
		readErr  error  // Expected read error
	}{
		// Latest protocol version with intentional read/write errors.
		// Force errors on empty string.
		{"", []byte{0x00}, 0, io.ErrShortWrite, io.EOF},
		// Force error on single byte varint + string.
		{"Test", []byte{0x04}, 2, io.ErrShortWrite, io.ErrUnexpectedEOF},
		// Force errors on 2-byte varint + string.
		{str256, []byte{0x80}, 1, io.ErrShortWrite, io.EOF},
	}

	t.Logf("Running %d tests", len(tests))
	for no, tt := range tests {
		w := util.NewFixedWriter(tt.max)
		err := WriteVarString(w, tt.in)
		assert.ErrorIs(t, err, tt.writeErr, "WriteVarString #%d", no)

		r := util.NewFixedReader(tt.max, tt.buf)
		_, err = ReadVarString(r)
		assert.ErrorIs(t, err, tt.readErr, "ReadVarString #%d wrong", no)
	}
}

// TestVarStringOverflowErrors performs tests to ensure deserializing variable
// length strings intentionally crafted to use large values for the string
// length are handled properly.  This could otherwise potentially be used as an
// attack vector.
func TestVarStringOverflowErrors(t *testing.T) {
	tests := []struct {
		buf []byte // Encoding bytes
	}{
		{[]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01}},
		{[]byte{0x80, 0x80, 0x80, 0x11}},
	}

	t.Logf("Running %d tests", len(tests))
	for no, tt := range tests {
		rbuf := bytes.NewReader(tt.buf)
		_, err := ReadVarString(rbuf)
		assert.Contains(t, err.Error(), "variable length string is too long", "ReadVarString #%d", no)
	}
}

// TestVarBytesEncoding tests encode and decode for variable length byte array.
func TestVarBytesEncoding(t *testing.T) {
	// bytes256 is a byte array that takes a 2-byte varint to encode.
	bytes256 := bytes.Repeat([]byte{0x01}, 256)

	tests := []struct {
		in  []byte // Byte Array to write
		buf []byte // Encoding bytes
	}{
		// Latest protocol version.
		// Empty byte array
		{[]byte{}, []byte{0x00}},
		// Single byte varint + byte array
		{[]byte{0x01}, []byte{0x01, 0x01}},
		// 2-byte varint + byte array
		{bytes256, append([]byte{0x80, 0x02}, bytes256...)},
	}

	t.Logf("Running %d tests", len(tests))
	for no, tt := range tests {
		var buf bytes.Buffer
		err := WriteVarBytes(&buf, tt.in)
		assert.NoError(t, err, "WriteVarBytes #%d", no)
		assert.Equal(t, buf.Bytes(), tt.buf, "WriteVarBytes #%d", no)

		rbuf := bytes.NewReader(tt.buf)
		val, err := ReadVarBytes(rbuf)
		assert.NoError(t, err, "ReadVarBytes #%d", no)
		assert.Equal(t, buf.Bytes(), tt.buf, "ReadVarBytes #%d", no)
		assert.Equal(t, val, tt.in, "ReadVarBytes #%d", no)
		assert.Equal(t, VarBytesSerializeSize(tt.in), len(tt.buf))
	}
}

// TestVarBytesEncodingErrors performs negative tests against encode and
// decode of variable length byte arrays to confirm error paths work correctly.
func TestVarBytesEncodingErrors(t *testing.T) {
	// bytes256 is a byte array that takes a 2-byte varint to encode.
	bytes256 := bytes.Repeat([]byte{0x01}, 256)

	tests := []struct {
		in       []byte // Byte Array to write
		buf      []byte // Encoding bytes
		max      int    // Max size of fixed buffer to induce errors
		writeErr error  // Expected write error
		readErr  error  // Expected read error
	}{
		// Latest protocol version with intentional read/write errors.
		// Force errors on empty byte array.
		{[]byte{}, []byte{0x00}, 0, io.ErrShortWrite, io.EOF},
		// Force error on single byte varint + byte array.
		{[]byte{0x01, 0x02, 0x03}, []byte{0x04}, 2, io.ErrShortWrite, io.ErrUnexpectedEOF},
		// Force errors on 2-byte varint + byte array.
		{bytes256, []byte{0x80}, 1, io.ErrShortWrite, io.EOF},
	}

	t.Logf("Running %d tests", len(tests))
	for no, tt := range tests {
		w := util.NewFixedWriter(tt.max)
		err := WriteVarBytes(w, tt.in)
		assert.ErrorIs(t, err, tt.writeErr, "WriteVarBytes #%d", no)

		r := util.NewFixedReader(tt.max, tt.buf)
		_, err = ReadVarBytes(r)
		assert.ErrorIs(t, err, tt.readErr, "ReadVarBytes #%d", no)
	}
}

// TestVarBytesOverflowErrors performs tests to ensure deserializing variable
// length byte arrays intentionally crafted to use large values for the array
// length are handled properly. This could otherwise potentially be used as an
// attack vector.
func TestVarBytesOverflowErrors(t *testing.T) {
	tests := []struct {
		buf []byte
	}{
		{[]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01}},
		{[]byte{0x80, 0x80, 0x80, 0x11}},
	}

	t.Logf("Running %d tests", len(tests))
	for no, tt := range tests {
		rbuf := bytes.NewReader(tt.buf)
		_, err := ReadVarBytes(rbuf)
		assert.Contains(t, err.Error(), "variable length byte array is too long", "ReadVarString #%d", no)
	}
}

// TestVarInt performs tests to ensure deserializing variable integers are
// handled properly. This could otherwise potentially be used as an attack
// vector.
func TestVarInt(t *testing.T) {
	tests := []struct {
		in  uint64 // Value to encode
		buf []byte // Encoded bytes
	}{
		{uint64(0x0), []byte{0x00}},
		{uint64(0xff), []byte{0xff, 0x01}},
		{uint64(0x7fff), []byte{0xff, 0xff, 0x01}},
		{uint64(0x3fffff), []byte{0xff, 0xff, 0xff, 0x01}},
		{uint64(0x1fffffff), []byte{0xff, 0xff, 0xff, 0xff, 0x01}},
		{uint64(0xfffffffff), []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0x01}},
		{uint64(0x7ffffffffff), []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01}},
		{uint64(0x3ffffffffffff), []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01}},
		{uint64(0x1ffffffffffffff), []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01}},
		{uint64(0xffffffffffffffff), []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01}},
		{uint64(0x200), []byte{0x80, 0x04}},
		{uint64(0x027f), []byte{0xff, 0x04}},
		{uint64(0xff00000000), []byte{0x80, 0x80, 0x80, 0x80, 0xf0, 0x1f}},
		{uint64(0xffffffff), []byte{0xff, 0xff, 0xff, 0xff, 0x0f}},
		{uint64(0x100000000), []byte{0x80, 0x80, 0x80, 0x80, 0x10}},
		{uint64(0x7ffffffff), []byte{0xff, 0xff, 0xff, 0xff, 0x7f}},
		{uint64(0x800000000), []byte{0x80, 0x80, 0x80, 0x80, 0x80, 0x1}},
	}
	for no, tt := range tests {
		var buf bytes.Buffer
		err := WriteVarInt(&buf, tt.in)
		assert.NoError(t, err, "WriteVarInt #%d", no)
		assert.Equal(t, buf.Bytes(), tt.buf, "WriteVarInt #%d", no)

		val, err := ReadVarInt(&buf)
		assert.NoError(t, err, "ReadVarInt #%d", no)
		assert.Equal(t, val, tt.in, "ReadVarInt #%d", no)
		assert.Equal(t, VarIntSerializeSize(tt.in), len(tt.buf))
	}
}

// TestVarIntError ensures variable length integers that are not encoded
// properly return the expected error.
func TestVarIntError(t *testing.T) {
	tests := []struct {
		in      []byte // Value to decode
		readErr error
	}{
		{
			[]byte{0x98, 0}, ErrNonCanonical,
		},
		{
			[]byte{0xFF}, io.EOF,
		},
		{
			[]byte{0x80, 0x00}, ErrNonCanonical,
		},
		{
			[]byte{0x80, 0xfe}, io.EOF,
		},
		{
			[]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x02}, ErrOverflow,
		},
		{
			[]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, ErrOverflow,
		},
	}

	t.Logf("Running %d tests", len(tests))
	for no, tt := range tests {
		rbuf := bytes.NewReader(tt.in)
		val, err := ReadVarInt(rbuf)
		assert.ErrorIs(t, err, tt.readErr, "ReadVarInt #%d", no)
		assert.Zero(t, val, "ReadVarInt #%d", no)
	}
}

func TestVarIntRandom(t *testing.T) {
	max := new(big.Int).SetUint64(^uint64(0)) // max uint64
	randIntBig, _ := rand.Int(rand.Reader, max)
	randInt1 := randIntBig.Uint64()

	var wBuf bytes.Buffer
	err := WriteVarInt(&wBuf, randInt1)
	require.NoError(t, err)

	rBuf := bytes.NewReader(wBuf.Bytes())
	randInt2, err := ReadVarInt(rBuf)
	require.NoError(t, err)

	assert.Equal(t, randInt1, randInt2)
}

func TestWriteElements(t *testing.T) {
	el1 := uint8(1)
	el2 := uint16(2)
	el3 := uint32(3)
	el4 := uint64(4)
	var buf bytes.Buffer
	err := WriteElements(&buf, &el1, &el2, &el3, &el4)
	assert.NoError(t, err)
	assert.Equal(t, buf.Bytes(), []byte{0x1, 0x2, 0x0, 0x3, 0x0, 0x0, 0x0, 0x4, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0})
}

func TestReadElements(t *testing.T) {
	el1 := uint8(1)
	el2 := uint16(2)
	el3 := uint32(3)
	el4 := uint64(4)
	r := bytes.NewReader([]byte{0x1, 0x2, 0x0, 0x3, 0x0, 0x0, 0x0, 0x4, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0})
	err := ReadElements(r, &el1, &el2, &el3, &el4)
	assert.NoError(t, err)
	assert.Equal(t, el1, uint8(1))
	assert.Equal(t, el2, uint16(2))
	assert.Equal(t, el3, uint32(3))
	assert.Equal(t, el4, uint64(4))
}
