package network

import (
	"bytes"
	"io"

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

func (mock *MockNetwork) SelfID() lp2ppeer.ID {
	return mock.ID
}

func (mock *MockNetwork) SendTo(data []byte, pid lp2pcore.PeerID) {
	net, exists := mock.OtherNets[pid]
	if exists {
		// direct message
		event := &StreamMessage{
			From:   mock.ID,
			Reader: io.NopCloser(bytes.NewReader(data)),
		}

		net.EventPipe.Send(event)
	}

	mock.PublishCh <- PublishData{
		Data:   data,
		Target: &pid,
	}
}

func (mock *MockNetwork) Broadcast(data []byte, _ TopicID) {
	for _, net := range mock.OtherNets {
		if net.SelfID() == mock.ID {
			continue
		}

		// Broadcast message
		event := &GossipMessage{
			From: mock.ID,
			Data: data,
		}
		net.EventPipe.Send(event)
	}

	mock.PublishCh <- PublishData{
		Data:   data,
		Target: nil, // Send to all
	}
}

func (mock *MockNetwork) CloseConnection(pid lp2ppeer.ID) {
	delete(mock.OtherNets, pid)
}

func (mock *MockNetwork) IsClosed(pid lp2ppeer.ID) bool {
	_, exists := mock.OtherNets[pid]

	return !exists
}

func (mock *MockNetwork) NumConnectedPeers() int {
	return len(mock.OtherNets)
}

func (mock *MockNetwork) AddAnotherNetwork(otherNet *MockNetwork) {
	mock.OtherNets[otherNet.SelfID()] = otherNet

	mock.EventPipe.Send(&ConnectEvent{
		PeerID:        otherNet.SelfID(),
		RemoteAddress: mock.RandMultiAddress(),
	})
	mock.EventPipe.Send(&ProtocolsEvents{
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

func (*MockNetwork) NumInbound() int {
	return 0
}

func (*MockNetwork) NumOutbound() int {
	return 0
}
