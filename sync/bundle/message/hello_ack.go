package message

import (
	"fmt"

	"github.com/pactus-project/pactus/network"
)

type HelloAckMessage struct {
	ResponseCode ResponseCode `cbor:"1,keyasint"`
	Reason       string       `cbor:"2,keyasint"`
	Height       uint32       `cbor:"3,keyasint"`
}

func NewHelloAckMessage(code ResponseCode, reason string, height uint32) *HelloAckMessage {
	return &HelloAckMessage{
		ResponseCode: code,
		Reason:       reason,
		Height:       height,
	}
}

func (*HelloAckMessage) BasicCheck() error {
	return nil
}

func (*HelloAckMessage) Type() Type {
	return TypeHelloAck
}

func (*HelloAckMessage) TopicID() network.TopicID {
	return network.TopicIDUnspecified
}

func (*HelloAckMessage) ShouldBroadcast() bool {
	return false
}

func (*HelloAckMessage) ConsensusHeight() uint32 {
	return 0
}

// LogString returns a concise string representation intended for use in logs.
func (m *HelloAckMessage) LogString() string {
	return fmt.Sprintf("{%s: %s %v}", m.ResponseCode, m.Reason, m.Height)
}
