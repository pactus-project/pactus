package network

import (
	"io"

	lp2pcore "github.com/libp2p/go-libp2p/core"
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
	EventTypeConnect    EventType = 1
	EventTypeDisconnect EventType = 2
	EventTypeProtocols  EventType = 3
	EventTypeGossip     EventType = 4
	EventTypeStream     EventType = 5
)

func (t EventType) String() string {
	switch t {
	case EventTypeConnect:
		return "connect"

	case EventTypeDisconnect:
		return "disconnect"

	case EventTypeProtocols:
		return "protocols"

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

// GossipMessage represents message from PubSub module.
// `From` is the ID of the peer that we received a message from.
type GossipMessage struct {
	From lp2pcore.PeerID
	Data []byte
}

func (*GossipMessage) Type() EventType {
	return EventTypeGossip
}

// StreamMessage represents message from Stream module.
// `From` is the ID of the peer that we received a message from.
type StreamMessage struct {
	From   lp2pcore.PeerID
	Reader io.ReadCloser
}

func (*StreamMessage) Type() EventType {
	return EventTypeStream
}

// ConnectEvent represents a peer connection event.
type ConnectEvent struct {
	PeerID        lp2pcore.PeerID
	RemoteAddress string
	Direction     string
}

func (*ConnectEvent) Type() EventType {
	return EventTypeConnect
}

// DisconnectEvent represents a peer disconnection event.
type DisconnectEvent struct {
	PeerID lp2pcore.PeerID
}

func (*DisconnectEvent) Type() EventType {
	return EventTypeDisconnect
}

// ProtocolsEvents represents updating protocols event.
type ProtocolsEvents struct {
	PeerID        lp2pcore.PeerID
	Protocols     []string
	SupportStream bool
}

func (*ProtocolsEvents) Type() EventType {
	return EventTypeProtocols
}

// ShouldPropagate determines whether a message should be disregarded:
// it will be neither delivered to the application nor forwarded to the network.
type ShouldPropagate func(*GossipMessage) bool

type Network interface {
	Start() error
	Stop()
	Protect(lp2pcore.PeerID, string)
	EventChannel() <-chan Event
	Broadcast([]byte, TopicID) error
	SendTo([]byte, lp2pcore.PeerID) error
	JoinGeneralTopic(shouldPropagate ShouldPropagate) error
	JoinConsensusTopic(shouldPropagate ShouldPropagate) error
	CloseConnection(pid lp2pcore.PeerID)
	SelfID() lp2pcore.PeerID
	NumConnectedPeers() int
	NumInbound() int
	NumOutbound() int
	ReachabilityStatus() string
	HostAddrs() []string
	Name() string
	Protocols() []string
}
