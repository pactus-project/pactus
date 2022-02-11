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

/// CallbackFn is a call back function to receive data from the network.
/// `from` is the ID of the peer that initiate the message.
type CallbackFn func(reader io.Reader, from peer.ID)

type Network interface {
	Start() error
	Stop()
	SetCallback(CallbackFn)
	Broadcast([]byte, TopicID) error
	SendTo([]byte, lp2pcore.PeerID) error
	JoinGeneralTopic() error
	JoinConsensusTopic() error
	CloseConnection(pid peer.ID)
	SelfID() peer.ID
	NumConnectedPeers() int
}
