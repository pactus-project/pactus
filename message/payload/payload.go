package payload

import (
	"fmt"
)

type PayloadType int

const (
	PayloadTypeSalam       = PayloadType(1) // Hello message
	PayloadTypeAleyk       = PayloadType(2) // Hello Ack message
	PayloadTypeBlocksReq   = PayloadType(3)
	PayloadTypeBlocks      = PayloadType(4)
	PayloadTypeTxsReq      = PayloadType(5)
	PayloadTypeTxs         = PayloadType(6)
	PayloadTypeProposalReq = PayloadType(7)
	PayloadTypeProposal    = PayloadType(8)
	PayloadTypeHeartBeat   = PayloadType(9)
	PayloadTypeVote        = PayloadType(10)
	PayloadTypeVoteSet     = PayloadType(11)
)

func (t PayloadType) String() string {
	switch t {
	case PayloadTypeSalam:
		return "salam"
	case PayloadTypeAleyk:
		return "aleyk"
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

type Payload interface {
	SanityCheck() error
	Type() PayloadType
	Fingerprint() string
}
