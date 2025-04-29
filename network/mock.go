package network

import (
	"bytes"
	"io"
	"sync"

	lp2pcore "github.com/libp2p/go-libp2p/core"
	lp2ppeer "github.com/libp2p/go-libp2p/core/peer"
	"github.com/pactus-project/pactus/util/pipeline"
	"github.com/pactus-project/pactus/util/testsuite"
)

var _ Network = &MockNetwork{}

type PublishData struct {
	Data   []byte
	Target *lp2pcore.PeerID
}

type MockNetwork struct {
	*testsuite.TestSuite

	lk        sync.RWMutex
	ID        lp2ppeer.ID
	PublishCh chan PublishData
	EventPipe pipeline.Pipeline[Event]
	OtherNets map[lp2ppeer.ID]*MockNetwork
}

func MockingNetwork(ts *testsuite.TestSuite, pid lp2ppeer.ID) *MockNetwork {
	return &MockNetwork{
		TestSuite: ts,
		PublishCh: make(chan PublishData, 100),
		EventPipe: pipeline.MockingPipeline[Event](),
		OtherNets: make(map[lp2ppeer.ID]*MockNetwork),
		ID:        pid,
	}
}

func (*MockNetwork) Start() error {
	return nil
}

func (*MockNetwork) Stop() {}

func (*MockNetwork) Protect(_ lp2pcore.PeerID, _ string) {}

func (*MockNetwork) JoinTopic(_ TopicID, _ PropagationEvaluator) error {
	return nil
}

func (m *MockNetwork) SelfID() lp2ppeer.ID {
	return m.ID
}

func (m *MockNetwork) SendTo(data []byte, pid lp2pcore.PeerID) {
	m.lk.RLock()
	defer m.lk.RUnlock()

	net, exists := m.OtherNets[pid]
	if exists {
		// direct message
		event := &StreamMessage{
			From:   m.ID,
			Reader: io.NopCloser(bytes.NewReader(data)),
		}

		net.EventPipe.Send(event)
	}

	m.PublishCh <- PublishData{
		Data:   data,
		Target: &pid,
	}
}

func (m *MockNetwork) Broadcast(data []byte, _ TopicID) {
	m.lk.RLock()
	defer m.lk.RUnlock()

	for _, net := range m.OtherNets {
		if net.SelfID() == m.ID {
			continue
		}

		// Broadcast message
		event := &GossipMessage{
			From: m.ID,
			Data: data,
		}
		net.EventPipe.Send(event)
	}

	m.PublishCh <- PublishData{
		Data:   data,
		Target: nil, // Send to all
	}
}

func (m *MockNetwork) CloseConnection(pid lp2ppeer.ID) {
	m.lk.RLock()
	defer m.lk.RUnlock()

	delete(m.OtherNets, pid)
}

func (m *MockNetwork) IsClosed(pid lp2ppeer.ID) bool {
	m.lk.RLock()
	defer m.lk.RUnlock()

	_, exists := m.OtherNets[pid]

	return !exists
}

func (m *MockNetwork) NumConnectedPeers() int {
	return len(m.OtherNets)
}

func (m *MockNetwork) AddAnotherNetwork(otherNet *MockNetwork) {
	m.lk.Lock()
	defer m.lk.Unlock()

	m.OtherNets[otherNet.SelfID()] = otherNet

	m.EventPipe.Send(&ConnectEvent{
		PeerID:        otherNet.SelfID(),
		RemoteAddress: m.RandMultiAddress(),
	})
	m.EventPipe.Send(&ProtocolsEvents{
		PeerID: otherNet.SelfID(),
	})
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

func (m *MockNetwork) NumInbound() int {
	return len(m.OtherNets)
}

func (m *MockNetwork) NumOutbound() int {
	return len(m.OtherNets)
}
