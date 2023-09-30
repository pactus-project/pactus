package payload

import (
	"io"
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBondType(t *testing.T) {
	pld := BondPayload{}
	assert.Equal(t, pld.Type(), TypeBond)
}

func TestBondDecoding(t *testing.T) {
	tests := []struct {
		raw      []byte
		value    int64
		readErr  error
		basicErr error
	}{
		{
			raw:      []byte{},
			value:    0,
			readErr:  io.EOF,
			basicErr: nil,
		},
		{
			raw: []byte{
				0x02, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
				0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10,
				0x11, 0x12, 0x13, 0x14, // sender
			},
			value:    0,
			readErr:  io.ErrUnexpectedEOF,
			basicErr: nil,
		},
		{
			raw: []byte{
				0x02, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
				0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10,
				0x11, 0x12, 0x13, 0x14, 0x15, // sender
				0x01, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
				0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F, 0x20,
				0x21, 0x12, 0x23, 0x24, // receiver
			},
			value:    0,
			readErr:  io.ErrUnexpectedEOF,
			basicErr: nil,
		},
		{
			raw: []byte{
				0x02, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
				0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10,
				0x11, 0x12, 0x13, 0x14, 0x15, // sender
				0x01, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
				0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F, 0x20,
				0x21, 0x12, 0x23, 0x24, 0x25, // receiver
			},
			value:    0,
			readErr:  io.EOF,
			basicErr: nil,
		},
		{
			raw: []byte{
				0x02, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
				0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10,
				0x11, 0x12, 0x13, 0x14, 0x15, // sender
				0x01, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
				0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F, 0x20,
				0x21, 0x12, 0x23, 0x24, 0x25, // receiver
				0x00,             // public key size
				0x80, 0x80, 0x80, // stake
			},
			value:    0,
			readErr:  io.EOF,
			basicErr: nil,
		},
		{
			raw: []byte{
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
				0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10,
				0x11, 0x12, 0x13, 0x14, 0x15, // sender
				0x01, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
				0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F, 0x20,
				0x21, 0x12, 0x23, 0x24, 0x25, // receiver
				0x00,                   // public key size
				0x80, 0x80, 0x80, 0x01, // stake
			},
			value:   0x200000,
			readErr: nil,
			basicErr: BasicCheckError{
				Reason: "sender is not an account address: pc1pqgpsgpgx",
			},
		},
		{
			raw: []byte{
				0x02, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
				0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10,
				0x11, 0x12, 0x13, 0x14, 0x15, // sender
				0x02, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
				0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F, 0x20,
				0x21, 0x12, 0x23, 0x24, 0x25, // receiver
				0x00,                   // public key size
				0x80, 0x80, 0x80, 0x01, // stake
			},
			value:   0x200000,
			readErr: nil,
			basicErr: BasicCheckError{
				Reason: "receiver is not a validator address: pc1zzgf3g9gk",
			},
		},
		{
			raw: []byte{
				0x02, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
				0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10,
				0x11, 0x12, 0x13, 0x14, 0x15, // sender
				0x01, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
				0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F, 0x20,
				0x21, 0x12, 0x23, 0x24, 0x25, // receiver
				0x00,                   // public key size
				0x80, 0x80, 0x80, 0x01, // stake
			},
			value:    0x200000,
			readErr:  nil,
			basicErr: nil,
		},
		{
			raw: []byte{
				0x02, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
				0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10,
				0x11, 0x12, 0x13, 0x14, 0x15, // sender
				0x01, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
				0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F, 0x20,
				0x21, 0x12, 0x23, 0x24, 0x25, // receiver
				0x00, // public key size
				0x00, // stake is zero
			},
			value:    0x0,
			readErr:  nil,
			basicErr: nil,
		},
		{
			raw: []byte{
				0x02, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
				0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10,
				0x11, 0x12, 0x13, 0x14, 0x15, // sender
				0x01, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
				0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F, 0x20,
				0x21, 0x12, 0x23, 0x24, 0x25, // receiver
				0x01, // public key size
			},
			value:    0x200000,
			readErr:  ErrInvalidPublicKeySize,
			basicErr: nil,
		},
		{
			raw: []byte{
				0x02, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
				0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10,
				0x11, 0x12, 0x13, 0x14, 0x15, // sender
				0x01, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
				0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F, 0x20,
				0x21, 0x12, 0x23, 0x24, 0x25, // receiver
				0x60, // public key size
				0xaf, // public key
			},
			value:    0x200000,
			readErr:  io.ErrUnexpectedEOF,
			basicErr: nil,
		},
		{
			raw: []byte{
				0x02, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
				0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10,
				0x11, 0x12, 0x13, 0x14, 0x15, // sender
				0x01, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
				0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F, 0x20,
				0x21, 0x12, 0x23, 0x24, 0x25, // receiver
				0x60, // public key size
				0xaf, 0x0f, 0x74, 0x91, 0x7f, 0x50, 0x65, 0xaf, 0x94, 0x72,
				0x7a, 0xe9, 0x54, 0x1b, 0x0d, 0xdc, 0xfb, 0x5b, 0x82, 0x8a,
				0x9e, 0x01, 0x6b, 0x02, 0x49, 0x8f, 0x47, 0x7e, 0xd3, 0x7f,
				0xb4, 0x4d, 0x5d, 0x88, 0x24, 0x95, 0xaf, 0xb6, 0xfd, 0x4f,
				0x97, 0x73, 0xe4, 0xea, 0x9d, 0xee, 0xe4, 0x36, 0x03, 0x0c,
				0x4d, 0x61, 0xc6, 0xe3, 0xa1, 0x15, 0x15, 0x85, 0xe1, 0xd8,
				0x38, 0xca, 0xe1, 0x44, 0x4a, 0x43, 0x8d, 0x08, 0x9c, 0xe7,
				0x7e, 0x10, 0xc4, 0x92, 0xa5, 0x5f, 0x69, 0x08, 0x12, 0x5c,
				0x5b, 0xe9, 0xb2, 0x36, 0xa2, 0x46, 0xe4, 0x08, 0x2d, 0x08,
				0xde, 0x56, 0x4e, 0x11, 0x1e, 0x65, // public key
				0x80, 0x80, 0x80, 0x01, // stake
			},
			value:   0x200000,
			readErr: nil,
			basicErr: crypto.AddressMismatchError{
				Expected: crypto.Address{
					0x01, 0xa1, 0x95, 0xd7, 0xfe, 0xcb, 0xa4, 0xc6,
					0x36, 0x83, 0x2f, 0x1d, 0xb0, 0xcd, 0x0e, 0xa1,
					0x4d, 0xb6, 0xdb, 0x8c, 0x71,
				},
				Got: crypto.Address{
					0x01, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
					0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F, 0x20,
					0x21, 0x12, 0x23, 0x24, 0x25,
				},
			},
		},
		{
			raw: []byte{
				0x02, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
				0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10,
				0x11, 0x12, 0x13, 0x14, 0x15, // sender
				0x01, 0xa1, 0x95, 0xd7, 0xfe, 0xcb, 0xa4, 0xc6,
				0x36, 0x83, 0x2f, 0x1d, 0xb0, 0xcd, 0x0e, 0xa1,
				0x4d, 0xb6, 0xdb, 0x8c, 0x71, // receiver
				0x60, // public key size
				0xaf, 0x0f, 0x74, 0x91, 0x7f, 0x50, 0x65, 0xaf, 0x94, 0x72,
				0x7a, 0xe9, 0x54, 0x1b, 0x0d, 0xdc, 0xfb, 0x5b, 0x82, 0x8a,
				0x9e, 0x01, 0x6b, 0x02, 0x49, 0x8f, 0x47, 0x7e, 0xd3, 0x7f,
				0xb4, 0x4d, 0x5d, 0x88, 0x24, 0x95, 0xaf, 0xb6, 0xfd, 0x4f,
				0x97, 0x73, 0xe4, 0xea, 0x9d, 0xee, 0xe4, 0x36, 0x03, 0x0c,
				0x4d, 0x61, 0xc6, 0xe3, 0xa1, 0x15, 0x15, 0x85, 0xe1, 0xd8,
				0x38, 0xca, 0xe1, 0x44, 0x4a, 0x43, 0x8d, 0x08, 0x9c, 0xe7,
				0x7e, 0x10, 0xc4, 0x92, 0xa5, 0x5f, 0x69, 0x08, 0x12, 0x5c,
				0x5b, 0xe9, 0xb2, 0x36, 0xa2, 0x46, 0xe4, 0x08, 0x2d, 0x08,
				0xde, 0x56, 0x4e, 0x11, 0x1e, 0x65, // public key
				0x80, 0x80, 0x80, 0x01, // stake
			},
			value:    0x200000,
			readErr:  nil,
			basicErr: nil,
		},
	}

	for n, test := range tests {
		pld := BondPayload{}
		r := util.NewFixedReader(len(test.raw), test.raw)
		err := pld.Decode(r)
		if test.readErr != nil {
			assert.ErrorIs(t, err, test.readErr)
		} else {
			assert.NoError(t, err)

			for i := 0; i < pld.SerializeSize(); i++ {
				w := util.NewFixedWriter(i)
				require.Error(t, pld.Encode(w), "encode %v failed", n)
			}
			w := util.NewFixedWriter(pld.SerializeSize())
			require.NoError(t, pld.Encode(w))
			assert.Equal(t, len(w.Bytes()), pld.SerializeSize())
			assert.Equal(t, w.Bytes(), test.raw)

			// Basic check
			if test.basicErr != nil {
				err := pld.BasicCheck()
				require.ErrorIs(t, err, test.basicErr, "basic check %v failed", n)
			} else {
				assert.NoError(t, pld.BasicCheck())

				// Check signer
				assert.Equal(t, pld.Signer(), crypto.Address(test.raw[:21]))
				assert.Equal(t, *pld.Receiver(), crypto.Address(test.raw[21:42]))
				assert.Equal(t, pld.Value(), test.value)
			}
		}
	}
}
