package message

import (
	"fmt"

	"github.com/fxamacker/cbor/v2"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/message/payload"
)

type Message struct {
	Initiator crypto.Address
	Target    crypto.Address
	Flags     int
	Type      payload.PayloadType
	Payload   payload.Payload
}

func (m *Message) SanityCheck() error {
	if err := m.Payload.SanityCheck(); err != nil {
		return err
	}
	if m.Type != m.Payload.Type() {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid message type")
	}
	if m.Flags != 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid flags")
	}
	return nil
}

func (m *Message) Fingerprint() string {
	return fmt.Sprintf("{%s %s}", m.Type, m.Payload.Fingerprint())
}

func (m *Message) PayloadType() payload.PayloadType {
	return m.Type
}

type _Message struct {
	Initiator   crypto.Address      `cbor:"1,keyasint,omitempty"`
	Target      crypto.Address      `cbor:"2,keyasint,omitempty"`
	Flags       int                 `cbor:"3,keyasint,omitempty"`
	PayloadType payload.PayloadType `cbor:"4,keyasint"`
	Payload     cbor.RawMessage     `cbor:"21,keyasint"`
}

func (m *Message) MarshalCBOR() ([]byte, error) {
	bs, err := cbor.Marshal(m.Payload)
	if err != nil {
		return nil, err
	}

	msg := &_Message{
		Initiator:   m.Initiator,
		Target:      m.Target,
		Flags:       m.Flags,
		PayloadType: m.Type,
		Payload:     bs,
	}

	return cbor.Marshal(msg)
}

func (m *Message) UnmarshalCBOR(bs []byte) error {
	var msg _Message
	err := cbor.Unmarshal(bs, &msg)
	if err != nil {
		return err
	}

	pld := makePayload(msg.PayloadType)
	if pld == nil {
		return errors.Errorf(errors.ErrInvalidMessage, "Invalid payload")
	}

	m.Type = msg.PayloadType
	m.Payload = pld
	return cbor.Unmarshal(msg.Payload, pld)
}
