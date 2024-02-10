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
	SendError error
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

func (mock *MockNetwork) Start() error {
	return nil
}

func (mock *MockNetwork) Stop() {}

func (mock *MockNetwork) Protect(_ lp2pcore.PeerID, _ string) {}

func (mock *MockNetwork) EventChannel() <-chan Event {
	return mock.EventCh
}

func (mock *MockNetwork) JoinGeneralTopic(_ ShouldPropagate) error {
	return nil
}

func (mock *MockNetwork) JoinConsensusTopic(_ ShouldPropagate) error {
	return nil
}

func (mock *MockNetwork) SelfID() lp2ppeer.ID {
	return mock.ID
}

func (mock *MockNetwork) SendTo(data []byte, pid lp2pcore.PeerID) error {
	if mock.SendError != nil {
		return mock.SendError
	}
	mock.PublishCh <- PublishData{
		Data:   data,
		Target: &pid,
	}

	return nil
}

func (mock *MockNetwork) Broadcast(data []byte, _ TopicID) error {
	mock.PublishCh <- PublishData{
		Data:   data,
		Target: nil, // Send to all
	}

	return nil
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

func (mock *MockNetwork) ReachabilityStatus() string {
	return "Unknown"
}

func (mock *MockNetwork) HostAddrs() []string {
	return []string{"localhost"}
}

func (mock *MockNetwork) Name() string {
	return "pactus"
}

func (mock *MockNetwork) Protocols() []string {
	return []string{"gossip"}
}

func (mock *MockNetwork) NumInbound() int {
	return 0
}

func (mock *MockNetwork) NumOutbound() int {
	return 0
}
