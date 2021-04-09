package sync

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/sync/message"
	"github.com/zarbchain/zarb-go/sync/message/payload"
)

type proposalHandler struct {
	*synchronizer
}

func newProposalHandler(sync *synchronizer) payloadHandler {
	return &proposalHandler{
		sync,
	}
}

func (handler *proposalHandler) ParsPayload(p payload.Payload, initiator peer.ID) error {
	pld := p.(*payload.ProposalPayload)
	handler.logger.Trace("Parsing proposal payload", "pld", pld)

	handler.cache.AddProposal(pld.Proposal)
	handler.consensus.SetProposal(pld.Proposal)

	return nil
}

func (handler *proposalHandler) PrepareMessage(p payload.Payload) *message.Message {
	return message.NewMessage(handler.SelfID(), p)
}
