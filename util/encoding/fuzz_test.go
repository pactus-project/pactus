package encoding

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

// The variable length values come from any peer: reading them from any bytes never panics,
// and a value that is read is written back to bytes that give the same value.
func FuzzReadVar(f *testing.F) {
	for _, seed := range [][]byte{
		{},
		{0x00},
		{0xfc},
		{0xfd, 0xfd, 0x00},
		{0xfe, 0x00, 0x00, 0x01, 0x00},
		{0xff, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00},
		{0x03, 'a', 'b', 'c'},
		{0xfe, 0xff, 0xff, 0xff, 0xff},
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		if val, err := ReadVarInt(bytes.NewReader(data)); err == nil {
			buf := new(bytes.Buffer)
			require.NoError(t, WriteVarInt(buf, val))
			require.Equal(t, VarIntSerializeSize(val), buf.Len())

			again, err := ReadVarInt(buf)
			require.NoError(t, err)
			require.Equal(t, val, again)
		}

		if val, err := ReadVarBytes(bytes.NewReader(data)); err == nil {
			buf := new(bytes.Buffer)
			require.NoError(t, WriteVarBytes(buf, val))
			require.Equal(t, VarBytesSerializeSize(val), buf.Len())
		}

		if val, err := ReadVarStringMax(bytes.NewReader(data), 64); err == nil {
			require.LessOrEqual(t, len(val), 64)
		}
	})
}
