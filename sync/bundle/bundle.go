package bundle

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/fxamacker/cbor/v2"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/util"
)

const (
	BundleFlagNetworkMainnet = 0x0001
	BundleFlagNetworkTestnet = 0x0002
	BundleFlagCarrierLibP2P  = 0x0010
	BundleFlagCompressed     = 0x0100
	BundleFlagBroadcasted    = 0x0200
	BundleFlagHandshaking    = 0x0400
)

// Custom type to enforce uint32 encoding as 4 bytes, ignoring zeros.
type fixedUint32 uint32

func (u fixedUint32) MarshalCBOR() ([]byte, error) {
	buf := make([]byte, 0, 5)

	// The header for a 4-byte integer is 0x1A followed by the 4 bytes of the uint32.
	buf = append(buf, 0x1A)

	// Append the uint32 in big-endian format
	buf = binary.BigEndian.AppendUint32(buf, uint32(u))

	return buf, nil
}

type Bundle struct {
	Flags           int
	Message         message.Message
	ConsensusHeight uint32
}

func NewBundle(msg message.Message) *Bundle {
	return &Bundle{
		Flags:           0,
		Message:         msg,
		ConsensusHeight: msg.ConsensusHeight(),
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

type _Bundle struct {
	Flags           int          `cbor:"1,keyasint"`
	MessageType     message.Type `cbor:"2,keyasint"`
	MessageData     []byte       `cbor:"3,keyasint"`
	ConsensusHeight fixedUint32  `cbor:"4,keyasint,omitempty"`
}

func (b *Bundle) Encode() ([]byte, error) {
	data, err := cbor.Marshal(b.Message)
	if err != nil {
		return nil, err
	}

	if util.IsFlagSet(b.Flags, BundleFlagCompressed) {
		c, err := util.CompressBuffer(data)
		if err != nil {
			return nil, err
		}
		data = c
	}

	msg := &_Bundle{
		Flags:           b.Flags,
		MessageType:     b.Message.Type(),
		MessageData:     data,
		ConsensusHeight: fixedUint32(b.ConsensusHeight),
	}

	return cbor.Marshal(msg)
}

func (b *Bundle) Decode(r io.Reader) (int, error) {
	var bdl _Bundle
	decOpts := cbor.DecOptions{}
	decOpts.MaxArrayElements = 65_536 // default in 131072
	decOpts.MaxMapPairs = 65_536      // default in 131072
	decMode, _ := decOpts.DecMode()
	d := decMode.NewDecoder(r)
	err := d.Decode(&bdl)
	bytesRead := d.NumBytesRead()
	if err != nil {
		return bytesRead, err
	}

	data := bdl.MessageData
	msg, err := message.MakeMessage(bdl.MessageType)
	if err != nil {
		return bytesRead, err
	}

	if util.IsFlagSet(bdl.Flags, BundleFlagCompressed) {
		c, err := util.DecompressBuffer(bdl.MessageData)
		if err != nil {
			return bytesRead, err
		}
		data = c
	}

	b.Flags = bdl.Flags
	b.Message = msg
	b.ConsensusHeight = uint32(bdl.ConsensusHeight)

	return bytesRead, cbor.Unmarshal(data, msg)
}
