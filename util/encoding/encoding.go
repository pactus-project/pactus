package encoding

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/pactus-project/pactus/crypto/hash"
)

const (
	// MaxPayloadSize is the maximum bytes a message can be regardless of other
	// individual limits imposed by messages themselves.
	MaxPayloadSize = 1024 * 1024 * 32 // 32MB
	// binaryFreeListMaxItems is the number of buffers to keep in the free
	// list to use for binary serialization and deserialization.
	binaryFreeListMaxItems = 1024
)

var (
	ErrOverflow     = errors.New("overflow")
	ErrNonCanonical = errors.New("non canonical")
)

// binaryFreeList defines a concurrent safe free list of byte slices (up to the
// maximum number defined by the binaryFreeListMaxItems constant) that have a
// cap of 8 (thus it supports up to a uint64).  It is used to provide temporary
// buffers for serializing and deserializing primitive numbers to and from their
// binary encoding in order to greatly reduce the number of allocations
// required.
//
// For convenience, functions are provided for each of the primitive unsigned
// integers that automatically obtain a buffer from the free list, perform the
// necessary binary conversion, read from or write to the given io.Reader or
// io.Writer, and return the buffer to the free list.
type binaryFreeList chan []byte

// Borrow returns a byte slice from the free list with a length of 8.  A new
// buffer is allocated if there are not any available on the free list.
func (l binaryFreeList) Borrow() []byte {
	var buf []byte
	select {
	case buf = <-l:
	default:
		buf = make([]byte, 8)
	}

	return buf[:8]
}

// Return puts the provided byte slice back on the free list.  The buffer MUST
// have been obtained via the Borrow function and therefore have a cap of 8.
func (l binaryFreeList) Return(buf []byte) {
	select {
	case l <- buf:
	default:
		// Let it go to the garbage collector.
	}
}

// Uint8 reads a single byte from the provided reader using a buffer from the
// free list.
func (l binaryFreeList) Uint8(r io.Reader, val *uint8) error {
	buf := l.Borrow()[:1]
	if _, err := io.ReadFull(r, buf); err != nil {
		l.Return(buf)

		return err
	}
	*val = buf[0]
	l.Return(buf)

	return nil
}

// Uint16 reads two bytes from the provided reader using a buffer from the
// free list, converts it to a number in little endian byte order.
func (l binaryFreeList) Uint16(r io.Reader, val *uint16) error {
	buf := l.Borrow()[:2]
	if _, err := io.ReadFull(r, buf); err != nil {
		l.Return(buf)

		return err
	}
	*val = binary.LittleEndian.Uint16(buf)
	l.Return(buf)

	return nil
}

// Uint32 reads four bytes from the provided reader using a buffer from the
// free list, converts it to a number in little endian byte order.
func (l binaryFreeList) Uint32(r io.Reader, val *uint32) error {
	buf := l.Borrow()[:4]
	if _, err := io.ReadFull(r, buf); err != nil {
		l.Return(buf)

		return err
	}
	*val = binary.LittleEndian.Uint32(buf)
	l.Return(buf)

	return nil
}

// Uint64 reads eight bytes from the provided reader using a buffer from the
// free list, converts it to a number in little endian byte order..
func (l binaryFreeList) Uint64(r io.Reader, val *uint64) error {
	buf := l.Borrow()[:8]
	if _, err := io.ReadFull(r, buf); err != nil {
		l.Return(buf)

		return err
	}
	*val = binary.LittleEndian.Uint64(buf)
	l.Return(buf)

	return nil
}

// PutUint8 copies the provided uint8 into a buffer from the free list and
// writes the resulting byte to the given writer.
func (l binaryFreeList) PutUint8(w io.Writer, val uint8) error {
	buf := l.Borrow()[:1]
	buf[0] = val
	_, err := w.Write(buf)
	l.Return(buf)

	return err
}

// PutUint16 serializes the provided uint16 using the given byte order into a
// buffer from the free list and writes the resulting two bytes to the given
// writer.
func (l binaryFreeList) PutUint16(w io.Writer, val uint16) error {
	buf := l.Borrow()[:2]
	binary.LittleEndian.PutUint16(buf, val)
	_, err := w.Write(buf)
	l.Return(buf)

	return err
}

