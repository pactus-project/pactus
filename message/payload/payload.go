package payload

import (
	"fmt"
)

type PayloadType int

const (
	PayloadTypeSalam               = PayloadType(1) // Hello message
	PayloadTypeAleyk               = PayloadType(2) // Hello Ack message
	PayloadTypeLatestBlocksRequest = PayloadType(3)
	PayloadTypeLatestBlocks        = PayloadType(4)
	PayloadTypeTransactionsRequest = PayloadType(5)
	PayloadTypeTransactions        = PayloadType(6)
	PayloadTypeProposalRequest     = PayloadType(7)
	PayloadTypeProposal            = PayloadType(8)
	PayloadTypeHeartBeat           = PayloadType(9)
	PayloadTypeVote                = PayloadType(10)
	PayloadTypeVoteSet             = PayloadType(11)
)

func (t PayloadType) String() string {
	switch t {
	case PayloadTypeSalam:
		return "salam"
	case PayloadTypeAleyk:
		return "aleyk"
	case PayloadTypeLatestBlocksRequest:
		return "blocks-req"
	case PayloadTypeLatestBlocks:
		return "blocks"
	case PayloadTypeTransactionsRequest:
		return "txs-req"
	case PayloadTypeTransactions:
		return "txs"
	case PayloadTypeProposalRequest:
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
