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
	PayloadTypeHello             = Type(1)
	PayloadTypeHeartBeat         = Type(2)
	PayloadTypeQueryTransactions = Type(3)
	PayloadTypeTransactions      = Type(4)
	PayloadTypeQueryProposal     = Type(5)
	PayloadTypeProposal          = Type(6)
	PayloadTypeQueryVotes        = Type(7)
	PayloadTypeVote              = Type(8)
	PayloadTypeBlockAnnounce     = Type(9)
	PayloadTypeBlocksRequest     = Type(10)
	PayloadTypeBlocksResponse    = Type(11)
)

func (t Type) TopicID() network.TopicID {
	switch t {
	case PayloadTypeHello,
		PayloadTypeHeartBeat,
		PayloadTypeQueryTransactions,
		PayloadTypeTransactions,
		PayloadTypeBlockAnnounce:
		return network.TopicIDGeneral

	case PayloadTypeQueryProposal,
		PayloadTypeProposal,
		PayloadTypeQueryVotes,
		PayloadTypeVote:
		return network.TopicIDConsensus

	default:
		panic("Invalid topic ID")
	}
}

func (t Type) String() string {
	switch t {
	case PayloadTypeHello:
		return "hello"
	case PayloadTypeHeartBeat:
		return "heart-beat"
	case PayloadTypeQueryTransactions:
		return "query-txs"
	case PayloadTypeTransactions:
		return "txs"
	case PayloadTypeQueryProposal:
		return "query-proposal"
	case PayloadTypeProposal:
		return "proposal"
	case PayloadTypeQueryVotes:
		return "query-votes"
	case PayloadTypeVote:
		return "vote"
	case PayloadTypeBlockAnnounce:
		return "block-announce"
	case PayloadTypeBlocksRequest:
		return "blocks-req"
	case PayloadTypeBlocksResponse:
		return "blocks-res"
	}
	return fmt.Sprintf("%d", t)
}

func MakePayload(t Type) Payload {
	switch t {
	case PayloadTypeHello:
		return &HelloPayload{}
	case PayloadTypeHeartBeat:
		return &HeartBeatPayload{}
	case PayloadTypeQueryTransactions:
		return &QueryTransactionsPayload{}
	case PayloadTypeTransactions:
		return &TransactionsPayload{}
	case PayloadTypeQueryProposal:
		return &QueryProposalPayload{}
	case PayloadTypeProposal:
		return &ProposalPayload{}
	case PayloadTypeQueryVotes:
		return &QueryVotesPayload{}
	case PayloadTypeVote:
		return &VotePayload{}
	case PayloadTypeBlockAnnounce:
		return &BlockAnnouncePayload{}
	case PayloadTypeBlocksRequest:
		return &BlocksRequestPayload{}
	case PayloadTypeBlocksResponse:
		return &BlocksResponsePayload{}
	}

	//
	return nil
}

type Payload interface {
	SanityCheck() error
	Type() Type
	Fingerprint() string
}