// PutUint32 serializes the provided uint32 using the given byte order into a
// buffer from the free list and writes the resulting four bytes to the given
// writer.
func (l binaryFreeList) PutUint32(w io.Writer, val uint32) error {
	buf := l.Borrow()[:4]
	binary.LittleEndian.PutUint32(buf, val)
	_, err := w.Write(buf)
	l.Return(buf)

	return err
}

// PutUint64 serializes the provided uint64 using the given byte order into a
// buffer from the free list and writes the resulting eight bytes to the given
// writer.
func (l binaryFreeList) PutUint64(w io.Writer, val uint64) error {
	buf := l.Borrow()[:8]
	binary.LittleEndian.PutUint64(buf, val)
	_, err := w.Write(buf)
	l.Return(buf)

	return err
}

// binarySerializer provides a free list of buffers to use for serializing and
// deserializing primitive integer values to and from io.Readers and io.Writers.
var binarySerializer binaryFreeList = make(chan []byte, binaryFreeListMaxItems)

// ReadElement reads the next sequence of bytes from r using little endian
// depending on the concrete type of element pointed to.
func ReadElement(r io.Reader, element interface{}) error {
	// Attempt to read the element based on the concrete type via fast
	// type assertions first.
	var err error
	switch e := element.(type) {
	case *bool:
		rv := uint8(0)
		err = binarySerializer.Uint8(r, &rv)
		if rv == 0x00 {
			*e = false
		} else {
			*e = true
		}
	case *int8:
		rv := uint8(0)
		err = binarySerializer.Uint8(r, &rv)
		*e = int8(rv)
	case *uint8:
		err = binarySerializer.Uint8(r, e)
	case *int16:
		rv := uint16(0)
		err = binarySerializer.Uint16(r, &rv)
		*e = int16(rv)
	case *uint16:
		err = binarySerializer.Uint16(r, e)
	case *int32:
		rv := uint32(0)
		err = binarySerializer.Uint32(r, &rv)
		*e = int32(rv)
	case *uint32:
		err = binarySerializer.Uint32(r, e)
	case *int64:
		rv := uint64(0)
		err = binarySerializer.Uint64(r, &rv)
		*e = int64(rv)
	case *uint64:
		err = binarySerializer.Uint64(r, e)
	case *hash.Hash:
		_, err = io.ReadFull(r, e[:])
	default:
		// Fall back to the slower binary.Read if a fast path was not available
		// above.
		err = binary.Read(r, binary.LittleEndian, element)
	}

	return err
}

// ReadElements reads multiple items from r.  It is equivalent to multiple
// calls to readElement.
func ReadElements(r io.Reader, elements ...interface{}) error {
	for _, element := range elements {
		err := ReadElement(r, element)
		if err != nil {
			return err
		}
	}

	return nil
}

// WriteElement writes the little endian representation of element to w.
func WriteElement(w io.Writer, element interface{}) error {
	// Attempt to write the element based on the concrete type via fast
	// type assertions first.
	var err error
	switch e := element.(type) {
	case bool:
		if e {
			err = binarySerializer.PutUint8(w, 0x01)
		} else {
			err = binarySerializer.PutUint8(w, 0x00)
		}
	case int8:
		err = binarySerializer.PutUint8(w, uint8(e))
	case uint8:
		err = binarySerializer.PutUint8(w, e)
	case int16:
		err = binarySerializer.PutUint16(w, uint16(e))
	case uint16:
		err = binarySerializer.PutUint16(w, e)
	case int32:
		err = binarySerializer.PutUint32(w, uint32(e))
	case uint32:
		err = binarySerializer.PutUint32(w, e)
	case int64:
		err = binarySerializer.PutUint64(w, uint64(e))
	case uint64:
		err = binarySerializer.PutUint64(w, e)
	case *hash.Hash:
		_, err = w.Write(e[:])
	default:
		// Fall back to the slower binary.Write if a fast path was not available
		// above.
		err = binary.Write(w, binary.LittleEndian, element)
	}

	return err
}

// WriteElements writes multiple items to w.  It is equivalent to multiple
// calls to writeElement.
func WriteElements(w io.Writer, elements ...interface{}) error {
	for _, element := range elements {
		err := WriteElement(w, element)
		if err != nil {
			return err
		}
	}

	return nil
}

