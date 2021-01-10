package payload

import (
	"fmt"
)

type ResponseCode int

const (
	ResponseCodeOK           = ResponseCode(0)
	ResponseCodeRejected     = ResponseCode(1)
	ResponseCodeBusy         = ResponseCode(2)
	ResponseCodeMoreBlocks   = ResponseCode(3)
	ResponseCodeNoMoreBlocks = ResponseCode(4)
)

type PayloadType int

const (
	PayloadTypeSalam                = PayloadType(1) // Hello message
	PayloadTypeAleyk                = PayloadType(2) // Hello Ack message
	PayloadTypeLatestBlocksRequest  = PayloadType(3)
	PayloadTypeLatestBlocksResponse = PayloadType(4)
	PayloadTypeQueryTransactions    = PayloadType(5)
	PayloadTypeTransactions         = PayloadType(6)
	PayloadTypeQueryProposal        = PayloadType(7)
	PayloadTypeProposal             = PayloadType(8)
	PayloadTypeHeartBeat            = PayloadType(9)
	PayloadTypeVote                 = PayloadType(10)
	PayloadTypeVoteSet              = PayloadType(11)
	PayloadTypeBlockAnnounce        = PayloadType(12)
	PayloadTypeDownloadRequest      = PayloadType(13)
	PayloadTypeDownloadResponse     = PayloadType(14)
)

func (t PayloadType) String() string {
	switch t {
	case PayloadTypeSalam:
		return "salam"
	case PayloadTypeAleyk:
		return "aleyk"
	case PayloadTypeLatestBlocksRequest:
		return "blocks-req"
	case PayloadTypeLatestBlocksResponse:
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
	case PayloadTypeVote:
		return "vote"
	case PayloadTypeVoteSet:
		return "vote-set"
	case PayloadTypeBlockAnnounce:
		return "block-announce"
	case PayloadTypeDownloadRequest:
		return "download-req"
	case PayloadTypeDownloadResponse:
		return "download-res"
	}
	return fmt.Sprintf("%d", t)
}

type Payload interface {
	SanityCheck() error
	Type() PayloadType
	Fingerprint() string
}
