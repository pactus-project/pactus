package network

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/sync/message"
	"github.com/zarbchain/zarb-go/sync/message/payload"
)

type MockNetwork struct {
	BroadcastCh chan *message.Message
	id          peer.ID
	CallbackFn  CallbackFn
	OtherNet    *MockNetwork
	Closed      bool
}

func MockingNetwork(id peer.ID) *MockNetwork {
	return &MockNetwork{
		BroadcastCh: make(chan *message.Message, 1000),
		id:          id,
	}
}
func (mock *MockNetwork) Start() error {
	return nil
}
func (mock *MockNetwork) Stop() {
}
func (mock *MockNetwork) JoinTopics(callbackFn CallbackFn) error {
	mock.CallbackFn = callbackFn
	return nil
}
func (mock *MockNetwork) JoinDownloadTopic() error {
	return nil
}
func (mock *MockNetwork) LeaveDownloadTopic() {}
func (mock *MockNetwork) SelfID() peer.ID {
	return mock.id
}
func (mock *MockNetwork) ReceivingMessageFromOtherPeer(id peer.ID, pld payload.Payload) {
	msg := message.NewMessage(id, pld)
	d, _ := msg.Encode()
	if d != nil {
		logger.Info("Parsing the message", "msg", msg)
		mock.CallbackFn(d, id)
	}
}
func (mock *MockNetwork) PublishMessage(msg *message.Message) error {
	mock.BroadcastCh <- msg
	return nil
}
func (mock *MockNetwork) SendMessageToOthePeer(msg *message.Message) {
	d, _ := msg.Encode()
	if d != nil {
		logger.Debug("Sending message to other peer", "msg", msg)
		mock.OtherNet.CallbackFn(d, mock.id)
	}
}
func (mock *MockNetwork) CloseConnection(pid peer.ID) {
	mock.Closed = true
}
