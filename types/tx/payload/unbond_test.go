package payload

import (
	"io"
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/util"
	"github.com/stretchr/testify/assert"
)

func TestUnbondType(t *testing.T) {
	pld := UnbondPayload{}
	assert.Equal(t, TypeUnbond, pld.Type())
}

func TestUnbondString(t *testing.T) {
	pld := UnbondPayload{}
	assert.Contains(t, pld.LogString(), "{Unbond ")
}

func TestUnbondDecoding(t *testing.T) {
	tests := []struct {
		raw      []byte
		value    amount.Amount
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
				0x01, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
				0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F, 0x20,
				0x21, 0x12, 0x23, 0x24,
			},
			value:    0,
			readErr:  io.ErrUnexpectedEOF,
			basicErr: nil,
		},
		{
			raw: []byte{
				0x02, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
				0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F, 0x20,
				0x21, 0x12, 0x23, 0x24, 0x25, // account as validator address
			},
			value:   0,
			readErr: nil,
			basicErr: BasicCheckError{
				Reason: "address is not a validator address",
			},
		},
		{
			raw: []byte{
				0x01, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
				0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F, 0x20,
				0x21, 0x12, 0x23, 0x24, 0x25, // validator
			},
			value:    0,
			readErr:  nil,
			basicErr: nil,
		},
	}

	for no, tt := range tests {
		pld := UnbondPayload{}
		r := util.NewFixedReader(len(tt.raw), tt.raw)
		err := pld.Decode(r)
		if tt.readErr != nil {
			assert.ErrorIs(t, err, tt.readErr)
		} else {
			assert.NoError(t, err)

			for i := 0; i < pld.SerializeSize(); i++ {
				w := util.NewFixedWriter(i)
				assert.Error(t, pld.Encode(w), "encode test %v failed", no)
			}
			w := util.NewFixedWriter(pld.SerializeSize())
			assert.NoError(t, pld.Encode(w))
			assert.Equal(t, pld.SerializeSize(), len(w.Bytes()))
			assert.Equal(t, tt.raw, w.Bytes())

			// Basic check
			if tt.basicErr != nil {
				assert.ErrorIs(t, pld.BasicCheck(), tt.basicErr)
			} else {
				assert.NoError(t, pld.BasicCheck())

				// Check signer
				if tt.raw[0] != 0 {
					assert.Equal(t, crypto.Address(tt.raw[:21]), pld.Signer())
				} else {
					assert.Equal(t, crypto.TreasuryAddress, pld.Signer())
				}
				assert.Equal(t, tt.value, pld.Value())
			}
		}
	}
}
