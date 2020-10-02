package message

import (
	"fmt"

	"github.com/fxamacker/cbor/v2"
	"github.com/zarbchain/zarb-go/errors"
)

type PayloadType int

const (
	PayloadTypeTx      = PayloadType(1)
	PayloadTypeRequest = PayloadType(2)
)

func (t PayloadType) String() string {
	switch t {
	case PayloadTypeTx:
		return "tx"
	case PayloadTypeRequest:
		return "request"
	}
	return "invalid message type"
}

type Message struct {
	Type    PayloadType
	Payload Payload
}

func (m *Message) SanityCheck() error {
	if err := m.Payload.SanityCheck(); err != nil {
		return err
	}
	if m.Type != m.Payload.Type() {
		errors.Errorf(errors.ErrInvalidMessage, "invalid message type")
	}
	return nil
}

func (m *Message) Fingerprint() string {
	return fmt.Sprintf("{%s%s}", m.Type, m.Payload.Fingerprint())
}

func (m *Message) PayloadType() PayloadType {
	return m.Type
}

type _Message struct {
	PayloadType PayloadType     `cbor:"1,keyasint,omitempty"`
	Payload     cbor.RawMessage `cbor:"10,keyasint,omitempty"`
}

func (m *Message) MarshalCBOR() ([]byte, error) {
	bs, err := cbor.Marshal(m.Payload)
	if err != nil {
		return nil, err
	}

	msg := &_Message{
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

	var payload Payload
	switch msg.PayloadType {
	case PayloadTypeTx:
		payload = &TxPayload{}
	case PayloadTypeRequest:
		payload = &RequestPayload{}
	default:
		return errors.Errorf(errors.ErrInvalidMessage, "Invalid payload")
	}

	m.Type = msg.PayloadType
	m.Payload = payload
	return cbor.Unmarshal(msg.Payload, payload)
}

type Payload interface {
	SanityCheck() error
	Type() PayloadType
	Fingerprint() string
}
