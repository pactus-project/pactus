package sync

import (
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/sync/bundle"
	"github.com/pactus-project/pactus/sync/bundle/message"
	"github.com/pactus-project/pactus/util/errors"
)

type queryVotesHandler struct {
	*synchronizer
}

func newQueryVotesHandler(sync *synchronizer) messageHandler {
	return &queryVotesHandler{
		sync,
	}
}

func (handler *queryVotesHandler) ParsMessage(m message.Message, initiator peer.ID) error {
	msg := m.(*message.QueryVotesMessage)
	handler.logger.Trace("parsing QueryVotes message", "message", msg)

	height, _ := handler.consensus.HeightRound()
	if msg.Height == height {
		if !handler.peerIsInTheCommittee(initiator) {
			return errors.Errorf(errors.ErrInvalidMessage, "peers is not in the commmittee")
		}
		v := handler.consensus.PickRandomVote()
		if v != nil {
			response := message.NewVoteMessage(v)
			handler.broadcast(response)
		}
	}

	return nil
}

func (handler *queryVotesHandler) PrepareBundle(m message.Message) *bundle.Bundle {
	if !handler.weAreInTheCommittee() {
		handler.logger.Debug("sending QueryVotes ignored. We are not in the committee")
		return nil
	}
	bdl := bundle.NewBundle(handler.SelfID(), m)

	return bdl
}
