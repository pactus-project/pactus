package payload

import (
	"fmt"

	"github.com/zarbchain/zarb-go/network"
)

type ResponseCode int

const (
	ResponseCodeNone         = ResponseCode(-1)
	ResponseCodeOK           = ResponseCode(0)
	ResponseCodeRejected     = ResponseCode(1)
	ResponseCodeBusy         = ResponseCode(2)
	ResponseCodeMoreBlocks   = ResponseCode(3)
	ResponseCodeNoMoreBlocks = ResponseCode(4)
	ResponseCodeSynced       = ResponseCode(5)
)

func (c ResponseCode) String() string {
	switch c {
	case ResponseCodeOK:
		return "ok"
	case ResponseCodeRejected:
		return "rejected"
	case ResponseCodeBusy:
		return "busy"
	case ResponseCodeMoreBlocks:
		return "more-blocks"
	case ResponseCodeNoMoreBlocks:
		return "no-more-blocks"
	case ResponseCodeSynced:
		return "synced"
	}
	return fmt.Sprintf("%d", c)
}

type Type int

const (
	PayloadTypeSalam             = Type(1) // Hello message
	PayloadTypeAleyk             = Type(2) // Hello Ack message
	PayloadTypeBlocksRequest     = Type(3)
	PayloadTypeBlocksResponse    = Type(4)
	PayloadTypeQueryTransactions = Type(5)
	PayloadTypeTransactions      = Type(6)
	PayloadTypeQueryProposal     = Type(7)
	PayloadTypeProposal          = Type(8)
	PayloadTypeHeartBeat         = Type(9)
	PayloadTypeQueryVotes        = Type(10)
	PayloadTypeVote              = Type(11)
	PayloadTypeBlockAnnounce     = Type(12)
)

func (t Type) TopicID() network.TopicID {
	switch t {
	case PayloadTypeSalam,
		PayloadTypeAleyk,
		PayloadTypeHeartBeat,
		PayloadTypeQueryTransactions,
		PayloadTypeTransactions,
		PayloadTypeBlockAnnounce:
		return network.GeneralTopic

	case PayloadTypeQueryProposal,
		PayloadTypeProposal,
		PayloadTypeQueryVotes,
		PayloadTypeVote:
		return network.ConsensusTopic

	default:
		panic("Invalid topic ID")
	}
}

func (t Type) String() string {
	switch t {
	case PayloadTypeSalam:
		return "salam"
	case PayloadTypeAleyk:
		return "aleyk"
	case PayloadTypeBlocksRequest:
		return "blocks-req"
	case PayloadTypeBlocksResponse:
		return "blocks-res"
	case PayloadTypeQueryTransactions:
		return "query-txs"
	case PayloadTypeTransactions:
		return "txs"
	case PayloadTypeQueryProposal:
		return "query-proposal"
	case PayloadTypeProposal:
		return "proposal"
	case PayloadTypeHeartBeat:
		return "heart-beat"
	case PayloadTypeQueryVotes:
		return "query-votes"
	case PayloadTypeVote:
		return "vote"
	case PayloadTypeBlockAnnounce:
		return "block-announce"
	}
	return fmt.Sprintf("%d", t)
}

func MakePayload(t Type) Payload {
	switch t {
	case PayloadTypeSalam:
		return &SalamPayload{}
	case PayloadTypeAleyk:
		return &AleykPayload{}
	case PayloadTypeBlocksRequest:
		return &BlocksRequestPayload{}
	case PayloadTypeBlocksResponse:
		return &BlocksResponsePayload{}
	case PayloadTypeQueryTransactions:
		return &QueryTransactionsPayload{}
	case PayloadTypeTransactions:
		return &TransactionsPayload{}
	case PayloadTypeQueryProposal:
		return &QueryProposalPayload{}
	case PayloadTypeProposal:
		return &ProposalPayload{}
	case PayloadTypeHeartBeat:
		return &HeartBeatPayload{}
	case PayloadTypeQueryVotes:
		return &QueryVotesPayload{}
	case PayloadTypeVote:
		return &VotePayload{}
	case PayloadTypeBlockAnnounce:
		return &BlockAnnouncePayload{}
	}

	//
	return nil
}

type Payload interface {
	SanityCheck() error
	Type() Type
	Fingerprint() string
}
