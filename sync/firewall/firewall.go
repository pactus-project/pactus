package firewall

import (
	"encoding/hex"

	"github.com/zarbchain/zarb-go/state"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/message"
	"github.com/zarbchain/zarb-go/message/payload"
	"github.com/zarbchain/zarb-go/sync/peerset"
	"github.com/zarbchain/zarb-go/util"
)

// Firewall hold statistic data about peers' behaviors
type Firewall struct {
	peerSet *peerset.PeerSet
	state   state.StateReader
}

func NewFirewall(peerSet *peerset.PeerSet, state state.StateReader) *Firewall {
	return &Firewall{
		peerSet: peerSet,
		state:   state,
	}
}

func (f *Firewall) ParsMessage(data []byte, from peer.ID) *message.Message {
	peer := f.peerSet.MustGetPeer(from)
	peer.IncreaseReceivedMsg()

	msg := new(message.Message)
	err := msg.Decode(data)
	if err != nil {
		peer.IncreaseInvalidMsg()
		logger.Debug("Error decoding message", "from", util.FingerprintPeerID(from), "data", hex.EncodeToString(data), "err", err)
		return nil
	}

	if err = msg.SanityCheck(); err != nil {
		peer.IncreaseInvalidMsg()
		logger.Debug("Peer sent us invalid msg", "peer", util.FingerprintPeerID(from), "msg", msg, "err", err)
		return nil
	}

	if f.badPeer(peer) {
		return nil
	}

	switch msg.PayloadType() {
	case payload.PayloadTypeSalam:
		pld := msg.Payload.(*payload.SalamPayload)
		f.peerSet.UpdateMaxClaimedHeight(pld.Height)

	case payload.PayloadTypeAleyk:
		pld := msg.Payload.(*payload.AleykPayload)
		peer.UpdateHeight(pld.Height)
		f.peerSet.UpdateMaxClaimedHeight(pld.Height)

	case payload.PayloadTypeHeartBeat:
		pld := msg.Payload.(*payload.HeartBeatPayload)
		peer.UpdateHeight(pld.Pulse.Height())
		f.peerSet.UpdateMaxClaimedHeight(pld.Pulse.Height() - 1)

	case payload.PayloadTypeProposal:
		pld := msg.Payload.(*payload.ProposalPayload)
		peer.UpdateHeight(pld.Proposal.Height())
		f.peerSet.UpdateMaxClaimedHeight(pld.Proposal.Height() - 1)

	case payload.PayloadTypeVote:
		pld := msg.Payload.(*payload.VotePayload)
		peer.UpdateHeight(pld.Vote.Height())
		f.peerSet.UpdateMaxClaimedHeight(pld.Vote.Height() - 1)

	case payload.PayloadTypeVoteSet:
		pld := msg.Payload.(*payload.VoteSetPayload)
		peer.UpdateHeight(pld.Height)
		f.peerSet.UpdateMaxClaimedHeight(pld.Height - 1)
	}

	return msg
}

func (f *Firewall) badPeer(peer *peerset.Peer) bool {
	ratio := (peer.InvalidMsg() * 100) / peer.ReceivedMsg()

	return ratio > 10
}
