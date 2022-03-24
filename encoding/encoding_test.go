package encoding

import (
	"bytes"
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/util"
)

// import (
// 	"bytes"
// 	"fmt"
// 	"io"
// 	"reflect"
// 	"strings"
// 	"testing"

// 	"github.com/btcsuite/btcd/chaincfg/chainhash"
// 	"github.com/davecgh/go-spew/spew"
// )

// // mainNetGenesisHash is the hash of the first block in the block chain for the
// // main network (genesis block).
// var mainNetGenesisHash = chainhash.Hash([chainhash.HashSize]byte{ // Make go vet happy.
// 	0x6f, 0xe2, 0x8c, 0x0a, 0xb6, 0xf1, 0xb3, 0x72,
// 	0xc1, 0xa6, 0xa2, 0x46, 0xae, 0x63, 0xf7, 0x4f,
// 	0x93, 0x1e, 0x83, 0x65, 0xe1, 0x5a, 0x08, 0x9c,
// 	0x68, 0xd6, 0x19, 0x00, 0x00, 0x00, 0x00, 0x00,
// })

// // mainNetGenesisMerkleRoot is the hash of the first transaction in the genesis
// // block for the main network.
// var mainNetGenesisMerkleRoot = chainhash.Hash([chainhash.HashSize]byte{ // Make go vet happy.
// 	0x3b, 0xa3, 0xed, 0xfd, 0x7a, 0x7b, 0x12, 0xb2,
// 	0x7a, 0xc7, 0x2c, 0x3e, 0x67, 0x76, 0x8f, 0x61,
// 	0x7f, 0xc8, 0x1b, 0xc3, 0x88, 0x8a, 0x51, 0x32,
// 	0x3a, 0x9f, 0xb8, 0xaa, 0x4b, 0x1e, 0x5e, 0x4a,
// })

// // fakeRandReader implements the io.Reader interface and is used to force
// // errors in the RandomUint64 function.
// type fakeRandReader struct {
// 	n   int
// 	err error
// }

// // Read returns the fake reader error and the lesser of the fake reader value
// // and the length of p.
// func (r *fakeRandReader) Read(p []byte) (int, error) {
// 	n := r.n
// 	if n > len(p) {
// 		n = len(p)
// 	}
// 	return n, r.err
// }

