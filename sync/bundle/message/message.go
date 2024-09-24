package message

import (
	"fmt"

	"github.com/pactus-project/pactus/network"
)

type ResponseCode int

const (
	ResponseCodeOK           = ResponseCode(0)
	ResponseCodeRejected     = ResponseCode(1)
	ResponseCodeMoreBlocks   = ResponseCode(2)
	ResponseCodeNoMoreBlocks = ResponseCode(3)
	ResponseCodeSynced       = ResponseCode(4)
)

func (c ResponseCode) String() string {
	switch c {
	case ResponseCodeOK:
		return "ok"

	case ResponseCodeRejected:
		return "rejected"

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
	TypeHello          = Type(1)
	TypeHelloAck       = Type(2)
	TypeTransaction    = Type(3)
	TypeQueryProposal  = Type(4)
	TypeProposal       = Type(5)
	TypeQueryVote      = Type(6)
	TypeVote           = Type(7)
	TypeBlockAnnounce  = Type(8)
	TypeBlocksRequest  = Type(9)
	TypeBlocksResponse = Type(10)
)

func (t Type) String() string {
	switch t {
	case TypeHello:
		return "hello"

	case TypeHelloAck:
		return "hello-ack"

	case TypeTransaction:
		return "transaction"

	case TypeQueryProposal:
		return "query-proposal"

	case TypeProposal:
		return "proposal"

	case TypeQueryVote:
		return "query-vote"

	case TypeVote:
		return "vote"

	case TypeBlockAnnounce:
		return "block-announce"

	case TypeBlocksRequest:
		return "blocks-request"

	case TypeBlocksResponse:
		return "blocks-response"

	default:
		return fmt.Sprintf("%d", t)
	}
}

func MakeMessage(t Type) (Message, error) {
	var msg Message
	switch t {
	case TypeHello:
		msg = &HelloMessage{}

	case TypeHelloAck:
		msg = &HelloAckMessage{}

	case TypeTransaction:
		msg = &TransactionsMessage{}

	case TypeQueryProposal:
		msg = &QueryProposalMessage{}

	case TypeProposal:
		msg = &ProposalMessage{}

	case TypeQueryVote:
		msg = &QueryVoteMessage{}

	case TypeVote:
		msg = &VoteMessage{}

	case TypeBlockAnnounce:
		msg = &BlockAnnounceMessage{}

	case TypeBlocksRequest:
		msg = &BlocksRequestMessage{}

	case TypeBlocksResponse:
		msg = &BlocksResponseMessage{}

	default:
		return nil, InvalidMessageTypeError{Type: int(t)}
	}

	//
	return msg, nil
}

type Message interface {
	BasicCheck() error
	Type() Type
	TopicID() network.TopicID
	ShouldBroadcast() bool
	String() string
}
