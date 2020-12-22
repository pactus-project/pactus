package sync

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	peer "github.com/libp2p/go-libp2p-peer"
	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/message"
	"github.com/zarbchain/zarb-go/message/payload"
)

type mockNetworkAPI struct {
	ch chan *message.Message
}

func mockingNetworkAPI() *mockNetworkAPI {
	return &mockNetworkAPI{
		ch: make(chan *message.Message, 10),
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
	return tSelfID
}

func (mock *mockNetworkAPI) waitingForMessage(t *testing.T, msg *message.Message) {
	timeout := time.NewTimer(1 * time.Second)

	for {
		select {
		case <-timeout.C:
			assert.NoError(t, fmt.Errorf("Timeout"))
			return
		case apiMsg := <-mock.ch:
			logger.Info("comparing messages", "apiMsg", apiMsg, "msg", msg)
			b1, _ := msg.MarshalCBOR()
			b2, _ := apiMsg.MarshalCBOR()

			tSync.ParsMessage(b2, tPeerID)
			if reflect.DeepEqual(b1, b2) {
				return
			}
		}
	}
}
func (mock *mockNetworkAPI) shouldReceiveMessageWithThisType(t *testing.T, pldType payload.PayloadType) {
	timeout := time.NewTimer(1 * time.Second)

	for {
		select {
		case <-timeout.C:
			assert.NoError(t, fmt.Errorf("Timeout"))
			return
		case apiMsg := <-mock.ch:
			logger.Info("comparing messages type", "apiMsg", apiMsg)
			b, _ := apiMsg.MarshalCBOR()

			tSync.ParsMessage(b, tPeerID)
			if apiMsg.PayloadType() == pldType {
				return
			}
		}
	}
}

func (mock *mockNetworkAPI) shouldNotReceiveAnyMessageWithThisType(t *testing.T, payloadType payload.PayloadType) {
	timeout := time.NewTimer(1 * time.Second)

	for {
		select {
		case <-timeout.C:
			return
		case apiMsg := <-mock.ch:
			logger.Info("comparing messages type", "apiMsg", apiMsg)
			b, _ := apiMsg.MarshalCBOR()

			tSync.ParsMessage(b, tPeerID)
			assert.NotEqual(t, apiMsg.PayloadType(), payloadType)
		}
	}
}
