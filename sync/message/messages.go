package message

import (
	"fmt"

	"github.com/fxamacker/cbor/v2"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/util"
)

const LastVersion = 1
const FlagCompressed = 0x1

type Message struct {
	Version int
	Flags   int
	Type    payload.PayloadType
	Payload payload.Payload
}

func (m *Message) SanityCheck() error {
	if err := m.Payload.SanityCheck(); err != nil {
		return err
	}
	if m.Type != m.Payload.Type() {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid message type")
	}
	if m.Flags|FlagCompressed != FlagCompressed {
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

func (m *Message) CompressIt() {
	m.Flags = util.SetFlag(m.Flags, FlagCompressed)
}

type _Message struct {
	Version     int                 `cbor:"1,keyasint"`
	Flags       int                 `cbor:"2,keyasint"`
	PayloadType payload.PayloadType `cbor:"3,keyasint"`
	Payload     []byte              `cbor:"4,keyasint"`
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
	m.Payload = pld
	return cbor.Unmarshal(data, pld)
}
