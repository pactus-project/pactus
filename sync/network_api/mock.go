package network_api

import (
	"fmt"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/sync/firewall"
	"github.com/zarbchain/zarb-go/sync/message"
	"github.com/zarbchain/zarb-go/sync/message/payload"
)

type MockNetworkAPI struct {
	PublishCh chan *message.Message
	id        peer.ID
	Firewall  *firewall.Firewall
	ParsFn    ParsMessageFn
	OtherAPI  *MockNetworkAPI
}

func MockingNetworkAPI(id peer.ID) *MockNetworkAPI {
	return &MockNetworkAPI{
		PublishCh: make(chan *message.Message, 1000),
		id:        id,
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
func (mock *MockNetworkAPI) SelfID() peer.ID {
	return mock.id
}
func (mock *MockNetworkAPI) PublishMessage(msg *message.Message) error {
	mock.PublishCh <- msg
	return nil
}
func (mock *MockNetworkAPI) CheckAndParsMessage(msg *message.Message, id peer.ID) bool {
	d, _ := msg.Encode()
	msg2 := mock.Firewall.ParsMessage(d, id)
	if msg2 != nil {
		logger.Info("Parsing the message", "msg", msg)
		mock.ParsFn(msg2, mock.id)
		return true
	}
	return false
}

func (mock *MockNetworkAPI) sendMessageToOtherPeer(m *message.Message) {
	mock.OtherAPI.CheckAndParsMessage(m, mock.id)
}

func (mock *MockNetworkAPI) ShouldPublishMessageWithThisType(t *testing.T, payloadType payload.PayloadType) *message.Message {
	timeout := time.NewTimer(2 * time.Second)

	for {
		select {
		case <-timeout.C:
			require.NoError(t, fmt.Errorf("ShouldPublishMessageWithThisType: Timeout"))
			return nil
		case msg := <-mock.PublishCh:
			logger.Info("shouldPublishMessageWithThisType", "id", mock.id, "msg", msg, "type", payloadType.String())
			mock.sendMessageToOtherPeer(msg)
			logger.Info("Nessage sent to other peer", "msg", msg)

			if msg.PayloadType() == payloadType {
				return msg
			}
		}
	}
}

func (mock *MockNetworkAPI) ShouldNotPublishMessageWithThisType(t *testing.T, payloadType payload.PayloadType) {
	timeout := time.NewTimer(300 * time.Millisecond)

	for {
		select {
		case <-timeout.C:
			return
		case msg := <-mock.PublishCh:
			logger.Info("shouldNotPublishMessageWithThisType", "id", mock.id, "msg", msg, "type", payloadType.String())
			mock.sendMessageToOtherPeer(msg)
			logger.Info("Nessage sent to other peer", "msg", msg)

			assert.NotEqual(t, msg.PayloadType(), payloadType)
		}
	}
}
