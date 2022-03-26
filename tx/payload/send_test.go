package payload

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/util"
)

func TestSendType(t *testing.T) {
	pld := SendPayload{}
	assert.Equal(t, pld.Type(), PayloadTypeSend)
}

func TestSendSanityCheck(t *testing.T) {
	invAddr := crypto.Address{0x02, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	pld := SendPayload{
		Sender:   invAddr,
		Receiver: invAddr,
		Amount:   1,
	}
	assert.Equal(t, errors.Code(pld.SanityCheck()), errors.ErrInvalidAddress)
	pld.Sender = crypto.GenerateTestAddress()
	assert.Equal(t, errors.Code(pld.SanityCheck()), errors.ErrInvalidAddress)
	pld.Receiver = crypto.GenerateTestAddress()
	assert.NoError(t, pld.SanityCheck())
}

func TestSendDecoding(t *testing.T) {
	bs := []byte{
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A,
		0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, // sender
		0x01, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1A,
		0x1B, 0x1C, 0x1D, 0x1E, 0x1F, 0x20, 0x21, 0x12, 0x23, 0x24, 0x25, // receiver
		0x80, 0x01, // amount
	}
	pld := SendPayload{}
	maxLen := []int{20, 40, 43}
	for i, len := range maxLen {
		r := util.NewFixedReader(len, bs)
		assert.Error(t, pld.Decode(r), "decode test %v failed", i)
	}

	r := util.NewFixedReader(44, bs)
	assert.NoError(t, pld.Decode(r))
	for i, len := range maxLen {
		w := util.NewFixedWriter(len)
		assert.Error(t, pld.Encode(w), "encode test %v failed", i)
	}
	w := util.NewFixedWriter(len(bs))
	assert.NoError(t, pld.Encode(w))
	assert.Equal(t, len(w.Bytes()), pld.SerializeSize())

	addr, _ := crypto.AddressFromBytes(bs[:21])
	assert.Equal(t, pld.Signer(), addr)
	assert.Equal(t, pld.Value(), int64(0x80))

	// covering fingerprint
	assert.NotEmpty(t, pld.Fingerprint())
}
