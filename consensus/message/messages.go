package message

import (
	"fmt"

	"github.com/fxamacker/cbor/v2"
	"gitlab.com/zarb-chain/zarb-go/errors"
)

type PayloadType int

const (
	PayloadTypeProposal = PayloadType(1)
	PayloadTypeBlock    = PayloadType(2)
	PayloadTypeStep     = PayloadType(3)
	PayloadTypeVote     = PayloadType(4)
	PayloadTypeVoteSet  = PayloadType(5)
)

func (t PayloadType) String() string {
	switch t {
	case PayloadTypeProposal:
		return "proposal"
	case PayloadTypeBlock:
		return "block"
	case PayloadTypeStep:
		return "round-step"
	case PayloadTypeVote:
		return "vote"
	case PayloadTypeVoteSet:
		return "vote-set"
	}
	return "invalid message type"
}

type Message struct {
	Height  int
	Type    PayloadType
	Payload Payload
}

func (m *Message) SanityCheck() error {
	if m.Height < 0 {
		return errors.Errorf(errors.ErrInvalidMessage, "invalid Height")
	}
	if err := m.Payload.SanityCheck(); err != nil {
		return err
	}
	if m.Type != m.Payload.Type() {
		errors.Errorf(errors.ErrInvalidMessage, "invalid message type")
	}
	return nil
}

func (m *Message) Fingerprint() string {
	return fmt.Sprintf("{%d %s %s}", m.Height, m.Type, m.Payload.Fingerprint())
}

func (m *Message) PayloadType() PayloadType {
	return m.Type
}

type _Message struct {
	PayloadType PayloadType     `cbor:"1,keyasint,omitempty"`
	Height      int             `cbor:"s,keyasint,omitempty"`
	Payload     cbor.RawMessage `cbor:"10,keyasint,omitempty"`
}

func (m *Message) MarshalCBOR() ([]byte, error) {
	bs, err := cbor.Marshal(m.Payload)
	if err != nil {
		return nil, err
	}

	msg := &_Message{
		PayloadType: m.Type,
		Height:      m.Height,
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
	case PayloadTypeProposal:
		payload = &ProposalPayload{}
	case PayloadTypeBlock:
		payload = &BlockPayload{}
	case PayloadTypeStep:
		payload = &StepPayload{}
	case PayloadTypeVote:
		payload = &VotePayload{}
	case PayloadTypeVoteSet:
		payload = &VoteSetPayload{}

	default:
		return errors.Errorf(errors.ErrInvalidMessage, "Invalid payload")
	}

	m.Type = msg.PayloadType
	m.Height = msg.Height
	m.Payload = payload
	return cbor.Unmarshal(msg.Payload, payload)
}

type Payload interface {
	SanityCheck() error
	Type() PayloadType
	Fingerprint() string
}
