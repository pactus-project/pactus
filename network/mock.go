package network

import (
	"bytes"
	"io/ioutil"

	lp2pcore "github.com/libp2p/go-libp2p/core"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/util"
)

var _ Network = &MockNetwork{}

type BroadcastData struct {
	Data   []byte
	Target *lp2pcore.PeerID
}

type MockNetwork struct {
	BroadcastCh chan BroadcastData
	EventCh     chan Event
	ID          peer.ID
	OtherNets   []*MockNetwork
}

func MockingNetwork(id peer.ID) *MockNetwork {
	return &MockNetwork{
		BroadcastCh: make(chan BroadcastData, 100),
		EventCh:     make(chan Event, 100),
		OtherNets:   make([]*MockNetwork, 0),
		ID:          id,
	}
}
func (mock *MockNetwork) Start() error {
	return nil
}
func (mock *MockNetwork) Stop() {
}

func (mock *MockNetwork) EventChannel() <-chan Event {
	return mock.EventCh
}
func (mock *MockNetwork) JoinGeneralTopic() error {
	return nil
}
func (mock *MockNetwork) JoinConsensusTopic() error {
	return nil
}
func (mock *MockNetwork) NumOfPeersInGeneralTopic() int {
	return len(mock.OtherNets)
}
func (mock *MockNetwork) NumOfPeersInConsensusTopic() int {
	return len(mock.OtherNets)
}
func (mock *MockNetwork) SelfID() peer.ID {
	return mock.ID
}
func (mock *MockNetwork) SendTo(data []byte, pid lp2pcore.PeerID) error {
	mock.BroadcastCh <- BroadcastData{
		Data:   data,
		Target: &pid,
	}
	return nil
}
func (mock *MockNetwork) Broadcast(data []byte, tid TopicID) error {
	mock.BroadcastCh <- BroadcastData{
		Data:   data,
		Target: nil, // Send to all
	}
	return nil
}

func (mock *MockNetwork) SendToOthers(data []byte, target *peer.ID) {
	for _, net := range mock.OtherNets {
		if target == nil {
			// Broadcast message
			event := &GossipMessage{
				Source: mock.ID,
				From:   mock.ID,
				Data:   data,
			}
			net.EventCh <- event
		} else if net.ID == *target {
			// direct message
			event := &StreamMessage{
				Source: mock.ID,
				Reader: ioutil.NopCloser(bytes.NewReader(data)),
			}
			net.EventCh <- event
		}
	}
}
func (mock *MockNetwork) CloseConnection(pid peer.ID) {
	for i, net := range mock.OtherNets {
		if net.ID == pid {
			mock.OtherNets = append(mock.OtherNets[:i], mock.OtherNets[i+1:]...)
		}
	}
}
func (mock *MockNetwork) IsClosed(pid peer.ID) bool {
	for _, net := range mock.OtherNets {
		if net.ID == pid {
			return false
		}
	}
	return true
}
func (mock *MockNetwork) NumConnectedPeers() int {
	return len(mock.OtherNets)
}
func (mock *MockNetwork) AddAnotherNetwork(net *MockNetwork) {
	mock.OtherNets = append(mock.OtherNets, net)
}

// TestRandomPeerID returns a random peer ID.
func TestRandomPeerID() peer.ID {
	s := util.Uint64ToSlice(util.RandUint64(0))
	id := [34]byte{0x12, 32}
	copy(id[2:], s[:])
	return peer.ID(id[:])
}
