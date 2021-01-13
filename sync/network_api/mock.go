package network_api

import (
	"fmt"
	"testing"
	"time"

	peer "github.com/libp2p/go-libp2p-peer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/sync/firewall"
	"github.com/zarbchain/zarb-go/sync/message"
	"github.com/zarbchain/zarb-go/sync/message/payload"
)

type MockNetworkAPI struct {
	ch       chan *message.Message
	id       peer.ID
	Firewall *firewall.Firewall
	ParsFn   ParsMessageFn
	OtherAPI *MockNetworkAPI
}

func MockingNetworkAPI(id peer.ID) *MockNetworkAPI {
	return &MockNetworkAPI{
		ch: make(chan *message.Message, 100),
		id: id,
	}
}
func (mock *MockNetworkAPI) Start() error {
	return nil
}
func (mock *MockNetworkAPI) Stop() {
}
func (mock *MockNetworkAPI) JoinDownloadTopic() error {
	return nil
}
func (mock *MockNetworkAPI) LeaveDownloadTopic() {}
func (mock *MockNetworkAPI) PublishMessage(msg *message.Message) error {
	mock.ch <- msg
	return nil
}
func (mock *MockNetworkAPI) SelfID() peer.ID {
	return mock.id
}
func (mock *MockNetworkAPI) CheckAndParsMessage(data []byte, id peer.ID) bool {
	msg := mock.Firewall.ParsMessage(data, id)
	if msg != nil {
		mock.ParsFn(msg, mock.id)
		return true
	}
	return false
}

func (mock *MockNetworkAPI) sendMessageToOtherPeer(m *message.Message) {
	data, _ := m.Encode()
	mock.OtherAPI.CheckAndParsMessage(data, mock.id)
}

func (mock *MockNetworkAPI) ShouldPublishThisMessage(t *testing.T, expectedMsg *message.Message) {
	timeout := time.NewTimer(2 * time.Second)

	for {
		select {
		case <-timeout.C:
			require.NoError(t, fmt.Errorf("Timeout"))
		case msg := <-mock.ch:
			logger.Info("shouldPublishMessageWithThisType", "id", mock.id, "msg", msg)
			mock.sendMessageToOtherPeer(msg)

			if msg.PayloadType() == expectedMsg.PayloadType() {
				logger.Info("Comparing two messages", "msg", msg, "expected", expectedMsg)
				bs1, _ := msg.Encode()
				bs2, _ := expectedMsg.Encode()
				assert.Equal(t, bs1, bs2)
				return
			}
		}
	}
}

func (mock *MockNetworkAPI) ShouldPublishMessageWithThisType(t *testing.T, payloadType payload.PayloadType) *message.Message {
	timeout := time.NewTimer(2 * time.Second)

	for {
		select {
		case <-timeout.C:
			require.NoError(t, fmt.Errorf("Timeout"))
			return nil
		case msg := <-mock.ch:
			logger.Info("shouldPublishMessageWithThisType", "id", mock.id, "msg", msg, "type", payloadType.String())
			mock.sendMessageToOtherPeer(msg)

			if msg.PayloadType() == payloadType {
				return msg
			}
		}
	}
}

func (mock *MockNetworkAPI) ShouldNotPublishMessageWithThisType(t *testing.T, payloadType payload.PayloadType) {
	timeout := time.NewTimer(1 * time.Second)

	for {
		select {
		case <-timeout.C:
			return
		case msg := <-mock.ch:
			logger.Info("shouldNotPublishMessageWithThisType", "id", mock.id, "msg", msg, "type", payloadType.String())
			mock.sendMessageToOtherPeer(msg)

			assert.NotEqual(t, msg.PayloadType(), payloadType)
		}
	}
}
