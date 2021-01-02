package message

import (
	"fmt"

	"github.com/fxamacker/cbor/v2"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/message/payload"
	"github.com/zarbchain/zarb-go/util"
)

const LastVersion = 1
const FlagCompressed = 0x1
const FlagHasSignature = 0x2

type Message struct {
	Version   int
	Flags     int
	Type      payload.PayloadType
	Payload   payload.Payload
	Signature *crypto.Signature
}

func (m *Message) SanityCheck() error {
	if err := m.Payload.SanityCheck(); err != nil {
		return err
	}
	if m.Type != m.Payload.Type() {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid message type")
	}
	if m.Flags|0x3 != 0x3 {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid flags")
	}
	if util.IsFlagSet(m.Flags, FlagHasSignature) && m.Signature == nil {
		return errors.Errorf(errors.ErrInvalidMessage, "should have signature")
	}
	return nil
}

func (m *Message) Fingerprint() string {
	return fmt.Sprintf("{%s %s}", m.Type, m.Payload.Fingerprint())
}

func (m *Message) PayloadType() payload.PayloadType {
	return m.Type
}

func (m *Message) SetSignature(sig *crypto.Signature) {
	m.Signature = sig
}

func (m *Message) CompressIt() {
	m.Flags = util.SetFlag(m.Flags, FlagCompressed)
}

type _Message struct {
	Version     int                 `cbor:"1,keyasint"`
	Flags       int                 `cbor:"2,keyasint"`
	PayloadType payload.PayloadType `cbor:"3,keyasint"`
	Payload     []byte              `cbor:"4,keyasint"`
	Signature   *crypto.Signature   `cbor:"21,keyasint,omitempty"`
}

func (m *Message) SignBytes() []byte {
	ms := new(Message)
	*ms = *m
	ms.Flags = 0
	ms.Signature = nil
	sb, _ := ms.Encode()
	return sb
}

func (m *Message) Encode() ([]byte, error) {
	pld, err := cbor.Marshal(m.Payload)
	if err != nil {
		return nil, err
	}

	if util.IsFlagSet(m.Flags, FlagCompressed) {
		c, err := util.CompressSlice(pld)
		if err == nil {
			pld = c
		}
	}

	msg := &_Message{
		Version:     m.Version,
		Flags:       m.Flags,
		PayloadType: m.Type,
		Payload:     pld,
		Signature:   m.Signature,
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
	pld := makePayload(msg.PayloadType)
	if pld == nil {
		return errors.Errorf(errors.ErrInvalidMessage, "Invalid payload")
	}

	if util.IsFlagSet(msg.Flags, FlagCompressed) {
		c, err := util.DecompressSlice(msg.Payload)
		if err != nil {
			return errors.Errorf(errors.ErrInvalidMessage, err.Error())
		}
		data = c
	}

	m.Version = msg.Version
	m.Flags = msg.Flags
	m.Type = msg.PayloadType
	m.Signature = msg.Signature
	m.Payload = pld
	return cbor.Unmarshal(data, pld)
}
