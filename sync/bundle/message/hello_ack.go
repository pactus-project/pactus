package message

import (
	"fmt"
)

type HelloAckMessage struct {
	ResponseCode ResponseCode `cbor:"1,keyasint"`
	Reason       string       `cbor:"2,keyasint"`
}

func NewHelloAckMessage(code ResponseCode, reason string) *HelloAckMessage {
	return &HelloAckMessage{
		ResponseCode: code,
		Reason:       reason,
	}
}

func (m *HelloAckMessage) BasicCheck() error {
	return nil
}

func (m *HelloAckMessage) Type() Type {
	return TypeHelloAck
}

func (m *HelloAckMessage) String() string {
	return fmt.Sprintf("{%s: %s}", m.ResponseCode, m.Reason)
}
