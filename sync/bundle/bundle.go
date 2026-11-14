package bundle

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/fxamacker/cbor/v2"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/encoding"
)

const (
	BundleFlagNetworkMainnet = 0x0001
	BundleFlagNetworkTestnet = 0x0002
	BundleFlagCarrierLibP2P  = 0x0010
	BundleFlagCompressed     = 0x0100 // Deprecated
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
	ConsensusHeight types.Height
}

func NewBundle(msg message.Message) *Bundle {
	return &Bundle{
		Flags:           0,
		Message:         msg,
		ConsensusHeight: msg.ConsensusHeight(),
	}
}

// BasicCheck performs basic validation checks on the bundle and its message.
func (b *Bundle) BasicCheck() error {
	return b.Message.BasicCheck()
}

// LogString returns a concise string representation intended for use in logs.
func (b *Bundle) LogString() string {
	return fmt.Sprintf("%s%s", b.Message.Type(), b.Message.LogString())
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

	msg := &_Bundle{
		Flags:           b.Flags,
		MessageType:     b.Message.Type(),
		MessageData:     data,
		ConsensusHeight: fixedUint32(b.ConsensusHeight),
	}

	return cbor.Marshal(msg)
}

func (b *Bundle) Decode(r io.Reader) (int, error) {
	return b.DecodeWithLimit(r, util.DefaultMaxDecompressedSize)
}

// DecodeWithLimit decodes a bundle. A compressed message is rejected if it
// expands to more than maxMessageSize bytes.
func (b *Bundle) DecodeWithLimit(r io.Reader, maxMessageSize int) (int, error) {
	// Read the whole message into memory first: the CBOR decoder below
	// stops after one item and, reading directly from r, may itself pull
	// more bytes from r than it logically consumes (its own internal
	// buffering), which would make any trailing-bytes check against r
	// unreliable. Decoding from a fresh, fully-read buffer lets bytesRead
	// be compared against the buffer's real length. Both callers (a gossip
	// message's in-memory bytes, and a stream message already buffered in
	// full, see sync.dispatchStreamMessage) already hand this function a
	// bounded, fully-available reader, so this never blocks.
	wire, err := io.ReadAll(io.LimitReader(r, encoding.MaxPayloadSize))
	if err != nil {
		return 0, err
	}

	var bdl _Bundle
	d := encoding.NewCBORDecoder(bytes.NewReader(wire))
	err = d.Decode(&bdl)
	bytesRead := d.NumBytesRead()
	if err != nil {
		return bytesRead, err
	}

	// The CBOR decoder above stops after one item; reject anything left
	// in wire instead of silently ignoring it. Two wire messages that
	// differ only by a trailing junk suffix would otherwise decode to the
	// identical Bundle while hashing to different gossip message IDs
	// (network.MessageIDFunc hashes the raw wire bytes), defeating
	// pubsub's own deduplication.
	if bytesRead != len(wire) {
		return bytesRead, errors.New("bundle: trailing bytes after decoded message")
	}

	data := bdl.MessageData
	msg, err := message.MakeMessage(bdl.MessageType)
	if err != nil {
		return bytesRead, err
	}

	if util.IsFlagSet(bdl.Flags, BundleFlagCompressed) {
		c, err := util.DecompressBuffer(bdl.MessageData, util.WithMaxDecompressedSize(maxMessageSize))
		if err != nil {
			return bytesRead, err
		}
		data = c
	}

	if err := encoding.UnmarshalCBOR(data, msg); err != nil {
		return bytesRead, err
	}

	b.Flags = bdl.Flags
	b.Message = msg
	b.ConsensusHeight = types.Height(bdl.ConsensusHeight)

	return bytesRead, nil
}
