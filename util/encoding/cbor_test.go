package encoding

import (
	"bytes"
	"testing"

	"github.com/fxamacker/cbor/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalCBORStrict(t *testing.T) {
	type item struct {
		Num  int    `cbor:"1,keyasint"`
		Data []byte `cbor:"2,keyasint"`
	}

	want := item{Num: 7, Data: []byte{0xaa, 0xbb}}
	canonical, err := cbor.Marshal(want)
	require.NoError(t, err)
	require.Equal(t, []byte{0xa2, 0x01, 0x07, 0x02, 0x42, 0xaa, 0xbb}, canonical)

	t.Run("the encoding of the encoder is accepted", func(t *testing.T) {
		var got item
		require.NoError(t, UnmarshalCBOR(canonical, &got))
		assert.Equal(t, want, got)

		var fromDecoder item
		require.NoError(t, NewCBORDecoder(bytes.NewReader(canonical)).Decode(&fromDecoder))
		assert.Equal(t, want, fromDecoder)
	})

	// Each of these is another encoding of the same value, that the default decoder accepts.
	variants := map[string][]byte{
		"map of indefinite length":         {0xbf, 0x01, 0x07, 0x02, 0x42, 0xaa, 0xbb, 0xff},
		"byte string of indefinite length": {0xa2, 0x01, 0x07, 0x02, 0x5f, 0x41, 0xaa, 0x41, 0xbb, 0xff},
		"duplicate map key":                {0xa3, 0x01, 0x07, 0x02, 0x42, 0xaa, 0xbb, 0x01, 0x07},
		"tagged byte string":               {0xa2, 0x01, 0x07, 0x02, 0xd8, 0x40, 0x42, 0xaa, 0xbb},
	}
	for name, data := range variants {
		t.Run(name, func(t *testing.T) {
			var lenient item
			require.NoError(t, cbor.Unmarshal(data, &lenient), "the variant is valid CBOR")

			var got item
			require.Error(t, UnmarshalCBOR(data, &got))
			require.Error(t, NewCBORDecoder(bytes.NewReader(data)).Decode(&got))
		})
	}
}