// TestElementWire tests wire encode and decode for various element types.  This
// is mainly to test the "fast" paths in readElement and writeElement which use
// type assertions to avoid reflection when possible.
func TestElementWire(t *testing.T) {
	type writeElementReflect int32

	tests := []struct {
		in  interface{} // Value to encode
		buf []byte      // Wire encoding
	}{
		{int8(-128), []byte{0x80}},
		{int8(127), []byte{0x7f}},
		{uint8(1), []byte{0x01}},
		{int32(-1), []byte{0xff, 0xff, 0xff, 0xff}},
		{int32(1), []byte{0x01, 0x00, 0x00, 0x00}},
		{uint32(256), []byte{0x00, 0x01, 0x00, 0x00}},
		{int64(-65536), []byte{0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}},
		{int64(65536), []byte{0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{uint64(4294967296), []byte{0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00}},
		{&hash.Stamp{
			0x01, 0x02, 0x03, 0x04,
		},
			[]byte{
				0x01, 0x02, 0x03, 0x04,
			},
		},
		{&hash.Hash{
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
		{&crypto.Address{
			0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
			0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
			0x11, 0x12, 0x13, 0x14, 0x15,
		},
			[]byte{
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
				0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
				0x11, 0x12, 0x13, 0x14, 0x15,
			},
		},

		// Type not supported by the "fast" path and requires reflection.
		{
			writeElementReflect(1),
			[]byte{0x01, 0x00, 0x00, 0x00},
		},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Write to wire format.
		var buf bytes.Buffer
		err := WriteElement(&buf, test.in)
		if err != nil {
			t.Errorf("writeElement #%d error %v", i, err)
			continue
		}
		if !bytes.Equal(buf.Bytes(), test.buf) {
			t.Errorf("writeElement #%d\n got: %s want: %s", i,
				spew.Sdump(buf.Bytes()), spew.Sdump(test.buf))
			continue
		}

		// Read from wire format.
		rbuf := bytes.NewReader(test.buf)
		val := test.in
		if reflect.ValueOf(test.in).Kind() != reflect.Ptr {
			val = reflect.New(reflect.TypeOf(test.in)).Interface()
		}
		err = ReadElement(rbuf, val)
		if err != nil {
			t.Errorf("readElement #%d error %v", i, err)
			continue
		}
		ival := val
		if reflect.ValueOf(test.in).Kind() != reflect.Ptr {
			ival = reflect.Indirect(reflect.ValueOf(val)).Interface()
		}
		if !reflect.DeepEqual(ival, test.in) {
			t.Errorf("readElement #%d\n got: %s want: %s", i,
				spew.Sdump(ival), spew.Sdump(test.in))
			continue
		}
	}
}

// TestElementWireErrors performs negative tests against wire encode and decode
// of various element types to confirm error paths work correctly.
func TestElementWireErrors(t *testing.T) {
	tests := []struct {
		in       interface{} // Value to encode
		max      int         // Max size of fixed buffer to induce errors
		writeErr error       // Expected write error
		readErr  error       // Expected read error
	}{
		{int8(127), 0, io.ErrShortWrite, io.EOF},
		{uint8(1), 0, io.ErrShortWrite, io.EOF},
		{int32(1), 0, io.ErrShortWrite, io.EOF},
		{uint32(256), 0, io.ErrShortWrite, io.EOF},
		{int64(65536), 0, io.ErrShortWrite, io.EOF},
		{uint64(4294967296), 0, io.ErrShortWrite, io.EOF},
		{&hash.Stamp{0x01, 0x02, 0x03, 0x04}, 0, io.ErrShortWrite, io.EOF},
		{&hash.Hash{
			0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
			0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
			0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
			0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20,
		},
			0, io.ErrShortWrite, io.EOF,
		},
		{&crypto.Address{
			0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
			0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
			0x11, 0x12, 0x13, 0x14, 0x15,
		},
			0, io.ErrShortWrite, io.EOF,
		},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Encode to wire format.
		w := util.NewFixedWriter(test.max)
		err := WriteElement(w, test.in)
		if err != test.writeErr {
			t.Errorf("writeElement #%d wrong error got: %v, want: %v",
				i, err, test.writeErr)
			continue
		}

		// Decode from wire format.
		r := util.NewFixedReader(test.max, nil)
		val := test.in
		if reflect.ValueOf(test.in).Kind() != reflect.Ptr {
			val = reflect.New(reflect.TypeOf(test.in)).Interface()
		}
		err = ReadElement(r, val)
		if err != test.readErr {
			t.Errorf("readElement #%d wrong error got: %v, want: %v",
				i, err, test.readErr)
			continue
		}
	}
}

// TestVarStringWire tests wire encode and decode for variable length strings.
func TestVarStringWire(t *testing.T) {

	// str256 is a string that takes a 2-byte varint to encode.
	str256 := strings.Repeat("test", 64)

	tests := []struct {
		in  string // String to encode
		out string // String to decoded value
		buf []byte // Wire encoding
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
	for i, test := range tests {
		// Encode to wire format.
		var buf bytes.Buffer
		err := WriteVarString(&buf, test.in)
		if err != nil {
			t.Errorf("WriteVarString #%d error %v", i, err)
			continue
		}
		if !bytes.Equal(buf.Bytes(), test.buf) {
			t.Errorf("WriteVarString #%d\n got: %s want: %s", i,
				spew.Sdump(buf.Bytes()), spew.Sdump(test.buf))
			continue
		}

		// Decode from wire format.
		rbuf := bytes.NewReader(test.buf)
		val, err := ReadVarString(rbuf)
		if err != nil {
			t.Errorf("ReadVarString #%d error %v", i, err)
			continue
		}
		if val != test.out {
			t.Errorf("ReadVarString #%d\n got: %s want: %s", i,
				val, test.out)
			continue
		}
	}
}

// TestVarStringWireErrors performs negative tests against wire encode and
// decode of variable length strings to confirm error paths work correctly.
func TestVarStringWireErrors(t *testing.T) {
	// str256 is a string that takes a 2-byte varint to encode.
	str256 := strings.Repeat("test", 64)

	tests := []struct {
		in       string // Value to encode
		buf      []byte // Wire encoding
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
	for i, test := range tests {
		// Encode to wire format.
		w := util.NewFixedWriter(test.max)
		err := WriteVarString(w, test.in)
		if err != test.writeErr {
			t.Errorf("WriteVarString #%d wrong error got: %v, want: %v",
				i, err, test.writeErr)
			continue
		}

		// Decode from wire format.
		r := util.NewFixedReader(test.max, test.buf)
		_, err = ReadVarString(r)
		if err != test.readErr {
			t.Errorf("ReadVarString #%d wrong error got: %v, want: %v",
				i, err, test.readErr)
			continue
		}
	}
}

// TestVarStringOverflowErrors performs tests to ensure deserializing variable
// length strings intentionally crafted to use large values for the string
// length are handled properly.  This could otherwise potentially be used as an
// attack vector.
func TestVarStringOverflowErrors(t *testing.T) {
	tests := []struct {
		buf []byte // Wire encoding
	}{
		{[]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01}},
		{[]byte{0x80, 0x80, 0x80, 0x11}},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Decode from wire format.
		rbuf := bytes.NewReader(test.buf)
		_, err := ReadVarString(rbuf)
		if !strings.Contains(err.Error(), "variable length string is too long") {
			t.Errorf("ReadVarString #%d wrong error", i)
			continue
		}
	}

}

// TestVarBytesWire tests wire encode and decode for variable length byte array.
func TestVarBytesWire(t *testing.T) {
	// bytes256 is a byte array that takes a 2-byte varint to encode.
	bytes256 := bytes.Repeat([]byte{0x01}, 256)

	tests := []struct {
		in  []byte // Byte Array to write
		buf []byte // Wire encoding
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
	for i, test := range tests {
		// Encode to wire format.
		var buf bytes.Buffer
		err := WriteVarBytes(&buf, test.in)
		if err != nil {
			t.Errorf("WriteVarBytes #%d error %v", i, err)
			continue
		}
		if !bytes.Equal(buf.Bytes(), test.buf) {
			t.Errorf("WriteVarBytes #%d\n got: %s want: %s", i,
				spew.Sdump(buf.Bytes()), spew.Sdump(test.buf))
			continue
		}

		// Decode from wire format.
		rbuf := bytes.NewReader(test.buf)
		val, err := ReadVarBytes(rbuf)
		if err != nil {
			t.Errorf("ReadVarBytes #%d error %v", i, err)
			continue
		}
		if !bytes.Equal(buf.Bytes(), test.buf) {
			t.Errorf("ReadVarBytes #%d\n got: %s want: %s", i,
				val, test.buf)
			continue
		}
	}
}

// TestVarBytesWireErrors performs negative tests against wire encode and
// decode of variable length byte arrays to confirm error paths work correctly.
func TestVarBytesWireErrors(t *testing.T) {
	// bytes256 is a byte array that takes a 2-byte varint to encode.
	bytes256 := bytes.Repeat([]byte{0x01}, 256)

	tests := []struct {
		in       []byte // Byte Array to write
		buf      []byte // Wire encoding
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
	for i, test := range tests {
		// Encode to wire format.
		w := util.NewFixedWriter(test.max)
		err := WriteVarBytes(w, test.in)
		if err != test.writeErr {
			t.Errorf("WriteVarBytes #%d wrong error got: %v, want: %v",
				i, err, test.writeErr)
			continue
		}

		// Decode from wire format.
		r := util.NewFixedReader(test.max, test.buf)
		_, err = ReadVarBytes(r)
		if err != test.readErr {
			t.Errorf("ReadVarBytes #%d wrong error got: %v, want: %v",
				i, err, test.readErr)
			continue
		}
	}
}

// TestVarBytesOverflowErrors performs tests to ensure deserializing variable
// length byte arrays intentionally crafted to use large values for the array
// length are handled properly.  This could otherwise potentially be used as an
// attack vector.
func TestVarBytesOverflowErrors(t *testing.T) {
	tests := []struct {
		buf []byte // Wire encoding
	}{
		{[]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01}},
		{[]byte{0x80, 0x80, 0x80, 0x11}},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Decode from wire format.
		rbuf := bytes.NewReader(test.buf)
		_, err := ReadVarBytes(rbuf)
		if !strings.Contains(err.Error(), "variable length byte array is too long") {
			t.Errorf("ReadVarString #%d wrong error", i)
			continue
		}
	}

}

// // TestRandomUint64 exercises the randomness of the random number generator on
// // the system by ensuring the probability of the generated numbers.  If the RNG
// // is evenly distributed as a proper cryptographic RNG should be, there really
// // should only be 1 number < 2^56 in 2^8 tries for a 64-bit number.  However,
// // use a higher number of 5 to really ensure the test doesn't fail unless the
// // RNG is just horrendous.
// func TestRandomUint64(t *testing.T) {
// 	tries := 1 << 8              // 2^8
// 	watermark := uint64(1 << 56) // 2^56
// 	maxHits := 5
// 	badRNG := "The random number generator on this system is clearly " +
// 		"terrible since we got %d values less than %d in %d runs " +
// 		"when only %d was expected"

// 	numHits := 0
// 	for i := 0; i < tries; i++ {
// 		nonce, err := RandomUint64()
// 		if err != nil {
// 			t.Errorf("RandomUint64 iteration %d failed - err %v",
// 				i, err)
// 			return
// 		}
// 		if nonce < watermark {
// 			numHits++
// 		}
// 		if numHits > maxHits {
// 			str := fmt.Sprintf(badRNG, numHits, watermark, tries, maxHits)
// 			t.Errorf("Random Uint64 iteration %d failed - %v %v", i,
// 				str, numHits)
// 			return
// 		}
// 	}
// }

// // TestRandomUint64Errors uses a fake reader to force error paths to be executed
// // and checks the results accordingly.
// func TestRandomUint64Errors(t *testing.T) {
// 	// Test short reads.
// 	fr := &fakeRandReader{n: 2, err: io.EOF}
// 	nonce, err := randomUint64(fr)
// 	if err != io.ErrUnexpectedEOF {
// 		t.Errorf("Error not expected value of %v [%v]",
// 			io.ErrUnexpectedEOF, err)
// 	}
// 	if nonce != 0 {
// 		t.Errorf("Nonce is not 0 [%v]", nonce)
// 	}
// }

func TestVarInt(t *testing.T) {
	tests := []struct {
		in  uint64 // Value to encode
		buf []byte // Wire encoding
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
	}
	for _, test := range tests {
		var buf bytes.Buffer
		err := WriteVarInt(&buf, test.in)
		assert.NoError(t, err)
		assert.Equal(t, buf.Bytes(), test.buf, "invalid write for %x", test.in)

		val, err := ReadVarInt(&buf)
		assert.NoError(t, err)
		assert.Equal(t, val, test.in, "invalid read for %x", test.in)

		assert.Equal(t, VarIntSerializeSize(test.in), len(test.buf), "invalid size for %x", test.in)
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
	for i, test := range tests {
		// Decode from wire format.
		rbuf := bytes.NewReader(test.in)
		val, err := ReadVarInt(rbuf)
		if err != test.readErr {
			t.Errorf("ReadVarInt #%d unexpected error %v", i, err)
			continue
		}
		if val != 0 {
			t.Errorf("ReadVarInt #%d \n got: %d want: 0", i, val)
			continue
		}
	}
}

func TestVarIntOverflow(t *testing.T) {
	var buf bytes.Buffer
	buf.Write([]byte{0xff, 0x00})
	_, err := ReadVarInt(&buf)
	assert.Error(t, err)

	buf.Reset()
	buf.Write([]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff})
	_, err = ReadVarInt(&buf)
	assert.Error(t, err)
}
