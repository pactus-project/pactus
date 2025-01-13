package zmq

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeTopicMsg(t *testing.T) {
	td := setup(t)

	_, addr := td.TestSuite.GenerateTestAccount()
	randHeight := td.TestSuite.RandHeight()

	tests := []struct {
		name      string
		parts     []any
		want      []byte
		wantPanic bool
	}{
		{
			name:  "single crypto.Address",
			parts: []any{addr},
			want:  addr.Bytes(),
		},
		{
			name:  "single Topic",
			parts: []any{BlockInfo},
			want:  BlockInfo.Bytes(),
		},
		{
			name:  "uint32 value",
			parts: []any{randHeight},
			want:  binary.BigEndian.AppendUint32([]byte{}, randHeight),
		},
		{
			name:  "uint16 value",
			parts: []any{uint16(0x0506)},
			want:  binary.BigEndian.AppendUint16([]byte{}, 0x0506),
		},
		{
			name:  "multiple types",
			parts: []any{addr, BlockInfo, uint32(0x0A0B0C0D), []byte{0x0E}},
			want: func() []byte {
				b := addr.Bytes()
				b = append(b, BlockInfo.Bytes()...)
				b = binary.BigEndian.AppendUint32(b, 0x0A0B0C0D)
				b = append(b, 0x0E)

				return b
			}(),
		},
		{
			name:      "unknown type",
			parts:     []any{"unknown"},
			wantPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				assert.Panics(t, func() {
					makeTopicMsg(tt.parts...)
				})

				return
			}
			got := makeTopicMsg(tt.parts...)
			assert.Equal(t, tt.want, got)
		})
	}
}
