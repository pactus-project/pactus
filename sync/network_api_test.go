package sync

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/message"
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
func (mock *mockNetworkAPI) Stop() error {
	return nil
}
func (mock *mockNetworkAPI) PublishMessage(msg *message.Message) error {
	mock.ch <- msg
	return nil
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

			if reflect.DeepEqual(b1, b2) {
				return
			}
		}
	}

}
