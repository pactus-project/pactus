package network

import (
	"io"

	lp2pcore "github.com/libp2p/go-libp2p-core"
)

type TopicID int

const (
	TopicIDGeneral   TopicID = 1
	TopicIDConsensus TopicID = 2
)

func (t TopicID) String() string {
	switch t {
	case TopicIDGeneral:
		return "general"
	case TopicIDConsensus:
		return "consensus"
	}
	return "invalid"
}

type EventType int

const (
	EventTypeGossip EventType = 1
	EventTypeStream EventType = 2
)

func (t EventType) String() string {
	switch t {
	case EventTypeGossip:
		return "gossip-msg"
	case EventTypeStream:
		return "stream-msg"
	}
	return "invalid"
}

type Event interface {
	Type() EventType
}

/// `GossipMessage` represents message from PubSub module.
/// `source` is the ID of the peer that initiate the message and
/// `from` is the ID of the peer that we received a message from.
/// They are not necessarily the same, especially in a decentralized network.
type GossipMessage struct {
	Source lp2pcore.PeerID
	From   lp2pcore.PeerID
	Data   []byte
}

func (*GossipMessage) Type() EventType {
	return EventTypeGossip
}

/// `GossipMessage` represents message from stream module.
type StreamMessage struct {
	Source lp2pcore.PeerID
	Reader io.ReadCloser
}

func (*StreamMessage) Type() EventType {
	return EventTypeStream
}

type Network interface {
	Start() error
	Stop()
	EventChannel() <-chan Event
	Broadcast([]byte, TopicID) error
	SendTo([]byte, lp2pcore.PeerID) error
	JoinGeneralTopic() error
	JoinConsensusTopic() error
	CloseConnection(pid lp2pcore.PeerID)
	SelfID() lp2pcore.PeerID
	NumConnectedPeers() int
}
