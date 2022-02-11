package network

import (
	"bytes"

	lp2pcore "github.com/libp2p/go-libp2p-core"
	"github.com/libp2p/go-libp2p-core/peer"
)

var _ Network = &MockNetwork{}

type BroadcastData struct {
	Data   []byte
	Target *lp2pcore.PeerID
}

type MockNetwork struct {
	BroadcastCh chan BroadcastData
	CallbackFn  CallbackFn
	ID          peer.ID
	OtherNet    []*MockNetwork
}

func MockingNetwork(id peer.ID) *MockNetwork {
	return &MockNetwork{
		BroadcastCh: make(chan BroadcastData, 100),
		OtherNet:    make([]*MockNetwork, 0),
		ID:          id,
	}
}
func (mock *MockNetwork) Start() error {
	return nil
}
func (mock *MockNetwork) SetCallback(callbackFn CallbackFn) {
	mock.CallbackFn = callbackFn
}
func (mock *MockNetwork) Stop() {
}
func (mock *MockNetwork) JoinGeneralTopic() error {
	return nil
}
func (mock *MockNetwork) JoinConsensusTopic() error {
	return nil
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
	for _, net := range mock.OtherNet {
		if target == nil || net.ID == *target {
			net.CallbackFn(bytes.NewReader(data), mock.ID, mock.ID)
		}
	}
}
func (mock *MockNetwork) CloseConnection(pid peer.ID) {
	for i, net := range mock.OtherNet {
		if net.ID == pid {
			mock.OtherNet = append(mock.OtherNet[:i], mock.OtherNet[i+1:]...)
		}
	}
}
func (mock *MockNetwork) IsClosed(pid peer.ID) bool {
	for _, net := range mock.OtherNet {
		if net.ID == pid {
			return false
		}
	}
	return true
}
func (mock *MockNetwork) NumConnectedPeers() int {
	return len(mock.OtherNet)
}
func (mock *MockNetwork) AddAnotherNetwork(net *MockNetwork) {
	mock.OtherNet = append(mock.OtherNet, net)
}
