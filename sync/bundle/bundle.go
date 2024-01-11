package bundle

import (
	"fmt"
	"io"

	"github.com/fxamacker/cbor/v2"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/errors"
)

const (
	BundleFlagNetworkMainnet = 0x0001
	BundleFlagNetworkTestnet = 0x0002
	BundleFlagCarrierLibP2P  = 0x0010
	BundleFlagCompressed     = 0x0100
	BundleFlagBroadcasted    = 0x0200
	BundleFlagHandshaking    = 0x0400
)

type Bundle struct {
	Flags      int
	SequenceNo int
	Message    message.Message
}

func NewBundle(msg message.Message) *Bundle {
	return &Bundle{
		Flags:   0,
		Message: msg,
	}
}

func (b *Bundle) BasicCheck() error {
	return b.Message.BasicCheck()
}

func (b *Bundle) String() string {
	return fmt.Sprintf("%s%s", b.Message.Type(), b.Message.String())
}

func (b *Bundle) CompressIt() {
	b.Flags = util.SetFlag(b.Flags, BundleFlagCompressed)
}

func (b *Bundle) SetSequenceNo(seqNo int) {
	b.SequenceNo = seqNo
}

type _Bundle struct {
	Flags       int          `cbor:"1,keyasint"`
	MessageType message.Type `cbor:"2,keyasint"`
	MessageData []byte       `cbor:"3,keyasint"`
	SequenceNo  int          `cbor:"4,keyasint"`
}

func (b *Bundle) Encode() ([]byte, error) {
	data, err := cbor.Marshal(b.Message)
	if err != nil {
		return nil, err
	}

	if util.IsFlagSet(b.Flags, BundleFlagCompressed) {
		c, err := util.CompressBuffer(data)
		if err == nil {
			data = c
		}
	}

	msg := &_Bundle{
		Flags:       b.Flags,
		MessageType: b.Message.Type(),
		MessageData: data,
		SequenceNo:  b.SequenceNo,
	}

	return cbor.Marshal(msg)
}

func (b *Bundle) Decode(r io.Reader) (int, error) {
	var bdl _Bundle
	d := cbor.NewDecoder(r)
	err := d.Decode(&bdl)
	bytesRead := d.NumBytesRead()
	if err != nil {
		return bytesRead, errors.Errorf(errors.ErrInvalidMessage, err.Error())
	}

	data := bdl.MessageData
	msg := message.MakeMessage(bdl.MessageType)
	if msg == nil {
		return bytesRead, errors.Errorf(errors.ErrInvalidMessage, "invalid data")
	}

	if util.IsFlagSet(bdl.Flags, BundleFlagCompressed) {
		c, err := util.DecompressBuffer(bdl.MessageData)
		if err != nil {
			return bytesRead, errors.Errorf(errors.ErrInvalidMessage, err.Error())
		}
		data = c
	}

	b.Flags = bdl.Flags
	b.SequenceNo = bdl.SequenceNo
	b.Message = msg

	return bytesRead, cbor.Unmarshal(data, msg)
}
