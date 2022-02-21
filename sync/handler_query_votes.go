package sync

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/sync/message"
	"github.com/zarbchain/zarb-go/sync/message/payload"
)

type queryVotesHandler struct {
	*synchronizer
}

func newQueryVotesHandler(sync *synchronizer) payloadHandler {
	return &queryVotesHandler{
		sync,
	}
}

func (handler *queryVotesHandler) ParsPayload(p payload.Payload, initiator peer.ID) error {
	pld := p.(*payload.QueryVotesPayload)
	handler.logger.Trace("parsing query votes payload", "pld", pld)

	height, _ := handler.consensus.HeightRound()
	if pld.Height == height {
		if !handler.peerIsInTheCommittee(initiator) {
			return errors.Errorf(errors.ErrInvalidMessage, "peers is not in the commmittee")
		}
		v := handler.consensus.PickRandomVote()
		if v != nil {
			response := payload.NewVotePayload(v)
			handler.broadcast(response)
		}
	}

	return nil
}

func (handler *queryVotesHandler) PrepareMessage(p payload.Payload) *message.Message {
	if !handler.weAreInTheCommittee() {
		handler.logger.Debug("sending QueryVotes ignored. We are not in the committee")
		return nil
	}
	msg := message.NewMessage(handler.SelfID(), p)

	return msg
}
