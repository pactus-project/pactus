package network

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/sync/message"
)

type MockNetwork struct {
	BroadcastCh chan *message.Message
	id          peer.ID
	CallbackFn  CallbackFn
	OtherNet    *MockNetwork
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
func (mock *MockNetwork) PublishMessage(msg *message.Message) error {
	mock.BroadcastCh <- msg
	return nil
}
func (mock *MockNetwork) ReceivedMessage(msg *message.Message, id peer.ID) {
	d, _ := msg.Encode()
	if d != nil {
		logger.Info("Parsing the message", "msg", msg)
		mock.CallbackFn(d, id)
	}
}
