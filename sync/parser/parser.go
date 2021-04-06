package parser

import (
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/sync/message"
	"github.com/zarbchain/zarb-go/sync/message/payload"
	"github.com/zarbchain/zarb-go/sync/parser/handler"
)

type Parser struct {
	handlers map[payload.PayloadType]handler.Handler
}

func NewParser(ctx *handler.HandlerContext) *Parser {
	handlers := make(map[payload.PayloadType]handler.Handler)
	handlers[payload.PayloadTypeSalam] = handler.NewSalamHandler(ctx)
	handlers[payload.PayloadTypeAleyk] = handler.NewAleykHandler(ctx)
	handlers[payload.PayloadTypeHeartBeat] = handler.NewHeartBeatHandler(ctx)
	handlers[payload.PayloadTypeVote] = handler.NewVoteHandler(ctx)
	handlers[payload.PayloadTypeTransactions] = handler.NewTransactionsHandler(ctx)
	handlers[payload.PayloadTypeHeartBeat] = handler.NewHeartBeatHandler(ctx)
	handlers[payload.PayloadTypeQueryVotes] = handler.NewQueryVotesHandler(ctx)
	handlers[payload.PayloadTypeQueryProposal] = handler.NewQueryProposalHandler(ctx)
	handlers[payload.PayloadTypeQueryTransactions] = handler.NewQueryTransactionsHandler(ctx)
	handlers[payload.PayloadTypeBlockAnnounce] = handler.NewBlockAnnounceHandler(ctx)
	handlers[payload.PayloadTypeDownloadRequest] = handler.NewDownloadRequestHandler(ctx)
	handlers[payload.PayloadTypeDownloadResponse] = handler.NewDownloadResponseHandler(ctx)
	handlers[payload.PayloadTypeLatestBlocksRequest] = handler.NewLatestBlocksRequestHandler(ctx)
	handlers[payload.PayloadTypeLatestBlocksResponse] = handler.NewLatestBlocksResponseHandler(ctx)

	return &Parser{handlers: handlers}
}

func (p *Parser) ParsMessage(msg *message.Message) error {
	h := p.handlers[msg.Payload.Type()]
	if h == nil {
		return errors.Errorf(errors.ErrInvalidMessage, "Invalid payload type: %v", msg.Payload.Type())
	}

	return h.ParsPayload(msg.Payload, msg.Initiator)
}
