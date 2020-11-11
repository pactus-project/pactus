package message

import (
	"fmt"

	"github.com/fxamacker/cbor/v2"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
)

type PayloadType int

const (
	PayloadTypeSalam       = PayloadType(1)
	PayloadTypeBlocksReq   = PayloadType(2)
	PayloadTypeBlocks      = PayloadType(3)
	PayloadTypeTxsReq      = PayloadType(4)
	PayloadTypeTxs         = PayloadType(5)
	PayloadTypeProposalReq = PayloadType(6)
	PayloadTypeProposal    = PayloadType(7)
	PayloadTypeHeartBeat   = PayloadType(8)
	PayloadTypeVote        = PayloadType(9)
	PayloadTypeVoteSet     = PayloadType(10)
)

func (t PayloadType) String() string {
	switch t {
	case PayloadTypeSalam:
		return "salam"
	case PayloadTypeBlocksReq:
		return "blocks-req"
	case PayloadTypeBlocks:
		return "blocks"
	case PayloadTypeTxsReq:
		return "txs-req"
	case PayloadTypeTxs:
		return "txs"
	case PayloadTypeProposalReq:
		return "proposal-req"
	case PayloadTypeProposal:
		return "proposal"
	case PayloadTypeHeartBeat:
		return "heart-beat"
	case PayloadTypeVote:
		return "vote"
	case PayloadTypeVoteSet:
		return "vote-set"
	}
	return fmt.Sprintf("%d", t)
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
	case PayloadTypeSalam:
		payload = &SalamPayload{}
	case PayloadTypeBlocksReq:
		payload = &BlocksReqPayload{}
	case PayloadTypeBlocks:
		payload = &BlocksPayload{}
	case PayloadTypeTxsReq:
		payload = &TxsReqPayload{}
	case PayloadTypeTxs:
		payload = &TxsPayload{}
	case PayloadTypeProposalReq:
		payload = &ProposalReqPayload{}
	case PayloadTypeProposal:
		payload = &ProposalPayload{}
	case PayloadTypeHeartBeat:
		payload = &HeartBeatPayload{}
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
