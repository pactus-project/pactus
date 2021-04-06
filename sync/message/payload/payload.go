package payload

import (
	"fmt"
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
	PayloadTypeQueryVotes           = PayloadType(10)
	PayloadTypeVote                 = PayloadType(11)
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
	case PayloadTypeQueryVotes:
		return "query-votes"
	case PayloadTypeVote:
		return "vote"
	case PayloadTypeBlockAnnounce:
		return "block-announce"
	case PayloadTypeDownloadRequest:
		return "download-req"
	case PayloadTypeDownloadResponse:
		return "download-res"
	}
	return fmt.Sprintf("%d", t)
}

func MakePayload(t PayloadType) Payload {
	switch t {
	case PayloadTypeSalam:
		return &SalamPayload{}
	case PayloadTypeAleyk:
		return &AleykPayload{}
	case PayloadTypeLatestBlocksRequest:
		return &LatestBlocksRequestPayload{}
	case PayloadTypeLatestBlocksResponse:
		return &LatestBlocksResponsePayload{}
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
	case PayloadTypeDownloadRequest:
		return &DownloadRequestPayload{}
	case PayloadTypeDownloadResponse:
		return &DownloadResponsePayload{}
	}

	//
	return nil
}

type Payload interface {
	SanityCheck() error
	Type() PayloadType
	Fingerprint() string
}
