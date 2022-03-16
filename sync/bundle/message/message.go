package message

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
	MessageTypeHello          = Type(1)
	MessageTypeHeartBeat      = Type(2)
	MessageTypeTransactions   = Type(3)
	MessageTypeQueryProposal  = Type(4)
	MessageTypeProposal       = Type(5)
	MessageTypeQueryVotes     = Type(6)
	MessageTypeVote           = Type(7)
	MessageTypeBlockAnnounce  = Type(8)
	MessageTypeBlocksRequest  = Type(9)
	MessageTypeBlocksResponse = Type(10)
)

func (t Type) TopicID() network.TopicID {
	switch t {
	case MessageTypeHello,
		MessageTypeHeartBeat,
		MessageTypeTransactions,
		MessageTypeBlockAnnounce:
		return network.TopicIDGeneral

	case MessageTypeQueryProposal,
		MessageTypeProposal,
		MessageTypeQueryVotes,
		MessageTypeVote:
		return network.TopicIDConsensus

	default:
		panic("Invalid topic ID")
	}
}

func (t Type) String() string {
	switch t {
	case MessageTypeHello:
		return "hello"
	case MessageTypeHeartBeat:
		return "heart-beat"
	case MessageTypeTransactions:
		return "txs"
	case MessageTypeQueryProposal:
		return "query-proposal"
	case MessageTypeProposal:
		return "proposal"
	case MessageTypeQueryVotes:
		return "query-votes"
	case MessageTypeVote:
		return "vote"
	case MessageTypeBlockAnnounce:
		return "block-announce"
	case MessageTypeBlocksRequest:
		return "blocks-req"
	case MessageTypeBlocksResponse:
		return "blocks-res"
	}
	return fmt.Sprintf("%d", t)
}

func MakeMessage(t Type) Message {
	switch t {
	case MessageTypeHello:
		return &HelloMessage{}
	case MessageTypeHeartBeat:
		return &HeartBeatMessage{}
	case MessageTypeTransactions:
		return &TransactionsMessage{}
	case MessageTypeQueryProposal:
		return &QueryProposalMessage{}
	case MessageTypeProposal:
		return &ProposalMessage{}
	case MessageTypeQueryVotes:
		return &QueryVotesMessage{}
	case MessageTypeVote:
		return &VoteMessage{}
	case MessageTypeBlockAnnounce:
		return &BlockAnnounceMessage{}
	case MessageTypeBlocksRequest:
		return &BlocksRequestMessage{}
	case MessageTypeBlocksResponse:
		return &BlocksResponseMessage{}
	}

	//
	return nil
}

type Message interface {
	SanityCheck() error
	Type() Type
	Fingerprint() string
}
