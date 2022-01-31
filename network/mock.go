package network

import (
	"bytes"

	lp2pcore "github.com/libp2p/go-libp2p-core"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/sync/message"
	"github.com/zarbchain/zarb-go/sync/message/payload"
)

var _ Network = &MockNetwork{}

type MockNetwork struct {
	BroadcastCh chan []byte
	id          peer.ID
	CallbackFn  CallbackFn
	OtherNet    *MockNetwork
	Closed      bool
}

func MockingNetwork(id peer.ID) *MockNetwork {
	return &MockNetwork{
		BroadcastCh: make(chan []byte, 1000),
		id:          id,
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
	return mock.id
}
func (mock *MockNetwork) ReceivingMessageFromOtherPeer(initiator peer.ID, pld payload.Payload) {
	msg := message.NewMessage(initiator, pld)
	d, _ := msg.Encode()
	if d != nil {
		logger.Info("Parsing the message", "msg", msg)
		mock.CallbackFn(bytes.NewReader(d), initiator)
	}
}
func (mock *MockNetwork) SendMessage(msg []byte, _ lp2pcore.PeerID) error {
	mock.BroadcastCh <- msg
	return nil
}
func (mock *MockNetwork) BroadcastMessage(msg []byte, _ TopicID) error {
	mock.BroadcastCh <- msg
	return nil
}
func (mock *MockNetwork) SendMessageToOthePeer(msg []byte) {
	logger.Debug("Sending message to other peer", "msg", msg)
	mock.OtherNet.CallbackFn(bytes.NewReader(msg), mock.id)
}
func (mock *MockNetwork) CloseConnection(pid peer.ID) {
	mock.Closed = true
}
func (mock *MockNetwork) NumConnectedPeers() int {
	return 0
}
