package message

import (
	"fmt"

	"github.com/fxamacker/cbor/v2"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/util"
)

const LastVersion = 1
const FlagCompressed = 0x1

type Message struct {
	Version   int
	Flags     int
	Initiator peer.ID
	Payload   payload.Payload
}

func NewMessage(initiator peer.ID, pld payload.Payload) *Message {
	return &Message{
		Version:   LastVersion,
		Flags:     0,
		Initiator: initiator,
		Payload:   pld,
	}
}

func (m *Message) SanityCheck() error {
	if err := m.Payload.SanityCheck(); err != nil {
		return err
	}
	if m.Flags|FlagCompressed != FlagCompressed {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid flags")
	}
	if err := m.Initiator.Validate(); err != nil {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid initiator peer id: %v", err)
	}
	return nil
}

func (m *Message) Fingerprint() string {
	return fmt.Sprintf("{%s: %s%s}", util.FingerprintPeerID(m.Initiator), m.Payload.Type(), m.Payload.Fingerprint())
}

func (m *Message) CompressIt() {
	m.Flags = util.SetFlag(m.Flags, FlagCompressed)
}

type _Message struct {
	Version     int                 `cbor:"1,keyasint"`
	Flags       int                 `cbor:"2,keyasint"`
	Initiator   peer.ID             `cbor:"3,keyasint"`
	PayloadType payload.PayloadType `cbor:"4,keyasint"`
	Payload     []byte              `cbor:"5,keyasint"`
	Signature   *crypto.Signature   `cbor:"6,keyasint,omitempty"`
}

func (m *Message) Encode() ([]byte, error) {
	data, err := cbor.Marshal(m.Payload)
	if err != nil {
		return nil, err
	}

	if util.IsFlagSet(m.Flags, FlagCompressed) {
		c, err := util.CompressBuffer(data)
		if err == nil {
			data = c
		}
	}

	msg := &_Message{
		Version:     m.Version,
		Flags:       m.Flags,
		Initiator:   m.Initiator,
		PayloadType: m.Payload.Type(),
		Payload:     data,
	}

	return cbor.Marshal(msg)
}

func (m *Message) Decode(bs []byte) error {
	var msg _Message
	err := cbor.Unmarshal(bs, &msg)
	if err != nil {
		return err
	}

	data := msg.Payload
	pld := payload.MakePayload(msg.PayloadType)
	if pld == nil {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid payload")
	}

	if util.IsFlagSet(msg.Flags, FlagCompressed) {
		c, err := util.DecompressBuffer(msg.Payload)
		if err != nil {
			return errors.Errorf(errors.ErrInvalidMessage, err.Error())
		}
		data = c
	}

	m.Version = msg.Version
	m.Flags = msg.Flags
	m.Initiator = msg.Initiator
	m.Payload = pld
	return cbor.Unmarshal(data, pld)
}
