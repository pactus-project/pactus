package firewall

import (
	"encoding/hex"

	"github.com/zarbchain/zarb-go/state"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/sync/message"
	"github.com/zarbchain/zarb-go/sync/peerset"
	"github.com/zarbchain/zarb-go/util"
)

// Firewall check packets before passing them to sync module
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
	msg := new(message.Message)
	err := msg.Decode(data)

	peer.IncreaseReceivedMessage()
	peer.IncreaseReceivedBytes(len(data))
	if err != nil {
		peer.IncreaseInvalidMessage()
		logger.Debug("Error decoding message", "from", util.FingerprintPeerID(from), "data", hex.EncodeToString(data), "err", err)
		return nil
	}

	if err = msg.SanityCheck(); err != nil {
		peer.IncreaseInvalidMessage()
		logger.Debug("Peer sent us invalid msg", "peer", util.FingerprintPeerID(from), "msg", msg, "err", err)
		return nil
	}

	if f.badPeer(peer) {
		return nil
	}

	return msg
}

func (f *Firewall) badPeer(peer *peerset.Peer) bool {
	ratio := (peer.InvalidMsg() * 100) / peer.ReceivedMsg()

	return ratio > 10
}
