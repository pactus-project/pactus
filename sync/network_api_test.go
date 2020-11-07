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

type mockNetworkApi struct {
	ch chan *message.Message
}

func mockingNetworkApi() *mockNetworkApi {
	return &mockNetworkApi{
		ch: make(chan *message.Message, 10),
	}
}
func (mock *mockNetworkApi) Start() error {
	return nil
}
func (mock *mockNetworkApi) Stop() error {
	return nil
}
func (mock *mockNetworkApi) PublishMessage(msg *message.Message) error {
	mock.ch <- msg
	return nil
}

func (mock *mockNetworkApi) waitingForMessage(t *testing.T, msg *message.Message) {
	timeout := time.NewTimer(1 * time.Second)

	for {
		select {
		case <-timeout.C:
			assert.NoError(t, fmt.Errorf("Timeout"))
			return
		case apiMsg := <-mock.ch:
			logger.Info("comparing messages", "apiMsg", apiMsg, "msg", msg)
			if reflect.DeepEqual(msg, apiMsg) {
				return
			}
		}
	}

}
