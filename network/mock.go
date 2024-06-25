package network

import (
	"bytes"
	"io"

	lp2pcore "github.com/libp2p/go-libp2p/core"
	lp2ppeer "github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/util/testsuite"
)

var _ Network = &MockNetwork{}

type PublishData struct {
	Data   []byte
	Target *lp2pcore.PeerID
}

type MockNetwork struct {
	*testsuite.TestSuite

	PublishCh chan PublishData
	EventCh   chan Event
	ID        lp2ppeer.ID
	OtherNets []*MockNetwork
}

func MockingNetwork(ts *testsuite.TestSuite, id lp2ppeer.ID) *MockNetwork {
	return &MockNetwork{
		TestSuite: ts,
		PublishCh: make(chan PublishData, 100),
		EventCh:   make(chan Event, 100),
		OtherNets: make([]*MockNetwork, 0),
		ID:        id,
	}
}

func (*MockNetwork) Start() error {
	return nil
}

func (*MockNetwork) Stop() {}

func (*MockNetwork) Protect(_ lp2pcore.PeerID, _ string) {}

func (mock *MockNetwork) EventChannel() <-chan Event {
	return mock.EventCh
}

func (*MockNetwork) JoinTopic(_ TopicID, _ ShouldPropagate) error {
	return nil
}

func (mock *MockNetwork) SelfID() lp2ppeer.ID {
	return mock.ID
}

func (mock *MockNetwork) SendTo(data []byte, pid lp2pcore.PeerID) {
	mock.PublishCh <- PublishData{
		Data:   data,
		Target: &pid,
	}
}

func (mock *MockNetwork) Broadcast(data []byte, _ TopicID) {
	mock.PublishCh <- PublishData{
		Data:   data,
		Target: nil, // Send to all
	}
}

func (mock *MockNetwork) SendToOthers(data []byte, target *lp2ppeer.ID) {
	for _, net := range mock.OtherNets {
		if target == nil {
			// Broadcast message
			event := &GossipMessage{
				From: mock.ID,
				Data: data,
			}
			net.EventCh <- event
		} else if net.ID == *target {
			// direct message
			event := &StreamMessage{
				From:   mock.ID,
				Reader: io.NopCloser(bytes.NewReader(data)),
			}
			net.EventCh <- event
		}
	}
}

func (mock *MockNetwork) CloseConnection(pid lp2ppeer.ID) {
	for i, net := range mock.OtherNets {
		if net.ID == pid {
			mock.OtherNets = append(mock.OtherNets[:i], mock.OtherNets[i+1:]...)
		}
	}
}

func (mock *MockNetwork) IsClosed(pid lp2ppeer.ID) bool {
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

func (*MockNetwork) ReachabilityStatus() string {
	return "Unknown"
}

func (*MockNetwork) HostAddrs() []string {
	return []string{"localhost"}
}

func (*MockNetwork) Name() string {
	return "pactus"
}

func (*MockNetwork) Protocols() []string {
	return []string{"gossip"}
}

func (*MockNetwork) NumInbound() int {
	return 0
}

func (*MockNetwork) NumOutbound() int {
	return 0
}
