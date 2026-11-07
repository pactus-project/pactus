package encoding

import (
	"bytes"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ReadVarString and ReadVarBytes only reject a length greater
// than MaxPayloadSize (32 MiB); they never compare it to the bytes the reader
// actually has left. A tiny wire message — a few varint bytes announcing a
// 32 MiB length, with no payload behind it at all — makes make([]byte, count)
// allocate the full 32 MiB before io.ReadFull ever fails on the short read.
// Every such malformed message (e.g. a transaction memo field) costs the
// receiver a 32 MiB allocation for a handful of bytes on the wire.
//
// This test expresses the desired behavior (reject before the large
// allocation, or allocate incrementally) and failed before the fix, because the
// allocation happens regardless of what the reader actually contains.
func TestReadVarStringAllocatesFullLengthBeforeReadFails(t *testing.T) {
	const announcedLen = 32 * 1024 * 1024 // == MaxPayloadSize, so the length check alone lets it through.
	const acceptableAllocationForATinyMessage = 1024

	buf := new(bytes.Buffer)
	require.NoError(t, WriteVarInt(buf, uint64(announcedLen)))
	wireSize := buf.Len() // a handful of bytes: the varint-encoded length, nothing else.

	var before, after runtime.MemStats
	runtime.GC() //nolint:revive // the test measures the memory that a decoding allocates.
	runtime.ReadMemStats(&before)

	_, err := ReadVarString(bytes.NewReader(buf.Bytes()))

	runtime.ReadMemStats(&after)
	allocated := after.TotalAlloc - before.TotalAlloc

	require.Error(t, err, "the reader has no payload at all, so the read must fail")
	t.Logf("a %d-byte message announcing a %d-byte string made ReadVarString allocate %d bytes before failing",
		wireSize, announcedLen, allocated)

	assert.Less(t, allocated, uint64(acceptableAllocationForATinyMessage),
		"ReadVarString must not allocate close to the announced length before it has confirmed that much data exists")
}

func TestReadVarBytesAllocatesFullLengthBeforeReadFails(t *testing.T) {
	const announcedLen = 32 * 1024 * 1024
	const acceptableAllocationForATinyMessage = 1024

	buf := new(bytes.Buffer)
	require.NoError(t, WriteVarInt(buf, uint64(announcedLen)))

	var before, after runtime.MemStats
	runtime.GC() //nolint:revive // the test measures the memory that a decoding allocates.
	runtime.ReadMemStats(&before)

	_, err := ReadVarBytes(bytes.NewReader(buf.Bytes()))

	runtime.ReadMemStats(&after)
	allocated := after.TotalAlloc - before.TotalAlloc

	require.Error(t, err)
	t.Logf("ReadVarBytes allocated %d bytes before failing on a message with no payload", allocated)

	assert.Less(t, allocated, uint64(acceptableAllocationForATinyMessage),
		"ReadVarBytes must not allocate close to the announced length before it has confirmed that much data exists")
}
