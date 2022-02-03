package network

import (
	"io"

	lp2pcore "github.com/libp2p/go-libp2p-core"
	"github.com/libp2p/go-libp2p-core/peer"
)

type TopicID int

const (
	GeneralTopic   TopicID = 1
	ConsensusTopic TopicID = 2
)

func (ti TopicID) String() string {
	switch ti {
	case GeneralTopic:
		return "General"
	case ConsensusTopic:
		return "consensus"
	}
	return "invalid"
}

type CallbackFn func(io.Reader, peer.ID)

type Network interface {
	Start() error
	Stop()
	SetCallback(CallbackFn)
	BroadcastMessage([]byte, TopicID) error
	SendMessage([]byte, lp2pcore.PeerID) error
	JoinGeneralTopic() error
	JoinConsensusTopic() error
	CloseConnection(pid peer.ID)
	SelfID() peer.ID
	NumConnectedPeers() int
}
