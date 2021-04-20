package sync

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/sync/message"
	"github.com/zarbchain/zarb-go/sync/message/payload"
)

type queryProposalHandler struct {
	*synchronizer
}

func newQueryProposalHandler(sync *synchronizer) payloadHandler {
	return &queryProposalHandler{
		sync,
	}
}

func (handler *queryProposalHandler) ParsPayload(p payload.Payload, initiator peer.ID) error {
	pld := p.(*payload.QueryProposalPayload)
	handler.logger.Trace("Parsing query proposal payload", "pld", pld)

	height, round := handler.consensus.HeightRound()
	if pld.Height == height && pld.Round == round {
		if !handler.peerIsInTheCommittee(initiator) {
			return errors.Errorf(errors.ErrInvalidMessage, "peers is not in the commmittee")
		}

		p := handler.consensus.RoundProposal(pld.Round)
		if p != nil {
			response := payload.NewProposalPayload(p)
			handler.broadcast(response)
		}
	}

	return nil
}

func (handler *queryProposalHandler) PrepareMessage(p payload.Payload) *message.Message {
	pld := p.(*payload.QueryProposalPayload)
	proposal := handler.consensus.RoundProposal(pld.Round)
	if proposal == nil {
		proposal = handler.cache.GetProposal(pld.Height, pld.Round)
		if proposal != nil {
			// We have the proposal inside the cache
			handler.consensus.SetProposal(proposal)
		} else {
			if handler.weAreInTheCommittee() {
				msg := message.NewMessage(handler.SelfID(), p)
				return msg
			} else {
				handler.logger.Debug("queryProposal ignored. Not an active validator")
			}
		}
	}

	return nil
}
