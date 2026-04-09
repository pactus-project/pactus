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
	"github.com/stretchr/testify/require"
)

func TestWithdrawType(t *testing.T) {
	pld := payload.WithdrawPayload{}
	assert.Equal(t, payload.TypeWithdraw, pld.Type())
}

func TestWithdrawString(t *testing.T) {
	pld := payload.WithdrawPayload{}
	assert.Contains(t, pld.LogString(), "{Withdraw ")
}

func TestWithdrawDecoding(t *testing.T) {
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
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
				0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10,
				0x11, 0x12, 0x13, 0x14, // from
			},
			value:    0,
			readErr:  io.ErrUnexpectedEOF,
			basicErr: nil,
		},
		{
			raw: []byte{
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
				0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10,
				0x11, 0x12, 0x13, 0x14, 0x15, // from
				0x02, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
				0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F, 0x20,
				0x21, 0x12, 0x23, 0x24, // to
			},
			value:    0,
			readErr:  io.ErrUnexpectedEOF,
			basicErr: nil,
		},
		{
			raw: []byte{
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
				0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10,
				0x11, 0x12, 0x13, 0x14, 0x15, // from
				0x02, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
				0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F, 0x20,
				0x21, 0x12, 0x23, 0x24, 0x25, // to
				0x80, 0x80, 0x80, // amount
			},
			value:    0,
			readErr:  io.EOF,
			basicErr: nil,
		},
		{
			raw: []byte{
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
				0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10,
				0x11, 0x12, 0x13, 0x14, 0x15, // from
				0x02, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
				0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F, 0x20,
				0x21, 0x12, 0x23, 0x24, 0x25, // to
				0x80, 0x80, 0x80, 0x01, // amount
			},
			value:    0x200000,
			readErr:  nil,
			basicErr: nil,
		},
		{
			raw: []byte{
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
				0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10,
				0x11, 0x12, 0x13, 0x14, 0x15, // from
				0x00,                   // to (treasury address)
				0x80, 0x80, 0x80, 0x01, // amount
			},
			value:    0x200000,
			readErr:  nil,
			basicErr: nil,
		},
	}

	for no, tt := range tests {
		pld := payload.WithdrawPayload{}
		r := util.NewFixedReader(len(tt.raw), tt.raw)
		err := pld.Decode(payload.DecodeContext{}, r)
		if tt.readErr != nil {
			require.ErrorIs(t, err, tt.readErr)
		} else {
			require.NoError(t, err)

			for i := 0; i < pld.SerializeSize(); i++ {
				w := util.NewFixedWriter(i)
				require.Error(t, pld.Encode(w), "encode test %v failed", no)
			}
			w := util.NewFixedWriter(pld.SerializeSize())
			require.NoError(t, pld.Encode(w))
			assert.Len(t, w.Bytes(), pld.SerializeSize())
			assert.Equal(t, tt.raw, w.Bytes())

			// Basic check
			if tt.basicErr != nil {
				require.ErrorIs(t, pld.BasicCheck(), tt.basicErr)
			} else {
				require.NoError(t, pld.BasicCheck())

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

func TestWithdrawBasicCheck(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	tests := []struct {
		pld payload.WithdrawPayload
		err string
	}{
		{
			pld: payload.WithdrawPayload{
				From: ts.RandAccAddress(),
				To:   ts.RandAccAddress(),
			},
			err: "sender is not a validator address",
		},
		{
			pld: payload.WithdrawPayload{
				From: ts.RandValAddress(),
				To:   ts.RandValAddress(),
			},
			err: "receiver is not an account address",
		},
	}

	for no, tt := range tests {
		assert.ErrorContains(t, tt.pld.BasicCheck(), tt.err, "test %v failed", no)
	}
}
