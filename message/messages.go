package message

import (
	"fmt"

	"github.com/fxamacker/cbor/v2"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
)

type PayloadType int

const (
	PayloadTypeStatusReq = PayloadType(1)
	PayloadTypeStatusRes = PayloadType(2)
	PayloadTypeBlocksReq = PayloadType(3)
	PayloadTypeBlocksRes = PayloadType(4)
	PayloadTypeTxReq     = PayloadType(5)
	PayloadTypeTxRes     = PayloadType(6)
	PayloadTypeProposal  = PayloadType(7)
	PayloadTypeBlock     = PayloadType(8)
	PayloadTypeHRS       = PayloadType(9)
	PayloadTypeVote      = PayloadType(10)
	PayloadTypeVoteSet   = PayloadType(11)
)

func (t PayloadType) String() string {
	switch t {
	case PayloadTypeStatusReq:
		return "status-req"
	case PayloadTypeStatusRes:
		return "status-res"
	case PayloadTypeBlocksReq:
		return "blocks-req"
	case PayloadTypeBlocksRes:
		return "blocks-res"
	case PayloadTypeTxReq:
		return "tx-req"
	case PayloadTypeTxRes:
		return "tx-res"
	case PayloadTypeProposal:
		return "proposal"
	case PayloadTypeBlock:
		return "block"
	case PayloadTypeHRS:
		return "hrs"
	case PayloadTypeVote:
		return "vote"
	case PayloadTypeVoteSet:
		return "vote-set"
	}
	return "invalid message type"
}

type Message struct {
	Initiator crypto.Address
	Target    crypto.Address
	Type      PayloadType
	Payload   Payload
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
	return fmt.Sprintf("{%s %s}", m.Type, m.Payload.Fingerprint())
}

func (m *Message) PayloadType() PayloadType {
	return m.Type
}

type _Message struct {
	Initiator   crypto.Address  `cbor:"1,keyasint,omitempty"`
	Target      crypto.Address  `cbor:"2,keyasint,omitempty"`
	PayloadType PayloadType     `cbor:"3,keyasint"`
	Payload     cbor.RawMessage `cbor:"10,keyasint"`
}

func (m *Message) MarshalCBOR() ([]byte, error) {
	bs, err := cbor.Marshal(m.Payload)
	if err != nil {
		return nil, err
	}

	msg := &_Message{
		Initiator:   m.Initiator,
		Target:      m.Target,
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
	case PayloadTypeStatusReq:
		payload = &StatusReqPayload{}
	case PayloadTypeStatusRes:
		payload = &StatusResPayload{}
	case PayloadTypeBlocksReq:
		payload = &BlocksReqPayload{}
	case PayloadTypeBlocksRes:
		payload = &BlocksResPayload{}
	case PayloadTypeTxReq:
		payload = &TxReqPayload{}
	case PayloadTypeTxRes:
		payload = &TxResPayload{}
	case PayloadTypeProposal:
		payload = &ProposalPayload{}
	case PayloadTypeBlock:
		payload = &BlockPayload{}
	case PayloadTypeHRS:
		payload = &HRSPayload{}
	case PayloadTypeVote:
		payload = &VotePayload{}
	case PayloadTypeVoteSet:
		payload = &VoteSetPayload{}

	default:
		return errors.Errorf(errors.ErrInvalidMessage, "Invalid payload")
	}

	m.Type = msg.PayloadType
	m.Payload = payload
	cbor.Unmarshal(msg.Payload, payload)
	return cbor.Unmarshal(msg.Payload, payload)
}

type Payload interface {
	SanityCheck() error
	Type() PayloadType
	Fingerprint() string
}
