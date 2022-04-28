package sync

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/sync/bundle"
	"github.com/zarbchain/zarb-go/sync/bundle/message"
	"github.com/zarbchain/zarb-go/util/errors"
)

type queryProposalHandler struct {
	*synchronizer
}

func newQueryProposalHandler(sync *synchronizer) messageHandler {
	return &queryProposalHandler{
		sync,
	}
}

func (handler *queryProposalHandler) ParsMessage(m message.Message, initiator peer.ID) error {
	msg := m.(*message.QueryProposalMessage)
	handler.logger.Trace("parsing QueryProposal message", "message", msg)

	height, round := handler.consensus.HeightRound()
	if msg.Height == height && msg.Round == round {
		if !handler.peerIsInTheCommittee(initiator) {
			return errors.Errorf(errors.ErrInvalidMessage, "peers is not in the commmittee")
		}

		p := handler.consensus.RoundProposal(msg.Round)
		if p != nil {
			response := message.NewProposalMessage(p)
			handler.broadcast(response)
		}
	}

	return nil
}

func (handler *queryProposalHandler) PrepareBundle(m message.Message) *bundle.Bundle {
	msg := m.(*message.QueryProposalMessage)
	proposal := handler.consensus.RoundProposal(msg.Round)
	if proposal == nil {
		proposal = handler.cache.GetProposal(msg.Height, msg.Round)
		if proposal != nil {
			// We have the proposal inside the cache
			handler.consensus.SetProposal(proposal)
		} else {
			if handler.weAreInTheCommittee() {
				msg := bundle.NewBundle(handler.SelfID(), m)
				return msg
			}
			handler.logger.Debug("not an active validator", "message", msg)
		}
	}

	return nil
}
