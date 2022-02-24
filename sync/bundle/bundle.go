package bundle

import (
	"fmt"
	"io"

	"github.com/fxamacker/cbor/v2"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/sync/bundle/message"
	"github.com/zarbchain/zarb-go/util"
)

const LastVersion = 1
const (
	BundleFlagNetworkLibP2P = 0x01
	BundleFlagCompressed    = 0x10
	BundleFlagBroadcasted   = 0x20
	BundleFlagHelloMessage  = 0x40
)

type Bundle struct {
	Version   int
	Flags     int
	Initiator peer.ID
	Message   message.Message
}

func NewBundle(initiator peer.ID, msg message.Message) *Bundle {
	return &Bundle{
		Version:   LastVersion,
		Flags:     BundleFlagNetworkLibP2P,
		Initiator: initiator,
		Message:   msg,
	}
}

func (b *Bundle) SanityCheck() error {
	if err := b.Message.SanityCheck(); err != nil {
		return err
	}
	if err := b.Initiator.Validate(); err != nil {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid initiator peer id: %v", err)
	}
	return nil
}

func (b *Bundle) Fingerprint() string {
	return fmt.Sprintf("{%s: %s%s}", util.FingerprintPeerID(b.Initiator), b.Message.Type(), b.Message.Fingerprint())
}

func (b *Bundle) CompressIt() {
	b.Flags = util.SetFlag(b.Flags, BundleFlagCompressed)
}

type _Message struct {
	Version     int          `cbor:"1,keyasint"`
	Flags       int          `cbor:"2,keyasint"`
	Initiator   peer.ID      `cbor:"3,keyasint"`
	MessageType message.Type `cbor:"4,keyasint"`
	Message     []byte       `cbor:"5,keyasint"`
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

	msg := &_Message{
		Version:     b.Version,
		Flags:       b.Flags,
		Initiator:   b.Initiator,
		MessageType: b.Message.Type(),
		Message:     data,
	}

	return cbor.Marshal(msg)
}

func (b *Bundle) Decode(r io.Reader) (int, error) {
	var msg _Message
	d := cbor.NewDecoder(r)
	err := d.Decode(&msg)
	bytesRead := d.NumBytesRead()
	if err != nil {
		return bytesRead, errors.Errorf(errors.ErrInvalidMessage, err.Error())
	}

	data := msg.Message
	pld := message.MakeMessage(msg.MessageType)
	if pld == nil {
		return bytesRead, errors.Errorf(errors.ErrInvalidMessage, "invalid payload")
	}

	if util.IsFlagSet(msg.Flags, BundleFlagCompressed) {
		c, err := util.DecompressBuffer(msg.Message)
		if err != nil {
			return bytesRead, errors.Errorf(errors.ErrInvalidMessage, err.Error())
		}
		data = c
	}

	b.Version = msg.Version
	b.Flags = msg.Flags
	b.Initiator = msg.Initiator
	b.Message = pld
	return bytesRead, cbor.Unmarshal(data, pld)
}
