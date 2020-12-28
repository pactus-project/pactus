package sync

import (
	"fmt"
	"testing"
	"time"

	peer "github.com/libp2p/go-libp2p-peer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/message"
	"github.com/zarbchain/zarb-go/message/payload"
)

type mockNetworkAPI struct {
	ch       chan *message.Message
	id       peer.ID
	peerSync *Synchronizer
}

func mockingNetworkAPI(id peer.ID) *mockNetworkAPI {
	return &mockNetworkAPI{
		ch: make(chan *message.Message, 100),
		id: id,
	}
}
func (mock *mockNetworkAPI) Start() error {
	return nil
}
func (mock *mockNetworkAPI) Stop() {
}
func (mock *mockNetworkAPI) PublishMessage(msg *message.Message) error {
	mock.ch <- msg
	return nil
}
func (mock *mockNetworkAPI) SelfID() peer.ID {
	return mock.id
}

func (mock *mockNetworkAPI) shouldPublishThisMessage(t *testing.T, expectedMsg *message.Message) {
	timeout := time.NewTimer(2 * time.Second)

	for {
		select {
		case <-timeout.C:
			require.NoError(t, fmt.Errorf("Timeout"))
		case msg := <-mock.ch:
			logger.Info("shouldPublishMessageWithThisType", "id", mock.id, "msg", msg)
			b, _ := msg.MarshalCBOR()
			mock.peerSync.ParsMessage(b, mock.id)

			if msg.PayloadType() == expectedMsg.PayloadType() {
				logger.Info("Comparing two messages", "msg", msg, "expected", expectedMsg)
				assert.Equal(t, msg.SignBytes(), expectedMsg.SignBytes())
				return
			}
		}
	}
}

func (mock *mockNetworkAPI) shouldPublishMessageWithThisType(t *testing.T, payloadType payload.PayloadType) *message.Message {
	timeout := time.NewTimer(2 * time.Second)

	for {
		select {
		case <-timeout.C:
			require.NoError(t, fmt.Errorf("Timeout"))
			return nil
		case msg := <-mock.ch:
			logger.Info("shouldPublishMessageWithThisType", "id", mock.id, "msg", msg, "type", payloadType.String())
			b, _ := msg.MarshalCBOR()
			mock.peerSync.ParsMessage(b, mock.id)

			if msg.PayloadType() == payloadType {
				return msg
			}
		}
	}
}

func (mock *mockNetworkAPI) shouldNotPublishMessageWithThisType(t *testing.T, payloadType payload.PayloadType) {
	timeout := time.NewTimer(1 * time.Second)

	for {
		select {
		case <-timeout.C:
			return
		case msg := <-mock.ch:
			logger.Info("shouldNotPublishMessageWithThisType", "id", mock.id, "msg", msg, "type", payloadType.String())

			b, _ := msg.MarshalCBOR()
			mock.peerSync.ParsMessage(b, mock.id)

			assert.NotEqual(t, msg.PayloadType(), payloadType)
		}
	}
}
