package message

import (
	"fmt"

	"github.com/pactus-project/pactus/network"
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

type Type int32

const (
	TypeUnspecified    = Type(0)
	TypeHello          = Type(1)
	TypeHeartBeat      = Type(2)
	TypeTransactions   = Type(3)
	TypeQueryProposal  = Type(4)
	TypeProposal       = Type(5)
	TypeQueryVotes     = Type(6)
	TypeVote           = Type(7)
	TypeBlockAnnounce  = Type(8)
	TypeBlocksRequest  = Type(9)
	TypeBlocksResponse = Type(10)
)

func (t Type) TopicID() network.TopicID {
	switch t {
	case TypeHello,
		TypeHeartBeat,
		TypeTransactions,
		TypeBlockAnnounce:
		return network.TopicIDGeneral

	case TypeQueryProposal,
		TypeProposal,
		TypeQueryVotes,
		TypeVote:
		return network.TopicIDConsensus

	default:
		panic("Invalid topic ID")
	}
}

func (t Type) String() string {
	switch t {
	case TypeHello:
		return "hello"
	case TypeHeartBeat:
		return "heart-beat"
	case TypeTransactions:
		return "txs"
	case TypeQueryProposal:
		return "query-proposal"
	case TypeProposal:
		return "proposal"
	case TypeQueryVotes:
		return "query-votes"
	case TypeVote:
		return "vote"
	case TypeBlockAnnounce:
		return "block-announce"
	case TypeBlocksRequest:
		return "blocks-req"
	case TypeBlocksResponse:
		return "blocks-res"
	}
	return fmt.Sprintf("%d", t)
}

func MakeMessage(t Type) Message {
	switch t {
	case TypeHello:
		return &HelloMessage{}
	case TypeHeartBeat:
		return &HeartBeatMessage{}
	case TypeTransactions:
		return &TransactionsMessage{}
	case TypeQueryProposal:
		return &QueryProposalMessage{}
	case TypeProposal:
		return &ProposalMessage{}
	case TypeQueryVotes:
		return &QueryVotesMessage{}
	case TypeVote:
		return &VoteMessage{}
	case TypeBlockAnnounce:
		return &BlockAnnounceMessage{}
	case TypeBlocksRequest:
		return &BlocksRequestMessage{}
	case TypeBlocksResponse:
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