// ReadVarInt reads a variable length integer from r and returns it as a uint64.
func ReadVarInt(r io.Reader) (uint64, error) {
	bits := 64
	write := uint64(0)
	for shift := 0; ; shift += 7 {
		b := uint8(0)
		err := binarySerializer.Uint8(r, &b)
		if err != nil {
			return 0, err
		}
		if shift+7 >= bits && b >= 1<<(bits-shift) {
			return uint64(0), ErrOverflow
		}
		if b == 0 && shift != 0 {
			return uint64(0), ErrNonCanonical
		}

		write |= uint64(b&0x7f) << shift // Does the actually placing into write, stripping the first bit

		// If there is no next
		if (b & 0x80) == 0 {
			break
		}
	}

	return write, nil
}

// WriteVarInt serializes val to w using a variable number of bytes depending
// on its value.
func WriteVarInt(w io.Writer, val uint64) error {
	// Make sure that there is one after this
	for val >= 0x80 {
		n := (uint8(val) & 0x7f) | 0x80
		err := binarySerializer.PutUint8(w, n)
		if err != nil {
			return err
		}
		val >>= 7 // It should be in multiples of 7, this should just get the next part
	}

	return binarySerializer.PutUint8(w, uint8(val))
}

// VarIntSerializeSize returns the number of bytes it would take to serialize
// val as a variable length integer.
func VarIntSerializeSize(val uint64) int {
	switch {
	case val >= 0x8000000000000000:
		return 10
	case val >= 0x100000000000000:
		return 9
	case val >= 0x2000000000000:
		return 8
	case val >= 0x40000000000:
		return 7
	case val >= 0x800000000:
		return 6
	case val >= 0x10000000:
		return 5
	case val >= 0x200000:
		return 4
	case val >= 0x4000:
		return 3
	case val >= 0x80:
		return 2
	default:
		return 1
	}
}

// VarStringSerializeSize returns the number of bytes it would take to serialize
// val as a string.
func VarStringSerializeSize(str string) int {
	return VarIntSerializeSize(uint64(len(str))) + len(str)
}

// VarBytesSerializeSize returns the number of bytes it would take to serialize
// val as a byte array.
func VarBytesSerializeSize(bytes []byte) int {
	return VarIntSerializeSize(uint64(len(bytes))) + len(bytes)
}

// ReadVarString reads a variable length string from r and returns it as a Go
// string. A variable length string is encoded as a variable length integer
// containing the length of the string followed by the bytes that represent the
// string itself. An error is returned if the length is greater than the
// maximum payload size since it helps protect against memory exhaustion
// attacks and forced panics through malformed messages.
func ReadVarString(r io.Reader) (string, error) {
	count, err := ReadVarInt(r)
	if err != nil {
		return "", err
	}

	// Prevent variable length strings that are larger than the maximum
	// payload size.  It would be possible to cause memory exhaustion and
	// panics without a sane upper bound on this count.
	if count > MaxPayloadSize {
		return "", fmt.Errorf("variable length string is too long "+
			"[count %d, max %d]", count, MaxPayloadSize)
	}

	buf := make([]byte, count)
	_, err = io.ReadFull(r, buf)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

// WriteVarString serializes str to w as a variable length integer containing
// the length of the string followed by the bytes that represent the string
// itself.
func WriteVarString(w io.Writer, str string) error {
	err := WriteVarInt(w, uint64(len(str)))
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(str))

	return err
}

// ReadVarBytes reads a variable length byte array.  A byte array is encoded
// as a varInt containing the length of the array followed by the bytes
// themselves.  An error is returned if the length is greater than the
// maximum payload size since it helps protect against memory exhaustion
// attacks and forced panics through malformed messages.
func ReadVarBytes(r io.Reader) ([]byte, error) {
	count, err := ReadVarInt(r)
	if err != nil {
		return nil, err
	}

	// Prevent byte array larger than the max message size. It would
	// be possible to cause memory exhaustion and panics without a sane
	// upper bound on this count.
	if count > uint64(MaxPayloadSize) {
		return nil, fmt.Errorf("variable length byte array is too long "+
			"[count %d, max %d]", count, MaxPayloadSize)
	}

	b := make([]byte, count)
	_, err = io.ReadFull(r, b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// WriteVarBytes serializes a variable length byte array to w as a varInt
// containing the number of bytes, followed by the bytes themselves.
func WriteVarBytes(w io.Writer, bytes []byte) error {
	slen := uint64(len(bytes))
	err := WriteVarInt(w, slen)
	if err != nil {
		return err
	}

	_, err = w.Write(bytes)

	return err
}
