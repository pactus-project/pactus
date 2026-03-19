package payload_test

import (
	"io"
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestUnbondType(t *testing.T) {
	pld := payload.UnbondPayload{}
	assert.Equal(t, payload.TypeUnbond, pld.Type())
}

func TestUnbondString(t *testing.T) {
	pld := payload.UnbondPayload{}
	assert.Contains(t, pld.LogString(), "{Unbond ")
}

func TestUnbondDecoding(t *testing.T) {
	tests := []struct {
		raw     []byte
		value   amount.Amount
		readErr error
	}{
		{
			raw:     []byte{},
			value:   0,
			readErr: io.EOF,
		},
		{
			raw: []byte{
				0x01, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
				0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F, 0x20,
				0x21, 0x12, 0x23, 0x24,
			},
			value:   0,
			readErr: io.ErrUnexpectedEOF,
		},
		{
			raw: []byte{
				0x01, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
				0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F, 0x20,
				0x21, 0x12, 0x23, 0x24, 0x25, // validator
			},
			value:   0,
			readErr: nil,
		},
	}

	for no, tt := range tests {
		pld := payload.UnbondPayload{}
		r := util.NewFixedReader(len(tt.raw), tt.raw)
		err := pld.Decode(payload.DecodeContext{}, r)
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

func TestUnbondBasicCheck(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	tests := []struct {
		pld payload.UnbondPayload
		err string
	}{
		{
			pld: payload.UnbondPayload{
				Validator: ts.RandAccAddress(),
			},
			err: "address is not a validator address",
		},
	}

	for no, tt := range tests {
		assert.ErrorContains(t, tt.pld.BasicCheck(), tt.err, "test %v failed", no)
	}
}
